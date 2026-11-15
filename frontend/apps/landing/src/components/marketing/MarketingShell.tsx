import { SiteSchema } from './SiteSchema'
import type { ReactNode } from 'react'
import { Navbar } from '@/components/Navbar'
import { Footer } from '@/components/Footer'
import { SkipLink } from '@/components/SkipLink'
import type { Locale } from '@/content/site'

export function MarketingShell({ children, locale }: { children: ReactNode; locale: Locale }) {
  return <><SkipLink locale={locale} /><Navbar locale={locale} /><SiteSchema /><main id="main-content" className="sv-main">{children}</main><Footer locale={locale} /></>
}
