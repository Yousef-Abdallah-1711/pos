import { HeaderSSR } from '@/components/header-ssr';
import { FooterSSR } from '@/components/footer-ssr';
import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server';
import TrackPageContent from './track-page-content';

interface PageProps {
  params: Promise<{ locale: 'en' | 'ar' }>;
}

export default async function TrackOrderPage({ params }: PageProps) {
  const { locale } = await params;

  // Fetch theme server-side
  const themeData = (await getThemeDataServer()) ?? getDefaultThemeData();

  return (
    <>
      <HeaderSSR themeData={themeData} cartItemsCount={0} locale={locale} />
      <main>
        <TrackPageContent />
      </main>
      <FooterSSR themeData={themeData} locale={locale} />
    </>
  );
}
