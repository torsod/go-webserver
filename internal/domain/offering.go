package domain

import "time"

// OfferingState represents the offering lifecycle states per CB-PDP-Gen-2 spec Section 4.4
type OfferingState string

const (
	OfferingStateNew          OfferingState = "NEW"
	OfferingStateUpcoming     OfferingState = "UPCOMING"
	OfferingStateOpen         OfferingState = "OPEN"
	OfferingStateClosePending OfferingState = "CLOSE_PENDING"
	OfferingStateClosing      OfferingState = "CLOSING"
	OfferingStateClearing     OfferingState = "CLEARING"
	OfferingStateClosed       OfferingState = "CLOSED"
	OfferingStateHalted       OfferingState = "HALTED"
	OfferingStateFrozen       OfferingState = "FROZEN"
	OfferingStateCanceled     OfferingState = "CANCELED"
)

// AssetType classification
type AssetType string

const (
	AssetTypeStock   AssetType = "STOCK"
	AssetTypeBond    AssetType = "BOND"
	AssetTypeFundETF AssetType = "FUND_ETF"
)

// SecurityType per spec
type SecurityType string

const (
	SecurityTypeCS   SecurityType = "CS"   // Common Stock
	SecurityTypePS   SecurityType = "PS"   // Preferred Stock
	SecurityTypeCORP SecurityType = "CORP" // Corporate Bond
	SecurityTypeMUNI SecurityType = "MUNI" // Municipal Bond
	SecurityTypeMF   SecurityType = "MF"   // Mutual Fund
	SecurityTypeETF  SecurityType = "ETF"  // Exchange Traded Fund
)

// YieldType for bonds
type YieldType string

const (
	YieldTypeYieldToMaturity   YieldType = "YIELD_TO_MATURITY"
	YieldTypeYieldToCall       YieldType = "YIELD_TO_CALL"
	YieldTypeSpreadToTreasury  YieldType = "SPREAD_TO_TREASURY"
	YieldTypeCurrentYield      YieldType = "CURRENT_YIELD"
)

// CpricePublishMode per spec Section 4.1.1
type CpricePublishMode string

const (
	CpricePublishModeMinIPOPrice          CpricePublishMode = "MIN_IPO_PRICE"
	CpricePublishModeNotPublishedAutoOn   CpricePublishMode = "NOT_PUBLISHED_AUTO_ON"
	CpricePublishModeNotPublishedManualOn CpricePublishMode = "NOT_PUBLISHED_MANUAL_ON"
	CpricePublishModePublishNow           CpricePublishMode = "PUBLISH_NOW"
)

// AllocationMethod per spec Section 7
type AllocationMethod string

const (
	AllocationMethodPriceTime             AllocationMethod = "PRICE_TIME"
	AllocationMethodProRata               AllocationMethod = "PRO_RATA"
	AllocationMethodPriorityGroupProRata  AllocationMethod = "PRIORITY_GROUP_PRO_RATA"
)

// PreferentialBidsAllowed per spec Section 3.2.2
type PreferentialBidsAllowed string

const (
	PreferentialBidsNo               PreferentialBidsAllowed = "NO"
	PreferentialBidsLeadManagerBDOnly PreferentialBidsAllowed = "LEAD_MANAGER_BD_ONLY"
	PreferentialBidsYes              PreferentialBidsAllowed = "YES"
)

// SecondaryOffersAllowed per spec Section 3.3.4
type SecondaryOffersAllowed string

const (
	SecondaryOffersNo               SecondaryOffersAllowed = "NO"
	SecondaryOffersLeadManagerBDOnly SecondaryOffersAllowed = "LEAD_MANAGER_BD_ONLY"
	SecondaryOffersYes              SecondaryOffersAllowed = "YES"
)

// TimeWindow per spec Section 3.9
type TimeWindow struct {
	WindowNumber          int       `json:"windowNumber" db:"window_number"`
	StartTime             time.Time `json:"startTime" db:"start_time"`
	SeasoningPeriodMinutes int      `json:"seasoningPeriodMinutes" db:"seasoning_period_minutes"`
	PriorityGroup         int       `json:"priorityGroup" db:"priority_group"`
}

// ClosingWindowsConfig per offering (spec Section 5.3)
type ClosingWindowsConfig struct {
	CW1DurationMinutes   int     `json:"cw1DurationMinutes" db:"cw1_duration_minutes"`
	CW2DurationMinutes   int     `json:"cw2DurationMinutes" db:"cw2_duration_minutes"`
	QuietDurationSeconds int     `json:"quietDurationSeconds" db:"quiet_duration_seconds"`
	VolumeChangeThreshold float64 `json:"volumeChangeThreshold" db:"volume_change_threshold"`
	StartTime            string  `json:"startTime" db:"start_time"` // HH:mm Eastern
}

// ClosingWindowsData runtime state
type ClosingWindowsData struct {
	CW1StartTime      *time.Time `json:"cw1StartTime"`
	CW2StartTime      *time.Time `json:"cw2StartTime"`
	CW2RandomOffset   float64    `json:"cw2RandomOffset"`
	LastClearingPrice  *float64   `json:"lastClearingPrice"`
	LastVolume         *float64   `json:"lastVolume"`
	QuietPeriodStart   *time.Time `json:"quietPeriodStart"`
	ForceCloseTime     *time.Time `json:"forceCloseTime"`
}

// ListingMinimums per spec Section 5.1
type ListingMinimums struct {
	MinMarketValue                    float64 `json:"minMarketValue"`
	MinPublicShares                   int64   `json:"minPublicShares"`
	MinPricePerShare                  float64 `json:"minPricePerShare"`
	MinHolders                        int     `json:"minHolders"`
	MinRoundLotHolders                int     `json:"minRoundLotHolders"`
	MaxPreferentialAllocationPerBidder float64 `json:"maxPreferentialAllocationPerBidder"`
	MaxTotalPreferentialAllocation     float64 `json:"maxTotalPreferentialAllocation"`
}

// OfferingChangeLog entry
type OfferingChangeLog struct {
	Field     string      `json:"field"`
	OldValue  interface{} `json:"oldValue"`
	NewValue  interface{} `json:"newValue"`
	ChangedBy string      `json:"changedBy"`
	ChangedAt time.Time   `json:"changedAt"`
}

// Offering represents an IPO/offering per CB-PDP-Gen-2 spec
type Offering struct {
	ID     string `json:"id" db:"id"`

	// Core identifiers
	Symbol string `json:"symbol" db:"symbol"`
	Name   string `json:"name" db:"name"`
	Issuer string `json:"issuer" db:"issuer"`
	Market string `json:"market" db:"market"`
	CUSIP  string `json:"cusip,omitempty" db:"cusip"`

	// Asset classification
	AssetType    AssetType    `json:"assetType" db:"asset_type"`
	SecurityType SecurityType `json:"securityType,omitempty" db:"security_type"`

	// Bond-specific
	YieldType      YieldType  `json:"yieldType,omitempty" db:"yield_type"`
	FaceValue      *float64   `json:"faceValue,omitempty" db:"face_value"`
	MaturityDate   *time.Time `json:"maturityDate,omitempty" db:"maturity_date"`
	OfferingCoupon *float64   `json:"offeringCoupon,omitempty" db:"offering_coupon"`

	// Fund/ETF
	NAV             *float64 `json:"nav,omitempty" db:"nav"`
	DistributionFee *float64 `json:"distributionFee,omitempty" db:"distribution_fee"`

	// State
	State         OfferingState `json:"state" db:"state"`
	PreviousState OfferingState `json:"previousState,omitempty" db:"previous_state"`

	// Dates
	AnnouncementDate   time.Time  `json:"announcementDate" db:"announcement_date"`
	BidPeriodStartDate time.Time  `json:"bidPeriodStartDate" db:"bid_period_start_date"`
	ScheduledCloseDate time.Time  `json:"scheduledCloseDate" db:"scheduled_close_date"`
	SettlementDate     *time.Time `json:"settlementDate,omitempty" db:"settlement_date"`

	// SEC Effectiveness
	SECEffectivenessDate         *time.Time `json:"secEffectivenessDate,omitempty" db:"sec_effectiveness_date"`
	SECEffectivenessDelayMinutes int        `json:"secEffectivenessDelayMinutes" db:"sec_effectiveness_delay_minutes"`

	// Price Collars (spec Section 3.6)
	MaxPriceAllowed float64 `json:"maxPriceAllowed" db:"max_price_allowed"`
	MinPriceAllowed float64 `json:"minPriceAllowed" db:"min_price_allowed"`

	// Price Range (spec Section 3.7)
	HighPriceRange float64 `json:"highPriceRange" db:"high_price_range"`
	LowPriceRange  float64 `json:"lowPriceRange" db:"low_price_range"`

	// Minimum Acceptable IPO Price (spec Section 3.8)
	MinAcceptableIPOPrice float64 `json:"minAcceptableIpoPrice" db:"min_acceptable_ipo_price"`

	// Dividend
	Dividend *float64 `json:"dividend,omitempty" db:"dividend"`

	// Offer Quantities (spec Sections 3.3.1-3.3.5)
	PrimaryQuantity            int64 `json:"primaryQuantity" db:"primary_quantity"`
	UpsizeQuantity             int64 `json:"upsizeQuantity" db:"upsize_quantity"`
	CommittedSecondaryQuantity int64 `json:"committedSecondaryQuantity" db:"committed_secondary_quantity"`
	LMShortQuantity            int64 `json:"lmShortQuantity" db:"lm_short_quantity"`

	// Order Constraints (spec Section 3.2.3)
	MinBidQuantity int64   `json:"minBidQuantity" db:"min_bid_quantity"`
	MaxBidQuantity int64   `json:"maxBidQuantity" db:"max_bid_quantity"`
	QtyIncrement   int64   `json:"qtyIncrement" db:"qty_increment"`
	PriceIncrement float64 `json:"priceIncrement" db:"price_increment"`

	// MinQty Configuration (spec Section 3.2.1.1)
	MinOrderSizeForMinQty int64    `json:"minOrderSizeForMinQty" db:"min_order_size_for_min_qty"`
	MaxMinQtyPercentage   float64  `json:"maxMinQtyPercentage" db:"max_min_qty_percentage"`
	MinQtyDeadline        *time.Time `json:"minQtyDeadline,omitempty" db:"min_qty_deadline"`

	// Allocation (spec Section 7)
	AllocationMethod        AllocationMethod `json:"allocationMethod" db:"allocation_method"`
	MinAllocationPerAccount int64            `json:"minAllocationPerAccount" db:"min_allocation_per_account"`

	// Preferential Bids (spec Section 3.2.2)
	PreferentialBidsAllowed PreferentialBidsAllowed `json:"preferentialBidsAllowed" db:"preferential_bids_allowed"`

	// Secondary Sell Offers (spec Section 3.3.4)
	SecondaryOffersAllowed SecondaryOffersAllowed `json:"secondaryOffersAllowed" db:"secondary_offers_allowed"`

	// Cprice publishing (spec Section 4.1.1)
	CpricePublishMode     CpricePublishMode `json:"cpricePublishMode" db:"cprice_publish_mode"`
	CpricePublishingActive bool             `json:"cpricePublishingActive" db:"cprice_publishing_active"`

	// Fees
	GrossUnderwritingSpread float64 `json:"grossUnderwritingSpread" db:"gross_underwriting_spread"`
	SellingConcession       float64 `json:"sellingConcession" db:"selling_concession"`

	// Priority Group Test Mode
	PrioGroupTest bool `json:"prioGroupTest" db:"prio_group_test"`

	// Time Windows / Seasoning Periods (spec Section 3.9, 3.10) - stored as JSONB
	TimeWindows []TimeWindow `json:"timeWindows" db:"time_windows"`

	// Closing Windows (spec Section 5.3) - stored as JSONB
	ClosingWindowsConfig ClosingWindowsConfig `json:"closingWindowsConfig" db:"closing_windows_config"`
	ClosingWindowsData   *ClosingWindowsData  `json:"closingWindowsData,omitempty" db:"closing_windows_data"`

	// Listing Exchange Minimums (spec Section 5.1) - stored as JSONB
	ListingMinimums *ListingMinimums `json:"listingMinimums,omitempty" db:"listing_minimums"`

	// Broker-Dealer Exclude List (spec Section 9.3.1) - stored as JSONB
	ExcludedBrokerDealers []string `json:"excludedBrokerDealers" db:"excluded_broker_dealers"`

	// Lead Manager BD
	LeadManagerBdID string `json:"leadManagerBdId,omitempty" db:"lead_manager_bd_id"`

	// Clearing/Offering Price (set after closing)
	ClearingPrice    *float64 `json:"clearingPrice,omitempty" db:"clearing_price"`
	OfferingPrice    *float64 `json:"offeringPrice,omitempty" db:"offering_price"`
	BondOfferingPrice *float64 `json:"bondOfferingPrice,omitempty" db:"bond_offering_price"`

	// Computed metrics
	IndicativeClearingPrice *float64 `json:"indicativeClearingPrice,omitempty" db:"indicative_clearing_price"`
	Cprice                  *float64 `json:"cprice,omitempty" db:"cprice"`
	TotalDemand             *int64   `json:"totalDemand,omitempty" db:"total_demand"`
	TotalOrders             *int64   `json:"totalOrders,omitempty" db:"total_orders"`

	// Change Log - stored as JSONB
	ChangeLog []OfferingChangeLog `json:"changeLog" db:"change_log"`

	// Timestamps
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

// GetOfferingSize returns total offering size for clearing price calculation (excludes LM Short)
func GetOfferingSize(o *Offering) int64 {
	return o.PrimaryQuantity + o.UpsizeQuantity + o.CommittedSecondaryQuantity
}

// GetAllocationSize returns total allocation size including LM Short
func GetAllocationSize(o *Offering) int64 {
	return GetOfferingSize(o) + o.LMShortQuantity
}
