/**
 * Optimized Image Component
 *
 * Wrapper around Next.js Image component with:
 * - Automatic WebP/AVIF format selection
 * - Responsive image sizes
 * - Lazy loading
 * - Placeholder support
 * - Performance monitoring
 */

import Image, { ImageProps } from 'next/image'
import { useState } from 'react'

interface OptimizedImageProps extends Omit<ImageProps, 'onLoadingComplete'> {
  variant?: 'thumbnail' | 'hero' | 'avatar' | 'default'
  containerClassName?: string
  onLoadComplete?: () => void
  fallback?: string
}

// Image variant configurations
const variantConfig = {
  thumbnail: {
    maxWidth: 200,
    maxHeight: 200,
    quality: 75,
    sizes: '(max-width: 640px) 100px, 200px',
  },
  hero: {
    maxWidth: 1920,
    maxHeight: 600,
    quality: 90,
    sizes: '(max-width: 640px) 100vw, (max-width: 1024px) 90vw, 1920px',
  },
  avatar: {
    maxWidth: 100,
    maxHeight: 100,
    quality: 80,
    sizes: '(max-width: 640px) 50px, 100px',
  },
  default: {
    maxWidth: 800,
    maxHeight: 600,
    quality: 85,
    sizes: '(max-width: 640px) 100vw, (max-width: 1024px) 90vw, 800px',
  },
}

export function OptimizedImage({
  src,
  alt,
  variant = 'default',
  width,
  height,
  quality = variantConfig[variant].quality,
  sizes = variantConfig[variant].sizes,
  containerClassName = '',
  onLoadComplete,
  fallback,
  ...props
}: OptimizedImageProps) {
  const [isLoading, setIsLoading] = useState(true)
  const [hasError, setHasError] = useState(false)

  const config = variantConfig[variant]

  // Use fallback if image fails to load
  const imageSrc = hasError && fallback ? fallback : src
  const maxWidth = width || config.maxWidth
  const maxHeight = height || config.maxHeight
  const ratio = width && height ? `${width}/${height}` : 'auto'

  return (
    <div
      className={`optimized-image-container relative overflow-hidden bg-gray-100 ${containerClassName}`}
    >
      <Image
        src={imageSrc}
        alt={alt}
        fill={!width || !height}
        width={width}
        height={height}
        quality={quality}
        sizes={sizes}
        priority={false}
        loading="lazy"
        onLoadingComplete={() => {
          setIsLoading(false)
          onLoadComplete?.()
        }}
        onError={() => {
          setIsLoading(false)
          setHasError(true)
        }}
        {...props}
      />

      {/* Loading skeleton */}
      {isLoading && (
        <div className="absolute inset-0 bg-gradient-to-r from-gray-100 via-gray-50 to-gray-100 animate-pulse" />
      )}

      {/* Error state */}
      {hasError && !fallback && (
        <div className="absolute inset-0 flex items-center justify-center bg-gray-100">
          <span className="text-xs text-gray-500">Failed to load image</span>
        </div>
      )}
      <style jsx>{`
        .optimized-image-container {
          max-width: ${maxWidth}px;
          max-height: ${maxHeight}px;
          aspect-ratio: ${ratio};
        }
      `}</style>
    </div>
  )
}

/**
 * Responsive Image Gallery Component
 */
export function ImageGallery({
  images,
  columns = 3,
  variant = 'default',
}: {
  images: Array<{ src: string; alt: string }>
  columns?: 2 | 3 | 4
  variant?: keyof typeof variantConfig
}) {
  const gridCols = {
    2: 'grid-cols-2',
    3: 'md:grid-cols-3',
    4: 'md:grid-cols-4',
  }

  return (
    <div className={`grid gap-4 grid-cols-1 ${gridCols[columns]}`}>
      {images.map((image, idx) => (
        <OptimizedImage
          key={idx}
          src={image.src}
          alt={image.alt}
          variant={variant}
          containerClassName="w-full h-auto"
        />
      ))}
    </div>
  )
}

/**
 * Image with aspect ratio wrapper
 */
export function AspectRatioImage({
  src,
  alt,
  ratio = '16/9',
  ...props
}: OptimizedImageProps & { ratio?: string }) {
  return (
    <div className="aspect-ratio-image relative w-full">
      <Image
        src={src}
        alt={alt}
        fill
        className="object-cover"
        {...props}
      />
      <style jsx>{`
        .aspect-ratio-image {
          aspect-ratio: ${ratio};
        }
      `}</style>
    </div>
  )
}

/**
 * Avatar image component
 */
export function Avatar({
  src,
  alt,
  size = 'md',
  fallback,
}: {
  src: string
  alt: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  fallback?: string
}) {
  const sizeMap = {
    sm: 'w-8 h-8',
    md: 'w-12 h-12',
    lg: 'w-16 h-16',
    xl: 'w-24 h-24',
  }

  return (
    <div className={`${sizeMap[size]} rounded-full overflow-hidden bg-gray-200`}>
      <OptimizedImage
        src={src}
        alt={alt}
        variant="avatar"
        width={parseInt(sizeMap[size].split('-')[1])}
        height={parseInt(sizeMap[size].split('-')[1])}
        fallback={fallback}
      />
    </div>
  )
}
