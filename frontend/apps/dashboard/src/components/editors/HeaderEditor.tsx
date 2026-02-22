'use client'

import React, { useState, useCallback, useMemo } from 'react'
import { Plus, Trash2, ChevronUp, ChevronDown, Eye, EyeOff } from 'lucide-react'
import type { ThemeComponent, HeaderConfig as ThemeHeaderConfig } from '@/types/theme'
import { useUpdateHeader, useThemeHeader } from '@/stores'

interface HeaderEditorProps {
  component?: ThemeComponent
  onChange?: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  className?: string
  useStore?: boolean
}

interface NavItem {
  id: string
  label: string
  href: string
  order: number
  isActive?: boolean
}

interface HeaderConfig extends ThemeHeaderConfig {
  transparentBackground?: boolean
  mobileMenuStyle?: 'hamburger' | 'drawer' | 'collapse'
}

/**
 * HeaderEditor Component
 * Configure website header with logo, navigation items, and styling
 * Supports both local state (component prop) and themeBuilderStore (useStore prop)
 */
export function HeaderEditor({
  component,
  onChange,
  onPreview,
  className = '',
  useStore = false,
}: HeaderEditorProps): JSX.Element {
  const storeHeader = useThemeHeader()
  const updateHeader = useUpdateHeader()

  const initialHeaderConfig: HeaderConfig = useMemo(() => {
    if (useStore && storeHeader) {
      return {
        ...storeHeader,
        transparentBackground: false,
        mobileMenuStyle: 'hamburger',
      }
    }
    return (component?.config as unknown as HeaderConfig) || {
      navigationItems: [],
      navPosition: 'right',
      navAlignment: 'horizontal',
      backgroundColor: '#ffffff',
      textColor: '#000000',
      height: 80,
      padding: 16,
      stickyHeader: false,
      transparentBackground: false,
      shadow: true,
      mobileMenuStyle: 'hamburger',
      hideNavOnMobile: false,
      showLogo: true,
      logoHeight: 40,
    }
  }, [useStore, storeHeader, component])

  const [headerConfig, setHeaderConfig] = useState<HeaderConfig>(initialHeaderConfig)

  const [previewMode, setPreviewMode] = useState(false)
  const [editingItemId, setEditingItemId] = useState<string | null>(null)
  const [editingLabel, setEditingLabel] = useState('')
  const [editingHref, setEditingHref] = useState('')
  const navJustify =
    headerConfig.navPosition === 'left'
      ? 'flex-start'
      : headerConfig.navPosition === 'center'
        ? 'center'
        : 'flex-end'

  const handleConfigChange = useCallback(
    (key: keyof HeaderConfig, value: any) => {
      const newConfig: HeaderConfig = {
        ...headerConfig,
        [key]: value,
      }

      setHeaderConfig(newConfig)

      if (useStore) {
        updateHeader({
          [key]: value,
        } as Partial<ThemeHeaderConfig>)
      } else if (component && onChange) {
        const updated: ThemeComponent = {
          ...component,
          config: newConfig,
        }
        onChange(updated)
        onPreview?.(updated)
      }
    },
    [headerConfig, component, onChange, onPreview, useStore, updateHeader]
  )

  const addNavItem = useCallback(() => {
    const newItem: NavItem = {
      id: `nav-${Date.now()}`,
      label: 'New Link',
      href: '#',
      order: (headerConfig.navigationItems?.length || 0) + 1,
    }

    const newConfig: HeaderConfig = {
      ...headerConfig,
      navigationItems: [...(headerConfig.navigationItems || []), newItem],
    }

    setHeaderConfig(newConfig)

    if (useStore) {
      updateHeader({ navigationItems: newConfig.navigationItems })
    } else if (component && onChange) {
      const updated: ThemeComponent = {
        ...component,
        config: newConfig,
      }
      onChange(updated)
    }
  }, [headerConfig, component, onChange, useStore, updateHeader])

  const deleteNavItem = useCallback(
    (id: string) => {
      const newConfig: HeaderConfig = {
        ...headerConfig,
        navigationItems: (headerConfig.navigationItems || []).filter((item) => item.id !== id),
      }

      setHeaderConfig(newConfig)

      if (useStore) {
        updateHeader({ navigationItems: newConfig.navigationItems })
      } else if (component && onChange) {
        const updated: ThemeComponent = {
          ...component,
          config: newConfig,
        }
        onChange(updated)
      }
    },
    [headerConfig, component, onChange, useStore, updateHeader]
  )

  const updateNavItem = useCallback(
    (id: string, label: string, href: string) => {
      const newConfig: HeaderConfig = {
        ...headerConfig,
        navigationItems: (headerConfig.navigationItems || []).map((item) =>
          item.id === id ? { ...item, label, href } : item
        ),
      }

      setHeaderConfig(newConfig)

      if (useStore) {
        updateHeader({ navigationItems: newConfig.navigationItems })
      } else if (component && onChange) {
        const updated: ThemeComponent = {
          ...component,
          config: newConfig,
        }
        onChange(updated)
      }
    },
    [headerConfig, component, onChange, useStore, updateHeader]
  )

  const moveNavItem = useCallback(
    (id: string, direction: 'up' | 'down') => {
      const items = [...(headerConfig.navigationItems || [])]
      const index = items.findIndex((item) => item.id === id)

      if (direction === 'up' && index > 0) {
        ;[items[index], items[index - 1]] = [items[index - 1], items[index]]
      } else if (direction === 'down' && index < items.length - 1) {
        ;[items[index], items[index + 1]] = [items[index + 1], items[index]]
      }

      const newConfig: HeaderConfig = {
        ...headerConfig,
        navigationItems: items.map((item, idx) => ({ ...item, order: idx + 1 })),
      }

      setHeaderConfig(newConfig)

      if (useStore) {
        updateHeader({ navigationItems: newConfig.navigationItems })
      } else if (component && onChange) {
        const updated: ThemeComponent = {
          ...component,
          config: newConfig,
        }
        onChange(updated)
      }
    },
    [headerConfig, component, onChange, useStore, updateHeader]
  )

  if (previewMode) {
    return (
      <div className={`w-full space-y-4 ${className}`}>
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Header Preview</h3>
          <button
            onClick={() => setPreviewMode(false)}
            className="px-3 py-2 text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition"
          >
            Back to Settings
          </button>
        </div>

        {/* Desktop Preview */}
        <div className="space-y-2">
          <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300">Desktop Preview</h4>
          <div className="header-preview-desktop border border-gray-200 dark:border-gray-700 rounded-lg flex items-center justify-between">
            {/* Logo */}
            {headerConfig.showLogo && (
              <div className="flex items-center gap-2">
                {headerConfig.logoUrl && (
                  <img
                    src={headerConfig.logoUrl}
                    alt="Logo"
                    className="header-logo-desktop object-contain"
                  />
                )}
                {headerConfig.logoText && <span className="font-bold">{headerConfig.logoText}</span>}
              </div>
            )}

            {/* Navigation */}
            <div className="header-preview-nav">
              {(headerConfig.navigationItems || []).map((item) => (
                <a key={item.id} href={item.href} className="text-sm hover:opacity-70">
                  {item.label}
                </a>
              ))}
            </div>
          </div>
        </div>

        {/* Mobile Preview */}
        <div className="space-y-2">
          <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300">Mobile Preview</h4>
          <div className="header-preview-mobile border border-gray-200 dark:border-gray-700 rounded-lg flex items-center justify-between w-full max-w-sm">
            {/* Logo Mobile */}
            {headerConfig.showLogo && (
              <div className="flex items-center gap-2">
                {headerConfig.logoUrl && (
                  <img
                    src={headerConfig.logoUrl}
                    alt="Logo"
                    className="header-logo-mobile object-contain"
                  />
                )}
              </div>
            )}

            {/* Mobile Menu Icon */}
            <div className="text-2xl">☰</div>
          </div>
        </div>
        <style jsx>{`
          .header-preview-desktop {
            background-color: ${headerConfig.backgroundColor || '#ffffff'};
            color: ${headerConfig.textColor || '#000000'};
            height: ${headerConfig.height || 80}px;
            padding: 0 ${headerConfig.padding || 16}px;
            box-shadow: ${headerConfig.shadow ? '0 1px 3px rgba(0,0,0,0.1)' : 'none'};
            position: ${headerConfig.stickyHeader ? 'sticky' : 'relative'};
          }

          .header-preview-mobile {
            background-color: ${headerConfig.backgroundColor || '#ffffff'};
            color: ${headerConfig.textColor || '#000000'};
            height: ${headerConfig.height || 80}px;
            padding: 0 ${headerConfig.padding || 16}px;
          }

          .header-logo-desktop,
          .header-logo-mobile {
            height: ${headerConfig.logoHeight || 40}px;
          }

          .header-preview-nav {
            display: flex;
            gap: 2rem;
            justify-content: ${navJustify};
          }
        `}</style>
      </div>
    )
  }

  return (
    <div className={`w-full space-y-6 ${className}`}>
      {/* Header */}
      <div>
        <h2 className="text-lg font-bold text-gray-900 dark:text-white mb-2">Header Configuration</h2>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          Configure your website header with logo, navigation items, and styling.
        </p>
      </div>

      {/* Logo Section */}
      <div className="space-y-4 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <h3 className="font-semibold text-gray-900 dark:text-white">Logo Settings</h3>

        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={headerConfig.showLogo || false}
            onChange={(e) => handleConfigChange('showLogo', e.target.checked)}
            className="rounded"
          />
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Show Logo</span>
        </label>

        {headerConfig.showLogo && (
          <>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Logo URL
              </label>
              <input
                type="url"
                value={headerConfig.logoUrl || ''}
                onChange={(e) => handleConfigChange('logoUrl', e.target.value)}
                placeholder="https://example.com/logo.png"
                className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Logo Height: {headerConfig.logoHeight || 40}px
              </label>
              <input
                type="range"
                min="20"
                max="100"
                value={headerConfig.logoHeight || 40}
                onChange={(e) => handleConfigChange('logoHeight', parseInt(e.target.value))}
                className="w-full"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Logo Text (Alternative)
              </label>
              <input
                type="text"
                value={headerConfig.logoText || ''}
                onChange={(e) => handleConfigChange('logoText', e.target.value)}
                placeholder="Your Brand Name"
                className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
              />
            </div>
          </>
        )}
      </div>

      {/* Navigation Section */}
      <div className="space-y-4 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold text-gray-900 dark:text-white">Navigation Items</h3>
          <button
            onClick={addNavItem}
            className="flex items-center gap-1.5 px-3 py-1.5 text-sm bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition"
          >
            <Plus className="h-4 w-4" />
            Add Item
          </button>
        </div>

        {/* Navigation Position */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Position
            </label>
            <select
              value={headerConfig.navPosition || 'right'}
              onChange={(e) => handleConfigChange('navPosition', e.target.value as any)}
              className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
            >
              <option value="left">Left</option>
              <option value="center">Center</option>
              <option value="right">Right</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Alignment
            </label>
            <select
              value={headerConfig.navAlignment || 'horizontal'}
              onChange={(e) => handleConfigChange('navAlignment', e.target.value as any)}
              className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
            >
              <option value="horizontal">Horizontal</option>
              <option value="vertical">Vertical</option>
            </select>
          </div>
        </div>

        {/* Navigation Items List */}
        <div className="space-y-2 max-h-64 overflow-y-auto">
          {(headerConfig.navigationItems || []).map((item, index) => (
            <div key={item.id} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3 space-y-2">
              {editingItemId === item.id ? (
                <>
                  <div className="grid grid-cols-2 gap-2">
                    <input
                      type="text"
                      value={editingLabel}
                      onChange={(e) => setEditingLabel(e.target.value)}
                      placeholder="Label"
                      className="px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white text-sm"
                    />
                    <input
                      type="text"
                      value={editingHref}
                      onChange={(e) => setEditingHref(e.target.value)}
                      placeholder="URL"
                      className="px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white text-sm"
                    />
                  </div>
                  <div className="flex gap-2 justify-end">
                    <button
                      onClick={() => {
                        updateNavItem(item.id, editingLabel, editingHref)
                        setEditingItemId(null)
                      }}
                      className="px-2 py-1 text-sm bg-green-600 hover:bg-green-700 text-white rounded transition"
                    >
                      Save
                    </button>
                    <button
                      onClick={() => setEditingItemId(null)}
                      className="px-2 py-1 text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white rounded transition"
                    >
                      Cancel
                    </button>
                  </div>
                </>
              ) : (
                <div className="flex items-center justify-between">
                  <div
                    className="flex-1 cursor-pointer"
                    onClick={() => {
                      setEditingItemId(item.id)
                      setEditingLabel(item.label)
                      setEditingHref(item.href)
                    }}
                  >
                    <div className="text-sm font-medium text-gray-900 dark:text-white">{item.label}</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">{item.href}</div>
                  </div>

                  <div className="flex items-center gap-1">
                    <button
                      onClick={() => moveNavItem(item.id, 'up')}
                      disabled={index === 0}
                      className="p-1 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded disabled:opacity-50 transition"
                    >
                      <ChevronUp className="h-4 w-4" />
                    </button>
                    <button
                      onClick={() => moveNavItem(item.id, 'down')}
                      disabled={index === (headerConfig.navigationItems?.length || 0) - 1}
                      className="p-1 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded disabled:opacity-50 transition"
                    >
                      <ChevronDown className="h-4 w-4" />
                    </button>
                    <button
                      onClick={() => deleteNavItem(item.id)}
                      className="p-1 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Styling Section */}
      <div className="space-y-4 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <h3 className="font-semibold text-gray-900 dark:text-white">Styling</h3>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Background Color
            </label>
            <input
              type="color"
              value={headerConfig.backgroundColor || '#ffffff'}
              onChange={(e) => handleConfigChange('backgroundColor', e.target.value)}
              className="w-full h-10 border border-gray-200 dark:border-gray-600 rounded-lg cursor-pointer"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Text Color
            </label>
            <input
              type="color"
              value={headerConfig.textColor || '#000000'}
              onChange={(e) => handleConfigChange('textColor', e.target.value)}
              className="w-full h-10 border border-gray-200 dark:border-gray-600 rounded-lg cursor-pointer"
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Height: {headerConfig.height || 80}px
            </label>
            <input
              type="range"
              min="40"
              max="200"
              value={headerConfig.height || 80}
              onChange={(e) => handleConfigChange('height', parseInt(e.target.value))}
              className="w-full"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Padding: {headerConfig.padding || 16}px
            </label>
            <input
              type="range"
              min="8"
              max="32"
              value={headerConfig.padding || 16}
              onChange={(e) => handleConfigChange('padding', parseInt(e.target.value))}
              className="w-full"
            />
          </div>
        </div>

        {/* Toggles */}
        <div className="space-y-2">
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={headerConfig.stickyHeader || false}
              onChange={(e) => handleConfigChange('stickyHeader', e.target.checked)}
              className="rounded"
            />
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Sticky Header</span>
          </label>

          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={headerConfig.shadow || false}
              onChange={(e) => handleConfigChange('shadow', e.target.checked)}
              className="rounded"
            />
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Show Shadow</span>
          </label>

          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={headerConfig.hideNavOnMobile || false}
              onChange={(e) => handleConfigChange('hideNavOnMobile', e.target.checked)}
              className="rounded"
            />
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Hide Nav on Mobile</span>
          </label>
        </div>
      </div>

      {/* Preview Button */}
      <div className="flex gap-3">
        <button
          onClick={() => setPreviewMode(true)}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition"
        >
          Preview Header
        </button>
      </div>
    </div>
  )
}
