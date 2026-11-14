'use client'

import Link from 'next/link'
import { ArrowUpRight } from 'lucide-react'
import type { Locale } from '@/content/site'
import { getSiteCopy } from '@/content/site'
import { trackFunnel } from '@/lib/funnel'

interface TrialLinkProps {
  locale: Locale
  source: string
  billing?: 'monthly' | 'annual'
  secondary?: boolean
}

export function TrialLink({ locale, source, billing = 'monthly', secondary }: TrialLinkProps) {
  return <Link className={`sv-button${secondary ? ' sv-button-secondary' : ''}`}
    href={`/${locale}/auth/register?billing=${billing}`}
    onClick={() => trackFunnel('landing_primary_cta_clicked', { locale, source, billing })}>
    {getSiteCopy(locale).trial}<ArrowUpRight size={18} aria-hidden="true" className="sv-directional" />
  </Link>
}
