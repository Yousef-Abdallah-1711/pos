import type { Locale } from '@/content/site'
import offer from '@/content/commercial-offer.json'

export function BillingDetails({ locale }: { locale: Locale }) {
  const ar = locale === 'ar'
  const format = (value: number) => new Intl.NumberFormat(ar ? 'ar-EG' : 'en-EG', { maximumFractionDigits: 2 }).format(value)
  const monthly = offer.monthly_price_cents / 100
  const annual = offer.annual_price_cents / 100
  const comparison = monthly * 12
  const items = ar ? [
    ['تبدأ التجربة عند إنشاء الحساب', `تُنشئ حسابًا للمطعم لمدة ${format(offer.trial_days)} يومًا. يحدد الخادم موعد انتهاء التجربة، ويُطبّق حد الوصول عند بلوغ هذا الموعد.`],
    ['الشهري والسنوي: الاشتراك نفسه', `${format(monthly)} جنيه شهريًا، أو ${format(annual)} جنيه تُفوَّتر سنويًا. لا يعني اختيار الفترة السنوية وجود فئة مزايا إضافية.`],
    ['كيف تُحسب المقارنة السنوية؟', `اثنتا عشرة دفعة شهرية تساوي ${format(comparison)} جنيه. الفرق مع السعر السنوي هو ${format(comparison - annual)} جنيه، أي حوالي ${format((comparison - annual) / comparison * 100)}٪. المتوسط الشهري المكافئ ${format(annual / 12)} جنيه، مع الفوترة السنوية.`],
    ['اختيار الفترة لا ينفّذ عملية دفع', 'يُحفظ تفضيل الفوترة مع حساب المطعم الجديد. الدفع المدفوع غير متاح حاليًا؛ وتُؤكَّد شروط الدفع والضرائب قبل أي تفعيل مدفوع. لا يغيّر هذا العرض تواريخ تجارب الحسابات الحالية أو سجلات مدفوعاتها.'],
  ] : [
    ['Your trial begins when the account is created', `Create a restaurant workspace with a ${format(offer.trial_days)}-day trial. The server determines the expiry time and enforces the access boundary when that time is reached.`],
    ['Monthly and annual are the same subscription', `${format(monthly)} EGP billed monthly, or ${format(annual)} EGP billed annually. Choosing annual billing does not imply an additional feature tier.`],
    ['How the annual comparison works', `Twelve monthly payments total ${format(comparison)} EGP. The annual price saves ${format(comparison - annual)} EGP, approximately ${format((comparison - annual) / comparison * 100)}%. Its equivalent monthly average is ${format(annual / 12)} EGP, billed annually.`],
    ['A billing preference is not a payment', 'Your preference is saved with the new restaurant account. Paid checkout is currently unavailable; payment and tax terms must be confirmed before paid activation. This offer does not rewrite existing trial dates or payment records.'],
  ]
  return <section className="sv-section sv-customization"><div className="sv-shell">
    <p className="sv-eyebrow">{ar ? 'تفاصيل الاشتراك' : 'Subscription details'}</p>
    <h2>{ar ? 'اعرف ما تبدأ به، وما يأتي بعده.' : 'Know what starts now, and what comes later.'}</h2>
    <div className="sv-feature-grid">{items.map(([title, body]) => <article className="sv-feature-card" key={title}><h3>{title}</h3><p>{body}</p></article>)}</div>
  </div></section>
}
