import type { Metadata } from 'next'
import { Alexandria, IBM_Plex_Sans_Arabic, Manrope } from 'next/font/google'
import './globals.css'
import './marketing.css'
import { ThemeProvider } from 'next-themes'

const alexandria = Alexandria({
  subsets: ['arabic', 'latin'],
  display: 'swap',
  variable: '--font-alexandria',
})
const plex = IBM_Plex_Sans_Arabic({
  weight: ['400', '500', '600', '700'],
  subsets: ['arabic', 'latin'],
  display: 'swap',
  variable: '--font-plex',
})
const manrope = Manrope({ subsets: ['latin'], display: 'swap', variable: '--font-manrope' })

interface RootLayoutProps {
  children: React.ReactNode
  params: Promise<{ locale?: string }>
}

export const metadata: Metadata = {
  metadataBase: process.env.NEXT_PUBLIC_LANDING_URL
    ? new URL(process.env.NEXT_PUBLIC_LANDING_URL)
    : undefined,
  title: 'SERVIO | إدارة متصلة لكل تفاصيل مطعمك',
  description: 'منصة واحدة لإدارة العمل في مطعمك من واجهة الخدمة إلى ما يدور خلفها.',
  applicationName: 'SERVIO',
  keywords: ['SERVIO', 'إدارة المطاعم', 'restaurant management'],
  openGraph: {
    title: 'SERVIO | إدارة متصلة لكل تفاصيل مطعمك',
    description: 'منصة واحدة لإدارة العمل في مطعمك من واجهة الخدمة إلى ما يدور خلفها.',
    siteName: 'SERVIO',
    type: 'website',
  },
}

export default async function RootLayout({ children, params }: RootLayoutProps) {
  const { locale: routeLocale } = await params
  const locale = routeLocale === 'ar' ? 'ar' : 'en'

  return (
    <html lang={locale} dir={locale === 'ar' ? 'rtl' : 'ltr'} suppressHydrationWarning>
      <body className={`${alexandria.variable} ${plex.variable} ${manrope.variable}`}>
        <ThemeProvider attribute="class" defaultTheme="light" forcedTheme="light">
          {children}
        </ThemeProvider>
      </body>
    </html>
  )
}
