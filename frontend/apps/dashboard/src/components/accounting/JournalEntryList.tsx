'use client'

import Link from 'next/link'
import { useMemo, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import type { JournalEntry, JournalEntryFilters } from '@/types/accounting'
import {
  formatCurrency,
  formatDate,
  getJournalStatusClasses,
} from './accounting-utils'

interface JournalEntryListProps {
  entries: JournalEntry[]
  filters: JournalEntryFilters
  locale: string
  loading?: boolean
  onFiltersChange: (filters: JournalEntryFilters) => void
}

export function JournalEntryList({
  entries,
  filters,
  locale,
  loading = false,
  onFiltersChange,
}: JournalEntryListProps) {
  const [visibleRows, setVisibleRows] = useState(150)
  const displayedEntries = useMemo(() => entries.slice(0, visibleRows), [entries, visibleRows])

  const updateFilter = (key: keyof JournalEntryFilters, value: string) => {
    onFiltersChange({
      ...filters,
      [key]: value || undefined,
    })
  }

  if (loading) {
    return (
      <div className="rounded-2xl border border-gray-200 bg-white p-6 dark:border-gray-700 dark:bg-gray-800">
        <div className="space-y-3">
          <div className="h-10 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
          <div className="h-14 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
          <div className="h-14 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
          <div className="h-14 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="rounded-2xl border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
        <div className="grid gap-3 lg:grid-cols-4">
          <Input
            type="date"
            value={filters.start_date || ''}
            onChange={(event) => updateFilter('start_date', event.target.value)}
          />
          <Input
            type="date"
            value={filters.end_date || ''}
            onChange={(event) => updateFilter('end_date', event.target.value)}
          />
          <select
            value={filters.status || ''}
            onChange={(event) => updateFilter('status', event.target.value)}
            className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
          >
            <option value="">All statuses</option>
            <option value="draft">Draft</option>
            <option value="posted">Posted</option>
            <option value="void">Void</option>
          </select>
          <Input
            value={filters.entry_type || ''}
            onChange={(event) => updateFilter('entry_type', event.target.value)}
            placeholder="Filter by entry type"
          />
        </div>
      </div>

      <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
        {entries.length === 0 ? (
          <div className="px-6 py-16 text-center">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
              No journal entries found
            </h3>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Create a manual entry or widen the filters to see more results.
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
              <thead className="bg-gray-50 dark:bg-gray-900/40">
                <tr>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Entry
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Date
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Description
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Type
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Debit
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Credit
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Status
                  </th>
                  <th className="px-4 py-3 text-right font-medium text-gray-500 dark:text-gray-400">
                    Details
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                {displayedEntries.map((entry) => (
                  <tr key={entry.id} className="hover:bg-gray-50 dark:hover:bg-gray-900/20">
                    <td className="px-4 py-4 font-mono text-sm font-semibold text-gray-900 dark:text-white">
                      {entry.entry_number}
                    </td>
                    <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                      {formatDate(entry.entry_date, locale)}
                    </td>
                    <td className="px-4 py-4">
                      <div className="font-medium text-gray-900 dark:text-white">
                        {entry.description}
                      </div>
                      {entry.reference_number ? (
                        <div className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                          Ref: {entry.reference_number}
                        </div>
                      ) : null}
                    </td>
                    <td className="px-4 py-4 capitalize text-gray-600 dark:text-gray-300">
                      {entry.entry_type}
                    </td>
                    <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                      {formatCurrency(entry.total_debit, locale)}
                    </td>
                    <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                      {formatCurrency(entry.total_credit, locale)}
                    </td>
                    <td className="px-4 py-4">
                      <span
                        className={cn(
                          'inline-flex rounded-full border px-2.5 py-1 text-xs font-medium capitalize',
                          getJournalStatusClasses(entry.status)
                        )}
                      >
                        {entry.status}
                      </span>
                    </td>
                    <td className="px-4 py-4 text-right">
                      <Button asChild variant="outline" size="sm">
                        <Link href={`/${locale}/dashboard/accounting/journal-entries/${entry.id}`}>
                          Open
                        </Link>
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {entries.length > displayedEntries.length ? (
        <div className="flex justify-center">
          <Button type="button" variant="outline" onClick={() => setVisibleRows((prev) => prev + 150)}>
            Show more entries
          </Button>
        </div>
      ) : null}
    </div>
  )
}
