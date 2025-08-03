package app

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"path/filepath"
	"tts-cache-manager-cli/properties"
	"tts-cache-manager-cli/resource_map"
	"tts-cache-manager-cli/tabletop_save"
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

func (app *CacheManagerApi) GetTabletopSaves() error {
	sm := tabletop_save.NewSaveManager(app.propertiesController.GetApplicationProperties())

	filesCh := make(chan tabletop_save.TSSaveFile)
	defer close(filesCh)
	errCh := make(chan error)
	defer close(errCh)

	doneCh, err := sm.GetTabletopSaveFilesAsync(filesCh, errCh)
	if err != nil {
		return err
	}

	for {
		select {
		case save := <-filesCh:
			view := save.ToView()
			runtime.EventsEmit(app.ctx, NewFileEvent, view)

		case err := <-errCh:
			return err
		case <-app.ctx.Done():
			return app.ctx.Err()
		case <-doneCh:
			return nil
		}
	}
}
