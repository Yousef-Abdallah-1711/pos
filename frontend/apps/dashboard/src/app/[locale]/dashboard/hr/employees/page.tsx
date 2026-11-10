'use client'

import { useState, useEffect, useMemo } from 'react'
import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { useHasPermission } from '@/hooks/useRbac'
import { EmployeeList } from '@/components/hr/EmployeeList'
import { EmployeeForm } from '@/components/hr/EmployeeForm'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { api } from '@/lib/api'
import { debugLog, debugError } from '@/lib/debug'
import { toast } from 'sonner'
import { Plus, AlertCircle, Lock, Loader2 } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { PageHeaderSkeleton, TableSkeleton } from '@/components/skeletons'
import { hasAuthenticatedSession } from '@/lib/authSession'

interface Employee {
  id: number
  first_name: string
  last_name: string
  email: string
  phone?: string
  employee_code: string
  position: string
  department?: string
  hire_date?: string
  status: string
  employment_type?: string
  monthly_salary?: number
  weekly_hours?: number
  is_active: boolean
  [key: string]: any
}

export default function EmployeesPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  // Check RBAC permissions for HR module - properly destructure the hook return
  const { hasPermission: canViewEmployees, loading: canViewLoading } = useHasPermission(
    MODULE_IDS.HR,
    PERMISSION_LEVELS.READ,
  )
  const { hasPermission: canCreateEmployees, loading: canCreateLoading } = useHasPermission(
    MODULE_IDS.HR,
    PERMISSION_LEVELS.WRITE,
  )
  const { hasPermission: canDeleteEmployees, loading: canDeleteLoading } = useHasPermission(
    MODULE_IDS.HR,
    PERMISSION_LEVELS.DELETE,
  )
  const permissionLoading = canViewLoading || canCreateLoading || canDeleteLoading

  const [employees, setEmployees] = useState<Employee[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedEmployee, setSelectedEmployee] = useState<Employee | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [debugInfo, setDebugInfo] = useState<string>('')

  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const hasSession = hasAuthenticatedSession(user, token)

  // Fetch employees - hook MUST be before early returns
  useEffect(() => {
    if (!hasSession || !canViewEmployees || permissionLoading) return

    const fetchEmployees = async () => {
      setLoading(true)
      setApiError(null)
      setDebugInfo(`Fetching employees at ${new Date().toLocaleTimeString()}`)

      try {
        debugLog('[EMPLOYEES_PAGE] Starting fetch...')
        const response = await api.get('/hr/employees')
        debugLog('[EMPLOYEES_PAGE] Response:', response.data)

        const employeesList = Array.isArray(response.data) ? response.data : (response.data?.data || response.data?.employees || [])
        setEmployees(employeesList)
        setDebugInfo(`Successfully loaded ${employeesList.length} employees at ${new Date().toLocaleTimeString()}`)
      } catch (error: any) {
        const errorMsg = error?.response?.data?.message || error?.message || 'Failed to load employees'
        debugError('[EMPLOYEES_PAGE] Error:', error)
        setApiError(errorMsg)
        setDebugInfo(`Error: ${errorMsg} at ${new Date().toLocaleTimeString()}`)
        toast.error(errorMsg)
      } finally {
        setLoading(false)
      }
    }

    fetchEmployees()
  }, [hasSession, canViewEmployees, permissionLoading])

  const handleCreateEmployee = async (formData: any) => {
    try {
      debugLog('[CREATE_EMPLOYEE] Submitting:', formData)
      const response = await api.post('/hr/employees', formData)
      debugLog('[CREATE_EMPLOYEE] Response:', response.data)

      const created = response.data?.data || response.data
      setEmployees([created, ...employees])
      setShowForm(false)
      toast.success('Employee created successfully')
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || 'Failed to create employee'
      debugError('[CREATE_EMPLOYEE] Error:', error)
      toast.error(errorMsg)
      throw error
    }
  }

  const handleUpdateEmployee = async (employeeId: number, formData: any) => {
    try {
      debugLog('[UPDATE_EMPLOYEE] Submitting:', formData)
      const response = await api.put(`/hr/employees/${employeeId}`, formData)
      debugLog('[UPDATE_EMPLOYEE] Response:', response.data)

      setEmployees(employees.map((e) => (e.id === employeeId ? { ...e, ...formData } : e)))
      setShowForm(false)
      setSelectedEmployee(null)
      setIsEditMode(false)
      toast.success('Employee updated successfully')
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || 'Failed to update employee'
      debugError('[UPDATE_EMPLOYEE] Error:', error)
      toast.error(errorMsg)
      throw error
    }
  }

  const handleDeleteEmployee = async (employeeId: number) => {
    try {
      debugLog('[DELETE_EMPLOYEE] Deleting:', employeeId)
      await api.delete(`/hr/employees/${employeeId}`)
      debugLog('[DELETE_EMPLOYEE] Success')

      setEmployees(employees.filter((e) => e.id !== employeeId))
      toast.success('Employee deleted successfully')
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || 'Failed to delete employee'
      debugError('[DELETE_EMPLOYEE] Error:', error)
      toast.error(errorMsg)
    }
  }

  const handleEditEmployee = (employee: Employee) => {
    setSelectedEmployee(employee)
    setIsEditMode(true)
    setShowForm(true)
  }

  // Show loading while permissions are being checked
  if (permissionLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  // Show permission denied if user can't view employees
  if (!canViewEmployees) {
    return (
      <div className="space-y-6">
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
          <div className="flex gap-3 items-start">
            <Lock className="h-5 w-5 text-red-600 dark:text-red-400 mt-0.5 flex-shrink-0" />
            <div>
              <h2 className="font-semibold text-red-900 dark:text-red-300">
                Access Denied
              </h2>
              <p className="text-sm text-red-800 dark:text-red-400 mt-1">
                You do not have permission to view employees. Please contact your administrator.
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {loading && employees.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <>
          {/* Header */}
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <div>
              <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">{t('hr.employees.title')}</h1>
              <p className="text-gray-500 dark:text-gray-400 mt-1">
                {t('hr.employees.manageEmployees')}
              </p>
            </div>
            {canCreateEmployees ? (
              <Button onClick={() => {
                setIsEditMode(false)
                setSelectedEmployee(null)
                setShowForm(true)
              }} className="gap-2 w-full sm:w-auto">
                <Plus className="h-4 w-4" />
                {t('hr.employees.addEmployee')}
              </Button>
            ) : (
              <Button disabled className="gap-2 w-full sm:w-auto">
                <Plus className="h-4 w-4" />
                {t('hr.employees.addEmployee')}
              </Button>
            )}
          </div>

          {/* API Error Alert */}
          {apiError && (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex gap-3">
              <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
              <div>
                <h3 className="font-semibold text-red-900 dark:text-red-300">{t('hr.employees.errorLoading')}</h3>
                <p className="text-red-700 dark:text-red-400 text-sm mt-1">{apiError}</p>
                <button
                  onClick={() => window.location.reload()}
                  className="text-red-600 dark:text-red-400 hover:underline text-sm mt-2 font-medium"
                >
                  {t('hr.employees.retry')}
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

          {/* Employee List */}
          <EmployeeList
            employees={employees}
            loading={loading}
            onEdit={handleEditEmployee}
            onDelete={canDeleteEmployees ? handleDeleteEmployee : undefined}
            canDelete={canDeleteEmployees}
          />
        </>
      )}

      {/* Employee Form Dialog */}
      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {isEditMode ? t('hr.employees.editEmployee') : t('hr.employees.createNewEmployee')}
            </DialogTitle>
          </DialogHeader>
          <EmployeeForm
            employee={isEditMode ? selectedEmployee : null}
            onSubmit={
              isEditMode && selectedEmployee
                ? (formData: any) => handleUpdateEmployee(selectedEmployee.id, formData)
                : handleCreateEmployee
            }
            onCancel={() => {
              setShowForm(false)
              setSelectedEmployee(null)
              setIsEditMode(false)
            }}
            existingEmails={employees.map((emp) => emp.email)}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}
