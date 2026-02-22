import { LanguageProvider } from '@/contexts/LanguageContext'
import { Toaster } from 'sonner'

type Locale = 'en' | 'ar'

interface LocaleLayoutProps {
  children: React.ReactNode
  params: Promise<{ locale: string }>
}

export function generateStaticParams() {
  return [{ locale: 'en' }, { locale: 'ar' }]
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
