'use client'

import { useState } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { Upload, ArrowLeft, AlertCircle, CheckCircle } from 'lucide-react'
import { ThemeData } from '@/types/theme'
import * as themeApi from '@/lib/api/themeApi'
import { toast } from 'sonner'
import { getLocaleFromPath } from '@/lib/translations'

export default function ImportThemePage() {
  const router = useRouter()
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)

  const [jsonInput, setJsonInput] = useState('')
  const [parsedTheme, setParsedTheme] = useState<Partial<ThemeData> | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    const reader = new FileReader()
    reader.onload = (e) => {
      try {
        const content = e.target?.result as string
        setJsonInput(content)
        parseJSON(content)
      } catch (err) {
        setError('Failed to read file')
        toast.error('Failed to read file')
      }
    }
    reader.readAsText(file)
  }

  const handlePaste = async () => {
    try {
      const text = await navigator.clipboard.readText()
      setJsonInput(text)
      parseJSON(text)
    } catch (err) {
      setError('Failed to read clipboard')
      toast.error('Failed to read clipboard')
    }
  }

  const parseJSON = (json: string) => {
    try {
      const theme = JSON.parse(json)

      // Validate theme structure
      if (!theme.colors || !theme.typography) {
        setError('Invalid theme format. Missing required fields.')
        setParsedTheme(null)
        return
      }

      setParsedTheme(theme as Partial<ThemeData>)
      setError(null)
    } catch (err) {
      setError('Invalid JSON format')
      setParsedTheme(null)
    }
  }

  const handleImport = async () => {
    if (!parsedTheme) {
      setError('No valid theme to import')
      return
    }
    if (!parsedTheme.colors || !parsedTheme.typography) {
      setError('Theme file must include colors and typography')
      return
    }

    try {
      setIsLoading(true)
      await themeApi.createTheme({
        name: parsedTheme.name || 'Imported Theme',
        description: parsedTheme.description,
        colors: parsedTheme.colors,
        typography: parsedTheme.typography,
        identity: parsedTheme.identity || parsedTheme.websiteIdentity,
        header: parsedTheme.header,
        footer: parsedTheme.footer,
        components: parsedTheme.components,
        customCss: parsedTheme.customCss
      })
      toast.success(`Theme "${parsedTheme.name || 'Imported'}" imported successfully!`)
      router.push(`/${locale}/dashboard/theme-builder`)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to import theme'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setIsLoading(false)
    }
  }

  const resetForm = () => {
    setJsonInput('')
    setParsedTheme(null)
    setError(null)
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Header */}
      <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
        <div className="max-w-4xl mx-auto px-6 py-6">
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
                <Upload className="h-6 w-6 text-blue-600" />
                Import Theme
              </h1>
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                Import a theme from a previously exported JSON file
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="max-w-4xl mx-auto px-6 py-8">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-8 space-y-6">
          {/* Upload/Paste Options */}
          <div className="space-y-4">
            <label className="block text-sm font-semibold text-gray-900 dark:text-white">
              Theme JSON
            </label>

            {/* Upload Buttons */}
            <div className="flex gap-2 flex-col sm:flex-row">
              <label className="flex-1 px-4 py-3 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg cursor-pointer hover:border-blue-500 hover:bg-blue-50 dark:hover:bg-blue-900/10 transition flex items-center justify-center gap-2 text-gray-700 dark:text-gray-300 font-medium">
                <Upload className="h-5 w-5" />
                <span>Upload JSON File</span>
                <input
                  type="file"
                  accept=".json"
                  onChange={handleFileUpload}
                  className="hidden"
                />
              </label>
              <button
                onClick={handlePaste}
                className="px-6 py-3 border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700/50 text-gray-900 dark:text-white rounded-lg transition font-medium"
              >
                Paste from Clipboard
              </button>
            </div>

            {/* JSON Textarea */}
            <div className="space-y-2">
              <label className="block text-sm text-gray-700 dark:text-gray-300">
                Or paste JSON directly:
              </label>
              <textarea
                value={jsonInput}
                onChange={(e) => {
                  setJsonInput(e.target.value)
                  if (e.target.value) {
                    parseJSON(e.target.value)
                  }
                }}
                placeholder="Paste your theme JSON here..."
                className="w-full h-48 p-4 border border-gray-300 dark:border-gray-600 rounded-lg bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          {/* Error Display */}
          {error && (
            <div className="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg flex gap-3">
              <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
              <div>
                <p className="text-sm font-medium text-red-900 dark:text-red-400">{error}</p>
              </div>
            </div>
          )}

          {/* Preview */}
          {parsedTheme && !error && (
            <div className="p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-700 rounded-lg">
              <div className="flex items-start gap-2 mb-4">
                <CheckCircle className="h-5 w-5 text-green-600 dark:text-green-400 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-green-900 dark:text-green-300">
                    ✓ Valid theme ready to import
                  </p>
                </div>
              </div>

              {/* Theme Details */}
              <div className="space-y-2 text-sm text-gray-800 dark:text-gray-200 bg-white dark:bg-gray-800/50 p-4 rounded">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase">Name</p>
                    <p className="text-base font-medium">{parsedTheme.name || 'Unnamed'}</p>
                  </div>
                  {parsedTheme.description && (
                    <div>
                      <p className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase">Description</p>
                      <p className="text-base font-medium">{parsedTheme.description}</p>
                    </div>
                  )}
                  {parsedTheme.components && (
                    <div>
                      <p className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase">Components</p>
                      <p className="text-base font-medium">{parsedTheme.components.length}</p>
                    </div>
                  )}
                  <div>
                    <p className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase">Colors</p>
                    <p className="text-base font-medium">{Object.keys(parsedTheme.colors || {}).length}</p>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Info Box */}
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-700 rounded-lg">
            <p className="text-sm text-blue-700 dark:text-blue-300">
              💡 <strong>Tip:</strong> Export a theme first from the theme list to get a valid JSON file. Then import it here to create a copy with modifications.
            </p>
          </div>

          {/* Actions */}
          <div className="flex gap-3 pt-4">
            <button
              onClick={handleImport}
              disabled={!parsedTheme || !!error || isLoading}
              className="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed text-white rounded-lg font-semibold transition flex items-center justify-center gap-2"
            >
              {isLoading ? (
                <>
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  Importing...
                </>
              ) : (
                <>
                  <Upload className="h-5 w-5" />
                  Import Theme
                </>
              )}
            </button>
            <button
              onClick={() => {
                resetForm()
                router.push(`/${locale}/dashboard/theme-builder`)
              }}
              className="px-6 py-3 border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700/50 text-gray-900 dark:text-white rounded-lg font-semibold transition"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
