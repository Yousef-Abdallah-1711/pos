'use client'

import { useCallback, useEffect, useState } from 'react'
import type { Resignation } from '@/types/hr-advanced'
import {
  acknowledgeResignation,
  calculateResignationSettlement,
  completeResignation,
  createResignation,
  deleteResignation,
  listResignations,
  updateResignation,
} from '@/lib/api/hrAdvancedApi'

export function useResignations() {
  const [items, setItems] = useState<Resignation[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const records = await listResignations()
      setItems(records)
    } catch (err: any) {
      setError(err?.response?.data?.error || err?.message || 'Failed to load resignations')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void refresh()
  }, [refresh])

  const createItem = async (payload: Partial<Resignation>) => {
    const created = await createResignation(payload)
    setItems((prev) => [created, ...prev])
  }

  const updateItem = async (id: number, payload: Partial<Resignation>) => {
    await updateResignation(id, payload)
    await refresh()
  }

  const removeItem = async (id: number) => {
    await deleteResignation(id)
    setItems((prev) => prev.filter((item) => item.id !== id))
  }

  const acknowledge = async (id: number) => {
    await acknowledgeResignation(id)
    await refresh()
  }

  const complete = async (id: number) => {
    await completeResignation(id)
    await refresh()
  }

  const calculateSettlement = async (id: number) => calculateResignationSettlement(id)

  return { items, loading, error, refresh, createItem, updateItem, removeItem, acknowledge, complete, calculateSettlement }
}
