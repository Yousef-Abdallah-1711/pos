'use client'
import { useEffect } from 'react'
import { trackFunnel } from '@/lib/funnel'
import type { Locale } from '@/content/site'
export function ViewEvent({ locale, source }: { locale: Locale; source: string }) {
  useEffect(() => { trackFunnel('feature_page_viewed', { locale, source }) }, [locale, source])
  return null
}
