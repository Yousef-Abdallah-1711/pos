import { GlobalColors, TypographySettings, HeaderConfig, FooterConfig } from '@/types/theme'

/**
 * Hex color validation
 */
export const isValidHexColor = (color: string): boolean => {
  const hexRegex = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/
  return hexRegex.test(color)
}

/**
 * Validate color with detailed message
 */
export const validateColorWithMessage = (color: string): { valid: boolean; message?: string } => {
  if (!color || typeof color !== 'string') {
    return { valid: false, message: 'Color must be a string' }
  }

  if (!isValidHexColor(color)) {
    return { valid: false, message: 'Color must be in hex format (e.g., #FF0000)' }
  }

  return { valid: true }
}

/**
 * Convert hex to RGB
 */
export const hexToRgb = (hex: string): { r: number; g: number; b: number } | null => {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  return result
    ? {
        r: parseInt(result[1], 16),
        g: parseInt(result[2], 16),
        b: parseInt(result[3], 16),
      }
    : null
}

/**
 * Convert RGB to hex
 */
export const rgbToHex = (r: number, g: number, b: number): string => {
  return '#' + [r, g, b].map((x) => {
    const hex = x.toString(16)
    return hex.length === 1 ? '0' + hex : hex
  }).join('')
}

/**
 * Convert hex to HSL
 */
export const hexToHsl = (hex: string): string => {
  const rgb = hexToRgb(hex)
  if (!rgb) return hex

  const r = rgb.r / 255
  const g = rgb.g / 255
  const b = rgb.b / 255

  const max = Math.max(r, g, b)
  const min = Math.min(r, g, b)
  let h = 0
  let s = 0
  const l = (max + min) / 2

  if (max !== min) {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)

    switch (max) {
      case r:
        h = ((g - b) / d + (g < b ? 6 : 0)) / 6
        break
      case g:
        h = ((b - r) / d + 2) / 6
        break
      case b:
        h = ((r - g) / d + 4) / 6
        break
    }
  }

  return `${Math.round(h * 360)} ${Math.round(s * 100)}% ${Math.round(l * 100)}%`
}

/**
 * Convert HSL to hex
 */
export const hslToHex = (h: number, s: number, l: number): string => {
  s /= 100
  l /= 100

  const c = (1 - Math.abs(2 * l - 1)) * s
  const x = c * (1 - ((h / 60) % 2 - 1))
  const m = l - c / 2
  let r = 0
  let g = 0
  let b = 0

  if (h >= 0 && h < 60) {
    r = c
    g = x
    b = 0
  } else if (h >= 60 && h < 120) {
    r = x
    g = c
    b = 0
  } else if (h >= 120 && h < 180) {
    r = 0
    g = c
    b = x
  } else if (h >= 180 && h < 240) {
    r = 0
    g = x
    b = c
  } else if (h >= 240 && h < 300) {
    r = x
    g = 0
    b = c
  } else if (h >= 300 && h < 360) {
    r = c
    g = 0
    b = x
  }

  return rgbToHex(
    Math.round((r + m) * 255),
    Math.round((g + m) * 255),
    Math.round((b + m) * 255)
  )
}

/**
 * Calculate relative luminance of a color
 * Used for WCAG contrast checking
 */
export const getRelativeLuminance = (hex: string): number => {
  const rgb = hexToRgb(hex)
  if (!rgb) return 0

  const [r, g, b] = [rgb.r / 255, rgb.g / 255, rgb.b / 255]

  const luminance = [r, g, b].map((channel) => {
    return channel <= 0.03928 ? channel / 12.92 : Math.pow((channel + 0.055) / 1.055, 2.4)
  })

  return 0.2126 * luminance[0] + 0.7152 * luminance[1] + 0.0722 * luminance[2]
}

/**
 * Calculate WCAG contrast ratio between two colors
 */
export const getContrastRatio = (color1: string, color2: string): number => {
  const lum1 = getRelativeLuminance(color1)
  const lum2 = getRelativeLuminance(color2)

  const lighter = Math.max(lum1, lum2)
  const darker = Math.min(lum1, lum2)

  return (lighter + 0.05) / (darker + 0.05)
}

/**
 * Check if contrast ratio meets WCAG AA standard (4.5:1 for normal text)
 */
export const meetsWCAGAA = (foreground: string, background: string): boolean => {
  return getContrastRatio(foreground, background) >= 4.5
}

/**
 * Check if contrast ratio meets WCAG AAA standard (7:1 for normal text)
 */
export const meetsWCAGAAA = (foreground: string, background: string): boolean => {
  return getContrastRatio(foreground, background) >= 7
}

/**
 * Validate color palette for accessibility
 */
export const validateColorPalette = (colors: GlobalColors): {
  valid: boolean
  issues: string[]
} => {
  const issues: string[] = []

  // Check text color contrast with background
  const textContrastRatio = getContrastRatio(colors.text, colors.background)
  if (textContrastRatio < 4.5) {
    issues.push(
      `Text color contrast is ${textContrastRatio.toFixed(2)}:1, needs to be at least 4.5:1 for WCAG AA`
    )
  }

  // Check primary color contrast with background
  const primaryContrastRatio = getContrastRatio(colors.primary, colors.background)
  if (primaryContrastRatio < 3) {
    issues.push(
      `Primary color contrast is ${primaryContrastRatio.toFixed(2)}:1, should be at least 3:1`
    )
  }

  return {
    valid: issues.length === 0,
    issues,
  }
}

/**
 * Validate font size
 */
export const isValidFontSize = (size: number): boolean => {
  return size > 0 && size < 100
}

/**
 * Validate line height
 */
export const isValidLineHeight = (height: number): boolean => {
  return height > 0 && height < 5
}

/**
 * Validate border radius
 */
export const isValidBorderRadius = (radius: number): boolean => {
  return radius >= 0 && radius <= 50
}

/**
 * Validate typography settings
 */
export const validateTypography = (typography: TypographySettings): {
  valid: boolean
  errors: string[]
} => {
  const errors: string[] = []

  if (!typography.fontFamily || typography.fontFamily.trim().length === 0) {
    errors.push('Font family is required')
  }

  if (!isValidFontSize(typography.baseFontSize)) {
    errors.push('Font size must be between 1 and 99 pixels')
  }

  if (!isValidLineHeight(typography.lineHeight)) {
    errors.push('Line height must be between 0 and 5')
  }

  if (!isValidBorderRadius(typography.borderRadius)) {
    errors.push('Border radius must be between 0 and 50 pixels')
  }

  return {
    valid: errors.length === 0,
    errors,
  }
}

/**
 * Validate header configuration
 */
export const validateHeader = (header: HeaderConfig): {
  valid: boolean
  errors: string[]
} => {
  const errors: string[] = []

  if (!isValidHexColor(header.backgroundColor)) {
    errors.push('Header background color must be a valid hex color')
  }

  if (!isValidHexColor(header.textColor)) {
    errors.push('Header text color must be a valid hex color')
  }

  if (header.height <= 0 || header.height > 200) {
    errors.push('Header height must be between 1 and 200 pixels')
  }

  if (!header.navigationItems || header.navigationItems.length === 0) {
    errors.push('Header must have at least one navigation item')
  }

  return {
    valid: errors.length === 0,
    errors,
  }
}

/**
 * Validate footer configuration
 */
export const validateFooter = (footer: FooterConfig): {
  valid: boolean
  errors: string[]
} => {
  const errors: string[] = []

  if (!footer.companyName || footer.companyName.trim().length === 0) {
    errors.push('Company name is required')
  }

  if (!isValidHexColor(footer.backgroundColor)) {
    errors.push('Footer background color must be a valid hex color')
  }

  if (!isValidHexColor(footer.textColor)) {
    errors.push('Footer text color must be a valid hex color')
  }

  if (footer.linkColor && !isValidHexColor(footer.linkColor)) {
    errors.push('Footer link color must be a valid hex color')
  }

  if (!footer.copyrightText || footer.copyrightText.trim().length === 0) {
    errors.push('Copyright text is required')
  }

  return {
    valid: errors.length === 0,
    errors,
  }
}

/**
 * Generate slug from theme name
 */
export const generateSlug = (name: string): string => {
  return name
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/--+/g, '-')
}

/**
 * Sanitize CSS string (basic)
 */
export const sanitizeCss = (css: string): string => {
  // Remove script tags
  let sanitized = css.replace(/<script[^>]*>.*?<\/script>/gi, '')

  // Remove event handlers
  sanitized = sanitized.replace(/on\w+\s*=\s*["'][^"']*["']/gi, '')

  return sanitized.trim()
}

/**
 * Merge two color objects
 */
export const mergeColors = (base: GlobalColors, override: Partial<GlobalColors>): GlobalColors => {
  return { ...base, ...override }
}

/**
 * Clone theme data deeply
 */
export const cloneThemeData = <T extends Record<string, any>>(data: T): T => {
  return JSON.parse(JSON.stringify(data))
}

/**
 * Get colors that changed between two color objects
 */
export const getColorChanges = (
  oldColors: GlobalColors,
  newColors: GlobalColors
): Partial<GlobalColors> => {
  const changes: Partial<GlobalColors> = {}

  Object.keys(newColors).forEach((key) => {
    const colorKey = key as keyof GlobalColors
    if (oldColors[colorKey] !== newColors[colorKey]) {
      changes[colorKey] = newColors[colorKey]
    }
  })

  return changes
}

/**
 * Create color palette from a single primary color
 */
export const generatePaletteFromPrimary = (primaryColor: string): Partial<GlobalColors> => {
  if (!isValidHexColor(primaryColor)) {
    return {}
  }

  const rgb = hexToRgb(primaryColor)
  if (!rgb) return {}

  // Simple palette generation
  return {
    primary: primaryColor,
    secondary: rgbToHex(Math.round(rgb.r * 0.8), Math.round(rgb.g * 0.8), Math.round(rgb.b * 0.8)),
    accent: rgbToHex(Math.min(255, Math.round(rgb.r * 1.1)), Math.round(rgb.g * 0.6), Math.round(rgb.b * 0.6)),
  }
}

/**
 * Validate entire theme data
 */
export const validateCompleteTheme = (colors: GlobalColors, typography: TypographySettings, header: HeaderConfig, footer: FooterConfig) => {
  const colorValidation = validateColorPalette(colors)
  const typographyValidation = validateTypography(typography)
  const headerValidation = validateHeader(header)
  const footerValidation = validateFooter(footer)

  return {
    valid:
      colorValidation.valid &&
      typographyValidation.valid &&
      headerValidation.valid &&
      footerValidation.valid,
    errors: [
      ...colorValidation.issues,
      ...typographyValidation.errors,
      ...headerValidation.errors,
      ...footerValidation.errors,
    ],
  }
}
