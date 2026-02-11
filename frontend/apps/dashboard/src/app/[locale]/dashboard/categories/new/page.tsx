'use client'

import { useRouter, usePathname } from 'next/navigation'
import { api } from '@/lib/api'
import { apiCache } from '@/lib/apiCache'
import { toast } from 'sonner'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { CategoryFormPage } from '@/components/categories/CategoryFormPage'
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbSeparator, BreadcrumbPage } from '@/components/ui/breadcrumb'
import { debugError } from '@/lib/debug'

export default function CreateCategoryPage() {
  const router = useRouter()
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  const handleCreateCategory = async (data: any) => {
    try {
      const response = await api.post('/categories', data)

      // Invalidate categories cache when creating new category
      apiCache.clear('categories')

      toast.success(t('categories.createSuccess') || 'Category created successfully')
      router.push(`/${locale}/dashboard/categories`)
    } catch (error) {
      toast.error(t('categories.createFailed') || 'Failed to create category')
      debugError(error)
    }
  }

  const handleCancel = () => {
    router.push(`/${locale}/dashboard/categories`)
  }

  const breadcrumbItems = [
    { href: `/${locale}/dashboard/categories`, label: t('categories.title') },
    { label: t('categories.createCategory') }
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
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{t('categories.createCategory')}</h1>
        <p className="text-gray-500 dark:text-gray-400 mt-2">
          Add a new category to organize your products
        </p>
      </div>

      {/* Form Container */}
      <div className="bg-white dark:bg-gray-900 rounded-lg shadow-md p-6">
        <CategoryFormPage
          category={null}
          onSubmit={handleCreateCategory}
          onCancel={handleCancel}
        />
      </div>
    </div>
  )
}
