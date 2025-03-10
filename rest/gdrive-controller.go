package rest

import (
	grdive2 "TTSBundler/grdive"
	"context"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"log"
)

func registerGDriveRoutes(app *fiber.App) {
	app.Get("/drive/auth/start", startAuth)
	app.Get("/drive/auth/callback", handleAuthCallback)
	app.Get("/drive/files", listFiles)
	app.Post("/drive/upload", testUpload)
}

// technically I don't need this one (open the web browser directly instead of hosting `/auth/start` endpoint)
func startAuth(ctx *fiber.Ctx) error {
	authURL := grdive2.AuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return ctx.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

func handleAuthCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	if code == "" {
		return fiber.ErrBadRequest
	}
	if grdive2.SetTokenFromCallback(code) != nil {
		return fiber.ErrInternalServerError
	}
	return nil
}

func listFiles(ctx *fiber.Ctx) error {
	var cnt context.Context = ctx.Context()
	res, err := grdive2.ListFiles(&cnt)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(res)
}

func testUpload(ctx *fiber.Ctx) error {
	var reqContext context.Context = ctx.Context()
	payload := ctx.Body() // Use body content as file data
	if len(payload) == 0 {
		log.Printf("Empty payload provided for file upload")
		return fiber.NewError(fiber.StatusBadRequest, "Supplied payload is empty")
	}
	res, err := grdive2.UploadFile(&reqContext, ctx.Query("filename", "file.txt"), ctx.Body())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(res)
}
