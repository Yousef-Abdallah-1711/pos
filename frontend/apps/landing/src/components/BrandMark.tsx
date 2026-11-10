import Link from 'next/link'

interface BrandMarkProps {
  href: string
  label: string
  inverse?: boolean
}

export function BrandMark({ href, label, inverse = false }: BrandMarkProps) {
  return (
    <Link href={href} className={`poos-brand${inverse ? ' poos-brand--inverse' : ''}`} aria-label={label}>
      <svg aria-hidden="true" viewBox="0 0 48 48" className="poos-brand__symbol" fill="none">
        <path d="M11 39V9h14a10 10 0 0 1 0 20H11m14 0h12" stroke="currentColor" strokeWidth="7" strokeLinecap="round" strokeLinejoin="round" />
        <circle cx="37" cy="29" r="3.5" fill="var(--poos-brand-node, #8B5CF6)" />
      </svg>
      <span>POOS</span>
    </Link>
  )
}
