import { headers } from 'next/headers';
import { DynamicHomePage } from '@/components/homepage/dynamic-home-page';
import { HeaderSSR } from '@/components/layout/header-ssr';
import { FooterSSR } from '@/components/layout/footer-ssr';
import { extractThemeConfig, getDefaultThemeData } from '@/lib/api/get-theme-server';
import { getApiBaseUrl } from '@/lib/config/api';
import { debugError } from '@/lib/debug';
import { AlertCircle } from 'lucide-react';
import type { Metadata } from 'next';

interface PageProps {
  params: Promise<{
    locale: 'en' | 'ar';
  }>;
}

/**
 * Generate dynamic metadata for SEO
 * Fetches restaurant theme data to create Open Graph tags for social sharing
 */
export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const resolvedParams = await params;
  const headersList = await headers();
  const restaurantSlug = headersList.get('x-restaurant-slug') || '';
  const locale = resolvedParams.locale || 'en';

  if (!restaurantSlug) {
    return {
      title: 'Restaurant Not Found',
      description: 'The requested restaurant could not be found.',
    };
  }

  try {
    // Fetch theme data for metadata (used for OG tags)
    const apiUrl = getApiBaseUrl();
    const themeResponse = await fetch(
      `${apiUrl}/public/restaurants/${restaurantSlug}/theme`,
      {
        next: { revalidate: 3600 }, // Cache for 1 hour
      }
    );

    if (!themeResponse.ok) {
      throw new Error('Failed to fetch theme data');
    }

    const themeData = await themeResponse.json();
    const theme = extractThemeConfig(themeData) || {};

    // Extract metadata from theme
    const siteTitle =
      theme.identity?.siteTitle ||
      theme.identity?.site_title ||
      theme.site_title ||
      'Restaurant';
    const description =
      theme.footer?.companyInfo?.description ||
      theme.footer?.companyDescription ||
      theme.footer?.company_description ||
      'Order delicious food online from our restaurant.';
    const logoUrl =
      theme.identity?.logoUrl ||
      theme.identity?.logo_url ||
      theme.logo_url ||
      '';

    const canonicalUrl = `https://${restaurantSlug}.example.com/${locale}`;

    return {
      title: `${siteTitle} - Order Food Online`,
      description: description,
      alternates: {
        // Hreflang for language alternates
        languages: {
          'en': `https://${restaurantSlug}.example.com/en`,
          'ar': `https://${restaurantSlug}.example.com/ar`,
        },
        canonical: canonicalUrl,
      },
      openGraph: {
        title: siteTitle,
        description: description,
        type: 'website',
        url: canonicalUrl,
        siteName: siteTitle,
        images: logoUrl
          ? [
              {
                url: logoUrl,
                width: 1200,
                height: 630,
                alt: siteTitle,
              },
            ]
          : [],
      },
      twitter: {
        card: 'summary_large_image',
        title: siteTitle,
        description: description,
        images: logoUrl ? [logoUrl] : [],
      },
      keywords: [siteTitle, 'restaurant', 'food delivery', 'order online'],
      authors: [{ name: siteTitle }],
      creator: siteTitle,
      // Additional metadata for search engines
      applicationName: siteTitle,
      generator: 'Next.js',
    };
  } catch (error) {
    debugError('Failed to generate metadata:', error);
    return {
      title: 'Restaurant',
      description: 'Order delicious food online.',
    };
  }
}

export default async function HomePage({ params }: PageProps) {
  const resolvedParams = await params;
  const headersList = await headers();
  const restaurantSlug = headersList.get('x-restaurant-slug') || '';
  const locale = resolvedParams.locale || 'en';

  // Validate restaurant slug
  if (!restaurantSlug) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-white dark:bg-gray-900">
        <div className="text-center max-w-md">
          <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
          <h1 className="text-2xl font-bold mb-2 text-gray-900 dark:text-white">Restaurant Not Found</h1>
          <p className="text-gray-600 dark:text-gray-400 mb-6">
            Restaurant not found. Please access via your restaurant subdomain.
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            If you own a restaurant, please contact support to set up your account.
          </p>
        </div>
      </div>
    );
  }

  try {
    // Fetch all data in PARALLEL (not sequential!)
    const apiUrl = getApiBaseUrl();

    const [homepageResponse, themeResponse, productsResponse] = await Promise.all([
      fetch(`${apiUrl}/public/restaurants/${restaurantSlug}/homepage`, {
        next: { revalidate: 300 }, // ISR: cache for 5 minutes
      }),
      fetch(`${apiUrl}/public/restaurants/${restaurantSlug}/theme`, {
        next: { revalidate: 3600 }, // Cache theme for 1 hour (rarely changes)
      }),
      fetch(`${apiUrl}/public/restaurants/${restaurantSlug}/products`, {
        next: { revalidate: 600 }, // Cache for 10 minutes
      }),
    ]);

    // Check if all responses are ok
    if (!homepageResponse.ok || !themeResponse.ok || !productsResponse.ok) {
      const failedResponse = !homepageResponse.ok
        ? homepageResponse
        : !themeResponse.ok
          ? themeResponse
          : productsResponse;

      if (failedResponse.status === 404) {
        return (
          <div className="flex items-center justify-center min-h-screen bg-white dark:bg-gray-900">
            <div className="text-center max-w-md">
              <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
              <h1 className="text-2xl font-bold mb-2 text-gray-900 dark:text-white">Restaurant Not Found</h1>
              <p className="text-gray-600 dark:text-gray-400 mb-6">
                The restaurant you are looking for does not exist.
              </p>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                If you own a restaurant, please contact support to set up your account.
              </p>
            </div>
          </div>
        );
      }

      throw new Error(`API Error: ${failedResponse.status}`);
    }

    // Parse all responses in parallel
    const [homepageData, themeData, productsData] = await Promise.all([
      homepageResponse.json(),
      themeResponse.json(),
      productsResponse.json(),
    ]);

    const themeDataExtracted = extractThemeConfig(themeData) || getDefaultThemeData();

    return (
      <div className="flex flex-col min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
        {/* Unified Header */}
        <HeaderSSR
          themeData={themeDataExtracted}
          locale={locale}
          cartItemsCount={0}
        />

        {/* Main Content */}
        <main className="flex-1 w-full">
          <DynamicHomePage
            initialHomepageData={homepageData}
            initialThemeData={themeDataExtracted}
            initialProductsData={productsData}
            restaurantSlug={restaurantSlug}
            locale={locale}
          />
        </main>

        {/* Unified Footer */}
        <FooterSSR
          themeData={themeDataExtracted}
          locale={locale}
        />
      </div>
    );
  } catch (error) {
    debugError('Failed to load restaurant data:', error);
    return (
      <div className="flex items-center justify-center min-h-screen bg-white dark:bg-gray-900">
        <div className="text-center max-w-md">
          <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
          <h1 className="text-2xl font-bold mb-2 text-gray-900 dark:text-white">Error Loading Restaurant</h1>
          <p className="text-gray-600 dark:text-gray-400 mb-6">
            Unable to load restaurant data. Please try again later.
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            If the problem persists, please contact support.
          </p>
        </div>
      </div>
    );
  }
}
