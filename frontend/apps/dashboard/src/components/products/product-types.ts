export interface ProductCategory {
  id: number
  name: string
}

export interface ProductSummary {
  id: number
  name?: string
  name_en?: string
  name_ar?: string
  price?: number
  category_id?: number | null
  main_image_url?: string
  is_available?: boolean
  quantity_in_stock?: number
}

export interface ProductDetails extends ProductSummary {
  [key: string]: unknown
  description_en?: string
  description_ar?: string
  cost?: number | null
  discount_price?: number | null
  low_stock_threshold?: number
  is_vegetarian?: boolean
  is_vegan?: boolean
  is_spicy?: boolean
  is_gluten_free?: boolean
  featured?: boolean
  track_inventory?: boolean
}

/** Values held by the product form. Numeric inputs stay strings until submit (see toNumber/toInt). */
export interface ProductFormValues {
  name_en: string
  name_ar: string
  category_id: string | number
  description_en: string
  description_ar: string
  price: string | number
  cost: string | number
  discount_price: string | number
  quantity_in_stock: string | number
  low_stock_threshold: string | number
  is_vegetarian: boolean
  is_vegan: boolean
  is_spicy: boolean
  is_gluten_free: boolean
  is_available: boolean
  featured: boolean
  track_inventory: boolean
  discount_percentage?: string | number
  display_order?: string | number
  calories?: string | number
}

export interface ProductFormProps {
  product?: ProductDetails | null
  categories: ProductCategory[]
  onSubmit: (formData: FormData) => Promise<void> | void
  onCancel: () => void
}

export interface ProductListFilters {
  searchText: string
  categoryId: number | null
  sortBy: string
  sortOrder: string
  page: number
  pageSize: number
}

export interface ProductListProps {
  products: ProductSummary[]
  categories: ProductCategory[]
  loading: boolean
  filters: ProductListFilters
  onFilterChange: (filters: ProductListFilters) => void
  onEdit?: (product: ProductSummary) => void
  onDelete?: (productId: number) => void | Promise<void>
  canDelete?: boolean
}
