'use client'

import { useEffect } from 'react'
import { ShoppingCart } from 'lucide-react'
import { usePathname, useRouter } from 'next/navigation'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { useCartStore } from '@/lib/store/cart-store'
import { useSubdomain } from '@/lib/subdomain-context'

interface CartButtonProps {
  onOpenCart?: () => void
}

export function CartButton({ onOpenCart }: CartButtonProps) {
  const router = useRouter()
  const pathname = usePathname()
  const itemCount = useCartStore((state) => state.totalItems)
  const setRestaurantSlug = useCartStore((state) => state.setRestaurantSlug)
  const { slug } = useSubdomain()
  const locale = pathname?.startsWith('/ar') ? 'ar' : 'en'

  useEffect(() => {
    setRestaurantSlug(slug || 'default')
  }, [setRestaurantSlug, slug])

  const handleClick = () => {
    if (onOpenCart) {
      onOpenCart()
      return
    }
    router.push(`/${locale}/cart`)
  }

  return (
    <Button
      type="button"
      variant="ghost"
      size="icon"
      onClick={handleClick}
      className="relative h-12 w-12 rounded-full"
      aria-label="Cart"
    >
      <ShoppingCart className="h-5 w-5" />
      {itemCount > 0 ? (
        <Badge className="absolute -end-1 -top-1 min-w-5 justify-center rounded-full bg-theme-error px-1 py-0 text-[10px] text-white">
          {itemCount}
        </Badge>
      ) : null}
    </Button>
  )
}
