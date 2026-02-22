'use client'

import { useLanguage } from '@/contexts/LanguageContext'
import { ChefHat } from 'lucide-react'
import Link from 'next/link'

export function Footer() {
  const { t, language, dir } = useLanguage()
  const locale = language
  const currentYear = new Date().getFullYear()

  return (
    <footer className="bg-gray-950 dark:bg-black text-gray-400 py-16 px-4 sm:px-6 lg:px-8" dir={dir}>
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mb-12">
          {/* Brand */}
          <div className="col-span-2 md:col-span-1">
            <Link href={`/${locale}`} className="flex items-center gap-2.5 mb-4">
              <div className="w-9 h-9 bg-gradient-to-br from-orange-500 to-red-600 rounded-xl flex items-center justify-center">
                <ChefHat className="h-5 w-5 text-white" />
              </div>
              <span className="font-extrabold text-xl text-white">
                Restaurant<span className="text-orange-500">POS</span>
              </span>
            </Link>
            <p className="text-sm text-gray-500 leading-relaxed max-w-xs">
              {t.footer.tagline}
            </p>
          </div>

          {/* Product */}
          <div>
            <h3 className="font-bold text-white mb-4 text-sm uppercase tracking-wider">{t.footer.product}</h3>
            <ul className="space-y-3 text-sm">
              <li><Link href="#features" className="hover:text-white transition-colors">{t.footer.features}</Link></li>
              <li><Link href="#pricing" className="hover:text-white transition-colors">{t.footer.pricing}</Link></li>
              <li><Link href={`/${locale}/demo`} className="hover:text-white transition-colors">{t.footer.demo}</Link></li>
            </ul>
          </div>

          {/* Company */}
          <div>
            <h3 className="font-bold text-white mb-4 text-sm uppercase tracking-wider">{t.footer.company}</h3>
            <ul className="space-y-3 text-sm">
              <li><Link href={`/${locale}/about`} className="hover:text-white transition-colors">{t.footer.about}</Link></li>
              <li><Link href={`/${locale}/blog`} className="hover:text-white transition-colors">{t.footer.blog}</Link></li>
              <li><Link href="#contact" className="hover:text-white transition-colors">{t.footer.contact}</Link></li>
            </ul>
          </div>

          {/* Support */}
          <div>
            <h3 className="font-bold text-white mb-4 text-sm uppercase tracking-wider">{t.footer.support}</h3>
            <ul className="space-y-3 text-sm">
              <li><Link href={`/${locale}/help`} className="hover:text-white transition-colors">{t.footer.helpCenter}</Link></li>
              <li><Link href={`/${locale}/docs`} className="hover:text-white transition-colors">{t.footer.documentation}</Link></li>
              <li><Link href={`/${locale}/api`} className="hover:text-white transition-colors">{t.footer.api}</Link></li>
            </ul>
          </div>
        </div>

        {/* Bottom */}
        <div className="pt-8 border-t border-gray-800 flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-gray-500">
          <p>© {currentYear} Restaurant POS SaaS. {t.footer.rights}</p>
          <div className="flex gap-6">
            <Link href={`/${locale}/privacy`} className="hover:text-white transition-colors">Privacy</Link>
            <Link href={`/${locale}/terms`} className="hover:text-white transition-colors">Terms</Link>
          </div>
        </div>
      </div>
    </footer>
  )
}
