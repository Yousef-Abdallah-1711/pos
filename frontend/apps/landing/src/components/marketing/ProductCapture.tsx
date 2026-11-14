import Image from 'next/image'
import captures from '@/content/captures.json'
import type { Locale } from '@/content/site'
import { getSiteCopy } from '@/content/site'

interface ProductCaptureProps {
  image: keyof typeof captures
  alt: string
  locale: Locale
  className?: string
  priority?: boolean
}

export function ProductCapture({ image, alt, locale, className = '', priority }: ProductCaptureProps) {
  const src = `/images/product/servio/${image}.png`
  return <figure className={`sv-capture ${className}`}>
    <div className="sv-capture-image">
      <Image src={src} alt={alt} width={captures[image].width} height={captures[image].height} priority={priority}
        sizes="(max-width: 768px) 100vw, (max-width: 1440px) 75vw, 1180px" />
    </div>
    <figcaption><span>{getSiteCopy(locale).englishCapture}</span>
      <a href={src} target="_blank" rel="noreferrer">{getSiteCopy(locale).openCapture}<span className="sr-only">: {alt}</span></a>
    </figcaption>
  </figure>
}
