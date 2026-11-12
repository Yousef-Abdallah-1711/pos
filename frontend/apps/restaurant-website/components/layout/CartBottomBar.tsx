'use client'

import { usePathname, useRouter } from 'next/navigation'
import { motion, AnimatePresence } from 'framer-motion'
import { useCartStore } from '@/lib/store/cart-store'
import { Button } from '@/components/ui/button'

function isHiddenRoute(pathname: string) {
  return pathname.includes('/cart') || pathname.includes('/checkout')
}

export function CartBottomBar() {
  const pathname = usePathname() || ''
  const router = useRouter()
  const itemCount = useCartStore((state) => state.getTotalItems())
  const total = useCartStore((state) => state.getTotal())
  const locale = pathname.startsWith('/ar') ? 'ar' : 'en'

  const shouldShow = itemCount > 0 && !isHiddenRoute(pathname)

  return (
    <AnimatePresence>
      {shouldShow ? (
        <motion.div
          initial={{ y: 80, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          exit={{ y: 80, opacity: 0 }}
          className="fixed inset-x-0 bottom-[calc(4rem+env(safe-area-inset-bottom))] z-50 px-3 md:hidden"
        >
          <div className="mx-auto flex max-w-xl items-center justify-between rounded-xl border border-theme-border/50 bg-[hsl(var(--theme-surface-elevated))] p-3 shadow-lg">
            <div>
              <p className="text-xs text-[hsl(var(--theme-text-muted))]">{itemCount} items</p>
              <p className="text-sm font-semibold">{total.toFixed(2)} EGP</p>
            </div>
            <Button
              onClick={() => router.push(`/${locale}/checkout`)}
              className="bg-[hsl(var(--theme-primary))] text-white hover:opacity-90"
            >
              Checkout ?
            </Button>
          </div>
        </motion.div>
      ) : null}
    </AnimatePresence>
  )
}
