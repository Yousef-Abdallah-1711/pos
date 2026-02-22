'use client'

import { useState } from 'react'
import { superAdminAPI } from '@/lib/api'
import type { SubscriptionHistoryItem } from '@/types'

interface ExtendTrialModalProps {
  tenantId: number
  tenantName: string
  currentTrialEndDate: string | null
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
}

export function ExtendTrialModal({
  tenantId,
  tenantName,
  currentTrialEndDate,
  isOpen,
  onClose,
  onSuccess,
}: ExtendTrialModalProps) {
  const [days, setDays] = useState<string>('15')
  const [reason, setReason] = useState<string>('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  if (!isOpen) return null

  const calculateNewDate = (): string => {
    if (!currentTrialEndDate) return 'N/A'
    const endDate = new Date(currentTrialEndDate)
    const daysNum = parseInt(days) || 0
    endDate.setDate(endDate.getDate() + daysNum)
    return endDate.toLocaleDateString()
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    const daysNum = parseInt(days)
    if (!daysNum || daysNum <= 0) {
      setError('Days must be greater than 0')
      return
    }

    setIsLoading(true)

    try {
      await superAdminAPI.extendTenantTrial(tenantId, daysNum, reason || undefined)

      onSuccess()
      onClose()
      setDays('15')
      setReason('')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-lg max-w-md w-full mx-4 p-6">
        <h2 className="text-xl font-bold mb-4">Extend Trial Period</h2>
        <p className="text-gray-600 mb-4">
          Tenant: <span className="font-semibold">{tenantName}</span>
        </p>

        {error && (
          <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Days to extend *
            </label>
            <input
              type="number"
              min="1"
              max="365"
              value={days}
              onChange={(e) => setDays(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Reason (optional)
            </label>
            <textarea
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows={3}
              placeholder="e.g., Customer requested extension"
            />
          </div>

          {currentTrialEndDate && (
            <div className="p-3 bg-blue-50 border border-blue-200 rounded">
              <p className="text-sm text-gray-600">
                New trial end date:{' '}
                <span className="font-semibold text-blue-600">
                  {calculateNewDate()}
                </span>
              </p>
            </div>
          )}

          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md transition-colors"
              disabled={isLoading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="flex-1 px-4 py-2 text-white bg-blue-600 hover:bg-blue-700 rounded-md transition-colors disabled:opacity-50"
              disabled={isLoading}
            >
              {isLoading ? 'Extending...' : 'Extend Trial'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
