import { AlertCircle } from 'lucide-react';
import { HeaderSSR } from '@/components/layout/header-ssr';
import { FooterSSR } from '@/components/layout/footer-ssr';
import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server';
import { CartPageContent } from './cart-page-content';

interface CartPageProps {
  params: Promise<{
    locale: 'en' | 'ar';
  }>;
}

export const metadata = {
  title: 'Shopping Cart - Order Food Online',
  description: 'Review your food order items and proceed to checkout',
};

export default async function CartPage({ params }: CartPageProps) {
  const resolvedParams = await params;
  const locale = resolvedParams.locale || 'en';

  // Fetch theme server-side
  const themeData = (await getThemeDataServer()) ?? getDefaultThemeData();

  return (
    <div className="flex flex-col min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
      {/* Header */}
      <HeaderSSR themeData={themeData} locale={locale} cartItemsCount={0} />

      {/* Main Content */}
      <main className="flex-1 w-full">
        <CartPageContent locale={locale} themeData={themeData} />
      </main>

      {/* Footer */}
      <FooterSSR themeData={themeData} locale={locale} />
    </div>
  );
}
