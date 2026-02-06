'use client'

import { useState } from 'react'

interface DisableTenantModalProps {
  tenantId: number
  tenantName: string
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
}

export function DisableTenantModal({
  tenantId,
  tenantName,
  isOpen,
  onClose,
  onSuccess,
}: DisableTenantModalProps) {
  const [reason, setReason] = useState<string>('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  if (!isOpen) return null

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!reason.trim()) {
      setError('Reason is required')
      return
    }

    setIsLoading(true)

    try {
      const response = await fetch(
        `/api/v1/super-admin/tenants/${tenantId}/disable`,
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ reason }),
        }
      )

      if (!response.ok) {
        throw new Error('Failed to disable tenant')
      }

      onSuccess()
      onClose()
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
        <h2 className="text-xl font-bold mb-2 text-red-600">Disable Tenant Access</h2>
        <p className="text-gray-600 mb-4">
          Tenant: <span className="font-semibold">{tenantName}</span>
        </p>

        <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded text-sm text-red-700">
          <p className="font-semibold mb-1">⚠️ Warning:</p>
          <ul className="list-disc list-inside space-y-1">
            <li>Tenant will immediately lose access to the dashboard</li>
            <li>Restaurant website will display a maintenance page</li>
            <li>This action can be reversed by enabling the tenant later</li>
          </ul>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Reason for disabling *
            </label>
            <textarea
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-red-500"
              rows={4}
              placeholder="e.g., Non-payment after multiple warnings, Terms of Service violation"
              required
            />
          </div>

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
              className="flex-1 px-4 py-2 text-white bg-red-600 hover:bg-red-700 rounded-md transition-colors disabled:opacity-50"
              disabled={isLoading || !reason.trim()}
            >
              {isLoading ? 'Disabling...' : 'Disable Tenant'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
