'use client'

import { useCallback, useEffect, useRef, useState } from 'react'

export type KDSOfflineActionType = 'bump' | 'recall'

export interface KDSOfflineAction {
  id: string
  type: KDSOfflineActionType
  orderId: number
  queuedAt: string
}

interface ReplayResult {
  processed: number
  failed: number
}

const DB_NAME = 'kds-offline-db'
const DB_VERSION = 1
const STORE_NAME = 'pending-actions'

const canUseIndexedDB = () => (
  typeof window !== 'undefined' &&
  typeof window.indexedDB !== 'undefined'
)

const canUseNavigatorOnline = () => (
  typeof navigator !== 'undefined' &&
  typeof navigator.onLine === 'boolean'
)

const createActionID = () => `${Date.now()}-${Math.random().toString(36).slice(2, 10)}`

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

const readAllActions = async (): Promise<KDSOfflineAction[]> => {
  const db = await openQueueDB()
  if (!db) {
    return []
  }

  return new Promise((resolve) => {
    const tx = db.transaction(STORE_NAME, 'readonly')
    const store = tx.objectStore(STORE_NAME)
    const request = store.getAll()

    request.onsuccess = () => {
      const rows = (request.result || []) as KDSOfflineAction[]
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

const upsertAction = async (action: KDSOfflineAction): Promise<boolean> => {
  const db = await openQueueDB()
  if (!db) {
    return false
  }

  return new Promise((resolve) => {
    const tx = db.transaction(STORE_NAME, 'readwrite')
    const store = tx.objectStore(STORE_NAME)
    const request = store.put(action)

    request.onsuccess = () => resolve(true)
    request.onerror = () => resolve(false)
    tx.oncomplete = () => db.close()
    tx.onerror = () => db.close()
  })
}

const removeAction = async (id: string): Promise<boolean> => {
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

const clearActions = async (): Promise<boolean> => {
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

export const useKDSOffline = () => {
  const [isOffline, setIsOffline] = useState(false)
  const [queuedActionsCount, setQueuedActionsCount] = useState(0)
  const [isReplayingQueue, setIsReplayingQueue] = useState(false)
  const replayingRef = useRef(false)

  const refreshQueuedCount = useCallback(async () => {
    const actions = await readAllActions()
    setQueuedActionsCount(actions.length)
    return actions.length
  }, [])

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

    updateOfflineState()
    void refreshQueuedCount()

    window.addEventListener('online', updateOfflineState)
    window.addEventListener('offline', updateOfflineState)

    return () => {
      window.removeEventListener('online', updateOfflineState)
      window.removeEventListener('offline', updateOfflineState)
    }
  }, [refreshQueuedCount])

  const enqueueAction = useCallback(async (type: KDSOfflineActionType, orderId: number) => {
    const action: KDSOfflineAction = {
      id: createActionID(),
      type,
      orderId,
      queuedAt: new Date().toISOString(),
    }

    const queued = await upsertAction(action)
    await refreshQueuedCount()
    return queued
  }, [refreshQueuedCount])

  const replayQueuedActions = useCallback(async (
    executor: (action: KDSOfflineAction) => Promise<void>
  ): Promise<ReplayResult> => {
    if (replayingRef.current) {
      return { processed: 0, failed: 0 }
    }

    if (canUseNavigatorOnline() && !navigator.onLine) {
      return { processed: 0, failed: 0 }
    }

    replayingRef.current = true
    setIsReplayingQueue(true)

    let processed = 0
    let failed = 0

    try {
      const actions = await readAllActions()
      for (const action of actions) {
        try {
          await executor(action)
          const removed = await removeAction(action.id)
          if (!removed) {
            failed += 1
            break
          }
          processed += 1
        } catch {
          failed += 1
          break
        }
      }
    } finally {
      await refreshQueuedCount()
      replayingRef.current = false
      setIsReplayingQueue(false)
    }

    return { processed, failed }
  }, [refreshQueuedCount])

  const clearQueuedActions = useCallback(async () => {
    const cleared = await clearActions()
    await refreshQueuedCount()
    return cleared
  }, [refreshQueuedCount])

  return {
    isOffline,
    queuedActionsCount,
    isReplayingQueue,
    enqueueAction,
    replayQueuedActions,
    clearQueuedActions,
  }
}
