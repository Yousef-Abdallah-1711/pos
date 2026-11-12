'use client'

import { usePathname } from 'next/navigation'
import { useForm, Controller } from 'react-hook-form'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Checkbox } from '@/components/ui/checkbox'
import { useState } from 'react'
import { errorText } from '@/lib/form-errors'
import type { CategoryFormPageProps, CategoryFormValues } from './category-types'

export function CategoryFormPage({ category, onSubmit, onCancel }: CategoryFormPageProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const { control, handleSubmit, reset, formState: { errors, isValid } } = useForm<CategoryFormValues>({
    mode: 'onChange',
    defaultValues: {
      name: category?.name || '',
      name_ar: category?.name_ar || '',
      description: category?.description || '',
      description_ar: category?.description_ar || '',
      display_order: category?.display_order || 0,
      is_active: category?.is_active !== false,
    },
  })

  const onFormSubmit = async (data: CategoryFormValues) => {
    setIsSubmitting(true)
    try {
      await onSubmit(data)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-8">
      {/* General Information Section */}
      <section className="space-y-6">
        <div className="pb-4 border-b">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">{t('categories.generalInformation')}</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">{t('categories.categoryNameAndBasicDetails')}</p>
        </div>

        {/* Names - English and Arabic */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('categories.nameEnglish')}</label>
            <Controller
              name="name"
              control={control}
              rules={{ required: t('categories.nameRequired') }}
              render={({ field }) => (
                <Input
                  {...field}
                  placeholder={t('categories.nameEnglishPlaceholder')}
                />
              )}
            />
            {errors.name && <span className="text-red-500 text-sm">{errorText(errors.name)}</span>}
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('categories.nameArabicLabel')}</label>
            <Controller
              name="name_ar"
              control={control}
              render={({ field }) => (
                <Input
                  {...field}
                  placeholder={t('categories.nameArabicPlaceholder')}
                  dir="rtl"
                />
              )}
            />
          </div>
        </div>

        {/* Descriptions - English and Arabic */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('categories.descriptionEnglish')}</label>
            <Controller
              name="description"
              control={control}
              render={({ field }) => (
                <Textarea
                  {...field}
                  placeholder={t('categories.descriptionEnglishPlaceholder')}
                  rows={4}
                />
              )}
            />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('categories.descriptionArabicLabel')}</label>
            <Controller
              name="description_ar"
              control={control}
              render={({ field }) => (
                <Textarea
                  {...field}
                  placeholder={t('categories.descriptionArabicPlaceholder')}
                  rows={4}
                  dir="rtl"
                />
              )}
            />
          </div>
        </div>
      </section>

      {/* Settings Section */}
      <section className="space-y-6">
        <div className="pb-4 border-b">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">{t('categories.settings')}</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">{t('categories.displayOrderAndVisibility')}</p>
        </div>

        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('categories.displayOrderLabel')}</label>
          <Controller
            name="display_order"
            control={control}
            render={({ field }) => (
              <Input
                {...field}
                type="number"
                min="0"
                placeholder="0"
              />
            )}
          />
          <p className="text-xs text-gray-500 dark:text-gray-400">{t('categories.lowerNumbersAppearFirst')}</p>
        </div>

        <div className="flex items-center space-x-3 p-3 bg-gray-50 dark:bg-gray-800 rounded">
          <Controller
            name="is_active"
            control={control}
            render={({ field }) => (
              <Checkbox checked={field.value} onCheckedChange={field.onChange} />
            )}
          />
          <div>
            <label className="text-sm font-medium text-gray-700 dark:text-gray-300">{t('categories.activeLabel')}</label>
            <p className="text-xs text-gray-500 dark:text-gray-400">{t('categories.categoryVisibleAndAvailable')}</p>
          </div>
        </div>
      </section>

      {/* Form Actions */}
      <div className="flex justify-end gap-3 pt-6 border-t">
        <Button type="button" variant="outline" onClick={onCancel}>
          {t('categories.cancel')}
        </Button>
        <Button type="submit" disabled={!isValid || isSubmitting}>
          {isSubmitting ? t('categories.savingButton') : category ? t('categories.updateCategory') : t('categories.createCategory')}
        </Button>
      </div>
    </form>
  )
}
