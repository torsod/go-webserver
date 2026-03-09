package domain

// OrderbookEntry represents a price level in the orderbook
type OrderbookEntry struct {
	Price                float64 `json:"price"`
	AccumulatedSize      int64   `json:"accumulatedSize"`
	TotalAccumulatedSize int64   `json:"totalAccumulatedSize"`
	DollarValue          float64 `json:"dollarValue"`
	TotalDollarValue     float64 `json:"totalDollarValue"`
	OrderCount           int     `json:"orderCount"`
}

// OrderbookSummary is a complete two-sided orderbook view
type OrderbookSummary struct {
	Symbol            string           `json:"symbol"`
	BidEntries        []OrderbookEntry `json:"bidEntries"`
	OfferEntries      []OrderbookEntry `json:"offerEntries"`
	ClearingPrice     *float64         `json:"clearingPrice"`
	Cprice            *float64         `json:"cprice"`
	AVWAP             *float64         `json:"aVWAP"`
	RVWAP             *float64         `json:"rVWAP"`
	SubscriptionRatio *float64         `json:"subscriptionRatio"`
	TotalBidQuantity  int64            `json:"totalBidQuantity"`
	TotalOfferQuantity int64           `json:"totalOfferQuantity"`
	TotalBidOrders    int              `json:"totalBidOrders"`
	TotalOfferOrders  int              `json:"totalOfferOrders"`
}
