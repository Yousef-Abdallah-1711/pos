'use client'

import Link from 'next/link'
import { useEffect, useState } from 'react'
import useEmblaCarousel from 'embla-carousel-react'
import Autoplay from 'embla-carousel-autoplay'
import { ProductCard, type Product } from '@/components/commerce/product-card'
import { useCartStore } from '@/lib/store/cart-store'
import type { ThemeLocale } from '@/types/theme'

interface ProductCarouselProps {
  locale: ThemeLocale
  products: Product[]
}

export function ProductCarousel({ locale, products }: ProductCarouselProps) {
  const addItem = useCartStore((state) => state.addItem)
  const [selectedIndex, setSelectedIndex] = useState(0)
  const [emblaRef, emblaApi] = useEmblaCarousel(
    { loop: true, align: 'start', direction: locale === 'ar' ? 'rtl' : 'ltr' },
    [Autoplay({ delay: 4000, stopOnMouseEnter: true })]
  )

  useEffect(() => {
    if (!emblaApi) return
    const onSelect = () => setSelectedIndex(emblaApi.selectedScrollSnap())
    onSelect()
    emblaApi.on('select', onSelect)
    return () => {
      emblaApi.off('select', onSelect)
    }
  }, [emblaApi])

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

        <div className="overflow-hidden" ref={emblaRef}>
          <div className="flex">
            {products.map((product) => (
              <div
                key={product.id}
                className="min-w-0 flex-[0_0_100%] px-2 sm:flex-[0_0_50%] lg:flex-[0_0_33.333%]"
              >
                <ProductCard
                  product={product}
                  locale={locale}
                  onAddToCart={handleAddToCart}
                />
              </div>
            ))}
          </div>
        </div>

        {products.length > 1 ? (
          <div className="mt-6 flex items-center justify-center gap-2">
            {products.map((_, index) => (
              <button
                key={index}
                type="button"
                className={`h-2.5 rounded-full transition-all ${selectedIndex === index ? 'w-8 bg-theme-primary' : 'w-2.5 bg-theme-border'}`}
                onClick={() => emblaApi?.scrollTo(index)}
                aria-label={`Go to product ${index + 1}`}
              />
            ))}
          </div>
        ) : null}
      </div>
    </section>
  )
}
