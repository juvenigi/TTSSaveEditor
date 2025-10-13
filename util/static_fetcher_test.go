package util

import (
	"testing"
)

func TestFetchImgFromYaml(t *testing.T) {
	FetchImgFromYaml("../../../resources/img", "imglist.yaml")
}
