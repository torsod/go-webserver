package domain

import "time"

// OrderSide: Bid (buy) or Offer (sell)
type OrderSide string

const (
	OrderSideBid   OrderSide = "BID"
	OrderSideOffer OrderSide = "OFFER"
)

// OrderType per spec Sections 3.2, 3.3
type OrderType string

const (
	OrderTypeCompetitive   OrderType = "COMPETITIVE"
	OrderTypePreferential  OrderType = "PREFERENTIAL"
	OrderTypeSecondarySell OrderType = "SECONDARY_SELL"
)

// OrderStatus lifecycle
type OrderStatus string

const (
	OrderStatusActive      OrderStatus = "ACTIVE"
	OrderStatusCanceled    OrderStatus = "CANCELED"
	OrderStatusHalted      OrderStatus = "HALTED"
	OrderStatusFilled      OrderStatus = "FILLED"
	OrderStatusPartialFill OrderStatus = "PARTIAL_FILL"
	OrderStatusRejected    OrderStatus = "REJECTED"
)

// ExecInst flags per spec
type ExecInst string

const (
	ExecInstNone         ExecInst = ""
	ExecInstPreferential ExecInst = "p"
	ExecInstDTCTracking  ExecInst = "t"
)

// OrderModification tracks modifications for audit per spec
type OrderModification struct {
	Field               string      `json:"field"`
	OldValue            interface{} `json:"oldValue"`
	NewValue            interface{} `json:"newValue"`
	ModifiedBy          string      `json:"modifiedBy"`
	ModifiedAt          time.Time   `json:"modifiedAt"`
	ImpactsTimePriority bool        `json:"impactsTimePriority"`
}

// Order per CB-PDP-Gen-2 spec
type Order struct {
	ID string `json:"id" db:"id"`

	// Core fields
	Symbol    string    `json:"symbol" db:"symbol"`
	Side      OrderSide `json:"side" db:"side"`
	OrderType OrderType `json:"orderType" db:"order_type"`
	Quantity  int64     `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`

	// MinQty (spec Section 3.2.1.1)
	MinQty *int64 `json:"minQty,omitempty" db:"min_qty"`

	// Account & Tracking (spec Sections 3.2.1.2, 3.2.2)
	Account  string   `json:"account,omitempty" db:"account"`
	ExecInst ExecInst `json:"execInst" db:"exec_inst"`

	// Priority & Time
	PriorityGroup     int       `json:"priorityGroup" db:"priority_group"`
	Timestamp         time.Time `json:"timestamp" db:"timestamp"`
	OriginalEntryTime time.Time `json:"originalEntryTime" db:"original_entry_time"`
	OrderSequence     int64     `json:"orderSequence" db:"order_sequence"`

	// User/Ownership
	UserID   string `json:"userId" db:"user_id"`
	BDFirmID string `json:"bdFirmId,omitempty" db:"bd_firm_id"`
	EnteredBy string `json:"enteredBy,omitempty" db:"entered_by"`

	// Status
	Status OrderStatus `json:"status" db:"status"`

	// Seasoning Period tracking (spec Section 3.9)
	SeasoningExpiresAt *time.Time `json:"seasoningExpiresAt,omitempty" db:"seasoning_expires_at"`
	TimeWindowAtEntry  int        `json:"timeWindowAtEntry" db:"time_window_at_entry"`

	// Halted Order tracking (spec Section 3.5)
	HaltReason       string     `json:"haltReason,omitempty" db:"halt_reason"`
	HaltedAt         *time.Time `json:"haltedAt,omitempty" db:"halted_at"`
	PreHaltTimestamp  *time.Time `json:"preHaltTimestamp,omitempty" db:"pre_halt_timestamp"`

	// Allocation results (populated after allocation)
	AllocatedQuantity   *int64  `json:"allocatedQuantity,omitempty" db:"allocated_quantity"`
	AllocationSessionID *string `json:"allocationSessionId,omitempty" db:"allocation_session_id"`

	// Cancellation
	CanceledAt   *time.Time `json:"canceledAt,omitempty" db:"canceled_at"`
	CancelReason string     `json:"cancelReason,omitempty" db:"cancel_reason"`

	// Audit
	CreatedAt           time.Time           `json:"createdAt" db:"created_at"`
	UpdatedAt           *time.Time          `json:"updatedAt,omitempty" db:"updated_at"`
	ModificationHistory []OrderModification `json:"modificationHistory,omitempty" db:"modification_history"`
}

// IsOrderCanceled checks if order is canceled
func IsOrderCanceled(o *Order) bool {
	return o.Status == OrderStatusCanceled
}

// IsOrderActive checks if order is active
func IsOrderActive(o *Order) bool {
	return o.Status == OrderStatusActive
}
