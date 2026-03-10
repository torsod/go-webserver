package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/torsod/go-webserver/internal/domain"
)

type offeringStore struct {
	pool *pgxpool.Pool
}

// NewOfferingStore creates a new PostgreSQL-backed offering store
func NewOfferingStore(pool *pgxpool.Pool) OfferingStore {
	return &offeringStore{pool: pool}
}

func (s *offeringStore) FindAll(ctx context.Context) ([]*domain.Offering, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+offeringColumns+` FROM offerings ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query offerings: %w", err)
	}
	defer rows.Close()

	var offerings []*domain.Offering
	for rows.Next() {
		o, err := scanOffering(rows)
		if err != nil {
			return nil, err
		}
		offerings = append(offerings, o)
	}
	return offerings, nil
}

func (s *offeringStore) FindByID(ctx context.Context, id string) (*domain.Offering, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+offeringColumns+` FROM offerings WHERE id = $1`, id)
	return scanOfferingRow(row)
}

func (s *offeringStore) FindBySymbol(ctx context.Context, symbol string) (*domain.Offering, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+offeringColumns+` FROM offerings WHERE symbol = $1`, symbol)
	return scanOfferingRow(row)
}

func (s *offeringStore) FindByStates(ctx context.Context, states []domain.OfferingState) ([]*domain.Offering, error) {
	stateStrs := make([]string, len(states))
	for i, st := range states {
		stateStrs[i] = string(st)
	}
	rows, err := s.pool.Query(ctx, `SELECT `+offeringColumns+` FROM offerings WHERE state = ANY($1)`, stateStrs)
	if err != nil {
		return nil, fmt.Errorf("query offerings by state: %w", err)
	}
	defer rows.Close()

	var offerings []*domain.Offering
	for rows.Next() {
		o, err := scanOffering(rows)
		if err != nil {
			return nil, err
		}
		offerings = append(offerings, o)
	}
	return offerings, nil
}

func (s *offeringStore) Insert(ctx context.Context, offering *domain.Offering) (string, error) {
	twJSON, _ := json.Marshal(offering.TimeWindows)
	cwcJSON, _ := json.Marshal(offering.ClosingWindowsConfig)
	cwdJSON, _ := json.Marshal(offering.ClosingWindowsData)
	lmJSON, _ := json.Marshal(offering.ListingMinimums)
	ebdJSON, _ := json.Marshal(offering.ExcludedBrokerDealers)
	clJSON, _ := json.Marshal(offering.ChangeLog)

	var id string
	err := s.pool.QueryRow(ctx, `
		INSERT INTO offerings (
			symbol, name, issuer, market, cusip, asset_type, security_type,
			yield_type, face_value, maturity_date, offering_coupon, nav, distribution_fee,
			state, previous_state, announcement_date, bid_period_start_date, scheduled_close_date, settlement_date,
			sec_effectiveness_date, sec_effectiveness_delay_minutes,
			max_price_allowed, min_price_allowed, high_price_range, low_price_range, min_acceptable_ipo_price,
			dividend, primary_quantity, upsize_quantity, committed_secondary_quantity, lm_short_quantity,
			min_bid_quantity, max_bid_quantity, qty_increment, price_increment,
			min_order_size_for_min_qty, max_min_qty_percentage, min_qty_deadline,
			allocation_method, min_allocation_per_account,
			preferential_bids_allowed, secondary_offers_allowed,
			cprice_publish_mode, cprice_publishing_active,
			gross_underwriting_spread, selling_concession, prio_group_test,
			time_windows, closing_windows_config, closing_windows_data, listing_minimums,
			excluded_broker_dealers, lead_manager_bd_id,
			clearing_price, offering_price, bond_offering_price,
			indicative_clearing_price, cprice, total_demand, total_orders,
			change_log, created_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
			$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,
			$39,$40,$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,$51,$52,$53,$54,$55,$56,$57,$58,$59,
			$60,$61,$62
		) RETURNING id`,
		offering.Symbol, offering.Name, offering.Issuer, offering.Market, offering.CUSIP,
		offering.AssetType, offering.SecurityType,
		offering.YieldType, offering.FaceValue, offering.MaturityDate, offering.OfferingCoupon,
		offering.NAV, offering.DistributionFee,
		offering.State, offering.PreviousState,
		offering.AnnouncementDate, offering.BidPeriodStartDate, offering.ScheduledCloseDate, offering.SettlementDate,
		offering.SECEffectivenessDate, offering.SECEffectivenessDelayMinutes,
		offering.MaxPriceAllowed, offering.MinPriceAllowed, offering.HighPriceRange, offering.LowPriceRange,
		offering.MinAcceptableIPOPrice,
		offering.Dividend, offering.PrimaryQuantity, offering.UpsizeQuantity,
		offering.CommittedSecondaryQuantity, offering.LMShortQuantity,
		offering.MinBidQuantity, offering.MaxBidQuantity, offering.QtyIncrement, offering.PriceIncrement,
		offering.MinOrderSizeForMinQty, offering.MaxMinQtyPercentage, offering.MinQtyDeadline,
		offering.AllocationMethod, offering.MinAllocationPerAccount,
		offering.PreferentialBidsAllowed, offering.SecondaryOffersAllowed,
		offering.CpricePublishMode, offering.CpricePublishingActive,
		offering.GrossUnderwritingSpread, offering.SellingConcession, offering.PrioGroupTest,
		twJSON, cwcJSON, cwdJSON, lmJSON, ebdJSON, offering.LeadManagerBdID,
		offering.ClearingPrice, offering.OfferingPrice, offering.BondOfferingPrice,
		offering.IndicativeClearingPrice, offering.Cprice, offering.TotalDemand, offering.TotalOrders,
		clJSON, offering.CreatedAt,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert offering: %w", err)
	}
	return id, nil
}

func (s *offeringStore) Update(ctx context.Context, id string, offering *domain.Offering) error {
	twJSON, _ := json.Marshal(offering.TimeWindows)
	cwcJSON, _ := json.Marshal(offering.ClosingWindowsConfig)
	cwdJSON, _ := json.Marshal(offering.ClosingWindowsData)
	lmJSON, _ := json.Marshal(offering.ListingMinimums)
	ebdJSON, _ := json.Marshal(offering.ExcludedBrokerDealers)
	clJSON, _ := json.Marshal(offering.ChangeLog)

	_, err := s.pool.Exec(ctx, `
		UPDATE offerings SET
			symbol=$1, name=$2, issuer=$3, market=$4, cusip=$5, asset_type=$6, security_type=$7,
			state=$8, previous_state=$9,
			announcement_date=$10, bid_period_start_date=$11, scheduled_close_date=$12, settlement_date=$13,
			max_price_allowed=$14, min_price_allowed=$15, high_price_range=$16, low_price_range=$17,
			min_acceptable_ipo_price=$18,
			primary_quantity=$19, upsize_quantity=$20, committed_secondary_quantity=$21, lm_short_quantity=$22,
			min_bid_quantity=$23, max_bid_quantity=$24, qty_increment=$25, price_increment=$26,
			allocation_method=$27, preferential_bids_allowed=$28, secondary_offers_allowed=$29,
			cprice_publish_mode=$30, cprice_publishing_active=$31,
			gross_underwriting_spread=$32, selling_concession=$33, prio_group_test=$34,
			time_windows=$35, closing_windows_config=$36, closing_windows_data=$37, listing_minimums=$38,
			excluded_broker_dealers=$39, lead_manager_bd_id=$40,
			clearing_price=$41, offering_price=$42, bond_offering_price=$43,
			cprice=$44, change_log=$45, updated_at=NOW()
		WHERE id=$46`,
		offering.Symbol, offering.Name, offering.Issuer, offering.Market, offering.CUSIP,
		offering.AssetType, offering.SecurityType,
		offering.State, offering.PreviousState,
		offering.AnnouncementDate, offering.BidPeriodStartDate, offering.ScheduledCloseDate, offering.SettlementDate,
		offering.MaxPriceAllowed, offering.MinPriceAllowed, offering.HighPriceRange, offering.LowPriceRange,
		offering.MinAcceptableIPOPrice,
		offering.PrimaryQuantity, offering.UpsizeQuantity, offering.CommittedSecondaryQuantity, offering.LMShortQuantity,
		offering.MinBidQuantity, offering.MaxBidQuantity, offering.QtyIncrement, offering.PriceIncrement,
		offering.AllocationMethod, offering.PreferentialBidsAllowed, offering.SecondaryOffersAllowed,
		offering.CpricePublishMode, offering.CpricePublishingActive,
		offering.GrossUnderwritingSpread, offering.SellingConcession, offering.PrioGroupTest,
		twJSON, cwcJSON, cwdJSON, lmJSON, ebdJSON, offering.LeadManagerBdID,
		offering.ClearingPrice, offering.OfferingPrice, offering.BondOfferingPrice,
		offering.Cprice, clJSON,
		id,
	)
	if err != nil {
		return fmt.Errorf("update offering: %w", err)
	}
	return nil
}

func (s *offeringStore) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return updateFields(ctx, s.pool, "offerings", id, fields)
}

// Column list for offerings table
const offeringColumns = `id, symbol, name, issuer, market, cusip, asset_type, security_type,
	yield_type, face_value, maturity_date, offering_coupon, nav, distribution_fee,
	state, previous_state, announcement_date, bid_period_start_date, scheduled_close_date, settlement_date,
	sec_effectiveness_date, sec_effectiveness_delay_minutes,
	max_price_allowed, min_price_allowed, high_price_range, low_price_range, min_acceptable_ipo_price,
	dividend, primary_quantity, upsize_quantity, committed_secondary_quantity, lm_short_quantity,
	min_bid_quantity, max_bid_quantity, qty_increment, price_increment,
	min_order_size_for_min_qty, max_min_qty_percentage, min_qty_deadline,
	allocation_method, min_allocation_per_account,
	preferential_bids_allowed, secondary_offers_allowed,
	cprice_publish_mode, cprice_publishing_active,
	gross_underwriting_spread, selling_concession, prio_group_test,
	time_windows, closing_windows_config, closing_windows_data, listing_minimums,
	excluded_broker_dealers, lead_manager_bd_id,
	clearing_price, offering_price, bond_offering_price,
	indicative_clearing_price, cprice, total_demand, total_orders,
	change_log, created_at, updated_at`

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanOfferingRow(row scannable) (*domain.Offering, error) {
	o := &domain.Offering{}
	var twJSON, cwcJSON, cwdJSON, lmJSON, ebdJSON, clJSON []byte

	err := row.Scan(
		&o.ID, &o.Symbol, &o.Name, &o.Issuer, &o.Market, &o.CUSIP,
		&o.AssetType, &o.SecurityType,
		&o.YieldType, &o.FaceValue, &o.MaturityDate, &o.OfferingCoupon, &o.NAV, &o.DistributionFee,
		&o.State, &o.PreviousState,
		&o.AnnouncementDate, &o.BidPeriodStartDate, &o.ScheduledCloseDate, &o.SettlementDate,
		&o.SECEffectivenessDate, &o.SECEffectivenessDelayMinutes,
		&o.MaxPriceAllowed, &o.MinPriceAllowed, &o.HighPriceRange, &o.LowPriceRange, &o.MinAcceptableIPOPrice,
		&o.Dividend, &o.PrimaryQuantity, &o.UpsizeQuantity, &o.CommittedSecondaryQuantity, &o.LMShortQuantity,
		&o.MinBidQuantity, &o.MaxBidQuantity, &o.QtyIncrement, &o.PriceIncrement,
		&o.MinOrderSizeForMinQty, &o.MaxMinQtyPercentage, &o.MinQtyDeadline,
		&o.AllocationMethod, &o.MinAllocationPerAccount,
		&o.PreferentialBidsAllowed, &o.SecondaryOffersAllowed,
		&o.CpricePublishMode, &o.CpricePublishingActive,
		&o.GrossUnderwritingSpread, &o.SellingConcession, &o.PrioGroupTest,
		&twJSON, &cwcJSON, &cwdJSON, &lmJSON, &ebdJSON, &o.LeadManagerBdID,
		&o.ClearingPrice, &o.OfferingPrice, &o.BondOfferingPrice,
		&o.IndicativeClearingPrice, &o.Cprice, &o.TotalDemand, &o.TotalOrders,
		&clJSON, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan offering: %w", err)
	}

	if twJSON != nil {
		json.Unmarshal(twJSON, &o.TimeWindows)
	}
	if cwcJSON != nil {
		json.Unmarshal(cwcJSON, &o.ClosingWindowsConfig)
	}
	if cwdJSON != nil {
		json.Unmarshal(cwdJSON, &o.ClosingWindowsData)
	}
	if lmJSON != nil {
		json.Unmarshal(lmJSON, &o.ListingMinimums)
	}
	if ebdJSON != nil {
		json.Unmarshal(ebdJSON, &o.ExcludedBrokerDealers)
	}
	if clJSON != nil {
		json.Unmarshal(clJSON, &o.ChangeLog)
	}

	return o, nil
}

func scanOffering(rows scannable) (*domain.Offering, error) {
	return scanOfferingRow(rows)
}
