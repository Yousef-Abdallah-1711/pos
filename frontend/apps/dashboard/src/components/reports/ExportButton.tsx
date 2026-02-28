'use client'

import { FileSpreadsheet, FileText } from 'lucide-react'
import { Button } from '@/components/ui/button'

type ExportFormat = 'csv' | 'pdf'

interface ExportButtonProps {
  format: ExportFormat
  label: string
  onClick: () => void
  disabled?: boolean
  variant?: 'default' | 'outline'
}

export function ExportButton({
  format,
  label,
  onClick,
  disabled,
  variant = format === 'pdf' ? 'default' : 'outline',
}: ExportButtonProps) {
  const Icon = format === 'pdf' ? FileText : FileSpreadsheet
  return (
    <Button type="button" variant={variant} onClick={onClick} disabled={disabled}>
      <Icon className="h-4 w-4" />
      {label}
    </Button>
  )
}

