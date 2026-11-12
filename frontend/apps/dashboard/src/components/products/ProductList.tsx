'use client'

import { useState } from 'react'
import Image from 'next/image'
import { usePathname } from 'next/navigation'
import Link from 'next/link'
import { getLocaleFromPath } from '@/lib/translations'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'
import { Skeleton } from '@/components/ui/skeleton'
import type { ProductListProps } from './product-types'

export function ProductList({
  products,
  categories,
  loading,
  filters,
  onFilterChange,
  onEdit,
  onDelete,
  canDelete = false,
}: ProductListProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const handleSearchChange = (value: string) => {
    onFilterChange({ ...filters, searchText: value, page: 1 })
  }

  const handleCategoryChange = (value: string) => {
    onFilterChange({
      ...filters,
      categoryId: value && value !== 'all' ? parseInt(value) : null,
      page: 1,
    })
  }

  const handleSortChange = (field: string) => {
    const newOrder =
      filters.sortBy === field && filters.sortOrder === 'ASC' ? 'DESC' : 'ASC'
    onFilterChange({
      ...filters,
      sortBy: field,
      sortOrder: newOrder,
      page: 1,
    })
  }

  const handlePageChange = (page: number) => {
    onFilterChange({ ...filters, page })
  }

  if (loading && products.length === 0) {
    return (
      <div className="space-y-2">
        {Array.from({ length: 5 }).map((_, i) => (
          <Skeleton key={i} className="h-12 w-full" />
        ))}
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-end">
        <div className="flex-1 w-full sm:w-auto">
          <label className="block text-sm font-medium mb-2">Search</label>
          <Input
            placeholder="Search products..."
            onChange={(e) => handleSearchChange(e.target.value)}
            defaultValue={filters.searchText}
          />
        </div>

        <div className="w-full sm:w-40">
          <label className="block text-sm font-medium mb-2">Category</label>
          <Select
            value={filters.categoryId?.toString() || 'all'}
            onValueChange={handleCategoryChange}
          >
            <SelectTrigger>
              <SelectValue placeholder="All categories" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All categories</SelectItem>
              {categories.map((cat) => (
                <SelectItem key={cat.id} value={cat.id.toString()}>
                  {cat.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Products Table */}
      <div className="border rounded-lg overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Image</TableHead>
              <TableHead
                onClick={() => handleSortChange('name_en')}
                className="cursor-pointer"
              >
                Name
              </TableHead>
              <TableHead>Category</TableHead>
              <TableHead
                onClick={() => handleSortChange('price')}
                className="cursor-pointer"
              >
                Price
              </TableHead>
              <TableHead>Stock</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {products.map((product) => {
              const category = categories.find((c) => c.id === product.category_id)
              return (
                <TableRow key={product.id}>
                  <TableCell>
                    {product.main_image_url && (
                      <Image
                        src={product.main_image_url}
                        alt={product.name_en ?? product.name ?? ''}
                        width={40}
                        height={40}
                        className="object-cover rounded"
                      />
                    )}
                  </TableCell>
                  <TableCell className="font-medium">{product.name_en}</TableCell>
                  <TableCell>{category?.name || 'N/A'}</TableCell>
                  <TableCell>{(product.price ?? 0).toFixed(2)} EGP</TableCell>
                  <TableCell>{product.quantity_in_stock}</TableCell>
                  <TableCell>
                    <span
                      className={`px-2 py-1 rounded text-sm ${product.is_available
                          ? 'bg-green-100 text-green-800'
                          : 'bg-red-100 text-red-800'
                        }`}
                    >
                      {product.is_available ? 'Available' : 'Unavailable'}
                    </span>
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-2">
                      <Link href={`/${locale}/dashboard/products/${product.id}/edit`}>
                        <Button
                          size="sm"
                          variant="outline"
                        >
                          Edit
                        </Button>
                      </Link>
                      {canDelete ? (
                        <Button
                          size="sm"
                          variant="destructive"
                          onClick={() => onDelete?.(product.id)}
                        >
                          Delete
                        </Button>
                      ) : (
                        <Button
                          size="sm"
                          variant="destructive"
                          disabled
                          title="You don't have permission to delete products"
                        >
                          Delete
                        </Button>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </div>

      {/* Pagination - Only show if products exist */}
      {products.length > 0 && (
        <Pagination>
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                onClick={() =>
                  handlePageChange(Math.max(1, filters.page - 1))
                }
                aria-disabled={filters.page === 1}
                className={filters.page === 1 ? 'pointer-events-none opacity-50' : undefined}
              />
            </PaginationItem>

            {Array.from({ length: 5 }).map((_, i) => {
              const page = i + 1
              return (
                <PaginationItem key={page}>
                  <PaginationLink
                    onClick={() => handlePageChange(page)}
                    isActive={filters.page === page}
                  >
                    {page}
                  </PaginationLink>
                </PaginationItem>
              )
            })}

            <PaginationItem>
              <PaginationNext
                onClick={() => handlePageChange(filters.page + 1)}
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      )}

      {/* Empty state message */}
      {products.length === 0 && !loading && (
        <div className="text-center py-12 text-gray-500 dark:text-gray-400">
          <p>No products found. Create your first product to get started.</p>
        </div>
      )}
    </div>
  )
}