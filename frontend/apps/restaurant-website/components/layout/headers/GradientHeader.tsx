'use client'

import { HeaderVariantShell, type HeaderVariantProps } from './shared'

export function GradientHeader(props: HeaderVariantProps) {
  return (
    <HeaderVariantShell
      {...props}
      className="sticky top-0 z-[60] border-transparent bg-gradient-to-r from-[hsl(var(--theme-primary))] to-[hsl(var(--theme-accent))] text-[hsl(var(--theme-header-text))]"
    />
  )
}
