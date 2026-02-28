'use client'

import { create } from 'zustand'
import { persist, type PersistStorage, type StorageValue } from 'zustand/middleware'

export type KDSOrderStatus = 'new' | 'in_progress' | 'ready' | 'completed' | 'recalled'
export type KDSFilter = 'all' | 'new' | 'in_progress' | 'ready' | 'completed'

export interface KDSItem {
  product_id?: number
  product_name: string
  quantity: number
  special_instructions?: string
  add_ons?: string
}

export interface KDSOrder {
  id: number
  tenant_id: number
  restaurant_id: number
  order_id: number
  order_number: string
  order_source?: string | null
  table_id?: number | null
  table_number?: number | null
  table_name?: string | null
  customer_name: string
  status: KDSOrderStatus
  priority: number
  station_id: string
  items: KDSItem[]
  notes?: string | null
  received_at: string
  started_at?: string | null
  completed_at?: string | null
  bumped_by?: number | null
  elapsed_seconds: number
  created_at: string
  updated_at: string
}

export interface KDSStats {
  active_orders: number
  avg_prep_minutes: number
  completed_today: number
  new_orders: number
  in_progress: number
  ready: number
}

interface KDSState {
  orders: KDSOrder[]
  stats: KDSStats | null
  filter: KDSFilter
  soundEnabled: boolean
  isFullscreen: boolean
  isConnected: boolean

  setOrders: (orders: KDSOrder[]) => void
  addOrder: (order: KDSOrder) => void
  upsertOrder: (order: KDSOrder) => void
  updateOrder: (id: number, updates: Partial<KDSOrder>) => void
  removeOrder: (id: number) => void

  setStats: (stats: KDSStats | null) => void
  setFilter: (filter: KDSFilter) => void
  setSoundEnabled: (enabled: boolean) => void
  toggleSound: () => void
  setFullscreen: (isFullscreen: boolean) => void
  toggleFullscreen: () => void
  setConnected: (isConnected: boolean) => void

  optimisticBump: (id: number) => KDSOrder | null
  optimisticRecall: (id: number) => KDSOrder | null
  rollbackOrder: (snapshot: KDSOrder | null) => void
}

type KDSPersistedState = Pick<KDSState, 'filter' | 'soundEnabled' | 'isFullscreen'>

const KDS_STORAGE_NAME = 'kds-storage'
const VALID_FILTERS: KDSFilter[] = ['all', 'new', 'in_progress', 'ready', 'completed']

const toTimestamp = (value: string) => {
  const timestamp = new Date(value).getTime()
  return Number.isFinite(timestamp) ? timestamp : 0
}

const isKDSFilter = (value: unknown): value is KDSFilter =>
  typeof value === 'string' && (VALID_FILTERS as string[]).includes(value)

const normalizePersistedState = (value: unknown): KDSPersistedState | null => {
  if (!value || typeof value !== 'object') {
    return null
  }

  const raw = value as Record<string, unknown>
  return {
    filter: isKDSFilter(raw.filter) ? raw.filter : 'all',
    soundEnabled: typeof raw.soundEnabled === 'boolean' ? raw.soundEnabled : true,
    isFullscreen: typeof raw.isFullscreen === 'boolean' ? raw.isFullscreen : false,
  }
}

const kdsPersistStorage: PersistStorage<KDSPersistedState> = {
  getItem: (name) => {
    if (typeof window === 'undefined') {
      return null
    }

    const raw = window.localStorage.getItem(name)
    if (!raw) {
      return null
    }

    try {
      const parsed = JSON.parse(raw) as StorageValue<unknown>
      const state = normalizePersistedState(parsed?.state)
      if (!state) {
        window.localStorage.removeItem(name)
        return null
      }

      return {
        state,
        version: typeof parsed?.version === 'number' ? parsed.version : 0,
      }
    } catch {
      // Corrupted payloads should not crash KDS boot.
      window.localStorage.removeItem(name)
      return null
    }
  },
  setItem: (name, value) => {
    if (typeof window === 'undefined') {
      return
    }

    try {
      window.localStorage.setItem(name, JSON.stringify(value))
    } catch {
      // Ignore storage write errors.
    }
  },
  removeItem: (name) => {
    if (typeof window === 'undefined') {
      return
    }
    window.localStorage.removeItem(name)
  },
}

const sortOrders = (orders: KDSOrder[]) =>
  [...orders].sort((left, right) => {
    if (left.priority !== right.priority) {
      return right.priority - left.priority
    }
    return toTimestamp(left.received_at) - toTimestamp(right.received_at)
  })

const cloneOrder = (order: KDSOrder): KDSOrder => ({
  ...order,
  items: Array.isArray(order.items) ? order.items.map((item) => ({ ...item })) : [],
})

const getNextBumpStatus = (status: KDSOrderStatus): KDSOrderStatus | null => {
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

export const useKDSStore = create<KDSState>()(
  persist(
    (set, get) => ({
      orders: [],
      stats: null,
      filter: 'all',
      soundEnabled: true,
      isFullscreen: false,
      isConnected: false,

      setOrders: (orders) => set({ orders: sortOrders(orders ?? []) }),

      addOrder: (order) =>
        set((state) => ({
          orders: sortOrders([order, ...state.orders.filter((item) => item.id !== order.id)]),
        })),

      upsertOrder: (order) =>
        set((state) => ({
          orders: sortOrders([
            ...state.orders.filter((item) => item.id !== order.id),
            order,
          ]),
        })),

      updateOrder: (id, updates) =>
        set((state) => ({
          orders: sortOrders(
            state.orders.map((order) =>
              order.id === id ? { ...order, ...updates } : order
            )
          ),
        })),

      removeOrder: (id) =>
        set((state) => ({
          orders: state.orders.filter((order) => order.id !== id),
        })),

      setStats: (stats) => set({ stats }),
      setFilter: (filter) => set({ filter }),
      setSoundEnabled: (enabled) => set({ soundEnabled: enabled }),
      toggleSound: () => set((state) => ({ soundEnabled: !state.soundEnabled })),
      setFullscreen: (isFullscreen) => set({ isFullscreen }),
      toggleFullscreen: () => set((state) => ({ isFullscreen: !state.isFullscreen })),
      setConnected: (isConnected) => set({ isConnected }),

      optimisticBump: (id) => {
        const current = get().orders.find((order) => order.id === id)
        if (!current) {
          return null
        }

        const nextStatus = getNextBumpStatus(current.status)
        if (!nextStatus) {
          return cloneOrder(current)
        }

        const snapshot = cloneOrder(current)
        const now = new Date().toISOString()

        get().updateOrder(id, {
          status: nextStatus,
          started_at:
            current.status === 'new'
              ? now
              : current.started_at ?? null,
          completed_at:
            nextStatus === 'completed'
              ? now
              : current.completed_at ?? null,
          updated_at: now,
        })

        return snapshot
      },

      optimisticRecall: (id) => {
        const current = get().orders.find((order) => order.id === id)
        if (!current) {
          return null
        }

        const snapshot = cloneOrder(current)
        get().updateOrder(id, {
          status: 'ready',
          completed_at: null,
          updated_at: new Date().toISOString(),
        })

        return snapshot
      },

      rollbackOrder: (snapshot) => {
        if (!snapshot) {
          return
        }
        get().upsertOrder(snapshot)
      },
    }),
    {
      name: KDS_STORAGE_NAME,
      storage: kdsPersistStorage,
      partialize: (state) => ({
        filter: state.filter,
        soundEnabled: state.soundEnabled,
        isFullscreen: state.isFullscreen,
      }),
    }
  )
)
