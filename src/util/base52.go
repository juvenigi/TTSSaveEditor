package util

import (
	"errors"
	"strings"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base = int64(len(alphabet))

// Decode converts a base52 string into an int64
func Decode(base52 string) (int64, error) {
	if base52 == "" {
		return 0, errors.New("empty string")
	}

	var value int64
	for _, r := range base52 {
		index := strings.IndexRune(alphabet, r)
		if index == -1 {
			return 0, errors.New("invalid character: " + string(r))
		}
		value = value*base + int64(index)
	}
	return value, nil
}

// Encode converts a non-negative int64 into a base52 string
func Encode(value int64) (string, error) {
	if value < 0 {
		return "", errors.New("negative values not supported")
	}
	if value == 0 {
		return string(alphabet[0]), nil
	}

	var result []byte
	for value > 0 {
		remainder := value % base
		result = append([]byte{alphabet[remainder]}, result...)
		value /= base
	}
	return string(result), nil
}
