'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@pos-saas/ui'
import { Check, Palette } from 'lucide-react'
import { createTranslator, LOCALES } from '@/lib/translations'
import { usePreferencesStore } from '@/lib/store/preferences-store'
import Link from 'next/link'
import type { ServerThemeData } from '@/lib/api/get-theme-server'

const THEME_OPTIONS = [
  { value: 'light', labelKey: 'settings.light' },
  { value: 'dark', labelKey: 'settings.dark' },
  { value: 'system', labelKey: 'settings.system' },
] as const

const COLOR_PRESETS = [
  { name: 'Blue', primary: '#3b82f6', secondary: '#10b981', accent: '#f59e0b' },
  { name: 'Indigo', primary: '#6366f1', secondary: '#8b5cf6', accent: '#ec4899' },
  { name: 'Rose', primary: '#e11d48', secondary: '#d946ef', accent: '#f59e0b' },
  { name: 'Green', primary: '#16a34a', secondary: '#0891b2', accent: '#f59e0b' },
]

interface SettingsPageContentProps {
  locale: 'en' | 'ar'
  themeData: ServerThemeData
}

export function SettingsPageContent({ locale, themeData }: SettingsPageContentProps) {
  const router = useRouter()
  const isRTL = locale === 'ar'
  const t = createTranslator(locale)

  const themePrimaryColor = themeData?.colors?.primary || '#f97316'
  const themeSecondaryColor = themeData?.colors?.secondary || '#0ea5e9'

  const { settings, updateLanguage, updateTheme, updateColors } = usePreferencesStore()
  const [saved, setSaved] = useState(false)
  const [primaryColor, setPrimaryColor] = useState(settings.primaryColor)
  const [secondaryColor, setSecondaryColor] = useState(settings.secondaryColor)
  const [accentColor, setAccentColor] = useState(settings.accentColor)

  useEffect(() => {
    if (saved) {
      const timer = setTimeout(() => setSaved(false), 3000)
      return () => clearTimeout(timer)
    }
  }, [saved])

  const handleLanguageChange = (newLocale: string) => {
    updateLanguage(newLocale as 'en' | 'ar')
    router.push(`/${newLocale}/settings`)
  }

  const handleThemeChange = (newTheme: 'light' | 'dark' | 'system') => {
    updateTheme(newTheme)
  }

  const handleColorPreset = (preset: (typeof COLOR_PRESETS)[0]) => {
    updateColors(preset.primary, preset.secondary, preset.accent)
    setPrimaryColor(preset.primary)
    setSecondaryColor(preset.secondary)
    setAccentColor(preset.accent)
  }

  const handleCustomColors = () => {
    updateColors(primaryColor, secondaryColor, accentColor)
    setSaved(true)
  }

  const presetStyles = COLOR_PRESETS.map(
    (preset, idx) => `
      .settings-preset-primary-${idx} { background-color: ${preset.primary}; }
      .settings-preset-secondary-${idx} { background-color: ${preset.secondary}; }
      .settings-preset-accent-${idx} { background-color: ${preset.accent}; }
    `
  ).join('\n')

  return (
    <>
      <div className="settings-page-header text-white py-12">
        <div className={`max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 ${isRTL ? 'text-right' : ''}`}>
          <h1 className="text-4xl font-bold mb-2">{t('settings.title')}</h1>
          <p className="text-gray-100">{t('settings.preferences')}</p>
        </div>
      </div>

      <div className={`max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 ${isRTL ? 'text-right' : ''}`}>
        <div className="space-y-8">
          {saved && (
            <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4 flex items-center gap-3">
              <Check className="h-5 w-5 text-green-600 dark:text-green-400 flex-shrink-0" />
              <p className="text-green-800 dark:text-green-100">{t('settings.saved')}</p>
            </div>
          )}

          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-2">
              <span>{t('settings.language')}</span>
            </h2>
            <p className="text-gray-600 dark:text-gray-400 mb-4">{t('settings.selectLanguage')}</p>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {LOCALES.map((loc) => (
                <button
                  key={loc}
                  onClick={() => handleLanguageChange(loc)}
                  className={`settings-selectable p-4 rounded-lg border-2 transition-all ${
                    locale === loc
                      ? 'is-selected bg-primary/10'
                      : 'border-gray-300 dark:border-gray-600 hover:border-primary'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className="font-semibold">{loc === 'en' ? 'English' : 'Ø§Ù„Ø¹Ø±Ø¨ÙŠØ©'}</span>
                    {locale === loc && <Check className="settings-check-icon h-5 w-5" />}
                  </div>
                </button>
              ))}
            </div>
          </div>

          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h2 className="text-2xl font-bold mb-4">{t('settings.theme')}</h2>
            <p className="text-gray-600 dark:text-gray-400 mb-4">{t('settings.selectTheme')}</p>

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              {THEME_OPTIONS.map(({ value, labelKey }) => (
                <button
                  key={value}
                  onClick={() => handleThemeChange(value as 'light' | 'dark' | 'system')}
                  className={`settings-selectable p-4 rounded-lg border-2 transition-all ${
                    settings.theme === value
                      ? 'is-selected bg-primary/10'
                      : 'border-gray-300 dark:border-gray-600 hover:border-primary'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className="font-semibold">{t(labelKey)}</span>
                    {settings.theme === value && <Check className="settings-check-icon h-5 w-5" />}
                  </div>
                </button>
              ))}
            </div>
          </div>

          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h2 className="text-2xl font-bold mb-2 flex items-center gap-2">
              <Palette className="h-6 w-6" />
              {t('settings.colors')}
            </h2>
            <p className="text-gray-600 dark:text-gray-400 mb-6">{t('settings.customizeColors')}</p>

            <div className="mb-8">
              <h3 className="text-lg font-semibold mb-4">{t('common.save')}</h3>
              <div className={`grid grid-cols-1 sm:grid-cols-2 gap-4 ${isRTL ? 'lg:grid-cols-2' : 'lg:grid-cols-4'}`}>
                {COLOR_PRESETS.map((preset, presetIndex) => (
                  <button
                    key={preset.name}
                    onClick={() => handleColorPreset(preset)}
                    className="p-4 rounded-lg border-2 border-gray-300 dark:border-gray-600 hover:border-primary transition-all"
                  >
                    <div className="flex items-center gap-3 mb-3">
                      <div className={`w-8 h-8 rounded-full settings-preset-primary-${presetIndex}`} />
                      <div className={`w-8 h-8 rounded-full settings-preset-secondary-${presetIndex}`} />
                      <div className={`w-8 h-8 rounded-full settings-preset-accent-${presetIndex}`} />
                    </div>
                    <p className="font-semibold">{preset.name}</p>
                  </button>
                ))}
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-4">{t('maintenance.customColors')}</h3>
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 mb-6">
                <div>
                  <label className="block text-sm font-medium mb-2">{t('settings.primaryColor')}</label>
                  <div className="flex items-center gap-3">
                    <input
                      type="color"
                      value={primaryColor}
                      onChange={(e) => setPrimaryColor(e.target.value)}
                      className="h-12 w-16 rounded cursor-pointer border border-gray-300 dark:border-gray-600"
                    />
                    <input
                      type="text"
                      value={primaryColor}
                      onChange={(e) => setPrimaryColor(e.target.value)}
                      className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white"
                      placeholder="#3b82f6"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium mb-2">{t('settings.secondaryColor')}</label>
                  <div className="flex items-center gap-3">
                    <input
                      type="color"
                      value={secondaryColor}
                      onChange={(e) => setSecondaryColor(e.target.value)}
                      className="h-12 w-16 rounded cursor-pointer border border-gray-300 dark:border-gray-600"
                    />
                    <input
                      type="text"
                      value={secondaryColor}
                      onChange={(e) => setSecondaryColor(e.target.value)}
                      className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white"
                      placeholder="#10b981"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium mb-2">{t('settings.accentColor')}</label>
                  <div className="flex items-center gap-3">
                    <input
                      type="color"
                      value={accentColor}
                      onChange={(e) => setAccentColor(e.target.value)}
                      className="h-12 w-16 rounded cursor-pointer border border-gray-300 dark:border-gray-600"
                    />
                    <input
                      type="text"
                      value={accentColor}
                      onChange={(e) => setAccentColor(e.target.value)}
                      className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white"
                      placeholder="#f59e0b"
                    />
                  </div>
                </div>
              </div>

              <Button onClick={handleCustomColors} className="settings-save-button w-full sm:w-auto">
                {t('settings.save')}
              </Button>
            </div>
          </div>

          <div className={`flex gap-4 ${isRTL ? 'flex-row-reverse' : ''}`}>
            <Link href={`/${locale}/menu`} className="flex-1 sm:flex-none">
              <Button variant="outline" className="w-full">
                {t('navigation.menu')}
              </Button>
            </Link>
            <Link href={`/${locale}/orders`} className="flex-1 sm:flex-none">
              <Button variant="outline" className="w-full">
                {t('navigation.orders')}
              </Button>
            </Link>
          </div>
        </div>
      </div>

      <style jsx>{`
        .settings-page-header {
          background: linear-gradient(135deg, ${themePrimaryColor} 0%, ${themeSecondaryColor} 100%);
        }

        .settings-selectable.is-selected {
          border-color: ${themePrimaryColor};
        }

        .settings-check-icon {
          color: ${themePrimaryColor};
        }

        .settings-save-button {
          background-color: ${themePrimaryColor};
        }

        ${presetStyles}
      `}</style>
    </>
  )
}
