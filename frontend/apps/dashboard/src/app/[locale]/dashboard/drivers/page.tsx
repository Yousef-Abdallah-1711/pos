'use client'

import { useState, useEffect } from 'react'
import { usePathname } from 'next/navigation'
import Link from 'next/link'
import { useAuthStore } from '@/stores/authStore'
import { api } from '@/lib/api'
import { hasAuthenticatedSession } from '@/lib/authSession'
import { Button } from '@/components/ui/button'
import { toast } from 'sonner'
import {
  Truck,
  Plus,
  Edit,
  Trash2,
  ToggleLeft,
  ToggleRight,
  AlertCircle,
  Loader,
} from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'

interface Driver {
  id: number
  first_name: string
  last_name: string
  email: string
  phone_number: string
  license_number: string
  vehicle_type?: string
  vehicle_number?: string
  status: string
  availability_status: string
  total_deliveries: number
  completed_deliveries: number
  average_rating: number
  active_orders_count: number
  is_verified: boolean
  created_at: string
}

export default function DriversPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'
  const t = createTranslator(locale)

  const [drivers, setDrivers] = useState<Driver[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<'all' | 'active' | 'available'>('all')
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const hasSession = hasAuthenticatedSession(user, token)

  useEffect(() => {
    if (!hasSession) {
      setLoading(false)
      return
    }
    fetchDrivers()
  }, [filter, hasSession])

  const fetchDrivers = async () => {
    setLoading(true)
    try {
      const response = await api.get('/admin/drivers')
      let driverList = response.data?.drivers || response.data || []

      if (filter === 'active') {
        driverList = driverList.filter((d: Driver) => d.status === 'active')
      } else if (filter === 'available') {
        driverList = driverList.filter((d: Driver) => d.availability_status === 'available')
      }

      setDrivers(driverList)
    } catch (error) {
      toast.error(t('drivers.loadFailed'))
      debugError('[DRIVERS_PAGE] Error:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleToggleStatus = async (driverId: number, currentStatus: string) => {
    try {
      const newStatus = currentStatus === 'active' ? 'inactive' : 'active'
      await api.put(`/admin/drivers/${driverId}/status`, { status: newStatus })
      toast.success(t('drivers.updateSuccess'))
      fetchDrivers()
    } catch (error) {
      toast.error(t('drivers.updateFailed'))
    }
  }

  const handleDeleteDriver = async (driverId: number) => {
    if (!confirm(t('drivers.confirmDelete'))) return

    try {
      await api.delete(`/admin/drivers/${driverId}`)
      toast.success(t('drivers.deleteSuccess'))
      fetchDrivers()
    } catch (error) {
      toast.error(t('drivers.deleteFailed'))
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className={`flex justify-between items-center ${isRTL ? 'flex-row-reverse' : ''}`}>
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <Truck className="h-8 w-8 text-primary" />
            {t('drivers.title')}
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">{t('drivers.manageDrivers')}</p>
        </div>
        <Link href={`/${locale}/dashboard/drivers/new`}>
          <Button>
            <Plus className={`h-4 w-4 ${isRTL ? 'ml-2' : 'mr-2'}`} />
            {t('drivers.addDriver')}
          </Button>
        </Link>
      </div>

      {/* Filters */}
      <div className={`flex gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
        {(['all', 'active', 'available'] as const).map((f) => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`px-4 py-2 rounded-lg transition font-medium ${
              filter === f
                ? 'bg-primary text-white'
                : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'
            }`}
          >
            {f === 'all'
              ? t('drivers.filters.all')
              : f === 'active'
              ? t('drivers.filters.active')
              : t('drivers.filters.available')}
          </button>
        ))}
      </div>

      {/* Drivers Table */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
        {loading ? (
          <div className="p-8 text-center">
            <Loader className="h-8 w-8 animate-spin text-primary mx-auto mb-2" />
            <p className="text-gray-600 dark:text-gray-400">{t('drivers.loading')}</p>
          </div>
        ) : drivers.length === 0 ? (
          <div className="p-8 text-center">
            <AlertCircle className="h-12 w-12 text-gray-400 mx-auto mb-2" />
            <p className="text-gray-600 dark:text-gray-400">{t('drivers.noDrivers')}</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600">
                <tr>
                  <th className={`px-6 py-3 text-left text-sm font-semibold text-gray-900 dark:text-white ${isRTL ? 'text-right' : ''}`}>
                    {t('drivers.columns.name')}
                  </th>
                  <th className={`px-6 py-3 text-left text-sm font-semibold text-gray-900 dark:text-white ${isRTL ? 'text-right' : ''}`}>
                    {t('drivers.columns.contact')}
                  </th>
                  <th className={`px-6 py-3 text-left text-sm font-semibold text-gray-900 dark:text-white ${isRTL ? 'text-right' : ''}`}>
                    {t('drivers.columns.vehicle')}
                  </th>
                  <th className={`px-6 py-3 text-left text-sm font-semibold text-gray-900 dark:text-white ${isRTL ? 'text-right' : ''}`}>
                    {t('drivers.columns.status')}
                  </th>
                  <th className={`px-6 py-3 text-left text-sm font-semibold text-gray-900 dark:text-white ${isRTL ? 'text-right' : ''}`}>
                    {t('drivers.columns.deliveries')}
                  </th>
                  <th className={`px-6 py-3 text-left text-sm font-semibold text-gray-900 dark:text-white ${isRTL ? 'text-right' : ''}`}>
                    {t('drivers.columns.rating')}
                  </th>
                  <th className={`px-6 py-3 text-left text-sm font-semibold text-gray-900 dark:text-white ${isRTL ? 'text-right' : ''}`}>
                    {t('drivers.columns.actions')}
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                {drivers.map((driver) => (
                  <tr key={driver.id} className="hover:bg-gray-50 dark:hover:bg-gray-750">
                    <td className="px-6 py-4">
                      <div className={isRTL ? 'text-right' : ''}>
                        <div className="font-semibold text-gray-900 dark:text-white">
                          {driver.first_name} {driver.last_name}
                        </div>
                        <div className="text-sm text-gray-600 dark:text-gray-400">
                          {driver.license_number}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className={isRTL ? 'text-right' : ''}>
                        <div className="text-sm text-gray-900 dark:text-white">{driver.phone_number}</div>
                        <div className="text-sm text-gray-600 dark:text-gray-400">{driver.email}</div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className={isRTL ? 'text-right' : ''}>
                        <div className="text-sm text-gray-900 dark:text-white">
                          {driver.vehicle_type || t('common.notAvailable', 'N/A')}
                        </div>
                        <div className="text-sm text-gray-600 dark:text-gray-400">
                          {driver.vehicle_number || t('common.notAvailable', '-')}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-col gap-1">
                        <span
                          className={`px-3 py-1 rounded-full text-xs font-semibold w-fit ${
                            driver.status === 'active'
                              ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300'
                              : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300'
                          }`}
                        >
                          {t(`drivers.statusLabels.${driver.status}`)}
                        </span>
                        <span
                          className={`px-3 py-1 rounded-full text-xs font-semibold w-fit ${
                            driver.availability_status === 'available'
                              ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300'
                              : 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300'
                          }`}
                        >
                          {t(`drivers.statusLabels.${driver.availability_status}`)}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className={isRTL ? 'text-right' : ''}>
                        <div className="text-sm font-medium text-gray-900 dark:text-white">
                          {driver.completed_deliveries}/{driver.total_deliveries}
                        </div>
                        <div className="text-xs text-gray-600 dark:text-gray-400">
                          {driver.total_deliveries > 0
                            ? Math.round(
                              (driver.completed_deliveries / driver.total_deliveries) * 100
                            )
                            : 0}
                          %
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className={isRTL ? 'text-right' : ''}>
                        <div className="text-sm font-semibold text-gray-900 dark:text-white">
                          * {driver.average_rating.toFixed(1)}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className={`flex gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
                        <Link href={`/${locale}/dashboard/drivers/${driver.id}`}>
                          <Button size="sm" variant="outline">
                            <Edit className="h-4 w-4" />
                          </Button>
                        </Link>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleToggleStatus(driver.id, driver.status)}
                          title={
                            driver.status === 'active'
                              ? t('drivers.actionLabels.deactivate')
                              : t('drivers.actionLabels.activate')
                          }
                        >
                          {driver.status === 'active' ? (
                            <ToggleRight className="h-4 w-4 text-green-600" />
                          ) : (
                            <ToggleLeft className="h-4 w-4 text-gray-400" />
                          )}
                        </Button>
                        <Button
                          size="sm"
                          variant="destructive"
                          onClick={() => handleDeleteDriver(driver.id)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Summary */}
      {drivers.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
            <p className="text-sm text-gray-600 dark:text-gray-400">{t('drivers.stats.totalDrivers')}</p>
            <p className="text-2xl font-bold text-gray-900 dark:text-white">{drivers.length}</p>
          </div>
          <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
            <p className="text-sm text-gray-600 dark:text-gray-400">{t('drivers.stats.activeDrivers')}</p>
            <p className="text-2xl font-bold text-green-600">
              {drivers.filter((d) => d.status === 'active').length}
            </p>
          </div>
          <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
            <p className="text-sm text-gray-600 dark:text-gray-400">{t('drivers.stats.availableDrivers')}</p>
            <p className="text-2xl font-bold text-blue-600">
              {drivers.filter((d) => d.availability_status === 'available').length}
            </p>
          </div>
          <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
            <p className="text-sm text-gray-600 dark:text-gray-400">{t('drivers.stats.activeOrders')}</p>
            <p className="text-2xl font-bold text-purple-600">
              {drivers.reduce((sum, d) => sum + d.active_orders_count, 0)}
            </p>
          </div>
        </div>
      )}
    </div>
  )
}
