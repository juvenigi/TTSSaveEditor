package tabletop_save

import (
	"fmt"
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

func (gr *GameResource) createJsonPatch() string {
	patchPath := strings.ReplaceAll(gr.JsonPointer, ".", "/")

	return fmt.Sprintf(`{"op":"replace","path":"%s", "value":"%s"}`, patchPath, gr.ResourceUrl)
}
