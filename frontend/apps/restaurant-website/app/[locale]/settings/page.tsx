import { HeaderSSR } from '@/components/header-ssr';
import { FooterSSR } from '@/components/footer-ssr';
import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server';
import { SettingsPageContent } from './settings-page-content';

interface SettingsPageProps {
  params: Promise<{
    locale: 'en' | 'ar';
  }>;
}

export const metadata = {
  title: 'Settings - Customize Your Preferences',
  description: 'Manage your language, theme, and color preferences',
};

export default async function SettingsPage({ params }: SettingsPageProps) {
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
        <SettingsPageContent locale={locale} themeData={themeData} />
      </main>

      {/* Footer */}
      <FooterSSR themeData={themeData} locale={locale} />
    </div>
  );
}
