'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import { ArrowUpRight, ChevronDown } from 'lucide-react'
import type { Locale } from '@/content/site'
import { getSiteCopy } from '@/content/site'
import { FEATURE_GROUPS, isNavigationRouteActive, localizedHref } from '@/content/navigation'
import { TrialLink } from '@/components/marketing/TrialLink'
import { navigationIcons } from './NavigationIcons'

export function DesktopMegaMenu({ locale, pathname }: { locale: Locale; pathname: string }) {
  const copy = getSiteCopy(locale)
  const [open, setOpen] = useState(false)
  const rootRef = useRef<HTMLDivElement>(null)
  const triggerRef = useRef<HTMLButtonElement>(null)
  const openedByPointer = useRef(false)
  const closeMenu = useCallback(() => {
    openedByPointer.current = false
    setOpen(false)
  }, [])

  useEffect(() => {
    if (!open) return
    const onPointerDown = (event: PointerEvent) => {
      if (event.target instanceof Node && !rootRef.current?.contains(event.target)) closeMenu()
    }
    const onFocusIn = (event: FocusEvent) => {
      if (event.target instanceof Node && !rootRef.current?.contains(event.target)) closeMenu()
    }
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key !== 'Escape') return
      event.preventDefault()
      closeMenu()
      triggerRef.current?.focus()
    }
    document.addEventListener('pointerdown', onPointerDown)
    document.addEventListener('focusin', onFocusIn)
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('pointerdown', onPointerDown)
      document.removeEventListener('focusin', onFocusIn)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [open, closeMenu])

  return <div
    className="sv-platform-menu"
    ref={rootRef}
    onPointerEnter={event => {
      if (event.pointerType === 'mouse' || event.pointerType === 'pen') {
        openedByPointer.current = true
        setOpen(true)
      }
    }}
    onPointerLeave={event => {
      const focusRemainsInPanel = document.activeElement !== triggerRef.current && rootRef.current?.contains(document.activeElement)
      if (event.pointerType === 'mouse' || event.pointerType === 'pen') {
        openedByPointer.current = false
        if (!focusRemainsInPanel) closeMenu()
      }
    }}>
    <button
      className="sv-platform-trigger"
      type="button"
      aria-expanded={open}
      aria-controls="sv-platform-mega-menu"
      aria-label={open ? copy.navigation.closePlatformMenu : copy.navigation.openPlatformMenu}
      aria-current={isNavigationRouteActive(pathname, 'features', 'features') ? 'page' : undefined}
      onClick={() => {
        if (openedByPointer.current) {
          openedByPointer.current = false
          setOpen(true)
        } else setOpen(value => !value)
      }}
      ref={triggerRef}>
      {copy.navigation.platform}<ChevronDown size={16} aria-hidden="true" />
    </button>
    <div className="sv-mega-menu" id="sv-platform-mega-menu" hidden={!open}>
      <div className="sv-mega-menu__groups">
        {FEATURE_GROUPS.map(group => <section className="sv-mega-group" key={group.id} aria-labelledby={`sv-mega-${group.id}`}>
          <h2 id={`sv-mega-${group.id}`}>{copy.navigation[group.label]}</h2>
          <ul>
            {group.items.map(item => {
              const Icon = navigationIcons[item.icon]
              const href = localizedHref(locale, `features/${item.id}`)
              return <li key={item.id}>
                <Link href={href} prefetch={false} aria-current={pathname === href ? 'page' : undefined} onClick={closeMenu}>
                  <span className="sv-mega-item__icon"><Icon size={19} aria-hidden="true" /></span>
                  <span className="sv-mega-item__copy">
                    <span className="sv-mega-item__title">{copy.navigation[item.label]}</span>
                    <span className="sv-mega-item__description">{copy.navigation[item.description]}</span>
                  </span>
                  <ArrowUpRight size={15} aria-hidden="true" className="sv-mega-item__arrow" />
                </Link>
              </li>
            })}
          </ul>
        </section>)}
      </div>
      <div className="sv-mega-menu__footer">
        <p>{copy.navigation.menuFooter}</p>
        <div>
          <Link href={localizedHref(locale, 'features')} prefetch={false} onClick={closeMenu}>{copy.navigation.allFeatures}<ArrowUpRight size={15} aria-hidden="true" /></Link>
          <Link href={localizedHref(locale, 'productTour')} prefetch={false} onClick={closeMenu}>{copy.navigation.productTour}<ArrowUpRight size={15} aria-hidden="true" /></Link>
          <TrialLink locale={locale} source="mega-menu" />
        </div>
      </div>
    </div>
  </div>
}
