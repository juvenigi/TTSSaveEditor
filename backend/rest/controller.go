package rest

import (
	"TTSBundler/backend/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func registerRoutes(app *fiber.App) {
	app.Static("/", "./res/angular/browser")
	app.Get("/api/savefile", handleGetSavefile)
	app.Get("/api/directory", handleGetDirectory)
}

func handleGetDirectory(ctx *fiber.Ctx) error {
	path := ctx.Query("path")
	log.Infof("requesting path: %s", path)
	data, err := service.GetEntries(path)
	if err != nil {
		if err := ctx.
			Status(fiber.StatusBadRequest).
			SendString("Bad Request"); err != nil {
			return err
		}
		return nil
	}
	if err := ctx.
		Type("application/json", "UTF-8").
		Send(data); err != nil {
		return err
	}
	return nil
}

func handleGetSavefile(ctx *fiber.Ctx) error {
	path := ctx.Query("path")
	log.Infof("requesting path: %s", path)
	if path == "" {
		if err := ctx.Status(fiber.StatusBadRequest).SendString("Bad Request"); err != nil {
			return err
		}
		return nil
	}
	if jsonBytes, err := service.GetSaveJson(path); err != nil {
		if err := ctx.Status(fiber.StatusBadRequest).SendString("Bad Request"); err != nil {
			return err
		}
		return nil
	} else {
		if err := ctx.
			Type("application/json", "UTF-8").
			Send(jsonBytes); err != nil {
			return err
		}
	}
	return nil
}