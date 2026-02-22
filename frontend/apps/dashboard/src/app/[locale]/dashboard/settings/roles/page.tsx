'use client'

import { useEffect, useState } from 'react'
import { useAuthStore } from '@/stores/authStore'
import { useHasPermission } from '@/hooks/useRbac'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Loader2, Plus, Edit2, Trash2, Lock } from 'lucide-react'
import { debugError } from '@/lib/debug'
import { getConfig } from '@/lib/config'

interface Role {
  id: number
  name: string
  description: string
  created_at: string
  permissions_count?: number
}

export default function RolesPage() {
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const pathname = usePathname()
  const locale = pathname.split('/')[1] || 'en'

  // Check admin permission
  const { hasPermission: isAdmin, loading: permissionLoading } = useHasPermission(MODULE_IDS.SETTINGS, PERMISSION_LEVELS.WRITE)

  const [roles, setRoles] = useState<Role[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [deleteConfirm, setDeleteConfirm] = useState<number | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  useEffect(() => {
    if (permissionLoading) {
      return
    }
    if (!isAdmin) {
      setIsLoading(false)
      return
    }
    fetchRoles()
  }, [isAdmin, permissionLoading, token])

  const fetchRoles = async () => {
    if (!token || !user) return
    setIsLoading(true)
    const { apiUrl } = getConfig()
    try {
      const response = await fetch(
        `${apiUrl}/rbac/roles`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      )

      if (response.ok) {
        const data = await response.json()
        // Handle both direct array and success response wrapper
        setRoles((data?.data || data) || [])
      } else {
        setError('Failed to fetch roles')
      }
    } catch (err) {
      setError('Error loading roles')
      debugError(err)
    } finally {
      setIsLoading(false)
    }
  }

  const handleDelete = async (roleId: number) => {
    if (!token || !user) return
    setIsDeleting(true)
    try {
      const { apiUrl } = getConfig()
      const response = await fetch(
        `${apiUrl}/rbac/roles/${roleId}`,
        {
          method: 'DELETE',
          headers: { Authorization: `Bearer ${token}` },
        }
      )

      if (response.ok) {
        setRoles(roles.filter((r) => r.id !== roleId))
        setDeleteConfirm(null)
      } else {
        setError('Failed to delete role')
      }
    } catch (err) {
      setError('Error deleting role')
      debugError(err)
    } finally {
      setIsDeleting(false)
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
          You don't have permission to manage roles. Please contact your administrator.
        </p>
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Role Management</h2>
          <p className="text-gray-600 dark:text-gray-400 text-sm mt-1">
            Create and manage roles with custom permissions
          </p>
        </div>
        <Link
          href={`/${locale}/dashboard/settings/roles/new`}
          className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
        >
          <Plus className="w-4 h-4" />
          New Role
        </Link>
      </div>

      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-lg">
          {error}
        </div>
      )}

      {roles.length === 0 ? (
        <div className="flex flex-col items-center justify-center h-64 border-2 border-dashed border-gray-300 dark:border-gray-700 rounded-lg">
          <p className="text-gray-500 dark:text-gray-400">No roles created yet</p>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-800">
                <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Name</th>
                <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Description</th>
                <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Permissions</th>
                <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Created</th>
                <th className="text-right py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Actions</th>
              </tr>
            </thead>
            <tbody>
              {roles.map((role) => (
                <tr key={role.id} className="border-b border-gray-200 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800">
                  <td className="py-3 px-4 text-gray-900 dark:text-white font-medium">{role.name}</td>
                  <td className="py-3 px-4 text-gray-600 dark:text-gray-400">{role.description}</td>
                  <td className="py-3 px-4 text-gray-600 dark:text-gray-400">{role.permissions_count || 0} permissions</td>
                  <td className="py-3 px-4 text-gray-600 dark:text-gray-400 text-sm">
                    {new Date(role.created_at).toLocaleDateString()}
                  </td>
                  <td className="py-3 px-4 text-right flex justify-end gap-2">
                    <Link
                      href={`/${locale}/dashboard/settings/roles/${role.id}/edit`}
                      className="inline-flex items-center gap-2 px-3 py-1 text-sm bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 rounded hover:bg-blue-100 dark:hover:bg-blue-900/40 transition-colors"
                    >
                      <Edit2 className="w-4 h-4" />
                      Edit
                    </Link>
                    <button
                      onClick={() => setDeleteConfirm(role.id)}
                      className="inline-flex items-center gap-2 px-3 py-1 text-sm bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded hover:bg-red-100 dark:hover:bg-red-900/40 transition-colors"
                    >
                      <Trash2 className="w-4 h-4" />
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Delete Confirmation Modal */}
      {deleteConfirm !== null && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-6 max-w-sm">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Delete Role?</h3>
            <p className="text-gray-600 dark:text-gray-400 mb-6">
              Are you sure you want to delete this role? This action cannot be undone.
            </p>
            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setDeleteConfirm(null)}
                className="px-4 py-2 text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-800 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={() => deleteConfirm && handleDelete(deleteConfirm)}
                disabled={isDeleting}
                className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 transition-colors flex items-center gap-2"
              >
                {isDeleting && <Loader2 className="w-4 h-4 animate-spin" />}
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
