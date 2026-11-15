'use client'

import Image from 'next/image'
import { useTheme } from 'next-themes'
import captures from '@/content/captures.json'
import type { Locale } from '@/content/site'

const MATRIX_IMAGES = ['pos', 'kitchen', 'inventory', 'staff', 'theme', 'dashboard', 'checkout'] as const
type CaptureName = keyof typeof captures
type CaptureCaption = 'enLight' | 'enDark' | 'arLight' | 'arDark'

interface CaptureMediaProps {
  image: CaptureName
  alt: string
  locale: Locale
  captions: Record<CaptureCaption, string>
  openLabel: string
  priority?: boolean
}

export function CaptureMedia({ image, alt, locale, captions, openLabel, priority }: CaptureMediaProps) {
  const { resolvedTheme } = useTheme()
  const hasMatrix = MATRIX_IMAGES.includes(image as (typeof MATRIX_IMAGES)[number])
  const theme = resolvedTheme === 'light' ? 'light' : 'dark'
  const sourceLocale = hasMatrix ? locale : 'en'
  const sourceTheme = hasMatrix ? theme : 'light'
  const source = hasMatrix
    ? `/images/product/servio/matrix/${image}-${sourceLocale}-${sourceTheme}.png`
    : `/images/product/servio/${image}.png`
  const caption = hasMatrix
    ? captions[`${sourceLocale}${sourceTheme === 'dark' ? 'Dark' : 'Light'}` as CaptureCaption]
    : captions.enLight
  const dimensions = image === 'checkout' && sourceLocale === 'ar'
    ? { width: 390, height: 2100 }
    : captures[image]

  return <>
    <div className="sv-capture-image" data-source-locale={sourceLocale} data-source-theme={sourceTheme}>
      <Image src={source} alt={`${alt} — ${caption}`} width={dimensions.width} height={dimensions.height}
        priority={priority} sizes="(max-width: 768px) 100vw, (max-width: 1440px) 75vw, 1180px" />
    </div>
    <figcaption><span>{caption}</span>
      <a href={source} target="_blank" rel="noreferrer">{openLabel}<span className="sr-only">: {alt}</span></a>
    </figcaption>
  </>
}
