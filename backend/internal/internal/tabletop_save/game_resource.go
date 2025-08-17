package tabletop_save

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
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

// ResourceStatus `RemoteCached` is a subset of `Remote` from a philosophical point of view, however
// it is not initialized from json save data for an obvious reason that it does not tell you if your resource is cached
// same goes for `Packed` as it required knowledge of PackData content
// `Steam` can become RemoteCached if the cached version of the resource exists on the disk
type ResourceStatus int

const (
	Failed ResourceStatus = iota
	Local
	Remote
	RemoteCached
	Steam
	Packed
)

var protocolPrefix = regexp.MustCompile(`^[^:]+:///?`)

func (gr *GameResource) createJsonPatch() string {
	patchPath := strings.ReplaceAll(gr.JsonPointer, ".", "/")
	url := gr.ResourceUrl

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
