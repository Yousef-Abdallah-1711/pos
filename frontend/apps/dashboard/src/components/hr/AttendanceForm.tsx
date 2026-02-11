'use client'

import { useState, useEffect } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { toast } from 'sonner'
import { debugLog as debugLogFn, debugError } from '@/lib/debug'
import { AlertCircle, Loader, Clock } from 'lucide-react'
import { api } from '@/lib/api'

interface AttendanceFormProps {
  attendance?: any
  employees?: any[]
  onSubmit: (formData: any) => Promise<void>
  onCancel: () => void
}

export function AttendanceForm({ attendance, employees = [], onSubmit, onCancel }: AttendanceFormProps) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [debugLog, setDebugLog] = useState<string[]>([])
  const [employeesList, setEmployeesList] = useState(employees)

  const { control, handleSubmit, formState: { errors }, watch, getValues } = useForm({
    defaultValues: {
      employee_id: attendance?.employee_id || '',
      clock_in: attendance?.clock_in || new Date().toISOString().slice(0, 16),
      clock_out: attendance?.clock_out || '',
      scheduled_clock_in: attendance?.scheduled_clock_in || '',
      scheduled_clock_out: attendance?.scheduled_clock_out || '',
      total_break_minutes: attendance?.total_break_minutes || '0',
      status: attendance?.status || 'present',
      notes: attendance?.notes || '',
      admin_notes: attendance?.admin_notes || '',
    },
  })

  // Fetch employees if not provided
  useEffect(() => {
    if (employees.length === 0) {
      const fetchEmployees = async () => {
        try {
          const response = await api.get('/hr/employees')
          const emps = Array.isArray(response.data) ? response.data : response.data?.employees || []
          setEmployeesList(emps)
          addDebugLog(`Loaded ${emps.length} employees`)
        } catch (error) {
          debugError('Failed to load employees:', error)
        }
      }
      fetchEmployees()
    }
  }, [employees])

  const addDebugLog = (message: string) => {
    setDebugLog(prev => [...prev, `[${new Date().toLocaleTimeString()}] ${message}`])
  }

  const onFormSubmit = async (data: any) => {
    setIsSubmitting(true)
    setApiError(null)
    setDebugLog([])

    try {
      addDebugLog('Starting form submission...')

      if (!data.employee_id) {
        throw new Error('Employee is required')
      }

      if (!attendance && !data.clock_in) {
        throw new Error('Clock in time is required')
      }

      addDebugLog('Validation passed')

      const cleanData = {
        ...data,
        employee_id: parseInt(data.employee_id.toString()),
        total_break_minutes: parseInt((data.total_break_minutes || 0).toString()),
        clock_in: data.clock_in ? new Date(data.clock_in).toISOString() : undefined,
        clock_out: data.clock_out ? new Date(data.clock_out).toISOString() : undefined,
      }

      addDebugLog(`Cleaned data: ${JSON.stringify(cleanData)}`)
      debugLogFn('DEBUG: Attendance Form Submission', { raw: data, cleaned: cleanData })

      await onSubmit(cleanData)
      addDebugLog('Form submission successful')
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || 'Failed to save attendance'
      setApiError(errorMsg)
      addDebugLog(`Error: ${errorMsg}`)
      debugError('Attendance Form Error:', error)
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
            <h3 className="font-semibold text-red-900 dark:text-red-300">Error Saving Attendance</h3>
            <p className="text-red-700 dark:text-red-400 text-sm mt-1">{apiError}</p>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-2">
          <label className="block text-sm font-medium">Employee *</label>
          <Controller
            name="employee_id"
            control={control}
            rules={{ required: 'Employee is required' }}
            render={({ field }) => (
              <select
                {...field}
                className={`w-full px-3 py-2 border rounded-md bg-white dark:bg-gray-800 ${
                  errors.employee_id ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'
                }`}
                disabled={!!attendance}
              >
                <option value="">Select employee...</option>
                {employeesList.map(emp => (
                  <option key={emp.id} value={emp.id}>
                    {emp.first_name} {emp.last_name}
                  </option>
                ))}
              </select>
            )}
          />
          {errors.employee_id && (
            <span className="text-red-500 text-sm flex items-center gap-1">
              <AlertCircle className="h-4 w-4" />
              {errors.employee_id.message}
            </span>
          )}
        </div>

        <div className="space-y-2">
          <label className="block text-sm font-medium">Status</label>
          <Controller
            name="status"
            control={control}
            render={({ field }) => (
              <select
                {...field}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800"
              >
                <option value="present">Present</option>
                <option value="absent">Absent</option>
                <option value="late">Late</option>
                <option value="early_leave">Early Leave</option>
                <option value="half_day">Half Day</option>
              </select>
            )}
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-2">
          <label className="block text-sm font-medium">Clock In Time</label>
          <Controller
            name="clock_in"
            control={control}
            render={({ field }) => (
              <Input {...field} type="datetime-local" />
            )}
          />
        </div>

        <div className="space-y-2">
          <label className="block text-sm font-medium">Clock Out Time</label>
          <Controller
            name="clock_out"
            control={control}
            render={({ field }) => (
              <Input {...field} type="datetime-local" />
            )}
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-2">
          <label className="block text-sm font-medium">Scheduled Clock In</label>
          <Controller
            name="scheduled_clock_in"
            control={control}
            render={({ field }) => (
              <Input {...field} type="time" placeholder="08:00" />
            )}
          />
        </div>

        <div className="space-y-2">
          <label className="block text-sm font-medium">Scheduled Clock Out</label>
          <Controller
            name="scheduled_clock_out"
            control={control}
            render={({ field }) => (
              <Input {...field} type="time" placeholder="17:00" />
            )}
          />
        </div>
      </div>

      <div className="space-y-2">
        <label className="block text-sm font-medium">Total Break Minutes</label>
        <Controller
          name="total_break_minutes"
          control={control}
          render={({ field }) => (
            <Input {...field} type="number" min="0" step="15" placeholder="0" />
          )}
        />
      </div>

      <div className="space-y-2">
        <label className="block text-sm font-medium">Employee Notes</label>
        <Controller
          name="notes"
          control={control}
          render={({ field }) => (
            <Textarea {...field} placeholder="Any notes about this attendance..." rows={3} />
          )}
        />
      </div>

      <div className="space-y-2">
        <label className="block text-sm font-medium">Admin Notes</label>
        <Controller
          name="admin_notes"
          control={control}
          render={({ field }) => (
            <Textarea {...field} placeholder="Administrative notes..." rows={3} />
          )}
        />
      </div>

      {/* Debug Log */}
      {debugLog.length > 0 && (
        <div className="bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
          <h4 className="text-xs font-semibold text-gray-600 dark:text-gray-400 mb-2">Debug Log</h4>
          <div className="space-y-1 max-h-40 overflow-y-auto">
            {debugLog.map((log, idx) => (
              <div key={idx} className="text-xs text-gray-600 dark:text-gray-400 font-mono">
                {log}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Form Actions */}
      <div className="flex justify-end gap-2 pt-4 border-t">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting} className="gap-2">
          {isSubmitting && <Loader className="h-4 w-4 animate-spin" />}
          {isSubmitting ? 'Saving...' : attendance ? 'Update Attendance' : 'Record Attendance'}
        </Button>
      </div>
    </form>
  )
}
