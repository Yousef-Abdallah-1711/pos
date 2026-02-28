'use client'

import { Button } from '@/components/ui/button'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Loader2 } from 'lucide-react'
import { PaymentHistoryItem } from '@/lib/api'

interface PaymentHistoryTableProps {
  payments: PaymentHistoryItem[]
  total: number
  page: number
  pageSize: number
  locale?: 'en' | 'ar'
  isRTL?: boolean
  isLoading?: boolean
  onPageChange: (page: number) => void
  labels?: {
    date?: string
    amount?: string
    status?: string
    method?: string
    card?: string
    reference?: string
    empty?: string
    pagination?: string
    previous?: string
    next?: string
    loading?: string
    cardFallback?: string
  }
}

function formatMoneyFromCents(amountCents: number, currency: string, locale: 'en' | 'ar') {
  return new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-US', {
    style: 'currency',
    currency: currency || 'EGP',
    maximumFractionDigits: 2,
  }).format((amountCents || 0) / 100)
}

function formatDate(value: string | null | undefined, locale: 'en' | 'ar') {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString(locale === 'ar' ? 'ar-EG' : 'en-US')
}

function getStatusPill(status?: string) {
  switch (status) {
    case 'success':
      return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300'
    case 'failed':
      return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
    case 'pending':
      return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-200'
    case 'refunded':
      return 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300'
    default:
      return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300'
  }
}

export function PaymentHistoryTable({
  payments,
  total,
  page,
  pageSize,
  locale = 'en',
  isRTL = false,
  isLoading = false,
  onPageChange,
  labels,
}: PaymentHistoryTableProps) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize))

  return (
    <div className="mt-4">
      <Table dir={isRTL ? 'rtl' : 'ltr'}>
        <TableHeader>
          <TableRow>
            <TableHead>{labels?.date || 'Date'}</TableHead>
            <TableHead>{labels?.amount || 'Amount'}</TableHead>
            <TableHead>{labels?.status || 'Status'}</TableHead>
            <TableHead>{labels?.method || 'Payment Method'}</TableHead>
            <TableHead>{labels?.card || 'Card'}</TableHead>
            <TableHead>{labels?.reference || 'Reference'}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading ? (
            <TableRow>
              <TableCell colSpan={6} className="text-center py-10 text-gray-500 dark:text-gray-400">
                <span className="inline-flex items-center gap-2">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  {labels?.loading || 'Loading payment history...'}
                </span>
              </TableCell>
            </TableRow>
          ) : payments.length === 0 ? (
            <TableRow>
              <TableCell colSpan={6} className="text-center py-8 text-gray-500 dark:text-gray-400">
                {labels?.empty || 'No payments yet'}
              </TableCell>
            </TableRow>
          ) : (
            payments.map((payment) => (
              <TableRow key={payment.id}>
                <TableCell>{formatDate(payment.created_at, locale)}</TableCell>
                <TableCell>{formatMoneyFromCents(payment.amount_cents, payment.currency, locale)}</TableCell>
                <TableCell>
                  <span
                    className={`inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium ${getStatusPill(payment.status)}`}
                  >
                    {payment.status.toUpperCase()}
                  </span>
                </TableCell>
                <TableCell>{payment.payment_method || '-'}</TableCell>
                <TableCell>
                  {payment.card_brand || payment.card_last_four
                    ? `${payment.card_brand || (labels?.cardFallback || 'Card')} **** ${payment.card_last_four || '----'}`
                    : '-'}
                </TableCell>
                <TableCell className="font-mono text-xs">{payment.merchant_order_id}</TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>

      {total > pageSize && (
        <div className={`mt-4 flex items-center justify-between ${isRTL ? 'flex-row-reverse' : ''}`}>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {labels?.pagination || 'Page'} {page} / {totalPages}
          </p>
          <div className={`flex items-center gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
            <Button
              variant="outline"
              size="sm"
              onClick={() => onPageChange(Math.max(1, page - 1))}
              disabled={page <= 1 || isLoading}
            >
              {labels?.previous || 'Previous'}
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => onPageChange(Math.min(totalPages, page + 1))}
              disabled={page >= totalPages || isLoading}
            >
              {labels?.next || 'Next'}
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}
