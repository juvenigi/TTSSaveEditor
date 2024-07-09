# General
- simple gui for asking where the tts savefiles are located
- REST API for savefile push / pull

# Design principles

- broker is the single source of truth
    - it knows if TTS is running a 'stale' file
    - it notifies subscribers of new changes
    - the subscribers need to handle updates: discard changes, override, create copy

Diagram:
<editor: web, intellij> -- <golang broker> -- <TTS>

# Journal

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
- don't create too generalised reducers, split them up if necessary
- navigation to savefile (actions, effect, reducer, query)
- expand model to include the actual tts savefile
- components, quill.js?
