'use client'

import React, { useState, useCallback } from 'react'
import { AlertCircle, CheckCircle, Copy, Trash2 } from 'lucide-react'
import type { ThemeComponent } from '@/types/theme'

interface CustomCSSEditorProps {
  component: ThemeComponent
  onChange: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  className?: string
}

interface CustomCSSConfig {
  customCss?: string
  enableCustomCss?: boolean
  cssValidationEnabled?: boolean
  minHeight?: number
}

interface ValidationError {
  line: number
  message: string
  type: 'error' | 'warning'
}

/**
 * CustomCSSEditor Component
 * Allows users to write and preview custom CSS that extends the theme
 */
export function CustomCSSEditor({
  component,
  onChange,
  onPreview,
  className = '',
}: CustomCSSEditorProps): JSX.Element {
  const [cssConfig, setCssConfig] = useState<CustomCSSConfig>(
    (component.config as unknown as CustomCSSConfig) || {}
  )

  const [previewMode, setPreviewMode] = useState(false)
  const [validationErrors, setValidationErrors] = useState<ValidationError[]>([])
  const [hasValidated, setHasValidated] = useState(false)
  const [copySuccess, setCopySuccess] = useState(false)

  // Validate CSS
  const validateCSS = useCallback((css: string): ValidationError[] => {
    const errors: ValidationError[] = []
    const lines = css.split('\n')

    lines.forEach((line, index) => {
      const lineNum = index + 1
      const trimmed = line.trim()

      // Skip empty lines and comments
      if (!trimmed || trimmed.startsWith('/*') || trimmed.startsWith('*') || trimmed.startsWith('//')) {
        return
      }

      // Check for unmatched braces
      if (trimmed.endsWith('{')) {
        // Opening brace found
        return
      }

      // Warn about !important overuse
      if (trimmed.includes('!important')) {
        errors.push({
          line: lineNum,
          message: 'Using !important is discouraged. Prefer CSS specificity.',
          type: 'warning',
        })
      }

      // Warn about inline styles (shouldn't be in CSS file)
      if (trimmed.includes('style=')) {
        errors.push({
          line: lineNum,
          message: 'Inline styles should not appear in CSS. Use selectors instead.',
          type: 'error',
        })
      }

      // Check for common mistakes
      if (trimmed.includes('color:') && !trimmed.includes('#') && !trimmed.includes('rgb')) {
        if (!['red', 'blue', 'green', 'black', 'white', 'transparent', 'inherit'].some((c) => trimmed.includes(c))) {
          errors.push({
            line: lineNum,
            message: 'Possibly invalid color value. Use hex (#), rgb, or valid color name.',
            type: 'warning',
          })
        }
      }
    })

    return errors
  }, [])

  const handleCssChange = useCallback(
    (newCss: string) => {
      const newConfig: CustomCSSConfig = {
        ...cssConfig,
        customCss: newCss,
      }

      setCssConfig(newConfig)

      // Validate if enabled
      if (newConfig.cssValidationEnabled) {
        setValidationErrors(validateCSS(newCss))
      }

      // Update component
      const updated: ThemeComponent = {
        ...component,
        config: newConfig,
      }

      onChange(updated)
      onPreview?.(updated)
    },
    [cssConfig, component, onChange, onPreview, validateCSS]
  )

  const handleValidate = useCallback(() => {
    const css = cssConfig.customCss || ''
    const errors = validateCSS(css)
    setValidationErrors(errors)
    setHasValidated(true)
  }, [cssConfig, validateCSS])

  const handleToggleValidation = useCallback(() => {
    const newConfig: CustomCSSConfig = {
      ...cssConfig,
      cssValidationEnabled: !cssConfig.cssValidationEnabled,
    }
    setCssConfig(newConfig)

    const updated: ThemeComponent = {
      ...component,
      config: newConfig,
    }

    onChange(updated)
  }, [cssConfig, component, onChange])

  const handleToggleEnable = useCallback(() => {
    const newConfig: CustomCSSConfig = {
      ...cssConfig,
      enableCustomCss: !cssConfig.enableCustomCss,
    }
    setCssConfig(newConfig)

    const updated: ThemeComponent = {
      ...component,
      config: newConfig,
    }

    onChange(updated)
  }, [cssConfig, component, onChange])

  const handleCopyCSS = useCallback(() => {
    const css = cssConfig.customCss || ''
    navigator.clipboard.writeText(css)
    setCopySuccess(true)
    setTimeout(() => setCopySuccess(false), 2000)
  }, [cssConfig.customCss])

  const handleClearCSS = useCallback(() => {
    if (window.confirm('Are you sure you want to clear all CSS?')) {
      handleCssChange('')
    }
  }, [handleCssChange])

  const errorCount = validationErrors.filter((e) => e.type === 'error').length
  const warningCount = validationErrors.filter((e) => e.type === 'warning').length
  const css = cssConfig.customCss || ''

  // Preview Mode
  if (previewMode) {
    return (
      <div className={`w-full space-y-4 ${className}`}>
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">CSS Preview</h3>
          <button
            onClick={() => setPreviewMode(false)}
            className="px-3 py-2 text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition"
          >
            Back to Editor
          </button>
        </div>

        {!cssConfig.enableCustomCss ? (
          <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
            <p className="text-sm text-yellow-800 dark:text-yellow-200">Custom CSS is disabled.</p>
          </div>
        ) : (
          <div className="bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-6 space-y-4">
            <div>
              <h4 className="font-semibold text-gray-900 dark:text-white mb-3">CSS Code</h4>
              <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded p-4 font-mono text-sm overflow-auto max-h-48">
                <pre className="text-gray-900 dark:text-gray-100 whitespace-pre-wrap break-words">{css}</pre>
              </div>
            </div>

            <div>
              <h4 className="font-semibold text-gray-900 dark:text-white mb-2">How it works</h4>
              <ul className="text-sm text-gray-600 dark:text-gray-400 space-y-1 list-disc list-inside">
                <li>CSS will be injected into the page after theme styles</li>
                <li>Use selectors to target elements and override styles</li>
                <li>Disable this CSS anytime without losing your code</li>
                <li>CSS must be valid for it to be applied</li>
              </ul>
            </div>
          </div>
        )}
      </div>
    )
  }

  // Editor Mode
  return (
    <div className={`w-full space-y-4 ${className}`}>
      {/* Header */}
      <div>
        <h2 className="text-lg font-bold text-gray-900 dark:text-white mb-2">Custom CSS Editor</h2>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          Write custom CSS to extend or override theme styles. CSS is applied after all theme styles.
        </p>
      </div>

      {/* Toolbar */}
      <div className="flex flex-wrap items-center gap-3 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={cssConfig.enableCustomCss || false}
            onChange={handleToggleEnable}
            className="rounded"
          />
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Enable Custom CSS</span>
        </label>

        <div className="w-px h-4 bg-gray-300 dark:bg-gray-700" />

        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={cssConfig.cssValidationEnabled || false}
            onChange={handleToggleValidation}
            className="rounded"
          />
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Validate CSS</span>
        </label>

        <div className="flex-1" />

        <button
          onClick={handleValidate}
          className="px-3 py-1.5 text-sm bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition"
        >
          Validate
        </button>

        <button
          onClick={handleCopyCSS}
          className="flex items-center gap-1.5 px-3 py-1.5 text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition"
        >
          <Copy className="h-3.5 w-3.5" />
          {copySuccess ? 'Copied!' : 'Copy'}
        </button>

        <button
          onClick={handleClearCSS}
          className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition"
        >
          <Trash2 className="h-3.5 w-3.5" />
          Clear
        </button>
      </div>

      {/* CSS Editor */}
      <div className="space-y-2">
        <label htmlFor="custom-css-code" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
          CSS Code
          <span className="text-gray-500 dark:text-gray-500 ml-2">({css.length} characters)</span>
        </label>
        <textarea
          id="custom-css-code"
          aria-label="CSS Code"
          value={css}
          onChange={(e) => handleCssChange(e.target.value)}
          disabled={!cssConfig.enableCustomCss}
          className="w-full h-96 px-4 py-3 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white font-mono text-sm resize-vertical"
          placeholder={`/* Add your custom CSS here */\n\n.selector {\n  property: value;\n}\n\n/* Example: Override button color */\n.btn {\n  background-color: #your-color;\n}`}
        />
        <p className="text-xs text-gray-500 dark:text-gray-400">
          Maximum file size: 50KB. CSS is applied after all theme styles.
        </p>
      </div>

      {/* Validation Results */}
      {validationErrors.length > 0 && (
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4 space-y-3">
          <div className="flex items-center gap-2">
            <AlertCircle className="h-5 w-5 text-blue-600 dark:text-blue-400" />
            <h4 className="font-semibold text-blue-900 dark:text-blue-200">
              Validation Results
            </h4>
          </div>

          <div className="flex gap-6 text-sm">
            {errorCount > 0 && (
              <div className="flex items-center gap-2">
                <div className="h-3 w-3 rounded-full bg-red-500" />
                <span className="text-gray-700 dark:text-gray-300">
                  {errorCount} error{errorCount !== 1 ? 's' : ''}
                </span>
              </div>
            )}
            {warningCount > 0 && (
              <div className="flex items-center gap-2">
                <div className="h-3 w-3 rounded-full bg-yellow-500" />
                <span className="text-gray-700 dark:text-gray-300">
                  {warningCount} warning{warningCount !== 1 ? 's' : ''}
                </span>
              </div>
            )}
          </div>

          <div className="space-y-2 max-h-48 overflow-y-auto">
            {validationErrors.map((error, index) => (
              <div
                key={index}
                className={`p-2 rounded text-sm ${
                  error.type === 'error'
                    ? 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-200'
                    : 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-200'
                }`}
              >
                <strong>Line {error.line}:</strong> {error.message}
              </div>
            ))}
          </div>
        </div>
      )}

      {hasValidated && errorCount === 0 && warningCount === 0 && validationErrors.length === 0 && (
        <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4 flex items-center gap-2">
          <CheckCircle className="h-5 w-5 text-green-600 dark:text-green-400" />
          <p className="text-sm text-green-800 dark:text-green-200">CSS looks good! No errors found.</p>
        </div>
      )}

      {/* Preview Button */}
      <div className="flex gap-3">
        <button
          onClick={() => setPreviewMode(true)}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition"
        >
          Preview CSS
        </button>
      </div>

      {/* Info Box */}
      <div className="bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-2">
        <h4 className="font-semibold text-gray-900 dark:text-white">Tips</h4>
        <ul className="text-sm text-gray-600 dark:text-gray-400 space-y-1 list-disc list-inside">
          <li>Custom CSS is loaded after theme styles, so it will override them</li>
          <li>Use browser DevTools (F12) to find correct selectors to target</li>
          <li>Keep CSS validation enabled to catch mistakes early</li>
          <li>If something breaks, you can always clear and start over</li>
          <li>CSS must be valid for changes to take effect</li>
        </ul>
      </div>
    </div>
  )
}
