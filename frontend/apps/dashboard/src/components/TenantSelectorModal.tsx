'use client'

import React, { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'

interface Tenant {
  tenant_id: number
  tenant_name: string
  user_id: number
  roles: string[]
  default_restaurant_id?: number
}

interface TenantSelectorModalProps {
  isOpen: boolean
  email: string
  /** Proof of password verification returned by POST /auth/login; required by /auth/login/confirm. */
  selectionToken: string
  tenants: Tenant[]
  onTenantSelected?: (tenantId: number) => void
}

export function TenantSelectorModal({
  isOpen,
  email,
  selectionToken,
  tenants,
  onTenantSelected,
}: TenantSelectorModalProps) {
  const router = useRouter()
  const setAuth = useAuthStore((state) => state.setAuth)
  const [selectedTenantId, setSelectedTenantId] = useState<number | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  if (!isOpen || !tenants || tenants.length === 0) {
    return null
  }

  const handleSelectTenant = async (tenantId: number) => {
    setSelectedTenantId(tenantId)
    setIsLoading(true)
    setError(null)

    try {
      const response = await fetch('/api/v1/auth/login/confirm', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email,
          tenant_id: tenantId,
          selection_token: selectionToken,
        }),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to select organization')
      }

      const data = await response.json()

      // Store token and user in the auth store (it persists them)
      if (data.token && data.user) {
        setAuth(data.user, data.token)
      }

      // Call callback if provided
      if (onTenantSelected) {
        onTenantSelected(tenantId)
      }

      // Redirect to dashboard
      router.push('/en/dashboard')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
      setIsLoading(false)
      setSelectedTenantId(null)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center">
      <div className="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-8 w-full max-w-md mx-4">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
            Select Organization
          </h2>
          <p className="text-gray-600 dark:text-gray-400 mt-2">
            You have accounts in multiple organizations. Please select which one you&apos;d like to access.
          </p>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-700 rounded-lg">
            <p className="text-red-700 dark:text-red-400 text-sm">{error}</p>
          </div>
        )}

        <div className="space-y-3">
          {tenants.map((tenant) => (
            <button
              key={tenant.tenant_id}
              onClick={() => handleSelectTenant(tenant.tenant_id)}
              disabled={isLoading}
              className={`w-full text-left p-4 border rounded-lg transition-all ${
                selectedTenantId === tenant.tenant_id
                  ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/30 dark:border-blue-400'
                  : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
              } ${
                isLoading
                  ? 'opacity-50 cursor-not-allowed'
                  : 'cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800'
              }`}
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-900 dark:text-white">
                    {tenant.tenant_name}
                  </h3>
                  <div className="mt-2 flex flex-wrap gap-1">
                    {tenant.roles && tenant.roles.map((role, idx) => (
                      <span
                        key={idx}
                        className="inline-block px-2 py-1 text-xs font-medium bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300 rounded"
                      >
                        {role}
                      </span>
                    ))}
                  </div>
                </div>
                {selectedTenantId === tenant.tenant_id && isLoading && (
                  <div className="ml-2 flex items-center">
                    <div className="animate-spin h-5 w-5 border-2 border-blue-500 border-t-transparent rounded-full"></div>
                  </div>
                )}
              </div>
            </button>
          ))}
        </div>

        <div className="mt-6 text-sm text-gray-600 dark:text-gray-400">
          <p>Logged in as: <strong>{email}</strong></p>
        </div>
      </div>
    </div>
  )
}
