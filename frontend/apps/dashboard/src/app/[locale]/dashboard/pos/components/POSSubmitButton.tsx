'use client'

import { Loader2 } from 'lucide-react'
import { selectPOSTotalItems, usePOSStore } from '@/stores/posStore'

interface POSSubmitButtonProps {
  t: (path: string, fallback?: string) => string
  isSubmitting: boolean
  onSubmit: () => void
}

export function POSSubmitButton({ t, isSubmitting, onSubmit }: POSSubmitButtonProps) {
  const totalItems = usePOSStore(selectPOSTotalItems)
  const disabled = totalItems <= 0 || isSubmitting

  return (
    <button
      type="button"
      disabled={disabled}
      onClick={onSubmit}
      className="inline-flex h-11 w-full items-center justify-center gap-2 rounded-lg bg-emerald-500 px-4 text-sm font-semibold text-white shadow transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-60"
    >
      {isSubmitting ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
      <span>
        {isSubmitting
          ? t('pos.processing', 'Processing...')
          : `${t('pos.placeOrder', 'Place Order')} (${totalItems})`}
      </span>
      {!isSubmitting && <span className="rounded bg-black/20 px-1.5 py-0.5 text-[10px]">Space</span>}
    </button>
  )
}

export default POSSubmitButton
