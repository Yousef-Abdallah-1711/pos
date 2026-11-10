import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { PaymentHistoryTable } from './PaymentHistoryTable'

const voidedPayment = {
  id: 1,
  amount_cents: 50000,
  currency: 'EGP',
  status: 'voided' as const,
  merchant_order_id: 'voided-order-1',
  created_at: '2026-11-09T12:00:00.000Z',
}

function renderPaymentHistory(locale: 'en' | 'ar', voidedLabel: string) {
  render(
    <PaymentHistoryTable
      payments={[voidedPayment]}
      total={1}
      page={1}
      pageSize={10}
      locale={locale}
      onPageChange={() => undefined}
      labels={{ voided: voidedLabel }}
    />,
  )
}

describe('PaymentHistoryTable', () => {
  it('shows the localized void status in English', () => {
    renderPaymentHistory('en', 'Voided')

    expect(screen.getByText('Voided')).toBeInTheDocument()
  })

  it('shows the localized void status in Arabic', () => {
    renderPaymentHistory('ar', 'ملغي')

    expect(screen.getByText('ملغي')).toBeInTheDocument()
  })
})
