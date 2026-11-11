import { parsePhoneNumberFromString } from 'libphonenumber-js'
import { z } from 'zod'

const codePointLength = (value: string) => Array.from(value).length

export function normalizeInternationalPhone(value: string): string | undefined {
  const trimmed = value.trim()
  if (!trimmed) return undefined
  // Match the backend's explicit international contract: only ASCII digits,
  // one leading plus, and common visual separators are accepted.
  if (!/^\+[0-9\s().-]+$/.test(trimmed)) return undefined
  if ((trimmed.match(/[0-9]/g)?.length ?? 0) > 15) return undefined
  const parsed = parsePhoneNumberFromString(trimmed)
  return parsed?.isValid() && parsed.number.replace(/\D/g, '').length <= 15 ? parsed.number : undefined
}

export function isSubdomainFormatValid(value: string): boolean {
  const slug = value.trim().toLowerCase()
  return slug.length >= 2 && slug.length <= 63
    && /^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/.test(slug)
    && !slug.includes('--')
}

const RESERVED_SUBDOMAINS = new Set([
  'api', 'admin', 'dashboard', 'www', 'mail', 'ftp', 'cdn', 'auth', 'app', 'docs', 'status', 'help', 'support',
  'blog', 'shop', 'store', 'dev', 'staging', 'test', 'localhost', '127', 'git', 'gitlab', 'github', 'jenkins', 'docker', 'kubernetes',
])

const nameField = (required: string, min: string, max: string) => z.string().trim()
  .min(1, required)
  .refine((value) => codePointLength(value) >= 2, min)
  .refine((value) => codePointLength(value) <= 255, max)

const emailField = (required: string, invalid: string, max: string) => z.string().trim().toLowerCase()
  .min(1, required)
  .email(invalid)
  .refine((value) => codePointLength(value) <= 255, max)

export function createRegisterSchema(messages: {
  required: string
  nameMin: string
  nameMax: string
  invalidEmail: string
  emailMax: string
  passwordMin: string
  passwordMax: string
  passwordMismatch: string
  invalidSubdomain: string
  reservedSubdomain: string
  invalidPhone: string
}) {
  return z.object({
    name: nameField(messages.required, messages.nameMin, messages.nameMax),
    email: emailField(messages.required, messages.invalidEmail, messages.emailMax),
    restaurantName: nameField(messages.required, messages.nameMin, messages.nameMax),
    restaurantSlug: z.string().trim().toLowerCase()
      .refine(isSubdomainFormatValid, messages.invalidSubdomain)
      .refine((value) => !RESERVED_SUBDOMAINS.has(value), messages.reservedSubdomain),
    phone: z.string().trim().transform((value, context) => {
      const normalized = normalizeInternationalPhone(value)
      if (value && !normalized) context.addIssue({ code: z.ZodIssueCode.custom, message: messages.invalidPhone })
      return normalized ?? ''
    }),
    password: z.string().refine((value) => codePointLength(value) >= 6, messages.passwordMin).refine((value) => new TextEncoder().encode(value).length <= 72, messages.passwordMax),
    confirmPassword: z.string().min(1, messages.required),
  }).refine((value) => value.password === value.confirmPassword, {
    message: messages.passwordMismatch,
    path: ['confirmPassword'],
  })
}

export function createLoginSchema(messages: {
  required: string
  invalidEmail: string
  emailMax: string
  passwordMin: string
  passwordMax: string
}) {
  return z.object({
    email: emailField(messages.required, messages.invalidEmail, messages.emailMax),
    password: z.string().min(1, messages.required)
      .refine((value) => new TextEncoder().encode(value).length <= 72, messages.passwordMax),
  })
}
