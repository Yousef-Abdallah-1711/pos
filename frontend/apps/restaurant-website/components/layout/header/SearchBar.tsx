'use client'

import { useEffect, useMemo, useState } from 'react'
import { Search } from 'lucide-react'
import { usePathname, useRouter } from 'next/navigation'
import { useQueryClient } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { useSubdomain } from '@/lib/subdomain-context'

interface SearchBarProps {
  maxResults?: number
}

type CachedProduct = {
  id: number
  name?: string
  name_en?: string
  category?: string
  category_id?: number
  price?: number
}

type CachedCategory = {
  id: number
  name?: string
}

export function SearchBar({ maxResults = 5 }: SearchBarProps) {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')
  const [debouncedQuery, setDebouncedQuery] = useState('')
  const router = useRouter()
  const pathname = usePathname()
  const locale = pathname?.startsWith('/ar') ? 'ar' : 'en'
  const { slug } = useSubdomain()
  const queryClient = useQueryClient()

  const products = useMemo(
    () => (queryClient.getQueryData(['products', slug || '']) as CachedProduct[] | undefined) || [],
    [queryClient, slug]
  )
  const categories = useMemo(
    () => (queryClient.getQueryData(['categories', slug || '']) as CachedCategory[] | undefined) || [],
    [queryClient, slug]
  )

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedQuery(query), 300)
    return () => clearTimeout(timer)
  }, [query])

  const lowered = debouncedQuery.trim().toLowerCase()
  const filteredProducts = useMemo(() => {
    if (!lowered) return []
    return products
      .filter((product) => {
        const name = String(product.name || product.name_en || '').toLowerCase()
        const category = String(product.category || '').toLowerCase()
        return name.includes(lowered) || category.includes(lowered)
      })
      .slice(0, maxResults)
  }, [lowered, products, maxResults])

  const filteredCategories = useMemo(() => {
    if (!lowered) return []
    return categories
      .filter((category) => String(category.name || '').toLowerCase().includes(lowered))
      .slice(0, maxResults)
  }, [lowered, categories, maxResults])

  const navigate = (target: string) => {
    setOpen(false)
    router.push(target)
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button type="button" variant="outline" aria-label={locale === 'ar' ? 'البحث عن المنتجات' : 'Search products'} className="h-12 min-w-12 gap-2 px-3 md:min-w-56">
          <Search className="h-4 w-4" />
          <span className="hidden md:inline">{query || (locale === 'ar' ? 'البحث عن المنتجات...' : 'Search products...')}</span>
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[22rem] p-0" align="start">
        <Command shouldFilter={false}>
          <CommandInput
            value={query}
            onValueChange={setQuery}
            placeholder={locale === 'ar' ? 'البحث عن المنتجات أو الفئات' : 'Search products or categories'}
          />
          <CommandList>
            <CommandEmpty>{locale === 'ar' ? 'لا توجد نتائج محفوظة. افتح القائمة أولاً.' : 'No cached results. Open Menu first.'}</CommandEmpty>
            {filteredProducts.length > 0 ? (
              <CommandGroup heading={locale === 'ar' ? 'المنتجات' : 'Products'}>
                {filteredProducts.map((product) => (
                  <CommandItem
                    key={`p-${product.id}`}
                    onSelect={() => navigate(`/${locale}/menu?search=${encodeURIComponent(String(product.name || product.name_en || ''))}`)}
                  >
                    <div className="flex w-full items-center justify-between gap-2">
                      <span className="truncate">{product.name || product.name_en || 'Product'}</span>
                      <span className="text-xs text-muted-foreground">{typeof product.price === 'number' ? `${product.price}` : ''}</span>
                    </div>
                  </CommandItem>
                ))}
              </CommandGroup>
            ) : null}
            {filteredCategories.length > 0 ? (
              <CommandGroup heading={locale === 'ar' ? 'الفئات' : 'Categories'}>
                {filteredCategories.map((category) => (
                  <CommandItem
                    key={`c-${category.id}`}
                    onSelect={() => navigate(`/${locale}/menu?category=${category.id}`)}
                  >
                    {category.name || 'Category'}
                  </CommandItem>
                ))}
              </CommandGroup>
            ) : null}
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
