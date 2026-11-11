import Link from 'next/link'
import { headers } from 'next/headers'
import { AlertCircle } from 'lucide-react'
import { HeaderSSR } from '@/components/layout/header-ssr'
import { FooterSSR } from '@/components/layout/footer-ssr'
import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server'
import type { ServerThemeData } from '@/lib/api/get-theme-server'
import { getServerApiBaseUrl } from '@/lib/config/api'
import { createTranslator } from '@/lib/translations'
import { getTableStatusNotice } from '@/lib/table-ordering'
import { MenuPageContent } from '../../../menu/menu-page-content'

interface TableMenuPageProps {
  params: Promise<{
    locale: 'en' | 'ar'
    tableId: string
  }>
}

interface TablePublicInfo {
  id: number
  table_number: number
  name?: string
  status?: string
  restaurant_slug?: string
}

interface TableInfoResponse {
  success?: boolean
  data?: TablePublicInfo
}

async function getPublicTableInfo(tableId: number): Promise<TablePublicInfo | null> {
  try {
    const apiBase = getServerApiBaseUrl()
    const response = await fetch(`${apiBase}/public/table/${tableId}/info`, {
      cache: 'no-store',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    if (!response.ok) {
      return null
    }

    const payload = (await response.json()) as TableInfoResponse
    if (!payload?.data?.id) {
      return null
    }

    return payload.data
  } catch {
    return null
  }
}

function InvalidTableState({
  locale,
  t,
  themeData,
  tableId,
}: {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  themeData: ServerThemeData
  tableId?: number
}) {
  return (
    <div className="flex flex-col min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
      <HeaderSSR themeData={themeData} locale={locale} cartItemsCount={0} />
      <main className="flex-1 w-full flex items-center justify-center px-4 py-8">
        <div className="max-w-md w-full rounded-xl border border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/20 p-6 text-center shadow-sm">
          <AlertCircle className="h-10 w-10 text-red-600 dark:text-red-400 mx-auto mb-3" />
          <h1 className="text-lg font-semibold text-red-900 dark:text-red-200">
            {t('tableOrdering.invalidTableTitle')}
          </h1>
          <p className="mt-2 text-sm text-red-800 dark:text-red-300">
            {t('tableOrdering.invalidTable')}
          </p>
          {tableId ? (
            <p className="mt-3 text-xs text-red-700/80 dark:text-red-200/80">
              {t('tableOrdering.table')} #{tableId}
            </p>
          ) : null}
          <div className="mt-6">
            <Link
              href={`/${locale}/menu`}
              className="inline-flex items-center justify-center rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
            >
              {t('tableOrdering.openMenu')}
            </Link>
          </div>
        </div>
      </main>
      <FooterSSR themeData={themeData} locale={locale} />
    </div>
  )
}

export default async function TableMenuPage({ params }: TableMenuPageProps) {
  const resolvedParams = await params
  const locale = resolvedParams.locale === 'ar' ? 'ar' : 'en'
  const tableId = Number.parseInt(resolvedParams.tableId, 10)
  const t = createTranslator(locale)

  const themeData = (await getThemeDataServer()) ?? getDefaultThemeData()

  if (!Number.isFinite(tableId) || tableId <= 0) {
    return <InvalidTableState locale={locale} t={t} themeData={themeData} />
  }

  const tableInfo = await getPublicTableInfo(tableId)
  const headersList = await headers()
  const fallbackSlug = headersList.get('x-restaurant-slug') || 'demo'

  if (!tableInfo) {
    return (
      <InvalidTableState
        locale={locale}
        t={t}
        themeData={themeData}
        tableId={tableId}
      />
    )
  }

  const restaurantSlug = tableInfo.restaurant_slug || fallbackSlug
  const statusNotice = getTableStatusNotice(tableInfo.status)

  return (
    <div className="flex flex-col min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
      <HeaderSSR themeData={themeData} locale={locale} cartItemsCount={0} />
      <main className="flex-1 w-full">
        <MenuPageContent
          locale={locale}
          themeData={themeData}
          restaurantSlug={restaurantSlug}
          tableContext={{
            tableId: tableInfo.id,
            tableNumber: tableInfo.table_number,
            tableName: tableInfo.name || null,
            tableStatus: tableInfo.status,
          }}
          tableStatusNotice={statusNotice}
        />
      </main>
      <FooterSSR themeData={themeData} locale={locale} />
    </div>
  )
}
