/**
 * Zod validation schemas for theme data
 * Mirrors backend validation rules from theme_validator.go
 */

import { z } from 'zod'

// ==========================================
// Color Validation
// ==========================================

/**
 * Validates hex color format (#RGB, #RRGGBB, or #RRGGBBAA)
 * Also accepts colors without hash prefix
 */
const hexColorSchema = z
  .string()
  .refine(
    (color) => {
      if (!color) return false
      // Remove hash if present
      const hex = color.replace(/^#/, '')
      // Check length (3, 6, or 8 chars)
      if (![3, 6, 8].includes(hex.length)) return false
      // Check valid hex characters
      return /^[0-9A-Fa-f]*$/.test(hex)
    },
    { message: 'Color must be in hex format (#RGB, #RRGGBB, or #RRGGBBAA)' }
  )
  .refine(
    (color) => {
      // Detect SQL injection in color
      const sqlPatterns = /(\b(DROP|DELETE|INSERT|UPDATE|SELECT)\b|--|\/\*|\*\/|;)/i
      return !sqlPatterns.test(color)
    },
    { message: 'Invalid color format detected' }
  )

// ==========================================
// Typography Validation
// ==========================================

const allowedFonts = [
  'Inter',
  'Roboto',
  'Open Sans',
  'Lato',
  'Montserrat',
  'Raleway',
  'Playfair Display',
  'Poppins',
  'Nunito',
  'Source Sans Pro',
  'Ubuntu',
  'Lora',
  'Merriweather',
  'Georgia',
  'Times New Roman',
  'Arial',
  'Helvetica',
  'Verdana',
  'Courier New',
  'Comic Sans MS',
]

const fontFamilySchema = z
  .string()
  .min(1, 'Font family is required')
  .refine((font) => allowedFonts.includes(font), {
    message: `Font must be one of: ${allowedFonts.join(', ')}`,
  })

const fontSizeSchema = z
  .number()
  .min(10, 'Font size must be at least 10px')
  .max(24, 'Font size must be at most 24px')

const borderRadiusSchema = z
  .number()
  .min(0, 'Border radius must be at least 0px')
  .max(50, 'Border radius must be at most 50px')

const lineHeightSchema = z
  .number()
  .min(1.0, 'Line height must be at least 1.0')
  .max(3.0, 'Line height must be at most 3.0')

export const typographySchema = z.object({
  fontFamily: fontFamilySchema,
  baseFontSize: fontSizeSchema,
  borderRadius: borderRadiusSchema,
  lineHeight: lineHeightSchema,
})

// ==========================================
// URL Validation
// ==========================================

const urlSchema = z
  .string()
  .refine(
    (url) => {
      if (!url) return true // Allow empty URLs if not required
      // Must start with http:// or https://
      return /^https?:\/\/.+/.test(url)
    },
    { message: 'URL must start with http:// or https://' }
  )
  .refine(
    (url) => {
      if (!url) return true
      // Block javascript: and data: URLs (XSS prevention)
      return !/^(javascript:|data:)/.test(url)
    },
    { message: 'Invalid URL protocol' }
  )

// ==========================================
// Text Field Validation
// ==========================================

const siteNameSchema = z
  .string()
  .min(1, 'Theme name is required')
  .max(100, 'Theme name must be at most 100 characters')
  .refine(
    (name) => {
      // SQL injection prevention
      const sqlPatterns = /(\b(DROP|DELETE|INSERT|UPDATE|SELECT)\b|--|\/\*|\*\/|xp_|sp_|;)/i
      return !sqlPatterns.test(name)
    },
    { message: 'Theme name contains invalid characters' }
  )

const siteTitleSchema = z
  .string()
  .min(1, 'Site title is required')
  .max(255, 'Site title must be at most 255 characters')

const descriptionSchema = z
  .string()
  .max(1000, 'Description must be at most 1000 characters')
  .optional()

// ==========================================
// Color Set Validation
// ==========================================

export const colorsSchema = z.object({
  primary: hexColorSchema.refine(
    (color) => color.length > 0,
    { message: 'Primary color is required' }
  ),
  secondary: hexColorSchema.refine(
    (color) => color.length > 0,
    { message: 'Secondary color is required' }
  ),
  accent: hexColorSchema.refine(
    (color) => color.length > 0,
    { message: 'Accent color is required' }
  ),
  background: hexColorSchema.refine(
    (color) => color.length > 0,
    { message: 'Background color is required' }
  ),
  text: hexColorSchema.refine(
    (color) => color.length > 0,
    { message: 'Text color is required' }
  ),
  border: hexColorSchema.refine(
    (color) => color.length > 0,
    { message: 'Border color is required' }
  ),
  shadow: hexColorSchema.refine(
    (color) => color.length > 0,
    { message: 'Shadow color is required' }
  ),
})

// ==========================================
// Website Identity Validation
// ==========================================

export const identitySchema = z.object({
  siteTitle: siteTitleSchema,
  logoUrl: urlSchema.optional().or(z.literal('')),
  faviconUrl: urlSchema.optional().or(z.literal('')),
  domain: z.string().optional(),
})

// ==========================================
// Theme Creation Validation
// ==========================================

export const createThemeSchema = z.object({
  name: siteNameSchema,
  description: descriptionSchema,
  colors: colorsSchema,
  typography: typographySchema,
  identity: identitySchema,
})

export type CreateThemeInput = z.infer<typeof createThemeSchema>

// ==========================================
// Theme Update Validation (partial)
// ==========================================

export const updateThemeSchema = z.object({
  name: siteNameSchema.optional(),
  description: descriptionSchema,
  colors: colorsSchema.partial().optional(),
  typography: typographySchema.partial().optional(),
  identity: identitySchema.partial().optional(),
})

export type UpdateThemeInput = z.infer<typeof updateThemeSchema>

// ==========================================
// Header Validation
// ==========================================

export const headerConfigSchema = z.object({
  logoUrl: urlSchema.optional().or(z.literal('')),
  showLogo: z.boolean(),
  navigationItems: z
    .array(
      z.object({
        label: z.string().min(1, 'Label is required'),
        href: z.string().min(1, 'URL is required'),
      })
    )
    .optional(),
  backgroundColor: hexColorSchema.optional(),
  textColor: hexColorSchema.optional(),
  height: z.number().min(40).max(200).optional(),
})

// ==========================================
// Footer Validation
// ==========================================

export const footerConfigSchema = z.object({
  companyName: z.string().min(1, 'Company name is required'),
  address: z.string().optional(),
  phone: z.string().optional(),
  email: z.string().email('Valid email required').optional().or(z.literal('')),
  copyrightText: z.string().optional(),
  socialLinks: z
    .array(
      z.object({
        platform: z.string(),
        url: urlSchema,
      })
    )
    .optional(),
  backgroundColor: hexColorSchema.optional(),
  textColor: hexColorSchema.optional(),
  showLinks: z.boolean().optional(),
})

// ==========================================
// Validation Error Type
// ==========================================

export interface ValidationError {
  field: string
  message: string
  code: string
}

export interface ValidationResult {
  isValid: boolean
  errors: ValidationError[]
}

// ==========================================
// Validation Helper Functions
// ==========================================

/**
 * Validates a theme creation request
 */
export function validateThemeCreate(data: unknown): ValidationResult {
  try {
    createThemeSchema.parse(data)
    return { isValid: true, errors: [] }
  } catch (error) {
    if (error instanceof z.ZodError) {
      const errors: ValidationError[] = error.issues.map((err) => ({
        field: err.path.join('.'),
        message: err.message,
        code: getErrorCode(err.path.join('.'), err.message),
      }))
      return { isValid: false, errors }
    }
    return {
      isValid: false,
      errors: [
        {
          field: 'unknown',
          message: 'Validation failed',
          code: 'VALIDATION_ERROR',
        },
      ],
    }
  }
}

/**
 * Validates a theme update request (partial)
 */
export function validateThemeUpdate(data: unknown): ValidationResult {
  try {
    updateThemeSchema.parse(data)
    return { isValid: true, errors: [] }
  } catch (error) {
    if (error instanceof z.ZodError) {
      const errors: ValidationError[] = error.issues.map((err) => ({
        field: err.path.join('.'),
        message: err.message,
        code: getErrorCode(err.path.join('.'), err.message),
      }))
      return { isValid: false, errors }
    }
    return {
      isValid: false,
      errors: [
        {
          field: 'unknown',
          message: 'Validation failed',
          code: 'VALIDATION_ERROR',
        },
      ],
    }
  }
}

/**
 * Maps validation messages to error codes (matching backend)
 */
function getErrorCode(field: string, message: string): string {
  if (message.includes('hex format') || message.includes('Invalid color')) {
    return 'INVALID_HEX_COLOR'
  }
  if (message.includes('Font must be')) {
    return 'FONT_NOT_ALLOWED'
  }
  if (message.includes('must be')) {
    if (message.includes('at least') || message.includes('at most')) {
      return 'OUT_OF_RANGE'
    }
  }
  if (message.includes('required')) {
    return 'REQUIRED_FIELD'
  }
  if (message.includes('URL')) {
    return 'INVALID_URL_PROTOCOL'
  }
  if (message.includes('invalid characters')) {
    return 'SQL_INJECTION_DETECTED'
  }
  return 'VALIDATION_ERROR'
}

/**
 * Validates individual fields
 */
export const fieldValidators = {
  color: (color: string) => {
    try {
      hexColorSchema.parse(color)
      return { isValid: true, error: null }
    } catch (error) {
      if (error instanceof z.ZodError) {
        return {
          isValid: false,
          error: error.issues[0]?.message || 'Invalid color',
        }
      }
      return { isValid: false, error: 'Validation failed' }
    }
  },

  fontSize: (size: number) => {
    try {
      fontSizeSchema.parse(size)
      return { isValid: true, error: null }
    } catch (error) {
      if (error instanceof z.ZodError) {
        return {
          isValid: false,
          error: error.issues[0]?.message || 'Invalid font size',
        }
      }
      return { isValid: false, error: 'Validation failed' }
    }
  },

  url: (url: string) => {
    try {
      urlSchema.parse(url)
      return { isValid: true, error: null }
    } catch (error) {
      if (error instanceof z.ZodError) {
        return {
          isValid: false,
          error: error.issues[0]?.message || 'Invalid URL',
        }
      }
      return { isValid: false, error: 'Validation failed' }
    }
  },

  font: (font: string) => {
    try {
      fontFamilySchema.parse(font)
      return { isValid: true, error: null }
    } catch (error) {
      if (error instanceof z.ZodError) {
        return {
          isValid: false,
          error: error.issues[0]?.message || 'Invalid font',
        }
      }
      return { isValid: false, error: 'Validation failed' }
    }
  },

  name: (name: string) => {
    try {
      siteNameSchema.parse(name)
      return { isValid: true, error: null }
    } catch (error) {
      if (error instanceof z.ZodError) {
        return {
          isValid: false,
          error: error.issues[0]?.message || 'Invalid name',
        }
      }
      return { isValid: false, error: 'Validation failed' }
    }
  },
}

// ==========================================
// Bilingual Content Validation
// ==========================================

/**
 * Validates bilingual text content with language-specific rules
 */
export const bilingualTextSchema = z.object({
  en: z
    .string()
    .min(1, 'English text is required')
    .max(500, 'English text must be at most 500 characters')
    .refine(
      (text) => {
        // SQL injection prevention
        const sqlPatterns = /(\b(DROP|DELETE|INSERT|UPDATE|SELECT)\b|--|\/\*|\*\/|xp_|sp_|;)/i
        return !sqlPatterns.test(text)
      },
      { message: 'English text contains invalid characters' }
    ),
  ar: z
    .string()
    .min(1, 'Arabic text is required')
    .max(500, 'Arabic text must be at most 500 characters')
    .refine(
      (text) => {
        // SQL injection prevention
        const sqlPatterns = /(\b(DROP|DELETE|INSERT|UPDATE|SELECT)\b|--|\/\*|\*\/|xp_|sp_|;)/i
        return !sqlPatterns.test(text)
      },
      { message: 'Arabic text contains invalid characters' }
    ),
})

/**
 * Optional bilingual text (nullable)
 */
export const optionalBilingualTextSchema = bilingualTextSchema.nullable().optional()

/**
 * Validates bilingual component structure
 */
export const bilingualComponentSchema = z.object({
  title: bilingualTextSchema,
  subtitle: optionalBilingualTextSchema,
  description: optionalBilingualTextSchema,
  buttonText: optionalBilingualTextSchema,
})

/**
 * Field validators for bilingual content
 */
export const bilingualFieldValidators = {
  /**
   * Validate English text field
   */
  englishText: (text: string) => {
    try {
      z.string()
        .min(1, 'English text is required')
        .max(500, 'English text must be at most 500 characters')
        .parse(text)
      return { isValid: true, error: null }
    } catch (error) {
      if (error instanceof z.ZodError) {
        return {
          isValid: false,
          error: error.issues[0]?.message || 'Invalid English text',
        }
      }
      return { isValid: false, error: 'Validation failed' }
    }
  },

  /**
   * Validate Arabic text field
   */
  arabicText: (text: string) => {
    try {
      z.string()
        .min(1, 'Arabic text is required')
        .max(500, 'Arabic text must be at most 500 characters')
        .parse(text)
      return { isValid: true, error: null }
    } catch (error) {
      if (error instanceof z.ZodError) {
        return {
          isValid: false,
          error: error.issues[0]?.message || 'Invalid Arabic text',
        }
      }
      return { isValid: false, error: 'Validation failed' }
    }
  },

  /**
   * Validate complete bilingual text object
   */
  bilingualText: (text: { en: string; ar: string }) => {
    try {
      bilingualTextSchema.parse(text)
      return { isValid: true, error: null }
    } catch (error) {
      if (error instanceof z.ZodError) {
        return {
          isValid: false,
          error: error.issues[0]?.message || 'Invalid bilingual text',
        }
      }
      return { isValid: false, error: 'Validation failed' }
    }
  },

  /**
   * Validate bilingual component
   */
  bilingualComponent: (component: any) => {
    try {
      bilingualComponentSchema.parse(component)
      return { isValid: true, error: null }
    } catch (error) {
      if (error instanceof z.ZodError) {
        return {
          isValid: false,
          error: error.issues[0]?.message || 'Invalid bilingual component',
        }
      }
      return { isValid: false, error: 'Validation failed' }
    }
  },
}

export default {
  createThemeSchema,
  updateThemeSchema,
  colorsSchema,
  typographySchema,
  identitySchema,
  bilingualTextSchema,
  bilingualComponentSchema,
  validateThemeCreate,
  validateThemeUpdate,
  fieldValidators,
  bilingualFieldValidators,
}
