package app

import (
	"context"
	"log"
	"path/filepath"
	"tts-cache-manager-cli/properties"
	"tts-cache-manager-cli/resource_map"
	"tts-cache-manager-cli/tabletop_save"
)

type CacheManagerApp struct {
	ctx                  context.Context
	propertiesController properties.Controller
	packData             resource_map.PackData
}

func (app *CacheManagerApp) startup(ctx context.Context) {
	app.ctx = ctx
}

func NewCacheManagerApp() CacheManagerApp {
	var instance CacheManagerApp
	instance.propertiesController = properties.NewPropertiesController()

	dir := instance.propertiesController.GetApplicationProperties().PackDataDir
	packData, err := resource_map.NewPackData(filepath.Join(dir, properties.PackDataYaml))
	if err != nil {
		panic(err)
	}
	instance.packData = packData

	return instance
}

// todo: WIP
func (app *CacheManagerApp) GetTabletopSaves() error {
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
			log.Println(save)
			log.Println("save file detected")
		case err := <-errCh:
			return err
		case <-app.ctx.Done():
			return app.ctx.Err()
		case <-doneCh:
			return nil
		}
	}
}
