export interface PerformanceReview {
  id: number
  employee_id: number
  reviewer_id: number
  review_period_start: string
  review_period_end: string
  review_date: string
  overall_rating: number
  goals_achievement: number
  technical_skills: number
  communication: number
  leadership: number
  strengths?: string
  areas_to_improve?: string
  action_plan?: string
  status: string
}

export interface TrainingRecord {
  id: number
  employee_id: number
  training_title: string
  training_provider?: string
  training_type: string
  start_date: string
  end_date?: string | null
  status: string
  certificate_number?: string
  certificate_expiry?: string | null
  notes?: string
}

export interface InsuranceTracking {
  id: number
  employee_id: number
  policy_type: string
  provider_name: string
  policy_number: string
  coverage_start: string
  coverage_end: string
  monthly_premium: number
  status: string
  notes?: string
}

export interface EmploymentContract {
  id: number
  employee_id: number
  contract_type: string
  contract_number: string
  start_date: string
  end_date?: string | null
  salary_amount: number
  currency: string
  status: string
  renewed_from_id?: number | null
  notes?: string
}

export interface Resignation {
  id: number
  employee_id: number
  resignation_date: string
  last_working_date: string
  reason: string
  status: string
  notice_period_days: number
  pending_salary: number
  unused_leave_balance: number
  settlement_amount: number
}

export interface HRReportDateFilters {
  date_from?: string
  date_to?: string
}

export interface HRHeadcountReport {
  filters: {
    department?: string
    status?: string
    employment_type?: string
  }
  summary: {
    total_employees: number
  }
  by_department: Array<{ department: string; count: number }>
  by_status: Array<{ status: string; count: number }>
  by_employment_type: Array<{ employment_type: string; count: number }>
  generated_at: string
}

export interface HRTurnoverReport {
  filters: {
    date_from: string
    date_to: string
  }
  summary: {
    hires: number
    departures: number
    average_headcount: number
    turnover_rate: number
  }
  trend: Array<{ month: string; hires: number; departures: number }>
  generated_at: string
}

export interface HRAbsenteeismReport {
  filters: {
    date_from: string
    date_to: string
  }
  summary: {
    total_records: number
    absent_records: number
    absenteeism_rate: number
  }
  patterns: Array<{ day_of_week: number; day_name: string; absent_days: number }>
  top_absentees: Array<{ employee_id: number; employee_name: string; absent_days: number }>
  generated_at: string
}

export interface HROvertimeReport {
  filters: {
    date_from: string
    date_to: string
  }
  summary: {
    total_overtime_hours: number
    employees_count: number
    total_overtime_cost: number
  }
  items: Array<{
    employee_id: number
    employee_name: string
    department: string
    overtime_hours: number
    hourly_rate: number
    overtime_cost: number
  }>
  generated_at: string
}

export interface HRLeaveBalancesReport {
  filters: {
    year: number
  }
  summary: {
    employees_count: number
    annual_leave_allowance: number
    total_remaining_annual: number
  }
  items: Array<{
    employee_id: number
    employee_name: string
    department: string
    approved_leave_days: number
    annual_leave_used: number
    annual_leave_allowance: number
    annual_leave_remaining: number
  }>
  generated_at: string
}

export interface HRTrainingComplianceReport {
  summary: {
    active_employees: number
    certified_employees: number
    compliance_rate: number
  }
  by_department: Array<{
    department: string
    total_employees: number
    certified_employees: number
    compliance_rate: number
  }>
  by_training_type: Array<{
    training_type: string
    total_records: number
    completed_records: number
    completion_rate: number
  }>
  generated_at: string
}
