package tabletop_save

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"tts-cache-manager-cli/backend/internal/internal/tabletop_save/internal"

	jsonpatch "github.com/evanphx/json-patch"
)

const steamApiUrlPrefix = "https://steamusercontent"

// TSSaveFile todo: use bundles
type TSSaveFile struct {
	savefileLocation string
	resourceBundle   []ResourceBundle
}

type SaveFileView struct {
	Filename              string `json:"filename"`
	Directory             string `json:"directory"`
	ObjectCount           int    `json:"objectCount"`
	UncachedResources     int    `json:"uncachedResources"`
	CachedRemoteResources int    `json:"cachedRemoteResources"`
	LocalResources        int    `json:"localResources"`
	PackedResources       int    `json:"packedResources"`
}

// note: may contain duplicates
func (s *TSSaveFile) GetAllResources() []*GameResource {
	var resources []*GameResource
	for _, bundle := range s.resourceBundle {
		for i := range bundle.resources {
			resources = append(resources, &bundle.resources[i])
		}
	}
	return resources
}

func (s *TSSaveFile) PutPackUrls() ([]byte, error) {
	file, err := os.ReadFile(s.savefileLocation)
	if err != nil {
		return nil, err
	}

	finalPatch := s.getJsonPatchBytes()
	patch, err := jsonpatch.DecodePatch(finalPatch)
	if err != nil {
		return nil, err
	}
	modified, err := patch.Apply(file)
	if err != nil {
		return nil, err
	}
	log.Printf("%s\n", modified)

	return modified, nil
}

func (s *TSSaveFile) getJsonPatchBytes() []byte {
	var patches []string
	for _, res := range s.GetAllResources() {
		if res.Status != Packed {
			continue
		}
		patches = append(patches, res.createJsonPatch())
	}
	var sb strings.Builder
	sb.WriteString("[")
	for _, patch := range patches {
		sb.WriteString(patch)
		sb.WriteString(",")
	}
	concatenated := []byte(sb.String())
	finalPatch := append(concatenated[:len(concatenated)-1], []byte("]")...)

	return finalPatch
}

func (s *TSSaveFile) ToView() SaveFileView {
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

	return SaveFileView{
		Filename:              filepath.Base(s.savefileLocation),
		Directory:             filepath.Dir(s.savefileLocation),
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

func deduceResourceStatus(url string, packDir string) ResourceStatus {
	if strings.HasPrefix(url, packDir) {
		return Packed
	} else if strings.HasPrefix(url, steamApiUrlPrefix) {
		return Steam
	} else if strings.HasPrefix(url, "http") {
		return Remote
	} else if strings.HasPrefix(url, "file:///") {
		return Local
	}

	return Failed
}

func getCachedResourceLoc(gameDir string, resourceUrl string) string {
	encodedFilename := GetCacheFilename(filepath.Base(resourceUrl))
	var cachedEntry = ""

	found := errors.New("found") // a slightly hacky approach to short-circuit fs walk after the file is found
	err := filepath.WalkDir(gameDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			sanitizedName := d.Name()
			sanitizedName = strings.TrimSuffix(sanitizedName, filepath.Ext(sanitizedName))
			if strings.Compare(encodedFilename, sanitizedName) == 0 {
				cachedEntry = path
				return found
			}
		}
		return nil
	})
	if err != nil && !errors.Is(err, found) {
		log.Println("error while searching for cached file:", err)
		return ""
	}

	return cachedEntry
}

func ParseResourcesBundles(blob []byte, packDir string) ([]ResourceBundle, error) {
	var data any
	if err := json.Unmarshal(blob, &data); err != nil {
		return nil, err
	}

	results := GetResourcesFromUnmarshalledJson(data, packDir)

	return results, nil
}

func GetResourcesFromUnmarshalledJson(data any, packDir string) []ResourceBundle {
	bundleMap := make(map[string]internal.Bundle)

	_ = internal.AggregateBundleRecur(data, "", bundleMap)

	var results []ResourceBundle
	for k, bundle := range bundleMap {
		results = append(results, MapToResourceBundle(&bundle, k, packDir))
	}
	return results
}

func MapToResourceBundle(bb *internal.Bundle, jsonPath string, packDataDir string) ResourceBundle {
	return ResourceBundle{
		JsonPointer: jsonPath,
		Guid:        normalize(bb.Guid),
		Name:        normalize(bb.Name),
		Nickname:    normalize(bb.Nickname),
		resources:   MapToResources(bb, packDataDir),
	}
}

func MapToResources(bb *internal.Bundle, packDataDir string) []GameResource {
	var urls []GameResource
	for k, v := range bb.ResUrls {
		urls = append(urls, GameResource{
			JsonPointer: k,
			ResourceUrl: v,
			Status:      deduceResourceStatus(v, packDataDir),
		})
	}
	return urls
}
