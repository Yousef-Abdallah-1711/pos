'use client'

import { useState, useEffect, useMemo } from 'react'
import { useAuthStore } from '@/stores/authStore'
import { api } from '@/lib/api'
import { Save, Loader2, Check, X, Store, Mail, Phone } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { usePathname } from 'next/navigation'
import { toast } from 'sonner'

export default function RestaurantSettingsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = useMemo(() => createTranslator(locale), [locale])
  const token = useAuthStore((state) => state.token)

  const [form, setForm] = useState({ email: '', phone: '', name: '' })
  const [saveStatus, setSaveStatus] = useState<'idle' | 'saving' | 'success' | 'error'>('idle')
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    if (!token) return

    const fetchContactInfo = async () => {
      try {
        const response = await api.get('/tenant/contact-info')
        const data = response.data?.data || response.data
        setForm({
          email: data.email || '',
          phone: data.phone || '',
          name: data.name || '',
        })
      } catch (error) {
        console.error('Failed to load contact info:', error)
      } finally {
        setIsLoading(false)
      }
    }

    fetchContactInfo()
  }, [token])

  const handleSave = async () => {
    setSaveStatus('saving')
    try {
      await api.put('/tenant/contact-info', {
        email: form.email,
        phone: form.phone,
      })
      setSaveStatus('success')
      toast.success('Contact information updated successfully')
      setTimeout(() => setSaveStatus('idle'), 2000)
    } catch (error) {
      setSaveStatus('error')
      toast.error('Failed to update contact information')
      setTimeout(() => setSaveStatus('idle'), 3000)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">
          Restaurant Settings
        </h1>
        <p className="text-gray-500 dark:text-gray-400 mt-1">
          Manage your restaurant contact information. This info is shown to customers when your site is under maintenance.
        </p>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6 space-y-6">
        <div className="flex items-center gap-3 mb-2">
          <Store className="h-5 w-5 text-primary-600" />
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            Contact Information
          </h2>
        </div>

        <p className="text-sm text-gray-500 dark:text-gray-400">
          This contact information will be displayed on your maintenance page so customers can reach you when your online ordering is unavailable.
        </p>

        <div className="space-y-4">
          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Restaurant Name
            </label>
            <div className="relative">
              <Store className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
              <Input
                value={form.name}
                disabled
                className="pl-10 bg-gray-50 dark:bg-gray-700 cursor-not-allowed"
              />
            </div>
            <p className="text-xs text-gray-400">Restaurant name is managed by your system administrator.</p>
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Contact Email
            </label>
            <div className="relative">
              <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
              <Input
                type="email"
                value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })}
                placeholder="contact@yourrestaurant.com"
                className="pl-10"
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Contact Phone
            </label>
            <div className="relative">
              <Phone className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
              <Input
                type="tel"
                value={form.phone}
                onChange={(e) => setForm({ ...form, phone: e.target.value })}
                placeholder="+1 (555) 123-4567"
                className="pl-10"
              />
            </div>
          </div>
        </div>

        <div className="pt-4 border-t border-gray-200 dark:border-gray-700">
          <Button
            onClick={handleSave}
            disabled={saveStatus === 'saving'}
            className="gap-2"
          >
            {saveStatus === 'saving' ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : saveStatus === 'success' ? (
              <Check className="h-4 w-4" />
            ) : saveStatus === 'error' ? (
              <X className="h-4 w-4" />
            ) : (
              <Save className="h-4 w-4" />
            )}
            {saveStatus === 'saving' ? 'Saving...' : saveStatus === 'success' ? 'Saved!' : 'Save Changes'}
          </Button>
        </div>
      </div>
    </div>
  )
}
