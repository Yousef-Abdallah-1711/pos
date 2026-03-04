import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { ReportTable } from './ReportTable'
import { LastUpdatedIndicator } from './LastUpdatedIndicator'

describe('reports components RTL rendering checks', () => {
  it('renders ReportTable with Arabic content inside RTL container', () => {
    render(
      <div dir="rtl">
        <ReportTable
          title="جدول المبيعات"
          columns={["المنتج", "الكمية", "الإيراد"]}
          rows={[["برجر", "10", "100 EGP"]]}
          emptyText="لا توجد بيانات"
        />
      </div>
    )

    expect(screen.getByText('جدول المبيعات')).toBeDefined()
    expect(screen.getByText('المنتج')).toBeDefined()
    expect(screen.getByText('برجر')).toBeDefined()
  })

  it('renders LastUpdatedIndicator using Arabic locale formatting path', () => {
    render(
      <div dir="rtl">
        <LastUpdatedIndicator
          locale="ar"
          generatedAtLabel="وقت التوليد"
          lastUpdatedLabel="آخر تحديث"
          generatedAt="2026-02-23T10:15:00Z"
          lastUpdatedAt={new Date('2026-02-23T10:16:00Z')}
        />
      </div>
    )

    expect(screen.getByText('وقت التوليد')).toBeDefined()
    expect(screen.getByText('آخر تحديث')).toBeDefined()
    // Ensure fallback dash is not shown when dates are provided.
    expect(screen.queryByText('—')).toBeNull()
  })
})
