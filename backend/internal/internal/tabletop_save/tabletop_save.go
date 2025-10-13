package tabletop_save

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"tts-cache-manager-cli/backend/internal/internal/resource_map"
	"tts-cache-manager-cli/backend/internal/internal/tabletop_save/internal"
	"tts-cache-manager-cli/util/base52"

	"github.com/bwmarrin/snowflake"
	jsonpatch "github.com/evanphx/json-patch"
)

type TSSaveFile struct {
	savefileLocation string
	savename         string
	revision         SaveRevision
	resourceBundle   []ResourceBundle
}

type SaveFileView struct {
	Filename              string `json:"filename"`
	Savename              string `json:"savename"`
	Directory             string `json:"directory"`
	PackId                string `json:"pack_id"`
	PackRevision          int    `json:"pack_revision"`
	ObjectCount           int    `json:"objectCount"`
	UncachedResources     int    `json:"uncachedResources"`
	CachedRemoteResources int    `json:"cachedRemoteResources"`
	LocalResources        int    `json:"localResources"`
	PackedResources       int    `json:"packedResources"`
}

type TagsAndSaveNamePartialJson struct {
	Tags     []string
	SaveName string
}

// GetTabletopSaveFile note: `seed` param is nillable
func GetTabletopSaveFile(loc string, packData resource_map.PackData) (*TSSaveFile, error) {
	blobBytes, err := os.ReadFile(loc)
	if err != nil {
		return nil, err
	}

	var unmarshalledData any
	if err := json.Unmarshal(blobBytes, &unmarshalledData); err != nil {
		return nil, err
	}

	var tagsAndSave TagsAndSaveNamePartialJson
	err = json.Unmarshal(blobBytes, &tagsAndSave)
	if err != nil {
		return nil, err
	}

	var revision SaveRevision
	if tagsAndSave.Tags != nil {
		if err := revision.Parse(tagsAndSave.Tags); err != nil {
			return nil, err
		}
	}

	bundles, err := ParseResourcesBundles(unmarshalledData, packData)
	if err != nil {
		return nil, err
	}

	return &TSSaveFile{
		savefileLocation: loc,
		savename:         tagsAndSave.SaveName,
		revision:         revision,
		resourceBundle:   bundles,
	}, nil
}

// GetAllResources note: may contain duplicates
func (s *TSSaveFile) GetAllResources() []*GameResource {
	var resources []*GameResource
	for _, bundle := range s.resourceBundle {
		for i := range bundle.resources {
			resources = append(resources, &bundle.resources[i])
		}
	}
	return resources
}

func (s *TSSaveFile) IntoPortable(seed *snowflake.Node, saveName string) ([]byte, error) {
	s.revision.IncrementOrGen(seed)

	file, err := os.ReadFile(s.savefileLocation)
	if err != nil {
		return nil, err
	}

	patchSet, err := jsonpatch.DecodePatch(s.GetJsonPatchBytesForPacked(false))
	if err != nil {
		return nil, err
	}
	if saveName != "" {
		namePatch, err := jsonpatch.DecodePatch([]byte(fmt.Sprintf(`[{"op":"replace","path":"/SaveName","value":"%s"}]`, saveName)))
		if err != nil {
			return nil, err
		}
		patchSet = append(patchSet, namePatch...)
	}
	modified, err := patchSet.Apply(file)
	if err != nil {
		return nil, err
	}

	return modified, nil
}

type SaveFileInfoJson struct {
	Name string
}

var tsSavePattern = regexp.MustCompile("^TS_Save_(\\d+)")

func getLargestSaveFileName(gameDir string) (string, error) {
	largest := 0
	err := filepath.WalkDir(filepath.Join(gameDir, "Saves")+string(filepath.Separator), func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(d.Name(), ".json") && !d.IsDir() {
			matches := tsSavePattern.FindStringSubmatch(d.Name())
			if len(matches) > 0 {
				if candidate, err := strconv.Atoi(matches[1]); err == nil && largest < candidate {
					largest = candidate
				} else {
					return err
				}
			}
		}
		return err
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("TS_Save_%d.json", largest+1), nil
}

func (s *TSSaveFile) WriteNewSaveToSavesDir(gameDir string, saveName string) error {
	patchset := s.GetJsonPatchBytesForPacked(true)
	patch, err := jsonpatch.DecodePatch(patchset)
	if err != nil {
		return err
	}
	if len(saveName) > 0 {
		saveNamePatch, err := jsonpatch.DecodePatch([]byte(fmt.Sprintf(`[{"op":"replace","path":"/SaveName","value":"%s"}]`, saveName)))
		if err != nil {
			return err
		}
		patch = append(patch, saveNamePatch...)
	}

	saveFileName, err := getLargestSaveFileName(gameDir)
	if err != nil {
		return err
	}
	srcBlob, err := os.ReadFile(s.savefileLocation)
	if err != nil {
		return err
	}
	dest, err := os.Create(filepath.Join(gameDir, "Saves", saveFileName))
	if err != nil {
		return err
	}
	defer dest.Close()

	patched, err := patch.Apply(srcBlob)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(patched)
	if _, err := io.Copy(dest, reader); err != nil {
		return err
	}
	if err := dest.Sync(); err != nil {
		return err
	}

	return nil
}

func (s *TSSaveFile) GetJsonPatchBytesForPacked(keepPaths bool) []byte {
	revisionBytes := s.revision.Serialize()

	var patches []string
	allResources := s.GetAllResources()
	for _, res := range allResources {
		//if !keepPaths && res.Status == PackPath {
		//	res.ResourceUrl = "pack://" + filepath.Base(res.ResourceUrl)
		//	res.Status = Packed
		//}
		//
		//if res.Status != Packed {
		//	continue
		//}
		patches = append(patches, res.createJsonPatch())
	}
	var sb strings.Builder
	sb.WriteString("[")
	if revisionBytes != nil {
		sb.WriteString(string(revisionBytes))
		sb.WriteString(",")
	}
	for _, patch := range patches {
		sb.WriteString(patch)
		sb.WriteString(",")
	}

	concatenated := []byte(sb.String())
	if len(concatenated) > 1 {
		// remove trailing comma
		return append(concatenated[:len(concatenated)-1], []byte("]")...)
	} else {
		return append(concatenated, []byte("]")...)
	}
}

func (s *TSSaveFile) ToView() *SaveFileView {
	objectCount := len(s.resourceBundle)
	uncachedRes := 0
	remoteCached := 0
	localResources := 0
	packedResources := 0
	for _, resourceBundle := range s.resourceBundle {
		for _, res := range resourceBundle.resources {
			switch res.Status {
			case Local:
				localResources++
			case Remote:
				uncachedRes++
			case PackPath:
				packedResources++
			case Packed:
				packedResources++
			case Steam:
				uncachedRes++
			case RemoteCached:
				remoteCached++
			default:
			}
		}
	}

	return &SaveFileView{
		Filename:              filepath.Base(s.savefileLocation),
		Savename:              s.savename,
		Directory:             filepath.Dir(s.savefileLocation),
		PackId:                s.revision.root,
		PackRevision:          s.revision.revision,
		ObjectCount:           objectCount,
		UncachedResources:     uncachedRes,
		CachedRemoteResources: remoteCached,
		LocalResources:        localResources,
		PackedResources:       packedResources,
	}
}

func (s *TSSaveFile) GetSaveName() string {
	return filepath.Base(s.savefileLocation)
}

func normalize(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

const fileLen = len("file:///")

func ParseResourcesBundles(unmarshalledJson any, packData resource_map.PackData) ([]ResourceBundle, error) {

	bundleMap := make(map[string]internal.Bundle)

	_ = internal.AggregateBundleRecur(unmarshalledJson, new(internal.JsonPointer), bundleMap)

	var results []ResourceBundle
	for k, bundle := range bundleMap {
		results = append(results, ResourceBundle{
			JsonPointer: k,
			Guid:        normalize(bundle.Guid),
			Name:        normalize(bundle.Name),
			Nickname:    normalize(bundle.Nickname),
			resources:   MapToResources(&bundle, packData),
		})
	}

	return results, nil
}

func MapToResources(bb *internal.Bundle, packData resource_map.PackData) []GameResource {
	var urls []GameResource
	for pointer, url := range bb.ResUrls {
		urls = append(urls, GameResource{
			JsonPointer: pointer,
			ResourceUrl: url,
			Status:      deduceResourceStatus(url, packData),
		})
	}
	return urls
}

type SaveRevision struct {
	original []string
	root     string
	revision int
}

var saveRevPattern = regexp.MustCompile("__packrev([^_]+)_?(.*)$")

func (rev *SaveRevision) IsBlank() bool {
	return len(rev.root) == 0
}

func (rev *SaveRevision) Parse(valueFromJson []string) error {
	dollarIdx := slices.IndexFunc(valueFromJson, func(s string) bool {
		return strings.HasPrefix(s, "__packrev")
	})
	if dollarIdx != -1 {
		submatches := saveRevPattern.FindStringSubmatch(valueFromJson[dollarIdx])
		rev.original = append(valueFromJson[:dollarIdx], valueFromJson[dollarIdx+1:]...)
		if len(submatches) != 3 {
			return nil
		}
		rev.root = submatches[1]
		if len(submatches[2]) > 0 {
			var err error
			int64Decode, err := base52.Decode(submatches[2])
			if err != nil {
				return err
			}
			rev.revision = int(int64Decode)
			return err
		}

	} else {
		rev.original = valueFromJson
	}
	return nil
}

func (rev *SaveRevision) Serialize() []byte {
	packrevStr, err := rev.String()
	if err != nil {
		return nil
	}
	allTags := append(rev.original, packrevStr)
	marshal, err := json.Marshal(allTags)
	if err != nil {
		return nil
	}
	return []byte(fmt.Sprintf(`{"op":"replace","path":"/Tags","value":%s}`, marshal))
}

func (rev *SaveRevision) String() (string, error) {
	if len(rev.root) == 0 {
		return "", errors.New("invalid save revision")
	}
	if rev.revision == 0 {
		return fmt.Sprintf("__packrev%s", rev.root), nil
	}
	encode, err := base52.Encode(int64(rev.revision))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("__packrev%s_%s", rev.root, encode), nil
}

func (rev *SaveRevision) IncrementOrGen(seed *snowflake.Node) {
	if len(rev.root) == 0 {
		rev.revision = 0
		encode, err := base52.Encode(seed.Generate().Int64())
		if err != nil {
			panic(err)
		}
		rev.root = encode
	} else {
		rev.revision++
	}
}
