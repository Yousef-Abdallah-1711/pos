'use client';

import Image from 'next/image';
import { formatPrice } from '@/lib/utils/format';

interface ProductCardProps {
  product: any;
  onAddToCart?: (product: any, quantity: number) => void;
  locale?: 'en' | 'ar';
}

export function ProductCard({ product, onAddToCart, locale = 'en' }: ProductCardProps) {
  const isRTL = locale === 'ar';
  const quantityInStock = Number(product.quantity_in_stock ?? 0);
  const isOutOfStock = product.track_inventory && quantityInStock <= 0;
  const isUnavailable = product.is_available === false || product.isAvailable === false;
  const isAvailable = !isOutOfStock && !isUnavailable;
  const outOfStockLabel = isRTL ? 'نفد المخزون' : 'Out of Stock';
  const unavailableLabel = isRTL ? 'غير متوفر' : 'Unavailable';
  const badgeLabel = isOutOfStock ? outOfStockLabel : unavailableLabel;
  const price = Number(product.price ?? 0);

  return (
    <div className="border rounded-lg p-4 hover:shadow-lg transition-all relative">
      {!isAvailable && (
        <span className="absolute top-3 right-3 z-10 rounded-full bg-red-600 text-white text-xs font-semibold px-3 py-1">
          {badgeLabel}
        </span>
      )}
      {/* Image */}
      {product.image && (
        <div className="relative h-48 mb-4 overflow-hidden rounded-md">
          <Image
            src={product.image}
            alt={product.name}
            fill
            className="object-cover hover:scale-105 transition-transform duration-300"
            sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
          />
        </div>
      )}

      {/* Name */}
      <h3 className="font-bold text-xl mb-1 text-gray-900 dark:text-gray-100">{product.name}</h3>

      {/* Description */}
      <p className="text-gray-600 dark:text-gray-400 text-sm mb-4 line-clamp-2 min-h-[40px]">{product.description}</p>

      {/* Price and Action */}
      <div className="flex items-center justify-between mt-auto">
        <p className="text-xl font-black text-orange-600">{formatPrice(price, locale)}</p>

        {/* Add to Cart */}
        <button
          onClick={() => {
            if (!isAvailable) return;
            onAddToCart?.(product, 1);
          }}
          disabled={!isAvailable}
          className={`px-4 py-2 rounded-lg font-bold transition-colors ${
            isAvailable
              ? 'bg-orange-600 text-white hover:bg-orange-700 active:scale-95'
              : 'bg-gray-300 text-gray-600 cursor-not-allowed'
          }`}
        >
          {isRTL ? 'إضافة' : 'Add'}
        </button>
      </div>
    </div>
  );
}
