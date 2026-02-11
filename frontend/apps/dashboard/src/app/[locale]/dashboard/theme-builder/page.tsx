'use client'

import { useState, useEffect, useMemo } from 'react'
import { usePathname, useRouter } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { Palette, Check, Loader2, AlertCircle, Search, Copy, Trash2, Eye, Download, Upload } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { ThemeData } from '@/types/theme'
import * as themeApi from '@/lib/api/themeApi'
import { toast } from 'sonner'
import { debugLog, debugWarn, debugError } from '@/lib/debug'
import { ExportDialog } from '@/components/theme/ExportDialog'


export default function ThemeBuilderPage() {
  const pathname = usePathname()
  const router = useRouter()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)

  const [themes, setThemes] = useState<ThemeData[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [debugInfo, setDebugInfo] = useState<string>('')
  const [searchQuery, setSearchQuery] = useState('')
  const [filterStatus, setFilterStatus] = useState<'all' | 'active' | 'inactive'>('all')
  const [sortBy, setSortBy] = useState<'name' | 'date' | 'status'>('date')

  // Dialog states
  const [exportThemeId, setExportThemeId] = useState<string | null>(null)

  // ===========================================
  // 🔍 DEBUG: Log component mount and auth status
  // ===========================================
  useEffect(() => {
    const debugMsg = `
    ===== THEME BUILDER DEBUG =====
    🕐 Mounted at: ${new Date().toLocaleTimeString()}
    👤 User Email: ${user?.email || 'NO USER'}
    👤 User Name: ${user?.name || 'NONE'}
    🆔 User ID: ${user?.id || 'NONE'}
    🏪 Restaurant ID: ${user?.restaurant_id || 'NONE'}
    🔑 Token exists: ${!!token}
    🔑 Token length: ${token?.length || 0}
    🌐 Locale: ${locale}
    📍 Pathname: ${pathname}
    ===============================
    `
    debugLog(debugMsg)
    setDebugInfo(debugMsg)
  }, [user?.email, user?.id, user?.restaurant_id, user?.name, token, locale, pathname])

  // ===========================================
  // 🔄 MAIN: Fetch themes from API
  // Purpose: Load all themes for the restaurant
  // Requires: user.restaurant_id, JWT token
  // ===========================================
  useEffect(() => {
    const fetchThemes = async () => {
      debugLog('📡 [fetchThemes] Starting fetch...')
      debugLog('📡 [fetchThemes] User restaurant_id:', user?.restaurant_id)
      debugLog('📡 [fetchThemes] Token:', token ? 'EXISTS' : 'MISSING')

      // Check if we have the required data
      if (!user?.restaurant_id) {
        debugWarn('⚠️  [fetchThemes] Missing restaurant_id! User object:', user)
        setError('No restaurant found. Please log in again.')
        setIsLoading(false)
        return
      }

      if (!token) {
        debugWarn('⚠️  [fetchThemes] Missing JWT token!')
        setError('No authentication token. Please log in again.')
        setIsLoading(false)
        return
      }

      try {
        setIsLoading(true)
        setError(null)
        debugLog('📡 [fetchThemes] Calling themeApi.getThemes()...')
        const fetchedThemes = await themeApi.getThemes()
        debugLog('✅ [fetchThemes] Success! Got themes:', fetchedThemes)
        setThemes(fetchedThemes)
      } catch (err) {
        debugError('❌ [fetchThemes] Error:', err)
        debugError('❌ [fetchThemes] Error type:', err instanceof Error ? 'Error object' : typeof err)
        debugError('❌ [fetchThemes] Full error:', JSON.stringify(err, null, 2))

        // More detailed error message
        let errorMessage = 'Failed to load themes'
        if (err instanceof Error) {
          errorMessage = err.message
          // Check if it's an auth error
          if (errorMessage.includes('401') || errorMessage.includes('Unauthorized')) {
            errorMessage = 'Authentication failed. Your session may have expired. Please log in again.'
          }
        }

        setError(errorMessage)
        setThemes([])
      } finally {
        setIsLoading(false)
      }
    }

    // Only fetch if we're ready
    if (user?.restaurant_id && token) {
      fetchThemes()
    } else {
      debugWarn('⚠️  [useEffect] Skipping fetch - waiting for auth. User:', !!user, 'Token:', !!token)
    }
  }, [user?.restaurant_id, token])


  const filteredAndSortedThemes = useMemo(() => {
    let filtered = themes

    if (filterStatus !== 'all') {
      filtered = filtered.filter((theme) =>
        filterStatus === 'active' ? theme.isActive : !theme.isActive
      )
    }

    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase()
      filtered = filtered.filter((theme) =>
        theme.name.toLowerCase().includes(query) ||
        theme.description?.toLowerCase().includes(query)
      )
    }

    const sorted = [...filtered]
    if (sortBy === 'name') {
      sorted.sort((a, b) => a.name.localeCompare(b.name))
    } else if (sortBy === 'date') {
      sorted.sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
    } else if (sortBy === 'status') {
      sorted.sort((a, b) => (b.isActive ? 1 : 0) - (a.isActive ? 1 : 0))
    }

    return sorted
  }, [themes, searchQuery, filterStatus, sortBy])

  const handleActivateTheme = async (themeId: string) => {
    try {
      setError(null)
      await themeApi.activateTheme(themeId)

      // Update local state
      setThemes(themes.map(t => ({
        ...t,
        isActive: t.id === themeId
      })))

      const themeName = themes.find(t => t.id === themeId)?.name || 'Theme'
      toast.success(`Theme "${themeName}" activated successfully`)
    } catch (err) {
      debugError('Failed to activate theme:', err)
      const errorMessage = err instanceof Error ? err.message : 'Failed to activate theme'
      setError(errorMessage)
      toast.error(errorMessage)
    }
  }

  const handleDeleteTheme = async (themeId: string) => {
    if (!window.confirm('Are you sure you want to delete this theme?')) {
      return
    }

    try {
      setError(null)
      await themeApi.deleteTheme(themeId)
      setThemes(themes.filter(t => t.id !== themeId))
      toast.success('Theme deleted successfully')
    } catch (err) {
      debugError('Failed to delete theme:', err)
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete theme'
      setError(errorMessage)
      toast.error(errorMessage)
    }
  }

  const handleDuplicateTheme = async (themeId: string, themeName: string) => {
    try {
      setError(null)
      const duplicatedTheme = await themeApi.duplicateTheme(themeId, `${themeName} (Copy)`)
      setThemes([...themes, duplicatedTheme])
      toast.success(`Theme duplicated as "${duplicatedTheme.name}"`)
    } catch (err) {
      debugError('Failed to duplicate theme:', err)
      const errorMessage = err instanceof Error ? err.message : 'Failed to duplicate theme'
      setError(errorMessage)
      toast.error(errorMessage)
    }
  }



  // Show loading state
  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <Loader2 className="h-12 w-12 animate-spin text-blue-600 mx-auto mb-4" />
          <p className="text-gray-600">Loading themes...</p>
        </div>
      </div>
    )
  }

  // Show error state
  if (error) {
    return (
      <div className="p-6">
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 flex gap-3">
          <AlertCircle className="h-5 w-5 text-red-600 flex-shrink-0 mt-0.5" />
          <div>
            <h3 className="font-semibold text-red-900">Error</h3>
            <p className="text-red-700 text-sm">{error}</p>
            <button
              onClick={() => window.location.reload()}
              className="mt-2 text-sm text-red-600 hover:text-red-700 font-medium"
            >
              Retry
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      {/* DEBUG INFO - Always visible for troubleshooting */}
      <div className="mb-6 p-4 bg-gray-100 border border-gray-300 rounded-lg font-mono text-xs text-gray-700 whitespace-pre-wrap">
        <div className="font-bold mb-2">🔍 DEBUG INFO (Copy for support):</div>
        {debugInfo}
      </div>

      {/* Header with Website Preview */}
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-2">
          <Palette className="h-8 w-8 text-blue-600" />
          <h1 className="text-3xl font-bold text-gray-900">Theme Builder</h1>
        </div>
        <p className="text-gray-600 mb-4">Create and manage website themes</p>

        {/* Website Subdomain Card */}
        {user?.restaurant_id && (
          <div className="bg-gradient-to-r from-blue-50 to-blue-100 border border-blue-200 rounded-lg p-4 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="h-10 w-10 bg-blue-600 rounded-lg flex items-center justify-center">
                <Eye className="h-5 w-5 text-white" />
              </div>
              <div>
                <p className="text-sm font-medium text-blue-900">Your Website</p>
                <p className="text-sm text-blue-700 font-mono">
                  {/*
                    Subdomain Format Options:
                    - Development: restaurant-{id}.localhost:3000
                    - Production: {slug}.yourdomain.com
                    - Multi-tenant: tenant-{id}.yourdomain.com
                    Customize this in environment variables or fetch restaurant slug from API
                  */}
                  restaurant-{user.restaurant_id}.localhost:3000
                </p>
              </div>
            </div>
            <a
              href={`http://restaurant-${user.restaurant_id}.localhost:3000`}
              target="_blank"
              rel="noopener noreferrer"
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition flex items-center gap-2 whitespace-nowrap"
            >
              <Eye className="h-4 w-4" />
              Preview Website
            </a>
          </div>
        )}
      </div>

      {/* Existing Themes */}
      {themes.length > 0 && (
        <div className="mb-12">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-bold text-gray-900">Your Themes ({filteredAndSortedThemes.length})</h2>
            <div className="flex gap-2">
              <button
                onClick={() => router.push(`/${locale}/dashboard/theme-builder/presets`)}
                className="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg font-medium transition flex items-center gap-2"
              >
                <Palette className="h-4 w-4" />
                From Template
              </button>
              <button
                onClick={() => router.push(`/${locale}/dashboard/theme-builder/import`)}
                className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg font-medium transition flex items-center gap-2"
              >
                <Upload className="h-4 w-4" />
                Import
              </button>
              <button
                onClick={() => router.push(`/${locale}/dashboard/theme-builder/editor/new`)}
                className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition"
              >
                + Create Theme
              </button>
            </div>
          </div>

          {/* Search and Filters */}
          <div className="bg-white border border-gray-200 rounded-lg p-4 mb-6 space-y-4">
            <div className="flex gap-4 flex-wrap">
              <div className="flex-1 min-w-xs">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                  <input
                    type="text"
                    placeholder="Search themes..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>

              <div className="flex gap-2">
                <select
                  value={filterStatus}
                  onChange={(e) => setFilterStatus(e.target.value as any)}
                  className="px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="all">All Themes</option>
                  <option value="active">Active Only</option>
                  <option value="inactive">Inactive</option>
                </select>

                <select
                  value={sortBy}
                  onChange={(e) => setSortBy(e.target.value as any)}
                  className="px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="date">Recent First</option>
                  <option value="name">Name (A-Z)</option>
                  <option value="status">Active First</option>
                </select>
              </div>
            </div>
          </div>

          {/* Themes Grid */}
          {filteredAndSortedThemes.length === 0 ? (
            <div className="text-center py-12 bg-gray-50 rounded-lg">
              <Palette className="h-12 w-12 text-gray-300 mx-auto mb-3" />
              <p className="text-gray-500 font-medium">No themes found</p>
              <p className="text-gray-400 text-sm mt-1">Try adjusting your search or filters</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {filteredAndSortedThemes.map((theme) => (
              <div
                key={theme.id}
                className="bg-white border border-gray-200 rounded-lg overflow-hidden hover:shadow-lg transition flex flex-col"
              >
                {/* Theme Preview */}
                <div className="h-28 flex gap-2 p-3"
                  style={{ backgroundColor: theme.colors?.background }}>
                  {['primary', 'secondary', 'accent'].map((color) => (
                    <div
                      key={color}
                      className="flex-1 rounded-md shadow-sm"
                      style={{
                        backgroundColor: theme.colors?.[color as keyof typeof theme.colors],
                      }}
                    />
                  ))}
                </div>

                {/* Theme Info */}
                <div className="p-4 flex-1 flex flex-col">
                  <div className="flex items-start justify-between mb-2">
                    <div>
                      <h3 className="font-bold text-gray-900">{theme.name}</h3>
                      <p className="text-xs text-gray-500 mt-0.5">
                        Updated {new Date(theme.updatedAt).toLocaleDateString()}
                      </p>
                    </div>
                    {theme.isActive && (
                      <div className="flex items-center gap-1 bg-green-100 text-green-700 px-2 py-1 rounded-full text-xs font-medium whitespace-nowrap">
                        <Check className="h-3 w-3" />
                        Active
                      </div>
                    )}
                  </div>
                  {theme.description && (
                    <p className="text-sm text-gray-600 mb-4 flex-1 line-clamp-2">{theme.description}</p>
                  )}

                  {/* Actions */}
                  <div className="grid grid-cols-3 gap-2">
                    <button
                      onClick={() => router.push(`/${locale}/dashboard/theme-builder/editor/${theme.id}`)}
                      className="px-3 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 text-sm font-medium transition flex items-center justify-center gap-1"
                      title="Edit theme"
                    >
                      <Eye className="h-3 w-3" />
                      Edit
                    </button>
                    <button
                      onClick={() => handleDuplicateTheme(theme.id, theme.name)}
                      className="px-3 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 text-sm font-medium transition flex items-center justify-center gap-1"
                      title="Duplicate theme"
                    >
                      <Copy className="h-3 w-3" />
                      Copy
                    </button>
                    <button
                      onClick={() => setExportThemeId(theme.id)}
                      className="px-3 py-2 bg-orange-100 text-orange-700 rounded hover:bg-orange-200 text-sm font-medium transition flex items-center justify-center gap-1"
                      title="Export theme as JSON"
                    >
                      <Download className="h-3 w-3" />
                      Export
                    </button>
                  </div>

                  {!theme.isActive && (
                    <button
                      onClick={() => handleActivateTheme(theme.id)}
                      className="mt-2 w-full px-3 py-2 bg-yellow-50 text-yellow-700 rounded hover:bg-yellow-100 text-sm font-medium transition"
                    >
                      Activate Theme
                    </button>
                  )}

                  {theme.isActive && (
                    <button
                      onClick={() => handleDeleteTheme(theme.id)}
                      className="mt-2 w-full px-3 py-2 bg-red-50 text-red-600 rounded hover:bg-red-100 text-sm font-medium transition flex items-center justify-center gap-1"
                    >
                      <Trash2 className="h-3 w-3" />
                      Delete
                    </button>
                  )}
                </div>
              </div>
            ))}
            </div>
          )}
        </div>
      )}

      {/* Export Dialog */}
      {exportThemeId && (
        <ExportDialog
          theme={themes.find((t) => t.id === exportThemeId)!}
          isOpen={!!exportThemeId}
          onClose={() => setExportThemeId(null)}
        />
      )}

    </div>
  )
}
