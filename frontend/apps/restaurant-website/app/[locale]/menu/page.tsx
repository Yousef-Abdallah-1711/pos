import { headers } from 'next/headers'
import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server'
import { MenuPageContent } from './menu-page-content'

interface MenuPageProps {
  params: Promise<{
    locale: 'en' | 'ar'
  }>
}

export const metadata = {
  title: 'Menu - Order Food Online',
  description: 'Browse our delicious menu and add items to your cart',
}

export default async function MenuPage({ params }: MenuPageProps) {
  const resolvedParams = await params
  const locale = resolvedParams.locale || 'en'

  const headersList = await headers()
  const restaurantSlug = headersList.get('x-restaurant-slug') || 'demo'

  const themeData = (await getThemeDataServer()) ?? getDefaultThemeData()

  return (
    <main className="min-h-screen w-full bg-white text-gray-900 dark:bg-gray-900 dark:text-gray-100">
      <MenuPageContent locale={locale} themeData={themeData} restaurantSlug={restaurantSlug} />
    </main>
  )
}
