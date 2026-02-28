import { redirect } from 'next/navigation'

export default async function ReportsIndexPage({
  params,
}: {
  params: Promise<{ locale: string }>
}) {
  const { locale } = await params
  redirect(`/${locale}/dashboard/reports/overview`)
}

