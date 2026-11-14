
import type { JSX } from 'react'
/**
 * ContactEditor Component
 * Specialized editor for contact form and information sections
 */

'use client'

import React, { useState, useCallback } from 'react'
import type { ThemeComponent } from '@/types/theme'

interface ContactEditorProps {
  component: ThemeComponent
  onChange: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  className?: string
}

interface ContactField {
  id: string
  name: string
  type: 'text' | 'email' | 'phone' | 'textarea' | 'select'
  label: string
  placeholder: string
  required: boolean
  order: number
}

interface ContactConfig {
  formTitle?: string
  formDescription?: string
  formFields?: ContactField[]
  submitButtonText?: string
  submitButtonColor?: string
  submitButtonHoverColor?: string
  successMessage?: string
  successMessageColor?: string
  errorMessage?: string
  errorMessageColor?: string
  formBackgroundColor?: string
  inputBackgroundColor?: string
  inputBorderColor?: string
  inputTextColor?: string
  labelColor?: string
  placeholderColor?: string
  showMap?: boolean
  mapEmbedUrl?: string
  mapHeight?: number
  contactInfoPosition?: 'left' | 'right' | 'bottom'
  showPhone?: boolean
  phone?: string
  showEmail?: boolean
  email?: string
  showAddress?: boolean
  address?: string
  showHours?: boolean
  businessHours?: string
  socialLinks?: Array<{ platform: string; url: string }>
  emailNotifications?: string
  autoReply?: boolean
  autoReplySubject?: string
  autoReplyMessage?: string
}

/**
 * ContactEditor Component
 */
export function ContactEditor({
  component,
  onChange,
  onPreview,
  className = '',
}: ContactEditorProps): JSX.Element {
  const [config, setConfig] = useState<ContactConfig>(
    (component.config as unknown as ContactConfig) || {}
  )

  const [previewMode, setPreviewMode] = useState(false)
  const [editingFieldId, setEditingFieldId] = useState<string | null>(null)

  /**
   * Handle config change
   */
  const handleConfigChange = useCallback(
    (newConfig: Partial<ContactConfig>) => {
      const updatedConfig = { ...config, ...newConfig }
      setConfig(updatedConfig)

      const updatedComponent: ThemeComponent = {
        ...component,
        config: updatedConfig,
      }

      onChange(updatedComponent)

      if (onPreview) {
        onPreview(updatedComponent)
      }
    },
    [config, component, onChange, onPreview]
  )

  const addFormField = () => {
    const newField: ContactField = {
      id: `field-${Date.now()}`,
      name: 'field_name',
      type: 'text',
      label: 'Field Label',
      placeholder: 'Placeholder text',
      required: false,
      order: (config.formFields || []).length,
    }

    const updated = [...(config.formFields || []), newField]
    handleConfigChange({ formFields: updated })
  }

  const updateField = (fieldId: string, updates: Partial<ContactField>) => {
    const updated = (config.formFields || []).map((f) =>
      f.id === fieldId ? { ...f, ...updates } : f
    )
    handleConfigChange({ formFields: updated })
  }

  const deleteField = (fieldId: string) => {
    const updated = (config.formFields || []).filter((f) => f.id !== fieldId)
    handleConfigChange({ formFields: updated })
  }

  const moveFieldUp = (fieldId: string) => {
    const fields = [...(config.formFields || [])]
    const index = fields.findIndex((f) => f.id === fieldId)
    if (index > 0) {
      [fields[index], fields[index - 1]] = [fields[index - 1], fields[index]]
      const reordered = fields.map((f, i) => ({ ...f, order: i }))
      handleConfigChange({ formFields: reordered })
    }
  }

  const moveFieldDown = (fieldId: string) => {
    const fields = [...(config.formFields || [])]
    const index = fields.findIndex((f) => f.id === fieldId)
    if (index < fields.length - 1) {
      [fields[index], fields[index + 1]] = [fields[index + 1], fields[index]]
      const reordered = fields.map((f, i) => ({ ...f, order: i }))
      handleConfigChange({ formFields: reordered })
    }
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Tabs */}
      <div className="flex gap-2 border-b border-gray-200 dark:border-gray-700">
        <button
          onClick={() => setPreviewMode(false)}
          className={`px-4 py-2 font-medium transition-colors ${
            !previewMode
              ? 'border-b-2 border-blue-600 text-blue-600'
              : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
          }`}
        >
          Settings
        </button>
        <button
          onClick={() => setPreviewMode(true)}
          className={`px-4 py-2 font-medium transition-colors ${
            previewMode
              ? 'border-b-2 border-blue-600 text-blue-600'
              : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
          }`}
        >
          Preview
        </button>
      </div>

      {!previewMode && (
        <div className="space-y-6">
          {/* Form Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Form Settings</h3>

            {/* Form Title */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Form Title
              </label>
              <input
                type="text"
                value={config.formTitle || ''}
                onChange={(e) => handleConfigChange({ formTitle: e.target.value })}
                placeholder="Get in Touch"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* Form Description */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Form Description
              </label>
              <textarea
                value={config.formDescription || ''}
                onChange={(e) => handleConfigChange({ formDescription: e.target.value })}
                placeholder="We'd love to hear from you"
                rows={2}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* Submit Button Text */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Submit Button Text
              </label>
              <input
                type="text"
                value={config.submitButtonText || 'Send Message'}
                onChange={(e) => handleConfigChange({ submitButtonText: e.target.value })}
                placeholder="Send Message"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* Success Message */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Success Message
              </label>
              <input
                type="text"
                value={config.successMessage || 'Thank you for your message!'}
                onChange={(e) => handleConfigChange({ successMessage: e.target.value })}
                placeholder="Thank you for your message!"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          {/* Form Fields Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="font-semibold text-gray-900 dark:text-white">Form Fields</h3>
              <button
                onClick={addFormField}
                className="px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-medium transition-colors"
              >
                + Add Field
              </button>
            </div>

            <div className="space-y-3 max-h-64 overflow-y-auto">
              {(config.formFields || []).map((field, index) => (
                <div
                  key={field.id}
                  className="border border-gray-300 dark:border-gray-600 rounded-lg p-3 space-y-2"
                >
                  {editingFieldId === field.id ? (
                    <>
                      <input
                        type="text"
                        value={field.label}
                        onChange={(e) => updateField(field.id, { label: e.target.value })}
                        placeholder="Field Label"
                        className="w-full px-2 py-1 border border-gray-300 dark:border-gray-600 rounded
                          bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm
                          focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <select
                        value={field.type}
                        onChange={(e) =>
                          updateField(field.id, {
                            type: e.target.value as
                              | 'text'
                              | 'email'
                              | 'phone'
                              | 'textarea'
                              | 'select',
                          })
                        }
                        className="w-full px-2 py-1 border border-gray-300 dark:border-gray-600 rounded
                          bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm
                          focus:outline-none focus:ring-2 focus:ring-blue-500"
                      >
                        <option value="text">Text</option>
                        <option value="email">Email</option>
                        <option value="phone">Phone</option>
                        <option value="textarea">Textarea</option>
                        <option value="select">Select</option>
                      </select>
                      <input
                        type="text"
                        value={field.placeholder}
                        onChange={(e) => updateField(field.id, { placeholder: e.target.value })}
                        placeholder="Placeholder"
                        className="w-full px-2 py-1 border border-gray-300 dark:border-gray-600 rounded
                          bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm
                          focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <label className="flex items-center gap-2 cursor-pointer">
                        <input
                          type="checkbox"
                          checked={field.required}
                          onChange={(e) => updateField(field.id, { required: e.target.checked })}
                          className="w-4 h-4 rounded"
                        />
                        <span className="text-sm text-gray-700 dark:text-gray-300">Required field</span>
                      </label>
                      <div className="flex gap-2">
                        <button
                          onClick={() => setEditingFieldId(null)}
                          className="flex-1 px-3 py-1 rounded text-sm bg-blue-600 hover:bg-blue-700 text-white font-medium transition-colors"
                        >
                          Done
                        </button>
                        <button
                          onClick={() => deleteField(field.id)}
                          className="flex-1 px-3 py-1 rounded text-sm bg-red-600 hover:bg-red-700 text-white font-medium transition-colors"
                        >
                          Delete
                        </button>
                      </div>
                    </>
                  ) : (
                    <>
                      <div className="flex items-center justify-between">
                        <div className="flex-1">
                          <p className="font-medium text-gray-900 dark:text-white">
                            {field.label || 'Untitled'}
                          </p>
                          <p className="text-xs text-gray-600 dark:text-gray-400">
                            {field.type}
                            {field.required && ' (required)'}
                          </p>
                        </div>
                        <div className="flex gap-1">
                          <button
                            onClick={() => moveFieldUp(field.id)}
                            className="px-2 py-1 rounded text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
                            title="Move up"
                          >
                            ↑
                          </button>
                          <button
                            onClick={() => moveFieldDown(field.id)}
                            className="px-2 py-1 rounded text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
                            title="Move down"
                          >
                            ↓
                          </button>
                          <button
                            onClick={() => setEditingFieldId(field.id)}
                            className="px-3 py-1 rounded text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
                          >
                            Edit
                          </button>
                        </div>
                      </div>
                    </>
                  )}
                </div>
              ))}
            </div>
          </div>

          {/* Colors Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Colors</h3>

            {/* Submit Button Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Submit Button Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.submitButtonColor || '#3b82f6'}
                  onChange={(e) => handleConfigChange({ submitButtonColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.submitButtonColor || '#3b82f6'}
                  onChange={(e) => handleConfigChange({ submitButtonColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Input Background */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Input Background Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.inputBackgroundColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ inputBackgroundColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.inputBackgroundColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ inputBackgroundColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Input Border Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Input Border Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.inputBorderColor || '#d1d5db'}
                  onChange={(e) => handleConfigChange({ inputBorderColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.inputBorderColor || '#d1d5db'}
                  onChange={(e) => handleConfigChange({ inputBorderColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>
          </div>

          {/* Contact Info Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Contact Information</h3>

            {/* Show Phone */}
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showPhone ?? true}
                onChange={(e) => handleConfigChange({ showPhone: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Phone Number</span>
            </label>

            {config.showPhone && (
              <input
                type="tel"
                value={config.phone || ''}
                onChange={(e) => handleConfigChange({ phone: e.target.value })}
                placeholder="+1 (555) 123-4567"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            )}

            {/* Show Email */}
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showEmail ?? true}
                onChange={(e) => handleConfigChange({ showEmail: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Email Address</span>
            </label>

            {config.showEmail && (
              <input
                type="email"
                value={config.email || ''}
                onChange={(e) => handleConfigChange({ email: e.target.value })}
                placeholder="contact@example.com"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            )}

            {/* Show Address */}
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showAddress ?? true}
                onChange={(e) => handleConfigChange({ showAddress: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Address</span>
            </label>

            {config.showAddress && (
              <textarea
                value={config.address || ''}
                onChange={(e) => handleConfigChange({ address: e.target.value })}
                placeholder="123 Main St, City, State 12345"
                rows={2}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            )}

            {/* Show Hours */}
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showHours ?? true}
                onChange={(e) => handleConfigChange({ showHours: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Business Hours</span>
            </label>

            {config.showHours && (
              <textarea
                value={config.businessHours || ''}
                onChange={(e) => handleConfigChange({ businessHours: e.target.value })}
                placeholder="Mon-Fri: 9AM-5PM&#10;Sat-Sun: Closed"
                rows={2}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            )}
          </div>

          {/* Map Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Map Settings</h3>

            {/* Show Map */}
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showMap ?? false}
                onChange={(e) => handleConfigChange({ showMap: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Map</span>
            </label>

            {config.showMap && (
              <>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Map Embed URL (Google Maps Embed)
                  </label>
                  <input
                    type="url"
                    value={config.mapEmbedUrl || ''}
                    onChange={(e) => handleConfigChange({ mapEmbedUrl: e.target.value })}
                    placeholder="https://www.google.com/maps/embed?pb=..."
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                      bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm
                      focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">
                    Get embed code from Google Maps
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Map Height (px)
                  </label>
                  <input
                    type="number"
                    value={config.mapHeight || 400}
                    onChange={(e) => handleConfigChange({ mapHeight: parseInt(e.target.value) })}
                    min="200"
                    max="600"
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                      bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                      focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </>
            )}
          </div>
        </div>
      )}

      {previewMode && (
        <div className="rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          <div className="bg-gradient-to-b from-gray-50 to-gray-100 dark:from-gray-800 dark:to-gray-900 p-8">
            <div className="max-w-4xl mx-auto">
              {/* Header */}
              <div className="text-center mb-8">
                <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
                  {config.formTitle || 'Contact Us'}
                </h2>
                {config.formDescription && (
                  <p className="text-gray-600 dark:text-gray-400">{config.formDescription}</p>
                )}
              </div>

              <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                {/* Form */}
                <div className="contact-preview-form rounded-lg p-6 shadow-lg">
                  <form className="space-y-4">
                    {(config.formFields || []).map((field) => (
                      <div key={field.id}>
                        <label className="contact-preview-label block text-sm font-medium mb-1">
                          {field.label}
                          {field.required && <span className="text-red-500">*</span>}
                        </label>
                        {field.type === 'textarea' ? (
                          <textarea
                            placeholder={field.placeholder}
                            rows={4}
                            className="contact-preview-input w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                          />
                        ) : (
                          <input
                            type={field.type}
                            placeholder={field.placeholder}
                            className="contact-preview-input w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                          />
                        )}
                      </div>
                    ))}

                    <button className="contact-preview-submit w-full text-white font-medium py-2 rounded-lg hover:opacity-90 transition-opacity">
                      {config.submitButtonText || 'Send Message'}
                    </button>
                  </form>
                </div>

                {/* Contact Info */}
                <div className="space-y-8">
                  {config.showPhone && (
                    <div>
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                        📞 Phone
                      </h3>
                      <p className="text-gray-600 dark:text-gray-400">{config.phone}</p>
                    </div>
                  )}

                  {config.showEmail && (
                    <div>
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                        ✉️ Email
                      </h3>
                      <p className="text-gray-600 dark:text-gray-400">{config.email}</p>
                    </div>
                  )}

                  {config.showAddress && (
                    <div>
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                        📍 Address
                      </h3>
                      <p className="text-gray-600 dark:text-gray-400 whitespace-pre-line">
                        {config.address}
                      </p>
                    </div>
                  )}

                  {config.showHours && (
                    <div>
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                        🕐 Hours
                      </h3>
                      <p className="text-gray-600 dark:text-gray-400 whitespace-pre-line">
                        {config.businessHours}
                      </p>
                    </div>
                  )}
                </div>
              </div>

              {/* Map */}
              {config.showMap && config.mapEmbedUrl && (
                <div className="mt-8">
                  <iframe
                    width="100%"
                    height={config.mapHeight || 400}
                    frameBorder="0"
                    className="contact-preview-map"
                    src={config.mapEmbedUrl}
                    allowFullScreen
                    loading="lazy"
                    referrerPolicy="no-referrer-when-downgrade"
                  />
                </div>
              )}
              <style jsx>{`
                .contact-preview-form {
                  background-color: ${config.formBackgroundColor || '#ffffff'};
                }

                .contact-preview-label {
                  color: ${config.labelColor || '#374151'};
                }

                .contact-preview-input {
                  background-color: ${config.inputBackgroundColor || '#ffffff'};
                  border-color: ${config.inputBorderColor || '#d1d5db'};
                  color: ${config.inputTextColor || '#111827'};
                }

                .contact-preview-submit {
                  background-color: ${config.submitButtonColor || '#3b82f6'};
                }

                .contact-preview-map {
                  border: 0;
                  border-radius: 0.5rem;
                }
              `}</style>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default ContactEditor
