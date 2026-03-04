'use client'

import { useMemo } from 'react'
import { buildPOSReceiptMarkup, type POSReceiptData } from '@/hooks/usePOSPrint'

interface POSReceiptProps {
  data: POSReceiptData | null
  receiptId?: string
}

export function POSReceipt({ data, receiptId = 'pos-receipt-template' }: POSReceiptProps) {
  const safeMarkup = useMemo(() => {
    if (!data) {
      return ''
    }
    const rawMarkup = buildPOSReceiptMarkup(data)
    if (typeof window === 'undefined') {
      return rawMarkup
    }

    const parser = new DOMParser()
    const doc = parser.parseFromString(rawMarkup, 'text/html')
    doc.querySelectorAll('script,iframe,object,embed,link,meta').forEach((el) => el.remove())
    doc.querySelectorAll('*').forEach((el) => {
      for (const attr of Array.from(el.attributes)) {
        const attrName = attr.name.toLowerCase()
        const attrValue = attr.value.toLowerCase()
        if (attrName.startsWith('on') || attrValue.includes('javascript:')) {
          el.removeAttribute(attr.name)
        }
      }
    })

    return doc.body.innerHTML
  }, [data])

  if (!data) {
    return null
  }

  return (
    <div
      id={receiptId}
      aria-hidden="true"
      style={{ display: 'none' }}
      dangerouslySetInnerHTML={{ __html: safeMarkup }}
    />
  )
}

export default POSReceipt
