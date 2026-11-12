'use client'

import { FooterShell, FooterLegal, type FooterVariantProps } from './shared'

export function MinimalFooter({ config, locale }: FooterVariantProps) {
  return (
    <FooterShell config={config} locale={locale} className="py-0">
      <div className="flex flex-col gap-2 py-4 text-center md:flex-row md:items-center md:justify-center md:gap-4">
        <p className="text-xs opacity-90">{config.copyrightText}</p>
        <FooterLegal config={config} locale={locale} />
      </div>
    </FooterShell>
  )
}
