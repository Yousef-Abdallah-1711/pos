#!/usr/bin/env python3
"""Scenario oracle fixtures SCN-01..SCN-24 (new file; does NOT modify deterministic_control.py).

Pure Python, imports only the standard library, NOTHING from the application. Integer piasters
(1 EGP = 100 piasters). Independent double-entry ledger. Run:  python -I oracle/scenarios.py

Writes oracle/scenarios_expected.json. Exit code 1 if any self-consistency assertion fails.

ASSUMPTION TAGS (every scenario carries them):
  PROMPT       = taken from the audit prompt / the existing deterministic control oracle
                 (EGP, VAT 14% exclusive, accrual basis, recoverable purchase VAT, no fees).
  PROVISIONAL  = UNAPPROVED accounting policy chosen only so a number can be computed
                 (see 32-accounting-policy-decision-register.md). NOT a claim about tax law,
                 NOT what the application does. If a policy is later decided differently the
                 scenario expectation must be regenerated.
Results are SELF-CONSISTENCY checks of this oracle only; they say nothing about the app.
"""
import json, os, sys
from collections import defaultdict

P = 100
def e(x):  # EGP integer/pair -> piasters (integer arithmetic only; pass ints or (egp, piasters))
    if isinstance(x, tuple):
        return x[0] * P + x[1]
    return x * P

VAT_BP = 1400
def vat(net_p):
    """PROVISIONAL ROUND-01: half-up to the piaster, per document."""
    return (net_p * VAT_BP + 5000) // 10000

# name: (normal side, class)
CHART = {
    'Cash': ('D', 'asset'), 'Bank': ('D', 'asset'),
    'Accounts Receivable': ('D', 'asset'), 'Courier Cash Receivable': ('D', 'asset'),
    'Wallet Clearing': ('D', 'asset'), 'Inventory': ('D', 'asset'), 'Input VAT': ('D', 'asset'),
    'Accounts Payable': ('C', 'liability'), 'Output VAT': ('C', 'liability'),
    'Payroll Payable': ('C', 'liability'), 'Withholding Payable': ('C', 'liability'),
    'Tips Payable': ('C', 'liability'),
    'Owner Capital': ('C', 'equity'), 'Opening Balance Equity': ('C', 'equity'),
    'Retained Earnings': ('C', 'equity'), 'Owner Drawings': ('D', 'equity'),   # contra-equity
    'Revenue': ('C', 'revenue'), 'Delivery Revenue': ('C', 'revenue'), 'Service Charge Revenue': ('C', 'revenue'),
    'Sales Discounts': ('D', 'revenue'), 'Sales Returns': ('D', 'revenue'),    # contra-revenue
    'COGS': ('D', 'expense'), 'Inventory Shrinkage': ('D', 'expense'), 'Waste Expense': ('D', 'expense'),
    'Utilities Expense': ('D', 'expense'), 'Rent Expense': ('D', 'expense'), 'Payroll Expense': ('D', 'expense'),
    'Income Summary': ('C', 'closing'),
}
CASHLIKE = ('Cash', 'Bank')   # PROVISIONAL cash-equivalent scope (CF-01)


class Ledger:
    def __init__(self):
        self.bal = defaultdict(int)       # signed debit - credit
        self.journals = []
        self.cf = defaultdict(int)        # direct-method categories
        self.ops = defaultdict(int)       # independent operational tallies (not derived from ledger)

    def post(self, name, lines, cf='operating'):
        d = sum(l[1] for l in lines); c = sum(l[2] for l in lines)
        assert d == c and d > 0, f'{name}: unbalanced {d} vs {c}'
        cash_delta = 0
        for a, dr, cr in lines:
            assert a in CHART, a
            assert dr >= 0 and cr >= 0 and not (dr and cr), f'{name}: bad line {a}'
            self.bal[a] += dr - cr
            if a in CASHLIKE:
                cash_delta += dr - cr
        if cash_delta:
            self.cf[cf] += cash_delta
        else:
            assert cf in ('operating', 'noncash'), 'non-cash journal given cash class'
        self.journals.append({'name': name, 'cf': cf if cash_delta else 'noncash',
                              'lines': [{'account': a, 'debit': dr, 'credit': cr} for a, dr, cr in lines]})

    def nat(self, a):
        side, _ = CHART[a]
        return self.bal[a] if side == 'D' else -self.bal[a]

    def cls(self, k):
        return sum(self.nat(a) for a in CHART if CHART[a][1] == k)

    def statements(self):
        rev_gross = sum(self.nat(a) for a in ('Revenue', 'Delivery Revenue', 'Service Charge Revenue'))
        contra = self.nat('Sales Discounts') + self.nat('Sales Returns')
        net_rev = rev_gross - contra
        cogs = self.nat('COGS')
        opex = sum(self.nat(a) for a in ('Inventory Shrinkage', 'Waste Expense', 'Utilities Expense',
                                         'Rent Expense', 'Payroll Expense'))
        net = net_rev - cogs - opex
        assets = self.cls('asset'); liabs = self.cls('liability')
        # equity: contra Owner Drawings is already negative via nat() (D normal on credit-class? handle explicitly)
        equity_booked = (self.nat('Owner Capital') + self.nat('Opening Balance Equity')
                         + self.nat('Retained Earnings') - self.nat('Owner Drawings'))
        income_summary = self.nat('Income Summary')
        equity = equity_booked + income_summary + net      # PROVISIONAL: current earnings presented in equity
        tb_d = sum(v for v in self.bal.values() if v > 0); tb_c = -sum(v for v in self.bal.values() if v < 0)
        return {
            'net_revenue': net_rev, 'gross_revenue': rev_gross, 'contra_revenue': contra, 'cogs': cogs,
            'gross_profit': net_rev - cogs, 'opex': opex, 'net_profit': net,
            'assets': assets, 'liabilities': liabs, 'equity_booked': equity_booked,
            'equity_incl_current_earnings': equity, 'bs_residual': assets - liabs - equity,
            'tb_debits': tb_d, 'tb_credits': tb_c,
            'cash_flow': dict(self.cf), 'cash_net_change_in_cf': sum(self.cf.values()),
            'closing_cash_and_bank': sum(self.bal[a] for a in CASHLIKE),
            'vat_input': self.nat('Input VAT'), 'vat_output': self.nat('Output VAT'),
            'vat_net_payable': self.nat('Output VAT') - self.nat('Input VAT'),
        }

    def invariants(self, opening_cash=None):
        s = self.statements()
        assert s['tb_debits'] == s['tb_credits'], 'trial balance'
        assert s['bs_residual'] == 0, f"BS residual {s['bs_residual']}"
        assert s['gross_profit'] - s['opex'] == s['net_profit']
        opening = self.cf.get('opening', 0)
        move = sum(v for k, v in self.cf.items() if k != 'opening')
        assert opening + move == s['closing_cash_and_bank'], 'cash flow reconciles to cash+bank'
        return s


def close_books(L):
    """PROVISIONAL CLOSE-01: close P&L accounts to Income Summary then to Retained Earnings."""
    lines = []
    for a, (side, k) in CHART.items():
        if k in ('revenue', 'expense') and L.bal[a] != 0:
            v = L.bal[a]
            lines.append((a, 0, v) if v > 0 else (a, -v, 0))
    net_dr = sum(l[1] for l in lines); net_cr = sum(l[2] for l in lines)
    if lines:
        # Income Summary takes the balancing side
        if net_dr > net_cr:
            lines.append(('Income Summary', 0, net_dr - net_cr))
        else:
            lines.append(('Income Summary', net_cr - net_dr, 0))
        L.post('Close P&L to Income Summary', lines)
    v = L.bal['Income Summary']
    if v:
        L.post('Close Income Summary to Retained Earnings',
               [('Income Summary', 0, v) if v > 0 else ('Income Summary', -v, 0),
                ('Retained Earnings', v, 0) if v > 0 else ('Retained Earnings', 0, -v)])


# ---------------------------------------------------------------------------------------------
# journal-building operations (shared by scenarios and the full month)
# ---------------------------------------------------------------------------------------------
def op_opening(L):
    L.post('Opening balances', [('Cash', e(50000), 0), ('Inventory', e(10000), 0),
                                ('Accounts Payable', 0, e(8000)), ('Opening Balance Equity', 0, e(52000))], cf='opening')

def op_capital(L, amt, to='Bank'):
    L.post('Capital contribution', [(to, amt, 0), ('Owner Capital', 0, amt)], cf='financing')

def op_withdrawal(L, amt, frm='Bank'):
    L.post('Owner withdrawal', [('Owner Drawings', amt, 0), (frm, 0, amt)], cf='financing')

def op_cash_purchase(L, net):
    v = vat(net)
    L.post('Cash purchase', [('Inventory', net, 0), ('Input VAT', v, 0), ('Cash', 0, net + v)])
    L.ops['purchases_net'] += net; L.ops['input_vat'] += v

def op_credit_purchase(L, net):
    v = vat(net)
    L.post('Credit purchase', [('Inventory', net, 0), ('Input VAT', v, 0), ('Accounts Payable', 0, net + v)], cf='operating')
    L.ops['purchases_net'] += net; L.ops['input_vat'] += v

def op_supplier_payment(L, amt, frm='Bank'):
    L.post('Supplier payment', [('Accounts Payable', amt, 0), (frm, 0, amt)])

def op_purchase_return(L, net):
    v = vat(net)
    L.post('Purchase return to supplier', [('Accounts Payable', net + v, 0), ('Inventory', 0, net), ('Input VAT', 0, v)])
    L.ops['purchases_net'] -= net; L.ops['input_vat'] -= v

def op_sale(L, net, tenders, cogs=0, name='Sale', discount=0, extra=None):
    """tenders: list of (account, amount). extra: list of (revenue_account, net) incl. VAT on each."""
    base = net - 0  # net is gross before discount; discount shown contra
    taxable = net - discount
    v = vat(taxable)
    lines = [('Revenue', 0, net), ('Output VAT', 0, v)]
    total = taxable + v
    if discount:
        lines.append(('Sales Discounts', discount, 0))
    for acct, amt in (extra or []):
        ve = vat(amt); lines.append((acct, 0, amt)); lines.append(('Output VAT', 0, ve)); total += amt + ve; v += ve
    assert sum(t[1] for t in tenders) == total, f'{name}: tenders {sum(t[1] for t in tenders)} != total {total}'
    for acct, amt in tenders:
        lines.append((acct, amt, 0))
    L.post(name, lines)
    L.ops['sales_net_of_discount'] += taxable + sum(a for _, a in (extra or []))
    L.ops['output_vat'] += v
    if cogs:
        L.post(name + ' COGS', [('COGS', cogs, 0), ('Inventory', 0, cogs)])
        L.ops['cogs'] += cogs
    return total, v

def total_for(net, discount=0, extra=()):
    t = net - discount; v = vat(t); tot = t + v
    for _, a in extra:
        tot += a + vat(a)
    return tot


SCENARIOS = {}
def scenario(sid, title, source_ids, assumptions):
    def deco(fn):
        SCENARIOS[sid] = dict(fn=fn, title=title, source_ids=source_ids, assumptions=assumptions)
        return fn
    return deco

PR = 'PROMPT'; PV = 'PROVISIONAL'
VAT_A = (PR, 'VAT 14% exclusive on net price; EGP (existing oracle test assumption, not asserted as law)')
ACC_A = (PR, 'Accrual basis; purchase (input) VAT recoverable (existing oracle assumption)')
RND_A = (PV, 'ROUND-01: VAT rounded half-up to the piaster per document; chosen only for computation')
NOFEE = (PR, 'No payment-processor fees, no discounts unless the scenario states otherwise')


@scenario('SCN-01', 'Opening balances', ['OPEN-01'], [
    (PV, 'OPEN-01: opening balances posted as ONE journal offset to Opening Balance Equity (evidence-backed)'),
    (PV, 'Inventory opening value 10,000 taken at book cost; AP opening 8,000'), ACC_A])
def s01():
    L = Ledger(); op_opening(L)
    s = L.invariants()
    assert L.nat('Cash') == e(50000) and s['assets'] == e(60000) and s['liabilities'] == e(8000)
    assert L.nat('Opening Balance Equity') == e(52000) and s['net_profit'] == 0
    return L

@scenario('SCN-02', 'Capital contribution', ['N4'], [
    (PV, 'Capital contribution of 100,000 into Bank is financing cash flow, equity credit to Owner Capital'), ACC_A])
def s02():
    L = Ledger(); op_capital(L, e(100000))
    s = L.invariants()
    assert s['cash_flow']['financing'] == e(100000) and s['equity_incl_current_earnings'] == e(100000)
    return L

@scenario('SCN-03', 'Cash purchase of inventory', ['INV-01', 'TAX-01'], [
    (PV, 'INV purchase capitalised to Inventory (not expensed); input VAT recoverable and held in Input VAT asset'), VAT_A, RND_A])
def s03():
    L = Ledger(); op_capital(L, e(10000), 'Cash'); op_cash_purchase(L, e(5000))
    s = L.invariants()
    assert L.nat('Inventory') == e(5000) and L.nat('Input VAT') == e(700) and L.nat('Cash') == e(10000) - e(5700)
    assert s['cash_flow']['operating'] == -e(5700) and s['net_profit'] == 0
    return L

@scenario('SCN-04', 'Credit purchase and part supplier payment', ['INV-01', 'TAX-01'], [
    (PV, 'Purchase capitalised to Inventory; AP carries gross incl. VAT; supplier paid by Bank'), VAT_A, ACC_A])
def s04():
    L = Ledger(); op_capital(L, e(100000)); op_credit_purchase(L, e(20000)); op_supplier_payment(L, e(10000))
    s = L.invariants()
    assert L.nat('Accounts Payable') == e(12800) and L.nat('Inventory') == e(20000) and L.nat('Input VAT') == e(2800)
    assert s['cash_flow']['operating'] == -e(10000)
    return L

@scenario('SCN-05', 'Purchase return to supplier', ['N16', 'N7'], [
    (PV, 'Return of goods reverses Inventory, Input VAT and AP at original cost (no restocking fee)'), VAT_A, ACC_A])
def s05():
    L = Ledger(); op_credit_purchase(L, e(20000)); op_purchase_return(L, e(4000))
    s = L.invariants()
    assert L.nat('Accounts Payable') == e(22800) - e(4560) and L.nat('Inventory') == e(16000) and L.nat('Input VAT') == e(2240)
    assert s['net_profit'] == 0
    return L

@scenario('SCN-06', 'Card sale', ['REV', 'N21'], [
    (PV, 'Card receipt settles directly into Bank at sale time, no processor fee / chargeback / settlement lag'), VAT_A, NOFEE, RND_A])
def s06():
    L = Ledger(); op_capital(L, e(5000), 'Cash'); L.post('Stock', [('Inventory', e(400), 0), ('Cash', 0, e(400))])
    tot, v = op_sale(L, e(1000), [('Bank', e(1140))], cogs=e(400), name='Card sale')
    s = L.invariants()
    assert (tot, v) == (e(1140), e(140)) and s['net_profit'] == e(600) and L.nat('Bank') == e(1140)
    return L

@scenario('SCN-07', 'Wallet sale with provider payout', ['N21'], [
    (PV, 'Wallet funds held in Wallet Clearing asset until payout; payout to Bank is operating cash; no fee'), VAT_A, NOFEE, RND_A])
def s07():
    L = Ledger()
    tot, v = op_sale(L, e(1000), [('Wallet Clearing', e(1140))], name='Wallet sale')
    assert s_cash(L) == 0
    L.post('Wallet payout', [('Bank', e(1140), 0), ('Wallet Clearing', 0, e(1140))])
    s = L.invariants()
    assert L.nat('Wallet Clearing') == 0 and s['cash_flow']['operating'] == e(1140)
    return L

def s_cash(L): return sum(L.bal[a] for a in CASHLIKE)

@scenario('SCN-08', 'Split cash + card sale', ['ROUND-01'], [
    (PV, 'Split tender: 1,000 cash + remainder card; whole-order VAT computed once on net, not per tender'), VAT_A, NOFEE, RND_A])
def s08():
    L = Ledger()
    tot, v = op_sale(L, e(2000), [('Cash', e(1000)), ('Bank', e(1280))], name='Split sale')
    s = L.invariants()
    assert tot == e(2280) and v == e(280) and L.nat('Cash') == e(1000) and L.nat('Bank') == e(1280)
    return L

@scenario('SCN-09', 'Online prepaid order (card via gateway)', ['N21', 'N22'], [
    (PV, 'Online capture recognises revenue at capture; delivery fee 0 (matches current app forced-0 fee, not a policy)'), VAT_A, NOFEE])
def s09():
    L = Ledger()
    op_sale(L, e(2000), [('Bank', e(2280))], name='Online prepaid order')
    s = L.invariants()
    assert s['net_revenue'] == e(2000) and s['vat_output'] == e(280)
    return L

@scenario('SCN-10', 'Online order cash-on-delivery with delivery fee', ['N20', 'N21'], [
    (PV, 'Delivery fee is separate revenue and VAT-able at 14%'), (PV, 'COD cash held by driver in Courier Cash Receivable until settled; revenue at delivery'),
    VAT_A, RND_A])
def s10():
    L = Ledger()
    tot, v = op_sale(L, e(3000), [('Courier Cash Receivable', total_for(e(3000), extra=[('', e(200))]))],
                     name='COD delivered', extra=[('Delivery Revenue', e(200))])
    assert tot == e(3648) and v == e(448)
    assert s_cash(L) == 0
    L.post('Courier cash settlement', [('Cash', tot, 0), ('Courier Cash Receivable', 0, tot)])
    s = L.invariants()
    assert L.nat('Courier Cash Receivable') == 0 and s['net_revenue'] == e(3200)
    return L

@scenario('SCN-11', 'Customer invoice and partial receipt', ['N22'], [
    (PV, 'Revenue recognised at invoice date (accrual); receipts clear AR; 5,400 remains in AR'), VAT_A, ACC_A])
def s11():
    L = Ledger()
    L.post('Customer invoice', [('Accounts Receivable', e(11400), 0), ('Revenue', 0, e(10000)), ('Output VAT', 0, e(1400))])
    L.post('Customer receipt (bank)', [('Bank', e(6000), 0), ('Accounts Receivable', 0, e(6000))])
    s = L.invariants()
    assert L.nat('Accounts Receivable') == e(5400) and s['net_revenue'] == e(10000) and s['cash_flow']['operating'] == e(6000)
    return L

@scenario('SCN-12', 'Order-level discount', ['N17', 'TAX-01'], [
    (PV, 'Discount ordering: VAT computed AFTER discount; discount shown in contra-revenue Sales Discounts (gross revenue 1,000)'),
    (PV, 'NOTE conflicting implementation: customer-invoice trigger taxes gross before discount (04 section 3)'), VAT_A, RND_A])
def s12():
    L = Ledger()
    tot, v = op_sale(L, e(1000), [('Cash', e(1026))], discount=e(100), name='Discounted sale')
    s = L.invariants()
    assert tot == e(1026) and v == e(126) and s['gross_revenue'] == e(1000) and s['net_revenue'] == e(900)
    return L

@scenario('SCN-13', 'Partial refund of a paid cash order', ['N2', 'REFUND-01'], [
    (PV, 'Refund is a linked reversal posted at refund date to Sales Returns (contra-revenue) and Output VAT'),
    (PV, 'Returned goods are resalable and restocked at original cost (no write-off)'), VAT_A, RND_A])
def s13():
    L = Ledger(); L.post('Stock', [('Inventory', e(400), 0), ('Owner Capital', 0, e(400))], cf='operating')
    op_sale(L, e(1000), [('Cash', e(1140))], cogs=e(400), name='Sale 2 items')
    L.post('Refund 1 of 2 items', [('Sales Returns', e(500), 0), ('Output VAT', e(70), 0), ('Cash', 0, e(570))])
    L.post('Restock refunded item', [('Inventory', e(200), 0), ('COGS', 0, e(200))])
    s = L.invariants()
    assert s['net_revenue'] == e(500) and s['cogs'] == e(200) and s['vat_output'] == e(70) and s['net_profit'] == e(300)
    assert L.nat('Cash') == e(570)
    return L

@scenario('SCN-14', 'Cancellation: unpaid order and paid-before-delivery order', ['N3', 'REFUND-01'], [
    (PV, 'Unpaid cancellation has NO journal (matches 11-posting-matrix current behaviour, but it is a policy default)'),
    (PV, 'Paid-order cancellation = full refund with stock return'), VAT_A, RND_A])
def s14():
    L = Ledger()
    assert not L.journals                              # unpaid cancel: ledger unchanged
    L.post('Stock', [('Inventory', e(400), 0), ('Owner Capital', 0, e(400))])
    op_sale(L, e(1000), [('Cash', e(1140))], cogs=e(400), name='Paid order')
    L.post('Cancel paid order - refund', [('Sales Returns', e(1000), 0), ('Output VAT', e(140), 0), ('Cash', 0, e(1140))])
    L.post('Cancel paid order - restock', [('Inventory', e(400), 0), ('COGS', 0, e(400))])
    s = L.invariants()
    assert s['net_revenue'] == 0 and s['cogs'] == 0 and s['vat_output'] == 0 and L.nat('Cash') == 0 and L.nat('Inventory') == e(400)
    return L

def wavg(qty_val_pairs):
    q = sum(a for a, _ in qty_val_pairs); v = sum(b for _, b in qty_val_pairs)
    return q, v

@scenario('SCN-15', 'Inventory consumption by sale (weighted average)', ['INV-01', 'COGS'], [
    (PV, 'INV-01: weighted-average cost, perpetual; COGS at sale = qty * (inventory value / qty on hand), half-up'),
    (PV, 'Purchases: 10 units at EGP 100 then 10 at EGP 120, both cash, VAT 14% recoverable'), VAT_A, RND_A])
def s15():
    L = Ledger(); op_capital(L, e(10000), 'Cash')
    op_cash_purchase(L, e(1000)); op_cash_purchase(L, e(1200))
    qty, val = 20, e(2200); assert L.nat('Inventory') == val
    cogs = (val * 5 * 2 + qty) // (2 * qty)        # half-up of val*5/qty
    assert cogs == e(550)
    op_sale(L, e(1000), [('Cash', e(1140))], cogs=cogs, name='Sale of 5 units @200')
    s = L.invariants()
    assert L.nat('Inventory') == e(1650) and s['cogs'] == e(550) and s['gross_profit'] == e(450)
    return L

@scenario('SCN-16', 'Inventory adjustment (count shortage)', ['N13', 'N14', 'INV-01'], [
    (PV, 'Physical count short by 2 units valued at weighted-average 110 each -> Inventory Shrinkage expense'),
    (PV, 'No VAT adjustment on shrinkage (accountant question)'), ACC_A])
def s16():
    L = Ledger(); L.post('Stock 15 units @110', [('Inventory', e(1650), 0), ('Owner Capital', 0, e(1650))])
    L.post('Count shortage 2 units', [('Inventory Shrinkage', e(220), 0), ('Inventory', 0, e(220))])
    s = L.invariants()
    assert L.nat('Inventory') == e(1430) and s['opex'] == e(220) and s['net_profit'] == -e(220)
    return L

@scenario('SCN-17', 'Waste / spoilage', ['N13', 'INV-01'], [
    (PV, '3 spoiled units at weighted-average 110 expensed to Waste Expense (separate from shrinkage)'),
    (PV, 'Input VAT previously claimed is NOT reversed (accountant question)'), ACC_A])
def s17():
    L = Ledger(); L.post('Stock 15 units @110', [('Inventory', e(1650), 0), ('Owner Capital', 0, e(1650))])
    L.post('Waste 3 units', [('Waste Expense', e(330), 0), ('Inventory', 0, e(330))])
    s = L.invariants()
    assert L.nat('Inventory') == e(1320) and s['net_profit'] == -e(330)
    return L

@scenario('SCN-18', 'Payroll accrual, settlement and withholding remittance', ['PAY-01', 'N11', 'N12'], [
    (PV, 'PAY-01: gross pay accrued to Payroll Expense when earned; employee deduction 800 (ILLUSTRATIVE, not a statutory rate) to Withholding Payable'),
    (PV, 'Net pay 9,200 settled by Bank; deduction remitted later; employer contributions and tips excluded'), ACC_A])
def s18():
    L = Ledger(); op_capital(L, e(50000))
    L.post('Payroll accrual', [('Payroll Expense', e(10000), 0), ('Payroll Payable', 0, e(9200)), ('Withholding Payable', 0, e(800))])
    s = L.invariants(); assert s['net_profit'] == -e(10000) and s['cash_flow'].get('operating', 0) == 0
    L.post('Net pay settlement', [('Payroll Payable', e(9200), 0), ('Bank', 0, e(9200))])
    L.post('Withholding remittance', [('Withholding Payable', e(800), 0), ('Bank', 0, e(800))])
    s = L.invariants()
    assert L.nat('Payroll Payable') == 0 and L.nat('Withholding Payable') == 0 and s['cash_flow']['operating'] == -e(10000)
    return L

@scenario('SCN-19', 'Operating expense (utilities) paid in cash', ['N8', 'TAX-01'], [
    (PV, 'Expense 1,000 + recoverable input VAT 140 (14%) paid in cash; expense recognised when incurred'), VAT_A, RND_A])
def s19():
    L = Ledger(); op_capital(L, e(5000), 'Cash')
    L.post('Utilities expense', [('Utilities Expense', e(1000), 0), ('Input VAT', e(140), 0), ('Cash', 0, e(1140))])
    s = L.invariants()
    assert s['net_profit'] == -e(1000) and L.nat('Input VAT') == e(140) and s['cash_flow']['operating'] == -e(1140)
    return L

@scenario('SCN-20', 'Bank transfer (cash to bank)', ['N9', 'CF-01'], [
    (PV, 'CF-01: Cash and Bank are both cash equivalents; internal transfer has zero cash-flow effect and no P&L effect; no fee'), ACC_A])
def s20():
    L = Ledger(); op_capital(L, e(20000), 'Cash')
    before = s_cash(L)
    L.post('Deposit cash to bank', [('Bank', e(5000), 0), ('Cash', 0, e(5000))], cf='operating')
    # transfer nets to zero within cash-like scope, so classify as noncash for CF purposes
    assert s_cash(L) == before and L.nat('Cash') == e(15000) and L.nat('Bank') == e(5000)
    s = L.invariants(); assert s['cash_flow']['financing'] == e(20000) and s['net_profit'] == 0
    return L

@scenario('SCN-21', 'Owner withdrawal', ['N4', 'CLOSE-01'], [
    (PV, 'Withdrawal debits Owner Drawings (contra-equity), financing outflow; not an expense; closed to Retained Earnings at close'), ACC_A])
def s21():
    L = Ledger(); op_capital(L, e(100000)); op_withdrawal(L, e(10000))
    s = L.invariants()
    assert s['equity_incl_current_earnings'] == e(90000) and s['cash_flow']['financing'] == e(90000) and s['net_profit'] == 0
    return L

@scenario('SCN-22', 'Service charge and card tip', ['N18', 'TAX-01'], [
    (PV, 'Service charge 100 is revenue and VAT-able at 14%; tip 50 is a liability (pass-through), no VAT, paid out in cash'),
    VAT_A, RND_A])
def s22():
    L = Ledger(); op_capital(L, e(1000), 'Cash')
    tot = total_for(e(1000), extra=[('', e(100))]) + e(50)
    assert tot == e(1304)
    lines = [('Revenue', 0, e(1000)), ('Output VAT', 0, e(140)), ('Service Charge Revenue', 0, e(100)), ('Output VAT', 0, e(14)),
             ('Tips Payable', 0, e(50)), ('Bank', tot, 0)]
    L.post('Card sale with service charge and tip', lines)
    L.post('Tip payout to staff', [('Tips Payable', e(50), 0), ('Cash', 0, e(50))])
    s = L.invariants()
    assert s['net_revenue'] == e(1100) and s['vat_output'] == e(154) and L.nat('Tips Payable') == 0
    return L

@scenario('SCN-23', 'Financial close (period closing to retained earnings)', ['CLOSE-01', 'EARN-01'], [
    (PV, 'CLOSE-01: P&L closed via Income Summary to Retained Earnings; hard-close; no VAT settlement journal'),
    (PV, 'Uses the same transactions as the existing deterministic control (T1..T8); close result must equal its -2,000 loss'),
    VAT_A, ACC_A])
def s23():
    L = Ledger()
    op_capital(L, e(100000), 'Cash'); op_credit_purchase(L, e(20000)); op_supplier_payment(L, e(10000), 'Cash')
    op_sale(L, e(10000), [('Cash', e(11400))], cogs=e(4000), name='Cash sale')
    L.post('Rent', [('Rent Expense', e(3000), 0), ('Cash', 0, e(3000))])
    L.post('Payroll accrual', [('Payroll Expense', e(5000), 0), ('Payroll Payable', 0, e(5000))])
    L.post('Payroll payment', [('Payroll Payable', e(3000), 0), ('Cash', 0, e(3000))])
    pre = L.invariants()
    assert pre['net_profit'] == -e(2000) and pre['equity_booked'] == e(100000) and pre['assets'] == e(114200)
    assert pre['equity_incl_current_earnings'] == e(98000)
    close_books(L)
    post = L.invariants()
    assert L.nat('Retained Earnings') == -e(2000) and post['net_profit'] == 0 and post['equity_booked'] == e(98000)
    assert post['assets'] == pre['assets'] and post['liabilities'] == pre['liabilities']
    for a in CHART:
        if CHART[a][1] in ('revenue', 'expense', 'closing'):
            assert L.bal[a] == 0, a
    return L

@scenario('SCN-24', 'Full operating month', ['OPEN-01', 'CLOSE-01', 'INV-01', 'PAY-01'], [
    (PV, 'Composite of SCN-01..22 policies on one ledger; every PROVISIONAL tag of those scenarios applies'),
    (PV, 'Bank funded by capital 30,000; month ends with hard close to Retained Earnings'), VAT_A, ACC_A, RND_A])
def s24():
    L = Ledger()
    op_opening(L); op_capital(L, e(30000))
    op_credit_purchase(L, e(20000)); op_cash_purchase(L, e(2000)); op_purchase_return(L, e(1000))
    op_supplier_payment(L, e(15000))
    op_sale(L, e(4000), [('Cash', e(4560))], cogs=e(1600), name='Cash sales day 1')
    op_sale(L, e(3000), [('Bank', e(3420))], cogs=e(1200), name='Card sales')
    op_sale(L, e(2000), [('Cash', e(1000)), ('Bank', e(1280))], cogs=e(800), name='Split sale')
    op_sale(L, e(1000), [('Cash', e(1026))], cogs=e(400), discount=e(100), name='Discounted sale')
    op_sale(L, e(3000), [('Courier Cash Receivable', e(3648))], cogs=e(1200), name='COD order',
            extra=[('Delivery Revenue', e(200))])
    L.post('Courier settlement', [('Cash', e(3648), 0), ('Courier Cash Receivable', 0, e(3648))])
    L.post('Customer invoice', [('Accounts Receivable', e(11400), 0), ('Revenue', 0, e(10000)), ('Output VAT', 0, e(1400))])
    L.post('Customer receipt', [('Bank', e(6000), 0), ('Accounts Receivable', 0, e(6000))])
    L.ops['sales_net_of_discount'] += e(10000); L.ops['output_vat'] += e(1400)
    L.post('Refund 1 item', [('Sales Returns', e(500), 0), ('Output VAT', e(70), 0), ('Cash', 0, e(570))])
    L.post('Restock', [('Inventory', e(200), 0), ('COGS', 0, e(200))])
    L.ops['refunds_net'] += e(500); L.ops['output_vat'] -= e(70); L.ops['cogs'] -= e(200)
    L.post('Count shortage', [('Inventory Shrinkage', e(220), 0), ('Inventory', 0, e(220))])
    L.post('Waste', [('Waste Expense', e(330), 0), ('Inventory', 0, e(330))])
    L.post('Rent', [('Rent Expense', e(4000), 0), ('Bank', 0, e(4000))])
    L.post('Utilities', [('Utilities Expense', e(1000), 0), ('Input VAT', e(140), 0), ('Cash', 0, e(1140))])
    L.ops['input_vat'] += e(140)
    L.post('Payroll accrual', [('Payroll Expense', e(10000), 0), ('Payroll Payable', 0, e(9200)), ('Withholding Payable', 0, e(800))])
    L.post('Net pay', [('Payroll Payable', e(9200), 0), ('Bank', 0, e(9200))])
    L.post('Deposit cash to bank', [('Bank', e(5000), 0), ('Cash', 0, e(5000))])
    op_withdrawal(L, e(2000))
    pre = L.invariants()
    # independent operational tallies must match the ledger
    gross_rev_ops = L.ops['sales_net_of_discount'] + e(100)   # discount added back to gross
    assert pre['gross_revenue'] == gross_rev_ops, (pre['gross_revenue'], gross_rev_ops)
    assert pre['vat_output'] == L.ops['output_vat'] and pre['vat_input'] == L.ops['input_vat']
    assert pre['cogs'] == L.ops['cogs']
    close_books(L)
    post = L.invariants()
    assert post['net_profit'] == 0 and L.nat('Retained Earnings') == pre['net_profit']
    assert post['assets'] == pre['assets'] and post['equity_incl_current_earnings'] == pre['equity_incl_current_earnings']
    return L


def run():
    out = {}; fails = []
    for sid, meta in SCENARIOS.items():
        try:
            L = meta['fn']()
            st = L.statements()
            out[sid] = {'title': meta['title'], 'policy_ids': meta['source_ids'],
                        'assumptions': [{'tag': t, 'text': x} for t, x in meta['assumptions']],
                        'journals': L.journals,
                        'ending_balances_piasters': {a: L.nat(a) for a in CHART if L.bal[a] != 0},
                        'statements_piasters': st, 'self_consistency': 'PASS'}
        except AssertionError as ex:
            import traceback; traceback.print_exc()
            fails.append((sid, str(ex))); out[sid] = {'title': meta['title'], 'self_consistency': 'FAIL', 'error': str(ex)}
        print(f"{sid} {meta['title']}: {out[sid]['self_consistency']}")
    here = os.path.dirname(os.path.abspath(__file__))
    with open(os.path.join(here, 'scenarios_expected.json'), 'w', encoding='utf-8') as f:
        json.dump(out, f, indent=1, ensure_ascii=False)
    print(f'{len(SCENARIOS)} scenarios, {len(SCENARIOS)-len(fails)} PASS, {len(fails)} FAIL')
    return 1 if fails else 0


if __name__ == '__main__':
    sys.exit(run())
