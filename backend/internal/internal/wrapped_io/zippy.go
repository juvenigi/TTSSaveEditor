package wrapped_io

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/bodgit/sevenzip"
)

func ExtractArchive(archive string, destination string) error {
	if err := os.MkdirAll(destination, 0755); err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil
		} else {
			return err
		}
	}

	r, err := sevenzip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if err = extractFile(f, destination); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(file *sevenzip.File, destination string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	dest, err := os.Create(filepath.Join(destination, file.Name))
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err = io.Copy(dest, rc); err != nil {
		return err
	}

	return dest.Sync()
}
