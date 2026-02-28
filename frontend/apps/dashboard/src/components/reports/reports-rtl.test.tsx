import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { ReportTable } from './ReportTable'
import { LastUpdatedIndicator } from './LastUpdatedIndicator'

describe('reports components RTL rendering checks', () => {
  it('renders ReportTable with Arabic content inside RTL container', () => {
    render(
      <div dir="rtl">
        <ReportTable
          title="???? ????????"
          columns={["??????", "??????", "???????"]}
          rows={[["????", "10", "100 EGP"]]}
          emptyText="?? ???? ??????"
        />
      </div>
    )

    expect(screen.getByText('???? ????????')).toBeDefined()
    expect(screen.getAllByText('??????').length).toBeGreaterThan(0)
    expect(screen.getByText('????')).toBeDefined()
  })

  it('renders LastUpdatedIndicator using Arabic locale formatting path', () => {
    render(
      <div dir="rtl">
        <LastUpdatedIndicator
          locale="ar"
          generatedAtLabel="?? ???????"
          lastUpdatedLabel="??? ?????"
          generatedAt="2026-02-23T10:15:00Z"
          lastUpdatedAt={new Date('2026-02-23T10:16:00Z')}
        />
      </div>
    )

    expect(screen.getByText('?? ???????')).toBeDefined()
    expect(screen.getByText('??? ?????')).toBeDefined()
    // Ensure fallback dash is not shown when dates are provided.
    expect(screen.queryByText('—')).toBeNull()
  })
})
