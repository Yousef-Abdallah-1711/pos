import { describe, expect, it } from 'vitest'
import { NextRequest } from 'next/server'
import { middleware } from './middleware'

describe('dashboard document locale', () => {
  it.each([['/ar/dashboard', 'ar'], ['/en/dashboard', 'en'], ['/auth/callback', 'en']])(
    'derives document metadata from %s', (pathname, locale) => {
      const request = new NextRequest(`https://pos.app.rakmyat.online${pathname}`, {
        headers: { 'x-dashboard-locale': locale === 'ar' ? 'en' : 'ar' },
      })
      expect(middleware(request).headers.get('x-middleware-request-x-dashboard-locale')).toBe(locale)
    }
  )
})
