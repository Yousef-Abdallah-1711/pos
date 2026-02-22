'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import { initiateSubscriptionPayment } from '@/lib/api'
import { useTenantStatusStore } from '@/stores/tenantStatusStore'
import { toast } from 'sonner'

export type PaymentStatus = 'idle' | 'loading' | 'success' | 'failed'

interface UsePaymentOptions {
  onSuccess?: () => void | Promise<void>
  onFailure?: () => void | Promise<void>
}

interface HandlePaymentCompletionArgs {
  status: 'success' | 'failed'
  notify?: boolean
}

function normalizeIframeURL(rawURL: string) {
  if (!rawURL) return ''
  if (rawURL.startsWith('http://') || rawURL.startsWith('https://')) {
    return rawURL
  }

  const apiBaseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'
  const apiOrigin = new URL(apiBaseURL).origin
  if (rawURL.startsWith('/')) {
    return `${apiOrigin}${rawURL}`
  }
  return `${apiOrigin}/${rawURL}`
}

export function usePayment(options?: UsePaymentOptions) {
  const searchParams = useSearchParams()
  const checkStatus = useTenantStatusStore((state) => state.checkStatus)
  const handledQueryRef = useRef<string | null>(null)

  const [isModalOpen, setIsModalOpen] = useState(false)
  const [iframeURL, setIframeURL] = useState('')
  const [paymentStatus, setPaymentStatus] = useState<PaymentStatus>('idle')
  const [error, setError] = useState<string | null>(null)

  const handlePaymentCompletion = useCallback(
    async ({ status, notify = true }: HandlePaymentCompletionArgs) => {
      setPaymentStatus(status)
      setIsModalOpen(false)

      if (status === 'success') {
        await checkStatus()
        if (notify) {
          toast.success('Payment completed successfully')
        }
        if (options?.onSuccess) {
          await options.onSuccess()
        }
        return
      }

      if (notify) {
        toast.error('Payment failed')
      }
      if (options?.onFailure) {
        await options.onFailure()
      }
    },
    [checkStatus, options]
  )

  const initiatePayment = useCallback(async (planID: string) => {
    setPaymentStatus('loading')
    setError(null)

    try {
      const response = await initiateSubscriptionPayment(planID)
      const normalizedURL = normalizeIframeURL(response.iframe_url)
      setIframeURL(normalizedURL)
      setIsModalOpen(true)
      setPaymentStatus('idle')
    } catch (err: any) {
      const backendMessage = err?.response?.data?.error
      const message = backendMessage || 'Failed to initiate payment'
      setError(message)
      setPaymentStatus('failed')
      toast.error(message)
    }
  }, [])

  const closeModal = useCallback(() => {
    setIsModalOpen(false)
  }, [])

  const resetPaymentState = useCallback(() => {
    setPaymentStatus('idle')
    setError(null)
  }, [])

  useEffect(() => {
    const queryStatus = searchParams.get('status')
    if (!queryStatus) return
    if (queryStatus !== 'success' && queryStatus !== 'failed') return

    const queryKey = `${queryStatus}|${searchParams.get('order_id') || ''}`
    if (handledQueryRef.current === queryKey) return
    handledQueryRef.current = queryKey

    void handlePaymentCompletion({ status: queryStatus, notify: true })
  }, [searchParams, handlePaymentCompletion])

  return {
    initiatePayment,
    handlePaymentCompletion,
    closeModal,
    resetPaymentState,
    isModalOpen,
    setIsModalOpen,
    iframeURL,
    paymentStatus,
    isLoading: paymentStatus === 'loading',
    error,
  }
}

