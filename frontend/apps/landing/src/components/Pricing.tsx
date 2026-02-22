'use client'

import { useLanguage } from '@/contexts/LanguageContext'
import { Button } from '@/components/ui/button'
import { Check, Sparkles, Clock } from 'lucide-react'
import Link from 'next/link'

export function Pricing() {
  const { t, language, dir } = useLanguage()
  const locale = language

  return (
    <section id="pricing" className="py-24 px-4 sm:px-6 lg:px-8 bg-white dark:bg-gray-900" dir={dir}>
      <div className="max-w-5xl mx-auto">
        {/* Header */}
        <div className="text-center mb-16">
          <span className="inline-block text-sm font-semibold text-orange-600 dark:text-orange-400 uppercase tracking-wider mb-3">
            Pricing
          </span>
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-black text-gray-900 dark:text-white mb-4 tracking-tight">
            {t.pricing.title}
          </h2>
          <p className="text-lg text-gray-500 dark:text-gray-400 max-w-2xl mx-auto">
            {t.pricing.subtitle}
          </p>
        </div>

        {/* Single Pricing Card */}
        <div className="max-w-2xl mx-auto">
          {/* Trial Banner */}
          <div className="bg-gradient-to-r from-orange-500 to-red-600 rounded-t-3xl px-8 py-5 text-center">
            <div className="flex items-center justify-center gap-2 text-white">
              <Clock className="h-5 w-5" />
              <span className="font-bold text-lg">{t.pricing.trialBadge}</span>
            </div>
            <p className="text-orange-100 text-sm mt-1">{t.pricing.trialNote}</p>
          </div>

          {/* Main Card */}
          <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 border-t-0 rounded-b-3xl p-8 shadow-2xl">
            {/* Price */}
            <div className="text-center mb-8 pb-8 border-b border-gray-100 dark:border-gray-700">
              <div className="flex items-center justify-center gap-1 mb-2">
                <Sparkles className="h-5 w-5 text-orange-500" />
                <span className="text-sm font-semibold text-gray-500 dark:text-gray-400">
                  {t.pricing.afterTrial}
                </span>
              </div>
              <div className="flex items-baseline justify-center gap-2">
                <span className="text-6xl font-black text-gray-900 dark:text-white">
                  {t.pricing.price}
                </span>
                <div className="text-start">
                  <div className="text-xl font-bold text-orange-500">{t.pricing.currency}</div>
                  <div className="text-gray-500 dark:text-gray-400 text-sm">{t.pricing.period}</div>
                </div>
              </div>
              <div className="mt-3 inline-flex items-center gap-2 bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-400 px-4 py-1.5 rounded-full text-sm font-medium">
                <Check className="h-4 w-4" />
                {t.pricing.allFeatures}
              </div>
            </div>

            {/* Features grid */}
            <div className="grid sm:grid-cols-2 gap-3 mb-8">
              {t.pricing.features.map((feature: string, index: number) => (
                <div key={index} className="flex items-start gap-3">
                  <div className="w-5 h-5 rounded-full bg-gradient-to-br from-orange-500 to-red-600 flex items-center justify-center flex-shrink-0 mt-0.5">
                    <Check className="h-3 w-3 text-white" />
                  </div>
                  <span className="text-gray-700 dark:text-gray-300 text-sm">{feature}</span>
                </div>
              ))}
            </div>

            {/* CTA */}
            <Button
              size="lg"
              asChild
              className="w-full h-14 text-base font-bold bg-gradient-to-r from-orange-500 to-red-600 hover:from-orange-600 hover:to-red-700 text-white shadow-lg border-0 transition-all hover:shadow-orange-200 dark:hover:shadow-orange-900/50"
            >
              <Link href={`/${locale}/auth/register`}>{t.pricing.cta}</Link>
            </Button>

            {/* Notes */}
            <div className="mt-4 flex flex-wrap items-center justify-center gap-4 text-sm text-gray-400 dark:text-gray-500">
              <span className="flex items-center gap-1">
                <Check className="h-3.5 w-3.5 text-green-500" />
                {t.pricing.noCard}
              </span>
              <span className="flex items-center gap-1">
                <Check className="h-3.5 w-3.5 text-green-500" />
                {t.pricing.cancelAnytime}
              </span>
            </div>
          </div>

          {/* After trial note */}
          <p className="text-center text-sm text-gray-400 dark:text-gray-500 mt-4">
            {t.pricing.afterTrialNote}
          </p>
        </div>
      </div>
    </section>
  )
}
