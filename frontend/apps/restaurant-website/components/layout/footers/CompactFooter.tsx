'use client'

import { FooterShell, FooterLegal, type FooterVariantProps } from './shared'

export function CompactFooter({ config, locale }: FooterVariantProps) {
  return (
    <FooterShell config={config} locale={locale} className="py-0">
      <div className="flex flex-col gap-4 border-b border-white/10 pb-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h3 className="text-base font-semibold">{config.companyName || 'Restaurant'}</h3>
          <p className="text-sm opacity-90">{config.phone || config.email || ''}</p>
        </div>
      </div>
      <div className="mt-4 flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
        <p className="text-xs opacity-90">{config.copyrightText}</p>
        <FooterLegal config={config} locale={locale} />
      </div>
    </FooterShell>
  )
}
