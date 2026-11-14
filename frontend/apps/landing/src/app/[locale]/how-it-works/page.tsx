import { MarketingShell } from '@/components/marketing/MarketingShell'
import { OrderJourney } from '@/components/marketing/OrderJourney'
import { GettingStarted, FinalCTA } from '@/components/marketing/ClosingSections'
import { getSiteCopy, type Locale } from '@/content/site'
import { marketingMetadata } from '@/lib/marketingMetadata'
type Props = { params: Promise<{ locale: Locale }> }
export async function generateMetadata({ params }: Props) {
 const { locale } = await params; const c = getSiteCopy(locale)
 return marketingMetadata(locale, '/how-it-works', c.how, c.journeyBody)
}
export default async function HowPage({ params }: Props) {
 const { locale } = await params; const c = getSiteCopy(locale)
 return <MarketingShell><section className="sv-section sv-feature-hero"><div className="sv-shell"><p className="sv-eyebrow">{c.how}</p><h1>{c.journeyTitle}</h1><p className="sv-lead">{c.journeyBody}</p></div></section><OrderJourney locale={locale} /><GettingStarted locale={locale} /><FinalCTA locale={locale} /></MarketingShell>
}
