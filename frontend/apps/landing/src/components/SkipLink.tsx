'use client'

import { useLanguage } from '@/contexts/LanguageContext'

export function SkipLink() {
  const { t } = useLanguage()
  return <a className="poos-skip-link" href="#main-content">{t.nav.skipToContent}</a>
}
