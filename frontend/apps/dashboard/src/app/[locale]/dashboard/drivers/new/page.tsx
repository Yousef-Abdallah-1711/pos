'use client'

import { useState } from 'react'
import { usePathname, useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuthStore } from '@/stores/authStore'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { toast } from 'sonner'
import { ArrowLeft, Loader } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'

interface FormData {
  first_name: string
  last_name: string
  email: string
  phone_number: string
  license_number: string
  vehicle_type: string
  vehicle_number: string
  date_of_birth?: string
  address?: string
  city?: string
  state?: string
  zip_code?: string
  notes?: string
}

interface FormErrors {
  [key: string]: string
}

export default function AddDriverPage() {
  const pathname = usePathname()
  const router = useRouter()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'
  const t = createTranslator(locale)

  const [loading, setLoading] = useState(false)
  const [errors, setErrors] = useState<FormErrors>({})
  const token = useAuthStore((state) => state.token)

  const [formData, setFormData] = useState<FormData>({
    first_name: '',
    last_name: '',
    email: '',
    phone_number: '',
    license_number: '',
    vehicle_type: '',
    vehicle_number: '',
    date_of_birth: '',
    address: '',
    city: '',
    state: '',
    zip_code: '',
    notes: '',
  })

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {}

    if (!formData.first_name.trim()) {
      newErrors.first_name = 'First name is required'
    }
    if (!formData.last_name.trim()) {
      newErrors.last_name = 'Last name is required'
    }
    if (!formData.email.trim()) {
      newErrors.email = 'Email is required'
    } else if (!formData.email.includes('@')) {
      newErrors.email = 'Invalid email format'
    }
    if (!formData.phone_number.trim()) {
      newErrors.phone_number = 'Phone number is required'
    }
    if (!formData.license_number.trim()) {
      newErrors.license_number = 'License number is required'
    }
    if (!formData.vehicle_type.trim()) {
      newErrors.vehicle_type = 'Vehicle type is required'
    }
    if (!formData.vehicle_number.trim()) {
      newErrors.vehicle_number = 'Vehicle number is required'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))
    // Clear error for this field when user starts typing
    if (errors[name]) {
      setErrors((prev) => ({
        ...prev,
        [name]: '',
      }))
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      toast.error('Please fill in all required fields')
      return
    }

    setLoading(true)
    try {
      await api.post('/admin/drivers', formData)
      toast.success('Driver added successfully')
      router.push(`/${locale}/dashboard/drivers`)
    } catch (error: any) {
      debugError('[ADD_DRIVER] Error:', error)
      const errorMessage = error.response?.data?.message || 'Failed to add driver'
      toast.error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = () => {
    router.push(`/${locale}/dashboard/drivers`)
  }

  const formFields = [
    {
      name: 'first_name',
      label: 'First Name',
      type: 'text',
      required: true,
      placeholder: 'Enter first name',
    },
    {
      name: 'last_name',
      label: 'Last Name',
      type: 'text',
      required: true,
      placeholder: 'Enter last name',
    },
    {
      name: 'email',
      label: 'Email',
      type: 'email',
      required: true,
      placeholder: 'Enter email address',
    },
    {
      name: 'phone_number',
      label: 'Phone Number',
      type: 'tel',
      required: true,
      placeholder: 'Enter phone number',
    },
    {
      name: 'license_number',
      label: 'License Number',
      type: 'text',
      required: true,
      placeholder: 'Enter license number',
    },
    {
      name: 'vehicle_type',
      label: 'Vehicle Type',
      type: 'select',
      required: true,
      options: ['car', 'bike', 'van', 'truck'],
    },
    {
      name: 'vehicle_number',
      label: 'Vehicle Number',
      type: 'text',
      required: true,
      placeholder: 'Enter vehicle license plate',
    },
    {
      name: 'date_of_birth',
      label: 'Date of Birth',
      type: 'date',
      required: false,
    },
    {
      name: 'address',
      label: 'Address',
      type: 'text',
      required: false,
      placeholder: 'Enter street address',
    },
    {
      name: 'city',
      label: 'City',
      type: 'text',
      required: false,
      placeholder: 'Enter city',
    },
    {
      name: 'state',
      label: 'State/Province',
      type: 'text',
      required: false,
      placeholder: 'Enter state',
    },
    {
      name: 'zip_code',
      label: 'Zip Code',
      type: 'text',
      required: false,
      placeholder: 'Enter zip code',
    },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Link href={`/${locale}/dashboard/drivers`}>
          <Button variant="outline" size="sm">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back
          </Button>
        </Link>
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Add New Driver</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Create a new driver account
          </p>
        </div>
      </div>

      {/* Form */}
      <div className="max-w-2xl bg-white dark:bg-gray-800 rounded-lg shadow p-6 md:p-8">
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Required Fields Section */}
          <div className="border-b border-gray-200 dark:border-gray-700 pb-6">
            <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
              Required Information
            </h2>

            <div className={`grid grid-cols-1 md:grid-cols-2 gap-6 ${isRTL ? 'rtl' : ''}`}>
              {formFields
                .filter((field) => field.required)
                .map((field) => (
                  <div key={field.name}>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      {field.label}
                      <span className="text-red-500 ml-1">*</span>
                    </label>

                    {field.type === 'select' ? (
                      <select
                        name={field.name}
                        value={formData[field.name as keyof FormData] || ''}
                        onChange={handleChange}
                        className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white ${
                          errors[field.name]
                            ? 'border-red-500 focus:ring-red-500'
                            : 'border-gray-300'
                        }`}
                      >
                        <option value="">Select {field.label.toLowerCase()}</option>
                        {(field.options || []).map((option) => (
                          <option key={option} value={option}>
                            {option.charAt(0).toUpperCase() + option.slice(1)}
                          </option>
                        ))}
                      </select>
                    ) : (
                      <input
                        type={field.type}
                        name={field.name}
                        value={formData[field.name as keyof FormData] || ''}
                        onChange={handleChange}
                        placeholder={field.placeholder}
                        className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white ${
                          errors[field.name]
                            ? 'border-red-500 focus:ring-red-500'
                            : 'border-gray-300'
                        }`}
                      />
                    )}

                    {errors[field.name] && (
                      <p className="text-red-500 text-sm mt-1">{errors[field.name]}</p>
                    )}
                  </div>
                ))}
            </div>
          </div>

          {/* Optional Fields Section */}
          <div>
            <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
              Additional Information
            </h2>

            <div className={`grid grid-cols-1 md:grid-cols-2 gap-6 ${isRTL ? 'rtl' : ''}`}>
              {formFields
                .filter((field) => !field.required)
                .map((field) => (
                  <div key={field.name}>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      {field.label}
                    </label>

                    <input
                      type={field.type}
                      name={field.name}
                      value={formData[field.name as keyof FormData] || ''}
                      onChange={handleChange}
                      placeholder={field.placeholder}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    />
                  </div>
                ))}

              {/* Notes Field - Full Width */}
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Notes
                </label>
                <textarea
                  name="notes"
                  value={formData.notes || ''}
                  onChange={handleChange}
                  placeholder="Enter any additional notes about the driver"
                  rows={4}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                />
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className={`flex gap-4 pt-6 border-t border-gray-200 dark:border-gray-700 ${isRTL ? 'flex-row-reverse' : ''}`}>
            <Button
              type="button"
              variant="outline"
              onClick={handleCancel}
              disabled={loading}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={loading}
              className="flex items-center gap-2"
            >
              {loading && <Loader className="h-4 w-4 animate-spin" />}
              {loading ? 'Adding Driver...' : 'Add Driver'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
