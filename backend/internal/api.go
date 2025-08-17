package internal

import (
	"context"
	"errors"
	"log"
	"path/filepath"
	"sync"
	"tts-cache-manager-cli/backend/internal/internal/properties"
	"tts-cache-manager-cli/backend/internal/internal/resource_map"
	"tts-cache-manager-cli/backend/internal/internal/tabletop_save"

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
	appProperties := api.propertiesController.GetApplicationProperties()
	gameDir := appProperties.PackDataDir
	packDataDir := appProperties.PackDataDir

	err, resChan := tabletop_save.GetTabletopSaveFilesAsync(gameDir, packDataDir)
	if err != nil {
		return err
	}

	for {
		select {
		case saveRes, ok := <-resChan:
			if !ok {
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

// WriteToPackData
//  1. mutating pack data
//  2. mutating save data
//
// todo: thread-safety
// todo: verify that the following mutations are done correctly:
// todo: return a view of PackData instead
func (api *CacheManagerApi) WriteToPackData(saveLocation string, writeSaveJson bool) error {
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

	if writeSaveJson {
		modifiedJsonBlob, err := saveData.GetPortableJsonBlob(api.seed)
		if err != nil {
			return err
		}

		if err = api.packData.WriteSaveToPackData(saveData.GetSaveName(), modifiedJsonBlob); err != nil {
			return err
		}
	}

	return nil
}

// todo: make this thread-safe
// todo: customize save name
func (api *CacheManagerApi) ConstructPackedSave(saveLocation string) error {
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

	// tmp debug
	packedCnt := 0
	for _, resource := range allResources {
		if resource.Status == tabletop_save.Packed {
			packedCnt++
		}
	}
	log.Println(packedCnt)

	if err = api.packData.LocalizePackedResources(allResources); err != nil {
		return err
	}

	return saveFile.WriteNewSaveToSavesDir(filepath.Join(appProperties.GameDir))
}
