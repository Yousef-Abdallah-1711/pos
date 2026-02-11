/**
 * Service Worker Provider Component
 *
 * Wraps app with service worker management and update notifications
 */

'use client'

import { useServiceWorker } from '@/hooks/useServiceWorker'
import { AlertCircle, Download, RefreshCw } from 'lucide-react'
import { useEffect, useState } from 'react'
import { debugLog } from '@/lib/debug'

export function ServiceWorkerProvider({ children }: { children: React.ReactNode }) {
  const [showUpdatePrompt, setShowUpdatePrompt] = useState(false)

  const sw = useServiceWorker({
    onUpdateAvailable: () => {
      setShowUpdatePrompt(true)
    },
    onInstalled: () => {
      debugLog('✅ Service Worker installed successfully')
    },
    enableAutoUpdate: true,
    checkInterval: 60000, // Check every minute
  })

  useEffect(() => {
    if (typeof window !== 'undefined' && 'serviceWorker' in navigator) {
      // Log status
      debugLog('[SW] Ready:', sw.isReady)
      debugLog('[SW] Offline:', sw.isOffline)
    }
  }, [sw.isReady, sw.isOffline])

  return (
    <>
      {children}

      {/* Update Available Notification */}
      {showUpdatePrompt && (
        <UpdatePrompt
          onUpdate={() => {
            sw.applyUpdate()
            setShowUpdatePrompt(false)
          }}
          onDismiss={() => setShowUpdatePrompt(false)}
        />
      )}

      {/* Offline Notification */}
      {sw.isOffline && (
        <OfflineIndicator />
      )}

      {/* Online Notification */}
      {sw.isReady && !sw.isOffline && (
        <OnlineIndicator />
      )}
    </>
  )
}

/**
 * Update Prompt Component
 */
function UpdatePrompt({
  onUpdate,
  onDismiss,
}: {
  onUpdate: () => void
  onDismiss: () => void
}) {
  return (
    <div className="fixed bottom-4 right-4 max-w-md z-50 animate-in slide-in-from-bottom-4">
      <div className="bg-blue-50 border border-blue-200 rounded-lg shadow-lg p-4">
        <div className="flex items-start gap-3">
          <div className="flex-shrink-0 mt-0.5">
            <Download className="h-5 w-5 text-blue-600" />
          </div>
          <div className="flex-1">
            <h3 className="font-semibold text-blue-900 text-sm">
              Update Available
            </h3>
            <p className="text-blue-700 text-xs mt-1">
              A new version of the app is available. Restart to get the latest features and fixes.
            </p>
            <div className="flex gap-2 mt-3">
              <button
                onClick={onUpdate}
                className="inline-flex items-center gap-1 px-3 py-2 bg-blue-600 text-white text-xs font-medium rounded hover:bg-blue-700 transition"
              >
                <RefreshCw className="h-3 w-3" />
                Restart Now
              </button>
              <button
                onClick={onDismiss}
                className="inline-flex px-3 py-2 bg-blue-100 text-blue-600 text-xs font-medium rounded hover:bg-blue-200 transition"
              >
                Later
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

/**
 * Offline Indicator
 */
function OfflineIndicator() {
  return (
    <div className="fixed top-0 left-0 right-0 z-40 bg-red-600 text-white py-2 px-4 flex items-center justify-between animate-in slide-in-from-top-2">
      <div className="flex items-center gap-2">
        <AlertCircle className="h-5 w-5 flex-shrink-0" />
        <span className="text-sm font-medium">
          You're offline. Some features may be limited.
        </span>
      </div>
    </div>
  )
}

/**
 * Online Indicator (shows briefly when coming back online)
 */
function OnlineIndicator() {
  const [show, setShow] = useState(false)

  useEffect(() => {
    const handleOnline = () => {
      setShow(true)
      const timeout = setTimeout(() => setShow(false), 3000)
      return () => clearTimeout(timeout)
    }

    window.addEventListener('online', handleOnline)
    return () => window.removeEventListener('online', handleOnline)
  }, [])

  if (!show) return null

  return (
    <div className="fixed top-0 left-0 right-0 z-40 bg-green-600 text-white py-2 px-4 flex items-center justify-between animate-in slide-in-from-top-2">
      <div className="flex items-center gap-2">
        <div className="h-2 w-2 bg-white rounded-full animate-pulse" />
        <span className="text-sm font-medium">
          Back online. Syncing data...
        </span>
      </div>
    </div>
  )
}
