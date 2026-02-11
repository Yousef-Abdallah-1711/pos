/**
 * Google Fonts Service
 * Manages fetching, caching, and managing Google Fonts data
 */

import { debugError } from '@/lib/debug'

export interface GoogleFont {
  id: string
  family: string
  category: 'serif' | 'sans-serif' | 'monospace' | 'display' | 'handwriting'
  variants: string[]
  weights: number[]
  italics: boolean
  popularity?: number
}

// Popular Google Fonts for quick access (cached locally)
const POPULAR_FONTS: GoogleFont[] = [
  {
    id: 'roboto',
    family: 'Roboto',
    category: 'sans-serif',
    variants: ['300', '400', '500', '700', '900'],
    weights: [300, 400, 500, 700, 900],
    italics: true,
    popularity: 1,
  },
  {
    id: 'open-sans',
    family: 'Open Sans',
    category: 'sans-serif',
    variants: ['300', '400', '600', '700', '800'],
    weights: [300, 400, 600, 700, 800],
    italics: true,
    popularity: 2,
  },
  {
    id: 'lato',
    family: 'Lato',
    category: 'sans-serif',
    variants: ['300', '400', '700', '900'],
    weights: [300, 400, 700, 900],
    italics: true,
    popularity: 3,
  },
  {
    id: 'poppins',
    family: 'Poppins',
    category: 'sans-serif',
    variants: ['400', '500', '600', '700', '800', '900'],
    weights: [400, 500, 600, 700, 800, 900],
    italics: false,
    popularity: 4,
  },
  {
    id: 'inter',
    family: 'Inter',
    category: 'sans-serif',
    variants: ['300', '400', '500', '600', '700', '800', '900'],
    weights: [300, 400, 500, 600, 700, 800, 900],
    italics: false,
    popularity: 5,
  },
  {
    id: 'playfair-display',
    family: 'Playfair Display',
    category: 'serif',
    variants: ['400', '700', '900'],
    weights: [400, 700, 900],
    italics: true,
    popularity: 6,
  },
  {
    id: 'merriweather',
    family: 'Merriweather',
    category: 'serif',
    variants: ['300', '400', '700', '900'],
    weights: [300, 400, 700, 900],
    italics: true,
    popularity: 7,
  },
  {
    id: 'montserrat',
    family: 'Montserrat',
    category: 'sans-serif',
    variants: ['300', '400', '500', '600', '700', '800', '900'],
    weights: [300, 400, 500, 600, 700, 800, 900],
    italics: true,
    popularity: 8,
  },
  {
    id: 'raleway',
    family: 'Raleway',
    category: 'sans-serif',
    variants: ['300', '400', '500', '600', '700', '800', '900'],
    weights: [300, 400, 500, 600, 700, 800, 900],
    italics: true,
    popularity: 9,
  },
  {
    id: 'ubuntu',
    family: 'Ubuntu',
    category: 'sans-serif',
    variants: ['300', '400', '500', '700'],
    weights: [300, 400, 500, 700],
    italics: true,
    popularity: 10,
  },
  {
    id: 'source-code-pro',
    family: 'Source Code Pro',
    category: 'monospace',
    variants: ['400', '500', '600', '700'],
    weights: [400, 500, 600, 700],
    italics: false,
    popularity: 11,
  },
  {
    id: 'fira-mono',
    family: 'Fira Mono',
    category: 'monospace',
    variants: ['400', '700'],
    weights: [400, 700],
    italics: false,
    popularity: 12,
  },
]

// Store for recent fonts
const STORAGE_KEY = 'theme-builder-recent-fonts'
const MAX_RECENT_FONTS = 5

export const googleFontsService = {
  /**
   * Get popular Google Fonts
   */
  getPopularFonts: (): GoogleFont[] => {
    return POPULAR_FONTS
  },

  /**
   * Get fonts by category
   */
  getFontsByCategory: (category: GoogleFont['category']): GoogleFont[] => {
    return POPULAR_FONTS.filter((font) => font.category === category)
  },

  /**
   * Search fonts by name
   */
  searchFonts: (query: string): GoogleFont[] => {
    const lowerQuery = query.toLowerCase()
    return POPULAR_FONTS.filter(
      (font) =>
        font.family.toLowerCase().includes(lowerQuery) || font.id.toLowerCase().includes(lowerQuery)
    )
  },

  /**
   * Get font by ID
   */
  getFont: (fontId: string): GoogleFont | undefined => {
    return POPULAR_FONTS.find((f) => f.id === fontId)
  },

  /**
   * Get Google Fonts CSS URL
   */
  getFontCssUrl: (fontFamily: string, weights: number[] = [400]): string => {
    const family = fontFamily.replace(/\s+/g, '+')
    const uniqueWeights = [...new Set(weights)].sort((a, b) => a - b)
    const weightParams = uniqueWeights.join(';')
    return `https://fonts.googleapis.com/css2?family=${family}:wght@${weightParams}&display=swap`
  },

  /**
   * Get font family for CSS
   */
  getFontFamilyCss: (fontFamily: string): string => {
    return `'${fontFamily}', sans-serif`
  },

  /**
   * Add font to recent list
   */
  addRecentFont: (font: GoogleFont): void => {
    try {
      const recent = JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]') as GoogleFont[]
      const filtered = recent.filter((f) => f.id !== font.id)
      const updated = [font, ...filtered].slice(0, MAX_RECENT_FONTS)
      localStorage.setItem(STORAGE_KEY, JSON.stringify(updated))
    } catch (error) {
      debugError('Failed to save recent font:', error)
    }
  },

  /**
   * Get recently used fonts
   */
  getRecentFonts: (): GoogleFont[] => {
    try {
      return JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]') as GoogleFont[]
    } catch (error) {
      debugError('Failed to load recent fonts:', error)
      return []
    }
  },

  /**
   * Load font CSS
   */
  loadFont: (fontFamily: string, weights: number[] = [400]): Promise<void> => {
    return new Promise((resolve, reject) => {
      const link = document.createElement('link')
      link.href = googleFontsService.getFontCssUrl(fontFamily, weights)
      link.rel = 'stylesheet'
      link.onload = () => resolve()
      link.onerror = () => reject(new Error(`Failed to load font: ${fontFamily}`))
      document.head.appendChild(link)
    })
  },
}
