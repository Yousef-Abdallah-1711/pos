import Link from 'next/link'

interface BrandMarkProps {
  href: string
  label: string
}

export function BrandMark({ href, label }: BrandMarkProps) {
  return (
    <Link href={href} className="poos-brand" aria-label={label}>
      <svg aria-hidden="true" viewBox="0 0 40 40" className="poos-brand__symbol">
        <rect x="4" y="4" width="32" height="32" rx="9" fill="currentColor" />
        <path d="M12 13h16M12 20h10M12 27h16" fill="none" stroke="white" strokeWidth="2.4" strokeLinecap="round" />
      </svg>
      <span>POOS</span>
    </Link>
  )
}
