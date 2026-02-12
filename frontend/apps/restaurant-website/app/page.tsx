import { headers } from 'next/headers';
import { DynamicHomePage } from '@/components/dynamic-home-page';
import { HeaderSSR } from '@/components/header-ssr';
import { FooterSSR } from '@/components/footer-ssr';
import { getDefaultThemeData } from '@/lib/api/get-theme-server';

/**
 * Root page (/) - renders the English homepage directly, no redirects
 */
export default async function RootPage() {
  const headersList = await headers();
  const restaurantSlug = headersList.get('x-restaurant-slug') || '';
  const locale = 'en';

  // If no valid restaurant slug, show a simple landing
  if (!restaurantSlug || restaurantSlug === 'localhost' || restaurantSlug === 'www') {
    const defaultTheme = getDefaultThemeData();
    return (
      <div className="flex flex-col min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
        <HeaderSSR themeData={defaultTheme} locale={locale} cartItemsCount={0} />
        <main className="flex-1 w-full flex items-center justify-center">
          <div className="text-center max-w-md p-8">
            <h1 className="text-3xl font-bold mb-4">Welcome</h1>
            <p className="text-gray-600 dark:text-gray-400 mb-6">
              Restaurant website is running. Access via your restaurant subdomain to see your restaurant.
            </p>
            <p className="text-sm text-gray-500">
              Development mode: localhost
            </p>
          </div>
        </main>
        <FooterSSR themeData={defaultTheme} locale={locale} />
      </div>
    );
  }

  // For real restaurant slugs, fetch data
  try {
    const apiUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1';

    const [homepageResponse, themeResponse, productsResponse] = await Promise.all([
      fetch(`${apiUrl}/public/restaurants/${restaurantSlug}/homepage`, { next: { revalidate: 300 } }),
      fetch(`${apiUrl}/public/restaurants/${restaurantSlug}/theme`, { next: { revalidate: 3600 } }),
      fetch(`${apiUrl}/public/restaurants/${restaurantSlug}/products`, { next: { revalidate: 600 } }),
    ]);

    if (!homepageResponse.ok || !themeResponse.ok || !productsResponse.ok) {
      throw new Error('Failed to fetch restaurant data');
    }

    const [homepageData, themeData, productsData] = await Promise.all([
      homepageResponse.json(),
      themeResponse.json(),
      productsResponse.json(),
    ]);

    const themeDataExtracted = themeData.data?.data?.config || themeData.data?.config || themeData;

    return (
      <div className="flex flex-col min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
        <HeaderSSR themeData={themeDataExtracted} locale={locale} cartItemsCount={0} />
        <main className="flex-1 w-full">
          <DynamicHomePage
            initialHomepageData={homepageData}
            initialThemeData={themeDataExtracted}
            initialProductsData={productsData}
            restaurantSlug={restaurantSlug}
            locale={locale}
          />
        </main>
        <FooterSSR themeData={themeDataExtracted} locale={locale} />
      </div>
    );
  } catch {
    const defaultTheme = getDefaultThemeData();
    return (
      <div className="flex flex-col min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
        <HeaderSSR themeData={defaultTheme} locale={locale} cartItemsCount={0} />
        <main className="flex-1 w-full flex items-center justify-center">
          <div className="text-center max-w-md p-8">
            <h1 className="text-2xl font-bold mb-4">Restaurant</h1>
            <p className="text-gray-600 dark:text-gray-400">
              Unable to load restaurant data. Please try again later.
            </p>
          </div>
        </main>
        <FooterSSR themeData={defaultTheme} locale={locale} />
      </div>
    );
  }
}
