import type { Locale } from '@/content/site'
import captures from '@/content/captures.json'
import { getSiteCopy } from '@/content/site'
import { CaptureMedia } from './CaptureMedia'

interface ProductCaptureProps {
  image: keyof typeof captures
  alt: string
  locale: Locale
  className?: string
  priority?: boolean
}

export function ProductCapture({ image, alt, locale, className = '', priority }: ProductCaptureProps) {
  const copy = getSiteCopy(locale)
  return <figure className={`sv-capture ${className}`}>
    <CaptureMedia image={image} alt={alt} locale={locale} captions={copy.captureMatrix} openLabel={copy.openCapture} priority={priority} />
  </figure>
}
