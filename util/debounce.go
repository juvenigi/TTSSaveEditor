package util

import (
	"context"
	"net/url"
	"sync"
	"time"
)

type Debouncer struct {
	lastExecution time.Time
	Pause         time.Duration
	lastHostCall  map[string]time.Time
	mu            sync.Mutex
}

func (debouncer *Debouncer) Init() {
	debouncer.lastHostCall = make(map[string]time.Time)
}

func (debouncer *Debouncer) Await(ctx context.Context) error {
	earliest := debouncer.lastExecution.Add(debouncer.Pause)
	return debouncer.await(ctx, earliest)
}

func (debouncer *Debouncer) AwaitHost(ctx context.Context, urllike string) error {
	parsed, err := url.Parse(urllike)
	if err != nil {
		// hotfix by slapping `http://` due to a sloppy upstream implementation
		parsed, err = url.Parse("http://" + urllike)
		if err != nil {
			// fallback option when url is weird
			return debouncer.Await(ctx)
		}
	}

	debouncer.mu.Lock()
	var lastCall time.Time
	lastCall, ok := debouncer.lastHostCall[parsed.Host]
	if !ok {
		debouncer.lastHostCall[parsed.Host] = time.Time{}
	}
	debouncer.mu.Unlock()

	return debouncer.await(ctx, lastCall)
}

func (debouncer *Debouncer) await(ctx context.Context, earliest time.Time) error {
	if !earliest.After(time.Now()) {
		debouncer.lastExecution = time.Now()
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(earliest.Sub(time.Now())):
		return nil
	}
}
