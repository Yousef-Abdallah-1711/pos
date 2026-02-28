export type SalesGranularity = 'day' | 'week' | 'month'

export interface SalesReportFilters {
  date_from: string
  date_to: string
  granularity: SalesGranularity
  top_limit?: number
}

export type ReportsExportType = 'sales' | 'profitability' | 'inventory'
export type ReportsExportFormat = 'csv' | 'pdf'

export interface ReportExportRequest {
  report_type: ReportsExportType
  format: ReportsExportFormat
  date_from: string
  date_to: string
  granularity?: SalesGranularity
  top_limit?: number
  locale?: 'en' | 'ar'
}

export interface ReportExportFile {
  blob: Blob
  filename: string
  contentType: string
}

export interface SalesSummary {
  total_revenue: number
  total_orders: number
  avg_ticket: number
}

export interface SalesTrendPoint {
  bucket: string
  revenue: number
  orders: number
}

export interface SalesHourlyPoint {
  hour: number
  revenue: number
  orders: number
}

export interface SalesTopItem {
  product_id: number
  name: string
  quantity: number
  revenue: number
}

export interface SalesTopCategory {
  category_id?: number
  name: string
  quantity: number
  revenue: number
}

export interface SalesReportData {
  filters: Record<string, unknown>
  summary: SalesSummary
  trend: SalesTrendPoint[]
  hourly: SalesHourlyPoint[]
  top_items: SalesTopItem[]
  top_categories: SalesTopCategory[]
  generated_at: string
}

export interface ProfitabilitySummary {
  total_revenue: number
  total_cogs: number
  gross_profit: number
  gross_margin_pct: number
  food_cost_pct: number
  units_missing_cost: number
  products_missing_cost_rows: number
}

export interface ProfitabilityItem {
  product_id: number
  name: string
  quantity: number
  revenue: number
  cogs: number
  gross_profit: number
  gross_margin_pct: number
  food_cost_pct: number
  cost_per_unit: number
  missing_cost_data: boolean
}

export interface ProfitabilityReportData {
  filters: Record<string, unknown>
  summary: ProfitabilitySummary
  top_items: ProfitabilityItem[]
  lowest_margin_items: ProfitabilityItem[]
  generated_at: string
}

export interface InventorySummary {
  tracked_products_count: number
  total_current_inventory_value: number
  total_estimated_cogs: number
  average_monthly_turnover: number
  slow_moving_count: number
}

export interface InventoryTurnoverItem {
  product_id: number
  name: string
  cost_per_unit: number
  current_stock_qty: number
  opening_stock_qty: number
  average_inventory_qty: number
  average_inventory_value: number
  sold_qty: number
  estimated_cogs: number
  waste_qty: number
  turnover_ratio: number
  monthly_turnover_ratio: number
  is_slow_moving: boolean
}

export interface InventoryReportData {
  filters: Record<string, unknown>
  summary: InventorySummary
  items: InventoryTurnoverItem[]
  slow_movers: InventoryTurnoverItem[]
  top_turnover: InventoryTurnoverItem[]
  generated_at: string
}

export interface ReportsOverviewData {
  filters: Record<string, unknown>
  sales: {
    summary: SalesSummary
    top_item?: SalesTopItem
  }
  profitability: {
    summary: ProfitabilitySummary
    most_profitable?: ProfitabilityItem
  }
  inventory: {
    summary: InventorySummary
    slowest_moving?: InventoryTurnoverItem
  }
  generated_at: string
}
