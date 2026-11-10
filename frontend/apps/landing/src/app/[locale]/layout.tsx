import { LanguageProvider } from '@/contexts/LanguageContext'
import { Toaster } from 'sonner'
import type { Metadata } from 'next'

type Locale = 'en' | 'ar'

interface LocaleLayoutProps {
  children: React.ReactNode
  params: Promise<{ locale: string }>
}

export function generateStaticParams() {
  return [{ locale: 'en' }, { locale: 'ar' }]
}

export async function generateMetadata({ params }: Pick<LocaleLayoutProps, 'params'>): Promise<Metadata> {
  const { locale } = await params
  const isArabic = locale === 'ar'
  const title = isArabic ? 'POOS | إدارة متصلة لكل تفاصيل مطعمك' : 'POOS | Connected restaurant management'
  const description = isArabic
    ? 'منصة واحدة لإدارة العمل في مطعمك من واجهة الخدمة إلى ما يدور خلفها.'
    : 'One platform for the work at the front and back of your restaurant.'

  return {
    title,
    description,
    openGraph: { title, description, siteName: 'POOS', type: 'website', locale: isArabic ? 'ar_EG' : 'en_US' },
  }
}

export default async function LocaleLayout({ children, params }: LocaleLayoutProps) {
  const { locale: rawLocale } = await params
  const locale: Locale = rawLocale === 'ar' ? 'ar' : 'en'

  return (
    <LanguageProvider locale={locale}>
      {children}
      <Toaster position={locale === 'ar' ? 'top-left' : 'top-right'} richColors />
    </LanguageProvider>
  )
}
