'use client'

import { useState, useEffect, memo } from 'react'
import { Button } from '@pos-saas/ui'
import { ShoppingCart, Plus, Minus } from 'lucide-react'
import Image from 'next/image'
import { useThemeColors } from '@/lib/hooks/use-theme'
import { getImageConfig, optimizeImageUrl } from '@/lib/image-optimization'

export interface Product {
  id: number
  name?: string
  name_en?: string
  description?: string
  description_en?: string
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
  const themeColors = useThemeColors() as any

  useEffect(() => {
    if (isAdded) {
      const timeoutId = setTimeout(() => setIsAdded(false), 2000)
      return () => clearTimeout(timeoutId)
    }
  }, [isAdded])

  const handleAddToCart = () => {
    onAddToCart?.(product, quantity)
    setIsAdded(true)
  }

  const quantityInStock = Number(product.quantity_in_stock ?? 0)
  const isOutOfStock = product.track_inventory && quantityInStock <= 0
  const isAvailable = product.isAvailable !== false && product.is_available !== false && !isOutOfStock
  const unavailableLabel = locale === 'en' ? 'Unavailable' : 'ØºÙŠØ± Ù…ØªÙˆÙØ±'
  const outOfStockLabel = locale === 'en' ? 'Out of Stock' : 'Ù†ÙØ¯ Ø§Ù„Ù…Ø®Ø²ÙˆÙ†'
  const badgeLabel = isOutOfStock ? outOfStockLabel : unavailableLabel

  const bgColor = themeColors?.background || '#ffffff'
  const textColor = themeColors?.text || '#1f2937'
  const primaryColor = themeColors?.primary || '#f97316'
  const borderColor = themeColors?.border || '#e5e7eb'
  const accentColor = themeColors?.accent || '#fbbf24'

  return (
    <div className="product-card rounded-lg overflow-hidden group shadow-sm hover:shadow-md transition-shadow duration-200">
      <div className="product-image-wrap relative h-48 overflow-hidden">
        {product.image || product.main_image_url ? (
          <Image
            src={optimizeImageUrl(product.image || product.main_image_url)}
            alt={product.name || product.name_en || 'Product'}
            fill
            loading="lazy"
            quality={80}
            sizes={getImageConfig('productCard').sizes}
            className="object-cover group-hover:scale-105 transition-transform duration-300"
          />
        ) : (
          <div className="product-image-fallback w-full h-full flex items-center justify-center">
            <svg className="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
          </div>
        )}

        {!isAvailable && (
          <div className="product-unavailable absolute inset-0 flex items-center justify-center z-10">
            <span className="product-unavailable-text">{badgeLabel}</span>
          </div>
        )}

        {product.rating && (
          <div className="product-rating absolute top-2 right-2 px-2 py-1 rounded-full text-sm font-semibold flex items-center gap-1">
            â­ {product.rating.toFixed(1)}
          </div>
        )}
      </div>

      <div className="p-4">
        <h3 className="product-title text-lg font-semibold mb-1 line-clamp-2">{product.name || product.name_en || 'Product'}</h3>

        {product.category && <p className="product-category text-xs mb-2">{product.category}</p>}

        <p className="product-description text-sm mb-4 line-clamp-2">{product.description || product.description_en || ''}</p>

        <div className="flex items-center justify-between">
          <span className="product-price text-2xl font-bold">${product.price.toFixed(2)}</span>

          {isAvailable ? (
            <div className="flex items-center gap-2">
              <div className="product-qty-wrap flex items-center gap-1 rounded-lg">
                <button
                  onClick={() => setQuantity(Math.max(1, quantity - 1))}
                  className="product-qty-btn p-1 transition-colors hover:bg-gray-100 dark:hover:bg-gray-700"
                  disabled={quantity === 1}
                >
                  <Minus className="h-4 w-4" />
                </button>
                <span className="product-qty-value px-2 min-w-[2rem] text-center text-sm font-semibold">{quantity}</span>
                <button
                  onClick={() => setQuantity(quantity + 1)}
                  className="product-qty-btn p-1 transition-colors hover:bg-gray-100 dark:hover:bg-gray-700"
                >
                  <Plus className="h-4 w-4" />
                </button>
              </div>

              <Button
                size="sm"
                variant={isAdded ? 'default' : 'outline'}
                onClick={handleAddToCart}
                className={`product-add-button ${isAdded ? 'is-added' : 'is-outline'}`}
              >
                {isAdded ? <>âœ“</> : <ShoppingCart className="h-4 w-4" />}
              </Button>
            </div>
          ) : (
            <Button disabled variant="outline" size="sm">
              {badgeLabel}
            </Button>
          )}
        </div>
      </div>

      <style jsx>{`
        .product-card {
          background-color: ${bgColor};
          border: 1px solid ${borderColor};
        }

        .product-image-wrap {
          background-color: hsl(var(--theme-background) / 0.5);
        }

        .product-image-fallback {
          color: hsl(var(--theme-text) / 0.3);
        }

        .product-unavailable {
          background-color: rgba(0, 0, 0, 0.5);
        }

        .product-unavailable-text {
          color: #ffffff;
          font-weight: 600;
        }

        .product-rating {
          background-color: ${accentColor};
          color: #1f2937;
        }

        .product-title {
          color: ${textColor};
        }

        .product-category {
          color: hsl(var(--theme-text) / 0.6);
        }

        .product-description {
          color: hsl(var(--theme-text) / 0.7);
        }

        .product-price {
          color: ${primaryColor};
        }

        .product-qty-wrap {
          background-color: hsl(var(--theme-background) / 0.8);
          border: 1px solid ${borderColor};
        }

        .product-qty-btn,
        .product-qty-value {
          color: ${textColor};
        }

        .product-add-button.is-added {
          background-color: #22c55e;
          color: #ffffff;
          border-color: #22c55e;
        }

        .product-add-button.is-outline {
          border-color: ${primaryColor};
          color: ${primaryColor};
        }
      `}</style>
    </div>
  )
}

export const ProductCard = memo(ProductCardComponent)
