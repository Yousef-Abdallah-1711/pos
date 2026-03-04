'use client'

import { useState, useEffect } from 'react'
import { usePathname } from 'next/navigation'
import Link from 'next/link'
import { useAuthStore } from '@/stores/authStore'
import { useUserPermissions } from '@/hooks/useRbac'
import { ProductList } from '@/components/products/ProductList'
import { Button } from '@/components/ui/button'
import { api } from '@/lib/api'
import { toast } from 'sonner'
import { Plus, Lock } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { PageHeaderSkeleton, TableSkeleton } from '@/components/skeletons'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { debugError } from '@/lib/debug'

interface Product {
  id: number
  name: string
  description?: string
  price: number
  category_id?: number
  [key: string]: any
}

interface Category {
  id: number
  name: string
  [key: string]: any
}

interface Filters {
  page: number
  pageSize: number
  categoryId: number | null
  searchText: string
  sortBy: string
  sortOrder: string
}

export default function ProductsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(false)
  const [filters, setFilters] = useState<Filters>({
    page: 1,
    pageSize: 20,
    categoryId: null,
    searchText: '',
    sortBy: 'created_at',
    sortOrder: 'DESC',
  })

  const token = useAuthStore((state) => state.token)

  // Check RBAC permissions for Products module
  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canViewProducts = checkPermission(MODULE_IDS.PRODUCTS, PERMISSION_LEVELS.READ)
  const canCreateProducts = checkPermission(MODULE_IDS.PRODUCTS, PERMISSION_LEVELS.WRITE)
  const canDeleteProducts = checkPermission(MODULE_IDS.PRODUCTS, PERMISSION_LEVELS.DELETE)

  // Fetch products
  useEffect(() => {
    const fetchProducts = async () => {
      if (permissionLoading || !canViewProducts) {
        setLoading(false)
        return
      }
      setLoading(true)
      try {
        const params: any = {
          page: filters.page,
          page_size: filters.pageSize,
          sort_by: filters.sortBy,
          sort_order: filters.sortOrder,
        }

        if (filters.categoryId) params.category_id = filters.categoryId
        if (filters.searchText) params.search = filters.searchText

        const response = await api.get('/products', { params })

        setProducts(response.data.products || [])
      } catch (error) {
        toast.error(t('products.loadFailed', 'Failed to load products'))
        debugError(error)
      } finally {
        setLoading(false)
      }
    }

    fetchProducts()
  }, [filters, token, permissionLoading, canViewProducts])

  // Fetch categories
  useEffect(() => {
    const fetchCategories = async () => {
      if (permissionLoading || !canViewProducts) {
        return
      }
      try {
        const response = await api.get('/categories')
        setCategories(response.data || [])
      } catch (error) {
        debugError('Failed to load categories:', error)
      }
    }

    fetchCategories()
  }, [token, permissionLoading, canViewProducts])


  const handleDeleteProduct = async (productId: number) => {
    // Check delete permission
    if (!canDeleteProducts) {
      toast.error('You do not have permission to delete products')
      return
    }

    if (!window.confirm(t('products.confirmDelete'))) {
      return
    }

    try {
      await api.delete(`/products/${productId}`)

      setProducts(products.filter((p) => p.id !== productId))
      toast.success(t('products.deleteSuccess'))
    } catch (error) {
      toast.error(t('products.deleteFailed'))
      debugError(error)
    }
  }


  if (permissionLoading) {
    return (
      <div className="space-y-6">
        <PageHeaderSkeleton />
        <TableSkeleton rows={8} />
      </div>
    )
  }

  // Show permission denied if user can't view products
  if (!canViewProducts) {
    return (
      <div className="space-y-6">
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
          <div className="flex gap-3 items-start">
            <Lock className="h-5 w-5 text-red-600 dark:text-red-400 mt-0.5 flex-shrink-0" />
            <div>
              <h2 className="font-semibold text-red-900 dark:text-red-300">
                {t('products.accessDeniedTitle', 'Access Denied')}
              </h2>
              <p className="text-sm text-red-800 dark:text-red-400 mt-1">
                {t('products.accessDeniedText', 'You do not have permission to view products. Please contact your administrator.')}
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {loading && products.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <>
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <div>
              <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">
                {t('products.title')}
              </h1>
              <p className="text-gray-500 dark:text-gray-400 mt-1">
                {t('products.manageProducts')}
              </p>
            </div>
            {canCreateProducts ? (
              <Link href={`/${locale}/dashboard/products/new`}>
                <Button className="gap-2 w-full sm:w-auto">
                  <Plus className="h-4 w-4" />
                  {t('products.addProduct')}
                </Button>
              </Link>
            ) : (
              <Button disabled className="gap-2 w-full sm:w-auto">
                <Plus className="h-4 w-4" />
                {t('products.addProduct')}
              </Button>
            )}
          </div>

          <ProductList
            products={products}
            categories={categories}
            loading={loading}
            filters={filters}
            onFilterChange={setFilters}
            onDelete={canDeleteProducts ? handleDeleteProduct : undefined}
            canDelete={canDeleteProducts}
          />
        </>
      )}
    </div>
  )
}
