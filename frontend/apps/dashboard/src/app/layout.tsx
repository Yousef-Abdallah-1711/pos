import type { Metadata } from 'next'
import { headers } from 'next/headers'
import { Analytics } from '@vercel/analytics/react'
import './globals.css'
import '@/styles/rtl.css'

export const metadata: Metadata = {
  title: 'Dashboard - Restaurant POS SaaS',
  description: 'Manage your restaurants with our powerful POS system',
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const locale = (await headers()).get('x-dashboard-locale') === 'ar' ? 'ar' : 'en'
  const enableVercelAnalytics =
    process.env.NEXT_PUBLIC_ENABLE_VERCEL_ANALYTICS === 'true' || process.env.VERCEL === '1'

  return (
    <html lang={locale} dir={locale === 'ar' ? 'rtl' : 'ltr'} suppressHydrationWarning>
      <body>
        {children}
        {enableVercelAnalytics ? <Analytics /> : null}
      </body>
    </html>
  )
}
