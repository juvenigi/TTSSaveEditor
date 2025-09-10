# KISS patch

- simplify app properties initialization (you are supposed to receive it from the frontend anyway)
    - make sure the fallback never fails
- simplify resource states
    - local : absolute path
    - localCache : cache hack
    - remoteCache : is a remote file but the cache can be found
    - remote : remote (http) without a cached entry
    - steam : uncached steam resource, cannot fetch due to no access
    - packed : remote file or checksum was found in packdata

# Release version 1

- config via yaml located next to exe
    - same goes for pack data
- fill up pack data
- savefile management
    - (might skip) view savefile entries (image preview, resources, amount of them cached)
    - create pack data
    - erase savefile

# BUG

- Configuration properties / Pack Data is not properly validated on api startup

# Backlog / Shelf

- make life harder with eslint (stricter rules)
- easier interaction with savedata (consolidation, deletion, etc...)
- fetch resources without needing to use tts
- handle 'root' resources (i.e. urls not associated with a GUID/Name/Nickname)

# The great bughunt

- ensure consistent init of pack data
- write down correct json pointer
- make all tests run
- eliminate all dead code at this point

# TODOs

- unclean impl.: currently using "Packed" as "Modified" - maybe I should name it appropriately or use a different
  status?

- jsonPointer copyFrom

- cleanup console spam

- smart routing ( this can be skipped entirely if we choose to restrict ourselves to overlays, collapsables and
  sidebars)

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
- [x] Initialize the Pack
    - [X] PackData init / core PackData functions
    - [x] gracefully handle url resources where the filename is not obvious
        - I learned that TTS is okay with files not having extensions, which is why I will write extensionless files in
          PackData
- [x] PackData utils
    - [x] savefile templates
        - [X] listing present savefile templates
        - [x] creating savefile template `.pack.json`
            - actually flush to disk
        - [x] transforming the save template into a valid savefile
    - [x] (skip) zip / unzip the pack
- [x] tabletop savefile utils
    - [x] json patch
    - [x] (skipped, as it is not necessary) modify save data infos
- [/] Populating PackData
    - [X] Rename save name on pack data entry/exit
    - [X] TS_Save pack json versioning
    - [ ] PackData merge (combine pack data with another packdata)
        - behavior: override, params: makeBackup
    - [ ] check all concurrency issues
    - [ ] do not be sloppy with pointers (make sure they are handled correctly)
- [X] File Delete (backend)
- [ ] The grand contextification (find places where it makes sense to cancel operations)
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
    - [X] TanStack Router
    - [ ] Make sure Singleton locks are tolerated
    - [ ] Provide some clue about progress when packing game savefile
    - [ ] Add creation / modification date to disambiguate saves 
    - [ ] Service / Command Layer
    - [ ] Store / Reducer
    - [X] selectors / ViewModel
        - [/] File List
- [x] Graceful panic handling (zenity)
    - actually, wails handles panics automatically
    - [ ] listen to error events and produce toasts
 