package app

import (
	"context"
	"path/filepath"
	"tts-cache-manager-cli/properties"
	"tts-cache-manager-cli/resource_map"
	"tts-cache-manager-cli/tabletop_save"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const NewFileEvent = "ttsc:newFile"

type CacheManagerApi struct {
	ctx                  context.Context
	propertiesController properties.Controller
	saveManager          tabletop_save.SaveManager
	packData             resource_map.PackData
}

func (api *CacheManagerApi) Startup(ctx context.Context) {
	api.ctx = ctx
}

func NewCacheManagerApi() *CacheManagerApi {
	var api = &CacheManagerApi{}
	api.propertiesController = properties.NewPropertiesController()
	api.saveManager = tabletop_save.NewSaveManager(api.propertiesController.GetApplicationProperties())
	dir := api.propertiesController.GetApplicationProperties().PackDataDir
	packData, err := resource_map.NewPackData(filepath.Join(dir))
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
	err, resChan := api.saveManager.GetTabletopSaveFilesAsync()
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

// todo: return a view of PackData instead
func (api *CacheManagerApi) WriteToPackData(saveLocation string) (error, bool) {
	gameDir := api.propertiesController.GetApplicationProperties().GameDir
	saveData, err := api.saveManager.GetTabletopSaveFile(saveLocation)
	if err != nil {
		return err, false
	}

	// todo: atm, this breaks the "clean code" principle of doing only one thing, as we are
	//  1. mutating pack data
	//  2. mutating save data
	//  on the other hand, doing so is hella convenient, I only need to make sure that I have a sync.Once error chan
	err = api.packData.AddResourcesFromSave(api.ctx, &saveData, gameDir)
	if err != nil {
		return err, false
	}
	return nil, true
}
