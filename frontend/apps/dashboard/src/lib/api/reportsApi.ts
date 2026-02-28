import { api } from '@/lib/api'
import type {
  InventoryReportData,
  ProfitabilityReportData,
  ReportExportFile,
  ReportExportRequest,
  ReportsOverviewData,
  SalesReportData,
  SalesReportFilters,
} from '@/types/reports'

export async function getSalesReport(filters: SalesReportFilters): Promise<SalesReportData> {
  const response = await api.get('/reports/sales', { params: filters })
  return response.data?.data as SalesReportData
}

export async function getProfitabilityReport(filters: SalesReportFilters): Promise<ProfitabilityReportData> {
  const response = await api.get('/reports/profitability', { params: filters })
  return response.data?.data as ProfitabilityReportData
}

export async function getInventoryReport(filters: SalesReportFilters): Promise<InventoryReportData> {
  const response = await api.get('/reports/inventory', { params: filters })
  return response.data?.data as InventoryReportData
}

export async function getReportsOverview(filters: Omit<SalesReportFilters, 'granularity'>): Promise<ReportsOverviewData> {
  const response = await api.get('/reports/overview', { params: filters })
  return response.data?.data as ReportsOverviewData
}

export async function exportReportFile(payload: ReportExportRequest): Promise<ReportExportFile> {
  const response = await api.post('/reports/export', payload, {
    responseType: 'blob',
  })

  const contentType = String(response.headers?.['content-type'] || 'application/octet-stream')
  const filename = parseFilenameFromDisposition(response.headers?.['content-disposition']) || buildFallbackReportFilename(payload)

  return {
    blob: response.data as Blob,
    filename,
    contentType,
  }
}

function parseFilenameFromDisposition(contentDisposition?: string): string | null {
  if (!contentDisposition) return null

  const utf8Match = contentDisposition.match(/filename\*=UTF-8''([^;]+)/i)
  if (utf8Match?.[1]) {
    try {
      return decodeURIComponent(utf8Match[1].trim())
    } catch {
      return utf8Match[1].trim()
    }
  }

  const plainMatch = contentDisposition.match(/filename="?([^";]+)"?/i)
  return plainMatch?.[1]?.trim() || null
}

function buildFallbackReportFilename(payload: ReportExportRequest): string {
  const ext = payload.format === 'pdf' ? 'pdf' : 'csv'
  return `${payload.report_type}-report-${payload.date_from}-to-${payload.date_to}.${ext}`
}
