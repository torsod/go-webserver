package domain

// AllocationEntry represents a priority level in allocation
type AllocationEntry struct {
	Priority        int     `json:"priority"`
	AccumulatedSize int64   `json:"accumulatedSize"`
	Percentage      float64 `json:"percentage,omitempty"`
}

// OrderAllocation represents the allocation for a single order
type OrderAllocation struct {
	OrderID           string `json:"orderId"`
	AllocatedQuantity int64  `json:"allocatedQuantity"`
	IsPartial         bool   `json:"isPartial"`
}

// UnallocatedOrder represents an order that was not allocated
type UnallocatedOrder struct {
	OrderID       string  `json:"orderId"`
	Quantity      int64   `json:"quantity"`
	Price         float64 `json:"price"`
	PriorityGroup int     `json:"priorityGroup"`
	UserID        string  `json:"userId"`
}

// AllocationResult contains the full allocation result
type AllocationResult struct {
	Allocations        []OrderAllocation `json:"allocations"`
	UnallocatedOrders  int               `json:"unallocatedOrders"`
	UnallocatedList    []UnallocatedOrder `json:"unallocatedOrdersList"`
	PartialAllocations int               `json:"partialAllocations"`
	TotalAllocated     int64             `json:"totalAllocated"`
}

// ClosingWindowsStatus for closing windows state
type ClosingWindowsStatus string

const (
	ClosingWindowsWaiting         ClosingWindowsStatus = "WAITING"
	ClosingWindowsInCW1           ClosingWindowsStatus = "IN_CW1"
	ClosingWindowsInCW2FirstHalf  ClosingWindowsStatus = "IN_CW2_FIRST_HALF"
	ClosingWindowsInCW2SecondHalf ClosingWindowsStatus = "IN_CW2_SECOND_HALF"
	ClosingWindowsTriggered       ClosingWindowsStatus = "TRIGGERED"
	ClosingWindowsHalted          ClosingWindowsStatus = "HALTED"
	ClosingWindowsFrozen          ClosingWindowsStatus = "FROZEN"
)
