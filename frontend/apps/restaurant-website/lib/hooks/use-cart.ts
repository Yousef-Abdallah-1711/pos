'use client'

import { useEffect } from 'react'
import { useSubdomain } from '@/lib/subdomain-context'
import { useCartStore } from '@/lib/store/cart-store'

export function useCart() {
  const { slug } = useSubdomain()
  const setRestaurantSlug = useCartStore((state) => state.setRestaurantSlug)
  const items = useCartStore((state) => state.items)
  const totalItems = useCartStore((state) => state.totalItems)
  const totalPrice = useCartStore((state) => state.totalPrice)
  const addItem = useCartStore((state) => state.addItem)
  const removeItem = useCartStore((state) => state.removeItem)
  const updateQty = useCartStore((state) => state.updateQuantity)
  const updateNotes = useCartStore((state) => state.updateNotes)
  const clearCart = useCartStore((state) => state.clearCart)
  const tableId = useCartStore((state) => state.tableId)
  const tableName = useCartStore((state) => state.tableName)
  const setTableContext = useCartStore((state) => state.setTableContext)
  const clearTableContext = useCartStore((state) => state.clearTableContext)

  useEffect(() => {
    setRestaurantSlug(slug || 'default')
  }, [setRestaurantSlug, slug])

  return {
    items,
    totalItems,
    totalPrice,
    addItem,
    removeItem,
    updateQty,
    updateNotes,
    clearCart,
    tableId,
    tableName,
    setTableContext,
    clearTableContext,
  }
}

