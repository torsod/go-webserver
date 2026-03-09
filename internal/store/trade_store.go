package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/torsod/go-webserver/internal/domain"
)

type tradeStore struct {
	pool *pgxpool.Pool
}

func NewTradeStore(pool *pgxpool.Pool) TradeStore {
	return &tradeStore{pool: pool}
}

const tradeColumns = `id, symbol, trade_type, quantity, leaves_qty, price,
	bid_price, bid_quantity, priority_group, order_type, account, exec_inst,
	user_id, bd_firm_id, order_id, allocation_id,
	is_dtc_tracked, selling_concession_amount, gross_spread_amount,
	status, "timestamp", busted_at, canceled_at`

func scanTrade(row scannable) (*domain.Trade, error) {
	t := &domain.Trade{}
	err := row.Scan(
		&t.ID, &t.Symbol, &t.TradeType, &t.Quantity, &t.LeavesQty, &t.Price,
		&t.BidPrice, &t.BidQuantity, &t.PriorityGroup, &t.OrderType, &t.Account, &t.ExecInst,
		&t.UserID, &t.BDFirmID, &t.OrderID, &t.AllocationID,
		&t.IsDtcTracked, &t.SellingConcessionAmount, &t.GrossSpreadAmount,
		&t.Status, &t.Timestamp, &t.BustedAt, &t.CanceledAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan trade: %w", err)
	}
	return t, nil
}

func (s *tradeStore) FindBySymbol(ctx context.Context, symbol string) ([]*domain.Trade, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+tradeColumns+` FROM trades WHERE symbol = $1 ORDER BY "timestamp" DESC`, symbol)
	if err != nil {
		return nil, fmt.Errorf("query trades by symbol: %w", err)
	}
	defer rows.Close()

	var trades []*domain.Trade
	for rows.Next() {
		t, err := scanTrade(rows)
		if err != nil {
			return nil, err
		}
		trades = append(trades, t)
	}
	return trades, nil
}

func (s *tradeStore) FindByAllocationID(ctx context.Context, allocationID string) ([]*domain.Trade, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+tradeColumns+` FROM trades WHERE allocation_id = $1 ORDER BY "timestamp"`, allocationID)
	if err != nil {
		return nil, fmt.Errorf("query trades by allocation ID: %w", err)
	}
	defer rows.Close()

	var trades []*domain.Trade
	for rows.Next() {
		t, err := scanTrade(rows)
		if err != nil {
			return nil, err
		}
		trades = append(trades, t)
	}
	return trades, nil
}

func (s *tradeStore) Insert(ctx context.Context, trade *domain.Trade) (string, error) {
	var id string
	err := s.pool.QueryRow(ctx, `
		INSERT INTO trades (
			symbol, trade_type, quantity, leaves_qty, price,
			bid_price, bid_quantity, priority_group, order_type, account, exec_inst,
			user_id, bd_firm_id, order_id, allocation_id,
			is_dtc_tracked, selling_concession_amount, gross_spread_amount,
			status, "timestamp"
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)
		RETURNING id`,
		trade.Symbol, trade.TradeType, trade.Quantity, trade.LeavesQty, trade.Price,
		trade.BidPrice, trade.BidQuantity, trade.PriorityGroup, trade.OrderType, trade.Account, trade.ExecInst,
		trade.UserID, trade.BDFirmID, trade.OrderID, trade.AllocationID,
		trade.IsDtcTracked, trade.SellingConcessionAmount, trade.GrossSpreadAmount,
		trade.Status, trade.Timestamp,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert trade: %w", err)
	}
	return id, nil
}

func (s *tradeStore) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return updateFields(ctx, s.pool, "trades", id, fields)
}

// Allocation Session Store

type allocationSessionStore struct {
	pool *pgxpool.Pool
}

func NewAllocationSessionStore(pool *pgxpool.Pool) AllocationSessionStore {
	return &allocationSessionStore{pool: pool}
}

const sessionColumns = `id, symbol, offering_price, offering_size, total_allocated, allocation_method, lm_user,
	preferential_allocated, preferential_bidders, group_summaries, min_qty_excluded,
	status, created_at, busted_at, confirmed_at`

func scanSession(row scannable) (*domain.AllocationSession, error) {
	s := &domain.AllocationSession{}
	var gsJSON []byte

	err := row.Scan(
		&s.ID, &s.Symbol, &s.OfferingPrice, &s.OfferingSize, &s.TotalAllocated,
		&s.AllocationMethod, &s.LMUser,
		&s.PreferentialAllocated, &s.PreferentialBidders, &gsJSON, &s.MinQtyExcluded,
		&s.Status, &s.CreatedAt, &s.BustedAt, &s.ConfirmedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan allocation session: %w", err)
	}

	if gsJSON != nil {
		json.Unmarshal(gsJSON, &s.GroupSummaries)
	}

	return s, nil
}

func (s *allocationSessionStore) FindAll(ctx context.Context) ([]*domain.AllocationSession, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+sessionColumns+` FROM allocation_sessions ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query allocation sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*domain.AllocationSession
	for rows.Next() {
		sess, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, sess)
	}
	return sessions, nil
}

func (s *allocationSessionStore) FindByID(ctx context.Context, id string) (*domain.AllocationSession, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+sessionColumns+` FROM allocation_sessions WHERE id = $1`, id)
	return scanSession(row)
}

func (s *allocationSessionStore) FindBySymbol(ctx context.Context, symbol string) ([]*domain.AllocationSession, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+sessionColumns+` FROM allocation_sessions WHERE symbol = $1 ORDER BY created_at DESC`, symbol)
	if err != nil {
		return nil, fmt.Errorf("query sessions by symbol: %w", err)
	}
	defer rows.Close()

	var sessions []*domain.AllocationSession
	for rows.Next() {
		sess, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, sess)
	}
	return sessions, nil
}

func (s *allocationSessionStore) Insert(ctx context.Context, session *domain.AllocationSession) (string, error) {
	gsJSON, _ := json.Marshal(session.GroupSummaries)

	var id string
	err := s.pool.QueryRow(ctx, `
		INSERT INTO allocation_sessions (
			symbol, offering_price, offering_size, total_allocated, allocation_method, lm_user,
			preferential_allocated, preferential_bidders, group_summaries, min_qty_excluded,
			status, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING id`,
		session.Symbol, session.OfferingPrice, session.OfferingSize, session.TotalAllocated,
		session.AllocationMethod, session.LMUser,
		session.PreferentialAllocated, session.PreferentialBidders, gsJSON, session.MinQtyExcluded,
		session.Status, session.CreatedAt,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert allocation session: %w", err)
	}
	return id, nil
}

func (s *allocationSessionStore) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return updateFields(ctx, s.pool, "allocation_sessions", id, fields)
}
