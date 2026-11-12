'use client'

import Link from 'next/link'
import { useCart } from '@/lib/hooks/use-cart'
import { Cart } from '@/components/commerce/cart'
import { createTranslator } from '@/lib/translations'
import { Button } from '@pos-saas/ui'
import { ShoppingBag } from 'lucide-react'
import type { ServerThemeData } from '@/lib/api/get-theme-server'

interface CartPageContentProps {
  locale: 'en' | 'ar'
  themeData: ServerThemeData
}

export function CartPageContent({ locale, themeData }: CartPageContentProps) {
  const isRTL = locale === 'ar'
  const t = createTranslator(locale)

  const primaryColor = themeData?.colors?.primary || '#f97316'
  const secondaryColor = themeData?.colors?.secondary || '#0ea5e9'

  const { items: cartItems, totalItems: cartItemsCount, totalPrice: orderSubtotal, tableId, tableName } = useCart()
  const isDineIn = tableId !== null
  const deliveryFee = isDineIn ? 0 : orderSubtotal > 50 ? 0 : 5.99
  const grandTotal = orderSubtotal + deliveryFee

  if (cartItems.length === 0) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className={`text-center max-w-md ${isRTL ? 'text-right' : ''}`}>
          <ShoppingBag className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h1 className="text-3xl font-bold mb-2">{t('cart.empty')}</h1>
          <p className="text-gray-600 dark:text-gray-400 mb-6">{t('cart.browseMenu')}</p>
          <Link href={`/${locale}/menu`}>
            <Button className="w-full">{t('navigation.menu')}</Button>
          </Link>
        </div>
      </div>
    )
  }

  return (
    <>
      <div className="cart-page-header text-white py-12">
        <div className={`max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 ${isRTL ? 'text-right' : ''}`}>
          <h1 className="text-4xl font-bold mb-2">{t('cart.title')}</h1>
          <p className="text-gray-100">
            {cartItemsCount} {t('cart.items')}
          </p>
        </div>
      </div>

      <div className={`max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 ${isRTL ? 'text-right' : ''}`}>
        {isDineIn && (
          <div
            className={`mb-6 rounded-lg border px-4 py-3 text-sm font-medium bg-orange-50 border-orange-200 text-orange-900 dark:bg-orange-900/20 dark:border-orange-700 dark:text-orange-200 ${
              isRTL ? 'text-right' : ''
            }`}
          >
            {t('tableOrdering.dineInOrder')} - {tableName || `${t('tableOrdering.table')} #${tableId}`}
          </div>
        )}

        <div className={`grid grid-cols-1 lg:grid-cols-3 gap-8 ${isRTL ? 'lg:grid-flow-col-dense' : ''}`}>
          <div className={`lg:col-span-2 ${isRTL ? 'lg:order-2' : ''}`}>
            <Cart locale={locale} themeColors={themeData?.colors} />
          </div>

          <div className={`lg:col-span-1 ${isRTL ? 'lg:order-1' : ''}`}>
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 lg:sticky lg:top-6">
              <h2 className="text-2xl font-bold mb-6">{t('checkout.orderSummary')}</h2>

              <div className={`space-y-2 mb-6 pb-6 border-b border-gray-200 dark:border-gray-700 ${isRTL ? 'text-right' : ''}`}>
                <div className={`flex justify-between ${isRTL ? 'flex-row-reverse' : ''}`}>
                  <span className="text-gray-600 dark:text-gray-400">{t('cart.items')}</span>
                  <span className="font-semibold">{cartItemsCount}</span>
                </div>
              </div>

              <div className={`space-y-3 mb-6 ${isRTL ? 'text-right' : ''}`}>
                <div className={`flex justify-between ${isRTL ? 'flex-row-reverse' : ''}`}>
                  <span className="text-gray-600 dark:text-gray-400">{t('common.subtotal')}</span>
                  <span className="font-semibold">${orderSubtotal.toFixed(2)}</span>
                </div>

                <div className={`flex justify-between ${isRTL ? 'flex-row-reverse' : ''}`}>
                  <span className="text-gray-600 dark:text-gray-400">{t('cart.deliveryFee')}</span>
                  <span className="text-sm text-green-600">
                    {deliveryFee === 0 ? t('common.free') : `$${deliveryFee.toFixed(2)}`}
                  </span>
                </div>

                <div className={`flex justify-between text-lg border-t border-gray-200 dark:border-gray-700 pt-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
                  <span className="font-bold">{t('common.total')}</span>
                  <span className="cart-total-value text-2xl font-bold">${grandTotal.toFixed(2)}</span>
                </div>
              </div>

              <Button asChild className="cart-checkout-button w-full py-3 text-lg">
                <Link href={`/${locale}/checkout`} className="block w-full mb-3">{t('cart.proceedToCheckout')}</Link>
              </Button>

              <Button asChild variant="outline" className="w-full py-3">
                <Link href={`/${locale}/menu`}>
                  {t('cart.continueShopping')}
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </div>

      <style jsx>{`
        .cart-page-header {
          background: linear-gradient(135deg, ${primaryColor} 0%, ${secondaryColor} 100%);
        }

        .cart-total-value {
          color: ${primaryColor};
        }

        .cart-checkout-button {
          background-color: ${primaryColor};
        }
      `}</style>
    </>
  )
}
