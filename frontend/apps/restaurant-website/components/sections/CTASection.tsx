'use client'

import Link from 'next/link'
import type { ThemeLocale, ThemeSectionConfig } from '@/types/theme'
import { readLocalized, toRecord, toStringValue } from '@/components/sections/section-utils'

interface CTASectionProps {
  locale: ThemeLocale
  sectionConfig?: ThemeSectionConfig
}

export function CTASection({ locale, sectionConfig }: CTASectionProps) {
  const settings = toRecord(sectionConfig?.settings)
  const title = readLocalized(sectionConfig?.title ?? settings.title, locale, locale === 'ar' ? 'جاهز للطلب؟' : 'Ready to order?')
  const subtitle = readLocalized(settings.subtitle ?? settings.description, locale, locale === 'ar' ? 'اطلب الآن واستمتع بوجبتك.' : 'Place your order now and enjoy your meal.')
  const whatsappPhone = toStringValue(settings.whatsapp_number ?? settings.whatsappNumber, '')
  const menuHref = `/${locale}/menu`
  const waText = encodeURIComponent(locale === 'ar' ? 'مرحبا، أريد الطلب من المطعم' : 'Hello, I want to place an order')
  const waHref = whatsappPhone ? `https://wa.me/${whatsappPhone}?text=${waText}` : menuHref

  return (
    <section className="px-4 py-16 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-6xl rounded-2xl bg-[linear-gradient(135deg,hsl(var(--theme-primary))_0%,hsl(var(--theme-accent))_100%)] px-6 py-12 text-white md:px-10">
        <h2 className={`text-3xl font-bold ${locale === 'ar' ? 'text-right' : 'text-left'}`}>{title}</h2>
        <p className={`mt-3 text-white/90 ${locale === 'ar' ? 'text-right' : 'text-left'}`}>{subtitle}</p>
        <div className={`mt-6 flex flex-wrap gap-3 ${locale === 'ar' ? 'justify-end' : 'justify-start'}`}>
          <Link href={menuHref} className="inline-flex min-h-12 items-center rounded-md bg-white px-6 py-3 text-sm font-semibold text-black animate-pulse">
            {locale === 'ar' ? 'اطلب الآن' : 'Order Now'}
          </Link>
          <Link
            href={waHref}
            target={whatsappPhone ? '_blank' : undefined}
            rel={whatsappPhone ? 'noreferrer noopener' : undefined}
            className="inline-flex min-h-12 items-center rounded-md bg-[#25D366] px-6 py-3 text-sm font-semibold text-white"
          >
            {locale === 'ar' ? 'اطلب عبر واتساب' : 'Order via WhatsApp'}
          </Link>
        </div>
      </div>
    </section>
  )
}

