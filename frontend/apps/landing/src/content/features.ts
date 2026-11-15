import type { Locale } from '@pos-saas/i18n'
import english from '../i18n/messages/en/features.json'
import arabic from '../i18n/messages/ar/features.json'
import captures from './captures.json'

export interface FeatureStep {
  title: string
  description: string
}

export interface Feature {
  slug: string
  title: string
  eyebrow: string
  headline: string
  description: string
  image: keyof typeof captures
  steps: FeatureStep[]
  boundary: string
  faqs: Array<{ question: string; answer: string }>
  related: string[]
}

export function getFeatures(locale: Locale): Feature[] {
  const localized = locale === 'ar' ? arabic : english
  return localized.map(feature => ({
    ...feature,
    image: feature.image as Feature['image'],
  }))
}
