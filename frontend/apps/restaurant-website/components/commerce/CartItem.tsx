'use client'

import Image from 'next/image'
import { Minus, Plus, Trash2 } from 'lucide-react'
import { formatPrice } from '@/lib/utils/format'
import type { CartItem as CartItemType } from '@/lib/store/cart-store'

interface CartItemProps {
  item: CartItemType
  locale: 'en' | 'ar'
  onIncrease: () => void
  onDecrease: () => void
  onRemove: () => void
}

export function CartItem({ item, locale, onIncrease, onDecrease, onRemove }: CartItemProps) {
  const isArabic = locale === 'ar'

  return (
    <article className={`flex gap-3 rounded-lg border border-[hsl(var(--theme-border))] bg-[hsl(var(--theme-surface))] p-3 ${isArabic ? 'flex-row-reverse text-right' : ''}`}>
      <div className="relative h-16 w-16 shrink-0 overflow-hidden rounded-md bg-[hsl(var(--theme-background))]">
        {item.image ? (
          <Image src={item.image} alt={item.name} fill className="object-cover" sizes="64px" />
        ) : null}
      </div>
      <div className="min-w-0 flex-1">
        <p className="truncate text-sm font-semibold text-[hsl(var(--theme-text))]">{item.name}</p>
        <p className="text-xs text-[hsl(var(--theme-text))]/70">{formatPrice(item.price, locale)}</p>
        {item.notes ? <p className="mt-1 line-clamp-1 text-xs text-[hsl(var(--theme-text))]/60">{item.notes}</p> : null}
        <div className={`mt-2 flex items-center gap-2 ${isArabic ? 'flex-row-reverse' : ''}`}>
          <button
            type="button"
            className="rounded border border-[hsl(var(--theme-border))] p-1.5"
            onClick={onDecrease}
            aria-label={isArabic ? 'تقليل الكمية' : 'Decrease quantity'}
          >
            <Minus className="h-3.5 w-3.5" />
          </button>
          <span className="min-w-6 text-center text-sm font-semibold">{item.quantity}</span>
          <button
            type="button"
            className="rounded border border-[hsl(var(--theme-border))] p-1.5"
            onClick={onIncrease}
            aria-label={isArabic ? 'زيادة الكمية' : 'Increase quantity'}
          >
            <Plus className="h-3.5 w-3.5" />
          </button>
          <span className="ms-auto text-sm font-semibold text-[hsl(var(--theme-primary))]">{formatPrice(item.total, locale)}</span>
          <button
            type="button"
            className="rounded p-1.5 text-red-600"
            onClick={onRemove}
            aria-label={isArabic ? 'حذف المنتج' : 'Remove item'}
          >
            <Trash2 className="h-4 w-4" />
          </button>
        </div>
      </div>
    </article>
  )
}
