'use client'

import { useEffect, useRef, useState } from 'react'
import type { Locale } from '@/content/site'
import { ProductCapture } from './ProductCapture'
import { startDesktopMotion } from './desktopMotion'

const copy = {
  en: {
    eyebrow: 'A clearer view of the work', title: 'See the day beyond one ticket.',
    body: 'Move between the dashboard, tracked inventory and team records in the same restaurant workspace.',
    sample: 'Sample workspace data, not customer results',
    orders: 'Sample orders', revenue: 'Sample paid revenue', currency: 'EGP',
    views: [
      { label: 'Overview', image: 'dashboard', alt: 'Actual SERVIO restaurant dashboard', detail: 'The dashboard brings operational areas into one view.' },
      { label: 'Inventory', image: 'inventory', alt: 'Actual SERVIO inventory screen', detail: 'Review tracked quantities and stock work.' },
      { label: 'Team', image: 'staff', alt: 'Actual SERVIO team records screen', detail: 'Find staff and people-related work in its own area.' },
    ],
  },
  ar: {
    eyebrow: 'رؤية أوضح للشغل', title: 'شوف اليوم كله، مش تذكرة واحدة.',
    body: 'انتقل بين لوحة التحكم والمخزون المسجل وبيانات الفريق من مساحة عمل المطعم.',
    sample: 'بيانات مساحة عمل تجريبية، وليست نتائج عملاء',
    orders: 'طلبات تجريبية', revenue: 'إيراد مدفوع تجريبي', currency: 'جنيه',
    views: [
      { label: 'نظرة عامة', image: 'dashboard', alt: 'لوحة تحكم فعلية من سيرفيو', detail: 'لوحة التحكم تجمع أجزاء الشغل في عرض واحد.' },
      { label: 'المخزون', image: 'inventory', alt: 'شاشة مخزون فعلية من سيرفيو', detail: 'راجع الكميات المسجلة وحركة المخزون.' },
      { label: 'الفريق', image: 'staff', alt: 'شاشة بيانات فريق فعلية من سيرفيو', detail: 'بيانات العاملين وشغل الفريق في مكان مخصص.' },
    ],
  },
} as const

export function ManagementExperience({ locale }: { locale: Locale }) {
  const root = useRef<HTMLElement>(null)
  const orders = useRef<HTMLElement>(null)
  const revenue = useRef<HTMLElement>(null)
  const [active, setActive] = useState(0)
  const c = copy[locale]
  const view = c.views[active]

  useEffect(() => startDesktopMotion(({ gsap }) => {
    if (!root.current || !orders.current || !revenue.current) return
    const context = gsap.context(() => {
      gsap.from([orders.current, revenue.current], {
        y: 14, duration: .8, stagger: .12, ease: 'power2.out',
        scrollTrigger: { trigger: root.current, start: 'top 70%', once: true },
      })
    }, root.current)
    return () => context.revert()
  }), [])

  return <section className="cv-section cv-management" aria-labelledby="cv-management-title" ref={root}>
    <div className="cv-wrap cv-management-layout">
      <div className="cv-management-copy">
        <p className="cv-eyebrow">{c.eyebrow}</p><h2 id="cv-management-title">{c.title}</h2><p>{c.body}</p>
        <div className="cv-management-metrics" aria-label={c.sample}>
          <div><strong ref={orders}>1</strong><span>{c.orders}</span></div>
          <div><strong><span ref={revenue}>55</span> <small>{c.currency}</small></strong><span>{c.revenue}</span></div>
        </div>
        <small>{c.sample}</small>
      </div>
      <div className="cv-management-product">
        <div className="cv-management-tabs" role="group" aria-label={c.title}>
          {c.views.map((item, index) => <button key={item.label} type="button" aria-pressed={active === index} onClick={() => setActive(index)}>{item.label}</button>)}
        </div>
        <ProductCapture image={view.image} alt={view.alt} locale={locale} />
        <p aria-live="polite">{view.detail}</p>
      </div>
    </div>
  </section>
}
