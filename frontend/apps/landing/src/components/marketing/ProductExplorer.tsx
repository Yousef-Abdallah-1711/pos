'use client'

import { useRef, useState, type KeyboardEvent } from 'react'
import Link from 'next/link'
import type { Feature } from '@/content/features'
import { getSiteCopy, type Locale } from '@/content/site'
import { ProductCapture } from './ProductCapture'

export function ProductExplorer({ locale, features }: { locale: Locale; features: Feature[] }) {
  const [selected, setSelected] = useState(0)
  const tabs = useRef<(HTMLButtonElement | null)[]>([])
  const c = getSiteCopy(locale)
  const navigate = (event: KeyboardEvent<HTMLButtonElement>, index: number) => {
    const forward = locale === 'ar' ? 'ArrowLeft' : 'ArrowRight'
    const back = locale === 'ar' ? 'ArrowRight' : 'ArrowLeft'
    const next = event.key === 'Home' ? 0 : event.key === 'End' ? features.length - 1
      : event.key === forward ? (index + 1) % features.length
      : event.key === back ? (index - 1 + features.length) % features.length : null
    if (next === null) return
    event.preventDefault()
    setSelected(next)
    tabs.current[next]?.focus()
  }
  return <div className="sv-explorer">
    <div role="tablist" aria-label={c.explore} className="sv-explorer-tabs">
      {features.map((feature, i) => <button key={feature.slug} ref={element => { tabs.current[i] = element }} type="button" role="tab" id={`explore-tab-${feature.slug}`} aria-controls={`explore-panel-${feature.slug}`} aria-selected={selected === i} tabIndex={selected === i ? 0 : -1} onClick={() => setSelected(i)} onKeyDown={event => navigate(event, i)}>{feature.title}</button>)}
    </div>
    {features.map((feature, i) => <div key={feature.slug} role="tabpanel" id={`explore-panel-${feature.slug}`} aria-labelledby={`explore-tab-${feature.slug}`} hidden={selected !== i} tabIndex={0}>
      <div className="sv-explorer-panel"><ProductCapture locale={locale} image={feature.image} alt={feature.title} />
        <div><p className="sv-eyebrow">{feature.eyebrow}</p><h3>{feature.headline}</h3><p>{feature.description}</p><Link className="sv-text-link" href={`/${locale}/features/${feature.slug}`}>{c.details}<span aria-hidden="true">↗</span></Link></div></div>
    </div>)}
  </div>
}
