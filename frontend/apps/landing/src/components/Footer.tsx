import Link from 'next/link'
import { getSiteCopy } from '@/content/site'
import type { Locale } from '@/content/site'
import { BrandMark } from './BrandMark'

export function Footer({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)
  const language = locale
  const links = [['features', c.platform], ['how-it-works', c.how], ['pricing', c.pricing], ['demo', c.demo], ['contact', c.contact], ['privacy', c.privacy]]
  return <footer className="sv-footer"><div className="sv-shell">
    <div className="sv-footer-top"><div><BrandMark href={`/${language}`} label="SERVIO" inverse /><p>{c.footerLine}</p></div>
      <nav aria-label={language === 'ar' ? 'روابط المنتج' : 'Product links'}>{links.map(([path, label]) => <Link key={path} href={`/${language}/${path}`}>{label}</Link>)}<Link href={`/${language}/auth/login`} prefetch={false}>{c.login}</Link><Link href={`/${language}/auth/register`} prefetch={false}>{c.trial}</Link></nav></div>
    <div className="sv-footer-bottom"><p>{c.sampleNote}</p><small>© {new Date().getFullYear()} SERVIO</small></div>
  </div></footer>
}
