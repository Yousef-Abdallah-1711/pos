import { describe, expect, it } from 'vitest'
import { defaultLocale, directionFor, formatCurrency, formatNumber, selectPlural, translate } from '@pos-saas/i18n'
import englishSubscription from '@pos-saas/i18n/locales/en/subscription.json'
import arabicSubscription from '@pos-saas/i18n/locales/ar/subscription.json'

describe('shared locale utilities', () => {
  it('uses Arabic and RTL as the public visitor defaults', () => {
    expect(defaultLocale).toBe('ar')
    expect(directionFor('ar')).toBe('rtl')
    expect(directionFor('en')).toBe('ltr')
  })

  it('formats amounts and plural categories using the active locale', () => {
    expect(formatNumber('en', 1250)).toContain('1')
    expect(formatCurrency('ar', 750)).toContain('ج.م')
    expect(selectPlural('en', { one: '{count} order', other: '{count} orders' }, 2)).toBe('2 orders')
  })

  it('loads JSON values and fails visibly when a key is missing', () => {
    expect(translate(englishSubscription, 'blocked.trialExpired.title')).toBe('Your trial has expired')
    expect(translate(arabicSubscription, 'blocked.trialExpired.title')).toBe('انتهت الفترة التجريبية')
    expect(() => translate(englishSubscription, 'blocked.trialExpired.missing' as never)).toThrow(/Missing translation/)
  })

  it('safely interpolates declared values without evaluating message text', () => {
    const dictionary = { greeting: 'Hello {name}' } as const
    expect(translate(dictionary, 'greeting', { name: '<b>Guest</b>' })).toBe('Hello <b>Guest</b>')
  })
})
