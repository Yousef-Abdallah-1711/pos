'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import { AxiosError } from 'axios'
import { useKDSOffline } from '@/hooks/useKDSOffline'
import { api } from '@/lib/api'
import { startPerformanceTimer } from '@/lib/performance'
import {
  KDSFilter,
  KDSOrder,
  KDSStats,
  useKDSStore,
} from '@/stores/kdsStore'

export type KDSDisplayDensity = 'compact' | 'comfortable' | 'spacious'

export interface KDSSettings {
  id: number
  tenant_id: number
  restaurant_id: number
  auto_bump_minutes: number
  alert_sound_enabled: boolean
  show_analytics_panel: boolean
  display_density: KDSDisplayDensity
  color_warning_minutes: number
  color_danger_minutes: number
  created_at?: string
  updated_at?: string
}

const DEFAULT_STATS: KDSStats = {
  active_orders: 0,
  avg_prep_minutes: 0,
  completed_today: 0,
  new_orders: 0,
  in_progress: 0,
  ready: 0,
}
const DEFAULT_SETTINGS: KDSSettings = {
  id: 0,
  tenant_id: 0,
  restaurant_id: 0,
  auto_bump_minutes: 0,
  alert_sound_enabled: true,
  show_analytics_panel: false,
  display_density: 'comfortable',
  color_warning_minutes: 10,
  color_danger_minutes: 20,
}
const POLL_INTERVAL_MS = 10_000
const POLL_FALLBACK_DELAY_MS = 5_000
const RETRY_ATTEMPTS = 2
const RETRY_BASE_DELAY_MS = 300
const FRONTEND_BUMP_DEBOUNCE_MS = 200
const KDS_ANALYTICS_PANEL_PREF_KEY = 'kds-show-analytics-panel'
const INITIAL_REQUEST_TIMEOUT_MS = 2_500
const STANDARD_REQUEST_TIMEOUT_MS = 5_000
const KDS_SNAPSHOT_KEY = 'kds-bootstrap-snapshot-v1'
const KDS_SNAPSHOT_MAX_AGE_MS = 5 * 60 * 1000

interface ApiResponse<T> {
  success?: boolean
  data?: T
  error?: string
  message?: string
}

interface FetchConfig {
  timeoutMs?: number
  attempts?: number
}

interface KDSSnapshot {
  orders: KDSOrder[]
  stats: KDSStats
  settings: KDSSettings
  savedAt: number
}

const extractData = <T>(payload: unknown, fallback: T): T => {
  if (payload && typeof payload === 'object' && 'data' in payload) {
    return ((payload as ApiResponse<T>).data ?? fallback) as T
  }
  return (payload as T) ?? fallback
}

const getErrorMessage = (error: unknown, fallback: string) => {
  if (!error || typeof error !== 'object') {
    return fallback
  }

  const axiosError = error as AxiosError<ApiResponse<unknown>>
  const responseError =
    axiosError.response?.data?.error ||
    axiosError.response?.data?.message

  return responseError || axiosError.message || fallback
}

const isTransientError = (error: unknown) => {
  if (!error || typeof error !== 'object') {
    return false
  }

  const axiosError = error as AxiosError<ApiResponse<unknown>>
  const statusCode = axiosError.response?.status
  if (typeof statusCode === 'number') {
    return statusCode >= 500 || statusCode === 429 || statusCode === 425 || statusCode === 408
  }

  return true
}

const sleep = (ms: number) =>
  new Promise<void>((resolve) => {
    setTimeout(resolve, ms)
  })

const readAnalyticsPanelPreference = () => {
  if (typeof window === 'undefined') {
    return null
  }
  try {
    const raw = window.localStorage.getItem(KDS_ANALYTICS_PANEL_PREF_KEY)
    if (raw === '1') {
      return true
    }
    if (raw === '0') {
      return false
    }
  } catch {
    // Ignore client storage access errors.
  }
  return null
}

const writeAnalyticsPanelPreference = (value: boolean) => {
  if (typeof window === 'undefined') {
    return
  }
  try {
    window.localStorage.setItem(KDS_ANALYTICS_PANEL_PREF_KEY, value ? '1' : '0')
  } catch {
    // Ignore client storage access errors.
  }
}

const resolveDefaultSettings = (): KDSSettings => {
  const localPreference = readAnalyticsPanelPreference()
  if (localPreference === null) {
    return DEFAULT_SETTINGS
  }

  return {
    ...DEFAULT_SETTINGS,
    show_analytics_panel: localPreference,
  }
}

const readKDSSnapshot = (): KDSSnapshot | null => {
  if (typeof window === 'undefined') {
    return null
  }

  try {
    const raw = window.localStorage.getItem(KDS_SNAPSHOT_KEY)
    if (!raw) {
      return null
    }

    const snapshot = JSON.parse(raw) as KDSSnapshot
    if (!snapshot || typeof snapshot !== 'object') {
      return null
    }

    if (!Array.isArray(snapshot.orders) || !snapshot.stats || !snapshot.settings) {
      return null
    }

    if (typeof snapshot.savedAt !== 'number') {
      return null
    }

    if (Date.now() - snapshot.savedAt > KDS_SNAPSHOT_MAX_AGE_MS) {
      window.localStorage.removeItem(KDS_SNAPSHOT_KEY)
      return null
    }

    return snapshot
  } catch {
    try {
      window.localStorage.removeItem(KDS_SNAPSHOT_KEY)
    } catch {
      // Ignore storage cleanup errors.
    }
    return null
  }
}

const writeKDSSnapshot = (snapshot: Omit<KDSSnapshot, 'savedAt'>) => {
  if (typeof window === 'undefined') {
    return
  }

  try {
    window.localStorage.setItem(
      KDS_SNAPSHOT_KEY,
      JSON.stringify({
        ...snapshot,
        savedAt: Date.now(),
      } satisfies KDSSnapshot)
    )
  } catch {
    // Ignore storage write errors.
  }
}

const withRetry = async <T>(operation: () => Promise<T>, attempts = RETRY_ATTEMPTS): Promise<T> => {
  let lastError: unknown

  for (let attempt = 0; attempt < attempts; attempt += 1) {
    try {
      return await operation()
    } catch (error) {
      lastError = error
      if (attempt === attempts - 1 || !isTransientError(error)) {
        throw error
      }

      const delay = RETRY_BASE_DELAY_MS * Math.pow(2, attempt)
      await sleep(delay)
    }
  }

  throw lastError
}

const fetchQueue = async (config: FetchConfig = {}): Promise<KDSOrder[]> => {
  const response = await withRetry(
    () => api.get('/kds/queue', { timeout: config.timeoutMs ?? STANDARD_REQUEST_TIMEOUT_MS }),
    config.attempts ?? RETRY_ATTEMPTS
  )
  const orders = extractData<KDSOrder[]>(response.data, [])
  return Array.isArray(orders) ? orders : []
}

const fetchStats = async (config: FetchConfig = {}): Promise<KDSStats> => {
  const response = await withRetry(
    () => api.get('/kds/stats', { timeout: config.timeoutMs ?? STANDARD_REQUEST_TIMEOUT_MS }),
    config.attempts ?? RETRY_ATTEMPTS
  )
  return extractData<KDSStats>(response.data, DEFAULT_STATS)
}

const fetchSettings = async (config: FetchConfig = {}): Promise<KDSSettings> => {
  const response = await withRetry(
    () => api.get('/kds/settings', { timeout: config.timeoutMs ?? STANDARD_REQUEST_TIMEOUT_MS }),
    config.attempts ?? RETRY_ATTEMPTS
  )
  const settings = extractData<KDSSettings>(response.data, DEFAULT_SETTINGS)
  const normalized = {
    ...DEFAULT_SETTINGS,
    ...settings,
  }

  const localPreference = readAnalyticsPanelPreference()
  if (localPreference !== null) {
    // Keep user-selected dashboard behavior stable even when backend settings drift
    // or when multiple backend versions are running in development.
    normalized.show_analytics_panel = localPreference
  }

  writeAnalyticsPanelPreference(normalized.show_analytics_panel)
  return normalized
}

export const useKDS = () => {
  const orders = useKDSStore((state) => state.orders)
  const stats = useKDSStore((state) => state.stats)
  const filter = useKDSStore((state) => state.filter)
  const soundEnabled = useKDSStore((state) => state.soundEnabled)
  const isFullscreen = useKDSStore((state) => state.isFullscreen)
  const isConnected = useKDSStore((state) => state.isConnected)

  const setOrders = useKDSStore((state) => state.setOrders)
  const upsertOrder = useKDSStore((state) => state.upsertOrder)
  const setStats = useKDSStore((state) => state.setStats)
  const setFilter = useKDSStore((state) => state.setFilter)
  const setSoundEnabled = useKDSStore((state) => state.setSoundEnabled)
  const toggleSound = useKDSStore((state) => state.toggleSound)
  const setFullscreen = useKDSStore((state) => state.setFullscreen)
  const toggleFullscreen = useKDSStore((state) => state.toggleFullscreen)
  const setConnected = useKDSStore((state) => state.setConnected)
  const optimisticBump = useKDSStore((state) => state.optimisticBump)
  const optimisticRecall = useKDSStore((state) => state.optimisticRecall)
  const rollbackOrder = useKDSStore((state) => state.rollbackOrder)
  const {
    isOffline,
    queuedActionsCount,
    isReplayingQueue,
    enqueueAction,
    replayQueuedActions,
  } = useKDSOffline()

  const [loading, setLoading] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [pendingOrderIds, setPendingOrderIds] = useState<number[]>([])
  const [settings, setSettings] = useState<KDSSettings>(() => resolveDefaultSettings())
  const mountedRef = useRef(true)
  const settingsHydratedRef = useRef(false)
  const lastBumpAttemptRef = useRef(0)
  const settingsRef = useRef<KDSSettings>(settings)
  const initialLoadMeasuredRef = useRef(false)
  const stopKDSOpenTimerRef = useRef<(() => void) | null>(null)

  if (stopKDSOpenTimerRef.current === null && typeof window !== 'undefined') {
    stopKDSOpenTimerRef.current = startPerformanceTimer('kds.open')
  }

  const markInitialLoadMeasured = useCallback(() => {
    if (initialLoadMeasuredRef.current) {
      return
    }
    initialLoadMeasuredRef.current = true
    stopKDSOpenTimerRef.current?.()
    stopKDSOpenTimerRef.current = null
  }, [])

  const applySettings = useCallback((nextSettings: KDSSettings) => {
    settingsRef.current = nextSettings
    setSettings(nextSettings)
    if (!settingsHydratedRef.current) {
      setSoundEnabled(nextSettings.alert_sound_enabled)
      settingsHydratedRef.current = true
    }
  }, [setSoundEnabled])

  const persistSnapshot = useCallback((
    nextOrders?: KDSOrder[],
    nextStats?: KDSStats,
    nextSettings: KDSSettings = settingsRef.current
  ) => {
    const currentState = useKDSStore.getState()
    writeKDSSnapshot({
      orders: nextOrders ?? currentState.orders,
      stats: nextStats ?? currentState.stats ?? DEFAULT_STATS,
      settings: nextSettings,
    })
  }, [])

  useEffect(() => {
    mountedRef.current = true
    return () => {
      mountedRef.current = false
    }
  }, [])

  useEffect(() => {
    let active = true

    const loadInitialData = async () => {
      const snapshot = readKDSSnapshot()
      if (snapshot && active && mountedRef.current) {
        setOrders(snapshot.orders)
        setStats(snapshot.stats)
        applySettings({
          ...DEFAULT_SETTINGS,
          ...snapshot.settings,
        })
        setError(null)
        setLoading(false)
        markInitialLoadMeasured()
      }

      if (isOffline) {
        if (active && mountedRef.current) {
          setLoading(false)
          markInitialLoadMeasured()
        }
        return
      }

      try {
        if (!snapshot) {
          setLoading(true)
        }
        const [queue, queueStats] = await Promise.all([
          fetchQueue({
            timeoutMs: INITIAL_REQUEST_TIMEOUT_MS,
            attempts: RETRY_ATTEMPTS,
          }),
          fetchStats({
            timeoutMs: INITIAL_REQUEST_TIMEOUT_MS,
            attempts: RETRY_ATTEMPTS,
          }),
        ])
        if (!active || !mountedRef.current) {
          return
        }
        setOrders(queue)
        setStats(queueStats)
        persistSnapshot(queue, queueStats)
        setError(null)
        markInitialLoadMeasured()
      } catch (err) {
        if (!active || !mountedRef.current) {
          return
        }

        if (snapshot) {
          setError('Live KDS data is temporarily unavailable. Showing last synced queue.')
          markInitialLoadMeasured()
        } else {
          setError(getErrorMessage(err, 'Failed to load kitchen display data'))
        }
      } finally {
        if (active && mountedRef.current) {
          setLoading(false)
          markInitialLoadMeasured()
        }
      }

      // Settings are non-critical for initial render; load in the background.
      try {
        const queueSettings = await fetchSettings({
          timeoutMs: STANDARD_REQUEST_TIMEOUT_MS,
          attempts: RETRY_ATTEMPTS,
        })
        if (!active || !mountedRef.current) {
          return
        }
        applySettings(queueSettings)
        persistSnapshot(undefined, undefined, queueSettings)
      } catch {
        // Keep queue visible even if settings fetch fails.
      }
    }

    void loadInitialData()
    return () => {
      active = false
    }
  }, [applySettings, isOffline, markInitialLoadMeasured, persistSnapshot, setOrders, setStats])

  useEffect(() => {
    if (isConnected || isOffline) {
      return
    }

    let pollId: number | null = null
    const runPoll = async () => {
      try {
        const [queue, queueStats] = await Promise.all([
          fetchQueue({ timeoutMs: STANDARD_REQUEST_TIMEOUT_MS, attempts: 1 }),
          fetchStats({ timeoutMs: STANDARD_REQUEST_TIMEOUT_MS, attempts: 1 }),
        ])
        if (!mountedRef.current) {
          return
        }
        setOrders(queue)
        setStats(queueStats)
        persistSnapshot(queue, queueStats)
        setError(null)
      } catch (err) {
        if (!mountedRef.current) {
          return
        }
        setError(getErrorMessage(err, 'Failed to sync kitchen display data'))
      }
    }

    const startPollTimeoutId = window.setTimeout(() => {
      void runPoll()
      pollId = window.setInterval(() => {
        void runPoll()
      }, POLL_INTERVAL_MS)
    }, POLL_FALLBACK_DELAY_MS)

    return () => {
      window.clearTimeout(startPollTimeoutId)
      if (pollId !== null) {
        window.clearInterval(pollId)
      }
    }
  }, [isConnected, isOffline, persistSnapshot, setOrders, setStats])

  useEffect(() => {
    if (isOffline) {
      return
    }

    let active = true

    const replayOfflineQueue = async () => {
      const replayResult = await replayQueuedActions(async (action) => {
        if (action.type === 'bump') {
          await withRetry(() => api.put(`/kds/orders/${action.orderId}/bump`))
          return
        }

        await withRetry(() => api.put(`/kds/orders/${action.orderId}/recall`))
      })

      if (!active || !mountedRef.current) {
        return
      }

      if (replayResult.processed === 0) {
        return
      }

      try {
        setIsRefreshing(true)
        const [queue, queueStats, queueSettings] = await Promise.all([
          fetchQueue({ timeoutMs: STANDARD_REQUEST_TIMEOUT_MS, attempts: RETRY_ATTEMPTS }),
          fetchStats({ timeoutMs: STANDARD_REQUEST_TIMEOUT_MS, attempts: RETRY_ATTEMPTS }),
          fetchSettings({ timeoutMs: STANDARD_REQUEST_TIMEOUT_MS, attempts: RETRY_ATTEMPTS }),
        ])
        if (!active || !mountedRef.current) {
          return
        }
        setOrders(queue)
        setStats(queueStats)
        applySettings(queueSettings)
        persistSnapshot(queue, queueStats, queueSettings)
        setError(null)
      } catch (err) {
        if (!active || !mountedRef.current) {
          return
        }
        setError(getErrorMessage(err, 'Failed to sync queued offline actions'))
      } finally {
        if (active && mountedRef.current) {
          setIsRefreshing(false)
        }
      }
    }

    void replayOfflineQueue()

    return () => {
      active = false
    }
  }, [applySettings, isOffline, persistSnapshot, replayQueuedActions, setOrders, setStats])

  useEffect(() => {
    if (loading || !stats) {
      return
    }
    persistSnapshot(orders, stats, settingsRef.current)
  }, [loading, orders, persistSnapshot, stats])

  const filteredOrders = orders.filter((order) => {
    if (filter === 'all') {
      return true
    }
    return order.status === filter
  })

  const markPending = (orderId: number) => {
    setPendingOrderIds((current) =>
      current.includes(orderId) ? current : [...current, orderId]
    )
  }

  const clearPending = (orderId: number) => {
    setPendingOrderIds((current) => current.filter((id) => id !== orderId))
  }

  const refresh = async () => {
    if (isOffline) {
      setError('Offline mode: unable to refresh until connection is restored')
      return
    }

    try {
      setIsRefreshing(true)
      const [queue, queueStats, queueSettings] = await Promise.all([
        fetchQueue({ timeoutMs: STANDARD_REQUEST_TIMEOUT_MS, attempts: RETRY_ATTEMPTS }),
        fetchStats({ timeoutMs: STANDARD_REQUEST_TIMEOUT_MS, attempts: RETRY_ATTEMPTS }),
        fetchSettings({ timeoutMs: STANDARD_REQUEST_TIMEOUT_MS, attempts: RETRY_ATTEMPTS }),
      ])
      if (!mountedRef.current) {
        return
      }
      setOrders(queue)
      setStats(queueStats)
      applySettings(queueSettings)
      persistSnapshot(queue, queueStats, queueSettings)
      setError(null)
    } catch (err) {
      if (!mountedRef.current) {
        return
      }
      setError(getErrorMessage(err, 'Failed to refresh kitchen display data'))
    } finally {
      if (mountedRef.current) {
        setIsRefreshing(false)
      }
    }
  }

  const bumpOrder = async (orderId: number) => {
    const now = Date.now()
    if (now - lastBumpAttemptRef.current < FRONTEND_BUMP_DEBOUNCE_MS) {
      return
    }
    lastBumpAttemptRef.current = now

    markPending(orderId)
    const snapshot = optimisticBump(orderId)

    if (isOffline) {
      const queued = await enqueueAction('bump', orderId)
      if (!queued) {
        rollbackOrder(snapshot)
        setError('Failed to queue bump action while offline')
      } else {
        setError('Offline mode: bump action queued and will sync automatically')
      }
      clearPending(orderId)
      return
    }

    try {
      const response = await withRetry(() => api.put(`/kds/orders/${orderId}/bump`))
      const updatedOrder = extractData<KDSOrder | null>(response.data, null)
      if (updatedOrder) {
        upsertOrder(updatedOrder)
        persistSnapshot()
      } else {
        const queue = await fetchQueue({ timeoutMs: STANDARD_REQUEST_TIMEOUT_MS, attempts: 1 })
        setOrders(queue)
        persistSnapshot(queue)
      }
      setError(null)
    } catch (err) {
      rollbackOrder(snapshot)
      const message = getErrorMessage(err, 'Failed to bump order')
      setError(message)
    } finally {
      clearPending(orderId)
    }
  }

  const recallOrder = async (orderId: number) => {
    markPending(orderId)
    const snapshot = optimisticRecall(orderId)

    if (isOffline) {
      const queued = await enqueueAction('recall', orderId)
      if (!queued) {
        rollbackOrder(snapshot)
        setError('Failed to queue recall action while offline')
      } else {
        setError('Offline mode: recall action queued and will sync automatically')
      }
      clearPending(orderId)
      return
    }

    try {
      const response = await withRetry(() => api.put(`/kds/orders/${orderId}/recall`))
      const updatedOrder = extractData<KDSOrder | null>(response.data, null)
      if (updatedOrder) {
        upsertOrder(updatedOrder)
        persistSnapshot()
      } else {
        const queue = await fetchQueue({ timeoutMs: STANDARD_REQUEST_TIMEOUT_MS, attempts: 1 })
        setOrders(queue)
        persistSnapshot(queue)
      }
      setError(null)
    } catch (err) {
      rollbackOrder(snapshot)
      const message = getErrorMessage(err, 'Failed to recall order')
      setError(message)
    } finally {
      clearPending(orderId)
    }
  }

  return {
    orders,
    filteredOrders,
    stats,
    settings,
    filter,
    soundEnabled,
    isFullscreen,
    isConnected,
    isOffline,
    queuedActionsCount,
    isReplayingOfflineQueue: isReplayingQueue,
    loading,
    isRefreshing,
    error,
    pendingOrderIds,

    setFilter: (nextFilter: KDSFilter) => setFilter(nextFilter),
    toggleSound,
    setFullscreen,
    toggleFullscreen,
    setConnected,
    refresh,
    bumpOrder,
    recallOrder,
    isOrderPending: (orderId: number) => pendingOrderIds.includes(orderId),
  }
}
