'use client'

import Link from 'next/link'
import { ArrowUpRight } from 'lucide-react'
import type { Locale } from '@/content/site'
import { getConversionCopy } from '@/content/conversionCopy'
import { trackFunnel } from '@/lib/funnel'

interface TrialLinkProps {
  locale: Locale
  source: string
  label?: string
  billing?: 'monthly' | 'annual'
  compact?: boolean
  secondary?: boolean
}

export function TrialLink({ locale, source, label, billing = 'monthly', secondary, compact }: TrialLinkProps) {
  return <Link prefetch={false} className={`sv-button${secondary ? ' sv-button-secondary' : ''}${compact ? ' sv-button-compact' : ''}`}
    href={`/${locale}/auth/register?billing=${billing}`}
    onClick={() => trackFunnel('landing_primary_cta_clicked', { locale, source, billing })}>
    {label ?? getConversionCopy(locale).createAccount}<ArrowUpRight size={18} aria-hidden="true" className="sv-directional" />
  </Link>
}
