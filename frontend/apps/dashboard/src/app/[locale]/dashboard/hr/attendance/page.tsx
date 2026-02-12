'use client'

import { useState, useEffect } from 'react'
import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { useHasPermission } from '@/hooks/useRbac'
import { AttendanceList } from '@/components/hr/AttendanceList'
import { AttendanceForm } from '@/components/hr/AttendanceForm'
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

interface Attendance {
  id: number
  employee_id: number
  employee_name?: string
  clock_in?: string
  clock_out?: string
  status: string
  total_hours?: number
  [key: string]: any
}

interface Employee {
  id: number
  first_name: string
  last_name: string
}

export default function AttendancePage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  // Check RBAC permissions for Attendance module (using HR module for attendance)
  const { hasPermission: canViewAttendance, loading: canViewLoading } = useHasPermission(
    MODULE_IDS.HR,
    PERMISSION_LEVELS.READ,
  )
  const { hasPermission: canCreateAttendance, loading: canCreateLoading } = useHasPermission(
    MODULE_IDS.HR,
    PERMISSION_LEVELS.WRITE,
  )
  const permissionLoading = canViewLoading || canCreateLoading

  // Move ALL hooks BEFORE early returns to avoid Hook order violation
  const [attendance, setAttendance] = useState<Attendance[]>([])
  const [employees, setEmployees] = useState<Employee[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedRecord, setSelectedRecord] = useState<Attendance | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [debugInfo, setDebugInfo] = useState<string>('')
  const [startDate, setStartDate] = useState(
    new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
  )
  const [endDate, setEndDate] = useState(new Date().toISOString().split('T')[0])

  const token = useAuthStore((state) => state.token)

  // Show permission denied if user can't view attendance
  if (permissionLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (!canViewAttendance) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            Attendance
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
                You do not have permission to view attendance records. Please contact your administrator.
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  // Fetch employees
  useEffect(() => {
    const fetchEmployees = async () => {
      try {
        const response = await api.get('/hr/employees')
        const empList = Array.isArray(response.data) ? response.data : response.data?.employees || []
        setEmployees(empList)
      } catch (error) {
        debugError('Failed to load employees:', error)
      }
    }

    if (token) {
      fetchEmployees()
    }
  }, [token])

  // Fetch attendance
  useEffect(() => {
    const fetchAttendance = async () => {
      setLoading(true)
      setApiError(null)
      setDebugInfo(`Fetching attendance records at ${new Date().toLocaleTimeString()}`)

      try {
        debugLog('[ATTENDANCE_PAGE] Starting fetch with dates:', { startDate, endDate })
        const response = await api.get('/hr/attendance', {
          params: {
            start_date: startDate,
            end_date: endDate,
          },
        })
        debugLog('[ATTENDANCE_PAGE] Response:', response.data)

        const attendanceList = Array.isArray(response.data) ? response.data : response.data?.attendance || []

        // Enrich with employee names
        const enrichedAttendance = attendanceList.map(record => ({
          ...record,
          employee_name: employees.find(e => e.id === record.employee_id)
            ? `${employees.find(e => e.id === record.employee_id)?.first_name} ${employees.find(e => e.id === record.employee_id)?.last_name}`
            : undefined,
        }))

        setAttendance(enrichedAttendance)
        setDebugInfo(`Successfully loaded ${enrichedAttendance.length} attendance records at ${new Date().toLocaleTimeString()}`)
      } catch (error: any) {
        const errorMsg = error?.response?.data?.message || error?.message || t('hr.attendance.loadFailed', 'Failed to load attendance')
        debugError('[ATTENDANCE_PAGE] Error:', error)
        setApiError(errorMsg)
        setDebugInfo(`Error: ${errorMsg} at ${new Date().toLocaleTimeString()}`)
        toast.error(errorMsg)
      } finally {
        setLoading(false)
      }
    }

    if (token) {
      fetchAttendance()
    }
  }, [token, startDate, endDate, employees, t])

  const handleCreateAttendance = async (formData: any) => {
    try {
      debugLog('[CREATE_ATTENDANCE] Submitting:', formData)
      const response = await api.post('/hr/attendance/clock-in', formData)
      debugLog('[CREATE_ATTENDANCE] Response:', response.data)

      setAttendance([response.data, ...attendance])
      setShowForm(false)
      toast.success(t('hr.attendance.successCreated', 'Attendance recorded successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.attendance.errorCreated', 'Failed to record attendance')
      debugError('[CREATE_ATTENDANCE] Error:', error)
      toast.error(errorMsg)
      throw error
    }
  }

  const handleUpdateAttendance = async (recordId: number, formData: any) => {
    try {
      debugLog('[UPDATE_ATTENDANCE] Submitting:', formData)
      const response = await api.put(`/hr/attendance/${recordId}`, formData)
      debugLog('[UPDATE_ATTENDANCE] Response:', response.data)

      setAttendance(attendance.map((r) => (r.id === recordId ? { ...r, ...formData } : r)))
      setShowForm(false)
      setSelectedRecord(null)
      setIsEditMode(false)
      toast.success(t('hr.attendance.successUpdated', 'Attendance updated successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.attendance.errorUpdated', 'Failed to update attendance')
      debugError('[UPDATE_ATTENDANCE] Error:', error)
      toast.error(errorMsg)
      throw error
    }
  }

  const handleDeleteAttendance = async (recordId: number) => {
    try {
      debugLog('[DELETE_ATTENDANCE] Deleting:', recordId)
      await api.delete(`/hr/attendance/${recordId}`)
      debugLog('[DELETE_ATTENDANCE] Success')

      setAttendance(attendance.filter((r) => r.id !== recordId))
      toast.success(t('hr.attendance.successDeleted', 'Attendance record deleted successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.attendance.errorDeleted', 'Failed to delete attendance')
      debugError('[DELETE_ATTENDANCE] Error:', error)
      toast.error(errorMsg)
    }
  }

  const handleEditRecord = (record: Attendance) => {
    setSelectedRecord(record)
    setIsEditMode(true)
    setShowForm(true)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">{t('hr.attendance.title')}</h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1">
            {t('hr.attendance.manageAttendance')}
          </p>
        </div>
        {canCreateAttendance ? (
          <Button onClick={() => {
            setIsEditMode(false)
            setSelectedRecord(null)
            setShowForm(true)
          }} className="gap-2 w-full sm:w-auto">
            <Plus className="h-4 w-4" />
            {t('hr.attendance.recordAttendance')}
          </Button>
        ) : (
          <Button disabled className="gap-2 w-full sm:w-auto opacity-50 cursor-not-allowed" title="You don't have permission to record attendance">
            <Plus className="h-4 w-4" />
            {t('hr.attendance.recordAttendance')}
          </Button>
        )}
      </div>

      {/* Filters */}
      <div className="flex flex-col md:flex-row gap-4 bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
        <div className="flex-1 space-y-2">
          <label className="block text-sm font-medium">{t('hr.attendance.fromDate')}</label>
          <Input
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
          />
        </div>
        <div className="flex-1 space-y-2">
          <label className="block text-sm font-medium">{t('hr.attendance.toDate')}</label>
          <Input
            type="date"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
          />
        </div>
      </div>

      {/* API Error Alert */}
      {apiError && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex gap-3">
          <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <div>
            <h3 className="font-semibold text-red-900 dark:text-red-300">{t('hr.attendance.errorLoading')}</h3>
            <p className="text-red-700 dark:text-red-400 text-sm mt-1">{apiError}</p>
            <button
              onClick={() => window.location.reload()}
              className="text-red-600 dark:text-red-400 hover:underline text-sm mt-2 font-medium"
            >
              {t('hr.attendance.retry')}
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

      {/* Attendance List */}
      {loading && attendance.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <AttendanceList
          attendance={attendance}
          loading={loading}
          onEdit={handleEditRecord}
          onDelete={handleDeleteAttendance}
        />
      )}

      {/* Attendance Form Dialog */}
      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {isEditMode ? t('hr.attendance.editRecord') : t('hr.attendance.recordAttendance')}
            </DialogTitle>
          </DialogHeader>
          <AttendanceForm
            attendance={isEditMode ? selectedRecord : null}
            employees={employees}
            onSubmit={
              isEditMode && selectedRecord
                ? (formData: any) => handleUpdateAttendance(selectedRecord.id, formData)
                : handleCreateAttendance
            }
            onCancel={() => {
              setShowForm(false)
              setSelectedRecord(null)
              setIsEditMode(false)
            }}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}
