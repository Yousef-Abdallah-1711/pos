import { headers } from 'next/headers';
import { HeaderSSR } from '@/components/layout/header-ssr';
import { FooterSSR } from '@/components/layout/footer-ssr';
import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server';
import { MenuPageContent } from './menu-page-content';

interface MenuPageProps {
  params: Promise<{
    locale: 'en' | 'ar';
  }>;
}

export const metadata = {
  title: 'Menu - Order Food Online',
  description: 'Browse our delicious menu and add items to your cart',
};

export default async function MenuPage({ params }: MenuPageProps) {
  const resolvedParams = await params;
  const locale = resolvedParams.locale || 'en';

  // Get restaurant slug from headers
  const headersList = await headers();
  const restaurantSlug = headersList.get('x-restaurant-slug') || 'demo';

  // Fetch theme server-side
  const themeData = (await getThemeDataServer()) ?? getDefaultThemeData();

  return (
    <div className="flex flex-col min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
      {/* Header */}
      <HeaderSSR themeData={themeData} locale={locale} cartItemsCount={0} />

      {/* Main Content */}
      <main className="flex-1 w-full">
        <MenuPageContent locale={locale} themeData={themeData} restaurantSlug={restaurantSlug} />
      </main>

      {/* Footer */}
      <FooterSSR themeData={themeData} locale={locale} />
    </div>
  );
}
