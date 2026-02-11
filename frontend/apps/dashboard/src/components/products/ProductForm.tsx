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
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { toast } from 'sonner'
import { debugError } from '@/lib/debug'

export function ProductForm({ product, categories, onSubmit, onCancel }) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const [imageFile, setImageFile] = useState(null)
  const [imagePreview, setImagePreview] = useState(product?.main_image_url || '')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const toNumber = (value, fallback = 0) => {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : fallback
  }
  const toOptionalNumber = (value) => {
    if (value === '' || value === null || value === undefined) {
      return undefined
    }
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : undefined
  }
  const toInt = (value, fallback = 0) => {
    const parsed = Number.parseInt(String(value), 10)
    return Number.isNaN(parsed) ? fallback : parsed
  }
  const toOptionalInt = (value) => {
    if (value === '' || value === null || value === undefined) {
      return undefined
    }
    const parsed = Number.parseInt(String(value), 10)
    return Number.isNaN(parsed) ? undefined : parsed
  }

  const { control, handleSubmit, formState: { errors }, watch } = useForm({
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

  const handleImageChange = (e) => {
    const file = e.target.files[0]
    if (file) {
      setImageFile(file)
      const reader = new FileReader()
      reader.onloadend = () => {
        setImagePreview(reader.result)
      }
      reader.readAsDataURL(file)
    }
  }

  const onFormSubmit = async (data) => {
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
        quantity_in_stock: toInt(data.quantity_in_stock, 0),
        low_stock_threshold: toInt(data.low_stock_threshold, 0),
        display_order: toInt(data.display_order, 0),
        calories: toOptionalInt(data.calories),
        // Ensure empty strings are handled as undefined or appropriate empty values for optional fields
        description_en: data.description_en || "",
        description_ar: data.description_ar || "",
        name_ar: data.name_ar || "",
      }



      // Add all data as JSON string in 'data' field (backend expects this)
      formData.append('data', JSON.stringify(cleanData))

      // Add image if selected
      if (imageFile) {
        formData.append('image', imageFile)
      }

      await onSubmit(formData)
    } catch (error) {
      toast.error('Failed to save product')
      debugError(error)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-8">
      <Tabs defaultValue="general" className="w-full">
        <TabsList className="grid w-full grid-cols-4 bg-gray-100 dark:bg-gray-800 p-1">
          <TabsTrigger value="general">General</TabsTrigger>
          <TabsTrigger value="pricing">Pricing</TabsTrigger>
          <TabsTrigger value="inventory">Inventory</TabsTrigger>
          <TabsTrigger value="details">Details</TabsTrigger>
        </TabsList>

        {/* General Tab */}
        <TabsContent value="general" className="space-y-4">
          {/* Image Upload */}
          <div className="space-y-2">
            <label className="block text-sm font-medium">Product Image</label>
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
                  alt="Preview"
                  fill
                  className="object-cover rounded"
                />
              </div>
            )}
          </div>

          {/* Names */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Name (English) *</label>
              <Controller
                name="name_en"
                control={control}
                rules={{ required: 'English name is required' }}
                render={({ field }) => <Input {...field} placeholder="e.g., Grilled Chicken" />}
              />
              {errors.name_en && <span className="text-red-500 text-sm">{errors.name_en.message}</span>}
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Name (Arabic)</label>
              <Controller
                name="name_ar"
                control={control}
                render={({ field }) => <Input {...field} placeholder={t('products.productNameArPlaceholder')} dir="rtl" />}
              />
              {errors.name_ar && <span className="text-red-500 text-sm">{errors.name_ar.message}</span>}
            </div>
          </div>

          {/* Category */}
          <div className="space-y-2">
            <label className="block text-sm font-medium">Category</label>
            <Controller
              name="category_id"
              control={control}
              render={({ field }) => (
                <Select value={field.value.toString()} onValueChange={field.onChange}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select category" />
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
              <span className="text-red-500 text-sm">{errors.category_id.message}</span>
            )}
          </div>

          {/* Descriptions */}
          <div className="space-y-2">
            <label className="block text-sm font-medium">Description (English)</label>
            <Controller
              name="description_en"
              control={control}
              render={({ field }) => (
                <Textarea {...field} placeholder="Describe the product..." rows={4} />
              )}
            />
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium">Description (Arabic)</label>
            <Controller
              name="description_ar"
              control={control}
              render={({ field }) => (
                <Textarea {...field} placeholder={t('products.descriptionArPlaceholder')} rows={4} dir="rtl" />
              )}
            />
          </div>

          {/* Availability */}
          <div className="flex items-center space-x-2">
            <Controller
              name="is_available"
              control={control}
              render={({ field }) => (
                <Checkbox checked={field.value} onCheckedChange={field.onChange} />
              )}
            />
            <label className="text-sm font-medium">Available</label>
          </div>

          <div className="flex items-center space-x-2">
            <Controller
              name="featured"
              control={control}
              render={({ field }) => (
                <Checkbox checked={field.value} onCheckedChange={field.onChange} />
              )}
            />
            <label className="text-sm font-medium">Featured Product</label>
          </div>
        </TabsContent>

        {/* Pricing Tab */}
        <TabsContent value="pricing" className="space-y-4">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Price (EGP) *</label>
              <Controller
                name="price"
                control={control}
                rules={{ required: 'Price is required', min: { value: 0.01, message: 'Price must be greater than 0' } }}
                render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />}
              />
              {errors.price && <span className="text-red-500 text-sm">{errors.price.message}</span>}
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Cost (EGP)</label>
              <Controller
                name="cost"
                control={control}
                render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Discount Price (EGP)</label>
              <Controller
                name="discount_price"
                control={control}
                render={({ field }) => <Input {...field} type="number" step="0.01" placeholder="0.00" />}
              />
            </div>
          </div>
        </TabsContent>

        {/* Inventory Tab */}
        <TabsContent value="inventory" className="space-y-4">
          <div className="flex items-center space-x-2 mb-4">
            <Controller
              name="track_inventory"
              control={control}
              render={({ field }) => (
                <Checkbox checked={field.value} onCheckedChange={field.onChange} />
              )}
            />
            <label className="text-sm font-medium">Track Inventory</label>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="block text-sm font-medium">Quantity in Stock</label>
              <Controller
                name="quantity_in_stock"
                control={control}
                render={({ field }) => <Input {...field} type="number" min="0" />}
              />
            </div>

            <div className="space-y-2">
              <label className="block text-sm font-medium">Low Stock Threshold</label>
              <Controller
                name="low_stock_threshold"
                control={control}
                render={({ field }) => <Input {...field} type="number" min="0" />}
              />
            </div>
          </div>
        </TabsContent>

        {/* Details Tab */}
        <TabsContent value="details" className="space-y-4">
          <div className="space-y-4">
            <h3 className="font-medium">Classification</h3>
            <div className="space-y-2">
              <div className="flex items-center space-x-2">
                <Controller
                  name="is_vegetarian"
                  control={control}
                  render={({ field }) => (
                    <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                  )}
                />
                <label className="text-sm">Vegetarian</label>
              </div>
              <div className="flex items-center space-x-2">
                <Controller
                  name="is_vegan"
                  control={control}
                  render={({ field }) => (
                    <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                  )}
                />
                <label className="text-sm">Vegan</label>
              </div>
              <div className="flex items-center space-x-2">
                <Controller
                  name="is_spicy"
                  control={control}
                  render={({ field }) => (
                    <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                  )}
                />
                <label className="text-sm">Spicy</label>
              </div>
              <div className="flex items-center space-x-2">
                <Controller
                  name="is_gluten_free"
                  control={control}
                  render={({ field }) => (
                    <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                  )}
                />
                <label className="text-sm">Gluten Free</label>
              </div>
            </div>
          </div>
        </TabsContent>
      </Tabs>

      {/* Form Actions */}
      <div className="flex justify-end gap-2 pt-4 border-t">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? 'Saving...' : product ? 'Update Product' : 'Create Product'}
        </Button>
      </div>
    </form>
  )
}
