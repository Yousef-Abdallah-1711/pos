'use client'

import { HeaderVariantShell, type HeaderVariantProps } from './shared'

export function SolidHeader(props: HeaderVariantProps) {
  return (
    <HeaderVariantShell
      {...props}
      className="sticky top-0 z-[60] border-theme-border/40 bg-[hsl(var(--theme-header-bg))] text-[hsl(var(--theme-header-text))]"
    />
  )
}
