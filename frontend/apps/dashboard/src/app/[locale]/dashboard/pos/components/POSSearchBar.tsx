'use client'

import { useEffect, useRef, useState } from 'react'
import { Search, X } from 'lucide-react'
import { selectPOSSearchQuery, usePOSStore } from '@/stores/posStore'

interface POSSearchBarProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
}

const DEBOUNCE_MS = 50

export function POSSearchBar({ locale, t }: POSSearchBarProps) {
  const searchQuery = usePOSStore(selectPOSSearchQuery)
  const setSearchQuery = usePOSStore((state) => state.setSearchQuery)
  const [draftQuery, setDraftQuery] = useState(searchQuery)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    setDraftQuery(searchQuery)
  }, [searchQuery])

  useEffect(() => {
    const timer = window.setTimeout(() => {
      setSearchQuery(draftQuery)
    }, DEBOUNCE_MS)

    return () => {
      window.clearTimeout(timer)
    }
  }, [draftQuery, setSearchQuery])

  useEffect(() => {
    const handleFocusSearch = () => {
      inputRef.current?.focus()
      inputRef.current?.select()
    }

    window.addEventListener('pos:focus-search', handleFocusSearch)
    return () => {
      window.removeEventListener('pos:focus-search', handleFocusSearch)
    }
  }, [])

  const isRTL = locale === 'ar'

  return (
    <div className="space-y-1">
      <div className="relative">
        <Search
          className={`pointer-events-none absolute top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400 ${
          isRTL ? 'right-3' : 'left-3'
          }`}
        />
        <input
          id="pos-product-search"
          ref={inputRef}
          type="search"
          value={draftQuery}
          onChange={(event) => setDraftQuery(event.target.value)}
          placeholder={t('pos.searchProducts', 'Search products...')}
          className={`h-10 w-full rounded-lg border border-gray-300 bg-white px-10 text-sm text-gray-900 outline-none transition placeholder:text-gray-400 focus:border-primary focus:ring-2 focus:ring-primary/30 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-100 ${
            isRTL ? 'text-right' : 'text-left'
          }`}
        />
        {draftQuery && (
          <button
            type="button"
            onClick={() => setDraftQuery('')}
            className={`absolute top-1/2 h-6 w-6 -translate-y-1/2 rounded-full text-gray-500 transition hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-gray-800 dark:hover:text-gray-200 ${
              isRTL ? 'left-2' : 'right-2'
            }`}
          >
            <X className="mx-auto h-4 w-4" />
          </button>
        )}
      </div>
      <p className={`text-xs text-gray-500 dark:text-gray-400 ${isRTL ? 'text-right' : 'text-left'}`}>
        {t('pos.searchShortcut', 'Shortcut: Ctrl+F or /')}
      </p>
    </div>
  )
}

export default POSSearchBar
