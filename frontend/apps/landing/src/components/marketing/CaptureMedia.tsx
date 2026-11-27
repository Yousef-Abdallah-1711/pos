'use client'

import Image from 'next/image'
import { useSyncExternalStore } from 'react'
import { useTheme } from 'next-themes'
import captures from '@/content/captures.json'
import type { Locale } from '@/content/site'

const MATRIX_IMAGES = ['pos', 'kitchen', 'inventory', 'staff', 'theme', 'dashboard', 'checkout'] as const
type CaptureName = keyof typeof captures
type CaptureCaption = 'enLight' | 'enDark' | 'arLight' | 'arDark'
const subscribeToClientMount = () => () => undefined
const getClientMountSnapshot = () => true
const getServerMountSnapshot = () => false

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
  const mounted = useSyncExternalStore(subscribeToClientMount, getClientMountSnapshot, getServerMountSnapshot)
  const hasMatrix = MATRIX_IMAGES.includes(image as (typeof MATRIX_IMAGES)[number])
  // Keep the server and first client render identical; next-themes reads saved
  // preferences only in the browser, then this image updates after hydration.
  const theme = !mounted || resolvedTheme === 'light' ? 'light' : 'dark'
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
        priority={priority} loading={priority ? 'eager' : 'lazy'}
        sizes="(max-width: 768px) 100vw, (max-width: 1440px) 75vw, 1180px" />
    </div>
    <figcaption><span>{caption}</span>
      <a href={source} target="_blank" rel="noreferrer">{openLabel}<span className="sr-only">: {alt}</span></a>
    </figcaption>
  </>
}
