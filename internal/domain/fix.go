package domain

import "time"

// FIXSession tracks FIX protocol connection sessions
type FIXSession struct {
	ID           string    `json:"id" db:"id"`
	SessionID    string    `json:"sessionId" db:"session_id"`
	Host         string    `json:"host" db:"host"`
	Port         int       `json:"port" db:"port"`
	SenderCompID string    `json:"senderCompId" db:"sender_comp_id"`
	TargetCompID string    `json:"targetCompId" db:"target_comp_id"`
	Connected    bool      `json:"connected" db:"connected"`
	Simulated    bool      `json:"simulated" db:"simulated"`
	MessagesSent int       `json:"messagesSent" db:"messages_sent"`
	MessagesRecv int       `json:"messagesReceived" db:"messages_received"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

// FIXOrder tracks orders submitted via FIX protocol
type FIXOrder struct {
	ID          string    `json:"id" db:"id"`
	ClOrdID     string    `json:"clOrdId" db:"cl_ord_id"`
	SessionID   string    `json:"sessionId" db:"session_id"`
	Symbol      string    `json:"symbol" db:"symbol"`
	Side        string    `json:"side" db:"side"`
	Quantity    int64     `json:"quantity" db:"quantity"`
	Price       float64   `json:"price" db:"price"`
	OrdType     string    `json:"ordType" db:"ord_type"`
	TimeInForce string    `json:"timeInForce" db:"time_in_force"`
	Account     string    `json:"account,omitempty" db:"account"`
	MainOrderID string    `json:"mainOrderId,omitempty" db:"main_order_id"`
	Status      string    `json:"status" db:"status"`
	TransactTime time.Time `json:"transactTime" db:"transact_time"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

// FIXLog message log entry
type FIXLog struct {
	ID        string    `json:"id" db:"id"`
	SessionID string    `json:"sessionId" db:"session_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Level     string    `json:"level" db:"level"`
	Message   string    `json:"message" db:"message"`
	Direction string    `json:"direction" db:"direction"` // IN or OUT
	RawData   string    `json:"rawData,omitempty" db:"raw_data"`
}
