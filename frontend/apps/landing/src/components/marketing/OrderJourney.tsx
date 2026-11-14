import { FeatureActions } from './FeatureActions'
import { getSiteCopy, type Locale } from '@/content/site'
import { SectionHead } from './SectionHead'

export function OrderJourney({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)
  return <section id="order-journey" className="sv-section sv-journey">
    <div className="sv-shell sv-journey-layout">
      <SectionHead eyebrow={c.journeyEyebrow} title={c.journeyTitle} body={c.journeyBody} />
      <div className="sv-journey-steps">{c.journey.map(([number, title, body], index) =>
        <details key={number} open={index === 0}>
          <summary><span className="sv-step-number">{number}</span><h3>{title}</h3><span aria-hidden="true" className="sv-expand">+</span></summary>
          <p>{body}</p>
        </details>)}
        <FeatureActions locale={locale} source="order-journey" destination="how-it-works" label={c.how} />
      </div>
    </div>
  </section>
}
