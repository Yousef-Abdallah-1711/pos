'use client'

import { useEffect, useRef, type ReactNode } from 'react'
import { startDesktopMotion } from './desktopMotion'

export function HeroCaptureMotion({ children }: { children: ReactNode }) {
  const root = useRef<HTMLDivElement>(null)

  useEffect(() => startDesktopMotion(({ gsap, ScrollTrigger }) => {
    if (!root.current) return
    const context = gsap.context(() => {
      const capture = root.current?.querySelector('.cv-hero-capture')
      if (!capture) return
      gsap.fromTo(capture,
        { rotateX: 11, scale: .93, transformOrigin: 'center top' },
        { rotateX: 0, scale: 1, ease: 'none', scrollTrigger: {
          trigger: root.current,
          start: 'top bottom',
          end: 'bottom top+=120',
          scrub: .65,
          invalidateOnRefresh: true,
        } })
    }, root.current)
    let active = true
    const refresh = () => { if (active) ScrollTrigger.refresh() }
    void document.fonts.ready.then(refresh)
    window.addEventListener('load', refresh, { once: true })
    return () => {
      active = false
      window.removeEventListener('load', refresh)
      context.revert()
    }
  }), [])

  return <div className="cv-hero-motion" ref={root}>{children}</div>
}
