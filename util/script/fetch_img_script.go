//go:build ignore
// +build ignore

package main

import "tts-cache-manager-cli/util"

func main() {
	util.FetchImgFromYaml("../../resources/img", "imglist.yaml")
}
