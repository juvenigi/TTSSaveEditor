package tabletop_save

import (
	"ReallyDumbCopyPaste/util"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TabletopSave struct {
	saveName  string
	resources []TabletopResource
}

type TabletopResource struct {
	JsonPointer string
	ResourceUrl string
}

type tabletopSaveJsonProjection struct {
	SaveName string
}

func (t *TabletopSave) Init(location string) error {
	blobBytes, err := os.ReadFile(location)
	if err != nil {
		return err
	}

	if err = t.readSaveName(err, blobBytes); err != nil {
		return err
	}

	var unmarshalledData any
	if err := json.Unmarshal(blobBytes, &unmarshalledData); err != nil {
		return err
	}

	t.extractResourceRecur(unmarshalledData, new(util.JsonPointer))

	return nil
}

func (t *TabletopSave) readSaveName(err error, blobBytes []byte) error {
	var saveNameSlice tabletopSaveJsonProjection
	if err = json.Unmarshal(blobBytes, &saveNameSlice); err != nil {
		return err
	}
	t.saveName = saveNameSlice.SaveName
	return nil
}

func (t *TabletopSave) extractResourceRecur(unmarshalledJson any, cursor *util.JsonPointer) {
	switch node := unmarshalledJson.(type) {
	case map[string]interface{}:
		for k, v := range node {
			pointerClone := cursor.Clone()
			pointerClone.Append(k)
			t.extractResourceRecur(v, pointerClone)
		}
	case []interface{}:
		for i, v := range node {
			pointerClone := cursor.Clone()
			pointerClone.Append(fmt.Sprintf("%d", i))
			t.extractResourceRecur(v, pointerClone)
		}
	case string:
		build := cursor.BuildPointer()
		if strings.HasSuffix(build, "URL") {
			t.resources = append(t.resources, TabletopResource{
				JsonPointer: build,
				ResourceUrl: node,
			})
		}
	}
}

func (t *TabletopSave) Patch(cache *GameCacheFinder) {
	for _, res := range t.resources {
		cache.MarkAsUsed(res.ResourceUrl)
	}
}
