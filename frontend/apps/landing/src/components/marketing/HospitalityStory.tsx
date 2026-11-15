import Image from 'next/image'
import { getSiteCopy, type Locale } from '@/content/site'

export function HospitalityStory({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)

  return <section className="sv-section sv-hospitality-story">
    <div className="sv-shell sv-hospitality-story__layout">
      <header className="sv-section-head">
        <p className="sv-eyebrow">{c.hospitalityEyebrow}</p>
        <h2>{c.hospitalityTitle}</h2>
        <p className="sv-lead">{c.hospitalityBody}</p>
      </header>
      <figure className="sv-hospitality-story__figure">
        <div className="sv-hospitality-story__image">
          <Image
            src="/images/marketing/cafe-interior.jpg"
            alt={c.hospitalityImageAlt}
            width={1400}
            height={2100}
            sizes="(max-width: 767px) calc(100vw - 40px), (max-width: 1359px) 45vw, 580px"
            loading="lazy"
          />
        </div>
        <figcaption>{c.hospitalityPhotoCredit} <a href="https://unsplash.com/photos/sunlight-streams-into-a-modern-cafe-interior-U_lZQfczRRw" target="_blank" rel="noreferrer">Unsplash</a></figcaption>
      </figure>
    </div>
  </section>
}
