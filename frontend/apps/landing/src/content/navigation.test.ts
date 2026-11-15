import { describe, expect, it } from 'vitest'
import { getFeatures } from '@/content/features'
import { getSiteCopy } from '@/content/site'
import { FEATURE_GROUPS, getFooterGroups, isNavigationRouteActive, localizedHref, SITE_ROUTES } from './navigation'

describe('public marketing navigation', () => {
  it('links every catalog feature exactly once in the mega menu and footer', () => {
    const configured = FEATURE_GROUPS.flatMap(group => group.items.map(item => item.id))
    const footerFeatures = getFooterGroups().find(group => group.id === 'product')?.links
      .filter(item => item.route.startsWith('features/'))
      .map(item => item.route.slice('features/'.length))

    expect(configured).toHaveLength(7)
    expect(new Set(configured).size).toBe(configured.length)
    expect(configured.sort()).toEqual(getFeatures('en').map(feature => feature.slug).sort())
    expect(configured.sort()).toEqual(getFeatures('ar').map(feature => feature.slug).sort())
    expect(footerFeatures?.sort()).toEqual(configured.sort())
  })

  it('uses only real public route destinations and localized navigation copy', () => {
    const routes = new Set(Object.values(SITE_ROUTES).map(route => localizedHref('en', route)))
    const footerRoutes = getFooterGroups().flatMap(group => group.links.map(link => localizedHref('en', link.route)))
    const featureRoutes = FEATURE_GROUPS.flatMap(group => group.items.map(item => localizedHref('en', `features/${item.id}`)))

    expect(footerRoutes.filter(route => !routes.has(route) && !featureRoutes.includes(route))).toEqual([])
    expect(getSiteCopy('en').navigation.platform).toBe('Platform')
    expect(getSiteCopy('ar').navigation.platform).toBe('المنصة')
    expect(Object.keys(getSiteCopy('en').navigation).sort()).toEqual(Object.keys(getSiteCopy('ar').navigation).sort())
  })

  it('marks feature detail routes as part of the platform section in both locales', () => {
    expect(isNavigationRouteActive('/en/features/pos', 'features', 'features')).toBe(true)
    expect(isNavigationRouteActive('/ar/features/analytics', 'features', 'features')).toBe(true)
    expect(isNavigationRouteActive('/ar/pricing', 'features', 'features')).toBe(false)
  })
})
