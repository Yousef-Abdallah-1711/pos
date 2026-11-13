'use client'

import { useState } from 'react'
import Image from 'next/image'
import { usePathname } from 'next/navigation'
import { useForm, Controller } from 'react-hook-form'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { toast } from 'sonner'
import { debugLog, debugError } from '@/lib/debug'
import { errorText } from '@/lib/form-errors'
import type { ProductFormProps, ProductFormValues } from './product-types'

export function ProductFormPage({ product, categories, onSubmit, onCancel }: ProductFormProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'
  const [imageFile, setImageFile] = useState<File | null>(null)
  const [imagePreview, setImagePreview] = useState(product?.main_image_url || '')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const toNumber = (value: unknown, fallback = 0) => {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : fallback
  }
  const toOptionalNumber = (value: unknown) => {
    if (value === '' || value === null || value === undefined) {
      return undefined
    }
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : undefined
  }
  const toInt = (value: unknown, fallback = 0) => {
    const parsed = Number.parseInt(String(value), 10)
    return Number.isNaN(parsed) ? fallback : parsed
  }
  const toOptionalInt = (value: unknown) => {
    if (value === '' || value === null || value === undefined) {
      return undefined
    }
    const parsed = Number.parseInt(String(value), 10)
    return Number.isNaN(parsed) ? undefined : parsed
  }

  const { control, handleSubmit, formState: { errors }, watch } = useForm<ProductFormValues>({
    defaultValues: {
      name_en: product?.name_en || '',
      name_ar: product?.name_ar || '',
      category_id: product?.category_id || '',
      description_en: product?.description_en || '',
      description_ar: product?.description_ar || '',
      price: product?.price || '',
      cost: product?.cost || '',
      discount_price: product?.discount_price || '',
      quantity_in_stock: product?.quantity_in_stock || 0,
      low_stock_threshold: product?.low_stock_threshold || 10,
      is_vegetarian: product?.is_vegetarian || false,
      is_vegan: product?.is_vegan || false,
      is_spicy: product?.is_spicy || false,
      is_gluten_free: product?.is_gluten_free || false,
      is_available: product?.is_available !== false,
      featured: product?.featured || false,
      track_inventory: product?.track_inventory !== false,
    },
  })

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setImageFile(file)
      const reader = new FileReader()
      reader.onloadend = () => {
        setImagePreview(typeof reader.result === 'string' ? reader.result : '')
      }
      reader.readAsDataURL(file)
    }
  }

  const onFormSubmit = async (data: ProductFormValues) => {
    setIsSubmitting(true)
    try {
      const formData = new FormData()

      // Sanitize data to match Go struct types
      const cleanData = {
        ...data,
        // Convert string inputs to proper numbers
        category_id: data.category_id ? toInt(data.category_id, 0) : 0,
        price: toNumber(data.price, 0),
        cost: toOptionalNumber(data.cost),
        discount_price: toOptionalNumber(data.discount_price),
        discount_percentage: toOptionalNumber(data.discount_percentage),
        quantity_in_stock: product && toInt(data.quantity_in_stock, 0) === product.quantity_in_stock
          ? undefined : toInt(data.quantity_in_stock, 0),
        low_stock_threshold: toInt(data.low_stock_threshold, 0),
        display_order: toInt(data.display_order, 0),
        calories: toOptionalInt(data.calories),
        // Ensure empty strings are handled as undefined or appropriate empty values for optional fields
        description_en: data.description_en || "",
        description_ar: data.description_ar || "",
        name_ar: data.name_ar || "",
      }

      debugLog('DEBUG: Form Submission Raw:', data);
      debugLog('DEBUG: Form Submission Cleaned:', cleanData);

      // Add all data as JSON string in 'data' field (backend expects this)
      formData.append('data', JSON.stringify(cleanData))

      // Add image if selected
      if (imageFile) {
        formData.append('image', imageFile)
      }

      await onSubmit(formData)
    } catch (error) {
      toast.error(t('products.saveFailed', 'Failed to save product'))
      debugError(error)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-8">
      {/* General Information Section */}
      <section className="space-y-6">
        <div className="pb-4 border-b">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">{t('products.generalInformation', 'General Information')}</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">{t('products.productBasicDetails', 'Basic product details and naming')}</p>
        </div>

        {/* Image Upload */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.productImage', 'Product Image')}</label>
          <Input
            type="file"
            accept="image/*"
            onChange={handleImageChange}
            className="w-full"
          />
          {imagePreview && (
            <div className="relative w-40 h-40 mt-2">
              <Image
                src={imagePreview}
                unoptimized
                alt={t('products.imagePreviewAlt', 'Preview')}
                fill
                className="object-cover rounded"
              />
            </div>
          )}
        </div>

        {/* Names - English and Arabic */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.nameEnglishRequired', 'Name (English) *')}</label>
            <Controller
              name="name_en"
              control={control}
              rules={{ required: t('products.nameEnglishRequiredError', 'English name is required') }}
              render={({ field }) => <Input {...field} placeholder={t('products.nameEnglishPlaceholder', 'e.g., Grilled Chicken')} />}
            />
            {errors.name_en && <span className="text-red-500 text-sm">{errorText(errors.name_en)}</span>}
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.nameArabicLabel', 'Name (Arabic)')}</label>
            <Controller
              name="name_ar"
              control={control}
              render={({ field }) => <Input {...field} placeholder={t('products.productNameArPlaceholder')} dir="rtl" />}
            />
            {errors.name_ar && <span className="text-red-500 text-sm">{errorText(errors.name_ar)}</span>}
          </div>
        </div>

        {/* Category */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.categoryLabel', 'Category')}</label>
          <Controller
            name="category_id"
            control={control}
            render={({ field }) => (
              <Select value={field.value.toString()} onValueChange={field.onChange}>
                <SelectTrigger>
                  <SelectValue placeholder={t('products.selectCategory', 'Select category')} />
                </SelectTrigger>
                <SelectContent>
                  {categories.map((cat) => (
                    <SelectItem key={cat.id} value={cat.id.toString()}>
                      {cat.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          />
          {errors.category_id && (
            <span className="text-red-500 text-sm">{errorText(errors.category_id)}</span>
          )}
        </div>

        {/* Descriptions - English and Arabic */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.descriptionEnglishLabel', 'Description (English)')}</label>
            <Controller
              name="description_en"
              control={control}
              render={({ field }) => (
                <Textarea {...field} placeholder={t('products.descriptionEnglishPlaceholder', 'Describe the product...')} rows={4} />
              )}
            />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.descriptionArabicLabel', 'Description (Arabic)')}</label>
            <Controller
              name="description_ar"
              control={control}
              render={({ field }) => (
                <Textarea {...field} placeholder={t('products.descriptionArPlaceholder')} rows={4} dir="rtl" />
              )}
            />
          </div>
        </div>

        {/* Availability Toggles */}
        <div className="space-y-4">
          <div className={`flex items-center space-x-3 ${isRTL ? 'space-x-reverse' : ''} p-3 bg-gray-50 dark:bg-gray-800 rounded`}>
            <Controller
              name="is_available"
              control={control}
              render={({ field }) => (
                <Checkbox checked={field.value} onCheckedChange={field.onChange} />
              )}
            />
            <div>
              <label className="text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.availableLabel', 'Available')}</label>
              <p className="text-xs text-gray-500 dark:text-gray-400">{t('products.availableHelp', 'Mark this product as available for purchase')}</p>
            </div>
          </div>

          <div className={`flex items-center space-x-3 ${isRTL ? 'space-x-reverse' : ''} p-3 bg-gray-50 dark:bg-gray-800 rounded`}>
            <Controller
              name="featured"
              control={control}
              render={({ field }) => (
                <Checkbox checked={field.value} onCheckedChange={field.onChange} />
              )}
            />
            <div>
              <label className="text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.featuredLabel', 'Featured Product')}</label>
              <p className="text-xs text-gray-500 dark:text-gray-400">{t('products.featuredHelp', 'Display prominently in menu')}</p>
            </div>
          </div>
        </div>
      </section>

      {/* Pricing Section */}
      <section className="space-y-6">
        <div className="pb-4 border-b">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">{t('products.pricingSectionTitle', 'Pricing')}</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">{t('products.pricingSectionHelp', 'Set prices and discounts')}</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.priceEgpRequired', 'Price (EGP) *')}</label>
            <Controller
              name="price"
              control={control}
              rules={{ required: t('products.priceRequired', 'Price is required'), min: { value: 0.01, message: t('products.priceMinError', 'Price must be greater than 0') } }}
              render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />}
            />
            {errors.price && <span className="text-red-500 text-sm">{errorText(errors.price)}</span>}
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.costEgpLabel', 'Cost (EGP)')}</label>
            <Controller
              name="cost"
              control={control}
              render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />}
            />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.discountPriceEgpLabel', 'Discount Price (EGP)')}</label>
            <Controller
              name="discount_price"
              control={control}
              render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />}
            />
          </div>
        </div>
      </section>

      {/* Inventory Section */}
      <section className="space-y-6">
        <div className="pb-4 border-b">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">{t('products.inventorySectionTitle', 'Inventory')}</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">{t('products.inventorySectionHelp', 'Manage stock and tracking')}</p>
        </div>

        <div className={`flex items-center space-x-3 ${isRTL ? 'space-x-reverse' : ''} p-3 bg-gray-50 dark:bg-gray-800 rounded mb-4`}>
          <Controller
            name="track_inventory"
            control={control}
            render={({ field }) => (
              <Checkbox checked={field.value} onCheckedChange={field.onChange} />
            )}
          />
          <div>
            <label className="text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.trackInventoryLabel', 'Track Inventory')}</label>
            <p className="text-xs text-gray-500 dark:text-gray-400">{t('products.trackInventoryHelp', 'Monitor stock levels for this product')}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.quantityInStockLabel', 'Quantity in Stock')}</label>
            <Controller
              name="quantity_in_stock"
              control={control}
              render={({ field }) => <Input {...field} type="number" min="0" />}
            />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.lowStockThresholdLabel', 'Low Stock Threshold')}</label>
            <Controller
              name="low_stock_threshold"
              control={control}
              render={({ field }) => <Input {...field} type="number" min="0" />}
            />
          </div>
        </div>
      </section>

      {/* Classification Section */}
      <section className="space-y-6">
        <div className="pb-4 border-b">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">{t('products.classificationSectionTitle', 'Classification')}</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">{t('products.classificationSectionHelp', 'Dietary and allergen information')}</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-3">
            <div className={`flex items-center space-x-3 ${isRTL ? 'space-x-reverse' : ''} p-3 bg-gray-50 dark:bg-gray-800 rounded`}>
              <Controller
                name="is_vegetarian"
                control={control}
                render={({ field }) => (
                  <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                )}
              />
              <label className="text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.vegetarian', 'Vegetarian')}</label>
            </div>

            <div className={`flex items-center space-x-3 ${isRTL ? 'space-x-reverse' : ''} p-3 bg-gray-50 dark:bg-gray-800 rounded`}>
              <Controller
                name="is_vegan"
                control={control}
                render={({ field }) => (
                  <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                )}
              />
              <label className="text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.vegan', 'Vegan')}</label>
            </div>
          </div>

          <div className="space-y-3">
            <div className={`flex items-center space-x-3 ${isRTL ? 'space-x-reverse' : ''} p-3 bg-gray-50 dark:bg-gray-800 rounded`}>
              <Controller
                name="is_spicy"
                control={control}
                render={({ field }) => (
                  <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                )}
              />
              <label className="text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.spicy', 'Spicy')}</label>
            </div>

            <div className={`flex items-center space-x-3 ${isRTL ? 'space-x-reverse' : ''} p-3 bg-gray-50 dark:bg-gray-800 rounded`}>
              <Controller
                name="is_gluten_free"
                control={control}
                render={({ field }) => (
                  <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                )}
              />
              <label className="text-sm font-medium text-gray-700 dark:text-gray-300">{t('products.glutenFree', 'Gluten Free')}</label>
            </div>
          </div>
        </div>
      </section>

      {/* Form Actions */}
      <div className="flex justify-end gap-3 pt-6 border-t">
        <Button type="button" variant="outline" onClick={onCancel}>
          {t('common.cancel', 'Cancel')}
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting
            ? t('common.saving', 'Saving...')
            : product
              ? t('products.updateProductButton', 'Update Product')
              : t('products.createProductButton', 'Create Product')}
        </Button>
      </div>
    </form>
  )
}
