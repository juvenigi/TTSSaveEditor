package wrapped_io

import (
	"crypto/sha3"
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

	return destF.Sync()
}

func GetFileAndChecksum(fileLoc string) ([]byte, []byte) {
	blob, err := os.ReadFile(fileLoc)
	if err != nil {
		return nil, nil
	}
	new256 := sha3.New256()
	_, err = new256.Write(blob)
	if err != nil {
		return nil, nil
	}
	return blob, new256.Sum(nil)
}
