'use client'

import { create } from 'zustand'
import { persist, type PersistStorage, type StorageValue } from 'zustand/middleware'

export type PaymentMethod = 'cash' | 'card' | 'wallet'
export type POSProductCardSize = 'small' | 'large'

export interface POSAddItemProduct {
  id: number
  name: string
  price: number
  image?: string
  track_inventory?: boolean
  quantity_in_stock?: number
  trackInventory?: boolean
  quantityInStock?: number
}

export interface POSCartItem {
  productId: number
  name: string
  price: number
  quantity: number
  notes?: string
  image?: string
  trackInventory?: boolean
  quantityInStock?: number
}

export interface HeldOrder {
  id: string
  items: POSCartItem[]
  tableId: number | null
  tableName: string | null
  customerName: string
  customerPhone: string
  paymentMethod: PaymentMethod
  orderNotes: string
  heldAt: string
}

export interface POSState {
  items: POSCartItem[]
  tableId: number | null
  tableName: string | null
  customerName: string
  customerPhone: string
  paymentMethod: PaymentMethod
  orderNotes: string
  heldOrders: HeldOrder[]
  selectedCategoryId: number | null
  searchQuery: string
  soundEnabled: boolean
  autoPrintEnabled: boolean
  focusModeEnabled: boolean
  hideOutOfStockProducts: boolean
  productCardSize: POSProductCardSize

  addItem: (product: POSAddItemProduct) => void
  removeItem: (productId: number) => void
  updateQuantity: (productId: number, quantity: number) => void
  updateItemNotes: (productId: number, notes: string) => void
  clearCart: () => void
  setTableContext: (tableId: number | null, tableName?: string | null) => void
  setCustomer: (name: string, phone: string) => void
  setPaymentMethod: (method: PaymentMethod) => void
  setOrderNotes: (notes: string) => void
  setSelectedCategoryId: (categoryId: number | null) => void
  setSearchQuery: (query: string) => void
  setSoundEnabled: (enabled: boolean) => void
  setAutoPrintEnabled: (enabled: boolean) => void
  setFocusModeEnabled: (enabled: boolean) => void
  setHideOutOfStockProducts: (enabled: boolean) => void
  setProductCardSize: (size: POSProductCardSize) => void
  toggleSound: () => void
  toggleAutoPrint: () => void
  toggleFocusMode: () => void
  holdCurrentOrder: () => string | null
  recallOrder: (id: string) => void
  deleteHeldOrder: (id: string) => void

  getSubtotal: () => number
  getTax: () => number
  getTotal: () => number
  getTotalItems: () => number
}

type POSPersistedState = Pick<
  POSState,
  | 'items'
  | 'tableId'
  | 'tableName'
  | 'customerName'
  | 'customerPhone'
  | 'paymentMethod'
  | 'heldOrders'
  | 'soundEnabled'
  | 'autoPrintEnabled'
  | 'focusModeEnabled'
  | 'hideOutOfStockProducts'
  | 'productCardSize'
>

const POS_STORAGE_NAME = 'pos-storage'
const TAX_RATE = 0.14
const VALID_PAYMENT_METHODS: PaymentMethod[] = ['cash', 'card', 'wallet']
const VALID_PRODUCT_CARD_SIZES: POSProductCardSize[] = ['small', 'large']

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null

const isPaymentMethod = (value: unknown): value is PaymentMethod =>
  typeof value === 'string' && (VALID_PAYMENT_METHODS as string[]).includes(value)

const isPOSProductCardSize = (value: unknown): value is POSProductCardSize =>
  typeof value === 'string' && (VALID_PRODUCT_CARD_SIZES as string[]).includes(value)

const toFiniteNumber = (value: unknown, fallback = 0): number => {
  const numeric = typeof value === 'number' ? value : Number(value)
  return Number.isFinite(numeric) ? numeric : fallback
}

const toPositiveInt = (value: unknown, fallback = 1): number =>
  Math.max(1, Math.floor(toFiniteNumber(value, fallback)))

const toNonNegativeInt = (value: unknown, fallback = 0): number =>
  Math.max(0, Math.floor(toFiniteNumber(value, fallback)))

const toPositiveIntOrNull = (value: unknown): number | null => {
  const numeric = toFiniteNumber(value, -1)
  if (numeric <= 0) {
    return null
  }
  return Math.floor(numeric)
}

const normalizeText = (value: unknown, fallback = ''): string =>
  typeof value === 'string' ? value.trim() : fallback

const normalizeNullableText = (value: unknown): string | null => {
  if (typeof value !== 'string') {
    return null
  }
  const trimmed = value.trim()
  return trimmed.length > 0 ? trimmed : null
}

const roundMoney = (value: number): number =>
  Math.round((value + Number.EPSILON) * 100) / 100

const clampQuantityByInventory = (
  quantity: number,
  trackInventory: boolean,
  quantityInStock: number
) => {
  if (!trackInventory) {
    return quantity
  }
  if (quantityInStock <= 0) {
    return 0
  }
  return Math.min(quantity, quantityInStock)
}

const cloneCartItem = (item: POSCartItem): POSCartItem => ({
  ...item,
})

const createHeldOrderId = () =>
  `held-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`

const normalizeCartItem = (value: unknown): POSCartItem | null => {
  if (!isRecord(value)) {
    return null
  }

  const productId = toPositiveIntOrNull(value.productId)
  const name = normalizeText(value.name)
  const price = toFiniteNumber(value.price, -1)
  const quantity = toPositiveInt(value.quantity, 1)
  const trackInventory = typeof value.trackInventory === 'boolean' ? value.trackInventory : false
  const quantityInStock = toNonNegativeInt(value.quantityInStock, 0)

  if (!productId || !name || price < 0) {
    return null
  }

  const normalizedQuantity = clampQuantityByInventory(quantity, trackInventory, quantityInStock)
  if (normalizedQuantity <= 0) {
    return null
  }

  const notes = normalizeText(value.notes)
  const image = normalizeText(value.image)

  return {
    productId,
    name,
    price: roundMoney(price),
    quantity: normalizedQuantity,
    notes: notes || undefined,
    image: image || undefined,
    trackInventory,
    quantityInStock:
      value.quantityInStock !== undefined || trackInventory
        ? quantityInStock
        : undefined,
  }
}

const normalizeHeldOrder = (value: unknown): HeldOrder | null => {
  if (!isRecord(value)) {
    return null
  }

  const id = normalizeText(value.id)
  if (!id) {
    return null
  }

  const rawItems = Array.isArray(value.items) ? value.items : []
  const items = rawItems
    .map((item) => normalizeCartItem(item))
    .filter((item): item is POSCartItem => item !== null)

  if (items.length === 0) {
    return null
  }

  return {
    id,
    items,
    tableId: toPositiveIntOrNull(value.tableId),
    tableName: normalizeNullableText(value.tableName),
    customerName: normalizeText(value.customerName),
    customerPhone: normalizeText(value.customerPhone),
    paymentMethod: isPaymentMethod(value.paymentMethod) ? value.paymentMethod : 'cash',
    orderNotes: normalizeText(value.orderNotes),
    heldAt:
      typeof value.heldAt === 'string' && value.heldAt
        ? value.heldAt
        : new Date().toISOString(),
  }
}

const normalizePersistedState = (value: unknown): POSPersistedState | null => {
  if (!isRecord(value)) {
    return null
  }

  const rawItems = Array.isArray(value.items) ? value.items : []
  const rawHeldOrders = Array.isArray(value.heldOrders) ? value.heldOrders : []

  return {
    items: rawItems
      .map((item) => normalizeCartItem(item))
      .filter((item): item is POSCartItem => item !== null),
    tableId: toPositiveIntOrNull(value.tableId),
    tableName: normalizeNullableText(value.tableName),
    customerName: normalizeText(value.customerName),
    customerPhone: normalizeText(value.customerPhone),
    paymentMethod: isPaymentMethod(value.paymentMethod) ? value.paymentMethod : 'cash',
    heldOrders: rawHeldOrders
      .map((order) => normalizeHeldOrder(order))
      .filter((order): order is HeldOrder => order !== null),
    soundEnabled: typeof value.soundEnabled === 'boolean' ? value.soundEnabled : true,
    autoPrintEnabled:
      typeof value.autoPrintEnabled === 'boolean' ? value.autoPrintEnabled : false,
    focusModeEnabled:
      typeof value.focusModeEnabled === 'boolean' ? value.focusModeEnabled : false,
    hideOutOfStockProducts:
      typeof value.hideOutOfStockProducts === 'boolean'
        ? value.hideOutOfStockProducts
        : false,
    productCardSize: isPOSProductCardSize(value.productCardSize)
      ? value.productCardSize
      : 'small',
  }
}

const posPersistStorage: PersistStorage<POSPersistedState> = {
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
      // Corrupted payloads should not crash POS boot.
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

const resolveInventory = (
  product: POSAddItemProduct,
  fallbackTrackInventory: boolean,
  fallbackQuantityInStock: number
) => {
  const trackInventory =
    typeof product.trackInventory === 'boolean'
      ? product.trackInventory
      : typeof product.track_inventory === 'boolean'
        ? product.track_inventory
        : fallbackTrackInventory

  const quantityInStock = toNonNegativeInt(
    product.quantityInStock ?? product.quantity_in_stock ?? fallbackQuantityInStock,
    0
  )

  return { trackInventory, quantityInStock }
}

export const usePOSStore = create<POSState>()(
  persist(
    (set, get) => ({
      items: [],
      tableId: null,
      tableName: null,
      customerName: '',
      customerPhone: '',
      paymentMethod: 'cash',
      orderNotes: '',
      heldOrders: [],
      selectedCategoryId: null,
      searchQuery: '',
      soundEnabled: true,
      autoPrintEnabled: false,
      focusModeEnabled: false,
      hideOutOfStockProducts: false,
      productCardSize: 'small',

      addItem: (product) => {
        set((state) => {
          const productId = toPositiveIntOrNull(product.id)
          const name = normalizeText(product.name)
          const price = toFiniteNumber(product.price, -1)

          if (!productId || !name || price < 0) {
            return state
          }

          const existingItem = state.items.find((item) => item.productId === productId)
          const inventory = resolveInventory(
            product,
            existingItem?.trackInventory ?? false,
            existingItem?.quantityInStock ?? 0
          )

          if (inventory.trackInventory && inventory.quantityInStock <= 0) {
            return state
          }

          if (existingItem) {
            const nextQuantity = clampQuantityByInventory(
              existingItem.quantity + 1,
              inventory.trackInventory,
              inventory.quantityInStock
            )

            if (nextQuantity === existingItem.quantity) {
              return state
            }

            return {
              items: state.items.map((item) =>
                item.productId === productId
                  ? {
                      ...item,
                      quantity: nextQuantity,
                      trackInventory: inventory.trackInventory,
                      quantityInStock: inventory.quantityInStock,
                    }
                  : item
              ),
            }
          }

          return {
            items: [
              ...state.items,
              {
                productId,
                name,
                price: roundMoney(price),
                quantity: 1,
                image: normalizeText(product.image) || undefined,
                trackInventory: inventory.trackInventory,
                quantityInStock: inventory.quantityInStock,
              },
            ],
          }
        })
      },

      removeItem: (productId) =>
        set((state) => ({
          items: state.items.filter((item) => item.productId !== productId),
        })),

      updateQuantity: (productId, quantity) => {
        set((state) => ({
          items: state.items
            .map((item) => {
              if (item.productId !== productId) {
                return item
              }

              const requestedQuantity = toFiniteNumber(quantity, 0)
              if (requestedQuantity <= 0) {
                return null
              }

              const normalizedQuantity = clampQuantityByInventory(
                Math.max(1, Math.floor(requestedQuantity)),
                item.trackInventory ?? false,
                toNonNegativeInt(item.quantityInStock, 0)
              )

              if (normalizedQuantity <= 0) {
                return null
              }

              return {
                ...item,
                quantity: normalizedQuantity,
              }
            })
            .filter((item): item is POSCartItem => item !== null),
        }))
      },

      updateItemNotes: (productId, notes) =>
        set((state) => ({
          items: state.items.map((item) =>
            item.productId === productId
              ? {
                  ...item,
                  notes: normalizeText(notes) || undefined,
                }
              : item
          ),
        })),

      clearCart: () =>
        set({
          items: [],
          tableId: null,
          tableName: null,
          customerName: '',
          customerPhone: '',
          paymentMethod: 'cash',
          orderNotes: '',
        }),

      setTableContext: (tableId, tableName) => {
        const nextTableId = toPositiveIntOrNull(tableId)
        set({
          tableId: nextTableId,
          tableName: nextTableId ? normalizeNullableText(tableName) : null,
        })
      },

      setCustomer: (name, phone) =>
        set({
          customerName: normalizeText(name),
          customerPhone: normalizeText(phone),
        }),

      setPaymentMethod: (method) =>
        set({
          paymentMethod: isPaymentMethod(method) ? method : 'cash',
        }),

      setOrderNotes: (notes) =>
        set({
          orderNotes: normalizeText(notes),
        }),

      setSelectedCategoryId: (categoryId) =>
        set({
          selectedCategoryId: toPositiveIntOrNull(categoryId),
        }),

      setSearchQuery: (query) =>
        set({
          searchQuery: normalizeText(query),
        }),

      setSoundEnabled: (enabled) => set({ soundEnabled: enabled }),
      setAutoPrintEnabled: (enabled) => set({ autoPrintEnabled: enabled }),
      setFocusModeEnabled: (enabled) => set({ focusModeEnabled: enabled }),
      setHideOutOfStockProducts: (enabled) =>
        set({ hideOutOfStockProducts: enabled }),
      setProductCardSize: (size) =>
        set({ productCardSize: isPOSProductCardSize(size) ? size : 'small' }),
      toggleSound: () => set((state) => ({ soundEnabled: !state.soundEnabled })),
      toggleAutoPrint: () =>
        set((state) => ({ autoPrintEnabled: !state.autoPrintEnabled })),
      toggleFocusMode: () =>
        set((state) => ({ focusModeEnabled: !state.focusModeEnabled })),

      holdCurrentOrder: () => {
        const currentState = get()
        if (currentState.items.length === 0) {
          return null
        }

        const heldOrder: HeldOrder = {
          id: createHeldOrderId(),
          items: currentState.items.map((item) => cloneCartItem(item)),
          tableId: currentState.tableId,
          tableName: currentState.tableName,
          customerName: currentState.customerName,
          customerPhone: currentState.customerPhone,
          paymentMethod: currentState.paymentMethod,
          orderNotes: currentState.orderNotes,
          heldAt: new Date().toISOString(),
        }

        set((state) => ({
          heldOrders: [heldOrder, ...state.heldOrders],
          items: [],
          tableId: null,
          tableName: null,
          customerName: '',
          customerPhone: '',
          paymentMethod: 'cash',
          orderNotes: '',
        }))

        return heldOrder.id
      },

      recallOrder: (id) =>
        set((state) => {
          const heldOrder = state.heldOrders.find((order) => order.id === id)
          if (!heldOrder) {
            return state
          }

          return {
            items: heldOrder.items.map((item) => cloneCartItem(item)),
            tableId: heldOrder.tableId,
            tableName: heldOrder.tableName,
            customerName: heldOrder.customerName,
            customerPhone: heldOrder.customerPhone,
            paymentMethod: heldOrder.paymentMethod,
            orderNotes: heldOrder.orderNotes,
            heldOrders: state.heldOrders.filter((order) => order.id !== id),
          }
        }),

      deleteHeldOrder: (id) =>
        set((state) => ({
          heldOrders: state.heldOrders.filter((order) => order.id !== id),
        })),

      getSubtotal: () =>
        roundMoney(
          get().items.reduce((total, item) => total + item.price * item.quantity, 0)
        ),
      getTax: () => roundMoney(get().getSubtotal() * TAX_RATE),
      getTotal: () => roundMoney(get().getSubtotal() + get().getTax()),
      getTotalItems: () =>
        get().items.reduce((total, item) => total + item.quantity, 0),
    }),
    {
      name: POS_STORAGE_NAME,
      storage: posPersistStorage,
      partialize: (state) => ({
        items: state.items,
        tableId: state.tableId,
        tableName: state.tableName,
        customerName: state.customerName,
        customerPhone: state.customerPhone,
        paymentMethod: state.paymentMethod,
        heldOrders: state.heldOrders,
        soundEnabled: state.soundEnabled,
        autoPrintEnabled: state.autoPrintEnabled,
        focusModeEnabled: state.focusModeEnabled,
        hideOutOfStockProducts: state.hideOutOfStockProducts,
        productCardSize: state.productCardSize,
      }),
    }
  )
)

export const selectPOSItems = (state: POSState) => state.items
export const selectPOSHeldOrders = (state: POSState) => state.heldOrders
export const selectPOSHeldOrdersCount = (state: POSState) => state.heldOrders.length
export const selectPOSTableId = (state: POSState) => state.tableId
export const selectPOSTableName = (state: POSState) => state.tableName
export const selectPOSCustomerName = (state: POSState) => state.customerName
export const selectPOSCustomerPhone = (state: POSState) => state.customerPhone
export const selectPOSPaymentMethod = (state: POSState) => state.paymentMethod
export const selectPOSOrderNotes = (state: POSState) => state.orderNotes
export const selectPOSSelectedCategoryId = (state: POSState) => state.selectedCategoryId
export const selectPOSSearchQuery = (state: POSState) => state.searchQuery
export const selectPOSSoundEnabled = (state: POSState) => state.soundEnabled
export const selectPOSAutoPrintEnabled = (state: POSState) => state.autoPrintEnabled
export const selectPOSFocusModeEnabled = (state: POSState) => state.focusModeEnabled
export const selectPOSHideOutOfStockProducts = (state: POSState) =>
  state.hideOutOfStockProducts
export const selectPOSProductCardSize = (state: POSState) => state.productCardSize
export const selectPOSSubtotal = (state: POSState) => state.getSubtotal()
export const selectPOSTax = (state: POSState) => state.getTax()
export const selectPOSTotal = (state: POSState) => state.getTotal()
export const selectPOSTotalItems = (state: POSState) => state.getTotalItems()
export const selectPOSItemByProductId = (productId: number) => (state: POSState) =>
  state.items.find((item) => item.productId === productId)
