'use client'

import { useCallback, useEffect, useState } from 'react'
import {
  getHRAbsenteeismReport,
  getHRHeadcountReport,
  getHRLeaveBalancesReport,
  getHROvertimeReport,
  getHRTrainingComplianceReport,
  getHRTurnoverReport,
} from '@/lib/api/hrAdvancedApi'
import type {
  HRAbsenteeismReport,
  HRHeadcountReport,
  HRLeaveBalancesReport,
  HROvertimeReport,
  HRTrainingComplianceReport,
  HRTurnoverReport,
} from '@/types/hr-advanced'

interface UseHRReportsFilters {
  date_from: string
  date_to: string
  year: number
}

interface HRReportsState {
  headcount: HRHeadcountReport | null
  turnover: HRTurnoverReport | null
  absenteeism: HRAbsenteeismReport | null
  overtime: HROvertimeReport | null
  leaveBalances: HRLeaveBalancesReport | null
  trainingCompliance: HRTrainingComplianceReport | null
}

export function useHRReports(filters: UseHRReportsFilters) {
  const [data, setData] = useState<HRReportsState>({
    headcount: null,
    turnover: null,
    absenteeism: null,
    overtime: null,
    leaveBalances: null,
    trainingCompliance: null,
  })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const [headcount, turnover, absenteeism, overtime, leaveBalances, trainingCompliance] = await Promise.all([
        getHRHeadcountReport(),
        getHRTurnoverReport({ date_from: filters.date_from, date_to: filters.date_to }),
        getHRAbsenteeismReport({ date_from: filters.date_from, date_to: filters.date_to }),
        getHROvertimeReport({ date_from: filters.date_from, date_to: filters.date_to }),
        getHRLeaveBalancesReport(filters.year),
        getHRTrainingComplianceReport(),
      ])

      setData({
        headcount,
        turnover,
        absenteeism,
        overtime,
        leaveBalances,
        trainingCompliance,
      })
    } catch (err: any) {
      setError(err?.response?.data?.error || err?.message || 'Failed to load HR reports')
    } finally {
      setLoading(false)
    }
  }, [filters.date_from, filters.date_to, filters.year])

  useEffect(() => {
    void refresh()
  }, [refresh])

  return {
    ...data,
    loading,
    error,
    refresh,
  }
}
