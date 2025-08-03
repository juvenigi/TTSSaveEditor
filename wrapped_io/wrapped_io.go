package wrapped_io

import (
	"io"
	"os"
)

func CopyFile(src string, dest string) error {
	destF, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destF.Close()

	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	if _, err = io.Copy(destF, srcF); err != nil {
		return err
	}

	return srcF.Sync()
}
