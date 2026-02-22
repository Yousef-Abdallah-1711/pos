'use client'

import { useEffect, useMemo, useRef, useState } from 'react'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { ExternalLink, Loader2, RotateCcw } from 'lucide-react'

interface PaymentModalProps {
  isOpen: boolean
  iframeURL: string
  title?: string
  description?: string
  openNewTabLabel?: string
  refreshLabel?: string
  unavailableLabel?: string
  onClose: () => void
  onRefresh: () => void | Promise<void>
  onComplete: (status: 'success' | 'failed') => void
}

function extractStatusFromURL(value: string): 'success' | 'failed' | null {
  if (!value) return null
  try {
    const parsed = new URL(value, window.location.origin)
    const status = parsed.searchParams.get('status')
    if (status === 'success' || status === 'failed') {
      return status
    }
  } catch (error) {
    return null
  }
  return null
}

export function PaymentModal({
  isOpen,
  iframeURL,
  title = 'Complete Payment',
  description = 'Finish checkout in the secure payment frame below.',
  openNewTabLabel = 'Open in New Tab',
  refreshLabel = 'Refresh Status',
  unavailableLabel = 'Checkout URL is unavailable',
  onClose,
  onRefresh,
  onComplete,
}: PaymentModalProps) {
  const iframeRef = useRef<HTMLIFrameElement | null>(null)
  const [isFrameLoading, setIsFrameLoading] = useState(true)
  const syncStatusInURL = (status: 'success' | 'failed') => {
    if (typeof window === 'undefined') return
    try {
      const url = new URL(window.location.href)
      url.searchParams.set('status', status)
      if (!url.searchParams.get('order_id')) {
        url.searchParams.set('order_id', String(Date.now()))
      }
      window.history.replaceState({}, '', url.toString())
    } catch (error) {
      // Ignore URL synchronization failures; payment completion still proceeds.
    }
  }
  const trustedOrigins = useMemo(() => {
    const origins = new Set<string>()
    if (typeof window !== 'undefined') {
      origins.add(window.location.origin)
      try {
        if (iframeURL) {
          const parsed = new URL(iframeURL, window.location.origin)
          if (parsed.protocol === 'http:' || parsed.protocol === 'https:') {
            origins.add(parsed.origin)
          }
        }
      } catch (error) {
        // Invalid iframe URL should not add any additional trusted origin.
      }
    }
    return origins
  }, [iframeURL])

  useEffect(() => {
    if (!isOpen) return
    setIsFrameLoading(true)
  }, [isOpen, iframeURL])

  useEffect(() => {
    if (!isOpen) return

    const handleMessage = (event: MessageEvent) => {
      const frameWindow = iframeRef.current?.contentWindow
      if (!frameWindow || event.source !== frameWindow) return
      if (!trustedOrigins.has(event.origin)) return

      const payload = event.data
      if (!payload) return

      if (typeof payload === 'object') {
        const type = String(payload.type || '').toLowerCase()
        const status = String(payload.status || '').toLowerCase()
        if (type.includes('payment') && (status === 'success' || status === 'failed')) {
          syncStatusInURL(status)
          onComplete(status)
          return
        }
      }

      if (typeof payload === 'string') {
        const lowerPayload = payload.toLowerCase()
        if (lowerPayload.includes('payment_success') || lowerPayload.includes('paymob:success')) {
          syncStatusInURL('success')
          onComplete('success')
        } else if (lowerPayload.includes('payment_failed') || lowerPayload.includes('paymob:failed')) {
          syncStatusInURL('failed')
          onComplete('failed')
        }
      }
    }

    window.addEventListener('message', handleMessage)
    return () => window.removeEventListener('message', handleMessage)
  }, [isOpen, onComplete, trustedOrigins])

  const frameLoadHandler = () => {
    setIsFrameLoading(false)

    try {
      const currentURL = iframeRef.current?.contentWindow?.location?.href || ''
      const status = extractStatusFromURL(currentURL)
      if (status) {
        syncStatusInURL(status)
        onComplete(status)
      }
    } catch (error) {
      // Cross-origin iframe in live mode; postMessage listener handles completion there.
    }
  }

  const isReady = useMemo(() => Boolean(iframeURL), [iframeURL])

  return (
    <Dialog open={isOpen} onOpenChange={(open) => { if (!open) onClose() }}>
      <DialogContent className="w-[95vw] max-w-[95vw] h-[92vh] max-h-[92vh] p-4 sm:p-6">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>

        <div className="relative flex-1 rounded-md border border-gray-200 dark:border-gray-700 overflow-hidden min-h-[65vh]">
          {isReady ? (
            <>
              {isFrameLoading && (
                <div className="absolute inset-0 z-10 flex items-center justify-center bg-white/85 dark:bg-gray-900/85">
                  <div className="flex items-center gap-2 text-gray-600 dark:text-gray-300">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    <span>Loading checkout...</span>
                  </div>
                </div>
              )}
              <iframe
                ref={iframeRef}
                src={iframeURL}
                title="payment-checkout"
                className="w-full h-full min-h-[65vh] bg-white"
                allow="payment *"
                onLoad={frameLoadHandler}
              />
            </>
          ) : (
            <div className="h-full min-h-[65vh] flex items-center justify-center text-gray-500 dark:text-gray-400">
              {unavailableLabel}
            </div>
          )}
        </div>

        <DialogFooter className="flex items-center gap-2">
          <Button
            variant="outline"
            onClick={() => window.open(iframeURL, '_blank')}
            disabled={!isReady}
            className="gap-2"
          >
            <ExternalLink className="h-4 w-4" />
            {openNewTabLabel}
          </Button>
          <Button variant="outline" onClick={() => { void onRefresh() }} className="gap-2">
            <RotateCcw className="h-4 w-4" />
            {refreshLabel}
          </Button>
          <Button variant="secondary" onClick={onClose}>
            Cancel
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
