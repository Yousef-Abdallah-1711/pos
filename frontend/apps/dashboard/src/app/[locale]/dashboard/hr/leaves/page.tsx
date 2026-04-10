'use client'

import { useState, useEffect, useMemo } from 'react'
import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { useUserPermissions } from '@/hooks/useRbac'
import { LeaveList } from '@/components/hr/LeaveList'
import { LeaveForm } from '@/components/hr/LeaveForm'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@/components/ui/dialog'
import { api } from '@/lib/api'
import { debugLog, debugError } from '@/lib/debug'
import { toast } from 'sonner'
import { Plus, AlertCircle, Lock, Loader2 } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { PageHeaderSkeleton, TableSkeleton } from '@/components/skeletons'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { hasAuthenticatedSession } from '@/lib/authSession'

interface Leave {
  id: number
  employee_id: number
  employee_name?: string
  start_date: string
  end_date: string
  leave_type: string
  status: string
  reason?: string
  [key: string]: any
}

interface Employee {
  id: number
  first_name: string
  last_name: string
}

export default function LeavesPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  // Check RBAC permissions - properly destructure the hook return
  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canViewLeaves = checkPermission(MODULE_IDS.LEAVES, PERMISSION_LEVELS.READ)
  const canCreateLeaves = checkPermission(MODULE_IDS.LEAVES, PERMISSION_LEVELS.WRITE)

  // ALL hooks MUST be before early returns to avoid Hook order violation
  const [leaves, setLeaves] = useState<Leave[]>([])
  const [employees, setEmployees] = useState<Employee[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedLeave, setSelectedLeave] = useState<Leave | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [debugInfo, setDebugInfo] = useState<string>('')
  const [actionDialog, setActionDialog] = useState<{
    type: 'approve' | 'reject' | null
    leaveId: number | null
    open: boolean
    input: string
  }>({
    type: null,
    leaveId: null,
    open: false,
    input: '',
  })
  const [startDate, setStartDate] = useState(
    new Date(Date.now() - 365 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
  )
  const [endDate, setEndDate] = useState(
    new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
  )

  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const hasSession = hasAuthenticatedSession(user, token)

  // Fetch employees - hook MUST be before early returns
  useEffect(() => {
    if (!hasSession || !canViewLeaves || permissionLoading) return

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
  }, [hasSession, canViewLeaves, permissionLoading])

  // Fetch leaves - hook MUST be before early returns
  useEffect(() => {
    if (!hasSession || !canViewLeaves || permissionLoading) return

    const fetchLeaves = async () => {
      setLoading(true)
      setApiError(null)
      setDebugInfo(`Fetching leave requests at ${new Date().toLocaleTimeString()}`)

      try {
        debugLog('[LEAVES_PAGE] Starting fetch with dates:', { startDate, endDate })
        const response = await api.get('/hr/leaves', {
          params: {
            start_date: startDate,
            end_date: endDate,
          },
        })
        debugLog('[LEAVES_PAGE] Response:', response.data)

        const leavesList = Array.isArray(response.data) ? response.data : (response.data?.data || response.data?.leaves || [])

        // Enrich with employee names
        const enrichedLeaves = leavesList.map(leave => ({
          ...leave,
          employee_name: employees.find(e => e.id === leave.employee_id)
            ? `${employees.find(e => e.id === leave.employee_id)?.first_name} ${employees.find(e => e.id === leave.employee_id)?.last_name}`
            : undefined,
        }))

        setLeaves(enrichedLeaves)
        setDebugInfo(`Successfully loaded ${enrichedLeaves.length} leave requests at ${new Date().toLocaleTimeString()}`)
      } catch (error: any) {
        const errorMsg = error?.response?.data?.message || error?.message || 'Failed to load leaves'
        debugError('[LEAVES_PAGE] Error:', error)
        setApiError(errorMsg)
        setDebugInfo(`Error: ${errorMsg} at ${new Date().toLocaleTimeString()}`)
        toast.error(errorMsg)
      } finally {
        setLoading(false)
      }
    }

    fetchLeaves()
  }, [hasSession, startDate, endDate, employees, canViewLeaves, permissionLoading])

  // Show loading while permissions are being checked
  if (permissionLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  // Show permission denied if user can't view leaves
  if (!canViewLeaves) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            {t('hr.leaves.title')}
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
                You do not have permission to view leave requests. Please contact your administrator.
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  const handleCreateLeave = async (formData: any) => {
    try {
      debugLog('[CREATE_LEAVE] Submitting:', formData)
      const response = await api.post('/hr/leaves', formData)
      debugLog('[CREATE_LEAVE] Response:', response.data)

      const created = response.data?.data || response.data
      setLeaves([created, ...leaves])
      setShowForm(false)
      toast.success(t('hr.leaves.successCreated', 'Leave request created successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.leaves.errorCreated', 'Failed to create leave')
      debugError('[CREATE_LEAVE] Error:', error)
      toast.error(errorMsg)
      throw error
    }
  }

  const handleUpdateLeave = async (leaveId: number, formData: any) => {
    try {
      debugLog('[UPDATE_LEAVE] Submitting:', formData)
      const response = await api.put(`/hr/leaves/${leaveId}`, formData)
      debugLog('[UPDATE_LEAVE] Response:', response.data)

      setLeaves(leaves.map((l) => (l.id === leaveId ? { ...l, ...formData } : l)))
      setShowForm(false)
      setSelectedLeave(null)
      setIsEditMode(false)
      toast.success(t('hr.leaves.successUpdated', 'Leave request updated successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.leaves.errorUpdated', 'Failed to update leave')
      debugError('[UPDATE_LEAVE] Error:', error)
      toast.error(errorMsg)
      throw error
    }
  }

  const handleDeleteLeave = async (leaveId: number) => {
    try {
      debugLog('[DELETE_LEAVE] Deleting:', leaveId)
      await api.delete(`/hr/leaves/${leaveId}`)
      debugLog('[DELETE_LEAVE] Success')

      setLeaves(leaves.filter((l) => l.id !== leaveId))
      toast.success(t('hr.leaves.successDeleted', 'Leave request deleted successfully'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.leaves.errorDeleted', 'Failed to delete leave')
      debugError('[DELETE_LEAVE] Error:', error)
      toast.error(errorMsg)
    }
  }

  const handleApproveLeave = async () => {
    if (!actionDialog.leaveId) return

    try {
      debugLog('[APPROVE_LEAVE] Approving:', actionDialog.leaveId)
      await api.post(`/hr/leaves/${actionDialog.leaveId}/approve`, {
        notes: actionDialog.input,
      })
      debugLog('[APPROVE_LEAVE] Success')

      setLeaves(leaves.map(l => l.id === actionDialog.leaveId ? { ...l, status: 'approved' } : l))
      setActionDialog({ type: null, leaveId: null, open: false, input: '' })
      toast.success(t('hr.leaves.successApproved', 'Leave request approved'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.leaves.errorApproved', 'Failed to approve leave')
      debugError('[APPROVE_LEAVE] Error:', error)
      toast.error(errorMsg)
    }
  }

  const handleRejectLeave = async () => {
    if (!actionDialog.leaveId) return

    try {
      debugLog('[REJECT_LEAVE] Rejecting:', actionDialog.leaveId)
      await api.post(`/hr/leaves/${actionDialog.leaveId}/reject`, {
        reason: actionDialog.input,
      })
      debugLog('[REJECT_LEAVE] Success')

      setLeaves(leaves.map(l => l.id === actionDialog.leaveId ? { ...l, status: 'rejected' } : l))
      setActionDialog({ type: null, leaveId: null, open: false, input: '' })
      toast.success(t('hr.leaves.successRejected', 'Leave request rejected'))
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.leaves.errorRejected', 'Failed to reject leave')
      debugError('[REJECT_LEAVE] Error:', error)
      toast.error(errorMsg)
    }
  }

  const handleEditLeave = (leave: Leave) => {
    setSelectedLeave(leave)
    setIsEditMode(true)
    setShowForm(true)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">{t('hr.leaves.title')}</h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1">
            {t('hr.leaves.manageLeaves')}
          </p>
        </div>
        {canCreateLeaves ? (
          <Button onClick={() => {
            setIsEditMode(false)
            setSelectedLeave(null)
            setShowForm(true)
          }} className="gap-2 w-full sm:w-auto">
            <Plus className="h-4 w-4" />
            {t('hr.leaves.requestLeave')}
          </Button>
        ) : (
          <Button disabled className="gap-2 w-full sm:w-auto opacity-50 cursor-not-allowed" title="You don't have permission to request leaves">
            <Plus className="h-4 w-4" />
            {t('hr.leaves.requestLeave')}
          </Button>
        )}
      </div>

      {/* Filters */}
      <div className="flex flex-col md:flex-row gap-4 bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
        <div className="flex-1 space-y-2">
          <label className="block text-sm font-medium">{t('hr.leaves.fromDate')}</label>
          <Input
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
          />
        </div>
        <div className="flex-1 space-y-2">
          <label className="block text-sm font-medium">{t('hr.leaves.toDate')}</label>
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
            <h3 className="font-semibold text-red-900 dark:text-red-300">{t('hr.leaves.errorLoading')}</h3>
            <p className="text-red-700 dark:text-red-400 text-sm mt-1">{apiError}</p>
            <button
              onClick={() => window.location.reload()}
              className="text-red-600 dark:text-red-400 hover:underline text-sm mt-2 font-medium"
            >
              {t('hr.leaves.retry')}
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

      {/* Leave List */}
      {loading && leaves.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <LeaveList
          leaves={leaves}
          loading={loading}
          onEdit={handleEditLeave}
          onDelete={handleDeleteLeave}
          onApprove={(id) => setActionDialog({ type: 'approve', leaveId: id, open: true, input: '' })}
          onReject={(id) => setActionDialog({ type: 'reject', leaveId: id, open: true, input: '' })}
          canApprove={true}
        />
      )}

      {/* Leave Form Dialog */}
      <Dialog open={showForm} onOpenChange={setShowForm}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {isEditMode ? t('hr.leaves.approveLeave') : t('hr.leaves.requestLeave')}
            </DialogTitle>
          </DialogHeader>
          <LeaveForm
            leave={isEditMode ? selectedLeave : null}
            employees={employees}
            onSubmit={
              isEditMode && selectedLeave
                ? (formData: any) => handleUpdateLeave(selectedLeave.id, formData)
                : handleCreateLeave
            }
            onCancel={() => {
              setShowForm(false)
              setSelectedLeave(null)
              setIsEditMode(false)
            }}
          />
        </DialogContent>
      </Dialog>

      {/* Approval/Rejection Dialog */}
      <Dialog open={actionDialog.open} onOpenChange={(open) => {
        if (!open) setActionDialog({ type: null, leaveId: null, open: false, input: '' })
      }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {actionDialog.type === 'approve' ? t('hr.leaves.approveLeave') : t('hr.leaves.rejectLeave')}
            </DialogTitle>
            <DialogDescription>
              <div className="space-y-4">
                <p>
                  {actionDialog.type === 'approve'
                    ? t('hr.leaves.confirmApprove')
                    : t('hr.leaves.confirmReject')}
                </p>
                <div>
                  <label className="block text-sm font-medium mb-2">
                    {actionDialog.type === 'approve' ? t('hr.leaves.notes') : t('hr.leaves.rejectionReason')}
                  </label>
                  <textarea
                    value={actionDialog.input}
                    onChange={(e) => setActionDialog({ ...actionDialog, input: e.target.value })}
                    placeholder={actionDialog.type === 'approve' ? t('hr.leaves.addNotes') : t('hr.leaves.provideReason')}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-sm"
                    rows={3}
                  />
                </div>
              </div>
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setActionDialog({ ...actionDialog, open: false })}
            >
              {t('common.cancel')}
            </Button>
            <Button
              onClick={actionDialog.type === 'approve' ? handleApproveLeave : handleRejectLeave}
              className={actionDialog.type === 'approve' ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700'}
            >
              {actionDialog.type === 'approve' ? t('hr.leaves.approveLeave') : t('hr.leaves.rejectLeave')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
