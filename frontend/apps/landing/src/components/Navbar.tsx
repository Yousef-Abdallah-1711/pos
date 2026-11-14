'use client'

import { useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import { Globe2, Menu, X } from 'lucide-react'
import { useLanguage } from '@/contexts/LanguageContext'
import { getSiteCopy } from '@/content/site'
import { BrandMark } from './BrandMark'
import { TrialLink } from './marketing/TrialLink'

export function Navbar() {
  const { language, setLanguage } = useLanguage()
  const [open, setOpen] = useState(false)
  const trigger = useRef<HTMLButtonElement>(null)
  const nav = useRef<HTMLDivElement>(null)
  const c = getSiteCopy(language)
  const links = [['features', c.platform], ['how-it-works', c.how], ['pricing', c.pricing]]
  useEffect(() => {
    if (!open) return
    nav.current?.querySelector<HTMLAnchorElement>('a')?.focus()
    const escape = (event: KeyboardEvent) => {
      if (event.key !== 'Escape') return
      setOpen(false)
      trigger.current?.focus()
    }
    document.addEventListener('keydown', escape)
    return () => document.removeEventListener('keydown', escape)
  }, [open])
  return <header className="sv-header"><nav className="sv-shell sv-nav" aria-label={language === 'ar' ? 'التنقل الرئيسي' : 'Main navigation'}>
    <BrandMark href={`/${language}`} label="SERVIO" />
    <div className="sv-desktop-links">{links.map(([path, label]) => <Link key={path} href={`/${language}/${path}`}>{label}</Link>)}</div>
    <div className="sv-nav-actions"><button className="sv-language" onClick={() => { setLanguage(language === 'en' ? 'ar' : 'en'); setOpen(false) }} aria-label={language === 'en' ? 'Switch to Arabic' : 'التغيير إلى الإنجليزية'}><Globe2 size={18} aria-hidden="true" /><span>{language === 'en' ? 'العربية' : 'English'}</span></button>
      <Link className="sv-login" href={`/${language}/auth/login`}>{c.login}</Link>
      <div className="sv-desktop-trial"><TrialLink locale={language} source="navigation" /></div>
      <button className="sv-menu-button" ref={trigger} type="button" aria-controls="sv-mobile-navigation" aria-expanded={open} aria-label={open ? (language === 'ar' ? 'إغلاق القائمة' : 'Close menu') : (language === 'ar' ? 'فتح القائمة' : 'Open menu')} onClick={() => setOpen(!open)}>{open ? <X aria-hidden="true" /> : <Menu aria-hidden="true" />}</button>
    </div>
  </nav><div hidden={!open} className="sv-mobile-navigation sv-shell" id="sv-mobile-navigation" ref={nav}>
    {links.map(([path, label]) => <Link key={path} href={`/${language}/${path}`} onClick={() => setOpen(false)}>{label}</Link>)}
    <Link href={`/${language}/auth/login`} onClick={() => setOpen(false)}>{c.login}</Link>
    <TrialLink locale={language} source="mobile-navigation" />
  </div></header>
}
