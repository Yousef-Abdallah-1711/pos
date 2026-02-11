'use client'

import { useState, useEffect } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { toast } from 'sonner'
import { debugLog as debugLogFn, debugError } from '@/lib/debug'
import { AlertCircle, Loader, FileText } from 'lucide-react'
import { api } from '@/lib/api'

interface LeaveFormProps {
  leave?: any
  employees?: any[]
  onSubmit: (formData: any) => Promise<void>
  onCancel: () => void
}

const LEAVE_TYPES = [
  { value: 'annual', label: 'Annual Leave' },
  { value: 'sick', label: 'Sick Leave' },
  { value: 'casual', label: 'Casual Leave' },
  { value: 'maternity', label: 'Maternity Leave' },
  { value: 'paternity', label: 'Paternity Leave' },
  { value: 'unpaid', label: 'Unpaid Leave' },
  { value: 'compensatory', label: 'Compensatory Leave' },
]

const LEAVE_CATEGORIES = [
  { value: 'full_day', label: 'Full Day' },
  { value: 'half_day_morning', label: 'Half Day (Morning)' },
  { value: 'half_day_afternoon', label: 'Half Day (Afternoon)' },
]

export function LeaveForm({ leave, employees = [], onSubmit, onCancel }: LeaveFormProps) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [debugLog, setDebugLog] = useState<string[]>([])
  const [employeesList, setEmployeesList] = useState(employees)

  const { control, handleSubmit, formState: { errors }, watch } = useForm({
    defaultValues: {
      employee_id: leave?.employee_id || '',
      start_date: leave?.start_date?.split('T')[0] || '',
      end_date: leave?.end_date?.split('T')[0] || '',
      leave_type: leave?.leave_type || 'annual',
      leave_category: leave?.leave_category || 'full_day',
      reason: leave?.reason || '',
      contact_number: leave?.contact_number || '',
      contact_address: leave?.contact_address || '',
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
      if (!data.start_date) {
        throw new Error('Start date is required')
      }
      if (!data.end_date) {
        throw new Error('End date is required')
      }
      if (!data.reason?.trim()) {
        throw new Error('Reason is required')
      }

      // Validate dates
      const startDate = new Date(data.start_date)
      const endDate = new Date(data.end_date)
      if (endDate < startDate) {
        throw new Error('End date must be after start date')
      }

      addDebugLog('Validation passed')

      const cleanData = {
        ...data,
        employee_id: parseInt(data.employee_id.toString()),
        start_date: new Date(data.start_date).toISOString(),
        end_date: new Date(data.end_date).toISOString(),
      }

      addDebugLog(`Cleaned data: ${JSON.stringify(cleanData)}`)
      debugLogFn('DEBUG: Leave Form Submission', { raw: data, cleaned: cleanData })

      await onSubmit(cleanData)
      addDebugLog('Form submission successful')
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || 'Failed to save leave'
      setApiError(errorMsg)
      addDebugLog(`Error: ${errorMsg}`)
      debugError('Leave Form Error:', error)
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
            <h3 className="font-semibold text-red-900 dark:text-red-300">Error Saving Leave</h3>
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
                disabled={!!leave}
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
          <label className="block text-sm font-medium">Leave Type *</label>
          <Controller
            name="leave_type"
            control={control}
            render={({ field }) => (
              <select
                {...field}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800"
              >
                {LEAVE_TYPES.map(type => (
                  <option key={type.value} value={type.value}>{type.label}</option>
                ))}
              </select>
            )}
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-2">
          <label className="block text-sm font-medium">Start Date *</label>
          <Controller
            name="start_date"
            control={control}
            rules={{ required: 'Start date is required' }}
            render={({ field }) => (
              <Input {...field} type="date" className={errors.start_date ? 'border-red-500' : ''} />
            )}
          />
          {errors.start_date && (
            <span className="text-red-500 text-sm flex items-center gap-1">
              <AlertCircle className="h-4 w-4" />
              {errors.start_date.message}
            </span>
          )}
        </div>

        <div className="space-y-2">
          <label className="block text-sm font-medium">End Date *</label>
          <Controller
            name="end_date"
            control={control}
            rules={{ required: 'End date is required' }}
            render={({ field }) => (
              <Input {...field} type="date" className={errors.end_date ? 'border-red-500' : ''} />
            )}
          />
          {errors.end_date && (
            <span className="text-red-500 text-sm flex items-center gap-1">
              <AlertCircle className="h-4 w-4" />
              {errors.end_date.message}
            </span>
          )}
        </div>
      </div>

      <div className="space-y-2">
        <label className="block text-sm font-medium">Leave Category</label>
        <Controller
          name="leave_category"
          control={control}
          render={({ field }) => (
            <select
              {...field}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800"
            >
              {LEAVE_CATEGORIES.map(cat => (
                <option key={cat.value} value={cat.value}>{cat.label}</option>
              ))}
            </select>
          )}
        />
      </div>

      <div className="space-y-2">
        <label className="block text-sm font-medium">Reason *</label>
        <Controller
          name="reason"
          control={control}
          rules={{ required: 'Reason is required' }}
          render={({ field }) => (
            <Textarea
              {...field}
              placeholder="State the reason for leave..."
              rows={4}
              className={errors.reason ? 'border-red-500' : ''}
            />
          )}
        />
        {errors.reason && (
          <span className="text-red-500 text-sm flex items-center gap-1">
            <AlertCircle className="h-4 w-4" />
            {errors.reason.message}
          </span>
        )}
      </div>

      <div className="border-t pt-4">
        <h3 className="font-semibold mb-4">Contact Information</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="block text-sm font-medium">Contact Number</label>
            <Controller
              name="contact_number"
              control={control}
              render={({ field }) => (
                <Input {...field} type="tel" placeholder="+20123456789" />
              )}
            />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium">Contact Address</label>
            <Controller
              name="contact_address"
              control={control}
              render={({ field }) => (
                <Input {...field} placeholder="Address during leave..." />
              )}
            />
          </div>
        </div>
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
          {isSubmitting ? 'Saving...' : leave ? 'Update Leave' : 'Request Leave'}
        </Button>
      </div>
    </form>
  )
}
