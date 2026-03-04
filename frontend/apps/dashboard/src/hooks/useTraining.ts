'use client'

import { useCallback, useEffect, useState } from 'react'
import type { TrainingRecord } from '@/types/hr-advanced'
import {
  createTrainingRecord,
  deleteTrainingRecord,
  listExpiringCertificates,
  listTrainingRecords,
  updateTrainingRecord,
} from '@/lib/api/hrAdvancedApi'

export function useTraining() {
  const [items, setItems] = useState<TrainingRecord[]>([])
  const [expiring, setExpiring] = useState<TrainingRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const [records, expiringRecords] = await Promise.all([
        listTrainingRecords(),
        listExpiringCertificates(30),
      ])
      setItems(records)
      setExpiring(expiringRecords)
    } catch (err: any) {
      setError(err?.response?.data?.error || err?.message || 'Failed to load training records')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void refresh()
  }, [refresh])

  const createItem = async (payload: Partial<TrainingRecord>) => {
    const created = await createTrainingRecord(payload)
    setItems((prev) => [created, ...prev])
  }

  const updateItem = async (id: number, payload: Partial<TrainingRecord>) => {
    await updateTrainingRecord(id, payload)
    await refresh()
  }

  const removeItem = async (id: number) => {
    await deleteTrainingRecord(id)
    setItems((prev) => prev.filter((item) => item.id !== id))
  }

  return { items, expiring, loading, error, refresh, createItem, updateItem, removeItem }
}
