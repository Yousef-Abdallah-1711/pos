import { describe, expect, it } from 'vitest'
import { canExportBalanceSheet } from './financial-integrity'

describe('canExportBalanceSheet', () => {
  it('allows export only when the explicit status and equation both reconcile', () => {
    expect(canExportBalanceSheet({ integrity_status: 'RECONCILED', equation_balanced: true })).toBe(true)
  })

  it('blocks unreconciled, inconsistent, and older unverified responses', () => {
    expect(canExportBalanceSheet({ integrity_status: 'UNRECONCILED', equation_balanced: false })).toBe(false)
    expect(canExportBalanceSheet({ integrity_status: 'RECONCILED', equation_balanced: false })).toBe(false)
    expect(canExportBalanceSheet(undefined)).toBe(false)
  })
})
