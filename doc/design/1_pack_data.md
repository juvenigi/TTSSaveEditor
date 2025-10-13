# 1: How to store local files in PackData?

`PackData` should provide performant mapping from tabletop savefile to the file located in PackData.

- Tabletop Simulator supports both remote resources (http\[s\]) and files on your local filesystem
- remote resource urls are unambiguous enough to be used directly
- local resources come with the complexity of being absolute paths of the current user's filesystem, reducing
  portability

I _could_ devise a way to convert local filesystem filepath into something discernible, but I find that introducing
checksums will increase overall robustness of all data involved.

## Design of PackData index (pack-data.yaml)

```go
package pack_data

type PackData struct {
	packFilenames        map[string]bool   // a quick set to lookup filenames
	sha3sumToResourceMap map[string]string // maps sha3 256 checksum to the filename contained in packdata (slow but reliable)
	urlToResourceMap     map[string]string // maps http urls to PackData
}
```

### Adding a new file to `PackData` means:

1. the checksum map is consulted to see if a file already exists
2. only if dealing with a remote resource:
    - urlToResourceMap is populated
3. packFilenames are extended
4. sha3sumToResourceMap and urlToResourceMap are persisted to yaml

### How game resource lookup using `PackData` is conducted:

- for 'packed'/'failed' resource (i.e. a resource pointing to a file in PackData): skip
- for 'local' resource: hash value is made for the file declared in the savefile
- for 'remote'/'remoteCached' resource: TTS cache is extracted (to eliminate unavailable resources): then
  urlToResourceMap is consulted
    - in parallel, the actual resource will be retrieved via httpClient, and a checksum comparison is done for both
      files

#### Implications for frontend/backend

- PopulatePackData yields intermediate results, which can be handled via channels
    - the channels can propagate events
    - I prefer channels over using wails event publisher directly, because I don't want to introduce a wails dependency
      in my 'inner shell', I want external code dependencies to live in the 'outer shell'.
- http fetching should be obviously done in parallel using goroutines
