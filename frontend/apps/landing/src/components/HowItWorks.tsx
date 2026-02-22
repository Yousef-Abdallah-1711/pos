'use client'

import { useLanguage } from '@/contexts/LanguageContext'
import { UserPlus, Settings, ShoppingBag } from 'lucide-react'

const stepIcons = [UserPlus, Settings, ShoppingBag]
const stepGradients = [
  'from-blue-500 to-blue-600',
  'from-orange-500 to-red-500',
  'from-green-500 to-emerald-600',
]

export function HowItWorks() {
  const { t, dir } = useLanguage()

  const steps = [t.howItWorks.step1, t.howItWorks.step2, t.howItWorks.step3]

  return (
    <section className="py-24 px-4 sm:px-6 lg:px-8 bg-gray-50 dark:bg-gray-800/50" dir={dir}>
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="text-center mb-16">
          <span className="inline-block text-sm font-semibold text-orange-600 dark:text-orange-400 uppercase tracking-wider mb-3">
            How It Works
          </span>
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-black text-gray-900 dark:text-white mb-4 tracking-tight">
            {t.howItWorks.title}
          </h2>
          <p className="text-lg text-gray-500 dark:text-gray-400 max-w-2xl mx-auto">
            {t.howItWorks.subtitle}
          </p>
        </div>

        {/* Steps */}
        <div className="grid md:grid-cols-3 gap-8 relative">
          {/* Connector line (desktop) */}
          <div className="hidden md:block absolute top-12 start-1/6 end-1/6 h-0.5 bg-gradient-to-r from-blue-200 via-orange-200 to-green-200 dark:from-blue-800 dark:via-orange-800 dark:to-green-800" />

          {steps.map((step, index) => {
            const Icon = stepIcons[index]
            const gradient = stepGradients[index]
            return (
              <div key={step.number} className="relative flex flex-col items-center text-center">
                {/* Step number + Icon */}
                <div className="relative mb-6">
                  <div className={`w-24 h-24 rounded-2xl bg-gradient-to-br ${gradient} flex items-center justify-center shadow-lg`}>
                    <Icon className="h-10 w-10 text-white" />
                  </div>
                  <div className="absolute -top-3 -end-3 w-8 h-8 rounded-full bg-white dark:bg-gray-900 border-2 border-gray-200 dark:border-gray-700 flex items-center justify-center">
                    <span className="text-xs font-black text-gray-600 dark:text-gray-300">{step.number}</span>
                  </div>
                </div>

                <h3 className="text-xl font-bold text-gray-900 dark:text-white mb-3">
                  {step.title}
                </h3>
                <p className="text-gray-500 dark:text-gray-400 text-sm leading-relaxed max-w-xs mx-auto">
                  {step.description}
                </p>
              </div>
            )
          })}
        </div>
      </div>
    </section>
  )
}
