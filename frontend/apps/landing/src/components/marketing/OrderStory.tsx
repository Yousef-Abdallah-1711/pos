'use client'

import { useEffect, useRef, useState } from 'react'
import type { Locale } from '@/content/site'
import { getConversionCopy } from '@/content/conversionCopy'
import { getOrderStoryCopy } from '@/content/orderStoryCopy'
import { ProductCapture } from './ProductCapture'
import { startDesktopMotion } from './desktopMotion'

export function OrderStory({ locale }: { locale: Locale }) {
  const root = useRef<HTMLElement>(null)
  const stage = useRef<HTMLDivElement>(null)
  const [active, setActive] = useState(0)
  const c = getConversionCopy(locale)
  const story = getOrderStoryCopy(locale)
  const current = story.stages[active]

  useEffect(() => startDesktopMotion(({ ScrollTrigger }) => {
    if (!root.current || !stage.current) return
    const trigger = ScrollTrigger.create({
      trigger: root.current,
      start: 'top top+=96',
      end: 'bottom bottom',
      invalidateOnRefresh: true,
      onUpdate: self => setActive(Math.min(story.stages.length - 1, Math.floor(self.progress * story.stages.length))),
    })
    return () => trigger.kill()
  }), [locale, story.stages.length])

  const selectStage = (index: number) => {
    setActive(index)
    root.current?.querySelectorAll('.cv-story-step')[index]?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }

  return <section ref={root} id="how-orders-move" className="cv-section cv-flow cv-story" aria-labelledby="cv-flow-title">
    <div className="cv-wrap">
      <div className="cv-section-heading">
        <div><p className="cv-eyebrow">{c.flowEyebrow}</p><h2 id="cv-flow-title">{c.flowTitle}</h2></div>
        <p>{c.flowBody}</p>
      </div>
      <div className="cv-story-layout">
        <div className="cv-story-narrative">
          {story.stages.map((item, index) => <article className="cv-story-step" key={item.title} data-active={index === active}>
            <button type="button" onClick={() => selectStage(index)} aria-current={index === active ? 'step' : undefined}>
              <span className="cv-index">{String(index + 1).padStart(2, '0')}</span>
              <span><strong>{item.title}</strong><span>{item.body}</span></span>
            </button>
          </article>)}
        </div>
        <div className="cv-story-stage-slot"><div className="cv-story-stage" ref={stage} aria-live="polite">
          <div className="cv-story-stage-head"><span>{story.sample}</span><span>{String(active + 1).padStart(2, '0')} / {String(story.stages.length).padStart(2, '0')}</span></div>
          <div className="cv-story-stage-screen"><ProductCapture image={current.image} alt={current.title} locale={locale} /></div>
          <div className="cv-story-status">
            <span><small>{story.orderStatus}</small><strong>{current.order}</strong></span>
            <span><small>{story.kitchenStatus}</small><strong>{current.kitchen}</strong></span>
          </div>
          <p className="cv-story-annotation"><strong>{story.annotation}:</strong> {current.note}</p>
          <p className="cv-story-screen-note">{story.screenNote}</p>
          <div className="cv-story-progress" role="progressbar" aria-label={c.flowTitle} aria-valuemin={1} aria-valuemax={story.stages.length} aria-valuenow={active + 1}>
            <span style={{ width: `${((active + 1) / story.stages.length) * 100}%` }} />
          </div>
        </div></div>
      </div>
    </div>
  </section>
}
