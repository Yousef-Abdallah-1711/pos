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
        <FooterLinkList title={locale === 'ar' ? 'روابط سريعة' : 'Quick Links'} links={(config.legalLinks || []).slice(0, 4)} locale={locale} />
        <div>
          <h3 className="mb-3 text-sm font-semibold">{locale === 'ar' ? 'ساعات العمل' : 'Hours'}</h3>
          <p className="text-sm opacity-90">{locale === 'ar' ? 'يرجى التواصل مع المطعم لمعرفة ساعات العمل' : 'Contact the restaurant for opening hours'}</p>
        </div>
        <div>
          <h3 className="mb-3 text-sm font-semibold">{locale === 'ar' ? 'تواصل معنا' : 'Contact'}</h3>
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
          <BackToTop locale={locale} />
        </div>
      </div>
    </FooterShell>
  )
}
