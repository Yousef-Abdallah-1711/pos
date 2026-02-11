/**
 * Validation Utilities for Theme Builder
 * Validates user input before API calls to catch errors early
 */

import { debugWarn } from '@/lib/debug'

/**
 * Validates hex color format
 * @param color - Color string to validate (e.g., "#FF0000", "FF0000", "#FFF", "FFF")
 * @returns true if valid hex color, false otherwise
 *
 * Examples:
 * - validateHexColor("#3b82f6") → true
 * - validateHexColor("#FF0") → true
 * - validateHexColor("rgb(255,0,0)") → false
 * - validateHexColor("invalid") → false
 */
export function validateHexColor(color: string): boolean {
  if (!color || typeof color !== 'string') return false

  // Remove spaces
  const trimmed = color.trim()

  // Match: #RGB, #RRGGBB, or color names
  const hexPattern = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/

  return hexPattern.test(trimmed)
}

/**
 * Validates URL format
 * @param url - URL string to validate
 * @param allowEmpty - Whether to allow empty string (default: true)
 * @returns true if valid URL or empty (if allowed), false otherwise
 *
 * Examples:
 * - validateUrl("https://example.com/image.png") → true
 * - validateUrl("http://localhost:3000") → true
 * - validateUrl("") → true (if allowEmpty=true)
 * - validateUrl("not a url") → false
 */
export function validateUrl(url: string, allowEmpty = true): boolean {
  // Allow empty if flag is set
  if (!url && allowEmpty) return true
  if (!url || typeof url !== 'string') return false

  try {
    const urlObj = new URL(url)
    // Only allow http and https
    return urlObj.protocol === 'http:' || urlObj.protocol === 'https:'
  } catch {
    return false
  }
}

/**
 * Validates if string is not empty
 * @param value - String to validate
 * @param fieldName - Name of field (for error messages)
 * @returns true if not empty, false otherwise
 */
export function validateRequired(value: string | undefined, fieldName: string): boolean {
  if (!value || typeof value !== 'string' || value.trim() === '') {
    debugWarn(`⚠️ [Validation] ${fieldName} is required`)
    return false
  }
  return true
}

/**
 * Validates font family is one of allowed values
 * @param fontFamily - Font family name
 * @returns true if valid, false otherwise
 */
export function validateFontFamily(fontFamily: string): boolean {
  const allowedFonts = [
    'Inter',
    'Roboto',
    'Open Sans',
    'Playfair Display',
    'Montserrat',
    'Lato',
    'Poppins',
    'Source Sans Pro',
  ]

  return allowedFonts.includes(fontFamily)
}

/**
 * Validates number is within range
 * @param value - Number to validate
 * @param min - Minimum allowed value
 * @param max - Maximum allowed value
 * @param fieldName - Field name for logging
 * @returns true if valid, false otherwise
 */
export function validateNumberRange(
  value: number | undefined,
  min: number,
  max: number,
  fieldName: string
): boolean {
  if (value === undefined || typeof value !== 'number') {
    debugWarn(`⚠️ [Validation] ${fieldName} must be a number`)
    return false
  }

  if (value < min || value > max) {
    debugWarn(`⚠️ [Validation] ${fieldName} must be between ${min} and ${max}, got ${value}`)
    return false
  }

  return true
}

/**
 * Validates complete theme object
 * @param theme - Theme object to validate
 * @returns object with isValid boolean and errors array
 *
 * Example:
 * ```typescript
 * const result = validateTheme(myTheme);
 * if (!result.isValid) {
 *   console.error(result.errors); // ["Primary color is invalid", ...]
 * }
 * ```
 */
export function validateTheme(theme: any): { isValid: boolean; errors: string[] } {
  const errors: string[] = []

  // Validate required fields
  if (!validateRequired(theme?.name, 'Theme name')) {
    errors.push('Theme name is required')
  }

  if (!validateRequired(theme?.identity?.siteTitle, 'Site title')) {
    errors.push('Site title is required')
  }

  // Validate colors
  if (theme?.colors) {
    const colorFields = ['primary', 'secondary', 'accent', 'background', 'text', 'border', 'shadow']
    for (const colorField of colorFields) {
      const color = theme.colors[colorField]
      if (color && !validateHexColor(color)) {
        errors.push(`Invalid ${colorField} color: ${color}`)
      }
    }
  }

  // Validate typography
  if (theme?.typography) {
    if (!validateFontFamily(theme.typography.fontFamily)) {
      errors.push(`Invalid font family: ${theme.typography.fontFamily}`)
    }

    if (!validateNumberRange(theme.typography.baseFontSize, 8, 72, 'Base font size')) {
      errors.push('Base font size must be between 8 and 72')
    }

    if (!validateNumberRange(theme.typography.borderRadius, 0, 32, 'Border radius')) {
      errors.push('Border radius must be between 0 and 32')
    }

    if (!validateNumberRange(theme.typography.lineHeight, 1, 3, 'Line height')) {
      errors.push('Line height must be between 1 and 3')
    }
  }

  // Validate URLs
  if (theme?.identity?.logoUrl && !validateUrl(theme.identity.logoUrl)) {
    errors.push(`Invalid logo URL: ${theme.identity.logoUrl}`)
  }

  if (theme?.identity?.faviconUrl && !validateUrl(theme.identity.faviconUrl)) {
    errors.push(`Invalid favicon URL: ${theme.identity.faviconUrl}`)
  }

  if (theme?.header?.logoUrl && !validateUrl(theme.header.logoUrl)) {
    errors.push(`Invalid header logo URL: ${theme.header.logoUrl}`)
  }

  // Validate components
  if (Array.isArray(theme?.components)) {
    for (let i = 0; i < theme.components.length; i++) {
      const component = theme.components[i]

      if (!validateRequired(component?.id, `Component ${i} ID`)) {
        errors.push(`Component ${i}: ID is required`)
      }

      if (!validateRequired(component?.type, `Component ${i} type`)) {
        errors.push(`Component ${i}: Type is required`)
      }

      // Validate displayOrder is a non-negative integer
      if (typeof component?.displayOrder !== 'number' || component.displayOrder < 0) {
        errors.push(`Component ${i}: Invalid display order`)
      }
    }
  }

  return {
    isValid: errors.length === 0,
    errors,
  }
}

/**
 * Validates color contrast ratio for accessibility
 * @param color1 - First color in hex format
 * @param color2 - Second color in hex format
 * @returns true if colors have sufficient contrast (WCAG AA), false otherwise
 *
 * Note: WCAG AA requires minimum 4.5:1 contrast ratio for normal text,
 * 3:1 for large text
 */
export function validateColorContrast(color1: string, color2: string): boolean {
  try {
    const lum1 = getRelativeLuminance(color1)
    const lum2 = getRelativeLuminance(color2)

    const lighter = Math.max(lum1, lum2)
    const darker = Math.min(lum1, lum2)

    const contrastRatio = (lighter + 0.05) / (darker + 0.05)

    // WCAG AA standard is 4.5:1
    return contrastRatio >= 4.5
  } catch (err) {
    debugWarn('⚠️ [Validation] Could not calculate contrast ratio:', err)
    return true // Allow if can't calculate
  }
}

/**
 * Helper function to calculate relative luminance for color contrast
 * Based on WCAG 2.0 formula
 */
function getRelativeLuminance(hexColor: string): number {
  const color = hexColor.replace('#', '')
  const r = parseInt(color.substring(0, 2), 16) / 255
  const g = parseInt(color.substring(2, 4), 16) / 255
  const b = parseInt(color.substring(4, 6), 16) / 255

  const luminance =
    0.2126 * (r <= 0.03928 ? r / 12.92 : Math.pow((r + 0.055) / 1.055, 2.4)) +
    0.7152 * (g <= 0.03928 ? g / 12.92 : Math.pow((g + 0.055) / 1.055, 2.4)) +
    0.0722 * (b <= 0.03928 ? b / 12.92 : Math.pow((b + 0.055) / 1.055, 2.4))

  return luminance
}

/**
 * Validates component configuration
 * @param component - Component object to validate
 * @returns object with isValid boolean and errors array
 */
export function validateComponent(component: any): { isValid: boolean; errors: string[] } {
  const errors: string[] = []

  if (!validateRequired(component?.id, 'Component ID')) {
    errors.push('Component ID is required')
  }

  if (!validateRequired(component?.type, 'Component type')) {
    errors.push('Component type is required')
  }

  const validTypes = ['hero', 'products', 'why_us', 'contact', 'testimonials', 'cta']
  if (component?.type && !validTypes.includes(component.type)) {
    errors.push(`Invalid component type: ${component.type}. Must be one of: ${validTypes.join(', ')}`)
  }

  if (!validateRequired(component?.title, 'Component title')) {
    errors.push('Component title is required')
  }

  if (typeof component?.displayOrder !== 'number' || component.displayOrder < 0) {
    errors.push('Component display order must be a non-negative number')
  }

  if (typeof component?.enabled !== 'boolean') {
    errors.push('Component enabled must be a boolean')
  }

  return {
    isValid: errors.length === 0,
    errors,
  }
}

/**
 * Formats validation errors for display to user
 * @param errors - Array of error messages
 * @returns User-friendly error message
 */
export function formatValidationErrors(errors: string[]): string {
  if (errors.length === 0) return ''
  if (errors.length === 1) return errors[0]
  return `${errors.length} errors:\n• ${errors.join('\n• ')}`
}

export default {
  validateHexColor,
  validateUrl,
  validateRequired,
  validateFontFamily,
  validateNumberRange,
  validateTheme,
  validateColorContrast,
  validateComponent,
  formatValidationErrors,
}
