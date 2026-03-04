'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Home, Menu, Receipt, Settings, ShoppingCart } from 'lucide-react'
import { cn } from '@/lib/utils'
import { useCartStore } from '@/lib/store/cart-store'

const items = [
  { key: 'home', href: '', icon: Home, label: 'Home' },
  { key: 'menu', href: '/menu', icon: Menu, label: 'Menu' },
  { key: 'cart', href: '/cart', icon: ShoppingCart, label: 'Cart' },
  { key: 'orders', href: '/orders', icon: Receipt, label: 'Orders' },
  { key: 'settings', href: '/settings', icon: Settings, label: 'Settings' },
] as const

export function MobileBottomNav() {
  const pathname = usePathname()
  const locale = pathname?.startsWith('/ar') ? 'ar' : 'en'
  const count = useCartStore((state) => state.getTotalItems())

  return (
    <nav className="fixed inset-x-0 bottom-0 z-50 border-t border-theme-border/40 bg-[hsl(var(--theme-surface-elevated))] pb-[env(safe-area-inset-bottom)] md:hidden">
      <ul className="mx-auto grid h-16 max-w-xl grid-cols-5">
        {items.map((item) => {
          const href = `/${locale}${item.href}`
          const active = pathname === href || (item.href === '' && pathname === `/${locale}`)
          const Icon = item.icon
          return (
            <li key={item.key}>
              <Link
                href={href}
                className={cn(
                  'relative flex h-full flex-col items-center justify-center gap-0.5 text-xs',
                  active ? 'text-[hsl(var(--theme-primary))]' : 'text-[hsl(var(--theme-text-muted))]'
                )}
              >
                <div className="relative">
                  <Icon className="h-4 w-4" />
                  {item.key === 'cart' && count > 0 ? (
                    <span className="absolute -right-2 -top-2 rounded-full bg-[hsl(var(--theme-error))] px-1 text-[10px] text-white">
                      {count}
                    </span>
                  ) : null}
                </div>
                <span>{item.label}</span>
                {active ? <span className="absolute bottom-1 h-1 w-1 rounded-full bg-[hsl(var(--theme-primary))]" /> : null}
              </Link>
            </li>
          )
        })}
      </ul>
    </nav>
  )
}
