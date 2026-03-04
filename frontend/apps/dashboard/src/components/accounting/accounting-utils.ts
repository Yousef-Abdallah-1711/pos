import type { Account, AccountType } from '@/types/accounting'

export const ACCOUNT_TYPE_OPTIONS: Array<{ value: AccountType; label: string }> = [
  { value: 'asset', label: 'Asset' },
  { value: 'liability', label: 'Liability' },
  { value: 'equity', label: 'Equity' },
  { value: 'revenue', label: 'Revenue' },
  { value: 'expense', label: 'Expense' },
]

export const DEFAULT_NORMAL_BALANCE: Record<AccountType, 'debit' | 'credit'> = {
  asset: 'debit',
  expense: 'debit',
  liability: 'credit',
  equity: 'credit',
  revenue: 'credit',
}

export const SUPPLIER_CATEGORY_OPTIONS = [
  'food',
  'beverage',
  'packaging',
  'equipment',
  'services',
]

export const SUPPLIER_PAYMENT_TERMS = [
  'cash_on_delivery',
  'net_7',
  'net_15',
  'net_30',
  'net_60',
]

export const CUSTOMER_TYPE_OPTIONS = [
  'individual',
  'corporate',
  'house_account',
  'delivery_partner',
]

export const CUSTOMER_PAYMENT_TERMS = [
  'immediate',
  'net_7',
  'net_15',
  'net_30',
]

export const BANK_ACCOUNT_TYPE_OPTIONS = [
  { value: 'checking', label: 'Checking' },
  { value: 'savings', label: 'Savings' },
  { value: 'credit_card', label: 'Credit Card' },
  { value: 'petty_cash', label: 'Petty Cash' },
] as const

export const INVOICE_STATUS_OPTIONS = [
  'draft',
  'sent',
  'approved',
  'partially_paid',
  'paid',
  'overdue',
  'void',
]

export const PAYMENT_METHOD_OPTIONS = [
  'cash',
  'card',
  'bank_transfer',
  'wallet',
  'cheque',
]

export function flattenAccounts(accounts: Account[]): Account[] {
  return accounts.flatMap((account) => [
    account,
    ...(account.children ? flattenAccounts(account.children) : []),
  ])
}

export function formatCurrency(
  value: number,
  locale: string = 'en',
  currency: string = 'EGP'
): string {
  return new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-US', {
    style: 'currency',
    currency,
    maximumFractionDigits: 2,
  }).format(Number.isFinite(value) ? value : 0)
}

export function formatDate(value?: string, locale: string = 'en'): string {
  if (!value) {
    return '-'
  }

  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime())) {
    return value
  }

  return new Intl.DateTimeFormat(locale === 'ar' ? 'ar-EG' : 'en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  }).format(parsed)
}

export function getAccountTypeClasses(accountType: string): string {
  const classes: Record<string, string> = {
    asset: 'border-blue-200 bg-blue-50 text-blue-700 dark:border-blue-900 dark:bg-blue-950/30 dark:text-blue-300',
    liability:
      'border-red-200 bg-red-50 text-red-700 dark:border-red-900 dark:bg-red-950/30 dark:text-red-300',
    equity:
      'border-violet-200 bg-violet-50 text-violet-700 dark:border-violet-900 dark:bg-violet-950/30 dark:text-violet-300',
    revenue:
      'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900 dark:bg-emerald-950/30 dark:text-emerald-300',
    expense:
      'border-amber-200 bg-amber-50 text-amber-700 dark:border-amber-900 dark:bg-amber-950/30 dark:text-amber-300',
  }

  return (
    classes[accountType] ||
    'border-gray-200 bg-gray-50 text-gray-700 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300'
  )
}

export function getJournalStatusClasses(status: string): string {
  const classes: Record<string, string> = {
    draft:
      'border-gray-200 bg-gray-100 text-gray-700 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300',
    posted:
      'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900 dark:bg-emerald-950/30 dark:text-emerald-300',
    void: 'border-red-200 bg-red-50 text-red-700 dark:border-red-900 dark:bg-red-950/30 dark:text-red-300',
  }

  return (
    classes[status] ||
    'border-gray-200 bg-gray-100 text-gray-700 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300'
  )
}

export function getInvoiceStatusClasses(status: string): string {
  const classes: Record<string, string> = {
    draft:
      'border-gray-200 bg-gray-100 text-gray-700 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300',
    sent: 'border-blue-200 bg-blue-50 text-blue-700 dark:border-blue-900 dark:bg-blue-950/30 dark:text-blue-300',
    approved:
      'border-indigo-200 bg-indigo-50 text-indigo-700 dark:border-indigo-900 dark:bg-indigo-950/30 dark:text-indigo-300',
    partially_paid:
      'border-amber-200 bg-amber-50 text-amber-700 dark:border-amber-900 dark:bg-amber-950/30 dark:text-amber-300',
    paid: 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900 dark:bg-emerald-950/30 dark:text-emerald-300',
    overdue: 'border-red-200 bg-red-50 text-red-700 dark:border-red-900 dark:bg-red-950/30 dark:text-red-300',
    void: 'border-gray-300 bg-gray-200 text-gray-700 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300',
  }

  return (
    classes[status] ||
    'border-gray-200 bg-gray-100 text-gray-700 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300'
  )
}

export function parseNumberInput(value: string): number {
  const parsed = Number(value)
  return Number.isFinite(parsed) ? parsed : 0
}

export function formatAccountLabel(account: Account): string {
  return `${account.account_code} - ${account.account_name}`
}
