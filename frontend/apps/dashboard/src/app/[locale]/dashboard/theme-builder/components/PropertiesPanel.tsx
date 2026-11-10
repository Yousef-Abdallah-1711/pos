'use client'

import { useState } from 'react'
import { Palette, Type, Globe } from 'lucide-react'
import { useThemeStore } from '@/hooks/useThemeStore'
import { ColorPicker } from './ColorPicker'

const FONT_FAMILIES = [
  'Inter',
  'Roboto',
  'Poppins',
  'Playfair Display',
  'Georgia',
  'Times New Roman',
  'Courier New',
  'Helvetica',
  'Arial',
]

export function PropertiesPanel() {
  const theme = useThemeStore((state) => state.currentTheme)
  const selectedComponentId = useThemeStore((state) => state.selectedComponentId)
  const updateColors = useThemeStore((state) => state.updateThemeColors)
  const updateTypography = useThemeStore((state) => state.updateTypography)
  const updateIdentity = useThemeStore((state) => state.updateIdentity)
  const updateComponent = useThemeStore((state) => state.updateComponent)

  const [activeTab, setActiveTab] = useState<'colors' | 'typography' | 'identity' | 'component'>(
    'colors'
  )

  if (!theme) {
    return (
      <div className="flex items-center justify-center h-full text-gray-500">
        Select a theme to edit properties
      </div>
    )
  }

  const selectedComponent = theme.components.find((c) => c.id === selectedComponentId)

  return (
    <div className="h-full flex flex-col bg-white dark:bg-gray-800 border-l border-gray-200 dark:border-gray-700">
      {/* Tabs */}
      <div className="flex border-b border-gray-200 dark:border-gray-700">
        <button
          onClick={() => setActiveTab('colors')}
          className={`flex-1 px-4 py-3 text-sm font-medium border-b-2 transition ${
            activeTab === 'colors'
              ? 'border-primary-600 text-primary-600'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-300'
          }`}
        >
          <Palette className="h-4 w-4 inline mr-2" />
          Colors
        </button>

        <button
          onClick={() => setActiveTab('typography')}
          className={`flex-1 px-4 py-3 text-sm font-medium border-b-2 transition ${
            activeTab === 'typography'
              ? 'border-primary-600 text-primary-600'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-300'
          }`}
        >
          <Type className="h-4 w-4 inline mr-2" />
          Type
        </button>

        <button
          onClick={() => setActiveTab('identity')}
          className={`flex-1 px-4 py-3 text-sm font-medium border-b-2 transition ${
            activeTab === 'identity'
              ? 'border-primary-600 text-primary-600'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-300'
          }`}
        >
          <Globe className="h-4 w-4 inline mr-2" />
          Identity
        </button>

        {selectedComponent && (
          <button
            onClick={() => setActiveTab('component')}
            className={`flex-1 px-4 py-3 text-sm font-medium border-b-2 transition ${
              activeTab === 'component'
                ? 'border-primary-600 text-primary-600'
                : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-300'
            }`}
          >
            Component
          </button>
        )}
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-6 space-y-6">
        {/* Colors Tab */}
        {activeTab === 'colors' && (
          <div className="space-y-6">
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-4">Global Colors</h3>
              <div className="space-y-4">
                <ColorPicker
                  label="Primary"
                  value={theme.colors.primary}
                  onChange={(color) => updateColors({ primary: color })}
                  description="Main brand color for buttons and CTAs"
                />
                <ColorPicker
                  label="Secondary"
                  value={theme.colors.secondary}
                  onChange={(color) => updateColors({ secondary: color })}
                  description="Complementary brand color"
                />
                <ColorPicker
                  label="Accent"
                  value={theme.colors.accent}
                  onChange={(color) => updateColors({ accent: color })}
                  description="Highlight color for special elements"
                />
                <ColorPicker
                  label="Background"
                  value={theme.colors.background}
                  onChange={(color) => updateColors({ background: color })}
                  description="Page background color"
                />
                <ColorPicker
                  label="Text"
                  value={theme.colors.text}
                  onChange={(color) => updateColors({ text: color })}
                  description="Default text color"
                />
                <ColorPicker
                  label="Border"
                  value={theme.colors.border}
                  onChange={(color) => updateColors({ border: color })}
                  description="Border and divider color"
                />
                <ColorPicker
                  label="Shadow"
                  value={theme.colors.shadow}
                  onChange={(color) => updateColors({ shadow: color })}
                  description="Shadow color"
                />
              </div>
            </div>
          </div>
        )}

        {/* Typography Tab */}
        {activeTab === 'typography' && (
          <div className="space-y-6">
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-4">Typography</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Font Family
                  </label>
                  <select
                    value={theme.typography.fontFamily}
                    onChange={(e) => updateTypography({ fontFamily: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
                  >
                    {FONT_FAMILIES.map((font) => (
                      <option key={font} value={font}>
                        {font}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Base Font Size: {theme.typography.baseFontSize}px
                  </label>
                  <input
                    type="range"
                    min="12"
                    max="20"
                    value={theme.typography.baseFontSize}
                    onChange={(e) =>
                      updateTypography({ baseFontSize: parseInt(e.target.value) })
                    }
                    className="w-full"
                  />
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">12px - 20px</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Border Radius: {theme.typography.borderRadius}px
                  </label>
                  <input
                    type="range"
                    min="0"
                    max="20"
                    value={theme.typography.borderRadius}
                    onChange={(e) =>
                      updateTypography({ borderRadius: parseInt(e.target.value) })
                    }
                    className="w-full"
                  />
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">0px - 20px</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Line Height: {theme.typography.lineHeight.toFixed(1)}
                  </label>
                  <input
                    type="range"
                    min="1.2"
                    max="2"
                    step="0.1"
                    value={theme.typography.lineHeight}
                    onChange={(e) =>
                      updateTypography({ lineHeight: parseFloat(e.target.value) })
                    }
                    className="w-full"
                  />
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">1.2 - 2.0</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Identity Tab */}
        {activeTab === 'identity' && (
          <div className="space-y-6">
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-4">Website Identity</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Site Title
                  </label>
                  <input
                    type="text"
                    value={theme.identity.siteTitle}
                    onChange={(e) => updateIdentity({ siteTitle: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Logo URL
                  </label>
                  <input
                    type="text"
                    value={theme.identity.logoUrl || ''}
                    onChange={(e) => updateIdentity({ logoUrl: e.target.value })}
                    placeholder="https://example.com/logo.png"
                    className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Favicon URL
                  </label>
                  <input
                    type="text"
                    value={theme.identity.faviconUrl || ''}
                    onChange={(e) => updateIdentity({ faviconUrl: e.target.value })}
                    placeholder="https://example.com/favicon.ico"
                    className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Domain
                  </label>
                  <input
                    type="text"
                    value={theme.identity.domain || ''}
                    disabled
                    className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-600 dark:text-gray-400 text-sm opacity-60"
                  />
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">Read-only field</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Component Tab */}
        {activeTab === 'component' && selectedComponent && (
          <div className="space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Component Settings</h3>
            <div className="p-4 rounded-lg bg-gray-50 dark:bg-gray-700/30 border border-gray-200 dark:border-gray-600">
              <p className="text-sm text-gray-700 dark:text-gray-300 mb-4">
                <strong>Type:</strong> {selectedComponent.type}
              </p>
              <p className="text-sm text-gray-700 dark:text-gray-300">
                <strong>Title:</strong> {selectedComponent.title}
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
