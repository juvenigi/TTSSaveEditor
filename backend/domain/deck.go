package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/gofiber/fiber/v2/log"
	"strconv"
	"strings"
)

type Transform struct {
	PosX   float64 `json:"posX"`
	PosY   float64 `json:"posY"`
	PosZ   float64 `json:"posZ"`
	RotX   float64 `json:"rotX"`
	RotY   float64 `json:"rotY"`
	RotZ   float64 `json:"rotZ"`
	ScaleX float64 `json:"scaleX"`
	ScaleY float64 `json:"scaleY"`
	ScaleZ float64 `json:"scaleZ"`
}

type AltLookAngle struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type ColorDiffuse struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
}

type CustomDeck struct {
	FaceURL      string `json:"FaceURL"`
	BackURL      string `json:"BackURL"`
	NumWidth     int    `json:"NumWidth"`
	NumHeight    int    `json:"NumHeight"`
	BackIsHidden bool   `json:"BackIsHidden"`
	UniqueBack   bool   `json:"UniqueBack"`
	Type         int    `json:"Type"`
}

type Card struct {
	GUID                 string                `json:"GUID"`
	Name                 string                `json:"Name"`
	Transform            Transform             `json:"Transform"`
	Nickname             string                `json:"Nickname"`
	Description          string                `json:"Description"`
	GMNotes              string                `json:"GMNotes"`
	AltLookAngle         AltLookAngle          `json:"AltLookAngle"`
	ColorDiffuse         ColorDiffuse          `json:"ColorDiffuse"`
	LayoutGroupSortIndex int                   `json:"LayoutGroupSortIndex"`
	Value                int                   `json:"Value"`
	Locked               bool                  `json:"Locked"`
	Grid                 bool                  `json:"Grid"`
	Snap                 bool                  `json:"Snap"`
	IgnoreFoW            bool                  `json:"IgnoreFoW"`
	MeasureMovement      bool                  `json:"MeasureMovement"`
	DragSelectable       bool                  `json:"DragSelectable"`
	Autoraise            bool                  `json:"Autoraise"`
	Sticky               bool                  `json:"Sticky"`
	Tooltip              bool                  `json:"Tooltip"`
	GridProjection       bool                  `json:"GridProjection"`
	HideWhenFaceDown     bool                  `json:"HideWhenFaceDown"`
	Hands                bool                  `json:"Hands"`
	CardID               int                   `json:"CardID"`
	SidewaysCard         bool                  `json:"SidewaysCard"`
	CustomDeck           map[string]CustomDeck `json:"CustomDeck"`
	LuaScript            string                `json:"LuaScript"`
	LuaScriptState       string                `json:"LuaScriptState"`
	XmlUI                string                `json:"XmlUI"`
}

type Deck struct {
	GUID                 string       `json:"GUID"`
	Name                 string       `json:"Name"`
	Transform            Transform    `json:"Transform"`
	Nickname             string       `json:"Nickname"`
	Description          string       `json:"Description"`
	GMNotes              string       `json:"GMNotes"`
	AltLookAngle         AltLookAngle `json:"AltLookAngle"`
	ColorDiffuse         ColorDiffuse `json:"ColorDiffuse"`
	LayoutGroupSortIndex int          `json:"LayoutGroupSortIndex"`
	Value                int          `json:"Value"`
	Locked               bool         `json:"Locked"`
	Grid                 bool         `json:"Grid"`
	Snap                 bool         `json:"Snap"`
	IgnoreFoW            bool         `json:"IgnoreFoW"`
	MeasureMovement      bool         `json:"MeasureMovement"`
	DragSelectable       bool         `json:"DragSelectable"`
	Autoraise            bool         `json:"Autoraise"`
	Sticky               bool         `json:"Sticky"`
	Tooltip              bool         `json:"Tooltip"`
	GridProjection       bool         `json:"GridProjection"`
	HideWhenFaceDown     bool         `json:"HideWhenFaceDown"`
	Hands                bool         `json:"Hands"`
	SidewaysCard         bool         `json:"SidewaysCard"`
	// A bit confusingly named, since it actually lists CardIDs of the cards in the deck, but supposedly helps with
	// integrity-checking CustomDeck (in case the deck consists of multiple Deck sprites)
	DeckIDs []int `json:"DeckIDs"`
	// The keys are the first 3 digits of a card's CardID representing the DeckID
	CustomDeck       map[string]CustomDeck `json:"CustomDeck"`
	LuaScript        string                `json:"LuaScript"`
	LuaScriptState   string                `json:"LuaScriptState"`
	XmlUI            string                `json:"XmlUI"`
	ContainedObjects []Card                `json:"ContainedObjects"` // Assuming the objects are of type Card
}

func (deckObject *Deck) AppendNewCard(luaScript string, luaScriptState string) error {
	var err error
	var maxCardID = len(deckObject.DeckIDs)

	// obtain DeckID prefix
	var deckKey string
	for key := range deckObject.CustomDeck {
		deckKey = key
		break
	}
	// effectively clone the first newCard in the deck to obtain the reference newCard
	var referenceCard Card
	err, _ = remarshal(deckObject.ContainedObjects[0], &referenceCard)
	if err != nil {
		return err
	}
	// now update GUID, CardID, and DeckID of the newCard
	referenceCard.GUID = "000000" //TODO
	cardIdNum, _ := strconv.Atoi(fmt.Sprintf("%s%02d", deckKey, maxCardID))
	referenceCard.CardID = cardIdNum
	// the values we care about
	referenceCard.LuaScript = luaScript
	referenceCard.LuaScriptState = luaScriptState
	// housekeeping
	deckObject.DeckIDs = append(deckObject.DeckIDs, cardIdNum)
	deckObject.ContainedObjects = append(deckObject.ContainedObjects, referenceCard)
	return nil
}

func (deckObject *Deck) RemoveCard(cardIdx int) error {
	var cardID = deckObject.ContainedObjects[cardIdx].CardID
	var idxToRemove = -1
	for idx, value := range deckObject.DeckIDs {
		if value == cardID {
			idxToRemove = idx
			break
		}
	}
	if idxToRemove == -1 {
		return errors.New("card not found in DeckIDs")
	} else {
		deckObject.DeckIDs = append(deckObject.DeckIDs[:idxToRemove], deckObject.DeckIDs[idxToRemove+1:]...)
	}
	deckObject.ContainedObjects = append(deckObject.ContainedObjects[:cardIdx], deckObject.ContainedObjects[cardIdx+1:]...)
	return nil
}

// AddCardsToDeck adds a card to the deck of a save file. The patch is incomplete as it does not include
// card id, deck id, as well as the card's position in the deck.
// therefore the operation is not idempotent at this level as the same card can be added multiple times
func (tt *Tabletop) AddCardsToDeck(luaScript string, luaScriptState string, savePath string, deckPath string) (error, []byte) {
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
	err = deckObject.AppendNewCard(luaScript, luaScriptState)
	if err != nil {
		return err, nil
	}
	return tt.finalizeChanges(deckObject, deckMarshalled, deckPath, originalBytes, savePath)
}

func (tt *Tabletop) RemoveCardFromDeck(savePath string, cardPath string) (error, []byte) {
	// parse cardPath to determine deck and card location within the deck
	var deckPath = cardPath[:strings.LastIndex(cardPath, "/ContainedObjects/")]
	cardIdx, err := strconv.Atoi(cardPath[strings.LastIndex(cardPath, "/")+1:])
	if err != nil {
		return err, nil
	}
	// obtain savefile
	originalBytes, err := tt.GetFile(savePath)
	if err != nil {
		return err, nil
	}
	// initialise deck
	var deckObject Deck
	deckMarshalled, err := getDeckFromSavefile(originalBytes, deckPath, &deckObject)
	if err != nil {
		return err, nil
	}
	// find CardID
	err = deckObject.RemoveCard(cardIdx)
	if err != nil {
		return err, nil
	}
	return tt.finalizeChanges(deckObject, deckMarshalled, deckPath, originalBytes, savePath)
}

func getDeckFromSavefile(originalBytes []byte, deckPath string, deckObj *Deck) ([]byte, error) {
	// unmarshal with loose typing in order to be able to navigate to the deck
	var originalFile interface{}
	err := json.Unmarshal(originalBytes, &originalFile)
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

func (tt *Tabletop) finalizeChanges(deckObject Deck, originalDeck []byte, deckPath string, originalSavefile []byte, savePath string) (error, []byte) {
	updatedDeck, err := json.Marshal(deckObject)
	if err != nil {
		return err, nil
	}
	patch, err := createDeckJsonPatch(originalDeck, updatedDeck, deckPath, "replace")
	if err != nil {
		return err, nil
	}
	// implement changes and return the updated file
	updated, err := tt.PatchSaveFile(savePath, originalSavefile, patch)
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
