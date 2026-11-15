import type { Locale, SiteCopy } from '@/content/site'

export type NavigationCopyKey = keyof SiteCopy['navigation']

export const SITE_ROUTES = {
  home: '',
  features: 'features',
  howItWorks: 'how-it-works',
  pricing: 'pricing',
  productTour: 'demo',
  demoOptions: 'contact',
  privacy: 'privacy',
  login: 'auth/login',
  register: 'auth/register',
} as const

export type SiteRouteId = keyof typeof SITE_ROUTES

export function localizedHref(locale: Locale, route: SiteRouteId | string): string {
  return route ? `/${locale}/${route}` : `/${locale}`
}

export const HEADER_LINKS = [
  { id: 'platform', route: 'features', label: 'platform', activePrefix: 'features' },
  { id: 'how-it-works', route: 'howItWorks', label: 'howItWorks' },
  { id: 'pricing', route: 'pricing', label: 'pricing' },
  { id: 'product-tour', route: 'productTour', label: 'productTour' },
  { id: 'demo-options', route: 'demoOptions', label: 'demoOptions' },
] as const satisfies ReadonlyArray<{
  id: string
  route: SiteRouteId
  label: NavigationCopyKey
  activePrefix?: string
}>

const FEATURE_NAVIGATION = {
  pos: { label: 'pos', description: 'posDescription', icon: 'pos' },
  'kitchen-display': { label: 'kitchen', description: 'kitchenDescription', icon: 'kitchen' },
  'online-ordering': { label: 'onlineOrdering', description: 'onlineOrderingDescription', icon: 'orders' },
  inventory: { label: 'inventory', description: 'inventoryDescription', icon: 'inventory' },
  hr: { label: 'staff', description: 'staffDescription', icon: 'staff' },
  analytics: { label: 'operations', description: 'operationsDescription', icon: 'operations' },
  'website-builder': { label: 'storefront', description: 'storefrontDescription', icon: 'storefront' },
} as const satisfies Record<string, { label: NavigationCopyKey; description: NavigationCopyKey; icon: string }>

type FeatureNavigationId = keyof typeof FEATURE_NAVIGATION

function featureNavigationItem(id: FeatureNavigationId) {
  return { id, ...FEATURE_NAVIGATION[id] }
}

export const FEATURE_GROUPS = [
  {
    id: 'orders',
    label: 'ordersGroup',
    items: [featureNavigationItem('pos'), featureNavigationItem('kitchen-display'), featureNavigationItem('online-ordering')],
  },
  {
    id: 'management',
    label: 'managementGroup',
    items: [featureNavigationItem('inventory'), featureNavigationItem('hr'), featureNavigationItem('analytics')],
  },
  {
    id: 'presence',
    label: 'presenceGroup',
    items: [featureNavigationItem('website-builder')],
  },
] as const satisfies ReadonlyArray<{
  id: string
  label: NavigationCopyKey
  items: ReadonlyArray<{ id: FeatureNavigationId; label: NavigationCopyKey; description: NavigationCopyKey; icon: string }>
}>

export interface FooterNavigationGroup {
  id: string
  label: NavigationCopyKey
  links: ReadonlyArray<{ id: string; route: SiteRouteId | string; label: NavigationCopyKey }>
}

export function getFooterGroups(): FooterNavigationGroup[] {
  const featureLinks = FEATURE_GROUPS.flatMap(group => group.items.map(item => ({
    id: item.id,
    route: `features/${item.id}`,
    label: item.label,
  })))
  return [
    {
      id: 'product',
      label: 'productGroup',
      links: [
        { id: 'all-features', route: 'features', label: 'allFeatures' },
        ...featureLinks,
      ],
    },
    {
      id: 'explore',
      label: 'exploreGroup',
      links: [
        { id: 'how-it-works', route: 'howItWorks', label: 'howItWorks' },
        { id: 'product-tour', route: 'productTour', label: 'productTour' },
        { id: 'demo-options', route: 'demoOptions', label: 'demoOptions' },
        { id: 'pricing', route: 'pricing', label: 'pricing' },
      ],
    },
    {
      id: 'account',
      label: 'accountGroup',
      links: [
        { id: 'login', route: 'login', label: 'login' },
        { id: 'register', route: 'register', label: 'startTrial' },
      ],
    },
    {
      id: 'company',
      label: 'companyGroup',
      links: [
        { id: 'privacy', route: 'privacy', label: 'privacy' },
        { id: 'contact', route: 'demoOptions', label: 'contact' },
      ],
    },
  ]
}

export function isNavigationRouteActive(pathname: string, route: string, activePrefix?: string): boolean {
  const localePath = pathname.replace(/^\/(?:en|ar)(?=\/|$)/, '') || '/'
  const target = `/${route}`
  if (localePath === target) return true
  return activePrefix ? localePath === `/${activePrefix}` || localePath.startsWith(`/${activePrefix}/`) : false
}
