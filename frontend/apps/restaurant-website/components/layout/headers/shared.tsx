'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Menu } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Sheet,
  SheetContent,
  SheetTrigger,
} from '@/components/ui/sheet'
import type { HeaderConfig } from '@pos-saas/ui-themes'
import { CartButton } from '@/components/layout/header/CartButton'
import { DarkModeToggle } from '@/components/layout/header/DarkModeToggle'
import { LanguageToggle } from '@/components/layout/header/LanguageToggle'
import { SearchBar } from '@/components/layout/header/SearchBar'
import { WhatsAppButton } from '@/components/layout/header/WhatsAppButton'
import { CartSheet } from '@/components/commerce/CartSheet'

export interface HeaderVariantProps {
  config: HeaderConfig
  locale: 'en' | 'ar'
  restaurantSlug?: string
  whatsAppNumber?: string
  showWhatsAppButton?: boolean
  className?: string
  innerClassName?: string
}

function hrefWithLocale(href: string, locale: 'en' | 'ar') {
  if (!href) return `/${locale}`
  if (href.startsWith('http')) return href
  if (/^\/(en|ar)(\/|$)/.test(href)) return href
  if (href === '/') return `/${locale}`
  return `/${locale}${href.startsWith('/') ? '' : '/'}${href}`
}

function NavLinks({
  config,
  locale,
  mobile = false,
}: {
  config: HeaderConfig
  locale: 'en' | 'ar'
  mobile?: boolean
}) {
  const navItems = config.navigationItems || []
  return (
    <nav className={cn('flex items-center gap-4', mobile && 'flex-col items-start gap-3')}>
      {navItems.map((item) => (
        <Link
          key={item.id || item.href}
          href={hrefWithLocale(item.href, locale)}
          className={cn(
            'text-sm font-medium transition-opacity hover:opacity-80',
            mobile && 'text-base'
          )}
        >
          {item.label}
        </Link>
      ))}
    </nav>
  )
}

export function HeaderVariantShell({
  config,
  locale,
  whatsAppNumber,
  showWhatsAppButton,
  className,
  innerClassName,
}: HeaderVariantProps) {
  const [isCartOpen, setIsCartOpen] = useState(false)
  const safeLogoText = config.logoText || 'Restaurant'
  const mobileSide = locale === 'ar' ? 'right' : 'left'

  return (
    <header className={cn('sticky top-0 z-[60] w-full border-b border-theme-border/30 text-theme-header-text', className)}>
      <div
        className={cn(
          'mx-auto flex h-16 w-full max-w-7xl items-center justify-between gap-3 px-3 sm:px-4 lg:px-6',
          innerClassName
        )}
      >
        <div className={cn('flex min-w-0 items-center gap-2', locale === 'ar' && 'flex-row-reverse')}>
          <Sheet>
            <SheetTrigger asChild>
              <Button type="button" variant="ghost" size="icon" className="h-12 w-12 md:hidden">
                <Menu className="h-5 w-5" />
              </Button>
            </SheetTrigger>
            <SheetContent side={mobileSide} className="w-[85vw] max-w-sm p-6">
              <div className="mb-6 text-lg font-semibold">{safeLogoText}</div>
              <NavLinks config={config} locale={locale} mobile />
            </SheetContent>
          </Sheet>
          <Link href={`/${locale}`} className="min-w-0 truncate text-base font-semibold sm:text-lg">
            {safeLogoText}
          </Link>
        </div>

        <div className="hidden md:flex md:items-center md:gap-6">
          <NavLinks config={config} locale={locale} />
        </div>

        <div className={cn('flex shrink-0 items-center gap-1', locale === 'ar' && 'flex-row-reverse')}>
          <SearchBar />
          <LanguageToggle />
          <DarkModeToggle />
          <CartButton onOpenCart={() => setIsCartOpen(true)} />
          <WhatsAppButton
            phoneNumber={whatsAppNumber}
            show={showWhatsAppButton !== false}
          />
        </div>
      </div>
      <CartSheet locale={locale} open={isCartOpen} onOpenChange={setIsCartOpen} />
    </header>
  )
}
