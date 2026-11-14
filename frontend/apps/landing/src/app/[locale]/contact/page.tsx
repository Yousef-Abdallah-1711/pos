import Link from 'next/link'
import { MarketingShell } from '@/components/marketing/MarketingShell'
import { TrialLink } from '@/components/marketing/TrialLink'
import { getSiteCopy, type Locale } from '@/content/site'
import { marketingMetadata } from '@/lib/marketingMetadata'
type Props = { params: Promise<{ locale: Locale }> }
const copy = { en: { title: 'Explore SERVIO with a product walkthrough', body: 'The walkthrough shows authentic screens and explains the supported workflow. You can explore it immediately or create a trial workspace to try the product yourself. Live demonstration scheduling is not available on this page.' }, ar: { title: 'تعرّف على سيرفيو من خلال جولة في المنتج', body: 'توضح الجولة شاشات فعلية ومسارات العمل المدعومة. يمكنك استكشافها الآن أو إنشاء مساحة تجريبية لتجربة المنتج بنفسك. لا تتوفر جدولة عروض مباشرة من هذه الصفحة.' } }
export async function generateMetadata({ params }: Props) { const { locale } = await params; return marketingMetadata(locale, '/contact', copy[locale].title, copy[locale].body) }
export default async function Page({ params }: Props) {
 const { locale } = await params; const c = getSiteCopy(locale)
 return <MarketingShell><section className="sv-section sv-feature-hero"><div className="sv-shell"><h1>{copy[locale].title}</h1><p className="sv-lead">{copy[locale].body}</p><div className="sv-actions"><Link className="sv-text-link" href={`/${locale}/demo`}>{c.demo}<span aria-hidden="true">↗</span></Link><TrialLink locale={locale} source="contact" /></div></div></section></MarketingShell>
}
