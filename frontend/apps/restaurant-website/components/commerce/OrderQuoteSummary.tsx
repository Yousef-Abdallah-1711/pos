'use client'
import type { useOrderQuote } from '@/lib/hooks/use-order-quote'
import { formatPrice } from '@/lib/utils/format'

type Props = { state: ReturnType<typeof useOrderQuote>; locale: 'en' | 'ar'; totalClassName?: string }
export function OrderQuoteSummary({ state, locale, totalClassName }: Props) {
  const arabic = locale === 'ar'
  if (state.isError) return <div role="alert"><p>{arabic ? 'تعذر تأكيد الإجمالي. يرجى المحاولة مرة أخرى.' : 'Unable to confirm totals. Please retry.'}</p><button type="button" onClick={() => void state.refetch()} className="min-h-12 underline">{arabic ? 'إعادة المحاولة' : 'Retry'}</button></div>
  if (!state.data) return <p role="status">{arabic ? 'جارٍ حساب الإجمالي...' : 'Calculating totals...'}</p>
  const quote = state.data
  const lines = [
    [arabic ? 'المجموع الفرعي' : 'Subtotal', quote.subtotal],
    [arabic ? 'الضريبة' : 'Tax', quote.tax_amount],
    [arabic ? 'التوصيل' : 'Delivery', quote.delivery_fee],
  ] as const
  return <div className="space-y-3" aria-busy={state.isFetching}>
    {lines.map(([label, amount]) => <div key={label} className="flex flex-wrap justify-between gap-2"><span>{label}</span><span>{formatPrice(amount, locale)}</span></div>)}
    <div className="flex flex-wrap justify-between gap-2 border-t pt-3 font-bold"><span>{arabic ? 'الإجمالي' : 'Total'}</span><span className={totalClassName}>{formatPrice(quote.total_amount, locale)}</span></div>
  </div>
}
