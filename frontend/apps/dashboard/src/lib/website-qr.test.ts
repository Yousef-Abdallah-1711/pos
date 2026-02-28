import { describe, expect, it } from 'vitest'
import {
  buildWebsiteQrTargetUrl,
  getWebsiteBaseUrl,
} from './website-qr'

describe('website-qr', () => {
  describe('getWebsiteBaseUrl', () => {
    it('keeps absolute https domains and trims trailing slash', () => {
      expect(getWebsiteBaseUrl('https://demo.example.com/')).toBe(
        'https://demo.example.com'
      )
    })

    it('uses http for localhost domains', () => {
      expect(getWebsiteBaseUrl('localhost:3003')).toBe('http://localhost:3003')
    })

    it('defaults to https for non-localhost domains', () => {
      expect(getWebsiteBaseUrl('tenant.example.com')).toBe(
        'https://tenant.example.com'
      )
    })
  })

  describe('buildWebsiteQrTargetUrl', () => {
    it('builds locale-aware home URL', () => {
      expect(buildWebsiteQrTargetUrl('home', 'en', 'tenant.example.com')).toBe(
        'https://tenant.example.com/en'
      )
    })

    it('builds locale-aware menu URL', () => {
      expect(buildWebsiteQrTargetUrl('menu', 'ar', 'localhost:3003')).toBe(
        'http://localhost:3003/ar/menu'
      )
    })

    it('handles locale with slashes safely', () => {
      expect(
        buildWebsiteQrTargetUrl('menu', '/en/', 'https://tenant.example.com/')
      ).toBe('https://tenant.example.com/en/menu')
    })
  })
})
