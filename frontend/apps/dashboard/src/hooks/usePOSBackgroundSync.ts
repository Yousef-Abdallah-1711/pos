'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import { debugError } from '@/lib/debug'

export const POS_ORDER_SYNC_TAG = 'pos-order-sync'

export interface POSBackgroundSyncMessage {
  type: 'pos-sync-complete' | 'pos-sync-error'
  processed?: number
  failed?: number
  reason?: string
}

interface UsePOSBackgroundSyncOptions {
  onSyncMessage?: (message: POSBackgroundSyncMessage) => void
}

type SyncRegistration = ServiceWorkerRegistration & {
  sync?: {
    register: (tag: string) => Promise<void>
  }
}

const hasSyncManager = (registration: ServiceWorkerRegistration | null): registration is SyncRegistration =>
  Boolean(
    registration &&
      (registration as SyncRegistration).sync &&
      typeof (registration as SyncRegistration).sync?.register === 'function'
  )

export const usePOSBackgroundSync = (options?: UsePOSBackgroundSyncOptions) => {
  const [syncSupported, setSyncSupported] = useState(false)
  const onSyncMessageRef = useRef(options?.onSyncMessage)

  useEffect(() => {
    onSyncMessageRef.current = options?.onSyncMessage
  }, [options?.onSyncMessage])

  useEffect(() => {
    if (typeof navigator === 'undefined' || !('serviceWorker' in navigator)) {
      setSyncSupported(false)
      return
    }

    let isMounted = true

    void (async () => {
      try {
        const registration = await navigator.serviceWorker.ready
        if (isMounted) {
          setSyncSupported(hasSyncManager(registration))
        }
      } catch {
        if (isMounted) {
          setSyncSupported(false)
        }
      }
    })()

    return () => {
      isMounted = false
    }
  }, [])

  useEffect(() => {
    if (typeof navigator === 'undefined' || !('serviceWorker' in navigator)) {
      return
    }

    const serviceWorker = navigator.serviceWorker
    if (!serviceWorker) {
      return
    }

    const handleMessage = (event: MessageEvent) => {
      const payload = event.data as POSBackgroundSyncMessage | undefined
      if (!payload || typeof payload !== 'object') {
        return
      }

      if (payload.type === 'pos-sync-complete' || payload.type === 'pos-sync-error') {
        onSyncMessageRef.current?.(payload)
      }
    }

    serviceWorker.addEventListener('message', handleMessage)
    return () => {
      serviceWorker.removeEventListener('message', handleMessage)
    }
  }, [])

  const registerOrderSync = useCallback(async () => {
    if (typeof navigator === 'undefined' || !('serviceWorker' in navigator)) {
      return false
    }

    try {
      const registration = await navigator.serviceWorker.ready
      if (!hasSyncManager(registration)) {
        setSyncSupported(false)
        return false
      }

      await registration.sync.register(POS_ORDER_SYNC_TAG)
      setSyncSupported(true)
      return true
    } catch (error) {
      debugError('[POS Background Sync] Failed to register sync event', error)
      return false
    }
  }, [])

  return {
    syncSupported,
    syncTag: POS_ORDER_SYNC_TAG,
    registerOrderSync,
  }
}

export default usePOSBackgroundSync
