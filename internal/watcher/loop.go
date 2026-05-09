package watcher

import (
	"context"
	"time"

	"github.com/Guliveer/twitch-miner-go/internal/logger"
)

// pollLoop runs the standard polling lifecycle used by CategoryWatcher and
// TeamWatcher: fire evaluate once, then evaluate on each tick until the context
// is cancelled, at which point shutdown runs before returning. Logs are emitted
// via log.Info using the provided startMsg and stopMsg.
func pollLoop(
	ctx context.Context,
	log *logger.Logger,
	interval time.Duration,
	startMsg, stopMsg string,
	startArgs []any,
	evaluate func(context.Context),
	shutdown func(),
) error {
	log.Info(startMsg, startArgs...)

	evaluate(ctx)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info(stopMsg)
			shutdown()
			return ctx.Err()
		case <-ticker.C:
			evaluate(ctx)
		}
	}
}
