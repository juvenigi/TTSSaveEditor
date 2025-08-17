package internal

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"tts-cache-manager-cli/backend/internal/internal/properties"
	"tts-cache-manager-cli/backend/internal/internal/resource_map"
	"tts-cache-manager-cli/backend/internal/internal/tabletop_save"
	"tts-cache-manager-cli/backend/internal/internal/wrapped_io"

	"github.com/bwmarrin/snowflake"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const NewFileEvent = "ttsc:newFile"

type CacheManagerApi struct {
	singletonLock        sync.Mutex
	seed                 *snowflake.Node
	ctx                  context.Context
	propertiesController properties.Controller
	packData             resource_map.PackData
}

func (api *CacheManagerApi) Startup(ctx context.Context) {
	api.ctx = ctx
}

func NewCacheManagerApi() *CacheManagerApi {
	var api = &CacheManagerApi{}

	api.propertiesController = properties.InitPropertiesController()
	dir := api.propertiesController.GetApplicationProperties().PackDataDir
	packData, err := resource_map.InitPackDataFromFile(filepath.Join(dir, properties.PackDataYaml))
	if err != nil {
		panic(err)
	}
	api.packData = packData

	seed, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	api.seed = seed

	return api
}

func (api *CacheManagerApi) GetProperties() properties.ApplicationPropertiesView {
	appProperties := api.propertiesController.GetApplicationProperties()

	return appProperties.ToView()
}

func (api *CacheManagerApi) GetTabletopSaves() error {
	pp := api.propertiesController.GetApplicationProperties()
	err, resChan := tabletop_save.GetTabletopSaveFilesAsync(pp.GameDir, pp.PackDataDir)
	if err != nil {
		return err
	}
	for {
		select {
		case saveRes, open := <-resChan:
			if !open {
				return nil
			}
			if saveRes.Err != nil {
				return saveRes.Err
			}
			view := saveRes.Result.ToView()
			runtime.EventsEmit(api.ctx, NewFileEvent, view)

		case <-api.ctx.Done():
			return api.ctx.Err()
		}
	}
}

func (api *CacheManagerApi) WriteToPackData(saveLocation string, savename string) error {
	if ok := api.singletonLock.TryLock(); !ok {
		return errors.New("api busy")
	}
	defer api.singletonLock.Unlock()

	appProperties := api.propertiesController.GetApplicationProperties()
	gameDir := appProperties.GameDir
	packDataDir := appProperties.PackDataDir

	saveData, err := tabletop_save.GetTabletopSaveFile(saveLocation, packDataDir)
	if err != nil {
		return err
	}
	if err = api.packData.ImportFromSave(api.ctx, &saveData, gameDir); err != nil {
		return err
	}

	if len(savename) > 0 {
		modifiedJsonBlob, err := saveData.GetPortableJsonBlob(api.seed, savename)
		if err != nil {
			return err
		}

		if err = api.packData.WriteSaveToPackData(saveData.GetSaveName(), modifiedJsonBlob); err != nil {
			return err
		}
	}

	return nil
}

// ConstructPackedSave note : saveName != savefile name, it's what the game will show you in the save selection
func (api *CacheManagerApi) ConstructPackedSave(saveLocation string, saveName string) error {
	if ok := api.singletonLock.TryLock(); !ok {
		return errors.New("api busy")
	}
	defer api.singletonLock.Unlock()

	appProperties := api.propertiesController.GetApplicationProperties()

	saveFile, err := tabletop_save.GetTabletopSaveFile(saveLocation, appProperties.PackDataDir)
	if err != nil {
		return err
	}
	allResources := saveFile.GetAllResources()

	if err = api.packData.LocalizePackedResources(allResources); err != nil {
		return err
	}

	return saveFile.WriteNewSaveToSavesDir(filepath.Join(appProperties.GameDir), saveName)
}

func (api *CacheManagerApi) Delete(fileLocation string) error {
	return os.Remove(fileLocation)
}

func (api *CacheManagerApi) MergePackDataWithAnother(otherLoc string, makeBackup bool) error {
	if ok := api.singletonLock.TryLock(); !ok {
		return errors.New("api busy")
	}
	defer api.singletonLock.Unlock()

	var srcLoc string
	if strings.HasSuffix(otherLoc, ".7z") {
		srcLoc, err := os.MkdirTemp("", "unzipped")
		if err != nil {
			return err
		}

		if err = wrapped_io.ExtractArchive(otherLoc, srcLoc); err != nil {
			return err
		}
	} else {
		srcLoc = otherLoc
	}

	return api.packData.MergeWith(srcLoc, makeBackup)
}
