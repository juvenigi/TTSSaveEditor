package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type JsonPatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func UnmarshalJsonPatch(opname string, data []byte) ([]JsonPatchOp, error) {
	var patch []JsonPatchOp
	var mergePatch map[string]interface{}
	err := json.Unmarshal(data, &mergePatch)
	if err != nil {
		return nil, err
	}
	for key, value := range mergePatch {
		patch = append(patch, JsonPatchOp{Op: opname, Path: key, Value: value})
	}
	return patch, nil
}

// UpdateDeckPath finds the number of elements in the array, then update the path to make it become the last new element
func (self *JsonPatchOp) UpdateDeckPath(baseJson []byte) error {
	// get the number of elements in the array
	count, err := jsonArrayLen(baseJson, "/data/decks")
	if err != nil {
		return err
	}
	// set the path to the last element
	self.Path = fmt.Sprintf("/data/decks/%d", count)
	return nil
}

// navigateToPath navigates a JSON document using a path and returns the element at that path.
func navigateToPath(json interface{}, path string) (interface{}, error) {
	// Split the path into parts
	parts := strings.Split(path, "/")[1:] // omits the leading slash
	var currentNode = json
	// Navigate through the json by following the path
	for _, part := range parts {
		switch cur := currentNode.(type) {
		case map[string]interface{}:
			// Navigate into objects
			currentNode = cur[part]
		case []interface{}:
			// Navigate into arrays
			index, err := strconv.Atoi(part)
			length := len(cur)
			if err != nil || index < 0 || index >= length {
				return nil, fmt.Errorf("invalid array index: %s, array size: %d", part, length)
			}
			currentNode = cur[index]
		default:
			return nil, fmt.Errorf("invalid path: %s", path)
		}
	}
	return currentNode, nil
}

func remarshal[T any](jsonObject interface{}, target *T) (error, []byte) {
	data, err := json.Marshal(jsonObject)
	if err != nil {
		return err, nil
	}
	err = json.Unmarshal(data, target)
	if err != nil {
		return err, nil
	}
	return nil, data
}

// jsonArrayLen counts the elements in the array at the given JSON Pointer path.
func jsonArrayLen(jsonDoc []byte, path string) (int, error) {
	var doc interface{}
	if err := json.Unmarshal(jsonDoc, &doc); err != nil {
		return 0, err
	}

	element, err := navigateToPath(doc, path)
	if err != nil {
		return 0, err
	}

	// Assert that the element is an array and count its elements.
	if arr, ok := element.([]interface{}); ok {
		return len(arr), nil
	}

	return 0, fmt.Errorf("element at path %s is not an array", path)
}
