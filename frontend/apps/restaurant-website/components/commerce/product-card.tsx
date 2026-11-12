'use client'

import { memo, useMemo, useState } from 'react'
import Image from 'next/image'
import { Minus, Plus, ShoppingCart, Star } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { formatPrice } from '@/lib/utils/format'

export interface Product {
  id: number
  name?: string
  name_en?: string
  name_ar?: string
  description?: string
  description_en?: string
  description_ar?: string
  price: number
  image?: string
  main_image_url?: string
  category?: string
  category_id?: number
  categoryId?: number
  isAvailable?: boolean
  is_available?: boolean
  track_inventory?: boolean
  quantity_in_stock?: number
  rating?: number
}

interface ProductCardProps {
  product: Product
  onAddToCart?: (product: Product, quantity: number) => void
  locale?: 'en' | 'ar'
}

function ProductCardComponent({ product, onAddToCart, locale = 'en' }: ProductCardProps) {
  const [quantity, setQuantity] = useState(1)
  const [isAdded, setIsAdded] = useState(false)
  const isArabic = locale === 'ar'
  const imageSrc = product.image || product.main_image_url || ''
  const quantityInStock = Number(product.quantity_in_stock ?? 0)
  const isOutOfStock = Boolean(product.track_inventory && quantityInStock <= 0)
  const isAvailable = product.isAvailable !== false && product.is_available !== false && !isOutOfStock
  const maxQty = product.track_inventory && quantityInStock > 0 ? quantityInStock : 99

  const name = useMemo(() => {
    if (isArabic) {
      return product.name_ar || product.name || product.name_en || 'منتج'
    }
    return product.name_en || product.name || product.name_ar || 'Product'
  }, [isArabic, product])

  const description = useMemo(() => {
    if (isArabic) {
      return product.description_ar || product.description || product.description_en || ''
    }
    return product.description_en || product.description || product.description_ar || ''
  }, [isArabic, product])

  const handleAdd = () => {
    onAddToCart?.(product, quantity)
    setIsAdded(true)
    setTimeout(() => setIsAdded(false), 1200)
  }

  return (
    <article className="overflow-hidden rounded-xl border border-[hsl(var(--theme-border))] bg-[hsl(var(--theme-surface))] shadow-sm transition-transform hover:-translate-y-0.5">
      <div className="relative h-44 w-full overflow-hidden">
        {imageSrc ? (
          <Image
            src={imageSrc}
            alt={name}
            fill
            className={`object-cover transition-transform duration-300 hover:scale-105 ${!isAvailable ? 'grayscale' : ''}`}
            sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-[hsl(var(--theme-background))] text-sm text-[hsl(var(--theme-text))]/60">
            {isArabic ? 'لا توجد صورة' : 'No image'}
          </div>
        )}
        {!isAvailable ? (
          <span className="absolute left-3 top-3 rounded-full bg-red-600 px-3 py-1 text-xs font-semibold text-white">
            {isArabic ? 'غير متاح' : 'Unavailable'}
          </span>
        ) : null}
        {typeof product.rating === 'number' ? (
          <span className="absolute right-3 top-3 inline-flex items-center gap-1 rounded-full bg-black/65 px-2 py-1 text-xs text-white">
            <Star className="h-3 w-3 fill-current text-amber-400" />
            {product.rating.toFixed(1)}
          </span>
        ) : null}
      </div>

      <div className="space-y-3 p-4">
        <div className={isArabic ? 'text-right' : 'text-left'}>
          <h3 className="line-clamp-1 text-base font-semibold text-[hsl(var(--theme-text))]">{name}</h3>
          <p className="line-clamp-2 text-sm text-[hsl(var(--theme-text))]/70">{description}</p>
        </div>

        <div className={`flex items-center justify-between ${isArabic ? 'flex-row-reverse' : ''}`}>
          <span className="rounded-full bg-[hsl(var(--theme-background))] px-3 py-1 text-sm font-bold text-[hsl(var(--theme-primary))]">
            {formatPrice(Number(product.price || 0), locale)}
          </span>
          <div className={`inline-flex items-center rounded-md border border-[hsl(var(--theme-border))] ${isArabic ? 'flex-row-reverse' : ''}`}>
            <button
              type="button"
              onClick={() => setQuantity((value) => Math.max(1, value - 1))}
              className="p-2 text-[hsl(var(--theme-text))]/80"
              aria-label={isArabic ? 'تقليل الكمية' : 'Decrease quantity'}
            >
              <Minus className="h-4 w-4" />
            </button>
            <span className="min-w-8 text-center text-sm font-semibold text-[hsl(var(--theme-text))]">{quantity}</span>
            <button
              type="button"
              onClick={() => setQuantity((value) => Math.min(maxQty, value + 1))}
              className="p-2 text-[hsl(var(--theme-text))]/80"
              aria-label={isArabic ? 'زيادة الكمية' : 'Increase quantity'}
              disabled={!isAvailable}
            >
              <Plus className="h-4 w-4" />
            </button>
          </div>
        </div>

        <Button
          type="button"
          onClick={handleAdd}
          disabled={!isAvailable}
          className="h-11 w-full bg-[hsl(var(--theme-primary))] text-white hover:bg-[hsl(var(--theme-primary))]/90"
        >
          {isAdded ? (
            isArabic ? 'تمت الإضافة' : 'Added'
          ) : (
            <span className="inline-flex items-center gap-2">
              <ShoppingCart className="h-4 w-4" />
              {isArabic ? 'أضف' : 'Add'}
            </span>
          )}
        </Button>
      </div>
    </article>
  )
}

export const ProductCard = memo(ProductCardComponent)

