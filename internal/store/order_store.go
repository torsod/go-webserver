package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/torsod/go-webserver/internal/domain"
)

type orderStore struct {
	pool *pgxpool.Pool
}

func NewOrderStore(pool *pgxpool.Pool) OrderStore {
	return &orderStore{pool: pool}
}

const orderColumns = `id, symbol, side, order_type, quantity, price, min_qty, account, exec_inst,
	priority_group, "timestamp", original_entry_time, order_sequence,
	user_id, bd_firm_id, entered_by, status,
	seasoning_expires_at, time_window_at_entry,
	halt_reason, halted_at, pre_halt_timestamp,
	allocated_quantity, allocation_session_id,
	canceled_at, cancel_reason,
	created_at, updated_at, modification_history`

func scanOrder(row scannable) (*domain.Order, error) {
	o := &domain.Order{}
	var modHistJSON []byte

	err := row.Scan(
		&o.ID, &o.Symbol, &o.Side, &o.OrderType, &o.Quantity, &o.Price, &o.MinQty, &o.Account, &o.ExecInst,
		&o.PriorityGroup, &o.Timestamp, &o.OriginalEntryTime, &o.OrderSequence,
		&o.UserID, &o.BDFirmID, &o.EnteredBy, &o.Status,
		&o.SeasoningExpiresAt, &o.TimeWindowAtEntry,
		&o.HaltReason, &o.HaltedAt, &o.PreHaltTimestamp,
		&o.AllocatedQuantity, &o.AllocationSessionID,
		&o.CanceledAt, &o.CancelReason,
		&o.CreatedAt, &o.UpdatedAt, &modHistJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("scan order: %w", err)
	}

	if modHistJSON != nil {
		json.Unmarshal(modHistJSON, &o.ModificationHistory)
	}

	return o, nil
}

func (s *orderStore) FindAll(ctx context.Context) ([]*domain.Order, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+orderColumns+` FROM orders ORDER BY order_sequence DESC LIMIT 5000`)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *orderStore) FindBySymbol(ctx context.Context, symbol string) ([]*domain.Order, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+orderColumns+` FROM orders WHERE symbol = $1 ORDER BY order_sequence DESC`, symbol)
	if err != nil {
		return nil, fmt.Errorf("query orders by symbol: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *orderStore) FindActive(ctx context.Context) ([]*domain.Order, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT `+orderColumns+` FROM orders WHERE status NOT IN ('CANCELED', 'REJECTED') ORDER BY order_sequence DESC LIMIT 5000`)
	if err != nil {
		return nil, fmt.Errorf("query active orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *orderStore) FindActiveBySymbol(ctx context.Context, symbol string) ([]*domain.Order, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT `+orderColumns+` FROM orders WHERE symbol = $1 AND status NOT IN ('CANCELED', 'REJECTED') ORDER BY order_sequence DESC`,
		symbol)
	if err != nil {
		return nil, fmt.Errorf("query active orders by symbol: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *orderStore) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+orderColumns+` FROM orders WHERE id = $1`, id)
	return scanOrder(row)
}

func (s *orderStore) FindForAllocation(ctx context.Context, symbol string, minPrice float64) ([]*domain.Order, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT `+orderColumns+` FROM orders
		WHERE symbol = $1 AND price >= $2 AND status NOT IN ('CANCELED', 'REJECTED') AND side != 'OFFER'
		ORDER BY priority_group ASC, price DESC, order_sequence ASC`,
		symbol, minPrice)
	if err != nil {
		return nil, fmt.Errorf("query orders for allocation: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *orderStore) Insert(ctx context.Context, order *domain.Order) (string, error) {
	modHistJSON, _ := json.Marshal(order.ModificationHistory)

	var id string
	err := s.pool.QueryRow(ctx, `
		INSERT INTO orders (
			symbol, side, order_type, quantity, price, min_qty, account, exec_inst,
			priority_group, "timestamp", original_entry_time, order_sequence,
			user_id, bd_firm_id, entered_by, status,
			seasoning_expires_at, time_window_at_entry,
			halt_reason, halted_at, pre_halt_timestamp,
			allocated_quantity, allocation_session_id,
			canceled_at, cancel_reason,
			created_at, updated_at, modification_history
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28)
		RETURNING id`,
		order.Symbol, order.Side, order.OrderType, order.Quantity, order.Price,
		order.MinQty, order.Account, order.ExecInst,
		order.PriorityGroup, order.Timestamp, order.OriginalEntryTime, order.OrderSequence,
		order.UserID, order.BDFirmID, order.EnteredBy, order.Status,
		order.SeasoningExpiresAt, order.TimeWindowAtEntry,
		order.HaltReason, order.HaltedAt, order.PreHaltTimestamp,
		order.AllocatedQuantity, order.AllocationSessionID,
		order.CanceledAt, order.CancelReason,
		order.CreatedAt, order.UpdatedAt, modHistJSON,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert order: %w", err)
	}
	return id, nil
}

func (s *orderStore) Update(ctx context.Context, id string, order *domain.Order) error {
	modHistJSON, _ := json.Marshal(order.ModificationHistory)

	_, err := s.pool.Exec(ctx, `
		UPDATE orders SET
			quantity=$1, price=$2, min_qty=$3, status=$4,
			"timestamp"=$5, priority_group=$6,
			seasoning_expires_at=$7, time_window_at_entry=$8,
			halt_reason=$9, halted_at=$10, pre_halt_timestamp=$11,
			allocated_quantity=$12, allocation_session_id=$13,
			canceled_at=$14, cancel_reason=$15,
			updated_at=NOW(), modification_history=$16
		WHERE id=$17`,
		order.Quantity, order.Price, order.MinQty, order.Status,
		order.Timestamp, order.PriorityGroup,
		order.SeasoningExpiresAt, order.TimeWindowAtEntry,
		order.HaltReason, order.HaltedAt, order.PreHaltTimestamp,
		order.AllocatedQuantity, order.AllocationSessionID,
		order.CanceledAt, order.CancelReason,
		modHistJSON, id,
	)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}
	return nil
}

func (s *orderStore) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return updateFields(ctx, s.pool, "orders", id, fields)
}

func (s *orderStore) NextSequence(ctx context.Context) (int64, error) {
	var seq int64
	err := s.pool.QueryRow(ctx, "SELECT nextval('order_sequence_seq')").Scan(&seq)
	if err != nil {
		return 0, fmt.Errorf("next order sequence: %w", err)
	}
	return seq, nil
}
