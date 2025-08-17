package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var urls = regexp.MustCompile("URL$")

var resourcePattern = regexp.MustCompile(`^(https?)|(file)|(pack).*`)

var metaKeys = []string{"Name", "Nickname", "GUID"}

type Bundle struct {
	Guid     *string
	Name     *string
	Nickname *string
	ResUrls  map[string]string
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

// returns a Bundle if it's a partial result
func AggregateBundleRecur(unmarshalledJsonNode any, jsonPointer *JsonPointer, result map[string]Bundle) *Bundle {
	partial := &Bundle{ResUrls: make(map[string]string)}
	pointerString := jsonPointer.BuildPointer()
	var pointerClone *JsonPointer
	switch v := unmarshalledJsonNode.(type) {
	case map[string]interface{}:
		for key, val := range v {
			pointerClone = jsonPointer.clone()
			pointerClone.Append(key)
			partial.reconcile(AggregateBundleRecur(val, pointerClone, result))
		}
	case []interface{}:
		for i, val := range v {
			pointerClone = jsonPointer.clone()
			pointerClone.Append(fmt.Sprintf("%d", i))
			partial.reconcile(AggregateBundleRecur(val, pointerClone, result))
		}
	case string:
		// a raw value
		for _, key := range metaKeys {
			if strings.HasSuffix(pointerString, key) {
				_ = partial.AddMetaKey(key, v)
				break
			}
		}

		// todo: I think that verifying the key or value is redundant here: I would simply check the value
		if urls.MatchString(pointerString) && resourcePattern.MatchString(strings.ToLower(v)) && len(pointerString) > 0 {
			partial.ResUrls[pointerString] = v
		}
	}

	if partial.isAggregatable() && len(pointerString) > 0 {
		result[pointerString] = *partial
		return nil
	} else {
		// pass downstream
		return partial
	}
}
