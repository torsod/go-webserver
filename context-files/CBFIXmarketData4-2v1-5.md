**ClearingBid, Inc.**

![Blue candlestick chart](media/image1.jpeg){width="6.49in"
height="5.22in"}FIX 4.2 Market Data API Implementation Manual

**\**

**FIX** (Financial Information eXchange) is an open, industry-standard
protocol for the electronic transmission of orders and related data. The
protocol is defined and maintained by an industry association, the FIX
trading Community.

This manual describes how ClearingBid (\"CBID\") implements the
industry-standard Financial Information Exchange (FIX™) Protocol,
specifically for CBID clients and vendors who connect to CBID using
\"FIX\". CBID currently supports FIX version 4.2.

For a complete guide to the FIX standard, please refer to the official
FIX homepage at the following URL: <https://www.fixprotocol.org>.

This document assumes the default definition of the FIX standard unless
otherwise noted. Only the message types and fields (tags and values)
supported by the CBID FIX engine are included in this manual.
Unsupported fields are not guaranteed to be included in FIX messages
sent by CBID. If CBID receives an unsupported field, it is not validated
for processing purposes. FIX messages with unknown tags are not rejected
if they include all necessary supported tags. Messages of types that are
not supported are rejected.

Use of FIX requires an integration process lead by ClearingBid's FIX
Engineering team.

**Network and Connectivity:**

**ClearingBid**

ClearingBid's order matching system is hosted with AWS. Internet
connections to ClearingBid are secured by Stunnel. ClearingBid offers
universal TLS/SSL tunneling via the open-source Stunnel tool.

For clients with leased lines ClearingBid can provide connectivity in
data centers that offer AWS Direct Connect service.

**Virtu (ITG Net)**

Clients connecting with Virtu (ITG Net) can route orders to ClearingBid
via their existing connection. ClearingBid has certified connectivity
with Virtu.

**FIX For OMS Vendors**

Client's OMS provider may contact ClearingBid for integration.

**Onboarding and Integration**

Client may schedule a supervised FIX certification test with a member of
the FIX team. The certification test will replicate all messages that
are planned to be sent in production. The test will also include
unsolicited cancels, simulated trade busts, and recovery testing
(simulated line failure).

**ClearingBid Implementation of Standard Market Data FIX tags**

FIX Market Data API

V1.5

Last Update: 12/20/2024

#  {#section .TOC-Heading}

# Table of Contents {#table-of-contents .TOC-Heading}

[Message Structure [5](#message-structure)](#message-structure)

[Header [5](#header)](#header)

[Body [6](#body)](#body)

[Trailer [6](#trailer)](#trailer)

[Session Messages [7](#session-messages)](#session-messages)

[Heartbeat [7](#heartbeat)](#heartbeat)

[Test Request [7](#test-request)](#test-request)

[Resend Request [7](#resend-request)](#resend-request)

[Reject [7](#reject)](#reject)

[Sequence Reset [8](#sequence-reset)](#sequence-reset)

[Logout [8](#logout)](#logout)

[Logon [9](#logon)](#logon)

[Application Messages [9](#application-messages)](#application-messages)

[Business Message Reject
[9](#business-message-reject)](#business-message-reject)

[Market Data Request [10](#market-data-request)](#market-data-request)

[Market Data Request Reject
[10](#market-data-request-reject)](#market-data-request-reject)

[Market Data Incremental Refresh
[11](#market-data-incremental-refresh)](#market-data-incremental-refresh)

[Security Definition Request
[12](#security-definition-request)](#security-definition-request)

[Security Definition [12](#security-definition)](#security-definition)

[Trading Session Status Request
[14](#trading-session-status-request)](#trading-session-status-request)

[Trading Session Status
[14](#trading-session-status)](#trading-session-status)

#  

**Summary**

This document describes ClearingBid's Market Data FIX API protocol.

The API is based on ***FIX.4.2*** version specification.

## Message Structure

All messages have the same structure: Header, Body, Trailer.

## Header

Header structure is the same for all messages.

  ---------------------------------------------------------------------------------
  **Tag**   **Field**   **Type**         **Required**
  --------- ----------- ---------------- ------------------------------------------
  8         Begin       String           Y
            String                       

  9         Body Length int              Y

  35        Msg Type    String           Y

  49        Sender Comp String           Y
            ID                           

  56        Target Comp String           Y
            ID                           

  34        Msg Seq Num int              Y

  50        Sender Sub  String           N
            ID                           

  57        Target Sub  String           N
            ID                           

  43        Poss Dup    Boolean          Y for resend
            Flag                         

  52        Sending     UTC Timestamp    Y
            Time                         

  122       Orig        UTC Timestamp    Y for resend
            Sending                      
            Time                         
  ---------------------------------------------------------------------------------

## Body

Body structure depends on message type.

## Trailer

Trailer structure is the same for all messages. Only Tag 10 is required.

  -----------------------------------------------------------------------------------
  **Tag**   **Field**   **Type**   **Required**
  --------- ----------- ---------- --------------------------------------------------
  10        Check Sum   String     Y

  -----------------------------------------------------------------------------------

## Session Messages

## Heartbeat

MsgType = 0

  ----------------------------------------------------------------------------------
  **Tag**   **Field**   **Type**   **Required**
  --------- ----------- ---------- -------------------------------------------------
  112       Test Req ID String     N

  ----------------------------------------------------------------------------------

## Test Request

MsgType = 1

  ----------------------------------------------------------------------------------
  **Tag**   **Field**   **Type**   **Required**
  --------- ----------- ---------- -------------------------------------------------
  112       Test Req ID String     Y

  ----------------------------------------------------------------------------------

## Resend Request

MsgType = 2

  -----------------------------------------------------------------------------------
  **Tag**   **Field**    **Type**   **Required**
  --------- ------------ ---------- -------------------------------------------------
  7         BeginSeqNo   int        Y

  16        EndSeqNo     int        Y
  -----------------------------------------------------------------------------------

## Reject

MsgType = 3

  --------------------------------------------------------------------------------------
  **Tag**   **Field**             **Type**   **Required**
  --------- --------------------- ---------- -------------------------------------------
  45        RefSeqNum             int        Y

  371       RefTagID              int        N

  372       RefMsgType            String     N

  373       SessionRejectReason   int        N

  58        Text                  String     N
  --------------------------------------------------------------------------------------

## Sequence Reset

MsgType = 4

  -----------------------------------------------------------------------------------
  **Tag**   **Field**     **Type**   **Required**
  --------- ------------- ---------- ------------------------------------------------
  123       GapFillFlag   Boolean    N

  36        NewSeqNo      int        Y
  -----------------------------------------------------------------------------------

## Logout

MsgType = 5

  ---------------------------------------------------------------------------------------
  **Tag**   **Field**   **Type**   **Required**
  --------- ----------- ---------- ------------------------------------------------------
  58        Text        String     N

  ---------------------------------------------------------------------------------------

## Logon

MsgType = A

  -----------------------------------------------------------------------------
  **Tag**   **Field**         **Type**   **Required**
  --------- ----------------- ---------- --------------------------------------
  98        EncryptMethod     int        Y

  108       HeartBeatInt      int        Y

  95        RawDataLength     Length     N

  96        RawData           data       N

  141       ResetSeqNumFlag   Boolean    N

  383       MaxMessageSize    int        N

  554       Password          String in  Y
                              MD5        
  -----------------------------------------------------------------------------

## Application Messages

## Business Message Reject

+---------+------------+----------+--------------+-----------------------------------------+
| **Tag** | **Field**  | **Type** | **Required** | **Notes**                               |
+=========+============+==========+==============+=========================================+
| 45      | RefSeqNum  | int      | N            |                                         |
+---------+------------+----------+--------------+-----------------------------------------+
| 372     | RefMsgType | String   | Y            |                                         |
+---------+------------+----------+--------------+-----------------------------------------+
| 380     | Business   | int      | Y            | - 4 = Application not available\        |
|         | Reject     |          |              |   In case incoming messages limit       |
|         | Reason     |          |              |   exceeding                             |
+---------+------------+----------+--------------+-----------------------------------------+
| 58      | Text       | String   | N            |                                         |
+---------+------------+----------+--------------+-----------------------------------------+

## Market Data Request

MsgType = V

Subscriptions will be done for both full order book and offering
changes.

Note: only symbol in each group needed

+---------+-------------------------+----------+--------------+----------------------------+
| **Tag** | **Field**               | **Type** | **Required** | **Notes**                  |
+=========+=========================+==========+==============+============================+
| 262     | MDReqId                 | String   | Y            | Unique per day and per     |
|         |                         |          |              | client                     |
+---------+-------------------------+----------+--------------+----------------------------+
| 263     | SubscriptionRequestType | char     | Y            | - 1 = Subscribe            |
|         |                         |          |              |                            |
|         |                         |          |              | - 2 = Unsubscribe          |
+---------+-------------------------+----------+--------------+----------------------------+
| 265     | MDUpdateType            | int      | Y if 263 = 1 | Always 1 = Incremental     |
|         |                         |          |              | Refresh                    |
+---------+-------------------------+----------+--------------+----------------------------+
| 267     | NoMDEntryTypes          | int      | Y            |                            |
+---------+-------------------------+----------+--------------+----------------------------+
| \> 269  | MDEntryType             | char     | Y            | Always 0 = BID             |
+---------+-------------------------+----------+--------------+----------------------------+
| 146     | NoRelatedSym            | int      | Y            | Number of symbols          |
|         |                         |          |              | requested                  |
+---------+-------------------------+----------+--------------+----------------------------+
| \> 55   | Symbol                  | String   | Y            |                            |
+---------+-------------------------+----------+--------------+----------------------------+

## Market Data Request Reject

MsgType = Y

+---------+----------------+----------+--------------+------------------------------------+
| **Tag** | **Field**      | **Type** | **Required** | **Notes**                          |
+=========+================+==========+==============+====================================+
| 262     | MDReqId        | String   | Y            | Unique per day and per client      |
+---------+----------------+----------+--------------+------------------------------------+
| 281     | MDReqRejReason | char     | N            | - 0 = UNKNOWN_SYMBOL               |
+---------+----------------+----------+--------------+------------------------------------+
| 58      | Text           | String   | N            | List of rejected symbols           |
+---------+----------------+----------+--------------+------------------------------------+

## Market Data Incremental Refresh

MsgType = X

+---------+----------------+--------------+--------------+----------------------------+
| **Tag** | **Field**      | **Type**     | **Required** | **Notes**                  |
+=========+================+==============+==============+============================+
| 262     | MDReqId        | String       | N            |                            |
+---------+----------------+--------------+--------------+----------------------------+
| 268     | NoMDEntries    | int          | Y            |                            |
+---------+----------------+--------------+--------------+----------------------------+
| \> 279  | MDUpdateAction | char         | Y            | - 0 = New                  |
|         |                |              |              |                            |
|         |                |              |              | - 1 = Change               |
|         |                |              |              |                            |
|         |                |              |              | - 2 = Delete               |
+---------+----------------+--------------+--------------+----------------------------+
| \> 269  | MDEntryType    | char         | Y            | BID(0)=Book Entry Price,   |
|         |                |              |              | OFFER(1)=Potential         |
|         |                |              |              | Clearing Price             |
|         |                |              |              |                            |
|         |                |              |              | note: in BOND case all     |
|         |                |              |              | prices are in \$ (100      |
|         |                |              |              | minus)                     |
+---------+----------------+--------------+--------------+----------------------------+
| \> 278  | MDEntryID      | String       | N            |                            |
+---------+----------------+--------------+--------------+----------------------------+
| \> 280  | MDEnttyRefID   | String       | N            |                            |
+---------+----------------+--------------+--------------+----------------------------+
| \> 55   | Symbol         | String       | Y            |                            |
+---------+----------------+--------------+--------------+----------------------------+
| \> 270  | MDEntryPx      | Price        | N            | Book Entry Price if        |
|         |                |              |              | MDEntryType=BID(0),        |
|         |                |              |              |                            |
|         |                |              |              | Clearing Price if          |
|         |                |              |              | MDEntryType=OFFER(1)       |
|         |                |              |              |                            |
|         |                |              |              | note: in BOND case all     |
|         |                |              |              | prices are in \$ (100      |
|         |                |              |              | minus)                     |
+---------+----------------+--------------+--------------+----------------------------+
| \> 271  | MDEntrySize    | Qty          | N            | qty on price level if      |
|         |                |              |              | MDEntryType=BID(0)         |
+---------+----------------+--------------+--------------+----------------------------+
| \> 273  | MDEntryTime    | UTCTimeStamp | N            |                            |
+---------+----------------+--------------+--------------+----------------------------+

## Security Definition Request

MsgType = c

+---------+---------------------+----------+--------------+--------------------------------------+
| **Tag** | **Field**           | **Type** | **Required** | **Notes**                            |
+=========+=====================+==========+==============+======================================+
| 320     | SecurityReqID       | String   | Y            |                                      |
+---------+---------------------+----------+--------------+--------------------------------------+
| 321     | SecurityRequestType | int      | Y            | - 3 = Request List Securities        |
+---------+---------------------+----------+--------------+--------------------------------------+
| 55      | Symbol              | String   | N            | - If specific security is requested, |
|         |                     |          |              |   Symbol should be specified         |
|         |                     |          |              |                                      |
|         |                     |          |              | - If all securities are requested,   |
|         |                     |          |              |   Symbol should be blank             |
+---------+---------------------+----------+--------------+--------------------------------------+

## Security Definition

MsgType = d

+---------+--------------------+--------------+--------------+---------------------------+
| **Tag** | **Field**          | **Type**     | **Required** | **Notes**                 |
+=========+====================+==============+==============+===========================+
| 320     | SecurityReqID      | String       | Y            |                           |
+---------+--------------------+--------------+--------------+---------------------------+
| 322     | SecurityResponseID | String       | Y            |                           |
+---------+--------------------+--------------+--------------+---------------------------+
| 393     | TotalNumSecurities | int          | Y            |                           |
+---------+--------------------+--------------+--------------+---------------------------+
| 55      | Symbol             | String       | Y            |                           |
+---------+--------------------+--------------+--------------+---------------------------+
| 167     | SecurityType       | String       | Y            | - CS = Common Stock       |
|         |                    |              |              |                           |
|         |                    |              |              | - PS = Preferred Stock    |
|         |                    |              |              |                           |
|         |                    |              |              | - CPS = Convertible       |
|         |                    |              |              |   Preferred Stock         |
|         |                    |              |              |                           |
|         |                    |              |              | - CORP = Corporate Bond   |
|         |                    |              |              |                           |
|         |                    |              |              | - GO = General Obligation |
|         |                    |              |              |   Bonds                   |
|         |                    |              |              |                           |
|         |                    |              |              | - TEB = Tax Exempt Bond   |
|         |                    |              |              |                           |
|         |                    |              |              | - MF = Mutual Fund        |
|         |                    |              |              |                           |
|         |                    |              |              | - ETF = Exchange Traded   |
|         |                    |              |              |   Fund                    |
+---------+--------------------+--------------+--------------+---------------------------+
| 107     | SecurityDesc       | String       | Y            | Name on UI                |
+---------+--------------------+--------------+--------------+---------------------------+
| 20200   | Offering Size      | Qty          | N            | Offering size             |
+---------+--------------------+--------------+--------------+---------------------------+
| 20201   | Bid Period Start   | UTCTimeStamp | N            | Bidding period start      |
|         |                    |              |              | datetime                  |
+---------+--------------------+--------------+--------------+---------------------------+
| 20202   | BidPeriodEnd       | UTCTimeStamp | N            | Bidding period end        |
|         |                    |              |              | datetime                  |
+---------+--------------------+--------------+--------------+---------------------------+
| 106     | Issuer             | String       | N            | Company issuer            |
+---------+--------------------+--------------+--------------+---------------------------+
| 20203   | Lead Agent         | String       | N            | Company lead agent        |
+---------+--------------------+--------------+--------------+---------------------------+
| 20204   | Launch Price       | Price        | N            | only for Closed           |
+---------+--------------------+--------------+--------------+---------------------------+
| 20205   | Offered Size       | Qty          | N            | only for Closed           |
+---------+--------------------+--------------+--------------+---------------------------+
| 20206   | Status             | String       | Y            | NEW, UPCOMING, OPEN,      |
|         |                    |              |              | FREEZE, HALT,             |
|         |                    |              |              | CLOSE_PENDING, CLOSING,   |
|         |                    |              |              | CLEARING, CLOSED,         |
|         |                    |              |              | CANCELLED                 |
+---------+--------------------+--------------+--------------+---------------------------+
| 20207   | Allows             | Boolean      | Y            | Y = Allows Preferential   |
|         | Preferential       |              |              | Bids                      |
|         |                    |              |              |                           |
|         |                    |              |              | N = Disallows             |
+---------+--------------------+--------------+--------------+---------------------------+
| 20208   | Round Lot          | Qty          | Y            | Round Lot value           |
+---------+--------------------+--------------+--------------+---------------------------+
| 20209   | Max Lot            | Qty          | Y            | Max Lot value             |
+---------+--------------------+--------------+--------------+---------------------------+
| 20210   | RequiresClientType | Boolean      | Y            | Y = Requires set client   |
|         |                    |              |              | type                      |
|         |                    |              |              | (retail/institutional) in |
|         |                    |              |              | the order                 |
|         |                    |              |              |                           |
|         |                    |              |              | N = Not requires          |
+---------+--------------------+--------------+--------------+---------------------------+

## Trading Session Status Request

MsgType = g

+---------+-------------------------+----------+----------+-----------------------------+
| **Tag** | **Field**               | **Type** | **Req.** | **Valid values**            |
+=========+=========================+==========+==========+=============================+
| 336     | TradingSessionID        | String   | Y        |                             |
+---------+-------------------------+----------+----------+-----------------------------+
| 335     | TradSesReqID            | String   | N        |                             |
+---------+-------------------------+----------+----------+-----------------------------+
| 263     | SubscriptionRequestType | char     | Y        | 0 = SNAPSHOT                |
|         |                         |          |          |                             |
|         |                         |          |          | 1 = SNAPSHOT_PLUS_UPDATES   |
|         |                         |          |          |                             |
|         |                         |          |          | 2 = DISABLE_PREVIOUS        |
+---------+-------------------------+----------+----------+-----------------------------+

## Trading Session Status

MsgType = h

+---------+---------------------+----------+----------+-----------------------------+
| **Tag** | **Field**           | **Type** | **Req.** | **Valid values**            |
+=========+=====================+==========+==========+=============================+
| 336     | TradingSessionID    | String   | Y        |                             |
+---------+---------------------+----------+----------+-----------------------------+
| 335     | TradSesReqID        | String   | N        |                             |
+---------+---------------------+----------+----------+-----------------------------+
| 340     | TradSesStatus       | char     | Y        | 1 = HALTED                  |
|         |                     |          |          |                             |
|         |                     |          |          | 2 = OPEN                    |
|         |                     |          |          |                             |
|         |                     |          |          | 3 = CLOSED                  |
|         |                     |          |          |                             |
|         |                     |          |          | 4 = PRE_OPEN                |
|         |                     |          |          |                             |
|         |                     |          |          | 5 = PRE_CLOSE               |
+---------+---------------------+----------+----------+-----------------------------+

**ClearingBid Markets**

**Six Landmark Square Suite 403**

**Stamford, Ct. 06901**
