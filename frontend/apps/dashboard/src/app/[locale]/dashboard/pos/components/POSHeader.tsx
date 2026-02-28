'use client'

import { useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import {
  Clock3,
  Eye,
  EyeOff,
  LogOut,
  Maximize2,
  Minimize2,
  Moon,
  Printer,
  Settings2,
  Sun,
  Volume2,
  VolumeX,
  Wifi,
  WifiOff,
} from 'lucide-react'
import { useTheme } from 'next-themes'
import { useAuthStore } from '@/stores/authStore'
import { usePreferencesStore } from '@/stores/preferencesStore'
import {
  selectPOSFocusModeEnabled,
  selectPOSHeldOrdersCount,
  selectPOSSoundEnabled,
  usePOSStore,
} from '@/stores/posStore'

interface POSHeaderProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  isOffline: boolean
  queuedOrderCount: number
  isReplayingQueue: boolean
  autoPrintEnabled: boolean
  onToggleAutoPrint: () => void
}

const formatClock = (date: Date, locale: 'en' | 'ar') =>
  date.toLocaleTimeString(locale === 'ar' ? 'ar-EG' : 'en-US', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })

export function POSHeader({
  locale,
  t,
  isOffline,
  queuedOrderCount,
  isReplayingQueue,
  autoPrintEnabled,
  onToggleAutoPrint,
}: POSHeaderProps) {
  const router = useRouter()
  const user = useAuthStore((state) => state.user)
  const heldOrdersCount = usePOSStore(selectPOSHeldOrdersCount)
  const soundEnabled = usePOSStore(selectPOSSoundEnabled)
  const focusModeEnabled = usePOSStore(selectPOSFocusModeEnabled)
  const toggleSound = usePOSStore((state) => state.toggleSound)
  const toggleFocusMode = usePOSStore((state) => state.toggleFocusMode)
  const { resolvedTheme, setTheme } = useTheme()
  const updateTheme = usePreferencesStore((state) => state.updateTheme)
  const [now, setNow] = useState<Date>(() => new Date())
  const [isFullscreen, setIsFullscreen] = useState(false)

  useEffect(() => {
    const timer = window.setInterval(() => setNow(new Date()), 1000)
    return () => {
      window.clearInterval(timer)
    }
  }, [])

  useEffect(() => {
    const syncFullscreen = () => {
      setIsFullscreen(Boolean(document.fullscreenElement))
    }
    syncFullscreen()
    document.addEventListener('fullscreenchange', syncFullscreen)
    return () => {
      document.removeEventListener('fullscreenchange', syncFullscreen)
    }
  }, [])

  const isRTL = locale === 'ar'
  const isDarkMode = resolvedTheme === 'dark'

  const onlineStatusText = useMemo(() => {
    if (isOffline) {
      return t('pos.offline', 'Offline')
    }
    if (isReplayingQueue) {
      return t('pos.syncingQueued', 'Syncing queued orders...')
    }
    return t('pos.online', 'Online')
  }, [isOffline, isReplayingQueue, t])

  const handleFullscreenToggle = async () => {
    if (!document.fullscreenEnabled) {
      return
    }

    try {
      if (!document.fullscreenElement) {
        await document.documentElement.requestFullscreen()
      } else {
        await document.exitFullscreen()
      }
    } catch {
      setIsFullscreen(Boolean(document.fullscreenElement))
    }
  }

  const handleThemeToggle = () => {
    const nextTheme = isDarkMode ? 'light' : 'dark'
    setTheme(nextTheme)
    updateTheme(nextTheme)
  }

  return (
    <header className="flex items-center justify-between gap-3 rounded-lg bg-gray-900 px-4 py-3 text-white">
      <div className={`flex items-center gap-4 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
        <div className={`flex items-center gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
          <Clock3 className="h-4 w-4 text-gray-300" />
          <span className="text-sm font-semibold tabular-nums">{formatClock(now, locale)}</span>
        </div>
        <div className={`hidden items-center gap-2 text-xs text-gray-300 sm:flex ${isRTL ? 'flex-row-reverse' : ''}`}>
          <span>{t('pos.cashier', 'Cashier')}:</span>
          <span className="font-medium text-white">{user?.name || '—'}</span>
        </div>
      </div>

      <div className={`flex items-center gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
        <div
          className={`inline-flex items-center gap-1 rounded-full px-3 py-1 text-xs font-semibold ${
            isOffline
              ? 'bg-red-500/20 text-red-200'
              : 'bg-emerald-500/20 text-emerald-200'
          }`}
        >
          {isOffline ? <WifiOff className="h-3.5 w-3.5" /> : <Wifi className="h-3.5 w-3.5" />}
          <span>{onlineStatusText}</span>
        </div>

        {heldOrdersCount > 0 && (
          <div className="rounded-full bg-amber-500/20 px-3 py-1 text-xs font-semibold text-amber-200">
            {heldOrdersCount} {t('pos.heldOrders', 'Held')}
          </div>
        )}

        {queuedOrderCount > 0 && (
          <div className="rounded-full bg-sky-500/20 px-3 py-1 text-xs font-semibold text-sky-200">
            {queuedOrderCount} {t('pos.queuedOrders', 'Queued')}
          </div>
        )}

        <button
          type="button"
          onClick={toggleSound}
          className="inline-flex h-9 w-9 items-center justify-center rounded-md border border-gray-700 bg-gray-800 transition hover:bg-gray-700"
          title={soundEnabled ? t('pos.soundOn', 'Sound on') : t('pos.soundOff', 'Sound off')}
        >
          {soundEnabled ? <Volume2 className="h-4 w-4" /> : <VolumeX className="h-4 w-4" />}
        </button>

        <button
          type="button"
          onClick={onToggleAutoPrint}
          className={`inline-flex h-9 w-9 items-center justify-center rounded-md border transition ${
            autoPrintEnabled
              ? 'border-emerald-500/70 bg-emerald-500/20 text-emerald-100 hover:bg-emerald-500/30'
              : 'border-gray-700 bg-gray-800 hover:bg-gray-700'
          }`}
          title={
            autoPrintEnabled
              ? t('pos.autoPrintOn', 'Auto print on')
              : t('pos.autoPrintOff', 'Auto print off')
          }
        >
          <Printer className="h-4 w-4" />
        </button>

        <button
          type="button"
          onClick={() => {
            void handleFullscreenToggle()
          }}
          className="inline-flex h-9 w-9 items-center justify-center rounded-md border border-gray-700 bg-gray-800 transition hover:bg-gray-700"
          title={isFullscreen ? t('pos.exitFullscreen', 'Exit fullscreen') : t('pos.fullscreen', 'Fullscreen')}
        >
          {isFullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
        </button>

        <button
          type="button"
          onClick={handleThemeToggle}
          className="inline-flex h-9 w-9 items-center justify-center rounded-md border border-gray-700 bg-gray-800 transition hover:bg-gray-700"
          title={isDarkMode ? t('common.lightMode', 'Switch to light mode') : t('common.darkMode', 'Switch to dark mode')}
        >
          {isDarkMode ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        </button>

        <button
          type="button"
          onClick={toggleFocusMode}
          className={`inline-flex h-9 w-9 items-center justify-center rounded-md border transition ${
            focusModeEnabled
              ? 'border-primary-500/70 bg-primary-500/20 text-primary-100 hover:bg-primary-500/30'
              : 'border-gray-700 bg-gray-800 hover:bg-gray-700'
          }`}
          title={
            focusModeEnabled
              ? t('pos.focusModeOn', 'Focus mode enabled')
              : t('pos.focusModeOff', 'Focus mode disabled')
          }
        >
          {focusModeEnabled ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
        </button>

        <Link
          href={`/${locale}/dashboard/pos/settings`}
          className="inline-flex h-9 w-9 items-center justify-center rounded-md border border-gray-700 bg-gray-800 transition hover:bg-gray-700"
          title={t('pos.settings', 'POS settings')}
        >
          <Settings2 className="h-4 w-4" />
        </Link>

        <button
          type="button"
          onClick={() => {
            router.push(`/${locale}/dashboard`)
          }}
          className="inline-flex h-9 items-center gap-2 rounded-md border border-gray-700 bg-gray-800 px-3 text-xs font-semibold transition hover:bg-gray-700"
        >
          <LogOut className="h-4 w-4" />
          <span className="hidden sm:inline">{t('pos.exitPOS', 'Exit POS')}</span>
        </button>
      </div>
    </header>
  )
}

export default POSHeader
