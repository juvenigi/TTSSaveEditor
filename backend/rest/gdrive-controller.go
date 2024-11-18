package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"log"
	"os"
)

var clientId string
var clientSecret string

var config oauth2.Config
var token oauth2.Token

func init() {
	clientSecret = os.Getenv("TTS_GDRIVE_SECRET")
	clientId = os.Getenv("TTS_GDIRVE_CLIENT")
	if clientSecret == "" || clientId == "" {
		panic("envvars not set!")
	}
	// Set up OAuth2 config
	config = oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:3000/callback",
		Scopes:       []string{drive.DriveScope},
		Endpoint:     google.Endpoint,
	}
}

func registerGDriveRoutes(app *fiber.App) {
	app.Get("/auth/start", startAuth)
	app.Get("/callback", handleCallback)
	app.Get("/drive/files", listFiles)
	app.Post("/drive/upload", testUpload)
}

// technically I don't need this one (open the web browser directly instead of hosting `/auth/start` endpoint)
func startAuth(ctx *fiber.Ctx) error {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return ctx.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

func handleCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	if code == "" {
		return fiber.ErrBadRequest
	}

	// Exchange authorization code for tokens
	myToken, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging code: %v", err)
		return fiber.ErrInternalServerError
	}
	token = *myToken // not sure if the variable gets shadowed, or is it simply a lifetime issue, skill issue detected

	// Store token in-memory (or securely in a real app)
	tokenJSON, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		log.Printf("Error marshaling token: %v", err)
		return fiber.ErrInternalServerError
	}

	return ctx.
		Status(fiber.StatusOK).
		JSON(fiber.Map{
			"message": "Authorization successful!",
			"token":   string(tokenJSON),
		})
}

func listFiles(ctx *fiber.Ctx) error {
	driveService, err := drive.NewService(context.Background(), option.WithCredentials(&google.Credentials{
		ProjectID:              "",
		TokenSource:            oauth2.StaticTokenSource(&token),
		JSON:                   nil, // hmm
		UniverseDomainProvider: nil,
	}))
	if err != nil {
		log.Printf("Error creating Drive client: %v", err)
		return fiber.ErrInternalServerError
	}

	r, err := driveService.Files.List().PageSize(10).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Printf("Error listing files: %v", err)
		return fiber.ErrInternalServerError
	}

	files := make([]fiber.Map, len(r.Files))
	for i, file := range r.Files {
		files[i] = fiber.Map{
			"name": file.Name,
			"id":   file.Id,
		}
	}

	return ctx.JSON(fiber.Map{
		"files": files,
	})
}

func testUpload(ctx *fiber.Ctx) error {
	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(oauth2.StaticTokenSource(&token)))
	if err != nil {
		log.Printf("Error creating Drive client: %v", err)
		return fiber.ErrInternalServerError
	}

	filename := ctx.Query("filename", "file.txt")
	payload := ctx.Body() // Use body content as file data

	fileMetadata := &drive.File{
		Name: filename,
	}

	file, err := driveService.Files.Create(fileMetadata).
		Media(bytes.NewReader(payload)).
		Do()
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return fiber.ErrInternalServerError
	}

	// Create a permission to make the file shareable
	_, err = driveService.Permissions.Create(file.Id, &drive.Permission{
		Type: "anyone", // Shareable by anyone with the link
		Role: "reader", // Read-only access
	}).Do()
	if err != nil {
		log.Printf("Error setting file permissions: %v", err)
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(fiber.Map{
		"message":     "File uploaded successfully",
		"fileId":      file.Id,
		"fileName":    file.Name,
		"mimeType":    file.MimeType,
		"webViewLink": file.WebViewLink,
	})
}
