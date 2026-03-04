'use client'

import { Landmark } from 'lucide-react'
import type { BankAccount, BankTransaction } from '@/types/invoices'
import { formatCurrency, formatDate } from './accounting-utils'

interface BankTransactionListProps {
  account?: BankAccount | null
  transactions: BankTransaction[]
  locale: string
  loading?: boolean
}

export function BankTransactionList({
  account,
  transactions,
  locale,
  loading = false,
}: BankTransactionListProps) {
  return (
    <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
      <div className="border-b border-gray-200 px-6 py-4 dark:border-gray-700">
        <div className="flex items-center justify-between gap-4">
          <div>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Bank Transactions
            </h2>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {account
                ? `${account.account_name} • ${account.bank_name}`
                : 'Select a bank account to review the journal-backed transaction history.'}
            </p>
          </div>
          {account ? (
            <div className="rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-2 text-sm font-medium text-emerald-700 dark:border-emerald-900 dark:bg-emerald-950/30 dark:text-emerald-300">
              Current balance {formatCurrency(account.current_balance, locale, account.currency)}
            </div>
          ) : null}
        </div>
      </div>

      {loading ? (
        <div className="space-y-3 p-6">
          <div className="h-12 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
          <div className="h-12 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
          <div className="h-12 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
        </div>
      ) : !account ? (
        <EmptyState message="Choose a bank account from the list to view its transaction stream." />
      ) : transactions.length === 0 ? (
        <EmptyState message="No bank transactions have been posted to this account yet." />
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 text-sm dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-900/40">
              <tr>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Date
                </th>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Entry
                </th>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Reference
                </th>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Description
                </th>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Debit
                </th>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Credit
                </th>
                <th className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  Running balance
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
              {transactions.map((transaction) => (
                <tr
                  key={`${transaction.journal_entry_id}-${transaction.entry_number}-${transaction.transaction_date}`}
                  className="hover:bg-gray-50 dark:hover:bg-gray-900/20"
                >
                  <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                    {formatDate(transaction.transaction_date, locale)}
                  </td>
                  <td className="px-4 py-4 font-mono text-sm font-semibold text-gray-900 dark:text-white">
                    {transaction.entry_number}
                  </td>
                  <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                    {transaction.reference_number || '-'}
                  </td>
                  <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                    {transaction.description}
                  </td>
                  <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                    {formatCurrency(transaction.debit_amount, locale, account.currency)}
                  </td>
                  <td className="px-4 py-4 text-gray-600 dark:text-gray-300">
                    {formatCurrency(transaction.credit_amount, locale, account.currency)}
                  </td>
                  <td className="px-4 py-4 font-semibold text-gray-900 dark:text-white">
                    {formatCurrency(transaction.running_balance, locale, account.currency)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

function EmptyState({ message }: { message: string }) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 px-6 py-16 text-center">
      <div className="rounded-full bg-gray-100 p-4 dark:bg-gray-900">
        <Landmark className="h-8 w-8 text-gray-400" />
      </div>
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
          No transactions yet
        </h3>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{message}</p>
      </div>
    </div>
  )
}
