package httpfetcher

import (
	"context"
	"crypto/sha3"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type UrlResource struct {
	Url      string
	Data     []byte
	Checksum []byte
}

// GetResourceAsync clarification: url comes from the savefile, which sloppily accepts links without a proper http/https scheme
// which is something that golang's http client does not tolerate. For the file to be properly substituted, I still need
// to write the 'malformed' url, but I need to prepend 'http://' for every incorrect url.
func GetResourceAsync(ctx context.Context, url string) (*UrlResource, error) {
	var requestUrl string
	if strings.HasPrefix(url, "http") {
		requestUrl = url
	} else {
		requestUrl = "http://" + url
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating request: %s", url)
	}
	var resp *http.Response
	var attempt = 0
	for attempt < 3 {
		resp, err = http.DefaultClient.Do(request)
		if err == nil {
			break
		}
		attempt += 1
	}

	if err != nil {
		return nil, errors.Wrapf(err, "error fetching resource %s", requestUrl)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non-ok status for: %s: %s", requestUrl, resp.Status)
	}

	respBody := resp.Body
	defer respBody.Close()

	result, err := InitUrlResult(url, respBody)
	if err != nil {
		return nil, errors.Wrapf(err, "could not hash: %s", requestUrl)
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
