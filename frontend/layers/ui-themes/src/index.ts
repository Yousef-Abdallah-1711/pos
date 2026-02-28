/**
 * UI Themes - Shared theme and component library
 *
 * Exports all theme types, components, and sections for use in both
 * dashboard and restaurant website applications.
 */

// Types
export type {
  GlobalColors,
  TypographySettings,
  WebsiteIdentity,
  ThemeComponent,
  Theme,
  ThemeProviderProps,
  SectionRendererProps,
  HeaderConfig,
  FooterConfig,
  ThemeData,
  BaseSectionProps,
  HeroSectionProps,
  ProductsSectionProps,
  WhyChooseUsSectionProps,
  ContactSectionProps,
  TestimonialsSectionProps,
  CTASectionProps,
  CustomSectionProps,
  HomepageData,
} from './types'

// Components
export { ThemeProvider, SectionRenderer, Header, Footer } from './components'

// Sections
export {
  HeroSection,
  ProductsSection,
  WhyChooseUsSection,
  ContactSection,
  TestimonialsSection,
  CTASection,
  CustomSection,
} from './sections'

// Component Registry System
export { componentRegistry, resolveDynamicComponent, ComponentLoadingFallback, ComponentErrorFallback } from './registry/componentRegistry'
export type { ComponentDefinition } from './registry/componentRegistry'

// Component Cache System
export { ComponentCache } from './registry/componentCache'
export type { CacheStats, CacheEntry } from './registry/componentCache'

// Component Discovery System
export { componentDiscovery } from './registry/componentDiscovery'
export type { DiscoveredComponent } from './registry/componentDiscovery'

// Hooks
export { useComponentRegistry } from './hooks/useComponentRegistry'
export { useComponentDiscovery } from './hooks/useComponentDiscovery'
export {
  useMockData,
  useHasMockData,
  useAllMockData,
  useIsDevelopmentMode,
  useMockDataStats,
  useMockDataContext,
} from './hooks/useMockData'
export { withMockData } from './hooks/withMockData'
export type { WithMockDataProps } from './hooks/withMockData'

// Mock Data Context
export { MockDataProvider, MockDataContext } from './context/MockDataContext'
export type { MockDataContextValue, MockDataProviderProps, MockDataStats } from './context/MockDataContext'
export { MOCK_DATA_REGISTRY } from './context/MockDataContext'
