'use client'

import { create } from 'zustand'
import { toast } from 'sonner'

export interface CartItem {
  id: number
  productId: number
  name: string
  price: number
  quantity: number
  trackInventory?: boolean
  quantityInStock?: number
  image?: string
  notes?: string
  total: number
}

export interface AddToCartProduct {
  id: number
  name: string
  price: number
  image?: string
  track_inventory?: boolean
  quantity_in_stock?: number
}

interface CartSnapshot {
  items: CartItem[]
  tableId: number | null
  tableName: string | null
}

interface CartStore {
  items: CartItem[]
  tableId: number | null
  tableName: string | null
  restaurantSlug: string
  totalPrice: number
  totalItems: number
  setRestaurantSlug: (restaurantSlug: string) => void
  setTableContext: (tableId: number, tableName?: string | null) => void
  clearTableContext: () => void
  addItem: (product: AddToCartProduct) => void
  removeItem: (productId: number) => void
  updateQuantity: (productId: number, quantity: number) => void
  updateNotes: (productId: number, notes: string) => void
  clearCart: () => void
  getTotal: () => number
  getTotalItems: () => number
}

const STORAGE_PREFIX = 'pos-cart-'
const SYNC_PREFIX = 'pos-cart-sync-'
const DEFAULT_SLUG = 'default'

let syncInitialized = false
let cartChannel: BroadcastChannel | null = null

function computeTotalPrice(items: CartItem[]): number {
  return items.reduce((acc, item) => acc + item.total, 0)
}

function computeTotalItems(items: CartItem[]): number {
  return items.reduce((acc, item) => acc + item.quantity, 0)
}

function getStorageKey(slug: string): string {
  return `${STORAGE_PREFIX}${slug || DEFAULT_SLUG}`
}

function getSyncKey(slug: string): string {
  return `${SYNC_PREFIX}${slug || DEFAULT_SLUG}`
}

function readSnapshot(slug: string): CartSnapshot {
  if (typeof window === 'undefined') {
    return { items: [], tableId: null, tableName: null }
  }

  try {
    const raw = localStorage.getItem(getStorageKey(slug))
    if (!raw) {
      return { items: [], tableId: null, tableName: null }
    }
    const parsed = JSON.parse(raw) as CartSnapshot
    return {
      items: Array.isArray(parsed.items) ? parsed.items : [],
      tableId: typeof parsed.tableId === 'number' ? parsed.tableId : null,
      tableName: typeof parsed.tableName === 'string' ? parsed.tableName : null,
    }
  } catch {
    return { items: [], tableId: null, tableName: null }
  }
}

function saveSnapshot(slug: string, snapshot: CartSnapshot) {
  if (typeof window === 'undefined') return
  localStorage.setItem(getStorageKey(slug), JSON.stringify(snapshot))
  localStorage.setItem(getSyncKey(slug), JSON.stringify({ at: Date.now() }))
  if (cartChannel) {
    cartChannel.postMessage({ type: 'cart-sync', slug })
  }
}

function syncStoreFromStorage(slug: string, set: (partial: Partial<CartStore>) => void) {
  const snapshot = readSnapshot(slug)
  set({
    items: snapshot.items,
    tableId: snapshot.tableId,
    tableName: snapshot.tableName,
    totalItems: computeTotalItems(snapshot.items),
    totalPrice: computeTotalPrice(snapshot.items),
  })
}

function initializeCrossTabSync(
  set: (partial: Partial<CartStore>) => void,
  get: () => CartStore
) {
  if (syncInitialized || typeof window === 'undefined') return
  syncInitialized = true

  if (typeof BroadcastChannel !== 'undefined') {
    cartChannel = new BroadcastChannel('pos-cart-sync')
    cartChannel.onmessage = () => {
      const currentSlug = get().restaurantSlug
      syncStoreFromStorage(currentSlug, set)
    }
  }

  window.addEventListener('storage', (event) => {
    if (!event.key || !event.key.startsWith(SYNC_PREFIX)) {
      return
    }
    const currentSlug = get().restaurantSlug
    if (event.key !== getSyncKey(currentSlug)) {
      return
    }
    syncStoreFromStorage(currentSlug, set)
  })
}

export const useCartStore = create<CartStore>((set, get) => ({
  items: [],
  tableId: null,
  tableName: null,
  restaurantSlug: DEFAULT_SLUG,
  totalPrice: 0,
  totalItems: 0,

  setRestaurantSlug: (restaurantSlug) => {
    const slug = (restaurantSlug || DEFAULT_SLUG).trim() || DEFAULT_SLUG
    initializeCrossTabSync(set, get)
    const snapshot = readSnapshot(slug)
    set({
      restaurantSlug: slug,
      items: snapshot.items,
      tableId: snapshot.tableId,
      tableName: snapshot.tableName,
      totalItems: computeTotalItems(snapshot.items),
      totalPrice: computeTotalPrice(snapshot.items),
    })
  },

  setTableContext: (tableId, tableName) => {
    const next = {
      tableId,
      tableName: tableName?.trim() || null,
    }
    set(next)
    const state = get()
    saveSnapshot(state.restaurantSlug, {
      items: state.items,
      tableId: tableId,
      tableName: next.tableName,
    })
  },

  clearTableContext: () => {
    set({ tableId: null, tableName: null })
    const state = get()
    saveSnapshot(state.restaurantSlug, {
      items: state.items,
      tableId: null,
      tableName: null,
    })
  },

  addItem: (product) => {
    const state = get()
    initializeCrossTabSync(set, get)

    const existingItem = state.items.find((item) => item.productId === product.id)
    const quantityInStock = Number(product.quantity_in_stock ?? existingItem?.quantityInStock ?? 0)
    const trackInventory =
      typeof product.track_inventory === 'boolean'
        ? product.track_inventory
        : existingItem?.trackInventory ?? false

    if (trackInventory && quantityInStock <= 0) {
      toast.warning('Product is out of stock')
      return
    }

    const nextItems = existingItem
      ? state.items.map((item) => {
          if (item.productId !== product.id) return item
          const requestedQuantity = item.quantity + 1
          const cappedQuantity =
            trackInventory && quantityInStock > 0
              ? Math.min(requestedQuantity, quantityInStock)
              : requestedQuantity
          if (cappedQuantity === item.quantity) {
            toast.warning('Maximum available stock reached')
            return item
          }
          return {
            ...item,
            quantity: cappedQuantity,
            total: cappedQuantity * item.price,
            quantityInStock,
            trackInventory,
          }
        })
      : [
          ...state.items,
          {
            id: Date.now(),
            productId: product.id,
            name: product.name,
            price: product.price,
            quantity: 1,
            trackInventory,
            quantityInStock,
            image: product.image,
            total: product.price,
          },
        ]

    const snapshot: CartSnapshot = {
      items: nextItems,
      tableId: state.tableId,
      tableName: state.tableName,
    }
    saveSnapshot(state.restaurantSlug, snapshot)
    set({
      items: nextItems,
      totalItems: computeTotalItems(nextItems),
      totalPrice: computeTotalPrice(nextItems),
    })
  },

  removeItem: (productId) => {
    const state = get()
    const nextItems = state.items.filter((item) => item.productId !== productId)
    saveSnapshot(state.restaurantSlug, {
      items: nextItems,
      tableId: state.tableId,
      tableName: state.tableName,
    })
    set({
      items: nextItems,
      totalItems: computeTotalItems(nextItems),
      totalPrice: computeTotalPrice(nextItems),
    })
  },

  updateQuantity: (productId, quantity) => {
    const state = get()
    const nextItems = state.items
      .map((item) => {
        if (item.productId !== productId) return item
        if (quantity <= 0) return null

        const available = Number(item.quantityInStock ?? 0)
        const trackInventory = item.trackInventory ?? false
        const cappedQuantity =
          trackInventory && available > 0 ? Math.min(quantity, available) : quantity

        if (trackInventory && cappedQuantity < quantity) {
          toast.warning('Maximum available stock reached')
        }

        if (cappedQuantity <= 0) return null
        return {
          ...item,
          quantity: cappedQuantity,
          total: cappedQuantity * item.price,
        }
      })
      .filter((item): item is CartItem => item !== null)

    saveSnapshot(state.restaurantSlug, {
      items: nextItems,
      tableId: state.tableId,
      tableName: state.tableName,
    })
    set({
      items: nextItems,
      totalItems: computeTotalItems(nextItems),
      totalPrice: computeTotalPrice(nextItems),
    })
  },

  updateNotes: (productId, notes) => {
    const state = get()
    const nextItems = state.items.map((item) =>
      item.productId === productId ? { ...item, notes } : item
    )
    saveSnapshot(state.restaurantSlug, {
      items: nextItems,
      tableId: state.tableId,
      tableName: state.tableName,
    })
    set({
      items: nextItems,
      totalItems: computeTotalItems(nextItems),
      totalPrice: computeTotalPrice(nextItems),
    })
  },

  clearCart: () => {
    const state = get()
    const nextItems: CartItem[] = []
    saveSnapshot(state.restaurantSlug, {
      items: nextItems,
      tableId: state.tableId,
      tableName: state.tableName,
    })
    set({
      items: nextItems,
      totalItems: 0,
      totalPrice: 0,
    })
  },

  getTotal: () => get().totalPrice,
  getTotalItems: () => get().totalItems,
}))

export const selectCartItems = (state: CartStore) => state.items
export const selectCartTotal = (state: CartStore) => state.totalPrice
export const selectCartItemsCount = (state: CartStore) => state.totalItems
export const selectCartTableId = (state: CartStore) => state.tableId
export const selectCartTableName = (state: CartStore) => state.tableName
export const selectCartItemByProductId = (productId: number) => (state: CartStore) =>
  state.items.find((item) => item.productId === productId)

