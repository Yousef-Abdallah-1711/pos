'use client'

import { useState, useEffect } from 'react'
import { usePathname } from 'next/navigation'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { toast } from 'sonner'
import { X, Loader, Star, TrendingUp, Truck } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'

interface Driver {
  id: number
  first_name: string
  last_name: string
  email: string
  phone_number: string
  status: string
  availability_status: string
  average_rating: number
  total_deliveries: number
  completed_deliveries: number
  active_orders_count: number
}

interface DriverAssignmentModalProps {
  isOpen: boolean
  orderId: number
  orderNumber: string
  onClose: () => void
  onAssignSuccess: () => void
}

export function DriverAssignmentModal({
  isOpen,
  orderId,
  orderNumber,
  onClose,
  onAssignSuccess,
}: DriverAssignmentModalProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const [drivers, setDrivers] = useState<Driver[]>([])
  const [loading, setLoading] = useState(true)
  const [assigning, setAssigning] = useState(false)
  const [selectedDriverId, setSelectedDriverId] = useState<number | null>(null)
  const [notes, setNotes] = useState('')

  useEffect(() => {
    if (isOpen) {
      fetchAvailableDrivers()
    }
  }, [isOpen])

  const fetchAvailableDrivers = async () => {
    setLoading(true)
    try {
      const response = await api.get('/admin/drivers/available')
      setDrivers(response.data?.drivers || [])
    } catch (error) {
      debugError('[FETCH_DRIVERS] Error:', error)
      toast.error(t('drivers.loadFailed'))
    } finally {
      setLoading(false)
    }
  }

  const handleAssign = async () => {
    if (!selectedDriverId) {
      toast.error(t('drivers.assignmentModal.selectDriver'))
      return
    }

    setAssigning(true)
    try {
      await api.post(`/admin/orders/${orderId}/assign`, {
        driver_id: selectedDriverId,
        assignment_notes: notes,
      })
      toast.success(t('drivers.assignmentModal.assignSuccess'))
      setSelectedDriverId(null)
      setNotes('')
      onAssignSuccess()
      onClose()
    } catch (error: any) {
      console.error('[ASSIGN_ORDER] Error:', error)
      const errorMessage = error.response?.data?.message || t('drivers.assignmentModal.assignFailed')
      toast.error(errorMessage)
    } finally {
      setAssigning(false)
    }
  }

  if (!isOpen) return null

  const completionRate = (driver: Driver) =>
    driver.total_deliveries > 0
      ? Math.round((driver.completed_deliveries / driver.total_deliveries) * 100)
      : 0

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] flex flex-col">
        {/* Header */}
        <div className="flex justify-between items-center p-6 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
            {t('drivers.assignmentModal.title')}
          </h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
          >
            <X className="h-6 w-6" />
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-6">
          <div className="mb-6">
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {t('orders.orderNumber')}: <span className="font-semibold text-gray-900 dark:text-white">{orderNumber}</span>
            </p>
          </div>

          {loading ? (
            <div className="flex justify-center items-center h-64">
              <Loader className="h-8 w-8 animate-spin text-primary" />
            </div>
          ) : drivers.length === 0 ? (
            <div className="flex justify-center items-center h-64">
              <p className="text-gray-500 dark:text-gray-400">{t('drivers.assignmentModal.noDriversAvailable')}</p>
            </div>
          ) : (
            <div className="space-y-3">
              {drivers.map((driver) => {
                const rate = completionRate(driver)
                return (
                  <div
                    key={driver.id}
                    onClick={() => setSelectedDriverId(driver.id)}
                    className={`p-4 border-2 rounded-lg cursor-pointer transition ${
                      selectedDriverId === driver.id
                        ? 'border-primary bg-primary/5 dark:bg-primary/10'
                        : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
                    }`}
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <h3 className="font-semibold text-gray-900 dark:text-white">
                          {driver.first_name} {driver.last_name}
                        </h3>
                        <p className="text-sm text-gray-600 dark:text-gray-400">
                          {driver.email}
                        </p>
                        <p className="text-sm text-gray-600 dark:text-gray-400">
                          {driver.phone_number}
                        </p>
                      </div>

                      {/* Radio button */}
                      <div className="mt-1">
                        <input
                          type="radio"
                          name="driver"
                          value={driver.id}
                          checked={selectedDriverId === driver.id}
                          onChange={() => setSelectedDriverId(driver.id)}
                          className="h-4 w-4 text-primary"
                        />
                      </div>
                    </div>

                    {/* Stats */}
                    <div className="mt-4 grid grid-cols-3 gap-3">
                      <div className="flex items-center gap-2">
                        <Star className="h-4 w-4 text-yellow-500" />
                        <div>
                          <p className="text-xs text-gray-600 dark:text-gray-400">{t('drivers.assignmentModal.rating')}</p>
                          <p className="font-semibold text-gray-900 dark:text-white">
                            {driver.average_rating.toFixed(1)}/5
                          </p>
                        </div>
                      </div>

                      <div className="flex items-center gap-2">
                        <TrendingUp className="h-4 w-4 text-green-500" />
                        <div>
                          <p className="text-xs text-gray-600 dark:text-gray-400">{t('drivers.assignmentModal.completion')}</p>
                          <p className="font-semibold text-gray-900 dark:text-white">
                            {rate}%
                          </p>
                        </div>
                      </div>

                      <div className="flex items-center gap-2">
                        <Truck className="h-4 w-4 text-primary" />
                        <div>
                          <p className="text-xs text-gray-600 dark:text-gray-400">{t('drivers.assignmentModal.activeOrders')}</p>
                          <p className="font-semibold text-gray-900 dark:text-white">
                            {driver.active_orders_count}
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          )}

          {/* Notes Section */}
          {!loading && drivers.length > 0 && (
            <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                {t('drivers.assignmentModal.notes')}
              </label>
              <textarea
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder={t('drivers.assignmentModal.notesPlaceholder')}
                rows={3}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              />
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex justify-end gap-3 p-6 border-t border-gray-200 dark:border-gray-700">
          <Button
            type="button"
            variant="outline"
            onClick={onClose}
            disabled={assigning}
          >
            {t('common.cancel')}
          </Button>
          <Button
            type="button"
            onClick={handleAssign}
            disabled={!selectedDriverId || loading || assigning}
            className="flex items-center gap-2"
          >
            {assigning && <Loader className="h-4 w-4 animate-spin" />}
            {assigning ? t('common.saving') : t('drivers.assignDriver')}
          </Button>
        </div>
      </div>
    </div>
  )
}
