import { getSiteCopy, type Locale } from '@/content/site'
import offer from '@/content/commercial-offer.json'

function format(locale: Locale, value: number): string {
  return new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-EG', { maximumFractionDigits: 2 }).format(value)
}

export function BillingDetails({ locale }: { locale: Locale }) {
  const copy = getSiteCopy(locale)
  const values = {
    trialDays: format(locale, offer.trial_days),
    monthly: format(locale, offer.monthly_price_cents / 100),
    annual: format(locale, offer.annual_price_cents / 100),
    comparison: format(locale, (offer.monthly_price_cents / 100) * 12),
    savings: format(locale, ((offer.monthly_price_cents / 100) * 12) - (offer.annual_price_cents / 100)),
    percent: format(locale, ((((offer.monthly_price_cents / 100) * 12) - (offer.annual_price_cents / 100)) / ((offer.monthly_price_cents / 100) * 12)) * 100),
    average: format(locale, (offer.annual_price_cents / 100) / 12),
  }
  const items = copy.billingDetails.map(item => ({
    title: item.title,
    body: item.body.replace(/\{(trialDays|monthly|annual|comparison|savings|percent|average)\}/g,
      (_token, key: keyof typeof values) => values[key]),
  }))
  return <section className="sv-section sv-customization"><div className="sv-shell">
    <p className="sv-eyebrow">{copy.billingDetailsEyebrow}</p>
    <h2>{copy.billingDetailsTitle}</h2>
    <div className="sv-feature-grid">{items.map(item => <article className="sv-feature-card" key={item.title}>
      <h3>{item.title}</h3><p>{item.body}</p>
    </article>)}</div>
  </div></section>
}
