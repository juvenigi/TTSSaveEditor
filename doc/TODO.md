# Release version 1

- config via yaml located next to exe
    - same goes for pack data
- fill up pack data
- savefile management
    - view savefile entries (image preview, resources, amount of them cached)
    - create pack data
    - erase savefile

# Backlog / Shelf

- easier interaction with savedata (consolidation, deletion, etc...)
- fetch resources without needing to use tts
- handle 'root' resources (i.e. urls not associated with a GUID/Name/Nickname)

# TODOs
- cleanup console spam

- smart routing

- note to self: it does not make sense to parallelize ResourceBundle retrieval
    - but it does make sense to parallelize http fetching / mass checksum calculation

- consider using a single package for everything, stop being a Java developer :)
- [X] get rid of overengineered `.yaml` vs `.yml` (pick one)
- code arch cleanup (a lot of things are unnecessarily public, some code is best put into `package_name/internal`)
    - use `properties/internal` for _true private_ methods.
- cleanup directory paths that come from tts, because they are a mess

# Planned features

- [X] Consider whether there is any benefit in having two separate dirs for save and cache data
    - right now, I will simply use tabletop root dir for orientation, rationale being that TTS won't work properly if it
      does not exit

- [X] yaml config
- [x] scan savefile for images
    - [X] url escaping used for cached items
    - [X] scan cached objects
        - [X] get the (properly-escaped) resource name comparator
            - [x] actually, the extension of the cached is sometime written twice (filepng.png or file.png)
        - [X] ignore steam items (ask if user wishes to convert it to local data using cache)
            - obsolete because we handle things differently now (no more trying to retrieve http resources, only check
              the pack and tts cache)
- [/] Initialize the Pack
    - [X] core PackData functions
    - [X] Populating PackData
    - [ ] gracefully handle url resources where the filename is not obvious
- [/] PackData utils
    - [/] savefile templates
        - [X] listing present savefile templates
        - [ ] creating savefile template `.pack.json`
        - [ ] transforming the save template into a valid savefile
    - [ ] zip / unzip the pack
- [ ] Graceful panic handling (zenity)
- [ ] tabletop savesfile utils
    - [ ] patch by jsonpath
    - [ ] modify save data infos
    - [ ] 
- [ ] The grand contextification
- [ ] Wails API
    - channel-to-event publishing
    - list savefiles
    - get savefile report
        - populate PackData
        - create new save template
        - use save template
        - convert to cache
        - convert to local
    - simple container format
        - OS-native file select
        - export / pack
        - import / unpack
- [/] React frontend
    - [ ] TanStack Router
    - [ ] Service / Command Layer
    - [ ] Store / Reducer
    - [X] selectors / ViewModel
       - [/] File List 