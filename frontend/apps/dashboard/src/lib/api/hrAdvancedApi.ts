import { api } from '@/lib/api'
import type {
  EmploymentContract,
  HRAbsenteeismReport,
  HRHeadcountReport,
  HRLeaveBalancesReport,
  HROvertimeReport,
  HRReportDateFilters,
  HRTrainingComplianceReport,
  HRTurnoverReport,
  InsuranceTracking,
  PerformanceReview,
  Resignation,
  TrainingRecord,
} from '@/types/hr-advanced'

export async function listPerformanceReviews(): Promise<PerformanceReview[]> {
  const response = await api.get('/hr/performance')
  return (response.data?.data || []) as PerformanceReview[]
}

export async function createPerformanceReview(payload: Partial<PerformanceReview>): Promise<PerformanceReview> {
  const response = await api.post('/hr/performance', payload)
  return response.data?.data as PerformanceReview
}

export async function updatePerformanceReview(id: number, payload: Partial<PerformanceReview>): Promise<void> {
  await api.put(`/hr/performance/${id}`, payload)
}

export async function deletePerformanceReview(id: number): Promise<void> {
  await api.delete(`/hr/performance/${id}`)
}

export async function acknowledgePerformanceReview(id: number): Promise<void> {
  await api.post(`/hr/performance/${id}/acknowledge`)
}

export async function completePerformanceReview(id: number): Promise<void> {
  await api.post(`/hr/performance/${id}/complete`)
}

export async function listTrainingRecords(): Promise<TrainingRecord[]> {
  const response = await api.get('/hr/training')
  return (response.data?.data || []) as TrainingRecord[]
}

export async function createTrainingRecord(payload: Partial<TrainingRecord>): Promise<TrainingRecord> {
  const response = await api.post('/hr/training', payload)
  return response.data?.data as TrainingRecord
}

export async function updateTrainingRecord(id: number, payload: Partial<TrainingRecord>): Promise<void> {
  await api.put(`/hr/training/${id}`, payload)
}

export async function deleteTrainingRecord(id: number): Promise<void> {
  await api.delete(`/hr/training/${id}`)
}

export async function listExpiringCertificates(days = 30): Promise<TrainingRecord[]> {
  const response = await api.get('/hr/training/expiring', { params: { days } })
  return (response.data?.data || []) as TrainingRecord[]
}

export async function listInsuranceRecords(): Promise<InsuranceTracking[]> {
  const response = await api.get('/hr/insurance')
  return (response.data?.data || []) as InsuranceTracking[]
}

export async function createInsuranceRecord(payload: Partial<InsuranceTracking>): Promise<InsuranceTracking> {
  const response = await api.post('/hr/insurance', payload)
  return response.data?.data as InsuranceTracking
}

export async function updateInsuranceRecord(id: number, payload: Partial<InsuranceTracking>): Promise<void> {
  await api.put(`/hr/insurance/${id}`, payload)
}

export async function deleteInsuranceRecord(id: number): Promise<void> {
  await api.delete(`/hr/insurance/${id}`)
}

export async function listExpiringInsurance(days = 30): Promise<InsuranceTracking[]> {
  const response = await api.get('/hr/insurance/expiring', { params: { days } })
  return (response.data?.data || []) as InsuranceTracking[]
}

export async function listContracts(): Promise<EmploymentContract[]> {
  const response = await api.get('/hr/contracts')
  return (response.data?.data || []) as EmploymentContract[]
}

export async function createContract(payload: Partial<EmploymentContract>): Promise<EmploymentContract> {
  const response = await api.post('/hr/contracts', payload)
  return response.data?.data as EmploymentContract
}

export async function updateContract(id: number, payload: Partial<EmploymentContract>): Promise<void> {
  await api.put(`/hr/contracts/${id}`, payload)
}

export async function deleteContract(id: number): Promise<void> {
  await api.delete(`/hr/contracts/${id}`)
}

export async function renewContract(id: number, payload: {
  start_date: string
  end_date?: string
  salary_amount: number
  currency: string
}): Promise<void> {
  await api.post(`/hr/contracts/${id}/renew`, payload)
}

export async function listExpiringContracts(days = 30): Promise<EmploymentContract[]> {
  const response = await api.get('/hr/contracts/expiring', { params: { days } })
  return (response.data?.data || []) as EmploymentContract[]
}

export async function listResignations(): Promise<Resignation[]> {
  const response = await api.get('/hr/resignations')
  return (response.data?.data || []) as Resignation[]
}

export async function createResignation(payload: Partial<Resignation>): Promise<Resignation> {
  const response = await api.post('/hr/resignations', payload)
  return response.data?.data as Resignation
}

export async function updateResignation(id: number, payload: Partial<Resignation>): Promise<void> {
  await api.put(`/hr/resignations/${id}`, payload)
}

export async function deleteResignation(id: number): Promise<void> {
  await api.delete(`/hr/resignations/${id}`)
}

export async function acknowledgeResignation(id: number): Promise<void> {
  await api.post(`/hr/resignations/${id}/acknowledge`)
}

export async function completeResignation(id: number): Promise<void> {
  await api.post(`/hr/resignations/${id}/complete`)
}

export async function calculateResignationSettlement(id: number): Promise<number> {
  const response = await api.get(`/hr/resignations/${id}/settlement`)
  return response.data?.data?.settlement_amount || 0
}

export async function getHRHeadcountReport(filters: {
  department?: string
  status?: string
  employment_type?: string
} = {}): Promise<HRHeadcountReport> {
  const response = await api.get('/reports/hr/headcount', { params: filters })
  return response.data?.data as HRHeadcountReport
}

export async function getHRTurnoverReport(filters: HRReportDateFilters): Promise<HRTurnoverReport> {
  const response = await api.get('/reports/hr/turnover', { params: filters })
  return response.data?.data as HRTurnoverReport
}

export async function getHRAbsenteeismReport(filters: HRReportDateFilters): Promise<HRAbsenteeismReport> {
  const response = await api.get('/reports/hr/absenteeism', { params: filters })
  return response.data?.data as HRAbsenteeismReport
}

export async function getHROvertimeReport(filters: HRReportDateFilters): Promise<HROvertimeReport> {
  const response = await api.get('/reports/hr/overtime', { params: filters })
  return response.data?.data as HROvertimeReport
}

export async function getHRLeaveBalancesReport(year?: number): Promise<HRLeaveBalancesReport> {
  const response = await api.get('/reports/hr/leave-balances', {
    params: year ? { year } : undefined,
  })
  return response.data?.data as HRLeaveBalancesReport
}

export async function getHRTrainingComplianceReport(): Promise<HRTrainingComplianceReport> {
  const response = await api.get('/reports/hr/training-compliance')
  return response.data?.data as HRTrainingComplianceReport
}
