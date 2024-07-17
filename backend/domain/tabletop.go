package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/gofiber/fiber/v2/log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// tabletop holds a map of save files with a mutex for concurrent access
type Tabletop struct {
	saves  map[string]*SaveFile
	rwLock sync.RWMutex
}

// SaveFile holds the byte data of a save file
type SaveFile struct {
	saveData []byte
	stale    bool
}

func NewTabletop() *Tabletop {
	return &Tabletop{saves: make(map[string]*SaveFile)}
}

func NewSavefile(fullPath string, data []byte) (*SaveFile, error) {
	err := os.WriteFile(fullPath, data, 0666)
	if err != nil {
		return nil, err
	}
	return &SaveFile{saveData: data}, err
}

// SetDirectory scans a directory and its subdirectories for JSON files,
// returning a map where the key is the full file path and the value is the SaveFile struct.
// TODO: do not remove previous SaveFiles unless we are force-reloading
func (tt *Tabletop) SetDirectory(path string) error {
	saveFiles := make(map[string]*SaveFile)
	visited := make(map[string]bool)
	// Queue for directories to process
	queue := []string{path}
	for len(queue) > 0 {
		currentDir := queue[0]
		queue = queue[1:]

		if entries, err := os.ReadDir(currentDir); err != nil {
			fmt.Printf("Error reading directory %s: %v\n", currentDir, err)
			continue
		} else {
			queue = append(queue, parseEntries(entries, currentDir, visited, saveFiles)...)
		}
	}
	tt.saves = saveFiles
	return nil
}

// GetFile checks for the file in the cache and reads it from disk if not found
func (tt *Tabletop) GetFile(path string) ([]byte, error) {
	if data := tt.trySaveFromCache(path); data != nil {
		return data, nil
	}
	return tt.getSaveFromFS(path)
}

func (tt *Tabletop) trySaveFromCache(path string) []byte {
	tt.rwLock.RLock()
	cached, exists := tt.saves[path]
	tt.rwLock.RUnlock()

	if exists {
		return cached.saveData
	}
	return nil
}

func (tt *Tabletop) getSaveFromFS(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tt.rwLock.Lock()
	tt.saves[path] = &SaveFile{saveData: data}
	tt.rwLock.Unlock()

	return data, nil
}

func (tt *Tabletop) PatchSaveFile(path string, original []byte, patches []byte) ([]byte, error) {
	patch, err := jsonpatch.DecodePatch(patches)
	if err != nil {
		return nil, err
	}
	modified, err := patch.Apply(original)
	if err != nil {
		return nil, err
	}
	err = tt.BackupSaveFile(path, modified)
	if err != nil {
		return nil, err
	}
	return modified, nil
}

// AddCardsToDeck adds a card to the deck of a save file. The patch is incomplete as it does not include
// card id, deck id, as well as the card's position in the deck.
// therefore the operation is not idempotent at this level as the same card can be added multiple times
func (tt *Tabletop) AddCardsToDeck(newCard PartialCard, savePath string, deckPath string) (error, []byte) {
	// obtain savefile
	originalBytes, err := tt.GetFile(savePath)
	if err != nil {
		return err, nil
	}
	// unmarshal with loose typing in order to be able to navigate to the deck
	var originalFile interface{}
	err = json.Unmarshal(originalBytes, &originalFile)
	if err != nil {
		return err, nil
	}
	// navigate to the deck and recast it to a Deck object via re-marshalling
	// we will use marshalled deckObject bytes later to obtain the patch
	tmpTypelessDeck, err := navigateToPath(originalFile, deckPath)
	if err != nil {
		return err, nil
	}
	var deckObject Deck
	err, deckMarshalled := remarshal(tmpTypelessDeck, &deckObject)
	if err != nil {
		return err, nil
	}
	err = deckObject.AppendNewCard(newCard.LuaScript, newCard.LuaScriptState)
	if err != nil {
		return err, nil
	}

	updatedDeck, err := json.Marshal(deckObject)
	if err != nil {
		return err, nil
	}
	// create a patch for the deck
	patch, err := createDeckJsonPatch(deckMarshalled, updatedDeck, deckPath, "replace")
	if err != nil {
		return err, nil
	}
	// implement changes and return the updated file
	updated, err := tt.PatchSaveFile(savePath, originalBytes, patch)
	if err != nil {
		return err, nil
	}
	return nil, updated
}

func getDeckFromSavefile(savePath string, deckPath string, deckObj *Deck) ([]byte, error) {
	// obtain savefile
	originalBytes, err := os.ReadFile(savePath)
	if err != nil {
		return nil, err
	}
	// unmarshal with loose typing in order to be able to navigate to the deck
	var originalFile interface{}
	err = json.Unmarshal(originalBytes, &originalFile)
	if err != nil {
		return nil, err
	}
	// navigate to the deck and recast it to a Deck object via re-marshalling
	// we will use marshalled deckObject bytes later to obtain the patch
	tmpTypelessDeck, err := navigateToPath(originalFile, deckPath)
	if err != nil {
		return nil, err
	}
	err, unmarshalledDeck := remarshal(tmpTypelessDeck, &deckObj)
	if err != nil {
		return nil, err
	}
	return unmarshalledDeck, nil
}

func (tt *Tabletop) RemoveCardFromDeck(savePath string, cardPath string) (error, []byte) {
	var deckPath = cardPath[:strings.LastIndex(cardPath, "/ContainedObjects/")]

	// obtain savefile
	originalBytes, err := tt.GetFile(savePath)
	if err != nil {
		return err, nil
	}
	cardIdx, err := strconv.Atoi(cardPath[strings.LastIndex(cardPath, "/")+1:])
	if err != nil {
		return err, nil
	}
	var originalFile interface{}
	err = json.Unmarshal(originalBytes, &originalFile)
	if err != nil {
		return err, nil
	}
	tmpTypelessDeck, err := navigateToPath(originalFile, deckPath)
	if err != nil {
		return err, nil
	}
	var deckObject Deck
	err, deckMarshalled := remarshal(tmpTypelessDeck, &deckObject)
	if err != nil {
		return err, nil
	}
	// find CardID
	var cardID = deckObject.ContainedObjects[cardIdx].CardID

	// remove deckID of th ecard
	var idxToRemove = -1
	for idx, value := range deckObject.DeckIDs {
		if value == cardID {
			idxToRemove = idx
			break
		}
	}
	if idxToRemove == -1 {
		return errors.New("Card not found in deck"), nil
	} else {
		deckObject.DeckIDs = append(deckObject.DeckIDs[:idxToRemove], deckObject.DeckIDs[idxToRemove+1:]...)
	}
	deckObject.ContainedObjects = append(deckObject.ContainedObjects[:cardIdx], deckObject.ContainedObjects[cardIdx+1:]...)
	updatedDeck, err := json.Marshal(deckObject)
	if err != nil {
		return err, nil
	}
	patch, err := createDeckJsonPatch(deckMarshalled, updatedDeck, deckPath, "replace")
	if err != nil {
		return err, nil
	}
	// implement changes and return the updated file
	updated, err := tt.PatchSaveFile(savePath, originalBytes, patch)
	if err != nil {
		return err, nil
	}
	return nil, updated
}

func createDeckJsonPatch(originalDeck, updatedDeck []byte, deckPath string, opname string) ([]byte, error) {
	// create a "relative" patch (relative to the deck, need to adjust the path parameter of each patch operation)

	patch, err := jsonpatch.CreateMergePatch(originalDeck, updatedDeck)
	if err != nil {
		return nil, err
	}
	log.Info(string(patch))
	patchObj, err := UnmarshalJsonPatch(opname, patch)
	if err != nil {
		return nil, err
	}
	// use absolute paths of the deck by prepending the deck's location to each patch operation
	for idx := range patchObj {
		patchObj[idx].Path = fmt.Sprintf("%s/%s", deckPath, patchObj[idx].Path)
		log.Info(patchObj[idx].Path)
	}
	patch, err = json.Marshal(patchObj)
	if err != nil {
		return nil, err
	}
	return patch, nil
}

// FIXME: this function does work, or the largeestBakNum does not work as advertised
func (tt *Tabletop) BackupSaveFile(path string, data []byte) error {
	save := tt.saves[path]
	if save == nil {
		return errors.New(path + " could not be found")
	}
	// backup previous file
	dir, name := filepath.Split(path)
	baknum := getLargestBakNum(dir)
	err := os.WriteFile(filepath.Join(dir, fmt.Sprint(name, ".", baknum, ".bak")), save.saveData, 0666)
	if err != nil {
		return err
	}
	// write new file
	err = os.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	// update cache
	save.saveData = data
	save.stale = true
	return nil
}

func getLargestBakNum(dir string) int {
	largestBak := 0
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		segments := strings.Split(entry.Name(), ".")
		slen := len(segments)
		if slen > 4 && slen-2 > 0 {
			bakNumStr := segments[slen-2]
			num, err := strconv.ParseInt(bakNumStr, 10, 8)
			if err != nil {
				continue
			}
			largestBak = int(num)
		}
	}
	return largestBak
}
