import type { BalanceSheetReport } from '@/types/financial-reports'

export function canExportBalanceSheet(
  report: Pick<BalanceSheetReport, 'integrity_status' | 'equation_balanced'> | null | undefined
): boolean {
  return report?.integrity_status === 'RECONCILED' && report.equation_balanced === true
}
