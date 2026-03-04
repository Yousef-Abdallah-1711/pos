import Link from 'next/link'
import { BarChart3, ChartNoAxesCombined, FileText, TrendingUp, Users } from 'lucide-react'

export default async function ReportsIndexPage({
  params,
}: {
  params: Promise<{ locale: string }>
}) {
  const { locale } = await params

  const links = [
    {
      href: `/${locale}/dashboard/reports/overview`,
      title: 'Overview',
      description: 'Quick snapshot across sales, profitability, and inventory.',
      icon: BarChart3,
    },
    {
      href: `/${locale}/dashboard/reports/sales`,
      title: 'Sales Analytics',
      description: 'Order and revenue trends for selected periods.',
      icon: TrendingUp,
    },
    {
      href: `/${locale}/dashboard/reports/profitability`,
      title: 'Profitability',
      description: 'Margin-focused analysis by product performance.',
      icon: ChartNoAxesCombined,
    },
    {
      href: `/${locale}/dashboard/reports/financial`,
      title: 'Financial Reports',
      description: 'P&L, balance sheet, cash flow, VAT, and related statements.',
      icon: FileText,
    },
    {
      href: `/${locale}/dashboard/reports/hr`,
      title: 'HR Reports',
      description: 'Headcount, turnover, absenteeism, overtime, leave balances, and training compliance.',
      icon: Users,
    },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Reports Center</h1>
        <p className="mt-1 text-gray-500 dark:text-gray-400">
          Access operational and financial reporting from one place.
        </p>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        {links.map((link) => (
          <Link
            key={link.href}
            href={link.href}
            className="rounded-2xl border border-gray-200 bg-white p-5 transition hover:border-primary-300 hover:shadow-sm dark:border-gray-700 dark:bg-gray-800 dark:hover:border-primary-800"
          >
            <div className="inline-flex rounded-xl bg-gray-100 p-2 text-gray-700 dark:bg-gray-900 dark:text-gray-200">
              <link.icon className="h-4 w-4" />
            </div>
            <h2 className="mt-4 text-lg font-semibold text-gray-900 dark:text-white">{link.title}</h2>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{link.description}</p>
          </Link>
        ))}
      </div>
    </div>
  )
}
