package domain

import "time"

// CBMClearingID is the ClearingBid Member Clearing firm identifier
const CBMClearingID = "CBMC"

// UTCTradeRecord represents an NSCC UTC locked-in trade record
type UTCTradeRecord struct {
	TradeID            string  `json:"tradeId"`
	Symbol             string  `json:"symbol"`
	CUSIP              string  `json:"cusip"`
	SettlementDate     string  `json:"settlementDate"`     // YYYYMMDD
	TradeDate          string  `json:"tradeDate"`          // YYYYMMDD
	BuyBrokerDealerID  string  `json:"buyBrokerDealerId"`
	SellBrokerDealerID string  `json:"sellBrokerDealerId"` // Always CBM Clearing
	Quantity           int64   `json:"quantity"`
	Price              float64 `json:"price"`
	NetAmount          float64 `json:"netAmount"`
	SellingConcession  float64 `json:"sellingConcession"`
	GrossSpread        float64 `json:"grossSpread"`
	AccountID          string  `json:"accountId"`
	ExecInst           string  `json:"execInst"`
}

// DTCTrackingRecord represents a DTC IPO Tracking System record
type DTCTrackingRecord struct {
	TradeID        string  `json:"tradeId"`
	Symbol         string  `json:"symbol"`
	CUSIP          string  `json:"cusip"`
	BDFirmID       string  `json:"bdFirmId"`
	Account        string  `json:"account"`
	Quantity       int64   `json:"quantity"`
	Price          float64 `json:"price"`
	IsDtcTracked   bool    `json:"isDtcTracked"`
	SettlementDate string  `json:"settlementDate"` // YYYYMMDD
}

// SettlementFileResult contains the complete settlement output
type SettlementFileResult struct {
	OfferingID          string             `json:"offeringId"`
	Symbol              string             `json:"symbol"`
	AllocationSessionID string             `json:"allocationSessionId"`
	UTCRecords          []UTCTradeRecord   `json:"utcRecords"`
	DTCRecords          []DTCTrackingRecord `json:"dtcRecords"`
	UTCCSV              string             `json:"utcCsv"`
	DTCCSV              string             `json:"dtcCsv"`
	GeneratedAt         time.Time          `json:"generatedAt"`
}
