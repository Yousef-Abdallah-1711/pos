'use client'

import React, { useState } from 'react'
import { useUserRoles, useUserPermissions } from '@/hooks/useRbac'
import { UserRoleAssignment } from '@/components/rbac/UserRoleAssignment'
import { PermissionBoundary } from '@/components/rbac/PermissionBoundary'
import { MODULE_IDS, PERMISSION_LEVELS, getModuleName, getPermissionDescription } from '@/lib/rbac'
import { AlertCircle, User, Lock, Shield } from 'lucide-react'
import { toast } from 'sonner'

interface RBACSettingsProps {
  userId: number
  tenantId: number
  userName?: string
}

/**
 * RBACSettings Component
 * Manage user roles and view permissions in Settings tab
 */
export const RBACSettings: React.FC<RBACSettingsProps> = ({
  userId,
  tenantId,
  userName,
}) => {
  const { roles, loading } = useUserRoles(userId)
  const { permissions, loading: permissionsLoading } = useUserPermissions()
  const [showRoleAssignment, setShowRoleAssignment] = useState(false)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
          <Shield className="h-5 w-5" />
          Access Control & Permissions
        </h3>
        <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
          Manage your roles and view your permissions across the system
        </p>
      </div>

      {/* Role Assignment Section */}
      <PermissionBoundary
        moduleId={MODULE_IDS.SETTINGS}
        permissionLevel={PERMISSION_LEVELS.ADMIN}
        fallback={
          <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
            <p className="text-sm text-blue-800 dark:text-blue-300">
              Contact your administrator to assign new roles
            </p>
          </div>
        }
      >
        <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6">
          <UserRoleAssignment
            userId={userId}
            tenantId={tenantId}
            userName={userName}
            onRolesChange={() => {
              toast.success('Roles updated successfully')
            }}
          />
        </div>
      </PermissionBoundary>

      {/* Your Roles Section */}
      <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6">
        <h4 className="font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
          <User className="h-5 w-5" />
          Your Roles
        </h4>

        {loading ? (
          <div className="flex items-center justify-center py-6">
            <div className="text-gray-500 dark:text-gray-400">Loading roles...</div>
          </div>
        ) : roles.length > 0 ? (
          <div className="space-y-2">
            {roles.map((role) => (
              <div
                key={role.id}
                className="flex items-start gap-3 p-3 bg-gray-50 dark:bg-gray-700 rounded-lg"
              >
                <div className="flex-shrink-0">
                  <div className="flex items-center justify-center h-8 w-8 rounded-full bg-blue-100 dark:bg-blue-900">
                    <Shield className="h-4 w-4 text-blue-600 dark:text-blue-400" />
                  </div>
                </div>
                <div className="flex-1">
                  <p className="font-medium text-gray-900 dark:text-white">
                    {role.role_name}
                  </p>
                  {role.description && (
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      {role.description}
                    </p>
                  )}
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-6 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg border border-yellow-200 dark:border-yellow-800">
            <AlertCircle className="h-8 w-8 text-yellow-600 dark:text-yellow-400 mx-auto mb-2" />
            <p className="text-gray-700 dark:text-gray-300 font-medium">
              No roles assigned yet
            </p>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
              Your account should have an Admin role. If this is unexpected, the RBAC backend service may need to be started.
            </p>
          </div>
        )}
      </div>

      {/* Your Permissions Summary */}
      <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6">
        <h4 className="font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
          <Lock className="h-5 w-5" />
          Permission Summary
        </h4>

        <div className="space-y-3">
          <div className="text-sm text-gray-600 dark:text-gray-400">
            <p className="mb-3">Your permissions across different modules:</p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {Object.entries(MODULE_IDS).map(([_name, moduleId]) => (
                <PermissionModuleSummary
                  key={moduleId}
                  moduleId={moduleId}
                  permissionLevel={permissions[moduleId]}
                  loading={permissionsLoading}
                />
              ))}
            </div>
          </div>

          {/* Info Box */}
          <div className="mt-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3">
            <p className="text-xs text-blue-800 dark:text-blue-300">
              💡 Your permissions are determined by the roles assigned to you.
              Each role provides access to different modules with specific permission levels.
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}

/**
 * PermissionModuleSummary Component
 * Display permission level for a single module
 */
const PermissionModuleSummary: React.FC<{
  moduleId: number
  permissionLevel?: string
  loading: boolean
}> = ({ moduleId, permissionLevel, loading }) => {
  const level = permissionLevel || PERMISSION_LEVELS.NONE

  return (
    <div className="p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
      <p className="font-medium text-gray-900 dark:text-white text-sm">
        {getModuleName(moduleId)}
      </p>
      <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
        {loading ? 'Loading permissions...' : `${level} - ${getPermissionDescription(level)}`}
      </p>
    </div>
  )
}

export default RBACSettings
