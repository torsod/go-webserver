package domain

import "time"

// UserType classification per spec
type UserType string

const (
	UserTypeLeadManager      UserType = "LM"
	UserTypeIssuer           UserType = "IS"
	UserTypeBrokerDealer     UserType = "BD"
	UserTypeMarketOperations UserType = "MO"
)

// BdRole sub-roles per spec Section 4.2
type BdRole string

const (
	BdRoleAdmin        BdRole = "ADMIN"
	BdRoleBroker       BdRole = "BROKER"
	BdRoleOMSUser      BdRole = "OMS_USER"
	BdRoleGroupManager BdRole = "GROUP_MANAGER"
	BdRoleMasterBroker BdRole = "MASTER_BROKER"
)

// FirmAccount per spec Section 4.2
type FirmAccount struct {
	AccountID string `json:"accountId"`
	Name      string `json:"name"`
	IsActive  bool   `json:"isActive"`
}

// User profile
type User struct {
	ID           string   `json:"id" db:"id"`
	Username     string   `json:"username" db:"username"`
	PasswordHash string   `json:"-" db:"password_hash"`
	UserType     UserType `json:"userType" db:"user_type"`
	IsLoggedIn   bool     `json:"isLoggedIn" db:"is_logged_in"`
	Disabled     bool     `json:"disabled" db:"disabled"`

	// BD-specific fields
	BdRole               BdRole        `json:"bdRole,omitempty" db:"bd_role"`
	BdFirmID             string        `json:"bdFirmId,omitempty" db:"bd_firm_id"`
	BdFirmName           string        `json:"bdFirmName,omitempty" db:"bd_firm_name"`
	QsrActive            bool          `json:"qsrActive,omitempty" db:"qsr_active"`
	FirmAccounts         []FirmAccount `json:"firmAccounts,omitempty" db:"firm_accounts"`
	AssignedAccounts     []string      `json:"assignedAccounts,omitempty" db:"assigned_accounts"`
	ReadOnlyAccounts     []string      `json:"readOnlyAccounts,omitempty" db:"read_only_accounts"`
	CanCancelGroupOrders bool          `json:"canCancelGroupOrders,omitempty" db:"can_cancel_group_orders"`
	CanWriteGroupOrders  bool          `json:"canWriteGroupOrders,omitempty" db:"can_write_group_orders"`

	// LM-specific fields
	LmFirmID string `json:"lmFirmId,omitempty" db:"lm_firm_id"`

	// Contact info
	Email       string `json:"email,omitempty" db:"email"`
	DisplayName string `json:"displayName,omitempty" db:"display_name"`

	// Timestamps
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty" db:"last_login_at"`
}
