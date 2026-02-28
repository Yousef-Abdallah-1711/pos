export type TableStatus = 'available' | 'occupied' | 'reserved' | 'maintenance'

export interface TableStatusNotice {
  status: TableStatus
  level: 'warning' | 'critical'
  translationKey: string
}

function isTableStatus(value: string): value is TableStatus {
  return (
    value === 'available' ||
    value === 'occupied' ||
    value === 'reserved' ||
    value === 'maintenance'
  )
}

export function normalizeTableStatus(status?: string | null): TableStatus | null {
  if (!status) {
    return null
  }

  const normalized = status.trim().toLowerCase()
  return isTableStatus(normalized) ? normalized : null
}

export function getTableStatusNotice(status?: string | null): TableStatusNotice | null {
  const normalizedStatus = normalizeTableStatus(status)
  if (!normalizedStatus || normalizedStatus === 'available') {
    return null
  }

  if (normalizedStatus === 'maintenance') {
    return {
      status: normalizedStatus,
      level: 'critical',
      translationKey: 'tableOrdering.statusWarnings.maintenance',
    }
  }

  if (normalizedStatus === 'occupied') {
    return {
      status: normalizedStatus,
      level: 'warning',
      translationKey: 'tableOrdering.statusWarnings.occupied',
    }
  }

  return {
    status: normalizedStatus,
    level: 'warning',
    translationKey: 'tableOrdering.statusWarnings.reserved',
  }
}
