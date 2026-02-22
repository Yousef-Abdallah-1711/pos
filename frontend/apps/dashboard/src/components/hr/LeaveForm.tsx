'use client'

import { useState, useEffect, useMemo } from 'react'
import { usePathname } from 'next/navigation'
import { useForm, Controller } from 'react-hook-form'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { toast } from 'sonner'
import { debugLog as debugLogFn, debugError } from '@/lib/debug'
import { AlertCircle, Loader } from 'lucide-react'
import { api } from '@/lib/api'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'

interface LeaveFormProps {
  leave?: any
  employees?: any[]
  onSubmit: (formData: any) => Promise<void>
  onCancel: () => void
}

export function LeaveForm({ leave, employees = [], onSubmit, onCancel }: LeaveFormProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])

  const LEAVE_TYPES = [
    { value: 'annual', label: t('hr.leaves.types.annual') },
    { value: 'sick', label: t('hr.leaves.types.sick') },
    { value: 'casual', label: t('hr.leaves.types.casual') },
    { value: 'maternity', label: t('hr.leaves.types.maternity') },
    { value: 'paternity', label: t('hr.leaves.types.paternity') },
    { value: 'unpaid', label: t('hr.leaves.types.unpaid') },
    { value: 'compensatory', label: t('hr.leaves.types.compensatory') },
  ]

  const LEAVE_CATEGORIES = [
    { value: 'full_day', label: t('hr.leaves.categories.fullDay') },
    { value: 'half_day_morning', label: t('hr.leaves.categories.halfDayMorning') },
    { value: 'half_day_afternoon', label: t('hr.leaves.categories.halfDayAfternoon') },
  ]

  const [isSubmitting, setIsSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [debugLog, setDebugLog] = useState<string[]>([])
  const [employeesList, setEmployeesList] = useState(employees)

  const { control, handleSubmit, formState: { errors } } = useForm({
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

  useEffect(() => {
    if (employees.length === 0) {
      const fetchEmployees = async () => {
        try {
          const response = await api.get('/hr/employees')
          const emps = Array.isArray(response.data) ? response.data : (response.data?.data || response.data?.employees || [])
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
        throw new Error(t('hr.leaves.employeeRequired'))
      }
      if (!data.start_date) {
        throw new Error(t('hr.leaves.startDateRequired'))
      }
      if (!data.end_date) {
        throw new Error(t('hr.leaves.endDateRequired'))
      }
      if (!data.reason?.trim()) {
        throw new Error(t('hr.leaves.reasonRequired'))
      }

      const startDate = new Date(data.start_date)
      const endDate = new Date(data.end_date)
      if (endDate < startDate) {
        throw new Error(t('hr.leaves.endDateAfterStart'))
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
      const errorMsg = error?.response?.data?.message || error?.message || t('hr.leaves.saveFailed')
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
      {apiError && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex gap-3">
          <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <div>
            <h3 className="font-semibold text-red-900 dark:text-red-300">{t('hr.leaves.errorSaving')}</h3>
            <p className="text-red-700 dark:text-red-400 text-sm mt-1">{apiError}</p>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-2">
          <label className="block text-sm font-medium">{t('hr.leaves.employee')} *</label>
          <Controller
            name="employee_id"
            control={control}
            rules={{ required: t('hr.leaves.employeeRequired') }}
            render={({ field }) => (
              <select
                {...field}
                className={`w-full px-3 py-2 border rounded-md bg-white dark:bg-gray-800 ${
                  errors.employee_id ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'
                }`}
                disabled={!!leave}
              >
                <option value="">{t('hr.leaves.selectEmployee')}</option>
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
          <label className="block text-sm font-medium">{t('hr.leaves.leaveType')} *</label>
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
          <label className="block text-sm font-medium">{t('hr.leaves.startDate')} *</label>
          <Controller
            name="start_date"
            control={control}
            rules={{ required: t('hr.leaves.startDateRequired') }}
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
          <label className="block text-sm font-medium">{t('hr.leaves.endDate')} *</label>
          <Controller
            name="end_date"
            control={control}
            rules={{ required: t('hr.leaves.endDateRequired') }}
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
        <label className="block text-sm font-medium">{t('hr.leaves.leaveCategory')}</label>
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
        <label className="block text-sm font-medium">{t('hr.leaves.reason')} *</label>
        <Controller
          name="reason"
          control={control}
          rules={{ required: t('hr.leaves.reasonRequired') }}
          render={({ field }) => (
            <Textarea
              {...field}
              placeholder={t('hr.leaves.reasonPlaceholder')}
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

      <div className="flex justify-end gap-2 pt-4 border-t">
        <Button type="button" variant="outline" onClick={onCancel}>
          {t('common.cancel')}
        </Button>
        <Button type="submit" disabled={isSubmitting} className="gap-2">
          {isSubmitting && <Loader className="h-4 w-4 animate-spin" />}
          {isSubmitting ? t('common.saving') : leave ? t('hr.leaves.update') : t('hr.leaves.submit')}
        </Button>
      </div>
    </form>
  )
}
