This document is a functional description of ClearingBid's price and
demand discovery platform. **This document should be read together with
section 5 to 14 of the "ClearingBid Participant System Guidelines"** as
the latter describes the functionality from a participating
broker-dealer's point of view, i.e., it describes the delivered user
functionality.

Overview

ClearingBid ("CB") provide equal access for all investors, via their
Broker-Dealer, to place Bids for new issue securities using a patented
method and system. Issuers and holders of securities use the PDP for
initial public offerings and follow-on public offerings, collectively
referred to as an "IPO".

PDP is designed to discover a market demand-based price of new issue
securities based on a transparent Auction Process providing equal access
to all interested investors.

Bids are sent by Broker-Dealers to the PDP as Good-'til-Cancelled
("GTC") limit orders in the same manner that secondary market GTC limit
orders are sent to an exchange. Bids are entered and managed as FIX
transactions or manually via the CB BidMan app (Bid Manager).

The Lead Manager monitors and manages an auction using the CB BatMan app
(Bidding and Auction Management). MarkOps refers to the CB PDP market
operations app.

After the auction close and the IPO is priced, allocations are created
on a bid-by-bid basis (i.e., individual order basis) and sent to NSCC
Universal Trade Capture ("UTC") as locked-in trades for participation in
standard NSCC/DTC/FED NSS continuous netting, clearing and settlement
processing. This requires CBM's clearing broker ("CBM Clearing") to have
a Qualified Special Representative ("QSR") relationship with all
participating Broker-Dealers' clearing brokers and the Lead
Underwriter's clearing broker. CB's clearing broker's net securities and
cash position will always be zero as all accepted Bids are
simultaneously accompanied by an equal contra IPO sell quantity of
primary and secondary shares/bonds. NSCC becomes the principal
counterparty to all locked-in trades after they have been accepted by
UTC.

**CB IPO process overview**:

- New offering registration (no info published)

- Upcoming offering announcement and detailed information

- Auction Process

  - Bidding Process

  - Closing Process

  - Pricing Process

  - Allocation Process

  - Trade Execution, Confirmation and Clearing Process

- Closed/Cancelled offering information

Table of Contents

[1 Definitions [4](#definitions)](#definitions)

[2 Development Phases [11](#development-phases)](#development-phases)

[3 Bidding Process [11](#bidding-process)](#bidding-process)

[3.1 Orders -- Bids and Offers
[11](#orders-bids-and-offers)](#orders-bids-and-offers)

[3.2 Bids [11](#bids)](#bids)

[[3.2.1]{.mark} [Competitive Bids / Competitive Orders (aka
standard/normal Bids)]{.mark}
[12](#competitive-bids-competitive-orders-aka-standardnormal-bids)](#competitive-bids-competitive-orders-aka-standardnormal-bids)

[[3.2.1.1]{.mark} [Minimum Quantity Accepted Condition (MinQty)]{.mark}
[13](#minimum-quantity-accepted-condition-minqty)](#minimum-quantity-accepted-condition-minqty)

[[3.2.1.2]{.mark} [Request for tracking by the DTC (ExecInst
't')]{.mark}
[14](#request-for-tracking-by-the-dtc-execinst-t)](#request-for-tracking-by-the-dtc-execinst-t)

[[3.2.2]{.mark} [Preferential Bids / Preferential Orders (ExecInst
'p')]{.mark}
[14](#preferential-bids-preferential-orders-execinst-p)](#preferential-bids-preferential-orders-execinst-p)

[[3.2.3]{.mark} [Minium / Maximum Bid Quantity Allowed]{.mark}
[16](#minium-maximum-bid-quantity-allowed)](#minium-maximum-bid-quantity-allowed)

[3.3 Offers [16](#offers)](#offers)

[[3.3.1]{.mark} [Primary Shares/Bonds]{.mark}
[17](#primary-sharesbonds)](#primary-sharesbonds)

[[3.3.2]{.mark} [Up/down-size of Primary Shares/Bonds]{.mark}
[17](#updown-size-of-primary-sharesbonds)](#updown-size-of-primary-sharesbonds)

[[3.3.3]{.mark} [Committed Secondary Shares/Bonds via Lead
Manager]{.mark}
[17](#committed-secondary-sharesbonds-via-lead-manager)](#committed-secondary-sharesbonds-via-lead-manager)

[[3.3.4]{.mark} [Limit Priced Secondary Shares Offers (two-sided
auctions)]{.mark}
[18](#limit-priced-secondary-shares-offers-two-sided-auctions)](#limit-priced-secondary-shares-offers-two-sided-auctions)

[[3.3.5]{.mark} [Lead Manager Short]{.mark}
[18](#lead-manager-short)](#lead-manager-short)

[[3.4]{.mark} [Order Cancellation]{.mark} and [Modification]{.mark}
[19](#order-cancellation-and-modification)](#order-cancellation-and-modification)

[[3.5]{.mark} [Halted Orders]{.mark}
[19](#halted-orders)](#halted-orders)

[[3.6]{.mark} [Price Collars -- Maximum Price Allowed and Minimum Price
Allowed]{.mark}
[20](#price-collars-maximum-price-allowed-and-minimum-price-allowed)](#price-collars-maximum-price-allowed-and-minimum-price-allowed)

[[3.7]{.mark} [Price Range -- High Price Range & Low Price Range]{.mark}
[21](#price-range-high-price-range-low-price-range)](#price-range-high-price-range-low-price-range)

[[3.8]{.mark} [Minimum Acceptable IPO Price]{.mark}
[21](#minimum-acceptable-ipo-price)](#minimum-acceptable-ipo-price)

[[3.9]{.mark} [Seasoning Periods]{.mark}
[22](#seasoning-periods)](#seasoning-periods)

[[3.10]{.mark} [Priority Groups]{.mark}
[23](#priority-groups)](#priority-groups)

[4 Limit Orderbook [23](#limit-orderbook)](#limit-orderbook)

[[4.1]{.mark} [Indicative Clearing Price aka Cprice]{.mark}
[23](#indicative-clearing-price-aka-cprice)](#indicative-clearing-price-aka-cprice)

[[4.1.1]{.mark} [Cprice/Cyield]{.mark}
[24](#cpricecyield)](#cpricecyield)

[[4.1.2]{.mark} [Single-side Auction]{.mark}
[25](#single-side-auction)](#single-side-auction)

[[4.1.3]{.mark} [Two-sided Auction]{.mark}
[25](#two-sided-auction)](#two-sided-auction)

[[4.2]{.mark} [BidMan app -- Bid Management app for
Broker-Dealers]{.mark}
[26](#bidman-app-bid-management-app-for-broker-dealers)](#bidman-app-bid-management-app-for-broker-dealers)

[[4.3]{.mark} [BatMan app -- Bidding and Auction Management app for Lead
Managers]{.mark}
[28](#batman-app-bidding-and-auction-management-app-for-lead-managers)](#batman-app-bidding-and-auction-management-app-for-lead-managers)

[[4.3.1]{.mark} [Broker-Dealer Exclude List (Blacklist)]{.mark}
[30](#broker-dealer-exclude-list-blacklist)](#broker-dealer-exclude-list-blacklist)

[4.4 [Orderbook States]{.mark}
[31](#orderbook-states)](#orderbook-states)

[4.4.1 Allowed Order Book State Transitions
[37](#allowed-order-book-state-transitions)](#allowed-order-book-state-transitions)

[[4.5]{.mark} [Scheduling]{.mark} [37](#scheduling)](#scheduling)

[[4.6]{.mark} [Orderbook Data Structure and Processing]{.mark}
[38](#orderbook-data-structure-and-processing)](#orderbook-data-structure-and-processing)

[[5]{.mark} [Closing Process / Close Pending state]{.mark}
[40](#closing-process-close-pending-state)](#closing-process-close-pending-state)

[[5.1]{.mark} [Closing Window]{.mark}
[40](#closing-window)](#closing-window)

[[5.2]{.mark} [Preferential Bids and Listing Exchange Minimum]{.mark}
[41](#preferential-bids-and-listing-exchange-minimum)](#preferential-bids-and-listing-exchange-minimum)

[[5.2.1]{.mark} [Preferential Bids Monitoring during the Auction]{.mark}
[43](#preferential-bids-monitoring-during-the-auction)](#preferential-bids-monitoring-during-the-auction)

[[5.2.1.1]{.mark} [Transition from Open to Close Pending State
Blocked]{.mark}
[45](#transition-from-open-to-close-pending-state-blocked)](#transition-from-open-to-close-pending-state-blocked)

[[5.2.2]{.mark} [Preferential Allocation Algorithm]{.mark}
[45](#preferential-allocation-algorithm)](#preferential-allocation-algorithm)

[6 Pricing Process [46](#pricing-process)](#pricing-process)

[[6.1]{.mark} [Priority Group Fill Level Simulation and
Monitoring]{.mark}
[46](#priority-group-fill-level-simulation-and-monitoring)](#priority-group-fill-level-simulation-and-monitoring)

[[6.2]{.mark} [Pro-rata Simulation and Monitoring]{.mark}
[47](#pro-rata-simulation-and-monitoring)](#pro-rata-simulation-and-monitoring)

[[6.3]{.mark} [Non-Pro-rata Allocation method Simulation]{.mark}
[47](#non-pro-rata-allocation-method-simulation)](#non-pro-rata-allocation-method-simulation)

[[6.4]{.mark} [Offering Price Submission and Publication]{.mark}
[48](#offering-price-submission-and-publication)](#offering-price-submission-and-publication)

[7 Allocations Process [48](#allocations-process)](#allocations-process)

[[7.1]{.mark} [Allocation to Preferential Bids]{.mark}
[48](#allocation-to-preferential-bids)](#allocation-to-preferential-bids)

[[7.2]{.mark} [Pro-rata Allocation]{.mark}
[49](#pro-rata-allocation)](#pro-rata-allocation)

[[7.3]{.mark} [Priority Group Pro-rata Allocation]{.mark}
[49](#priority-group-pro-rata-allocation)](#priority-group-pro-rata-allocation)

[[7.4]{.mark} [Price-Time Priority Allocation]{.mark}
[50](#price-time-priority-allocation)](#price-time-priority-allocation)

[[7.5]{.mark} [Time Priority Allocation]{.mark}
[51](#time-priority-allocation)](#time-priority-allocation)

[[7.6]{.mark} [High Price Before Time Priority Allocation]{.mark}
[51](#high-price-before-time-priority-allocation)](#high-price-before-time-priority-allocation)

[[7.7]{.mark} [Allocation Integrity Verification]{.mark}
[51](#allocation-integrity-verification)](#allocation-integrity-verification)

[8 Trade Execution, Confirmation and Clearing Process
[52](#trade-execution-confirmation-and-clearing-process)](#trade-execution-confirmation-and-clearing-process)

[[8.1]{.mark} [Auction Close Day Executions]{.mark}
[52](#auction-close-day-executions)](#auction-close-day-executions)

[8.2 Next Trading Day Executions
[53](#next-trading-day-executions)](#next-trading-day-executions)

[8.3 Clearing and Settlement File Generation
[54](#clearing-and-settlement-file-generation)](#clearing-and-settlement-file-generation)

[8.4 FIX Trade Execution and Order Cancellation Messages
[54](#fix-trade-execution-and-order-cancellation-messages)](#fix-trade-execution-and-order-cancellation-messages)

[8.5 Error Handling and Reversal/Bust Management
[54](#error-handling-and-reversalbust-management)](#error-handling-and-reversalbust-management)

[9 Parameters Attributes
[55](#parameters-attributes)](#parameters-attributes)

[[9.1]{.mark} [PDP System Wide Parameters and Attributes]{.mark}
[55](#pdp-system-wide-parameters-and-attributes)](#pdp-system-wide-parameters-and-attributes)

[9.2 Asset Types [55](#asset-types)](#asset-types)

[[9.3]{.mark} [Common Attributes for All Asset Types (Stock, Bond, and
Fund/ETF)]{.mark}
[56](#common-attributes-for-all-asset-types-stock-bond-and-fundetf)](#common-attributes-for-all-asset-types-stock-bond-and-fundetf)

[[9.4]{.mark} [Stock Attributes]{.mark}
[59](#stock-attributes)](#stock-attributes)

[[9.5]{.mark} [Bond Attributes]{.mark}
[62](#bond-attributes)](#bond-attributes)

[[9.6]{.mark} [Funds/ETF Attributes]{.mark}
[66](#fundsetf-attributes)](#fundsetf-attributes)

[[9.7]{.mark} [Upcoming Offering Information]{.mark}
[70](#upcoming-offering-information)](#upcoming-offering-information)

[[9.8]{.mark} [Closed Offering Information]{.mark}
[73](#closed-offering-information)](#closed-offering-information)

[10 System Features and Parameters
[77](#system-features-and-parameters)](#system-features-and-parameters)

[10.1 User Management and Authorization
[77](#user-management-and-authorization)](#user-management-and-authorization)

[10.2 Default Values Table
[80](#default-values-table)](#default-values-table)

[11 Interfaces [84](#interfaces)](#interfaces)

[11.1 FIX [84](#fix)](#fix)

[11.2 ClearingBid Web-Portal API
[87](#clearingbid-web-portal-api)](#clearingbid-web-portal-api)

[11.3 BidMan -- Broker-dealer App
[91](#bidman-broker-dealer-app)](#bidman-broker-dealer-app)

[11.4 BatMan -- Lead Manager App
[94](#batman-lead-manager-app)](#batman-lead-manager-app)

[11.5 MarkOps -- PDP System and Market Operations App
[98](#markops-pdp-system-and-market-operations-app)](#markops-pdp-system-and-market-operations-app)

[11.6 Market Surveillance
[103](#market-surveillance)](#market-surveillance)

[11.7 Network Monitoring and Connectivity Management
[104](#network-monitoring-and-connectivity-management)](#network-monitoring-and-connectivity-management)

[12 Notes, Remaining Questions, Issues and TBDs
[108](#notes-remaining-questions-issues-and-tbds)](#notes-remaining-questions-issues-and-tbds)

[13 Version history [109](#version-history)](#version-history)

# Definitions

**This section only contains definitions in addition to the definitions
in the "ClearingBid Participant System Guidelines" document**. If the
same definition is contained in both documents, then the Participant
Guidelines definition prevails. References in superscript square
brackets (e.g., ^[\[1\]\[2\]]{.underline}^) refers source documents.

**aVWAP**

- Volume Weighted Average Price (VWAP) of **all** orders in the order
  book^[\[1\]\[2\]]{.underline}^. See also rVWAP.

**Allocation Priority Group / Priority Group**

- In allocation algorithms (Time Priority Pro Rata and Price Priority
  Pro Rata), orders are grouped by entry day or price. Lower group
  numbers (e.g., Prio 1) receive higher allocation
  preference^[\[1\]\[2\]]{.underline}^.

**Asset Types**

- The system recognizes three main asset types:

  - **Stock:** Price discovery by bidding on price per unit.

  - **Bond:** Price discovery by bidding on yield equivalent (defined by
    Yield Type).

  - **Fund/ETF:** No price discovery; demand discovery by price level;
    issue quantity varies by net asset value ("NAV"); bidding on price
    per unit^[\[1\]\[2\]]{.underline}^.

**BatMan**

- The Bidding & Auction Management application provided by ClearingBid
  for Underwriter Participants and Lead Manager to manage securities
  offerings.

**Bid Period / Bidding Period**

- The period during which orders can be placed, modified, or cancelled
  for an offering. The period begins when the orderbook transition to
  the Open state, continues through the Pre-open/Close Pending states
  and ends at transition to the Pricing state. When in the Pre-open
  state order entry, modification and cancellation are temporary
  unavailable.

**BidMan**

- The Bid & Offer Management application provided by CB to Broker-Dealer
  Participants for management of Bids and Offers.

**Bond Offering Price**

- See Offering Price.

**Clearing Price / Clearing Yield**[^1]

- See the Guidelines for the definition of Clearing Price. In this
  document the Clearing Price and Clearing Yield are synonyms, i.e.,
  when Clearing Price is referred it means Clearing Yield when the
  context is such that the price referred is expressed as a yield. For
  example, the statement "the Offering Price cannot exceed the Clearing
  Price" in the context of yield pricing is to be interpreted as "the
  Offering Yield cannot be lower than the Clearing Yield".

**Close / Closing**

- The state following the Bidding Period, i.e., the state following the
  Close Pending state. No further orders, modifications, or
  cancellations are allowed. The closing process is designed to minimize
  price manipulation risks and occurs within a randomly selected time
  window^[\[3\]\[1\]\[2\]]{.underline}^.

**Cprice**

- An abbreviation for Indicative Clearing Price.

**Current Price Range / Current Yield Range**

- The price or yield range set by the Lead Manager on behalf of the
  Issuer, which may be adjusted during the Bid
  Period^[\[3\]\[1\]\[2\]]{.underline}^.

**Distribution Fee**

- The compensation (in percentage points) that the Lead Manager receives
  from the issuer for Fund/ETF new issues. Paid directly by the
  issuer^[\[1\]\[2\]\[4\]]{.underline}^.

> **Dividend** and **Dividend Yield** and **Indicative Dividend Yield**

- Optional reference dividend value in dollars for stocks. The
  Indicative Dividend Yield is the Dividend divided by the Indicative
  Clearing Price and the Dividend Yield is the Dividend divided by
  Offering Price.

**Face Value / Par Value**

- The amount paid to a bondholder at maturity (typically \$1,000 per
  bond)^[\[1\]\[2\]]{.underline}^.

**FIX**

- Financial Information eXchange protocol. Used for electronic
  communication of orders and market data between Broker-Dealers and the
  CB platform^[\[5\]]{.underline}^.

**Gross Underwriting Spread**

- The total compensation (expressed in percentage points) that the
  Underwriters receive from the issuer for stock and bond offerings.
  Includes the Selling Concession^[\[1\]\[2\]]{.underline}^. For
  example, if the Gross Underwriting Spread is 5%, and a stock is issued
  at \$20, then the investor pays \$20 but the issuer only receives
  \$19.00. The Lead Manager/underwriter receives \$1.00 as compensation.
  The **Selling Concession (or commission)** to be paid to the
  Broker-Dealer firm whose brokers sell the stock is part of the Gross
  Underwriting Spread. It is retained by the Lead Manager/underwriter if
  the broker (who is selling the stock to their investor/client) works
  for the Lead Manager. For example, if the Selling Concession is 40%,
  then in the example above, \$0.40 is paid to the Broker-Dealer that
  sells the stock. If the Lead Manager sets the Bond Offering Price for
  a bond with a Face Value of \$1000 at \$995, then the investor pays
  \$995; however, the issuer would only receive \$975.10 per bond,
  assuming a 2% underwriting fee.

**Guidelines**

- The "ClearingBid Participant System Guidelines" document.

**High Price Before Time Priority Allocation**

- See section 7.5.

**High Price Range / High Yield Range**

- See Price Range.

**Indicative Clearing Price / Yield**

- See the Guidelines for definition. See also CPrice.

**Lead Manager**

- The underwriter selected by the issuer to lead the offering process,
  including setting price ranges and managing
  allocations^[\[3\]\[1\]\[2\]]{.underline}^.

**Lead Manager Short Quantity**

- See Primary Quantity

**Limit Order**

- An order to buy or sell at a specified price or better. In PDP, all
  Limit Orders are Good 'Til Canceled orders
  ("GTC"))^[\[1\]\[2\]]{.underline}^.

**Listing Exchange aka as Market**

- The exchange at which a security is listed.

**Low Price Range / Low Yield Range**

- See Price Range.

**Market aka as Listing Exchange**

- The exchange at which a security is listed.

**Market Order**

- An order to buy or sell immediately at the best available price. Not
  currently implemented; however, a feature that may be allowed in the
  future^[\[3\]\[1\]]{.underline}^.

**Minimum Allocation**

- The minimum allocation per account during the allocation process.
  Reserved for future use^[\[3\]\[1\]\[2\]]{.underline}^.

**Minimum Quantity Condition (MinQty)**

- An investor-specified minimum quantity allocation in number of units
  per order that must be met when allocations are determined. If not
  met, then the order will be excluded from
  allocations^[\[3\]\[1\]\[2\]]{.underline}^. There are multiple
  instrument level parameters restricting the use of the MinQty
  condition -- minimum order size that are allowed to set MinQty,
  maximum allowed MinQty as percentage of order quantity and depending
  on the order's assigned Priority Group.

**NAV (Net Asset Value)**

- The per-share value of a Fund/ETF, calculated by an external source
  and used as the execution price for Fund/ETF
  offerings^[\[1\]\[2\]]{.underline}^.

**New Issue Announcement**

- Information about a new issue made available to investors and the
  public, including Bid Period and Current Price
  Range^[\[3\]\[1\]\[2\]]{.underline}^.

**Offering Coupon**

- For bonds, the Offering Coupon rate (%) in combination with the bond
  Offering Price (\$) are set by the Lead Manager and Issuer so that the
  offering yield is equal to or greater than the Clearing
  Yield^[\[1\]\[2\]]{.underline}^.

**Offering Price / Bond Offering Price**

- The price to be paid at settlement for stocks or bonds, set by the
  Lead Manager and Issuer^[\[3\]\[1\]\[2\]]{.underline}^. The (Bond)
  Offering Price will be based on the Clearing Price (Yield) and other
  factors, such as market conditions, the strength of the order book and
  the performance of other comparable new issues. The Offering Price
  cannot exceed the Clearing Price. The Bond Offering Price combined
  with Offering Coupon cannot be lower than the Clearing Yield. The Bond
  Offering Price might differ from the bond Face Value.

**Opening**

- The state transition to the Open state, i.e., the start of the Bidding
  Period.

**Order Book**

- The list of all buy orders (and, in the future sell orders) for an
  offering, organized by price and
  time^[\[3\]\[1\]\[5\]\[2\]]{.underline}^.

**Orderbook States**

- The lifecycle of an order book includes states such as NEW, UPCOMING,
  OPEN, CLOSE PENDING, PRICING, EXECUTION, CLOSED, CANCELLED, FROZEN,
  HALTED and MODIFICATION. Each state governs what actions and
  operations are allowed on orders.

**Price Collars**

- Consists of **Minimum Price Allowed** and **Maximum Price Allowed**
  applicable for new or modified orders (existing orders in the order
  book might be outside the Price Collars). See also Hight Price Range,
  Low Price Range and Minimum Acceptable IPO Price.

**Price-Demand Discovery**

- The real-time process of disseminating aggregated order and price data
  during the Bid Period to facilitate transparent price
  formation^[\[3\]\[1\]\[2\]]{.underline}^.

**Price Range, High Price Range, Low Price Range / Yield Range, Low
Yield Range, High Yield Range**

- The range of prices (for stocks) or yields (for bonds) within which
  the Lead Manager at the time expects the offering to
  price^[\[3\]\[1\]\[2\]]{.underline}^. The Price/Yield Range consists
  of High Price Range (Low Yield Range) and Low Price Range (High Yield
  Range), and is initially specified in the Preliminary Prospectus and
  then further optionally modified intraday. See also Price Collars.

**Price-Time Priority Allocation**

- See section 7.4.

**Pro-Rata Allocation**

- See section 7.2.

**Primary Quantity, Up-/Down-size Quantity, Secondary Quantity, Lead
Manager Short Quantity**

- **Primary Quantity**: Number of new shares or bonds initially offered
  by the issuer.

- **Up-/Down-size**: Incremental change to Primary Quantity

- **Secondary Quantity**: Number of existing shares or bonds offered by
  existing holders^[\[1\]\[2\]]{.underline}^.

- **Lead Manger Short Quantity** ("LM Short"): Greenshoe covered and/or
  naked short. This quantity is not part of the calculation of the
  Cprice nor Clearing Price but is included in the allocated quantity.
  If the LM Short is larger than the excess demand at the Offering Price
  this will show up as a fill rate larger than 100%. Then the LM must
  either decrease the LM Short or the Offering Price.

**Prio 1 Order, Prio 2 Order, etc.**

- Orders assigned to allocation Priority Group, with Prio 1 having the
  highest allocation preference^[\[1\]\[2\]]{.underline}^.

**Priority Group Pro-Rata Allocation**

- See section 7.3.

**rVWAP**

- VWAP of orders within the current **Price
  Range**^[\[1\]\[2\]]{.underline}^. See also aVWAP.

**Secondary Quantity**

- See Primary Quantity

**Security Types**

- Standardized codes for asset types, e.g., CS (Common Stock), CORP
  (Corporate Bond), MF (Mutual Fund), ETF (Exchange Traded
  Fund)^[\[1\]\[2\]]{.underline}^.

**Selling Concession**

- The portion of the Gross Underwriting Spread paid as commission to the
  Broker-Dealer that submits the order^[\[1\]\[2\]]{.underline}^. For
  example, if the Gross Underwriting Spread is 5% and the Selling
  Concession is 40%, then the order submitting Broker-Dealer receives 2%
  of the 5% Gross Underwriting Spread. See further Gross Underwriting
  Spread.

**Settlement**

- The process of generating files with allocation instructions for the
  Lead Manager, Broker-Dealers, and Clearing Agent, enabling the
  transfer of cash and securities^[\[3\]\[1\]\[2\]]{.underline}^.

**Symbol (or Ticker Symbol or just Ticker)**

- A unique, short abbreviation of letters (and sometimes numbers) that
  identifies a specific security.

**Ticker or Ticker Symbol** -- see Symbol

**Trading Day, Trading Hours, Trading Day Start Time** and **Trading Day
End Time** -- see section 4.5.

**Up-/Down-size Quantity**

- See Primary Quantity

**VWAP (Volume Weighted Average Price)**

- The average price of all orders, weighted by
  volume^[\[1\]\[2\]]{.underline}^.

**Yield Range**

- See Price Range.

**Yield Type**

- For bonds, defines how the Lead Manager sets the Offering Price and
  Coupon, and what type of yield the investor is bidding for (e.g.,
  Yield to Maturity, Yield to Call, Spread to
  Treasury)^[\[1\]\[2\]]{.underline}^.

References withing square brackets "\[x\]" in this section:

1.  ClearingBid-Backend-Functions-CMC-v2-7.docx

2.  ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx

3.  CB-Backend-Specification-v1-4-8-comments.docx

4.  Fund-ETF-Summary-v01.docx

5.  BE-Terminal-Reference-v004-2.docx

# Development Phases

PDP Gen 2 is developed in phases. Requirements in this document are
highlighted to indicate development phase and PDP version:

> [V2.0 -- Demo only version]{.mark}
>
> [V2.1 -- Minimal Viable Product (MVP) for single seller equity
> IPO]{.mark}
>
> [V2.2 -- Version 2.2]{.mark}
>
> [V2.x -- Version 2.3 or later]{.mark}
>
> No highlight means either:

- the version is inherited from the section/heading hierarchy; or

- the information is for context.

When possible, technical design of data structures and algorithms should
consider existing requirements specified for future development. This
does not apply to V2.0 development which modules might be discarded and
developed anew in a later version, if that is considered more efficient.

# Bidding Process

The Bidding Process begins when the orderbook transitions to the Open
State, continues through the Close Pending State and ends at transition
to the Pricing State.

## Orders -- Bids and Offers

During an offering's Bidding Period/Process Orders can be placed,
modified, or cancelled. There are two types of Orders; Bids to buy and
Offers to sell.

> Only Limit Orders with time-in-force Good-'til-Canceled (GTC) are
> accepted; other order types or time-in-force instructions are
> rejected.

## Bids

> Bids are entered into the PDP orders via the use of the BidMan
> application, another CB approved trading app or via FIX messages in
> the same manner as secondary market orders are sent to an exchange.
> [For new issue securities with price discovery, e.g., equities and
> bonds, only Bids with limit price and the GTC (Good-'til-Cancelled)
> time-in-force attribute are accepted. All other order types and
> time-in-force are rejected]{.mark}. [For new issue securities without
> price discovery (e.g., ETFs and funds) Market-at-Close bids are also
> accepted]{.mark}.

There are two types of Bids; Competitive Bids (standard/normal Bids) and
Preferential Bids.

> **Summary Table: Competitive vs. Preferential Bids**

+----------------+--------------------------+-----------------------+
| **Feature**    | **Competitive Bids**     | **Preferential Bids** |
+================+==========================+=======================+
| Definition     | All bids not designated  | Bids designed to      |
|                | as Preferential Bids     | fulfill exchange      |
|                |                          | listing minimum       |
|                |                          | requirements          |
+----------------+--------------------------+-----------------------+
| Allocation     | After Preferential Bids, | Receive allocation    |
| Priority       | based on price and/or    | ahead of Competitive  |
|                | time depending on        | Bids                  |
|                | auction allocation       |                       |
|                | algorithm type           |                       |
+----------------+--------------------------+-----------------------+
| Minimum        | Allowed, subject to      | Not allowed           |
| Quantity       | thresholds and time      |                       |
| Condition;     |                          |                       |
| MinQty         |                          |                       |
+----------------+--------------------------+-----------------------+
| Eligible for   | Yes                      | No                    |
| DTCC Tracking  |                          |                       |
+----------------+--------------------------+-----------------------+
| Modification/  | Allowed before auction   | Stricter withdrawal   |
|                | close, subject to        | rules after reg.      |
| /Withdrawal    | seasoning periods        | statement effective   |
+----------------+--------------------------+-----------------------+

### [Competitive Bids / Competitive Orders (aka standard/normal Bids)]{.mark}

> Competitive Bids compete for allocation using order entry limit price
> and/or time depending on the selected allocation algorithm for the
> auction.
>
> **Key Features of Competitive Bids:**

- **Eligibility:** Broker-Dealers with the QSR Active attribute true and
  that are not on the Exclude List can place Competitive Bids.

- **Order attributes**: Symbol, Order ID (ClOrdID), Account, Side,
  Price, Quantity, MinQty (order contingent on minimum allocation),
  Timestamps, Type, Time In Force, ExecInst (Preferential Bid flag and
  Tracking flag), OrderCapacity (Instiutional flag, Retail flag),
  ClearingFirm and ClearingAccount.

- **Time priority**: PDP time of entry or last modification of price,
  quantity or MinQty of the bid. Other modifications such as, e.g.,
  Account, Preferential Bid flag, Retail Bid flag, Institutional Bid
  flag, tracking Flag, ClearingFirm and ClearingAccount do not impact
  time priority.

- **Modification/Withdrawal:** Bids can be modified or withdrawn before
  the auction closes, subject to \"Seasoning Periods\" that restrict how
  soon after entry/modification a Bid can be canceller or have its
  price, quantity or MinQty can changed.

- [**Minimum Quantity Condition:** Investors may specify that their bid
  is only valid if they receive at least a minimum number of securities,
  subject to specific rules and thresholds.]{.mark}

- [**Tracking:** Only Competitive Bids are eligible to be flagged for
  tracking by the DTCC IPO Tracking System, using the Execution
  Instruction field set to \"t\".]{.mark}

#### [Minimum Quantity Accepted Condition (MinQty)]{.mark}

> The Minimum Quantity Accepted Condition (MinQty) is a feature in
> ClearingBid\'s IPO auction process that allows investors to specify a
> minimum number of securities they are willing to accept in the
> allocation. If the pro-rata allocation for an order falls below this
> threshold, the order receives no allocation. Key aspects include:
>
> **Key Characteristics of MinQty**

1.  **Purpose**: Ensures investors only receive allocations meeting
    their minimum size requirement, preventing fractional fills below
    their threshold\[1\]\[3\].

2.  **Eligibility Restrictions**: Only **Competitive Bids** may include
    MinQty conditions.

3.  **Order Size Requirements**:

    a.  Minimum order size to set MinQty: Default **100,000 shares /
        10,000 bonds** intraday configurable.

    b.  Maximum Allowable MinQty: Default **50% of the order size**
        intraday configurable.

4.  **Minimum Quantity Entry Deadline (MinQty Deadline)**

    a.  Default 11:59:59 PM EST the day before the scheduled IPO close;
        intraday configurable.

    b.  For an order with MinQty, any request for modification of said
        orders price or quantity, received after the MinQty Deadline,
        must also remove the MinQty condition. If said modification
        request does not remove the MinQty condition, the modification
        request is rejected with an error message informing that the
        MinQty condition must be removed or the order cancelled because
        the order modification time is after the MinQty order entry
        deadline.

#### [Request for tracking by the DTC (ExecInst 't')]{.mark}

> The "Request for tracking by the DTC [(ExecInst 't')]{.mark}" refers
> to the process by which IPO share allocations are specifically
> designated to be monitored and settled through the DTC (Depository
> Trust Company) IPO Tracking System or the DTC Institutional Delivery
> System, as opposed to standard net settlement processing by NSCC CNS
> (Continuous Net Settlement) and DTC/FED NSS net settlement. The
> tracking mechanism is used for allocations that require more granular
> tracking, typically for regulatory, compliance, or issuer
> requirements, such as ensuring minimum distribution standards or
> tracking specific investors.
>
> **How Tracking is Requested and Implemented**

- **Eligibility**: Only **Competitive Bids** (not Preferential Bids) are
  eligible for DTC IPO Tracking System tracking. DTC IPO Tracking System
  is for shares only, not bonds.

- **Flagging a Bid for Tracking**: To request tracking, the Participant
  must flag the relevant Bid by setting the ExecInst (Execution
  Instruction) field to t (for tracking) and the "OrderCapacity" field
  to B for institutional clients or C for retail clients in the FIX
  message or manually by checking the corresponding DTC Tracking
  checkboxes in the BidMan application.

- **Clearing Firm**: For tracked Bids it is mandatory to populate FIX
  tag 439 "ClearingFirm" (the Broker-Dealer's clearing firm's NSCC
  number) and FIX tag 440 "ClearingAccount" (account for the Bid). These
  fields can also be entered/modified in BidMan.

- **Clearing and Settlement** (for more details see paragraph
  **EXECUTION:** in section 4.4)

### [Preferential Bids / Preferential Orders (ExecInst 'p')]{.mark}

> **Preferential Bids** are a specific type of bid in the ClearingBid
> IPO auction process, designed to help issuers comply with stock
> exchange minimum distribution requirements, such as minimum numbers of
> shareholders or round lot holders at the time of listing.
>
> **Definition and Purpose**

- **Purpose:** Preferential Bids are intended to facilitate compliance
  with exchange listing minimums, such as a required number of holders
  or round lot holders, which are often prerequisites for a new issue to
  be listed on an exchange. Preferential Bids are only valid for shares,
  not bonds.

- **Eligibility:** These bids are eligible for allocation of shares
  **ahead of competitive bids**, but only up to the extent needed to
  satisfy the exchange's minimum listing requirements.

> **Key Features of Preferential Bids:**

- **Eligibility:** When **Preferential Bids allowed** is [not set to
  No]{.underline}, then the following Broker-Dealers with the QSR Active
  attribute true and not on the Exclude List can place Preferential
  Bids:

  - When set to Lead Manager BD Only, then the Lead Manager BD only.

  - When set to Yes, then all eligible Broker-Dealers.

  <!-- -->

  - **Minimum Preferential Quantity** (MinPrefSec): Number of
    shares/bonds set aside for allocation to Preferential Bids. If
    MinPrefSec is set too low (so that the Exchange Listing Minimum
    requirement cannot be fulfilled), then the PDP calculated minimum
    quantity to fulfill the Exchange Listing Minimum requirement will be
    allotted to Preferential Bids. Hence, MinPrefSec will only increase
    the allotted quantity above what is required by the exchange.

  - **Order Type:** Just like any order, only GTC limit orders are
    accepted; other order types or time-in-force instructions are
    rejected. A Preferential Bid is recognized as such by the PDP when
    the Participant has populated the Account field with an investors
    account identification and the ExecInst field with the value 'p'
    (Preferential Bid). The Account field and the ExecInst field can be
    populated either through inclusion in the FIX message or manually
    via the of BidMan.

- **Time priority**: PDP time of entry or last modification of price,
  quantity or MinQty of the bid. Other modifications do not impact time
  priority.

- **Modification/Withdrawal:** Bids can be modified or withdrawn until
  ***one hour after the registration statement becomes effective***,
  subject to \"seasoning periods\" that restrict how soon after
  entry/modification a bid can be changed.

- **Minimum Quantity Condition:** Not allowed.

- **Tracking:** ***Not eligible*** to be flagged for tracking by the
  DTCC IPO Tracking System.

- **Any part of a Preferential Bid that do not receive Preferential
  Allocation participate in the auction equally together with all
  Competitive Bids.**

> **Information for context**

- Allocation Priority: Only a limited number of shares are allocated to
  Preferential Bids, and those allocations occur before competitive bids
  are considered. The allocation is first-come, first-served based on
  the latest entry or modification time of the bid. Details in section
  7.1 Allocation to Preferential Bids.

- If a Preferential Bid is for more shares than the preferential
  allocation, the excess is treated as a Competitive Bid.

- All Preferential Bids from the same investor must be placed with the
  same broker-dealer (Broker-Dealer terms and conditions).

- The broker-dealer must use the same value in the Account field for all
  Preferential Bids from the same investor (CB Participant terms and
  conditions).

- Investors must consent not to place Preferential Bids for the same
  offering with any other broker-dealer unless all existing Preferential
  Bids are withdrawn first (Broker-Dealer terms and conditions).

- **I**nvestors can place multiple Preferential Bids at different
  prices. Preferential allocation is per bidder, not per individual bid.

### [Minium / Maximum Bid Quantity Allowed]{.mark}

> **Minimum Bid Quantity Allowed** refers to the smallest bid size that
> can be submitted for a new issue offering on the ClearingBid platform.
> **Maximum Bid Quantity Allowed** is the largest permissible bid size
> for a single order or account in a given offering. These parameters
> are implemented to ensure orderly price discovery, prevent market
> manipulation, and comply with regulatory or issuer-specific
> requirements.
>
> **Default values**:

- The Minimum Bid Quantity default is **10**, and **configurable
  intraday**.

- The Maximum Bid Quantity default is **10,000 for Bonds** and
  **1,000,000 for Stocks and ETF/Funds**, both configurable intraday.

- Bids below/above these thresholds are automatically rejected with a
  message explaining the reason.

- **Application:** If the minimum/maximum is changed after bids have
  already been entered, those existing bids remain valid in the book as
  long as not modified; only new or modified bids are subject to the
  updated minimum.

## Offers

The Offers to sell shares/bonds are quantified in five groups. The Limit
Priced Secondary Shares Offers are managed by Broker-Dealers, the other
groups are managed by the Lead Manager:

1.  [Primary Shares/Bonds (market-at-close order)]{.mark}

2.  [Up/down-size of Primary Shares/Bonds (market-at-close
    order)]{.mark}

3.  [Committed Secondary Shares/Bonds via Lead Manager (market-at-close
    order)]{.mark}

4.  [Limit Priced Secondary Shares Offers (only applicable for two-sided
    auctions)]{.mark}

5.  [Lead Manager Short]{.mark}

**Primary Shares/Bonds** refer to the **number of new securities
initially offered by the issuer in a securities offering**. These are
created and sold by the company itself, with proceeds going directly to
the issuer. The term is distinct from **Secondary Shares/Bonds**, which
are existing securities sold by current holders. The Lead Manager Short
are a quantity sold short and normally covered by an overallotment of
primary shares (aka greenshoe option).

The Lead Manager Short quantity is not included in the Indicative
Clearing Price or Clearing Price calculation, all other quantities are
included.

### [Primary Shares/Bonds]{.mark}

The total quantity of primary shares or bonds initially offered by the
issuer.

- Quantity entered and updated by the Lead Manager.

- Can be updated intraday, but only before transitioning to the Pricing
  state.

- PDP keeps a log with date and time of updates.

- For context: changed quantity should normally be done via
  Up/Down-size, but the primary quantity can change due to for example
  some securities being set aside for a price taker strategic
  investor[^2].

### [Up/down-size of Primary Shares/Bonds]{.mark}

> **Up-/Down-size** refers to the ability to **increase or decrease the
> number of primary securities offered**. This mechanism allows the Lead
> Manager and issuer to adjust the total offering size in response to
> investor demand observed during the Bidding Process (price discovery).

- Quantity entered and updated by the Lead Manager.

- Updated intraday, but only before transitioning to the Pricing state.

- Quantity [is]{.underline} included in the Indicative Clearing Price or
  Clearing Price calculation

- PDP keeps a log with date and time of updates.

### [Committed Secondary Shares/Bonds via Lead Manager]{.mark}

The total quantity of secondary securities offers that the Lead Manager
includes in the auction as a lump sum:

- Quantity entered and updated by the Lead Manager.

- Can be updated intraday, but only before transitioning to the Pricing
  state.

- Quantity [is]{.underline} included in the Indicative Clearing Price or
  Clearing Price calculation

- PDP keeps a log with date and time of updates.

- For context: The quantity might be made up of offers from several
  investors, but that information is not known to or kept by the PDP.

### [Limit Priced Secondary Shares Offers (two-sided auctions)]{.mark}

> Limit Priced Secondary Shares Offers are sell orders placed by
> authorized Broker-Dealers on behalf of individual holders of existing
> secondary shares.

- **Eligibility**: When Limit Priced Secondary Shares Offers Allowed is
  not set to No, then the following Broker-Dealers with the QSR Active
  attribute true and not on the Exclude List can place Preferential
  Bids:

<!-- -->

- When set to Lead Manager BD Only, then the Lead Manager BD only.

- When set to Yes, then all eligible Broker-Dealers.

<!-- -->

- Only limit price Offers with GTC time-in-force attribute are accepted.

- Offers can be modified or cancelled subject to the same Seasoning
  Periods as for Preferential Bids.

- Offers are subject to current Price Collars (Maximum/Minimum Price
  Allowed).

- Request for DTC IPO Tracking System tracking is prohibited.

- There may be IPOs with only Limit Priced Secondary Shares Offers.

- The settlement input to CB's clearing broker must include fields with
  the amount to be deducted from the payment to cover the Gross
  Underwriting Spread attributable to the shares sold by the respective
  sell order. Said amount is part of the UTC input.

- For context: The actual shares must be held in the seller's
  Broker-Dealer's clearing broker's DTC free account. Short selling is
  not allowed, but not checked by PDP.

### [Lead Manager Short]{.mark}

> **Lead Manager Short Shares** - also referred to as **LM short**, or
> **over-allotment/greenshoe short** - represent shares that the lead
> manager sells in excess of the participating primary and secondary
> offered quantities. These shares are not part of the actual shares
> issued by the company or sold by existing holders at the time of the
> offering but are instead sold short by the lead manager as part of the
> allocation process.

- Quantity entered and updated by the Lead Manager.

- Updated intraday, but only before Lead Manager submits the final
  Offering Price and trigger the Allocation Process.

- PDP keeps a log with date and time of updates.

- LM Short Quantity [is not]{.underline} included in the Indicative
  Clearing Price or Clearing Price calculation.

- LM Short Quantity [is]{.underline} included in the quantity allocated
  to investors.

- LM Short Quantity is normally not used for Bonds, but the PDP does not
  prevent/prohibit this.

## [Order Cancellation]{.mark} and [Modification]{.mark}

> **Orders can generally be cancelled or modified at any time during the
> Bidding Period with the following exceptions**:

1.  The Seasoning Period, which is period of time following Order entry
    or modification as further described in section 3.9.

2.  Bids entered or modified at or after the later of **Regular Trading
    Close Time** (default 4:00 pm EST) and **SEC Effectiveness Delay**
    after **SEC Effectiveness Date/Time** (the latter only if SEC
    Effectiveness is not N/A).

3.  Preferential Bids and Limit Priced Secondary Shares Offers cannot be
    cancelled or modified at or after SEC Effectiveness Delay after SEC
    Effectiveness Date/Time.

> Order Modification within PDP refers to change of price, quantity
> and/or MinQty. Other modifications are not considered an Order
> Modification and hence allowed without restriction, for example a
> change of client Account.
>
> For the avoidance of doubt, order modification restrictions do not
> apply to Lead Manager managed Offer quantities.

## [Halted Orders]{.mark}

Orders are assigned the Halted Order condition when:

1.  a Broker-Dealer with active orders is added to the Exclude List of
    the new issue security

2.  a Broker-Dealer's QSR Active attribute is changed to False

3.  manually assigned in the MarketOps app

The Halted Order condition is not communicated via FIX but shows in
BidMan. Halted Orders can be cancelled but not modified.

Halted Orders are kept in PDP, but do not participate in any actions,
i.e., are removed from order depth, indicative clearing price
calculation, market data, listing exchange minimum monitoring, Priority
Group Fill Simulation, etc.

If/when the cause of the halt is resolved prior to transition to the
Pricing state, then Halted Orders are re-activated keeping the time
priority they had when halted.

At transition to the Pricing state, before any allocations are
determined, all Halted Orders are cancelled. The cancellation is
confirmed via FIX messages when applicable.

When there are Halted Orders, the MarkOps app displays a reminder every
5 minutes (system level configurable) to resolve the cause or cancel the
Orders.

## [Price Collars -- Maximum Price Allowed and Minimum Price Allowed]{.mark}

> Price Collars consist of **Maximum Price Allowed** and **Minimum Price
> Allowed** and are used for rejection of new or modified orders with a
> price outside the Price Collar range. Existing orders in the orderbook
> are kept when the Price Collars are changed and therefore might have a
> price outside the then current Price Collars.
>
> Orders outside the Price Collars are rejected with an explanation.
>
> See also Hight Price Range, Low Price Range and Minimum Acceptable IPO
> Price which are related to but not part of the Price Collars. Relation
> between the different price thresholds:

- Maximum Price Allowed (e.g., \$45.00)

  - High Price Range (e.g., \$25.00)

  - Low Price Range (e.g., \$20.00)

  - Minimum Acceptable IPO Price (e.g., \$15.00)

- Minimum Price Allowed (e.g., \$12.50)

> Price Collars and all the above thresholds are on orderbook level and
> configurable intraday by the Lead Manager.
>
> [Price Collars for Bonds are expressed in yield and hence the relation
> between the terms is reversed:]{.mark}

- [Minimum Yield Allowed (e.g., 3.500%)]{.mark}

  - [Low Yield Range (e.g., 4.000%)]{.mark}

  - [High Yield Range (e.g., 4.800%)]{.mark}

  - [Maximum Acceptable IPO Yield (e.g., 4.999%)]{.mark}

- [Maximum Yield Allowed (e.g., 5.500%)]{.mark}

> **For context:**

- An IPO's final Offering Price might be set as high as 80% above the
  final High Price Range. The High Price Range might be adjusted upwards
  during the Bid Period. Therefore, the Maximum Price Allowed is set
  high above the initial High Price Range to avoid rejecting Bids that
  later might well be within 180% of the final High Price Range. The
  goal is to allow for genuine price discovery while still filtering out
  orders that are unreasonably high.

- The Minimum Price Allowed, which is set just below the low end of the
  Minimum Acceptable IPO Price to block unrealistically low bids.

- The difference between the minimum and maximum is that the latter
  contribute to price discovery as it is expected to execute, whereas
  the former lower will not participate and hence there is no cost for
  placing misleading orders.

## [Price Range -- High Price Range & Low Price Range]{.mark}

> **Definition of \"Price Range,\" \"High Price Range,\" and \"Low Price
> Range\"**

- **Price Range** refers to the expected range of prices for an
  offering, with the initial low and high value (the \"Low Price Range\"
  and \"High Price Range\") specified in the Preliminary Prospectus.
  This range is set by the Lead Manager on behalf of the issuer and may
  be adjusted intraday during the Bid Period. The Price Range is
  intended to guide investors on where the offering is likely to be
  priced and to facilitate transparent price discovery during the
  auction process.

- **High Price Range** is the upper bound of the currently expected
  price range for the offering.

- **Low Price Range** is the lower bound of the currently expected price
  range for the offering.

> **How Price Range is Used in the Process**

- The current Price Range is published on the offering information page,
  BidMan and BatMan. It can be revised during intraday during the Bid
  Period.

## [Minimum Acceptable IPO Price]{.mark}

> The **Minimum Acceptable IPO Price** is the lowest price at which an
> issuer is willing to sell the securities.
>
> The Minimum Acceptable IPO Price is discretionary decide by the issuer
> but often the initial Minium Acceptable IPO is set by the Lead Manger
> to the **lower end of the Preliminary Prospectus price range minus 20%
> of the upper end** of that range.
>
> *Example*: If the price range is \$16--\$20, the calculation is \$16 -
> (20% × \$20) = \$16 - \$4 = **\$12**.
>
> The **Minimum Acceptable IPO Price** intraday updateable.
>
> **Purpose**:

a.  Serves as a **fallback value** for the Indicative Clearing Price
    (Cprice) when there is insufficient demand to cover the minimum
    offering size at or above this price\[1\]\[3\].

b.  If the current minimum offering size (Primary securities +
    Up/down-size + Secondary securities via Lead Manager + Lead Manager
    Short) is not fully subscribed at the Minimum Acceptable IPO Price,
    the IPO is **postponed or canceled**. The postponement or
    cancellation is not automatic by the PDP, but PDP prohibits
    transition to the Close Pending state and to the Pricing state.

## [Seasoning Periods]{.mark} 

> **Seasoning Periods** in the ClearingBid auction process are defined
> as mandatory waiting times after an Order is entered or modified,
> during which the Order cannot be withdrawn/canceled, or further have
> its price, quantity or MinQty modified. The purpose of these periods
> is to prevent manipulative or disruptive bidding behavior and to
> ensure the integrity of the price discovery and allocation process.
> Other modifications such as, e.g., Account, Preferential Bid flag,
> Retail Bid flag, Institutional Bid flag, tracking Flag, ClearingFirm
> and ClearingAccount are allowed.
>
> **The duration of a Seasoning Period is determined by the Time Window
> associated with the Order's PDP assigned timestamp. Said timestamp is
> assigned by the PDP limit orderbook sequencer at ingestion of a new
> order and updated at any modification of its price, quantity or
> MinQty.**
>
> **[Seasoning Period]{.mark}**: Each Seasoning Period, specified as
> hh:mm, is a time period during which an Order cannot be cancelled or
> have its price, quantity or MinQty modified.

- [If the specified Seasoning Period is 12:00 or greater, then this is
  the real-time clock time]{.mark}

- [If the specified Seasoning Period is less than 12:00, then the time
  only counts down during Trading Hours]{.mark}

- [**Default**: Three (3) Time Windows (Time Windows can be deleted and
  added)]{.mark}

  - [**Time Window 1**: From 00:00:00 on the date of Open state;
    Seasoning Period = 24:00]{.mark}

  - [**Time Window 2**: From 00:00:00 on the date of Pricing -1;
    Seasoning Period = 4:00]{.mark}

  - [**Time Window 3**: From 00:00:00 on the date of Pricing; Seasoning
    Period = 0:05]{.mark}

> **Time Windows:**

- [**Start time** is specified as YYMMDD-HH:MM:SS and can be
  modified]{.mark}

  - Note the start time can be at any time but actual Order entry or
    modification is only available during Trading Hours when in the
    orderbook is in Open or Close Pending state.

- [When adding Time Windows, the default start time is 00:00:00 on the
  day before the Time Window with the next higher numberFor example,
  adding a fourth Time Window to the above would change Time Window 3
  default start to Pricing -1 and Time Window 2 to Pricing -2.]{.mark}

- [Time Window 1 always has a default start of 00:00:00 on the date of
  the scheduled transition to Open.]{.mark}

- [The last Time Window always has a default Seasoning Period =
  0:05]{.mark}

- [The next to last Time Window always has a default Seasoning Period =
  4:00]{.mark}

- [All other Time Windows have a default Seasoning Period =
  24:00]{.mark}

- [When there is only one Time Window, the default Seasoning Period =
  0:05]{.mark}

## [Priority Groups]{.mark}

When using the **Priority Group Pro-rata** allocation method (aka time
priority pro-rata), **Competitive Bids are assigned Priority Group
depending on its associated Time Window**.

Priority Groups with lower numbers receive a relatively larger
allocation than Priority Groups with higher numbers because said Bids
start contributing to the transparent demand and price discovery process
earlier in the marketing period. Example with 3 groups:

  ------------------------------------------------------------------
  Priority Group Time Window 1 Allocation ≥ Priority Group 2
  1                            
  -------------- ------------- -------------------------------------
  Priority Group Time Window 2 Allocation ≥ Priority Group 3
  2                            

  Priority Group Time Window 3 Remaining quantity after Group 1 & 2
  3                            allocations
  ------------------------------------------------------------------

For further details see section 7.3 Priority Group Pro-rata (aka time
priority pro-rata).

# Limit Orderbook

## [Indicative Clearing Price aka Cprice]{.mark}

> **Indicative Clearing Price: Definition and Process**
>
> The **Indicative Clearing Price** is a real-time, market-derived price
> that reflects the current equilibrium between buy (bid) and sell
> (offer) orders during the auction phase of a new issue or IPO on the
> ClearingBid platform. It is continually calculated and published
> during the Bid Period, providing transparency to all.
>
> **Key features and calculation:**

- **Real-time calculation:** The Indicative Clearing Price is updated
  continually as new orders are entered, modified, or canceled during
  the auction process. Note that the Indicative Clearing Price is
  continually, not continuously, calculated. The frequence of the
  calculation is per a system parameter which typically is set to 1
  second or less (default 1s), but could be larger.

- **Market-driven:** It is based strictly on the current order book at
  time of calculation. There are two different calculations depending on
  whether [all securities are offered at the Minimum Acceptable IPO
  Price (Single-side Auction)]{.mark} or if there are [multiple offers
  at different minimum price levels (Two-sided Auction)]{.mark}.

- **Indicative nature:** The Indicative Clearing Price is *indicative*
  until the auction closes; it is not binding until the end of the Bid
  Period, at which point the **Clearing Price** is set to the final
  Indicative Clearing Price. The Clearing Price determines the maximum
  permissible Offering Price.

- **Relation to Offering Price:** The final Offering Price, set by the
  issuer and lead manager, cannot exceed the Clearing Price but may be
  set lower, depending on market conditions and other factors.

### [Cprice/Cyield]{.mark} 

> **Cprice/Cyield: Definition and Process**

- The **Cprice** is simply the larger of the Minimumum Acceptable IPO
  Price and the Indicative Clearing Price. The purpose of the Cprise is
  to provide investors with the current best indication of what the
  final Clearing Price might be.

  - For **Stocks and Bonds**, the Cprice calculated based on the
    orderbook.

  - For **FundETFs**, the Cprice is an externally provided Net Asset
    Value (NAV).

- **Transparency:** The Cprice is published on the ClearingBid website
  and made available to participants and the public, enhancing price
  discovery and market transparency.

- The **Cprice is disseminated in real-time via ClearingBid\'s public
  API, FIX market data, BidMan and BatMan.**

- The **[Cprice publishing method is configurable]{.mark}** and can be
  set to one of several modes:

  - **[Min IPO Price]{.mark}**: The Cprise is published from the start
    of the Bid Period, i.e., when the full offering quantity is not
    subscribed at the Minimum Acceptable IPO Price, said price is
    published. At times when the full quantity is subscribed for at the
    Minimum Acceptable IPO Price or higher the published Cprise is the
    Indicative Clearing Price.

  - **[Not Published + Auto On]{.mark}:** No Cprice is published
    initially; publishing is automatically turned on when the full
    quantity is subscribed at the Minimum Acceptable IPO Price.
    Publishing stays on even if the full quantity falls back below the
    Minimum Acceptable IPO Price.

  - **[Not Published + Manual On]{.mark}:** No Cprice is published
    initially; publishing is manually turned on by the Lead Manager by
    selecting "Publish Cprise NOW". If publishing is changed to this
    method after it has been turned on, then it is turned off.

  - **[Publish Cprice NOW]{.mark}:** Cprice publishing is turned on.

### [Single-side Auction]{.mark} 

> A **Single-side Auction** is an auction format where all securities
> are offered conditioned only by the Minimum Acceptable IPO Price. The
> offering can consist of Primary Securities and/or Secondary
> Securities. The participating Secondary Securities can be offered by
> one or multiple sellers.

- **Order Book Structure:** In a Single-side auction, the Orderbook
  consists only of buy orders (bids) from investors and a total sell
  quantity comprising all shares offered (see section 3.3) conditioned
  by the Minimum Acceptable IPO Price.

- **Price Discovery:** The (Indicative) Clearing Price is the highest
  price at which all offered shares can be sold to investors. The
  process is similar to a Dutch auction, where the price is set so that
  the entire offering is subscribed. The (Indicative) Clearing Price is
  determined by traversing the order book in price-time priority,
  accumulating volumes from the highest bid down until the offering size
  is met excluding any bids where the Minimum Quantity condition cannot
  be met.

- **Offering Price:** The issuer and underwriters may set the final
  Offering Price at or below the Clearing Price, but not above it.

### [Two-sided Auction]{.mark}

> In a Two-sided Auction there are multiple sell offers allowed on
> multiple price levels. On the sell side there is one Main Offering
> order (i.e., the total of Primary Shares, Up/down-size of Primary
> Shares and Committed Secondary shares) with an implied price at the
> Minimum Acceptable IPO Price (which might include Committed Secondary
> Shares and up/down-sizing), plus a number of individual Limit Priced
> Secondary Shares Offers.
>
> **Indicative Clearing Price and Clearing Price**
>
> The Clearing Price[^3] is the price which generates the maximum number
> of shares crossed. If multiple price levels have the same maximum,
> then the Clearing Price is the price level with the lowest imbalance.
> If there are multiple price levels with the lowest imbalance, then the
> Clearing Price is determined by the Imbalance Pressure. The Imbalance
> Pressure is the accumulated buy quantity minus the accumulated sell
> quantity for the price level. The Clearing Price is determined as
> follows:

(i) The price level that maximizes the number of participating shares;

(ii) if (i) produce more than one price, then select the price level
     that minimize the absolute imbalance;

(iii) if (ii) produce more than one price, then use the imbalance
      pressure to select the price:

(iv) if all imbalances on the in (ii) identified price levels are on the
     sell side, i.e., the sell quantity is larger than the buy quantity,
     then use the lowest price (because there is higher pressure from
     sellers and hence the most favorable price for the buyer prevails);

(v) if all imbalances are on the buy side, then use the highest price;

(vi) if the pressure are on both the buy and sell side, use the rounded
     two decimal mid-price between the highest buy pressure price level
     and the lowest sell pressure price level (for yield order books,
     round to three yield decimals)

(vii) if all of the lowest imbalance levels on the price levels
      identified in (ii) are neutral, i.e., zero, then use the midpoint
      of the in (ii) identified price levels rounded to two decimal
      points (for yield order books, round to three yield decimals)

> **Pricing Simulation Additions**
>
> The Limit Priced Secondary quantity at the Offering Price shall also
> be included in the available quantity.
>
> If the Offering Price is set below the Clearing Price the available
> quantity might be smaller than the quantity at the Clearing Price
> level due to the Offering Price being lower than the Limit Price of
> some or all Secondary Shares Sell Orders.
>
> To inform the Lead Manager about the Offering Price impact on number
> of shares sold the simulation during the Pricing Process also returns
> a table with available quantity by price level at the Clearing Price
> and below.
>
> When run before the start time of the lowest Priority Group, the
> Offering Price Simulation function and the Quantity by Level
> Simulation function treats all orders as Priority Group 1, i.e., the
> allocation is first to Preferential Bids and then to all Competitive
> Bids as if there were only one Priority Group. The simulation should
> include the impact of MinQty bids.
>
> **Allocation Algorithm Additions**
>
> Limit Priced Secondary Shares Offers are not participating Pro-rata,
> they are filled in strict price-time priority.

## [BidMan app -- Bid Management app for Broker-Dealers]{.mark}

> This section is an overview only. For details on the BidMan app see
> the document "PDP MVP v\<n n\>.pdf/docx".
>
> **Broker-Dealer Orderbook management, view and information**
>
> Broker-Dealers (BDs) mainly interact with ClearingBid\'s Price
> Discovery Platform (PDP) through their existing Order Management
> Systems (OMS). The OMS interact with PDP using the FIX protocol.
> BidMan is an app complementing the OMS. The BidMan interface provides
> order management capabilities, real-time visibility into new-issue
> auctions, and market data. Below is a synthesized description based on
> the provided documents, followed by identified contradictions.
>
> **Key Features of the BidMan app**

- **Order Management**:

  a.  BDs can enter, modify, and cancel **Good-'till-Canceled (GTC)
      limit orders** for stocks, bonds, or funds/ETFs. For stocks and
      bonds only limit order bids with GTC are allowed. For funds/ETFs
      also Market-at-Close orders are allowed.

  b.  Order entry and modification include a **DTC Tracking checkbox**
      with info pop-up that settlement will be by the Lead Manager, not
      via NSCC-CNS and FED NSS, when checked. When checked a
      retail/institutional switch must also be set.

  c.  Price for bonds/fixed income are entered and expressed as yield.

  d.  Orders are displayed with details: **Limit Price, Quantity,
      Timestamp, Minimum Quantity (MinQty)** condition and
      **Status**(active/cancelled/halted/filled). The display can be
      expanded to cover all order attributes.

- **Market Data and Analytics**:

  a.  **Real-time order book**: Displays price levels, bid sizes,
      accumulated volume, and key metrics:

      i.  **Cprice:** Visualized via price level price displayed in
          green

      ii. **Low/High Price Range**: Visualized via green lines

      iii. **Priority Group Metrics**: \"My\" columns (e.g., MyCLRPOrd,
           My LowR.Size) show BD-specific data like orders/shares at
           Cprice or price ranges.

  b.  [**Aggregation and Display Freeze Functions**: BDs can group price
      levels or freeze the display for analysis]{.mark}

- **Workflow Integration**:

  a.  **FIX Compatibility**: Orders can be entered via FIX or manually
      through BidMan. MinQty orders may be restricted to BidMan-only
      entry per instrument rules.

      i.  Orders entered via OMS FIX can be displayed and cancelled via
          BidMan

      ii. [Further integration is desired, for example update Account
          ID, but any implementation will not be decided until potential
          complexities are investigated]{.mark}

  b.  **Allocation Visibility**: Post-execution, BDs see execution
      details (e.g., fill quantity, price) per order.

- **BD Users and Permissions**:

  a.  [Admin]{.mark} -- Each BD firm must have at least one Admin user
      which is a user that can create/manage users and [Firm
      Accounts]{.mark}[^4] within the BD firm. The Admin cannot place
      orders or view orderbooks.

  b.  [Broker]{.mark} -- a Broker can enter, manage and view orders and
      data belonging to Firm Accounts for which he is assigned write
      privileges and can view orders and data of Firm Accounts for which
      he is assigned read privileges.

  c.  [OMS User]{.mark} -- the login identity and credentials for an OMS
      System. An OMS user can only be associated with one Firm Account.

  d.  [Group Manager]{.mark} -- a Group Manager are assigned read
      privileges for several Brokers and/or Group Managers. The Group
      Manager can view orders and aggregated information for all users
      within the group. A Group Manager can optionally be assigned
      cancellation and/or write privileges for orders belonging to the
      group.

  e.  [Master Broker]{.mark} -- a Master Broker can view orders and
      aggregated information for all users within the firm. A Master
      Broker can optionally be assigned cancellation and/or write
      privileges all orders within the firm.

## [BatMan app -- Bidding and Auction Management app for Lead Managers]{.mark}

> This section is an overview only. For details on the BatMan app see
> the document "PDP MVP v\<nn n\>.pdf/docx".
>
> **Lead Manager Orderbook View and Information**
>
> The **Lead Manager (LM) Orderbook View** in the ClearingBid platform
> is a specialized interface designed for underwriters (Lead Managers)
> to monitor, manage, and allocate new security offerings. This view
> provides comprehensive, real-time access to order flow, pricing, and
> allocation tools, supporting the LM's central role in the price
> discovery and distribution process.
>
> **Key Features and Information Available in the Lead Manager Orderbook
> View**

- **Access and Permissions**

  - LMs use BatMan to create, monitor, and allocate new offerings.

  - LMs have access to overview orderbook data, including both
    aggregated and per-broker-dealer breakdowns, as well as allocation
    simulation and execution tools. LMs does not have access to
    individual orders, only total quantity and number of orders per
    broker-dealer and price level.

- **Orderbook Display**

  - **Aggregated and Per-Broker-Dealer Data:**LMs can see accumulated
    and aggregated volumes for all price levels at all times, including
    breakdowns by broker-dealer and selling concession by
    broker-dealer\[1\]\[2\].

  - **Real-Time Price Discovery:** The view displays real-time
    indicative clearing prices (CLRP), average weighted prices (aWAP,
    rWAP), and accumulated order sizes at key price points (low/high
    price range, clearing price)\[1\]\[2\].

  - **Orderbook States:** The orderbook transitions through well-defined
    states with LMs able to monitor and, in some cases, request state
    changes.

- **Offering and Order Attributes**

  - LMs can view and edit (where permitted) offering attributes such as
    symbol, market (i.e, listing exchange), issuer, asset class, bid
    period dates, price ranges, minimum/maximum order sizes, allocation
    algorithms, and more\[1\]\[4\]\[2\]\[3\].

  - **Allocation Simulation:** LMs can simulate allocation outcomes by
    adjusting offering price, quantity, and secondary shares, and see
    pro-rata fills, demand by priority group, and unfilled statistics
    before submitting final allocations\[1\]\[2\].

  - **Allocation Algorithms:** Supported algorithms include Pro-rata,
    Priority Group Pro-rata, Price-Time Priority, and High Price Before
    Time Priority.

- **Closing and Pricing Process**

  - The Closing process (i.e., Close Pending state) is designed to
    minimize manipulation risk by using random time windows and
    stability checks on price and volume.

  - During the Close Pending and Pricing state the LM can simulate the
    outcome of different Offering Prices and, for Priority Group
    Pro-rata Allocation, different Fill Rates. This can also be done
    during the Open state but is only relevant close to the end of said
    state.

  - During the Pricing state and when satisfied with the simulation, the
    LM submit the final Offering Price (and Fill Rates for Priority
    Group Pro-rata Allocation).

- **Reporting and Export**

  - LMs and Market Operations can export detailed reports (e.g., CSV)
    from the orderbook, including allocations, commissions, and fill
    percentages for reconciliation and compliance purposes\[5\].

  - The system supports filtering and exporting data by symbol, broker,
    priority group, and more, with configurable thresholds to prevent
    performance issues\[5\].

<!-- -->

- **Lead Manager Users and Permissions**:

  - [Admin]{.mark} -- Each firm acting as a Lead Manager must have at
    least one Admin user which is a user that can create/manage
    Underwriter users. The Admin cannot manage new issues or view
    orderbooks.

  - [LM aka Lead Manger]{.mark} -- a user that can create, monitor, and
    allocate new offerings.

  - [Issuer]{.mark} -- a user that have read permissions to everything
    of the LM.

> **Information Uniquely Available to Lead Managers**

- **Full Market Depth:** LMs see the entire orderbook, not just their
  own firm's orders, but only to a detail level of total quantity and
  number of orders per price level and broker-dealer.

- **Broker-Dealer Breakdown:** LMs can number of orders and allocations
  by broker-dealer, including selling concession
  calculations\[1\]\[2\]\[5\].

- **Allocation Tools:** LMs have access to simulation tools to test
  different allocation scenarios before final submission\[1\]\[2\].

- **Commission and Fee Calculations:** LMs can view and verify selling
  concessions and gross spreads per broker-dealer, ensuring the total
  matches the expected value\[5\].

### [Broker-Dealer Exclude List (Blacklist)]{.mark}

> The **Exclude List** functionality allows the Lead Manager to
> selectively restrict specific Broker-Dealers from participating in an
> IPO auction. Below are a detailed description and analysis of
> contradictions in the documents.
>
> **Purpose and Mechanism**:
>
> The Lead Manager can **block specific Broker-Dealers** from entering
> orders for an offering or accessing market data.
>
> This is implemented via the **\"Exclude Broker-Dealers\" function** in
> the Offering Widget in the BatMan application. The list contains all
> PDP participating Broker-Dealers and on the top, there are two radio
> buttons - Block All and Block None. Block All put checks the box in
> front of all Broker-Dealers and Block None uncheck all boxes. Default
> is Block None. The Lead Manager can check/uncheck individual boxes.
>
> **When blocked, Broker-Dealers cannot:**

- Submit new orders or modify existing orders.

- Receive market data (e.g., indicative clearing price, order book
  updates).

- The feature operates as a **blacklist (exclude list)**, not a
  whitelist. There is no explicit \"include list\".

> **Configuration**:

- Excluding/Blocking is configured per offering and can be applied
  intraday.

- If a Broker-Dealer has Orders in the Orderbook when excluded the Lead
  Manager gets a warning and is asked to confirm to proceed, and if
  proceeding the affected Orders will be assigned the Halted Orders
  state.

## [Orderbook States]{.mark}

> The Orderbook States define the lifecycle of an offering within PDP.
>
> The PDP manages offerings through distinct states, each governing
> specific operational rules (e.g., order entry, cancellations,
> allocation procedures). The core states include:

1.  **NEW**

    a.  Active until the offering announcement.

    b.  Orderbook is only visible to the Lead Manager firms and Market
        Operations.

    c.  Orderbook is empty. No order entry is allowed.

    d.  Manual transition from Frozen by MarkOps is possible but only
        after MarkOps has cancelled all orders.

2.  **UPCOMING**

    a.  Normally active from announcement until transition to Open (Bid
        Period start). Orderbook publicly available.

    b.  Normally the orderbook is empty. No order entry allowed.

    c.  In case state has transitioned manually from Frozen by MarkOps,
        then orders might persist and can be cancelled, but not modified
        or entered.

3.  **PRE-OPEN**

    a.  Orderbooks in the OPEN state transition to PRE-OPEN at the
        Trading Day End Time of Trading Days except on the scheduled day
        of closing/Pricing. There is no other automatic transition to
        the PRE-OPEN state but MarkOps can manually make the transition.

    b.  All order entry/modifications/cancellations are prohibited.

    c.  The state is active during the Bidding Period from Trading Day
        End Time until next Trading Day Open Time at which point the
        orderbook state transitions back to the OPEN state.

4.  **OPEN**

    a.  Active during Trading Hours during the Bid Period. The initial
        transition to the OPEN state for an orderbook occurs at the
        Trading Day Start Time of the set Auction Opening Day.

    b.  Orders can be entered, modified, or cancelled during Trading
        Hours subject to restrictions described in section 3.4 [Order
        Cancellation]{.mark} and [Modification]{.mark}.

    c.  Transitions to PRE-OPEN at the Trading Day End Time of Trading
        Days. For the avoidance of doubt, the transition is applicable
        only to orderbooks in the OPEN state, not when in the CLOSE
        PENDING state.

    d.  Transitions to **CLOSE PENDING** at the set Closing Process
        Start Time for the orderbook.

    e.  Manual transition to FROZEN, HALTED or PRE-OPEN by MarkOps is
        possible.

5.  **CLOSE PENDING**

    a.  Close Pending triggers the start of the Closing Window and is
        part of the Bidding Period. For details of Closing Window see
        the section 5.

    b.  Transition to Close Pending is prohibited if the current minimum
        offering size (Primary securities + Up/down-size + Secondary
        securities via Lead Manager + Lead Manager Short) is not fully
        subscribed at the Minimum Acceptable IPO Price.

    c.  If SEC Effectiveness is [not set to N/A]{.underline}, transition
        to Close Pending is prohibited before SEC Effectiveness Delay
        after SEC Effectiveness Date/Time.

    d.  If SEC Effectiveness is null, i.e. not yet set, at Regular
        Trading Close Time or at Scheduled transition to Close Pending,
        then the LM and MarkOps must be notified that the SEC
        Effectiveness must be set for Close Pending to be initiated.
        Close Pending start must be delayed until SEC Effectiveness is
        set.

    e.  If SEC Effectiveness is [not set to N/A]{.underline} and
        Preferential Bids Allowed is [not set to No]{.underline} then
        the Lead Manager must confirm that the exchange listing minimum
        requirements can be fulfilled before transition to Close Pending
        (there is an exchange listing minimum requirements fulfillment
        monitoring and simulation tool in BatMan).

    f.  Closing Window countdown begins.

    g.  Manual transition to Frozen or Halted by MarkOps is also
        possible.

6.  **PRICING**

    a.  Triggered automatically within the Closing Window that includes
        price/volume stabilization checks in accordance with the
        automated Closing Process (see section 5). Can also be triggered
        manually by MarkOps.

    b.  Transition to PRICING is prohibited if the current minimum
        offered quantity (see section 3.3) is not fully subscribed at
        the current Minimum Acceptable IPO Price, i.e., the Indicative
        Clearing Price is lower than the current Minimum Acceptable IPO
        Price.

    c.  Order entry and cancellations are prohibited.

    d.  Primary, primary up/down-size and secondary securities
        quantities change prohibited.

    e.  Snapshot of orderbook for potential future restore enabling
        correction or errors and rerun of pricing and downstream
        processing.

    f.  Final allocation simulation and calculations performed by Lead
        Manager. Lead Manager Short updates allowed.

    g.  When satisfied with the simulation result, the Lead Manager
        submits final Offering Price, Lead Manager Short and for Fill
        Rate per Priority Group (when applicable) which triggers
        generation of an Allocations. Lead Manager Short updates are
        prohibited after said submission. Note that the Allocation
        Process only assigns/distributes shares among eligible orders.
        The actual trades are generated by the Execution state.

    h.  After allocation generation an Allocation Verification Process
        is executed producing an Allocation Verification Report which is
        analyzed by the Lead Manager and MarkOps

        i.  if the verification is successful the MarkOps releases the
            allocations for execution (the PRICING state persists until
            start of EXECUTION state), or

        ii. if the verification is unsuccessful MarkOps triggers
            transition to the MODIFICATION state.

7.  **EXECUTION**

    a.  Within the EXECUTION state, executions are triggered:

        i.  Immediately for securities with the Execution parameter set
            to "Auction Close Day Executions", or

        ii. In the early morning of the Trading Day following the
            Auction Close for securities with the Execution parameter
            set to "Next Trading Day Executions". The transition is at
            the time set by the system wide parameter "Next Trading
            Day Execution Time" which reflects the earliest time CB's
            clearing broker and NSCC UTC can ingest trade confirmations.

    b.  The system generates:

        i.  Trade executions (fills and partial fills) and sends
            corresponding FIX messages and drop copies to participants
            and optionally participant's clearing firms (drop copies).

            1.  Securities with Next Trading Day Executions and the
                SameDaySettlmnt parameter checked include FIX tag 63
                "SettlmntTyp" with a value of 1 (1=cash). Securities
                with SameDaySettlmnt not checked use standard settlement
                cycle for the security, i.e., does not include FIX tag
                63.

            2.  [Trade executions for orders requested to be tracked by
                the DTC IPO Tracking System or the DTC Institutional
                Delivery System and delivered to the respective
                Broker-Dealer's clearing firm's DTC IPO Control
                Account]{.mark}[^5] [must include FIX tag 439
                "ClearingFirm" populated with the NSCC number of the
                Broker-Dealer's clearing firm and *FIX tag 440
                "ClearingAccount" indicate the DTC clearing account
                related to the Bid*.]{.mark}

        ii. A Trade Execution file for orders that are not requested to
            be tracked by the DTC IPO Tracking System or the DTC
            Institutional Delivery System. The file is uploaded to CB's
            clearing broker for further delivery to NSCC Universal Trade
            Capture (UTC) for further Continuous Net Settlement (CNS)
            processing and DTC/FED NSS net settlement.

        iii. [A Trade Execution file with all executions for Bids
             requested to be tracked by the DTC IPO Tracking System. The
             executions in the file are sorted by Broker-Dealer. The
             file also contains quantity sub-totals per Broker-Dealer.
             The file is sent to the Lead Manager for further settlement
             processing by the Lead Manager.]{.mark}

        iv. [A Trade Execution file with all executions for Bids
            requested to be tracked by the DTC Institutional Delivery
            System. The executions in the file are sorted by
            Broker-Dealer and within the Broker-Dealer by
            ClearingAccount (FIX tag 440 which can also be
            entered/modified in BidMan). The file also contains quantity
            sub-totals per Broker-Dealer and within the Broker-Dealer by
            ClearingAccount. The file is sent to the Lead Manager for
            further settlement processing by the Lead Manager.]{.mark}

        v.  [One Trade Execution file for every Broker-Dealer with
            executions for Bids requested to be tracked by the DTC IPO
            Tracking System. The files are sent to the respective
            Broker-Dealer.]{.mark}

        vi. [One Trade Execution file for every Broker-Dealer with
            executions for Bids requested to be tracked and delivered
            through the DTC Institutional Delivery System. The
            executions in the file are sorted by ClearingAccount (FIX
            tag 440 which can also be entered/modified in BidMan). The
            files are sent to the respective Broker-Dealer]{.mark}

    c.  State persists until:

        i.  CB's clearing broker confirms acceptance of all settlement
            transactions by NSCC UTC and no other abnormal issues have
            been identified. MarkOps then manually confirm that Trade
            Execution is successful which triggers release of
            cancellations of unfilled or partially unfilled orders, or

        ii. in case of any issue, MarkOps trigger change to the
            Modification state.

8.  **CLOSED**

    a.  Final state after successful generation of all trade executions
        and cancellations. Orderbook empty.

9.  **CANCELLED**

    a.  Manually assigned; cancels all orders.

    b.  Irreversible final state.

**Substates**

Two substates modify **OPEN** and **CLOSE PENDING**:

- **FROZEN** -- Prohibits all order entry/modifications/cancellations.
  Freezes Closing Window clock.

- **HALTED** -- Prohibits order entry and modifications but allows
  cancellations. Freezes Closing Window clock.\
  *Source: \[3\]\[4\]\[6\]\[8\]*

Substate **MODFICATION** is used to go back from **EXECUTION** and
restart **PRICING:**

- **MODIFICATION**

> A temporary state to deal with any issues discovered, for example:

- Unsuccessful offline verification of allocations

- Desire to change Offering Price or Fill% per priority group discovered
  during offline verification

- Bust of trade executions necessary due to defaulting participant or
  participant withdrawal of QSR relationship with CB's clearing broker.

- Bust necessary due to system/programming error

- Bust necessary due to CB or Lead Manager error

- When bust is required, the orderbook is restored to the state at the
  beginning of the PRICING state. Any corrections to orders applied.
  Reversals for trade executions and NSCC UTC input generated and
  disseminated.

> When necessary modifications has been performed MarkOps triggers
> transition to the PRICING state.
>
> [A set of tools must be developed to speed up management of predicted
> modification scenarios like orderbook restore, generation of reversals
> for order management systems and NSCC UTC, etc.]{.mark}

**State Transition Rules**

- **Manual Changes**:

  - Market Operations can manually trigger allowed state changes except
    to **EXECUTION** and to/from **CLOSED** or **CANCELLED**.

  - Transitions to **NEW**/**UPCOMING** from **FROZEN** require order
    cancellations.

- **Resetting Process**:

  - Moving from **CLOSE PENDING** to **OPEN** aborts the Closing Process
    and a new Closing Process Start Time must be set.

### Allowed Order Book State Transitions

![](media/image1.emf){width="6.300755686789151in"
height="3.481570428696413in"}

\[***TBD: Tor S. to update this diagram to reflect updates in the
4.4***\]

## [Scheduling]{.mark}

> **Scheduling** in the context of the ClearingBid platform refers to
> the parameters and system mechanisms that define system availability
> and automated state transitions.
>
> **Key Elements of Scheduling**
>
> **Calendar** -- PDP has a plurality of calendars and every orderbook
> are assigned to one calendar. Each calendar has the following
> attributes:

- **Trading Day**: Monday to Friday except Holidays. During these days
  bids are accepted and orderbooks can transition to the OPEN state.

- **Holidays**: A list of specific dates when trading is closed, i.e.,
  Trading Days exceptions.

- **Partial Holidays**: A list of specific dates when Trading Day Early
  End Time overrides the Trading Day End Time.

- **Trading Day Start Time**: The time when applicable orderbooks
  transition to the OPEN state. Default 8:00 AM ET.

- **Trading Day End Time**: The time when applicable orderbooks
  transition to the PRE-OPEN state. Default 8:00 PM ET.

- **Trading Day Early End Time**: A time that overrides the Trading Day
  End Time on Partial Holidays. Default 3:00 PM ET.

- **Trading Hours**: When the term Trading Hours is used in this
  document it refers to the time from Trading Day Start Time to the
  Trading Day (Early) End Time.

> **Initial Calendars:**
>
> Calendar 1 = Equities (normally follows the NYSE/Nasdaq trading
> calendar)
>
> Calendar 2 = Bonds (normally follows the SIFMA trading calendar)
>
> **System Open Day**: A day that is a Trading Day in any of the
> Calendars.
>
> **System Open Time**: The time from which Participants can expect the
> system to be available for logon on System Open Days. However, the
> system is normally available for logon and transactions also before
> this time and outside Trading Hours.
>
> **System Close Time**: The time at which Participants might be logged
> out. However, the system does normally not log out Participants
> between Trading Days.
>
> MarkOps is responsible for entering and maintaining the calendar.

## [Orderbook Data Structure and Processing]{.mark}

> **Determinism and Fault Tolerance**
>
> State Machine Replication should be used for fault tolerance and
> recovery must not drop connected sessions or result in response time
> that cause issues. Low latency is not required.
>
> The orderbook data structure and processing must be designed to manage
> short transaction peaks that are many times normal.
>
> MarkOps commands and internal management events like state changes
> should have a priority queue which must always be empty before next
> order is accepted.
>
> **Orderbook Data Structure**
>
> PDP's orderbook data structure should be optimized for order entry,
> cancel/replace with changed timestamp (price, quantity or MinQty
> updates), fast control of accumulated volume followed by traverse of
> selected price level to eliminate MinQty fallouts. Unlike a continuous
> trading system, PDP does to be optimized for executions or fast
> control of best bid/ask.
>
> Note that the sell side is just one price level for the first couple
> of releases of PDP v2. [The two-sided auction will not be implemented
> until v2.3 or later.]{.mark}
>
> A matching engine expert has proposed the following (however, the
> developer should consider other solutions if they believe such
> solution to be more efficient):

- Red-Black Tree (RBT) for price levels provides consistent fast price
  level access

- FIFO double-linked list (DLL) for orders is efficient for keeping
  track of time priority and managing cancel/replace that update time
  priority (price, quantity and MinQty). In Java an ArrayDeque is
  probably more efficient than an outright DLL.

- Hash Map (HM) for direct access to ClOrdID for cancel and
  cancel/replace that preserve time priority

> Price Tree (RBT)
>
> ├ Price Level (contains FIFO queue)
>
> │ ├ Order 1 (DLL node)
>
> │ ├ Order 2
>
> │ └── \...
>
> Order Lookup Hash Map: order_id → (side, price, DLL node)
>
> **Total Quantity per price level, MinQty and Preferential Bids**
>
> A possible implementation for fast indicative Clearing Price updates
> and listing exchange minimum fulfillment monitoring might be to keep
> following data at the RBT node; (i) total order quantity of the price
> level, (ii) flag indicating if the price level includes MinQty and
> (iii) flag indicating if the price level includes Preferential Bids.
>
> **Halted Orders**
>
> When implementing management of Halted Orders consider that they are
> seldom used.
>
> **[Bonds and Fixed Income]{.mark}**
>
> [Bonds and Fixed Income orderbooks trade on yield, i.e., the orderbook
> is inversed because a higher price when expressed as yield is a lower
> yield. For example, a Bid for 3% is higher than a Bid for 4% because
> when bidding on a bond you are bidding for the expected rate of
> return. To receive the lower return of 3% you must pay a higher
> purchase price for the same rate of return. For example, if you bid 3%
> for a \$1000 par value bond with one year to maturity and zero coupon
> you pay \$1000/1.03=\$970.87 and if you bid 4% you pay
> \$1000/1.04=\$961.54. In other words, a lower yield (higher price) is
> better for the issuer, and a higher yield (lower price) is better for
> the investor (bidder).]{.mark}
>
> [The implementation should consider if it is most efficient to use
> yield for fixed income RBT nodes or a price converted to 100-yield to
> "un-inverse" the orderbook.]{.mark}
>
> [**FIX 4.4 and higher supports price express as yield (PriceType
> \<423\> field = 9)**. However, FIX 4.2 and lower only supports
> PriceType 1 to 3 where 1 is Percentage, often called "dollar price".
> This is not a yield, it is percentage of par value, e.g., a bid for
> 98.55 is a bid to pay \$985.50 per bond for a \$1000 par value
> bond.]{.mark}
>
> \[*Note to ClearingBid: We must probably require that Broker-Dealers'
> OMS supports expressing price in yield and FIX 4.4 or higher for
> participation in fixed income auctions.*\]

# [Closing Process / Close Pending state]{.mark}

The **Closing Process** begins at the transition to the Close Pending
state, ends at the transition to the Pricing State and is designed to
minimize the risk of price manipulation and ensure fairness for all
participants.

**End of Bidding Period**

The Bid Period, during which orders can be placed, modified, or
cancelled, ends at the transition from the Close Pending state to the
PRICING state, i.e., the Closing Pending state is part of the Bidding
Period and the actual Close of the orderbook is when the Close Pending
state ends. No new orders, cancellations or modifications are allowed
after this point (note that Seasoning Periods also restricts
cancellations and modifications during the Bidding Period).

## [Closing Window]{.mark}

> The orderbook Close, i.e., the transition to the PRICING state occurs
> within a Closing Window consisting of two time windows where the
> second time window starts at a randomized time within the first:
>
> **Closing Time Window 1 (CTW1):** Typically starts at 4:10 PM Eastern
> Time and lasts 20 minutes.
>
> **Closing Time Window 2 (CTW2):** A 20-minute window that begins at a
> random time within CTW1.
>
> The system will trigger the transition to the PRICING state within
> CTW2 when there has been no change to the indicative Clearing Price
> (or, for ETFs, no significant change to participating volume) for a
> rolling 10-second period, the "**Stable Period**". If this stability
> does not occur in the first half of CTW2, the transition to the
> Pricing state is triggered at a random time in the second half of CTW2
> or at the first stable 10-second period.
>
> All time periods are configurable. Default time periods are system
> wide per asset class.
>
> **Manual Intervention**: MarkOps can manually halt or freeze the
> closing process in case of technical issues or unusual order book
> activity, or after a request from the Lead Manager. MarkOps can also
> transition back to the Open state which requires setting a new Closing
> Process Start Time. When an issue include Preferential Bids there is
> also a blocking and confirmation mechanism described in section
> 5.2.1.1.
>
> **State Transitions**:
>
> The offering moves through the following states:
>
> **OPEN** → **CLOSE PENDING** (start of CTW1) → **PRICING** (triggered
> within CW2) → **EXECUTION** → **CLOSED**.
>
> Substates like **FROZEN** and **HALTED** allow for additional control,
> such as prohibiting order entry or cancellations\[1\]\[2\]\[3\].
>
> **Key Parameters (Configurable)**
>
> **Closing Process Start** **Time**: Default is 4:10 PM or 4:00 PM.
>
> **CTW1 Duration**: Default 20 minutes.
>
> **CTW2 Duration**: Default 20 minutes\[1\]\[2\]\[3\].
>
> **Stable Period**: Length of the stable interval required to Close.
> Default 10 seconds.
>
> **Volume Change Threshold**: For Funds/ETFs, maximum allowed volume
> change during the Stable Period. Default 10,000 units (shares/bonds).

## [Preferential Bids and Listing Exchange Minimum]{.mark}

> **Preferential Bids** are designed to facilitate compliance with
> exchange listing minimum requirements and are eligible for an
> Allocation of shares ahead of Competitive Bids, in exchange for
> accepting certain restrictions. Any Preferential Bids, or portions of
> Preferential Bids, that do not receive a preferential allocation,
> participate in the auction as Competitive Bids, and are assigned a
> Priority Group according to the time the Bid was entered or last
> modified.
>
> When **Preferential Bids allowed** is [not set to No]{.underline},
> then the following Broker-Dealers with the QSR Active attribute true
> and not on the Exclude List can place Preferential Bids:

- When set to Lead Manager BD Only, then the Lead Manager BD only.

- When set to Yes, then all eligible Broker-Dealers.

> **Minimum Preferential Quantity** (MinPrefSec): Number of shares/bonds
> set aside for allocation to Preferential Bids. Can be updated
> intraday. If MinPrefSec is set too low, then the PDP calculated
> minimum quantity to fulfill the Exchange Listing Minimum requirement
> will be allotted to Preferential Bids. Hence, MinPrefSec will only
> increase the allotted quantity above what is required by the exchange.
>
> The Participant Broker-Dealer is required to populate all Preferential
> Bids from the same investor with the one and same value in the Account
> field (an investor can place multiple Preferential Bids the investor
> is required to place all bids with the same Broker-Dealer).
>
> A Preferential Bid is recognized as such by the CB System when the
> Broker-Dealer has populated the Account field with an investors
> account identification and the ExecInst field (aka Execution
> Instruction) with the value 'p' (Preferential Bid). The Account field
> and the ExecInst field can be populated either through inclusion in
> the FIX message or manually using the BidMan.
>
> Successful bidders that have placed Preferential Bids will be
> allocated preferential Allocation until the larger of the Minimum
> Preferential Quantity or the quantity that satisfy the exchange's
> minimum listing requirements. Preferential Allocation is 'first come
> first served' based on the latest of the bidders last Preferential Bid
> entry time/last Preferential Bid modification time. If the exchange's
> applicable minimum listing requirement include a certain number of
> bidders with a minimum holding larger than one round lot, bidders with
> Preferential Bids for a total quantity exceeding said minimum holding
> will receive Allocation ahead of other Preferential Bids until the
> requirement is fulfilled. For the avoidance of doubt, preferential
> Allocation is per bidder, not per Bid. For Preferential Bids that do
> not receive preferential Allocation, or only receive partial
> preferential Allocation, the remaining quantity of said Preferential
> Bids will be treated as Competitive Bids in accordance with their time
> of entry or last modification.
>
> The PDP supports multiple different Listing Exchange Minimum
> specifications. This example is the first to be implemented:

1.  Market value of unrestricted publicly-held shares at time of
    listing, e.g., \$40,000,000

2.  Number of unrestricted publicly-held shares at time of listing,
    e.g., 1,250,000

3.  Price per share at time of listing, e.g., \$4

4.  Trading volume or holders' minimums at time of listing, e.g.,

    (A) at least 550 total holders and an average monthly trading volume
        over the prior 12 months of at least 1,100,000 shares per month;
        or

    (B) at least 2,200 total holders; or

    (C) a minimum of 450 round lot holders among which at least 50% of
        such round lot holders must each hold unrestricted securities
        with a market value of at least \$2,500

> Requirements of type 1-3 are verified based on the Minimum Acceptable
> IPO Price and the number of offered unrestricted shares. If any of
> requirements 4(A)-(C) can be verified and confirmed before the
> ClearingBid IPO Auction Process, then no further action is required.
> If none of requirements 4(A)-(C) can be verified beforehand, then the
> ClearingBid IPO Auction Process are expected to include Preferential
> Bids to comply with Listing Exchange Minimums 4(B) or 4(C), whichever
> of 4(B) or 4(C) requires the least number of shares allocated to
> Preferential Bids.

### [Preferential Bids Monitoring during the Auction]{.mark}

The PDP tracks:

- Number of Preferential Bids and Bidders

- Number of Preferential Bidders with at least a round lot or large lot

- Quantities required to meet minimum holder or round lot thresholds

- The lead manager's BatMan app displays these metrics in real time and
  allows simulation of different offering prices to assess compliance.

Below is a model process for implementing the example described in
section 5.2 above. The developer can choose another implementation if it
is deemed more efficient.

**Parameters** per security (offering) for Listing Exchange Minimum
management (all changeable intraday):

- RndLot (default 100)

- MinBidQty - Minimum Bid Quantity allowed (default 100)

- MinPrefAlloc - the minimum number of securities assigned in the first
  iteration to each preferential account. MinPrefAlloc \<= MinBidQty
  (intraday change possible).

- MinPrefSec - minimum number of securities allotted to be preferential
  allocation (default 100)

- MaxLot - maximum number of securities allotted to one and the same
  bidder (default 200).

- MinHolders - the 4(B) total number of holders requirement (default
  2,200)

- MinRndLots - the 4(C) minimum round lots requirement (default 450)

- MinMktVal - the 4(C) market value requirement (default \$2,500)

- MinLrgPct - the 4(C) percentage of MinRndLots to exceed MinMktVal
  requirement (default 50%)

- System internal parameter:

  - PreBidders - current number of Preferential Bidders with bid limit
    price \>= Cprice

  - RndBidders - current number of Preferential Bidders where total bid
    quantity of bids with limit price \>= Cprice is \>= RndLot

  - LrgLot = if Cprice\*RndLot\<MinMktVal then LrgLot = MinMktVal/Cprice
    else LrgLot = RndLot (Cprice is the current Indicative Clearing
    Price)

  - LrgBidders - current number of Preferential Bidders where total bid
    quantity of bids with limit price \>= Cprice is \>= LrgLot

> For context: If before the auction start there are already some
> existing holders, then the Lead Manager adjust the 4(B) and 4(C)
> parameters to reflect only the delta required to meet the Listing
> Exchange Minimum requerements.
>
> **Listing Exchange Minimum Fulfillment Monitoring**

The Lead Manager terminal (BatMan) includes a monitoring window
displaying the following values at the Indicative Clearing Price
(Cprice):

(a) Number of Preferential Bids

(b) Number of Preferential Bidders (PreBidders) and
    PreBidders/MinHolders percentage

(c) Number of Preferential Bidders with total 'in the money' quantity of
    at least RndLot (RndBidders) and RndBidders/MinRndLots percentage

(d) Number of Preferential Bidders with total 'in the money' quantity of
    at least LrgLot (LrgBidders) and LrgBidders/(MinRndLots\*MinLrgPct)
    percentage if LrgLot\>RndLot, else display N/A

(e) If PreBidders\<MinHolders display "Preferential bidders PreBidders
    \< min MinHolders", else display quantity of Preferential Bid
    securities at (i) minimum required (i.e.,
    MinHolders\*MinPrefAlloc), (ii) max potential allocation
    quantity[^6] for first MinHolders number of preferential bidders 'in
    the money' bids (i.e., sum of the smaller of the total 'in the
    money' bid quantity and MaxLot for each of the first MinHolders
    number of preferential bidders), and (iii) max allocation quantity
    for all preferential bidders

(f) If RndBidders\<MinRndLots display "Preferential Round Lot bidders
    RndBidders \< min MinRndLots", else display quantity of Preferential
    Bid Round Lot securities at (i) minimum required (i.e.,
    MinRndLots\*RndLot), (ii) max round lot qualified potential
    allocation^3^ quantity for first[^7] MinRndLots number of qualified
    preferential bidders with total 'in the money' bid quantity\>=RndLot
    (i.e., sum of the smaller of the total bid quantity and MaxLot for
    each of the first MinRndLots number of bidders, each having total
    'in the money' quantity\>=RndLot), and (iii) all max round lot
    qualified allocation bid quantity for all preferential bidders with
    'in the money' quantity\>=RndLot

(g) If LrgLot\<=RndLot display N/A else if
    LrgBidders\<MinRndLots\*MinLrgPct display "Preferential Large Lot
    bidders LrgBidders \< min MinRndLots\*MinLrgPct", else display
    quantity of Preferential Bid LrgLot securities at (i) minimum
    required (i.e., MinLrgPct\*MinRndLots\*LrgLot), (ii) max large lot
    qualified potential allocation^3^ quantity for first^4^
    MinRndLots\*MinLrgPct number of qualified preferential bidders with
    total 'in the money' bid quantity\>=LrgLot (i.e., sum of the smaller
    of the total bid quantity and MaxLot for each of the first
    MinLrgPct\*MinRndLots number of bidders, each having total 'in the
    money' quantity\>=LrgLot), and (iii) all max large lot qualified
    allocation bid quantity for all preferential Bidders with 'in the
    money' quantity\>=LrgLot

(h) Calculate the minimum price (Mprice) at which there are
    MinLrgPct\*MinRndLots bidders with a market value \>= MinMktVal and
    display "At Mprice there are #bidders with a bid value \>=
    MinMktVal". If there are less than MinLrgPct\*MinRndLots bidders
    with a market value \>= MinMktVal at the low Price Collar, use the
    low Price Collar as Mprice and display "There are only #bidders with
    a bid \>= MinMktVal at Mprice which is less than the required
    MinLrgPct\*MinRndLots".

The BatMan monitoring window also includes a simulation function where
the Lead Manager can enter a tentative Offering Price to calculate the
above values.

#### [Transition from Open to Close Pending State Blocked]{.mark}

> If at any time during 60 minutes before Close Pending start
> (configurable) the Listing Exchange Minimum requirements are not
> fulfilled at the Cprice, the MarkOps and BatMan app displays a
> question asking for confirmation to proceed with transition to Close
> Pending at the set time. The same request for confirmation is also
> displayed 5 minutes before Close Pending start and if at any time
> during Close Pending the Listing Exchange Minimums are not fulfilled
> at the Cprice. In this case the Close is blocked until confirmation.
> **The Lead Manager must only make such confirmation when confident
> that an Offering Price can be set that satisfy the Exchange Listing
> Minimum requirements**.
>
> The test is as follows:

If at Cprice PreBidders\>=MinHolders or (RndBidders\>=MinRndLots and
LrgBidders\>=MinRndLots\*MinLrgPct) then proceed with allocation else
the exchange's minimum listing requirement is not fulfilled and the
transition to Close Pending is blocked.

### [Preferential Allocation Algorithm]{.mark}

The PDP allocates shares to eligible Preferential Bids until the larger
of the Minimum Preferential Quantity or the quantity that satisfy the
exchange's minimum listing

Below is a model process for implementing the example described in
section 5.2.1 above. The developer can choose another implementation if
it is deemed more efficient.

1.  Select rule 4(B) or 4(C) depending on the smallest quantity of
    securities needed to fulfill the minimum listing requirement

2.  If rule 4(B), allocate MinPrefAlloc to preferential bidders in
    'first come first served' order based on the latest of the bidders
    last Preferential Bid entry time and last Preferential Bid
    modification time, until MinHolders is reached. Then if MaxLot is
    not null allocate securities up to MaxLot in the same 'first come
    first served' order until MinPrefSec is reached. This second
    allocation round starts over from the first bidder, i.e., bidders
    that in the first allocation round received MinPrefAlloc might
    receive additional preferential allocation.

3.  If rule 4(C) and LrgLot\>RndLot, allocate one LrgLot to preferential
    bidders with a total successful bid quantity of at least LrgLot
    until the number of bidders that receive allocation equals
    MinRndLots\*MinLrgPct, then proceed with step 4 (keeping priority as
    in this step 3)

4.  If rule 4(C) allocate one RndLot to preferential bidders that have
    not received any potential LrgLot allocation in step 3 and has a
    total successful bid quantity of at least RndLot until MinRndLots is
    reached (priority as in step 2). Then allocate additional securities
    up to MaxLot in oldest to newest order until MinPrefSec is reached.
    This additional allocation starts over from the first preferential
    bidder, regardless of bid quantity, i.e., bidders with bids smaller
    than one RndLot and bidders that already received RndLot or
    potential LrgLot might receive additional preferential allocation.
    If MinPrefSec is set too high or MaxLot is set too low so that
    allocation MaxLot to all bids does not fulfill MinPreSec, then the
    PDP displays a warning and ask the LM to adjust the parameters.

5.  Proceed with allocation according to the allocation method set for
    the auction after reducing the total issue quantity to allocate with
    the quantity already allocated above. In the Allocation Windows,
    separately display the three quantities; Offering Quantity,
    Preferential Quantity and Offering Quantity minus Preferential
    Quantity; the latter being the quantity for allocation to
    Competitive Bids. Preferential Bids, or part of Preferential Bids,
    that have not received allocation participate as Competitive Bids.
    Preferential Bids might receive up to three executions, i.e., one or
    more of a preferential Allocation, an additional preferential
    allocation and a Competitive Bid allocation.

> \[*[**Note to draft**: ClearingBid has a concept called "**implied
> holders**", but this would need to be approved by the exchange and SEC
> so for now it is not included in this specification.]{.mark}*\]

# Pricing Process

**Finalization at Closing:** At the end of the Bid Period (Closing), no
further orders or modifications are allowed. The final Clearing Price is
determined, and the Lead Manager sets the Offering Price, which cannot
exceed the Clearing Price but may be set lower\[1\]\[2\]\[3\].

The Pricing Process is when the Lead Manager determines:

i.  the final **Offering Price**,

ii. the size of any optional **Lead Manager Short** quantity, and

iii. if Priority Group Pro-rata Allocation is used the **Fill Rates**
     for the different Priority Groups

## [Priority Group Fill Level Simulation and Monitoring]{.mark} 

Priority Group Fill Level Monitoring and Simulation enables the Lead
Manager to optimize the Offering Price and Priority Group Fill Rates
when the auction is using the Priority Group Pro-rata Allocation method.
This process involves:

- **Simulation**: Iterative adjustment of fill targets for each group
  using a Pricing and Allocation Tool (PA Tool), allowing LMs to model
  outcomes before finalizing pricing. The simulation can be run and
  re-run at any time until the LM hits the submit button during the
  Pricing State.

- **Monitoring**: Tracking of Fill Rates a quantities per Priority
  Group.

**Key Workflow**

1.  PDP displays the current Indicative Clearing Price (note that any
    optional Lead Manager Short quantity is not included in the
    calculation of the current Indicative Clearing Price).

2.  The Lead Manager establishes/changes a tentative Offering Price and
    optionally add/change a Lead Manager Short quantity which is added
    to the available quantity for the auction. The available quantity is
    the sum of Primary quantity, Up/down-size of Primary quantity and
    Committed Secondary quantity minus the current quantity reserved for
    Preferential Bids. The PA Tool then calculates a baseline straight
    pro-rata fill rate for the Offering Price where all Priority Groups
    receive the same relative allocation, for example a 62.48% Fill
    Rate.

3.  The Lead Manager sets Fill Rates for all Priority Groups except the
    last. For example, if there are 3 groups, set Priority Group 1 Fill
    rate to 100% and Priority Group 2 Fill Rate to 70%.

4.  The PA Tool then generates the Fill Rate for the last Priority Group
    and a list of the 10 (configurable) largest orders per Priority
    Group that has been excluded due to their MinQty condition. See
    section 7.3 for a detailed description of the allocation algorithm.

5.  Towards the end of the Bidding Period and during Pricing, the Lead
    Manager iteratively runs step 2-4 until a satisfactory Offering
    Price and Fill Rate distribution is achieved.

6.  The submit button is grayed out when not in the Pricing state.
    During the Pricing state the Lead Manager can hit the submit button
    and the system then generates the allocations as further described
    in section 7.3.

When implementing Two-sided auction, see additions to the simulation
described in section 4.1.3.

## [Pro-rata Simulation and Monitoring]{.mark} 

When the auction uses the Pro-rata Allocation method, the simulation and
monitoring is simple as there are no target fill rate and monitoring is
of Fill Rate and MinQty exclusions.

The workflow is similar to the Priority Group Level Simulation workflow
with the exception that the Lead Manager does not set any Fill Rate,
only the Offering Price and Lead Manager short quantity.

## [Non-Pro-rata Allocation method Simulation]{.mark}

For the non-pro-rata allocation methods, i.e., Price-Time Priority, Time
Priority and High Price Before Time Priority, the simulation provides
the following metrics (Lead Manager set Offering Price and Lead Manager
short quantity):

- Number of orders receiving full allocation

- Fill rate of the order receiving partial allocation

- Number of orders receiving no allocation despite having a limit price
  at or higher than the Offering Price

When implementing Two-sided auction, see additions to the simulation
described in section 4.1.3.

## [Offering Price Submission and Publication]{.mark}

When satisfied with the simulation result, the Lead Manager hits the
submit button which enables MarkOps to start Allocation Integrity
Verification (see section 7.7).

After successful Allocation Integrity Verification has been confirmed by
MarkOps, the Lead Manager hit the release button which triggers
publishing of the Offering Price and transition to the Execution state.

Notes for context:

- The Offering Price can be set above or below the published price
  range, subject to certain regulatory and disclosure requirements, but
  there is on PDP functionality associated with this.

- The Offering Price is normally published by 6:00 p.m. ET on pricing
  day and later added to the effective registration statement and final
  prospectus.

# Allocations Process

The Allocation Process distributes the shares among eligible bids, i.e.,
bids with a limit price at or higher than the Offering Price. For each
new issue one allocation method is selected. The allocation methods are:

- Pro-rata

- Priority Group Pro-rata

- Price-Time Priority

- Time Priority

- High Price Before Time Priority

If the offering includes Preferential Bids, allocation to Preferential
Bids is performed before using one of the above methods for allocation
to Competitive Bids.

## [Allocation to Preferential Bids]{.mark}

If the offering includes Preferential Bids, allocation to Preferential
Bids is performed before allocation to Competitive Bids. The quantity
available for Competitive Bids is then reduced by the quantity allocated
to Preferential Bids. The detailed description of Preferential Bids can
be found in section 3.2.2 and of allocation to Preferential Bids in
section 5.2.2.

## [Pro-rata Allocation]{.mark}

**Pro-rata allocation** means that all qualifying orders receive a
proportional share of the available securities, calculated as a
percentage of their bid relative to the total quantity offered. For
example, if the total demand at or above the offering price is 25,000
shares but only 10,000 are available, each order would receive 40% of
its requested quantity.

**Minimum Quantity Condition (MinQty):** If an order specifies a minimum
quantity that must be allocated for the bidder to accept the allocation,
and the pro-rata allocation would result in less than this minimum, the
order is excluded from allocation. The system then recalculates the
pro-rata percentages for the remaining orders. Note that this may result
in MinQty orders that receive no allocation even if the final pro-rata
percentage is higher than their specified minimum, due to redistribution
of quantities from excluded MinQty orders.

**Allocation Algorithm**

The Allocation Algorithm is similar to Priority Group Pro-rata
Allocation excluding step 2-5.

## [Priority Group Pro-rata Allocation]{.mark}

See also the Guidelines. Orders are assigned Priority Group according to
the order entry or last modification time. Groups with a lower Priority
Group number receive higher allocations. For Priority Group Pro Rata
there are typically 3 Priority Groups, but the system can be configured
for more or less groups. Priority Groups are named 1, 2, 3, etc. The
term Prio 1 Order is used to denote an order assigned Priority Group 1,
representing the highest allocation preference.

See also section 3.10 ([Priority Groups]{.mark}) and 3.9 ([Seasoning
Periods]{.mark}).

**Allocation Methodology**:

- After allocation to any **Preferential Bids**, the remaining
  securities are allocated among successful competitive bids according
  to their Priority Group. Preferential Bids or part of Preferential
  Bids that did not receive preferential allocation are included among
  competitive bids with Priority Group as of the bid timestamp.

- **Pro rata allocation percentages** are set by the issuer and
  underwriters for each group, with Priority Group 1 always receiving an
  allocation percentage equal to or greater than Group 2, which in turn
  is equal to or greater than Group 3.

- Example allocation percentages might be: Group 1 receives **Fill
  Rate** (**Fill%**)of 100%, Group 2 receives 70%, and Group 3 receives
  36.375% of their eligible bid quantities.

- Within each group, pro rata allocation is applied, and bids with a
  Minimum Quantity Condition (MinQty) only receive shares if their
  minimum is met. Elimination of MinQty orders from the last group might
  result in a fina Fill Rate that is higher than the eliminated MinQty
  orders due to redistribution.

**Allocation Algorithm**

Below is an example algorithm for allocation. Development of the PDP are
free to improve on or replace this algorithm.

1.  Reduce the available quantity with the quantity reserved for
    Preferential Bids (MinPrefSec).

2.  For all orders at or above the Offering Price in Priority Group 1,
    traverse the book and multiply order size with the target fill rate.
    Allocation is rounded to integer, but not less than 1.

3.  If resulting allocation is less than the MinQty, exclude order and
    move to next.

4.  Repeat 2 and 3 for all Priority Groups except the last.

5.  Decrease the total volume by the quantity allocated to all Priority
    Groups except the last.

6.  Divide the remaining quantity with the total qualified quantity of
    the Last Priority Group to produce a Last Group Fill Target.

7.  Traverse the Last Priority Group's qualified orders and identify all
    orders with MinQty larger than Last Group Fill Target.

8.  Exclude the MinQty order with the largest relative gap to the Last
    Group Target, re-calculate Last Group Fill Target and restart from 7
    until there are no orders with MinQty larger than Last Group Fill
    Target.

9.  Apply Last Group Fill Target to the remaining orders in the Last
    Priority Group. Note that this might result in order(s) with MinQty
    smaller than Last Group Fill Target receiving no allocation.

10. If the allocated quantity is larger than the available quantity,
    reduce allocation with 1 share/bond to orders in re-reverse time
    priority, i.e., from latest to oldest, except for orders that have
    been allocated 1 share/bond.

11. If there are remaining un-allocated quantity, allocate 1 share/bond
    to orders in time priority, i.e., from oldest to newest.

When implementing Two-sided auction, see additions to the simulation
described in section 4.1.3.

## [Price-Time Priority Allocation]{.mark}

An allocation method where orders receive allocation in strict
price-time order, subject to minimum quantity conditions (see MinQty).
Orders receive full fill in strict price-time order until there a not
enough quantity to fully fill the next Order. Said next Order then
receive a partial allocation consisting of the remaining quantity if
passing the MinQty condition. All remaining Order receive no allocation
even when they have a price that is equal or higher than the Offering
Price.

This is [not a pro-rata method]{.underline}. Note that when the Offering
Price is set below the Clearing Price, the full quantity will still be
allocated at the Clearing Price level so [orders with a limit price at
or above the Offering Price but below the Clearing Price will not
receive any allocation]{.underline}.

This is the traditional secondary market allocation model.

## [Time Priority Allocation]{.mark}

Note, the Time Priority Allocation method is [not a pro-rata
method]{.underline}. All bids at or above the Offering Price will not
receive allocation. Allocation depends on the time of the bid. Orders
with a limit price above the Clearing Price receive full allocation.
Orders where the limit price equals the Clearing Price receive full fill
in strict time order until there is not enough quantity to fully fill
the next Order. Said next Order then receive a partial allocation if
passing the MinQty condition. All remaining Orders receive no allocation
even when they have a price that is equal or higher than the Offering
Price.

The Time Priority Allocation method is designed to first recognize
investors that are willing to pay a higher price, then investors that
make a larger price discovery contribution.

## [High Price Before Time Priority Allocation]{.mark}

Note, the High Price Before Time Priority Allocation method is [not a
pro-rata method]{.underline}. All bids at or above the Offering Price
are not eligible for allocation. Allocation eligibility depends on both
the price and time of the bid. Orders with a price above the Clearing
Price receive full allocation. Orders [at the Clearing Price together
with orders at or above the Offering Price]{.underline} but below the
Clearing Price receives allocation in strict time priority contingent on
if passing the MinQty condition.

The High Price Before Time Priority is designed to recognize investors
that are prepared to pay a higher price while including a time component
to recognize orders with larger price discovery contribution.

## [Allocation Integrity Verification]{.mark}

Allocation Integrity Verification is designed to prevent operational and
reputational risks by ensuring allocation accuracy before Trade
Execution and subsequent submission to ClearingBid's clearing broker and
NSCC. It involves automated checks to validate allocation correctness,
with discrepancies flagged for investigation. Below are the draft rules
for verifying the generated allocations:

- List of orders with no allocation despite having a price at or above
  the offering price (including info on MinQty). There should only be
  MinQty orders on this list. There should only be a few MinQty orders
  with a MinQty that is slightly below the Fill Rate for the Priority
  Group (this can happen due to that the quantity freed up by an order
  that is excluded is redistributed among the other orders hence
  increasing the Fill Rate). If there are MinQty orders on the list with
  a MinQty% much larger than the fill% then this must be investigated
  for correctness.

- List of orders with allocation despite having a price lower than the
  Offering Price (should be none).

- List of allocations that differ more than 1 share from the priority
  group fill%. The details of the check need to be worked out
  considering the rounding rules of the allocation and the "round-robin"
  allocation of any residual quantity in step 10 or 11 of section 7.3.
  There should be no orders on the list if there are no errors. For
  orders that have received Preferential Allocation this verification is
  only applicable for the part of the order that participate as a
  Competitive Bid.

- Verify that total number of allocated shares equals issue total
  quantity.

- Statistics e.g., number of orders, allocations and allocation
  quantities per priority group; number of preferential orders and
  preferential orders in-the-money; number of preferential allocations;
  etc.

- Highlight very large orders, e.g., orders receiving more than 3% of
  available quantity (configurable).

**Key workflow**:

1.  The Lead Manager allocation simulations.

2.  When in the Pricing state and the Lead Manager is satisfied with the
    simulation the Lead Manger hit the Submit button. A snapshot of the
    orderbook is saved, the allocations are generated internally in the
    system and the Allocation Verification Report is generated.

3.  When Market Operation and the Lead Manager conclude that the
    verification is successful then Market Operations submit the
    Allocation Verified command (or hit the said button) which triggers
    the transition to the Execution state.

4.  After submit the EOD file for Wedbush is created and executions are
    disseminated via FIX, but cancellations are held back

5.  Market Operations is presented a on-top window with a countdown for
    sending out FIX cancellations. Cancellations are sent out
    automatically at the end of the countdown if Market Operations do
    not hit ABORT before. Default countdown is 30 minutes
    (configurable).

# Trade Execution, Confirmation and Clearing Process

## [Auction Close Day Executions]{.mark}

**Auction Close Day Executions** refer to the process by which trade
executions are generated and distributed to participants on the day of
the Auction Close and Pricing.

**Definition and Timing**: Auction Close Day Executions are trade
executions generated immediately after pricing.

**For context**:

- This execution method is primarily used when the auctioned securities
  have a recent secondary market trade price available. It enables rapid
  confirmation.

- The standard settlement period T+1 is used (trade date plus one
  trading day), unless allocations are requested to be tracked by the
  DTC IPO Tracking System, in which case the settlement follows the
  process designated by the Lead Underwriter.

- The FIX execution message for standard settlement does not include any
  special settlement instructions.

**Distribution**:

Trade execution confirmation messages are disseminated via the CB FIX
gateway to connected participants.

Confirmations can also be accessed through BidMan.

## Next Trading Day Executions

**Next Trading Day Executions** refer to the process in which **trade
executions are generated and disseminated on the morning of the Trading
Day following the Auction Close and Pricing**, rather than on the day of
the Auction Close itself. This process is particularly relevant for
securities where a recent secondary market trade price is not available,
such as newly listed equities\[1\]\[2\]\[3\].

**Key Features of Next Trading Day Executions**

**Timing:** Trade executions are distributed before the market opens on
the Trading Day after Auction Close.

**Settlement:** Except for allocations tracked by the DTC IPO Tracking
System, these executions normally have **same-day settlement (T+0, or
\"cash\")**, and the corresponding FIX execution message includes a
same-day-settlement instruction (FIX Tag 63 = 1). For allocations marked
to be tracked by the DTC IPO Tracking System, no settlement instruction
is included, and settlement follows the process set by the Lead
Underwriter.

**Distribution:** Execution confirmations are sent via the CB FIX
gateway to participants, and can also be accessed via the BidMan
application (for Broker-Dealers) and BatMan application (for
Underwriters). For allocations not marked to be tracked by the DTC IPO
Tracking System, trade executions are also confirmed via NSCC UTC output
messages.

## Clearing and Settlement File Generation

ClearingBid\'s clearing and settlement process involves generating files
that facilitate the delivery of securities and cash transfers after an
auction has closed. The system distinguishes between **non-tracked
allocations** (handled like secondary market trades) and **tracked
allocations** (managed by the Lead Manager via DTCC's IPO Tracking
System). Key steps include:

**Allocation Separation**:

**Non-tracked allocations**: PDP generates an "EOD file" (End-Of-Day)
with all execution that are not to be tracked and sends the file to
ClearingBid's clearing broker. The file format is per the clearing
brokers specification. For context, the clearing broker converts the
executions to UTC Input for transmission to the DTCC's Universal Trade
Capture (UTC) system and further Continuous Net Settlement (CNS) and
National Settlement Service (NSS) processing.

**Tracked allocations**: PDP generates one EOD file for the Lead Manager
with all tracked allocations and one EOD file per broker-dealer with
their respective allocations, if any. The allocations in the EOD files
are grouped by institutional (to be handled via DTCC's Institutional
Delivery System) and retail (delivered to broker-dealers' DTC IPO
Control Accounts).

## FIX Trade Execution and Order Cancellation Messages

Fix Execution Messages are sent directly after creation.

FIX order cancellation messages for partially executed orders and
un-executed orders are held back and not sent out until ClearingBid's
clearing broker confirms acceptance of all settlement transactions by
NSCC UTC and no other abnormal issues have been identified. The FIX
order cancellation messages are then released by MarkOps manual command.

## Error Handling and Reversal/Bust Management

If an error is discovered after closing (e.g., incorrect parameters,
system error), the system allows for a reversal or \"bust\" of the
closing. This restores the order book to its state prior to the saved
snapshot, allowing corrections and a redo Pricing and Allocation. This
process includes sending reversal messages via FIX and new EOD clearing
files. The bust process is complex and may require coordination with
clearing brokers and external systems and may not always be executable
on the same day.

Some examples of events requiring reversal:

- The QSR between a broker-dealer and ClearingBid's clearing broker no
  longer in place

- Defaulting participant

- Lead Manager or MarkOps process mistake

- PDP bug

**Initiation:** The process is typically initiated by a Market
Operations Manager or Lead Manager via a dedicated \"Bust\" button.
Initiating the Bust process is not possible after MarkOps has manually
confirmed successful Trade Executions because FIX order cancellation
messages has then been sent out and any restore of the master orderbook
will thereby be out of sync with broker-dealer order management systems.

**Bust process workflow:**

1.  **Orderbook state transition to MODIFICATION**

2.  **FIX reversals:** PDP generates and sends reversal (bust)
    transactions corresponding to all previously distributed FIX Trade
    Execution messages.

3.  **Master Orderbook Restore:** The before Lead Manager submittal of
    Offering Price Master Orderbook snapshot is restored.

4.  **UTC Input reversal**: If NSCC have accepted UTC Input that are
    effected by the bust UTC Input reversals must be generated and sent
    to NSCC.

5.  **Corrective Actions applied**: Actions applied to eliminate the
    reason for the bust.

6.  **Orderbook state transition to PRICING**

7.  **Rerun Pricing, Allocation, and Trade Execution**

# Parameters Attributes

## [PDP System Wide Parameters and Attributes]{.mark}

**\[TBD: Tor S. to review and update this section\]**

**SEC Effectiveness Delay**: Default 60 minutes

## Asset Types

**\[TBD: Tor S. to review and update this section\]**

The ClearingBid platform recognizes **three primary asset types** for
offerings, each with distinct attributes and behaviors:

  -----------------------------------------------------------------------------------------
  **Asset Type** **Price       **Issue      **Bidding     **Key Attributes**
                 Discovery**   Quantity**   Mechanism**   
  -------------- ------------- ------------ ------------- ---------------------------------
  **Stock**      Yes           Fixed        Price per     Symbol, CUSIP, Name, Issuer,
                                            unit          Market, Lead Manager, Number of
                                                          decimals (default 2), Price
                                                          increment (default 0.01),
                                                          Minimum/Maximum Order Quantity,
                                                          Primary/Secondary Quantity, Price
                                                          Range, Gross Underwriting Spread,
                                                          Selling Concession, Allocation
                                                          Algorithm, Order Book
                                                          States\[1\]\[2\]\[3\]\[4\]

  **Bond**       Yes           Fixed        Yield         Symbol, CUSIP, Name, Issuer,
                                            (various      Market, Lead Manager, Number of
                                            types)        decimals (default 3), Price
                                                          increment (default 0.001), Face
                                                          Value (default 1000), Issue Size,
                                                          Yield Type (YTM-FV, YTM-CP,
                                                          YTC-FV, YTC-CP, SpT-SM, SpT-FR),
                                                          Gross Underwriting Spread,
                                                          Selling Concession, Offering
                                                          Price, Offering Coupon, Offering
                                                          Yield, Allocation Algorithm,
                                                          Order Book
                                                          States\[1\]\[2\]\[3\]\[4\]

  **Fund/ETF**   No (demand    Variable (by Price per     Symbol, CUSIP, Name, Issuer,
                 discovery     NAV)         unit (NAV)    Market, Lead Manager, Number of
                 only)                                    decimals (default 2), Price
                                                          increment (default 0.01), CNAV
                                                          (Clearing NAV), Target Volume,
                                                          Distribution Fee, Latest NAV
                                                          Update, Allocation Algorithm
                                                          (NAV), Order Book
                                                          States\[1\]\[2\]\[3\]\[5\]\[4\]
  -----------------------------------------------------------------------------------------

## [Common Attributes for All Asset Types (Stock, Bond, and Fund/ETF)]{.mark}

**\[TBD: Tor S. to review and update this section\]**

Based on the provided documents, here are the common attributes for all
asset types (Stock, Bond, and Fund/ETF), along with identified
contradictions:

**Common Attributes for All Asset Types**

**OrderBookID**: A unique identifier for the offering (e.g., symbol +
date) \[1\]\[2\]\[3\].

**Symbol**: Ticker or CUSIP identifier \[1\]\[2\]\[3\]\[5\].

**Issuer**: Name of the entity issuing the security
\[1\]\[2\]\[3\]\[5\].

**Market/Listing**: Trading venue where the security will be listed
(e.g., NYSE) or \"Not Listed\" \[1\]\[2\]\[3\]\[5\].

**Lead Manager**: Underwriter selected by the issuer to manage the
offering\[1\]\[2\]\[3\]\[4\].

**Lead Manager BD**: Lead Manager related Broker-Dealer for placing
orders.

**Type**: Asset class (Stock, Bond, Fund/ETF) \[1\]\[3\]\[5\].

**SEC Effectiveness Data/Time**: No default, must be set by LM to N/A or
the actual time of effectiveness. If not set, then transition to Pricing
state is prohibited.

**Regular Trading Close Time**: Default 4:00 pm EST, but can be changed
by LM. Note, that this normally refers to the closing time for regular
trading in the market of the security, but it can also be used to set
specific closing time for the individual auction because it is on the
security level.

**State**: Lifecycle stage (e.g., NEW, UPCOMING, OPEN, CLOSED)
\[1\]\[2\]\[3\]\[5\].

**Number of Decimals**: Price precision (default: 2 for Stocks/Funds, 3
for Bonds) \[1\]\[2\]\[3\].

**Price Increment**: Minimum price movement (default: \$0.01\$ for
Stocks/Funds, \$0.001\$ for Bonds) \[1\]\[2\]\[3\].

**Quantity Increment**: Smallest allowable order size (configurable,
default: 1 unit) \[2\]\[3\].

**Min/Max Order Quantity**: Configurable limits per order
\[1\]\[2\]\[3\].

**Min/Max Price Allowed**: Dynamic price collars to reject orders
outside bounds \[1\]\[2\]\[3\].

**Target Orders**: Minimum number of orders required to close the
offering \[1\]\[2\]\[3\].

**Allocation Algorithm**: Method for allocating shares (e.g., Pro-Rata,
Price-Time Priority) \[1\]\[2\]\[3\].

**Contradictions in Documents**

**\"Type\" Attribute**:

\[1\]\[3\] define \"Type\" strictly as **Stock**, **Bond**, or
**Fund/ETF**.

\[5\] suggests it includes **security subtypes** (e.g., \"Common
Stock,\" \"Corporate Bond\"), which conflicts with \[1\]\[3\].

**\"State\" Definitions**:

\[1\] lists states as NEW, UPCOMING, OPEN, CLOSE PENDING, CLOSING,
CLEARING, CLOSED, CANCELLED, FROZEN, HALTED.

\[5\] uses Registering(unmentioned in \[1\]\[3\]) and omits CLEARING and
CLOSE PENDING, creating inconsistency.

**Allocation Algorithms**:

\[1\]\[2\] list four algorithms: **Time Priority Pro-Rata**,
**Pro-Rata**, **Price Priority Pro-Rata**, and **Price-Time Priority**.

\[3\] states *\"Price-priority Pro-rata do not have to be
implemented\"*, contradicting \[1\]\[2\].

**Blocked Broker-Dealers Function**:

\[1\]\[2\] explicitly mention this as a Lead Manager capability.

\[3\]\[5\] omit this feature entirely.

**Order Entry/Cancellation Rules**:

\[1\]\[2\] specify intraday configurable \"Cancel Time Limits\" per
Priority Group.

\[3\] describes these limits but omits intraday configurability, while
\[5\] lacks details.

**Summary**

Common attributes focus on **identification**, **pricing rules**,
**order constraints**, and **allocation methods**. Key contradictions
involve granularity of \"Type,\" lifecycle states, allocation algorithm
implementation, and feature consistency (e.g., Blocked Broker-Dealers).
These discrepancies suggest evolving specifications or document version
mismatches.

Sources \[1\] CB-PDP-Gen-2-specification-updated-definitions-04.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/d0735754-a7b5-4d48-9c69-5fd13fba6142/CB-PDP-Gen-2-specification-updated-definitions-04.docx>\[2\]
ClearingBid-Backend-Functions-CMC-v2-7.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/60bd3f30-70db-4db8-9031-579093dfb8b8/ClearingBid-Backend-Functions-CMC-v2-7.docx>
\[3\] ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>
\[4\] CB-Backend-Specification-v1-4-8-comments.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/571d1572-1928-4edd-a5d2-fd3b20d2ba8d/CB-Backend-Specification-v1-4-8-comments.docx>
\[5\] BE-Terminal-Reference-v004-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/bb17a992-9932-4b23-b225-8bb487beab97/BE-Terminal-Reference-v004-2.docx>
\[6\] Fund-ETF-Summary-v01.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/03b9046e-a9fd-4ffa-9ca1-ea087a21e91a/Fund-ETF-Summary-v01.docx>
\[7\] CB-PDP-Gen-2-specification-v2-0r01.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/f0a7c1b2-8ef0-428d-b64c-f20597bfe366/CB-PDP-Gen-2-specification-v2-0r01.docx>
\[8\] Holiday-schedule-010725-143132.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/68096c24-9212-4aa6-b3a0-fe6457226119/Holiday-schedule-010725-143132.pdf>
\[9\] SD-PRD\_-Allocation-Integrity-Verification-300625-211230.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/434a11b8-fa30-4fe2-85c8-13e2062fb3d1/SD-PRD_-Allocation-Integrity-Verification-300625-211230.pdf>\[10\]
SD-PRD\_-PDP-Changes-for-Fixed-Income-Prototype-010725-151230.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/ffaad82f-212c-4181-bfa8-634619e12be7/SD-PRD_-PDP-Changes-for-Fixed-Income-Prototype-010725-151230.pdf>
\[11\] BidMan-User-Guide-v0-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/e5243159-49ba-46a2-9c45-c585d64ce1a1/BidMan-User-Guide-v0-2.docx>
\[12\]
ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/25825a46-ec4c-4544-93e0-5e46d77e1d5a/ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx>
\[13\] ClearingBid-IPO-Price-Discovery-and-Allocation-v0-0-8.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/4cb22cd3-2d63-4d83-bd30-d28c1f0d2bde/ClearingBid-IPO-Price-Discovery-and-Allocation-v0-0-8.docx>

## [Stock Attributes]{.mark}

**\[TBD: Tor S. to review and update this section\]**

**Stock Attributes** in the ClearingBid platform are a defined set of
data fields and parameters that describe and control the behavior of
stock offerings during the auction and allocation process. These
attributes are consistently described across the provided documents,
with only minor differences in terminology or emphasis.

**Core Stock Attributes**

The following are the main attributes for stocks, as extracted from
multiple documents:

**Name**: Full name of the offering\[1\]\[2\]\[3\]\[4\].

**Issuer**: Name of the company issuing the stock\[1\]\[2\]\[3\]\[4\].

**Symbol**: Exchange-assigned symbol or CUSIP
identifier\[1\]\[2\]\[3\]\[4\].

**CUSIP**: Unique code for the offering\[1\]\[2\]\[3\].

**Market/Listing**: Trading venue where the stock will be
listed\[1\]\[2\]\[3\]\[4\].

**Offering Size/Total Quantity**: Number of shares offered in the
auction\[1\]\[2\]\[3\]\[4\].

**Primary Quantity**: Number of new shares offered by the
issuer\[2\]\[3\].

**Secondary Quantity**: Number of existing shares offered by current
holders\[2\]\[3\].

**Target Volume**: Minimum total volume required for the offering to
close\[1\]\[2\]\[3\].

**Target Orders**: Minimum number of orders required for the offering to
close\[1\]\[2\]\[3\].

**Price Range (Low/High)**: Estimated price range set by the Lead
Manager, can be updated during the Bid Period\[1\]\[2\]\[3\]\[4\].

**Clearing Price**: The price at which the full offering is subscribed;
calculated in real-time during the Bid Period\[1\]\[2\]\[3\]\[4\].

**Offering Price**: Final price set by the Lead Manager and Issuer,
cannot exceed the Clearing Price\[1\]\[2\]\[3\].

**Minimum/Maximum Order Quantity**: Limits on order sizes, configurable
per offering\[2\]\[3\].

**Minimum/Maximum Price Allowed**: Price collars for order
entry\[2\]\[3\].

**Minimum Allocation**: Minimum allocation per account at allocation
(reserved for future use)\[1\]\[2\]\[3\].

**Gross Underwriting Spread**: Underwriter\'s fee as a percentage of the
offering price\[1\]\[2\]\[3\].

**Selling Concession**: Portion of the underwriting spread paid to the
Broker-Dealer submitting the order\[1\]\[2\]\[3\].

**DTC IPO Tracking allowed**: Yes/No, default No

**Limit Priced Secondary Shares Offers allowed**: Yes/Lead Manager BD
Only/No, default No

**Preferential Bids allowed**: Yes/Lead Manager BD Only/No, default No

**Dividend/Dividend Yield**: Optional reference values for stocks\[5\].

**Order Book States**: Lifecycle states such as NEW, UPCOMING, OPEN,
CLOSE PENDING, CLOSING, CLEARING, CLOSED, CANCELLED, FROZEN,
HALTED\[2\]\[3\]\[4\].

**Additional Attributes and Parameters**

**aVWAP/rVWAP/VWAP**: Volume-weighted average prices for all orders or
those within the price range\[2\]\[5\]\[3\]\[4\].

**Order Entry Type**: Only Good-Til-Canceled (GTC) limit orders are
accepted for stocks\[2\]\[3\].

**Allocation Algorithms**: Several are supported, including Pro-rata,
Price-Time Priority, and Priority Pro Rata, with specific rules for
each\[1\]\[2\]\[3\].

**Account Maximum Order Size**: To prevent a single account from
dominating the allocation\[1\].

**Minimum Number of Individual Accounts**: To ensure broad participation
and liquidity\[1\]\[2\]\[3\].

**Instrument Reference Data**: All parameters are stored in the backend
and are available via APIs\[1\]\[2\]\[3\].

**Contradicting Information**

The documents are largely consistent in their description of Stock
Attributes. Notable points of potential contradiction or ambiguity
include:

**Terminology for Allocation Algorithms**: Some documents refer to
\"Priority Group Pro Rata\" and \"Price Priority Pro Rata\"
interchangeably, but clarify that the number of groups and their
definitions (by time or price) can be configured\[2\]\[5\]\[3\]. The
distinction is subtle but could cause confusion if not standardized.

**Offering Price Field**: In some places, the \"Offering Price\" is
described as a field to be set in the Allocation Widget, but later
documents state it should be set in the Edit Offering Widget, with the
Allocation Widget only displaying it\[2\]\[3\]. This is a minor UI/UX
specification difference.

**Minimum Allocation**: This field is described as \"for future use\" or
\"reserved\" in several places, indicating it is not currently
implemented but planned\[1\]\[2\]\[5\]\[3\]. There is no contradiction,
but the implementation status should be clarified.

**Order Book State Names**: There are slight differences in the naming
and explanation of order book states (e.g., \"CLOSING\" vs. \"CLOSE
PENDING\"), but the overall lifecycle is consistent\[2\]\[3\]\[4\].

**Dividend Attributes**: Only one document explicitly lists Dividend and
Dividend Yield as stock attributes\[5\], while others do not mention
them, suggesting they are optional or not universally implemented.

**In summary:**\
Stock Attributes on the ClearingBid platform are well-defined and highly
consistent across the documentation, covering all key data fields,
limits, and process controls required for a managed auction of equity
securities. Minor differences are present in terminology and UI field
placement, but there are no material contradictions regarding the
substance or function of Stock Attributes\[1\]\[2\]\[5\]\[3\]\[4\].

Sources \[1\] CB-Backend-Specification-v1-4-8-comments.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/571d1572-1928-4edd-a5d2-fd3b20d2ba8d/CB-Backend-Specification-v1-4-8-comments.docx>
\[2\] ClearingBid-Backend-Functions-CMC-v2-7.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/60bd3f30-70db-4db8-9031-579093dfb8b8/ClearingBid-Backend-Functions-CMC-v2-7.docx>
\[3\] ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>
\[4\] BE-Terminal-Reference-v004-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/bb17a992-9932-4b23-b225-8bb487beab97/BE-Terminal-Reference-v004-2.docx>
\[5\] CB-PDP-Gen-2-specification-updated-definitions-04.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/d0735754-a7b5-4d48-9c69-5fd13fba6142/CB-PDP-Gen-2-specification-updated-definitions-04.docx>\[6\]
BidMan-User-Guide-v0-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/e5243159-49ba-46a2-9c45-c585d64ce1a1/BidMan-User-Guide-v0-2.docx>
\[7\] Auction-Process-ClearingBid-2025-02-09-v2-Draft.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/2723ad3f-6e9d-4c22-a998-b292e737e4e7/Auction-Process-ClearingBid-2025-02-09-v2-Draft.docx>
\[8\] Fund-ETF-Summary-v01.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/03b9046e-a9fd-4ffa-9ca1-ea087a21e91a/Fund-ETF-Summary-v01.docx>
\[9\] Outstanding-Questions-and-Issues-v-0-3.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/968218fe-4eac-4014-b466-41fee1591018/Outstanding-Questions-and-Issues-v-0-3.docx>

## [Bond Attributes]{.mark}

**\[TBD: Tor S. to review and update this section\]**

**Bond Attributes** are the key data fields and characteristics that
define a bond offering on the ClearingBid platform. These attributes are
consistently described across the provided documents, with some minor
differences in terminology and detail. The core bond attributes include:

**Face Value (Par Value):** The amount paid to the bondholder at
maturity, typically \$1,000 per bond\[1\]\[2\]\[3\].

**Bond Offering Price:** The price set by the Lead Manager and Issuer to
be paid at settlement; it may differ from Face Value and is determined
based on the Clearing Price/Yield and various market factors. The
Offering Price cannot exceed the Clearing Price, and combined with the
Offering Coupon, cannot be lower than the Clearing Yield\[1\]\[2\]\[3\].

**Offering Coupon:** The interest rate set for the bond, determined by
the Lead Manager and Issuer so that the offering yield is equal to or
greater than the Clearing Yield\[1\]\[2\]\[3\].

**Yield Type:** Defines how the Offering Price and Coupon are set and
what type of yield the investor is bidding for. Supported types include:

YTM-FV (Yield to Maturity at Face Value)

YTM-CP (Yield to Maturity by Coupon and Price)

YTC-FV (Yield to Call at Face Value)

YTC-CP (Yield to Call by Coupon and Price)

SpT-SM (Spread to Treasury with Similar Maturity)

SpT-FR (Spread to Treasury with Floating Rate)\[1\]\[2\]\[3\].

**Offering Yield:** Calculated and entered by the Lead Manager based on
the Offering Price and Coupon\[1\]\[2\].

**Issue Size:** Calculated as Face Value × Offering Quantity\[1\]\[2\].

**Offering Quantity:** The number of bonds offered on the
platform\[1\]\[2\].

**Gross Underwriting Spread:** The total compensation (in percentage
points) received by underwriters, including the Selling
Concession\[1\]\[2\]\[3\].

**Selling Concession:** The portion of the underwriting spread paid as
commission to the Broker-Dealer that submits the order\[1\]\[2\]\[3\].

**Minimum/Maximum Order Size:**Parameters to control the smallest and
largest order sizes allowed\[1\]\[2\].

**Minimum/Maximum Yield Allowed:**Limits set for acceptable bid
yields\[1\]\[2\].

**Order Book States:** States such as NEW, UPCOMING, OPEN, CLOSE
PENDING, CLOSING, CLEARING, CLOSED, CANCELLED, FROZEN, and HALTED, which
govern the lifecycle of orders\[1\]\[2\]\[3\].

**Allocation Algorithm:** The method used to allocate bonds among
bidders (e.g., Pro-rata, Price-Time Priority, Time Priority
Pro-rata)\[1\]\[2\].

**Bid Expression:** Bids are expressed in yield, and the order book is
ascending (higher yield is better for investors, lower yield is better
for the issuer)\[1\]\[2\]\[3\].

**Price Format:** Bonds can be displayed in yield (percentage), dollar
price (percentage of par), or, for prototypes, in basis points
(bps)\[4\].

**Additional attributes** sometimes listed:

**CUSIP:** Unique identifier for the bond\[1\]\[2\].

**Ratings:** Ratings from agencies\[5\].

**Maturity Date:** The date the bond matures\[5\].

**Redemption Features:** Such as call protection or call premium\[5\].

**Contradicting Information**

A thorough review of the documents reveals **no material
contradictions** in the definition or handling of Bond Attributes.
Instead, the documents are highly consistent, often repeating the same
definitions and logic, sometimes with slightly different wording or
additional details.

**Minor differences or clarifications:**

**Terminology:** Some documents use \"Allocation Priority Group\" while
others use \"Priority Group,\" but both refer to the same concept of
grouping orders for allocation\[1\]\[2\]\[3\].

**Yield Type Implementation:** All documents agree on the supported
yield types, though some provide more detail on calculation or FIX
protocol handling\[1\]\[2\]\[3\].

**Order Book Display:** The possibility to display bond prices in basis
points (bps) is mentioned as a prototype feature in one document\[4\],
whereas others focus on yield or dollar price. This is not a
contradiction but an extension for prototyping purposes.

**Minimum Allocation:** Several documents note that \"Minimum Allocation
per account\" is reserved for future use and is not currently
implemented\[1\]\[2\]\[3\].

**Gross Underwriting Spread Calculation:** Both \[1\] and \[2\] note
uncertainty (TBD) about whether the Gross Underwriting Spread for bonds
is based on Face Value or Bond Offering Price. This is a known open
question, not a contradiction.

**Summary Table: Key Bond Attributes Across Documents**

  -------------------------------------------------------------
  **Attribute**   **Description/Notes**   **Contradictions?**
  --------------- ----------------------- ---------------------
  Face Value      Amount paid at maturity None
                  (usually \$1,000)       

  Bond Offering   Set by Lead Manager;    None
  Price           may differ from Face    
                  Value                   

  Offering Coupon Set to ensure yield ≥   None
                  Clearing Yield          

  Yield Type      YTM-FV, YTM-CP, YTC-FV, None
                  YTC-CP, SpT-SM, SpT-FR  

  Offering Yield  Calculated from Price   None
                  and Coupon              

  Issue Size      Face Value × Offering   None
                  Quantity                

  Offering        Number of bonds offered None
  Quantity                                

  Gross           Total underwriter       Open question: basis?
  Underwriting    compensation            
  Spread                                  

  Selling         Broker-Dealer           None
  Concession      commission              

  Order Size      Min/Max per             None
  Limits          order/account           

  Yield Limits    Min/Max yield allowed   None

  Allocation      Pro-rata, Price-Time    None
  Algorithm       Priority, etc.          

  Bid Expression  Yield (ascending order  None
                  book)                   

  Price Format    Yield, dollar price,    None
                  basis points            
                  (prototype)             
  -------------------------------------------------------------

**Conclusion:**\
The documents present a **coherent and consistent definition of Bond
Attributes** for the ClearingBid platform. Any differences are due to
evolving features (such as basis point display) or open questions (such
as the calculation basis for underwriting spread) that are explicitly
flagged and do not constitute contradictions\[1\]\[2\]\[5\]\[3\]\[4\].

Sources \[1\] ClearingBid-Backend-Functions-CMC-v2-7.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/60bd3f30-70db-4db8-9031-579093dfb8b8/ClearingBid-Backend-Functions-CMC-v2-7.docx>
\[2\] ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>
\[3\] CB-PDP-Gen-2-specification-updated-definitions-04.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/d0735754-a7b5-4d48-9c69-5fd13fba6142/CB-PDP-Gen-2-specification-updated-definitions-04.docx>\[4\]
SD-PRD\_-PDP-Changes-for-Fixed-Income-Prototype-010725-151230.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/ffaad82f-212c-4181-bfa8-634619e12be7/SD-PRD_-PDP-Changes-for-Fixed-Income-Prototype-010725-151230.pdf>
\[5\] CB-Backend-Specification-v1-4-8-comments.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/571d1572-1928-4edd-a5d2-fd3b20d2ba8d/CB-Backend-Specification-v1-4-8-comments.docx>
\[6\] BidMan-User-Guide-v0-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/e5243159-49ba-46a2-9c45-c585d64ce1a1/BidMan-User-Guide-v0-2.docx>
\[7\] BE-Terminal-Reference-v004-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/bb17a992-9932-4b23-b225-8bb487beab97/BE-Terminal-Reference-v004-2.docx>
\[8\] Fund-ETF-Summary-v01.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/03b9046e-a9fd-4ffa-9ca1-ea087a21e91a/Fund-ETF-Summary-v01.docx>
\[9\]
ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/25825a46-ec4c-4544-93e0-5e46d77e1d5a/ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx>

## [Funds/ETF Attributes]{.mark}

**\[TBD: Tor S. to review and update this section\]**

**Funds/ETF Attributes** in the ClearingBid system are a defined set of
characteristics and parameters that distinguish Fund/ETF offerings from
stocks and bonds. These attributes are consistently described across the
documents, with only minor differences or open questions. Below is a
structured summary, followed by an analysis of any contradictions found.

**Key Attributes of Funds/ETFs**

  ------------------------------------------------------------------
  **Attribute**    **Description**            **Source(s)**
  ---------------- -------------------------- ----------------------
  **No Price       Price is set by the **Net  \[1\]\[2\]\[3\]\[4\]
  Discovery**      Asset Value (NAV)**,       
                   calculated externally and  
                   updated in real time.      

  **Issue Size**   There is **no fixed        \[1\]\[2\]\[3\]\[4\]
                   offering size**; the       
                   number of shares is        
                   determined by demand at    
                   NAV.                       

  **Order Types    **Limit GTC**(Good Til     \[1\]\[2\]\[3\]
  Allowed**        Cancelled) and             
                   **MarketOnClose GTC**      
                   orders are accepted.       

  **Allocation**   All orders at or above NAV \[1\]\[2\]\[3\]
                   and all MarketOnClose      
                   orders at closing are      
                   **fully filled**.          

  **Orders Below   Limit orders below NAV are \[1\]\[2\]\[3\]
  NAV**            **cancelled** at closing.  

  **Distribution   A fee (in percentage       \[1\]\[2\]\[3\]\[4\]
  Fee**            points) paid by the issuer 
                   to the Lead Manager, not   
                   deducted from NAV.         

  **Minimum Volume A **TargetVolume**(minimum \[1\]\[2\]\[3\]
  to Close**       volume) must be met for    
                   the offering to close;     
                   otherwise, an alert is     
                   triggered.                 

  **NAV Display**  **CNAV** (Clearing NAV) is \[1\]\[2\]\[3\]\[4\]
                   shown in place of a        
                   Clearing Price.            

  **Order Book     Columns like Low/High      \[1\]\[2\]\[3\]
  Columns**        Price Range Orders and     
                   Sizes are **blank** for    
                   Fund/ETF offerings.        

  **Allocation     Allocation is **not        \[1\]\[2\]\[3\]
  Algorithm**      simulated**; all           
                   qualifying orders receive  
                   100% allocation at NAV.    

  **Distribution   Only visible to the Lead   \[1\]\[2\]\[3\]
  Fee Visibility** Manager.                   

  **Latest NAV     Timestamp of the latest    \[1\]\[2\]\[3\]
  Update**         NAV update is displayed.   
  ------------------------------------------------------------------

**Additional Details**

**Security Types**: Funds/ETFs are identified by standard codes (e.g.,
MF for Mutual Fund, ETF for Exchange Traded Fund) in the SecurityType
167 FIX field\[1\]\[2\]\[3\]\[4\].

**Minimum and Maximum Order Quantities**: Defaults and configuration for
minimum and maximum order sizes are described for all asset types,
including Funds/ETFs\[1\]\[2\].

**Graph and Reporting**: Volume graphs for Funds/ETFs default to showing
volume at CNAV. Market orders and accumulated market orders should be
displayed at the same price level as the iNAV; if iNAV changes, these
move accordingly\[3\].

**Closing Process**: For Funds/ETFs, closing occurs at the Window 1
Start Time, with other closing settings fields disabled\[1\]\[2\]\[3\].

**Contradicting Information or Open Questions**

**1. Allocation Algorithm Naming/Implementation:**

Some sections refer to allocation algorithms (e.g., Pro-rata, Price-Time
Priority), but for Funds/ETFs, simulation and complex allocation are
explicitly *not*performed---all qualifying orders are filled at
NAV\[1\]\[2\]\[3\]. There is a question in one document about whether
more than one Priority Group is needed for Funds/ETFs, but the
implemented logic is that all qualifying orders are filled, making
Priority Groups unnecessary\[3\].

**2. Terminology: \"Clearing Price\" vs. \"CNAV\":**

The term \"Clearing Price\" is replaced by \"CNAV\" for Funds/ETFs, but
in some places, the documentation generically refers to \"Clearing
Price\" for all asset types, which could be confusing. However, all
detailed sections clarify that for Funds/ETFs, CNAV is used
instead\[1\]\[2\]\[3\]\[4\].

**3. Distribution Fee Calculation/Visibility:**

All documents agree that the Distribution Fee is paid directly by the
issuer and is not deducted from NAV, and is only visible to the Lead
Manager\[1\]\[2\]\[3\]\[4\]. There is no contradiction here, but this is
an area where the system\'s transparency differs from stocks/bonds
(where the underwriting spread is more visible).

**4. Order Book Columns:**

There is consistent treatment that Low/High Price Range and related
columns are blank for Funds/ETFs, but the explicit mention of this
varies between documents\[1\]\[2\]\[3\].

**5. Security Type Expansion:**

There is an open question in one document about whether additional
Security Types are needed for Funds, but this is not a
contradiction---just an item for further consideration\[3\].

**Conclusion**

The documents present a **consistent and unified definition** of
Funds/ETF attributes within the ClearingBid platform. The only minor
discrepancies are open questions about future enhancements (such as
additional Security Types or Priority Groups), not contradictions in the
current implementation. All key operational and reference data
attributes---such as price setting by NAV, no price discovery, full
allocation at NAV, and the role of the Distribution Fee---are aligned
across the documentation\[1\]\[2\]\[3\]\[4\].

Sources \[1\]
ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>
\[2\] ClearingBid-Backend-Functions-CMC-v2-7.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/60bd3f30-70db-4db8-9031-579093dfb8b8/ClearingBid-Backend-Functions-CMC-v2-7.docx>
\[3\] Fund-ETF-Summary-v01.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/03b9046e-a9fd-4ffa-9ca1-ea087a21e91a/Fund-ETF-Summary-v01.docx>
\[4\] CB-PDP-Gen-2-specification-updated-definitions-04.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/d0735754-a7b5-4d48-9c69-5fd13fba6142/CB-PDP-Gen-2-specification-updated-definitions-04.docx>\[5\]
CB-Backend-Specification-v1-4-8-comments.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/571d1572-1928-4edd-a5d2-fd3b20d2ba8d/CB-Backend-Specification-v1-4-8-comments.docx>
\[6\] BE-Terminal-Reference-v004-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/bb17a992-9932-4b23-b225-8bb487beab97/BE-Terminal-Reference-v004-2.docx>
\[7\] Auction-Process-ClearingBid-2025-02-09-v2-Draft.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/2723ad3f-6e9d-4c22-a998-b292e737e4e7/Auction-Process-ClearingBid-2025-02-09-v2-Draft.docx>
\[8\] CB-PDP-Gen-2-specification-v2-0r01.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/f0a7c1b2-8ef0-428d-b64c-f20597bfe366/CB-PDP-Gen-2-specification-v2-0r01.docx>
\[9\] SD-PRD\_-Allocation-Integrity-Verification-300625-211230.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/434a11b8-fa30-4fe2-85c8-13e2062fb3d1/SD-PRD_-Allocation-Integrity-Verification-300625-211230.pdf>

## [Upcoming Offering Information]{.mark}

**\[TBD: Tor S. to review and update this section\]**

**Upcoming Offering Information** refers to the set of data, parameters,
and user interface elements that describe and present details about new
securities offerings that are scheduled but not yet open for bidding or
subscription on the ClearingBid platform. This information is crucial
for investors, broker-dealers, issuers, and lead managers to understand
the terms, timeline, and mechanics of an upcoming offering before the
order book opens.

**Key Elements of \"Upcoming Offering Information\"**

**Basic Attributes Displayed**:\
Upcoming offerings are shown in dedicated sections or widgets in both
professional workstations and the public website. The displayed columns
and fields typically include:

**Issuer** (the entity offering the security)

**Name** (of the security/offering)

**Symbol** (Ticker or CUSIP)

**Type** (Stock, Bond with Yield Type, ETF, Fund)

**Announcement Date** (date of public announcement)

**Bidding Opens** (date/time when order entry begins)

**Offering Size** (number of shares, bonds, or units available)

**Low Range / High Range** (expected price or yield range)

**Bidding Ends** (date/time when order entry closes)

**Lead Manager** (underwriter managing the offering)\[1\]\[2\].

**Additional Parameters**:

**Target Volume** (minimum volume for the offering to close)

**Target Orders** (minimum number of orders required)

**Links to Prospectus** and detailed offering documents

**Order Book State** (e.g., NEW, UPCOMING, OPEN, etc.)

**Allocation Algorithm** (set at offering definition, cannot be changed
after opening)\[1\]\[2\].

**User Access and Presentation**:

**Professional Users** (lead managers, broker-dealers) see upcoming
offerings in the BatMan/BatMan app, with the ability to drill down into
offering details and see parameters relevant to their roles\[3\]\[1\].

**Public Website** displays upcoming offerings with limited but
essential information for potential investors, including the ability to
subscribe for updates\[3\]\[1\].

Clicking on an upcoming offering typically opens a window or widget
showing all attributes the user is permitted to see\[1\]\[2\].

**State Management**:

The system tracks offerings through states: NEW (pre-announcement),
UPCOMING (announced but not open), OPEN (bidding open), CLOSE PENDING,
CLOSING, etc. The \"Upcoming\" state is specifically for offerings that
are announced but not yet accepting orders\[1\]\[2\].

**Notifications and Updates**:

Investors and the public can subscribe to receive updates (email, SMS)
about upcoming offerings, including changes to the bid period, price
range, or other material terms\[3\]\[1\]\[4\].

**Contradicting Information Identified**

**1. Field Names and Display Order**

The field names and column order for upcoming offerings differ slightly
between documents:

One specification lists columns as: Issuer, Name, Symbol, Type,
Announcement, Bidding Opens, OfferingSize, Low Range, High Range,
Bidding Ends, Lead Manager\[1\].

Another document uses: Name, Ticker, Market, Time to Close,
Price/Clearing Prices, Market Value, etc., and includes \"Upcoming
Offerings\" as a menu item but does not specify the same column order or
exact field names\[3\].

The Auction Process document (prospectus) does not specify the exact
data fields in the same way but refers to the public website with
hyperlinks to underwriters and offering-specific pages, suggesting a
different presentation\[4\].

**2. Allocation Algorithm Terminology and Setting**

Some documents refer to \"Allocation Algo\" as being set in the Edit
Offering Widget and not changeable after opening\[1\]\[2\].

Others use slightly different terminology (\"Allocation Algorithm type
is set during issue definition and cannot be changed after
Opening\")\[2\].

The number and naming of allocation algorithms (e.g., \"Time-priority
Pro-rata,\" \"Pro-rata,\" \"Price-time Priority,\" \"Price-priority
Pro-rata\") are not always consistent or fully aligned in their
descriptions across documents\[1\]\[2\].

**3. Minimum Quantity and Order Size Parameters**

The minimum and maximum order size parameters, as well as the rules for
\"MinQty\" (Minimum Quantity Condition), are described with different
default values and configuration options in different
documents\[1\]\[2\]. For example:

Default minimum MinQty order size is 100,000 in some
documents\[1\]\[2\].

The Auction Process/prospectus describes the minimum bid size as 100
shares for IPOs\[4\].

The terminology for these parameters (\"MinQty,\" \"Min Order Qty,\"
\"Minimum Quantity Condition\") is not always consistent.

**4. Timing and Notification Details**

The timing for \"Bidding Opens\" and \"Bidding Ends\" can be described
differently:

Technical specifications refer to configurable parameters and default
times (e.g., 4:05 PM for Window 1 Start Time)\[1\]\[2\].

The prospectus/auction process describes the notification process for
opening and closing bidding in terms of Eastern Time and relative to SEC
effectiveness, which may not align exactly with the technical
defaults\[4\].

**5. Scope of Information Shown to Users**

Some documents suggest that all attributes are visible to users
permitted to see them, while others imply that certain fields (e.g.,
Gross Underwriting Spread, Selling Concession) are only visible to lead
managers and not to all users\[1\]\[2\].

**Conclusion**

**Upcoming Offering Information** on the ClearingBid platform
encompasses a comprehensive set of parameters, states, and user
interface elements that describe new securities offerings before the
order book opens. While the core concepts are consistent across
documentation, there are **minor contradictions** in field names,
display order, allocation algorithm terminology, default parameter
values, and the scope of information visible to different user types.
These inconsistencies should be harmonized in future documentation and
user interface design to ensure clarity and a seamless user
experience\[3\]\[1\]\[4\]\[2\].

Sources \[1\] ClearingBid-Backend-Functions-CMC-v2-7.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/60bd3f30-70db-4db8-9031-579093dfb8b8/ClearingBid-Backend-Functions-CMC-v2-7.docx>
\[2\] ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>
\[3\] CB-Backend-Specification-v1-4-8-comments.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/571d1572-1928-4edd-a5d2-fd3b20d2ba8d/CB-Backend-Specification-v1-4-8-comments.docx>
\[4\] Auction-Process-ClearingBid-2025-02-09-v2-Draft.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/2723ad3f-6e9d-4c22-a998-b292e737e4e7/Auction-Process-ClearingBid-2025-02-09-v2-Draft.docx>
\[5\] BE-Terminal-Reference-v004-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/bb17a992-9932-4b23-b225-8bb487beab97/BE-Terminal-Reference-v004-2.docx>
\[6\] SD-PRD\_-PDP-Reporting-for-Offering-Close-300625-211034.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/8bb11d92-8e11-4996-9592-e104e7b21a66/SD-PRD_-PDP-Reporting-for-Offering-Close-300625-211034.pdf>
\[7\] BidMan-User-Guide-v0-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/e5243159-49ba-46a2-9c45-c585d64ce1a1/BidMan-User-Guide-v0-2.docx>
\[8\]
SD-PRD\_-Offering-Dashboard-For-Issuer-and-Lead-Manager-010725-151542.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/8978878e-2a16-4319-8a5d-6e0fe6c6b162/SD-PRD_-Offering-Dashboard-For-Issuer-and-Lead-Manager-010725-151542.pdf>

## [Closed Offering Information]{.mark}

**\[TBD: Tor S. to review and update this section\]**

**Closed Offering Information** in the ClearingBid platform refers to
the data and reporting available after a securities offering has
completed its auction and allocation process. This information is
presented to different user types (Lead Managers, Issuers,
Broker-Dealers, Market Operations) through the BatMan/BidMan/MarkOps app
and is also used for post-offering processes such as commission
calculation and reporting.

**Key Elements of Closed Offering Information**

**1. Data Fields Displayed for Closed Offerings**

For both Lead Managers/Issuers and Broker-Dealers, the BatMan/BidMan app
displays the following fields for each closed offering:

**Issuer**: Name of the company issuing the security

**Name**: Full name of the offering

**Symbol**: Exchange-assigned symbol

**Market**: Primary listing exchange (sometimes referred to as
\"Market\" or \"Listing\"; terminology is still under review)

**Type**: Type of offering (e.g., common stock, corporate bond)

**Offering Size**: Number of units or dollar value sold

**Last Sale**: *TBD if this is the same as Offering Price*

**Offering Price**: The price at which the security was sold to
investors

**CLRP (Clearing Price)**: The price at which the full available
quantity was subscribed

**Low Range/High Range**: The estimated price range for the offering

**Lead Manager**: Lead manager for the offering

**State**: Closed or cancelled

**Offering Closed**: Date the offering was closed

**Action**: Allocation information\[1\]\[2\]\[3\]

The number of visible columns can be customized by the user.

**2. Reporting and Export**

After a deal is closed, Market Operations and Lead Managers can generate
CSV reports summarizing allocations, trades, commissions, and other key
parameters.

Reports include fields such as symbol, quantity executed, price per
share, timestamp, client order ID, broker name, account, total order
quantity filled, leaves (unfilled quantity), priority group, fill
percentages, and MPID of the order submitter.

Calculations for commissions (selling concession, gross underwriting
spread) are included, and reports are used for reconciliation and
payment processing\[4\].

**3. Allocation and State Management**

The allocation algorithm and offering parameters used during the auction
are displayed for reference but cannot be changed after the offering is
open or closed.

The system tracks and displays the final allocation per order, including
fill percentages and priority group assignments.

The order book for a closed offering is locked, and no further order
entry, modification, or cancellation is allowed\[1\]\[2\]\[3\].

**4. User Access**

Lead Managers and Issuers have read-only access to closed offering
details.

Broker-Dealers can view allocations and their own order outcomes.

Market Operations has full access for oversight and reporting\[1\].

**Contradicting or Inconsistent Information**

**1. \"Last Sale\" vs. \"Offering Price\"**

In several places, the field \"Last Sale\" is marked as \"TBD\" or with
a question about whether it is the same as \"Offering
Price\"\[1\]\[2\]\[3\]. The documentation does not clearly resolve
whether these fields are always equivalent or if \"Last Sale\" might
represent a different value (e.g., a final trade in a secondary market
or a different settlement price).

**2. Terminology: \"Market\" vs. \"Listing\"**

The terminal documentation notes uncertainty over whether to use
\"Market\" or \"Listing\" to describe the primary exchange field,
indicating a lack of standardization in terminology\[1\]\[2\]\[3\].

**3. Allocation Algorithm Field**

There is inconsistency in where the allocation algorithm is set and
displayed. Some documents state the allocation algorithm is set in the
\"New/Edit Offering\" widget and is only informational in the Allocation
widget, and cannot be changed after opening\[2\]\[3\]. However, this is
not always clearly enforced or described in the workflow documentation.

**4. Reporting and Export**

There is some ambiguity about which columns and fields are mandatory in
exports and which are optional or configurable, as well as how to handle
extremely large datasets (e.g., over 30,000 rows)\[4\].

**5. Commission Eligibility**

There is uncertainty in the documentation about determining commission
eligibility for certain participants (e.g., RIAs vs. Broker-Dealers),
which may affect the accuracy of closed offering reports and payment
processing\[4\].

**6. State Definitions**

While the state transitions for offerings (e.g., CLOSED, CANCELLED,
FROZEN, HALTED) are generally consistent, there are minor differences in
how some documents describe allowable state transitions and
exceptions\[1\]\[2\]\[3\].

**Summary Table: Closed Offering Information Fields**

  ---------------------------------------------------------------
  **Field**        **Description**     **Notes/Contradictions**
  ---------------- ------------------- --------------------------
  Issuer           Company issuing the 
                   security            

  Name             Full name of the    
                   offering            

  Symbol           Exchange-assigned   
                   symbol              

  Market/Listing   Primary listing     \"Market\" vs. \"Listing\"
                   exchange            terminology TBD

  Type             Type of security    
                   (stock, bond, etc.) 

  Offering Size    Total units or      
                   dollar value sold   

  Last Sale        *TBD if same as     Not clearly defined
                   Offering Price*     

  Offering Price   Price paid by       
                   investors           

  CLRP             Clearing Price      

  Low/High Range   Estimated price     
                   range               

  Lead Manager     Lead underwriter    

  State            Closed/Cancelled    

  Offering Closed  Date of closing     

  Action           Allocation          
                   information         
  ---------------------------------------------------------------

**In conclusion:**\
Closed Offering Information in ClearingBid consists of a standardized
set of fields summarizing the outcome of an offering, with reporting and
export functionality for reconciliation and payment. However, there are
minor inconsistencies in field definitions (notably \"Last Sale\" vs.
\"Offering Price\"), terminology (\"Market\" vs. \"Listing\"), and some
workflow details regarding allocation algorithms and reporting, which
should be resolved for full clarity and operational
consistency\[1\]\[2\]\[3\]\[4\].

Sources \[1\] BE-Terminal-Reference-v004-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/bb17a992-9932-4b23-b225-8bb487beab97/BE-Terminal-Reference-v004-2.docx>
\[2\] ClearingBid-Backend-Functions-CMC-v2-7.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/60bd3f30-70db-4db8-9031-579093dfb8b8/ClearingBid-Backend-Functions-CMC-v2-7.docx>
\[3\] ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>
\[4\] SD-PRD\_-PDP-Reporting-for-Offering-Close-300625-211034.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/8bb11d92-8e11-4996-9592-e104e7b21a66/SD-PRD_-PDP-Reporting-for-Offering-Close-300625-211034.pdf>
\[5\] CB-Backend-Specification-v1-4-8-comments.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/571d1572-1928-4edd-a5d2-fd3b20d2ba8d/CB-Backend-Specification-v1-4-8-comments.docx>
\[6\] Auction-Process-ClearingBid-2025-02-09-v2-Draft.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/2723ad3f-6e9d-4c22-a998-b292e737e4e7/Auction-Process-ClearingBid-2025-02-09-v2-Draft.docx>
\[7\] BidMan-User-Guide-v0-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/e5243159-49ba-46a2-9c45-c585d64ce1a1/BidMan-User-Guide-v0-2.docx>
\[8\]
SD-PRD\_-Offering-Dashboard-For-Issuer-and-Lead-Manager-010725-151542.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/8978878e-2a16-4319-8a5d-6e0fe6c6b162/SD-PRD_-Offering-Dashboard-For-Issuer-and-Lead-Manager-010725-151542.pdf>

# System Features and Parameters 

**Core System Features**

**Broker-Dealer Integration:** Broker-dealers connect via standard
protocols (FIX), with options for manual entry through the dedicated
BidMan app. The system supports both automated and manual workflows,
including order management and allocation downloads\[1\]\[3\]\[2\].

## User Management and Authorization

**\[TBD: Tor S. to review and update this section. Note that there are
some wording about user management also in the BidMan and BatMan
section\]**

**User Management and Authorization in ClearingBid\'s System**

User management and authorization in ClearingBid\'s platform is
structured around distinct user roles with specific permissions,
authentication protocols, and administrative controls. The system
supports multiple user types, each with tailored access rights to ensure
secure and efficient operations. Below is a synthesis of the key
components based on the provided documentation:

**Key Components**

**User Roles and Permissions**:

**Lead Managers (LMs)**:

Full access to create, monitor, and allocate new offerings.

Can simulate allocations, define offering parameters (e.g., price
ranges, bid periods), and submit closing commands\[1\]\[2\]\[3\].

Access to real-time order book data, selling concession calculations,
and allocation tools\[2\]\[5\].

**Issuers (ISs)**:

Read-only access to LM functions (e.g., viewing offering details, order
book states)\[1\]\[3\].

**Broker-Dealers (BDs)**:

Order entry, cancellation, and market data access.

Restricted to managing their own orders and accounts; cannot modify
offerings\[1\]\[4\].

Trade Admins (within BD firms) can create user accounts and assign
permissions (read, read-write, or no access)\[4\].

**ClearingBid Market Operations (MO)**:

Full system control, including manual state changes (e.g.,
halting/freezing order books), canceling orders, and managing
offerings\[1\]\[3\]\[6\].

Can override automated processes (e.g., restart closing
windows)\[3\]\[6\].

**Super Users**:

**Market Surveillance Users**: Access all data for monitoring.

**LM/Broker-Dealer Super Users**: Permissions for all offerings under
their domain\[2\]\[3\].

**Authentication and Security**:

**Login**: Initial credentials provided by MO; users can change
passwords and enable two-factor authentication (2FA) via the hamburger
menu\[1\]\[4\].

**Session Management**: Username displayed in the UI; logout/password
changes initiated from the same menu\[1\]\[4\].

**2FA**: Requires a mobile device with a compatible app (e.g., Google
Authenticator)\[4\].

**Account and Access Management**:

**Trade Admins (BD Firms)**:

Create/edit accounts and assign user permissions (e.g., read-write for
traders, no access for restricted roles)\[4\].

Manage visibility of accounts to traders within their firm\[4\].

**Blocked Broker-Dealers**: LMs can block specific BDs from order entry
or market data access\[2\]\[3\].

**State and Permission Controls**:

**Order Book States**: User actions are restricted by system states
(e.g., no cancellations during CLOSING; no order entry in FROZEN
state)\[3\]\[6\].

**Authorization During Closing**: MO can manually halt closing
processes; LMs can only act within predefined states (e.g., cannot move
from CLOSED to OPEN)\[3\]\[6\].

**Contradictions in Documentation**

**FIX MinQty Orders**:

\[1\] and \[3\] state that MinQty orders are allowed only via the
workstation (not FIX) when the setting is NO. However, \[2\] mentions
FIX field support for MinQty without clarifying this restriction,
creating ambiguity about FIX compatibility\[1\]\[2\]\[3\].

**Allocation Algorithm Flexibility**:

\[2\] claims the allocation algorithm (e.g., Price-Time Priority) can be
set during offering definition and \"cannot be changed after Opening.\"
Conversely, \[3\] implies LM can simulate allocations with different
algorithms post-opening, suggesting flexibility\[2\]\[3\].

**Bond Price Display**:

\[3\] specifies that bond yields can be displayed as percentages, dollar
prices, or per FIX PriceType. \[4\]'s BidMan guide only mentions dollar
prices, omitting percentage/yield options\[3\]\[4\].

**Super User Roles**:

\[2\] and \[3\] define \"Lead Manager Super Users\" with permissions for
all offerings. \[1\]'s terminal (BatMan) reference lacks this role,
implying a documentation gap between backend and frontend
specs\[1\]\[2\]\[3\].

**Summary**

ClearingBid's user management enforces strict role-based access, with
Trade Admins handling BD-level permissions and MO overseeing global
controls. Contradictions arise in FIX order handling, allocation
algorithm flexibility, bond display formats, and super user definitions.
These inconsistencies highlight areas where documentation alignment is
needed, particularly between functional specifications (e.g.,
\[2\]\[3\]) and user guides (e.g., \[4\]).

Sources \[1\] BE-Terminal-Reference-v004-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/bb17a992-9932-4b23-b225-8bb487beab97/BE-Terminal-Reference-v004-2.docx>
\[2\] ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>
\[3\] ClearingBid-Backend-Functions-CMC-v2-7.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/60bd3f30-70db-4db8-9031-579093dfb8b8/ClearingBid-Backend-Functions-CMC-v2-7.docx>
\[4\] BidMan-User-Guide-v0-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/e5243159-49ba-46a2-9c45-c585d64ce1a1/BidMan-User-Guide-v0-2.docx>
\[5\] CB-system-high-level-description-v4-04.04.22.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/a869d44e-1ca0-40fb-a4a6-07b3b15461c3/CB-system-high-level-description-v4-04.04.22.docx>
\[6\] Bust-Management-010725-142822.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/c81baec6-73d0-49b2-8f55-728a24511839/Bust-Management-010725-142822.pdf>
\[7\]
SD-PRD\_-Gross-Spread-and-Selling-Concession-Calculations-010725-143301.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/79b09419-72ed-42dc-9015-5e01f567d190/SD-PRD_-Gross-Spread-and-Selling-Concession-Calculations-010725-143301.pdf>
\[8\] SD-PRD\_-PDP-Reporting-for-Offering-Close-300625-211034.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/8bb11d92-8e11-4996-9592-e104e7b21a66/SD-PRD_-PDP-Reporting-for-Offering-Close-300625-211034.pdf>

## Default Values Table

**\[TBD: Tor S. to review and update this section\]**

**Default Values Table** refers to a structured method for managing
system-wide default parameters for new offerings in the ClearingBid
platform. This table is designed to be maintained by ClearingBid Market
Operations and contains columns for each asset type: **Stock**,
**Bond**, and **FundETF**. Each column specifies the default value for a
range of parameters relevant to that asset class. The intention is to
centralize and standardize the configuration of key parameters, ensuring
consistency and simplifying operational management\[1\]\[2\].

**Key Parameters in the Default Values Table:**

**Number of Decimals:** Default is 2 for Stocks, ETFs, and Funds; 3 for
Bonds; maximum is 8.

**Price Increment:** Default is 0.01 for Stocks, ETFs, and Funds; 0.001
for Bonds.

**Minimum MinQty Order Size:** Default is 100,000 (configurable
intraday).

**Maximum MinQty Order Size (% of order size):** Default is 50 for all
Allocation/Priority Groups except the last (which defaults to 0).

**FIX MinQty Orders Allowed:** Default is NO (configurable intraday).

**Cancel Time Limit:** Configurable per Allocation/Priority Group;
default for Prio 1 is 1440 minutes (24 hours), others default to 0.

**Minimum Bid Quantity Allowed:**Default is 10 (configurable intraday).

**Maximum Bid Quantity Allowed:**Configurable intraday.

**Minimum Price/Maximum Yield Allowed:** Configurable intraday.

**Maximum Price/Minimum Yield Allowed:** Configurable intraday.

**Quantity Increment:** Default is 1 (configurable intraday).

**Target Orders (min number of orders):**Default is 0 (no minimum).

**Max Quantity per order/account:**Issuer-defined, configurable.

**Min number of individual order accounts:** Issuer-defined, must be met
for closing to start.

**Opening Clearing Price Method:**Options include Not Published Manual
On (default), Not Published Auto On, Low Price Range, and Price-Time
Priority.

**Allocation Algorithm:** Set during issue definition; cannot be changed
after opening\[1\]\[2\].

**Purpose and Usage:**

The Default Values Table is intended to ensure that every new offering
starts with a consistent set of parameters, which can then be adjusted
as needed by Market Operations.

These defaults are especially important for operational efficiency and
regulatory compliance, as they help prevent errors and ensure that
offerings meet predefined standards\[1\]\[2\].

**Contradicting Information Identified:**

**Terminology for Priority Groups:**

In some documents, the term **\"Allocation Priority Group\"** is used,
while in others it is referred to as **\"Priority Group\"**. The
functional meaning is the same, but the inconsistent terminology could
cause confusion\[1\]\[2\].

**Default Values for Cancel Time Limit:**

Both documents specify that the default for the Time-priority Pro-rata
Prio 1 group is 1440 minutes (24 hours), and for other groups it is 0.
However, in one document, the field is called \"Cancel Prohibited\" and
in another, \"Cancel Time Limit\"\[1\]\[2\]. The meaning is consistent,
but the label differs.

**Price Increment for Bonds:**

The main documents specify a default price increment of 0.001 for
Bonds\[1\]\[2\]. However, the Fixed Income Prototype document suggests
that when using basis points as the price format, the default price
increment might be 1 (i.e., 1 basis point), or possibly configurable
depending on the format\[3\]. This could lead to different default
increments for Bonds depending on the selected price display mode.

**Order Quantity Defaults for Bonds:**

The Fixed Income Prototype document sets the **Order Min Qty Size**
default for bonds to 10,000, while the main documents set the **Minimum
Order Bid Quantity Allowed**default to 10 for all asset
types\[1\]\[2\]\[3\]. This is a significant discrepancy and should be
clarified.

**Target Volume/Target Orders:**

The main documents define **TargetOrders** (minimum number of orders)
with a default of 0, and **TargetVolume** (minimum volume to close) for
FundETF with a default of 0\[1\]\[2\]. The Fixed Income Prototype sets
**Target Volume**default to 100,000,000 and **Total Qty** default to
100,000 for bonds\[3\]. These are much higher than the general defaults
and may reflect prototype-specific values rather than system-wide
defaults.

**Yield Type Default:**

The Fixed Income Prototype sets the **Yield Type** default to SPT-SM
(Spread to Treasury with Similar Maturity)\[3\], while the main
documents do not specify a default Yield Type\[1\]\[2\]. This could lead
to confusion about what the system should default to for new bond
offerings.

**Summary Table of Key Default Values (with Contradictions
Highlighted):**

  ---------------------------------------------------------------------------------
  **Parameter**   **Main Docs Default      **Fixed       **Notes/Contradictions**
                  (Stock/Bond/FundETF)**   Income        
                                           Prototype     
                                           Default       
                                           (Bond)**      
  --------------- ------------------------ ------------- --------------------------
  Number of       2 / 3 / 2                Not specified Consistent
  Decimals                                               

  Price Increment 0.01 / 0.001 / 0.01      1 (bps mode), Contradicts for Bond in
                                           0.01 (other)  bps mode

  Min Order Qty   10                       10,000        Contradicts for Bond
  Allowed                                                

  Max Order Qty   Configurable             100,000       Prototype sets explicit
  Allowed                                                default

  Target Orders   0                        Not specified 

  Target Volume   0 (FundETF)              100,000,000   Contradicts for Bond

  Cancel Time     1440 min (24h)           Not specified Consistent
  Limit (Prio 1)                                         

  Yield Type      Not specified            SPT-SM        Prototype sets explicit
                                                         default
  ---------------------------------------------------------------------------------

**Conclusion:**\
The **Default Values Table** is a central configuration mechanism for
setting initial parameters for new offerings, with defaults defined per
asset type and managed by Market Operations. While the core
documentation is consistent on most defaults, the Fixed Income Prototype
introduces some conflicting default values for bonds, particularly
regarding order size, price increment (when using basis points), target
volume, and yield type. These discrepancies should be reconciled to
avoid operational confusion and ensure system
consistency\[1\]\[2\]\[3\].

Sources \[1\] ClearingBid-Backend-Functions-CMC-v2-7.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/60bd3f30-70db-4db8-9031-579093dfb8b8/ClearingBid-Backend-Functions-CMC-v2-7.docx>
\[2\] ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>
\[3\] SD-PRD\_-PDP-Changes-for-Fixed-Income-Prototype-010725-151230.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/ffaad82f-212c-4181-bfa8-634619e12be7/SD-PRD_-PDP-Changes-for-Fixed-Income-Prototype-010725-151230.pdf>
\[4\] BidMan-User-Guide-v0-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/e5243159-49ba-46a2-9c45-c585d64ce1a1/BidMan-User-Guide-v0-2.docx>
\[5\] Preferential-Bids-v-1-2-redlined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/62e01028-b0a8-48b7-afed-0355e3213d05/Preferential-Bids-v-1-2-redlined.docx>
\[6\] SD-PRD\_-Investor-Account-Fix-to-PDP-Mapping-010725-143448.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/b94ec060-0019-4d13-abb3-1f910f16bb8c/SD-PRD_-Investor-Account-Fix-to-PDP-Mapping-010725-143448.pdf>
\[7\] CB-Backend-Specification-v1-4-8-comments.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/571d1572-1928-4edd-a5d2-fd3b20d2ba8d/CB-Backend-Specification-v1-4-8-comments.docx>
\[8\] CB-PDP-Gen-2-specification-v2-0r01.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/f0a7c1b2-8ef0-428d-b64c-f20597bfe366/CB-PDP-Gen-2-specification-v2-0r01.docx>
\[9\] CB-BE-Functions-CMC-MinQTy-add-v4-3.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/2fec3715-234a-4477-9d70-156b95a27452/CB-BE-Functions-CMC-MinQTy-add-v4-3.docx>
\[10\] ClearingBid-IPO-Price-Discovery-and-Allocation-v0-0-8.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/4cb22cd3-2d63-4d83-bd30-d28c1f0d2bde/ClearingBid-IPO-Price-Discovery-and-Allocation-v0-0-8.docx>

# Interfaces

## FIX

**\[TBD: Tor S. to review and update this section\]**

**Key Features and Usage in ClearingBid**

**Order Entry and Management**: Broker-dealers submit orders to the PDP
using FIX messages. The FIX protocol supports various message types
(e.g., New Order Single, Execution Report, Order Cancel
Request)\[1\]\[2\]\[5\]\[3\].

**Market Data**: Market data, including order book updates and offering
changes, is disseminated via FIX messages\[2\].

**Session Management**: FIX sessions are established using standard
parameters (SenderCompID, TargetCompID, network address). Authentication
and session messages (Logon, Heartbeat, Logout, etc.) follow FIX 4.2
conventions\[1\]\[2\].

**Connectivity Options**: Participants can connect via VPN, secure
internet, dedicated lines, or colocation for reduced
latency\[1\]\[2\]\[3\].

**Integration and Certification**: Onboarding involves a certification
process with the FIX engineering team to ensure compatibility and
message integrity\[1\]\[2\].

**FIX in the ClearingBid Workflow**

**Broker-Dealer Participation**: Broker-dealers can connect directly,
through connectivity providers, or via execution/clearing service
providers. Each method uses specific FIX tags to identify parties
(SenderCompID, OnBehalfOfCompID, SenderSubID, etc.)\[5\].

**Order and Execution Flow**: Orders flow from investors →
broker-dealers → ClearingBid PDP via FIX. Allocations and confirmations
are sent back using FIX messages\[5\]\[3\].

**Investor Identification**: The Account field (FIX Tag 1) and
SenderSubID are used for identifying the origin of orders, especially
for features like Preferential Bids and tracked orders\[6\]\[7\]\[4\].

**Summary Table: FIX in ClearingBid**

  -------------------------------------------------------------
  **Aspect**       **Description**       **Reference**
  ---------------- --------------------- ----------------------
  Protocol Version FIX 4.2               \[1\]\[2\]

  Core Use Cases   Order entry,          \[1\]\[2\]\[5\]\[3\]
                   execution reporting,  
                   market data           
                   dissemination         

  Message Types    New Order Single,     \[1\]\[2\]\[5\]
  Supported        Execution Report,     
                   Order Cancel, Market  
                   Data Request, etc.    

  Connectivity     VPN, Internet,        \[1\]\[2\]\[3\]
  Options          dedicated lines,      
                   colocation            

  Participant      SenderCompID,         \[6\]\[5\]\[7\]\[4\]
  Identification   OnBehalfOfCompID,     
                   SenderSubID, Account  
                   (Tag 1)               

  Integration      Certification test    \[1\]\[2\]
  Process          with FIX team,        
                   supports standard FIX 
                   engines               
  -------------------------------------------------------------

**Contradicting or Inconsistent Information**

**1. Use of the Account Field (Tag 1) and Investor Identification**

Some documents state that the **Account field (Tag 1)** in FIX is used
to identify bids from the same investor, particularly for Preferential
Bids and tracked orders\[6\]\[7\]\[4\].

However, it is clarified in multiple places that there is **no direct
relationship between the Account field used in BidMan (the GUI tool) and
the SenderSubID used in FIX messages**. The Broker-Dealer must choose
either BidMan or their own OMS for order entry, and the mapping of
investor identity differs between these methods\[6\]\[7\].

The Participant System Guidelines imply the Account field in FIX is
optional and only used for certain scenarios, while other documents
emphasize its necessity for investor grouping in Preferential
Bids\[6\]\[7\]\[4\].

**2. FIX Tag Mapping and Order Flow**

The mapping of FIX tags (SenderCompID, OnBehalfOfCompID, SenderSubID,
etc.) varies slightly between documents describing direct broker-dealer
participation, connectivity services, and execution/clearing
providers\[5\]. While the core logic is consistent, the
required/optional status of certain tags (e.g., SenderSubID,
OnBehalfOfSubID) is not always clearly aligned.

**3. Message Validation and Handling of Unsupported Fields**

Both the trade and market data API manuals state that **unsupported
fields in FIX messages are not validated and do not cause rejection if
all necessary tags are present**\[1\]\[2\]. However, there is no
explicit documentation on how this might affect downstream processes or
data integrity, and some sections suggest stricter validation may apply
for certain message types.

**4. Certification and Integration Details**

While most documents mention a **certification process** for FIX
integration, the specifics (such as which message types are tested, what
constitutes a pass/fail, etc.) are not consistently detailed across all
sources\[1\]\[2\].

**In summary:**\
The documentation consistently describes FIX as the core protocol for
order and data exchange in ClearingBid, with version 4.2 as the
standard. Minor inconsistencies exist in the treatment of the Account
field and the mapping of identity-related FIX tags between order entry
methods. These should be clarified to avoid confusion, especially for
broker-dealers integrating with the
system\[1\]\[2\]\[6\]\[5\]\[3\]\[7\]\[4\].

Sources \[1\] CBFIXtrade4-2v1-4.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/a37dba8b-f50b-49f3-8128-f1b2c5be5fba/CBFIXtrade4-2v1-4.docx>
\[2\] FIXmarketData4-2v1-4.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/f78db1de-58c3-4189-bcf0-11b67f1ad96b/FIXmarketData4-2v1-4.docx>\[3\]
CB-system-high-level-description-v4-04.04.22.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/a869d44e-1ca0-40fb-a4a6-07b3b15461c3/CB-system-high-level-description-v4-04.04.22.docx>
\[4\]
ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/25825a46-ec4c-4544-93e0-5e46d77e1d5a/ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx>
\[5\] SD-FIX-Message-Flow-010725-150358.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/b7cd962b-13de-4b2c-b592-ded864024d62/SD-FIX-Message-Flow-010725-150358.pdf>
\[6\] SD-PRD\_-Investor-Account-Fix-to-PDP-Mapping-010725-143448.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/b94ec060-0019-4d13-abb3-1f910f16bb8c/SD-PRD_-Investor-Account-Fix-to-PDP-Mapping-010725-143448.pdf>
\[7\]
SD-Investor-identification-for-Preferential-Bids-and-Tracked-Orders-010725-143948.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/6c70c9d9-e314-475e-9b8c-159f27f7e0df/SD-Investor-identification-for-Preferential-Bids-and-Tracked-Orders-010725-143948.pdf>
\[8\] ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>

## ClearingBid Web-Portal API

**\[TBD: Tor S. to review and update this section\]**

**ClearingBid Web-Portal API Overview**

**Purpose and Role**

The ClearingBid Web-Portal AP is a proprietary API designed to provide
real-time information about new issue securities offerings managed on
the ClearingBid platform. Its core purpose is to provide ClearingBid's
web portal and mobile applications with live data on IPOs and other
offerings.

**Key Features**

**Real-Time Market Data:** Publishes bids, indicative clearing prices,
demand curves, and other offering metrics in real time, allowing all
market participants to monitor price discovery and demand throughout the
marketing period\[2\]\[1\].

**Public Accessibility:** The API is accessible via the Internet,
enabling public users (not just broker-dealers or institutional
participants) to view offering information, statistics, and updates
without requiring secure network access\[1\].

**Integration with Web and Mobile:** The API powers the ClearingBid
public web portal and mobile app, providing a consistent data feed for
both platforms\[1\].

**Offering Information:** Provides detailed data on current, upcoming,
and past offerings, including instrument details, order book status,
price ranges, and allocation statistics\[1\].

**Subscription and Alerts:** Enables investors and the public to
subscribe to updates, such as new issue announcements, bid period
changes, and indicative Clearing Price alerts, via email or text
message\[1\].

**Reference Data:** Supplies prospectuses, research, roadshow materials,
and other documents relevant to the offering, ensuring that all users
have access to the same information\[1\].

**System Architecture**

**Separation of Networks:** The Web-Portal API operates over a public
network, distinct from the secure VPN and dedicated connections used for
broker-dealer FIX transactions and BidMan/BatMan access\[2\].

**Frontend-Backend Integration:** The API serves as the bridge between
the backend order book and the public-facing frontend, ensuring that all
published data accurately reflects the state of the master order
book\[1\].

**Cloud Hosting:** The public frontend application and its API are
hosted by a cloud provider, supporting scalability and broad
accessibility\[1\].

**Supported Interfaces**

**Web Portal and Mobile App:** Main consumers of the Public API,
providing dashboards and interactive views for retail and institutional
investors\[1\].

**API for Developers:** While primarily intended for internal use by
ClearingBid's web and mobile platforms, the API may also support
third-party integrations for partners or data vendors, subject to
ClearingBid's policies\[1\].

**Contradicting Information Identified**

**1. API Scope and Access**

**Public vs. Restricted Access:** Some documents describe the Public API
as strictly for public, read-only data dissemination, while others
mention the possibility of API access for registered users or partners
to exchange information or even submit data (e.g., user registration,
alerts)\[1\]. The degree of openness and write-access capabilities is
not always clear.

*Example:* One source states, \"Users of the site will be invited to
register to be able to exchange information with ClearingBid and other
users,\" which could imply more interactive API features, while others
restrict the API to data publishing only\[1\].

**2. Data Consistency and Timing**

**Real-Time vs. Delayed Data:** All sources agree the Public API
delivers real-time data, but some technical descriptions note that
certain metrics (such as indicative clearing prices or order book depth)
may be subject to update intervals or lead manager controls, potentially
resulting in brief delays or manual intervention\[2\]\[1\].

*Example:* The lead manager may choose to turn on public publishing of
the indicative clearing price either automatically or manually, which
could lead to discrepancies in what is visible through the API at any
given moment\[2\].

**3. API Functionality for Order Information**

**Order Detail Granularity:** Some documentation suggests that the
Public API provides only aggregated order book data (e.g., price levels,
total demand), while others indicate that more granular information,
such as individual order details or broker-dealer breakdowns, may be
available to certain users or under specific conditions\[1\].

*Example:* The backend specification lists \"individual orders or firm,
per order\" as available data, but it is unclear if this is exposed via
the public API or only to authenticated broker-dealer terminals\[1\].

**Summary Table: Public API Capabilities and Document Contradictions**

  --------------------------------------------------------------------
  **Feature/Aspect**   **Description   **Contradiction/Uncertainty**
                       (Consensus)**   
  -------------------- --------------- -------------------------------
  Data Scope           Real-time,      Write-access and registration
                       public offering features unclear\[1\]
                       data            

  Data Granularity     Aggregated      Individual order details
                       order book,     possibly available\[1\]
                       price levels    

  Data Timing          Real-time       Manual lead manager controls
                       updates         may introduce delays\[2\]\[1\]

  Access Method        Web portal,     Third-party/partner API access
                       mobile app, API not consistently defined\[1\]
                       endpoints       

  Security/Network     Public internet Some references to registration
                       access          or restricted access\[1\]
  --------------------------------------------------------------------

**References:**\
\[2\]: CB-system-high-level-description-v4-04.04.22.docx\
\[1\]: CB-Backend-Specification-v1-4-8-comments.docx

Sources \[1\] CB-Backend-Specification-v1-4-8-comments.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/571d1572-1928-4edd-a5d2-fd3b20d2ba8d/CB-Backend-Specification-v1-4-8-comments.docx>
\[2\] CB-system-high-level-description-v4-04.04.22.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/a869d44e-1ca0-40fb-a4a6-07b3b15461c3/CB-system-high-level-description-v4-04.04.22.docx>
\[3\] CB-PDP-Gen-2-specification-v2-0r01.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/f0a7c1b2-8ef0-428d-b64c-f20597bfe366/CB-PDP-Gen-2-specification-v2-0r01.docx>
\[4\] ClearingBid-New-Issue-Auction-Process-v-1-2-8.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/22cd9570-1ab0-4c3f-b33a-16a854e18019/ClearingBid-New-Issue-Auction-Process-v-1-2-8.docx>
\[5\] ClearingBid-Backend-Functions-CMC-v2-7.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/60bd3f30-70db-4db8-9031-579093dfb8b8/ClearingBid-Backend-Functions-CMC-v2-7.docx>
\[6\]
ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/25825a46-ec4c-4544-93e0-5e46d77e1d5a/ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx>
\[7\] ClearingBid-IPO-Price-Discovery-and-Allocation-v0-0-8.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/4cb22cd3-2d63-4d83-bd30-d28c1f0d2bde/ClearingBid-IPO-Price-Discovery-and-Allocation-v0-0-8.docx>
\[8\] ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>

## BidMan -- Broker-dealer App

**\[TBD: Tor S. use this a checklist and see if anything needs to be
included in section 4.2 -- BidMan. Then delete this section\]**

**Overview of the BidMan - Broker Dealer App**

The **BidMan application** is designed for broker-dealers as either a
complement to, or replacement for, their existing order management
systems. It serves as a user interface to ClearingBid's backend systems,
allowing users to view market data, manage orders, and access reports.
The app can function as a read-only viewer or as the primary trading
interface, depending on the level of integration with the
broker-dealer's
infrastructure[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx).

**Key Features**

**Login & Security**

Access via browser at a specified URL.

Supports all modern browsers.

Two-factor authentication (2FA) can be enabled by users.

Password management and reset options.

**Market Data**

Clicking on a symbol displays its order book and related orders.

Order book widget visualizes price ranges and clearing price.

Offering details are accessible via the symbol's name.

**Order Entry & Management**

Orders can be placed by clicking on a bid price in the order book.

Suggested quantity defaults to the bid size; preset quantities are
editable.

Orders can be canceled directly from the order book or the \"My Orders\"
widget.

Widget reference sections provide detailed order entry and cancellation
instructions.

**Reporting**

Orders and trades reports can be printed from the \"My Orders\" widget
(noted as \[TBD\] in the documentation).

**User Interface Customization**

Widgets can be added, moved, and customized within the dashboard.

Tabs allow for different widget layouts.

Widgets are color-coded and linked for contextual navigation.

**Account & User Management**

Trader Admins can create and edit accounts and manage user access
rights.

User rights include no access, read access, or read-write access.

**Other Utilities**

Online status indicator.

Full screen toggle.

Drop-down menus for navigation and account management.

**Data Displayed**

**Open Offerings**

Issuer, Name, Symbol, Listing, Type, Offering Size, Priority Bid End,
Closing, Clearing Price, MyCLRP#Ord, and other metrics.

**Closed Offerings**

Issuer, Name, Symbol, Market/Listing (terminology under review), Type,
Offering Size, Last Sale (definition TBD), Offering Price, Clearing
Price, Low/High Range, Lead Manager, State, Offering Closed, Action.

**Upcoming Offerings**

Issuer, Name, Symbol, Market/Listing (terminology under review), Type,
Bidding Opens, Low/High Range, Bidding Ends, Offering Size, Issue Size
(noted as N.A. and TBD), Lead Manager, State.

**Order Book & My Orders**

Displays buy orders, bid sizes, accumulated volumes, and allows for
order cancellation.

\"My Orders\" includes offering symbol, order ID, account, limit price,
share count, minimum quantity, executed quantity/price, order date,
order type, value, and status.

**Contradicting or Unresolved Information**

**Terminology: \"Market\" vs. \"Listing\"**

Both documents note uncertainty about whether to use \"Market\" or
\"Listing\" for the primary exchange column in Closed and Upcoming
Offerings sections. This is explicitly marked as \"\[TBD\]\" in both
versions[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx).

**\"Last Sale\" vs. \"Offering Price\"**

In the Closed Offerings section, \"Last Sale\" is marked as \"(TBD.?
Same as Offering price?)\", indicating ambiguity about whether these
fields are synonymous or
distinct[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx).

**\"Issue Size\" Field**

In Upcoming Offerings, \"Issue Size\" is marked as \"N.A.? \[TBD\]\",
suggesting it may not be applicable or its definition is not
finalized[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx).

**Reports Functionality**

The ability to print Orders and Trades reports from the \"My Orders\"
widget is marked as \"\[TBD\]\", indicating this feature may not be
fully implemented or its status is
unresolved[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx).

**Chart Description**

The \"Chart\" feature in the \"My Orders\" section is labeled \"\[TDB.
Description\]\", suggesting its functionality or documentation is
incomplete[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx).

**Summary Table of Contradictions**

  --------------------------------------------------------------------------------------------------------------------
  **Area**       **Description of **Location in Docs**
                 Contradiction or 
                 TBD**            
  -------------- ---------------- ------------------------------------------------------------------------------------
  Market vs.     Unclear which    Closed/Upcoming
  Listing        term to use      Offerings[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx)

  Last Sale      Ambiguous if     Closed
                 same as Offering Offerings[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx)
                 Price            

  Issue Size     Marked as N.A.   Upcoming
                 or TBD           Offerings[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx)

  Reports        Status marked as My Orders
  Printing       TBD              Widget[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx)

  Chart          Description is   My Orders
  Description    missing (TDB)    Widget[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx)
  --------------------------------------------------------------------------------------------------------------------

**Conclusion:**\
The BidMan - Broker Dealer App is a flexible, widget-based trading and
order management platform for broker-dealers, offering comprehensive
access to market data, order entry, and account management. Several
sections of the documentation contain unresolved terminology, feature
definitions, or incomplete descriptions, but no direct contradictions
between the documents were identified---rather, both documents
consistently flag the same areas as needing clarification or
completion[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx).

Sources
[1](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-1.docx)
BidMan-User-Guide-v0-1.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/dc82727f-f275-48ce-bd11-14fc076c025f/BidMan-User-Guide-v0-1.docx>
[2](https://app.markdownword.barinbritva.com/BidMan-User-Guide-v0-2.docx)
BidMan-User-Guide-v0-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/e5243159-49ba-46a2-9c45-c585d64ce1a1/BidMan-User-Guide-v0-2.docx>

## BatMan -- Lead Manager App

**\[TBD: Tor S. use this a checklist and see if anything needs to be
included in section 4.3 -- BatMan. Then delete this section\]**

**Overview: BatMan -- Lead Manager App**

**BatMan** is the Lead Manager application within the ClearingBid Price
Discovery Platform (PDP) ecosystem. Its primary function is to enable
Lead Managers to monitor, manage, and execute IPO auctions and related
allocation processes. BatMan is positioned as the central tool for Lead
Managers to interact with the auction, analyze order book data, simulate
pricing scenarios, and finalize allocations.

**Key Features and Capabilities**

**Auction Monitoring and Management**

Real-time monitoring of bids, including Preferential Bids and
Competitive Bids.

Visualization of key metrics: number of preferential bids, number of
unique preferential bidders, round lot holders, large lot holders, and
compliance with exchange minimums.

Simulation functionality: Lead Managers can enter tentative offering
prices to see how allocations and compliance metrics would change\[1\].

**Allocation and Compliance**

Automated checks for compliance with exchange listing requirements
(e.g., minimum number of holders, round lots, market value).

Preferential Bids are tracked and allocated ahead of Competitive Bids if
necessary to meet listing requirements.

Allocation algorithms prioritize first-come, first-served for
preferential allocations, with further allocation rounds as needed to
fulfill requirements\[1\].

**Reporting and Export**

BatMan is intended to include a Reports widget, allowing Lead Managers
to view commissions and allocations for their offerings.

Planned functionality includes the ability to export allocation and
commission data as CSV for reconciliation and transfer of funds.

Some requirements suggest this reporting/export functionality may be
handled in the Backoffice app rather than BatMan itself, or that
implementation details are still under discussion.

**Integration and Workflow**

BatMan is integrated into the broader PDP environment, operating
alongside BidMan (for Broker-Dealers) and MarkOps (for Market
Operations).

Used during critical phases: auction monitoring, closing, allocation
submission, and error/bust management if needed.

**Workflow Example**

**Auction Monitoring:** Lead Manager uses BatMan to monitor order book
and compliance metrics in real time.

**Simulation:** Tentative prices can be entered to simulate allocations
and compliance outcomes.

**Allocation Submission:** Once auction closes and compliance is
verified, allocations are submitted via BatMan.

**Reporting/Export:** Allocation and commission reports are generated
(either in BatMan or Backoffice) for reconciliation and payment
processing.

**Contradictory or Ambiguous Information**

**CSV Export and Reporting Functionality**

**Planned in BatMan:** Some requirements specify that BatMan should have
a Reports widget and CSV export button for Lead Managers to export
allocation and commission data.

**Planned in Backoffice:** Other sections indicate that this
reporting/export functionality may instead be handled in the Backoffice
app, with the idea of adding it to BatMan considered \"old\" or not
proceeding.

**Current Status:** There is ambiguity as to whether BatMan will
ultimately provide direct CSV export/reporting, or if Lead Managers will
need to rely on Backoffice for these tasks.

**Allocation and Compliance Logic**

**Allocation Algorithm:** The documents consistently describe BatMan as
the tool for monitoring compliance and managing allocations, but some
details about the allocation algorithm (e.g., handling of \"implied
holders\" or future changes to rules) are marked as \"to be revised\" or
subject to further updates\[1\].

**Simulation Function:** The simulation capability is described as part
of BatMan, but the specifics of how it interacts with other system
components (e.g., BidMan and MarkOps) are not always fully detailed,
leading to some uncertainty about the exact workflow\[1\].

**Summary Table: BatMan - Lead Manager App**

  ------------------------------------------------------------------------------
  **Feature/Functionality**   **Description**   **Contradictions/Ambiguities**
  --------------------------- ----------------- --------------------------------
  Auction Monitoring          Real-time         None
                              metrics,          
                              compliance        
                              checks,           
                              simulation of     
                              allocations       

  Allocation Management       Preferential and  Details of allocation logic may
                              competitive bid   change\[1\]
                              allocation,       
                              compliance with   
                              listing           
                              requirements      

  Reporting/CSV Export        Planned Reports   Location (BatMan vs. Backoffice)
                              widget and CSV    unclear
                              export for        
                              allocations and   
                              commissions       

  Integration                 Works with BidMan BatMan as app vs. terminal
                              (BDs), SuperMan   module unclear
                              (Ops), BE         
                              Terminal widgets  

  Error/Bust Management       Supports reversal None
                              of closes and     
                              bust scenarios in 
                              conjunction with  
                              Market Ops        
  ------------------------------------------------------------------------------

**Conclusion**

BatMan is designed as the Lead Manager's central tool for monitoring IPO
auctions, ensuring compliance, managing allocations, and (potentially)
exporting reports. The core auction and allocation features are clearly
described and consistent across documents. However, there is some
ambiguity regarding where reporting/export functionality will reside and
whether BatMan is a standalone app or a module within the BE Terminal.
These points may require clarification as the platform evolves\[1\].

Sources \[1\] Holiday-schedule-010725-143132.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/68096c24-9212-4aa6-b3a0-fe6457226119/Holiday-schedule-010725-143132.pdf>
\[2\] Preferential-Bids-v-1-2-redlined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/62e01028-b0a8-48b7-afed-0355e3213d05/Preferential-Bids-v-1-2-redlined.docx>
\[3\] SD-PRD\_-PDP-Reporting-for-Offering-Close-300625-211034.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/8bb11d92-8e11-4996-9592-e104e7b21a66/SD-PRD_-PDP-Reporting-for-Offering-Close-300625-211034.pdf>
\[4\] BE-Terminal-Reference-v004-2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/bb17a992-9932-4b23-b225-8bb487beab97/BE-Terminal-Reference-v004-2.docx>
\[5\] BE-Terminal-Reference-v004-1.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/3daaa907-bef8-48cf-acbc-55eeeb5cb83d/BE-Terminal-Reference-v004-1.docx>
\[6\] Bust-Management-010725-142822.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/c81baec6-73d0-49b2-8f55-728a24511839/Bust-Management-010725-142822.pdf>
\[7\] ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>
\[8\]
SD-PRD\_-Offering-Dashboard-For-Issuer-and-Lead-Manager-010725-151542.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/8978878e-2a16-4319-8a5d-6e0fe6c6b162/SD-PRD_-Offering-Dashboard-For-Issuer-and-Lead-Manager-010725-151542.pdf>
\[9\] CB-PDP-Gen-2-specification-v2-0r01.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/f0a7c1b2-8ef0-428d-b64c-f20597bfe366/CB-PDP-Gen-2-specification-v2-0r01.docx>

## MarkOps -- PDP System and Market Operations App

**\[TBD: Tor S. to review and update this section\]**

**MarkOps -- PDP System and Market Operations App**

**Overview**

The **MarkOps - PDP System and Market Operations App** is designed to
facilitate and oversee the operational workflow for IPOs and new issue
securities. Its primary functions include system state management,
allocation verification, reporting, market surveillance, and ensuring
the integrity and efficiency of the closing and settlement processes.

**Key Features and Functions**

**1. System State Management**

**State Controls:** Market Operations can set the system to various
states (Halted, Pre-Open, Open, Close Pending, Closing, Clearing,
Closed, Canceled, Frozen) directly from the Backoffice app.

**Operational Impacts:** In the Halted state, order entry/cancellation
and offering closings are disabled. Transitions between states are
controlled to comply with regulatory requirements, such as ensuring a
minimum time between resuming from halt and auction close\[1\].

**2. Allocation Integrity Verification**

**Verification Workflow:** Before finalizing allocations, the Operations
Manager uses the app to:

Run fill simulations (submit button is initially greyed out).

Generate a snapshot and comprehensive report for allocation
verification.

Approve the offering after successful verification, enabling the Lead
Manager to submit the close.

**Checks Performed:** The system checks for errors such as:

Orders with no allocation despite qualifying price.

Orders with allocation below the offering price.

Allocation mismatches by more than one share from priority group fill.

Total allocated shares not matching the offering size\[2\]\[3\].

**Audit Trail:** Snapshots and reports are created for audit purposes,
supporting both manual and automated checks.

**3. Reporting and Reconciliation**

**CSV Export:** Market Operations can export comprehensive reports from
the Backoffice, including all relevant fields for reconciliation and
commission calculations.

**Commission Calculations:** Reports include selling concession
calculations, gross spread, and detailed trade breakdowns for each
broker-dealer.

**Threshold Controls:** Exports are limited by configurable row
thresholds to prevent system lag\[4\].

**4. Market Surveillance and Controls**

**Order Suspension and Rejection:** The app enables manual or rule-based
rejection/suspension of suspicious orders, with alerts for prolonged
suspensions.

**State Freezing:** The system can freeze order entry/cancellation
during critical periods (e.g., closing windows) to prevent
manipulation\[1\].

**Surveillance Alerts:** Integrated with market surveillance tools to
notify Market Operations of disruptive or manipulative activity\[5\].

**5. Closing and Settlement Operations**

**Multi-Stage Approval:** The closing process requires verification and
approval by both Market Operations and the Lead Manager before
allocations are finalized and files are sent to the clearing firm.

**Automated and Manual Controls:**Execution reports and end-of-day files
can be generated and sent based on configurable parameters for same-day
or next-day settlement.

**Countdown for Cancellations:** After closing, a countdown window
allows Market Operations to abort or proceed with sending FIX
cancellation messages, defaulting to 30 minutes\[2\]\[6\].

**6. Security and Access**

**Role-Based Permissions:** Access to critical functions (e.g., state
changes, allocation approval) is restricted to authorized Market
Operations personnel.

**Audit and Logging:** All actions are logged for compliance and review.

**Contradictory or Inconsistent Information**

**1. Location of Allocation Integrity Verification**

**Initial Approach:** Verification was to be performed outside the PDP
(e.g., using Excel), with plans to internalize the process in the PDP
after initial launches\[2\]\[3\].

**Current Approach:** Later documents and workflows indicate that
verification is now expected to be run within the PDP on the actual
generated allocations, eliminating errors from data
export/import\[2\]\[3\].

**Contradiction:** There is inconsistency about whether verification is
external (offline) or internal (within PDP), but the prevailing
direction is toward internalization.

**2. Visibility and Access to Reports**

**Backoffice vs. Lead Manager:** Some documents specify that detailed
allocation verification reports are only for Market Operations in the
Backoffice, not for Lead Managers in the BatMan app\[2\]. Others suggest
Lead Managers should have access to certain reports for their
offerings\[4\].

**Contradiction:** The scope of report visibility and export
functionality between roles is not fully harmonized.

**3. Order Suspension and Surveillance**

**Manual vs. Automated Actions:**Requirements for manual order
rejection/suspension coexist with proposals for automated rule-based
surveillance and rejection\[7\]\[8\]. The implementation status and
priority of these features are not always clear.

**Contradiction:** There is some ambiguity about how much is manual
versus automated, and which features are currently live.

**4. State Definitions and Transitions**

**Use of Pre-Open and Close Pending States:** Some documents question
the necessity of all defined states (Pre-Open, Open, Pre-Close, etc.),
while others lay out detailed workflows using all these states\[1\].

**Contradiction:** The rationale and operational use of each state may
not be uniformly agreed upon or documented.

**Summary Table: MarkOps PDP App Core Functions**

  -------------------------------------------------------
  **Functionality**   **Description**   **Contradictory
                                        Details?**
  ------------------- ----------------- -----------------
  System State        Controls          State definitions
  Management          system/market     not uniform
                      states, halts,    
                      and transitions   

  Allocation          Multi-stage       Internal vs.
  Verification        verification and  external process
                      approval of       
                      allocations       

  Reporting &         CSV exports,      Report access by
  Reconciliation      commission        role unclear
                      calculations,     
                      audit trails      

  Market Surveillance Manual and        Manual vs.
                      automated order   automated status
                      controls, alerts  

  Closing &           Multi-stage       None noted
  Settlement          approvals, EOD    
                      file generation,  
                      cancellation      
                      countdowns        

  Security & Access   Role-based        None noted
                      permissions,      
                      audit logging     
  -------------------------------------------------------

**References**

\[2\]\[1\]\[4\]\[6\]\[3\]\[8\] (see inline citations above for precise
sourcing)

Sources \[1\] SD-PRD\_-System-States-010725-150251.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/6efba2df-f560-4928-8718-88aa7a25f616/SD-PRD_-System-States-010725-150251.pdf>
\[2\] SD-PRD\_-Allocation-Integrity-Verification-300625-211230.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/434a11b8-fa30-4fe2-85c8-13e2062fb3d1/SD-PRD_-Allocation-Integrity-Verification-300625-211230.pdf>\[3\]
Holiday-schedule-010725-143132.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/68096c24-9212-4aa6-b3a0-fe6457226119/Holiday-schedule-010725-143132.pdf>
\[4\] SD-PRD\_-PDP-Reporting-for-Offering-Close-300625-211034.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/8bb11d92-8e11-4996-9592-e104e7b21a66/SD-PRD_-PDP-Reporting-for-Offering-Close-300625-211034.pdf>
\[5\] CB-system-high-level-description-v4-04.04.22.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/a869d44e-1ca0-40fb-a4a6-07b3b15461c3/CB-system-high-level-description-v4-04.04.22.docx>
\[6\] SD-PRD\_-PDP-Trade-Creation-Settlement-Updates-010725-150916.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/3e006b25-bdfb-4068-9603-78b401b4c13e/SD-PRD_-PDP-Trade-Creation-Settlement-Updates-010725-150916.pdf>
\[7\] SD-PRD-Future-PDP-enhancements-010725-151449.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/d0d0d9b9-1776-43d7-b5d7-f9c4b630a86c/SD-PRD-Future-PDP-enhancements-010725-151449.pdf>
\[8\] SD-PRD\_-Investor-Account-Fix-to-PDP-Mapping-010725-143448.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/b94ec060-0019-4d13-abb3-1f910f16bb8c/SD-PRD_-Investor-Account-Fix-to-PDP-Mapping-010725-143448.pdf>
\[9\] ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/0309247e-53ee-437d-a84a-5d35a8aa15c2/ClearingBid-PDP-Functional-Specification-v2-0-01-combined.docx>
\[10\] CB-PDP-Gen-2-specification-v2-0r01.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/f0a7c1b2-8ef0-428d-b64c-f20597bfe366/CB-PDP-Gen-2-specification-v2-0r01.docx>
\[11\] CB-Backend-Specification-v1-4-8-comments.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/571d1572-1928-4edd-a5d2-fd3b20d2ba8d/CB-Backend-Specification-v1-4-8-comments.docx>
\[12\] ClearingBid-New-Issue-Auction-Process-v-1-2-8.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/22cd9570-1ab0-4c3f-b33a-16a854e18019/ClearingBid-New-Issue-Auction-Process-v-1-2-8.docx>
\[13\]
ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/25825a46-ec4c-4544-93e0-5e46d77e1d5a/ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx>
\[14\] SD-PDP-Database-Data-010725-150114.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/e30db78e-bb47-469d-8d1a-3dac94510dec/SD-PDP-Database-Data-010725-150114.pdf>

## Market Surveillance

[ClearingBid IPO Price Discovery and Allocation v0 0
8.docx](https://clearingbid.sharepoint.com/:w:/s/ClearingBidGeneral/IQAsMJOe2OUwS7yZc5X7SB5RARq7qN3y1-6a-a2Mzza6lXc?e=fcKjTt)

**Market Surveillance: Overview and Functionality**

**Market Surveillance** in the ClearingBid ecosystem refers to the set
of tools, alerts, and operational processes designed to maintain the
integrity of the order book and detect manipulative or disruptive
trading behavior during IPOs and other offerings. Its main objectives
are to ensure fair price discovery, prevent market manipulation, and
enable timely intervention by market operations staff.

**Core Features and Alerts**

Market Surveillance is primarily described as a combination of automated
alerts and manual tools available to Market Operations Managers. Key
features include:

**1. Misleading Order Entry Alert**

Detects significant orders (or a series of significant orders) that are
cancelled shortly after the cancellation prevention period, potentially
indicating manipulative intent.

Default parameters:

Large order: 0.1 of total issue quantity

Very large order: 0.5 of total issue quantity

Significant number of orders: 5

Cancellation window: 60 minutes

**2. Crowd Swarm Alert**

Triggers when a large number of smaller orders or cancellations, within
a short period, significantly impact the clearing price.

The definition of \"large number\" and \"significant impact\" is
configurable\[1\].

**3. Significant Clearing Price Change Alert**

Identifies substantial clearing price changes across multiple rolling
time periods (e.g., 5, 10, 60, 240, 1440 minutes).

Lists the orders that most impacted the price change.

**4. Operational Tools**

Freeze new order entry (user or broker-dealer)

Mass cancel of orders (user or broker-dealer)

Suspend individual orders (orders remain in order book but are excluded
from clearing price calculation)

Suspended orders can be released or cancelled; if not resolved, system
alerts are triggered and closing cannot proceed\[1\]\[2\].

**5. Automated Remediation**

The system may automatically delay closing if alerts are triggered and
significant suspicious activity is detected\[1\].

**6. Notification and Audit**

Alerts are delivered via email to designated recipients.

## Network Monitoring and Connectivity Management

**\[TBD: Tor S. to review and update this section\]**

**Network Monitoring and Connectivity Management**

**Overview**

ClearingBid\'s platform incorporates robust **network monitoring and
connectivity management** to ensure secure, reliable, and efficient
operation for all participants---broker-dealers, lead managers, issuers,
and market operations. The system is designed to support real-time
trading, order management, and market surveillance, leveraging both
public and secure network channels.

**Connectivity Management**

**Physical Connectivity Options:**

**Encrypted VPN/Direct Internet:**Clients can connect via encrypted VPN
tunnels or secure internet connections.

**Dedicated Lines:** Option for dedicated leased lines for higher
reliability and lower latency.

**Colocation:** Clients may collocate infrastructure within the same
data centers as the network provider to minimize latency\[1\]\[2\].

**FIX Protocol Integration:**

The system uses the industry-standard FIX (Financial Information
eXchange) protocol (version 4.2) for electronic order and data
transmission.

Clients configure their FIX engines (e.g., QuickFIX, FIXEdge) to
communicate with the network gateway, setting parameters such as
SenderCompID, TargetCompID, IP, and port.

Onboarding includes supervised FIX certification, testing order types,
recovery scenarios, and simulated line failures\[1\]\[2\].

**Access Portals:**

**ClearingBid's web portal (CB=Web):** For general access to offering
information, educational materials, and real-time market data.

**Secure Portal:** For FIX transactions and BatMan/BidMan access,
restricted to lead managers, syndicate members, and broker-dealers,
using controlled data connections.

**Workstation and API Access:**

Broker-dealers, lead managers, and issuers can use browser-based
workstations or APIs for order management and offering oversight.

The system supports both manual and automated (FIX) order entry and
management.

**Network Monitoring**

**Market Surveillance Platform (MSP):**

Continuously monitors order flow and trading activity for potential
manipulation or disruptive behaviors.

Generates alerts for:

**Misleading Order Entry:**Detects significant orders entered and
quickly canceled after the cancellation prevention period.

**Crowd Swarm Alerts:**Identifies clusters of smaller orders or
cancellations that impact the indicative clearing price.

**Significant Price Change Alerts:** Monitors for large price swings
over configurable rolling time periods\[3\].

**Operational Controls:**

Tools to freeze new order entry, mass cancel orders, or suspend
individual orders for investigation.

Suspended orders are kept in the order book but do not participate in
price calculations until released or canceled.

Alerts can trigger automatic remediation or delays in closing windows if
suspicious activity is detected\[3\].

**Network Health and Recovery:**

Certification and integration tests include simulated line failures and
recovery scenarios to ensure resilience.

Heartbeat and test request messages are part of the FIX session to
monitor connection health\[1\]\[2\].

**Contradictory Information**

**No direct contradictions** were found in the reviewed documentation
regarding network monitoring and connectivity management. The documents
are consistent in describing:

The use of multiple connectivity options (VPN, direct internet,
dedicated lines, colocation).

The central role of the FIX protocol for secure, standardized
communications.

The presence of both public and secure network portals.

The implementation of comprehensive market surveillance and operational
controls.

**Minor differences** in terminology or emphasis (e.g., detailed
configuration steps or which monitoring tools are prioritized) reflect
evolving feature sets or document perspectives, but do not constitute
substantive contradictions.

**Summary Table**

  ---------------------------------------------------
  **Aspect**       **Description**   **References**
  ---------------- ----------------- ----------------
  Connectivity     VPN, direct       \[1\]\[2\]
  Options          internet,         
                   dedicated lines,  
                   colocation        

  Protocols        FIX 4.2 for order \[1\]\[2\]
                   and data          
                   transmission      

  Access           Public portal     
                   (info), secure    
                   portal            
                   (transactions,    
                   terminals)        

  Monitoring Tools Market            \[3\]
                   Surveillance      
                   Platform, alerts  
                   for manipulation, 
                   operational       
                   controls for      
                   order management  

  Health &         FIX session       \[1\]\[2\]
  Recovery         monitoring,       
                   certification     
                   testing,          
                   simulated         
                   failures          

  Contradictions   None found; minor ---
                   differences in    
                   detail only       
  ---------------------------------------------------

Sources \[1\] FIXmarketData4-2v1-4.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/f78db1de-58c3-4189-bcf0-11b67f1ad96b/FIXmarketData4-2v1-4.docx>\[2\]
Preferential-Bids-v-1-2-redlined.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/62e01028-b0a8-48b7-afed-0355e3213d05/Preferential-Bids-v-1-2-redlined.docx>
\[3\] SD-Preferential-Bids-010725-143823.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/8cc8f496-8010-409a-950f-9978f1fb74e1/SD-Preferential-Bids-010725-143823.pdf>
\[4\] CB-system-high-level-description-v4-04.04.22.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/a869d44e-1ca0-40fb-a4a6-07b3b15461c3/CB-system-high-level-description-v4-04.04.22.docx>
\[5\] SD-FIX-Message-Flow-010725-150358.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/b7cd962b-13de-4b2c-b592-ded864024d62/SD-FIX-Message-Flow-010725-150358.pdf>
\[6\]
ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/25825a46-ec4c-4544-93e0-5e46d77e1d5a/ClearingBid_Participant_System_Guidelines_v0_1_3_2-9-25_vs_0_1_2.docx>
\[7\] Bust-Management-010725-142822.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/c81baec6-73d0-49b2-8f55-728a24511839/Bust-Management-010725-142822.pdf>
\[8\] ClearingBid-IPO-Price-Discovery-and-Allocation-v0-0-8.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/4cb22cd3-2d63-4d83-bd30-d28c1f0d2bde/ClearingBid-IPO-Price-Discovery-and-Allocation-v0-0-8.docx>
\[9\] CB-Backend-Specification-v1-4-8-comments.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/571d1572-1928-4edd-a5d2-fd3b20d2ba8d/CB-Backend-Specification-v1-4-8-comments.docx>
\[10\] CB-PDP-Gen-2-specification-v2-0r01.docx
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/f0a7c1b2-8ef0-428d-b64c-f20597bfe366/CB-PDP-Gen-2-specification-v2-0r01.docx>
\[11\] SD-PRD-Market-Surveillance-010725-150731.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/a1564578-dd83-4b41-abb5-d24641bb63f3/SD-PRD-Market-Surveillance-010725-150731.pdf>
\[12\] SD-PRD-Future-PDP-enhancements-010725-151449.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/d0d0d9b9-1776-43d7-b5d7-f9c4b630a86c/SD-PRD-Future-PDP-enhancements-010725-151449.pdf>
\[13\] SD-PRD\_-PDP-Reporting-for-Offering-Close-300625-211034.pdf
<https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/collection_1c4dc5ea-7fe6-4b4e-ab69-232b6a416ee3/8bb11d92-8e11-4996-9592-e104e7b21a66/SD-PRD_-PDP-Reporting-for-Offering-Close-300625-211034.pdf>

# Notes, Remaining Questions, Issues and TBDs

- Gross Underwriting Spread - It is currently unclear if the **Gross
  Underwriting Spread** for bonds is **based on the Face Value or Bond
  Offering Price**?

- Must CBM IPO Trade Executions be reported to a **trade reporting
  facility**, e.g., FINRA-ADF, even if the security at the time of the
  trade is not an NMS security (because it is still not listed)? If so,
  a FINRA Transparency Services Uniform Executing Broker Agreement is
  required to enable CBM to report trade information to FINRA on behalf
  of another member, even if the parties have a QSR agreement in effect.
  Under the executing party trade reporting structure, members can
  continue to agree to allow another member to report and lock-in trades
  on their behalf. If not, include wording on how the Trade Executions
  should be reported, or that they do not need to be reported.

# Version history

+-----------+--------------+-------------------------------------------------------+
| **Date    | **Version#** | **Change Summary**                                    |
| created** |              |                                                       |
+-----------+--------------+-------------------------------------------------------+
| 6/8/25    | V2.0r01      | Document structure                                    |
+-----------+--------------+-------------------------------------------------------+
| 6/12/25   | V2.0r02      | Added content from Perplexity output. Some references |
|           |              | are duplicated.                                       |
+-----------+--------------+-------------------------------------------------------+
| 6/24/25   | V2.0r03      | New version to preserve original Perplexity output.   |
|           |              | Tracked changes turned on. Going forward, all tracked |
|           |              | changes should be accepted when a new named version   |
|           |              | is created.                                           |
|           |              |                                                       |
|           |              | Correction of Perplexity output plus updates from     |
|           |              | Source References documents.                          |
+-----------+--------------+-------------------------------------------------------+
| 6/29/25   | V2.0r04      | New headings based on Tor's approach to add           |
|           |              | Perplexity output on lowest sub-heading level. \[TS}: |
|           |              | continued to develop this with relevant information   |
|           |              | described on each level. I.e. not only on sub-level.  |
|           |              | Additional documents added to references, PRDs,       |
|           |              | guidelines, auction process, new definitions          |
+-----------+--------------+-------------------------------------------------------+
| 2/28/26   | V2.0r05      | Lots of updates during last 6 months by Roland        |
+-----------+--------------+-------------------------------------------------------+
| 3/1/26    | V2.0r06      | Current draft                                         |
+-----------+--------------+-------------------------------------------------------+

[^1]: The bid can be expressed in Yield (%) or per unit (\$) for bonds.
    For all other asset types, the bid is per unit (\$). Sometimes when
    the term price is used it should be understood as yield if it is in
    the context of a bond. For example, Price-Time Priority Price
    Calculation can be read as Yield-Time Priority Yield Calculation
    when it relates to a bond.

[^2]: Shares reserved for \"price takers\" may be excluded from the
    auctioned primary shares. As a consequence the number of primary
    shares in the auction may not match the number on the cover of the
    prospectus.

[^3]: For convenience "Indicative Clearing Price and Clearing Price" is
    shortened to "Clearing Price" in this section.

[^4]: Firm Accounts should not be confused with the investor Account.
    The investor Account identifies the beneficial account holder and is
    the identification used on BidMan and FIX entered orders when the
    investor Account identification is required.

[^5]: For context: This description is simplified. Banks do not actually
    have IPO Control Accounts. DTC monitors broker-dealers and banks
    differently. Banks are monitored via DTC's Institutional Delivery
    (ID) System. Both IPO Initial Distribution and IPO Secondary Market
    shares are delivered to the banks\' General Free Accounts, but they
    are tagged in the IPO database as IPO Initial Distribution and IPO
    Secondary Market shares, respectively. Hence, the process is similar
    and, therefore, this document has simplified the description. The
    full description of the DTCC IPO Tracking System, in *The DTC
    Service Guide "IPO.pdf",* comprises 143 pages. The detailed
    transaction and workflow for delivery to the IPIO Control Accounts
    and via the DTC Institutional Delivery System must be worked out
    with DTCC experts.

[^6]: Note that max potential allocation can be larger than MinPrefSec
    and is then an indication of the amount of Preferential Bid
    oversubscription at the set minimum number of required bidders.

[^7]: When the qualifying 'in the money' quantity for a Round Lot or
    Large Lot bidder is made up of multiple bids the time priority of
    the last bid that qualified the bidder as a Round Lot or Large Lot
    bidder is used to determine the "Qualification Time Priority" of the
    bidder. The Qualification Time Priority is used to determine the
    priority of additional 'max potential allocation' for the 'first'
    MinRndLots or MinLrgPct\*MinRndLots number of bidders respectively.
