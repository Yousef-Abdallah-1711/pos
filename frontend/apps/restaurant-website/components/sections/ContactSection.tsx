'use client'

import Link from 'next/link'
import { Clock3, MapPin, MessageCircle, Phone } from 'lucide-react'
import type { ThemeLocale, ThemeSectionConfig } from '@/types/theme'
import { readLocalized, toRecord, toStringValue } from '@/components/sections/section-utils'

interface ContactSectionProps {
  locale: ThemeLocale
  sectionConfig?: ThemeSectionConfig
}

export function ContactSection({ locale, sectionConfig }: ContactSectionProps) {
  const settings = toRecord(sectionConfig?.settings)
  const title = readLocalized(sectionConfig?.title ?? settings.title, locale, locale === 'ar' ? 'تواصل معنا' : 'Contact')
  const phone = toStringValue(settings.phone, '+201000000000')
  const whatsapp = toStringValue(settings.whatsapp_number ?? settings.whatsappNumber, phone)
  const address = readLocalized(settings.address, locale, locale === 'ar' ? 'عنوان المطعم' : 'Restaurant Address')
  const hours = readLocalized(settings.hours, locale, locale === 'ar' ? 'يوميا من 10 ص إلى 12 م' : 'Daily 10 AM - 12 AM')
  const mapsUrl = toStringValue(settings.maps_url ?? settings.mapUrl, `https://maps.google.com/?q=${encodeURIComponent(address)}`)
  const waText = encodeURIComponent(locale === 'ar' ? 'مرحبا، أريد الاستفسار عن الطلب' : 'Hello, I want to ask about my order')

  const cards = [
    {
      icon: Phone,
      label: locale === 'ar' ? 'الهاتف' : 'Phone',
      value: phone,
      href: `tel:${phone}`,
      action: locale === 'ar' ? 'اتصال' : 'Call',
    },
    {
      icon: MessageCircle,
      label: locale === 'ar' ? 'واتساب' : 'WhatsApp',
      value: whatsapp,
      href: `https://wa.me/${whatsapp}?text=${waText}`,
      action: locale === 'ar' ? 'محادثة' : 'Chat',
    },
    {
      icon: MapPin,
      label: locale === 'ar' ? 'الموقع' : 'Location',
      value: address,
      href: mapsUrl,
      action: locale === 'ar' ? 'الخريطة' : 'Map',
    },
    {
      icon: Clock3,
      label: locale === 'ar' ? 'ساعات العمل' : 'Hours',
      value: hours,
      href: `/${locale}/menu`,
      action: locale === 'ar' ? 'اطلب الآن' : 'Order',
    },
  ]

  return (
    <section className="bg-[hsl(var(--theme-surface))] px-4 py-16 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-6xl">
        <h2 className={`mb-8 text-3xl font-bold text-[hsl(var(--theme-text))] ${locale === 'ar' ? 'text-right' : 'text-left'}`}>{title}</h2>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {cards.map((card) => (
            <article key={card.label} className="rounded-xl border border-[hsl(var(--theme-border))] bg-[hsl(var(--theme-background))] p-5">
              <div className="mb-4 flex items-center gap-3">
                <card.icon className="h-5 w-5 text-[hsl(var(--theme-primary))]" />
                <h3 className="text-base font-semibold text-[hsl(var(--theme-text))]">{card.label}</h3>
              </div>
              <p className="mb-4 text-sm text-[hsl(var(--theme-text))]/80">{card.value}</p>
              <Link
                href={card.href}
                target={card.href.startsWith('http') ? '_blank' : undefined}
                rel={card.href.startsWith('http') ? 'noreferrer noopener' : undefined}
                className="inline-flex min-h-10 items-center rounded-md bg-[hsl(var(--theme-primary))] px-4 py-2 text-sm font-medium text-white"
              >
                {card.action}
              </Link>
            </article>
          ))}
        </div>
      </div>
    </section>
  )
}

