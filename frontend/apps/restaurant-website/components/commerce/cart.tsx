'use client'

import { useCartStore } from '@/lib/store/cart-store'
import { Button } from '@pos-saas/ui'
import { Minus, Plus, Trash2, ShoppingCart } from 'lucide-react'
import Image from 'next/image'
import Link from 'next/link'

interface CartProps {
  locale?: 'en' | 'ar'
  themeColors?: Partial<{
    background: string
    text: string
    border: string
    primary: string
  }>
}

/**
 * Shopping cart component
 * Uses theme CSS variables for colors
 */
export function Cart({ locale = 'en', themeColors = {} }: CartProps) {
  const items = useCartStore((state) => state.items)
  const removeItem = useCartStore((state) => state.removeItem)
  const updateQuantity = useCartStore((state) => state.updateQuantity)
  const updateNotes = useCartStore((state) => state.updateNotes)
  const getTotal = useCartStore((state) => state.getTotal)
  const getTotalItems = useCartStore((state) => state.getTotalItems)
  const clearCart = useCartStore((state) => state.clearCart)
  const isRTL = locale === 'ar'
  const total = getTotal()
  const itemCount = getTotalItems()

  const bgColor = themeColors.background || '#ffffff'
  const textColor = themeColors.text || '#1f2937'
  const borderColor = themeColors.border || '#e5e7eb'
  const primaryColor = themeColors.primary || '#f97316'

  if (items.length === 0) {
    return (
      <div className="text-center py-12">
        <ShoppingCart className="cart-empty-icon h-16 w-16 mx-auto mb-4" />
        <h3 className="cart-empty-title text-xl font-semibold mb-2">{locale === 'en' ? 'Your cart is empty' : 'Ø³Ù„ØªÙƒ ÙØ§Ø±ØºØ©'}</h3>
        <p className="cart-empty-subtitle mb-6">
          {locale === 'en' ? 'Add some delicious items to get started' : 'Ø£Ø¶Ù Ø¨Ø¹Ø¶ Ø§Ù„Ø¹Ù†Ø§ØµØ± Ø§Ù„Ù„Ø°ÙŠØ°Ø© Ù„Ù„Ø¨Ø¯Ø¡'}
        </p>
        <Link href={`/${locale}/menu`}>
          <Button variant="default">{locale === 'en' ? 'Continue Shopping' : 'Ù…ØªØ§Ø¨Ø¹Ø© Ø§Ù„ØªØ³ÙˆÙ‚'}</Button>
        </Link>

        <style jsx>{`
          .cart-empty-icon {
            color: hsl(var(--theme-text) / 0.4);
          }

          .cart-empty-title {
            color: ${textColor};
          }

          .cart-empty-subtitle {
            color: hsl(var(--theme-text) / 0.6);
          }
        `}</style>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="cart-shell rounded-lg shadow-md overflow-hidden">
        <div>
          {items.map((item) => (
            <div
              key={item.id}
              className={`cart-item-row p-6 flex gap-4 ${isRTL ? 'flex-row-reverse' : ''}`}
            >
              {item.image ? (
                <div className="cart-item-media flex-shrink-0 w-24 h-24 rounded-lg overflow-hidden">
                  <Image
                    src={item.image}
                    alt={item.name}
                    width={96}
                    height={96}
                    className="w-full h-full object-cover"
                  />
                </div>
              ) : (
                <div className="cart-item-media flex-shrink-0 w-24 h-24 rounded-lg flex items-center justify-center">
                  <ShoppingCart className="cart-item-fallback-icon h-8 w-8" />
                </div>
              )}

              <div className="flex-1">
                <h3 className="cart-item-name font-semibold text-lg mb-1">{item.name}</h3>
                <p className="cart-item-price text-sm mb-4">
                  ${item.price.toFixed(2)} {locale === 'en' ? 'each' : 'Ù„ÙƒÙ„ ÙˆØ§Ø­Ø¯'}
                </p>

                <div className="flex items-center gap-3 mb-4">
                  <button
                    onClick={() => updateQuantity(item.productId, Math.max(1, item.quantity - 1))}
                    className="cart-qty-btn p-1 rounded transition-colors"
                  >
                    <Minus className="h-4 w-4" />
                  </button>
                  <span className="cart-qty-value px-4 py-1 rounded font-semibold min-w-[3rem] text-center">
                    {item.quantity}
                  </span>
                  <button
                    onClick={() => updateQuantity(item.productId, item.quantity + 1)}
                    className="cart-qty-btn p-1 rounded transition-colors"
                  >
                    <Plus className="h-4 w-4" />
                  </button>
                </div>

                <textarea
                  value={item.notes || ''}
                  onChange={(e) => updateNotes(item.productId, e.target.value)}
                  placeholder={locale === 'en' ? 'Add special notes...' : 'Ø£Ø¶Ù Ù…Ù„Ø§Ø­Ø¸Ø§Øª Ø®Ø§ØµØ©...'}
                  className="cart-notes w-full p-2 text-sm border rounded resize-none focus:outline-none"
                  rows={2}
                />
              </div>

              <div className={`flex flex-col items-end gap-4 ${isRTL ? 'items-start' : ''}`}>
                <div className={isRTL ? 'text-left' : 'text-right'}>
                  <p className="cart-total-label text-sm">{locale === 'en' ? 'Total' : 'Ø§Ù„Ø¥Ø¬Ù…Ø§Ù„ÙŠ'}</p>
                  <p className="cart-total-value text-2xl font-bold">${item.total.toFixed(2)}</p>
                </div>
                <button onClick={() => removeItem(item.productId)} className="cart-delete-btn p-2 rounded transition-colors">
                  <Trash2 className="h-5 w-5" />
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="cart-summary rounded-lg shadow-md p-6 space-y-4">
        <div className="space-y-2">
          <div className={`flex justify-between ${isRTL ? 'flex-row-reverse' : ''} cart-summary-row`}>
            <span>{locale === 'en' ? 'Subtotal' : 'Ø§Ù„Ù…Ø¬Ù…ÙˆØ¹'}</span>
            <span>${total.toFixed(2)}</span>
          </div>
          <div className={`flex justify-between ${isRTL ? 'flex-row-reverse' : ''} cart-summary-row`}>
            <span>{locale === 'en' ? 'Delivery' : 'Ø§Ù„ØªÙˆØµÙŠÙ„'}</span>
            <span>${(total > 50 ? 0 : 5).toFixed(2)}</span>
          </div>
          {total > 50 && (
            <p className="cart-free-delivery text-sm">
              {locale === 'en' ? 'Free delivery on orders over $50!' : 'ØªÙˆØµÙŠÙ„ Ù…Ø¬Ø§Ù†ÙŠ Ø¹Ù„Ù‰ Ø§Ù„Ø·Ù„Ø¨Ø§Øª ÙÙˆÙ‚ 50 Ø¯ÙˆÙ„Ø§Ø±!'}
            </p>
          )}
        </div>

        <div className="cart-summary-total pt-4">
          <div className={`flex justify-between text-lg font-bold mb-6 ${isRTL ? 'flex-row-reverse' : ''}`}>
            <span className="cart-total-title">{locale === 'en' ? 'Total' : 'Ø§Ù„Ø¥Ø¬Ù…Ø§Ù„ÙŠ'}</span>
            <span className="cart-total-primary">${(total + (total > 50 ? 0 : 5)).toFixed(2)}</span>
          </div>

          <div className="space-y-3">
            <Button asChild className="w-full" size="lg">
              <Link href={`/${locale}/checkout`}>
                {locale === 'en' ? 'Proceed to Checkout' : 'Ù…ØªØ§Ø¨Ø¹Ø© Ø§Ù„Ø¯ÙØ¹'}
              </Link>
            </Button>
            <Button asChild variant="outline" className="w-full" size="lg">
              <Link href={`/${locale}/menu`}>
                {locale === 'en' ? 'Continue Shopping' : 'Ù…ØªØ§Ø¨Ø¹Ø© Ø§Ù„ØªØ³ÙˆÙ‚'}
              </Link>
            </Button>
          </div>
        </div>

        <button
          onClick={() => {
            if (window.confirm(locale === 'en' ? 'Clear cart?' : 'Ù…Ø³Ø­ Ø§Ù„Ø³Ù„Ø©ØŸ')) {
              clearCart()
            }
          }}
          className="cart-clear-btn w-full font-semibold py-2 rounded transition-colors"
        >
          {locale === 'en' ? 'Clear Cart' : 'Ù…Ø³Ø­ Ø§Ù„Ø³Ù„Ø©'}
        </button>
      </div>

      <style jsx>{`
        .cart-shell,
        .cart-summary {
          background-color: ${bgColor};
        }

        .cart-item-row {
          border-bottom: 1px solid ${borderColor};
        }

        .cart-item-media {
          background-color: hsl(var(--theme-background) / 0.5);
        }

        .cart-item-fallback-icon {
          color: hsl(var(--theme-text) / 0.4);
        }

        .cart-item-name {
          color: ${textColor};
        }

        .cart-item-price {
          color: hsl(var(--theme-text) / 0.7);
        }

        .cart-qty-btn {
          color: ${textColor};
        }

        .cart-qty-btn:hover {
          background-color: hsl(var(--theme-text) / 0.05);
        }

        .cart-qty-value {
          background-color: hsl(var(--theme-background) / 0.8);
          color: ${textColor};
          border: 1px solid ${borderColor};
        }

        .cart-notes {
          background-color: hsl(var(--theme-background) / 0.8);
          color: ${textColor};
          border-color: ${borderColor};
        }

        .cart-notes:focus {
          border-color: ${primaryColor};
          box-shadow: 0 0 0 2px ${primaryColor}33;
        }

        .cart-total-label,
        .cart-summary-row {
          color: hsl(var(--theme-text) / 0.6);
        }

        .cart-total-value,
        .cart-total-primary {
          color: ${primaryColor};
        }

        .cart-delete-btn {
          color: #ef4444;
        }

        .cart-delete-btn:hover {
          background-color: #fee2e2;
        }

        .cart-free-delivery {
          color: #16a34a;
        }

        .cart-summary-total {
          border-top: 1px solid ${borderColor};
        }

        .cart-total-title {
          color: ${textColor};
        }

        .cart-clear-btn {
          color: #dc2626;
        }

        .cart-clear-btn:hover {
          background-color: #fee2e2;
        }
      `}</style>
    </div>
  )
}
