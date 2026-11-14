'use client'

import React, { ReactNode } from 'react'
import { useHasPermission } from '@/hooks/useRbac'
import { AlertCircle } from 'lucide-react'
import { PERMISSION_LEVELS } from '@/lib/rbac'

interface RoleGuardProps {
  moduleId: number
  permissionLevel?: string
  children: ReactNode
  fallback?: ReactNode
  showError?: boolean
}

/**
 * RoleGuard Component
 * Wraps components/pages that require specific RBAC permissions
 * Only renders children if user has the required permission
 */
export const RoleGuard: React.FC<RoleGuardProps> = ({
  moduleId,
  permissionLevel = PERMISSION_LEVELS.READ,
  children,
  fallback,
  showError = true,
}) => {
  const { hasPermission, loading } = useHasPermission(moduleId, permissionLevel)

  if (loading) {
    return (
      <div className="flex items-center justify-center py-8">
        <div className="text-gray-500">Checking permissions...</div>
      </div>
    )
  }

  if (!hasPermission) {
    if (fallback) {
      return <>{fallback}</>
    }

    if (showError) {
      return (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6 flex gap-3">
          <AlertCircle className="h-6 w-6 text-red-600 dark:text-red-400 flex-shrink-0" />
          <div>
            <h3 className="font-semibold text-red-900 dark:text-red-300">
              Access Denied
            </h3>
            <p className="text-red-700 dark:text-red-400 text-sm mt-1">
              You don&apos;t have the required permissions to access this content.
              {permissionLevel &&
                permissionLevel !== PERMISSION_LEVELS.READ &&
                ` (Required: ${permissionLevel})`}
            </p>
          </div>
        </div>
      )
    }

    return null
  }

  return <>{children}</>
}

export default RoleGuard
