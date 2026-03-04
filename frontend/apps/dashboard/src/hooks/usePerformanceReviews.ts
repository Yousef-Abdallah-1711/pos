'use client'

import { useCallback, useEffect, useState } from 'react'
import type { PerformanceReview } from '@/types/hr-advanced'
import {
  acknowledgePerformanceReview,
  completePerformanceReview,
  createPerformanceReview,
  deletePerformanceReview,
  listPerformanceReviews,
  updatePerformanceReview,
} from '@/lib/api/hrAdvancedApi'

export function usePerformanceReviews() {
  const [items, setItems] = useState<PerformanceReview[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await listPerformanceReviews()
      setItems(data)
    } catch (err: any) {
      setError(err?.response?.data?.error || err?.message || 'Failed to load performance reviews')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void refresh()
  }, [refresh])

  const createItem = async (payload: Partial<PerformanceReview>) => {
    const created = await createPerformanceReview(payload)
    setItems((prev) => [created, ...prev])
  }

  const updateItem = async (id: number, payload: Partial<PerformanceReview>) => {
    await updatePerformanceReview(id, payload)
    await refresh()
  }

  const removeItem = async (id: number) => {
    await deletePerformanceReview(id)
    setItems((prev) => prev.filter((item) => item.id !== id))
  }

  const acknowledge = async (id: number) => {
    await acknowledgePerformanceReview(id)
    await refresh()
  }

  const complete = async (id: number) => {
    await completePerformanceReview(id)
    await refresh()
  }

  return { items, loading, error, refresh, createItem, updateItem, removeItem, acknowledge, complete }
}
