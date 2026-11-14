import Link from 'next/link'
import { getSiteCopy, type Locale } from '@/content/site'
import { ProductCapture } from './ProductCapture'
import { SectionHead } from './SectionHead'

export function ManagementSections({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)
  return <>
    <section id="customization" className="sv-section sv-customization"><div className="sv-shell">
      <div className="sv-editorial-row"><SectionHead eyebrow={c.themeEyebrow} title={c.themeTitle} body={c.themeBody} />
        <Link className="sv-text-link" href={`/${locale}/features/website-builder`}>{c.themeLink}<span aria-hidden="true">↗</span></Link></div>
      <ProductCapture locale={locale} image="theme" alt={locale === 'ar' ? 'محرر القالب مع اللون البنفسجي المحفوظ ومعاينة تجريبية' : 'Theme editor with a saved violet color and sample preview'} className="sv-wide-capture" />
      <p className="sv-note">{c.themeNote}</p>
    </div></section>
    <section id="inventory" className="sv-section sv-stock"><div className="sv-shell sv-stock-grid">
      <div><SectionHead eyebrow={c.inventoryEyebrow} title={c.inventoryTitle} body={c.inventoryBody} />
        <Link className="sv-text-link" href={`/${locale}/features/inventory`}>{c.inventoryLink}<span aria-hidden="true">↗</span></Link></div>
      <div><ol className="sv-stock-ledger">{[20, 19, 20].map((quantity, i) => <li key={i}><span>{c.stockStates[i]}</span><strong>{quantity}</strong><span>{locale === 'ar' ? 'وحدة' : 'units'}</span></li>)}</ol><p className="sv-note">{c.stockNote}</p></div>
    </div></section>
    <section id="team" className="sv-section sv-team"><div className="sv-shell">
      <SectionHead eyebrow={c.teamEyebrow} title={c.teamTitle} body={c.teamBody} />
      <div className="sv-team-layout"><ProductCapture locale={locale} image="staff" alt={c.teamTitle} />
        <div className="sv-team-notes">{c.teamItems.map(([title, body]) => <article key={title}><h3>{title}</h3><p>{body}</p></article>)}
          <Link className="sv-text-link" href={`/${locale}/features/hr`}>{c.teamLink}<span aria-hidden="true">↗</span></Link></div></div>
    </div></section>
    <section id="visibility" className="sv-section sv-visibility"><div className="sv-shell sv-visibility-grid">
      <div><SectionHead eyebrow={c.visibilityEyebrow} title={c.visibilityTitle} body={c.visibilityBody} />
        <p className="sv-note">{c.visibilityNote}</p><Link className="sv-text-link" href={`/${locale}/features/analytics`}>{c.visibilityLink}<span aria-hidden="true">↗</span></Link></div>
      <ProductCapture locale={locale} image="dashboard" alt={locale === 'ar' ? 'لوحة المطعم التجريبي لعرض الطلبات وحالة مساحة العمل' : 'Synthetic restaurant overview with orders and workspace status'} />
    </div></section>
  </>
}
