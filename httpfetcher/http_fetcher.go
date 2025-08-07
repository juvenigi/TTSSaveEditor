package httpfetcher

import (
	"context"
	"crypto/sha3"
	"io"
	"net/http"
	"sync"
	"tts-cache-manager-cli/tabletop_save"
	"tts-cache-manager-cli/util"
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

type HttpFetcher struct {
	UrlMap map[string]*UrlResource
}

func NewHttpFetcher(ctx context.Context, resources []tabletop_save.GameResource) *util.Promise[HttpFetcher] {
	return util.RunAsync(ctx, newFetcherRunnable(resources))
}

func newFetcherRunnable(resources []tabletop_save.GameResource) func(ctx context.Context) (HttpFetcher, error) {
	return func(ctx context.Context) (HttpFetcher, error) {
		var mu = new(sync.Mutex)
		var wg sync.WaitGroup

		seen := make(map[string]struct{})
		UrlMap := make(map[string]*UrlResource)
		for _, res := range resources {
			if res.Status != tabletop_save.RemoteCached {
				continue
			}
			if _, already := seen[res.ResourceUrl]; already {
				continue
			} else {
				seen[res.ResourceUrl] = struct{}{}
			}
			wg.Add(1)
			go GetRemoteResourceToMap(ctx, &wg, mu, &res, UrlMap)
		}
		wg.Wait()
		return HttpFetcher{UrlMap: UrlMap}, nil
	}
}

func GetRemoteResourceToMap(ctx context.Context, wg *sync.WaitGroup, mu *sync.Mutex, res *tabletop_save.GameResource, UrlMap map[string]*UrlResource) {
	defer wg.Done()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, res.ResourceUrl, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	respBody := resp.Body
	defer respBody.Close()

	result, err := readResult(res.ResourceUrl, respBody)
	if err != nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	UrlMap[res.ResourceUrl] = result
}

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
