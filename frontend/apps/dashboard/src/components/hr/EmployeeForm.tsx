'use client'

import { useState } from 'react'
import { usePathname } from 'next/navigation'
import { useForm, Controller } from 'react-hook-form'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { toast } from 'sonner'
import { debugLog } from '@/lib/debug'
import { AlertCircle, CheckCircle, Loader } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'

interface EmployeeFormProps {
  employee?: any
  onSubmit: (formData: any) => Promise<void>
  onCancel: () => void
  existingEmails?: string[]
}

export function EmployeeForm({ employee, onSubmit, onCancel, existingEmails = [] }: EmployeeFormProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)

  const { control, handleSubmit, formState: { errors }, watch } = useForm({
    defaultValues: {
      first_name: employee?.first_name || '',
      first_name_ar: employee?.first_name_ar || '',
      last_name: employee?.last_name || '',
      last_name_ar: employee?.last_name_ar || '',
      email: employee?.email || '',
      phone: employee?.phone || '',
      employee_code: employee?.employee_code || '',
      position: employee?.position || '',
      position_ar: employee?.position_ar || '',
      department: employee?.department || '',
      department_ar: employee?.department_ar || '',
      hire_date: employee?.hire_date || '',
      status: employee?.status || 'active',
      employment_type: employee?.employment_type || 'full_time',
      salary_currency: employee?.salary_currency || 'EGP',
      monthly_salary: employee?.monthly_salary || '',
      weekly_hours: employee?.weekly_hours || '40',
      address: employee?.address || '',
      address_ar: employee?.address_ar || '',
      city: employee?.city || '',
      country: employee?.country || '',
      emergency_contact_name: employee?.emergency_contact_name || '',
      emergency_contact_phone: employee?.emergency_contact_phone || '',
      is_active: employee?.is_active !== false,
    },
  })

  const onFormSubmit = async (data: any) => {
    setIsSubmitting(true)
    setApiError(null)

    try {
      // Validate required fields - minimal validation
      if (!data.first_name?.trim()) {
        throw new Error('First name is required')
      }
      if (!data.email?.trim()) {
        throw new Error('Email is required')
      }

      const cleanData = {
        ...data,
        monthly_salary: data.monthly_salary ? parseFloat(data.monthly_salary.toString()) : 0,
        weekly_hours: parseInt((data.weekly_hours || 40).toString()),
        hire_date: data.hire_date ? new Date(data.hire_date).toISOString() : new Date().toISOString(),
      }

      // DEBUG: Log the request data
      debugLog('[CREATE_EMPLOYEE] DEBUG - Submitting employee data')
      debugLog('Raw form data:', data)
      debugLog('Cleaned data being sent:', cleanData)
      debugLog('Timestamp:', new Date().toISOString())
      // groupEnd (debug)

      await onSubmit(cleanData)

      debugLog('[CREATE_EMPLOYEE] SUCCESS - Employee created successfully')
    } catch (error: any) {
      const errorMsg = error?.response?.data?.message || error?.message || 'Failed to save employee'

      // DEBUG: Detailed error logging
      debugLog('[CREATE_EMPLOYEE] ERROR - Detailed Debug Info')
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
            <h3 className="font-semibold text-red-900 dark:text-red-300">Error Saving Employee</h3>
            <p className="text-red-700 dark:text-red-400 text-sm mt-1">{apiError}</p>
          </div>
        </div>
      )}

      {/* Tabs */}
      <Tabs defaultValue="personal" className="w-full">
        <TabsList className="grid w-full grid-cols-4 bg-gray-100 dark:bg-gray-800 p-1">
          <TabsTrigger value="personal">Personal</TabsTrigger>
          <TabsTrigger value="employment">Employment</TabsTrigger>
          <TabsTrigger value="compensation">Compensation</TabsTrigger>
          <TabsTrigger value="contact">Contact</TabsTrigger>
        </TabsList>

        {/* Personal Tab */}
        <TabsContent value="personal" className="space-y-4 mt-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">First Name (English) *</label>
              <Controller
                name="first_name"
                control={control}
                rules={{ required: 'First name is required' }}
                render={({ field }) => (
                  <Input {...field} placeholder="e.g., Ahmed" className={errors.first_name ? 'border-red-500' : ''} />
                )}
              />
              {errors.first_name && (
                <span className="text-red-500 text-sm flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.first_name.message}
                </span>
              )}
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">First Name (Arabic)</label>
              <Controller
                name="first_name_ar"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder={t('hr.employees.firstNameArPlaceholder')} dir="rtl" />
                )}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Last Name (English)</label>
              <Controller
                name="last_name"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder="e.g., Mohamed" />
                )}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Last Name (Arabic)</label>
              <Controller
                name="last_name_ar"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder={t('hr.employees.lastNameArPlaceholder')} dir="rtl" />
                )}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Email *</label>
              <Controller
                name="email"
                control={control}
                rules={{
                  required: 'Email is required',
                  pattern: {
                    value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                    message: 'Invalid email address',
                  },
                  validate: (value) => {
                    const emailLower = value.toLowerCase()
                    const originalEmail = employee?.email?.toLowerCase()
                    // Allow the same email if editing (comparing against original)
                    if (employee && originalEmail === emailLower) {
                      return true
                    }
                    // Check against existing emails (excluding current employee's email when editing)
                    if (existingEmails.some((email) => email.toLowerCase() === emailLower)) {
                      return 'This email is already in use'
                    }
                    return true
                  },
                }}
                render={({ field }) => (
                  <Input {...field} type="email" placeholder="email@example.com" className={errors.email ? 'border-red-500' : ''} />
                )}
              />
              {errors.email && (
                <span className="text-red-500 text-sm flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.email.message}
                </span>
              )}
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Phone</label>
              <Controller
                name="phone"
                control={control}
                render={({ field }) => (
                  <Input {...field} type="tel" placeholder="+20123456789" />
                )}
              />
            </div>
          </div>
        </TabsContent>

        {/* Employment Tab */}
        <TabsContent value="employment" className="space-y-4 mt-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Employee Code (Auto-generated)</label>
              <Controller
                name="employee_code"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder="Leave blank for auto-generation" disabled />
                )}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Hire Date</label>
              <Controller
                name="hire_date"
                control={control}
                render={({ field }) => (
                  <Input {...field} type="date" />
                )}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Position (English)</label>
              <Controller
                name="position"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder="e.g., Chef (defaults to Staff)" />
                )}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Position (Arabic)</label>
              <Controller
                name="position_ar"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder={t('hr.employees.positionArPlaceholder')} dir="rtl" />
                )}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Department (English)</label>
              <Controller
                name="department"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder="e.g., Kitchen" />
                )}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Department (Arabic)</label>
              <Controller
                name="department_ar"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder={t('hr.employees.departmentArPlaceholder')} dir="rtl" />
                )}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Employment Type</label>
              <Controller
                name="employment_type"
                control={control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="full_time">Full Time</SelectItem>
                      <SelectItem value="part_time">Part Time</SelectItem>
                      <SelectItem value="contract">Contract</SelectItem>
                      <SelectItem value="seasonal">Seasonal</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Status</label>
              <Controller
                name="status"
                control={control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="active">Active</SelectItem>
                      <SelectItem value="inactive">Inactive</SelectItem>
                      <SelectItem value="on_leave">On Leave</SelectItem>
                      <SelectItem value="terminated">Terminated</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </div>
          </div>

          <div className="flex items-center space-x-2">
            <Controller
              name="is_active"
              control={control}
              render={({ field }) => (
                <Checkbox checked={field.value} onCheckedChange={field.onChange} />
              )}
            />
            <label className="text-sm font-medium">Active</label>
          </div>
        </TabsContent>

        {/* Compensation Tab */}
        <TabsContent value="compensation" className="space-y-4 mt-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Monthly Salary</label>
              <Controller
                name="monthly_salary"
                control={control}
                render={({ field }) => (
                  <Input {...field} type="number" step="0.01" placeholder="0.00" />
                )}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Currency</label>
              <Controller
                name="salary_currency"
                control={control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="EGP">Egyptian Pound (EGP)</SelectItem>
                      <SelectItem value="USD">US Dollar (USD)</SelectItem>
                      <SelectItem value="EUR">Euro (EUR)</SelectItem>
                      <SelectItem value="SAR">Saudi Riyal (SAR)</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium">Weekly Hours</label>
            <Controller
              name="weekly_hours"
              control={control}
              render={({ field }) => (
                <Input {...field} type="number" step="0.5" min="0" placeholder="40" />
              )}
            />
          </div>
        </TabsContent>

        {/* Contact Tab */}
        <TabsContent value="contact" className="space-y-4 mt-6">
          <div className="space-y-2">
            <label className="block text-sm font-medium">Address (English)</label>
            <Controller
              name="address"
              control={control}
              render={({ field }) => (
                <Textarea {...field} placeholder="Street address..." rows={3} />
              )}
            />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium">Address (Arabic)</label>
            <Controller
              name="address_ar"
              control={control}
              render={({ field }) => (
                <Textarea {...field} placeholder={t('hr.employees.addressArPlaceholder')} rows={3} dir="rtl" />
              )}
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">City</label>
              <Controller
                name="city"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder="e.g., Cairo" />
                )}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Country</label>
              <Controller
                name="country"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder="e.g., Egypt" />
                )}
              />
            </div>
          </div>

          <div className="border-t pt-4">
            <h3 className="font-semibold mb-4">Emergency Contact</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="block text-sm font-medium">Name</label>
                <Controller
                  name="emergency_contact_name"
                  control={control}
                  render={({ field }) => (
                    <Input {...field} placeholder="Contact name" />
                  )}
                />
              </div>

              <div className="space-y-2">
                <label className="block text-sm font-medium">Phone</label>
                <Controller
                  name="emergency_contact_phone"
                  control={control}
                  render={({ field }) => (
                    <Input {...field} type="tel" placeholder="+20123456789" />
                  )}
                />
              </div>
            </div>
          </div>
        </TabsContent>
      </Tabs>

      {/* Form Actions */}
      <div className="flex justify-end gap-2 pt-4 border-t">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting} className="gap-2">
          {isSubmitting && <Loader className="h-4 w-4 animate-spin" />}
          {isSubmitting ? 'Saving...' : employee ? 'Update Employee' : 'Create Employee'}
        </Button>
      </div>
    </form>
  )
}
