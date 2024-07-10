# General
- Card / TTS Object editing inside the `SaveData` (Tabletop Simulator JSON save file).
- REST API for savefile push / pull
  - editing save file
  - hot-reloading scripts

# Design principles

- broker is the single source of truth
    - it knows if TTS is running a 'stale' file
    - it notifies subscribers of new changes
    - the subscribers need to handle updates: discard changes, override, create copy

Diagram:
<editor: web, intellij> -- <golang broker> -- <TTS>

# Journal

## 10.07.2024

## progress
- learned ngrx and angular18 shenanigans
- refactored most of the ugly code

## complications
- lack of caching is annoying
- lack of default directory is more annoying
- `BUG` hardcoded path delimiter (currently using Windows-style "\\" )
- `WARN` tabletop simulator save file is not validated (meaning that SaveData types are fictional / unverified)

## tasks
- [ ] return default directory if no path is specified
- [ ] list bags/decks of a savefile
- [ ] filter objects by bag
- [ ] create a formGroup for a card
- [ ] wire IO events
- [ ] quill.js textarea?
- [ ] backend: patch savefile json (individual cards, entire decks)
```go
package main 

import (
"fmt"
jsonpatch "github.com/evanphx/json-patch"
"log"
)

func main() {
	original := []byte(`{"nested": [1,2,3,4]}`)
	patchJSON := []byte(`[
		
		{"op": "add", "path": "/nested/1", "value": {"different":1}}
	]`)

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		panic(err)
	}

	modified, err := patch.Apply(original)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Original document: %s\n", original)
	fmt.Printf("Modified document: %s\n", modified)
}
```
## 06.07.2024

### progress
- load list of tts save directories (backend + frontend)
- basic router config & fire an effect

### complications
- still rusty with ngrx
- got myself into a pickle: I think I should just fire an event with a savefile path, then find the tab id based on it
  - reducer simply checks if the savefile exists and sets the status to loading
  - effect does the heavy lifting: check tabs, load data, append new tab (preferably via firing actions)

### tasks
- [x] don't create too generalised reducers, split them up if necessary
- [x] navigation to savefile (actions, effect, reducer, query)
- [x] expand model to include the actual tts savefile
