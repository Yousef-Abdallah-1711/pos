import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server'
import { OrdersPageContent } from './orders-page-content'

interface OrdersPageProps {
  params: Promise<{
    locale: 'en' | 'ar'
  }>
}

export const metadata = {
  title: 'My Orders - Track Your Orders',
  description: 'View and track all your food orders',
}

export default async function OrdersPage({ params }: OrdersPageProps) {
  const resolvedParams = await params
  const locale = resolvedParams.locale || 'en'
  const themeData = (await getThemeDataServer()) ?? getDefaultThemeData()

  return (
    <main className="min-h-screen w-full bg-white text-gray-900 dark:bg-gray-900 dark:text-gray-100">
      <OrdersPageContent locale={locale} themeData={themeData} />
    </main>
  )
}
