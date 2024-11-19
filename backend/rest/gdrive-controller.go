package rest

import (
	"TTSBundler/backend/grdive"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

func registerGDriveRoutes(app *fiber.App) {
	app.Get("/auth/start", startAuth)
	app.Get("/drive/start", startDrive)
	app.Get("/callback", handleCallback)
	app.Get("/drive/files", listFiles)
	app.Post("/drive/upload", testUpload)
}

// technically I don't need this one (open the web browser directly instead of hosting `/auth/start` endpoint)
func startAuth(ctx *fiber.Ctx) error {
	authURL := grdive.AuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return ctx.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

func startDrive(ctx *fiber.Ctx) error {

	var err = browser.OpenURL(grdive.AuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline))
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return nil
}

func handleCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	if code == "" {
		return fiber.ErrBadRequest
	}
	if grdive.SetTokenFromCallback(code) != nil {
		return fiber.ErrInternalServerError
	}
	return nil
}

func listFiles(ctx *fiber.Ctx) error {
	return grdive.ListFiles(ctx)
}

func testUpload(ctx *fiber.Ctx) error {
	return grdive.UploadFile(ctx)
}
