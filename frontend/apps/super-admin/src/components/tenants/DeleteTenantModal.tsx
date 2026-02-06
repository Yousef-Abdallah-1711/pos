'use client'

import { useState } from 'react'

interface DeleteTenantModalProps {
  tenantId: number
  tenantName: string
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
}

export function DeleteTenantModal({
  tenantId,
  tenantName,
  isOpen,
  onClose,
  onSuccess,
}: DeleteTenantModalProps) {
  const [confirmText, setConfirmText] = useState<string>('')
  const [reason, setReason] = useState<string>('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  if (!isOpen) return null

  const isConfirmed = confirmText === tenantName

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!isConfirmed) {
      setError('Confirmation text does not match tenant name')
      return
    }

    setIsLoading(true)

    try {
      const response = await fetch(
        `/api/v1/super-admin/tenants/${tenantId}?confirm=true`,
        {
          method: 'DELETE',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            reason: reason || undefined,
          }),
        }
      )

      if (!response.ok) {
        throw new Error('Failed to delete tenant')
      }

      onSuccess()
      onClose()
      setConfirmText('')
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
        <h2 className="text-xl font-bold mb-2 text-red-700">Delete Tenant Permanently</h2>
        <p className="text-gray-600 mb-4">
          Tenant: <span className="font-semibold text-red-600">{tenantName}</span>
        </p>

        <div className="mb-4 p-3 bg-red-50 border-2 border-red-300 rounded text-sm text-red-700">
          <p className="font-bold mb-2">🚨 CRITICAL WARNING:</p>
          <ul className="list-disc list-inside space-y-1">
            <li>This action is IRREVERSIBLE</li>
            <li>All tenant data will be permanently deleted</li>
            <li>Backups should be created before proceeding</li>
            <li>Users will lose all access immediately</li>
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
              Type tenant name to confirm: <span className="font-bold">{tenantName}</span> *
            </label>
            <input
              type="text"
              value={confirmText}
              onChange={(e) => setConfirmText(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-red-500"
              placeholder={tenantName}
              required
            />
            {confirmText && (
              <p
                className={`mt-2 text-sm ${
                  isConfirmed ? 'text-green-600' : 'text-red-600'
                }`}
              >
                {isConfirmed ? '✓ Confirmed' : '✗ Does not match'}
              </p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Reason for deletion (optional)
            </label>
            <textarea
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-red-500"
              rows={3}
              placeholder="e.g., Customer requested account closure"
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
              className="flex-1 px-4 py-2 text-white bg-red-700 hover:bg-red-800 rounded-md transition-colors disabled:opacity-50"
              disabled={isLoading || !isConfirmed}
            >
              {isLoading ? 'Deleting...' : 'Delete Permanently'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
