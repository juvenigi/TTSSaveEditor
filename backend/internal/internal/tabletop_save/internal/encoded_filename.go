package internal

import (
	"regexp"
)

var excludeChars = regexp.MustCompile(`[^a-zA-Z0-9]`)

// GetCacheFilename note: this trims file extension
// TTS decides to turn `file.png` into `filepng.png` when caching, therefore we need to account this in our code.
// this is why I chose to remove the file extension completely if it is present
// note2: sometimes people put dropbox links which work like this: `https://www.dropbox.com/s/<id>/image.png?dl=1`
// this breaks a lot of previously held assumptions, meaning I do need to handle the entire input string
func GetCacheFilename(raw string) string {
	return excludeChars.ReplaceAllString(raw, "")
}
