package grdive

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/drive/v3"
	"log"
)

func ListFiles(ctx *fiber.Ctx) error {
	service, err := GetDriveService(ctx.Context())
	if err != nil {
		log.Printf("Error creating Drive client: %v", err)
		return err
	}

	r, err := service.Files.List().PageSize(10).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Printf("Error listing files: %v", err)
		return err
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

func UploadFile(ctx *fiber.Ctx) error {
	// Retrieve the Google Drive service
	service, err := GetDriveService(ctx.Context())
	if err != nil {
		log.Printf("Error creating Drive client: %v", err)
		return fiber.ErrInternalServerError
	}

	// Extract query parameters and request body
	filename := ctx.Query("filename", "file.txt")
	payload := ctx.Body() // Use body content as file data
	if len(payload) == 0 {
		log.Printf("Empty payload provided for file upload")
		return fiber.NewError(fiber.StatusBadRequest, "File content is empty")
	}

	// Create file metadata
	// Upload file to Google Drive
	file, err := service.Files.Create(&drive.File{Name: filename}).
		Media(bytes.NewReader(payload)).
		Do()
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return fiber.ErrInternalServerError
	}

	// Set permissions for the uploaded file
	_, err = service.Permissions.Create(file.Id, &drive.Permission{
		Type: "anyone", // Shareable by anyone with the link
		Role: "reader", // Read-only access
	}).Do()
	if err != nil {
		log.Printf("Error setting file permissions: %v", err)
		return fiber.ErrInternalServerError
	}

	// Prepare the response
	return ctx.JSON(fiber.Map{
		"message":  "File uploaded successfully",
		"fileId":   file.Id,
		"fileName": file.Name,
		"mimeType": file.MimeType,
	})
}
