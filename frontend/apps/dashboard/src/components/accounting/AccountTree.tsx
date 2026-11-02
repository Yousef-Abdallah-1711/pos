'use client'

import { useEffect, useMemo, useState } from 'react'
import { ChevronDown, ChevronRight, Edit2, FolderTree, Search, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import type { Account } from '@/types/accounting'
import {
  ACCOUNT_TYPE_OPTIONS,
  formatCurrency,
  getAccountTypeClasses,
} from './accounting-utils'

interface AccountTreeProps {
  accounts: Account[]
  loading?: boolean
  locale: string
  searchTerm: string
  typeFilter: string
  onSearchChange: (value: string) => void
  onTypeFilterChange: (value: string) => void
  onEdit?: (account: Account) => void
  onDelete?: (account: Account) => void
  canWrite?: boolean
  canDelete?: boolean
}

export function AccountTree({
  accounts,
  loading = false,
  locale,
  searchTerm,
  typeFilter,
  onSearchChange,
  onTypeFilterChange,
  onEdit,
  onDelete,
  canWrite = false,
  canDelete = false,
}: AccountTreeProps) {
  const [expandedIds, setExpandedIds] = useState<Record<number, boolean>>({})

  useEffect(() => {
    if (accounts.length === 0) {
      return
    }

    setExpandedIds((current) => {
      if (Object.keys(current).length > 0) {
        return current
      }

      const next: Record<number, boolean> = {}
      accounts.forEach((account) => {
        next[account.id] = true
      })
      return next
    })
  }, [accounts])

  const filteredAccounts = useMemo(
    () => filterAccounts(accounts, searchTerm.trim().toLowerCase(), typeFilter),
    [accounts, searchTerm, typeFilter]
  )

  const toggleExpanded = (accountId: number) => {
    setExpandedIds((current) => ({
      ...current,
      [accountId]: !current[accountId],
    }))
  }

  if (loading) {
    return (
      <div className="rounded-2xl border border-gray-200 bg-white p-8 dark:border-gray-700 dark:bg-gray-800">
        <div className="space-y-3">
          <div className="h-10 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
          <div className="h-16 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
          <div className="h-16 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
          <div className="h-16 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700" />
        </div>
      </div>
    )
  }

  return (
    <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
      <div className="border-b border-gray-200 p-4 dark:border-gray-700">
        <div className="grid gap-3 lg:grid-cols-[minmax(0,1fr),220px]">
          <div className="relative">
            <Search className="pointer-events-none absolute left-3 top-3.5 h-4 w-4 text-gray-400" />
            <Input
              value={searchTerm}
              onChange={(event) => onSearchChange(event.target.value)}
              placeholder="Search accounts by code or name"
              className="pl-9"
            />
          </div>
          <select
            value={typeFilter}
            onChange={(event) => onTypeFilterChange(event.target.value)}
            className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
          >
            <option value="">All account types</option>
            {ACCOUNT_TYPE_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </div>
      </div>

      {filteredAccounts.length === 0 ? (
        <div className="flex flex-col items-center justify-center gap-3 px-6 py-16 text-center">
          <div className="rounded-full bg-gray-100 p-4 dark:bg-gray-900">
            <FolderTree className="h-8 w-8 text-gray-400" />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
              No accounts found
            </h3>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Adjust the search or filters to find a matching account.
            </p>
          </div>
        </div>
      ) : (
        <div className="divide-y divide-gray-100 dark:divide-gray-700">
          {filteredAccounts.map((account) => (
            <AccountNode
              key={account.id}
              account={account}
              depth={0}
              locale={locale}
              expandedIds={expandedIds}
              onToggle={toggleExpanded}
              onEdit={onEdit}
              onDelete={onDelete}
              canWrite={canWrite}
              canDelete={canDelete}
            />
          ))}
        </div>
      )}
    </div>
  )
}

interface AccountNodeProps {
  account: Account
  depth: number
  locale: string
  expandedIds: Record<number, boolean>
  onToggle: (accountId: number) => void
  onEdit?: (account: Account) => void
  onDelete?: (account: Account) => void
  canWrite: boolean
  canDelete: boolean
}

function AccountNode({
  account,
  depth,
  locale,
  expandedIds,
  onToggle,
  onEdit,
  onDelete,
  canWrite,
  canDelete,
}: AccountNodeProps) {
  const hasChildren = Boolean(account.children && account.children.length > 0)
  const isExpanded = expandedIds[account.id] ?? false

  return (
    <div>
      <div
        className={cn(
          'flex flex-col gap-3 px-4 py-4 lg:flex-row lg:items-center lg:justify-between',
          depth > 0 && 'bg-gray-50/50 dark:bg-gray-900/20'
        )}
        style={{ paddingLeft: `${16 + depth * 20}px` }}
      >
        <div className="flex min-w-0 items-start gap-3">
          {hasChildren ? (
            <button
              type="button"
              onClick={() => onToggle(account.id)}
              className="mt-0.5 rounded-md border border-gray-200 p-1 text-gray-500 transition hover:bg-gray-100 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-gray-700"
            >
              {isExpanded ? (
                <ChevronDown className="h-4 w-4" />
              ) : (
                <ChevronRight className="h-4 w-4" />
              )}
            </button>
          ) : (
            <span className="mt-1 block h-6 w-6 rounded-md border border-dashed border-gray-200 dark:border-gray-700" />
          )}

          <div className="min-w-0">
            <div className="flex flex-wrap items-center gap-2">
              <span className="font-mono text-sm font-semibold text-gray-900 dark:text-white">
                {account.account_code}
              </span>
              <span className="text-base font-semibold text-gray-900 dark:text-white">
                {account.account_name}
              </span>
              {account.account_name_ar ? (
                <span className="text-sm text-gray-500 dark:text-gray-400">
                  {account.account_name_ar}
                </span>
              ) : null}
            </div>
            <div className="mt-2 flex flex-wrap items-center gap-2">
              <span
                className={cn(
                  'inline-flex rounded-full border px-2.5 py-1 text-xs font-medium capitalize',
                  getAccountTypeClasses(account.account_type)
                )}
              >
                {account.account_type}
              </span>
              {account.account_subtype ? (
                <span className="inline-flex rounded-full border border-gray-200 bg-gray-50 px-2.5 py-1 text-xs text-gray-600 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300">
                  {account.account_subtype}
                </span>
              ) : null}
              {account.is_header ? (
                <span className="inline-flex rounded-full border border-sky-200 bg-sky-50 px-2.5 py-1 text-xs text-sky-700 dark:border-sky-900 dark:bg-sky-950/30 dark:text-sky-300">
                  Header
                </span>
              ) : null}
              {account.is_system ? (
                <span className="inline-flex rounded-full border border-purple-200 bg-purple-50 px-2.5 py-1 text-xs text-purple-700 dark:border-purple-900 dark:bg-purple-950/30 dark:text-purple-300">
                  System
                </span>
              ) : null}
              {!account.is_active ? (
                <span className="inline-flex rounded-full border border-red-200 bg-red-50 px-2.5 py-1 text-xs text-red-700 dark:border-red-900 dark:bg-red-950/30 dark:text-red-300">
                  Inactive
                </span>
              ) : null}
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-3 lg:items-end">
          <div className="text-left lg:text-right">
            <div className="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">
              Current balance
            </div>
            <div className="text-sm font-semibold text-gray-900 dark:text-white">
              {formatCurrency(account.current_balance, locale, account.currency)}
            </div>
          </div>

          {(canWrite || canDelete) && (
            <div className="flex flex-wrap gap-2">
              {canWrite && onEdit ? (
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => onEdit(account)}
                >
                  <Edit2 className="mr-2 h-4 w-4" />
                  Edit
                </Button>
              ) : null}
              {canDelete && onDelete ? (
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => onDelete(account)}
                  className="text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
                  disabled={account.is_system}
                >
                  <Trash2 className="mr-2 h-4 w-4" />
                  Delete
                </Button>
              ) : null}
            </div>
          )}
        </div>
      </div>

      {hasChildren && isExpanded && (
        <div className="divide-y divide-gray-100 dark:divide-gray-700">
          {account.children!.map((child) => (
            <AccountNode
              key={child.id}
              account={child}
              depth={depth + 1}
              locale={locale}
              expandedIds={expandedIds}
              onToggle={onToggle}
              onEdit={onEdit}
              onDelete={onDelete}
              canWrite={canWrite}
              canDelete={canDelete}
            />
          ))}
        </div>
      )}
    </div>
  )
}

function filterAccounts(accounts: Account[], searchTerm: string, typeFilter: string): Account[] {
  return accounts
    .map((account): Account | null => {
      const children = filterAccounts(account.children || [], searchTerm, typeFilter)
      const matchesSearch =
        searchTerm === '' ||
        [
          account.account_code,
          account.account_name,
          account.account_name_ar,
          account.account_subtype,
        ]
          .filter(Boolean)
          .some((value) => value!.toLowerCase().includes(searchTerm))
      const matchesType = typeFilter === '' || account.account_type === typeFilter

      if ((matchesSearch && matchesType) || children.length > 0) {
        return { ...account, children }
      }

      return null
    })
    .filter((account): account is Account => account !== null)
}
