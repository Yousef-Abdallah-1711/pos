'use client'

import { memo } from 'react'
import type { POSCategory } from '@/lib/api/posApi'

interface POSCategoryTabsProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  categories: POSCategory[]
  selectedCategoryId: number | null
  onSelectCategory: (categoryId: number | null) => void
}

function POSCategoryTabsComponent({
  locale,
  t,
  categories,
  selectedCategoryId,
  onSelectCategory,
}: POSCategoryTabsProps) {
  const isRTL = locale === 'ar'
  return (
    <nav className="overflow-x-auto pb-1">
      <div className={`flex min-w-max items-center gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
        <button
          type="button"
          onClick={() => onSelectCategory(null)}
          className={`rounded-full border px-3 py-1.5 text-xs font-semibold transition ${
            selectedCategoryId === null
              ? 'border-primary bg-primary text-white'
              : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200 dark:hover:bg-gray-800'
          }`}
        >
          {t('pos.allCategories', 'All')}
        </button>

        {categories.map((category) => {
          const isActive = selectedCategoryId === category.id
          const label =
            locale === 'ar' && category.name_ar ? category.name_ar : category.name

          return (
            <button
              key={category.id}
              type="button"
              onClick={() => onSelectCategory(isActive ? null : category.id)}
              className={`rounded-full border px-3 py-1.5 text-xs font-semibold transition ${
                isActive
                  ? 'border-primary bg-primary text-white'
                  : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200 dark:hover:bg-gray-800'
              }`}
            >
              {label}
            </button>
          )
        })}
      </div>
    </nav>
  )
}

export const POSCategoryTabs = memo(POSCategoryTabsComponent)
POSCategoryTabs.displayName = 'POSCategoryTabs'

export default POSCategoryTabs
