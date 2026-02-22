'use client'

import { useState, useEffect, useMemo } from 'react'
import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { RoleList } from '@/components/hr/RoleList'
import { RoleForm } from '@/components/hr/RoleForm'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { api } from '@/lib/api'
import { debugLog, debugError } from '@/lib/debug'
import { toast } from 'sonner'
import { Plus, AlertCircle } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { PageHeaderSkeleton, TableSkeleton } from '@/components/skeletons'

interface Role {
  id: number
  name: string
  code: string
  access_level: number
  description?: string
  can_manage_users?: boolean
  can_manage_roles?: boolean
  permissions?: string | string[]
  [key: string]: any
}

export default function RolesPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedRole, setSelectedRole] = useState<Role | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [debugInfo, setDebugInfo] = useState<string>('')

  const token = useAuthStore((state) => state.token)

  // Fetch roles
  useEffect(() => {
    const fetchRoles = async () => {
      setLoading(true)
      setApiError(null)
      setDebugInfo(`Fetching roles at ${new Date().toLocaleTimeString()}`)

      try {
        debugLog('[ROLES_PAGE] Starting fetch...')
        const response = await api.get('/hr/roles')
        debugLog('[ROLES_PAGE] Response:', response.data)

        const rolesList = Array.isArray(response.data) ? response.data : (response.data?.data || response.data?.roles || [])
        setRoles(rolesList)
        setDebugInfo(`Successfully loaded ${rolesList.length} roles at ${new Date().toLocaleTimeString()}`)
      } catch (error: any) {
        const errorMsg = error?.response?.data?.message || error?.message || 'Failed to load roles'
        debugError('[ROLES_PAGE] Error:', error)
        setApiError(errorMsg)
        setDebugInfo(`Error: ${errorMsg} at ${new Date().toLocaleTimeString()}`)
        toast.error(errorMsg)
      } finally {
        setLoading(false)
      }
    }

    if (token) {
      fetchRoles()
    }
  }, [token])

  const handleCreateRole = async (formData: any) => {
    try {
      debugLog('[CREATE_ROLE] Submitting:', formData)
      const response = await api.post('/hr/roles', formData)
      debugLog('[CREATE_ROLE] Response:', response.data)

      const created = response.data?.data || response.data
      setRoles([created, ...roles])
      setShowForm(false)
      toast.success(t('hr.roles.createSuccess', 'Role created successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.roles.createFailed', 'Failed to create role')
      debugError('[CREATE_ROLE] Error:', error)
      toast.error(errorMsg)
      throw error
    }
  }

  const handleUpdateRole = async (roleId: number, formData: any) => {
    try {
      debugLog('[UPDATE_ROLE] Submitting:', formData)
      const response = await api.put(`/hr/roles/${roleId}`, formData)
      debugLog('[UPDATE_ROLE] Response:', response.data)

      setRoles(roles.map((r) => (r.id === roleId ? { ...r, ...formData } : r)))
      setShowForm(false)
      setSelectedRole(null)
      setIsEditMode(false)
      toast.success(t('hr.roles.updateSuccess', 'Role updated successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.roles.updateFailed', 'Failed to update role')
      debugError('[UPDATE_ROLE] Error:', error)
      toast.error(errorMsg)
      throw error
    }
  }

  const handleDeleteRole = async (roleId: number) => {
    try {
      debugLog('[DELETE_ROLE] Deleting:', roleId)
      await api.delete(`/hr/roles/${roleId}`)
      debugLog('[DELETE_ROLE] Success')

      setRoles(roles.filter((r) => r.id !== roleId))
      toast.success(t('hr.roles.deleteSuccess', 'Role deleted successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.roles.deleteFailed', 'Failed to delete role')
      debugError('[DELETE_ROLE] Error:', error)
      toast.error(errorMsg)
    }
  }

  const handleEditRole = (role: Role) => {
    setSelectedRole(role)
    setIsEditMode(true)
    setShowForm(true)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">{t('hr.roles.title')}</h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1">
            {t('hr.roles.manageRoles')}
          </p>
        </div>
        <Button onClick={() => {
          setIsEditMode(false)
          setSelectedRole(null)
          setShowForm(true)
        }} className="gap-2 w-full sm:w-auto">
          <Plus className="h-4 w-4" />
          {t('hr.roles.addRole')}
        </Button>
      </div>

      {/* API Error Alert */}
      {apiError && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex gap-3">
          <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <div>
            <h3 className="font-semibold text-red-900 dark:text-red-300">{t('hr.roles.errorLoading')}</h3>
            <p className="text-red-700 dark:text-red-400 text-sm mt-1">{apiError}</p>
            <button
              onClick={() => window.location.reload()}
              className="text-red-600 dark:text-red-400 hover:underline text-sm mt-2 font-medium"
            >
              {t('hr.roles.retry')}
            </button>
          </div>
        </div>
      )}

      {/* Debug Info */}
      {debugInfo && (
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3">
          <p className="text-xs text-blue-600 dark:text-blue-400 font-mono">{debugInfo}</p>
        </div>
      )}

      {/* Role List */}
      {loading && roles.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <RoleList
          roles={roles}
          loading={loading}
          onEdit={handleEditRole}
          onDelete={handleDeleteRole}
        />
      )}

      {/* Role Form Dialog */}
      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {isEditMode ? t('hr.roles.editRole') : t('hr.roles.createNewRole')}
            </DialogTitle>
          </DialogHeader>
          <RoleForm
            role={isEditMode ? selectedRole : null}
            onSubmit={
              isEditMode && selectedRole
                ? (formData: any) => handleUpdateRole(selectedRole.id, formData)
                : handleCreateRole
            }
            onCancel={() => {
              setShowForm(false)
              setSelectedRole(null)
              setIsEditMode(false)
            }}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}
