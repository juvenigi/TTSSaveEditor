package internal

import (
	"context"
	"path/filepath"
	"tts-cache-manager-cli/backend/internal/internal/properties"
	"tts-cache-manager-cli/backend/internal/internal/resource_map"
	"tts-cache-manager-cli/backend/internal/internal/tabletop_save"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const NewFileEvent = "ttsc:newFile"

type CacheManagerApi struct {
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

	return api
}

func (api *CacheManagerApi) GetProperties() properties.ApplicationPropertiesView {
	applicationProperties := api.propertiesController.GetApplicationProperties()

	return applicationProperties.ToView()
}

func (api *CacheManagerApi) GetTabletopSaves() error {
	applicationProperties := api.propertiesController.GetApplicationProperties()
	gameDir := applicationProperties.PackDataDir
	packDataDir := applicationProperties.PackDataDir

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
// todo: verify that the following mutations are done correctly:
// todo: return a view of PackData instead
func (api *CacheManagerApi) WriteToPackData(saveLocation string) error {
	applicationProperties := api.propertiesController.GetApplicationProperties()
	gameDir := applicationProperties.GameDir
	packDataDir := applicationProperties.PackDataDir

	saveData, err := tabletop_save.GetTabletopSaveFile(saveLocation, packDataDir)
	if err != nil {
		return err
	}

	if err = api.packData.ImportFromSave(api.ctx, &saveData, gameDir); err != nil {
		return err
	}

	return nil
}
