import { registerEditorLoader } from '@pos-saas/ui-components/registry'

export const registerDashboardThemeEditors = () => {
  registerEditorLoader('HeroEditor', () => import('./components/HeroEditor'))
  registerEditorLoader('ProductsEditor', () => import('./components/ProductsEditor'))
  registerEditorLoader('WhyChooseUsEditor', () => import('./components/WhyChooseUsEditor'))
  registerEditorLoader('ContactEditor', () => import('./components/ContactEditor'))
  registerEditorLoader('TestimonialsEditor', () => import('./components/TestimonialsEditor'))
  registerEditorLoader('CTAEditor', () => import('./components/CTAEditor'))
}
