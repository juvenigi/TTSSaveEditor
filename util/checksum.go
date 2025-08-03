package util

import (
	"crypto/sha3"
	"os"
)

func GetFileChecksum(fileLoc string) []byte {
	blob, err := os.ReadFile(fileLoc)
	if err != nil {
		return nil
	}
	new256 := sha3.New256()
	_, err = new256.Write(blob)
	if err != nil {
		return nil
	}
	return new256.Sum(nil)
}
