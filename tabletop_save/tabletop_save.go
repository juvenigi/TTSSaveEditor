package tabletop_save

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

func (gr *GameResource) createJsonPatch() string {
	patchPath := strings.ReplaceAll(gr.JsonPath, ".", "/")

	return fmt.Sprintf(`{"op":"replace","path":"%s", "value":"%s"}`, patchPath, gr.ResourceUrl)
}

func (s *TSSaveFile) SaveAsPackTemplate() error {
	file, err := os.ReadFile(s.savefileLocation)
	if err != nil {
		return err
	}

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
	patch, err := jsonpatch.DecodePatch(finalPatch)
	if err != nil {
		return err
	}
	modified, err := patch.Apply(file)
	if err != nil {
		return err
	}
	log.Printf("%s\n", modified)

	return nil
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

type ResourceBundle struct {
	jsonPath  string
	Name      string
	Nickname  string
	Guid      string
	resources []GameResource
}

type GameResource struct {
	JsonPath    string
	ResourceUrl string
	Status      ResourceStatus
}

type ResourceStatus int

const (
	Failed ResourceStatus = iota
	Local
	Remote
	RemoteCached
	Steam
	Packed // this is needed to prevent redundant 're-packing'
)

func normalize(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func ParseResourcesBundles(blob []byte, packDir string) ([]ResourceBundle, error) {
	var data any

	bundleMap := make(map[string]Bundle)

	if err := json.Unmarshal(blob, &data); err != nil {
		return nil, err
	}

	_ = walkJsonRecur2(data, "", bundleMap)

	var results []ResourceBundle
	for k, bundle := range bundleMap {
		results = append(results, bundle.MapToResourceBundle(k, packDir))
	}

	return results, nil
}

func (bb *Bundle) MapToResourceBundle(jsonPath string, packDataDir string) ResourceBundle {
	return ResourceBundle{
		jsonPath:  jsonPath,
		Guid:      normalize(bb.Guid),
		Name:      normalize(bb.Name),
		Nickname:  normalize(bb.Nickname),
		resources: bb.MapToResources(packDataDir),
	}
}

func (bb *Bundle) MapToResources(packDataDir string) []GameResource {
	var urls []GameResource
	for k, v := range bb.ResUrls {
		urls = append(urls, GameResource{
			JsonPath:    k,
			ResourceUrl: v,
			Status:      deduceResourceStatus(v, packDataDir),
		})
	}
	return urls
}

type GameCache struct {
	ResourceCached map[string]string // stores cached filenames -> complete path
}

func NewCacheScanner(gameDir string) (GameCache, error) {
	var result GameCache
	cachedEntries := make(map[string]string)

	err := filepath.WalkDir(gameDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			name := d.Name()
			name = strings.TrimSuffix(name, filepath.Ext(name))

			cachedEntries[name] = path
		}
		return nil
	})
	if err != nil {
		return result, err
	}

	result.ResourceCached = cachedEntries
	return result, nil
}

func (rv *GameCache) SetCacheStatus(res *GameResource) {

	if res.Status == Local || res.Status == RemoteCached || res.Status == Packed {
		return
	}

	base := filepath.Base(res.ResourceUrl)
	cacheFilename := GetCacheFilename(base)
	if _, ok := rv.ResourceCached[cacheFilename]; ok {
		res.Status = RemoteCached
	}

	return
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

var urls = regexp.MustCompile("URL$")

var resourcePattern = regexp.MustCompile(`^(https?)|(file)://.*`)

var metaKeys = []string{"Name", "Nickname", "GUID"}

type Bundle struct {
	Guid     *string
	Name     *string
	Nickname *string
	ResUrls  map[string]string
}

func NewBundle() *Bundle {
	return &Bundle{ResUrls: make(map[string]string)}
}

func (bb *Bundle) reconcile(b *Bundle) {
	if b == nil {
		return
	}

	if b.Nickname != nil {
		bb.Nickname = b.Nickname
	}
	if b.Name != nil {
		bb.Name = b.Name
	}
	if b.Guid != nil {
		bb.Guid = b.Guid
	}
	if b.ResUrls != nil {
		for k, v := range b.ResUrls {
			bb.ResUrls[k] = v
		}
	}
}

func (bb *Bundle) AddMetaKey(metaKey string, metaValue string) error {
	switch metaKey {
	case "Name":
		bb.Name = &metaValue
		return nil
	case "Nickname":
		bb.Nickname = &metaValue
		return nil
	case "GUID":
		bb.Guid = &metaValue
		return nil
	default:
		return errors.New("invalid meta key")
	}
}

func (bb *Bundle) isAggregatable() bool {
	hasIdentifier := bb.Guid != nil || bb.Name != nil || bb.Nickname != nil
	hasMap := len(bb.ResUrls) > 0

	return hasMap && hasIdentifier
}

// returns a bundle if it's a partial result
func walkJsonRecur2(node interface{}, path string, result map[string]Bundle) *Bundle {
	partial := NewBundle()

	switch v := node.(type) {
	case map[string]interface{}:
		for key, val := range v {
			partial.reconcile(walkJsonRecur2(val, path+"."+key, result))
		}
	case []interface{}:
		for i, val := range v {
			partial.reconcile(walkJsonRecur2(val, fmt.Sprintf("%s[%d]", path, i), result))
		}
	case string:
		for _, key := range metaKeys {
			if strings.HasSuffix(path, key) {
				_ = partial.AddMetaKey(key, v)
				break
			}
		}

		if urls.MatchString(path) && resourcePattern.MatchString(strings.ToLower(v)) && len(path) > 0 {
			partial.ResUrls[path[1:]] = v
		}
	}

	if partial.isAggregatable() && len(path) > 0 {
		result[path[1:]] = *partial
		return nil
	} else {
		// pass downstream
		return partial
	}
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
