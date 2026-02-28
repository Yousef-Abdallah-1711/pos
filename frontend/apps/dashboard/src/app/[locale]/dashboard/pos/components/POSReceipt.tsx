'use client'

import { buildPOSReceiptMarkup, type POSReceiptData } from '@/hooks/usePOSPrint'

interface POSReceiptProps {
  data: POSReceiptData | null
  receiptId?: string
}

export function POSReceipt({ data, receiptId = 'pos-receipt-template' }: POSReceiptProps) {
  if (!data) {
    return null
  }

  return (
    <div
      id={receiptId}
      aria-hidden="true"
      style={{ display: 'none' }}
      dangerouslySetInnerHTML={{ __html: buildPOSReceiptMarkup(data) }}
    />
  )
}

export default POSReceipt
