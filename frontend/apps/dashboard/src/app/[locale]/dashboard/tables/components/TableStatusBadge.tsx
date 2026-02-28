import { type TableStatus } from '@/lib/api/tablesApi'

const statusBadgeStyles: Record<TableStatus, string> = {
  available:
    'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300',
  occupied: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300',
  reserved:
    'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300',
  maintenance:
    'bg-gray-200 text-gray-800 dark:bg-gray-700 dark:text-gray-200',
}

interface TableStatusBadgeProps {
  status: TableStatus
  label: string
}

export function TableStatusBadge({ status, label }: TableStatusBadgeProps) {
  return (
    <span
      className={`inline-flex rounded-full px-3 py-1 text-xs font-semibold ${statusBadgeStyles[status]}`}
    >
      {label}
    </span>
  )
}

