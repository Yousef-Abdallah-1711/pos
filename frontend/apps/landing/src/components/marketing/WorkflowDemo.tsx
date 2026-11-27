'use client'

import { useRef, useState } from 'react'
import { useGSAP } from '@gsap/react'
import gsap from 'gsap'
import type { Locale } from '@/content/site'
import { ProductCapture } from './ProductCapture'

type Kind = 'online' | 'kitchen'

const copy = {
  en: {
    sample: 'Illustrative workflow · sample data',
    online: {
      title: 'From storefront request to staff confirmation',
      alt: 'Actual SERVIO restaurant checkout screen',
      stages: [
        { label: 'Browse', status: 'Choosing from the menu', detail: 'The guest explores the restaurant menu before submitting anything.' },
        { label: 'Review', status: 'Cash on delivery', detail: 'The guest checks items and delivery details. No card payment is implied.' },
        { label: 'Submit', status: 'Pending staff review', detail: 'The order request reaches staff; it has not entered kitchen preparation.' },
        { label: 'Confirm', status: 'Confirmed by staff', detail: 'A team member confirms the online order before kitchen preparation.' },
      ],
    },
    kitchen: {
      title: 'A ticket moves through preparation',
      alt: 'Actual SERVIO kitchen display screen',
      stages: [
        { label: 'New', status: 'New ticket', detail: 'A confirmed order appears in the kitchen queue.' },
        { label: 'In progress', status: 'Preparing', detail: 'The kitchen team moves the ticket into progress.' },
        { label: 'Ready', status: 'Ready', detail: 'The team marks preparation ready. Fulfillment remains separate.' },
        { label: 'Completed', status: 'Completed ticket', detail: 'The kitchen ticket is complete; the order follows its own fulfillment state.' },
      ],
    },
  },
  ar: {
    sample: 'عرض توضيحي · بيانات تجريبية',
    online: {
      title: 'من طلب الموقع إلى تأكيد الفريق',
      alt: 'شاشة فعلية لإتمام طلب المطعم في سيرفيو',
      stages: [
        { label: 'تصفح', status: 'اختيار من المنيو', detail: 'الضيف يستكشف منيو المطعم قبل إرسال أي طلب.' },
        { label: 'مراجعة', status: 'الدفع عند الاستلام', detail: 'الضيف يراجع الأصناف والتوصيل؛ لا يوجد دفع إلكتروني بالبطاقة.' },
        { label: 'إرسال', status: 'بانتظار مراجعة الفريق', detail: 'الطلب يصل للفريق ولم يدخل مرحلة تجهيز المطبخ.' },
        { label: 'تأكيد', status: 'تأكيد من الفريق', detail: 'أحد أفراد الفريق يؤكد الطلب الأونلاين قبل بدء التحضير.' },
      ],
    },
    kitchen: {
      title: 'حركة التذكرة داخل المطبخ',
      alt: 'شاشة فعلية لمطبخ سيرفيو',
      stages: [
        { label: 'جديدة', status: 'تذكرة جديدة', detail: 'الطلب المؤكد يظهر في قائمة المطبخ.' },
        { label: 'جارٍ التحضير', status: 'جارٍ التحضير', detail: 'الفريق ينقل التذكرة إلى مرحلة التحضير.' },
        { label: 'جاهزة', status: 'جاهزة', detail: 'الفريق يحدد أن التحضير انتهى؛ التسليم خطوة منفصلة.' },
        { label: 'مكتملة', status: 'تذكرة مكتملة', detail: 'تذكرة المطبخ اكتملت؛ للطلب حالة تسليم مستقلة.' },
      ],
    },
  },
} as const

export function WorkflowDemo({ locale, kind }: { locale: Locale; kind: Kind }) {
  const root = useRef<HTMLDivElement>(null)
  const [active, setActive] = useState(0)
  const c = copy[locale][kind]
  const current = c.stages[active]

  useGSAP(() => {
    if (!root.current || matchMedia('(prefers-reduced-motion: reduce)').matches) return
    const state = root.current.querySelector('.cv-demo-state')
    if (!state) return
    gsap.fromTo(state, { y: 8, opacity: .65 }, { y: 0, opacity: 1, duration: .32, ease: 'power2.out' })
  }, { scope: root, dependencies: [active, locale, kind] })

  return <div className={`cv-workflow-demo cv-workflow-demo--${kind}`} ref={root}>
    <div className="cv-demo-screen"><ProductCapture image={kind === 'online' ? 'checkout' : 'kitchen'} alt={c.alt} locale={locale} /></div>
    <div className="cv-demo-controls">
      <p className="cv-demo-label">{copy[locale].sample}</p>
      <h3>{c.title}</h3>
      <div className="cv-demo-tabs" role="group" aria-label={c.title}>
        {c.stages.map((item, index) => <button key={item.label} type="button" aria-pressed={active === index} onClick={() => setActive(index)}>{item.label}</button>)}
      </div>
      <div className="cv-demo-state" aria-live="polite"><strong>{current.status}</strong><p>{current.detail}</p></div>
      <div className="cv-demo-progress" aria-hidden="true"><span style={{ width: `${((active + 1) / c.stages.length) * 100}%` }} /></div>
    </div>
  </div>
}
