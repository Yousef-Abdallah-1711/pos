export interface CategoryRecord {
  name?: string
  name_ar?: string
  description?: string
  description_ar?: string
  display_order?: number
  is_active?: boolean
  [key: string]: unknown
}

export interface CategoryFormValues {
  name: string
  name_ar: string
  description: string
  description_ar: string
  display_order: number | string
  is_active: boolean
}

export interface CategoryFormPageProps {
  category?: CategoryRecord | null
  // The page-level handlers decide how to persist; they receive the raw form values.
  onSubmit: (values: CategoryFormValues) => Promise<void> | void
  onCancel: () => void
}
