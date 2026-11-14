'use client'

import type { JSX } from 'react'

import React, { useState, useCallback } from 'react'
import type { ThemeComponent } from '@/types/theme'

interface GenericEditorProps {
  component: ThemeComponent
  onChange: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  className?: string
}

/**
 * Generic/Fallback Editor
 * Used for unknown component types or custom components
 * Provides basic JSON editing capabilities
 */
export function GenericEditor({
  component,
  onChange,
  onPreview,
  className = '',
}: GenericEditorProps): JSX.Element {
  const [config, setConfig] = useState<Record<string, any>>(component.config || {})
  const [previewMode, setPreviewMode] = useState(false)
  const [configJson, setConfigJson] = useState(JSON.stringify(config, null, 2))

  const handleConfigChange = useCallback(
    (newConfig: Record<string, any>) => {
      setConfig(newConfig)
      setConfigJson(JSON.stringify(newConfig, null, 2))
      onChange({
        ...component,
        config: newConfig,
      })
      onPreview?.({
        ...component,
        config: newConfig,
      })
    },
    [component, onChange, onPreview]
  )

  const handleJsonChange = (json: string) => {
    setConfigJson(json)
    try {
      const parsed = JSON.parse(json)
      setConfig(parsed)
      handleConfigChange(parsed)
    } catch (error) {
      // Invalid JSON, just update display without updating config
      setConfigJson(json)
    }
  }

  if (previewMode) {
    return (
      <div className={`w-full space-y-4 ${className}`}>
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Preview</h3>
          <button
            onClick={() => setPreviewMode(false)}
            className="px-3 py-2 text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition"
          >
            Back to Settings
          </button>
        </div>

        <div className="bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
          <div className="text-center text-gray-500 dark:text-gray-400">
            <p className="mb-2">Generic Component Preview</p>
            <p className="text-sm">Type: {component.type}</p>
            {Object.keys(config).length > 0 && (
              <div className="mt-4 text-left bg-gray-100 dark:bg-gray-800 p-4 rounded font-mono text-xs overflow-auto max-h-96">
                <pre>{JSON.stringify(config, null, 2)}</pre>
              </div>
            )}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className={`w-full space-y-6 ${className}`}>
      {/* Header */}
      <div>
        <h2 className="text-lg font-bold text-gray-900 dark:text-white mb-2">
          Edit {component.type}
        </h2>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          Type: {component.type} (using generic editor)
        </p>
      </div>

      {/* Basic Info */}
      <div className="space-y-4">
        <h3 className="font-semibold text-gray-900 dark:text-white">Basic Information</h3>
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Component Title
          </label>
          <input
            type="text"
            value={component.title}
            onChange={(e) =>
              onChange({
                ...component,
                title: e.target.value,
              })
            }
            className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white"
          />
        </div>

        <div>
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={component.enabled}
              onChange={(e) =>
                onChange({
                  ...component,
                  enabled: e.target.checked,
                })
              }
              className="rounded"
            />
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
              Enable this component
            </span>
          </label>
        </div>
      </div>

      {/* Configuration JSON Editor */}
      <div className="space-y-4">
        <h3 className="font-semibold text-gray-900 dark:text-white">Configuration</h3>
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Configuration JSON
          </label>
          <textarea
            value={configJson}
            onChange={(e) => handleJsonChange(e.target.value)}
            className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white font-mono text-sm"
            rows={12}
            placeholder={`{\n  "key": "value"\n}`}
          />
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">
            Edit the JSON configuration directly. Changes are applied in real-time.
          </p>
        </div>
      </div>

      {/* Preview Button */}
      <div className="flex gap-3">
        <button
          onClick={() => setPreviewMode(true)}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition"
        >
          Preview
        </button>
      </div>
    </div>
  )
}
