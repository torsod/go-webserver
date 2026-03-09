package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/torsod/go-webserver/internal/domain"
)

type userStore struct {
	pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) UserStore {
	return &userStore{pool: pool}
}

const userColumns = `id, username, password_hash, user_type, is_logged_in, disabled,
	bd_role, bd_firm_id, bd_firm_name, qsr_active,
	firm_accounts, assigned_accounts, read_only_accounts,
	can_cancel_group_orders, can_write_group_orders,
	lm_firm_id, email, display_name,
	created_at, updated_at, last_login_at`

func scanUser(row scannable) (*domain.User, error) {
	u := &domain.User{}
	var faJSON, aaJSON, roaJSON []byte

	err := row.Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.UserType, &u.IsLoggedIn, &u.Disabled,
		&u.BdRole, &u.BdFirmID, &u.BdFirmName, &u.QsrActive,
		&faJSON, &aaJSON, &roaJSON,
		&u.CanCancelGroupOrders, &u.CanWriteGroupOrders,
		&u.LmFirmID, &u.Email, &u.DisplayName,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}

	if faJSON != nil {
		json.Unmarshal(faJSON, &u.FirmAccounts)
	}
	if aaJSON != nil {
		json.Unmarshal(aaJSON, &u.AssignedAccounts)
	}
	if roaJSON != nil {
		json.Unmarshal(roaJSON, &u.ReadOnlyAccounts)
	}

	return u, nil
}

func (s *userStore) FindAll(ctx context.Context) ([]*domain.User, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+userColumns+` FROM users ORDER BY username`)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *userStore) FindByID(ctx context.Context, id string) (*domain.User, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (s *userStore) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE username = $1`, username)
	return scanUser(row)
}

func (s *userStore) Insert(ctx context.Context, user *domain.User) (string, error) {
	faJSON, _ := json.Marshal(user.FirmAccounts)
	aaJSON, _ := json.Marshal(user.AssignedAccounts)
	roaJSON, _ := json.Marshal(user.ReadOnlyAccounts)

	var id string
	err := s.pool.QueryRow(ctx, `
		INSERT INTO users (
			username, password_hash, user_type, is_logged_in, disabled,
			bd_role, bd_firm_id, bd_firm_name, qsr_active,
			firm_accounts, assigned_accounts, read_only_accounts,
			can_cancel_group_orders, can_write_group_orders,
			lm_firm_id, email, display_name, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
		RETURNING id`,
		user.Username, user.PasswordHash, user.UserType, user.IsLoggedIn, user.Disabled,
		user.BdRole, user.BdFirmID, user.BdFirmName, user.QsrActive,
		faJSON, aaJSON, roaJSON,
		user.CanCancelGroupOrders, user.CanWriteGroupOrders,
		user.LmFirmID, user.Email, user.DisplayName, user.CreatedAt,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert user: %w", err)
	}
	return id, nil
}

func (s *userStore) Update(ctx context.Context, id string, user *domain.User) error {
	faJSON, _ := json.Marshal(user.FirmAccounts)
	aaJSON, _ := json.Marshal(user.AssignedAccounts)
	roaJSON, _ := json.Marshal(user.ReadOnlyAccounts)

	_, err := s.pool.Exec(ctx, `
		UPDATE users SET
			user_type=$1, disabled=$2,
			bd_role=$3, bd_firm_id=$4, bd_firm_name=$5, qsr_active=$6,
			firm_accounts=$7, assigned_accounts=$8, read_only_accounts=$9,
			can_cancel_group_orders=$10, can_write_group_orders=$11,
			lm_firm_id=$12, email=$13, display_name=$14, updated_at=NOW()
		WHERE id=$15`,
		user.UserType, user.Disabled,
		user.BdRole, user.BdFirmID, user.BdFirmName, user.QsrActive,
		faJSON, aaJSON, roaJSON,
		user.CanCancelGroupOrders, user.CanWriteGroupOrders,
		user.LmFirmID, user.Email, user.DisplayName, id,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (s *userStore) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return updateFields(ctx, s.pool, "users", id, fields)
}
