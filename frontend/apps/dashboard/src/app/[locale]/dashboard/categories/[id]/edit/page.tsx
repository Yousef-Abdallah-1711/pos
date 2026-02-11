'use client'

import { useState, useEffect } from 'react'
import { useRouter, usePathname, useParams } from 'next/navigation'
import { api } from '@/lib/api'
import { apiCache } from '@/lib/apiCache'
import { toast } from 'sonner'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { CategoryFormPage } from '@/components/categories/CategoryFormPage'
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbSeparator, BreadcrumbPage } from '@/components/ui/breadcrumb'
import { FormSkeleton } from '@/components/skeletons'
import { debugError } from '@/lib/debug'

interface Category {
  id: number
  name: string
  [key: string]: any
}

export default function EditCategoryPage() {
  const router = useRouter()
  const pathname = usePathname()
  const params = useParams()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  const categoryId = params.id as string
  const [category, setCategory] = useState<Category | null>(null)
  const [loading, setLoading] = useState(true)

  // Fetch category
  useEffect(() => {
    const fetchCategory = async () => {
      try {
        setLoading(true)
        const response = await api.get(`/categories/${categoryId}`)
        setCategory(response.data)
      } catch (error) {
        toast.error('Failed to load category')
        debugError('Failed to load category:', error)
        router.push(`/${locale}/dashboard/categories`)
      } finally {
        setLoading(false)
      }
    }

    if (categoryId) {
      fetchCategory()
    }
  }, [categoryId, locale, router])

  const handleUpdateCategory = async (data: any) => {
    try {
      const response = await api.put(`/categories/${categoryId}`, data)

      // Invalidate categories cache when updating
      apiCache.clear('categories')

      toast.success(t('categories.updateSuccess') || 'Category updated successfully')
      router.push(`/${locale}/dashboard/categories`)
    } catch (error) {
      toast.error(t('categories.updateFailed') || 'Failed to update category')
      debugError(error)
    }
  }

  const handleCancel = () => {
    router.push(`/${locale}/dashboard/categories`)
  }

  const breadcrumbItems = [
    { href: `/${locale}/dashboard/categories`, label: t('categories.title') },
    { href: `/${locale}/dashboard/categories/${categoryId}`, label: category?.name || 'Category' },
    { label: t('categories.editCategory') }
  ]

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
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{t('categories.editCategory')}</h1>
        <p className="text-gray-500 dark:text-gray-400 mt-2">
          {category?.name && `Editing: ${category.name}`}
        </p>
      </div>

      {/* Form Container */}
      <div className="bg-white dark:bg-gray-900 rounded-lg shadow-md p-6">
        {loading && !category ? (
          <FormSkeleton />
        ) : category ? (
          <CategoryFormPage
            category={category}
            onSubmit={handleUpdateCategory}
            onCancel={handleCancel}
          />
        ) : (
          <div className="text-center py-12">
            <p className="text-red-500">Category not found</p>
          </div>
        )}
      </div>
    </div>
  )
}
