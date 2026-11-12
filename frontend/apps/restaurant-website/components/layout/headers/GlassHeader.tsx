'use client'

import { HeaderVariantShell, type HeaderVariantProps } from './shared'

export function GlassHeader(props: HeaderVariantProps) {
  return (
    <HeaderVariantShell
      {...props}
      className="sticky top-0 z-[60] border-white/20 bg-[hsl(var(--theme-header-bg)/0.1)] text-[hsl(var(--theme-header-text))] backdrop-blur-xl"
      innerClassName=""
    />
  )
}
