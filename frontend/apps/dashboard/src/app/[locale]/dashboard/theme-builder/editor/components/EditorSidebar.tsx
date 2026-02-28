'use client'

import { useState, useEffect } from 'react'
import { usePathname } from 'next/navigation'
import { Palette, Home, Settings, Layout, Copyright, Code, Edit3, ChevronLeft, Clock } from 'lucide-react'
import { useThemeStore } from '@/hooks/useThemeStore'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { ColorPicker } from '../../components/ColorPicker'
import { ComponentLibrary } from './ComponentLibrary'
import { ComponentEditor } from './ComponentEditor'
import { ThemeHistory } from './ThemeHistory'
import { HeaderEditor } from '@/components/editors/HeaderEditor'
import { FooterEditor } from '@/components/editors/FooterEditor'
import { CustomCSSEditor } from '@/components/editors/CustomCSSEditor'
import { FontPicker } from '@/components/editors/FontPicker'
import { ContrastChecker } from '@/components/editors/ContrastChecker'
import { ColorPalettes } from '@/components/editors/ColorPalettes'
import { TypographyPreview } from '@/components/editors/TypographyPreview'
import { EditorErrorBoundary } from '@/components/editors/EditorErrorBoundary'
import { HeaderConfig, FooterConfig } from '@/types/theme'

type TabType = 'colors' | 'homepage' | 'header' | 'footer' | 'settings' | 'history' | 'component-edit'

// Helper function to extract label from translation objects
function extractLabel(label: any): string {
  if (typeof label === 'string') return label
  if (label && typeof label === 'object') {
    // Handle translation objects like {en: "...", ar: "..."}
    return label.en || label.ar || Object.values(label)[0] || ''
  }
  return ''
}

/**
 * Click-based sidebar with tab navigation
 * - Colors: Color picker for all theme colors
 * - Homepage: Components and layout
 * - Settings: Typography and other settings
 *
 * When a tab is clicked, its content displays in the sidebar (not collapsible)
 */
export function EditorSidebar() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'
  const theme = useThemeStore((state) => state.currentTheme)
  const selectedComponentId = useThemeStore((state) => state.selectedComponentId)
  const updateColors = useThemeStore((state) => state.updateThemeColors)
  const updateHeader = useThemeStore((state) => state.updateHeader)
  const updateFooter = useThemeStore((state) => state.updateFooter)
  const updateTypography = useThemeStore((state) => state.updateTypography)
  const updateCustomCss = useThemeStore((state) => state.updateCustomCss)
  const [activeTab, setActiveTab] = useState<TabType>('colors')
  const [previousTab, setPreviousTab] = useState<TabType>('colors')

  // Watch selectedComponentId and switch to component-edit tab when component is selected
  useEffect(() => {
    if (selectedComponentId) {
      setPreviousTab(activeTab)
      setActiveTab('component-edit')
    }
  }, [selectedComponentId, activeTab])

  if (!theme) {
    return (
        <div className="flex items-center justify-center h-full text-gray-500">
        {t('themeBuilderEditor.noThemeSelected', 'No theme selected')}
      </div>
    )
  }

  // Default values - use these if theme properties are missing or invalid
  const defaultColors = {
    primary: '#3b82f6',
    secondary: '#1e40af',
    accent: '#0ea5e9',
    background: '#ffffff',
    text: '#1f2937',
    border: '#e5e7eb',
    shadow: '#1f2937',
  }

  const defaultTypography = {
    fontFamily: 'Inter',
    baseFontSize: 14,
    borderRadius: 8,
    lineHeight: 1.5,
    fontWeights: [400, 700] as number[],
    fontSize: 14,
    fontStyle: 'normal',
  }

  const defaultHeader = {
    logoUrl: '',
    logoText: 'Restaurant',
    logoHeight: 40,
    showLogo: true,
    navigationItems: [],
    navPosition: 'right' as const,
    navAlignment: 'horizontal' as const,
    backgroundColor: '#3b82f6',
    textColor: '#ffffff',
    height: 64,
    padding: 16,
    stickyHeader: false,
    showShadow: true,
    hideNavOnMobile: false,
  }

  const defaultFooter = {
    companyName: 'Restaurant',
    companyDescription: '',
    address: '',
    phone: '',
    email: '',
    copyrightText: '',
    socialLinks: [],
    footerSections: [],
    legalLinks: [],
    backgroundColor: '#1f2937',
    textColor: '#ffffff',
    linkColor: '#3b82f6',
    padding: 48,
    showLinks: true,
    showLegal: true,
    showBackToTop: true,
    columns: 3,
    layout: 'expanded' as const,
  }

  // Ensure colors object exists with defaults - merge theme colors with defaults
  const colors = {
    ...defaultColors,
    ...(theme.colors && typeof theme.colors === 'object' ? theme.colors : {}),
  }

  // Ensure typography object exists with defaults
  const typography = {
    ...defaultTypography,
    ...(theme.typography && typeof theme.typography === 'object' ? theme.typography : {}),
  }

  // Ensure header object exists with defaults
  const header = {
    ...defaultHeader,
    ...(theme.header && typeof theme.header === 'object' ? theme.header : {}),
  }

  // Ensure footer object exists with defaults
  const footer = {
    ...defaultFooter,
    ...(theme.footer && typeof theme.footer === 'object' ? theme.footer : {}),
  }

  // Find the selected component in theme.components
  const selectedComponent = selectedComponentId && Array.isArray(theme.components)
    ? theme.components.find((comp) => comp.id === selectedComponentId)
    : null

  return (
    <div className="h-full flex flex-col bg-white dark:bg-gray-800">
      {/* Sidebar Navigation - Tab Buttons */}
      <div className="border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-4 py-3 space-y-2">
        {activeTab === 'component-edit' && selectedComponent && (
          <button
            onClick={() => setActiveTab(previousTab)}
            className="w-full flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm font-medium transition-colors bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 border border-gray-200 dark:border-gray-600"
            title={t('themeBuilderEditor.backToPreviousTab', 'Go back to previous tab')}
          >
            <ChevronLeft className={`h-4 w-4 flex-shrink-0 ${isRTL ? 'rotate-180' : ''}`} />
            <span>{t('common.back', 'Back')}</span>
          </button>
        )}
        <div className="grid grid-cols-2 gap-2 pt-2">
          {[
            { id: 'colors' as TabType, label: t('themeBuilderEditor.tabs.colors', 'Colors'), icon: <Palette className="h-4 w-4" /> },
            { id: 'homepage' as TabType, label: t('themeBuilderEditor.tabs.homepage', 'Homepage'), icon: <Home className="h-4 w-4" /> },
            { id: 'header' as TabType, label: t('themeBuilderEditor.tabs.header', 'Header'), icon: <Layout className="h-4 w-4" /> },
            { id: 'footer' as TabType, label: t('themeBuilderEditor.tabs.footer', 'Footer'), icon: <Copyright className="h-4 w-4" /> },
            { id: 'settings' as TabType, label: t('themeBuilderEditor.tabs.settings', 'Settings'), icon: <Settings className="h-4 w-4" /> },
            { id: 'history' as TabType, label: t('themeBuilderEditor.tabs.history', 'History'), icon: <Clock className="h-4 w-4" /> },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`flex flex-col items-center justify-center gap-1.5 px-3 py-3 rounded-lg text-xs font-medium transition-all ${
                activeTab === tab.id
                  ? 'bg-primary-600 text-white shadow-md ring-2 ring-primary-400 ring-opacity-50'
                  : 'bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 border border-gray-200 dark:border-gray-600'
              }`}
              title={t('themeBuilderEditor.openTabSettings', 'Open settings')}
            >
              {tab.icon}
              <span className="leading-tight">{tab.label}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Content Area - Displays Selected Tab Content */}
      <div className="flex-1 overflow-y-auto p-6 space-y-6">
        {/* Colors Tab */}
        {activeTab === 'colors' && (
          <>
            {/* Color Palettes */}
            <EditorErrorBoundary componentType="color-palettes">
              <ColorPalettes onApply={updateColors} />
            </EditorErrorBoundary>

            {/* Individual Color Pickers */}
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-4 text-sm flex items-center gap-2">
                <Palette className="h-4 w-4 text-primary-600" />
                {t('themeBuilderEditor.globalColors', 'Global Colors')}
              </h3>
              <div className="space-y-4">
                <ColorPicker
                  label={t('themeBuilderEditor.colors.primary', 'Primary')}
                  value={colors.primary}
                  onChange={(color) => updateColors({ primary: color })}
                  description={t('themeBuilderEditor.colors.primaryDesc', 'Main brand color')}
                />
                <ColorPicker
                  label={t('themeBuilderEditor.colors.secondary', 'Secondary')}
                  value={colors.secondary}
                  onChange={(color) => updateColors({ secondary: color })}
                  description={t('themeBuilderEditor.colors.secondaryDesc', 'Complementary color')}
                />
                <ColorPicker
                  label={t('themeBuilderEditor.colors.accent', 'Accent')}
                  value={colors.accent}
                  onChange={(color) => updateColors({ accent: color })}
                  description={t('themeBuilderEditor.colors.accentDesc', 'Highlight color')}
                />
                <ColorPicker
                  label={t('themeBuilderEditor.colors.background', 'Background')}
                  value={colors.background}
                  onChange={(color) => updateColors({ background: color })}
                  description={t('themeBuilderEditor.colors.backgroundDesc', 'Page background')}
                />
                <ColorPicker
                  label={t('themeBuilderEditor.colors.text', 'Text')}
                  value={colors.text}
                  onChange={(color) => updateColors({ text: color })}
                  description={t('themeBuilderEditor.colors.textDesc', 'Text color')}
                />
                <ColorPicker
                  label={t('themeBuilderEditor.colors.border', 'Border')}
                  value={colors.border}
                  onChange={(color) => updateColors({ border: color })}
                  description={t('themeBuilderEditor.colors.borderDesc', 'Border color')}
                />
                <ColorPicker
                  label={t('themeBuilderEditor.colors.shadow', 'Shadow')}
                  value={colors.shadow}
                  onChange={(color) => updateColors({ shadow: color })}
                  description={t('themeBuilderEditor.colors.shadowDesc', 'Shadow color')}
                />
              </div>
            </div>

            {/* Contrast Checker */}
            <EditorErrorBoundary componentType="contrast-checker">
              <ContrastChecker colors={colors} />
            </EditorErrorBoundary>
          </>
        )}

        {/* Homepage Tab */}
        {activeTab === 'homepage' && (
          <>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-4 text-sm flex items-center gap-2">
                <Home className="h-4 w-4 text-primary-600" />
                {t('themeBuilderEditor.pageComponents', 'Page Components')}
              </h3>
              <ComponentLibrary />
            </div>
          </>
        )}

        {/* Header Tab */}
        {activeTab === 'header' && (
          <>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-4 text-sm flex items-center gap-2">
                <Layout className="h-4 w-4 text-primary-600" />
                {t('themeBuilderEditor.headerConfiguration', 'Header Configuration')}
              </h3>
              <EditorErrorBoundary componentType="header">
                <HeaderEditor
                  component={{
                    id: 'header',
                    type: 'custom',
                    title: t('themeBuilderEditor.tabs.header', 'Header'),
                    displayOrder: 0,
                    enabled: true,
                    config: header,
                  }}
                  onChange={(component) => {
                    const config = component.config as Partial<HeaderConfig>
                    updateHeader(config)
                  }}
                  className="w-full"
                />
              </EditorErrorBoundary>
            </div>
          </>
        )}

        {/* Footer Tab */}
        {activeTab === 'footer' && (
          <>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-4 text-sm flex items-center gap-2">
                <Copyright className="h-4 w-4 text-primary-600" />
                {t('themeBuilderEditor.footerConfiguration', 'Footer Configuration')}
              </h3>
              <EditorErrorBoundary componentType="footer">
                <FooterEditor
                  component={{
                    id: 'footer',
                    type: 'custom',
                    title: t('themeBuilderEditor.tabs.footer', 'Footer'),
                    displayOrder: 0,
                    enabled: true,
                    config: footer,
                  }}
                  onChange={(component) => {
                    const config = component.config as Partial<FooterConfig>
                    updateFooter(config)
                  }}
                  className="w-full"
                />
              </EditorErrorBoundary>
            </div>
          </>
        )}

        {/* Settings Tab */}
        {activeTab === 'settings' && (
          <>
            {/* Typography & Font Settings */}
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-4 text-sm flex items-center gap-2">
                <Settings className="h-4 w-4 text-primary-600" />
                {t('themeBuilderEditor.typography', 'Typography')}
              </h3>
              <EditorErrorBoundary componentType="font-picker">
                <FontPicker
                  component={{
                    id: 'typography',
                    type: 'custom',
                    title: t('themeBuilderEditor.fontSelection', 'Font Selection'),
                    displayOrder: 0,
                    enabled: true,
                    config: {
                      fontFamily: typography.fontFamily,
                      fontWeight: 400,
                      fontSize: typography.baseFontSize,
                      lineHeight: typography.lineHeight,
                      fontWeights: typography.fontWeights || [400, 700],
                      fontStyle: typography.fontStyle || 'normal',
                    },
                  }}
                  onChange={(component) => {
                    const config = component.config as any
                    updateTypography({
                      fontFamily: config.fontFamily as string,
                      baseFontSize: config.fontSize as number || typography.baseFontSize,
                      lineHeight: config.lineHeight as number || typography.lineHeight,
                      borderRadius: typography.borderRadius,
                      fontWeights: config.fontWeights as number[],
                      fontStyle: config.fontStyle as string,
                    })
                  }}
                  className="w-full"
                />
              </EditorErrorBoundary>
            </div>

            {/* Typography Preview */}
            <EditorErrorBoundary componentType="typography-preview">
              <TypographyPreview typography={typography} textColor={colors.text} />
            </EditorErrorBoundary>

            {/* Custom CSS */}
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-4 text-sm">
                {t('themeBuilderEditor.customCss', 'Custom CSS')}
              </h3>
              <EditorErrorBoundary componentType="custom-css-editor">
                <CustomCSSEditor
                  component={{
                    id: 'custom-css',
                    type: 'custom',
                    title: t('themeBuilderEditor.customCss', 'Custom CSS'),
                    displayOrder: 0,
                    enabled: true,
                    config: {
                      customCss: theme.customCss || '',
                      isEnabled: true,
                      validationEnabled: true,
                    },
                  }}
                  onChange={(component) => {
                    const config = component.config as any
                    updateCustomCss(config.customCss as string || '')
                  }}
                  className="w-full"
                />
              </EditorErrorBoundary>
            </div>
          </>
        )}

        {/* Component Edit Tab */}
        {activeTab === 'component-edit' && selectedComponent && (
          <EditorErrorBoundary componentType="component-editor">
            <>
              <div className="p-4 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
                <p className="text-sm text-blue-700 dark:text-blue-400">
                  <Edit3 className={`h-4 w-4 inline ${isRTL ? 'ml-2' : 'mr-2'}`} />
                  {t('themeBuilderEditor.editingComponent', 'Editing component')}:{' '}
                  <strong>{extractLabel(selectedComponent.title)}</strong>
                </p>
              </div>
              <ComponentEditor component={selectedComponent} />
            </>
          </EditorErrorBoundary>
        )}

        {/* History Tab */}
        {activeTab === 'history' && theme && (
          <>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-4 text-sm flex items-center gap-2">
                <Clock className="h-4 w-4 text-primary-600" />
                {t('themeBuilderEditor.themeVersionHistory', 'Theme Version History')}
              </h3>
              <ThemeHistory themeId={theme.id} />
            </div>
          </>
        )}
      </div>
    </div>
  )
}
