'use client'

import { useCallback } from 'react'
import { debugError } from '@/lib/debug'

export interface POSReceiptLineItem {
  name: string
  quantity: number
  unitPrice: number
  totalPrice: number
  notes?: string
}

export interface POSReceiptData {
  locale: 'en' | 'ar'
  restaurantName: string
  orderNumber: string
  createdAtISO: string
  paymentMethod: string
  subtotal: number
  tax: number
  total: number
  items: POSReceiptLineItem[]
  tableName?: string | null
  customerName?: string | null
  cashierName?: string | null
}

interface UsePOSPrintOptions {
  cleanupDelayMs?: number
}

type ReceiptLocale = POSReceiptData['locale']

interface ReceiptLabels {
  date: string
  order: string
  table: string
  customer: string
  cashier: string
  payment: string
  subtotal: string
  tax: string
  total: string
  qty: string
  item: string
  amount: string
  footer: string
}

const RECEIPT_LABELS: Record<ReceiptLocale, ReceiptLabels> = {
  en: {
    date: 'Date',
    order: 'Order',
    table: 'Table',
    customer: 'Customer',
    cashier: 'Cashier',
    payment: 'Payment',
    subtotal: 'Subtotal',
    tax: 'Tax (14%)',
    total: 'TOTAL',
    qty: 'Qty',
    item: 'Item',
    amount: 'Amount',
    footer: 'Thank you! / شكراً لزيارتكم',
  },
  ar: {
    date: 'التاريخ',
    order: 'الطلب',
    table: 'الطاولة',
    customer: 'العميل',
    cashier: 'الكاشير',
    payment: 'الدفع',
    subtotal: 'الإجمالي الفرعي',
    tax: 'الضريبة (14%)',
    total: 'الإجمالي',
    qty: 'الكمية',
    item: 'الصنف',
    amount: 'المبلغ',
    footer: 'Thank you! / شكراً لزيارتكم',
  },
}

const PRINT_CSS = `
  @page {
    size: 80mm auto;
    margin: 4mm;
  }

  html, body {
    margin: 0;
    padding: 0;
    background: #fff;
    color: #000;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
    font-size: 12px;
    line-height: 1.35;
  }

  .receipt {
    width: 72mm;
    margin: 0 auto;
  }

  .receipt__center {
    text-align: center;
  }

  .receipt__title {
    margin: 0;
    font-size: 15px;
    font-weight: 700;
  }

  .receipt__rule {
    border-top: 1px dashed #000;
    margin: 6px 0;
  }

  .receipt__meta-row {
    display: flex;
    justify-content: space-between;
    gap: 8px;
  }

  .receipt__meta-row > span:first-child {
    opacity: 0.85;
  }

  .receipt__items {
    width: 100%;
    border-collapse: collapse;
    table-layout: fixed;
  }

  .receipt__items th,
  .receipt__items td {
    padding: 2px 0;
    vertical-align: top;
  }

  .receipt__items th {
    font-weight: 700;
    opacity: 0.95;
  }

  .receipt__qty {
    width: 16%;
  }

  .receipt__item {
    width: 56%;
    word-break: break-word;
  }

  .receipt__amount {
    width: 28%;
    text-align: end;
  }

  .receipt__note {
    display: block;
    font-size: 10px;
    opacity: 0.8;
    margin-top: 1px;
  }

  .receipt__totals {
    display: grid;
    gap: 2px;
  }

  .receipt__total-row {
    display: flex;
    justify-content: space-between;
    gap: 8px;
  }

  .receipt__total-row--grand {
    font-weight: 700;
    font-size: 13px;
  }

  .receipt__footer {
    text-align: center;
    margin-top: 6px;
  }

  @media print {
    body {
      -webkit-print-color-adjust: exact;
      print-color-adjust: exact;
    }
  }
`

const escapeHTML = (value: string) =>
  value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')

const toLocale = (locale: ReceiptLocale) =>
  locale === 'ar' ? 'ar-EG' : 'en-US'

const formatMoney = (amount: number, locale: ReceiptLocale) =>
  new Intl.NumberFormat(toLocale(locale), {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(Number.isFinite(amount) ? amount : 0)

const formatDateTime = (dateISO: string, locale: ReceiptLocale) => {
  const date = new Date(dateISO)
  if (Number.isNaN(date.getTime())) {
    return dateISO
  }
  return new Intl.DateTimeFormat(toLocale(locale), {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

const paymentLabel = (method: string, locale: ReceiptLocale) => {
  const normalized = method.trim().toLowerCase()
  const labels =
    locale === 'ar'
      ? { cash: 'نقداً', card: 'بطاقة', wallet: 'محفظة', online: 'أونلاين' }
      : { cash: 'Cash', card: 'Card', wallet: 'Wallet', online: 'Online' }
  return labels[normalized as keyof typeof labels] ?? method
}

export const buildPOSReceiptMarkup = (data: POSReceiptData) => {
  const locale = data.locale === 'ar' ? 'ar' : 'en'
  const labels = RECEIPT_LABELS[locale]
  const direction = locale === 'ar' ? 'rtl' : 'ltr'

  const safeRestaurant = escapeHTML(data.restaurantName || 'Restaurant')
  const safeOrder = escapeHTML(data.orderNumber || '-')
  const safeDate = escapeHTML(formatDateTime(data.createdAtISO, locale))
  const safePayment = escapeHTML(paymentLabel(data.paymentMethod, locale))
  const safeCashier = escapeHTML(data.cashierName || '-')

  const tableRow = data.tableName
    ? `<div class="receipt__meta-row"><span>${labels.table}:</span><span>${escapeHTML(
        data.tableName
      )}</span></div>`
    : ''

  const customerRow = data.customerName
    ? `<div class="receipt__meta-row"><span>${labels.customer}:</span><span>${escapeHTML(
        data.customerName
      )}</span></div>`
    : ''

  const itemRows = data.items
    .map((item) => {
      const safeName = escapeHTML(item.name)
      const safeNotes = item.notes
        ? `<span class="receipt__note">${escapeHTML(item.notes)}</span>`
        : ''

      return `
        <tr>
          <td class="receipt__qty">${Math.max(1, Math.floor(item.quantity))}</td>
          <td class="receipt__item">${safeName}${safeNotes}</td>
          <td class="receipt__amount">${escapeHTML(formatMoney(item.totalPrice, locale))}</td>
        </tr>
      `
    })
    .join('')

  return `
    <article class="receipt" dir="${direction}">
      <header class="receipt__center">
        <h1 class="receipt__title">${safeRestaurant}</h1>
      </header>

      <div class="receipt__rule"></div>

      <section>
        <div class="receipt__meta-row"><span>${labels.date}:</span><span>${safeDate}</span></div>
        <div class="receipt__meta-row"><span>${labels.order}:</span><span>#${safeOrder}</span></div>
        ${tableRow}
        ${customerRow}
      </section>

      <div class="receipt__rule"></div>

      <section>
        <table class="receipt__items" aria-label="receipt-items">
          <thead>
            <tr>
              <th class="receipt__qty">${labels.qty}</th>
              <th class="receipt__item">${labels.item}</th>
              <th class="receipt__amount">${labels.amount}</th>
            </tr>
          </thead>
          <tbody>
            ${itemRows}
          </tbody>
        </table>
      </section>

      <div class="receipt__rule"></div>

      <section class="receipt__totals">
        <div class="receipt__total-row"><span>${labels.subtotal}:</span><span>${escapeHTML(
          formatMoney(data.subtotal, locale)
        )}</span></div>
        <div class="receipt__total-row"><span>${labels.tax}:</span><span>${escapeHTML(
          formatMoney(data.tax, locale)
        )}</span></div>
        <div class="receipt__total-row receipt__total-row--grand"><span>${labels.total}:</span><span>${escapeHTML(
          formatMoney(data.total, locale)
        )}</span></div>
      </section>

      <div class="receipt__rule"></div>

      <section>
        <div class="receipt__meta-row"><span>${labels.payment}:</span><span>${safePayment}</span></div>
        <div class="receipt__meta-row"><span>${labels.cashier}:</span><span>${safeCashier}</span></div>
      </section>

      <div class="receipt__rule"></div>

      <footer class="receipt__footer">${labels.footer}</footer>
    </article>
  `
}

export const usePOSPrint = (options?: UsePOSPrintOptions) => {
  const cleanupDelayMs = options?.cleanupDelayMs ?? 5000

  const printReceipt = useCallback(
    async (data: POSReceiptData) => {
      if (typeof window === 'undefined') {
        return false
      }

      const iframe = document.createElement('iframe')
      iframe.setAttribute('aria-hidden', 'true')
      iframe.style.position = 'fixed'
      iframe.style.width = '0'
      iframe.style.height = '0'
      iframe.style.opacity = '0'
      iframe.style.pointerEvents = 'none'
      iframe.style.border = '0'

      const cleanup = () => {
        if (iframe.parentNode) {
          iframe.parentNode.removeChild(iframe)
        }
      }

      document.body.appendChild(iframe)

      try {
        const printWindow = iframe.contentWindow
        if (!printWindow) {
          cleanup()
          return false
        }

        const markup = buildPOSReceiptMarkup(data)
        const documentTitle = `${data.restaurantName || 'POS'} Receipt #${data.orderNumber}`

        printWindow.document.open()
        printWindow.document.write(`
          <!doctype html>
          <html lang="${data.locale}">
            <head>
              <meta charset="utf-8" />
              <meta name="viewport" content="width=device-width, initial-scale=1" />
              <title>${escapeHTML(documentTitle)}</title>
              <style>${PRINT_CSS}</style>
            </head>
            <body>
              ${markup}
            </body>
          </html>
        `)
        printWindow.document.close()

        await new Promise((resolve) => {
          window.setTimeout(resolve, 90)
        })

        let cleaned = false
        const cleanupOnce = () => {
          if (cleaned) {
            return
          }
          cleaned = true
          cleanup()
        }

        printWindow.onafterprint = cleanupOnce
        window.setTimeout(cleanupOnce, cleanupDelayMs)
        printWindow.focus()
        printWindow.print()
        return true
      } catch (error) {
        cleanup()
        debugError('[POS Print] Failed to print receipt', error)
        return false
      }
    },
    [cleanupDelayMs]
  )

  return {
    printReceipt,
  }
}

export default usePOSPrint
