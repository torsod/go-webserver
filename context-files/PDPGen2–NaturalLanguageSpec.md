# PDP Gen 2 – Natural Language Spec

## 1. System overview

PDP Gen 2 is an auction-based **Price Discovery Platform** for new issues (IPOs and follow-ons), providing price and demand discovery, allocation, and settlement outputs for multiple asset types (stocks, bonds, funds/ETFs).[file:102] The platform exposes capabilities to Lead Managers (LM), Broker-Dealers (BD), Market Operations (MarkOps), and Issuers via apps (BatMan, BidMan, MarkOps), FIX, and REST APIs.[file:102][file:103]

Core responsibilities:

- Accept, validate, and manage GTC limit orders (bids and offers where enabled) during configured bidding periods, enforcing eligibility, Min/Max quantity, MinQty, and DTC tracking rules.[file:102][file:103]
- Maintain an offering lifecycle state machine and enforce state-dependent permissions for order actions and configuration changes.[file:102][file:103]
- Compute Indicative Clearing Price and Cprice/Cyield under single-side (now) and two-sided (later) auction modes, respecting supply composition, MinQty, and price controls.[file:102][file:103]
- Run allocation algorithms (Price-time, Pro-rata, Priority Group Pro-rata) and produce auditable allocation and settlement outputs.[file:102][file:103]
- Integrate with NSCC/DTC/FED NSS flows via UTC and DTC IPO Tracking files, and with BD OMSs via FIX for orders and executions.[file:102][file:103]

---

## 2. Roles and access

The system defines these primary roles.[file:102][file:103]

- **Lead Manager (LM)**  
  - Creates and configures offerings (asset attributes, supply, price ranges/collars, Min/Max quantities, allocation methods).[file:102][file:103]
  - Monitors orderbook, runs allocation simulations, sets Offering Price, and confirms allocations.[file:102][file:103]

- **Broker-Dealer (BD)**  
  - Submits and manages orders via FIX and BidMan.[file:102][file:103]
  - Sub-roles: Admin (no trading), Broker, OMS User (machine identity), Group Manager, Master Broker, with role-based constraints on order entry, view, and cancel/write.[file:102][file:103]

- **Market Operations (MarkOps)**  
  - Operates state transitions within allowed rules, halts/freezes offerings, configures system-wide parameters (calendar, market hours, defaults), and monitors system health.[file:102][file:103]

- **Issuer / Issuer-facing LM**  
  - Accesses read-only dashboards showing offering status, indicative clearing statistics, and allocation summaries.[file:102][file:103]

All mutating actions are subject to method-level RBAC checks, including orders, offering updates, and settlement generation.[file:103] QSR Active and BD Exclude List are enforced for BD eligibility at order entry.[file:102][file:103]

---

## 3. Offerings and lifecycle

### 3.1 Offer entities and supply

Each offering models a single new-issue event with:[file:102][file:103]

- Asset and security type (stock, bond, fund/ETF, CS, PS, CORP, MUNI, MF, ETF).  
- Bid period times and state.[file:102][file:103]
- Supply broken into offer groups:
  1. Primary Shares/Bonds  
  2. Up-/Down-size of Primary  
  3. Committed Secondary Shares/Bonds via LM  
  4. Limit-priced Secondary Shares Offers (BD offers, if enabled)  
  5. Lead Manager Short (LM Short)[file:102][file:103]

Rules:

- Primary/Up/Down/Committed Secondary and LM Short are maintained by LM and logged with timestamp on each change.[file:102][file:103]
- LM Short is excluded from Indicative/Clearing Price but included in allocation supply.[file:102][file:103]
- Limit-priced Secondary Offers are GTC limit offers, subject to collars and seasoning, cannot carry DTC Tracking flags, and are allowed only when offering parameters enable them.[file:102][file:103]

### 3.2 States and transitions

Offering states:[file:102][file:103]

- NEW: Pre-announcement; no order entry or cancel.  
- UPCOMING: Post-announcement; existing orders may be cancelled; no new entry/modify.  
- OPEN: Bidding Period; full order entry/modify/cancel allowed, subject to seasoning and cutoffs.  
- CLOSE_PENDING: Opening of Closing Window while still in Bidding Period; subject to minimum size and SEC conditions.  
- CLOSING: No order entry/cancel; closing process running; LM Short may be updated until Offering Price submission.  
- CLEARING: Allocation algorithms execute; LM sets Offering Price and confirms allocations; LM Short then frozen.  
- CLOSED: Settlement outputs generated; all trading changes disabled.  
- CANCELED: All orders canceled; offering terminates.[file:102][file:103]

Substates/modifiers:[file:102][file:103]

- HALTED: No new orders; cancellations allowed; price-up-only modifications allowed; closing timers conceptually paused.  
- FROZEN: No order entry or cancellation (except defined MO exceptions); closing timers paused.

Transition rules (high level):

- Automatic transitions (e.g., NEW→UPCOMING→OPEN→CLOSE_PENDING) may run on a scheduler using configured timestamps.[file:103]
- CLOSE_PENDING entry requires:
  - Demand at or above Minimum Acceptable IPO Price sufficient to cover minimum offering size, and  
  - SEC Effectiveness confirmed if required; otherwise block transition and notify LM/MarkOps until resolved.[file:102][file:103]
- CLOSING entry requires the same minimum size conditions as CLOSE_PENDING.[file:102][file:103]
- Transition CLOSE_PENDING→OPEN resets Closing Window timers to initial values (e.g., Window 1).[file:102][file:103]
- CANCELED requires all active orders to have been cancelled.[file:103]

---

## 4. Orders, seasoning, and constraints

### 4.1 Order types and eligibility

Order types:[file:102][file:103]

- Only GTC limit orders are accepted for both bids and offers; any other type or TIF must be rejected.  
- Bid types:
  - Competitive Bids (standard)  
  - Preferential Bids (ExecInst = 'p')  
- Offer types:
  - System-managed supply groups (Primary, Up/Down, Committed Secondary, LM Short).  
  - BD-entered Limit-priced Secondary Shares Offers, when enabled.[file:102][file:103]

Eligibility:

- Only BDs with QSR Active = true and not on the BD Exclude List can submit Competitive Bids.[file:102][file:103]
- Preferential Bids:
  - Enabled by offering-level parameter (No, LM BD only, Yes/all).[file:102][file:103]
  - Shares only (no bonds); MinQty not allowed; not eligible for DTC Tracking.[file:102][file:103]
- DTC IPO Tracking:
  - Only Competitive share bids may set ExecInst = 't'.  
  - Account field required when `execInst = 't'`.[file:102][file:103]

Min/Max quantities:

- Each offering defines Min/Max Bid Quantity and quantity increments; orders outside limits or not respecting increments are rejected with explicit errors.[file:102][file:103]
- Existing orders that become outside new thresholds remain valid until modified, at which point new limits apply.[file:102][file:103]

### 4.2 Seasoning, Time Windows, and MinQty

Time Windows:

- System defines TW1, TW2, TW3, …, each with start time and seasoning duration. Defaults include:
  - TW1 from date-of-OPEN 00:00; 24-hour seasoning.  
  - TW2 from date-of-CLOSING−1 day 00:00; 4-hour seasoning.  
  - TW3 from date-of-CLOSING 00:00; 5-minute seasoning.[file:102][file:103]

Seasoning semantics:

- If seasoning ≥ 12 hours, the clock runs in real time; if < 12 hours, the clock runs only during market hours.[file:102]
- After entry or modification affecting price, quantity, or MinQty, the order enters seasoning; cancel/modify is blocked until seasoningExpiresAt.[file:102][file:103]

Priority Groups:

- Competitive Bids are assigned Priority Groups based on Time Window at entry; groups are used by Priority Group Pro-rata algorithms and for certain cancel restrictions.[file:102][file:103]

MinQty:

- Only Competitive Bids may carry MinQty.[file:102][file:103]
- Offering-level parameters control:
  - Minimum order size to allow MinQty (e.g., 100,000 shares), and  
  - Max MinQty as percentage of order quantity (e.g., 50%).[file:102][file:103]
- MinQty Deadline (default: 23:59:59 EST day before close):  
  - After this deadline, new MinQty or modifications that would retain MinQty are rejected; or MinQty is stripped according to implementation rules.[file:102][file:103]
- During allocation, any order whose computed allocation falls below MinQty must receive zero allocation.[file:102][file:103]

### 4.3 Cancellation and modification

Order actions are constrained by:

- Offering state (e.g., UPCOMING: cancel-only; CLOSING: no entry/cancel).  
- Seasoning (no cancel/modify while under seasoning).  
- Type-specific cutoffs (e.g., Preferential Bids and Limit-priced Secondary Offers cannot be cancelled/modified at or after defined regulatory times post-SEC Effectiveness).[file:102]

Modification semantics:

- Price, quantity, and MinQty changes reset time priority and (re)start seasoning; account changes do not affect time priority.[file:102][file:103]
- In HALTED, only price-up modifications are allowed; all other modifications are rejected.[file:102][file:103]

---

## 5. Price controls and price discovery

Price controls per offering:[file:102][file:103]

- Price (or Yield) Range: High and Low, adjustable intraday.  
- Minimum Acceptable IPO Price (or maximum acceptable IPO yield).  
- Price collars: Maximum Price Allowed and Minimum Price Allowed.

Rules:

- Orders outside collars are rejected with explicit error messages; existing out-of-collar orders remain but cannot be newly entered or modified to those prices.[file:102][file:103]
- Minimum Acceptable IPO Price acts as a floor for Indicative Clearing Price and as a gating condition for CLOSE_PENDING/CLOSING (demand must meet minimum size at this price).[file:102][file:103]

Indicative Clearing Price (single-side auction):

- Supply is composed of Primary + Up/Down + Committed Secondary (excluding LM Short).[file:102][file:103]
- Bids are traversed in price-time order from highest to lowest, respecting MinQty, to find the highest price where cumulative demand ≥ supply.[file:102][file:103]

Two-sided auction (extended):

- Includes limit-priced secondary offers on the sell side; the target is the price that maximizes matched quantity, with defined tie-breaking (minimal imbalance, imbalance pressure).[file:102]

Cprice/Cyield and dissemination:

- Cprice/Cyield = max(Minimum Acceptable IPO Price, Indicative Clearing Price), or appropriately converted to yield.[file:102][file:103]
- Publish modes:
  - Min IPO Price  
  - Not Published + Auto On  
  - Not Published + Manual On  
  - Publish NOW[file:102][file:103]
- Cprice and related orderbook aggregates are exposed via:
  - BatMan and BidMan UIs.  
  - REST APIs (e.g., `/api/cprice/:symbol`, `/api/orderbook/:symbol`).  
  - FIX and public site, according to publish mode.[file:102][file:103]

---

## 6. Allocation, execution, and settlement

### 6.1 Allocation algorithms

Supported algorithms:[file:102][file:103]

- **Price-time Priority**: sort by price desc, timestamp asc; fully fill until supply exhausted, then partial fill next order, remainder zero.  
- **Pro-rata**: allocate proportionally to order size for all orders at or better than Offering Price, with two-pass floor + remainder distribution and integer arithmetic.[file:102][file:103]
- **Priority Group Pro-rata**: allocate by Priority Group in order (Group 1, then 2, etc.) with configured group weights, using pro-rata logic within each group.[file:102][file:103]

Preferential allocation phase:

- Phase 1 allocates a configured preferential quantity to Preferential Bids using FCFS by time, respecting per-investor and total caps.[file:102][file:103]
- Any excess preferential demand feeds into the competitive pool for subsequent phases.[file:102][file:103]

MinQty enforcement:

- Orders whose allocation would fall below MinQty must receive zero allocation; this applies in all algorithms.[file:102][file:103]

Allocation integrity:

- System must expose validation that:
  - Total allocated ≤ total supply (including LM Short where appropriate).  
  - MinQty handling, Priority Group ordering, and LM Short treatment follow spec.[file:102][file:103]

### 6.2 Execution and post-trade

After LM confirms allocations:[file:102][file:103]

- Generate BUY trades per investor order and a single aggregate SELL trade for LM.  
- Generate:
  - NSCC UTC settlement file for locked-in trades.  
  - DTC IPO Tracking file containing only DTC-tracked trades derived from ExecInst 't'.  
- Compute commission and selling concession amounts per trade based on offering spread/concession settings and expose them via reports.[file:102][file:103]
- Support bust/reallocation workflows that reverse trades and allow re-running allocation, with full audit log of changes.[file:102][file:103]

---

## 7. Parameters, reporting, and audit

System-wide parameters:

- Trading calendar (holidays), market hours, asset-type defaults (sizes, increments, thresholds), seasoning defaults, closing windows configuration.[file:102][file:103]

Offering-level parameters:

- Asset/security type, state, timing, Min/Max order quantity, ranges/collars, Minimum Acceptable IPO Price, allowed bid/offer types, DTC policy, allocation algorithm, Time Windows, Closing Windows, listing minimums, spreads/concessions.[file:102][file:103]

BD Exclude List:

- Per-offering list of BD identifiers; when a BD is excluded, the system must block order entry and, where configured, market data access for that BD.[file:102][file:103]

Reporting and audit:

- Exports:
  - Orderbook snapshots.  
  - Allocation details per order, per BD, and per Priority Group.  
  - Commission/selling concession per BD.[file:102][file:103]
- Audit logs:
  - Offering state changes and parameter changes.  
  - Order lifecycle events (entry, modification, cancellation, allocation, bust).  
  - User identity and timestamp for each change.[file:102][file:103]
