'use client'

import { FooterShell, FooterLegal, type FooterVariantProps } from './shared'

export function CardsFooter({ config, locale }: FooterVariantProps) {
  const cards = [
    { title: 'About', body: config.companyDescription || '' },
    { title: 'Contact', body: [config.phone, config.email].filter(Boolean).join(' | ') },
    { title: 'Menu', body: 'Explore categories and featured products' },
  ]

  return (
    <FooterShell config={config} locale={locale}>
      <div className="grid gap-4 md:grid-cols-3">
        {cards.map((card) => (
          <div key={card.title} className="rounded-xl border border-white/15 bg-white/5 p-4">
            <h3 className="text-sm font-semibold">{card.title}</h3>
            <p className="mt-2 text-sm opacity-90">{card.body}</p>
          </div>
        ))}
      </div>
      <div className="mt-6 flex flex-col gap-2 border-t border-white/10 pt-4 md:flex-row md:items-center md:justify-between">
        <p className="text-xs opacity-90">{config.copyrightText}</p>
        <FooterLegal config={config} locale={locale} />
      </div>
    </FooterShell>
  )
}
