'use client'

import { useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { Globe2, Menu, Moon, Sun, X } from 'lucide-react'
import { useTheme } from 'next-themes'
import { getSiteCopy, type Locale } from '@/content/site'
import { BrandMark } from './BrandMark'
import { TrialLink } from './marketing/TrialLink'

export function Navbar({ locale: language }: { locale: Locale }) {
  const { resolvedTheme, setTheme } = useTheme()
  const pathname = usePathname()
  const router = useRouter()
  const [mounted, setMounted] = useState(false)
  const [open, setOpen] = useState(false)
  const trigger = useRef<HTMLButtonElement>(null)
  const nav = useRef<HTMLDivElement>(null)
  const c = getSiteCopy(language)
  const nextLocale = language === 'en' ? 'ar' : 'en'
  const localizedPath = pathname.replace(/^\/(?:en|ar)(?=\/|$)/, `/${nextLocale}`)
  const theme = mounted ? resolvedTheme : 'dark'
  useEffect(() => setMounted(true), [])
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
  return <header className="sv-header"><nav className="sv-shell sv-nav" aria-label={c.navLabel}>
    <BrandMark href={`/${language}`} label="SERVIO" />
    <div className="sv-desktop-links">{links.map(([path, label]) => <Link key={path} href={`/${language}/${path}`}>{label}</Link>)}</div>
    <div className="sv-nav-actions"><Link className="sv-language" href={localizedPath} aria-label={c.switchLanguage} onClick={event => {
      event.preventDefault()
      setOpen(false)
      router.push(`${localizedPath}${window.location.search}${window.location.hash}`)
    }}><Globe2 size={18} aria-hidden="true" /><span>{language === 'en' ? 'العربية' : 'English'}</span></Link>
      <button className="sv-theme-toggle" type="button" onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')} aria-label={theme === 'dark' ? c.switchToLight : c.switchToDark} title={theme === 'dark' ? c.switchToLight : c.switchToDark}>{theme === 'dark' ? <Sun size={18} aria-hidden="true" /> : <Moon size={18} aria-hidden="true" />}</button>
      <Link className="sv-login" href={`/${language}/auth/login`} prefetch={false}>{c.login}</Link>
      <div className="sv-desktop-trial"><TrialLink locale={language} source="navigation" /></div>
      <button className="sv-menu-button" ref={trigger} type="button" aria-controls="sv-mobile-navigation" aria-expanded={open} aria-label={open ? c.closeMenu : c.openMenu} onClick={() => setOpen(!open)}>{open ? <X aria-hidden="true" /> : <Menu aria-hidden="true" />}</button>
    </div>
  </nav><div hidden={!open} className="sv-mobile-navigation sv-shell" id="sv-mobile-navigation" ref={nav}>
    {links.map(([path, label]) => <Link key={path} href={`/${language}/${path}`} onClick={() => setOpen(false)}>{label}</Link>)}
    <Link href={`/${language}/auth/login`} prefetch={false} onClick={() => setOpen(false)}>{c.login}</Link>
    <TrialLink locale={language} source="mobile-navigation" />
  </div></header>
}
