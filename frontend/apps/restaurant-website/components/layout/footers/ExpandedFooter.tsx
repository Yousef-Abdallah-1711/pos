'use client'

import { FooterShell, FooterLinkList, FooterSocial, FooterLegal, BackToTop, type FooterVariantProps } from './shared'

export function ExpandedFooter({ config, locale }: FooterVariantProps) {
  const sections = config.footerSections || []
  return (
    <FooterShell config={config} locale={locale}>
      <div className="grid gap-8 md:grid-cols-4">
        <div>
          <h3 className="text-lg font-semibold">{config.companyName || 'Restaurant'}</h3>
          <p className="mt-3 text-sm opacity-90">{config.companyDescription || ''}</p>
          <div className="mt-4"><FooterSocial config={config} /></div>
        </div>
        <div>
          <h3 className="mb-3 text-sm font-semibold">Quick Links</h3>
          <ul className="space-y-2 text-sm opacity-90">
            {(config.legalLinks || []).slice(0, 4).map((link) => (
              <li key={`${link.href}-${link.label}`}>{link.label}</li>
            ))}
          </ul>
        </div>
        <div>
          <h3 className="mb-3 text-sm font-semibold">Hours</h3>
          <p className="text-sm opacity-90">Daily 10:00 AM - 12:00 AM</p>
        </div>
        <div>
          <h3 className="mb-3 text-sm font-semibold">Contact</h3>
          <p className="text-sm opacity-90">{config.phone || ''}</p>
          <p className="text-sm opacity-90">{config.email || ''}</p>
          <p className="text-sm opacity-90">{config.address || ''}</p>
        </div>
      </div>

      {sections.length > 0 ? (
        <div className="mt-8 grid gap-6 border-t border-white/10 pt-6 md:grid-cols-3">
          {sections.slice(0, 3).map((section) => (
            <FooterLinkList key={section.id} title={section.title} links={section.links} locale={locale} />
          ))}
        </div>
      ) : null}

      <div className="mt-8 flex flex-col gap-3 border-t border-white/10 pt-4 md:flex-row md:items-center md:justify-between">
        <p className="text-xs opacity-90">{config.copyrightText}</p>
        <div className="flex items-center gap-4">
          <FooterLegal config={config} locale={locale} />
          <BackToTop />
        </div>
      </div>
    </FooterShell>
  )
}
