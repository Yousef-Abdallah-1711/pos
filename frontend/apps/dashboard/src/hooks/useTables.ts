import { useCallback, useEffect, useMemo, useState } from 'react'
import {
  type CreateTablePayload,
  type RestaurantTable,
  type TableStats,
  type TableStatus,
  type UpdateTablePayload,
  tablesApi,
} from '@/lib/api/tablesApi'

export type TableFilter = 'all' | TableStatus

const emptyStats: TableStats = {
  total_tables: 0,
  available: 0,
  occupied: 0,
  reserved: 0,
  maintenance: 0,
  total_capacity: 0,
}

export function useTables(initialFilter: TableFilter = 'all') {
  const [tables, setTables] = useState<RestaurantTable[]>([])
  const [allTables, setAllTables] = useState<RestaurantTable[]>([])
  const [stats, setStats] = useState<TableStats>(emptyStats)
  const [statusFilter, setStatusFilter] = useState<TableFilter>(initialFilter)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    setLoading(true)
    try {
      const selectedStatus = statusFilter === 'all' ? undefined : statusFilter
      const [tablesData, statsData, allTablesData] = await Promise.all([
        tablesApi.list(selectedStatus),
        tablesApi.getStats(),
        selectedStatus ? tablesApi.list() : Promise.resolve([]),
      ])

      setTables(tablesData)
      setAllTables(selectedStatus ? allTablesData : tablesData)
      setStats(statsData)
      setError(null)
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : 'Failed to load tables'
      setError(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [statusFilter])

  useEffect(() => {
    void refresh()
  }, [refresh])

  const createTable = useCallback(
    async (payload: CreateTablePayload) => {
      await tablesApi.create(payload)
      await refresh()
    },
    [refresh]
  )

  const updateTable = useCallback(
    async (tableID: number, payload: UpdateTablePayload) => {
      await tablesApi.update(tableID, payload)
      await refresh()
    },
    [refresh]
  )

  const deleteTable = useCallback(
    async (tableID: number) => {
      await tablesApi.remove(tableID)
      await refresh()
    },
    [refresh]
  )

  const updateStatus = useCallback(
    async (tableID: number, status: TableStatus) => {
      await tablesApi.updateStatus(tableID, status)
      await refresh()
    },
    [refresh]
  )

  const nextTableNumber = useMemo(() => {
    const source = allTables.length > 0 ? allTables : tables
    const maxNumber = source.reduce((max, table) => {
      return table.table_number > max ? table.table_number : max
    }, 0)
    return maxNumber + 1
  }, [allTables, tables])

  return {
    tables,
    stats,
    statusFilter,
    setStatusFilter,
    loading,
    error,
    refresh,
    createTable,
    updateTable,
    deleteTable,
    updateStatus,
    getNextTableNumber: () => nextTableNumber,
  }
}
