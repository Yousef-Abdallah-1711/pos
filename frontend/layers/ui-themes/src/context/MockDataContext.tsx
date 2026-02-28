import React, { createContext, ReactNode, useMemo } from 'react'

// Mock data imports
import { HERO_MOCK_DATA } from '../sections/marketing/HeroSection/mockData'
import { PRODUCTS_MOCK_DATA } from '../sections/commerce/ProductsSection/mockData'
import { WHY_CHOOSE_US_MOCK_DATA } from '../sections/content/WhyChooseUsSection/mockData'
import { TESTIMONIALS_MOCK_DATA } from '../sections/marketing/TestimonialsSection/mockData'
import { CTA_MOCK_DATA } from '../sections/marketing/CTASection/mockData'
import { CONTACT_MOCK_DATA } from '../sections/forms/ContactSection/mockData'
import { CUSTOM_MOCK_DATA } from '../sections/custom/CustomSection/mockData'

/**
 * Mock data registry mapping component IDs to their mock data
 */
export const MOCK_DATA_REGISTRY: Record<string, Record<string, any>> = {
  hero: HERO_MOCK_DATA,
  products: PRODUCTS_MOCK_DATA,
  why_us: WHY_CHOOSE_US_MOCK_DATA,
  testimonials: TESTIMONIALS_MOCK_DATA,
  cta: CTA_MOCK_DATA,
  contact: CONTACT_MOCK_DATA,
  custom: CUSTOM_MOCK_DATA,
}

/**
 * Mock data context interface
 */
export interface MockDataContextValue {
  // Check if mock data is available for a component
  hasMockData: (componentId: string) => boolean

  // Get mock data for a component
  getMockData: (componentId: string) => Record<string, any> | null

  // Get all available mock data
  getAllMockData: () => Record<string, Record<string, any>>

  // Get mock data with fallback
  getMockDataWithFallback: (componentId: string, fallback?: Record<string, any>) => Record<string, any>

  // Check if development mode is enabled
  isDevelopmentMode: boolean

  // Get statistics about mock data
  getStats: () => MockDataStats
}

/**
 * Mock data statistics
 */
export interface MockDataStats {
  totalComponents: number
  componentsWithMockData: number
  componentIds: string[]
  developmentMode: boolean
}

/**
 * Create the MockDataContext
 */
export const MockDataContext = createContext<MockDataContextValue | undefined>(undefined)

/**
 * Props for MockDataProvider
 */
export interface MockDataProviderProps {
  children: ReactNode
  developmentMode?: boolean
  customMockData?: Record<string, Record<string, any>>
}

/**
 * Mock Data Provider Component
 *
 * Provides mock data to all child components via context
 *
 * @example
 * ```tsx
 * <MockDataProvider developmentMode={true}>
 *   <App />
 * </MockDataProvider>
 * ```
 */
export const MockDataProvider: React.FC<MockDataProviderProps> = ({
  children,
  developmentMode = process.env.NODE_ENV === 'development',
  customMockData = {},
}) => {
  const value = useMemo<MockDataContextValue>(() => {
    // Merge custom mock data with default registry
    const mergedRegistry = {
      ...MOCK_DATA_REGISTRY,
      ...customMockData,
    }

    return {
      hasMockData: (componentId: string) => {
        return componentId in mergedRegistry
      },

      getMockData: (componentId: string) => {
        const normalizedId = componentId.toLowerCase().replace(/[_-]/g, '_')
        return mergedRegistry[normalizedId] || null
      },

      getAllMockData: () => {
        return mergedRegistry
      },

      getMockDataWithFallback: (componentId: string, fallback = {}) => {
        const data = mergedRegistry[componentId]
        return data || fallback
      },

      isDevelopmentMode: developmentMode,

      getStats: () => ({
        totalComponents: Object.keys(mergedRegistry).length,
        componentsWithMockData: Object.keys(mergedRegistry).length,
        componentIds: Object.keys(mergedRegistry),
        developmentMode,
      }),
    }
  }, [customMockData, developmentMode])

  return <MockDataContext.Provider value={value}>{children}</MockDataContext.Provider>
}

export default MockDataProvider
