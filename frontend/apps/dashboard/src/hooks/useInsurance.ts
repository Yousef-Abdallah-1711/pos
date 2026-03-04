'use client'

import { useCallback, useEffect, useState } from 'react'
import type { InsuranceTracking } from '@/types/hr-advanced'
import {
  createInsuranceRecord,
  deleteInsuranceRecord,
  listExpiringInsurance,
  listInsuranceRecords,
  updateInsuranceRecord,
} from '@/lib/api/hrAdvancedApi'

export function useInsurance() {
  const [items, setItems] = useState<InsuranceTracking[]>([])
  const [expiring, setExpiring] = useState<InsuranceTracking[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const [records, expiringRecords] = await Promise.all([
        listInsuranceRecords(),
        listExpiringInsurance(30),
      ])
      setItems(records)
      setExpiring(expiringRecords)
    } catch (err: any) {
      setError(err?.response?.data?.error || err?.message || 'Failed to load insurance records')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void refresh()
  }, [refresh])

  const createItem = async (payload: Partial<InsuranceTracking>) => {
    const created = await createInsuranceRecord(payload)
    setItems((prev) => [created, ...prev])
  }

  const updateItem = async (id: number, payload: Partial<InsuranceTracking>) => {
    await updateInsuranceRecord(id, payload)
    await refresh()
  }

  const removeItem = async (id: number) => {
    await deleteInsuranceRecord(id)
    setItems((prev) => prev.filter((item) => item.id !== id))
  }

  return { items, expiring, loading, error, refresh, createItem, updateItem, removeItem }
}
