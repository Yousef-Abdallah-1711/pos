'use client'

import { useLanguage } from '@/contexts/LanguageContext'
import { Button } from '@/components/ui/button'
import { ArrowRight, Play, CheckCircle2, TrendingUp, Zap } from 'lucide-react'
import Link from 'next/link'

export function Hero() {
  const { t, language, dir } = useLanguage()
  const locale = language
  const chartHeights = [40, 65, 45, 80, 55, 90, 70]
  const chartHeightRules = chartHeights
    .map((height, index) => `.chart-bar-${index} { height: ${height}%; }`)
    .join('\n')

  return (
    <section className="landing-hero relative pt-24 pb-16 px-4 sm:px-6 lg:px-8 overflow-hidden" dir={dir}>
      {/* Background */}
      <div className="absolute inset-0 -z-10">
        <div className="absolute inset-0 bg-gradient-to-br from-orange-50 via-white to-red-50 dark:from-gray-950 dark:via-gray-900 dark:to-gray-950" />
        <div className="absolute top-0 start-1/4 w-96 h-96 bg-orange-200/30 dark:bg-orange-900/10 rounded-full blur-3xl" />
        <div className="absolute bottom-0 end-1/4 w-96 h-96 bg-red-200/30 dark:bg-red-900/10 rounded-full blur-3xl" />
        <div className="hero-grid-pattern absolute inset-0 opacity-[0.03] dark:opacity-[0.05]" />
      </div>

      <div className="max-w-7xl mx-auto">
        <div className="grid lg:grid-cols-2 gap-12 lg:gap-16 items-center min-h-[calc(100vh-6rem)]">
          {/* Left Content */}
          <div className="flex flex-col justify-center text-center lg:text-start">
            {/* Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-300 text-sm font-medium mb-6 self-center lg:self-start border border-orange-200 dark:border-orange-800">
              <Zap className="h-4 w-4 fill-current" />
              {t.hero.badge}
            </div>

            {/* Heading */}
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-black text-gray-900 dark:text-white mb-6 leading-[1.1] tracking-tight">
              {t.hero.title}
              <br />
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-orange-500 to-red-600">
                {t.hero.titleHighlight}
              </span>
            </h1>

            {/* Description */}
            <p className="text-lg sm:text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-xl mx-auto lg:mx-0 leading-relaxed">
              {t.hero.description}
            </p>

            {/* CTAs */}
            <div className="flex flex-col sm:flex-row gap-3 justify-center lg:justify-start mb-10">
              <Button
                size="lg"
                asChild
                className="group bg-gradient-to-r from-orange-500 to-red-600 hover:from-orange-600 hover:to-red-700 text-white shadow-lg hover:shadow-orange-200 dark:hover:shadow-orange-900/50 border-0 text-base px-8 h-12 transition-all"
              >
                <Link href={`/${locale}/auth/register`}>
                  {t.hero.cta}
                  <ArrowRight className="ms-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
                </Link>
              </Button>
              <Button
                size="lg"
                variant="outline"
                asChild
                className="group text-base px-8 h-12 border-gray-300 dark:border-gray-700"
              >
                <Link href="#demo">
                  <Play className="me-2 h-5 w-5 fill-current opacity-70" />
                  {t.hero.demo}
                </Link>
              </Button>
            </div>

            {/* Trust indicator */}
            <div className="flex items-center gap-2 justify-center lg:justify-start text-sm text-gray-500 dark:text-gray-400">
              <CheckCircle2 className="h-4 w-4 text-green-500 flex-shrink-0" />
              <span>{t.hero.customers}</span>
            </div>

            {/* Stats */}
            <div className="mt-10 grid grid-cols-3 gap-4 pt-8 border-t border-gray-200 dark:border-gray-800">
              {[
                { value: '500+', label: t.hero.stats.restaurants },
                { value: '10K+', label: t.hero.stats.orders },
                { value: '99.9%', label: t.hero.stats.uptime },
              ].map((stat) => (
                <div key={stat.label} className="text-center lg:text-start">
                  <div className="text-2xl sm:text-3xl font-black text-gray-900 dark:text-white">
                    {stat.value}
                  </div>
                  <div className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                    {stat.label}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Right — Dashboard Mockup */}
          <div className="relative flex items-center justify-center">
            <div className="relative w-full max-w-lg mx-auto">
              {/* Main mockup */}
              <div className="relative rounded-2xl overflow-hidden shadow-2xl border border-white/20 dark:border-gray-700/50 bg-gradient-to-br from-gray-900 to-gray-800 aspect-[4/3]">
                {/* Window top bar */}
                <div className="flex items-center gap-2 px-4 py-3 bg-gray-800/80 border-b border-gray-700/50">
                  <div className="flex gap-1.5">
                    <div className="w-3 h-3 rounded-full bg-red-500/80" />
                    <div className="w-3 h-3 rounded-full bg-yellow-500/80" />
                    <div className="w-3 h-3 rounded-full bg-green-500/80" />
                  </div>
                  <div className="flex-1 text-center">
                    <span className="text-xs text-gray-400 font-mono">Restaurant POS Dashboard</span>
                  </div>
                </div>

                {/* Dashboard content */}
                <div className="p-4 space-y-3">
                  {/* Stats row */}
                  <div className="grid grid-cols-3 gap-2">
                    {[
                      { label: "Revenue", value: '12,450 EGP', color: 'text-green-400' },
                      { label: 'Orders', value: '84', color: 'text-blue-400' },
                      { label: 'Avg. Order', value: '148 EGP', color: 'text-orange-400' },
                    ].map((item) => (
                      <div key={item.label} className="bg-gray-700/50 rounded-xl p-3">
                        <div className={`text-lg font-bold ${item.color}`}>{item.value}</div>
                        <div className="text-xs text-gray-400 mt-0.5">{item.label}</div>
                      </div>
                    ))}
                  </div>

                  {/* Orders list */}
                  <div className="bg-gray-700/30 rounded-xl p-3 space-y-2">
                    <div className="text-xs font-semibold text-gray-300 mb-2">Recent Orders</div>
                    {[
                      { id: '#2341', table: 'Table 5', items: '3 items', status: 'Ready', color: 'bg-green-500/20 text-green-400' },
                      { id: '#2340', table: 'Table 12', items: '2 items', status: 'Cooking', color: 'bg-orange-500/20 text-orange-400' },
                      { id: '#2339', table: 'Takeaway', items: '5 items', status: 'Done', color: 'bg-blue-500/20 text-blue-400' },
                    ].map((order) => (
                      <div key={order.id} className="flex items-center justify-between text-xs">
                        <span className="text-gray-300 font-mono">{order.id}</span>
                        <span className="text-gray-400">{order.table}</span>
                        <span className="text-gray-500">{order.items}</span>
                        <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${order.color}`}>{order.status}</span>
                      </div>
                    ))}
                  </div>

                  {/* Chart bar visual */}
                  <div className="bg-gray-700/30 rounded-xl p-3">
                    <div className="text-xs text-gray-400 mb-2">Weekly Sales</div>
                    <div className="flex items-end gap-1.5 h-10">
                      {chartHeights.map((_, i) => (
                        <div
                          key={i}
                          className={`chart-bar chart-bar-${i} flex-1 bg-gradient-to-t from-orange-500 to-red-500 rounded-sm opacity-80`}
                        />
                      ))}
                    </div>
                  </div>
                </div>
              </div>

              {/* Floating card — order status */}
              <div className="absolute -bottom-4 -start-6 bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-3.5 border border-gray-100 dark:border-gray-700 max-w-[180px]">
                <div className="flex items-center gap-2.5">
                  <div className="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-xl flex items-center justify-center flex-shrink-0">
                    <CheckCircle2 className="h-5 w-5 text-green-600 dark:text-green-400" />
                  </div>
                  <div>
                    <div className="text-xs font-bold text-gray-900 dark:text-white">{t.hero.card1Title}</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">{t.hero.card1Sub}</div>
                  </div>
                </div>
              </div>

              {/* Floating card — revenue */}
              <div className="absolute -top-4 -end-4 bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-3.5 border border-gray-100 dark:border-gray-700">
                <div className="flex items-center gap-2">
                  <TrendingUp className="h-5 w-5 text-orange-500" />
                  <div>
                    <div className="text-sm font-black text-gray-900 dark:text-white">12,450 EGP</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">{t.hero.card2Title}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <style jsx>{`
        .hero-grid-pattern {
          background-image:
            linear-gradient(#374151 1px, transparent 1px),
            linear-gradient(to right, #374151 1px, transparent 1px);
          background-size: 64px 64px;
        }

        ${chartHeightRules}
      `}</style>
    </section>
  )
}
