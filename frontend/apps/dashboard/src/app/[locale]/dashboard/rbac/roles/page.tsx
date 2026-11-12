'use client'

import { useState, useEffect } from 'react'
import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { hasAuthenticatedSession } from '@/lib/authSession'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { useRoles } from '@/hooks/useRbac'
import { PermissionMatrix } from '@/components/rbac/PermissionMatrix'
import { Plus, AlertCircle, Edit2, Trash2, Loader2, ChevronDown } from 'lucide-react'
import { toast } from 'sonner'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { PageHeaderSkeleton, TableSkeleton } from '@/components/skeletons'

interface RoleWithPermissions {
  id: number
  role_name: string
  description?: string
  is_system: boolean
  created_at?: string
  updated_at?: string
  expanded?: boolean
}

export default function RBACRolesPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  const { roles, loading, error, fetchRoles, createRole, updateRole, deleteRole } = useRoles()
  const [displayRoles, setDisplayRoles] = useState<RoleWithPermissions[]>([])
  const [showForm, setShowForm] = useState(false)
  const [selectedRole, setSelectedRole] = useState<RoleWithPermissions | null>(null)
  const [isEditMode, setIsEditMode] = useState(false)
  const [formData, setFormData] = useState({ role_name: '', description: '' })
  const [submitting, setSubmitting] = useState(false)
  const [expandedRole, setExpandedRole] = useState<number | null>(null)

  const token = useAuthStore((state) => state.token)
  const user = useAuthStore((state) => state.user)
  const hasSession = hasAuthenticatedSession(user, token)
  const tenantId = user?.tenant_id || user?.tenantId || null

  // Fetch roles on mount
  useEffect(() => {
    if (hasSession && tenantId) {
      fetchRoles(tenantId)
    }
  }, [fetchRoles, hasSession, tenantId])

  // Update display roles when roles change
  useEffect(() => {
    setDisplayRoles(roles.map((r) => ({ ...r, expanded: false })))
  }, [roles])

  // Handle role creation
  const handleCreateRole = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.role_name) {
      toast.error('Role name is required')
      return
    }

    try {
      setSubmitting(true)
      await createRole({
        role_name: formData.role_name,
        description: formData.description,
      })
      toast.success('Role created successfully')
      setShowForm(false)
      setFormData({ role_name: '', description: '' })
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to create role')
    } finally {
      setSubmitting(false)
    }
  }

  // Handle role update
  const handleUpdateRole = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedRole || !formData.role_name) {
      toast.error('Role name is required')
      return
    }

    try {
      setSubmitting(true)
      await updateRole(selectedRole.id, {
        role_name: formData.role_name,
        description: formData.description,
      })
      toast.success('Role updated successfully')
      setShowForm(false)
      setSelectedRole(null)
      setFormData({ role_name: '', description: '' })
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to update role')
    } finally {
      setSubmitting(false)
    }
  }

  // Handle role deletion
  const handleDeleteRole = async (role: RoleWithPermissions) => {
    if (role.is_system) {
      toast.error('Cannot delete system roles')
      return
    }

    if (!confirm(`Are you sure you want to delete "${role.role_name}"?`)) return

    try {
      await deleteRole(role.id)
      toast.success('Role deleted successfully')
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to delete role')
    }
  }

  // Handle edit role
  const handleEditRole = (role: RoleWithPermissions) => {
    setSelectedRole(role)
    setFormData({ role_name: role.role_name, description: role.description || '' })
    setIsEditMode(true)
    setShowForm(true)
  }

  // Handle new role
  const handleNewRole = () => {
    setSelectedRole(null)
    setFormData({ role_name: '', description: '' })
    setIsEditMode(false)
    setShowForm(true)
  }

  // Toggle role details expansion
  const toggleExpanded = (roleId: number) => {
    setExpandedRole(expandedRole === roleId ? null : roleId)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">
            RBAC Role Management
          </h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1">
            Create and manage roles with granular permissions
          </p>
        </div>
        <Button onClick={handleNewRole} className="gap-2 w-full sm:w-auto">
          <Plus className="h-4 w-4" />
          Create Role
        </Button>
      </div>

      {/* Error Alert */}
      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex gap-3">
          <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <div>
            <h3 className="font-semibold text-red-900 dark:text-red-300">
              Error Loading Roles
            </h3>
            <p className="text-red-700 dark:text-red-400 text-sm mt-1">
              {error.message}
            </p>
          </div>
        </div>
      )}

      {/* Roles List */}
      {loading && displayRoles.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : displayRoles.length > 0 ? (
        <div className="space-y-3">
          {displayRoles.map((role) => (
            <div
              key={role.id}
              className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
            >
              {/* Role Header */}
              <div
                className="p-4 flex items-center justify-between cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700/50"
                onClick={() => toggleExpanded(role.id)}
              >
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-900 dark:text-white">
                    {role.role_name}
                    {role.is_system && (
                      <span className="ml-2 inline-block px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-xs rounded-full">
                        System
                      </span>
                    )}
                  </h3>
                  {role.description && (
                    <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                      {role.description}
                    </p>
                  )}
                </div>
                <div className="flex items-center gap-2">
                  {!role.is_system && (
                    <>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={(e) => {
                          e.stopPropagation()
                          handleEditRole(role)
                        }}
                        className="gap-1"
                      >
                        <Edit2 className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={(e) => {
                          e.stopPropagation()
                          handleDeleteRole(role)
                        }}
                        className="gap-1 text-red-600"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </>
                  )}
                  <ChevronDown
                    className={`h-5 w-5 text-gray-400 transition-transform ${
                      expandedRole === role.id ? 'rotate-180' : ''
                    }`}
                  />
                </div>
              </div>

              {/* Role Details - Permissions */}
              {expandedRole === role.id && (
                <div className="border-t border-gray-200 dark:border-gray-700 p-4 bg-gray-50 dark:bg-gray-900/50">
                  <div className="space-y-4">
                    <div>
                      <h4 className="font-semibold text-gray-900 dark:text-white mb-3">
                        Permissions
                      </h4>
                      <PermissionMatrix
                        roleId={role.id}
                        readOnly={role.is_system}
                      />
                    </div>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-12 bg-gray-50 dark:bg-gray-800 rounded-lg">
          <p className="text-gray-500 dark:text-gray-400">No roles found</p>
          <Button onClick={handleNewRole} variant="outline" className="mt-4">
            Create your first role
          </Button>
        </div>
      )}

      {/* Role Form Dialog */}
      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>
              {isEditMode ? 'Edit Role' : 'Create New Role'}
            </DialogTitle>
          </DialogHeader>

          <form
            onSubmit={isEditMode ? handleUpdateRole : handleCreateRole}
            className="space-y-4"
          >
            {/* Role Name */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Role Name *
              </label>
              <input
                type="text"
                value={formData.role_name}
                onChange={(e) =>
                  setFormData({ ...formData, role_name: e.target.value })
                }
                disabled={submitting}
                placeholder="e.g., Manager, Editor, Viewer"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-500 dark:placeholder-gray-400 disabled:opacity-50"
              />
            </div>

            {/* Description */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Description
              </label>
              <textarea
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
                disabled={submitting}
                placeholder="Describe this role's purpose..."
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-500 dark:placeholder-gray-400 disabled:opacity-50 resize-none"
              />
            </div>

            {/* Actions */}
            <div className="flex gap-3 pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => setShowForm(false)}
                disabled={submitting}
                className="flex-1"
              >
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={submitting}
                className="flex-1"
              >
                {submitting ? (
                  <>
                    <Loader2 className="h-4 w-4 animate-spin mr-2" />
                    {isEditMode ? 'Updating...' : 'Creating...'}
                  </>
                ) : isEditMode ? (
                  'Update'
                ) : (
                  'Create'
                )}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
