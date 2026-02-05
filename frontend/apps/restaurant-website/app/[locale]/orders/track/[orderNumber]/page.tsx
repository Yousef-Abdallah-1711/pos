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
  const themeDataExtracted = themeData.data?.data?.config || themeData.data?.config || themeData;
  const initialOrderNumber = decodeURIComponent(orderNumber);

  return (
    <>
      <HeaderSSR themeData={themeDataExtracted} cartItemsCount={0} locale={locale} />
      <main>
        <TrackPageContent initialOrderNumber={initialOrderNumber} />
      </main>
      <FooterSSR themeData={themeDataExtracted} locale={locale} />
    </>
  );
}
