'use client'

import { useEffect, useMemo, useState } from 'react'
import { Loader2, Lock } from 'lucide-react'
import { PageHeaderSkeleton } from '@/components/skeletons'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useUserPermissions } from '@/hooks/useRbac'
import { getAccounts } from '@/lib/api/accountingApi'
import { debugError } from '@/lib/debug'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { usePathname } from 'next/navigation'
import dynamic from 'next/dynamic'
import type { Account } from '@/types/accounting'

const LedgerView = dynamic(
  () => import('@/components/accounting/LedgerView').then((mod) => mod.LedgerView),
  { ssr: false, loading: () => <PageHeaderSkeleton /> }
)
const TrialBalanceTable = dynamic(
  () => import('@/components/accounting/TrialBalanceTable').then((mod) => mod.TrialBalanceTable),
  { ssr: false, loading: () => <PageHeaderSkeleton /> }
)

export default function GeneralLedgerPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])
  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canView = checkPermission(MODULE_IDS.ACCOUNTING, PERMISSION_LEVELS.READ)

  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (permissionLoading || !canView) {
      return
    }

    const loadAccounts = async () => {
      setLoading(true)
      try {
        const data = await getAccounts({ tree: true, include_inactive: false })
        setAccounts(data)
      } catch (error) {
        debugError('[ACCOUNTING] load ledger accounts failed', error)
      } finally {
        setLoading(false)
      }
    }

    void loadAccounts()
  }, [canView, permissionLoading])

  if (permissionLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (!canView) {
    return (
      <div className="rounded-xl border border-red-200 bg-red-50 p-6 dark:border-red-900 dark:bg-red-950/30">
        <div className="flex items-start gap-3">
          <Lock className="mt-0.5 h-5 w-5 text-red-600 dark:text-red-400" />
          <div>
            <h2 className="font-semibold text-red-900 dark:text-red-300">
              {t('accounting.accessDeniedTitle', 'Access denied')}
            </h2>
            <p className="mt-1 text-sm text-red-700 dark:text-red-300">
              {t(
                'accounting.accessDeniedText',
                'You do not have permission to view accounting records.'
              )}
            </p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {loading ? <PageHeaderSkeleton /> : null}

      <div className="space-y-1">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
          {t('accounting.generalLedger', 'General Ledger')}
        </h1>
        <p className="text-gray-500 dark:text-gray-400">
          {t(
            'accounting.generalLedgerSubtitle',
            'Review posted accounting activity by account and validate the trial balance.'
          )}
        </p>
      </div>

      <Tabs defaultValue="ledger" className="space-y-6">
        <TabsList>
          <TabsTrigger value="ledger">
            {t('accounting.ledgerTab', 'Ledger')}
          </TabsTrigger>
          <TabsTrigger value="trial-balance">
            {t('accounting.trialBalanceTab', 'Trial Balance')}
          </TabsTrigger>
        </TabsList>

        <TabsContent value="ledger">
          <LedgerView accounts={accounts} locale={locale} />
        </TabsContent>

        <TabsContent value="trial-balance">
          <TrialBalanceTable locale={locale} />
        </TabsContent>
      </Tabs>
    </div>
  )
}


