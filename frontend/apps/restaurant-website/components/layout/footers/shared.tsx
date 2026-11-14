'use client'

import Link from 'next/link'
import { ChevronUp, Facebook, Globe, Instagram, MessageCircle } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import type { FooterConfig } from '@pos-saas/ui-themes'

export interface FooterVariantProps {
  config: FooterConfig
  locale: 'en' | 'ar'
  className?: string
}

function socialIcon(platform: string) {
  const p = platform.toLowerCase()
  if (p.includes('facebook')) return Facebook
  if (p.includes('instagram')) return Instagram
  if (p.includes('whatsapp')) return MessageCircle
  return Globe
}

function withLocale(href: string, locale: 'en' | 'ar') {
  if (!href) return '#'
  if (href.startsWith('http')) return href
  if (/^\/(en|ar)(\/|$)/.test(href)) return href
  return `/${locale}${href.startsWith('/') ? '' : '/'}${href}`
}

export function FooterSocial({ config }: { config: FooterConfig }) {
  if (!config.socialLinks || config.socialLinks.length === 0) return null
  return (
    <div className="flex items-center gap-3">
      {config.socialLinks.map((link) => {
        const Icon = socialIcon(link.platform)
        return (
          <Link
            key={link.id || link.platform}
            href={link.url}
            target="_blank"
            rel="noreferrer"
            className="rounded-full border border-white/20 p-2 transition hover:bg-white/10"
            aria-label={link.platform}
          >
            <Icon className="h-4 w-4" />
          </Link>
        )
      })}
    </div>
  )
}

export function FooterLegal({ config, locale }: { config: FooterConfig; locale: 'en' | 'ar' }) {
  if (!config.legalLinks || config.legalLinks.length === 0) return null
  return (
    <div className="flex flex-wrap items-center gap-3 text-xs opacity-90">
      {config.legalLinks.map((link) => (
        <Link key={`${link.href}-${link.label}`} href={withLocale(link.href, locale)}>
          {link.label}
        </Link>
      ))}
    </div>
  )
}

export function BackToTop({ locale }: { locale: 'en' | 'ar' }) {
  return (
    <Button
      type="button"
      variant="outline"
      size="sm"
      onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
      className="h-9 gap-1 border-white/20 bg-transparent text-current hover:bg-white/10"
    >
      <ChevronUp className="h-4 w-4" />
      {locale === 'ar' ? 'الأعلى' : 'Top'}
    </Button>
  )
}

export function FooterShell({ config, locale, className, children }: FooterVariantProps & { children: React.ReactNode }) {
  return (
    <footer
      className={cn(
        'mt-auto w-full border-t border-theme-border/40 bg-[hsl(var(--theme-footer-bg))] text-[hsl(var(--theme-footer-text))]',
        className
      )}
      dir={locale === 'ar' ? 'rtl' : 'ltr'}
    >
      <div className="mx-auto w-full max-w-7xl px-4 py-10 sm:px-6 lg:px-8">{children}</div>
    </footer>
  )
}

export function FooterLinkList({
  title,
  links,
  locale,
}: {
  title: string
  links: Array<{ label: string; href: string }>
  locale: 'en' | 'ar'
}) {
  return (
    <div>
      <h3 className="mb-3 text-sm font-semibold">{title}</h3>
      <ul className="space-y-2 text-sm opacity-90">
        {links.map((link) => (
          <li key={`${link.href}-${link.label}`}>
            <Link href={withLocale(link.href, locale)} className="hover:underline">
              {link.label}
            </Link>
          </li>
        ))}
      </ul>
    </div>
  )
}
