package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/store"
)

// Valid state transitions per CB-PDP-Gen-2 spec Section 4.4
var validTransitions = map[domain.OfferingState][]domain.OfferingState{
	domain.OfferingStateNew:          {domain.OfferingStateUpcoming, domain.OfferingStateCanceled},
	domain.OfferingStateUpcoming:     {domain.OfferingStateOpen, domain.OfferingStateCanceled},
	domain.OfferingStateOpen:         {domain.OfferingStateHalted, domain.OfferingStateFrozen, domain.OfferingStateClosePending},
	domain.OfferingStateHalted:       {domain.OfferingStateOpen, domain.OfferingStateFrozen, domain.OfferingStateClosePending, domain.OfferingStateCanceled},
	domain.OfferingStateFrozen:       {domain.OfferingStateOpen, domain.OfferingStateHalted, domain.OfferingStateClosePending, domain.OfferingStateCanceled},
	domain.OfferingStateClosePending: {domain.OfferingStateClosing},
	domain.OfferingStateClosing:      {domain.OfferingStateClearing},
	domain.OfferingStateClearing:     {domain.OfferingStateClosed},
	domain.OfferingStateClosed:       {},
	domain.OfferingStateCanceled:     {},
}

// OfferingStateService manages offering state transitions
type OfferingStateService struct {
	offerings store.OfferingStore
	orders    store.OrderStore
}

func NewOfferingStateService(offerings store.OfferingStore, orders store.OrderStore) *OfferingStateService {
	return &OfferingStateService{offerings: offerings, orders: orders}
}

// GetValidTransitions returns valid target states from the current state
func GetValidTransitions(state domain.OfferingState) []domain.OfferingState {
	return validTransitions[state]
}

// IsOrderEntryAllowed checks if new orders can be entered
func IsOrderEntryAllowed(state domain.OfferingState) bool {
	return state == domain.OfferingStateOpen
}

// IsOrderModificationAllowed checks if order modifications are allowed
func IsOrderModificationAllowed(state domain.OfferingState, isPriceUp bool) bool {
	if state == domain.OfferingStateOpen {
		return true
	}
	// In HALTED state, only price-up modifications are allowed
	if state == domain.OfferingStateHalted && isPriceUp {
		return true
	}
	return false
}

// AreCancellationsAllowed checks if cancellations are allowed
func AreCancellationsAllowed(state domain.OfferingState, isMO bool) bool {
	switch state {
	case domain.OfferingStateOpen, domain.OfferingStateHalted:
		return true
	case domain.OfferingStateFrozen, domain.OfferingStateClosePending:
		return isMO // Only MO can cancel in these states
	}
	return false
}

// IsReadOnlyState checks if the state is read-only
func IsReadOnlyState(state domain.OfferingState) bool {
	switch state {
	case domain.OfferingStateNew, domain.OfferingStateUpcoming,
		domain.OfferingStateClosePending, domain.OfferingStateClosing,
		domain.OfferingStateClearing, domain.OfferingStateClosed,
		domain.OfferingStateCanceled:
		return true
	}
	return false
}

// ChangeState performs a manual state transition
func (s *OfferingStateService) ChangeState(ctx context.Context, id string, targetState domain.OfferingState, reason string, userID string) error {
	offering, err := s.offerings.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("offering not found: %w", err)
	}

	// Validate transition
	valid := false
	for _, t := range validTransitions[offering.State] {
		if t == targetState {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid state transition from %s to %s", offering.State, targetState)
	}

	now := time.Now()
	changeLog := domain.OfferingChangeLog{
		Field:     "state",
		OldValue:  string(offering.State),
		NewValue:  string(targetState),
		ChangedBy: userID,
		ChangedAt: now,
	}

	offering.PreviousState = offering.State
	offering.State = targetState
	offering.ChangeLog = append(offering.ChangeLog, changeLog)
	offering.UpdatedAt = &now

	return s.offerings.Update(ctx, id, offering)
}

// CheckAndTransitionOfferings checks for date-based automatic transitions
func (s *OfferingStateService) CheckAndTransitionOfferings(ctx context.Context) {
	now := time.Now()

	offerings, err := s.offerings.FindAll(ctx)
	if err != nil {
		slog.Error("failed to query offerings for state transition", "error", err)
		return
	}

	for _, o := range offerings {
		var targetState domain.OfferingState
		var shouldTransition bool

		switch o.State {
		case domain.OfferingStateNew:
			if !o.AnnouncementDate.IsZero() && now.After(o.AnnouncementDate) {
				targetState = domain.OfferingStateUpcoming
				shouldTransition = true
			}
		case domain.OfferingStateUpcoming:
			if !o.BidPeriodStartDate.IsZero() && now.After(o.BidPeriodStartDate) {
				targetState = domain.OfferingStateOpen
				shouldTransition = true
			}
		case domain.OfferingStateOpen, domain.OfferingStateHalted, domain.OfferingStateFrozen:
			if !o.ScheduledCloseDate.IsZero() && now.After(o.ScheduledCloseDate) {
				targetState = domain.OfferingStateClosePending
				shouldTransition = true
			}
		}

		if shouldTransition {
			slog.Info("auto-transitioning offering",
				"symbol", o.Symbol,
				"from", o.State,
				"to", targetState,
			)
			err := s.ChangeState(ctx, o.ID, targetState, "automatic date-based transition", "SYSTEM")
			if err != nil {
				slog.Error("failed to auto-transition offering",
					"symbol", o.Symbol,
					"error", err,
				)
			}
		}
	}
}
