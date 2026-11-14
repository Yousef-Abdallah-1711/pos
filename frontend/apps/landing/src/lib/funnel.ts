export type FunnelEvent =
  | 'landing_primary_cta_clicked' | 'feature_page_viewed' | 'pricing_viewed'
  | 'billing_period_selected' | 'trial_registration_started'
  | 'trial_registration_completed' | 'onboarding_started' | 'onboarding_completed'
  | 'demo_contact_initiated'

interface FunnelContext {
  locale: 'en' | 'ar'
  source: string
  billing?: 'monthly' | 'annual'
}

// Local event bridge only. No personal fields, identifiers or external tracking.
// An approved analytics adapter can subscribe without changing product controls.
export function trackFunnel(event: FunnelEvent, context: FunnelContext) {
  if (typeof window === 'undefined') return
  window.dispatchEvent(new CustomEvent('servio:funnel', {
    detail: { event, locale: context.locale, source: context.source, billing: context.billing },
  }))
}
