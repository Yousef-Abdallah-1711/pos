'use client'

import { useState, useEffect } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { Palette, ArrowLeft, Loader2 } from 'lucide-react'
import { ThemeData } from '@/types/theme'
import * as themeApi from '@/lib/api/themeApi'
import { toast } from 'sonner'
import { getLocaleFromPath } from '@/lib/translations'
import { debugLog, debugError } from '@/lib/debug'
import { PresetsGallery } from '@/components/theme/PresetsGallery'

export default function PresetsPage() {
  const router = useRouter()
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)

  const [presets, setPresets] = useState<ThemeData[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isSelecting, setIsSelecting] = useState(false)

  useEffect(() => {
    const loadPresets = async () => {
      try {
        setIsLoading(true)
        debugLog('📡 Loading theme presets...')
        const presetsData = await themeApi.getThemePresets(100)
        debugLog('✅ Got presets:', presetsData)
        setPresets(presetsData)
      } catch (err) {
        debugError('Failed to load presets:', err)
        toast.error('Failed to load theme presets')
        setPresets([])
      } finally {
        setIsLoading(false)
      }
    }

    loadPresets()
  }, [])

  const handleSelectPreset = async (preset: ThemeData) => {
    try {
      setIsSelecting(true)
      debugLog('Creating theme from preset:', preset.name)
      const newTheme = await themeApi.duplicateTheme(preset.id, `${preset.name} (Custom)`)
      toast.success(`Theme created from "${preset.name}"`)
      router.push(`/${locale}/dashboard/theme-builder/editor/${newTheme.id}`)
    } catch (err) {
      debugError('Failed to create theme from preset:', err)
      const errorMessage = err instanceof Error ? err.message : 'Failed to create theme'
      toast.error(errorMessage)
      setIsSelecting(false)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <Loader2 className="h-12 w-12 animate-spin text-blue-600 mx-auto mb-4" />
          <p className="text-gray-600">Loading theme presets...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Header */}
      <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
        <div className="max-w-6xl mx-auto px-6 py-6">
          <div className="flex items-center gap-4 mb-4">
            <button
              onClick={() => router.push(`/${locale}/dashboard/theme-builder`)}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition"
              title="Go back"
            >
              <ArrowLeft className="h-5 w-5 text-gray-600 dark:text-gray-400" />
            </button>
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
                <Palette className="h-6 w-6 text-blue-600" />
                Theme Templates
              </h1>
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                Choose a pre-designed theme and customize it to match your brand
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="max-w-6xl mx-auto px-6 py-8">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-8">
          {isSelecting && (
            <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
              <div className="bg-white dark:bg-gray-800 rounded-lg p-6 flex gap-4 items-center">
                <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
                <p className="text-gray-900 dark:text-white">Creating theme from preset...</p>
              </div>
            </div>
          )}

          <PresetsGallery
            presets={presets}
            onSelectPreset={handleSelectPreset}
            isLoading={false}
          />
        </div>
      </div>
    </div>
  )
}
