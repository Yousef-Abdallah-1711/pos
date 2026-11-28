# 35 — Financial Scenario Test Specification (NOT EXECUTED)

Source of expected values: `oracle/scenarios.py` → `oracle/scenarios_expected.json` (independent integer-piaster double-entry oracle; imports nothing from the application). Results so far = **oracle self-consistency only (24/24 PASS)**; nothing has been executed against the app.

**Assumption tags:** PROMPT = from the audit prompt / existing S13 control; PROVISIONAL = unapproved policy chosen only to compute a number (see [32](32-accounting-policy-decision-register.md)). If a policy is decided differently, regenerate the expectation. The existing S13 control is in `oracle/deterministic_control.py` (reviewed in [35a](35a-oracle-independent-review.md)).

**Capability-dependent scenarios:** where the workflow depends on a capability classified missing/partial in [33](33-capability-completeness-matrix.csv), the first execution step is a *capability-verification test* (can the UI/API do it at all?). If not, record NOT IMPLEMENTED and execute only the parts that exist (e.g. manual journal), without pretending the workflow is native.

Common UI route: Landing 13111 login → Dashboard 13112 (pages per [06](06-complete-page-inventory.md): POS, Orders, Accounting/Journal Entries, Suppliers/Invoices, Inventory, HR/Salaries, Reports). Fresh tenant per scenario group; reset by golden dump.


## SCN-01 — Opening balances

- Policy/capability links: OPEN-01
- Capability check first: none (core flow)
- Assumptions:
  - [PROVISIONAL] OPEN-01: opening balances posted as ONE journal offset to Opening Balance Equity (evidence-backed)
  - [PROVISIONAL] Inventory opening value 10,000 taken at book cost; AP opening 8,000
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
- Expected journals (EGP; Dr / Cr):
  - Opening balances: Dr Cash 50,000.00; Dr Inventory 10,000.00; Cr Accounts Payable 8,000.00; Cr Opening Balance Equity 52,000.00
- Expected ending balances: Cash 50,000.00; Inventory 10,000.00; Accounts Payable 8,000.00; Opening Balance Equity 52,000.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 0.00; net profit 0.00; assets 60,000.00; liabilities 8,000.00; equity incl. current earnings 52,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-02 — Capital contribution

- Policy/capability links: N4
- Capability check first: N4 (capital (manual journal))
- Assumptions:
  - [PROVISIONAL] Capital contribution of 100,000 into Bank is financing cash flow, equity credit to Owner Capital
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
- Expected journals (EGP; Dr / Cr):
  - Capital contribution: Dr Bank 100,000.00; Cr Owner Capital 100,000.00
- Expected ending balances: Bank 100,000.00; Owner Capital 100,000.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 0.00; net profit 0.00; assets 100,000.00; liabilities 0.00; equity incl. current earnings 100,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-03 — Cash purchase of inventory

- Policy/capability links: INV-01, TAX-01
- Capability check first: none (core flow)
- Assumptions:
  - [PROVISIONAL] INV purchase capitalised to Inventory (not expensed); input VAT recoverable and held in Input VAT asset
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - Capital contribution: Dr Cash 10,000.00; Cr Owner Capital 10,000.00
  - Cash purchase: Dr Inventory 5,000.00; Dr Input VAT 700.00; Cr Cash 5,700.00
- Expected ending balances: Cash 4,300.00; Inventory 5,000.00; Input VAT 700.00; Owner Capital 10,000.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 0.00; net profit 0.00; assets 10,000.00; liabilities 0.00; equity incl. current earnings 10,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-04 — Credit purchase and part supplier payment

- Policy/capability links: INV-01, TAX-01
- Capability check first: none (core flow)
- Assumptions:
  - [PROVISIONAL] Purchase capitalised to Inventory; AP carries gross incl. VAT; supplier paid by Bank
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
- Expected journals (EGP; Dr / Cr):
  - Capital contribution: Dr Bank 100,000.00; Cr Owner Capital 100,000.00
  - Credit purchase: Dr Inventory 20,000.00; Dr Input VAT 2,800.00; Cr Accounts Payable 22,800.00
  - Supplier payment: Dr Accounts Payable 10,000.00; Cr Bank 10,000.00
- Expected ending balances: Bank 90,000.00; Inventory 20,000.00; Input VAT 2,800.00; Accounts Payable 12,800.00; Owner Capital 100,000.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 0.00; net profit 0.00; assets 112,800.00; liabilities 12,800.00; equity incl. current earnings 100,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-05 — Purchase return to supplier

- Policy/capability links: N16, N7
- Capability check first: N16 (credit note), N7 (supplier→inventory)
- Assumptions:
  - [PROVISIONAL] Return of goods reverses Inventory, Input VAT and AP at original cost (no restocking fee)
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
- Expected journals (EGP; Dr / Cr):
  - Credit purchase: Dr Inventory 20,000.00; Dr Input VAT 2,800.00; Cr Accounts Payable 22,800.00
  - Purchase return to supplier: Dr Accounts Payable 4,560.00; Cr Inventory 4,000.00; Cr Input VAT 560.00
- Expected ending balances: Inventory 16,000.00; Input VAT 2,240.00; Accounts Payable 18,240.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 0.00; net profit 0.00; assets 18,240.00; liabilities 18,240.00; equity incl. current earnings 0.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-06 — Card sale

- Policy/capability links: REV, N21
- Capability check first: N21 (order payment records/gateway)
- Assumptions:
  - [PROVISIONAL] Card receipt settles directly into Bank at sale time, no processor fee / chargeback / settlement lag
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROMPT] No payment-processor fees, no discounts unless the scenario states otherwise
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - Capital contribution: Dr Cash 5,000.00; Cr Owner Capital 5,000.00
  - Stock: Dr Inventory 400.00; Cr Cash 400.00
  - Card sale: Cr Revenue 1,000.00; Cr Output VAT 140.00; Dr Bank 1,140.00
  - Card sale COGS: Dr COGS 400.00; Cr Inventory 400.00
- Expected ending balances: Cash 4,600.00; Bank 1,140.00; Output VAT 140.00; Owner Capital 5,000.00; Revenue 1,000.00; COGS 400.00
- Statements: net revenue 1,000.00; COGS 400.00; gross profit 600.00; opex 0.00; net profit 600.00; assets 5,740.00; liabilities 140.00; equity incl. current earnings 5,600.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-07 — Wallet sale with provider payout

- Policy/capability links: N21
- Capability check first: N21 (order payment records/gateway)
- Assumptions:
  - [PROVISIONAL] Wallet funds held in Wallet Clearing asset until payout; payout to Bank is operating cash; no fee
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROMPT] No payment-processor fees, no discounts unless the scenario states otherwise
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - Wallet sale: Cr Revenue 1,000.00; Cr Output VAT 140.00; Dr Wallet Clearing 1,140.00
  - Wallet payout: Dr Bank 1,140.00; Cr Wallet Clearing 1,140.00
- Expected ending balances: Bank 1,140.00; Output VAT 140.00; Revenue 1,000.00
- Statements: net revenue 1,000.00; COGS 0.00; gross profit 1,000.00; opex 0.00; net profit 1,000.00; assets 1,140.00; liabilities 140.00; equity incl. current earnings 1,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-08 — Split cash + card sale

- Policy/capability links: ROUND-01
- Capability check first: none (core flow)
- Assumptions:
  - [PROVISIONAL] Split tender: 1,000 cash + remainder card; whole-order VAT computed once on net, not per tender
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROMPT] No payment-processor fees, no discounts unless the scenario states otherwise
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - Split sale: Cr Revenue 2,000.00; Cr Output VAT 280.00; Dr Cash 1,000.00; Dr Bank 1,280.00
- Expected ending balances: Cash 1,000.00; Bank 1,280.00; Output VAT 280.00; Revenue 2,000.00
- Statements: net revenue 2,000.00; COGS 0.00; gross profit 2,000.00; opex 0.00; net profit 2,000.00; assets 2,280.00; liabilities 280.00; equity incl. current earnings 2,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-09 — Online prepaid order (card via gateway)

- Policy/capability links: N21, N22
- Capability check first: N21 (order payment records/gateway), N22 (order AR link)
- Assumptions:
  - [PROVISIONAL] Online capture recognises revenue at capture; delivery fee 0 (matches current app forced-0 fee, not a policy)
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROMPT] No payment-processor fees, no discounts unless the scenario states otherwise
- Expected journals (EGP; Dr / Cr):
  - Online prepaid order: Cr Revenue 2,000.00; Cr Output VAT 280.00; Dr Bank 2,280.00
- Expected ending balances: Bank 2,280.00; Output VAT 280.00; Revenue 2,000.00
- Statements: net revenue 2,000.00; COGS 0.00; gross profit 2,000.00; opex 0.00; net profit 2,000.00; assets 2,280.00; liabilities 280.00; equity incl. current earnings 2,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-10 — Online order cash-on-delivery with delivery fee

- Policy/capability links: N20, N21
- Capability check first: N20 (delivery/COD), N21 (order payment records/gateway)
- Assumptions:
  - [PROVISIONAL] Delivery fee is separate revenue and VAT-able at 14%
  - [PROVISIONAL] COD cash held by driver in Courier Cash Receivable until settled; revenue at delivery
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - COD delivered: Cr Revenue 3,000.00; Cr Output VAT 420.00; Cr Delivery Revenue 200.00; Cr Output VAT 28.00; Dr Courier Cash Receivable 3,648.00
  - Courier cash settlement: Dr Cash 3,648.00; Cr Courier Cash Receivable 3,648.00
- Expected ending balances: Cash 3,648.00; Output VAT 448.00; Revenue 3,000.00; Delivery Revenue 200.00
- Statements: net revenue 3,200.00; COGS 0.00; gross profit 3,200.00; opex 0.00; net profit 3,200.00; assets 3,648.00; liabilities 448.00; equity incl. current earnings 3,200.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-11 — Customer invoice and partial receipt

- Policy/capability links: N22
- Capability check first: N22 (order AR link)
- Assumptions:
  - [PROVISIONAL] Revenue recognised at invoice date (accrual); receipts clear AR; 5,400 remains in AR
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
- Expected journals (EGP; Dr / Cr):
  - Customer invoice: Dr Accounts Receivable 11,400.00; Cr Revenue 10,000.00; Cr Output VAT 1,400.00
  - Customer receipt (bank): Dr Bank 6,000.00; Cr Accounts Receivable 6,000.00
- Expected ending balances: Bank 6,000.00; Accounts Receivable 5,400.00; Output VAT 1,400.00; Revenue 10,000.00
- Statements: net revenue 10,000.00; COGS 0.00; gross profit 10,000.00; opex 0.00; net profit 10,000.00; assets 11,400.00; liabilities 1,400.00; equity incl. current earnings 10,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-12 — Order-level discount

- Policy/capability links: N17, TAX-01
- Capability check first: N17 (discount)
- Assumptions:
  - [PROVISIONAL] Discount ordering: VAT computed AFTER discount; discount shown in contra-revenue Sales Discounts (gross revenue 1,000)
  - [PROVISIONAL] NOTE conflicting implementation: customer-invoice trigger taxes gross before discount (04 section 3)
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - Discounted sale: Cr Revenue 1,000.00; Cr Output VAT 126.00; Dr Sales Discounts 100.00; Dr Cash 1,026.00
- Expected ending balances: Cash 1,026.00; Output VAT 126.00; Revenue 1,000.00; Sales Discounts 100.00
- Statements: net revenue 900.00; COGS 0.00; gross profit 900.00; opex 0.00; net profit 900.00; assets 1,026.00; liabilities 126.00; equity incl. current earnings 900.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-13 — Partial refund of a paid cash order

- Policy/capability links: N2, REFUND-01
- Capability check first: N2 (refund)
- Assumptions:
  - [PROVISIONAL] Refund is a linked reversal posted at refund date to Sales Returns (contra-revenue) and Output VAT
  - [PROVISIONAL] Returned goods are resalable and restocked at original cost (no write-off)
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - Stock: Dr Inventory 400.00; Cr Owner Capital 400.00
  - Sale 2 items: Cr Revenue 1,000.00; Cr Output VAT 140.00; Dr Cash 1,140.00
  - Sale 2 items COGS: Dr COGS 400.00; Cr Inventory 400.00
  - Refund 1 of 2 items: Dr Sales Returns 500.00; Dr Output VAT 70.00; Cr Cash 570.00
  - Restock refunded item: Dr Inventory 200.00; Cr COGS 200.00
- Expected ending balances: Cash 570.00; Inventory 200.00; Output VAT 70.00; Owner Capital 400.00; Revenue 1,000.00; Sales Returns 500.00; COGS 200.00
- Statements: net revenue 500.00; COGS 200.00; gross profit 300.00; opex 0.00; net profit 300.00; assets 770.00; liabilities 70.00; equity incl. current earnings 700.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-14 — Cancellation: unpaid order and paid-before-delivery order

- Policy/capability links: N3, REFUND-01
- Capability check first: none (core flow)
- Assumptions:
  - [PROVISIONAL] Unpaid cancellation has NO journal (matches 11-posting-matrix current behaviour, but it is a policy default)
  - [PROVISIONAL] Paid-order cancellation = full refund with stock return
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - Stock: Dr Inventory 400.00; Cr Owner Capital 400.00
  - Paid order: Cr Revenue 1,000.00; Cr Output VAT 140.00; Dr Cash 1,140.00
  - Paid order COGS: Dr COGS 400.00; Cr Inventory 400.00
  - Cancel paid order - refund: Dr Sales Returns 1,000.00; Dr Output VAT 140.00; Cr Cash 1,140.00
  - Cancel paid order - restock: Dr Inventory 400.00; Cr COGS 400.00
- Expected ending balances: Inventory 400.00; Owner Capital 400.00; Revenue 1,000.00; Sales Returns 1,000.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 0.00; net profit 0.00; assets 400.00; liabilities 0.00; equity incl. current earnings 400.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-15 — Inventory consumption by sale (weighted average)

- Policy/capability links: INV-01, COGS
- Capability check first: none (core flow)
- Assumptions:
  - [PROVISIONAL] INV-01: weighted-average cost, perpetual; COGS at sale = qty * (inventory value / qty on hand), half-up
  - [PROVISIONAL] Purchases: 10 units at EGP 100 then 10 at EGP 120, both cash, VAT 14% recoverable
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - Capital contribution: Dr Cash 10,000.00; Cr Owner Capital 10,000.00
  - Cash purchase: Dr Inventory 1,000.00; Dr Input VAT 140.00; Cr Cash 1,140.00
  - Cash purchase: Dr Inventory 1,200.00; Dr Input VAT 168.00; Cr Cash 1,368.00
  - Sale of 5 units @200: Cr Revenue 1,000.00; Cr Output VAT 140.00; Dr Cash 1,140.00
  - Sale of 5 units @200 COGS: Dr COGS 550.00; Cr Inventory 550.00
- Expected ending balances: Cash 8,632.00; Inventory 1,650.00; Input VAT 308.00; Output VAT 140.00; Owner Capital 10,000.00; Revenue 1,000.00; COGS 550.00
- Statements: net revenue 1,000.00; COGS 550.00; gross profit 450.00; opex 0.00; net profit 450.00; assets 10,590.00; liabilities 140.00; equity incl. current earnings 10,450.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-16 — Inventory adjustment (count shortage)

- Policy/capability links: N13, N14, INV-01
- Capability check first: N13 (waste/adjust)
- Assumptions:
  - [PROVISIONAL] Physical count short by 2 units valued at weighted-average 110 each -> Inventory Shrinkage expense
  - [PROVISIONAL] No VAT adjustment on shrinkage (accountant question)
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
- Expected journals (EGP; Dr / Cr):
  - Stock 15 units @110: Dr Inventory 1,650.00; Cr Owner Capital 1,650.00
  - Count shortage 2 units: Dr Inventory Shrinkage 220.00; Cr Inventory 220.00
- Expected ending balances: Inventory 1,430.00; Owner Capital 1,650.00; Inventory Shrinkage 220.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 220.00; net profit -220.00; assets 1,430.00; liabilities 0.00; equity incl. current earnings 1,430.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-17 — Waste / spoilage

- Policy/capability links: N13, INV-01
- Capability check first: N13 (waste/adjust)
- Assumptions:
  - [PROVISIONAL] 3 spoiled units at weighted-average 110 expensed to Waste Expense (separate from shrinkage)
  - [PROVISIONAL] Input VAT previously claimed is NOT reversed (accountant question)
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
- Expected journals (EGP; Dr / Cr):
  - Stock 15 units @110: Dr Inventory 1,650.00; Cr Owner Capital 1,650.00
  - Waste 3 units: Dr Waste Expense 330.00; Cr Inventory 330.00
- Expected ending balances: Inventory 1,320.00; Owner Capital 1,650.00; Waste Expense 330.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 330.00; net profit -330.00; assets 1,320.00; liabilities 0.00; equity incl. current earnings 1,320.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-18 — Payroll accrual, settlement and withholding remittance

- Policy/capability links: PAY-01, N11, N12
- Capability check first: none (core flow)
- Assumptions:
  - [PROVISIONAL] PAY-01: gross pay accrued to Payroll Expense when earned; employee deduction 800 (ILLUSTRATIVE, not a statutory rate) to Withholding Payable
  - [PROVISIONAL] Net pay 9,200 settled by Bank; deduction remitted later; employer contributions and tips excluded
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
- Expected journals (EGP; Dr / Cr):
  - Capital contribution: Dr Bank 50,000.00; Cr Owner Capital 50,000.00
  - Payroll accrual: Dr Payroll Expense 10,000.00; Cr Payroll Payable 9,200.00; Cr Withholding Payable 800.00
  - Net pay settlement: Dr Payroll Payable 9,200.00; Cr Bank 9,200.00
  - Withholding remittance: Dr Withholding Payable 800.00; Cr Bank 800.00
- Expected ending balances: Bank 40,000.00; Owner Capital 50,000.00; Payroll Expense 10,000.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 10,000.00; net profit -10,000.00; assets 40,000.00; liabilities 0.00; equity incl. current earnings 40,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-19 — Operating expense (utilities) paid in cash

- Policy/capability links: N8, TAX-01
- Capability check first: none (core flow)
- Assumptions:
  - [PROVISIONAL] Expense 1,000 + recoverable input VAT 140 (14%) paid in cash; expense recognised when incurred
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - Capital contribution: Dr Cash 5,000.00; Cr Owner Capital 5,000.00
  - Utilities expense: Dr Utilities Expense 1,000.00; Dr Input VAT 140.00; Cr Cash 1,140.00
- Expected ending balances: Cash 3,860.00; Input VAT 140.00; Owner Capital 5,000.00; Utilities Expense 1,000.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 1,000.00; net profit -1,000.00; assets 4,000.00; liabilities 0.00; equity incl. current earnings 4,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-20 — Bank transfer (cash to bank)

- Policy/capability links: N9, CF-01
- Capability check first: N9 (bank transfer)
- Assumptions:
  - [PROVISIONAL] CF-01: Cash and Bank are both cash equivalents; internal transfer has zero cash-flow effect and no P&L effect; no fee
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
- Expected journals (EGP; Dr / Cr):
  - Capital contribution: Dr Cash 20,000.00; Cr Owner Capital 20,000.00
  - Deposit cash to bank: Dr Bank 5,000.00; Cr Cash 5,000.00
- Expected ending balances: Cash 15,000.00; Bank 5,000.00; Owner Capital 20,000.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 0.00; net profit 0.00; assets 20,000.00; liabilities 0.00; equity incl. current earnings 20,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-21 — Owner withdrawal

- Policy/capability links: N4, CLOSE-01
- Capability check first: N4 (capital (manual journal))
- Assumptions:
  - [PROVISIONAL] Withdrawal debits Owner Drawings (contra-equity), financing outflow; not an expense; closed to Retained Earnings at close
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
- Expected journals (EGP; Dr / Cr):
  - Capital contribution: Dr Bank 100,000.00; Cr Owner Capital 100,000.00
  - Owner withdrawal: Dr Owner Drawings 10,000.00; Cr Bank 10,000.00
- Expected ending balances: Bank 90,000.00; Owner Capital 100,000.00; Owner Drawings 10,000.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 0.00; net profit 0.00; assets 90,000.00; liabilities 0.00; equity incl. current earnings 90,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-22 — Service charge and card tip

- Policy/capability links: N18, TAX-01
- Capability check first: N18 (tips/service)
- Assumptions:
  - [PROVISIONAL] Service charge 100 is revenue and VAT-able at 14%; tip 50 is a liability (pass-through), no VAT, paid out in cash
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - Capital contribution: Dr Cash 1,000.00; Cr Owner Capital 1,000.00
  - Card sale with service charge and tip: Cr Revenue 1,000.00; Cr Output VAT 140.00; Cr Service Charge Revenue 100.00; Cr Output VAT 14.00; Cr Tips Payable 50.00; Dr Bank 1,304.00
  - Tip payout to staff: Dr Tips Payable 50.00; Cr Cash 50.00
- Expected ending balances: Cash 950.00; Bank 1,304.00; Output VAT 154.00; Owner Capital 1,000.00; Revenue 1,000.00; Service Charge Revenue 100.00
- Statements: net revenue 1,100.00; COGS 0.00; gross profit 1,100.00; opex 0.00; net profit 1,100.00; assets 2,254.00; liabilities 154.00; equity incl. current earnings 2,100.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-23 — Financial close (period closing to retained earnings)

- Policy/capability links: CLOSE-01, EARN-01
- Capability check first: none (core flow)
- Assumptions:
  - [PROVISIONAL] CLOSE-01: P&L closed via Income Summary to Retained Earnings; hard-close; no VAT settlement journal
  - [PROVISIONAL] Uses the same transactions as the existing deterministic control (T1..T8); close result must equal its -2,000 loss
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
- Expected journals (EGP; Dr / Cr):
  - Capital contribution: Dr Cash 100,000.00; Cr Owner Capital 100,000.00
  - Credit purchase: Dr Inventory 20,000.00; Dr Input VAT 2,800.00; Cr Accounts Payable 22,800.00
  - Supplier payment: Dr Accounts Payable 10,000.00; Cr Cash 10,000.00
  - Cash sale: Cr Revenue 10,000.00; Cr Output VAT 1,400.00; Dr Cash 11,400.00
  - Cash sale COGS: Dr COGS 4,000.00; Cr Inventory 4,000.00
  - Rent: Dr Rent Expense 3,000.00; Cr Cash 3,000.00
  - Payroll accrual: Dr Payroll Expense 5,000.00; Cr Payroll Payable 5,000.00
  - Payroll payment: Dr Payroll Payable 3,000.00; Cr Cash 3,000.00
  - Close P&L to Income Summary: Dr Revenue 10,000.00; Cr COGS 4,000.00; Cr Rent Expense 3,000.00; Cr Payroll Expense 5,000.00; Dr Income Summary 2,000.00
  - Close Income Summary to Retained Earnings: Cr Income Summary 2,000.00; Dr Retained Earnings 2,000.00
- Expected ending balances: Cash 95,400.00; Inventory 16,000.00; Input VAT 2,800.00; Accounts Payable 12,800.00; Output VAT 1,400.00; Payroll Payable 2,000.00; Owner Capital 100,000.00; Retained Earnings -2,000.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 0.00; net profit 0.00; assets 114,200.00; liabilities 16,200.00; equity incl. current earnings 98,000.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.

## SCN-24 — Full operating month

- Policy/capability links: OPEN-01, CLOSE-01, INV-01, PAY-01
- Capability check first: none (core flow)
- Assumptions:
  - [PROVISIONAL] Composite of SCN-01..22 policies on one ledger; every PROVISIONAL tag of those scenarios applies
  - [PROVISIONAL] Bank funded by capital 30,000; month ends with hard close to Retained Earnings
  - [PROMPT] VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)
  - [PROMPT] Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)
  - [PROVISIONAL] ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation
- Expected journals (EGP; Dr / Cr):
  - Opening balances: Dr Cash 50,000.00; Dr Inventory 10,000.00; Cr Accounts Payable 8,000.00; Cr Opening Balance Equity 52,000.00
  - Capital contribution: Dr Bank 30,000.00; Cr Owner Capital 30,000.00
  - Credit purchase: Dr Inventory 20,000.00; Dr Input VAT 2,800.00; Cr Accounts Payable 22,800.00
  - Cash purchase: Dr Inventory 2,000.00; Dr Input VAT 280.00; Cr Cash 2,280.00
  - Purchase return to supplier: Dr Accounts Payable 1,140.00; Cr Inventory 1,000.00; Cr Input VAT 140.00
  - Supplier payment: Dr Accounts Payable 15,000.00; Cr Bank 15,000.00
  - Cash sales day 1: Cr Revenue 4,000.00; Cr Output VAT 560.00; Dr Cash 4,560.00
  - Cash sales day 1 COGS: Dr COGS 1,600.00; Cr Inventory 1,600.00
  - Card sales: Cr Revenue 3,000.00; Cr Output VAT 420.00; Dr Bank 3,420.00
  - Card sales COGS: Dr COGS 1,200.00; Cr Inventory 1,200.00
  - Split sale: Cr Revenue 2,000.00; Cr Output VAT 280.00; Dr Cash 1,000.00; Dr Bank 1,280.00
  - Split sale COGS: Dr COGS 800.00; Cr Inventory 800.00
  - Discounted sale: Cr Revenue 1,000.00; Cr Output VAT 126.00; Dr Sales Discounts 100.00; Dr Cash 1,026.00
  - Discounted sale COGS: Dr COGS 400.00; Cr Inventory 400.00
  - COD order: Cr Revenue 3,000.00; Cr Output VAT 420.00; Cr Delivery Revenue 200.00; Cr Output VAT 28.00; Dr Courier Cash Receivable 3,648.00
  - COD order COGS: Dr COGS 1,200.00; Cr Inventory 1,200.00
  - Courier settlement: Dr Cash 3,648.00; Cr Courier Cash Receivable 3,648.00
  - Customer invoice: Dr Accounts Receivable 11,400.00; Cr Revenue 10,000.00; Cr Output VAT 1,400.00
  - Customer receipt: Dr Bank 6,000.00; Cr Accounts Receivable 6,000.00
  - Refund 1 item: Dr Sales Returns 500.00; Dr Output VAT 70.00; Cr Cash 570.00
  - Restock: Dr Inventory 200.00; Cr COGS 200.00
  - Count shortage: Dr Inventory Shrinkage 220.00; Cr Inventory 220.00
  - Waste: Dr Waste Expense 330.00; Cr Inventory 330.00
  - Rent: Dr Rent Expense 4,000.00; Cr Bank 4,000.00
  - Utilities: Dr Utilities Expense 1,000.00; Dr Input VAT 140.00; Cr Cash 1,140.00
  - Payroll accrual: Dr Payroll Expense 10,000.00; Cr Payroll Payable 9,200.00; Cr Withholding Payable 800.00
  - Net pay: Dr Payroll Payable 9,200.00; Cr Bank 9,200.00
  - Deposit cash to bank: Dr Bank 5,000.00; Cr Cash 5,000.00
  - Owner withdrawal: Dr Owner Drawings 2,000.00; Cr Bank 2,000.00
  - Close P&L to Income Summary: Dr Revenue 23,000.00; Dr Delivery Revenue 200.00; Cr Sales Discounts 100.00; Cr Sales Returns 500.00; Cr COGS 5,000.00; Cr Inventory Shrinkage 220.00; Cr Waste Expense 330.00; Cr Utilities Expense 1,000.00; Cr Rent Expense 4,000.00; Cr Payroll Expense 10,000.00; Cr Income Summary 2,050.00
  - Close Income Summary to Retained Earnings: Dr Income Summary 2,050.00; Cr Retained Earnings 2,050.00
- Expected ending balances: Cash 51,244.00; Bank 15,500.00; Accounts Receivable 5,400.00; Inventory 25,450.00; Input VAT 3,080.00; Accounts Payable 14,660.00; Output VAT 3,164.00; Withholding Payable 800.00; Owner Capital 30,000.00; Opening Balance Equity 52,000.00; Retained Earnings 2,050.00; Owner Drawings 2,000.00
- Statements: net revenue 0.00; COGS 0.00; gross profit 0.00; opex 0.00; net profit 0.00; assets 100,674.00; liabilities 18,624.00; equity incl. current earnings 82,050.00; BS residual 0.00
- self_consistency: "PASS"
- App-side comparison: TB difference = 0; BS residual = 0; per-account balance vs oracle; report figures (P&L, BS, CF, VAT) vs above. Status: **NOT EXECUTED**.
