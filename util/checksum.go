package util

import (
	"crypto/sha3"
	"os"
)

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
