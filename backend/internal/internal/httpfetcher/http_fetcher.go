package httpfetcher

import (
	"context"
	"crypto/sha3"
	"io"
	"net/http"
)

type UrlResource struct {
	Url      string
	Data     []byte
	Checksum []byte
}

func GetResourceAsync(ctx context.Context, url string) (*UrlResource, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	respBody := resp.Body
	defer respBody.Close()

	result, err := InitUrlResult(url, respBody)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func InitUrlResult(url string, respBody io.ReadCloser) (*UrlResource, error) {
	hash := sha3.New256()
	tee := io.TeeReader(respBody, hash)
	body, err := io.ReadAll(tee)
	if err != nil {
		return nil, err
	}

	return &UrlResource{url, body, hash.Sum(nil)}, nil
}
