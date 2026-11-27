import type { Metadata } from 'next'
import { Inter, Manrope, Fraunces, IBM_Plex_Sans_Arabic } from 'next/font/google'
import { Suspense } from 'react'
import '../globals.css'
import '@/styles/rtl.css'
import { ThemeProvider } from 'next-themes'
import { Toaster } from 'sonner'
import { NextIntlClientProvider } from 'next-intl'
import AuthProvider from '@/components/providers/AuthProvider'
import { PerformanceMonitor } from '@/components/providers/PerformanceMonitor'
import { WebVitalsProvider } from '@/components/providers/WebVitalsProvider'
import subscriptionEn from '@pos-saas/i18n/locales/en/subscription.json'
import subscriptionAr from '@pos-saas/i18n/locales/ar/subscription.json'

const inter = Inter({ subsets: ['latin'] })
const manrope = Manrope({ subsets: ['latin'], display: 'swap', variable: '--font-workspace-manrope' })
const fraunces = Fraunces({ subsets: ['latin'], display: 'swap', variable: '--font-workspace-fraunces' })
const plex = IBM_Plex_Sans_Arabic({ subsets: ['arabic', 'latin'], weight: ['400', '600', '700'], display: 'swap', preload: false, variable: '--font-workspace-arabic' })

export const metadata: Metadata = {
  title: 'Dashboard - Restaurant POS SaaS',
  description: 'Manage your restaurants with our powerful POS system',
}

type Props = {
  children: React.ReactNode
  params: Promise<{
    locale: string
  }>
}

export function generateStaticParams() {
  return [{ locale: 'en' }, { locale: 'ar' }]
}

async function getLocaleMessages(locale: 'en' | 'ar') {
  if (locale === 'ar') {
    return (await import('@/i18n/messages/ar.json')).default
  }
  return (await import('@/i18n/messages/en.json')).default
}

export default async function LocaleLayout({
  children,
  params,
}: Props) {
  const { locale: paramLocale } = await params
  const locale = paramLocale === 'ar' ? 'ar' : 'en'
  const isRTL = locale === 'ar'
  const dir = isRTL ? 'rtl' : 'ltr'
  const appMessages = await getLocaleMessages(locale)
  const messages = { ...appMessages, subscription: locale === 'ar' ? subscriptionAr : subscriptionEn }

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      <Suspense
        fallback={<div lang={locale} dir={dir} className={`${inter.className} min-h-screen bg-gray-50 dark:bg-gray-900`} />}
      >
        <AuthProvider>
          <ThemeProvider attribute="class" defaultTheme="light" enableSystem>
            <WebVitalsProvider />
            <PerformanceMonitor />
            <div lang={locale} dir={dir} className={`${inter.className} ${manrope.variable} ${fraunces.variable} ${plex.variable}`}>
              {children}
            </div>
            <Toaster position={isRTL ? 'top-left' : 'top-right'} richColors />
          </ThemeProvider>
        </AuthProvider>
      </Suspense>
    </NextIntlClientProvider>
  )
}
