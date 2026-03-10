# PDP Gen 2 – Scenarios

This file defines example scenarios used as external verification and regression checks against the PDP spec.[file:102][file:103]

---

## Scenario 1 – Competitive Bid with MinQty and seasoning

**Objective:** Validate MinQty constraints, seasoning clocks, and time-priority behaviour for a standard Competitive Bid.[file:102][file:103]

1. LM creates an equity IPO offering in state OPEN with:
   - Min order size for MinQty = 100,000 shares.
   - Max MinQty = 50% of order quantity.
   - Default Time Windows TW1/TW2/TW3 and trading hours enabled.[file:102][file:103]
2. A BD Broker (QSR Active = true, not on Exclude List) submits a Competitive Bid:
   - Quantity = 300,000 shares.
   - Price = 20.00.
   - MinQty = 150,000.
   - Time of entry places order in TW1.[file:102][file:103]
3. System behaviour:
   - Accepts the order and stores MinQty.
   - Sets time priority to entry timestamp.
   - Computes seasoningExpiresAt = entryTime + TW1 seasoning (24h real-time).
   - Assigns Priority Group based on TW1.[file:102][file:103]
4. During seasoning, Broker attempts to:
   - Modify quantity to 200,000 and keep MinQty 150,000, or
   - Cancel the order entirely.[file:102][file:103]
5. System rejects both actions because the order is still under seasoning.[file:102][file:103]
6. After seasoningExpiresAt:
   - Broker modifies price to 21.00, quantity to 220,000, MinQty to 110,000.
   - System validates MinQty ≤ 50% of qty and accepts.
   - Time priority is reset to modification timestamp; seasoning restarts with the active Time Window (possibly TW2).[file:102][file:103]

**Expected results:**

- All attempts to cancel/modify during seasoning are rejected with explicit seasoning error codes/messages.
- Post-seasoning modification updates time priority and Priority Group (if window changed) while maintaining MinQty within configured thresholds.[file:102][file:103]

---

## Scenario 2 – OPEN → CLOSE_PENDING with minimum size and SEC effectiveness

**Objective:** Validate CLOSE_PENDING transition preconditions: minimum offering size coverage at Minimum Acceptable IPO Price and SEC Effectiveness status.[file:102][file:103]

1. Offering is in state OPEN with:
   - Primary, Up/Down, Committed Secondary, and LM Short quantities configured.
   - Minimum Acceptable IPO Price = 18.00.
   - SEC Effectiveness required.[file:102][file:103]
2. Current orderbook demand at or above 18.00 is 90% of minimum offering size.[file:102][file:103]
3. SEC Effectiveness flag remains null.[file:102][file:103]
4. LM attempts to transition to CLOSE_PENDING.[file:102][file:103]
5. System behaviour:
   - Checks total demand at or above 18.00 versus required minimum offering size.
   - Detects demand < 100% and SEC Effectiveness null.
   - Rejects transition and records:
     - A reason code for insufficient demand at Min Acceptable IPO Price.
     - A notification event for LM/MarkOps about SEC Effectiveness missing.[file:102][file:103]
6. Later, after:
   - Demand at or above 18.00 ≥ 100% of minimum offering size, and
   - SEC Effectiveness is set to confirmed and delay has elapsed, LM retries the transition.[file:102][file:103]
7. System allows transition to CLOSE_PENDING and initializes Closing Window 1 timers.[file:102][file:103]

**Expected results:**

- CLOSE_PENDING is unreachable until both demand and SEC requirements are satisfied.
- Notifications/log entries exist for blocked transitions, including reasons and affected offering.[file:102][file:103]

---

## Scenario 3 – Preferential Bids, cutoff, and Phase 1 allocation

**Objective:** Exercise Preferential Bid rules, cancellation cutoff, and the dedicated preferential allocation phase.[file:102][file:103]

1. LM configures an offering with:
   - Preferential Bids = YES (all eligible BDs).
   - Total preferential capacity and per-investor cap in listingMinimums.
   - Defined preferential cancel/modify cutoff relative to SEC Effectiveness.[file:102][file:103]
2. Before SEC Effectiveness:
   - BD1 submits a Preferential Bid: 50,000 shares, price 17.00, ExecInst = 'p', Account populated.
   - BD2 submits a Preferential Bid: 80,000 shares, price 17.00, ExecInst = 'p', Account populated.[file:102][file:103]
3. After SEC Effectiveness Delay but before cutoff:
   - BD2 cancels its original preferential order and submits a new one at 100,000 shares, price 17.50.[file:102][file:103]
4. At the preferential cutoff time:
   - System marks all Preferential Bids as frozen with respect to cancel/modify.
   - Competitive Bids remain cancellable/modifiable subject to their own seasoning/regulatory rules.[file:102]
5. In CLEARING, Allocation Phase 1:
   - Allocates preferential capacity FCFS by time among Preferential Bids, respecting per-investor and total caps.
   - Any unfilled demand above caps is transformed into Competitive demand for subsequent allocation phases.[file:102][file:103]

**Expected results:**

- Preferential Bids cannot be cancelled/modified after cutoff and generate clear error messages when attempted.
- Preferential allocation phase produces an auditable mapping from bids to fills and from excess preferential demand to competitive pool entries.[file:102][file:103]

---

## Scenario 4 – Price collars, Indicative Clearing Price, and Cprice dissemination

**Objective:** Verify enforcement of price collars, Indicative Clearing Price computation (single-side), and Cprice publication via APIs and UIs.[file:102][file:103]

1. LM configures an equity offering with:
   - Minimum Acceptable IPO Price = 15.00.
   - Low/High Range = 16.00–20.00.
   - Collars: Minimum Price Allowed = 14.00, Maximum Price Allowed = 25.00.[file:102][file:103]
2. BDs submit Competitive Bids at prices 14.50–23.00. Another BD attempts a bid at 26.00.[file:102][file:103]
3. System behaviour:
   - Accepts bids within collars.
   - Rejects the 26.00 bid with a clear out-of-collar error.[file:102][file:103]
4. At a specific time, system computes:
   - Indicative Clearing Price for a single-side auction using Primary + Up/Down + Committed Secondary as supply, ignoring LM Short.
   - Traverses bids in price-time order high→low and respects MinQty constraints when determining the equilibrium price.[file:102][file:103]
5. Cprice is set to max(15.00, Indicative Clearing Price) and:
   - Displayed in BatMan and BidMan.
   - Exposed via `/api/cprice/:symbol`, controlled by publish mode.[file:102][file:103]
6. LM changes publish mode from “Not Published + Auto On” to “Publish NOW”; clients begin seeing actual Cprice instead of Min IPO Price.[file:102][file:103]

**Expected results:**

- Collars are enforced on entry/modify.
- Indicative Clearing Price behaves as single-side equilibrium price with MinQty.
- Cprice computation and publication mode transitions are visible across UI and API surfaces.[file:102][file:103]

---

## Scenario 5 – Priority Group Pro-rata, MinQty, and DTC tracking in settlement

**Objective:** Validate allocation integrity under Priority Group Pro-rata, MinQty behaviour, and correct DTC tracking outputs.[file:102][file:103]

1. Offering reaches CLEARING with:
   - Final Offering Price set.
   - Allocation method = Priority Group Pro-rata with configured group percentages.
   - Mix of Competitive and Preferential Bids, some with MinQty, some with ExecInst = 't' and Account populated.[file:102][file:103]
2. Allocation engine runs:
   - Phase 1: Preferential allocation per scenario 3.
   - Phase 2: For Competitive Bids at or above Offering Price:
     - Applies Priority Group ordering (Group 1 first, then 2, then 3, etc.).
     - Uses two-pass Pro-rata within each group (floor + remainder, integer arithmetic).
     - Drops any order whose allocation would fall below MinQty to zero.[file:102][file:103]
3. System validates:
   - Total allocated quantity ≤ total supply including LM Short where applicable.
   - No order with MinQty > 0 receives an allocation below MinQty.
   - Group weights and ordering match configuration.[file:102][file:103]
4. System generates:
   - NSCC UTC settlement file for all trades.
   - DTC IPO Tracking file containing only trades originating from ExecInst 't' orders, with Account preserved as passthrough identifier.[file:102][file:103]
   - Per-trade commission and selling concession values based on spread/concession parameters.[file:102][file:103]

**Expected results:**

- Allocation logs show correct Priority Group handling and MinQty enforcement.
- UTC and DTC files reflect allocated trades and DTC-tracked subset exactly.
- Commission and concession amounts reconcile with offering parameters and trade-level data.[file:102][file:103]

