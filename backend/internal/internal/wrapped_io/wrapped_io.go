package wrapped_io

import (
	"bytes"
	"crypto/sha3"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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

func AttemptToWriteExtLess(dir string, name string, tryCount int, content []byte) (string, []string, error) {
	tries := 0
	var previousAttempts []string
	attemptName := name
	previousAttempts = append(previousAttempts, attemptName)
	for {
		file, err := os.OpenFile(filepath.Join(dir, attemptName), os.O_WRONLY|os.O_CREATE|os.O_EXCL, os.ModePerm)
		if err == nil {
			defer file.Close()
			_, err := io.Copy(file, bytes.NewReader(content))
			if err != nil {
				return "", previousAttempts, err
			}
			if err := file.Sync(); err != nil {
				return "", previousAttempts, err
			}
			break
		} else if !errors.Is(err, fs.ErrExist) {
			return "", previousAttempts, err
		}

		// try again after failed write attempt due to file already existing
		if tries++; tries > tryCount {
			break
		} else {
			previousAttempts = append(previousAttempts, attemptName)
			attemptName = fmt.Sprintf("%s-%d", attemptName, tryCount+tries)
		}
	}
	return attemptName, previousAttempts, nil
}
