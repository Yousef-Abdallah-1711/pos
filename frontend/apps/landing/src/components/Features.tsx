'use client'

import { useLanguage } from '@/contexts/LanguageContext'
import {
  CreditCard,
  ChefHat,
  Globe,
  BarChart3,
  Store,
  Smartphone,
} from 'lucide-react'

const features = [
  {
    key: 'pos' as const,
    icon: CreditCard,
    gradient: 'from-blue-500 to-blue-600',
    bgLight: 'bg-blue-50 dark:bg-blue-900/20',
    border: 'group-hover:border-blue-200 dark:group-hover:border-blue-800',
  },
  {
    key: 'kitchen' as const,
    icon: ChefHat,
    gradient: 'from-orange-500 to-red-500',
    bgLight: 'bg-orange-50 dark:bg-orange-900/20',
    border: 'group-hover:border-orange-200 dark:group-hover:border-orange-800',
  },
  {
    key: 'online' as const,
    icon: Globe,
    gradient: 'from-green-500 to-emerald-600',
    bgLight: 'bg-green-50 dark:bg-green-900/20',
    border: 'group-hover:border-green-200 dark:group-hover:border-green-800',
  },
  {
    key: 'analytics' as const,
    icon: BarChart3,
    gradient: 'from-purple-500 to-violet-600',
    bgLight: 'bg-purple-50 dark:bg-purple-900/20',
    border: 'group-hover:border-purple-200 dark:group-hover:border-purple-800',
  },
  {
    key: 'multi' as const,
    icon: Store,
    gradient: 'from-pink-500 to-rose-600',
    bgLight: 'bg-pink-50 dark:bg-pink-900/20',
    border: 'group-hover:border-pink-200 dark:group-hover:border-pink-800',
  },
  {
    key: 'mobile' as const,
    icon: Smartphone,
    gradient: 'from-indigo-500 to-blue-600',
    bgLight: 'bg-indigo-50 dark:bg-indigo-900/20',
    border: 'group-hover:border-indigo-200 dark:group-hover:border-indigo-800',
  },
]

export function Features() {
  const { t, dir } = useLanguage()

  return (
    <section id="features" className="py-24 px-4 sm:px-6 lg:px-8 bg-white dark:bg-gray-900" dir={dir}>
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="text-center mb-16">
          <span className="inline-block text-sm font-semibold text-orange-600 dark:text-orange-400 uppercase tracking-wider mb-3">
            Features
          </span>
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-black text-gray-900 dark:text-white mb-4 tracking-tight">
            {t.features.title}
          </h2>
          <p className="text-lg sm:text-xl text-gray-500 dark:text-gray-400 max-w-3xl mx-auto">
            {t.features.subtitle}
          </p>
        </div>

        {/* Features Grid */}
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map(({ key, icon: Icon, gradient, bgLight, border }) => (
            <div
              key={key}
              className={`group relative p-6 rounded-2xl border border-gray-200 dark:border-gray-800 hover:shadow-xl transition-all duration-300 bg-white dark:bg-gray-800/50 ${border} cursor-default`}
            >
              {/* Icon */}
              <div className={`w-14 h-14 rounded-xl ${bgLight} flex items-center justify-center mb-5 transition-transform duration-300 group-hover:scale-110`}>
                <div className={`w-9 h-9 rounded-lg bg-gradient-to-br ${gradient} flex items-center justify-center shadow-sm`}>
                  <Icon className="h-5 w-5 text-white" />
                </div>
              </div>

              <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">
                {t.features[key].title}
              </h3>
              <p className="text-gray-500 dark:text-gray-400 text-sm leading-relaxed">
                {t.features[key].description}
              </p>

              {/* Subtle gradient overlay on hover */}
              <div className={`absolute inset-0 rounded-2xl bg-gradient-to-br ${gradient} opacity-0 group-hover:opacity-[0.03] transition-opacity duration-300`} />
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
