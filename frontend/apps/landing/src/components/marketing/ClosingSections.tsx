import Link from 'next/link'
import { getSiteCopy, type Locale } from '@/content/site'
import { TrialLink } from './TrialLink'
import { SectionHead } from './SectionHead'
import { FeatureActions } from './FeatureActions'

export function GettingStarted({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)
  return <article id="getting-started" className="sv-section sv-setup"><div className="sv-shell">
    <SectionHead eyebrow={c.setupEyebrow} title={c.setupTitle} />
    <ol>{c.setup.map(([title, body], i) => <li key={title}><span>0{i + 1}</span><h3>{title}</h3><p>{body}</p></li>)}</ol>
    <div className="sv-setup-foot"><p className="sv-note">{c.setupNote}<br /><Link className="sv-text-link" href={`/${locale}/how-it-works`}>{c.how}</Link></p><TrialLink locale={locale} source="getting-started" /></div>
  </div></article>
}

export function MarketingFAQ({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)
  const questions = c.faqItems
  return <section id="faq" className="sv-section sv-faq"><div className="sv-shell sv-faq-layout"><SectionHead eyebrow={c.faqEyebrow} title={c.faqTitle} />
    <div>{questions.map(([question, answer]) => <details key={question}><summary>{question}<span aria-hidden="true">+</span></summary><p>{answer}</p></details>)}<FeatureActions locale={locale} source="faq" destination="contact" label={c.contact} /></div>
  </div></section>
}

export function FinalCTA({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)
  return <section id="start" className="sv-section sv-final"><div className="sv-shell">
    <p className="sv-eyebrow">{c.finalEyebrow}</p><h2>{c.finalTitle}</h2><div className="sv-final-bottom"><p>{c.finalBody}</p><div><TrialLink locale={locale} source="final" secondary /><br /><Link className="sv-text-link" href={`/${locale}/pricing`}>{c.pricingLink}</Link></div></div>
  </div></section>
}
