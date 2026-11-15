import {
  BarChart3,
  ChefHat,
  Package,
  ShoppingCart,
  Store,
  UsersRound,
  UtensilsCrossed,
  type LucideIcon,
} from 'lucide-react'

export const navigationIcons = {
  pos: ShoppingCart,
  kitchen: ChefHat,
  orders: UtensilsCrossed,
  inventory: Package,
  staff: UsersRound,
  operations: BarChart3,
  storefront: Store,
} satisfies Record<string, LucideIcon>
