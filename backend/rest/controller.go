package rest

import (
	"TTSBundler/backend/domain"
	"TTSBundler/backend/service"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func registerRoutes(app *fiber.App) {
	app.Static("/", "./res/angular/browser")
	app.Get("/api/savefile", getSavefile)
	app.Patch("/api/savefile", patchSavefile)
	app.Get("/api/directory", getDirectory)
	app.Post("/api/card", addNewCard)
	app.Delete("/api/card", deleteCard)
}

func deleteCard(ctx *fiber.Ctx) error {
	savefileLocation := ctx.Query("path")
	deckJsonPath := ctx.Query("jsonPath")
	if savefileLocation == "" || deckJsonPath == "" {
		return fiber.ErrBadRequest
	}
	err, result := service.DeleteCard(savefileLocation, deckJsonPath)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	if err := ctx.
		Type("application/json", "UTF-8").
		Send(result); err != nil {
		return fiber.ErrInternalServerError
	}
	return nil
}

func addNewCard(ctx *fiber.Ctx) error {
	var patchJson domain.PartialCard
	savefileLocation := ctx.Query("path")
	deckJsonPath := ctx.Query("jsonPath")
	if savefileLocation == "" || deckJsonPath == "" {
		return fiber.ErrBadRequest
	}
	err := json.Unmarshal(ctx.Body(), &patchJson)
	if err != nil {
		return fiber.ErrBadRequest
	}
	err, result := service.AddNewCard(patchJson, savefileLocation, deckJsonPath)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	if err := ctx.
		Type("application/json", "UTF-8").
		Send(result); err != nil {
		return fiber.ErrInternalServerError
	}
	return nil
}

func patchSavefile(ctx *fiber.Ctx) error {
	path := ctx.Query("path")
	body := ctx.Body()
	log.Infof("requesting path: %s", path)
	if path == "" || body == nil {
		return fiber.ErrBadRequest
	}
	if jsonBytes, err := service.PatchSavefile(path, body); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		if err := ctx.
			Type("application/json", "UTF-8").
			Send(jsonBytes); err != nil {
			return fiber.ErrInternalServerError
		}
	}

	return nil
}

func getDirectory(ctx *fiber.Ctx) error {
	path := ctx.Query("path")
	log.Infof("requesting path: %s", path)
	if len(path) == 0 {
		if def, err := domain.GetDefaultTTSAbsPath(); err == nil {
			path = def
		} else {
			return fiber.ErrBadRequest
		}
	}
	data, err := service.GetEntries(path)
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := ctx.JSON(data); err != nil {
		return err
	}
	return nil
}

func getSavefile(ctx *fiber.Ctx) error {
	path := ctx.Query("path")
	log.Infof("requesting path: %s", path)
	if path == "" {
		return fiber.ErrBadRequest
	}
	if jsonBytes, err := service.GetSaveJson(path); err != nil {
		return fiber.ErrInternalServerError
	} else {
		if err := ctx.
			Type("application/json", "UTF-8").
			Send(jsonBytes); err != nil {
			return fiber.ErrInternalServerError
		}
	}
	return nil
}
