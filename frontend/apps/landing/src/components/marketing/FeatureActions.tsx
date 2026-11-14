import Link from 'next/link'
import type { Locale } from '@/content/site'
import { TrialLink } from './TrialLink'

export function FeatureActions({ locale, source, destination, label }: { locale: Locale; source: string; destination: string; label: string }) {
  return <div className="sv-section-actions"><TrialLink locale={locale} source={source} compact />
    <Link className="sv-text-link" href={`/${locale}/${destination}`}>{label}<span aria-hidden="true">↗</span></Link>
  </div>
}
