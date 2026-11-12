'use client'

import { HeaderVariantShell, type HeaderVariantProps } from './shared'

export function MinimalHeader(props: HeaderVariantProps) {
  return (
    <HeaderVariantShell
      {...props}
      className="sticky top-0 z-[60] border-b border-[hsl(var(--theme-border))] bg-transparent text-[hsl(var(--theme-header-text))]"
    />
  )
}
