import Link from 'next/link'

interface BrandMarkProps { href: string; label: string; inverse?: boolean }

export function BrandMark({ href, label, inverse = false }: BrandMarkProps) {
  return <Link href={href} className={`sv-brand${inverse ? ' sv-brand-inverse' : ''}`} aria-label={label}>
    <span className="sv-brand-mark" aria-hidden="true"><i /><i /><i /></span><span>SERVIO</span>
  </Link>
}
