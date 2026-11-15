import type { Metadata } from 'next'
import type { ReactNode } from 'react'
import { Toaster } from 'sonner'
import { LanguageProvider } from '@/contexts/LanguageContext'
import type { Locale } from '@/content/site'

export const metadata: Metadata = { robots: { index: false, follow: false } }

interface AuthLayoutProps {
  children: ReactNode
  params: Promise<{ locale: string }>
}

export default async function AuthLayout({ children, params }: AuthLayoutProps) {
  const { locale: value } = await params
  const locale: Locale = value === 'ar' ? 'ar' : 'en'
  return <LanguageProvider locale={locale}>{children}<Toaster position={locale === 'ar' ? 'top-left' : 'top-right'} richColors /></LanguageProvider>
}
