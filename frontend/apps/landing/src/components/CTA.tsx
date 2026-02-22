'use client'

import { useLanguage } from '@/contexts/LanguageContext'
import { Button } from '@/components/ui/button'
import { ArrowRight, CheckCircle2 } from 'lucide-react'
import Link from 'next/link'

export function CTA() {
  const { t, language, dir } = useLanguage()
  const locale = language

  return (
    <section className="landing-cta py-24 px-4 sm:px-6 lg:px-8 relative overflow-hidden" dir={dir}>
      {/* Background */}
      <div className="absolute inset-0 bg-gradient-to-br from-orange-600 via-red-600 to-red-700" />
      <div className="absolute inset-0 opacity-10">
        <div className="cta-grid-pattern w-full h-full" />
      </div>
      <div className="absolute top-0 start-0 w-64 h-64 bg-white/10 rounded-full blur-3xl" />
      <div className="absolute bottom-0 end-0 w-64 h-64 bg-white/10 rounded-full blur-3xl" />

      <div className="relative max-w-4xl mx-auto text-center">
        <h2 className="text-3xl sm:text-4xl lg:text-5xl font-black text-white mb-4 tracking-tight">
          {t.cta.title}
        </h2>
        <p className="text-xl text-white/80 mb-10 max-w-2xl mx-auto">
          {t.cta.subtitle}
        </p>

        <Button
          size="lg"
          asChild
          className="group bg-white text-red-600 hover:bg-gray-50 shadow-2xl border-0 text-base font-bold px-10 h-14 transition-all hover:-translate-y-0.5"
        >
          <Link href={`/${locale}/auth/register`}>
            {t.cta.button}
            <ArrowRight className="ms-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
          </Link>
        </Button>

        <div className="mt-6 flex items-center justify-center gap-2 text-white/70 text-sm">
          <CheckCircle2 className="h-4 w-4 text-white/90 flex-shrink-0" />
          <span>{t.cta.noCard}</span>
        </div>
      </div>
      <style jsx>{`
        .cta-grid-pattern {
          background-image:
            linear-gradient(white 1px, transparent 1px),
            linear-gradient(to right, white 1px, transparent 1px);
          background-size: 48px 48px;
        }
      `}</style>
    </section>
  )
}
