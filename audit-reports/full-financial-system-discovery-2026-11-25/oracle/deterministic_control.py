#!/usr/bin/env python3
"""Independent financial oracle for the Section 13 deterministic control scenario.

Pure Python + exact integer piasters (1 EGP = 100 piasters). It imports NOTHING from the
application. It builds a double-entry ledger from the scenario transactions, then derives
Trial Balance, P&L, Balance Sheet, Cash Flow (direct method) and VAT position, and checks
every invariant. Output: oracle/deterministic_control_expected.json and
16-financial-oracle-results.csv (expected side only; 'actual' columns are filled during runtime tests).

Policy assumptions (TEST ASSUMPTIONS from the prompt, not claims about the app or tax law):
accrual basis, recoverable purchase VAT, VAT 14%, no opening inventory, no fees/discounts.
"""
import csv, json, os, sys
from collections import defaultdict

P = 100  # piasters per EGP


def egp(x):  # EGP (int or str decimal) -> piasters, exact
    from decimal import Decimal
    return int((Decimal(str(x)) * P).to_integral_exact())


def fmt(p):  # piasters -> 'EGP 1,234.56'
    sign = '-' if p < 0 else ''
    p = abs(p)
    return f"{sign}{p // P:,}.{p % P:02d}"


# account -> (normal side, statement class)
ACCOUNTS = {
    'Cash':            ('D', 'asset'),
    'Inventory':       ('D', 'asset'),
    'Input VAT':       ('D', 'asset'),
    'Accounts Payable':('C', 'liability'),
    'Output VAT':      ('C', 'liability'),
    'Payroll Payable': ('C', 'liability'),
    'Owner Capital':   ('C', 'equity'),
    'Revenue':         ('C', 'revenue'),
    'COGS':            ('D', 'expense'),
    'Rent Expense':    ('D', 'expense'),
    'Payroll Expense': ('D', 'expense'),
}

VAT_BP = 1400  # 14.00% in basis points
def vat(net_p):
    # exact: net * 14% ; require exact result for the scenario (asserted below)
    num = net_p * VAT_BP
    assert num % 10000 == 0, 'VAT not exact in piasters: rounding policy would be needed'
    return num // 10000

# ---- scenario ----------------------------------------------------------------------------
capital   = egp(100000)
purchase  = egp(20000);  purchase_vat = vat(purchase)
sup_pay   = egp(10000)
sale      = egp(10000);  sale_vat     = vat(sale)
cogs      = egp(4000)
rent      = egp(3000)
pay_accr  = egp(5000)
pay_paid  = egp(3000)

assert purchase_vat == egp(2800), purchase_vat
assert sale_vat == egp(1400), sale_vat

JOURNALS = [
    ('T1 Capital',          [('Cash', capital, 0), ('Owner Capital', 0, capital)]),
    ('T2 Purchase on credit',[('Inventory', purchase, 0), ('Input VAT', purchase_vat, 0), ('Accounts Payable', 0, purchase + purchase_vat)]),
    ('T3 Supplier payment', [('Accounts Payable', sup_pay, 0), ('Cash', 0, sup_pay)]),
    ('T4 Cash sale',        [('Cash', sale + sale_vat, 0), ('Revenue', 0, sale), ('Output VAT', 0, sale_vat)]),
    ('T5 COGS',             [('COGS', cogs, 0), ('Inventory', 0, cogs)]),
    ('T6 Rent paid cash',   [('Rent Expense', rent, 0), ('Cash', 0, rent)]),
    ('T7 Payroll accrual',  [('Payroll Expense', pay_accr, 0), ('Payroll Payable', 0, pay_accr)]),
    ('T8 Payroll payment',  [('Payroll Payable', pay_paid, 0), ('Cash', 0, pay_paid)]),
]

# ---- ledger ------------------------------------------------------------------------------
bal = defaultdict(int)  # signed: debit - credit
for name, lines in JOURNALS:
    d = sum(l[1] for l in lines); c = sum(l[2] for l in lines)
    assert d == c, f'{name} unbalanced {d} != {c}'          # double-entry integrity
    for acct, dr, cr in lines:
        assert acct in ACCOUNTS
        bal[acct] += dr - cr

def natural(acct):
    side, _ = ACCOUNTS[acct]
    return bal[acct] if side == 'D' else -bal[acct]

ending = {a: natural(a) for a in ACCOUNTS}

# ---- statements --------------------------------------------------------------------------
tb_debit  = sum(v for a, v in bal.items() if v > 0)
tb_credit = -sum(v for a, v in bal.items() if v < 0)
revenue = ending['Revenue']; cogs_t = ending['COGS']
gross = revenue - cogs_t
opex  = ending['Rent Expense'] + ending['Payroll Expense']
net   = gross - opex                                  # negative = loss
assets = sum(ending[a] for a in ACCOUNTS if ACCOUNTS[a][1] == 'asset')
liabs  = sum(ending[a] for a in ACCOUNTS if ACCOUNTS[a][1] == 'liability')
equity = ending['Owner Capital'] + net                # capital + current-period earnings
residual = assets - liabs - equity

# direct-method cash flow: classify each cash-touching journal
cf = {'operating': 0, 'investing': 0, 'financing': 0}
cash_class = {'T1 Capital': 'financing', 'T3 Supplier payment': 'operating', 'T4 Cash sale': 'operating',
              'T6 Rent paid cash': 'operating', 'T8 Payroll payment': 'operating'}
for name, lines in JOURNALS:
    cash_delta = sum(dr - cr for a, dr, cr in lines if a == 'Cash')
    if cash_delta:
        cf[cash_class[name]] += cash_delta
net_cash = sum(cf.values())

vat_in, vat_out = ending['Input VAT'], ending['Output VAT']

expected = {
  'ending_balances_egp': {a: fmt(v) for a, v in ending.items()},
  'net_profit_egp': fmt(net),
  'pnl': {'revenue': fmt(revenue), 'cogs': fmt(cogs_t), 'gross_profit': fmt(gross), 'opex': fmt(opex), 'net': fmt(net)},
  'balance_sheet': {'assets': fmt(assets), 'liabilities': fmt(liabs), 'equity_incl_current_period': fmt(equity), 'residual': fmt(residual)},
  'trial_balance': {'debits': fmt(tb_debit), 'credits': fmt(tb_credit), 'difference': fmt(tb_debit - tb_credit)},
  'cash_flow': {k: fmt(v) for k, v in cf.items()} | {'net_change': fmt(net_cash), 'closing_cash': fmt(ending['Cash'])},
  'vat': {'input': fmt(vat_in), 'output': fmt(vat_out), 'net_recoverable': fmt(vat_in - vat_out)},
}

# ---- compare against the figures stated in the prompt (checks the PROMPT, not the app) ----
PROMPT = {
  'Cash': 95400, 'Inventory': 16000, 'Input VAT': 2800, 'Accounts Payable': 12800, 'Output VAT': 1400,
  'Payroll Payable': 2000, 'Owner Capital': 100000, 'Revenue': 10000, 'COGS': 4000,
  'Rent Expense': 3000, 'Payroll Expense': 5000,
}
checks = []
def chk(name, exp_p, act_p): checks.append((name, exp_p, act_p, 'PASS' if exp_p == act_p else 'FAIL'))
for a, e in PROMPT.items(): chk(f'ending {a}', egp(e), ending[a])
chk('net profit (loss)', egp(-2000), net)
chk('gross profit', egp(6000), gross); chk('opex', egp(8000), opex)
chk('total assets', egp(114200), assets); chk('total liabilities', egp(16200), liabs)
chk('total equity incl. current-period', egp(98000), equity); chk('BS residual', 0, residual)
chk('TB debits', egp(126200), tb_debit); chk('TB credits', egp(126200), tb_credit); chk('TB difference', 0, tb_debit - tb_credit)
chk('CF operating', egp(-4600), cf['operating']); chk('CF investing', 0, cf['investing']); chk('CF financing', egp(100000), cf['financing'])
chk('CF net change', egp(95400), net_cash); chk('CF closing cash', egp(95400), ending['Cash'])
chk('VAT input', egp(2800), vat_in); chk('VAT output', egp(1400), vat_out); chk('VAT net recoverable', egp(1400), vat_in - vat_out)
chk('invariant: opening cash(0)+net movement = closing', ending['Cash'], 0 + net_cash)

if __name__ == '__main__':
    here = os.path.dirname(os.path.abspath(__file__))
    root = os.path.dirname(here)
    with open(os.path.join(here, 'deterministic_control_expected.json'), 'w', encoding='utf-8') as f:
        json.dump(expected, f, indent=2, ensure_ascii=False)
    out = os.path.join(root, '16-financial-oracle-results.csv')
    with open(out, 'w', newline='', encoding='utf-8') as f:
        w = csv.writer(f)
        w.writerow(['test_id', 'scenario', 'metric', 'prompt_stated_piasters', 'oracle_expected_piasters', 'oracle_vs_prompt', 'app_actual_piasters', 'difference_piasters', 'status', 'notes'])
        for i, (name, exp_p, act_p, st) in enumerate(checks, 1):
            w.writerow([f'ORC-{i:03d}', 'S13-deterministic-control', name, exp_p, act_p, st, '', '', 'NOT TESTED', 'app actual pending runtime test'])
    fails = [c for c in checks if c[3] == 'FAIL']
    print(json.dumps(expected, indent=2))
    print(f'\n{len(checks)} oracle checks vs prompt figures: {len(checks)-len(fails)} PASS, {len(fails)} FAIL')
    for c in fails: print('FAIL', c)
    sys.exit(1 if fails else 0)
