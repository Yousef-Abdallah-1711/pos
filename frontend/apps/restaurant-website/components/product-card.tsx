"use client";

import { useState, useEffect, memo } from "react";
import { Button } from "@pos-saas/ui";
import { ShoppingCart, Plus, Minus } from "lucide-react";
import Image from "next/image";
import { useThemeColors } from "@/lib/hooks/use-theme";
import { getImageConfig, optimizeImageUrl } from "@/lib/image-optimization";

export interface Product {
  id: number;
  name?: string;
  name_en?: string;
  description?: string;
  description_en?: string;
  price: number;
  image?: string;
  main_image_url?: string;
  category?: string;
  category_id?: number;
  categoryId?: number;
  isAvailable?: boolean;
  is_available?: boolean;
  track_inventory?: boolean;
  quantity_in_stock?: number;
  rating?: number;
}

interface ProductCardProps {
  product: Product;
  onAddToCart?: (product: Product, quantity: number) => void;
  locale?: "en" | "ar";
}

/**
 * Product card component for displaying menu items
 * Uses theme CSS variables for colors
 */
function ProductCardComponent({
  product,
  onAddToCart,
  locale = "en",
}: ProductCardProps) {
  const [quantity, setQuantity] = useState(1);
  const [isAdded, setIsAdded] = useState(false);
  const themeColors = useThemeColors() as any;

  useEffect(() => {
    if (isAdded) {
      const timeoutId = setTimeout(() => setIsAdded(false), 2000);
      return () => clearTimeout(timeoutId);
    }
  }, [isAdded]);

  const handleAddToCart = () => {
    onAddToCart?.(product, quantity);
    setIsAdded(true);
  };

  const quantityInStock = Number(product.quantity_in_stock ?? 0);
  const isOutOfStock = product.track_inventory && quantityInStock <= 0;
  const isAvailable = product.isAvailable !== false && product.is_available !== false && !isOutOfStock;
  const unavailableLabel = locale === "en" ? "Unavailable" : "غير متوفر";
  const outOfStockLabel = locale === "en" ? "Out of Stock" : "نفد المخزون";
  const badgeLabel = isOutOfStock ? outOfStockLabel : unavailableLabel;

  // Fallback colors if theme not loaded
  const bgColor = themeColors?.background || "#ffffff";
  const textColor = themeColors?.text || "#1f2937";
  const primaryColor = themeColors?.primary || "#f97316";
  const borderColor = themeColors?.border || "#e5e7eb";
  const accentColor = themeColors?.accent || "#fbbf24";

  return (
    <div
      className="rounded-lg overflow-hidden group shadow-sm hover:shadow-md transition-shadow duration-200"
      style={{
        backgroundColor: bgColor,
        borderColor: borderColor,
        border: `1px solid ${borderColor}`,
      }}
    >
      {/* Image Container */}
      <div
        className="relative h-48 overflow-hidden"
        style={{
          backgroundColor: `hsl(var(--theme-background) / 0.5)`,
        }}
      >
        {(product.image || product.main_image_url) ? (
          <Image
            src={optimizeImageUrl(product.image || product.main_image_url)}
            alt={product.name || product.name_en || "Product"}
            fill
            loading="lazy"
            quality={80}
            sizes={getImageConfig('productCard').sizes}
            className="object-cover group-hover:scale-105 transition-transform duration-300"
          />
        ) : (
          <div
            className="w-full h-full flex items-center justify-center"
            style={{ color: `hsl(var(--theme-text) / 0.3)` }}
          >
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

        {/* Availability Badge */}
        {!isAvailable && (
          <div
            className="absolute inset-0 flex items-center justify-center z-10"
            style={{ backgroundColor: "rgba(0, 0, 0, 0.5)" }}
          >
            <span style={{ color: "#ffffff", fontWeight: "600" }}>
              {badgeLabel}
            </span>
          </div>
        )}

        {/* Rating Badge */}
        {product.rating && (
          <div
            className="absolute top-2 right-2 px-2 py-1 rounded-full text-sm font-semibold flex items-center gap-1"
            style={{
              backgroundColor: accentColor,
              color: "#1f2937",
            }}
          >
            ⭐ {product.rating.toFixed(1)}
          </div>
        )}
      </div>

      {/* Content */}
      <div className="p-4">
        {/* Title */}
        <h3
          className="text-lg font-semibold mb-1 line-clamp-2"
          style={{ color: textColor }}
        >
          {product.name || product.name_en || "Product"}
        </h3>

        {/* Category */}
        {product.category && (
          <p
            className="text-xs mb-2"
            style={{ color: `hsl(var(--theme-text) / 0.6)` }}
          >
            {product.category}
          </p>
        )}

        {/* Description */}
        <p
          className="text-sm mb-4 line-clamp-2"
          style={{ color: `hsl(var(--theme-text) / 0.7)` }}
        >
          {product.description || product.description_en || ""}
        </p>

        {/* Price and Actions */}
        <div className="flex items-center justify-between">
          <span
            className="text-2xl font-bold"
            style={{ color: primaryColor }}
          >
            ${product.price.toFixed(2)}
          </span>

          {isAvailable ? (
            <div className="flex items-center gap-2">
              {/* Quantity Selector */}
              <div
                className="flex items-center gap-1 rounded-lg"
                style={{
                  backgroundColor: `hsl(var(--theme-background) / 0.8)`,
                  borderColor: borderColor,
                  border: `1px solid ${borderColor}`,
                }}
              >
                <button
                  onClick={() => setQuantity(Math.max(1, quantity - 1))}
                  className="p-1 transition-colors hover:bg-gray-100 dark:hover:bg-gray-700"
                  style={{
                    color: textColor,
                  }}
                  disabled={quantity === 1}
                >
                  <Minus className="h-4 w-4" />
                </button>
                <span
                  className="px-2 min-w-[2rem] text-center text-sm font-semibold"
                  style={{ color: textColor }}
                >
                  {quantity}
                </span>
                <button
                  onClick={() => setQuantity(quantity + 1)}
                  className="p-1 transition-colors hover:bg-gray-100 dark:hover:bg-gray-700"
                  style={{
                    color: textColor,
                  }}
                >
                  <Plus className="h-4 w-4" />
                </button>
              </div>

              {/* Add to Cart Button */}
              <Button
                size="sm"
                variant={isAdded ? "default" : "outline"}
                onClick={handleAddToCart}
                style={
                  isAdded
                    ? {
                      backgroundColor: "#22c55e",
                      color: "#ffffff",
                      borderColor: "#22c55e",
                    }
                    : {
                      borderColor: primaryColor,
                      color: primaryColor,
                    }
                }
              >
                {isAdded ? (
                  <>✓</>
                ) : (
                  <>
                    <ShoppingCart className="h-4 w-4" />
                  </>
                )}
              </Button>
            </div>
          ) : (
            <Button disabled variant="outline" size="sm">
              {badgeLabel}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}

export const ProductCard = memo(ProductCardComponent);
