'use client'

import { Download } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { Button } from '@/components/ui/button'
import { exportFinancialReportToPDF } from '@/lib/pdf-export'

interface ReportExportButtonProps {
  filename: string
  title: string
  rows: Array<Record<string, string | number>>
  template?: 'generic' | 'profit-loss' | 'balance-sheet' | 'vat'
}

export function ReportExportButton({ filename, title, rows, template = 'generic' }: ReportExportButtonProps) {
  const t = useTranslations('reports.financial')
  const onExport = async () => {
    await exportFinancialReportToPDF(title, rows, filename, template)
  }

  return (
    <Button type="button" variant="outline" onClick={() => void onExport()} className="gap-2">
      <Download className="h-4 w-4" />
      {t('common.exportPdf')}
    </Button>
  )
}
