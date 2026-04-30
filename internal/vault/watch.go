package vault

import (
	"context"
	"time"
)

// WatchEvent represents a change detected while watching a secret path.
type WatchEvent struct {
	Path       string
	OldVersion int
	NewVersion int
	ChangedAt  time.Time
}

// WatchOptions configures the polling behaviour of a Watcher.
type WatchOptions struct {
	Mount    string
	Interval time.Duration
}

// DefaultWatchOptions returns sensible defaults for WatchOptions.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Mount:    "secret",
		Interval: 30 * time.Second,
	}
}

// Watcher polls a Vault secret path and emits events when the latest version changes.
type Watcher struct {
	client  *Client
	options WatchOptions
}

// NewWatcher creates a Watcher using the provided client and options.
func NewWatcher(client *Client, opts WatchOptions) *Watcher {
	return &Watcher{client: client, options: opts}
}

// Watch polls the given path at the configured interval, sending a WatchEvent to
// the returned channel whenever the latest version number increases. The channel
// is closed when ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context, path string) (<-chan WatchEvent, error) {
	history, err := GetVersionHistory(w.client, w.options.Mount, path)
	if err != nil {
		return nil, err
	}

	events := make(chan WatchEvent, 4)
	known := history.Latest()

	go func() {
		defer close(events)
		ticker := time.NewTicker(w.options.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				h, err := GetVersionHistory(w.client, w.options.Mount, path)
				if err != nil {
					continue
				}
				latest := h.Latest()
				if latest > known {
					events <- WatchEvent{
						Path:       path,
						OldVersion: known,
						NewVersion: latest,
						ChangedAt:  time.Now().UTC(),
					}
					known = latest
				}
			}
		}
	}()

	return events, nil
}
