'use client'

import { useState, useEffect } from 'react'
import { useRouter, usePathname, useParams } from 'next/navigation'
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

interface Product {
  id: number
  name: string
  [key: string]: any
}

interface Category {
  id: number
  name: string
  [key: string]: any
}

export default function EditProductPage() {
  const router = useRouter()
  const pathname = usePathname()
  const params = useParams()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  // Check RBAC permissions for Products module
  const canUpdateProducts = useHasPermission(
    MODULE_IDS.PRODUCTS,
    PERMISSION_LEVELS.WRITE,
  )

  const productId = params.id as string
  const [product, setProduct] = useState<Product | null>(null)
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const token = useAuthStore((state) => state.token)

  // Fetch product and categories with caching
  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true)

        // Check cache for categories
        let cachedCategories = apiCache.get('categories')
        if (!cachedCategories) {
          const categoriesResponse = await api.get('/categories')
          cachedCategories = categoriesResponse.data || []
          // Cache for 5 minutes
          apiCache.set('categories', cachedCategories, 300000)
        }
        setCategories(cachedCategories)

        // Fetch product (always fresh, don't cache single products)
        const productResponse = await api.get(`/products/${productId}`)
        setProduct(productResponse.data)
      } catch (error) {
        toast.error('Failed to load data')
        debugError('Failed to load data:', error)
        router.push(`/${locale}/dashboard/products`)
      } finally {
        setLoading(false)
      }
    }

    if (token && productId) {
      fetchData()
    }
  }, [token, productId, locale, router])

  const handleUpdateProduct = async (formData: any) => {
    try {
      const response = await api.put(`/products/${productId}`, formData)

      // Invalidate products cache when updating
      apiCache.clear('products')
      apiCache.clear('products_list')

      toast.success(t('products.updateSuccess'))
      router.push(`/${locale}/dashboard/products`)
    } catch (error) {
      toast.error(t('products.updateFailed'))
      debugError(error)
    }
  }

  const handleCancel = () => {
    router.push(`/${locale}/dashboard/products`)
  }

  const breadcrumbItems = [
    { href: `/${locale}/dashboard/products`, label: t('products.title') },
    { href: `/${locale}/dashboard/products/${productId}`, label: product?.name_en || 'Product' },
    { label: t('products.editProduct') }
  ]

  // Show permission denied if user can't update products
  if (!canUpdateProducts) {
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
                Access Denied
              </h2>
              <p className="text-sm text-red-800 dark:text-red-400 mt-1">
                You do not have permission to edit products. Please contact your administrator.
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
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{t('products.editProduct')}</h1>
        <p className="text-gray-500 dark:text-gray-400 mt-2">
          {product?.name_en && `Editing: ${product.name_en}`}
        </p>
      </div>

      {/* Form Container */}
      <div className="bg-white dark:bg-gray-900 rounded-lg shadow-md p-6">
        {loading && !product ? (
          <FormSkeleton />
        ) : product ? (
          <ProductFormPage
            product={product}
            categories={categories}
            onSubmit={handleUpdateProduct}
            onCancel={handleCancel}
          />
        ) : (
          <div className="text-center py-12">
            <p className="text-red-500">Product not found</p>
          </div>
        )}
      </div>
    </div>
  )
}
