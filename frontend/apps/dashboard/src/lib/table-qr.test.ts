import { describe, expect, it } from 'vitest'
import { buildTableMenuUrl, getRestaurantBaseUrl, sanitizeQrFileName } from './table-qr'

describe('table-qr', () => {
  describe('getRestaurantBaseUrl', () => {
    it('keeps absolute https domain and trims trailing slash', () => {
      expect(getRestaurantBaseUrl('https://restaurant.example.com/')).toBe(
        'https://restaurant.example.com'
      )
    })

    it('uses http for localhost domains', () => {
      expect(getRestaurantBaseUrl('localhost:3003')).toBe('http://localhost:3003')
    })

    it('defaults to https for bare domains', () => {
      expect(getRestaurantBaseUrl('demo.example.com')).toBe('https://demo.example.com')
    })
  })

  describe('buildTableMenuUrl', () => {
    it('builds locale-aware table menu URL', () => {
      expect(buildTableMenuUrl(12, 'en', 'demo.example.com')).toBe(
        'https://demo.example.com/en/table/12/menu'
      )
      expect(buildTableMenuUrl(5, 'ar', 'localhost:3003')).toBe(
        'http://localhost:3003/ar/table/5/menu'
      )
    })
  })

  describe('sanitizeQrFileName', () => {
    it('replaces unsafe characters and trims separators', () => {
      expect(sanitizeQrFileName(' table #12 / qr ')).toBe('table-12-qr')
    })
  })
})
