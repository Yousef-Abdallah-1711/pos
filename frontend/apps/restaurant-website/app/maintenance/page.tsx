import { Metadata } from 'next'
import { MaintenancePage } from '@/components/MaintenancePage'
import { headers } from 'next/headers'

export const metadata: Metadata = {
  title: 'Temporarily Unavailable | Restaurant',
  description: 'Our online ordering system is temporarily unavailable. Please try again later.',
  robots: {
    index: false,
    follow: false,
  },
}

export default async function MaintenancePageRoute() {
  const headerList = await headers()
  const slug = headerList.get('x-restaurant-slug') || ''
  const blockReason = headerList.get('x-block-reason') || 'maintenance'

  // Map reason from middleware header
  const reason = blockReason === 'trial_expired' || blockReason === 'tenant_disabled' ? blockReason : 'maintenance'

  return (
    <MaintenancePage
      restaurantName="Our Restaurant"
      reason={reason as 'trial_expired' | 'tenant_disabled' | 'maintenance' | 'error'}
      contactEmail="contact@restaurant.com"
      contactPhone="(555) 123-4567"
    />
  )
}
