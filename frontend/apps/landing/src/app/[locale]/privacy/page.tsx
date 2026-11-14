import Link from 'next/link'
import { MarketingShell } from '@/components/marketing/MarketingShell'
import { TrialLink } from '@/components/marketing/TrialLink'
import { getSiteCopy, type Locale } from '@/content/site'
import { marketingMetadata } from '@/lib/marketingMetadata'
type Props = { params: Promise<{ locale: Locale }> }
const copy = { en: { title: 'Privacy information for this website', body: 'Public product captures contain synthetic data. The website’s funnel controls dispatch local browser events without personal fields or external analytics requests. Trial registration sends the account and restaurant details you enter to the SERVIO service to create your workspace. Authentication uses service-managed cookies. This page describes the implemented website behavior; it does not establish additional contractual terms.' }, ar: { title: 'معلومات الخصوصية لهذا الموقع', body: 'تحتوي صور المنتج العامة على بيانات تجريبية. ترسل أدوات مسار التسجيل أحداثًا محلية في المتصفح دون بيانات شخصية أو طلبات تحليلات خارجية. يرسل التسجيل بيانات الحساب والمطعم التي تدخلها إلى خدمة سيرفيو لإنشاء مساحة العمل. تستخدم المصادقة ملفات تعريف ارتباط تديرها الخدمة. تصف هذه الصفحة سلوك الموقع المنفذ، ولا تضع شروطًا تعاقدية إضافية.' } }
export async function generateMetadata({ params }: Props) { const { locale } = await params; return marketingMetadata(locale, '/privacy', copy[locale].title, copy[locale].body) }
export default async function Page({ params }: Props) {
 const { locale } = await params; const c = getSiteCopy(locale)
 return <MarketingShell><section className="sv-section sv-feature-hero"><div className="sv-shell"><h1>{copy[locale].title}</h1><p className="sv-lead">{copy[locale].body}</p><div className="sv-actions"><Link className="sv-text-link" href={`/${locale}/demo`}>{c.demo}<span aria-hidden="true">↗</span></Link><TrialLink locale={locale} source="privacy" /></div></div></section></MarketingShell>
}
