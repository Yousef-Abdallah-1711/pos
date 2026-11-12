'use client'

import { useThemeStore } from '@/hooks/useThemeStore'
import { Smartphone, Monitor } from 'lucide-react'
import { useMemo, useState, type ComponentProps } from 'react'
import { ThemeProvider, SectionRenderer } from '@pos-saas/ui-themes'

/**
 * LivePreview: Real-time theme preview in the dashboard
 *
 * IMPORTANT: This component now uses the SAME components as the restaurant website,
 * ensuring the live preview accurately reflects the production rendering.
 *
 * Key improvement: Now renders actual section components instead of just titles
 * This fixes the issue where dashboard preview didn't match the real website
 */
export function LivePreview() {
  const theme = useThemeStore((state) => state.currentTheme)
  const [viewMode, setViewMode] = useState<'desktop' | 'mobile'>('desktop')

  if (!theme) {
    return (
      <div className="flex items-center justify-center h-full bg-gray-100 dark:bg-gray-900 text-gray-500">
        Select a theme to see preview
      </div>
    )
  }

  const paletteItems = [
    { name: 'Primary', color: theme.colors.primary },
    { name: 'Secondary', color: theme.colors.secondary },
    { name: 'Accent', color: theme.colors.accent },
    { name: 'Background', color: theme.colors.background },
    { name: 'Text', color: theme.colors.text },
    { name: 'Border', color: theme.colors.border },
  ]

  const paletteRules = useMemo(
    () =>
      paletteItems
        .map(
          (item, index) => `
            .palette-item-${index} .palette-swatch {
              background-color: ${item.color};
              border-color: ${theme.colors.border};
            }
          `
        )
        .join('\n'),
    [paletteItems, theme.colors.border]
  )

  return (
    <div className="live-preview-themed h-full flex flex-col bg-gray-100 dark:bg-gray-900">
      {/* Controls */}
      <div className="flex items-center justify-between p-4 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
        <h3 className="text-sm font-semibold text-gray-900 dark:text-white">Live Preview</h3>

        <div className="flex items-center gap-2 bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
          <button
            onClick={() => setViewMode('desktop')}
            className={`px-3 py-1 rounded-md text-sm transition flex items-center gap-1 ${
              viewMode === 'desktop'
                ? 'bg-white dark:bg-gray-600 text-gray-900 dark:text-white'
                : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
            }`}
          >
            <Monitor className="h-4 w-4" />
            Desktop
          </button>
          <button
            onClick={() => setViewMode('mobile')}
            className={`px-3 py-1 rounded-md text-sm transition flex items-center gap-1 ${
              viewMode === 'mobile'
                ? 'bg-white dark:bg-gray-600 text-gray-900 dark:text-white'
                : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
            }`}
          >
            <Smartphone className="h-4 w-4" />
            Mobile
          </button>
        </div>
      </div>

      {/* Preview Area */}
      <div className="flex-1 overflow-auto p-6 flex items-center justify-center">
        <div
          className={`bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden ${
            viewMode === 'mobile' ? 'w-96 max-w-full' : 'w-full'
          }`}
        >
          {/* Preview Content - Using shared ThemeProvider and SectionRenderer */}
          <ThemeProvider theme={theme}>
            <div className="bg-white">
              {/* Header */}
              <div className="preview-header text-white">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    {theme.identity.logoUrl && (
                      <img
                        src={theme.identity.logoUrl}
                        alt="Logo"
                        className="h-10 w-10 rounded"
                      />
                    )}
                    <h1 className="text-2xl font-bold">{theme.identity.siteTitle}</h1>
                  </div>
                  <nav className="flex gap-4 text-sm">
                    <a href="#" className="hover:opacity-80 transition">
                      Home
                    </a>
                    <a href="#" className="hover:opacity-80 transition">
                      Menu
                    </a>
                    <a href="#" className="hover:opacity-80 transition">
                      About
                    </a>
                    <a href="#" className="hover:opacity-80 transition">
                      Contact
                    </a>
                  </nav>
                </div>
              </div>

              {/* Components Preview - NOW RENDERS ACTUAL COMPONENTS */}
              <div>
                {theme.components && theme.components.length > 0 ? (
                  <>
                    {theme.components
                      .filter((c) => c.enabled)
                      .sort((a, b) => a.displayOrder - b.displayOrder)
                      .map((component) => (
                        <SectionRenderer
                          key={component.id}
                          // ComponentConfig (string id) vs ThemeComponent (numeric id): two parallel theme models, see PHASE_2B_HANDOFF.md
                          component={component as unknown as ComponentProps<typeof SectionRenderer>['component']}
                          isArabic={false}
                        />
                      ))}
                  </>
                ) : (
                  // Default preview if no components
                  <div className="py-16 px-8 text-center bg-gray-50">
                    <h2 className="preview-title text-3xl font-bold mb-4">
                      Welcome to {theme.identity.siteTitle}
                    </h2>
                    <p className="text-gray-600 mb-6 max-w-2xl mx-auto">
                      Start adding components to see your website come to life.
                      Create sections for Hero, Products, Contact, and more!
                    </p>
                    <button className="preview-cta-button px-8 py-3 rounded-lg font-semibold hover:opacity-90 transition">
                      Add Components
                    </button>
                  </div>
                )}
              </div>

              {/* Color Palette Display */}
              <div className="palette-section p-8 bg-gray-50 border-t">
                <h3 className="text-lg font-semibold mb-4">Theme Colors</h3>
                <div className="grid grid-cols-2 gap-4 max-w-sm">
                  {paletteItems.map((item, index) => (
                    <div key={item.name} className={`palette-item-${index} flex items-center gap-3`}>
                      <div className="palette-swatch w-10 h-10 rounded-lg border-2" />
                      <span className="text-sm font-medium">{item.name}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Footer */}
              <div className="preview-footer p-6 text-center text-sm">
                <p>© 2025 {theme.identity.siteTitle}. All rights reserved.</p>
              </div>
            </div>
          </ThemeProvider>
        </div>
      </div>

      <style jsx>{`
        .preview-header {
          background-color: ${theme.colors.primary};
          padding: 2rem;
        }

        .preview-title {
          color: ${theme.colors.primary};
        }

        .preview-cta-button {
          background-color: ${theme.colors.accent};
          color: #ffffff;
        }

        .palette-section {
          border-color: ${theme.colors.border};
        }

        .preview-footer {
          background-color: ${theme.colors.text};
          color: ${theme.colors.background};
        }

        ${paletteRules}
      `}</style>
    </div>
  )
}