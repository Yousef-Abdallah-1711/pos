import { describe, it, expect, beforeEach, vi } from 'vitest'

vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

import { api } from '@/lib/api'
import {
  exportReportFile,
  getInventoryReport,
  getProfitabilityReport,
  getReportsOverview,
  getSalesReport,
} from './reportsApi'

const mockedApi = api as unknown as {
  get: ReturnType<typeof vi.fn>
  post: ReturnType<typeof vi.fn>
}

describe('reportsApi', () => {
  beforeEach(() => {
    mockedApi.get.mockReset()
    mockedApi.post.mockReset()
  })

  it('requests sales report and returns nested data payload', async () => {
    const filters = { date_from: '2026-02-01', date_to: '2026-02-23', granularity: 'day' as const }
    const payload = { summary: { total_revenue: 100 } }
    mockedApi.get.mockResolvedValue({ data: { success: true, data: payload } })

    const result = await getSalesReport(filters)

    expect(mockedApi.get).toHaveBeenCalledWith('/reports/sales', { params: filters })
    expect(result).toEqual(payload)
  })

  it('calls profitability/inventory/overview endpoints with params', async () => {
    mockedApi.get
      .mockResolvedValueOnce({ data: { data: { summary: { gross_margin_pct: 60 } } } })
      .mockResolvedValueOnce({ data: { data: { summary: { average_monthly_turnover: 3.2 } } } })
      .mockResolvedValueOnce({ data: { data: { sales: { summary: { total_orders: 10 } } } } })

    await getProfitabilityReport({ date_from: '2026-02-01', date_to: '2026-02-23', granularity: 'day' })
    await getInventoryReport({ date_from: '2026-02-01', date_to: '2026-02-23', granularity: 'day' })
    await getReportsOverview({ date_from: '2026-02-01', date_to: '2026-02-23' })

    expect(mockedApi.get).toHaveBeenNthCalledWith(1, '/reports/profitability', {
      params: { date_from: '2026-02-01', date_to: '2026-02-23', granularity: 'day' },
    })
    expect(mockedApi.get).toHaveBeenNthCalledWith(2, '/reports/inventory', {
      params: { date_from: '2026-02-01', date_to: '2026-02-23', granularity: 'day' },
    })
    expect(mockedApi.get).toHaveBeenNthCalledWith(3, '/reports/overview', {
      params: { date_from: '2026-02-01', date_to: '2026-02-23' },
    })
  })

  it('parses utf-8 filename and content-type from export response headers', async () => {
    const blob = new Blob(['csv-data'], { type: 'text/csv' })
    mockedApi.post.mockResolvedValue({
      data: blob,
      headers: {
        'content-type': 'text/csv; charset=utf-8',
        'content-disposition': "attachment; filename*=UTF-8''sales%20report%20feb.csv",
      },
    })

    const result = await exportReportFile({
      report_type: 'sales',
      format: 'csv',
      date_from: '2026-02-01',
      date_to: '2026-02-23',
      granularity: 'day',
      top_limit: 10,
    })

    expect(mockedApi.post).toHaveBeenCalledWith(
      '/reports/export',
      expect.objectContaining({ report_type: 'sales', format: 'csv' }),
      { responseType: 'blob' },
    )
    expect(result.filename).toBe('sales report feb.csv')
    expect(result.contentType).toBe('text/csv; charset=utf-8')
    expect(result.blob).toBe(blob)
  })

  it('falls back to generated filename when content-disposition is missing', async () => {
    mockedApi.post.mockResolvedValue({
      data: new Blob(['pdf-bytes'], { type: 'application/pdf' }),
      headers: {
        'content-type': 'application/pdf',
      },
    })

    const result = await exportReportFile({
      report_type: 'inventory',
      format: 'pdf',
      date_from: '2026-02-01',
      date_to: '2026-02-23',
      top_limit: 10,
    })

    expect(result.filename).toBe('inventory-report-2026-02-01-to-2026-02-23.pdf')
    expect(result.contentType).toBe('application/pdf')
  })
})

