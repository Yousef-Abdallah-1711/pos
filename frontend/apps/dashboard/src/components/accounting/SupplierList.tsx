'use client'

import { useMemo, useState } from 'react'
import { Edit2, Search, Store, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import type { Supplier } from '@/types/accounting'
import {
  SUPPLIER_CATEGORY_OPTIONS,
  formatCurrency,
} from './accounting-utils'

interface SupplierListProps {
  suppliers: Supplier[]
  locale: string
  loading?: boolean
  onEdit: (supplier: Supplier) => void
  onDelete?: (supplier: Supplier) => void
  canDelete?: boolean
}

export function SupplierList({
  suppliers,
  locale,
  loading = false,
  onEdit,
  onDelete,
  canDelete = false,
}: SupplierListProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')

  const filteredSuppliers = useMemo(
    () =>
      suppliers.filter((supplier) => {
        const matchesSearch =
          searchTerm.trim() === '' ||
          [
            supplier.company_name,
            supplier.company_name_ar,
            supplier.supplier_code,
            supplier.contact_person,
            supplier.email,
          ]
            .filter(Boolean)
            .some((value) =>
              value!.toLowerCase().includes(searchTerm.trim().toLowerCase())
            )

        const matchesCategory =
          categoryFilter === '' || supplier.category === categoryFilter

        return matchesSearch && matchesCategory
      }),
    [categoryFilter, searchTerm, suppliers]
  )

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
        <div className="grid gap-3 lg:grid-cols-[minmax(0,1fr),220px]">
          <div className="relative">
            <Search className="pointer-events-none absolute left-3 top-3.5 h-4 w-4 text-gray-400" />
            <Input
              value={searchTerm}
              onChange={(event) => setSearchTerm(event.target.value)}
              placeholder="Search suppliers by code, name, or contact"
              className="pl-9"
            />
          </div>
          <select
            value={categoryFilter}
            onChange={(event) => setCategoryFilter(event.target.value)}
            className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
          >
            <option value="">All categories</option>
            {SUPPLIER_CATEGORY_OPTIONS.map((category) => (
              <option key={category} value={category}>
                {category}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
        {filteredSuppliers.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-3 px-6 py-16 text-center">
            <div className="rounded-full bg-gray-100 p-4 dark:bg-gray-900">
              <Store className="h-8 w-8 text-gray-400" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                No suppliers found
              </h3>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                Add a supplier or relax the filters to see matching records.
              </p>
            </div>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
              <thead className="bg-gray-50 dark:bg-gray-900/40">
                <tr>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Supplier
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Category
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Contact
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Terms
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Balance
                  </th>
                  <th className="px-4 py-3 text-right font-medium text-gray-500 dark:text-gray-400">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                {filteredSuppliers.map((supplier) => (
                  <tr key={supplier.id} className="hover:bg-gray-50 dark:hover:bg-gray-900/20">
                    <td className="px-4 py-4">
                      <div className="font-semibold text-gray-900 dark:text-white">
                        {supplier.company_name}
                      </div>
                      <div className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                        {supplier.supplier_code}
                      </div>
                    </td>
                    <td className="px-4 py-4 capitalize text-gray-600 dark:text-gray-300">
                      {supplier.category || '-'}
                    </td>
                    <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                      <div>{supplier.contact_person || '-'}</div>
                      <div className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                        {supplier.email || supplier.phone || '-'}
                      </div>
                    </td>
                    <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                      {supplier.payment_terms || '-'}
                    </td>
                    <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                      {formatCurrency(supplier.current_balance, locale, supplier.currency)}
                    </td>
                    <td className="px-4 py-4">
                      <div className="flex justify-end gap-2">
                        <Button type="button" variant="outline" size="sm" onClick={() => onEdit(supplier)}>
                          <Edit2 className="mr-2 h-4 w-4" />
                          Edit
                        </Button>
                        {canDelete && onDelete ? (
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={() => onDelete(supplier)}
                            className="text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        ) : null}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
