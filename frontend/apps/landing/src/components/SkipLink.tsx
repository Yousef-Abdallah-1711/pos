import type { Locale } from '@/content/site'
import english from '@/i18n/messages/en/legacy.json'
import arabic from '@/i18n/messages/ar/legacy.json'

export function SkipLink({ locale }: { locale: Locale }) {
  const label = locale === 'ar' ? arabic.nav.skipToContent : english.nav.skipToContent
  return <a className="poos-skip-link" href="#main-content">{label}</a>
}
