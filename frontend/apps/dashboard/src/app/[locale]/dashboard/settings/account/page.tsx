'use client'

import { useState, useEffect } from 'react'
import { useAuthStore } from '@/stores/authStore'
import { Save, Loader2, Check, X } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { usePathname } from 'next/navigation'
import { getConfig } from '@/lib/config'

export default function AccountSettingsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)

  // Form states
  const [profileForm, setProfileForm] = useState({ name: '', avatarUrl: '' })
  const [passwordForm, setPasswordForm] = useState({ currentPassword: '', newPassword: '', confirmPassword: '' })

  // Status states
  const [saveStatus, setSaveStatus] = useState<'idle' | 'saving' | 'success' | 'error'>('idle')
  const [errorMessage, setErrorMessage] = useState('')
  const [isLoading, setIsLoading] = useState(true)

  // Initialize forms when user is loaded
  useEffect(() => {
    if (user) {
      setProfileForm({ name: user.name || '', avatarUrl: '' })
      setIsLoading(false)
    }
  }, [user])

  const handleSaveProfile = async () => {
    if (!user?.id || !token) return

    setSaveStatus('saving')
    try {
      const { apiUrl } = getConfig()
      const response = await fetch(`${apiUrl}/users/${user.id}/profile`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          
        },
        body: JSON.stringify({
          name: profileForm.name,
          avatar_url: profileForm.avatarUrl || null,
        }),
      })

      if (response.ok) {
        setSaveStatus('success')
        setTimeout(() => setSaveStatus('idle'), 2000)
      } else {
        const error = await response.json()
        setErrorMessage(error.error || t('settings.updateProfile'))
        setSaveStatus('error')
      }
    } catch (error) {
      setErrorMessage('Network error')
      setSaveStatus('error')
    }
  }

  const handleChangePassword = async () => {
    if (!user?.id || !token) return

    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      setErrorMessage(t('settings.passwordMismatch'))
      setSaveStatus('error')
      return
    }

    if (passwordForm.newPassword.length < 8) {
      setErrorMessage(t('settings.passwordTooShort'))
      setSaveStatus('error')
      return
    }

    setSaveStatus('saving')
    try {
      const { apiUrl } = getConfig()
      const response = await fetch(`${apiUrl}/users/${user.id}/settings/change-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          
        },
        body: JSON.stringify({
          current_password: passwordForm.currentPassword,
          new_password: passwordForm.newPassword,
          confirm_password: passwordForm.confirmPassword,
        }),
      })

      if (response.ok) {
        setSaveStatus('success')
        setPasswordForm({ currentPassword: '', newPassword: '', confirmPassword: '' })
        setTimeout(() => setSaveStatus('idle'), 2000)
      } else {
        const error = await response.json()
        setErrorMessage(error.error || t('settings.changePassword'))
        setSaveStatus('error')
      }
    } catch (error) {
      setErrorMessage('Network error')
      setSaveStatus('error')
    }
  }

  const renderSaveButton = (onClick: () => void, label: string = t('settings.save')) => (
    <button
      onClick={onClick}
      disabled={saveStatus === 'saving'}
      className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
    >
      {saveStatus === 'saving' ? (
        <Loader2 className="w-4 h-4 animate-spin" />
      ) : saveStatus === 'success' ? (
        <Check className="w-4 h-4" />
      ) : (
        <Save className="w-4 h-4" />
      )}
      {saveStatus === 'saving' ? t('settings.saving') : saveStatus === 'success' ? t('settings.success') : label}
    </button>
  )

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Profile Section */}
      <div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-6">{t('settings.profileTitle')}</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('settings.name')}
            </label>
            <input
              type="text"
              value={profileForm.name}
              onChange={(e) => setProfileForm({ ...profileForm, name: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('settings.email')}
            </label>
            <input
              type="email"
              value={user?.email || ''}
              disabled
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 cursor-not-allowed"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('settings.role')}
            </label>
            <input
              type="text"
              value={user?.role || ''}
              disabled
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 cursor-not-allowed capitalize"
            />
          </div>
        </div>
        <div className="flex justify-end mt-6">
          {renderSaveButton(handleSaveProfile, t('settings.updateProfile'))}
        </div>
      </div>

      {/* Divider */}
      <div className="border-t border-gray-200 dark:border-gray-800 pt-8" />

      {/* Security Section */}
      <div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-6">{t('settings.securityTitle')}</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('settings.currentPassword')}
            </label>
            <input
              type="password"
              value={passwordForm.currentPassword}
              onChange={(e) => setPasswordForm({ ...passwordForm, currentPassword: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('settings.newPassword')}
            </label>
            <input
              type="password"
              value={passwordForm.newPassword}
              onChange={(e) => setPasswordForm({ ...passwordForm, newPassword: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {t('settings.passwordTooShort')}
            </p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('settings.confirmPassword')}
            </label>
            <input
              type="password"
              value={passwordForm.confirmPassword}
              onChange={(e) => setPasswordForm({ ...passwordForm, confirmPassword: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
          </div>
        </div>
        {saveStatus === 'error' && (
          <div className="mt-4 flex items-center gap-2 p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-lg">
            <X className="w-4 h-4" />
            {errorMessage}
          </div>
        )}
        <div className="flex justify-end mt-6">
          {renderSaveButton(handleChangePassword, t('settings.changePassword'))}
        </div>
      </div>
    </div>
  )
}
