'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import { type CreatePOSOrderPayload, posApi } from '@/lib/api/posApi'
import { debugError } from '@/lib/debug'

export interface POSQueuedOrder {
  id: string
  payload: CreatePOSOrderPayload
  queuedAt: string
  retryCount: number
}

interface ReplayResult {
  processed: number
  failed: number
}

interface POSOfflineResult {
  isOffline: boolean
  queuedOrderCount: number
  isReplaying: boolean
  enqueueOrder: (payload: CreatePOSOrderPayload) => Promise<boolean>
  replayQueuedOrders: () => Promise<ReplayResult>
  clearQueue: () => Promise<boolean>
}

const DB_NAME = 'pos-offline-db'
const DB_VERSION = 1
const STORE_NAME = 'pending-orders'

const canUseIndexedDB = () =>
  typeof window !== 'undefined' && typeof window.indexedDB !== 'undefined'

const canUseNavigatorOnline = () =>
  typeof navigator !== 'undefined' && typeof navigator.onLine === 'boolean'

const createQueueID = () =>
  `pos-${Date.now()}-${Math.random().toString(36).slice(2, 10)}`

const openQueueDB = (): Promise<IDBDatabase | null> => {
  if (!canUseIndexedDB()) {
    return Promise.resolve(null)
  }

  return new Promise((resolve) => {
    try {
      const request = window.indexedDB.open(DB_NAME, DB_VERSION)

      request.onupgradeneeded = () => {
        const db = request.result
        if (!db.objectStoreNames.contains(STORE_NAME)) {
          db.createObjectStore(STORE_NAME, { keyPath: 'id' })
        }
      }

      request.onsuccess = () => resolve(request.result)
      request.onerror = () => resolve(null)
    } catch {
      resolve(null)
    }
  })
}

const readAllPendingOrders = async (): Promise<POSQueuedOrder[]> => {
  const db = await openQueueDB()
  if (!db) {
    return []
  }

  return new Promise((resolve) => {
    const tx = db.transaction(STORE_NAME, 'readonly')
    const store = tx.objectStore(STORE_NAME)
    const request = store.getAll()

    request.onsuccess = () => {
      const rows = (request.result || []) as POSQueuedOrder[]
      rows.sort((left, right) => {
        if (left.queuedAt === right.queuedAt) {
          return left.id.localeCompare(right.id)
        }
        return left.queuedAt.localeCompare(right.queuedAt)
      })
      resolve(rows)
    }

    request.onerror = () => resolve([])
    tx.oncomplete = () => db.close()
    tx.onerror = () => db.close()
  })
}

const upsertPendingOrder = async (order: POSQueuedOrder): Promise<boolean> => {
  const db = await openQueueDB()
  if (!db) {
    return false
  }

  return new Promise((resolve) => {
    const tx = db.transaction(STORE_NAME, 'readwrite')
    const store = tx.objectStore(STORE_NAME)
    const request = store.put(order)

    request.onsuccess = () => resolve(true)
    request.onerror = () => resolve(false)
    tx.oncomplete = () => db.close()
    tx.onerror = () => db.close()
  })
}

const removePendingOrder = async (id: string): Promise<boolean> => {
  const db = await openQueueDB()
  if (!db) {
    return false
  }

  return new Promise((resolve) => {
    const tx = db.transaction(STORE_NAME, 'readwrite')
    const store = tx.objectStore(STORE_NAME)
    const request = store.delete(id)

    request.onsuccess = () => resolve(true)
    request.onerror = () => resolve(false)
    tx.oncomplete = () => db.close()
    tx.onerror = () => db.close()
  })
}

const clearPendingOrders = async (): Promise<boolean> => {
  const db = await openQueueDB()
  if (!db) {
    return false
  }

  return new Promise((resolve) => {
    const tx = db.transaction(STORE_NAME, 'readwrite')
    const store = tx.objectStore(STORE_NAME)
    const request = store.clear()

    request.onsuccess = () => resolve(true)
    request.onerror = () => resolve(false)
    tx.oncomplete = () => db.close()
    tx.onerror = () => db.close()
  })
}

export const usePOSOffline = (): POSOfflineResult => {
  const [isOffline, setIsOffline] = useState(false)
  const [queuedOrderCount, setQueuedOrderCount] = useState(0)
  const [isReplaying, setIsReplaying] = useState(false)
  const replayingRef = useRef(false)

  const refreshQueuedCount = useCallback(async () => {
    const queue = await readAllPendingOrders()
    setQueuedOrderCount(queue.length)
    return queue.length
  }, [])

  const replayQueuedOrders = useCallback(async (): Promise<ReplayResult> => {
    if (replayingRef.current) {
      return { processed: 0, failed: 0 }
    }

    if (canUseNavigatorOnline() && !navigator.onLine) {
      return { processed: 0, failed: 0 }
    }

    replayingRef.current = true
    setIsReplaying(true)

    let processed = 0
    let failed = 0

    try {
      const queue = await readAllPendingOrders()
      for (const pending of queue) {
        try {
          const order = await posApi.createOrder(pending.payload)
          if (pending.payload.payment_method === 'cash' && order.id > 0) {
            await posApi.updateOrderPaymentStatus(order.id, 'paid')
          }
          await removePendingOrder(pending.id)
          processed += 1
        } catch (error) {
          failed += 1
          debugError('[usePOSOffline] Failed to replay queued POS order:', error)
          await upsertPendingOrder({
            ...pending,
            retryCount: pending.retryCount + 1,
          })
          break
        }
      }
    } finally {
      await refreshQueuedCount()
      replayingRef.current = false
      setIsReplaying(false)
    }

    return { processed, failed }
  }, [refreshQueuedCount])

  const enqueueOrder = useCallback(
    async (payload: CreatePOSOrderPayload) => {
      const queued = await upsertPendingOrder({
        id: createQueueID(),
        payload,
        queuedAt: new Date().toISOString(),
        retryCount: 0,
      })
      await refreshQueuedCount()
      return queued
    },
    [refreshQueuedCount]
  )

  const clearQueue = useCallback(async () => {
    const cleared = await clearPendingOrders()
    await refreshQueuedCount()
    return cleared
  }, [refreshQueuedCount])

  useEffect(() => {
    if (typeof window === 'undefined') {
      return
    }

    const updateOfflineState = () => {
      if (!canUseNavigatorOnline()) {
        setIsOffline(false)
        return
      }
      setIsOffline(!navigator.onLine)
    }

    const handleOnline = () => {
      updateOfflineState()
      void replayQueuedOrders()
    }

    updateOfflineState()
    void refreshQueuedCount()

    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', updateOfflineState)

    return () => {
      window.removeEventListener('online', handleOnline)
      window.removeEventListener('offline', updateOfflineState)
    }
  }, [refreshQueuedCount, replayQueuedOrders])

  return {
    isOffline,
    queuedOrderCount,
    isReplaying,
    enqueueOrder,
    replayQueuedOrders,
    clearQueue,
  }
}

export default usePOSOffline
