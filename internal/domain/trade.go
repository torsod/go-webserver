package domain

import "time"

// TradeType
type TradeType string

const (
	TradeTypeSell TradeType = "SELL"
	TradeTypeBuy  TradeType = "BUY"
)

// TradeStatus
type TradeStatus string

const (
	TradeStatusNew         TradeStatus = "NEW"
	TradeStatusFilled      TradeStatus = "FILLED"
	TradeStatusPartialFill TradeStatus = "PARTIAL_FILL"
	TradeStatusBusted      TradeStatus = "BUSTED"
	TradeStatusCanceled    TradeStatus = "CANCELED"
)

// SessionStatus for allocation sessions
type SessionStatus string

const (
	SessionStatusTemporary SessionStatus = "TEMPORARY"
	SessionStatusConfirmed SessionStatus = "CONFIRMED"
	SessionStatusBusted    SessionStatus = "BUSTED"
)

// Trade per CB-PDP-Gen-2 spec Section 8
type Trade struct {
	ID     string `json:"id" db:"id"`
	Symbol string `json:"symbol" db:"symbol"`

	TradeType TradeType `json:"tradeType" db:"trade_type"`
	Quantity  int64     `json:"quantity" db:"quantity"`
	LeavesQty int64    `json:"leavesQty" db:"leaves_qty"`
	Price     float64  `json:"price" db:"price"`

	// Original order reference
	BidPrice      *float64 `json:"bidPrice,omitempty" db:"bid_price"`
	BidQuantity   *int64   `json:"bidQuantity,omitempty" db:"bid_quantity"`
	PriorityGroup *int     `json:"priorityGroup,omitempty" db:"priority_group"`
	OrderType     string   `json:"orderType,omitempty" db:"order_type"`
	Account       string   `json:"account,omitempty" db:"account"`
	ExecInst      string   `json:"execInst,omitempty" db:"exec_inst"`

	// User/Ownership
	UserID       string `json:"userId" db:"user_id"`
	BDFirmID     string `json:"bdFirmId,omitempty" db:"bd_firm_id"`
	OrderID      string `json:"orderId,omitempty" db:"order_id"`
	AllocationID string `json:"allocationId" db:"allocation_id"`

	// Settlement (spec Section 8.3)
	IsDtcTracked             *bool    `json:"isDtcTracked,omitempty" db:"is_dtc_tracked"`
	SellingConcessionAmount  *float64 `json:"sellingConcessionAmount,omitempty" db:"selling_concession_amount"`
	GrossSpreadAmount        *float64 `json:"grossSpreadAmount,omitempty" db:"gross_spread_amount"`

	// Status
	Status    TradeStatus `json:"status" db:"status"`
	Timestamp time.Time   `json:"timestamp" db:"timestamp"`
	BustedAt  *time.Time  `json:"bustedAt,omitempty" db:"busted_at"`
	CanceledAt *time.Time `json:"canceledAt,omitempty" db:"canceled_at"`
}

// GroupSummary for allocation session reporting
type GroupSummary struct {
	PriorityGroup  int     `json:"priorityGroup"`
	TotalDemand    int64   `json:"totalDemand"`
	TotalAllocated int64   `json:"totalAllocated"`
	FillPercentage float64 `json:"fillPercentage"`
	OrderCount     int     `json:"orderCount"`
}

// AllocationSession per spec Section 7
type AllocationSession struct {
	ID             string `json:"id" db:"id"`
	Symbol         string `json:"symbol" db:"symbol"`
	OfferingPrice  float64 `json:"offeringPrice" db:"offering_price"`
	OfferingSize   int64   `json:"offeringSize" db:"offering_size"`
	TotalAllocated int64   `json:"totalAllocated" db:"total_allocated"`
	AllocationMethod string `json:"allocationMethod" db:"allocation_method"`
	LMUser         string  `json:"lmUser" db:"lm_user"`

	// Preferential allocation summary
	PreferentialAllocated *int64 `json:"preferentialAllocated,omitempty" db:"preferential_allocated"`
	PreferentialBidders   *int   `json:"preferentialBidders,omitempty" db:"preferential_bidders"`

	// Priority group summaries - stored as JSONB
	GroupSummaries []GroupSummary `json:"groupSummaries,omitempty" db:"group_summaries"`

	// MinQty impact
	MinQtyExcluded *int `json:"minQtyExcluded,omitempty" db:"min_qty_excluded"`

	// Status
	Status      SessionStatus `json:"status" db:"status"`
	CreatedAt   time.Time     `json:"createdAt" db:"created_at"`
	BustedAt    *time.Time    `json:"bustedAt,omitempty" db:"busted_at"`
	ConfirmedAt *time.Time    `json:"confirmedAt,omitempty" db:"confirmed_at"`
}
