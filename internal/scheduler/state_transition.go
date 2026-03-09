package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/torsod/go-webserver/internal/service"
)

// StartStateTransitionScheduler runs the offering state transition checker
func StartStateTransitionScheduler(ctx context.Context, stateService *service.OfferingStateService) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	slog.Info("started offering state transition scheduler (30s interval)")

	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping state transition scheduler")
			return
		case <-ticker.C:
			stateService.CheckAndTransitionOfferings(ctx)
		}
	}
}
