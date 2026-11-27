import type { Metadata } from 'next'
import { Fraunces, IBM_Plex_Sans_Arabic, Manrope } from 'next/font/google'
import { ThemeProvider } from 'next-themes'
import '@/app/globals.css'
import '@/app/marketing.css'
import '@/styles/marketing-coral.css'
import '@/styles/marketing-interactive.css'

const plex = IBM_Plex_Sans_Arabic({
  weight: ['400', '600', '700'],
  subsets: ['arabic', 'latin'],
  display: 'swap',
  variable: '--font-plex',
})
const manrope = Manrope({ subsets: ['latin'], display: 'swap', variable: '--font-manrope' })
const fraunces = Fraunces({ subsets: ['latin'], display: 'swap', variable: '--font-fraunces' })

interface RootDocumentProps {
  children: React.ReactNode
  locale: 'en' | 'ar'
}

export const rootMetadata: Metadata = {
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

export function RootDocument({ children, locale }: RootDocumentProps) {

  return (
    <html lang={locale} dir={locale === 'ar' ? 'rtl' : 'ltr'} className="light" suppressHydrationWarning>
      <body className={`${plex.variable} ${manrope.variable} ${fraunces.variable}`}>
        <ThemeProvider attribute="class" defaultTheme="light" enableSystem={false} storageKey="servio-marketing-theme">
          {children}
        </ThemeProvider>
      </body>
    </html>
  )
}
