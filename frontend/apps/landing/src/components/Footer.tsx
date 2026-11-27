import Link from 'next/link'
import { getSiteCopy } from '@/content/site'
import type { Locale } from '@/content/site'
import { getFooterGroups, localizedHref } from '@/content/navigation'
import { BrandMark } from './BrandMark'
import { TrialLink } from './marketing/TrialLink'
import { getConversionCopy } from '@/content/conversionCopy'

export function Footer({ locale }: { locale: Locale }) {
  const copy = getSiteCopy(locale)
  const groups = getFooterGroups()
  return <footer className="sv-footer">
    <div className="sv-shell">
      <div className="sv-footer-top">
        <div className="sv-footer-brand">
          <BrandMark href={localizedHref(locale, 'home')} label="SERVIO" inverse />
          <p>{copy.navigation.brandDescription}</p>
          <TrialLink locale={locale} source="footer" label={getConversionCopy(locale).footerCta} />
        </div>
        <nav aria-label={copy.navigation.footerNavigation}>
          <div className="sv-footer-groups">
            {groups.map(group => <section className={`sv-footer-group sv-footer-group--${group.id}`} key={group.id}>
              <h2>{copy.navigation[group.label]}</h2>
              <ul>
                {group.links.map(item => <li key={item.id}>
                  <Link href={localizedHref(locale, item.route)} prefetch={item.route === 'login' || item.route === 'register' ? false : undefined}>
                    {copy.navigation[item.label]}
                  </Link>
                </li>)}
              </ul>
            </section>)}
          </div>
        </nav>
      </div>
      <div className="sv-footer-signature" aria-hidden="true">SERVIO</div>
      <div className="sv-footer-bottom">
        <p>{copy.sampleNote}</p>
        <div className="sv-footer-utility">
          <Link href={localizedHref(locale, 'privacy')}>{copy.navigation.privacy}</Link>
          <small>© {new Date().getFullYear()} SERVIO</small>
        </div>
      </div>
    </div>
  </footer>
}
