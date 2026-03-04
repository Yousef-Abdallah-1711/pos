'use client'

import Link from 'next/link'
import { ProductCard, type Product } from '@/components/commerce/product-card'
import { useCartStore } from '@/lib/store/cart-store'
import type { ThemeLocale } from '@/types/theme'

interface ProductGridProps {
  locale: ThemeLocale
  products: Product[]
}

export function ProductGrid({ locale, products }: ProductGridProps) {
  const addItem = useCartStore((state) => state.addItem)

  const handleAddToCart = (product: Product, quantity: number) => {
    for (let i = 0; i < quantity; i += 1) {
      addItem({
        id: product.id,
        name: product.name || product.name_en || 'Product',
        price: Number(product.price || 0),
        image: product.image || product.main_image_url,
        track_inventory: product.track_inventory,
        quantity_in_stock: product.quantity_in_stock,
      })
    }
  }

  return (
    <section className="px-4 py-16 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-6xl">
        <div className={`mb-8 flex items-center justify-between ${locale === 'ar' ? 'flex-row-reverse' : ''}`}>
          <h2 className="text-3xl font-bold text-[hsl(var(--theme-text))]">{locale === 'ar' ? 'منتجات مميزة' : 'Featured Products'}</h2>
          <Link href={`/${locale}/menu`} className="text-sm font-medium text-[hsl(var(--theme-primary))] hover:underline">
            {locale === 'ar' ? 'عرض القائمة كاملة' : 'View Full Menu'}
          </Link>
        </div>

        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {products.map((product) => (
            <ProductCard
              key={product.id}
              product={product}
              locale={locale}
              onAddToCart={handleAddToCart}
            />
          ))}
        </div>
      </div>
    </section>
  )
}

