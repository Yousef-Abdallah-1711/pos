import { redirect } from 'next/navigation'

type SettingsPageProps = {
  params: Promise<{ locale: string }>
}

export default async function SettingsPage({ params }: SettingsPageProps) {
  const { locale: requestedLocale } = await params
  const locale = requestedLocale === 'ar' ? 'ar' : 'en'

  redirect(`/${locale}/dashboard/settings/account`)
}
