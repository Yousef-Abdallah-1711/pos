'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import { ArrowUpRight, ChevronDown, Globe2, Moon, Sun, X } from 'lucide-react'
import type { RefObject } from 'react'
import type { Locale } from '@/content/site'
import { getSiteCopy } from '@/content/site'
import { FEATURE_GROUPS, HEADER_LINKS, isNavigationRouteActive, localizedHref } from '@/content/navigation'
import { TrialLink } from '@/components/marketing/TrialLink'

interface MobileNavigationProps {
  locale: Locale
  pathname: string
  open: boolean
  onClose: () => void
  triggerRef: RefObject<HTMLButtonElement | null>
  theme: string
  onToggleTheme: () => void
  onSwitchLocale: () => void
}

export function MobileNavigation({ locale, pathname, open, onClose, triggerRef, theme, onToggleTheme, onSwitchLocale }: MobileNavigationProps) {
  const copy = getSiteCopy(locale)
  const panelRef = useRef<HTMLDivElement>(null)
  const [featuresOpen, setFeaturesOpen] = useState(false)
  const closeAndRestoreFocus = useCallback(() => {
    onClose()
    window.requestAnimationFrame(() => triggerRef.current?.focus())
  }, [onClose, triggerRef])

  useEffect(() => {
    if (!open) return
    const previousOverflow = document.body.style.overflow
    const inactiveElements = [
      document.querySelector<HTMLElement>('.sv-nav'),
      document.querySelector<HTMLElement>('main'),
      document.querySelector<HTMLElement>('.sv-footer'),
    ].filter((element): element is HTMLElement => element !== null)
    const previousInertState = inactiveElements.map(element => ({
      element,
      inert: element.inert,
      ariaHidden: element.getAttribute('aria-hidden'),
    }))
    document.body.style.overflow = 'hidden'
    for (const element of inactiveElements) {
      element.inert = true
      element.setAttribute('aria-hidden', 'true')
    }
    panelRef.current?.querySelector<HTMLElement>('a, button')?.focus()
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        event.preventDefault()
        closeAndRestoreFocus()
        return
      }
      if (event.key !== 'Tab' || !panelRef.current) return
      const controls = [...panelRef.current.querySelectorAll<HTMLElement>('a[href], button:not([disabled])')]
      const first = controls[0]
      const last = controls.at(-1)
      if (event.shiftKey && document.activeElement === first) {
        event.preventDefault()
        last?.focus()
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault()
        first?.focus()
      }
    }
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.body.style.overflow = previousOverflow
      for (const { element, inert, ariaHidden } of previousInertState) {
        element.inert = inert
        if (ariaHidden === null) element.removeAttribute('aria-hidden')
        else element.setAttribute('aria-hidden', ariaHidden)
      }
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [open, closeAndRestoreFocus])

  return <>
    <button className="sv-mobile-backdrop" type="button" aria-hidden="true" tabIndex={-1} hidden={!open} onClick={closeAndRestoreFocus} />
    <div className="sv-mobile-navigation" id="sv-mobile-navigation" ref={panelRef} role="dialog" aria-modal="true" aria-label={copy.navigation.mobileNavigation} hidden={!open}>
      <div className="sv-mobile-navigation__heading">
        <span>{copy.navigation.mobileNavigation}</span>
        <button className="sv-mobile-navigation__close" type="button" onClick={closeAndRestoreFocus} aria-label={copy.closeMenu}><X size={20} aria-hidden="true" /></button>
      </div>
      <nav aria-label={copy.navigation.mobileNavigation}>
        <div className="sv-mobile-navigation__group">
          <button
            className="sv-mobile-navigation__accordion"
            type="button"
            aria-expanded={featuresOpen}
            aria-label={featuresOpen ? copy.navigation.collapseFeatures : copy.navigation.expandFeatures}
            aria-current={isNavigationRouteActive(pathname, 'features', 'features') ? 'page' : undefined}
            aria-controls="sv-mobile-feature-links"
            onClick={() => setFeaturesOpen(value => !value)}>
            <span>{copy.navigation.platform}</span><ChevronDown size={18} aria-hidden="true" />
          </button>
          <div className="sv-mobile-feature-links" id="sv-mobile-feature-links" hidden={!featuresOpen}>
            <Link href={localizedHref(locale, 'features')} aria-current={isNavigationRouteActive(pathname, 'features') ? 'page' : undefined} onClick={closeAndRestoreFocus}>{copy.navigation.allFeatures}</Link>
            {FEATURE_GROUPS.map(group => <section key={group.id}>
              <h2>{copy.navigation[group.label]}</h2>
              {group.items.map(item => <Link key={item.id} href={localizedHref(locale, `features/${item.id}`)} aria-current={pathname === localizedHref(locale, `features/${item.id}`) ? 'page' : undefined} onClick={closeAndRestoreFocus}>
                <span>{copy.navigation[item.label]}</span><span className="sv-mobile-feature-links__description">{copy.navigation[item.description]}</span>
              </Link>)}
            </section>)}
          </div>
        </div>
        {HEADER_LINKS.filter(item => item.route !== 'features').map(item => <Link className="sv-mobile-navigation__link" key={item.id} href={localizedHref(locale, item.route)} onClick={closeAndRestoreFocus}>
          {copy.navigation[item.label]}<ArrowUpRight size={15} aria-hidden="true" />
        </Link>)}
        <Link className="sv-mobile-navigation__link" href={localizedHref(locale, 'privacy')} onClick={closeAndRestoreFocus}>{copy.navigation.privacy}</Link>
        <Link className="sv-mobile-navigation__link" href={localizedHref(locale, 'login')} prefetch={false} onClick={closeAndRestoreFocus}>{copy.navigation.login}</Link>
      </nav>
      <div className="sv-mobile-navigation__controls">
        <button type="button" onClick={() => { onSwitchLocale(); closeAndRestoreFocus() }} aria-label={locale === 'en' ? copy.switchToArabic : copy.switchToEnglish}>
          <Globe2 size={18} aria-hidden="true" /><span>{locale === 'en' ? 'العربية' : 'English'}</span>
        </button>
        <button type="button" onClick={onToggleTheme} aria-label={theme === 'dark' ? copy.switchToLight : copy.switchToDark}>
          {theme === 'dark' ? <Sun size={18} aria-hidden="true" /> : <Moon size={18} aria-hidden="true" />}<span>{theme === 'dark' ? copy.switchToLight : copy.switchToDark}</span>
        </button>
      </div>
      <div onClick={closeAndRestoreFocus}><TrialLink locale={locale} source="mobile-navigation" /></div>
    </div>
  </>
}
