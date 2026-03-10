**ClearingBid, Inc.**

FIX 4.2 Trade API Implementation Manual

![Blue candlestick chart](media/image1.jpeg){width="6.49in"
height="5.22in"}

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
the FIX team. The certification test will replicate all messages planned
to be sent in production. The test will also include unsolicited
cancels, simulated trade busts, and recovery testing (simulated line
failure).

**ClearingBid Implementation of Standard FIX tags**

FIX Trade API

V1.5

Last Update: 12/20/2024 

 

# Table of Contents {#table-of-contents .TOC-Heading}

[Messages Throttling [5](#_Toc176967436)](#_Toc176967436)

[Message Structure [5](#_Toc176967437)](#_Toc176967437)

[Authorization [5](#_Toc176967438)](#_Toc176967438)

[Header [6](#_Toc176967439)](#_Toc176967439)

[Body [7](#_Toc176967440)](#_Toc176967440)

[Trailer [7](#_Toc176967441)](#_Toc176967441)

[Session Messages [8](#_Toc176967442)](#_Toc176967442)

[Heartbeat [8](#_Toc176967443)](#_Toc176967443)

[Test Request [8](#_Toc176967444)](#_Toc176967444)

[Resend Request [8](#_Toc176967445)](#_Toc176967445)

[Reject [8](#_Toc176967446)](#_Toc176967446)

[Sequence Reset [9](#_Toc176967447)](#_Toc176967447)

[Logout [9](#_Toc176967448)](#_Toc176967448)

[Logon [9](#_Toc176967449)](#_Toc176967449)

[Messages [11](#_Toc176967450)](#_Toc176967450)

[New Order Single [11](#_Toc176967451)](#_Toc176967451)

[Execution Report [12](#_Toc176967452)](#_Toc176967452)

[Order Cancel Replace Request [14](#_Toc176967453)](#_Toc176967453)

[Order Cancel Request [15](#_Toc176967454)](#_Toc176967454)

[Order Cancel Reject [15](#_Toc176967455)](#_Toc176967455)

[Trading Session Status Request [16](#_Toc176967456)](#_Toc176967456)

[Trading Session Status [16](#_Toc176967457)](#_Toc176967457)

**Summary **

This document describes ClearingBid's Trade FIX API protocol. 

The API is based on ***FIX.4.**2* version specification. 

 

[]{#_Toc176967436 .anchor}**Messages Throttling **

Client is not allowed to exceed N incoming messages per second. Messages
exceeding the limit will be rejected with a Business Message Reject. 

Current max value is 1000 messages per 1 seconds. The result of
exceeding the limit is:

- Session will be logged out

- All client\'s orders will be canceled

- Comp Id will be deactivated. 

 

[]{#_Toc176967437 .anchor}**Message Structure **

All messages have the same structure: Header, Body, Trailer. 

 

[]{#_Toc176967438 .anchor}**Authorization **

Sender Comp ID - Sender Comp ID from user details 

The **Logon** message expects an extra tag with a password

The logon message has to be the first message.  

  -----------------------------------------
  554    Password     String in MD5 
  ------ ------------ ---------------------

  -----------------------------------------

[]{#_Toc176967439 .anchor}**Header **

** **Header structure is the same for all messages. 

+-----------+----------------+------------+----------------+--------------------------+
| **Tag **  | **Field **     | **Type **  | **Required **  | **Valid values /         |
|           |                |            |                | Comments**               |
+===========+================+============+================+==========================+
| 8         | Begin String   | String     | Y              | FIX.4.2                  |
+-----------+----------------+------------+----------------+--------------------------+
| 9         | Body Length    | int        | Y              |                          |
+-----------+----------------+------------+----------------+--------------------------+
| 35        | Msg Type       | String     | Y              |                          |
+-----------+----------------+------------+----------------+--------------------------+
| 49        | Sender Comp    | String     | Y              | Identifies the           |
|           | ID             |            |                | originator of the        |
|           |                |            |                | message. Eg: Institution |
|           |                |            |                | or third party.          |
+-----------+----------------+------------+----------------+--------------------------+
| 56        | Target Comp    | String     | Y              |                          |
|           | ID             |            |                |                          |
+-----------+----------------+------------+----------------+--------------------------+
| 115       | On Behalf Of   | String     | Y if HUB       | Assigned value used to   |
|           | Comp ID        |            |                | identify a firm          |
|           |                |            |                | originating message if   |
|           |                |            |                | the message was          |
|           |                |            |                | delivered by a third     |
|           |                |            |                | party. The third party   |
|           |                |            |                | would then be displayed  |
|           |                |            |                | in the Sender Comp ID.   |
|           |                |            |                | The information in this  |
|           |                |            |                | field on the order is    |
|           |                |            |                | returned in the Deliver  |
|           |                |            |                | To Comp ID (tag 128)     |
|           |                |            |                | field on each execution  |
|           |                |            |                | related to that order.   |
+-----------+----------------+------------+----------------+--------------------------+
| 128       | Deliver To     | String     | Y if HUB       | Assigned value used to   |
|           | Comp ID        |            |                | identify the firm        |
|           |                |            |                | targeted to receive the  |
|           |                |            |                | message if the message   |
|           |                |            |                | is delivered by a third  |
|           |                |            |                | party. The third-party   |
|           |                |            |                | firm identifier would be |
|           |                |            |                | delivered in the Target  |
|           |                |            |                | Comp ID (56) field. This |
|           |                |            |                | tag is populated on      |
|           |                |            |                | execution messages with  |
|           |                |            |                | the information that was |
|           |                |            |                | contained in the On      |
|           |                |            |                | Behalf Of Comp ID tag    |
|           |                |            |                | (115) of the order       |
|           |                |            |                | message.                 |
+-----------+----------------+------------+----------------+--------------------------+
| 34        | Msg Seq Num    | int        | Y              |                          |
+-----------+----------------+------------+----------------+--------------------------+
| 52        | Sending Time   | UTC Time   | Y              |                          |
|           |                | stamp      |                |                          |
+-----------+----------------+------------+----------------+--------------------------+
| 50        | Sender Sub ID  | String     | N              | Assigned value used to   |
|           |                |            |                | identify specific        |
|           |                |            |                | message originator       |
|           |                |            |                | (desk, OMS trader,       |
|           |                |            |                | etc.). The value in this |
|           |                |            |                | tag is returned on an    |
|           |                |            |                | execution in the Target  |
|           |                |            |                | Sub ID field (tag 57).   |
+-----------+----------------+------------+----------------+--------------------------+
| 57        | Target Sub ID  | String     | N              | Assigned value used to   |
|           |                |            |                | identify specific        |
|           |                |            |                | individual or unit       |
|           |                |            |                | intended to receive      |
|           |                |            |                | message. This tag is     |
|           |                |            |                | populated on execution   |
|           |                |            |                | messages with the value  |
|           |                |            |                | contained in the Sender  |
|           |                |            |                | Sub ID tag (50) of the   |
|           |                |            |                | order message.           |
+-----------+----------------+------------+----------------+--------------------------+
| 116       | On Behalf of   | String     | N              | Assigned value used to   |
|           | Sub ID         |            |                | identify specific        |
|           |                |            |                | message originator (i.e. |
|           |                |            |                | trader) if the message   |
|           |                |            |                | was delivered by a third |
|           |                |            |                | party. The information   |
|           |                |            |                | in this field on the     |
|           |                |            |                | order is returned in the |
|           |                |            |                | Deliver To Sub ID (tag   |
|           |                |            |                | 129) field on each       |
|           |                |            |                | execution related to     |
|           |                |            |                | that order.              |
+-----------+----------------+------------+----------------+--------------------------+
| 129       | Deliver To Sub | String     | N              | Assigned value used to   |
|           | ID             |            |                | identify specific        |
|           |                |            |                | message recipient (i.e.  |
|           |                |            |                | trader) if the message   |
|           |                |            |                | is delivered by a third  |
|           |                |            |                | party.                   |
|           |                |            |                |                          |
|           |                |            |                | This tag is populated on |
|           |                |            |                | execution messages with  |
|           |                |            |                | the information that was |
|           |                |            |                | contained in the On      |
|           |                |            |                | Behalf Of Sub ID tag     |
|           |                |            |                | (116) of the order       |
+-----------+----------------+------------+----------------+--------------------------+

** **

[]{#_Toc176967440 .anchor}**Body ** 

Body structure depends on message type. 

[]{#_Toc176967441 .anchor}**Trailer **

Trailer structure is the same for all messages. Tag 10 is the only
required.

  -------------------------------------------------------
  **Tag **    **Field **    **Type **    **Required ** 
  ----------- ------------- ------------ ----------------
  10          Check Sum     String       Y 

  -------------------------------------------------------

** **

**\**

[]{#_Toc176967442 .anchor}

**Session Messages **

[]{#_Toc176967443 .anchor}**Heartbeat **

Msg Type = 0 

  ----------------------------------------------------------
  **Tag**    **Field**               **Type**    **Req.** 
  ---------- ----------------------- ----------- -----------
  112        Test Req ID             String      N 

  ----------------------------------------------------------

  

[]{#_Toc176967444 .anchor}**Test Request **

Msg Type = 1 

  ----------------------------------------------------------
  **Tag**    **Field**               **Type**    **Req.** 
  ---------- ----------------------- ----------- -----------
  112        Test Req ID             String      Y 

  ----------------------------------------------------------

 

[]{#_Toc176967445 .anchor}**Resend Request **

Msg Type = 2 

  ----------------------------------------------------------
  **Tag**    **Field**               **Type**    **Req.** 
  ---------- ----------------------- ----------- -----------
  7          Begin Seq No            Int         Y 

  16         End Seq No              Int         Y 
  ----------------------------------------------------------

 

[]{#_Toc176967446 .anchor}**Reject **

Msg Type = 3 

  ----------------------------------------------------------
  **Tag**    **Field**               **Type**    **Req.** 
  ---------- ----------------------- ----------- -----------
  45         Ref Seq Num             Int         Y 

  371        Ref Tag ID              Int         N 

  372        Ref Msg Type            String      N 

  373        Session Reject Reason   Int         N 

  58         Text                    String      N 
  ----------------------------------------------------------

 

[]{#_Toc176967447 .anchor}**Sequence Reset **

Msg Type = 4 

  ----------------------------------------------------------
  **Tag**    **Field**               **Type**    **Req.** 
  ---------- ----------------------- ----------- -----------
  123        Gap Fill Flag           Boolean     N 

  36         New Seq No              Int         Y 
  ----------------------------------------------------------

[]{#_Toc176967448 .anchor}**Logout **

Msg Type = 5 

  ----------------------------------------------------------
  ** Tag**    **Field**              **Type**    **Req.** 
  ----------- ---------------------- ----------- -----------
  58          Text                   String      N 

  ----------------------------------------------------------

  

[]{#_Toc176967449 .anchor}**Logon **

Msg Type = A 

  -----------------------------------------------------------
  ** Tag**    **Field**               **Type**    **Req.** 
  ----------- ----------------------- ----------- -----------
  98          Encrypt Method          Int         Y 

  108         Heartbeat Int           Int         Y 

  95          Raw Data Length         Length      N 

  96          Raw Data                Data        N 

  141         Reset Seq Num Flag      Boolean     N 

  383         Max Message Size        Int         N 

  554         Password                String In   Y 
                                      MD5         
  -----------------------------------------------------------

**\**

[]{#_Toc176967450 .anchor}**Messages  **

[]{#_Toc176967451 .anchor}**New Order Single **

 Msg Type = D 

+-----------+----------------+------------+----------------+-------------------------+
| **Tag **  | **Field **     | **Type **  | **Required **  | **Valid Values **       |
+===========+================+============+================+=========================+
| 1         | Account        | String     | N              | User account name       |
+-----------+----------------+------------+----------------+-------------------------+
| 11        | Cl Ord ID      | String     | Y              | Unique for specific     |
|           |                |            |                | Comp ID                 |
+-----------+----------------+------------+----------------+-------------------------+
| 21        | HandlInst      | char       | Y              | - 2 = Automated public  |
+-----------+----------------+------------+----------------+-------------------------+
| 110       | Min Qty        | Qty        | N              |                         |
+-----------+----------------+------------+----------------+-------------------------+
| 55        | Symbol         | String     | Y              |                         |
+-----------+----------------+------------+----------------+-------------------------+
| 54        | Side           | char       | Y              | - 1 = Buy               |
+-----------+----------------+------------+----------------+-------------------------+
| 60        | Transact Time  | UTC        | Y              |                         |
|           |                | Timestamp  |                |                         |
+-----------+----------------+------------+----------------+-------------------------+
| 38        | Order Qty      | Qty        | Y              |                         |
+-----------+----------------+------------+----------------+-------------------------+
| 40        | Ord Type       | char       | Y              | - 2 = Limit             |
|           |                |            |                |                         |
|           |                |            |                | <!-- -->                |
|           |                |            |                |                         |
|           |                |            |                | - 5 = Market on Close   |
|           |                |            |                |                         |
|           |                |            |                | <!-- -->                |
|           |                |            |                |                         |
|           |                |            |                | - P = Pegged            |
|           |                |            |                |                         |
|           |                |            |                | Notes:                  |
|           |                |            |                |                         |
|           |                |            |                | Ord Type=5, Ord Type=P  |
|           |                |            |                | applicable only for     |
|           |                |            |                | ETFs.                   |
|           |                |            |                |                         |
|           |                |            |                | Pegged order requires   |
|           |                |            |                | Price (44) field.       |
|           |                |            |                |                         |
|           |                |            |                | Market On Close with    |
|           |                |            |                | Price (44) field will   |
|           |                |            |                | be rejected.            |
+-----------+----------------+------------+----------------+-------------------------+
| 44        | Price          | Price      | Y              | Note: for BOND the      |
|           |                |            |                | price must be expressed |
|           |                |            |                | as 100 minus the        |
|           |                |            |                | bidder's minimum        |
|           |                |            |                | accepted Yield.         |
|           |                |            |                |                         |
|           |                |            |                | For example, an         |
|           |                |            |                | investor that is        |
|           |                |            |                | willing to buy a bond   |
|           |                |            |                | if the Yield is 3.537%  |
|           |                |            |                | or higher, enters a     |
|           |                |            |                | limit buy order with a  |
|           |                |            |                | dollar price of         |
|           |                |            |                | 96.463.                 |
+-----------+----------------+------------+----------------+-------------------------+
| 59        | Time In Force  | char       | N              | - 1 = Good till Cancel  |
|           |                |            |                |   (GTC)                 |
+-----------+----------------+------------+----------------+-------------------------+
| 18        | ExecInst       | char       | N              | - p = Preferential bid  |
|           |                |            |                |   (PREFERENTIAL)        |
+-----------+----------------+------------+----------------+-------------------------+
| 58        | Text           | String     | N              |                         |
+-----------+----------------+------------+----------------+-------------------------+
| 20101     | Order          | char       | N              | - B = INSTITUTIONAL     |
|           | Capacity       |            |                |                         |
|           |                |            |                | <!-- -->                |
|           |                |            |                |                         |
|           |                |            |                | - C = RETAIL            |
+-----------+----------------+------------+----------------+-------------------------+

** **[]{#_Toc176967452 .anchor}**Execution Report **

Msg Type = 8 

+-----------+-------------+------------+----------------+--------------------+
| **Tag **  | **Field **  | **Type **  | **Required **  | **Valid Values **  |
+===========+=============+============+================+====================+
| 1         | Account     | String     | N              | User account name  |
+-----------+-------------+------------+----------------+--------------------+
| 63        | SettlmntTyp | int        | N              | Indicates order    |
|           |             |            |                | settlement period. |
|           |             |            |                | Absence of this    |
|           |             |            |                | field is           |
|           |             |            |                | interpreted as     |
|           |             |            |                | Regular            |
|           |             |            |                |                    |
|           |             |            |                | - 0 = Regular      |
|           |             |            |                |                    |
|           |             |            |                | - 1 = Cash         |
+-----------+-------------+------------+----------------+--------------------+
| 37        | Order ID    | String     | Y              | Unique             |
+-----------+-------------+------------+----------------+--------------------+
| 11        | Cl Ord ID   | String     | Y              |                    |
+-----------+-------------+------------+----------------+--------------------+
| 41        | Orig Cl Ord | String     | N              |                    |
|           | ID          |            |                |                    |
+-----------+-------------+------------+----------------+--------------------+
| 17        | Exec ID     | String     | Y              | Unique             |
+-----------+-------------+------------+----------------+--------------------+
| 18        | Exec Inst   | String     | N              | - p = Preferential |
|           |             |            |                |   bid              |
|           |             |            |                |   (PREFERENTIAL)   |
+-----------+-------------+------------+----------------+--------------------+
| 15        | Currency    | String     | Y              | - USD              |
+-----------+-------------+------------+----------------+--------------------+
| 20        | Exec Trans  | char       | Y              | - 0 = New          |
|           | Type        |            |                |                    |
+-----------+-------------+------------+----------------+--------------------+
| 19        | Exec Ref    | String     | N              |                    |
|           | ID          |            |                |                    |
+-----------+-------------+------------+----------------+--------------------+
|           | Exec Type   | char       | Y              | - 0 = New          |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - 1 = Partial      |
|           |             |            |                |   fill             |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - 2 = Fill         |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - 4 = Canceled     |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - 8 = Rejected     |
+-----------+-------------+------------+----------------+--------------------+
| 39        | Ord Status  | char       | Y              | - 0 = New          |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - 1 = Partial      |
|           |             |            |                |   fill             |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - 2 = Fill         |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - 4 = Canceled     |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - 8 = Rejected     |
+-----------+-------------+------------+----------------+--------------------+
| 103       | Ord Rej     | int        | N              | - 1 = Unknown      |
|           | Reason      |            |                |   symbol           |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - 2 = Exchange     |
|           |             |            |                |   closed           |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - 6 = Duplicate    |
|           |             |            |                |   Order            |
+-----------+-------------+------------+----------------+--------------------+
| 55        | Symbol      | String     | Y              |                    |
+-----------+-------------+------------+----------------+--------------------+
| 54        | Side        | char       | Y              |                    |
+-----------+-------------+------------+----------------+--------------------+
| 38        | Order Qty   | Qty        | Y              |                    |
+-----------+-------------+------------+----------------+--------------------+
| 40        | Ord Type    | char       | Y              | - 2 = Limit        |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - 5 = Market on    |
|           |             |            |                |   Close            |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - P = Pegged       |
+-----------+-------------+------------+----------------+--------------------+
| 44        | Price       | Price      | N              |                    |
+-----------+-------------+------------+----------------+--------------------+
| 59        | Time In     | char       | Y              | - 1 = Good till    |
|           | Force       |            |                |   Cancel (GTC)     |
+-----------+-------------+------------+----------------+--------------------+
| 32        | Last        | Qty        | Y              |                    |
|           | Shares      |            |                |                    |
+-----------+-------------+------------+----------------+--------------------+
| 31        | Last Px     | Price      | Y              | Expressed as 100   |
|           |             |            |                | minus the bidder's |
|           |             |            |                | minimum accepted   |
|           |             |            |                | Yield.             |
|           |             |            |                |                    |
|           |             |            |                | For example, an    |
|           |             |            |                | investor that is   |
|           |             |            |                | willing to buy a   |
|           |             |            |                | bond if the Yield  |
|           |             |            |                | is 3.537% or       |
|           |             |            |                | higher, enters a   |
|           |             |            |                | limit buy order    |
|           |             |            |                | with a dollar      |
|           |             |            |                | price of 96.463.   |
+-----------+-------------+------------+----------------+--------------------+
| 151       | Leaves Qty  | Qty        | Y              |                    |
+-----------+-------------+------------+----------------+--------------------+
| 14        | Cum Qty     | Qty        | Y              |                    |
+-----------+-------------+------------+----------------+--------------------+
| 6         | Avg Px      | Price      | Y              | Expressed as 100   |
|           |             |            |                | minus the bidder's |
|           |             |            |                | minimum accepted   |
|           |             |            |                | Yield.  For        |
|           |             |            |                | example, an        |
|           |             |            |                | investor that is   |
|           |             |            |                | willing to buy a   |
|           |             |            |                | bond if the Yield  |
|           |             |            |                | is 3.537% or       |
|           |             |            |                | higher, enters a   |
|           |             |            |                | limit buy order    |
|           |             |            |                | with a dollar      |
|           |             |            |                | price of 96.463.   |
+-----------+-------------+------------+----------------+--------------------+
| 60        | Transact    | UTC Time   | Y              |                    |
|           | Time        | stamp      |                |                    |
+-----------+-------------+------------+----------------+--------------------+
| 58        | Text        | String     | N              |                    |
+-----------+-------------+------------+----------------+--------------------+
| 20101     | Order       | char       | N              | - B =              |
|           | Capacity    |            |                |   INSTITUTIONAL    |
|           |             |            |                |                    |
|           |             |            |                | <!-- -->           |
|           |             |            |                |                    |
|           |             |            |                | - C = RETAIL       |
|           |             |            |                |                    |
|           |             |            |                | ClearingBid may    |
|           |             |            |                | request that this  |
|           |             |            |                | categorization is  |
|           |             |            |                | made.              |
+-----------+-------------+------------+----------------+--------------------+

** **

[]{#_Toc176967453 .anchor}**Order Cancel Replace Request **

Msg Type = G 

  -----------------------------------------------------------
  **Tag **    **Field **      **Type **      **Required ** 
  ----------- --------------- -------------- ----------------
  11          Cl Ord ID       String         Y 

  1           Account         String         N 

  41          Orig Cl Ord ID  String         Y 

  21          Handl Inst      char           Y 

  55          Symbol          String         Y 

  54          Side            char           Y 

  60          Transact Time   UTC Time       Y 
                              stamp          

  38          Order Qty       Qty            Y 

  40          Ord Type        char           Y 

  59          Time In Force   char           N 

  58          Text            String         N 
  -----------------------------------------------------------

** **

[]{#_Toc176967454 .anchor}**Order Cancel Request **

Msg Type = F 

  -----------------------------------------------------------
  **Tag **    **Field **      **Type **      **Required ** 
  ----------- --------------- -------------- ----------------
  1           Account         String         N 

  11          Cl Ord ID       String         Y 

  41          Orig Cl Ord ID  String         Y 

  55          Symbol          String         Y 

  54          Side            char           Y 

  60          Transact Time   UTC Time       Y 
                              stamp          

  38          Order Qty       Qty            Y 

  40          Ord Type        char           Y 

  59          Time In Force   char           N 

  58          Text            String         N 
  -----------------------------------------------------------

** **

[]{#_Toc176967455 .anchor}**Order Cancel Reject **

Msg Type = 9 

+-----------+--------------+------------+----------------+---------------------------+
| **Tag **  | **Field **   | **Type **  | **Required **  | **Valid values **         |
+===========+==============+============+================+===========================+
| 37        | Order ID     | String     | Y              |                           |
+-----------+--------------+------------+----------------+---------------------------+
| 11        | Cl Ord ID    | String     | Y              |                           |
+-----------+--------------+------------+----------------+---------------------------+
| 41        | Orig Cl Ord  | String     | Y              |                           |
|           | ID           |            |                |                           |
+-----------+--------------+------------+----------------+---------------------------+
| 39        | Ord Status   | char       | Y              | - 8 = Rejected            |
+-----------+--------------+------------+----------------+---------------------------+
| 60        | Transact     | UTC        | N              |                           |
|           | Time         | Timestamp  |                |                           |
+-----------+--------------+------------+----------------+---------------------------+
| 434       | Cxl Rej      | char       | Y              | - 1 = Order Cancel        |
|           | Response To  |            |                |   Request                 |
|           |              |            |                |                           |
|           |              |            |                | <!-- -->                  |
|           |              |            |                |                           |
|           |              |            |                | - 2 -- Order Cancel       |
|           |              |            |                |   Replace Request         |
+-----------+--------------+------------+----------------+---------------------------+
| 102       | Cxl Rej      | int        | N              | - 0 - Too late            |
|           | Reason       |            |                |                           |
|           |              |            |                | <!-- -->                  |
|           |              |            |                |                           |
|           |              |            |                | - 1 - Unknown order       |
|           |              |            |                |                           |
|           |              |            |                | <!-- -->                  |
|           |              |            |                |                           |
|           |              |            |                | - 2 - Broker option       |
+-----------+--------------+------------+----------------+---------------------------+
| 58        | Text         | String     | N              |                           |
+-----------+--------------+------------+----------------+---------------------------+

 

[]{#_Toc176967456 .anchor}**Trading Session Status Request **

Msg Type = g 

+-----------+----------------------+-----------+-----------+----------------------------+
| ** Tag**  | **Field**            | **Type**  | **Req.**  | **Valid values**           |
+===========+======================+===========+===========+============================+
| 336       | Trading Session ID   | String    | Y         |                            |
+-----------+----------------------+-----------+-----------+----------------------------+
| 335       | Trad Ses Req ID      | String    | N         |                            |
+-----------+----------------------+-----------+-----------+----------------------------+
| 263       | Subscription Request | char      | Y         | 0 = SNAPSHOT               |
|           | Type                 |           |           |                            |
|           |                      |           |           | 1 = SNAPSHOT_PLUS_UPDATES  |
|           |                      |           |           |                            |
|           |                      |           |           | 2 = DISABLE_PREVIOUS       |
+-----------+----------------------+-----------+-----------+----------------------------+

[]{#_Toc176967457 .anchor}**Trading Session Status **

Msg Type = h 

+-----------+----------------------+-----------+-----------+---------------------------+
| ** Tag**  | **Field**            | **Type**  | **Req.**  | **Valid values**          |
+===========+======================+===========+===========+===========================+
| 336       | Trading Session ID   | String    | Y         |                           |
+-----------+----------------------+-----------+-----------+---------------------------+
| 335       | Trad Ses Req ID      | String    | N         |                           |
+-----------+----------------------+-----------+-----------+---------------------------+
| 340       | Trad Ses Status      | char      | Y         | 1 = HALTED                |
|           |                      |           |           |                           |
|           |                      |           |           | 2 = OPEN                  |
|           |                      |           |           |                           |
|           |                      |           |           | 3 = CLOSED                |
|           |                      |           |           |                           |
|           |                      |           |           | 4 = PRE_OPEN              |
|           |                      |           |           |                           |
|           |                      |           |           | 5 = PRE_CLOSE             |
+-----------+----------------------+-----------+-----------+---------------------------+

** **

**ClearingBid Markets**

**Six Landmark Square Suite 403**

**Stamford, Ct. 06901**
