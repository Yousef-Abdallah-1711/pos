import { MarketingShell } from '@/components/marketing/MarketingShell'
import { ProductExplorer } from '@/components/marketing/ProductExplorer'
import { FinalCTA } from '@/components/marketing/ClosingSections'
import { getFeatures } from '@/content/features'
import { getSiteCopy, type Locale } from '@/content/site'
import { marketingMetadata } from '@/lib/marketingMetadata'
type Props = { params: Promise<{ locale: Locale }> }
export async function generateMetadata({ params }: Props) {
 const { locale } = await params; const c = getSiteCopy(locale)
 return marketingMetadata(locale, '/demo', c.demo, c.explorerBody)
}
export default async function DemoPage({ params }: Props) {
 const { locale } = await params; const c = getSiteCopy(locale)
 return <MarketingShell><section className="sv-section sv-feature-hero"><div className="sv-shell"><p className="sv-eyebrow">{c.demo}</p><h1>{c.explorerTitle}</h1><p className="sv-lead">{c.explorerBody}</p><ProductExplorer headingLevel={2} locale={locale} features={getFeatures(locale)} /></div></section><FinalCTA locale={locale} /></MarketingShell>
}
