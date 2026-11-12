'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useMemo, useState } from 'react'
import { toast } from 'sonner'
import { ArrowLeft, Loader2, Lock, Moon, Save, Sun } from 'lucide-react'
import { useTheme } from 'next-themes'
import { useUserPermissions } from '@/hooks/useRbac'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { usePreferencesStore } from '@/stores/preferencesStore'
import {
  selectPOSAutoPrintEnabled,
  selectPOSFocusModeEnabled,
  selectPOSHideOutOfStockProducts,
  selectPOSProductCardSize,
  selectPOSSoundEnabled,
  usePOSStore,
} from '@/stores/posStore'

export default function POSSettingsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'
  const { resolvedTheme, setTheme } = useTheme()
  const updateTheme = usePreferencesStore((state) => state.updateTheme)

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canManagePOS = checkPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.WRITE)

  const soundEnabled = usePOSStore(selectPOSSoundEnabled)
  const autoPrintEnabled = usePOSStore(selectPOSAutoPrintEnabled)
  const focusModeEnabled = usePOSStore(selectPOSFocusModeEnabled)
  const hideOutOfStockProducts = usePOSStore(selectPOSHideOutOfStockProducts)
  const productCardSize = usePOSStore(selectPOSProductCardSize)
  const setSoundEnabled = usePOSStore((state) => state.setSoundEnabled)
  const setAutoPrintEnabled = usePOSStore((state) => state.setAutoPrintEnabled)
  const setFocusModeEnabled = usePOSStore((state) => state.setFocusModeEnabled)
  const setHideOutOfStockProducts = usePOSStore(
    (state) => state.setHideOutOfStockProducts
  )
  const setProductCardSize = usePOSStore((state) => state.setProductCardSize)

  const isDarkMode = resolvedTheme === 'dark'
  const [draftSoundEnabled, setDraftSoundEnabled] = useState(soundEnabled)
  const [draftAutoPrintEnabled, setDraftAutoPrintEnabled] = useState(autoPrintEnabled)
  const [draftFocusModeEnabled, setDraftFocusModeEnabled] = useState(focusModeEnabled)
  const [draftHideOutOfStockProducts, setDraftHideOutOfStockProducts] = useState(
    hideOutOfStockProducts
  )
  const [draftProductCardSize, setDraftProductCardSize] = useState(productCardSize)
  const [draftThemeMode, setDraftThemeMode] = useState<'light' | 'dark'>(isDarkMode ? 'dark' : 'light')

  const hasChanges = useMemo(
    () =>
      draftSoundEnabled !== soundEnabled ||
      draftAutoPrintEnabled !== autoPrintEnabled ||
      draftFocusModeEnabled !== focusModeEnabled ||
      draftHideOutOfStockProducts !== hideOutOfStockProducts ||
      draftProductCardSize !== productCardSize ||
      (draftThemeMode === 'dark') !== isDarkMode,
    [
      autoPrintEnabled,
      draftAutoPrintEnabled,
      draftFocusModeEnabled,
      draftHideOutOfStockProducts,
      draftProductCardSize,
      draftSoundEnabled,
      draftThemeMode,
      focusModeEnabled,
      hideOutOfStockProducts,
      isDarkMode,
      productCardSize,
      soundEnabled,
    ]
  )

  const applyTheme = (theme: 'light' | 'dark') => {
    setTheme(theme)
    updateTheme(theme)
  }

  const handleSave = () => {
    setSoundEnabled(draftSoundEnabled)
    setAutoPrintEnabled(draftAutoPrintEnabled)
    setFocusModeEnabled(draftFocusModeEnabled)
    setHideOutOfStockProducts(draftHideOutOfStockProducts)
    setProductCardSize(draftProductCardSize)
    applyTheme(draftThemeMode)
    toast.success(t('pos.settingsSaved', 'Settings saved'))
  }

  if (permissionLoading) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (!canManagePOS) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center px-4">
        <div className="w-full max-w-lg rounded-xl border border-red-400/40 bg-red-500/10 p-6">
          <div className={`flex items-start gap-3 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
            <Lock className="mt-0.5 h-5 w-5 flex-shrink-0 text-red-300" />
            <div>
              <h1 className="text-lg font-semibold text-red-100">
                {t('pos.accessDenied', 'Access denied')}
              </h1>
              <p className="mt-1 text-sm text-red-200">
                {t('pos.settingsAccessDenied', 'You do not have permission to manage POS settings.')}
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="mx-auto w-full max-w-3xl space-y-6">
      <div className={`flex flex-wrap items-center justify-between gap-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
        <div className={isRTL ? 'text-right' : ''}>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            {t('pos.settings', 'POS Settings')}
          </h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            {t('pos.settingsSubtitle', 'Configure cashier workspace behavior and display preferences')}
          </p>
        </div>

        <div className={`flex items-center gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
          <button
            type="button"
            onClick={handleSave}
            disabled={!hasChanges}
            className="inline-flex items-center gap-2 rounded-lg bg-primary px-3 py-2 text-sm font-medium text-white transition hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
          >
            <Save className="h-4 w-4" />
            {t('common.save', 'Save')}
          </button>
          <Link
            href={`/${locale}/dashboard/pos`}
            className="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200 dark:hover:bg-gray-800"
          >
            <ArrowLeft className="h-4 w-4" />
            {t('pos.backToPOS', 'Back to POS')}
          </Link>
        </div>
      </div>

      <section className="space-y-5 rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-800 dark:bg-gray-900">
        <div className={isRTL ? 'text-right' : ''}>
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            {t('pos.workspacePreferences', 'Workspace Preferences')}
          </h2>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            {t('pos.workspacePreferencesText', 'Tune POS to match cashier workflow and reduce distractions')}
          </p>
        </div>

        <label className={`flex items-center gap-3 rounded-lg border border-gray-200 bg-gray-50 px-3 py-3 dark:border-gray-700 dark:bg-gray-800 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
          <input
            type="checkbox"
            checked={draftFocusModeEnabled}
            onChange={(event) => {
              setDraftFocusModeEnabled(event.target.checked)
            }}
            className="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
          />
          <div>
            <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
              {t('pos.focusMode', 'Focus mode')}
            </p>
            <p className="text-xs text-gray-600 dark:text-gray-400">
              {t('pos.focusModeText', 'Hide non-essential overlays such as the POS performance badge')}
            </p>
          </div>
        </label>

        <label className={`flex items-center gap-3 rounded-lg border border-gray-200 bg-gray-50 px-3 py-3 dark:border-gray-700 dark:bg-gray-800 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
          <input
            type="checkbox"
            checked={draftHideOutOfStockProducts}
            onChange={(event) => {
              setDraftHideOutOfStockProducts(event.target.checked)
            }}
            className="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
          />
          <div>
            <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
              {t('pos.hideOutOfStockProducts', 'Hide out-of-stock products')}
            </p>
            <p className="text-xs text-gray-600 dark:text-gray-400">
              {t(
                'pos.hideOutOfStockProductsText',
                'Keep only available products visible for faster cashier flow'
              )}
            </p>
          </div>
        </label>

        <div className={`rounded-lg border border-gray-200 bg-gray-50 px-3 py-3 dark:border-gray-700 dark:bg-gray-800 ${isRTL ? 'text-right' : ''}`}>
          <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
            {t('pos.productCardSize', 'Product card size')}
          </p>
          <p className="mt-1 text-xs text-gray-600 dark:text-gray-400">
            {t('pos.productCardSizeText', 'Choose smaller or larger cards for your screen')}
          </p>
          <div className={`mt-3 flex gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
            <button
              type="button"
              onClick={() => setDraftProductCardSize('small')}
              className={`inline-flex min-h-10 items-center gap-2 rounded-md border px-3 py-2 text-sm font-medium transition ${
                draftProductCardSize === 'small'
                  ? 'border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-200'
                  : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200 dark:hover:bg-gray-800'
              }`}
            >
              {t('pos.productCardSizeSmall', 'Small')}
            </button>
            <button
              type="button"
              onClick={() => setDraftProductCardSize('large')}
              className={`inline-flex min-h-10 items-center gap-2 rounded-md border px-3 py-2 text-sm font-medium transition ${
                draftProductCardSize === 'large'
                  ? 'border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-200'
                  : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200 dark:hover:bg-gray-800'
              }`}
            >
              {t('pos.productCardSizeLarge', 'Large')}
            </button>
          </div>
        </div>

        <label className={`flex items-center gap-3 rounded-lg border border-gray-200 bg-gray-50 px-3 py-3 dark:border-gray-700 dark:bg-gray-800 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
          <input
            type="checkbox"
            checked={draftSoundEnabled}
            onChange={(event) => {
              setDraftSoundEnabled(event.target.checked)
            }}
            className="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
          />
          <div>
            <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
              {t('pos.soundAlerts', 'Sound alerts')}
            </p>
            <p className="text-xs text-gray-600 dark:text-gray-400">
              {t('pos.soundAlertsText', 'Play sounds for add/remove, success, and error actions')}
            </p>
          </div>
        </label>

        <label className={`flex items-center gap-3 rounded-lg border border-gray-200 bg-gray-50 px-3 py-3 dark:border-gray-700 dark:bg-gray-800 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
          <input
            type="checkbox"
            checked={draftAutoPrintEnabled}
            onChange={(event) => {
              setDraftAutoPrintEnabled(event.target.checked)
            }}
            className="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
          />
          <div>
            <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
              {t('pos.autoPrint', 'Auto print receipts')}
            </p>
            <p className="text-xs text-gray-600 dark:text-gray-400">
              {t('pos.autoPrintText', 'Automatically print receipt after successful order placement')}
            </p>
          </div>
        </label>

        <div className={`rounded-lg border border-gray-200 bg-gray-50 px-3 py-3 dark:border-gray-700 dark:bg-gray-800 ${isRTL ? 'text-right' : ''}`}>
          <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
            {t('pos.themeMode', 'Theme mode')}
          </p>
          <p className="mt-1 text-xs text-gray-600 dark:text-gray-400">
            {t('pos.themeModeText', 'Switch cashier interface between light and dark themes')}
          </p>
          <div className={`mt-3 flex gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
            <button
              type="button"
              onClick={() => setDraftThemeMode('light')}
              className={`inline-flex min-h-10 items-center gap-2 rounded-md border px-3 py-2 text-sm font-medium transition ${
                draftThemeMode === 'light'
                  ? 'border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-200'
                  : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200 dark:hover:bg-gray-800'
              }`}
            >
              <Sun className="h-4 w-4" />
              {t('common.lightMode', 'Light Mode')}
            </button>
            <button
              type="button"
              onClick={() => setDraftThemeMode('dark')}
              className={`inline-flex min-h-10 items-center gap-2 rounded-md border px-3 py-2 text-sm font-medium transition ${
                draftThemeMode === 'dark'
                  ? 'border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-200'
                  : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200 dark:hover:bg-gray-800'
              }`}
            >
              <Moon className="h-4 w-4" />
              {t('common.darkMode', 'Dark Mode')}
            </button>
          </div>
        </div>
      </section>
    </div>
  )
}
