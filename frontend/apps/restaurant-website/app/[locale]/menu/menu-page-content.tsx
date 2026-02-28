'use client'

import { useCallback, useEffect } from 'react'
import { ProductCard } from '@/components/commerce/product-card-simple'
import { useCartStore } from '@/lib/store/cart-store'
import { createTranslator } from '@/lib/translations'
import Link from 'next/link'
import { Loader } from 'lucide-react'
import { useProducts } from '@/lib/hooks/use-api-queries'
import type { ServerThemeData } from '@/lib/api/get-theme-server'
import type { TableStatusNotice } from '@/lib/table-ordering'

export interface TableContext {
  tableId: number
  tableNumber: number
  tableName?: string | null
  tableStatus?: string
}

interface MenuPageContentProps {
  locale: 'en' | 'ar'
  themeData: ServerThemeData
  restaurantSlug: string
  tableContext?: TableContext | null
  tableStatusNotice?: TableStatusNotice | null
}

export function MenuPageContent({
  locale,
  themeData,
  restaurantSlug,
  tableContext = null,
  tableStatusNotice = null,
}: MenuPageContentProps) {
  const isRTL = locale === 'ar'
  const t = createTranslator(locale)

  const primaryColor = themeData?.colors?.primary || '#f97316'
  const secondaryColor = themeData?.colors?.secondary || '#0ea5e9'

  const addToCart = useCartStore((state) => state.addItem)
  const setTableContext = useCartStore((state) => state.setTableContext)
  const cartItems = useCartStore((state) => state.getTotalItems())

  const { data: products = [], isLoading } = useProducts(restaurantSlug)

  useEffect(() => {
    if (!tableContext?.tableId) {
      return
    }

    const label = tableContext.tableName || `Table #${tableContext.tableNumber}`
    setTableContext(tableContext.tableId, label)
  }, [setTableContext, tableContext])

  const handleAddToCart = useCallback(
    (product: any, quantity: number) => {
      for (let i = 0; i < quantity; i++) {
        addToCart({
          id: product.id,
          name: product.name_en || product.name || 'Product',
          price: product.price || 0,
          image: product.main_image_url || product.image || '',
          track_inventory: product.track_inventory,
          quantity_in_stock: product.quantity_in_stock,
        })
      }
    },
    [addToCart]
  )

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <Loader className="h-12 w-12 text-primary mx-auto mb-4 animate-spin" />
          <p className="text-gray-600 dark:text-gray-400">{t('common.loading')}</p>
        </div>
      </div>
    )
  }

  return (
    <>
      <div className="menu-page-header text-white py-12">
        <div className={`max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 ${isRTL ? 'text-right' : ''}`}>
          <h1 className="text-4xl font-bold mb-2">{t('menu.title')}</h1>
          <p className="text-gray-100">{t('menu.subtitle')}</p>
        </div>
      </div>

      <div className={`max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 ${isRTL ? 'text-right' : ''}`}>
        {tableContext && (
          <div
            className={`mb-6 rounded-lg border px-4 py-3 text-sm font-medium bg-orange-50 border-orange-200 text-orange-900 dark:bg-orange-900/20 dark:border-orange-700 dark:text-orange-200 ${
              isRTL ? 'text-right' : ''
            }`}
          >
            {t('tableOrdering.dineInOrder')} - {t('tableOrdering.table')} #{tableContext.tableNumber}
          </div>
        )}

        {tableStatusNotice && (
          <div
            className={`mb-6 rounded-lg border px-4 py-3 text-sm font-medium ${
              tableStatusNotice.level === 'critical'
                ? 'border-red-200 bg-red-50 text-red-900 dark:border-red-800 dark:bg-red-900/20 dark:text-red-200'
                : 'border-amber-200 bg-amber-50 text-amber-900 dark:border-amber-700 dark:bg-amber-900/20 dark:text-amber-200'
            } ${isRTL ? 'text-right' : ''}`}
          >
            {t(tableStatusNotice.translationKey)}
          </div>
        )}

        <div className={`mb-6 ${isRTL ? 'text-right' : ''}`}>
          <h2 className="text-2xl font-semibold">{t('menu.allItems')}</h2>
          <p className="text-gray-600 dark:text-gray-400">
            {products.length} {t('menu.items')}
          </p>
        </div>

        {products.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
            {products.map((product: any) => (
              <ProductCard
                key={product.id}
                product={{
                  id: product.id,
                  name: product.name_en || product.name || '',
                  description: product.description_en || product.description || '',
                  price: product.price || 0,
                  image: product.main_image_url || product.image || '',
                  track_inventory: product.track_inventory,
                  quantity_in_stock: product.quantity_in_stock,
                  is_available: product.is_available,
                }}
                onAddToCart={handleAddToCart}
                locale={locale}
              />
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-gray-500">{t('menu.noItems')}</p>
          </div>
        )}

        {cartItems > 0 && (
          <div
            className={`menu-cart-banner text-white p-6 rounded-lg flex items-center justify-between ${
              isRTL ? 'flex-row-reverse' : ''
            }`}
          >
            <div>
              <p className="text-lg font-semibold">
                {cartItems} {t('menu.itemsInCart')}
              </p>
            </div>
            <Link href={`/${locale}/cart`}>
              <button className="menu-cart-banner-button px-6 py-2 rounded-lg font-semibold">{t('common.cart')}</button>
            </Link>
          </div>
        )}
      </div>

      <style jsx>{`
        .menu-page-header {
          background: linear-gradient(135deg, ${primaryColor} 0%, ${secondaryColor} 100%);
        }

        .menu-cart-banner {
          background-color: ${primaryColor};
        }

        .menu-cart-banner-button {
          background-color: #ffffff;
          color: ${primaryColor};
        }
      `}</style>
    </>
  )
}
