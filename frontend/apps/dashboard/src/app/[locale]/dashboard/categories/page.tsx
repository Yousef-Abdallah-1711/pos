'use client'

import { useState, useEffect } from 'react'
import { usePathname } from 'next/navigation'
import Link from 'next/link'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-react'
import { toast } from 'sonner'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { PageHeaderSkeleton, TableSkeleton } from '@/components/skeletons'
import { debugError } from '@/lib/debug'

interface Category {
  id: number
  name: string
  name_ar?: string
  description?: string
  description_ar?: string
  icon_url?: string
  display_order: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export default function CategoriesPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadCategories()
  }, [])

  const loadCategories = async () => {
    try {
      setLoading(true)
      const response = await api.get('/categories')
      setCategories(response.data || [])
    } catch (error: any) {
      const errorMsg = error.response?.data?.error || t('categories.loadFailed', 'Failed to load categories')
      toast.error(errorMsg)
      debugError('Error loading categories:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id: number) => {
    const confirmMsg = t('categories.confirmDelete', 'Are you sure you want to delete this category?')

    if (!confirm(confirmMsg)) return

    try {
      await api.delete(`/categories/${id}`)
      const successMsg = t('categories.deleteSuccess', 'Category deleted successfully')
      toast.success(successMsg)
      loadCategories()
    } catch (error: any) {
      const errorMsg = error.response?.data?.error || t('categories.deleteFailed', 'Failed to delete category')
      toast.error(errorMsg)
      debugError('Error deleting category:', error)
    }
  }

  return (
    <div className="space-y-6">
      {loading && categories.length === 0 ? (
        <>
          <PageHeaderSkeleton />
          <TableSkeleton rows={8} />
        </>
      ) : (
        <>
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <div>
              <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">{t('categories.title')}</h1>
              <p className="text-gray-500 dark:text-gray-400 mt-1">
                {t('categories.manageCategories')}
              </p>
            </div>
            <Link href={`/${locale}/dashboard/categories/new`}>
              <Button className="gap-2 w-full sm:w-auto">
                <Plus className="h-4 w-4" />
                {t('categories.addCategory')}
              </Button>
            </Link>
          </div>

          <div className="bg-white dark:bg-gray-800 rounded-lg shadow">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t('categories.name')}</TableHead>
                  <TableHead>{t('categories.nameArabic')}</TableHead>
                  <TableHead>{t('categories.description')}</TableHead>
                  <TableHead>{t('categories.displayOrder')}</TableHead>
                  <TableHead>{t('categories.status')}</TableHead>
                  <TableHead>{t('categories.actions')}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loading ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center py-8">
                      {t('categories.loading')}
                    </TableCell>
                  </TableRow>
                ) : categories.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center py-8 text-gray-500 dark:text-gray-400">
                      {t('categories.noCategories')}
                    </TableCell>
                  </TableRow>
                ) : (
                  categories.map((category) => (
                    <TableRow key={category.id}>
                      <TableCell className="font-medium">{category.name}</TableCell>
                      <TableCell>{category.name_ar || '-'}</TableCell>
                      <TableCell className="max-w-xs truncate">{category.description || '-'}</TableCell>
                      <TableCell>{category.display_order}</TableCell>
                      <TableCell>
                        <span
                          className={`px-2 py-1 rounded text-sm ${
                            category.is_active
                              ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
                              : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
                          }`}
                        >
                          {category.is_active ? t('categories.active') : t('categories.inactive')}
                        </span>
                      </TableCell>
                      <TableCell>
                        <div className="flex gap-2">
                          <Link href={`/${locale}/dashboard/categories/${category.id}/edit`}>
                            <Button
                              size="sm"
                              variant="outline"
                            >
                              {t('categories.edit')}
                            </Button>
                          </Link>
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={() => handleDelete(category.id)}
                          >
                            {t('categories.delete')}
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </>
      )}
    </div>
  )
}
