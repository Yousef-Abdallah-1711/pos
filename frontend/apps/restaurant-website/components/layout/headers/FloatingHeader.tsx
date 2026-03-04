'use client'

import { HeaderVariantShell, type HeaderVariantProps } from './shared'

export function FloatingHeader(props: HeaderVariantProps) {
  return (
    <HeaderVariantShell
      {...props}
      className="fixed inset-x-3 top-3 z-[60] rounded-b-2xl border-theme-border/40 bg-[hsl(var(--theme-header-bg))] text-[hsl(var(--theme-header-text))] shadow-2xl"
      innerClassName="max-w-none"
    />
  )
}
