import { ViewEvent } from '@/components/marketing/ViewEvent'
import { MarketingShell } from '@/components/marketing/MarketingShell'
import { ProductExplorer } from '@/components/marketing/ProductExplorer'
import { FinalCTA } from '@/components/marketing/ClosingSections'
import { getFeatures } from '@/content/features'
import { getSiteCopy, type Locale } from '@/content/site'
import { marketingMetadata } from '@/lib/marketingMetadata'
type Props = { params: Promise<{ locale: Locale }> }
export async function generateMetadata({ params }: Props) {
 const { locale } = await params; const c = getSiteCopy(locale)
 return marketingMetadata(locale, '/demo', c.demo, locale === 'ar' ? 'استكشف شاشات سيرفيو الفعلية لنقاط البيع وتذاكر المطبخ والطلبات النقدية والقوالب والمخزون وسجلات الموظفين، ضمن جولة ببيانات تجريبية.' : 'Explore authentic SERVIO screens for cashier POS, kitchen tickets, guest cash orders, storefront themes, inventory and staff records.')
}
export default async function DemoPage({ params }: Props) {
 const { locale } = await params; const c = getSiteCopy(locale)
 return <MarketingShell><ViewEvent locale={locale} source="demo" event="demo_contact_initiated" /><section className="sv-section sv-feature-hero"><div className="sv-shell"><p className="sv-eyebrow">{c.demo}</p><h1>{c.explorerTitle}</h1><p className="sv-lead">{c.explorerBody}</p><ProductExplorer headingLevel={2} locale={locale} features={getFeatures(locale)} /></div></section><FinalCTA locale={locale} /></MarketingShell>
}
