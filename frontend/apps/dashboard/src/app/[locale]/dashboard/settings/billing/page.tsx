'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import { usePathname } from 'next/navigation'
import {
  cancelSubscription,
  CurrentSubscriptionResponse,
  PaymentHistoryItem,
  SubscriptionPlan,
  getCurrentSubscription,
  getPaymentHistory,
  getSubscriptionPlans,
} from '@/lib/api'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { Loader2 } from 'lucide-react'
import { PlanCard } from '@/components/billing/PlanCard'
import { PaymentHistoryTable } from '@/components/billing/PaymentHistoryTable'
import { PaymentModal } from '@/components/billing/PaymentModal'
import { usePayment } from '@/hooks/usePayment'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { toast } from 'sonner'

const PAGE_SIZE = 10

function formatDate(value: string | null | undefined, locale: 'en' | 'ar') {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString(locale === 'ar' ? 'ar-EG' : 'en-US')
}

function getSubscriptionPill(status?: string) {
  switch (status) {
    case 'active':
      return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300'
    case 'trial':
      return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
    case 'expired':
      return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
    case 'disabled':
      return 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-300'
    default:
      return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300'
  }
}

function formatPlanPrice(plan: SubscriptionPlan, locale: 'en' | 'ar') {
  return new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-US', {
    style: 'currency',
    currency: plan.currency || 'EGP',
    maximumFractionDigits: 2,
  }).format((plan.price_cents || 0) / 100)
}

function isAnnualPlan(plan: SubscriptionPlan) {
  const period = (plan.period || '').toLowerCase()
  const id = (plan.id || '').toLowerCase()
  return period === 'year' || period === 'annual' || id.includes('annual')
}

export default function BillingSettingsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  const [plans, setPlans] = useState<SubscriptionPlan[]>([])
  const [currentSubscription, setCurrentSubscription] = useState<CurrentSubscriptionResponse | null>(null)
  const [history, setHistory] = useState<PaymentHistoryItem[]>([])
  const [totalPayments, setTotalPayments] = useState(0)
  const [page, setPage] = useState(1)
  const [selectedPlanID, setSelectedPlanID] = useState('')
  const [isLoading, setIsLoading] = useState(true)
  const [isHistoryLoading, setIsHistoryLoading] = useState(false)
  const [errorMessage, setErrorMessage] = useState('')
  const [isCancelDialogOpen, setIsCancelDialogOpen] = useState(false)
  const [isCancelling, setIsCancelling] = useState(false)
  const isRTL = locale === 'ar'

  const loadCurrentAndPlans = useCallback(async () => {
    setErrorMessage('')
    const [plansResult, currentResult] = await Promise.allSettled([
      getSubscriptionPlans(),
      getCurrentSubscription(),
    ])

    let resolvedPlans: SubscriptionPlan[] = []
    if (plansResult.status === 'fulfilled') {
      resolvedPlans = plansResult.value
      setPlans(resolvedPlans)
    } else {
      setErrorMessage(t('settings.billingPage.loadPlansFailed', 'Failed to load plans'))
    }

    let resolvedCurrent: CurrentSubscriptionResponse | null = null
    if (currentResult.status === 'fulfilled') {
      resolvedCurrent = currentResult.value
      setCurrentSubscription(resolvedCurrent)
    } else {
      setCurrentSubscription(null)
    }

    if (resolvedPlans.length > 0) {
      const preferredPlanID =
        resolvedCurrent?.current_plan?.id ||
        resolvedCurrent?.plan?.id ||
        ''
      setSelectedPlanID((previousID) => {
        if (preferredPlanID && resolvedPlans.some((plan) => plan.id === preferredPlanID)) {
          return preferredPlanID
        }
        if (previousID && resolvedPlans.some((plan) => plan.id === previousID)) {
          return previousID
        }
        return resolvedPlans[0].id
      })
    }
  }, [t])

  const loadPaymentHistory = useCallback(
    async (targetPage: number) => {
      setIsHistoryLoading(true)
      try {
        const response = await getPaymentHistory(targetPage, PAGE_SIZE)
        setHistory(response.payments)
        setTotalPayments(response.total)
      } catch (error) {
        setHistory([])
        setTotalPayments(0)
      } finally {
        setIsHistoryLoading(false)
      }
    },
    []
  )

  const refreshBillingData = useCallback(async () => {
    await Promise.all([loadCurrentAndPlans(), loadPaymentHistory(page)])
  }, [loadCurrentAndPlans, loadPaymentHistory, page])

  const handlePaymentSuccess = useCallback(async () => {
    await refreshBillingData()
  }, [refreshBillingData])

  const handlePaymentFailure = useCallback(async () => {
    await loadPaymentHistory(page)
  }, [loadPaymentHistory, page])

  const handleCancelSubscription = useCallback(async () => {
    setIsCancelling(true)
    try {
      await cancelSubscription()
      toast.success(
        t(
          'settings.billingPage.cancelSuccess',
          'Subscription cancellation has been scheduled for period end.'
        )
      )
      setIsCancelDialogOpen(false)
      await refreshBillingData()
    } catch (error: any) {
      const backendMessage = error?.response?.data?.error
      const message =
        backendMessage ||
        t('settings.billingPage.cancelFailed', 'Failed to cancel subscription')
      toast.error(message)
    } finally {
      setIsCancelling(false)
    }
  }, [refreshBillingData, t])

  const paymentFlow = usePayment({
    onSuccess: handlePaymentSuccess,
    onFailure: handlePaymentFailure,
  })

  const loadPageData = useCallback(async () => {
    setIsLoading(true)
    await Promise.all([loadCurrentAndPlans(), loadPaymentHistory(1)])
    setPage(1)
    setIsLoading(false)
  }, [loadCurrentAndPlans, loadPaymentHistory])

  useEffect(() => {
    void loadPageData()
  }, [loadPageData])

  useEffect(() => {
    if (page === 1) return
    void loadPaymentHistory(page)
  }, [page, loadPaymentHistory])

  const isActiveSubscription = currentSubscription?.subscription_status === 'active'
  const isCancellationScheduled = Boolean(
    currentSubscription?.cancel_at_period_end || currentSubscription?.subscription_cancelled_at
  )
  const subscriptionExpiryText = formatDate(currentSubscription?.subscription_expires_at, locale)
  const cancellationScheduledText = t(
    'settings.billingPage.cancellationScheduled',
    'Subscription will expire on {date}.'
  ).replace('{date}', subscriptionExpiryText)
  const cancelConfirmDescriptionText = t(
    'settings.billingPage.cancelConfirmDescription',
    'Your subscription will remain active until {date}.'
  ).replace('{date}', subscriptionExpiryText)
  const currentPlanID = currentSubscription?.current_plan?.id || currentSubscription?.plan?.id || ''
  const selectedPlan = plans.find((plan) => plan.id === selectedPlanID) || plans[0] || null
  const plansToRender = selectedPlan ? [selectedPlan] : plans

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">
          {t('settings.billingPage.title', 'Billing & Subscription')}
        </h1>
        <p className="text-gray-500 dark:text-gray-400">
          {t('settings.billingPage.description', 'Manage your subscription plan and payment history.')}
        </p>
      </div>

      {errorMessage && (
        <div className="rounded-lg border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-900/40 dark:bg-red-900/20 dark:text-red-300">
          {errorMessage}
        </div>
      )}

      {paymentFlow.error && (
        <div className="rounded-lg border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-900/40 dark:bg-red-900/20 dark:text-red-300">
          {paymentFlow.error}
        </div>
      )}

      <section className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 p-6">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
          <div>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              {t('settings.billingPage.currentPlan', 'Current Subscription')}
            </h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
              {t('settings.billingPage.currentHint', 'Your current subscription access details')}
            </p>
          </div>
          <span className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-medium ${getSubscriptionPill(currentSubscription?.subscription_status)}`}>
            {(currentSubscription?.subscription_status || 'unknown').toUpperCase()}
          </span>
        </div>

        <div className="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
          <div className="rounded-lg bg-gray-50 dark:bg-gray-800 p-4">
            <p className="text-gray-500 dark:text-gray-400">{t('settings.billingPage.daysRemaining', 'Days Remaining')}</p>
            <p className="text-xl font-semibold text-gray-900 dark:text-white">
              {typeof currentSubscription?.days_remaining === 'number' ? currentSubscription.days_remaining : '-'}
            </p>
          </div>
          <div className="rounded-lg bg-gray-50 dark:bg-gray-800 p-4">
            <p className="text-gray-500 dark:text-gray-400">{t('settings.billingPage.expiresAt', 'Subscription Expires')}</p>
            <p className="text-xl font-semibold text-gray-900 dark:text-white">
              {subscriptionExpiryText}
            </p>
          </div>
        </div>

        {isCancellationScheduled && (
          <div className="mt-4 rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800 dark:border-amber-900/40 dark:bg-amber-900/20 dark:text-amber-200">
            {cancellationScheduledText}
          </div>
        )}

        {isActiveSubscription && !isCancellationScheduled && (
          <div className="mt-4 flex justify-end">
            <Button
              type="button"
              variant="outline"
              className="border-red-300 text-red-700 hover:bg-red-50 hover:text-red-800 dark:border-red-900/50 dark:text-red-300 dark:hover:bg-red-900/20"
              onClick={() => setIsCancelDialogOpen(true)}
            >
              {t('settings.billingPage.cancelSubscription', 'Cancel Subscription')}
            </Button>
          </div>
        )}
      </section>

      <section className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 p-6">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
          {t('settings.billingPage.planDetails', 'Available Plan')}
        </h2>
        {plans.length > 1 && (
          <div className="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-3">
            {plans.map((plan) => {
              const selected = selectedPlanID === plan.id
              return (
                <button
                  key={plan.id}
                  type="button"
                  onClick={() => setSelectedPlanID(plan.id)}
                  className={`rounded-lg border p-4 transition ${isRTL ? 'text-right' : 'text-left'} ${
                    selected
                      ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                      : 'border-gray-200 dark:border-gray-700 hover:border-primary-300'
                  }`}
                >
                  <p className="text-sm font-semibold text-gray-900 dark:text-white">{plan.name}</p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    {formatPlanPrice(plan, locale)} / {locale === 'ar' ? (isAnnualPlan(plan) ? 'سنة' : 'شهر') : plan.period}
                  </p>
                  {isAnnualPlan(plan) && (
                    <span className="inline-flex mt-2 rounded-full bg-emerald-100 px-2 py-0.5 text-[11px] font-medium text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300">
                      {t('settings.billingPage.saveAnnual', 'Save 17%')}
                    </span>
                  )}
                </button>
              )
            })}
          </div>
        )}
        <div className="mt-4 grid grid-cols-1 gap-4">
          {plansToRender.map((plan) => (
            <PlanCard
              key={plan.id}
              plan={plan}
              isActive={isActiveSubscription && currentPlanID === plan.id}
              isLoading={paymentFlow.isLoading}
              locale={locale}
              subscribeLabel={t('settings.billingPage.subscribeNow', 'Subscribe Now')}
              activeLabel={t('settings.billingPage.activePlan', 'Already Active')}
              annualSaveLabel={t('settings.billingPage.saveAnnual', 'Save 17%')}
              translateFeature={(feature) => {
                const featureMap: Record<string, string> = {
                  'Full POS System': t('settings.billingPage.features.fullPosSystem', 'Full POS System'),
                  'Kitchen Display System': t('settings.billingPage.features.kitchenDisplaySystem', 'Kitchen Display System'),
                  'Online Ordering Website': t('settings.billingPage.features.onlineOrderingWebsite', 'Online Ordering Website'),
                  'Analytics and Reports': t('settings.billingPage.features.analyticsAndReports', 'Analytics and Reports'),
                  'Multi-Restaurant Management': t('settings.billingPage.features.multiRestaurantManagement', 'Multi-Restaurant Management'),
                  'Mobile Apps': t('settings.billingPage.features.mobileApps', 'Mobile Apps'),
                  'Unlimited Staff': t('settings.billingPage.features.unlimitedStaff', 'Unlimited Staff'),
                  'Priority Support': t('settings.billingPage.features.prioritySupport', 'Priority Support'),
                  'Free Updates': t('settings.billingPage.features.freeUpdates', 'Free Updates'),
                  'Driver Management': t('settings.billingPage.features.driverManagement', 'Driver Management'),
                }
                return featureMap[feature] || feature
              }}
              onSubscribe={(planID) => {
                void paymentFlow.initiatePayment(planID)
              }}
            />
          ))}
        </div>
      </section>

      <section className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 p-6">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            {t('settings.billingPage.historyTitle', 'Payment History')}
          </h2>
          {isHistoryLoading && <Loader2 className="h-4 w-4 animate-spin text-gray-500" />}
        </div>

        <PaymentHistoryTable
          payments={history}
          total={totalPayments}
          page={page}
          pageSize={PAGE_SIZE}
          locale={locale}
          isRTL={isRTL}
          isLoading={isHistoryLoading}
          onPageChange={setPage}
          labels={{
            date: t('settings.billingPage.tableDate', 'Date'),
            amount: t('settings.billingPage.tableAmount', 'Amount'),
            status: t('settings.billingPage.tableStatus', 'Status'),
            voided: t('settings.billingPage.voided', 'Voided'),
            method: t('settings.billingPage.tableMethod', 'Method'),
            card: t('settings.billingPage.tableCard', 'Card'),
            reference: t('settings.billingPage.tableReference', 'Reference'),
            empty: t('settings.billingPage.noPayments', 'No payments yet'),
            pagination: t('settings.billingPage.pagination', 'Page'),
            previous: t('settings.billingPage.previous', 'Previous'),
            next: t('settings.billingPage.next', 'Next'),
            loading: t('settings.billingPage.loadingHistory', 'Loading payment history...'),
            cardFallback: t('settings.billingPage.cardFallback', 'Card'),
          }}
        />
      </section>

      <PaymentModal
        isOpen={paymentFlow.isModalOpen}
        iframeURL={paymentFlow.iframeURL}
        title={t('settings.billingPage.checkoutTitle', 'Complete Payment')}
        description={t('settings.billingPage.checkoutDescription', 'Finish checkout in the secure payment frame below.')}
        openNewTabLabel={t('settings.billingPage.openNewTab', 'Open in New Tab')}
        refreshLabel={t('settings.billingPage.refreshStatus', 'Refresh Status')}
        unavailableLabel={t('settings.billingPage.checkoutUnavailable', 'Checkout URL is unavailable')}
        onClose={() => {
          paymentFlow.closeModal()
          void refreshBillingData()
        }}
        onRefresh={() => refreshBillingData()}
        onComplete={(status) => {
          void paymentFlow.handlePaymentCompletion({ status, notify: true })
        }}
      />

      <Dialog open={isCancelDialogOpen} onOpenChange={setIsCancelDialogOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>
              {t('settings.billingPage.cancelConfirmTitle', 'Cancel subscription?')}
            </DialogTitle>
            <DialogDescription>
              {cancelConfirmDescriptionText}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setIsCancelDialogOpen(false)}
              disabled={isCancelling}
            >
              {t('common.cancel', 'Cancel')}
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={() => {
                void handleCancelSubscription()
              }}
              disabled={isCancelling}
            >
              {isCancelling
                ? t('settings.billingPage.cancelling', 'Cancelling...')
                : t('settings.billingPage.confirmCancel', 'Confirm Cancellation')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

