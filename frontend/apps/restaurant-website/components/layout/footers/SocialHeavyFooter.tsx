'use client'

import { FooterShell, FooterSocial, FooterLegal, type FooterVariantProps } from './shared'

export function SocialHeavyFooter({ config, locale }: FooterVariantProps) {
  return (
    <FooterShell config={config} locale={locale}>
      <div className="flex flex-col items-center gap-4 text-center">
        <h3 className="text-lg font-semibold">{config.companyName || 'Restaurant'}</h3>
        <p className="max-w-xl text-sm opacity-90">{config.companyDescription || ''}</p>
        <div className="[&_a]:p-3"><FooterSocial config={config} /></div>
      </div>
      <div className="mt-6 flex flex-col gap-2 border-t border-white/10 pt-4 md:flex-row md:items-center md:justify-between">
        <p className="text-xs opacity-90">{config.copyrightText}</p>
        <FooterLegal config={config} locale={locale} />
      </div>
    </FooterShell>
  )
}
