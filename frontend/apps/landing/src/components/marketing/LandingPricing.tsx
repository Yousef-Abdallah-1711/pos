import Link from 'next/link'
import { ArrowUpRight } from 'lucide-react'
import type { Locale } from '@/content/site'
import { TrialLink } from './TrialLink'

const copy = {
  en: { eyebrow: 'Clear terms before you begin', title: 'See the plan, then make it yours.', body: 'Review the current trial and billing details on the pricing page before creating a restaurant workspace.', pricing: 'View pricing', account: 'Create your restaurant account' },
  ar: { eyebrow: 'التفاصيل قبل البداية', title: 'راجع الخطة وابدأ على راحتك.', body: 'اطلع على تفاصيل التجربة والفواتير الحالية في صفحة الأسعار قبل إنشاء مساحة عمل مطعمك.', pricing: 'شوف الأسعار', account: 'أنشئ حساب مطعمك' },
} as const

export function LandingPricing({ locale }: { locale: Locale }) {
  const c = copy[locale]
  return <section className="cv-section cv-pricing" aria-labelledby="cv-pricing-title">
    <div className="cv-wrap cv-pricing-layout">
      <div><p className="cv-eyebrow">{c.eyebrow}</p><h2 id="cv-pricing-title">{c.title}</h2></div>
      <div><p>{c.body}</p><div className="cv-pricing-actions"><Link href={`/${locale}/pricing`}>{c.pricing}<ArrowUpRight size={18} aria-hidden="true" /></Link><TrialLink locale={locale} source="home-pricing" label={c.account} /></div></div>
    </div>
  </section>
}
