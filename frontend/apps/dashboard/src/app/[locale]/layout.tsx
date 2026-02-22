import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import '../globals.css'
import '@/styles/rtl.css'
import { ThemeProvider } from 'next-themes'
import { Toaster } from 'sonner'
import AuthProvider from '@/components/providers/AuthProvider'
import { PerformanceMonitor } from '@/components/providers/PerformanceMonitor'
import { WebVitalsProvider } from '@/components/providers/WebVitalsProvider'

const inter = Inter({ subsets: ['latin'] })

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

export default async function LocaleLayout({
  children,
  params,
}: Props) {
  const { locale: paramLocale } = await params
  const locale = paramLocale === 'ar' ? 'ar' : 'en'
  const isRTL = locale === 'ar'
  const dir = isRTL ? 'rtl' : 'ltr'

  return (
    <AuthProvider>
      <ThemeProvider attribute="class" defaultTheme="light" enableSystem>
        <WebVitalsProvider />
        <PerformanceMonitor />
        <div lang={locale} dir={dir} className={inter.className}>
          {children}
        </div>
        <Toaster position={isRTL ? 'top-left' : 'top-right'} richColors />
      </ThemeProvider>
    </AuthProvider>
  )
}
