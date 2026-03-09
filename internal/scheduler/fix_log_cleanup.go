package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/torsod/go-webserver/internal/store"
)

// StartFIXLogCleanup runs periodic FIX log rotation
func StartFIXLogCleanup(ctx context.Context, fixLogs store.FIXLogStore) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	slog.Info("started FIX log cleanup scheduler (1h interval, 7-day retention)")

	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping FIX log cleanup scheduler")
			return
		case <-ticker.C:
			deleted, err := fixLogs.DeleteOlderThan(ctx, 7)
			if err != nil {
				slog.Error("FIX log cleanup failed", "error", err)
			} else if deleted > 0 {
				slog.Info("cleaned up old FIX logs", "deleted", deleted)
			}
		}
	}
}
