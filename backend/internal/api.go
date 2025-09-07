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
	singletonLock         sync.Mutex
	seed                  *snowflake.Node
	ctx                   context.Context
	applicationProperties properties.ApplicationPropertiesView
	packData              resource_map.PackData
}

func NewCacheManagerApi() *CacheManagerApi {
	var api = &CacheManagerApi{}
	applicationProperties, err := properties.SetupDefaultAppProperties()
	if err != nil {
		panic(err)
	}

	err = applicationProperties.Validate()
	if err != nil {
		panic(err)
	}
	api.applicationProperties = *applicationProperties.ToView()
	packData, err := resource_map.InitPackDataFromFile(filepath.Join(api.applicationProperties.PackDataDir, properties.PackDataJson))
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

func (api *CacheManagerApi) Startup(ctx context.Context) {
	api.ctx = ctx
}

func (api *CacheManagerApi) SetProperties(view properties.ApplicationPropertiesView) error {
	defer noPanicBeHappy(api.ctx)

	var props properties.ApplicationProperties

	if err := props.InitFrom(&view); err != nil {
		return err
	} else {
		api.applicationProperties = view

		if packData, err := resource_map.InitPackDataFromFile(filepath.Join(view.PackDataDir, properties.PackDataJson)); err != nil {
			return err
		} else {
			api.packData = packData
			return nil
		}
	}
}

func (api *CacheManagerApi) GetProperties() properties.ApplicationPropertiesView {
	defer noPanicBeHappy(api.ctx)

	return api.applicationProperties
}

func (api *CacheManagerApi) GetTabletopSaves(scanDir string) error {
	defer noPanicBeHappy(api.ctx)

	pp := api.applicationProperties
	err, resChan := tabletop_save.GetTabletopSaveFilesAsync(scanDir, pp.PackDataDir)
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

func (api *CacheManagerApi) WriteToPackData(saveLocation string, saveName string) error {
	defer noPanicBeHappy(api.ctx)

	if ok := api.singletonLock.TryLock(); !ok {
		return errors.New("api busy")
	}
	defer api.singletonLock.Unlock()

	appProperties := api.applicationProperties
	gameDir := appProperties.GameDir
	packDataDir := appProperties.PackDataDir

	saveData, err := tabletop_save.GetTabletopSaveFile(saveLocation, packDataDir)
	if err != nil {
		return err
	}
	if err = api.packData.ImportFromSave(api.ctx, saveData, gameDir); err != nil {
		return err
	}

	if len(saveName) > 0 {
		modifiedJsonBlob, err := saveData.IntoPortable(api.seed, saveName)
		if err != nil {
			return err
		}

		if err = api.packData.WriteSaveToPackData(saveData.GetSaveName(), modifiedJsonBlob); err != nil {
			return err
		}
	}

	return nil
}

// CreateSingleplayerSave note : saveName != savefile name, save name is what the game will show you in the save selection
func (api *CacheManagerApi) CreateSingleplayerSave(saveLocation string, saveName string) error {
	defer noPanicBeHappy(api.ctx)

	if ok := api.singletonLock.TryLock(); !ok {
		return errors.New("api busy")
	}
	defer api.singletonLock.Unlock()

	appProperties := api.applicationProperties

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

func (api *CacheManagerApi) ConstructPackCachedSave(saveLocation string, saveName string) error {
	defer noPanicBeHappy(api.ctx)

	if ok := api.singletonLock.TryLock(); !ok {
		return errors.New("api busy")
	}
	defer api.singletonLock.Unlock()

	appProperties := api.applicationProperties

	saveFile, err := tabletop_save.GetTabletopSaveFile(saveLocation, appProperties.PackDataDir)
	if err != nil {
		return err
	}
	allResources := saveFile.GetAllResources()

	if err := api.packData.CachePackedResources(allResources, appProperties.GameDir); err != nil {
		return err
	}

	return saveFile.WriteNewSaveToSavesDir(filepath.Join(appProperties.GameDir), saveName)
}

func (api *CacheManagerApi) Delete(fileLocation string) error {
	defer noPanicBeHappy(api.ctx)

	return os.Remove(fileLocation)
}

func (api *CacheManagerApi) MergePackDataWithAnother(otherLoc string, makeBackup bool) error {
	defer noPanicBeHappy(api.ctx)

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

func (api *CacheManagerApi) PanicButton() error {
	defer noPanicBeHappy(api.ctx)

	panic("test panic")
	return nil
}

func noPanicBeHappy(ctx context.Context) {
	if err := recover(); err != nil {
		runtime.EventsEmit(ctx, "runtime:error", err)
	}
}
