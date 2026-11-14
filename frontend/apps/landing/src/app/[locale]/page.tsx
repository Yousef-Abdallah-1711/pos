import { MarketingShell } from '@/components/marketing/MarketingShell'
import { HomeHero } from '@/components/marketing/HomeHero'
import { OrderJourney } from '@/components/marketing/OrderJourney'
import { OperationsSections } from '@/components/marketing/OperationsSections'
import { ManagementSections } from '@/components/marketing/ManagementSections'
import { ProductExplorer } from '@/components/marketing/ProductExplorer'
import { PricingOffer } from '@/components/marketing/PricingOffer'
import { GettingStarted, MarketingFAQ, FinalCTA } from '@/components/marketing/ClosingSections'
import { SectionHead } from '@/components/marketing/SectionHead'
import { getFeatures } from '@/content/features'
import { getSiteCopy, type Locale } from '@/content/site'

export default async function HomePage({ params }: { params: Promise<{ locale: Locale }> }) {
  const { locale } = await params
  const c = getSiteCopy(locale)
  return <MarketingShell><HomeHero locale={locale} /><OrderJourney locale={locale} />
    <OperationsSections locale={locale} /><ManagementSections locale={locale} />
    <section id="product-explorer" className="sv-section"><div className="sv-shell">
      <SectionHead eyebrow={c.explorerEyebrow} title={c.explorerTitle} />
      <ProductExplorer locale={locale} features={getFeatures(locale)} />
    </div></section><GettingStarted locale={locale} /><PricingOffer locale={locale} />
    <MarketingFAQ locale={locale} /><FinalCTA locale={locale} /></MarketingShell>
}
