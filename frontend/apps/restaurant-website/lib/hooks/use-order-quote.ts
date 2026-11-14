'use client'
import { useQuery } from '@tanstack/react-query'
import { useSubdomain } from '../subdomain-context'
import { fetchOrderQuote, type QuoteCartItem } from '../order-quote'

export function useOrderQuote(items: readonly QuoteCartItem[]) {
  const { slug, loading } = useSubdomain()
  const quantities = items.map(({ productId, quantity }) => ({ productId, quantity }))
  return useQuery({
    queryKey: ['order-quote', slug, quantities],
    queryFn: () => fetchOrderQuote(quantities),
    enabled: !loading && !!slug && quantities.length > 0,
    staleTime: 0,
    retry: 1,
    retryDelay: 1000,
  })
}
