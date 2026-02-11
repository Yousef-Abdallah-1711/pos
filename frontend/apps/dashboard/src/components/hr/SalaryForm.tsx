'use client'

import { useState, useEffect } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { toast } from 'sonner'
import { debugLog as debugLogFn, debugError } from '@/lib/debug'
import { AlertCircle, Loader } from 'lucide-react'
import { api } from '@/lib/api'

interface SalaryFormProps {
  salary?: any
  employees?: any[]
  onSubmit: (formData: any) => Promise<void>
  onCancel: () => void
}

export function SalaryForm({ salary, employees = [], onSubmit, onCancel }: SalaryFormProps) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [debugLog, setDebugLog] = useState<string[]>([])
  const [employeesList, setEmployeesList] = useState(employees)

  const { control, handleSubmit, formState: { errors }, watch } = useForm({
    defaultValues: {
      employee_id: salary?.employee_id || '',
      pay_period_start: salary?.pay_period_start?.split('T')[0] || '',
      pay_period_end: salary?.pay_period_end?.split('T')[0] || '',
      base_salary: salary?.base_salary || '',
      currency: salary?.currency || 'EGP',
      overtime_hours: salary?.overtime_hours || '0',
      overtime_rate: salary?.overtime_rate || '0',
      bonus: salary?.bonus || '0',
      commission: salary?.commission || '0',
      allowances: salary?.allowances || '0',
      tips: salary?.tips || '0',
      tax: salary?.tax || '0',
      social_insurance: salary?.social_insurance || '0',
      health_insurance: salary?.health_insurance || '0',
      pension: salary?.pension || '0',
      loan_deduction: salary?.loan_deduction || '0',
      advance_deduction: salary?.advance_deduction || '0',
      days_worked: salary?.days_worked || '20',
      days_absent: salary?.days_absent || '0',
      total_hours_worked: salary?.total_hours_worked || '160',
      payment_method: salary?.payment_method || 'bank_transfer',
      status: salary?.status || 'pending',
    },
  })

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
      if (!data.pay_period_start) {
        throw new Error('Pay period start date is required')
      }
      if (!data.pay_period_end) {
        throw new Error('Pay period end date is required')
      }

      addDebugLog('Validation passed')

      const cleanData = {
        ...data,
        employee_id: parseInt(data.employee_id.toString()),
        base_salary: parseFloat(data.base_salary.toString()) || 0,
        overtime_hours: parseFloat(data.overtime_hours.toString()) || 0,
        overtime_rate: parseFloat(data.overtime_rate.toString()) || 0,
        bonus: parseFloat(data.bonus.toString()) || 0,
        commission: parseFloat(data.commission.toString()) || 0,
        allowances: parseFloat(data.allowances.toString()) || 0,
        tips: parseFloat(data.tips.toString()) || 0,
        tax: parseFloat(data.tax.toString()) || 0,
        social_insurance: parseFloat(data.social_insurance.toString()) || 0,
        health_insurance: parseFloat(data.health_insurance.toString()) || 0,
        pension: parseFloat(data.pension.toString()) || 0,
        loan_deduction: parseFloat(data.loan_deduction.toString()) || 0,
        advance_deduction: parseFloat(data.advance_deduction.toString()) || 0,
        days_worked: parseInt(data.days_worked.toString()) || 0,
        days_absent: parseInt(data.days_absent.toString()) || 0,
        total_hours_worked: parseFloat(data.total_hours_worked.toString()) || 0,
        pay_period_start: new Date(data.pay_period_start).toISOString(),
        pay_period_end: new Date(data.pay_period_end).toISOString(),
      }

      addDebugLog(`Cleaned data: ${JSON.stringify(cleanData)}`)
      debugLogFn('DEBUG: Salary Form Submission', { raw: data, cleaned: cleanData })

      await onSubmit(cleanData)
      addDebugLog('Form submission successful')
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || 'Failed to save salary'
      setApiError(errorMsg)
      addDebugLog(`Error: ${errorMsg}`)
      debugError('Salary Form Error:', error)
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
            <h3 className="font-semibold text-red-900 dark:text-red-300">Error Saving Salary</h3>
            <p className="text-red-700 dark:text-red-400 text-sm mt-1">{apiError}</p>
          </div>
        </div>
      )}

      <Tabs defaultValue="general" className="w-full">
        <TabsList className="grid w-full grid-cols-4 bg-gray-100 dark:bg-gray-800 p-1">
          <TabsTrigger value="general">General</TabsTrigger>
          <TabsTrigger value="earnings">Earnings</TabsTrigger>
          <TabsTrigger value="deductions">Deductions</TabsTrigger>
          <TabsTrigger value="records">Records</TabsTrigger>
        </TabsList>

        <TabsContent value="general" className="space-y-4 mt-6">
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
                    className={`w-full px-3 py-2 border rounded-md bg-white dark:bg-gray-800 ${errors.employee_id ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'}`}
                    disabled={!!salary}
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
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Status</label>
              <Controller
                name="status"
                control={control}
                render={({ field }) => (
                  <select {...field} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800">
                    <option value="pending">Pending</option>
                    <option value="approved">Approved</option>
                    <option value="paid">Paid</option>
                  </select>
                )}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Pay Period Start *</label>
              <Controller
                name="pay_period_start"
                control={control}
                rules={{ required: 'Start date is required' }}
                render={({ field }) => (
                  <Input {...field} type="date" />
                )}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Pay Period End *</label>
              <Controller
                name="pay_period_end"
                control={control}
                rules={{ required: 'End date is required' }}
                render={({ field }) => (
                  <Input {...field} type="date" />
                )}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Base Salary</label>
              <Controller
                name="base_salary"
                control={control}
                render={({ field }) => (
                  <Input {...field} type="number" step="0.01" placeholder="0.00" />
                )}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Currency</label>
              <Controller
                name="currency"
                control={control}
                render={({ field }) => (
                  <select {...field} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800">
                    <option value="EGP">EGP</option>
                    <option value="USD">USD</option>
                    <option value="EUR">EUR</option>
                  </select>
                )}
              />
            </div>
          </div>
        </TabsContent>

        <TabsContent value="earnings" className="space-y-4 mt-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Overtime Hours</label>
              <Controller name="overtime_hours" control={control} render={({ field }) => <Input {...field} type="number" step="0.5" placeholder="0" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Overtime Rate</label>
              <Controller name="overtime_rate" control={control} render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Bonus</label>
              <Controller name="bonus" control={control} render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Commission</label>
              <Controller name="commission" control={control} render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Allowances</label>
              <Controller name="allowances" control={control} render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Tips</label>
              <Controller name="tips" control={control} render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />} />
            </div>
          </div>
        </TabsContent>

        <TabsContent value="deductions" className="space-y-4 mt-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Tax</label>
              <Controller name="tax" control={control} render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Social Insurance</label>
              <Controller name="social_insurance" control={control} render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Health Insurance</label>
              <Controller name="health_insurance" control={control} render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Pension</label>
              <Controller name="pension" control={control} render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Loan Deduction</label>
              <Controller name="loan_deduction" control={control} render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Advance Deduction</label>
              <Controller name="advance_deduction" control={control} render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />} />
            </div>
          </div>
        </TabsContent>

        <TabsContent value="records" className="space-y-4 mt-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Days Worked</label>
              <Controller name="days_worked" control={control} render={({ field }) => <Input {...field} type="number" min="0" max="31" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Days Absent</label>
              <Controller name="days_absent" control={control} render={({ field }) => <Input {...field} type="number" min="0" max="31" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Total Hours Worked</label>
              <Controller name="total_hours_worked" control={control} render={({ field }) => <Input {...field} type="number" step="0.5" placeholder="0" />} />
            </div>
            <div className="space-y-2">
              <label className="block text-sm font-medium">Payment Method</label>
              <Controller
                name="payment_method"
                control={control}
                render={({ field }) => (
                  <select {...field} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800">
                    <option value="bank_transfer">Bank Transfer</option>
                    <option value="cash">Cash</option>
                    <option value="cheque">Cheque</option>
                  </select>
                )}
              />
            </div>
          </div>
        </TabsContent>
      </Tabs>

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
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting} className="gap-2">
          {isSubmitting && <Loader className="h-4 w-4 animate-spin" />}
          {isSubmitting ? 'Saving...' : salary ? 'Update Salary' : 'Create Salary'}
        </Button>
      </div>
    </form>
  )
}
