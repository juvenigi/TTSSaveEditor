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
	packData             resource_map.PackData
}

func (app *CacheManagerApi) Startup(ctx context.Context) {
	app.ctx = ctx
}

func NewCacheManagerApi() *CacheManagerApi {
	var instance = &CacheManagerApi{}
	instance.propertiesController = properties.NewPropertiesController()

	dir := instance.propertiesController.GetApplicationProperties().PackDataDir
	packData, err := resource_map.NewPackData(filepath.Join(dir))
	if err != nil {
		panic(err)
	}
	instance.packData = packData

	return instance
}

func (app *CacheManagerApi) GetProperties() properties.ApplicationPropertiesView {
	applicationProperties := app.propertiesController.GetApplicationProperties()

	return applicationProperties.ToView()
}

func (app *CacheManagerApi) GetTabletopSaves() error {
	sm := tabletop_save.NewSaveManager(app.propertiesController.GetApplicationProperties())

	err, resChan := sm.GetTabletopSaveFilesAsync()
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
			runtime.EventsEmit(app.ctx, NewFileEvent, view)

		case <-app.ctx.Done():
			return app.ctx.Err()
		}
	}
}
