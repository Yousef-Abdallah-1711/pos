'use client'

import { useState } from 'react'
import { usePathname } from 'next/navigation'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Pagination,
  PaginationContent,
  PaginationItem,
} from '@/components/ui/pagination'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Skeleton } from '@/components/ui/skeleton'
import { Edit2, Trash2, Search, AlertCircle, Eye } from 'lucide-react'
import { toast } from 'sonner'

interface Employee {
  id: number
  first_name: string
  last_name: string
  email: string
  employee_code: string
  position: string
  department?: string
  status: string
  is_active: boolean
}

interface EmployeeListProps {
  employees: Employee[]
  loading: boolean
  onEdit: (employee: Employee) => void
  onDelete?: (id: number) => void
  canDelete?: boolean
}

export function EmployeeList({ employees, loading, onEdit, onDelete, canDelete = false }: EmployeeListProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'
  const [searchText, setSearchText] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const [itemsPerPage] = useState(10)
  const [deleteConfirm, setDeleteConfirm] = useState<{ open: boolean; id: number | null }>({
    open: false,
    id: null,
  })
  const [viewMode, setViewMode] = useState<'table' | 'card'>('table')

  const filteredEmployees = employees.filter(
    (emp) =>
      emp.first_name.toLowerCase().includes(searchText.toLowerCase()) ||
      emp.last_name.toLowerCase().includes(searchText.toLowerCase()) ||
      emp.email.toLowerCase().includes(searchText.toLowerCase()) ||
      emp.employee_code.toLowerCase().includes(searchText.toLowerCase())
  )

  const totalPages = Math.ceil(filteredEmployees.length / itemsPerPage)
  const startIdx = (currentPage - 1) * itemsPerPage
  const paginatedEmployees = filteredEmployees.slice(startIdx, startIdx + itemsPerPage)

  const handleDeleteClick = (id: number) => {
    setDeleteConfirm({ open: true, id })
  }

  const handleConfirmDelete = async () => {
    if (deleteConfirm.id) {
      await onDelete?.(deleteConfirm.id)
      setDeleteConfirm({ open: false, id: null })
      setCurrentPage(1)
    }
  }

  if (loading) {
    return (
      <div className="space-y-4">
        {[...Array(5)].map((_, i) => (
          <Skeleton key={i} className="h-16 w-full" />
        ))}
      </div>
    )
  }

  if (employees.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 bg-gray-50 dark:bg-gray-900/50 rounded-lg">
        <AlertCircle className="h-12 w-12 text-gray-400 mb-4" />
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">No Employees</h3>
        <p className="text-gray-600 dark:text-gray-400">Get started by creating your first employee</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Controls */}
      <div className="flex flex-col md:flex-row gap-4 items-center justify-between">
        <div className="w-full md:w-80 relative">
          <Search className={`absolute ${isRTL ? 'right-3' : 'left-3'} top-3 h-4 w-4 text-gray-400`} />
          <Input
            placeholder="Search by name, email, or code..."
            value={searchText}
            onChange={(e) => {
              setSearchText(e.target.value)
              setCurrentPage(1)
            }}
            className={isRTL ? 'pr-10' : 'pl-10'}
          />
        </div>

        <div className="flex gap-2">
          <Button
            variant={viewMode === 'table' ? 'default' : 'outline'}
            onClick={() => setViewMode('table')}
            size="sm"
          >
            Table
          </Button>
          <Button
            variant={viewMode === 'card' ? 'default' : 'outline'}
            onClick={() => setViewMode('card')}
            size="sm"
          >
            Cards
          </Button>
        </div>
      </div>

      {filteredEmployees.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-8 bg-gray-50 dark:bg-gray-900/50 rounded-lg">
          <AlertCircle className="h-10 w-10 text-gray-400 mb-3" />
          <p className="text-gray-600 dark:text-gray-400">No employees match your search</p>
        </div>
      ) : viewMode === 'table' ? (
        <>
          {/* Table View */}
          <div className="overflow-x-auto bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
            <Table>
              <TableHeader>
                <TableRow className="bg-gray-50 dark:bg-gray-900">
                  <TableHead>Name</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Position</TableHead>
                  <TableHead>Code</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {paginatedEmployees.map((emp) => (
                  <TableRow key={emp.id} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                    <TableCell className="font-medium">
                      <div className="flex items-center gap-2">
                        <div
                          className="h-8 w-8 rounded-full bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white text-xs font-bold"
                        >
                          {emp.first_name[0]}{emp.last_name[0]}
                        </div>
                        <div>
                          <div className="font-medium">{emp.first_name} {emp.last_name}</div>
                          {emp.department && <div className="text-xs text-gray-500">{emp.department}</div>}
                        </div>
                      </div>
                    </TableCell>
                    <TableCell className="text-sm text-gray-600 dark:text-gray-400">{emp.email}</TableCell>
                    <TableCell>{emp.position}</TableCell>
                    <TableCell className="text-sm text-gray-600">{emp.employee_code}</TableCell>
                    <TableCell>
                      <span className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium ${emp.status === 'active'
                        ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300'
                        : emp.status === 'inactive'
                          ? 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300'
                          : 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300'
                        }`}>
                        <span className={`h-2 w-2 rounded-full ${emp.status === 'active' ? 'bg-green-600' : emp.status === 'inactive' ? 'bg-gray-600' : 'bg-yellow-600'
                          }`} />
                        {emp.status}
                      </span>
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex justify-end gap-2">
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => onEdit(emp)}
                          className="gap-1"
                        >
                          <Edit2 className="h-4 w-4" />
                          Edit
                        </Button>
                        {canDelete ? (
                          <Button
                            size="sm"
                            variant="outline"
                            className="text-red-600 hover:text-red-700 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
                            onClick={() => handleDeleteClick(emp.id)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        ) : (
                          <Button
                            size="sm"
                            variant="outline"
                            disabled
                            title="You don't have permission to delete employees"
                            className="text-red-600 hover:text-red-700 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                    disabled={currentPage === 1}
                    className="gap-1"
                  >
                    Previous
                  </Button>
                </PaginationItem>
                <PaginationItem>
                  <span className="text-sm font-medium px-4">
                    Page {currentPage} of {totalPages}
                  </span>
                </PaginationItem>
                <PaginationItem>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
                    disabled={currentPage === totalPages}
                    className="gap-1"
                  >
                    Next
                  </Button>
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          )}
        </>
      ) : (
        <>
          {/* Card View */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {paginatedEmployees.map((emp) => (
              <div
                key={emp.id}
                className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <div className="h-10 w-10 rounded-full bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white font-bold text-sm">
                      {emp.first_name[0]}{emp.last_name[0]}
                    </div>
                    <div>
                      <h3 className="font-semibold">{emp.first_name} {emp.last_name}</h3>
                      <p className="text-xs text-gray-500">{emp.employee_code}</p>
                    </div>
                  </div>
                  <span className={`inline-flex items-center gap-1 px-2 py-1 rounded text-xs font-medium ${emp.status === 'active'
                    ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300'
                    : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300'
                    }`}>
                    {emp.status}
                  </span>
                </div>

                <div className="space-y-2 mb-4 text-sm">
                  <div className="flex items-center text-gray-600 dark:text-gray-400">
                    <span className="font-medium mr-2">Position:</span>
                    <span>{emp.position}</span>
                  </div>
                  {emp.department && (
                    <div className="flex items-center text-gray-600 dark:text-gray-400">
                      <span className="font-medium mr-2">Dept:</span>
                      <span>{emp.department}</span>
                    </div>
                  )}
                  <div className="flex items-center text-gray-600 dark:text-gray-400 break-all">
                    <span className="font-medium mr-2">Email:</span>
                    <span>{emp.email}</span>
                  </div>
                </div>

                <div className="flex gap-2 pt-4 border-t">
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => onEdit(emp)}
                    className="flex-1 gap-1"
                  >
                    <Edit2 className="h-4 w-4" />
                    Edit
                  </Button>
                  {canDelete ? (
                    <Button
                      size="sm"
                      variant="outline"
                      className="text-red-600 hover:text-red-700 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
                      onClick={() => handleDeleteClick(emp.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  ) : (
                    <Button
                      size="sm"
                      variant="outline"
                      disabled
                      title="You don't have permission to delete employees"
                      className="text-red-600 hover:text-red-700 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              </div>
            ))}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                    disabled={currentPage === 1}
                    className="gap-1"
                  >
                    Previous
                  </Button>
                </PaginationItem>
                <PaginationItem>
                  <span className="text-sm font-medium px-4">
                    Page {currentPage} of {totalPages}
                  </span>
                </PaginationItem>
                <PaginationItem>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
                    disabled={currentPage === totalPages}
                    className="gap-1"
                  >
                    Next
                  </Button>
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          )}
        </>
      )}

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteConfirm.open} onOpenChange={(open) => setDeleteConfirm({ ...deleteConfirm, open })}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Employee</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this employee? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteConfirm({ ...deleteConfirm, open: false })}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleConfirmDelete}
              className="bg-red-600 hover:bg-red-700"
            >
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
