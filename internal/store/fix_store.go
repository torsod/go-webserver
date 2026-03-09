package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/torsod/go-webserver/internal/domain"
)

// FIX Session Store

type fixSessionStore struct {
	pool *pgxpool.Pool
}

func NewFIXSessionStore(pool *pgxpool.Pool) FIXSessionStore {
	return &fixSessionStore{pool: pool}
}

func (s *fixSessionStore) FindByID(ctx context.Context, id string) (*domain.FIXSession, error) {
	fs := &domain.FIXSession{}
	err := s.pool.QueryRow(ctx, `
		SELECT id, session_id, host, port, sender_comp_id, target_comp_id,
			connected, simulated, messages_sent, messages_received, created_at, updated_at
		FROM fix_sessions WHERE id = $1`, id).Scan(
		&fs.ID, &fs.SessionID, &fs.Host, &fs.Port, &fs.SenderCompID, &fs.TargetCompID,
		&fs.Connected, &fs.Simulated, &fs.MessagesSent, &fs.MessagesRecv, &fs.CreatedAt, &fs.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("find fix session: %w", err)
	}
	return fs, nil
}

func (s *fixSessionStore) FindActive(ctx context.Context) (*domain.FIXSession, error) {
	fs := &domain.FIXSession{}
	err := s.pool.QueryRow(ctx, `
		SELECT id, session_id, host, port, sender_comp_id, target_comp_id,
			connected, simulated, messages_sent, messages_received, created_at, updated_at
		FROM fix_sessions WHERE connected = true ORDER BY created_at DESC LIMIT 1`).Scan(
		&fs.ID, &fs.SessionID, &fs.Host, &fs.Port, &fs.SenderCompID, &fs.TargetCompID,
		&fs.Connected, &fs.Simulated, &fs.MessagesSent, &fs.MessagesRecv, &fs.CreatedAt, &fs.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("find active fix session: %w", err)
	}
	return fs, nil
}

func (s *fixSessionStore) Insert(ctx context.Context, session *domain.FIXSession) (string, error) {
	var id string
	err := s.pool.QueryRow(ctx, `
		INSERT INTO fix_sessions (session_id, host, port, sender_comp_id, target_comp_id, connected, simulated)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		session.SessionID, session.Host, session.Port, session.SenderCompID, session.TargetCompID,
		session.Connected, session.Simulated,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert fix session: %w", err)
	}
	return id, nil
}

func (s *fixSessionStore) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return updateFields(ctx, s.pool, "fix_sessions", id, fields)
}

// FIX Order Store

type fixOrderStore struct {
	pool *pgxpool.Pool
}

func NewFIXOrderStore(pool *pgxpool.Pool) FIXOrderStore {
	return &fixOrderStore{pool: pool}
}

func (s *fixOrderStore) Insert(ctx context.Context, order *domain.FIXOrder) (string, error) {
	var id string
	err := s.pool.QueryRow(ctx, `
		INSERT INTO fix_orders (cl_ord_id, session_id, symbol, side, quantity, price,
			ord_type, time_in_force, account, main_order_id, status, transact_time)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id`,
		order.ClOrdID, order.SessionID, order.Symbol, order.Side, order.Quantity, order.Price,
		order.OrdType, order.TimeInForce, order.Account, order.MainOrderID,
		order.Status, order.TransactTime,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert fix order: %w", err)
	}
	return id, nil
}

func (s *fixOrderStore) FindBySessionID(ctx context.Context, sessionID string) ([]*domain.FIXOrder, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, cl_ord_id, session_id, symbol, side, quantity, price,
			ord_type, time_in_force, account, main_order_id, status, transact_time, created_at
		FROM fix_orders WHERE session_id = $1 ORDER BY created_at DESC`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("query fix orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.FIXOrder
	for rows.Next() {
		fo := &domain.FIXOrder{}
		err := rows.Scan(
			&fo.ID, &fo.ClOrdID, &fo.SessionID, &fo.Symbol, &fo.Side, &fo.Quantity, &fo.Price,
			&fo.OrdType, &fo.TimeInForce, &fo.Account, &fo.MainOrderID, &fo.Status,
			&fo.TransactTime, &fo.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan fix order: %w", err)
		}
		orders = append(orders, fo)
	}
	return orders, nil
}

// FIX Log Store

type fixLogStore struct {
	pool *pgxpool.Pool
}

func NewFIXLogStore(pool *pgxpool.Pool) FIXLogStore {
	return &fixLogStore{pool: pool}
}

func (s *fixLogStore) Insert(ctx context.Context, log *domain.FIXLog) (string, error) {
	var id string
	err := s.pool.QueryRow(ctx, `
		INSERT INTO fix_logs (session_id, "timestamp", level, message, direction, raw_data)
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		log.SessionID, log.Timestamp, log.Level, log.Message, log.Direction, log.RawData,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert fix log: %w", err)
	}
	return id, nil
}

func (s *fixLogStore) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM fix_logs WHERE "timestamp" < NOW() - INTERVAL '1 day' * $1`, days)
	if err != nil {
		return 0, fmt.Errorf("delete old fix logs: %w", err)
	}
	return tag.RowsAffected(), nil
}
