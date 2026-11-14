'use client'

import { useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import { getSiteCopy, type Locale } from '@/content/site'
import offer from '@/content/commercial-offer.json'
import { trackFunnel } from '@/lib/funnel'
import { TrialLink } from './TrialLink'

export function PricingOffer({ locale, detailed = false }: { locale: Locale; detailed?: boolean }) {
  const [billing, setBilling] = useState<'monthly' | 'annual'>('monthly')
  const section = useRef<HTMLElement>(null)
  const c = getSiteCopy(locale)
  const annual = billing === 'annual'
  const amount = (annual ? offer.annual_price_cents : offer.monthly_price_cents) / 100
  const price = new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-EG').format(amount)
  useEffect(() => {
    const element = section.current
    if (!element || typeof IntersectionObserver === 'undefined') return
    const observer = new IntersectionObserver(([entry]) => {
      if (!entry.isIntersecting) return
      trackFunnel('pricing_viewed', { locale, source: detailed ? 'pricing-page' : 'homepage' })
      observer.disconnect()
    }, { threshold: .25 })
    observer.observe(element)
    return () => observer.disconnect()
  }, [locale, detailed])
  const select = (value: 'monthly' | 'annual') => {
    if (value === billing) return
    setBilling(value)
    trackFunnel('billing_period_selected', { locale, source: detailed ? 'pricing-page' : 'homepage', billing: value })
  }
  return <section id="pricing" className="sv-section sv-pricing" ref={section}>
    <div className="sv-shell sv-pricing-layout"><div className="sv-price-story">
      <p className="sv-eyebrow">{c.offerEyebrow}</p>{detailed ? <h1>{c.offerTitle}</h1> : <h2>{c.offerTitle}</h2>}
      <p className="sv-lead">{c.offerBody}</p><p className="sv-note">{c.offerScope}</p>
      {!detailed && <Link className="sv-text-link" href={`/${locale}/pricing`}>{c.pricingLink}<span aria-hidden="true">↗</span></Link>}
    </div><div className="sv-offer">
      <fieldset className="sv-billing-selector"><legend className="sr-only">{c.billingPreference}</legend>
        {(['monthly', 'annual'] as const).map(value => <label key={value} data-selected={billing === value}><input type="radio" name={detailed ? 'page-billing' : 'home-billing'} value={value} checked={billing === value} onChange={() => select(value)} />{value === 'annual' ? c.annual : c.monthly}</label>)}
      </fieldset>
      <div className="sv-price-value" aria-live="polite" aria-atomic="true"><p><strong><bdi>{price}</bdi></strong><span>{locale === 'ar' ? 'جنيه مصري' : 'EGP'}</span></p><span>{annual ? c.yearUnit : c.monthUnit}</span>
        <p className="sv-billing-description">{annual ? c.annualBilling : c.monthlyBilling}</p>
        {annual && <div className="sv-savings"><strong>{c.annualSavings}</strong><p>{c.annualEquivalent}</p><p>{c.savingsBasis}</p></div>}
      </div>
      <TrialLink locale={locale} source={detailed ? 'pricing-page' : 'homepage-pricing'} billing={billing} />
      <p className="sv-offer-limit">{c.checkoutLimit}</p>
    </div></div>
  </section>
}
