'use client'
import { useEffect } from 'react'
import { trackFunnel } from '@/lib/funnel'
import type { Locale } from '@/content/site'
export function ViewEvent({ locale, source, event = 'feature_page_viewed' }: { locale: Locale; source: string; event?: 'feature_page_viewed' | 'demo_contact_initiated' }) {
  useEffect(() => { trackFunnel(event, { locale, source }) }, [locale, source, event])
  return null
}
