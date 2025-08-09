package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var urls = regexp.MustCompile("URL$")

var resourcePattern = regexp.MustCompile(`^(https?)|(file)://.*`)

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
func AggregateBundleRecur(unmarshalledJsonNode any, jsonPointer string, result map[string]Bundle) *Bundle {
	partial := &Bundle{ResUrls: make(map[string]string)}

	switch v := unmarshalledJsonNode.(type) {
	case map[string]interface{}:
		for key, val := range v {
			partial.reconcile(AggregateBundleRecur(val, jsonPointer+"."+key, result))
		}
	case []interface{}:
		for i, val := range v {
			partial.reconcile(AggregateBundleRecur(val, fmt.Sprintf("%s[%d]", jsonPointer, i), result))
		}
	case string:
		for _, key := range metaKeys {
			if strings.HasSuffix(jsonPointer, key) {
				_ = partial.AddMetaKey(key, v)
				break
			}
		}

		if urls.MatchString(jsonPointer) && resourcePattern.MatchString(strings.ToLower(v)) && len(jsonPointer) > 0 {
			partial.ResUrls[jsonPointer[1:]] = v
		}
	}

	if partial.isAggregatable() && len(jsonPointer) > 0 {
		result[jsonPointer[1:]] = *partial
		return nil
	} else {
		// pass downstream
		return partial
	}
}
