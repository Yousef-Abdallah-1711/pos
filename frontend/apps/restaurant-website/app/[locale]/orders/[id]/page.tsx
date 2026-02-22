import { HeaderSSR } from '@/components/header-ssr';
import { FooterSSR } from '@/components/footer-ssr';
import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server';
import OrderDetailContent from './order-detail-content';

interface PageProps {
  params: Promise<{ locale: 'en' | 'ar'; id: string }>;
}

export default async function OrderDetailsPage({ params }: PageProps) {
  const { locale } = await params;

  // Fetch theme server-side
  const themeData = (await getThemeDataServer()) ?? getDefaultThemeData();

  return (
    <>
      <HeaderSSR themeData={themeData} cartItemsCount={0} locale={locale} />
      <main>
        <OrderDetailContent />
      </main>
      <FooterSSR themeData={themeData} locale={locale} />
    </>
  );
}
