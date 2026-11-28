# 35a — Independent review of the S13 oracle (`oracle/deterministic_control.py`)

Method: recomputed by hand from the prompt figures (not by running the script's own functions). The original files were **not edited**.

## Recomputation
| Item | Hand result | Oracle | Match |
|---|---|---|---|
| Purchase VAT 14% × 20,000 | 2,800 | 2,800 | ✔ |
| Sale VAT 14% × 10,000 | 1,400 | 1,400 | ✔ |
| Cash: +100,000 −10,000 +11,400 −3,000 −3,000 | 95,400 | 95,400 | ✔ |
| Inventory: 20,000 − 4,000 | 16,000 | 16,000 | ✔ |
| Input VAT | 2,800 | 2,800 | ✔ |
| A/P: 22,800 − 10,000 | 12,800 | 12,800 | ✔ |
| Output VAT / Payroll payable (5,000−3,000) | 1,400 / 2,000 | same | ✔ |
| Net result: 10,000 − 4,000 − 3,000 − 5,000 | −2,000 (loss) | −2,000 | ✔ |
| Assets 95,400+16,000+2,800 | 114,200 | 114,200 | ✔ |
| Liabilities 12,800+1,400+2,000 | 16,200 | 16,200 | ✔ |
| Equity 100,000 − 2,000 | 98,000 | 98,000 | ✔ |
| A = L + E: 16,200 + 98,000 | 114,200 | residual 0 | ✔ |
| TB debits (95,400+16,000+2,800+4,000+3,000+5,000) / credits (12,800+1,400+2,000+100,000+10,000) | 126,200 / 126,200 | same | ✔ |
| Cash flow: operating −10,000+11,400−3,000−3,000 = −4,600; financing +100,000; net 95,400 | ✔ | ✔ | ✔ |

## Observations (corrections to documentation, not to values)
1. **Self-consistency only.** The 30 checks compare the script to figures *stated in the prompt* (and the script derives them from the same journals). It proves arithmetic, not that the journals are the right accounting.
2. **Policy assumptions baked in** (PROMPT-tagged, unapproved): accrual basis; purchase VAT recoverable as an asset; no withholding/insurance on payroll; COGS fixed at 4,000 (no costing method); VAT exact at 14% with no rounding case (the script asserts exactness, so it cannot catch rounding defects — covered by SCN fixtures with half-up, [35](35-financial-scenario-test-specification.md)).
3. **Cash-flow classification:** all cash in/out other than capital is "operating" (VAT receipt/payment included). That is the direct-method convention used here; the application uses an indirect/prefix heuristic (VER-002), so an exact match of *classification* is not guaranteed even if totals agree.
4. **Equity includes current-period earnings** (capital + net loss) while the application's Balance Sheet cannot (VER-001). Expected app result: BS will not balance by exactly the net result (2,000 here) unless a closing entry exists. This is a **prediction (INFERRED)**, not a verified fact.
5. **T4 sale is cash**: with the POS flow, only cash orders are auto-marked paid (VER-008); the scenario must be entered so the order reaches `paid`, otherwise no revenue journal exists.
6. **Date basis not specified** in the prompt: all transactions are assumed within one period; period boundaries are covered separately.
7. No errors found in inputs, formulas, debits, credits, closing balances or statement relationships. No value was changed.
