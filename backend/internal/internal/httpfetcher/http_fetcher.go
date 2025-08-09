package httpfetcher

import (
	"context"
	"crypto/sha3"
	"io"
	"net/http"
	"tts-cache-manager-cli/backend/internal/internal/tabletop_save"
)

type UrlResource struct {
	Url      string
	Data     []byte
	Checksum []byte
}

// previous blocking code
func GetResource(url string) (*UrlResource, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	respBody := resp.Body
	defer resp.Body.Close()

	return readResult(url, respBody)
}

func readResult(url string, respBody io.ReadCloser) (*UrlResource, error) {
	hash := sha3.New256()
	tee := io.TeeReader(respBody, hash)
	body, err := io.ReadAll(tee)
	if err != nil {
		return nil, err
	}

	return &UrlResource{url, body, hash.Sum(nil)}, nil
}

// todo: write test for the async version
func GetResourceAsync(ctx context.Context, res *tabletop_save.GameResource) (*UrlResource, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, res.ResourceUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	respBody := resp.Body
	defer respBody.Close()

	result, err := readResult(res.ResourceUrl, respBody)
	if err != nil {
		return nil, err
	}

	return result, nil
}
