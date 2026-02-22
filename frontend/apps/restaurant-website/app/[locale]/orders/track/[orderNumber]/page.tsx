import { HeaderSSR } from '@/components/header-ssr';
import { FooterSSR } from '@/components/footer-ssr';
import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server';
import TrackPageContent from '../track-page-content';

interface PageProps {
  params: Promise<{ locale: 'en' | 'ar'; orderNumber: string }>;
}

export default async function TrackOrderNumberPage({ params }: PageProps) {
  const { locale, orderNumber } = await params;

  const themeData = (await getThemeDataServer()) ?? getDefaultThemeData();
  const initialOrderNumber = decodeURIComponent(orderNumber);

  return (
    <>
      <HeaderSSR themeData={themeData} cartItemsCount={0} locale={locale} />
      <main>
        <TrackPageContent initialOrderNumber={initialOrderNumber} />
      </main>
      <FooterSSR themeData={themeData} locale={locale} />
    </>
  );
}
