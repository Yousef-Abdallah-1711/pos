'use client'

import { useState } from 'react'

interface EnableTenantModalProps {
  tenantId: number
  tenantName: string
  disabledReason?: string | null
  disabledAt?: string | null
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
}

export function EnableTenantModal({
  tenantId,
  tenantName,
  disabledReason,
  disabledAt,
  isOpen,
  onClose,
  onSuccess,
}: EnableTenantModalProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  if (!isOpen) return null

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setIsLoading(true)

    try {
      const response = await fetch(
        `/api/v1/super-admin/tenants/${tenantId}/enable`,
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
        }
      )

      if (!response.ok) {
        throw new Error('Failed to enable tenant')
      }

      onSuccess()
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-lg max-w-md w-full mx-4 p-6">
        <h2 className="text-xl font-bold mb-4 text-green-600">Re-enable Tenant</h2>
        <p className="text-gray-600 mb-4">
          Tenant: <span className="font-semibold">{tenantName}</span>
        </p>

        {disabledAt && (
          <div className="mb-4 p-3 bg-gray-50 border border-gray-200 rounded text-sm">
            <p className="text-gray-600">
              <span className="font-semibold">Disabled at:</span>{' '}
              {new Date(disabledAt).toLocaleString()}
            </p>
            {disabledReason && (
              <p className="text-gray-600 mt-2">
                <span className="font-semibold">Reason:</span> {disabledReason}
              </p>
            )}
          </div>
        )}

        <div className="mb-4 p-3 bg-green-50 border border-green-200 rounded text-sm text-green-700">
          <p className="font-semibold mb-1">✓ Enabling this tenant will:</p>
          <ul className="list-disc list-inside space-y-1">
            <li>Restore access to the dashboard</li>
            <li>Restaurant website will be back online</li>
            <li>All features will be re-enabled</li>
          </ul>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        <div className="flex gap-3">
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
            onClick={handleSubmit}
            className="flex-1 px-4 py-2 text-white bg-green-600 hover:bg-green-700 rounded-md transition-colors disabled:opacity-50"
            disabled={isLoading}
          >
            {isLoading ? 'Enabling...' : 'Enable Tenant'}
          </button>
        </div>
      </div>
    </div>
  )
}
