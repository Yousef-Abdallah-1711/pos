import { render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import { BalanceSheet } from './BalanceSheet'

vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
}))

describe('BalanceSheet integrity messaging', () => {
  it('keeps unreconciled figures visible and exposes the warning to assistive technology', () => {
    render(
      <BalanceSheet
        locale="en"
        report={{
          as_of_date: '2026-11-10',
          assets: [{ account_code: '1110', account_name: 'Cash', amount: 100 }],
          liabilities: [],
          equity: [],
          total_assets: 100,
          total_liabilities: 0,
          total_equity: 0,
          equation_balanced: false,
          integrity_status: 'UNRECONCILED',
          equation_difference: 100,
          integrity_check_basis: 'existing_round2_equation_check',
        }}
      />
    )

    expect(screen.getByRole('alert')).toHaveTextContent('integrity.unreconciledTitle')
    expect(screen.getByText('metrics.totalAssets')).toBeInTheDocument()
    expect(screen.getAllByText(/100\.00/).length).toBeGreaterThan(0)
  })
})
