import { SiteSchema } from './SiteSchema'
import type { ReactNode } from 'react'
import { Navbar } from '@/components/Navbar'
import { Footer } from '@/components/Footer'
import { SkipLink } from '@/components/SkipLink'

export function MarketingShell({ children }: { children: ReactNode }) {
  return <><SkipLink /><Navbar /><SiteSchema /><main id="main-content" className="sv-main">{children}</main><Footer /></>
}
