'use client';

import { useCallback, useMemo } from 'react';
import { ProductCard } from './product-card';
import { useCartStore } from '@/lib/store/cart-store';
import { useProducts } from '@/lib/hooks/use-api-queries';
import { AlertCircle, Loader } from 'lucide-react';

interface Product {
  id: number;
  name_en?: string;
  name?: string;
  description_en?: string;
  description?: string;
  price: number;
  main_image_url?: string;
  image?: string;
  category_id?: number;
  category?: string;
  is_available?: boolean;
  isAvailable?: boolean;
  track_inventory?: boolean;
  quantity_in_stock?: number;
  rating?: number;
}

interface FeaturedProductsSectionProps {
  restaurantSlug: string;
  locale?: 'en' | 'ar';
  limit?: number;
}

/**
 * Featured Products Section
 * Fetches and displays real products from the API with Add to Cart functionality
 */
export function FeaturedProductsSection({
  restaurantSlug,
  locale = 'en',
  limit = 6,
}: FeaturedProductsSectionProps) {
  const addToCart = useCartStore((state) => state.addItem);
  const isRTL = locale === 'ar';

  // Use React Query for data fetching with built-in caching
  const { data: allProducts = [], isLoading, error } = useProducts(restaurantSlug);

  // Limit featured products and memoize
  const products = useMemo(() => {
    return Array.isArray(allProducts) ? allProducts.slice(0, limit) : [];
  }, [allProducts, limit]);

  // Memoize handleAddToCart to preserve reference
  const handleAddToCart = useCallback((product: Product, quantity: number) => {
    // Add to cart store
    for (let i = 0; i < quantity; i++) {
      addToCart({
        id: product.id,
        name: product.name || product.name_en || 'Product',
        price: product.price || 0,
        image: product.image || product.main_image_url || '',
        track_inventory: product.track_inventory,
        quantity_in_stock: product.quantity_in_stock,
      });
    }
  }, [addToCart]);

  // Loading state
  if (isLoading) {
    return (
      <section className="w-full py-16 px-4 bg-gray-50 dark:bg-gray-800/50">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center min-h-96">
            <Loader className="h-8 w-8 animate-spin text-gray-400" />
          </div>
        </div>
      </section>
    );
  }

  // Error state
  if (error && products.length === 0) {
    const errorMessage = error instanceof Error ? error.message : 'Failed to load products. Please try again later.';
    return (
      <section className="w-full py-16 px-4 bg-gray-50 dark:bg-gray-800/50">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center gap-3 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
            <AlertCircle className="h-5 w-5 text-yellow-600 dark:text-yellow-400 flex-shrink-0" />
            <p className="text-sm text-yellow-800 dark:text-yellow-300">{errorMessage}</p>
          </div>
        </div>
      </section>
    );
  }

  // No products
  if (!isLoading && products.length === 0) {
    return (
      <section className="w-full py-16 px-4 bg-gray-50 dark:bg-gray-800/50">
        <div className="max-w-6xl mx-auto text-center">
          <p className="text-gray-600 dark:text-gray-400">No products available at this time.</p>
        </div>
      </section>
    );
  }

  // Render products
  return (
    <section className="w-full py-16 px-4 bg-gray-50 dark:bg-gray-800/50">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className={`text-center mb-12 ${isRTL ? 'text-right' : 'text-left'}`}>
          <h2 className="text-4xl font-bold mb-4 text-gray-900 dark:text-white">
            {locale === 'en' ? 'Featured Products' : 'المنتجات المميزة'}
          </h2>
          <p className="text-gray-600 dark:text-gray-400">
            {locale === 'en'
              ? 'Check out our selection of delicious items'
              : 'تصفح مختاراتنا من الأطباق اللذيذة'}
          </p>
        </div>

        {/* Products Grid */}
        <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3">
          {products.map((product) => (
            <ProductCard
              key={product.id}
              product={product}
              onAddToCart={handleAddToCart}
              locale={locale as 'en' | 'ar'}
            />
          ))}
        </div>

        {/* View All Link */}
        <div className="text-center mt-12">
          <a
            href={`/${locale}/menu`}
            className="inline-block px-8 py-3 bg-blue-600 dark:bg-blue-600 text-white rounded-lg hover:bg-blue-700 dark:hover:bg-blue-700 transition font-semibold"
          >
            {locale === 'en' ? 'View All Products' : 'عرض جميع المنتجات'}
          </a>
        </div>
      </div>
    </section>
  );
}
