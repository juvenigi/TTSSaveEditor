package tabletop_save

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"tts-cache-manager-cli/backend/internal/internal/resource_map"
)

// ResourceBundle todo: JsonPointer in ResourceBundle is currently unused -- remove
type ResourceBundle struct {
	JsonPointer string
	Name        string
	Nickname    string
	Guid        string
	resources   []GameResource
}

type GameResource struct {
	JsonPointer string
	ResourceUrl string
	Status      ResourceStatus
}

func (res *GameResource) UpdateToPacked(written string) {
	res.ResourceUrl = "pack://" + written
	res.Status = Packed
}

// ResourceStatus `RemoteCached` is a subset of `Remote` from a philosophical point of view, however
// it is not initialized from json save data for an obvious reason that it does not tell you if your resource is cached
// same goes for `Packed` as it required knowledge of PackData content
// `Steam` can become RemoteCached if the cached version of the resource exists on the disk
// PackPath : a local reference to a Packed file
type ResourceStatus int

const (
	Failed ResourceStatus = iota
	Local
	Remote
	RemoteCached
	Steam
	Packed
	PackPath
	PackCached
)

const steamApiUrlPrefix = "https://steamusercontent"

func deduceResourceStatus(url string, packData resource_map.PackData) ResourceStatus {
	if strings.HasPrefix(url, "file:///") {
		if strings.HasPrefix(url[fileLen:], packData.PackDataDir) {
			return PackPath
		} else {
			return Local
		}
	} else if strings.HasPrefix(url, steamApiUrlPrefix) {
		return Steam
	} else if _, present := packData.UrlArchive[url]; present {
		return Packed
	} else {
		return Remote
	}
}

var protocolPrefix = regexp.MustCompile(`^[^:]+:///?`)

func (res *GameResource) createJsonPatch() string {
	patchPath := strings.ReplaceAll(res.JsonPointer, ".", "/")
	url := res.ResourceUrl

	separatorString := string(filepath.Separator)
	var escapedUrl string
	if strings.Compare(separatorString, "\\") == 0 {
		idx := 0
		if match := protocolPrefix.FindStringIndex(url); match != nil {
			idx = match[1]
		}

		protocol := url[:idx]
		url = url[idx:]

		escapedSeparator := "\\" + separatorString

		escapedUrl = protocol + strings.ReplaceAll(url, separatorString, escapedSeparator)
	} else {
		escapedUrl = url
	}

	return fmt.Sprintf(`{"op":"replace","path":"%s", "value":"%s"}`, patchPath, escapedUrl)
}
