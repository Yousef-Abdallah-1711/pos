import { describe, expect, it } from 'vitest'
import { createLoginSchema, createRegisterSchema, normalizeInternationalPhone } from '@pos-saas/form-foundation'
import { applyApiFieldErrors, parseApiError } from '@pos-saas/form-foundation'

const messages = {
  required: 'Required', nameMin: 'Name too short', nameMax: 'Name too long', invalidEmail: 'Invalid email', emailMax: 'Email too long',
  passwordMin: 'Password too short', passwordMax: 'Password too long', passwordMismatch: 'Passwords do not match',
  invalidSubdomain: 'Invalid subdomain', reservedSubdomain: 'Reserved subdomain', invalidPhone: 'Invalid phone',
}

describe('shared registration contract', () => {
  const schema = createRegisterSchema(messages)
  const valid = { name: 'يوسف عبد الله', email: ' OWNER@Example.COM ', password: '123456', confirmPassword: '123456', restaurantName: 'مطعم النيل', restaurantSlug: 'nile-food', phone: '+1 (650) 555-0123' }

  it('trims and canonicalizes international and Unicode inputs', () => {
    const result = schema.parse(valid)
    expect(result.name).toBe('يوسف عبد الله')
    expect(result.email).toBe('owner@example.com')
    expect(result.phone).toBe('+16505550123')
  })

  it('treats blank phone as optional and rejects local or invalid numbers', () => {
    expect(normalizeInternationalPhone('  ')).toBeUndefined()
    expect(schema.safeParse({ ...valid, phone: '' }).success).toBe(true)
    expect(schema.safeParse({ ...valid, phone: '6505550123' }).success).toBe(false)
    expect(schema.safeParse({ ...valid, phone: '+999123' }).success).toBe(false)
    expect(schema.safeParse({ ...valid, phone: '+1234567890123456' }).success).toBe(false)
  })

  it('validates Unicode code point boundaries and subdomain format', () => {
    expect(schema.safeParse({ ...valid, name: '名' }).success).toBe(false)
    expect(schema.safeParse({ ...valid, name: '名'.repeat(256) }).success).toBe(false)
    expect(schema.safeParse({ ...valid, restaurantSlug: 'bad--slug' }).success).toBe(false)
    expect(schema.safeParse({ ...valid, restaurantSlug: 'admin' }).success).toBe(false)
    expect(schema.safeParse({ ...valid, restaurantSlug: 'A-good-2' }).success).toBe(true)
  })

  it('attaches password mismatch to confirmation and enforces bcrypt byte limit', () => {
    const mismatch = schema.safeParse({ ...valid, confirmPassword: '654321' })
    expect(mismatch.success).toBe(false)
    if (!mismatch.success) expect(mismatch.error.issues.some((issue) => issue.path.join('.') === 'confirmPassword')).toBe(true)
    expect(schema.safeParse({ ...valid, password: 'é'.repeat(37), confirmPassword: 'é'.repeat(37) }).success).toBe(false)
  })
})

describe('shared login and API error handling', () => {
  it('normalizes email and validates required login values', () => {
    const schema = createLoginSchema({ required: 'Required', invalidEmail: 'Invalid email', emailMax: 'Too long', passwordMin: 'Too short', passwordMax: 'Too long' })
    expect(schema.parse({ email: ' USER@EXAMPLE.COM ', password: 'validpw' }).email).toBe('user@example.com')
    expect(schema.safeParse({ email: 'bad', password: '' }).success).toBe(false)
    expect(schema.safeParse({ email: 'user@example.com', password: 'old' }).success).toBe(true)
  })

  it('maps documented error codes without exposing unknown internal messages', () => {
    const fallback = { validation: 'Review fields', invalidCredentials: 'Sign-in failed', rateLimited: 'Wait', unavailable: 'Try later', network: 'Offline', unexpected: 'Unexpected', conflict: 'Conflict' }
    expect(parseApiError({ response: { status: 500, data: { code: 'INTERNAL_ERROR', error: 'pq: private schema' } } }, fallback).message).toBe('Unexpected')
    expect(parseApiError({ response: { status: 409, data: { code: 'SUBDOMAIN_TAKEN', fields: { restaurant_slug: 'Taken' } } } }, fallback).fields.restaurant_slug).toBe('Taken')
    expect(parseApiError({ code: 'ERR_NETWORK' }, fallback).retryable).toBe(true)
  })

  it('maps API field aliases to form names', () => {
    const mapped: Record<string, string> = {}
    applyApiFieldErrors({ restaurant_slug: 'Taken', phone: 'Invalid' }, (field, message) => { mapped[field] = message })
    expect(mapped).toEqual({ restaurantSlug: 'Taken', phone: 'Invalid' })
  })
})
