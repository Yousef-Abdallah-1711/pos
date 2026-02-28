'use client'

import { useState, useRef, useEffect } from 'react'
import { usePathname } from 'next/navigation'
import { Copy, Check } from 'lucide-react'
import { debugLog } from '@/lib/debug'
import { getLocaleFromPath } from '@/lib/translations'

interface ColorPickerProps {
  label: string
  value: string
  onChange: (color: string) => void
  description?: string
}

export function ColorPicker({ label, value, onChange, description }: ColorPickerProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'
  const tx = (en: string, ar: string) => (isRTL ? ar : en)
  const [copied, setCopied] = useState(false)
  const [inputValue, setInputValue] = useState(value)
  const [showPicker, setShowPicker] = useState(false)
  const pickerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    setInputValue(value)
  }, [value])

  const validateHex = (color: string): boolean => {
    const hexRegex = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/
    return hexRegex.test(color)
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value
    setInputValue(newValue)
    if (validateHex(newValue)) {
      onChange(newValue)
    }
  }

  const handleColorInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newColor = e.target.value
    debugLog(`🎨 [ColorPicker] Color selected: ${newColor}`)
    onChange(newColor)
    setInputValue(newColor)
    // Close picker after selection
    setShowPicker(false)
  }

  const copyToClipboard = () => {
    navigator.clipboard.writeText(value)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (pickerRef.current && !pickerRef.current.contains(e.target as Node)) {
        setShowPicker(false)
      }
    }

    if (showPicker) {
      document.addEventListener('mousedown', handleClickOutside)
      return () => document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [showPicker])

  return (
    <div dir={isRTL ? 'rtl' : 'ltr'} className={`space-y-2 ${isRTL ? 'text-right' : ''}`}>
      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
        {label}
      </label>

      <div className="flex items-center gap-3">
        {/* Color Preview */}
        <div
          className="relative w-16 h-16 rounded-lg border-2 border-gray-200 dark:border-gray-600 cursor-pointer hover:shadow-lg transition overflow-hidden"
          onClick={() => setShowPicker(!showPicker)}
        >
          <input
            type="color"
            value={value}
            onChange={handleColorInputChange}
            className="w-full h-full cursor-pointer"
            title={value}
          />
        </div>

        {/* Color Input and Actions */}
        <div className="flex-1 space-y-2">
          <input
            type="text"
            value={inputValue}
            onChange={handleInputChange}
            placeholder={tx('#000000', '#000000')}
            className="w-full px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm font-mono"
          />

          <button
            onClick={copyToClipboard}
            className="text-xs text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 flex items-center gap-1 transition"
          >
            {copied ? (
              <>
                <Check className="h-3 w-3" /> {tx('Copied', 'تم النسخ')}
              </>
            ) : (
              <>
                <Copy className="h-3 w-3" /> {tx('Copy', 'نسخ')}
              </>
            )}
          </button>
        </div>
      </div>

      {description && (
        <p className="text-xs text-gray-500 dark:text-gray-400">{description}</p>
      )}
    </div>
  )
}
