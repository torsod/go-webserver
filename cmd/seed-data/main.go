package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/torsod/go-webserver/internal/config"
	"github.com/torsod/go-webserver/internal/domain"
	"github.com/torsod/go-webserver/internal/store"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()

	db, err := store.New(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.RunMigrations("migrations"); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	offeringStore := store.NewOfferingStore(db.Pool)
	orderStore := store.NewOrderStore(db.Pool)

	// ===== OFFERINGS =====
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	nextWeek := now.Add(7 * 24 * time.Hour)
	twoWeeks := now.Add(14 * 24 * time.Hour)

	offerings := []domain.Offering{
		{
			Symbol:                "LNRG",
			Name:                  "LunarEnergy Inc.",
			Issuer:                "LunarEnergy",
			Market:                "CB",
			AssetType:             domain.AssetTypeStock,
			SecurityType:          domain.SecurityTypeCS,
			State:                 domain.OfferingStateOpen,
			AnnouncementDate:      now,
			BidPeriodStartDate:    now,
			ScheduledCloseDate:    nextWeek,
			MinPriceAllowed:       10.00,
			MaxPriceAllowed:       50.00,
			MinAcceptableIPOPrice: 15.00,
			LowPriceRange:        18.00,
			HighPriceRange:        22.00,
			PriceIncrement:        0.01,
			PrimaryQuantity:       1000000,
			UpsizeQuantity:        150000,
			MinBidQuantity:        100,
			MaxBidQuantity:        500000,
			QtyIncrement:          100,
			MinAllocationPerAccount: 100,
			AllocationMethod:      domain.AllocationMethodPriceTime,
			CpricePublishMode:     domain.CpricePublishModeMinIPOPrice,
			PreferentialBidsAllowed: domain.PreferentialBidsNo,
			SecondaryOffersAllowed:  domain.SecondaryOffersNo,
			GrossUnderwritingSpread: 0.07,
			SellingConcession:       0.04,
		},
		{
			Symbol:                "SML",
			Name:                  "SmallCap Technologies",
			Issuer:                "SmallCap Tech",
			Market:                "CB",
			AssetType:             domain.AssetTypeStock,
			SecurityType:          domain.SecurityTypeCS,
			State:                 domain.OfferingStateOpen,
			AnnouncementDate:      now,
			BidPeriodStartDate:    now,
			ScheduledCloseDate:    twoWeeks,
			MinPriceAllowed:       5.00,
			MaxPriceAllowed:       30.00,
			MinAcceptableIPOPrice: 8.00,
			LowPriceRange:        12.00,
			HighPriceRange:        16.00,
			PriceIncrement:        0.01,
			PrimaryQuantity:       500000,
			UpsizeQuantity:        75000,
			MinBidQuantity:        100,
			MaxBidQuantity:        250000,
			QtyIncrement:          100,
			MinAllocationPerAccount: 100,
			AllocationMethod:      domain.AllocationMethodProRata,
			CpricePublishMode:     domain.CpricePublishModePublishNow,
			PreferentialBidsAllowed: domain.PreferentialBidsYes,
			SecondaryOffersAllowed:  domain.SecondaryOffersNo,
			GrossUnderwritingSpread: 0.065,
			SellingConcession:       0.035,
		},
		{
			Symbol:                "GRNB",
			Name:                  "GreenBond Corp 5Y Note",
			Issuer:                "GreenBond Corp",
			Market:                "CB",
			AssetType:             domain.AssetTypeBond,
			SecurityType:          domain.SecurityTypeCORP,
			YieldType:             domain.YieldTypeYieldToMaturity,
			State:                 domain.OfferingStateUpcoming,
			AnnouncementDate:      now,
			BidPeriodStartDate:    tomorrow,
			ScheduledCloseDate:    twoWeeks,
			MinPriceAllowed:       95.00,
			MaxPriceAllowed:       105.00,
			MinAcceptableIPOPrice: 98.00,
			LowPriceRange:        99.00,
			HighPriceRange:        101.00,
			PriceIncrement:        0.125,
			PrimaryQuantity:       200000,
			MinBidQuantity:        10,
			MaxBidQuantity:        50000,
			QtyIncrement:          10,
			MinAllocationPerAccount: 10,
			AllocationMethod:      domain.AllocationMethodPriorityGroupProRata,
			CpricePublishMode:     domain.CpricePublishModeNotPublishedAutoOn,
			PreferentialBidsAllowed: domain.PreferentialBidsLeadManagerBDOnly,
			SecondaryOffersAllowed:  domain.SecondaryOffersNo,
			GrossUnderwritingSpread: 0.02,
			SellingConcession:       0.01,
		},
	}

	offeringIDs := map[string]string{}
	for _, o := range offerings {
		// Check if symbol exists
		existing, _ := offeringStore.FindBySymbol(ctx, o.Symbol)
		if existing != nil {
			fmt.Printf("Offering %s already exists (id: %s), skipping\n", o.Symbol, existing.ID)
			offeringIDs[o.Symbol] = existing.ID
			continue
		}

		id, err := offeringStore.Insert(ctx, &o)
		if err != nil {
			slog.Error("failed to create offering", "symbol", o.Symbol, "error", err)
			continue
		}
		offeringIDs[o.Symbol] = id
		fmt.Printf("Created offering %s - %s (id: %s, state: %s)\n", o.Symbol, o.Name, id, o.State)
	}

	// ===== ORDERS for LNRG (OPEN offering) =====
	fmt.Println("\n--- Seeding orders for LNRG ---")

	users := []struct {
		userID   string
		bdFirmID string
	}{
		{"bd1", "BD_FIRM_1"},
		{"bd2", "BD_FIRM_2"},
		{"bd3", "BD_FIRM_3"},
	}

	// Generate a spread of bid orders
	bidOrders := []struct {
		price    float64
		qty      int64
		priority int
		userIdx  int
	}{
		{22.00, 5000, 1, 0},
		{21.50, 10000, 1, 0},
		{21.00, 8000, 1, 1},
		{21.00, 12000, 2, 1},
		{20.75, 3000, 1, 2},
		{20.50, 15000, 1, 0},
		{20.50, 7000, 2, 2},
		{20.00, 20000, 1, 1},
		{20.00, 10000, 1, 2},
		{19.50, 25000, 2, 0},
		{19.00, 8000, 1, 1},
		{18.50, 30000, 3, 2},
		{18.00, 5000, 1, 0},
		{18.00, 12000, 2, 1},
		{17.50, 10000, 1, 2},
	}

	for _, bo := range bidOrders {
		u := users[bo.userIdx]
		seq, _ := orderStore.NextSequence(ctx)
		order := &domain.Order{
			Symbol:            "LNRG",
			Side:              domain.OrderSideBid,
			OrderType:         domain.OrderTypeCompetitive,
			Quantity:          bo.qty,
			Price:             bo.price,
			PriorityGroup:     bo.priority,
			UserID:            u.userID,
			BDFirmID:          u.bdFirmID,
			Status:            domain.OrderStatusActive,
			Timestamp:         now.Add(time.Duration(rand.Intn(3600)) * time.Second),
			OriginalEntryTime: now,
			OrderSequence:     seq,
		}

		id, err := orderStore.Insert(ctx, order)
		if err != nil {
			slog.Error("failed to create order", "error", err)
			continue
		}
		fmt.Printf("  BID %s: %d @ $%.2f (PG=%d, user=%s, id=%s)\n",
			order.Symbol, order.Quantity, order.Price, order.PriorityGroup, order.UserID, id)
	}

	// Generate a few offer orders
	offerOrders := []struct {
		price    float64
		qty      int64
		priority int
		userIdx  int
	}{
		{22.50, 3000, 1, 0},
		{23.00, 5000, 1, 1},
		{23.50, 8000, 1, 2},
		{24.00, 10000, 2, 0},
	}

	for _, oo := range offerOrders {
		u := users[oo.userIdx]
		seq, _ := orderStore.NextSequence(ctx)
		order := &domain.Order{
			Symbol:            "LNRG",
			Side:              domain.OrderSideOffer,
			OrderType:         domain.OrderTypeCompetitive,
			Quantity:          oo.qty,
			Price:             oo.price,
			PriorityGroup:     oo.priority,
			UserID:            u.userID,
			BDFirmID:          u.bdFirmID,
			Status:            domain.OrderStatusActive,
			Timestamp:         now.Add(time.Duration(rand.Intn(3600)) * time.Second),
			OriginalEntryTime: now,
			OrderSequence:     seq,
		}

		id, err := orderStore.Insert(ctx, order)
		if err != nil {
			slog.Error("failed to create order", "error", err)
			continue
		}
		fmt.Printf("  OFFER %s: %d @ $%.2f (PG=%d, user=%s, id=%s)\n",
			order.Symbol, order.Quantity, order.Price, order.PriorityGroup, order.UserID, id)
	}

	// ===== ORDERS for SML (OPEN offering) =====
	fmt.Println("\n--- Seeding orders for SML ---")

	smlBids := []struct {
		price    float64
		qty      int64
		priority int
		userIdx  int
	}{
		{16.00, 3000, 1, 0},
		{15.50, 5000, 1, 1},
		{15.00, 8000, 1, 2},
		{14.50, 10000, 2, 0},
		{14.00, 12000, 1, 1},
		{13.50, 7000, 1, 2},
		{13.00, 15000, 2, 0},
		{12.50, 20000, 3, 1},
		{12.00, 10000, 1, 2},
	}

	for _, bo := range smlBids {
		u := users[bo.userIdx]
		seq, _ := orderStore.NextSequence(ctx)
		order := &domain.Order{
			Symbol:            "SML",
			Side:              domain.OrderSideBid,
			OrderType:         domain.OrderTypeCompetitive,
			Quantity:          bo.qty,
			Price:             bo.price,
			PriorityGroup:     bo.priority,
			UserID:            u.userID,
			BDFirmID:          u.bdFirmID,
			Status:            domain.OrderStatusActive,
			Timestamp:         now.Add(time.Duration(rand.Intn(3600)) * time.Second),
			OriginalEntryTime: now,
			OrderSequence:     seq,
		}

		id, err := orderStore.Insert(ctx, order)
		if err != nil {
			slog.Error("failed to create order", "error", err)
			continue
		}
		fmt.Printf("  BID %s: %d @ $%.2f (PG=%d, user=%s, id=%s)\n",
			order.Symbol, order.Quantity, order.Price, order.PriorityGroup, order.UserID, id)
	}

	fmt.Printf("\nSeed data complete!\n")
	fmt.Printf("  Offerings: %d\n", len(offerings))
	fmt.Printf("  LNRG orders: %d bids + %d offers\n", len(bidOrders), len(offerOrders))
	fmt.Printf("  SML orders: %d bids\n", len(smlBids))
	fmt.Println("\nView at: http://localhost:3000/manage.html")
}
