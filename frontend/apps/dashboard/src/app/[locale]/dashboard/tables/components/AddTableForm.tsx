'use client'

import { useEffect, useMemo, useState } from 'react'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import {
  type CreateTablePayload,
  type TableLocation,
  type TableShape,
} from '@/lib/api/tablesApi'

const LOCATION_OPTIONS: TableLocation[] = [
  'indoor',
  'outdoor',
  'vip',
  'rooftop',
  'private',
]

const SHAPE_OPTIONS: TableShape[] = ['square', 'circle', 'rectangle']

interface FormState {
  table_number: string
  name: string
  capacity: string
  location: TableLocation
  shape: TableShape
  description: string
}

interface AddTableFormProps {
  initialValues?: Partial<CreateTablePayload>
  suggestedTableNumber?: number
  submitting?: boolean
  showCancel?: boolean
  submitLabel?: string
  cancelLabel?: string
  t: (path: string, fallback?: string) => string
  isRTL: boolean
  onSubmit: (payload: CreateTablePayload) => Promise<void> | void
  onCancel?: () => void
}

const createInitialState = (
  initialValues: Partial<CreateTablePayload> | undefined,
  suggestedTableNumber?: number
): FormState => ({
  table_number: String(initialValues?.table_number ?? suggestedTableNumber ?? 1),
  name: initialValues?.name ?? '',
  capacity: String(initialValues?.capacity ?? 4),
  location: initialValues?.location ?? 'indoor',
  shape: initialValues?.shape ?? 'square',
  description: initialValues?.description ?? '',
})

export function AddTableForm({
  initialValues,
  suggestedTableNumber,
  submitting = false,
  showCancel = true,
  submitLabel,
  cancelLabel,
  t,
  isRTL,
  onSubmit,
  onCancel,
}: AddTableFormProps) {
  const [formState, setFormState] = useState<FormState>(() =>
    createInitialState(initialValues, suggestedTableNumber)
  )

  const initialValuesKey = useMemo(() => {
    return JSON.stringify({
      initialValues,
      suggestedTableNumber,
    })
  }, [initialValues, suggestedTableNumber])

  useEffect(() => {
    setFormState(createInitialState(initialValues, suggestedTableNumber))
  }, [initialValuesKey, initialValues, suggestedTableNumber])

  const handleFormChange = (
    event: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = event.target
    setFormState((prev) => ({
      ...prev,
      [name]: value,
    }))
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const tableNumber = Number.parseInt(formState.table_number, 10)
    if (!Number.isFinite(tableNumber) || tableNumber <= 0) {
      toast.error(t('tables.invalidTableNumber'))
      return
    }

    const capacity = Number.parseInt(formState.capacity, 10)
    if (!Number.isFinite(capacity) || capacity < 1 || capacity > 50) {
      toast.error(t('tables.invalidCapacity'))
      return
    }

    await onSubmit({
      table_number: tableNumber,
      name: formState.name.trim(),
      capacity,
      location: formState.location,
      shape: formState.shape,
      description: formState.description.trim(),
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div>
          <label className="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('tables.tableNumber')}
          </label>
          <input
            type="number"
            min={1}
            name="table_number"
            value={formState.table_number}
            onChange={handleFormChange}
            className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary dark:border-gray-600 dark:bg-gray-700 dark:text-white"
          />
        </div>

        <div>
          <label className="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('tables.capacity')}
          </label>
          <input
            type="number"
            min={1}
            max={50}
            name="capacity"
            value={formState.capacity}
            onChange={handleFormChange}
            className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary dark:border-gray-600 dark:bg-gray-700 dark:text-white"
          />
        </div>

        <div>
          <label className="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('tables.tableName')}
          </label>
          <input
            type="text"
            name="name"
            value={formState.name}
            onChange={handleFormChange}
            className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary dark:border-gray-600 dark:bg-gray-700 dark:text-white"
          />
        </div>

        <div>
          <label className="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('tables.location')}
          </label>
          <select
            name="location"
            value={formState.location}
            onChange={handleFormChange}
            className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary dark:border-gray-600 dark:bg-gray-700 dark:text-white"
          >
            {LOCATION_OPTIONS.map((location) => (
              <option key={location} value={location}>
                {t(`tables.locations.${location}`)}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('tables.shape')}
          </label>
          <select
            name="shape"
            value={formState.shape}
            onChange={handleFormChange}
            className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary dark:border-gray-600 dark:bg-gray-700 dark:text-white"
          >
            {SHAPE_OPTIONS.map((shape) => (
              <option key={shape} value={shape}>
                {t(`tables.shapes.${shape}`)}
              </option>
            ))}
          </select>
        </div>

        <div className="md:col-span-2">
          <label className="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('tables.description')}
          </label>
          <textarea
            name="description"
            rows={3}
            value={formState.description}
            onChange={handleFormChange}
            className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary dark:border-gray-600 dark:bg-gray-700 dark:text-white"
          />
        </div>
      </div>

      <div className="flex justify-end gap-2 pt-2">
        {showCancel && (
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
            disabled={submitting}
          >
            {cancelLabel || t('tables.cancel')}
          </Button>
        )}
        <Button type="submit" disabled={submitting}>
          {submitting && (
            <Loader2
              className={
                isRTL ? 'ml-2 h-4 w-4 animate-spin' : 'mr-2 h-4 w-4 animate-spin'
              }
            />
          )}
          {submitLabel || t('tables.save')}
        </Button>
      </div>
    </form>
  )
}
