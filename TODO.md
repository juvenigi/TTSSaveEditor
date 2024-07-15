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

## 13.07.2024

## progress
- none (did not have time or motivation)

## tasks
- [ ] create a formGroup for a card
  - font size / max characters autoformatter 
  - [ ] quill.js textarea? bbcode? -> use textarea + autodecreasing font size for now
- [ ] wire IO events
- [ ] improve backup functionality
  - [ ] check if the backup + patch == current savefile
  - [ ] keep a patchlist instead of saving snapshots (when patches do not result in the savefile, create a new snapshot)
- [ ] Card NGRX Store
  - Action: submit text / discard text -> (effect) use patch endpoint & sync states
  - Action: add new card -> (util) create new card, luascript management, guid generation
  - Action: remove card from deck
  - (you can now move cards between decks, since you have add and remove)

## 10.07.2024 (after sleeping)

## progress
- worked mostly on functionality

## complications
- I think my backend code is showing some ugly sides (sloppy error-checking)
- Add tests at some point

## tasks
- [X] return default directory if no path is specified
- [X] list bags/decks of a savefile
- [X] filter objects by bag

- [X] backend: patch savefile json (individual cards, entire decks)

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
