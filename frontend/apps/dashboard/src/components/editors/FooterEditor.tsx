'use client'

import type { JSX } from 'react'

import React, { useState, useCallback, useMemo } from 'react'
import { usePathname } from 'next/navigation'
import { Plus, Trash2, ChevronUp, ChevronDown, Mail, MapPin, Phone } from 'lucide-react'
import type { ThemeComponent, FooterConfig as ThemeFooterConfig } from '@/types/theme'
import { useUpdateFooter, useThemeFooter } from '@/stores'
import { getLocaleFromPath } from '@/lib/translations'

interface FooterEditorProps {
  component?: ThemeComponent
  onChange?: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  className?: string
  useStore?: boolean
}

interface LinkItem {
  label: string
  href: string
}

interface SocialLink {
  id: string
  platform: 'facebook' | 'twitter' | 'instagram' | 'linkedin' | 'youtube'
  url: string
  order: number
}

interface FooterSection {
  id: string
  title: string
  links: LinkItem[]
  order: number
}

interface FooterConfig extends Omit<ThemeFooterConfig, 'footerSections' | 'socialLinks'> {
  footerSections?: FooterSection[]
  socialLinks: SocialLink[]
  layout?: 'compact' | 'expanded' | 'minimal'
  legalLinks?: LinkItem[]
}

/** Sections and social links saved without ids cannot be edited (lookups are by id): give them stable ones. */
function withFooterIds(config: ThemeFooterConfig): Omit<ThemeFooterConfig, 'footerSections' | 'socialLinks'> & { footerSections: FooterSection[]; socialLinks: SocialLink[] } {
  return {
    ...config,
    footerSections: (config.footerSections ?? []).map((section, index) => ({
      ...section,
      id: section.id ?? `section-${index + 1}`,
      order: section.order ?? index + 1,
    })),
    socialLinks: (config.socialLinks ?? []).map((link, index) => ({
      ...link,
      id: link.id ?? `social-${index + 1}`,
      order: link.order ?? index + 1,
      platform: link.platform as SocialLink['platform'],
    })),
  }
}

const SOCIAL_PLATFORMS = ['facebook', 'twitter', 'instagram', 'linkedin', 'youtube'] as const

/**
 * FooterEditor Component
 * Configure website footer with company info, sections, social links, and styling
 * Supports both local state (component prop) and themeBuilderStore (useStore prop)
 */
export function FooterEditor({
  component,
  onChange,
  onPreview,
  className = '',
  useStore = false,
}: FooterEditorProps): JSX.Element {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'
  const tx = (en: string, ar: string) => (isRTL ? ar : en)

  const storeFooter = useThemeFooter()
  const updateFooter = useUpdateFooter()

  const initialFooterConfig: FooterConfig = useMemo(() => {
    if (useStore && storeFooter) {
      return {
        ...withFooterIds(storeFooter),
        layout: 'expanded',
        legalLinks: [
          { label: 'Privacy Policy', href: '/privacy' },
          { label: 'Terms of Service', href: '/terms' },
        ],
      }
    }
    if (component?.config) {
      return withFooterIds(component.config as unknown as ThemeFooterConfig) as FooterConfig
    }
    return {
      companyName: '',
      showLinks: true,
      footerSections: [],
      socialLinks: [],
      backgroundColor: '#1f2937',
      textColor: '#ffffff',
      linkColor: '#3b82f6',
      padding: 48,
      copyrightText: '© 2024 Your Company. All rights reserved.',
      showLegal: true,
      legalLinks: [
        { label: 'Privacy Policy', href: '/privacy' },
        { label: 'Terms of Service', href: '/terms' },
      ],
      columns: 3,
      layout: 'expanded',
      showBackToTop: true,
    }
  }, [useStore, storeFooter, component])

  const [footerConfig, setFooterConfig] = useState<FooterConfig>(initialFooterConfig)

  const [previewMode, setPreviewMode] = useState(false)
  const [editingSection, setEditingSection] = useState<string | null>(null)

  const notifyConfigChange = useCallback((nextConfig: FooterConfig) => {
    if (useStore || !component || !onChange) return

    const updated: ThemeComponent = { ...component, config: nextConfig }
    onChange(updated)
    onPreview?.(updated)
  }, [component, onChange, onPreview, useStore])

  const handleConfigChange = useCallback(
    (key: keyof FooterConfig, value: any) => {
      const newConfig: FooterConfig = {
        ...footerConfig,
        [key]: value,
      }

      setFooterConfig(newConfig)

      if (useStore) {
        updateFooter({
          [key]: value,
        } as Partial<ThemeFooterConfig>)
      } else if (component && onChange) {
        const updated: ThemeComponent = {
          ...component,
          config: newConfig,
        }
        onChange(updated)
        onPreview?.(updated)
      }
    },
    [footerConfig, component, onChange, onPreview, useStore, updateFooter]
  )

  const addSection = useCallback(() => {
    const newSection: FooterSection = {
      id: `section-${Date.now()}`,
      title: tx('New Section', 'قسم جديد'),
      links: [],
      order: (footerConfig.footerSections?.length || 0) + 1,
    }

    const updatedSections = [...(footerConfig.footerSections || []), newSection]
    setFooterConfig({
      ...footerConfig,
      footerSections: updatedSections,
    })
    notifyConfigChange({ ...footerConfig, footerSections: updatedSections })

    if (useStore) {
      updateFooter({ footerSections: updatedSections })
    }
  }, [footerConfig, useStore, updateFooter, notifyConfigChange])

  const deleteSection = useCallback(
    (id: string) => {
      const updatedSections = (footerConfig.footerSections || []).filter((s) => s.id !== id)
      setFooterConfig({
        ...footerConfig,
        footerSections: updatedSections,
      })
      notifyConfigChange({ ...footerConfig, footerSections: updatedSections })

      if (useStore) {
        updateFooter({ footerSections: updatedSections })
      }
    },
    [footerConfig, useStore, updateFooter, notifyConfigChange]
  )

  const updateSection = useCallback(
    (id: string, title: string, links: LinkItem[]) => {
      const updatedSections = (footerConfig.footerSections || []).map((s) =>
        s.id === id ? { ...s, title, links } : s
      )
      setFooterConfig({
        ...footerConfig,
        footerSections: updatedSections,
      })
      notifyConfigChange({ ...footerConfig, footerSections: updatedSections })

      if (useStore) {
        updateFooter({ footerSections: updatedSections })
      }
    },
    [footerConfig, useStore, updateFooter, notifyConfigChange]
  )

  const addSocialLink = useCallback(() => {
    const newLink: SocialLink = {
      id: `social-${Date.now()}`,
      platform: 'facebook',
      url: '',
      order: (footerConfig.socialLinks?.length || 0) + 1,
    }

    const updatedLinks = [...(footerConfig.socialLinks || []), newLink]
    setFooterConfig({
      ...footerConfig,
      socialLinks: updatedLinks,
    })
    notifyConfigChange({ ...footerConfig, socialLinks: updatedLinks })

    if (useStore) {
      updateFooter({ socialLinks: updatedLinks })
    }
  }, [footerConfig, useStore, updateFooter, notifyConfigChange])

  const deleteSocialLink = useCallback(
    (id: string) => {
      const updatedLinks = (footerConfig.socialLinks || []).filter((link) => link.id !== id)
      setFooterConfig({
        ...footerConfig,
        socialLinks: updatedLinks,
      })
      notifyConfigChange({ ...footerConfig, socialLinks: updatedLinks })

      if (useStore) {
        updateFooter({ socialLinks: updatedLinks })
      }
    },
    [footerConfig, useStore, updateFooter, notifyConfigChange]
  )

  const updateSocialLink = useCallback(
    (id: string, platform: string, url: string) => {
      const updatedLinks = (footerConfig.socialLinks || []).map((link) =>
        link.id === id ? { ...link, platform: platform as any, url } : link
      )
      setFooterConfig({
        ...footerConfig,
        socialLinks: updatedLinks,
      })
      notifyConfigChange({ ...footerConfig, socialLinks: updatedLinks })

      if (useStore) {
        updateFooter({ socialLinks: updatedLinks })
      }
    },
    [footerConfig, useStore, updateFooter, notifyConfigChange]
  )

  if (previewMode) {
    return (
      <div dir={isRTL ? 'rtl' : 'ltr'} className={`w-full space-y-4 ${className} ${isRTL ? 'text-right' : ''}`}>
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{tx('Footer Preview', 'معاينة التذييل')}</h3>
          <button
            onClick={() => setPreviewMode(false)}
            className="px-3 py-2 text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition"
          >
            {tx('Back to Settings', 'العودة إلى الإعدادات')}
          </button>
        </div>

        <div
          className="footer-preview-root border border-gray-200 dark:border-gray-700 rounded-lg"
        >
          <div className="grid grid-cols-3 gap-8 max-w-6xl mx-auto mb-8">
            {/* Company Info */}
            <div>
              {footerConfig.companyName && <h4 className="font-semibold mb-2">{footerConfig.companyName}</h4>}
              {footerConfig.companyDescription && (
                <p className="text-sm opacity-80 mb-4">{footerConfig.companyDescription}</p>
              )}
              <div className="space-y-2 text-sm opacity-75">
                {footerConfig.phone && (
                  <div className="flex items-center gap-2">
                    <Phone className="h-4 w-4" />
                    {footerConfig.phone}
                  </div>
                )}
                {footerConfig.email && (
                  <div className="flex items-center gap-2">
                    <Mail className="h-4 w-4" />
                    {footerConfig.email}
                  </div>
                )}
                {footerConfig.address && (
                  <div className="flex items-start gap-2">
                    <MapPin className="h-4 w-4 mt-0.5" />
                    {footerConfig.address}
                  </div>
                )}
              </div>
            </div>

            {/* Footer Sections */}
            {(footerConfig.footerSections || []).slice(0, 2).map((section) => (
              <div key={section.id}>
                <h4 className="font-semibold mb-4">{section.title}</h4>
                <ul className="space-y-2">
                  {section.links.map((link, idx) => (
                    <li key={idx}>
                      <a
                        href={link.href}
                        className="footer-preview-link text-sm hover:opacity-80"
                      >
                        {link.label}
                      </a>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>

          {/* Social Links */}
          {(footerConfig.socialLinks || []).length > 0 && (
            <div className="border-t border-gray-700 pt-6 mb-6">
              <div className="flex gap-4 justify-center">
                {(footerConfig.socialLinks || []).map((link) => (
                  <a
                    key={link.id}
                    href={link.url}
                    className="footer-preview-link text-sm hover:opacity-80"
                  >
                    {link.platform.toUpperCase()}
                  </a>
                ))}
              </div>
            </div>
          )}

          {/* Copyright */}
          <div className="border-t border-gray-700 pt-6">
            <p className="text-sm opacity-75 text-center">{footerConfig.copyrightText}</p>
            {footerConfig.showLegal && (footerConfig.legalLinks || []).length > 0 && (
              <div className="flex gap-4 justify-center mt-4 text-xs opacity-60">
                {footerConfig.legalLinks?.map((link, idx) => (
                  <a key={idx} href={link.href} className="hover:opacity-100">
                    {link.label}
                  </a>
                ))}
              </div>
            )}
          </div>
        </div>
        <style jsx>{`
          .footer-preview-root {
            background-color: ${footerConfig.backgroundColor || '#1f2937'};
            color: ${footerConfig.textColor || '#ffffff'};
            padding: ${footerConfig.padding || 48}px 24px;
          }

          .footer-preview-link {
            color: ${footerConfig.linkColor || '#3b82f6'};
          }
        `}</style>
      </div>
    )
  }

  return (
    <div dir={isRTL ? 'rtl' : 'ltr'} className={`w-full space-y-6 ${className} ${isRTL ? 'text-right' : ''}`}>
      {/* Header */}
      <div>
        <h2 className="text-lg font-bold text-gray-900 dark:text-white mb-2">{tx('Footer Configuration', 'إعدادات التذييل')}</h2>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          {tx('Configure your website footer with company info, sections, and social links.', 'قم بإعداد تذييل موقعك مع معلومات الشركة والأقسام وروابط التواصل.')}
        </p>
      </div>

      {/* Company Information */}
      <div className="space-y-4 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <h3 className="font-semibold text-gray-900 dark:text-white">{tx('Company Information', 'معلومات الشركة')}</h3>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            {tx('Company Name', 'اسم الشركة')}
          </label>
          <input
            type="text"
            value={footerConfig.companyName || ''}
            onChange={(e) => handleConfigChange('companyName', e.target.value)}
            placeholder="Your Company"
            className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            {tx('Description', 'الوصف')}
          </label>
          <textarea
            value={footerConfig.companyDescription || ''}
            onChange={(e) => handleConfigChange('companyDescription', e.target.value)}
            placeholder="Brief company description"
            rows={2}
            className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
          />
        </div>

        <div className="grid grid-cols-1 gap-3">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {tx('Phone', 'الهاتف')}
            </label>
            <input
              type="tel"
              value={footerConfig.phone || ''}
              onChange={(e) => handleConfigChange('phone', e.target.value)}
              placeholder="+1 (555) 123-4567"
              className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {tx('Email', 'البريد الإلكتروني')}
            </label>
            <input
              type="email"
              value={footerConfig.email || ''}
              onChange={(e) => handleConfigChange('email', e.target.value)}
              placeholder="contact@example.com"
              className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {tx('Address', 'العنوان')}
            </label>
            <textarea
              value={footerConfig.address || ''}
              onChange={(e) => handleConfigChange('address', e.target.value)}
              placeholder="123 Main St, City, State 12345"
              rows={2}
              className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
            />
          </div>
        </div>
      </div>

      {/* Footer Sections */}
      <div className="space-y-4 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold text-gray-900 dark:text-white">{tx('Footer Sections', 'أقسام التذييل')}</h3>
          <button
            onClick={addSection}
            className="flex items-center gap-1.5 px-3 py-1.5 text-sm bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition"
          >
            <Plus className="h-4 w-4" />
            {tx('Add Section', 'إضافة قسم')}
          </button>
        </div>

        <div className="space-y-2 max-h-64 overflow-y-auto">
          {(footerConfig.footerSections || []).map((section) => (
            <div key={section.id} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
              {editingSection === section.id ? (
                <div className="space-y-2">
                  <input
                    type="text"
                    value={section.title}
                    onChange={(e) => updateSection(section.id, e.target.value, section.links)}
                    className="w-full px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white text-sm"
                  />
                  <button
                    onClick={() => setEditingSection(null)}
                    className="px-2 py-1 text-sm bg-green-600 hover:bg-green-700 text-white rounded transition"
                  >
                    {tx('Done', 'تم')}
                  </button>
                </div>
              ) : (
                <div className="flex items-center justify-between">
                  <div
                    className="flex-1 cursor-pointer"
                    onClick={() => setEditingSection(section.id)}
                  >
                    <div className="font-medium text-gray-900 dark:text-white">{section.title}</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">{section.links.length} {tx('links', 'روابط')}</div>
                  </div>
                  <button
                    onClick={() => deleteSection(section.id)}
                    aria-label={`Delete footer section ${section.title}`}
                    className="p-1 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Social Links */}
      <div className="space-y-4 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold text-gray-900 dark:text-white">{tx('Social Media Links', 'روابط التواصل الاجتماعي')}</h3>
          <button
            onClick={addSocialLink}
            className="flex items-center gap-1.5 px-3 py-1.5 text-sm bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition"
          >
            <Plus className="h-4 w-4" />
            {tx('Add Link', 'إضافة رابط')}
          </button>
        </div>

        <div className="space-y-2">
          {(footerConfig.socialLinks || []).map((link) => (
            <div key={link.id} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3 flex gap-2">
              <select
                value={link.platform}
                aria-label="Social platform"
                onChange={(e) => updateSocialLink(link.id, e.target.value, link.url)}
                className="px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white text-sm"
              >
                {SOCIAL_PLATFORMS.map((p) => (
                  <option key={p} value={p}>
                    {p.charAt(0).toUpperCase() + p.slice(1)}
                  </option>
                ))}
              </select>
              <input
                type="url"
                value={link.url}
                onChange={(e) => updateSocialLink(link.id, link.platform, e.target.value)}
                placeholder="https://..."
                aria-label={`${link.platform} URL`}
                className="flex-1 px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white text-sm"
              />
              <button
                onClick={() => deleteSocialLink(link.id)}
                aria-label={`Delete ${link.platform} link`}
                className="p-1 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Legal Section */}
      <div className="space-y-4 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <h3 className="font-semibold text-gray-900 dark:text-white">{tx('Legal & Copyright', 'القانونية وحقوق النشر')}</h3>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            {tx('Copyright Text', 'نص حقوق النشر')}
          </label>
          <input
            type="text"
            value={footerConfig.copyrightText || ''}
            onChange={(e) => handleConfigChange('copyrightText', e.target.value)}
            placeholder="© 2024 Your Company. All rights reserved."
            className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm"
          />
        </div>

        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={footerConfig.showLegal || false}
            onChange={(e) => handleConfigChange('showLegal', e.target.checked)}
            className="rounded"
          />
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">{tx('Show Legal Links', 'إظهار الروابط القانونية')}</span>
        </label>
      </div>

      {/* Styling */}
      <div className="space-y-4 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <h3 className="font-semibold text-gray-900 dark:text-white">{tx('Styling', 'التنسيق')}</h3>

        <div className="grid grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {tx('Background', 'الخلفية')}
            </label>
            <input
              type="color"
              value={footerConfig.backgroundColor || '#1f2937'}
              onChange={(e) => handleConfigChange('backgroundColor', e.target.value)}
              className="w-full h-10 border border-gray-200 dark:border-gray-600 rounded-lg cursor-pointer"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {tx('Text Color', 'لون النص')}
            </label>
            <input
              type="color"
              value={footerConfig.textColor || '#ffffff'}
              onChange={(e) => handleConfigChange('textColor', e.target.value)}
              className="w-full h-10 border border-gray-200 dark:border-gray-600 rounded-lg cursor-pointer"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {tx('Link Color', 'لون الروابط')}
            </label>
            <input
              type="color"
              value={footerConfig.linkColor || '#3b82f6'}
              onChange={(e) => handleConfigChange('linkColor', e.target.value)}
              className="w-full h-10 border border-gray-200 dark:border-gray-600 rounded-lg cursor-pointer"
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            {tx('Padding', 'الحشو')}: {footerConfig.padding || 48}px
          </label>
          <input
            type="range"
            min="24"
            max="80"
            value={footerConfig.padding || 48}
            onChange={(e) => handleConfigChange('padding', parseInt(e.target.value))}
            className="w-full"
          />
        </div>

        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={footerConfig.showBackToTop || false}
            onChange={(e) => handleConfigChange('showBackToTop', e.target.checked)}
            className="rounded"
          />
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">{tx('Show Back to Top', 'إظهار العودة للأعلى')}</span>
        </label>
      </div>

      {/* Preview Button */}
      <div className="flex gap-3">
        <button
          onClick={() => setPreviewMode(true)}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition"
        >
          {tx('Preview Footer', 'معاينة التذييل')}
        </button>
      </div>
    </div>
  )
}
