package rest

import (
	"TTSBundler/domain"
	service2 "TTSBundler/service"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func noCache(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")
	return c.Next()
}

func registerRoutes(app *fiber.App) {
	app.Use("/", noCache)
	app.Static("/", "./res/browser")
	app.Get("/api/savefile", getSavefile)
	app.Patch("/api/savefile", patchSavefile)
	app.Get("/api/directory", getDirectory)
	app.Post("/api/card", addNewCard)
	app.Delete("/api/card", deleteCard)
	app.Get("/api/urls", getSaveResources)
}

func deleteCard(ctx *fiber.Ctx) error {
	savefileLocation := ctx.Query("path")
	deckJsonPath := ctx.Query("jsonPath")
	if savefileLocation == "" || deckJsonPath == "" {
		return fiber.ErrBadRequest
	}
	err, result := service2.DeleteCard(savefileLocation, deckJsonPath)
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
	var patchJson service2.PartialCard
	savefileLocation := ctx.Query("path")
	deckJsonPath := ctx.Query("jsonPath")
	if savefileLocation == "" || deckJsonPath == "" {
		return fiber.ErrBadRequest
	}
	err := json.Unmarshal(ctx.Body(), &patchJson)
	if err != nil {
		return fiber.ErrBadRequest
	}
	err, result := service2.AddNewCard(patchJson, savefileLocation, deckJsonPath)
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
	if jsonBytes, err := service2.PatchSavefile(path, body); err != nil {
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
	data, err := service2.GetEntries(path)
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
	if path == "" {
		return fiber.ErrBadRequest
	}
	log.Infof("requesting path: %s", path)
	jsonBytes, err := service2.GetSaveJson(path)
	if err != nil {
		return err
	}
	return ctx.
		Type("application/json", "UTF-8").
		Send(jsonBytes)
}

func getSaveResources(ctx *fiber.Ctx) error {
	path := ctx.Query("path")
	if path == "" {
		return fiber.ErrBadRequest
	}
	cachePath := ctx.Query("cachePath")

	log.Infof("requesting resource URLs: %s cache: %s", path, cachePath)
	jsonBytes, err := service2.GetSaveResources(path, cachePath)
	if err != nil {
		return err
	}
	return ctx.
		Type("application/json", "UTF-8").
		Send(jsonBytes)
}
