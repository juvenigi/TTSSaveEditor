package grdive

import (
	"bytes"
	"context"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/drive/v3"
	"log"
)

func ListFiles(ctx *context.Context) (fiber.Map, error) {
	service, err := GetDriveService(*ctx)
	if err != nil {
		log.Printf("Error creating Drive client: %v", err)
		return nil, err
	}

	r, err := service.Files.List().PageSize(10).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Printf("Error listing files: %v", err)
		return nil, err
	}

	files := make([]fiber.Map, len(r.Files))
	for i, file := range r.Files {
		files[i] = fiber.Map{
			"name": file.Name,
			"id":   file.Id,
		}
	}

	return fiber.Map{
		"files": files,
	}, nil
}

func UploadFile(ctx *context.Context, filename string, payload []byte) (fiber.Map, error) {
	// Retrieve the Google Drive service
	service, err := GetDriveService(*ctx)
	if err != nil {
		log.Printf("Error creating Drive client: %v", err)
		return nil, err
	}
	if len(payload) == 0 {
		log.Printf("Empty payload provided for file upload")
		return nil, err
	}

	file, err := service.Files.Create(&drive.File{Name: filename}).
		Media(bytes.NewReader(payload)).
		Do()
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return nil, err
	}

	// Set permissions for the uploaded file
	_, err = service.Permissions.Create(file.Id, &drive.Permission{
		Type: "anyone", // Shareable by anyone with the link
		Role: "reader", // Read-only access
	}).Do()
	if err != nil {
		log.Printf("Error setting file permissions: %v", err)
		return nil, err
	}

	return fiber.Map{
		"message":  "File uploaded successfully",
		"fileId":   file.Id,
		"fileName": file.Name,
		"mimeType": file.MimeType,
	}, nil
}
