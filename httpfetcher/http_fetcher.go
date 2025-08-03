package httpfetcher

import (
	"crypto/sha3"
	"io"
	"net/http"
)

type UrlResource struct {
	Url      string
	Data     []byte
	checksum []byte
}

func GetResource(url string) (*UrlResource, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	respBody := resp.Body
	defer resp.Body.Close()

	hash := sha3.New256()
	tee := io.TeeReader(respBody, hash)
	body, err := io.ReadAll(tee)
	if err != nil {
		return nil, err
	}

	return &UrlResource{url, body, hash.Sum(nil)}, nil
}
