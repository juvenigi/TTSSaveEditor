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

# Skipped items

- [ ] SaveData JSON validation
    - (would be nice to avoid potential savefile corruption of the format ever changes / becomes incompatible)
- [ ] Rich text editing with BBCode
    - considering that I want to create a graphical card editor, this is not necessary at all
- [ ] Save Editor API integration

# Journal

## 16.07.2024

## progress

- started late again, because decided to drink beer instead
- you can now edit cards outside of tabletop simulator

## complications

- I highly suspect that I have overengineered the backend because I want to have integration with Save Editor API,
  something which I am only scaffolding

## tasks
- [X] create a formGroup for a card
- [X] wire IO events
- [X] Card NGRX Store (partially)
  - [X] Action: submit text / discard text -> (effect) use patch endpoint & sync states
- [X] small backend and frontend code quality improvements
- [X] Save Editor API listening port
- [X] bug: filters are not reset when opening a different savefile
- [ ] bug/qol: sort directory entries topologically
- [ ] bug: don't reset card form if changes are pending
- [ ] bug: don't discard card changes if the other deck is selected, make the edited card autosave
- [ ] bug: backup files don't go beyond 0.bak (if 0.bak exists, 1.bak should be created instead of overriding 0.bak)
- [ ] Card NGRX Store (remaining)
    - Action: submit text / discard text -> (effect) use patch endpoint & sync states
    - Action: add new card -> (util) create new card, luascript management, guid generation
    - Action: remove card from deck
    - (you can now move cards between decks, since you have add and remove)
- [ ] feature: deck/bag selector, refile cards
- [ ] feature: max characters
- [ ] refactor backend (deobfuscate logic)
- [ ] improve backup functionality
  - [ ] check if the backup + patch == current savefile
  - [ ] keep a patchlist instead of saving snapshots (when patches do not result in the savefile, create a new
    snapshot)
- [ ] UI upgrades
  - [ ] improve header (fancy savefile name formatting)
  - [ ] improve directory (visual improvement)
  - [ ] improve directory (savefile names and other metadata, needs backend extensions)
- [ ]  savefile selector optimisations

## 13.07.2024

## progress

- none (did not have time or motivation)

## tasks

(moved)

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
