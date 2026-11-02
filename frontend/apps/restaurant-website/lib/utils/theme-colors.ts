'use client'

/**
 * Color conversion utilities for theme system
 * Converts hex colors to HSL format for CSS variables
 */

/**
 * Convert hex color to RGB
 * @param hex - Hex color string (e.g., "#b45309")
 * @returns [R, G, B] array with values 0-255
 */
export function hexToRgb(hex: string): [number, number, number] {
	const normalized = hex.trim().replace(/^#/, '')
	const expanded = normalized.length === 3
		? normalized.split('').map((digit) => `${digit}${digit}`).join('')
		: normalized
	const result = /^([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(expanded)
  if (!result) {
    throw new Error(`Invalid hex color: ${hex}`)
  }
  return [
    parseInt(result[1], 16),
    parseInt(result[2], 16),
    parseInt(result[3], 16),
  ]
}

/**
 * Convert RGB to HSL
 * @param r - Red (0-255)
 * @param g - Green (0-255)
 * @param b - Blue (0-255)
 * @returns HSL string in format "H S% L%" (e.g., "0 100% 50%")
 */
export function rgbToHsl(r: number, g: number, b: number): string {
  r /= 255
  g /= 255
  b /= 255

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

  const hue = Math.round(h * 360)
  const saturation = Math.round(s * 100)
  const lightness = Math.round(l * 100)

  return `${hue} ${saturation}% ${lightness}%`
}

/**
 * Convert hex to HSL
 * @param hex - Hex color string (e.g., "#b45309")
 * @returns HSL string (e.g., "30 100% 35%")
 */
export function hexToHsl(hex: string): string {
  const [r, g, b] = hexToRgb(hex)
  return rgbToHsl(r, g, b)
}

/**
 * Calculate relative luminance per WCAG spec
 * Used for contrast ratio calculation
 * @param hex - Hex color
 * @returns Luminance value 0-1
 */
export function getLuminance(hex: string): number {
  const [r, g, b] = hexToRgb(hex)

  // Convert to 0-1 range and apply gamma
  const [rs, gs, bs] = [r, g, b].map((c) => {
    c = c / 255
    return c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4)
  })

  // WCAG luminance formula
  return 0.2126 * rs + 0.7152 * gs + 0.0722 * bs
}

/**
 * Calculate contrast ratio between two colors
 * WCAG formula: (L1 + 0.05) / (L2 + 0.05) where L1 is lighter
 * @param foreground - Foreground color hex
 * @param background - Background color hex
 * @returns Contrast ratio 1-21
 */
export function getContrastRatio(foreground: string, background: string): number {
  const l1 = getLuminance(foreground)
  const l2 = getLuminance(background)
  const lighter = Math.max(l1, l2)
  const darker = Math.min(l1, l2)
  return (lighter + 0.05) / (darker + 0.05)
}

/**
 * Check if color pair meets WCAG AA standard (4.5:1)
 * @param foreground - Foreground color hex
 * @param background - Background color hex
 * @returns true if contrast >= 4.5:1
 */
export function meetsWCAGAA(foreground: string, background: string): boolean {
  return getContrastRatio(foreground, background) >= 4.5
}

/**
 * Check if color pair meets WCAG AAA standard (7:1)
 * @param foreground - Foreground color hex
 * @param background - Background color hex
 * @returns true if contrast >= 7:1
 */
export function meetsWCAGAAA(foreground: string, background: string): boolean {
  return getContrastRatio(foreground, background) >= 7
}

/**
 * Lighten a hex color by percentage
 * @param hex - Hex color
 * @param percent - Percentage to lighten (0-100)
 * @returns Lighter hex color
 */
export function lighten(hex: string, percent: number): string {
  const [r, g, b] = hexToRgb(hex)
  const hsl = rgbToHsl(r, g, b).split(' ')
  const l = parseInt(hsl[2]) + percent
  const newHsl = `hsl(${hsl[0]}, ${hsl[1]}, ${Math.min(100, l)}%)`
  return newHsl
}

/**
 * Darken a hex color by percentage
 * @param hex - Hex color
 * @param percent - Percentage to darken (0-100)
 * @returns Darker hex color
 */
export function darken(hex: string, percent: number): string {
  const [r, g, b] = hexToRgb(hex)
  const hsl = rgbToHsl(r, g, b).split(' ')
  const l = parseInt(hsl[2]) - percent
  const newHsl = `hsl(${hsl[0]}, ${hsl[1]}, ${Math.max(0, l)}%)`
  return newHsl
}

/**
 * Check if hex color is valid
 * @param hex - Hex color string
 * @returns true if valid hex color
 */
export function isValidHex(hex: string): boolean {
  try {
    hexToRgb(hex)
    return true
  } catch {
    return false
  }
}
