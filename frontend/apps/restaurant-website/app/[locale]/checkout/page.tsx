import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server'
import { CheckoutPageContent } from './checkout-page-content'

interface CheckoutPageProps {
  params: Promise<{
    locale: 'en' | 'ar'
  }>
}

export const metadata = {
  title: 'Checkout - Complete Your Order',
  description: 'Review your order and complete the checkout process',
}

export default async function CheckoutPage({ params }: CheckoutPageProps) {
  const resolvedParams = await params
  const locale = resolvedParams.locale || 'en'
  const themeData = (await getThemeDataServer()) ?? getDefaultThemeData()

  return (
    <main className="min-h-screen w-full bg-white text-gray-900 dark:bg-gray-900 dark:text-gray-100">
      <CheckoutPageContent locale={locale} themeData={themeData} />
    </main>
  )
}
