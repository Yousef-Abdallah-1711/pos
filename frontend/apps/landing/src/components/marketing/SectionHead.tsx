interface SectionHeadProps {
  eyebrow: string
  title: string
  body?: string
}

export function SectionHead({ eyebrow, title, body }: SectionHeadProps) {
  return <header className="sv-section-head">
    <p className="sv-eyebrow">{eyebrow}</p><h2>{title}</h2>{body && <p className="sv-lead">{body}</p>}
  </header>
}
