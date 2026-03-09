package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/store"
)

// SettlementService generates settlement files
type SettlementService struct {
	sessions  store.AllocationSessionStore
	trades    store.TradeStore
	offerings store.OfferingStore
}

func NewSettlementService(sessions store.AllocationSessionStore, trades store.TradeStore, offerings store.OfferingStore) *SettlementService {
	return &SettlementService{sessions: sessions, trades: trades, offerings: offerings}
}

// GenerateSettlementFiles generates UTC and DTC settlement files for an allocation session
func (s *SettlementService) GenerateSettlementFiles(ctx context.Context, sessionID string) (*domain.SettlementFileResult, error) {
	session, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("allocation session not found: %s", sessionID)
	}
	if session.Status != domain.SessionStatusConfirmed {
		return nil, fmt.Errorf("session must be CONFIRMED (current: %s)", session.Status)
	}

	offering, err := s.offerings.FindBySymbol(ctx, session.Symbol)
	if err != nil {
		return nil, fmt.Errorf("offering not found: %s", session.Symbol)
	}

	trades, err := s.trades.FindByAllocationID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find trades: %w", err)
	}

	settlementDate := ""
	if offering.SettlementDate != nil {
		settlementDate = formatDate(*offering.SettlementDate)
	}
	tradeDate := formatDate(time.Now())

	// Build UTC records
	var utcRecords []domain.UTCTradeRecord
	for _, t := range trades {
		if t.TradeType != domain.TradeTypeBuy || t.Quantity <= 0 || t.Status == domain.TradeStatusBusted {
			continue
		}
		sellingConc := float64(0)
		if t.SellingConcessionAmount != nil {
			sellingConc = *t.SellingConcessionAmount
		}
		grossSpread := float64(0)
		if t.GrossSpreadAmount != nil {
			grossSpread = *t.GrossSpreadAmount
		}

		bdID := t.BDFirmID
		if bdID == "" {
			bdID = t.UserID
		}

		utcRecords = append(utcRecords, domain.UTCTradeRecord{
			TradeID:            t.ID,
			Symbol:             t.Symbol,
			CUSIP:              offering.CUSIP,
			SettlementDate:     settlementDate,
			TradeDate:          tradeDate,
			BuyBrokerDealerID:  bdID,
			SellBrokerDealerID: domain.CBMClearingID,
			Quantity:           t.Quantity,
			Price:              t.Price,
			NetAmount:          float64(t.Quantity) * t.Price,
			SellingConcession:  sellingConc,
			GrossSpread:        grossSpread,
			AccountID:          t.Account,
			ExecInst:           t.ExecInst,
		})
	}

	// Build DTC records
	var dtcRecords []domain.DTCTrackingRecord
	for _, t := range trades {
		if t.TradeType != domain.TradeTypeBuy || t.Quantity <= 0 || t.Status == domain.TradeStatusBusted {
			continue
		}
		if t.IsDtcTracked == nil || !*t.IsDtcTracked {
			continue
		}
		bdID := t.BDFirmID
		if bdID == "" {
			bdID = t.UserID
		}
		dtcRecords = append(dtcRecords, domain.DTCTrackingRecord{
			TradeID:        t.ID,
			Symbol:         t.Symbol,
			CUSIP:          offering.CUSIP,
			BDFirmID:       bdID,
			Account:        t.Account,
			Quantity:       t.Quantity,
			Price:          t.Price,
			IsDtcTracked:   true,
			SettlementDate: settlementDate,
		})
	}

	return &domain.SettlementFileResult{
		OfferingID:          offering.ID,
		Symbol:              session.Symbol,
		AllocationSessionID: sessionID,
		UTCRecords:          utcRecords,
		DTCRecords:          dtcRecords,
		UTCCSV:              utcToCSV(utcRecords),
		DTCCSV:              dtcToCSV(dtcRecords),
		GeneratedAt:         time.Now(),
	}, nil
}

func formatDate(t time.Time) string {
	return t.UTC().Format("20060102")
}

func escapeCSV(val interface{}) string {
	s := fmt.Sprintf("%v", val)
	if strings.ContainsAny(s, ",\"\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}

func utcToCSV(records []domain.UTCTradeRecord) string {
	header := "TradeId,Symbol,CUSIP,SettlementDate,TradeDate,BuyBD,SellBD,Quantity,Price,NetAmount,SellingConcession,GrossSpread,AccountId,ExecInst"
	var rows []string
	for _, r := range records {
		rows = append(rows, fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s",
			escapeCSV(r.TradeID), escapeCSV(r.Symbol), escapeCSV(r.CUSIP),
			escapeCSV(r.SettlementDate), escapeCSV(r.TradeDate),
			escapeCSV(r.BuyBrokerDealerID), escapeCSV(r.SellBrokerDealerID),
			escapeCSV(r.Quantity), escapeCSV(r.Price), escapeCSV(r.NetAmount),
			escapeCSV(r.SellingConcession), escapeCSV(r.GrossSpread),
			escapeCSV(r.AccountID), escapeCSV(r.ExecInst),
		))
	}
	return header + "\n" + strings.Join(rows, "\n")
}

func dtcToCSV(records []domain.DTCTrackingRecord) string {
	header := "TradeId,Symbol,CUSIP,BDFirmId,Account,Quantity,Price,IsDtcTracked,SettlementDate"
	var rows []string
	for _, r := range records {
		rows = append(rows, fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s",
			escapeCSV(r.TradeID), escapeCSV(r.Symbol), escapeCSV(r.CUSIP),
			escapeCSV(r.BDFirmID), escapeCSV(r.Account),
			escapeCSV(r.Quantity), escapeCSV(r.Price),
			escapeCSV(r.IsDtcTracked), escapeCSV(r.SettlementDate),
		))
	}
	return header + "\n" + strings.Join(rows, "\n")
}
