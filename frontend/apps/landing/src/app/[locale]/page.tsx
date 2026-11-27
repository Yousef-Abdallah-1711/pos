import { MarketingShell } from '@/components/marketing/MarketingShell'
import { ConversionLanding } from '@/components/marketing/ConversionLanding'
import { marketingMetadata } from '@/lib/marketingMetadata'
import { getConversionCopy } from '@/content/conversionCopy'
import type { Locale } from '@/content/site'

export async function generateMetadata({ params }: { params: Promise<{ locale: Locale }> }) {
  const { locale } = await params
  const c = getConversionCopy(locale)
  return marketingMetadata(locale, '', c.heroTitle, c.heroBody)
}

export default async function HomePage({ params }: { params: Promise<{ locale: Locale }> }) {
  const { locale } = await params
  return <MarketingShell locale={locale}><ConversionLanding locale={locale} /></MarketingShell>
}
