import React from 'react'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const mockUsePathname = vi.fn()
const mockUseHasPermission = vi.fn()
const mockUseReports = vi.fn()
const mockExportReportFile = vi.fn()
const mockToastSuccess = vi.fn()
const mockToastError = vi.fn()
const mockGetSalesReport = vi.fn()
const mockRefetch = vi.fn()

vi.mock('next/navigation', () => ({
  usePathname: () => mockUsePathname(),
}))

vi.mock('sonner', () => ({
  toast: {
    success: (...args: any[]) => mockToastSuccess(...args),
    error: (...args: any[]) => mockToastError(...args),
  },
}))

vi.mock('@/hooks/useRbac', () => ({
  useHasPermission: (...args: any[]) => mockUseHasPermission(...args),
}))

vi.mock('@/hooks/useReports', () => ({
  useReports: (...args: any[]) => mockUseReports(...args),
}))

vi.mock('@/lib/api/reportsApi', () => ({
  getSalesReport: (...args: any[]) => mockGetSalesReport(...args),
  exportReportFile: (...args: any[]) => mockExportReportFile(...args),
}))

vi.mock('recharts', () => {
  const Wrap = ({ children }: any) => React.createElement('div', null, children)
  const Nil = () => React.createElement('div', null)
  return {
    ResponsiveContainer: Wrap,
    LineChart: Wrap,
    Line: Nil,
    XAxis: Nil,
    YAxis: Nil,
    CartesianGrid: Nil,
    Tooltip: Nil,
    Legend: Nil,
    BarChart: Wrap,
    Bar: Nil,
    PieChart: Wrap,
    Pie: Wrap,
    Cell: Nil,
  }
})

// Import after mocks
import SalesReportsPage from './page'

describe('SalesReportsPage (Phase 10 frontend coverage)', () => {
  const originalCreateObjectURL = URL.createObjectURL
  const originalRevokeObjectURL = URL.revokeObjectURL

  beforeEach(() => {
    mockUsePathname.mockReturnValue('/en/dashboard/reports/sales')
    mockUseHasPermission.mockReturnValue({ hasPermission: true, loading: false })
    mockUseReports.mockReset()
    mockExportReportFile.mockReset()
    mockGetSalesReport.mockReset()
    mockToastSuccess.mockReset()
    mockToastError.mockReset()
    mockRefetch.mockReset()

    URL.createObjectURL = vi.fn(() => 'blob:mock') as any
    URL.revokeObjectURL = vi.fn() as any

    mockUseReports.mockImplementation((opts: any) => ({
      data: {
        generated_at: '2026-02-23T00:00:00Z',
        summary: { total_revenue: 1000, total_orders: 10, avg_ticket: 100 },
        trend: [{ bucket: '2026-02-23', revenue: 1000, orders: 10 }],
        hourly: [{ hour: 12, revenue: 1000, orders: 10 }],
        top_items: [{ name: 'Burger', quantity: 10, revenue: 1000 }],
        top_categories: [{ name: 'Main', quantity: 10, revenue: 1000 }],
      },
      loading: false,
      refreshing: false,
      error: null,
      lastUpdatedAt: new Date('2026-02-23T00:01:00Z'),
      refetch: mockRefetch,
      setData: vi.fn(),
    }))
  })

  afterEach(() => {
    URL.createObjectURL = originalCreateObjectURL
    URL.revokeObjectURL = originalRevokeObjectURL
  })

  it('updates filter state and passes updated deps into useReports', async () => {
    render(<SalesReportsPage />)

    const dateInputs = screen.getAllByDisplayValue(/\d{4}-\d{2}-\d{2}/)
    fireEvent.change(dateInputs[0], { target: { value: '2026-02-01' } })
    fireEvent.change(dateInputs[1], { target: { value: '2026-02-15' } })

    const select = screen.getByRole('combobox')
    fireEvent.change(select, { target: { value: 'week' } })

    await waitFor(() => {
      const latestCall = mockUseReports.mock.calls.at(-1)?.[0]
      expect(latestCall.deps).toEqual(['2026-02-01', '2026-02-15', 'week'])
    })
  })

  it('exports CSV using current filters and shows error toast when export fails', async () => {
    mockExportReportFile.mockRejectedValueOnce(new Error('export failed'))

    render(<SalesReportsPage />)

    const dateInputs = screen.getAllByDisplayValue(/\d{4}-\d{2}-\d{2}/)
    fireEvent.change(dateInputs[0], { target: { value: '2026-02-01' } })
    fireEvent.change(dateInputs[1], { target: { value: '2026-02-22' } })
    fireEvent.change(screen.getByRole('combobox'), { target: { value: 'month' } })

    fireEvent.click(screen.getByRole('button', { name: /csv/i }))

    await waitFor(() => {
      expect(mockExportReportFile).toHaveBeenCalledWith(
        expect.objectContaining({
          report_type: 'sales',
          format: 'csv',
          date_from: '2026-02-01',
          date_to: '2026-02-22',
          granularity: 'month',
        })
      )
    })

    await waitFor(() => {
      expect(mockToastError).toHaveBeenCalledTimes(1)
    })
    expect(mockToastSuccess).not.toHaveBeenCalled()
  })
})
