import { MarketingShell } from '@/components/marketing/MarketingShell'
import { PricingOffer } from '@/components/marketing/PricingOffer'
import { MarketingFAQ } from '@/components/marketing/ClosingSections'
import { marketingMetadata } from '@/lib/marketingMetadata'
import { getSiteCopy, type Locale } from '@/content/site'
type Props = { params: Promise<{ locale: Locale }> }
export async function generateMetadata({ params }: Props) {
 const { locale } = await params; const c = getSiteCopy(locale)
 return marketingMetadata(locale, '/pricing', c.pricing, c.offerBody)
}
export default async function PricingPage({ params }: Props) {
 const { locale } = await params
 return <MarketingShell><PricingOffer locale={locale} detailed /><MarketingFAQ locale={locale} /></MarketingShell>
}
