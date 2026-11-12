'use client'

import Link from 'next/link'
import { ShoppingCart } from 'lucide-react'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { CartItem } from '@/components/commerce/CartItem'
import { formatPrice } from '@/lib/utils/format'
import { useCart } from '@/lib/hooks/use-cart'

interface CartSheetProps {
  locale: 'en' | 'ar'
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CartSheet({ locale, open, onOpenChange }: CartSheetProps) {
  const isArabic = locale === 'ar'
  const { items, totalItems, totalPrice, updateQty, removeItem } = useCart()
  const deliveryFee = totalPrice > 500 ? 0 : 45
  const grandTotal = totalPrice + deliveryFee

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side={isArabic ? 'left' : 'right'} className="flex w-full max-w-md flex-col p-0">
        <SheetHeader className="px-4 py-4">
          <SheetTitle className={isArabic ? 'text-right' : 'text-left'}>
            {isArabic ? `سلة التسوق (${totalItems})` : `Cart (${totalItems})`}
          </SheetTitle>
        </SheetHeader>
        <Separator />

        {items.length === 0 ? (
          <div className="flex flex-1 flex-col items-center justify-center px-6 text-center">
            <ShoppingCart className="mb-3 h-14 w-14 text-[hsl(var(--theme-text))]/35" />
            <h3 className="text-lg font-semibold text-[hsl(var(--theme-text))]">
              {isArabic ? 'سلتك فارغة' : 'Your cart is empty'}
            </h3>
            <p className="mt-2 text-sm text-[hsl(var(--theme-text))]/65">
              {isArabic ? 'ابدأ بإضافة منتجات من القائمة' : 'Start adding products from the menu'}
            </p>
            <Button asChild className="mt-5 bg-[hsl(var(--theme-primary))] text-white">
              <Link href={`/${locale}/menu`} onClick={() => onOpenChange(false)}>
                {isArabic ? 'تصفح القائمة' : 'Browse Menu'}
              </Link>
            </Button>
          </div>
        ) : (
          <>
            <ScrollArea className="flex-1 px-4">
              <div className="space-y-3 py-4">
                {items.map((item) => (
                  <CartItem
                    key={item.id}
                    item={item}
                    locale={locale}
                    onIncrease={() => updateQty(item.productId, item.quantity + 1)}
                    onDecrease={() => updateQty(item.productId, item.quantity - 1)}
                    onRemove={() => removeItem(item.productId)}
                  />
                ))}
              </div>
            </ScrollArea>
            <Separator />
            <div className="space-y-3 px-4 py-4">
              <Textarea
                placeholder={isArabic ? 'ملاحظات الطلب (اختياري)' : 'Order notes (optional)'}
                className="min-h-20"
              />
              <div className="space-y-1 text-sm">
                <div className={`flex items-center justify-between ${isArabic ? 'flex-row-reverse' : ''}`}>
                  <span className="text-[hsl(var(--theme-text))]/70">{isArabic ? 'المجموع' : 'Subtotal'}</span>
                  <span className="font-medium">{formatPrice(totalPrice, locale)}</span>
                </div>
                <div className={`flex items-center justify-between ${isArabic ? 'flex-row-reverse' : ''}`}>
                  <span className="text-[hsl(var(--theme-text))]/70">{isArabic ? 'التوصيل' : 'Delivery'}</span>
                  <span className="font-medium">{deliveryFee === 0 ? (isArabic ? 'مجاني' : 'Free') : formatPrice(deliveryFee, locale)}</span>
                </div>
                <div className={`flex items-center justify-between text-base font-bold ${isArabic ? 'flex-row-reverse' : ''}`}>
                  <span>{isArabic ? 'الإجمالي' : 'Total'}</span>
                  <span className="text-[hsl(var(--theme-primary))]">{formatPrice(grandTotal, locale)}</span>
                </div>
              </div>
              <Button asChild className="h-11 w-full bg-[linear-gradient(130deg,hsl(var(--theme-primary)),hsl(var(--theme-accent)))] text-white">
                <Link href={`/${locale}/checkout`} onClick={() => onOpenChange(false)}>
                  {isArabic ? 'إتمام الطلب' : 'Checkout'}
                </Link>
              </Button>
              <Button variant="ghost" asChild className="h-10 w-full">
                <Link href={`/${locale}/menu`} onClick={() => onOpenChange(false)}>
                  {isArabic ? 'متابعة التسوق' : 'Continue Shopping'}
                </Link>
              </Button>
            </div>
          </>
        )}
      </SheetContent>
    </Sheet>
  )
}

