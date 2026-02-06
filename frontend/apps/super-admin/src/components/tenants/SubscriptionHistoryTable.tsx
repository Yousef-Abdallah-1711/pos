'use client'

import { useState, useEffect } from 'react'
import type { SubscriptionHistoryItem } from '@/types'

interface SubscriptionHistoryTableProps {
  tenantId: number
  className?: string
}

export function SubscriptionHistoryTable({
  tenantId,
  className = '',
}: SubscriptionHistoryTableProps) {
  const [history, setHistory] = useState<SubscriptionHistoryItem[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const pageSize = 10

  useEffect(() => {
    fetchHistory()
  }, [tenantId, page])

  const fetchHistory = async () => {
    setIsLoading(true)
    setError(null)

    try {
      const response = await fetch(
        `/api/v1/super-admin/tenants/${tenantId}/subscription-history?page=${page}&page_size=${pageSize}`
      )

      if (!response.ok) {
        throw new Error('Failed to fetch history')
      }

      const data = await response.json()
      setHistory(data.data?.items || [])
      setTotal(data.data?.total || 0)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  const getActionIcon = (action: string): string => {
    switch (action) {
      case 'TRIAL_STARTED':
        return '🎬'
      case 'TRIAL_EXTENDED':
        return '📅'
      case 'TRIAL_EXPIRED':
        return '⏰'
      case 'TENANT_DISABLED':
        return '🔒'
      case 'TENANT_ENABLED':
        return '✓'
      case 'DELETION_SCHEDULED':
        return '📋'
      case 'DELETION_CANCELLED':
        return '❌'
      case 'TENANT_DELETED':
        return '🗑️'
      case 'STATUS_CHANGED':
        return '🔄'
      default:
        return '•'
    }
  }

  const getActionLabel = (action: string): string => {
    return action
      .split('_')
      .map((word) => word.charAt(0) + word.slice(1).toLowerCase())
      .join(' ')
  }

  const totalPages = Math.ceil(total / pageSize)

  if (isLoading && history.length === 0) {
    return (
      <div className={`text-center py-8 ${className}`}>
        <p className="text-gray-500">Loading history...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className={`text-center py-8 ${className}`}>
        <p className="text-red-500">Error: {error}</p>
      </div>
    )
  }

  if (history.length === 0) {
    return (
      <div className={`text-center py-8 ${className}`}>
        <p className="text-gray-500">No history available</p>
      </div>
    )
  }

  return (
    <div className={className}>
      <div className="overflow-x-auto">
        <table className="w-full border-collapse">
          <thead>
            <tr className="border-b-2 border-gray-200">
              <th className="text-left px-4 py-3 font-semibold text-gray-700">
                Date & Time
              </th>
              <th className="text-left px-4 py-3 font-semibold text-gray-700">
                Action
              </th>
              <th className="text-left px-4 py-3 font-semibold text-gray-700">
                Details
              </th>
              <th className="text-left px-4 py-3 font-semibold text-gray-700">
                Reason
              </th>
            </tr>
          </thead>
          <tbody>
            {history.map((item) => (
              <tr key={item.id} className="border-b border-gray-100 hover:bg-gray-50">
                <td className="px-4 py-3 text-sm text-gray-600">
                  {new Date(item.created_at).toLocaleString()}
                </td>
                <td className="px-4 py-3 text-sm font-medium">
                  <span className="mr-2">{getActionIcon(item.action)}</span>
                  {getActionLabel(item.action)}
                </td>
                <td className="px-4 py-3 text-sm text-gray-600">
                  {item.old_status && item.new_status && (
                    <span>
                      {item.old_status} → {item.new_status}
                    </span>
                  )}
                  {item.days_added && <span>{item.days_added} days added</span>}
                  {!item.old_status && !item.new_status && !item.days_added && (
                    <span className="text-gray-400">—</span>
                  )}
                </td>
                <td className="px-4 py-3 text-sm text-gray-600">
                  {item.reason || <span className="text-gray-400">—</span>}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-between mt-4">
          <div className="text-sm text-gray-600">
            Showing {(page - 1) * pageSize + 1} to{' '}
            {Math.min(page * pageSize, total)} of {total} records
          </div>

          <div className="flex gap-2">
            <button
              onClick={() => setPage(Math.max(1, page - 1))}
              disabled={page === 1}
              className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Previous
            </button>

            <div className="flex items-center gap-2">
              {Array.from({ length: totalPages }, (_, i) => i + 1)
                .slice(
                  Math.max(0, page - 2),
                  Math.min(totalPages, page + 1)
                )
                .map((pageNum) => (
                  <button
                    key={pageNum}
                    onClick={() => setPage(pageNum)}
                    className={`px-2 py-1 text-sm border rounded ${
                      pageNum === page
                        ? 'bg-blue-500 text-white border-blue-500'
                        : 'border-gray-300 hover:bg-gray-50'
                    }`}
                  >
                    {pageNum}
                  </button>
                ))}
            </div>

            <button
              onClick={() => setPage(Math.min(totalPages, page + 1))}
              disabled={page === totalPages}
              className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
