import Link from 'next/link'
import Image from 'next/image'
import { ArrowUpRight, Check } from 'lucide-react'
import type { Locale } from '@/content/site'
import { getConversionCopy } from '@/content/conversionCopy'
import { ProductCapture } from './ProductCapture'
import { TrialLink } from './TrialLink'
import { OrderStory } from './OrderStory'
import { HeroCaptureMotion } from './HeroCaptureMotion'
import { WorkflowDemo } from './WorkflowDemo'
import { FeatureShowcase } from './FeatureShowcase'
import { LayeredTypography } from './LayeredTypography'
import { BusinessEcosystem } from './BusinessEcosystem'
import { ManagementExperience } from './ManagementExperience'
import { StackedSetupCards } from './StackedSetupCards'
import { LandingPricing } from './LandingPricing'
import { OperationalBento } from './OperationalBento'

export function ConversionLanding({ locale }: { locale: Locale }) {
  const c = getConversionCopy(locale)

  return <div className="cv-page">
    <section className="cv-hero" aria-labelledby="cv-hero-title">
      <Image className="cv-hero-photo" src="/images/marketing/cafe-interior.jpg" alt="" fill priority fetchPriority="high" quality={50} sizes="100vw" aria-hidden="true" />
      <div className="cv-wrap cv-hero-inner">
        <p className="cv-eyebrow">{c.heroEyebrow}</p>
        <h1 id="cv-hero-title">{c.heroTitle}</h1>
        <p className="cv-hero-body">{c.heroBody}</p>
        <div className="cv-actions">
          <TrialLink locale={locale} source="home-hero" label={c.createAccount} />
          <a className="cv-button-outline" href="#how-orders-move">{c.seeFlow}<ArrowUpRight size={17} aria-hidden="true" /></a>
        </div>
        <HeroCaptureMotion><ProductCapture image="dashboard" alt={c.dashboardAlt} locale={locale} className="cv-hero-capture" priority /></HeroCaptureMotion>
        <p className="cv-photo-credit">{locale === 'ar' ? 'صورة المقهى: Julen Rey Azcona' : 'Café photo: Julen Rey Azcona'}</p>
      </div>
    </section>

    <OperationalBento locale={locale} />

    <OrderStory locale={locale} />

    <section id="online-ordering" className="cv-section cv-online" aria-labelledby="cv-online-title">
      <div className="cv-wrap cv-split">
        <div className="cv-split-copy">
          <p className="cv-eyebrow">{c.onlineEyebrow}</p>
          <h2 id="cv-online-title">{c.onlineTitle}</h2>
          <p>{c.onlineBody}</p>
          <ul className="cv-checklist">
            {[c.onlinePointOne, c.onlinePointTwo, c.onlinePointThree].map(point => <li key={point}><Check size={19} aria-hidden="true" />{point}</li>)}
          </ul>
          <Link className="cv-text-link" href={`/${locale}/features/online-ordering`}>{c.onlineLink}<ArrowUpRight size={18} aria-hidden="true" /></Link>
        </div>
        <WorkflowDemo locale={locale} kind="online" />
      </div>
    </section>

    <section id="kitchen" className="cv-section cv-kitchen" aria-labelledby="cv-kitchen-title">
      <div className="cv-wrap cv-split cv-split-reverse">
        <div className="cv-split-copy">
          <p className="cv-eyebrow">{c.kitchenEyebrow}</p>
          <h2 id="cv-kitchen-title">{c.kitchenTitle}</h2>
          <p>{c.kitchenBody}</p>
          <Link className="cv-text-link" href={`/${locale}/features/kitchen-display`}>{c.kitchenLink}<ArrowUpRight size={18} aria-hidden="true" /></Link>
        </div>
        <WorkflowDemo locale={locale} kind="kitchen" />
      </div>
    </section>

    <LayeredTypography locale={locale} />
    <FeatureShowcase locale={locale} />
    <BusinessEcosystem locale={locale} />
    <ManagementExperience locale={locale} />

    <StackedSetupCards locale={locale} />
    <LandingPricing locale={locale} />

    <section className="cv-section cv-faq" aria-labelledby="cv-faq-title">
      <div className="cv-wrap cv-faq-layout">
        <div><p className="cv-eyebrow">{c.faqEyebrow}</p><h2 id="cv-faq-title">{c.faqTitle}</h2></div>
        <div className="cv-faq-list">{c.faqs.map(item => <details key={item.question}>
          <summary>{item.question}</summary><p>{item.answer}</p>
        </details>)}</div>
      </div>
    </section>

    <section className="cv-close" aria-labelledby="cv-close-title">
      <div className="cv-wrap"><h2 id="cv-close-title">{c.closeTitle}</h2><p>{c.closeBody}</p>
        <TrialLink locale={locale} source="home-final" label={c.createAccount} /></div>
    </section>
  </div>
}
