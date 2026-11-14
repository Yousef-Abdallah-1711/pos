import { RootDocument, rootMetadata } from '@/components/RootDocument'
import { LanguageProvider } from '@/contexts/LanguageContext'
import { Toaster } from 'sonner'
import type { Metadata } from 'next'
import { notFound } from 'next/navigation'

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
  if (locale !== 'en' && locale !== 'ar') notFound()

  const isArabic = locale === 'ar'
  const title = isArabic ? 'SERVIO | إدارة متصلة لكل تفاصيل مطعمك' : 'SERVIO | Connected restaurant management'
  const description = isArabic
    ? 'منصة واحدة لإدارة العمل في مطعمك من واجهة الخدمة إلى ما يدور خلفها.'
    : 'One platform for the work at the front and back of your restaurant.'

  return {
    ...rootMetadata,
    title,
    description,
    alternates: {
      canonical: `/${locale}`,
      languages: { en: '/en', ar: '/ar' },
    },
    openGraph: { title, description, siteName: 'SERVIO', type: 'website', locale: isArabic ? 'ar_EG' : 'en_US' },
  }
}

export default async function LocaleLayout({ children, params }: LocaleLayoutProps) {
  const { locale: rawLocale } = await params
  if (rawLocale !== 'en' && rawLocale !== 'ar') notFound()
  const locale: Locale = rawLocale === 'ar' ? 'ar' : 'en'

  return (
    <RootDocument locale={locale}><LanguageProvider locale={locale}>
      {children}
      <Toaster position={locale === 'ar' ? 'top-left' : 'top-right'} richColors />
    </LanguageProvider></RootDocument>
  )
}
