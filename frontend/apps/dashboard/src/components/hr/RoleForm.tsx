'use client'

import { useState } from 'react'
import { usePathname } from 'next/navigation'
import { useForm, Controller } from 'react-hook-form'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Checkbox } from '@/components/ui/checkbox'
import { toast } from 'sonner'
import { debugLog } from '@/lib/debug'
import { AlertCircle, Loader, Plus, X } from 'lucide-react'

interface RoleFormProps {
  role?: any
  onSubmit: (formData: any) => Promise<void>
  onCancel: () => void
}

const AVAILABLE_PERMISSIONS = [
  { id: 'employees.view', label: 'View Employees' },
  { id: 'employees.create', label: 'Create Employee' },
  { id: 'employees.edit', label: 'Edit Employee' },
  { id: 'employees.delete', label: 'Delete Employee' },
  { id: 'roles.view', label: 'View Roles' },
  { id: 'roles.create', label: 'Create Role' },
  { id: 'roles.edit', label: 'Edit Role' },
  { id: 'roles.delete', label: 'Delete Role' },
  { id: 'attendance.view', label: 'View Attendance' },
  { id: 'attendance.edit', label: 'Edit Attendance' },
  { id: 'leaves.view', label: 'View Leaves' },
  { id: 'leaves.approve', label: 'Approve Leaves' },
  { id: 'salaries.view', label: 'View Salaries' },
  { id: 'salaries.edit', label: 'Edit Salary' },
  { id: 'salaries.process', label: 'Process Payroll' },
]

export function RoleForm({ role, onSubmit, onCancel }: RoleFormProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>(
    role?.permissions ? JSON.parse(typeof role.permissions === 'string' ? role.permissions : JSON.stringify(role.permissions)) : []
  )

  const { control, handleSubmit, formState: { errors }, watch } = useForm({
    defaultValues: {
      role_name: role?.name || role?.role_name || '',
      role_name_ar: role?.name_ar || role?.role_name_ar || '',
      description: role?.description || '',
      description_ar: role?.description_ar || '',
      // Only populate role_code if editing an existing role, leave empty for new roles (backend will auto-generate)
      role_code: role && (role?.code || role?.role_code) ? (role.code || role.role_code) : '',
      access_level: role?.access_level?.toString() || 'basic',
    },
  })

  const togglePermission = (permId: string) => {
    setSelectedPermissions(prev =>
      prev.includes(permId)
        ? prev.filter(p => p !== permId)
        : [...prev, permId]
    )
  }

  const onFormSubmit = async (data: any) => {
    setIsSubmitting(true)
    setApiError(null)

    try {
      // Convert permissions array to JSON object
      const permissionsObj: Record<string, boolean> = {}
      selectedPermissions.forEach(perm => {
        permissionsObj[perm] = true
      })

      const cleanData: any = {
        ...data,
        permissions: permissionsObj,
        access_level: data.access_level.toString(),
      }

      // Only send role_code if it's not empty AND contains valid characters (ASCII only)
      // Let backend auto-generate if empty or contains invalid characters
      if (!cleanData.role_code ||
          cleanData.role_code.trim() === '' ||
          !/^[a-zA-Z0-9_-]+$/.test(cleanData.role_code.trim())) {
        delete cleanData.role_code
      }

      // DEBUG: Log the request data
      debugLog('[CREATE_ROLE] DEBUG - Submitting role data')
      debugLog('Raw form data:', data)
      debugLog('Selected permissions:', selectedPermissions)
      debugLog('Cleaned data being sent:', cleanData)
      debugLog('Timestamp:', new Date().toISOString())
      // groupEnd (debug)

      await onSubmit(cleanData)

      debugLog('[CREATE_ROLE] SUCCESS - Role created successfully')
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || 'Failed to save role'

      // DEBUG: Detailed error logging
      debugLog('[CREATE_ROLE] ERROR - Detailed Debug Info')
      debugLog('Error Type:', error?.constructor?.name || typeof error)
      debugLog('Error Message:', error?.message)
      debugLog('Status Code:', error?.response?.status)
      debugLog('Response Data:', error?.response?.data)
      debugLog('Error Stack:', error?.stack)
      debugLog('Full Error Object:', error)
      // groupEnd (debug)

      setApiError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-6">
      {/* API Error Display */}
      {apiError && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex gap-3">
          <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <div>
            <h3 className="font-semibold text-red-900 dark:text-red-300">Error Saving Role</h3>
            <p className="text-red-700 dark:text-red-400 text-sm mt-1">{apiError}</p>
          </div>
        </div>
      )}

      {/* System Role Warning */}
      {role && ['admin', 'manager', 'employee'].includes(role.code?.toLowerCase()) && (
        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
          <p className="text-yellow-800 dark:text-yellow-300 text-sm">
            This is a system role and cannot be edited or deleted.
          </p>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-2">
          <label className="block text-sm font-medium">Role Name (English)</label>
          <Controller
            name="role_name"
            control={control}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="e.g., Kitchen Manager"
              />
            )}
          />
        </div>

        <div className="space-y-2">
          <label className="block text-sm font-medium">Role Name (Arabic)</label>
          <Controller
            name="role_name_ar"
            control={control}
            render={({ field }) => (
              <Input {...field} placeholder={t('hr.roles.roleNameArPlaceholder')} dir="rtl" />
            )}
          />
        </div>

        <div className="space-y-2">
          <label className="block text-sm font-medium">Role Code (Optional - auto-generated if left empty)</label>
          <Controller
            name="role_code"
            control={control}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="e.g., kitchen_manager (leave empty for auto-generation)"
                onChange={(e) => {
                  // Only allow ASCII characters (a-z, A-Z, 0-9, underscore, hyphen)
                  const value = e.target.value.replace(/[^a-zA-Z0-9_-]/g, '');
                  field.onChange(value);
                }}
              />
            )}
          />
        </div>

        <div className="space-y-2">
          <label className="block text-sm font-medium">Access Level</label>
          <Controller
            name="access_level"
            control={control}
            render={({ field }) => (
              <select
                {...field}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800"
              >
                <option value="basic">Basic</option>
                <option value="supervisor">Supervisor</option>
                <option value="manager">Manager</option>
                <option value="admin">Admin</option>
                <option value="owner">Owner</option>
              </select>
            )}
          />
        </div>
      </div>

      <div className="space-y-2">
        <label className="block text-sm font-medium">Description (English)</label>
        <Controller
          name="description"
          control={control}
          render={({ field }) => (
            <Textarea {...field} placeholder="Describe this role..." rows={3} />
          )}
        />
      </div>

      <div className="space-y-2">
        <label className="block text-sm font-medium">Description (Arabic)</label>
        <Controller
          name="description_ar"
          control={control}
          render={({ field }) => (
            <Textarea {...field} placeholder={t('hr.roles.descriptionArPlaceholder')} rows={3} dir="rtl" />
          )}
        />
      </div>

      <div className="border-t pt-6">
        <h3 className="font-semibold mb-4">Permissions</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
          {AVAILABLE_PERMISSIONS.map(permission => (
            <div key={permission.id} className="flex items-center space-x-2">
              <Checkbox
                checked={selectedPermissions.includes(permission.id)}
                onCheckedChange={() => togglePermission(permission.id)}
              />
              <label className="text-sm">{permission.label}</label>
            </div>
          ))}
        </div>
      </div>

      {/* Form Actions */}
      <div className="flex justify-end gap-2 pt-4 border-t">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting} className="gap-2">
          {isSubmitting && <Loader className="h-4 w-4 animate-spin" />}
          {isSubmitting ? 'Saving...' : role ? 'Update Role' : 'Create Role'}
        </Button>
      </div>
    </form>
  )
}
