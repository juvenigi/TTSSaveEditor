package tabletop_save

import (
	"bytes"
	"context"
	"crypto/sha3"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
	"tts-cache-manager-cli/backend/internal/internal/httpfetcher"
	"tts-cache-manager-cli/backend/internal/internal/resource_map"
	"tts-cache-manager-cli/backend/internal/internal/wrapped_io"
	"tts-cache-manager-cli/util"
)

func AddLocalFile(rm *resource_map.PackData, res *GameResource) (error, *ApplicationFile) {
	var location = strings.TrimPrefix(res.ResourceUrl, "file:///")
	blob, sum := wrapped_io.GetFileAndChecksum(location)
	if sum == nil {
		return errors.New("file not found"), nil
	}

	packDataFile, exists := rm.Sha3Archive[hex.EncodeToString(sum)]
	if exists {
		res.ResourceUrl = "pack:///" + packDataFile
		res.Status = Packed
		return nil, nil
	}

	var destFilename = GetCacheFilename(filepath.Base(res.ResourceUrl))
	return nil, &ApplicationFile{
		blob:     blob,
		filepath: destFilename,
		res:      res,
	}
}

func AddRemoteFile(rm *resource_map.PackData, ctx context.Context, res *GameResource) (error, *ApplicationFile) {
	urlRes, err := httpfetcher.GetResourceAsync(ctx, res.ResourceUrl)
	if err != nil {
		res.Status = Failed
		return err, nil
	}

	packFileName, ok := rm.Sha3Archive[hex.EncodeToString(urlRes.Checksum)]
	if ok {
		if file, err := os.ReadFile(filepath.Join(rm.PackDataDir, packFileName)); err == nil {
			res.Status = Failed // todo: re-check pack data consistency if such error is detected
			return err, nil

		} else if bytes.Compare(urlRes.Data, file) != 0 {
			res.Status = Failed
			return errors.New("hash collision / file inconsistency detected (you are a unicorn)"), nil
		}

		res.ResourceUrl = "pack://" + packFileName
		res.Status = Packed
		return nil, nil
	}
	// first we check the entries ahead of time, then we rely on the singleton nature of the syscalls
	escapedName := GetCacheFilename(urlRes.Url)
	return nil, &ApplicationFile{urlRes.Data, escapedName, res}
}

type ApplicationFile struct {
	blob     []byte
	filepath string
	res      *GameResource
}

func FlushResources(rm *resource_map.PackData, res *GameResource, file []byte, filenames *resource_map.DirectorySnapshot) error {
	escapedName := GetCacheFilename(res.ResourceUrl)
	written, failedAttempts, criticalErr := wrapped_io.AttemptToWriteExtLess(rm.PackDataDir, escapedName, 10, file)
	if criticalErr != nil {
		res.Status = Failed
		return criticalErr
	}

	filenames.Entries[written] = struct{}{}
	for _, fAttempt := range failedAttempts {
		filenames.Entries[fAttempt] = struct{}{}
	}

	new256 := sha3.New256()
	if _, criticalErr = new256.Write(file); criticalErr != nil {
		return criticalErr
	}

	rm.UrlArchive[res.ResourceUrl] = written
	rm.Sha3Archive[hex.EncodeToString(new256.Sum(nil))] = written

	res.UpdateToPacked(written)

	return nil
}

func AddRemoteCachedFile(rm *resource_map.PackData, res *GameResource, scanner *GameCacheFinder) (error, *ApplicationFile) {
	filePath, ok := scanner.ResourceCached[GetCacheFilename(res.ResourceUrl)]
	if !ok {
		return errors.New("cached file not found"), nil
	}

	file, fileChecksum := wrapped_io.GetFileAndChecksum(filePath)
	if fileChecksum == nil {
		return errors.New("cached file not found"), nil
	}

	packDatafile, ok := rm.Sha3Archive[hex.EncodeToString(fileChecksum)]
	if ok {
		res.ResourceUrl = "pack://" + packDatafile
		res.Status = Packed
		return nil, nil
	} else {
		destinationFilename := GetCacheFilename(res.ResourceUrl)
		return nil, &ApplicationFile{
			blob:     file,
			filepath: destinationFilename,
			res:      res,
		}
	}
}

func CachePackedResources(rm *resource_map.PackData, resList []*GameResource, gameDir string) error {
	for _, res := range resList {
		if res.Status == Packed {
			if err := copyPackedToGameCache(res, gameDir, rm.PackDataDir); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyPackedToGameCache(res *GameResource, gameDir string, packDataDir string) error {
	if res.Status != Packed {
		return fmt.Errorf("unexpected resource type %d", res.Status)
	}

	var lastSlash = strings.LastIndex(res.JsonPointer, "/")
	lastJsonPointerSegment := res.JsonPointer[lastSlash+1:]

	var filename = strings.TrimPrefix(res.ResourceUrl, "pack://")

	var gameCacheFolder string
	switch lastJsonPointerSegment {
	case "AssetBundleURL":
		gameCacheFolder = "Assetbundles"

	case "ImageURL", "ImageSecondaryURL", "DiffuseURL", "NormalURL", "SpecularURL", "SkyURL", "FaceURL", "BackURL":
		gameCacheFolder = "Images"

	case "MeshURL", "ColliderURL":
		gameCacheFolder = "Models"

	case "AudioURL":
		gameCacheFolder = "Audio"

	case "XmlUIURL":
		gameCacheFolder = "UI"

	default:
		return fmt.Errorf("unknown json pointer: %s", lastJsonPointerSegment)
	}

	srcPath := filepath.Join(packDataDir, filename)

	finalPath := filepath.Join(gameDir, "Mods", gameCacheFolder, filename)

	if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	if err := wrapped_io.CopyFile(srcPath, finalPath); err != nil {
		return err
	}

	res.ResourceUrl = filename

	return nil
}

func AddGameResources(rm *resource_map.PackData, ctx context.Context, data *TSSaveFile, gameDir string) ([]string, error) {
	var sem *util.Semaphore = util.NewSemaphore(5)
	var err error
	var minorErr []string
	var seen = make(map[string]struct{})
	var fileChan = make(chan *ApplicationFile)
	var wg sync.WaitGroup
	var cache GameCacheFinder

	filenames, err := resource_map.NewDirectorySnapshot(rm.PackDataDir)
	if err != nil {
		return nil, err
	}

	if err = cache.InitGameCacheFinder(gameDir); err != nil {
		return nil, err
	}
	resources := data.GetAllResources()
	for _, resource := range resources {
		cache.Update(resource)
	}

	for idx := range resources {
		var resourceRef = resources[idx]
		if _, already := seen[resourceRef.ResourceUrl]; already {
			continue
		}
		seen[resourceRef.ResourceUrl] = struct{}{}

		wg.Add(1)
		go func(idx int, fileChan chan *ApplicationFile) {
			defer wg.Done()
			sem.Acquire()
			defer sem.Release()

			var errSw error
			var file *ApplicationFile
			switch resourceRef.Status {
			case Remote, Steam, RemoteCached:
				if packUrl, ok := rm.UrlArchive[resourceRef.ResourceUrl]; ok {
					resourceRef.UpdateToPacked(packUrl)
				} else if resourceRef.Status == Remote || resourceRef.Status == Steam {
					timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
					defer cancel()
					errSw, file = AddRemoteFile(rm, timeoutCtx, resourceRef)
				} else {
					errSw, file = AddRemoteCachedFile(rm, resourceRef, &cache)
				}
			case Local:
				errSw, file = AddLocalFile(rm, resourceRef)
			default:
				errSw = fmt.Errorf("invalid resource status: %d url: %s", resourceRef.Status, resourceRef.ResourceUrl)
			}
			if errSw != nil {
				log.Println(errSw)
				minorErr = append(minorErr, errSw.Error())
				return
			}
			if file != nil {
				fileChan <- file
			}
		}(idx, fileChan)
	}
	promise := make(chan struct{})
	go func() {
	loop:
		for {
			select {
			case file := <-fileChan:
				file.filepath = filenames.GetUniqueFilename(file.filepath)
				if err := FlushResources(rm, file.res, file.blob, filenames); err != nil {
					minorErr = append(minorErr, err.Error())
				}
			case <-ctx.Done():
				break loop
			case <-promise:
				break loop
			}
		}
		promise <- struct{}{}
	}()
	wg.Wait()
	promise <- struct{}{}
	<-promise
	if err := rm.FlushToDisk(); err != nil {
		return minorErr, err
	}

	return minorErr, nil
}

const packLen = len("pack://")

func LocalizePackedResources(rm *resource_map.PackData, resources []*GameResource) error {
	packJsonLocs, err := rm.GetAllFilesInPackData()
	if err != nil {
		return err
	}

	for idx := range resources {
		var res = resources[idx]
		if res.Status == Packed && strings.HasPrefix(res.ResourceUrl, "pack://") {
			filename := res.ResourceUrl[packLen:]
			if slices.Contains(packJsonLocs, filename) {
				res.ResourceUrl = "file:///" + filepath.Join(rm.PackDataDir, filename)
			} else {
				return fmt.Errorf("could not find resource %s in PackData", res.ResourceUrl)
			}
		} else if res.Status == Remote {
			if packDataLoc, ok := rm.UrlArchive[res.ResourceUrl]; ok {
				res.Status = Packed
				res.ResourceUrl = "file:///" + filepath.Join(rm.PackDataDir, packDataLoc)
			} else {
				continue
			}
		}
	}

	return nil
}
