import React, { useEffect, useRef, useState, useMemo, useId } from 'react'
import { useAssetUrl } from '../hooks/useAssetUrl'

export interface OptimizedImageProps {
  assetId: number
  alt: string
  width?: number
  height?: number
  className?: string
  style?: React.CSSProperties
  lazy?: boolean
  placeholder?: 'blur' | 'empty' | 'color'
  placeholderColor?: string
  priority?: boolean
  onLoad?: (img: HTMLImageElement) => void
  onError?: (error: Error) => void
  sizes?: string
  quality?: number
  formats?: ('webp' | 'avif' | 'jpeg' | 'png')[]
  fit?: 'contain' | 'cover' | 'fill' | 'scale-down'
}

const camelToKebab = (value: string): string =>
  value.replace(/[A-Z]/g, (match) => `-${match.toLowerCase()}`)

const serializeStyleObject = (style: React.CSSProperties): string =>
  Object.entries(style)
    .filter(([, val]) => val !== undefined && val !== null)
    .map(([prop, val]) => `${camelToKebab(prop)}: ${String(val)};`)
    .join(' ')

export const OptimizedImage: React.FC<OptimizedImageProps> = ({
  assetId,
  alt,
  width,
  height,
  className = '',
  style = {},
  lazy = true,
  placeholder = 'blur',
  placeholderColor = '#f0f0f0',
  priority = false,
  onLoad,
  onError,
  sizes = '100vw',
  quality = 85,
  formats = ['webp', 'avif', 'jpeg'],
  fit = 'cover',
}) => {
  const [isLoaded, setIsLoaded] = useState(!lazy)
  const [hasError, setHasError] = useState(false)
  const imgRef = useRef<HTMLImageElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)
  const reactId = useId()

  // Get CDN URLs for different formats
  const { assetUrl, isLoading } = useAssetUrl(assetId, { formats, quality })

  // Track lazy loading state
  const [isInView, setIsInView] = useState(!lazy)

  // Intersection Observer for lazy loading
  useEffect(() => {
    if (!lazy || priority || !containerRef.current) return

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setIsInView(true)
            observer.unobserve(entry.target)
          }
        })
      },
      {
        rootMargin: '50px',
      }
    )

    observer.observe(containerRef.current)

    return () => {
      observer.disconnect()
    }
  }, [lazy, priority])

  // Handle image load
  const handleLoad = () => {
    setIsLoaded(true)
    if (imgRef.current && onLoad) {
      onLoad(imgRef.current)
    }
  }

  // Handle image error
  const handleError = () => {
    setHasError(true)
    if (onError) {
      onError(new Error(`Failed to load image: ${alt}`))
    }
  }

  // Generate srcset for responsive images
  const srcSet = useMemo(() => {
    if (!assetUrl) return ''

    const baseUrl = assetUrl.split('?')[0]
    const widths = [320, 640, 960, 1280, 1920]
    const queryParams = new URLSearchParams(assetUrl.split('?')[1] || '')

    return widths
      .map((w) => {
        const params = new URLSearchParams(queryParams)
        params.set('w', w.toString())
        return `${baseUrl}?${params.toString()} ${w}w`
      })
      .join(', ')
  }, [assetUrl])

  // Aspect ratio padding for placeholder
  const aspectRatio = width && height ? (height / width) * 100 : 0
  const hasAspectSpacer = aspectRatio > 0 && !width && !height
  const encodedPlaceholderColor = placeholderColor.replace('#', '%23')
  const instanceClass = useMemo(
    () => `optimized-image-${reactId.replace(/[:]/g, '')}`,
    [reactId]
  )
  const customContainerCss = useMemo(() => serializeStyleObject(style), [style])

  return (
    <div
      ref={containerRef}
      className={`${instanceClass} optimized-image ${className}`.trim()}
      data-asset-id={assetId}
      data-loaded={isLoaded}
      data-error={hasError}
    >
      {/* Placeholder */}
      {!isLoaded && placeholder === 'blur' && (
        <div className="oi-placeholder oi-placeholder-blur" />
      )}
      {!isLoaded && placeholder === 'color' && (
        <div className="oi-placeholder oi-placeholder-color" />
      )}

      {/* Aspect ratio preserver */}
      {hasAspectSpacer && <div className="oi-aspect-preserver" />}

      {/* Image */}
      {(isInView || !lazy) && assetUrl && !hasError && (
        <picture>
          {formats.includes('avif') && (
            <source
              type="image/avif"
              srcSet={srcSet.replace(/webp/g, 'avif')}
              sizes={sizes}
            />
          )}
          {formats.includes('webp') && (
            <source
              type="image/webp"
              srcSet={srcSet.replace(/jpeg/g, 'webp')}
              sizes={sizes}
            />
          )}
          <img
            ref={imgRef}
            src={assetUrl}
            srcSet={srcSet}
            sizes={sizes}
            alt={alt}
            width={width}
            height={height}
            className={`oi-image ${isLoaded ? 'is-loaded' : 'is-loading'}`}
            onLoad={handleLoad}
            onError={handleError}
            loading={lazy && !priority ? 'lazy' : 'eager'}
            decoding="async"
          />
        </picture>
      )}

      {/* Error state */}
      {hasError && (
        <div className="oi-error">
          <span>Failed to load image</span>
        </div>
      )}

      {/* Loading indicator (optional) */}
      {!isLoaded && isLoading && (
        <div className="oi-loading-indicator">
          <div className="oi-spinner" />
        </div>
      )}

      <style>{`
        .${instanceClass} {
          position: relative;
          overflow: hidden;
          ${customContainerCss}
        }

        .${instanceClass} .oi-placeholder {
          position: absolute;
          inset: 0;
          background-color: ${placeholderColor};
          opacity: ${isLoaded ? 0 : 1};
          transition: opacity 0.3s ease-in-out;
          z-index: 0;
        }

        .${instanceClass} .oi-placeholder-blur {
          background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 ${width || 100} ${height || 100}'%3E%3Cfilter id='b'%3E%3CfeGaussianBlur stdDeviation='20'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' fill='${encodedPlaceholderColor}'/%3E%3C/filter%3E%3C/svg%3E");
          background-size: cover;
        }

        .${instanceClass} .oi-aspect-preserver {
          width: 100%;
          padding-bottom: ${hasAspectSpacer ? `${aspectRatio}%` : '0'};
        }

        .${instanceClass} .oi-image {
          width: 100%;
          height: 100%;
          object-fit: ${fit};
          object-position: center;
          z-index: 1;
          position: relative;
        }

        .${instanceClass} .oi-image.is-loading {
          display: none;
        }

        .${instanceClass} .oi-image.is-loaded {
          display: block;
        }

        .${instanceClass} .oi-error {
          position: absolute;
          inset: 0;
          display: flex;
          align-items: center;
          justify-content: center;
          background-color: #f0f0f0;
          color: #999;
          font-size: 14px;
          text-align: center;
          padding: 16px;
          z-index: 1;
        }

        .${instanceClass} .oi-loading-indicator {
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          z-index: 2;
        }

        .${instanceClass} .oi-spinner {
          width: 24px;
          height: 24px;
          border: 2px solid #f0f0f0;
          border-top: 2px solid #999;
          border-radius: 50%;
          animation: oi-spin 0.6s linear infinite;
        }

        @keyframes oi-spin {
          to {
            transform: rotate(360deg);
          }
        }
      `}</style>
    </div>
  )
}

export default OptimizedImage
