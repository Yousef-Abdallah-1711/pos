type MotionModules = {
  gsap: typeof import('gsap').default
  ScrollTrigger: typeof import('gsap/ScrollTrigger').ScrollTrigger
}

const desktopMotionQuery = '(min-width: 901px) and (prefers-reduced-motion: no-preference)'

export function startDesktopMotion(setup: (modules: MotionModules) => void | (() => void)) {
  if (typeof window.matchMedia !== 'function') return () => undefined
  const media = window.matchMedia(desktopMotionQuery)
  let generation = 0
  let cleanup: (() => void) | undefined

  const update = async () => {
    const current = ++generation
    cleanup?.()
    cleanup = undefined
    if (!media.matches) return

    const [{ default: gsap }, { ScrollTrigger }] = await Promise.all([
      import('gsap'),
      import('gsap/ScrollTrigger'),
    ])
    if (current !== generation || !media.matches) return
    gsap.registerPlugin(ScrollTrigger)
    cleanup = setup({ gsap, ScrollTrigger }) || undefined
  }

  media.addEventListener('change', update)
  void update()
  return () => {
    generation++
    media.removeEventListener('change', update)
    cleanup?.()
  }
}
