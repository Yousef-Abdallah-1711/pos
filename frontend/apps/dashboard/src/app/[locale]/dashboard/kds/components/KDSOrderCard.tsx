'use client'

import { AlertTriangle, RotateCcw, SquareArrowOutUpRight } from 'lucide-react'
import { KDSOrder, KDSOrderStatus } from '@/stores/kdsStore'
import { KDSTimer } from './KDSTimer'

interface KDSOrderCardProps {
  order: KDSOrder
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  canBump: boolean
  canRecall: boolean
  isBusy?: boolean
  warningMinutes?: number
  dangerMinutes?: number
  onBump: (orderId: number) => void | Promise<void>
  onRecall: (orderId: number) => void | Promise<void>
}

const getStatusBadgeClass = (status: KDSOrderStatus) => {
  switch (status) {
    case 'new':
      return 'bg-sky-100 text-sky-800 border-sky-300 dark:bg-sky-500/20 dark:text-sky-300 dark:border-sky-400/40'
    case 'in_progress':
      return 'bg-amber-100 text-amber-800 border-amber-300 dark:bg-amber-500/20 dark:text-amber-300 dark:border-amber-400/40'
    case 'ready':
      return 'bg-emerald-100 text-emerald-800 border-emerald-300 dark:bg-emerald-500/20 dark:text-emerald-300 dark:border-emerald-400/40'
    case 'completed':
      return 'bg-gray-100 text-gray-700 border-gray-300 dark:bg-slate-500/20 dark:text-slate-300 dark:border-slate-400/40'
    case 'recalled':
      return 'bg-violet-100 text-violet-800 border-violet-300 dark:bg-violet-500/20 dark:text-violet-300 dark:border-violet-400/40'
    default:
      return 'bg-gray-100 text-gray-700 border-gray-300 dark:bg-slate-500/20 dark:text-slate-300 dark:border-slate-400/40'
  }
}

const getStatusLabel = (
  status: KDSOrderStatus,
  t: (path: string, fallback?: string) => string
) => {
  switch (status) {
    case 'new':
      return t('kds.status.new', 'New')
    case 'in_progress':
      return t('kds.status.in_progress', 'In Progress')
    case 'ready':
      return t('kds.status.ready', 'Ready')
    case 'completed':
      return t('kds.status.completed', 'Completed')
    case 'recalled':
      return t('kds.status.recalled', 'Recalled')
    default:
      return status
  }
}

const getNextStatus = (status: KDSOrderStatus): KDSOrderStatus | null => {
  switch (status) {
    case 'new':
      return 'in_progress'
    case 'in_progress':
      return 'ready'
    case 'ready':
      return 'completed'
    case 'recalled':
      return 'ready'
    default:
      return null
  }
}

export function KDSOrderCard({
  order,
  locale,
  t,
  canBump,
  canRecall,
  isBusy = false,
  warningMinutes = 10,
  dangerMinutes = 20,
  onBump,
  onRecall,
}: KDSOrderCardProps) {
  const isRTL = locale === 'ar'
  const nextStatus = getNextStatus(order.status)
  const canShowBump = canBump && nextStatus !== null
  const canShowRecall = canRecall && order.status === 'completed'
  const hasNotes = typeof order.notes === 'string' && order.notes.trim().length > 0
  const orderSource = order.order_source?.trim().toLowerCase() || ''
  const notesLabel =
    orderSource === 'pos'
      ? t('pos.orderNotes', 'Order notes')
      : t('kds.specialInstructionsOptional', 'Special Instructions (Optional)')
  const customerName = order.customer_name?.trim() || t('kds.unknownCustomer', 'Unknown customer')
  const hasTableContext = typeof order.table_number === 'number' && order.table_number > 0
  const tableLabel = hasTableContext
    ? `${t('navigation.tables', 'Tables')} #${order.table_number}${order.table_name ? ` - ${order.table_name}` : ''}`
    : null

  return (
    <article className="h-full rounded-xl border border-gray-200 bg-white p-4 shadow-sm dark:border-slate-700 dark:bg-slate-900/90">
      <header className={`mb-3 flex items-start justify-between gap-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
        <div className={`min-w-0 ${isRTL ? 'text-right' : 'text-left'}`}>
          <div className={`mb-1 flex items-center gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
            <h3 className="truncate text-lg font-bold text-gray-900 dark:text-white">
              #{order.order_number}
            </h3>
            {order.priority > 0 && (
              <span className="inline-flex items-center rounded-full border border-red-300 bg-red-100 px-2 py-0.5 text-xs font-semibold uppercase tracking-wide text-red-800 dark:border-red-400/40 dark:bg-red-500/20 dark:text-red-200">
                {t('kds.rush', 'Rush')}
              </span>
            )}
          </div>
          <p className="truncate text-sm text-gray-600 dark:text-slate-300">
            {t('kds.customer', 'Customer')}: {customerName}
          </p>
          {tableLabel && (
            <p className="mt-1 truncate text-xs font-semibold text-primary-700 dark:text-primary-300">
              {tableLabel}
            </p>
          )}
        </div>

        <div className={`flex flex-col gap-2 ${isRTL ? 'items-start' : 'items-end'}`}>
          <KDSTimer
            elapsedSeconds={order.elapsed_seconds}
            warningMinutes={warningMinutes}
            dangerMinutes={dangerMinutes}
            isStatic={order.status === 'completed'}
          />
          <span
            className={`inline-flex rounded-full border px-2 py-0.5 text-xs font-medium ${getStatusBadgeClass(order.status)}`}
          >
            {getStatusLabel(order.status, t)}
          </span>
        </div>
      </header>

      <section className="mb-4 space-y-2">
        {(order.items || []).map((item, index) => (
          <div key={`${item.product_name}-${index}`} className={`${isRTL ? 'text-right' : 'text-left'}`}>
            <p className="text-sm text-gray-900 dark:text-slate-100">
              <span className="font-semibold">{item.quantity}x</span>{' '}
              {item.product_name}
            </p>
            {item.special_instructions && (
              <p className="mt-1 rounded-md border border-amber-300 bg-amber-50 px-2 py-1 text-xs text-amber-800 dark:border-amber-400/40 dark:bg-amber-500/15 dark:text-amber-200">
                {item.special_instructions}
              </p>
            )}
            {item.add_ons && (
              <p className="mt-1 text-xs text-gray-500 dark:text-slate-400">
                {item.add_ons}
              </p>
            )}
          </div>
        ))}

        {hasNotes && (
          <div className={`${isRTL ? 'text-right' : 'text-left'}`}>
            <p className="mb-1 text-xs font-semibold uppercase tracking-wide text-amber-700 dark:text-amber-300">
              {notesLabel}
            </p>
            <p className="rounded-md border border-amber-300 bg-amber-50 px-2 py-1 text-xs text-amber-800 dark:border-amber-400/40 dark:bg-amber-500/15 dark:text-amber-100">
              {order.notes}
            </p>
          </div>
        )}
      </section>

      <footer className="space-y-2">
        {canShowBump && (
          <button
            type="button"
            onClick={() => {
              void onBump(order.id)
            }}
            disabled={isBusy}
            className="touch-manipulation flex min-h-14 w-full items-center justify-center gap-2 rounded-lg bg-primary-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-60"
          >
            <SquareArrowOutUpRight className="h-4 w-4" />
            {isBusy ? t('kds.bumping', 'Bumping...') : t('kds.bump', 'Bump')}
          </button>
        )}

        {canShowRecall && (
          <button
            type="button"
            onClick={() => {
              void onRecall(order.id)
            }}
            disabled={isBusy}
            className="touch-manipulation flex min-h-12 w-full items-center justify-center gap-2 rounded-lg border border-gray-300 bg-white px-4 py-3 text-sm font-semibold text-gray-800 transition hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700"
          >
            <RotateCcw className="h-4 w-4" />
            {isBusy ? t('kds.recalling', 'Recalling...') : t('kds.recall', 'Recall')}
          </button>
        )}

        {nextStatus && (
          <p className={`flex items-center gap-1 text-xs text-gray-500 dark:text-slate-400 ${isRTL ? 'flex-row-reverse justify-end' : ''}`}>
            <AlertTriangle className="h-3.5 w-3.5" />
            {getStatusLabel(order.status, t)} {'->'} {getStatusLabel(nextStatus, t)}
          </p>
        )}
      </footer>
    </article>
  )
}
