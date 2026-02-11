'use client'

import { useState, useEffect } from 'react'
import { usePathname, useRouter, useParams } from 'next/navigation'
import Link from 'next/link'
import { useAuthStore } from '@/stores/authStore'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { toast } from 'sonner'
import {
  ArrowLeft,
  Loader,
  Truck,
  Star,
  TrendingUp,
  MapPin,
  AlertCircle,
} from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'

interface Driver {
  id: number
  first_name: string
  last_name: string
  email: string
  phone_number: string
  license_number: string
  vehicle_type?: string
  vehicle_number?: string
  status: string
  availability_status: string
  total_deliveries: number
  completed_deliveries: number
  cancelled_deliveries: number
  average_rating: number
  active_orders_count: number
  is_verified: boolean
  date_of_birth?: string
  address?: string
  city?: string
  state?: string
  zip_code?: string
  notes?: string
  created_at: string
  updated_at: string
}

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

export default function EditDriverPage() {
  const params = useParams()
  const pathname = usePathname()
  const router = useRouter()
  const driverId = params.id as string
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'
  const t = createTranslator(locale)

  const [driver, setDriver] = useState<Driver | null>(null)
  const [loading, setLoading] = useState(true)
  const [updating, setUpdating] = useState(false)
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

  useEffect(() => {
    fetchDriver()
  }, [driverId, token])

  const fetchDriver = async () => {
    setLoading(true)
    try {
      const response = await api.get(`/admin/drivers/${driverId}`)
      const driverData = response.data
      setDriver(driverData)

      // Pre-fill form
      setFormData({
        first_name: driverData.first_name || '',
        last_name: driverData.last_name || '',
        email: driverData.email || '',
        phone_number: driverData.phone_number || '',
        license_number: driverData.license_number || '',
        vehicle_type: driverData.vehicle_type || '',
        vehicle_number: driverData.vehicle_number || '',
        date_of_birth: driverData.date_of_birth || '',
        address: driverData.address || '',
        city: driverData.city || '',
        state: driverData.state || '',
        zip_code: driverData.zip_code || '',
        notes: driverData.notes || '',
      })
    } catch (error: any) {
      debugError('[DRIVER_DETAILS] Error:', error)
      toast.error('Failed to load driver details')
      router.push(`/${locale}/dashboard/drivers`)
    } finally {
      setLoading(false)
    }
  }

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

    setUpdating(true)
    try {
      await api.put(`/admin/drivers/${driverId}`, formData)
      toast.success('Driver updated successfully')
      fetchDriver()
    } catch (error: any) {
      debugError('[UPDATE_DRIVER] Error:', error)
      const errorMessage = error.response?.data?.message || 'Failed to update driver'
      toast.error(errorMessage)
    } finally {
      setUpdating(false)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-96">
        <Loader className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!driver) {
    return (
      <div className="flex justify-center items-center h-96">
        <AlertCircle className="h-8 w-8 text-red-500 mr-2" />
        <span>Driver not found</span>
      </div>
    )
  }

  const completionRate =
    driver.total_deliveries > 0
      ? Math.round((driver.completed_deliveries / driver.total_deliveries) * 100)
      : 0

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
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            {driver.first_name} {driver.last_name}
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">Driver ID: {driver.id}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column - Form */}
        <div className="lg:col-span-2">
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6 md:p-8">
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Required Fields Section */}
              <div className="border-b border-gray-200 dark:border-gray-700 pb-6">
                <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
                  Personal Information
                </h2>

                <div className={`grid grid-cols-1 md:grid-cols-2 gap-6 ${isRTL ? 'rtl' : ''}`}>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      First Name
                    </label>
                    <input
                      type="text"
                      name="first_name"
                      value={formData.first_name}
                      onChange={handleChange}
                      className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white ${
                        errors.first_name
                          ? 'border-red-500 focus:ring-red-500'
                          : 'border-gray-300'
                      }`}
                    />
                    {errors.first_name && (
                      <p className="text-red-500 text-sm mt-1">{errors.first_name}</p>
                    )}
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Last Name
                    </label>
                    <input
                      type="text"
                      name="last_name"
                      value={formData.last_name}
                      onChange={handleChange}
                      className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white ${
                        errors.last_name
                          ? 'border-red-500 focus:ring-red-500'
                          : 'border-gray-300'
                      }`}
                    />
                    {errors.last_name && (
                      <p className="text-red-500 text-sm mt-1">{errors.last_name}</p>
                    )}
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Email
                    </label>
                    <input
                      type="email"
                      name="email"
                      value={formData.email}
                      onChange={handleChange}
                      className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white ${
                        errors.email ? 'border-red-500 focus:ring-red-500' : 'border-gray-300'
                      }`}
                    />
                    {errors.email && (
                      <p className="text-red-500 text-sm mt-1">{errors.email}</p>
                    )}
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Phone Number
                    </label>
                    <input
                      type="tel"
                      name="phone_number"
                      value={formData.phone_number}
                      onChange={handleChange}
                      className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white ${
                        errors.phone_number
                          ? 'border-red-500 focus:ring-red-500'
                          : 'border-gray-300'
                      }`}
                    />
                    {errors.phone_number && (
                      <p className="text-red-500 text-sm mt-1">{errors.phone_number}</p>
                    )}
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Date of Birth
                    </label>
                    <input
                      type="date"
                      name="date_of_birth"
                      value={formData.date_of_birth}
                      onChange={handleChange}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      License Number
                    </label>
                    <input
                      type="text"
                      name="license_number"
                      value={formData.license_number}
                      onChange={handleChange}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    />
                  </div>
                </div>
              </div>

              {/* Vehicle Information */}
              <div className="border-b border-gray-200 dark:border-gray-700 pb-6">
                <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
                  Vehicle Information
                </h2>

                <div className={`grid grid-cols-1 md:grid-cols-2 gap-6 ${isRTL ? 'rtl' : ''}`}>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Vehicle Type
                    </label>
                    <select
                      name="vehicle_type"
                      value={formData.vehicle_type}
                      onChange={handleChange}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    >
                      <option value="">Select vehicle type</option>
                      <option value="car">Car</option>
                      <option value="bike">Bike</option>
                      <option value="van">Van</option>
                      <option value="truck">Truck</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Vehicle Number
                    </label>
                    <input
                      type="text"
                      name="vehicle_number"
                      value={formData.vehicle_number}
                      onChange={handleChange}
                      placeholder="License plate"
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    />
                  </div>
                </div>
              </div>

              {/* Address Information */}
              <div className="border-b border-gray-200 dark:border-gray-700 pb-6">
                <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
                  Address Information
                </h2>

                <div className={`grid grid-cols-1 md:grid-cols-2 gap-6 ${isRTL ? 'rtl' : ''}`}>
                  <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Street Address
                    </label>
                    <input
                      type="text"
                      name="address"
                      value={formData.address}
                      onChange={handleChange}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      City
                    </label>
                    <input
                      type="text"
                      name="city"
                      value={formData.city}
                      onChange={handleChange}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      State/Province
                    </label>
                    <input
                      type="text"
                      name="state"
                      value={formData.state}
                      onChange={handleChange}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Zip Code
                    </label>
                    <input
                      type="text"
                      name="zip_code"
                      value={formData.zip_code}
                      onChange={handleChange}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    />
                  </div>
                </div>
              </div>

              {/* Notes */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Notes
                </label>
                <textarea
                  name="notes"
                  value={formData.notes}
                  onChange={handleChange}
                  rows={4}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                />
              </div>

              {/* Action Buttons */}
              <div className={`flex gap-4 pt-6 border-t border-gray-200 dark:border-gray-700 ${isRTL ? 'flex-row-reverse' : ''}`}>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => router.push(`/${locale}/dashboard/drivers`)}
                  disabled={updating}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={updating}
                  className="flex items-center gap-2"
                >
                  {updating && <Loader className="h-4 w-4 animate-spin" />}
                  {updating ? 'Updating...' : 'Update Driver'}
                </Button>
              </div>
            </form>
          </div>
        </div>

        {/* Right Column - Statistics */}
        <div className="space-y-6">
          {/* Status Badge */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
              Status
            </h3>
            <div className="space-y-3">
              <div>
                <p className="text-sm text-gray-600 dark:text-gray-400">Account Status</p>
                <span
                  className={`px-3 py-1 rounded-full text-sm font-semibold inline-block mt-1 ${
                    driver.status === 'active'
                      ? 'bg-green-100 text-green-800'
                      : 'bg-gray-100 text-gray-800'
                  }`}
                >
                  {driver.status.toUpperCase()}
                </span>
              </div>
              <div>
                <p className="text-sm text-gray-600 dark:text-gray-400">Availability</p>
                <span
                  className={`px-3 py-1 rounded-full text-sm font-semibold inline-block mt-1 ${
                    driver.availability_status === 'available'
                      ? 'bg-blue-100 text-blue-800'
                      : driver.availability_status === 'busy'
                      ? 'bg-yellow-100 text-yellow-800'
                      : 'bg-gray-100 text-gray-800'
                  }`}
                >
                  {driver.availability_status.replace('_', ' ').toUpperCase()}
                </span>
              </div>
              <div>
                <p className="text-sm text-gray-600 dark:text-gray-400">Verified</p>
                <p className="text-sm font-medium text-gray-900 dark:text-white mt-1">
                  {driver.is_verified ? '✓ Yes' : '✗ No'}
                </p>
              </div>
            </div>
          </div>

          {/* Performance Stats */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6 space-y-4">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Performance
            </h3>

            <div className="space-y-3">
              <div className="flex items-center gap-3">
                <Truck className="h-5 w-5 text-gray-400" />
                <div className="flex-1">
                  <p className="text-sm text-gray-600 dark:text-gray-400">Total Deliveries</p>
                  <p className="text-2xl font-bold text-gray-900 dark:text-white">
                    {driver.total_deliveries}
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <div className="h-5 w-5 text-green-500">✓</div>
                <div className="flex-1">
                  <p className="text-sm text-gray-600 dark:text-gray-400">Completed</p>
                  <p className="text-2xl font-bold text-green-600">
                    {driver.completed_deliveries}
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <div className="h-5 w-5 text-red-500">✕</div>
                <div className="flex-1">
                  <p className="text-sm text-gray-600 dark:text-gray-400">Cancelled</p>
                  <p className="text-2xl font-bold text-red-600">
                    {driver.cancelled_deliveries}
                  </p>
                </div>
              </div>

              <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">
                  Completion Rate
                </p>
                <div className="flex items-center gap-2">
                  <div className="flex-1 h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-green-500 transition-all"
                      style={{ width: `${completionRate}%` }}
                    ></div>
                  </div>
                  <span className="text-sm font-semibold text-gray-900 dark:text-white">
                    {completionRate}%
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Rating */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
              <Star className="h-5 w-5 text-yellow-500" />
              Rating
            </h3>
            <div className="mt-4">
              <p className="text-4xl font-bold text-yellow-500">
                {driver.average_rating.toFixed(1)}
              </p>
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
                Out of 5.0
              </p>
            </div>
          </div>

          {/* Active Orders */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
              <MapPin className="h-5 w-5 text-primary" />
              Active Orders
            </h3>
            <p className="text-4xl font-bold text-primary mt-4">
              {driver.active_orders_count}
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
