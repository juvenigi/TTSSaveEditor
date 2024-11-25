package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

type URLEntry struct {
	JsonPath string   `json:"jsonpath"`
	GUID     string   `json:"GUID"`
	Names    []string `json:"Names"`
	Value    string   `json:"value"`
}

type URLEntryList struct {
	SavefilePath string     `json:"savefilePath"`
	CachePath    string     `json:"cachePath"`
	Entries      []URLEntry `json:"entries"`
}

// GetSaveResources TODO: check cache for file existence
func GetSaveResources(savefilePath string, cachePath string) ([]byte, error) {
	fileBytes, err := GetSaveJson(savefilePath)
	if err != nil {
		return nil, err
	}

	var root interface{}
	if err := json.Unmarshal(fileBytes, &root); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var entries []URLEntry
	var dummyEntry = URLEntry{
		JsonPath: "",
		GUID:     "",
		Names:    nil,
		Value:    "",
	}
	err = traverseJSON(root, "", dummyEntry, &entries)
	if err != nil {
		return nil, err
	}

	result := URLEntryList{
		SavefilePath: savefilePath,
		CachePath:    cachePath,
		Entries:      entries,
	}
	log.Print(result)
	return json.Marshal(result)
}

// TODO add maxDepth condition to avoid evil json
func traverseJSON(node interface{}, path string, currentEntry URLEntry, entries *[]URLEntry) error {
	localEntry := currentEntry
	var urls []string
	var lastPathSegment []string

	switch val := node.(type) {
	case map[string]interface{}:
		// current-level inspection
		for key, value := range val {
			// Update inherited properties in the local copy
			if key == "GUID" {
				if guidValue, ok := value.(string); ok {
					localEntry.GUID = guidValue
				}
			} else if key == "Nickname" || key == "Name" {
				if nameValue, ok := value.(string); ok {
					if len(nameValue) > 0 {
						localEntry.Names = append(localEntry.Names, nameValue)
					}
				}
			} else if strings.HasSuffix(key, "URL") {
				if urlValue, ok := value.(string); ok {
					if len(urlValue) > 0 {
						urls = append(urls, urlValue)
						lastPathSegment = append(lastPathSegment, key)
					}
				}
			}
		}
		// prove lower level
		for key, value := range val {
			if _, ok := value.(string); !ok {
				// could be a node, attempt to traverse it
				currentPath := fmt.Sprintf("%s/%s", path, key)
				if err := traverseJSON(value, currentPath, localEntry, entries); err != nil {
					return err
				}
			}
		}
		// add all urls belonging to this level
		if len(lastPathSegment) != len(urls) {
			return errors.New("something went wrong")
		}
		for i, urlValue := range urls {
			newEntry := localEntry // Copy inherited properties
			newEntry.JsonPath = fmt.Sprintf("%s/%s", path, lastPathSegment[i])
			newEntry.Value = urlValue
			*entries = append(*entries, newEntry)
		}
	case []interface{}: // check array
		for idx, arrayItem := range val {
			if err := traverseJSON(arrayItem, fmt.Sprintf("%s/%d", path, idx), localEntry, entries); err != nil {
				return err
			}
		}
	}

	return nil
}
