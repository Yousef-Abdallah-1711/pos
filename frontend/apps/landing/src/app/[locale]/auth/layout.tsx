import type { Metadata } from 'next'
import type { ReactNode } from 'react'
import { Toaster } from 'sonner'

export const metadata: Metadata = { robots: { index: false, follow: false } }

interface AuthLayoutProps {
  children: ReactNode
  params: Promise<{ locale: string }>
}

export default async function AuthLayout({ children, params }: AuthLayoutProps) {
  const { locale } = await params
  return <>{children}<Toaster position={locale === 'ar' ? 'top-left' : 'top-right'} richColors /></>
}
