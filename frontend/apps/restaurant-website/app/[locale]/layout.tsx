import { NextIntlClientProvider } from 'next-intl'
import { notFound } from 'next/navigation'
import enMessages from '@/i18n/messages/en.json'
import arMessages from '@/i18n/messages/ar.json'
import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server'
import { LayoutChrome } from '@/components/layout/layout-chrome'
import '../../styles/globals.css'

type Locale = 'en' | 'ar'

const messages = {
  en: enMessages,
  ar: arMessages,
}

export default async function LocaleLayout({
  children,
  params,
}: {
  children: React.ReactNode
  params: Promise<{ locale: string }>
}) {
  const { locale: localeParam } = await params

  if (localeParam !== 'en' && localeParam !== 'ar') {
    notFound()
  }

  const locale = localeParam as Locale
  const direction = locale === 'ar' ? 'rtl' : 'ltr'
  const themeData = (await getThemeDataServer()) || getDefaultThemeData()

  return (
    <NextIntlClientProvider locale={locale} messages={messages[locale]}>
      <div lang={locale} dir={direction} className={locale === 'ar' ? 'font-arabic' : ''}>
        <LayoutChrome locale={locale} themeData={themeData}>
          {children}
        </LayoutChrome>
      </div>
    </NextIntlClientProvider>
  )
}
