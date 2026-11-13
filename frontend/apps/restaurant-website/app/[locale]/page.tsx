import type { Metadata } from 'next'
import { headers } from 'next/headers'
import { AlertCircle } from 'lucide-react'
import { DynamicHomePage } from '@/components/homepage/dynamic-home-page'
import { extractThemeConfig } from '@/lib/api/get-theme-server'
import { getServerApiBaseUrl } from '@/lib/config/api'
import { normalizeThemeConfig } from '@/lib/utils/theme-normalizer'
import type { Product } from '@/components/commerce/product-card'
import { TEMPLATE_REGISTRY } from '@/lib/utils/theme-templates'

interface PageProps {
  params: Promise<{
    locale: 'en' | 'ar'
  }>
  searchParams?: Promise<{
    template?: string
  }>
}

interface CategoryLike {
  id?: number | string
  name?: string
  name_en?: string
  name_ar?: string
  item_count?: number
  count?: number
  emoji?: string
}

export const revalidate = 300

async function fetchJson(path: string, revalidateSeconds: number) {
  const response = await fetch(path, { next: { revalidate: revalidateSeconds } })
  if (!response.ok) {
    throw new Error(`Failed request (${response.status}) for ${path}`)
  }
  return response.json()
}

function parseProducts(payload: unknown): Product[] {
  const data = (payload as Record<string, unknown>)?.data as Record<string, unknown> | undefined
  const productsFromData = data?.products
  if (Array.isArray(productsFromData)) {
    return productsFromData as Product[]
  }
  const productsFromRoot = (payload as Record<string, unknown>)?.products
  return Array.isArray(productsFromRoot) ? (productsFromRoot as Product[]) : []
}

function parseCategories(payload: unknown): CategoryLike[] {
  const data = (payload as Record<string, unknown>)?.data as Record<string, unknown> | undefined
  const categoriesFromData = data?.categories
  if (Array.isArray(categoriesFromData)) {
    return categoriesFromData as CategoryLike[]
  }
  const categoriesFromRoot = (payload as Record<string, unknown>)?.categories
  return Array.isArray(categoriesFromRoot) ? (categoriesFromRoot as CategoryLike[]) : []
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { locale } = await params
  const headersList = await headers()
  const restaurantSlug = headersList.get('x-restaurant-slug') || ''

  if (!restaurantSlug) {
    return {
      title: 'Restaurant',
      description: 'Order food online',
    }
  }

  try {
    const apiUrl = getServerApiBaseUrl()
    const themePayload = await fetchJson(
      `${apiUrl}/public/restaurants/${restaurantSlug}/theme`,
      3600
    )
    const themeConfig = normalizeThemeConfig(extractThemeConfig(themePayload) ?? {}, locale)
    const siteTitle = String((themeConfig.identity as Record<string, unknown>)?.site_title ?? themeConfig.name ?? 'Restaurant')
    const description = locale === 'ar'
      ? `اطلب اونلاين من ${siteTitle}`
      : `Order online from ${siteTitle}`

    return {
      title: `${siteTitle} | ${locale === 'ar' ? 'اطلب اونلاين' : 'Order Online'}`,
      description,
      alternates: {
        languages: {
          en: `/${'en'}`,
          ar: `/${'ar'}`,
        },
      },
    }
  } catch {
    return {
      title: 'Restaurant',
      description: 'Order food online',
    }
  }
}

export default async function HomePage({ params, searchParams }: PageProps) {
  const { locale } = await params
  const headersList = await headers()
  const restaurantSlug = headersList.get('x-restaurant-slug') || ''

  if (!restaurantSlug) {
    return (
      <div className="flex min-h-screen items-center justify-center px-4">
        <div className="max-w-md rounded-lg border border-red-300 bg-red-50 p-6 text-center">
          <AlertCircle className="mx-auto mb-3 h-10 w-10 text-red-600" />
          <h1 className="mb-2 text-xl font-semibold text-red-900">
            {locale === 'ar' ? 'المطعم غير موجود' : 'Restaurant Not Found'}
          </h1>
          <p className="text-sm text-red-800">
            {locale === 'ar'
              ? 'يرجى الدخول عبر رابط المطعم الصحيح.'
              : 'Please access this page through a valid restaurant subdomain.'}
          </p>
        </div>
      </div>
    )
  }

  try {
    const resolvedSearchParams = (await searchParams) || {}
    const requestedTemplate = typeof resolvedSearchParams.template === 'string' ? resolvedSearchParams.template : ''
    const templateOverride = TEMPLATE_REGISTRY[requestedTemplate]
    const apiUrl = getServerApiBaseUrl()
    const [themePayload, homepagePayload, productsPayload, categoriesPayload] = await Promise.all([
      fetchJson(`${apiUrl}/public/restaurants/${restaurantSlug}/theme`, 3600),
      fetchJson(`${apiUrl}/public/restaurants/${restaurantSlug}/homepage`, 300),
      fetchJson(`${apiUrl}/public/restaurants/${restaurantSlug}/products?lang=${locale}`, 300),
      fetchJson(`${apiUrl}/public/restaurants/${restaurantSlug}/categories?lang=${locale}`, 300),
    ])

    const themeSource = extractThemeConfig(themePayload) || extractThemeConfig(homepagePayload) || {}
    const themeConfig = normalizeThemeConfig(themeSource, locale)
    if (templateOverride) {
      themeConfig.templateId = templateOverride.id
      themeConfig.headerVariant = templateOverride.headerVariant
      themeConfig.footerVariant = templateOverride.footerVariant
      ;(themeConfig as unknown as Record<string, unknown>).heroVariant = templateOverride.heroVariant
      ;(themeConfig as unknown as Record<string, unknown>).productsVariant = templateOverride.productsVariant
    }
    const products = parseProducts(productsPayload)
    const categories = parseCategories(categoriesPayload)

    return (
      <DynamicHomePage
        locale={locale}
        restaurantSlug={restaurantSlug}
        themeConfig={themeConfig}
        products={products}
        categories={categories}
      />
    )
  } catch {
    return (
      <div className="flex min-h-screen items-center justify-center px-4">
        <div className="max-w-md rounded-lg border border-red-300 bg-red-50 p-6 text-center">
          <AlertCircle className="mx-auto mb-3 h-10 w-10 text-red-600" />
          <h1 className="mb-2 text-xl font-semibold text-red-900">
            {locale === 'ar' ? 'تعذر تحميل الصفحة' : 'Unable to Load Homepage'}
          </h1>
          <p className="text-sm text-red-800">
            {locale === 'ar'
              ? 'حدث خطأ أثناء تحميل بيانات المطعم.'
              : 'There was an error loading restaurant data.'}
          </p>
        </div>
      </div>
    )
  }
}
