'use client'

import Link from 'next/link'
import { useLanguage } from '@/contexts/LanguageContext'
import { getSiteCopy } from '@/content/site'
import { BrandMark } from './BrandMark'

export function Footer() {
  const { language } = useLanguage()
  const c = getSiteCopy(language)
  const links = [['features', c.platform], ['how-it-works', c.how], ['pricing', c.pricing], ['demo', c.demo], ['contact', c.contact], ['privacy', c.privacy]]
  return <footer className="sv-footer"><div className="sv-shell">
    <div className="sv-footer-top"><div><BrandMark href={`/${language}`} label="SERVIO" inverse /><p>{c.footerLine}</p></div>
      <nav aria-label={language === 'ar' ? 'روابط المنتج' : 'Product links'}>{links.map(([path, label]) => <Link key={path} href={`/${language}/${path}`}>{label}</Link>)}<Link href={`/${language}/auth/login`}>{c.login}</Link><Link href={`/${language}/auth/register`}>{c.trial}</Link></nav></div>
    <div className="sv-footer-bottom"><p>{c.sampleNote}</p><small>© {new Date().getFullYear()} SERVIO</small></div>
  </div></footer>
}
