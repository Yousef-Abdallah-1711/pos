'use client'

import { useState, useEffect } from 'react'
import { useAuthStore } from '@/stores/authStore'
import { usePreferencesStore } from '@/stores/preferencesStore'
import { Save, Loader2, Check } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { usePathname } from 'next/navigation'
import * as themeApi from '@/lib/api/themeApi'
import { getConfig } from '@/lib/config'

export default function AppearanceSettingsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const { settings, setSettings, setLoading, isLoading } = usePreferencesStore()

  const [appearanceForm, setAppearanceForm] = useState({
    theme: 'system' as 'light' | 'dark' | 'system',
    primaryColor: '#3B82F6',
    secondaryColor: '#6366F1',
    accentColor: '#10B981',
  })

  const [saveStatus, setSaveStatus] = useState<'idle' | 'saving' | 'success' | 'error'>('idle')
  const [errorMessage, setErrorMessage] = useState('')
  const previewColorRules = `
    .appearance-preview-primary { background-color: ${appearanceForm.primaryColor}; }
    .appearance-preview-secondary { background-color: ${appearanceForm.secondaryColor}; }
    .appearance-preview-accent { background-color: ${appearanceForm.accentColor}; }
  `

  const syncBrandColorsToActiveTheme = async () => {
    const themes = await themeApi.getThemes()
    const activeTheme = themes.find((theme) => theme.isActive)
    if (!activeTheme) return

    const mergedColors = {
      ...activeTheme.colors,
      primary: appearanceForm.primaryColor,
      secondary: appearanceForm.secondaryColor,
      accent: appearanceForm.accentColor,
    }

    const components = (activeTheme.components || []).map((component, index) => {
      const order = typeof component.displayOrder === 'number' ? component.displayOrder : index
      const config = component.config || component.settings || {}
      return {
        id: String(component.id),
        type: component.type,
        title: {
          en: config.title_en || config.title || component.title || '',
          ar: config.title_ar || '',
        },
        order,
        displayOrder: order,
        enabled: component.enabled !== false,
        config,
        settings: config,
      }
    })

    const config = {
      colors: mergedColors,
      typography: activeTheme.typography,
      identity: activeTheme.identity,
      header: activeTheme.header,
      footer: activeTheme.footer,
      components,
      customCss: activeTheme.customCss || '',
    }

    await themeApi.updateTheme(activeTheme.id, {
      name: activeTheme.name,
      slug: activeTheme.slug,
      description: activeTheme.description,
      colors: mergedColors,
      typography: activeTheme.typography,
      identity: activeTheme.identity,
      header_config: activeTheme.header,
      footer_config: activeTheme.footer,
      components: components as any,
      customCss: activeTheme.customCss || '',
      config,
    } as any)
  }

  // Initialize form when settings are loaded
  useEffect(() => {
    if (settings) {
      setAppearanceForm({
        theme: settings.theme,
        primaryColor: settings.primaryColor,
        secondaryColor: settings.secondaryColor,
        accentColor: settings.accentColor,
      })
    }
  }, [settings])

  const handleSaveTheme = async () => {
    if (!user?.id || !token) return

    setSaveStatus('saving')
    try {
      const response = await fetch(`${getConfig().apiUrl}/users/${user.id}/settings/theme`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          
        },
        body: JSON.stringify({ theme: appearanceForm.theme }),
      })

      if (response.ok) {
        const data = await response.json()
        if (settings) {
          const updatedSettings = { ...settings, theme: data.theme }
          setSettings(updatedSettings)
        }
        setSaveStatus('success')
        setTimeout(() => setSaveStatus('idle'), 2000)
      } else {
        const error = await response.json()
        setErrorMessage(error.error || t('settings.appearanceTitle'))
        setSaveStatus('error')
      }
    } catch (error) {
      setErrorMessage('Network error')
      setSaveStatus('error')
    }
  }

  const handleSaveColors = async () => {
    if (!user?.id || !token) return

    setSaveStatus('saving')
    try {
      const response = await fetch(`${getConfig().apiUrl}/users/${user.id}/settings/theme-colors`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          
        },
        body: JSON.stringify({
          primary_color: appearanceForm.primaryColor,
          secondary_color: appearanceForm.secondaryColor,
          accent_color: appearanceForm.accentColor,
        }),
      })

        if (response.ok) {
          const data = await response.json()
          if (settings) {
            const updatedSettings = {
              ...settings,
            primaryColor: data.primary_color,
            secondaryColor: data.secondary_color,
            accentColor: data.accent_color,
          }
          setSettings(updatedSettings)
          }
          try {
            await syncBrandColorsToActiveTheme()
          } catch (syncError) {
            console.warn('Failed to sync colors to active website theme', syncError)
          }
          setSaveStatus('success')
          setTimeout(() => setSaveStatus('idle'), 2000)
        } else {
        const error = await response.json()
        setErrorMessage(error.error || t('settings.appearanceTitle'))
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
      {/* Theme Section */}
      <div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-6">{t('settings.theme')}</h2>
        <div className="grid grid-cols-3 gap-4">
          {(['light', 'dark', 'system'] as const).map((theme) => (
            <button
              key={theme}
              onClick={() => setAppearanceForm({ ...appearanceForm, theme })}
              className={`p-4 border-2 rounded-lg text-center capitalize transition-colors ${
                appearanceForm.theme === theme
                  ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                  : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
              }`}
            >
              <div
                className={`w-8 h-8 mx-auto mb-2 rounded-full ${
                  theme === 'light' ? 'bg-yellow-400' : theme === 'dark' ? 'bg-gray-800' : 'bg-gradient-to-r from-yellow-400 to-gray-800'
                }`}
              />
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300">{theme}</span>
            </button>
          ))}
        </div>
        <div className="mt-6 flex justify-end">
          {renderSaveButton(handleSaveTheme, t('settings.theme'))}
        </div>
      </div>

      {/* Divider */}
      <div className="border-t border-gray-200 dark:border-gray-800 pt-8" />

      {/* Colors Section */}
      <div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-6">{t('settings.appearanceTitle')}</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('settings.primaryColor')}
            </label>
            <div className="flex items-center gap-3">
              <input
                type="color"
                value={appearanceForm.primaryColor}
                onChange={(e) => setAppearanceForm({ ...appearanceForm, primaryColor: e.target.value })}
                className="w-12 h-12 rounded-lg cursor-pointer border-0"
              />
              <input
                type="text"
                value={appearanceForm.primaryColor}
                onChange={(e) => setAppearanceForm({ ...appearanceForm, primaryColor: e.target.value })}
                className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm"
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('settings.secondaryColor')}
            </label>
            <div className="flex items-center gap-3">
              <input
                type="color"
                value={appearanceForm.secondaryColor}
                onChange={(e) => setAppearanceForm({ ...appearanceForm, secondaryColor: e.target.value })}
                className="w-12 h-12 rounded-lg cursor-pointer border-0"
              />
              <input
                type="text"
                value={appearanceForm.secondaryColor}
                onChange={(e) => setAppearanceForm({ ...appearanceForm, secondaryColor: e.target.value })}
                className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm"
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('settings.accentColor')}
            </label>
            <div className="flex items-center gap-3">
              <input
                type="color"
                value={appearanceForm.accentColor}
                onChange={(e) => setAppearanceForm({ ...appearanceForm, accentColor: e.target.value })}
                className="w-12 h-12 rounded-lg cursor-pointer border-0"
              />
              <input
                type="text"
                value={appearanceForm.accentColor}
                onChange={(e) => setAppearanceForm({ ...appearanceForm, accentColor: e.target.value })}
                className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm"
              />
            </div>
          </div>
        </div>

        {/* Color Preview */}
        <div className="mt-6 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
          <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">Preview</h3>
          <div className="flex flex-col sm:flex-row gap-4">
            <button className="appearance-preview-primary px-4 py-2 rounded-lg text-white text-sm font-medium">
              Primary
            </button>
            <button className="appearance-preview-secondary px-4 py-2 rounded-lg text-white text-sm font-medium">
              Secondary
            </button>
            <button className="appearance-preview-accent px-4 py-2 rounded-lg text-white text-sm font-medium">
              Accent
            </button>
          </div>
        </div>

        <div className="mt-6 flex justify-end">
          {renderSaveButton(handleSaveColors, t('settings.appearanceTitle'))}
        </div>
      </div>
      <style jsx>{`
        ${previewColorRules}
      `}</style>
    </div>
  )
}
