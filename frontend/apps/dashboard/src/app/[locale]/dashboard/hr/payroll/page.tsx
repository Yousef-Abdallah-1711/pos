'use client'

import { useState, useEffect, useMemo } from 'react'
import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { useUserPermissions } from '@/hooks/useRbac'
import { SalaryList } from '@/components/hr/SalaryList'
import { SalaryForm } from '@/components/hr/SalaryForm'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { api } from '@/lib/api'
import { debugLog, debugError } from '@/lib/debug'
import { toast } from 'sonner'
import { Plus, AlertCircle, Lock, Loader2 } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { PageHeaderSkeleton, TableSkeleton } from '@/components/skeletons'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'

interface Salary {
  id: number
  employee_id: number
  employee_name?: string
  base_salary: number
  status: string
  pay_period_start: string
  pay_period_end: string
  [key: string]: any
}

interface Employee {
  id: number
  first_name: string
  last_name: string
}

export default function PayrollPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  // Check RBAC permissions for Salaries module
  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canViewSalaries = checkPermission(MODULE_IDS.SALARIES, PERMISSION_LEVELS.READ)
  const canCreateSalaries = checkPermission(MODULE_IDS.SALARIES, PERMISSION_LEVELS.WRITE)

  // ALL hooks MUST be before early returns to avoid Hook order violation
  const [salaries, setSalaries] = useState<Salary[]>([])
  const [employees, setEmployees] = useState<Employee[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedSalary, setSelectedSalary] = useState<Salary | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [debugInfo, setDebugInfo] = useState<string>('')
  const [startDate, setStartDate] = useState(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0])
  const [endDate, setEndDate] = useState(new Date().toISOString().split('T')[0])

  const token = useAuthStore((state) => state.token)

  // Fetch employees - hook MUST be before early returns
  useEffect(() => {
    if (!token || !canViewSalaries || permissionLoading) return

    const fetchEmployees = async () => {
      try {
        const response = await api.get('/hr/employees')
        const empList = Array.isArray(response.data) ? response.data : (response.data?.data || response.data?.employees || [])
        setEmployees(empList)
      } catch (error) {
        debugError('Failed to load employees:', error)
      }
    }
    fetchEmployees()
  }, [token, canViewSalaries, permissionLoading])

  // Fetch salaries - hook MUST be before early returns
  useEffect(() => {
    if (!token || !canViewSalaries || permissionLoading) return

    const fetchSalaries = async () => {
      setLoading(true)
      setApiError(null)
      setDebugInfo(`Fetching salaries at ${new Date().toLocaleTimeString()}`)

      try {
        debugLog('[PAYROLL_PAGE] Starting fetch...')
        const response = await api.get('/hr/salaries', { params: { start_date: startDate, end_date: endDate } })
        debugLog('[PAYROLL_PAGE] Response:', response.data)

        const salariesList = Array.isArray(response.data) ? response.data : (response.data?.data || response.data?.salaries || [])
        const enrichedSalaries = salariesList.map(s => ({
          ...s,
          employee_name: employees.find(e => e.id === s.employee_id)
            ? `${employees.find(e => e.id === s.employee_id)?.first_name} ${employees.find(e => e.id === s.employee_id)?.last_name}`
            : undefined,
        }))

        setSalaries(enrichedSalaries)
        setDebugInfo(`Successfully loaded ${enrichedSalaries.length} salary records`)
      } catch (error: any) {
        const errorMsg = error?.response?.data?.message || error?.message || 'Failed to load salaries'
        debugError('[PAYROLL_PAGE] Error:', error)
        setApiError(errorMsg)
        setDebugInfo(`Error: ${errorMsg}`)
        toast.error(errorMsg)
      } finally {
        setLoading(false)
      }
    }

    fetchSalaries()
  }, [token, startDate, endDate, employees, canViewSalaries, permissionLoading])

  // Show permission denied if user can't view salaries
  if (permissionLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (!canViewSalaries) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            {t('hr.payroll.title')}
          </h1>
        </div>

        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
          <div className="flex gap-3 items-start">
            <Lock className="h-5 w-5 text-red-600 dark:text-red-400 mt-0.5 flex-shrink-0" />
            <div>
              <h2 className="font-semibold text-red-900 dark:text-red-300">
                Access Denied
              </h2>
              <p className="text-sm text-red-800 dark:text-red-400 mt-1">
                You do not have permission to view salary information. Please contact your administrator.
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  const handleCreateSalary = async (formData: any) => {
    try {
      debugLog('[CREATE_SALARY] Submitting:', formData)
      const response = await api.post('/hr/salaries', formData)
      debugLog('[CREATE_SALARY] Response:', response.data)
      const created = response.data?.data || response.data
      setSalaries([created, ...salaries])
      setShowForm(false)
      toast.success(t('hr.payroll.successCreated', 'Salary record created successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.payroll.errorCreated', 'Failed to create salary')
      debugError('[CREATE_SALARY] Error:', error)
      toast.error(errorMsg)
      throw error
    }
  }

  const handleUpdateSalary = async (salaryId: number, formData: any) => {
    try {
      debugLog('[UPDATE_SALARY] Submitting:', formData)
      const response = await api.put(`/hr/salaries/${salaryId}`, formData)
      debugLog('[UPDATE_SALARY] Response:', response.data)
      setSalaries(salaries.map(s => s.id === salaryId ? { ...s, ...formData } : s))
      setShowForm(false)
      setSelectedSalary(null)
      setIsEditMode(false)
      toast.success(t('hr.payroll.successUpdated', 'Salary record updated successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.payroll.errorUpdated', 'Failed to update salary')
      debugError('[UPDATE_SALARY] Error:', error)
      toast.error(errorMsg)
      throw error
    }
  }

  const handleDeleteSalary = async (salaryId: number) => {
    try {
      debugLog('[DELETE_SALARY] Deleting:', salaryId)
      await api.delete(`/hr/salaries/${salaryId}`)
      setSalaries(salaries.filter(s => s.id !== salaryId))
      toast.success(t('hr.payroll.successDeleted', 'Salary record deleted successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.payroll.errorDeleted', 'Failed to delete salary')
      debugError('[DELETE_SALARY] Error:', error)
      toast.error(errorMsg)
    }
  }

  const handleMarkPaid = async (salaryId: number) => {
    try {
      debugLog('[MARK_PAID] Marking salary as paid:', salaryId)
      await api.post(`/hr/salaries/${salaryId}/mark-paid`, { payment_reference: `PAY-${Date.now()}` })
      setSalaries(salaries.map(s => s.id === salaryId ? { ...s, status: 'paid' } : s))
      toast.success(t('hr.payroll.successMarkedPaid', 'Salary marked as paid'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.payroll.errorMarkedPaid', 'Failed to mark salary as paid')
      debugError('[MARK_PAID] Error:', error)
      toast.error(errorMsg)
    }
  }

  const handleEditSalary = (salary: Salary) => {
    setSelectedSalary(salary)
    setIsEditMode(true)
    setShowForm(true)
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">{t('hr.payroll.title')}</h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1">{t('hr.payroll.managePayroll')}</p>
        </div>
        {canCreateSalaries ? (
          <Button onClick={() => { setIsEditMode(false); setSelectedSalary(null); setShowForm(true) }} className="gap-2 w-full sm:w-auto">
            <Plus className="h-4 w-4" />
            {t('hr.payroll.createSalary')}
          </Button>
        ) : (
          <Button disabled className="gap-2 w-full sm:w-auto opacity-50 cursor-not-allowed" title="You don't have permission to create salaries">
            <Plus className="h-4 w-4" />
            {t('hr.payroll.createSalary')}
          </Button>
        )}
      </div>

      <div className="flex flex-col md:flex-row gap-4 bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
        <div className="flex-1 space-y-2">
          <label className="block text-sm font-medium">{t('hr.payroll.fromDate')}</label>
          <Input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} />
        </div>
        <div className="flex-1 space-y-2">
          <label className="block text-sm font-medium">{t('hr.payroll.toDate')}</label>
          <Input type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} />
        </div>
      </div>

      {apiError && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex gap-3">
          <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <div>
            <h3 className="font-semibold text-red-900 dark:text-red-300">{t('hr.payroll.errorLoading')}</h3>
            <p className="text-red-700 dark:text-red-400 text-sm mt-1">{apiError}</p>
            <button onClick={() => window.location.reload()} className="text-red-600 dark:text-red-400 hover:underline text-sm mt-2 font-medium">{t('hr.payroll.retry')}</button>
          </div>
        </div>
      )}

      {debugInfo && (
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3">
          <p className="text-xs text-blue-600 dark:text-blue-400 font-mono">{debugInfo}</p>
        </div>
      )}

      {loading && salaries.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <SalaryList
          salaries={salaries}
          loading={loading}
          onEdit={handleEditSalary}
          onDelete={handleDeleteSalary}
          onMarkPaid={handleMarkPaid}
        />
      )}

      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{isEditMode ? t('hr.payroll.editSalary') : t('hr.payroll.createSalary')}</DialogTitle>
          </DialogHeader>
          <SalaryForm
            salary={isEditMode ? selectedSalary : null}
            employees={employees}
            onSubmit={isEditMode && selectedSalary ? (formData: any) => handleUpdateSalary(selectedSalary.id, formData) : handleCreateSalary}
            onCancel={() => { setShowForm(false); setSelectedSalary(null); setIsEditMode(false) }}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}
