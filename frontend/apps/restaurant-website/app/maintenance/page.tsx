import { Metadata } from 'next'
import { MaintenancePage } from '@/components/MaintenancePage'

export const metadata: Metadata = {
  title: 'Temporarily Unavailable | Restaurant',
  description: 'Our online ordering system is temporarily unavailable. Please try again later.',
  robots: {
    index: false,
    follow: false,
  },
}

export default async function MaintenancePageRoute({
  searchParams,
}: {
  searchParams: Promise<{ reason?: string; name?: string; email?: string; phone?: string }>
}) {
  const params = await searchParams
  const blockReason = params.reason || 'maintenance'

  // Map reason from query parameter
  const reason = blockReason === 'trial_expired' || blockReason === 'tenant_disabled' ? blockReason : 'maintenance'

  // Use tenant contact info from query params (set by middleware from API)
  const restaurantName = params.name || 'Our Restaurant'
  const contactEmail = params.email || undefined
  const contactPhone = params.phone || undefined

  return (
    <MaintenancePage
      restaurantName={restaurantName}
      reason={reason as 'trial_expired' | 'tenant_disabled' | 'maintenance' | 'error'}
      contactEmail={contactEmail}
      contactPhone={contactPhone}
    />
  )
}
