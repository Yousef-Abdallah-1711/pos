'use client'

import { useMemo, useState } from 'react'
import { Edit2, FileText, Search, Trash2, Users } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import type { CustomerAccount } from '@/types/accounting'
import { formatCurrency } from './accounting-utils'

interface CustomerListProps {
  customers: CustomerAccount[]
  locale: string
  loading?: boolean
  onEdit: (customer: CustomerAccount) => void
  onDelete?: (customer: CustomerAccount) => void
  onViewStatement?: (customer: CustomerAccount) => void
  canDelete?: boolean
}

export function CustomerList({
  customers,
  locale,
  loading = false,
  onEdit,
  onDelete,
  onViewStatement,
  canDelete = false,
}: CustomerListProps) {
  const [searchTerm, setSearchTerm] = useState('')

  const filteredCustomers = useMemo(
    () =>
      customers.filter((customer) => {
        if (searchTerm.trim() === '') {
          return true
        }

        return [
          customer.customer_name,
          customer.customer_name_ar,
          customer.customer_code,
          customer.contact_person,
          customer.email,
        ]
          .filter(Boolean)
          .some((value) => value!.toLowerCase().includes(searchTerm.trim().toLowerCase()))
      }),
    [customers, searchTerm]
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
        <div className="relative">
          <Search className="pointer-events-none absolute left-3 top-3.5 h-4 w-4 text-gray-400" />
          <Input
            value={searchTerm}
            onChange={(event) => setSearchTerm(event.target.value)}
            placeholder="Search customers by code, name, or contact"
            className="pl-9"
          />
        </div>
      </div>

      <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
        {filteredCustomers.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-3 px-6 py-16 text-center">
            <div className="rounded-full bg-gray-100 p-4 dark:bg-gray-900">
              <Users className="h-8 w-8 text-gray-400" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                No customers found
              </h3>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                Add a customer account or use a broader search query.
              </p>
            </div>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
              <thead className="bg-gray-50 dark:bg-gray-900/40">
                <tr>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Customer
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                    Type
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
                {filteredCustomers.map((customer) => (
                  <tr key={customer.id} className="hover:bg-gray-50 dark:hover:bg-gray-900/20">
                    <td className="px-4 py-4">
                      <div className="font-semibold text-gray-900 dark:text-white">
                        {customer.customer_name}
                      </div>
                      <div className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                        {customer.customer_code}
                      </div>
                    </td>
                    <td className="px-4 py-4 capitalize text-gray-600 dark:text-gray-300">
                      {customer.customer_type}
                    </td>
                    <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                      <div>{customer.contact_person || '-'}</div>
                      <div className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                        {customer.email || customer.phone || '-'}
                      </div>
                    </td>
                    <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                      {customer.payment_terms || '-'}
                    </td>
                    <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                      {formatCurrency(customer.current_balance, locale, customer.currency)}
                    </td>
                    <td className="px-4 py-4">
                      <div className="flex justify-end gap-2">
                        {onViewStatement ? (
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={() => onViewStatement(customer)}
                          >
                            <FileText className="mr-2 h-4 w-4" />
                            Statement
                          </Button>
                        ) : null}
                        <Button type="button" variant="outline" size="sm" onClick={() => onEdit(customer)}>
                          <Edit2 className="mr-2 h-4 w-4" />
                          Edit
                        </Button>
                        {canDelete && onDelete ? (
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={() => onDelete(customer)}
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
