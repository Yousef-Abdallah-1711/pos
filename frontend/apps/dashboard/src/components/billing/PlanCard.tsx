'use client'

import { CheckCircle2, CreditCard, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SubscriptionPlan } from '@/lib/api'

interface PlanCardProps {
  plan: SubscriptionPlan
  disabled?: boolean
  isActive: boolean
  isLoading: boolean
  locale?: 'en' | 'ar'
  subscribeLabel?: string
  activeLabel?: string
  annualSaveLabel?: string
  translateFeature?: (feature: string) => string
  onSubscribe: (planID: string) => void
}

function formatMoneyFromCents(amountCents: number, currency: string, locale: 'en' | 'ar') {
  return new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-US', {
    style: 'currency',
    currency: currency || 'EGP',
    maximumFractionDigits: 2,
  }).format((amountCents || 0) / 100)
}

function isAnnualPlan(plan: SubscriptionPlan) {
  const period = (plan.period || '').toLowerCase()
  const id = (plan.id || '').toLowerCase()
  return period === 'year' || period === 'annual' || id.includes('annual')
}

export function PlanCard({
  plan,
  isActive,
  disabled = false,
  isLoading,
  locale = 'en',
  subscribeLabel = 'Subscribe Now',
  activeLabel = 'Already Active',
  annualSaveLabel = 'Save approximately 44.4%',
  translateFeature,
  onSubscribe,
}: PlanCardProps) {
  const periodLabel =
    locale === 'ar'
      ? ((plan.period || '').toLowerCase() === 'year' || (plan.period || '').toLowerCase() === 'annual' ? 'سنة' : 'شهر')
      : (plan.period || '')

  return (
    <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-5">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <p className="text-lg font-semibold text-gray-900 dark:text-white">{plan.name}</p>
            {isAnnualPlan(plan) && (
              <span className="inline-flex rounded-full bg-emerald-100 px-2 py-0.5 text-[11px] font-medium text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300">
                {annualSaveLabel}
              </span>
            )}
          </div>
          <p className="text-gray-500 dark:text-gray-400 text-sm mt-1">
            {formatMoneyFromCents(plan.price_cents, plan.currency, locale)} / {periodLabel}
          </p>
        </div>
        <Button
          onClick={() => onSubscribe(plan.id)}
          disabled={disabled || isLoading || isActive}
          className="gap-2"
        >
          {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <CreditCard className="h-4 w-4" />}
          {isActive ? activeLabel : subscribeLabel}
        </Button>
      </div>

      {plan.features && plan.features.length > 0 && (
        <ul className="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-2 text-sm text-gray-600 dark:text-gray-300">
          {plan.features.map((feature) => (
            <li key={feature} className="flex items-start gap-2">
              <CheckCircle2 className="h-4 w-4 mt-0.5 text-green-500" />
              <span>{translateFeature ? translateFeature(feature) : feature}</span>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
