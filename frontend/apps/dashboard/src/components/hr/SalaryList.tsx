'use client'

import { useState } from 'react'
import { usePathname } from 'next/navigation'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Skeleton } from '@/components/ui/skeleton'
import { Edit2, Trash2, Search, AlertCircle, DollarSign } from 'lucide-react'

interface Salary {
  id: number
  employee_id: number
  employee_name?: string
  base_salary: number
  status: string
  pay_period_start: string
  pay_period_end: string
  [key: string]: any
}

interface SalaryListProps {
  salaries: Salary[]
  loading: boolean
  onEdit: (salary: Salary) => void
  onDelete: (id: number) => void
  onMarkPaid?: (id: number) => void
}

export function SalaryList({ salaries, loading, onEdit, onDelete, onMarkPaid }: SalaryListProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'
  const [searchText, setSearchText] = useState('')

  const filteredSalaries = salaries.filter(s => s.employee_name?.toLowerCase().includes(searchText.toLowerCase()))

  const formatDate = (dateString: string) => new Date(dateString).toLocaleDateString()
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'paid': return 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300'
      case 'approved': return 'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300'
      case 'pending': return 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300'
      default: return 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300'
    }
  }

  if (loading) {
    return <div className="space-y-4">{[...Array(5)].map((_, i) => <Skeleton key={i} className="h-16 w-full" />)}</div>
  }

  if (salaries.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 bg-gray-50 dark:bg-gray-900/50 rounded-lg">
        <DollarSign className="h-12 w-12 text-gray-400 mb-4" />
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">No Salary Records</h3>
        <p className="text-gray-600 dark:text-gray-400">Create your first salary record</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="w-full md:w-80 relative">
        <Search className={`absolute ${isRTL ? 'right-3' : 'left-3'} top-3 h-4 w-4 text-gray-400`} />
        <Input placeholder="Search by employee name..." value={searchText} onChange={(e) => setSearchText(e.target.value)} className={isRTL ? 'pr-10' : 'pl-10'} />
      </div>

      {filteredSalaries.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-8 bg-gray-50 dark:bg-gray-900/50 rounded-lg">
          <AlertCircle className="h-10 w-10 text-gray-400 mb-3" />
          <p className="text-gray-600 dark:text-gray-400">No salary records match your search</p>
        </div>
      ) : (
        <div className="overflow-x-auto bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
          <Table>
            <TableHeader>
              <TableRow className="bg-gray-50 dark:bg-gray-900">
                <TableHead>Employee</TableHead>
                <TableHead>Period</TableHead>
                <TableHead>Base Salary</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredSalaries.map(salary => (
                <TableRow key={salary.id} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                  <TableCell className="font-medium">{salary.employee_name || `Employee #${salary.employee_id}`}</TableCell>
                  <TableCell className="text-sm text-gray-600">{formatDate(salary.pay_period_start)} - {formatDate(salary.pay_period_end)}</TableCell>
                  <TableCell className="font-medium">{salary.base_salary?.toFixed(2) || '0.00'}</TableCell>
                  <TableCell>
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(salary.status)}`}>
                      {t(`hr.payroll.${salary.status}`, salary.status)}
                    </span>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      {salary.status !== 'paid' && (
                        <Button size="sm" variant="outline" onClick={() => onEdit(salary)} className="gap-1">
                          <Edit2 className="h-4 w-4" />
                          Edit
                        </Button>
                      )}
                      {salary.status !== 'paid' && onMarkPaid && (
                        <Button size="sm" className="bg-green-600 hover:bg-green-700" onClick={() => onMarkPaid(salary.id)}>
                          Mark Paid
                        </Button>
                      )}
                      <Button size="sm" variant="outline" className="text-red-600 hover:text-red-700 hover:bg-red-50 dark:text-red-400" onClick={() => onDelete(salary.id)}>
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  )
}
