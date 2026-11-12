import { getThemeDataServer, getDefaultThemeData } from '@/lib/api/get-theme-server'
import { SettingsPageContent } from './settings-page-content'

interface SettingsPageProps {
  params: Promise<{
    locale: 'en' | 'ar'
  }>
}

export const metadata = {
  title: 'Settings - Customize Your Preferences',
  description: 'Manage your language, theme, and color preferences',
}

export default async function SettingsPage({ params }: SettingsPageProps) {
  const resolvedParams = await params
  const locale = resolvedParams.locale || 'en'
  const themeData = (await getThemeDataServer()) ?? getDefaultThemeData()

  return (
    <main className="min-h-screen w-full bg-white text-gray-900 dark:bg-gray-900 dark:text-gray-100">
      <SettingsPageContent locale={locale} themeData={themeData} />
    </main>
  )
}
