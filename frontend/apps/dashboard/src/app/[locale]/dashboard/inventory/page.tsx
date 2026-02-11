'use client'

import { useState, useEffect } from 'react'
import { usePathname } from 'next/navigation'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { toast } from 'sonner'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Textarea } from '@/components/ui/textarea'
import { useForm, Controller } from 'react-hook-form'
import { Package, AlertTriangle } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { debugError } from '@/lib/debug'

interface Product {
  id: number
  name_en: string
  quantity_in_stock: number
  low_stock_threshold: number
}

interface InventoryLog {
  id: number
  product_id: number
  product_name: string
  quantity_change: number
  quantity_after: number
  reason: string
  notes?: string
  created_at: string
}

export default function InventoryPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  const [products, setProducts] = useState<Product[]>([])
  const [inventoryLogs, setInventoryLogs] = useState<InventoryLog[]>([])
  const [loading, setLoading] = useState(true)
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)

  const { control, handleSubmit, reset, formState: { errors } } = useForm({
    defaultValues: {
      product_id: '',
      quantity_change: 0,
      reason: 'adjustment',
      notes: '',
    },
  })

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)
      const [productsRes] = await Promise.all([
        api.get('/products'),
      ])
      setProducts(productsRes.data?.products || [])
      setInventoryLogs([])
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to load data')
      debugError('Error loading data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAdjustInventory = (product: Product) => {
    reset({
      product_id: product.id.toString(),
      quantity_change: 0,
      reason: 'adjustment',
      notes: '',
    })
    setSelectedProduct(product)
    setIsDialogOpen(true)
  }

  const onSubmit = async (data: any) => {
    try {
      await api.post('/inventory/adjust', {
        ...data,
        product_id: parseInt(data.product_id),
        quantity_change: parseInt(data.quantity_change),
      })
      toast.success('Inventory adjusted successfully')
      setIsDialogOpen(false)
      loadData()
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to adjust inventory')
      debugError('Error adjusting inventory:', error)
    }
  }

  const lowStockProducts = products.filter(
    (p) => p.quantity_in_stock <= p.low_stock_threshold
  )

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{t('inventory.title')}</h1>
        <p className="text-gray-500 dark:text-gray-400 mt-1">
          {t('inventory.trackInventory')}
        </p>
      </div>

      {/* Low Stock Alert */}
      {lowStockProducts.length > 0 && (
        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
          <div className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-yellow-600 dark:text-yellow-400" />
            <h3 className="font-medium text-yellow-800 dark:text-yellow-200">
              {t('inventory.lowStockAlert')}
            </h3>
          </div>
          <p className="text-sm text-yellow-700 dark:text-yellow-300 mt-1">
            {lowStockProducts.length} {t('inventory.productsLow')}
          </p>
          <ul className="mt-2 space-y-1">
            {lowStockProducts.map((product) => (
              <li key={product.id} className="text-sm text-yellow-700 dark:text-yellow-300">
                • {product.name_en}: {product.quantity_in_stock} {t('inventory.units')} ({t('inventory.threshold')}: {product.low_stock_threshold})
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Current Stock Levels */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow">
        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
            <Package className="h-5 w-5" />
            {t('inventory.currentStock')}
          </h2>
        </div>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t('inventory.product')}</TableHead>
              <TableHead>{t('inventory.currentStockCol')}</TableHead>
              <TableHead>{t('inventory.lowStockThreshold')}</TableHead>
              <TableHead>{t('inventory.status')}</TableHead>
              <TableHead>{t('inventory.actions')}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center py-8">
                  {t('inventory.loading')}
                </TableCell>
              </TableRow>
            ) : products.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center py-8 text-gray-500 dark:text-gray-400">
                  {t('inventory.noProducts')}
                </TableCell>
              </TableRow>
            ) : (
              products.map((product) => {
                const isLowStock = product.quantity_in_stock <= product.low_stock_threshold
                return (
                  <TableRow key={product.id}>
                    <TableCell className="font-medium">{product.name_en}</TableCell>
                    <TableCell>{product.quantity_in_stock}</TableCell>
                    <TableCell>{product.low_stock_threshold}</TableCell>
                    <TableCell>
                      <span
                        className={`px-2 py-1 rounded text-sm ${
                          isLowStock
                            ? 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
                            : 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
                        }`}
                      >
                        {isLowStock ? t('inventory.lowStock') : t('inventory.inStock')}
                      </span>
                    </TableCell>
                    <TableCell>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => handleAdjustInventory(product)}
                      >
                        {t('inventory.adjustStock')}
                      </Button>
                    </TableCell>
                  </TableRow>
                )
              })
            )}
          </TableBody>
        </Table>
      </div>

      {/* Inventory Adjustment Dialog */}
      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t('inventory.adjustInventory')}</DialogTitle>
          </DialogHeader>
          {selectedProduct && (
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div className="bg-gray-50 dark:bg-gray-800/50 p-4 rounded-lg">
                <p className="text-sm font-medium text-gray-900 dark:text-white">
                  {selectedProduct.name_en}
                </p>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  {t('inventory.currentStockCol')}: {selectedProduct.quantity_in_stock} {t('inventory.units')}
                </p>
              </div>

              <div className="space-y-2">
                <label className="block text-sm font-medium">{t('inventory.quantityChange')} *</label>
                <Controller
                  name="quantity_change"
                  control={control}
                  rules={{
                    required: t('inventory.quantityRequired'),
                    validate: (value) => value !== 0 || t('inventory.cannotBeZero')
                  }}
                  render={({ field }) => (
                    <Input
                      {...field}
                      type="number"
                      placeholder="Enter positive for addition, negative for reduction"
                    />
                  )}
                />
                <p className="text-xs text-gray-500 dark:text-gray-400">
                  {t('inventory.helperText')}
                </p>
                {errors.quantity_change && (
                  <span className="text-red-500 text-sm">{errors.quantity_change.message}</span>
                )}
              </div>

              <div className="space-y-2">
                <label className="block text-sm font-medium">{t('inventory.reason')} *</label>
                <Controller
                  name="reason"
                  control={control}
                  rules={{ required: t('inventory.reasonRequired') }}
                  render={({ field }) => (
                    <Select value={field.value} onValueChange={field.onChange}>
                      <SelectTrigger>
                        <SelectValue placeholder={t('inventory.reasonPlaceholder')} />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="purchase">{t('inventory.reasons.purchase')}</SelectItem>
                        <SelectItem value="sale">{t('inventory.reasons.sale')}</SelectItem>
                        <SelectItem value="adjustment">{t('inventory.reasons.adjustment')}</SelectItem>
                        <SelectItem value="waste">{t('inventory.reasons.waste')}</SelectItem>
                        <SelectItem value="return">{t('inventory.reasons.return')}</SelectItem>
                      </SelectContent>
                    </Select>
                  )}
                />
                {errors.reason && <span className="text-red-500 text-sm">{errors.reason.message}</span>}
              </div>

              <div className="space-y-2">
                <label className="block text-sm font-medium">{t('inventory.notes')}</label>
                <Controller
                  name="notes"
                  control={control}
                  render={({ field }) => (
                    <Textarea {...field} placeholder={t('inventory.notesPlaceholder')} rows={3} />
                  )}
                />
              </div>

              <div className="flex justify-end gap-2 pt-4 border-t">
                <Button type="button" variant="outline" onClick={() => setIsDialogOpen(false)}>
                  {t('inventory.cancel')}
                </Button>
                <Button type="submit">
                  {t('inventory.adjustInventory')}
                </Button>
              </div>
            </form>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}
