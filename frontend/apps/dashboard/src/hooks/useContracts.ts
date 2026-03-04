'use client'

import { useCallback, useEffect, useState } from 'react'
import type { EmploymentContract } from '@/types/hr-advanced'
import {
  createContract,
  deleteContract,
  listContracts,
  listExpiringContracts,
  renewContract,
  updateContract,
} from '@/lib/api/hrAdvancedApi'

export function useContracts() {
  const [items, setItems] = useState<EmploymentContract[]>([])
  const [expiring, setExpiring] = useState<EmploymentContract[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const [records, expiringRecords] = await Promise.all([
        listContracts(),
        listExpiringContracts(30),
      ])
      setItems(records)
      setExpiring(expiringRecords)
    } catch (err: any) {
      setError(err?.response?.data?.error || err?.message || 'Failed to load contracts')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void refresh()
  }, [refresh])

  const createItem = async (payload: Partial<EmploymentContract>) => {
    const created = await createContract(payload)
    setItems((prev) => [created, ...prev])
  }

  const updateItem = async (id: number, payload: Partial<EmploymentContract>) => {
    await updateContract(id, payload)
    await refresh()
  }

  const removeItem = async (id: number) => {
    await deleteContract(id)
    setItems((prev) => prev.filter((item) => item.id !== id))
  }

  const renewItem = async (id: number, payload: { start_date: string; end_date?: string; salary_amount: number; currency: string }) => {
    await renewContract(id, payload)
    await refresh()
  }

  return { items, expiring, loading, error, refresh, createItem, updateItem, removeItem, renewItem }
}
