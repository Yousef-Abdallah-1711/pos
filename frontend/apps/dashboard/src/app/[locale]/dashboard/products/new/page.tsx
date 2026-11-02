'use client'

import { useState, useEffect } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { useHasPermission } from '@/hooks/useRbac'
import { api } from '@/lib/api'
import { apiCache } from '@/lib/apiCache'
import { toast } from 'sonner'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { ProductFormPage } from '@/components/products/ProductFormPage'
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbSeparator, BreadcrumbPage } from '@/components/ui/breadcrumb'
import { FormSkeleton } from '@/components/skeletons'
import { Lock } from 'lucide-react'
import { debugError } from '@/lib/debug'
import type { ProductCategory as Category } from '@/components/products/product-types'

export default function CreateProductPage() {
  const router = useRouter()
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  // Check RBAC permissions for Products module
  const canCreateProductsPermission = useHasPermission(
    MODULE_IDS.PRODUCTS,
    PERMISSION_LEVELS.WRITE,
  )
  const canCreateProducts = canCreateProductsPermission.loading || canCreateProductsPermission.hasPermission

  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const token = useAuthStore((state) => state.token)

  // Fetch categories with caching
  useEffect(() => {
    const fetchCategories = async () => {
      try {
        // Check cache first (categories don't change often)
        const cachedCategories = apiCache.get('categories')
        if (cachedCategories) {
          setCategories(cachedCategories)
          setLoading(false)
          return
        }

        const response = await api.get('/categories')
        const categoryData = response.data || []

        // Cache for 5 minutes (300000ms)
        apiCache.set('categories', categoryData, 300000)
        setCategories(categoryData)
      } catch (error) {
        toast.error(t('categories.loadFailed', 'Failed to load categories'))
        debugError('Failed to load categories:', error)
        setCategories([])
      } finally {
        setLoading(false)
      }
    }

    if (token) {
      fetchCategories()
    }
  }, [token])

  const handleCreateProduct = async (formData: any) => {
    try {
      const response = await api.post('/products', formData)

      // Invalidate products cache when creating new product
      apiCache.clear('products')
      apiCache.clear('products_list')

      toast.success(t('products.createSuccess'))
      router.push(`/${locale}/dashboard/products`)
    } catch (error) {
      toast.error(t('products.createFailed'))
      debugError(error)
    }
  }

  const handleCancel = () => {
    router.push(`/${locale}/dashboard/products`)
  }

  const breadcrumbItems = [
    { href: `/${locale}/dashboard/products`, label: t('products.title') },
    { label: t('products.createNewProduct') }
  ]

  // Show permission denied if user can't create products
  if (!canCreateProducts) {
    return (
      <div className="space-y-6">
        <Breadcrumb>
          {breadcrumbItems.map((item, index) => (
            <div key={index} className="flex items-center gap-2">
              {index > 0 && <BreadcrumbSeparator />}
              {item.href ? (
                <BreadcrumbLink href={item.href} className="text-blue-600 hover:text-blue-700">
                  {item.label}
                </BreadcrumbLink>
              ) : (
                <BreadcrumbPage className="text-gray-600">{item.label}</BreadcrumbPage>
              )}
            </div>
          ))}
        </Breadcrumb>

        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
          <div className="flex gap-3 items-start">
            <Lock className="h-5 w-5 text-red-600 dark:text-red-400 mt-0.5 flex-shrink-0" />
            <div>
              <h2 className="font-semibold text-red-900 dark:text-red-300">
                {t('products.accessDeniedTitle', 'Access Denied')}
              </h2>
              <p className="text-sm text-red-800 dark:text-red-400 mt-1">
                {t('products.accessDeniedCreateText', 'You do not have permission to create products. Please contact your administrator.')}
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Breadcrumb Navigation */}
      <Breadcrumb>
        {breadcrumbItems.map((item, index) => (
          <div key={index} className="flex items-center gap-2">
            {index > 0 && <BreadcrumbSeparator />}
            {item.href ? (
              <BreadcrumbLink href={item.href} className="text-blue-600 hover:text-blue-700">
                {item.label}
              </BreadcrumbLink>
            ) : (
              <BreadcrumbPage className="text-gray-600">{item.label}</BreadcrumbPage>
            )}
          </div>
        ))}
      </Breadcrumb>

      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{t('products.createNewProduct')}</h1>
        <p className="text-gray-500 dark:text-gray-400 mt-2">
          {t('products.addProductDescription', 'Add a new product to your menu')}
        </p>
      </div>

      {/* Form Container */}
      <div className="bg-white dark:bg-gray-900 rounded-lg shadow-md p-6">
        {loading && categories.length === 0 ? (
          <FormSkeleton />
        ) : (
          <ProductFormPage
            product={null}
            categories={categories}
            onSubmit={handleCreateProduct}
            onCancel={handleCancel}
          />
        )}
      </div>
    </div>
  )
}
