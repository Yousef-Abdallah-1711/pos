import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server'
import { CartPageContent } from './cart-page-content'

interface CartPageProps {
  params: Promise<{
    locale: 'en' | 'ar'
  }>
}

export const metadata = {
  title: 'Shopping Cart - Order Food Online',
  description: 'Review your food order items and proceed to checkout',
}

export default async function CartPage({ params }: CartPageProps) {
  const resolvedParams = await params
  const locale = resolvedParams.locale || 'en'
  const themeData = (await getThemeDataServer()) ?? getDefaultThemeData()

  return (
    <main className="min-h-screen w-full bg-white text-gray-900 dark:bg-gray-900 dark:text-gray-100">
      <CartPageContent locale={locale} themeData={themeData} />
    </main>
  )
}
