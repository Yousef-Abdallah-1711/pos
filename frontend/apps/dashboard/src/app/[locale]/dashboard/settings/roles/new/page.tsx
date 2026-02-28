'use client'

import { useState } from 'react'
import { useAuthStore } from '@/stores/authStore'
import { useHasPermission } from '@/hooks/useRbac'
import { MODULE_IDS, PERMISSION_LEVELS, getPermissionDescription } from '@/lib/rbac'
import { useRouter, usePathname } from 'next/navigation'
import { Loader2, Lock, ArrowLeft } from 'lucide-react'
import Link from 'next/link'
import { debugError } from '@/lib/debug'
import { getConfig } from '@/lib/config'

interface PermissionModule {
  moduleId: number
  moduleName: string
  permissions: {
    level: string
    name: string
    description: string
  }[]
}

const PERMISSION_OPTIONS = [
  {
    level: PERMISSION_LEVELS.READ,
    name: 'Read',
    description: getPermissionDescription(PERMISSION_LEVELS.READ),
  },
  {
    level: PERMISSION_LEVELS.WRITE,
    name: 'Write',
    description: getPermissionDescription(PERMISSION_LEVELS.WRITE),
  },
  {
    level: PERMISSION_LEVELS.DELETE,
    name: 'Delete',
    description: getPermissionDescription(PERMISSION_LEVELS.DELETE),
  },
  {
    level: PERMISSION_LEVELS.ADMIN,
    name: 'Admin',
    description: getPermissionDescription(PERMISSION_LEVELS.ADMIN),
  },
]

const MODULES: PermissionModule[] = [
  {
    moduleId: MODULE_IDS.PRODUCTS,
    moduleName: 'Products',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.CATEGORIES,
    moduleName: 'Categories',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.ORDERS,
    moduleName: 'Orders & POS',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.KDS,
    moduleName: 'Kitchen Display (KDS)',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.TABLES,
    moduleName: 'Tables',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.HR,
    moduleName: 'HR & Payroll',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.LEAVES,
    moduleName: 'Leaves',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.SALARIES,
    moduleName: 'Payroll',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.NOTIFICATIONS,
    moduleName: 'Notifications',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.SETTINGS,
    moduleName: 'Settings',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.USERS,
    moduleName: 'Users',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.ROLES,
    moduleName: 'Roles',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.REPORTS,
    moduleName: 'Reports',
    permissions: PERMISSION_OPTIONS,
  },
  {
    moduleId: MODULE_IDS.THEMES,
    moduleName: 'Themes',
    permissions: PERMISSION_OPTIONS,
  },
]

export default function NewRolePage() {
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const router = useRouter()
  const pathname = usePathname()
  const locale = pathname.split('/')[1] || 'en'

  // Check admin permission
  const { hasPermission: isAdmin, loading: permissionLoading } = useHasPermission(
    MODULE_IDS.SETTINGS,
    PERMISSION_LEVELS.WRITE
  )

  const [formData, setFormData] = useState({
    name: '',
    description: '',
  })

  const [selectedPermissions, setSelectedPermissions] = useState<Record<number, string>>({})
  const [isSaving, setIsSaving] = useState(false)
  const [error, setError] = useState('')

  const handlePermissionChange = (moduleId: number, level: string) => {
    if (selectedPermissions[moduleId] === level) {
      // Deselect if clicking the same level
      const newPerms = { ...selectedPermissions }
      delete newPerms[moduleId]
      setSelectedPermissions(newPerms)
    } else {
      setSelectedPermissions({ ...selectedPermissions, [moduleId]: level })
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!formData.name.trim()) {
      setError('Role name is required')
      return
    }

    if (Object.keys(selectedPermissions).length === 0) {
      setError('At least one permission must be selected')
      return
    }

    if (!token || !user) return

    setIsSaving(true)
    try {
      const roleData = {
        roleName: formData.name,
        description: formData.description,
        permissions: Object.entries(selectedPermissions).map(([moduleId, level]) => ({
          moduleId: Number(moduleId),
          level,
        })),
      }

      const response = await fetch(`${getConfig().apiUrl}/rbac/roles`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(roleData),
      })

      if (response.ok) {
        router.push(`/${locale}/dashboard/settings/roles`)
      } else {
        const err = await response.json()
        setError(err.error || 'Failed to create role')
      }
    } catch (err) {
      setError('Error creating role')
      debugError(err)
    } finally {
      setIsSaving(false)
    }
  }

  if (permissionLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (!isAdmin) {
    return (
      <div className="flex flex-col items-center justify-center h-96">
        <Lock className="w-16 h-16 text-gray-400 mb-4" />
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">Access Denied</h2>
        <p className="text-gray-600 dark:text-gray-400 text-center">
          You don't have permission to create roles. Please contact your administrator.
        </p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Link
          href={`/${locale}/dashboard/settings/roles`}
          className="inline-flex items-center gap-2 px-3 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          Back to Roles
        </Link>
      </div>

      <div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Create New Role</h2>
        <p className="text-gray-600 dark:text-gray-400 text-sm mt-1">
          Define a new role and assign permissions to modules
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Basic Info */}
        <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-6 space-y-4">
          <h3 className="font-semibold text-gray-900 dark:text-white">Basic Information</h3>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Role Name *
            </label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="e.g., Manager, Supervisor"
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Description
            </label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              placeholder="Describe what this role does..."
              rows={3}
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
          </div>
        </div>

        {/* Permissions */}
        <div className="space-y-4">
          <h3 className="font-semibold text-gray-900 dark:text-white">Assign Permissions</h3>

          {MODULES.map((module) => (
            <div key={module.moduleId} className="border border-gray-200 dark:border-gray-800 rounded-lg p-4">
              <h4 className="font-medium text-gray-900 dark:text-white mb-4">{module.moduleName}</h4>

              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                {module.permissions.map((perm) => (
                  <button
                    key={`${module.moduleId}-${perm.level}`}
                    type="button"
                    onClick={() => handlePermissionChange(module.moduleId, perm.level)}
                    className={`p-3 rounded-lg border-2 transition-colors text-left ${
                      selectedPermissions[module.moduleId] === perm.level
                        ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                        : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
                    }`}
                  >
                    <p className="font-medium text-sm text-gray-900 dark:text-white">{perm.name}</p>
                    <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">{perm.description}</p>
                  </button>
                ))}
              </div>
            </div>
          ))}
        </div>

        {error && (
          <div className="p-4 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-lg">
            {error}
          </div>
        )}

        {/* Actions */}
        <div className="flex gap-3 justify-end">
          <Link
            href={`/${locale}/dashboard/settings/roles`}
            className="px-4 py-2 text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-800 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
          >
            Cancel
          </Link>
          <button
            type="submit"
            disabled={isSaving}
            className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 transition-colors flex items-center gap-2"
          >
            {isSaving && <Loader2 className="w-4 h-4 animate-spin" />}
            Create Role
          </button>
        </div>
      </form>
    </div>
  )
}
