'use client';

import { useState } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@pos-saas/ui';
import { AlertCircle, CheckCircle, Loader } from 'lucide-react';
import Link from 'next/link';
import { toast } from 'sonner';
import { getLocaleFromPath, createTranslator } from '@/lib/translations';
import { useCartStore } from '@/lib/store/cart-store';
import { useCreateOrder } from '@/lib/hooks/use-api-queries';
import type { ServerThemeData } from '@/lib/api/get-theme-server';

// Validation schema
const checkoutSchema = z.object({
  customerName: z.string().min(2, 'Name is required'),
  customerEmail: z.string().email('Invalid email'),
  customerPhone: z.string().min(10, 'Phone is required'),
  deliveryAddress: z.string().min(5, 'Address is required'),
  specialInstructions: z.string().optional(),
  paymentMethod: z.literal('cash'),
});

type CheckoutFormData = z.infer<typeof checkoutSchema>;

interface CheckoutPageContentProps {
  locale: 'en' | 'ar';
  themeData: ServerThemeData;
}

export function CheckoutPageContent({ locale, themeData }: CheckoutPageContentProps) {
  const pathname = usePathname();
  const router = useRouter();
  const localeFromPath = getLocaleFromPath(pathname);
  const isRTL = locale === 'ar';
  const t = createTranslator(locale);

  // Extract theme colors with fallbacks
  const primaryColor = themeData?.colors?.primary || '#f97316';
  const secondaryColor = themeData?.colors?.secondary || '#0ea5e9';

  const cartItems = useCartStore((state) => state.items);
  const cartTotal = useCartStore((state) => state.getTotal());
  const clearCart = useCartStore((state) => state.clearCart);

  const { mutate: createOrder, isPending: isLoading } = useCreateOrder();
  const [error, setError] = useState<string | null>(null);
  const [successOrderNumber, setSuccessOrderNumber] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CheckoutFormData>({
    resolver: zodResolver(checkoutSchema),
    defaultValues: {
      paymentMethod: 'cash',
    },
  });

  const deliveryFee = cartTotal > 50 ? 0 : 5.99;
  const grandTotal = cartTotal + deliveryFee;

  if (successOrderNumber) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-green-50 to-green-100 dark:from-green-900/20 dark:to-green-800/20">
        <div className={`max-w-md w-full mx-4 bg-white dark:bg-gray-800 rounded-lg shadow-xl p-8 text-center ${isRTL ? 'text-right' : ''}`}>
          <CheckCircle className="h-16 w-16 text-green-600 mx-auto mb-4" />
          <h1 className="text-3xl font-bold text-green-600 mb-2">{t('checkout.success')}</h1>
          <p className="text-gray-600 dark:text-gray-400 mb-4">{t('checkout.orderNumber')}:</p>
          <p className="text-2xl font-bold mb-6" style={{ color: primaryColor }}>
            {successOrderNumber}
          </p>
          <div className={`flex gap-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
            <Link href={`/${locale}/menu`} className="flex-1">
              <Button className="w-full">{t('common.menu')}</Button>
            </Link>
            <Link href={`/${locale}/orders/track/${successOrderNumber}`} className="flex-1">
              <Button variant="outline" className="w-full">
                {t('checkout.trackOrder')}
              </Button>
            </Link>
          </div>
        </div>
      </div>
    );
  }

  if (cartItems.length === 0 && !successOrderNumber) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className={`text-center ${isRTL ? 'text-right' : ''}`}>
          <AlertCircle className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h1 className="text-2xl font-bold mb-2">{t('cart.empty')}</h1>
          <Link href={`/${locale}/menu`}>
            <Button className="mt-4">{t('menu.readyToCheckout')}</Button>
          </Link>
        </div>
      </div>
    );
  }

  const onSubmit = (data: CheckoutFormData) => {
    setError(null);

    const insufficientItems = cartItems.filter(
      (item) =>
        (item.trackInventory ?? false) &&
        (item.quantityInStock ?? 0) > -1 &&
        (item.quantityInStock ?? 0) < item.quantity
    );
    if (insufficientItems.length > 0) {
      const inventoryMessage = t('errors.insufficientInventory') || 'Insufficient inventory for one or more items.';
      toast.error(inventoryMessage);
      setError(inventoryMessage);
      return;
    }

    // Transform cart items to backend format (snake_case)
    const orderItems = cartItems.map((item) => ({
      product_id: item.productId,
      quantity: item.quantity,
      special_instructions: item.notes || '',
    }));

    // Debug payload
    console.log('[Checkout] Submitting order items:', orderItems);

    createOrder(
      {
        customer_name: data.customerName,
        customer_email: data.customerEmail,
        customer_phone: data.customerPhone,
        delivery_address: data.deliveryAddress,
        payment_method: data.paymentMethod,
        notes: data.specialInstructions || '',
        items: orderItems,
      } as any,
      {
        onSuccess: (response) => {
          const orderNumber =
            response?.data?.orderNumber ||
            response?.data?.order_number ||
            response?.orderNumber ||
            response?.order_number;
          if (orderNumber) {
            setSuccessOrderNumber(orderNumber);
          }
          clearCart();
          reset();
        },
        onError: (err: any) => {
          const rawMessage = err.response?.data?.error || err.response?.data?.message;
          const message = typeof rawMessage === 'string' && rawMessage.trim().length > 0
            ? rawMessage
            : t('errors.somethingWentWrong');
          const normalized = message.toLowerCase();

          if (normalized.includes('insufficient inventory')) {
            const inventoryMessage = t('errors.insufficientInventory');
            toast.error(inventoryMessage);
            setError(inventoryMessage);
            return;
          }

          setError(message);
        },
      }
    );
  };

  return (
    <>
      {/* Page Header */}
      <div
        className="text-white py-12"
        style={{
          background: `linear-gradient(135deg, ${primaryColor} 0%, ${secondaryColor} 100%)`,
        }}
      >
        <div className={`max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 ${isRTL ? 'text-right' : ''}`}>
          <h1 className="text-4xl font-bold mb-2">{t('checkout.title')}</h1>
        </div>
      </div>

      {/* Main Content */}
      <div className={`max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 ${isRTL ? 'text-right' : ''}`}>
        <div className={`grid grid-cols-1 lg:grid-cols-3 gap-8 ${isRTL ? 'lg:grid-flow-col-dense' : ''}`}>
          {/* Form Section */}
          <div className={`lg:col-span-2 ${isRTL ? 'lg:order-2' : ''}`}>
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              {error && (
                <div className="mb-6 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex items-start gap-3">
                  <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
                  <p className="text-red-800 dark:text-red-200">{error}</p>
                </div>
              )}

              <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
                {/* Customer Information */}
                <div>
                  <h2 className="text-2xl font-bold mb-4">{t('checkout.customerInformation')}</h2>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium mb-1">{t('checkout.fullName')}</label>
                      <input
                        {...register('customerName')}
                        type="text"
                        className={`w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 dark:bg-gray-700 dark:text-white ${isRTL ? 'text-right' : ''} ${errors.customerName ? 'border-red-500' : ''}`}
                        style={{ borderColor: errors.customerName ? undefined : '#e5e7eb' }}
                        placeholder={t('checkout.fullName')}
                      />
                      {errors.customerName && (
                        <p className="text-red-500 text-sm mt-1">{errors.customerName.message}</p>
                      )}
                    </div>

                    <div>
                      <label className="block text-sm font-medium mb-1">{t('checkout.email')}</label>
                      <input
                        {...register('customerEmail')}
                        type="email"
                        className={`w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 dark:bg-gray-700 dark:text-white ${isRTL ? 'text-right' : ''} ${errors.customerEmail ? 'border-red-500' : ''}`}
                        style={{}}
                        placeholder={t('checkout.email')}
                      />
                      {errors.customerEmail && (
                        <p className="text-red-500 text-sm mt-1">{errors.customerEmail.message}</p>
                      )}
                    </div>

                    <div>
                      <label className="block text-sm font-medium mb-1">{t('checkout.phone')}</label>
                      <input
                        {...register('customerPhone')}
                        type="tel"
                        className={`w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 dark:bg-gray-700 dark:text-white ${isRTL ? 'text-right' : ''} ${errors.customerPhone ? 'border-red-500' : ''}`}
                        style={{}}
                        placeholder={t('checkout.phone')}
                      />
                      {errors.customerPhone && (
                        <p className="text-red-500 text-sm mt-1">{errors.customerPhone.message}</p>
                      )}
                    </div>
                  </div>
                </div>

                {/* Delivery Information */}
                <div>
                  <h2 className="text-2xl font-bold mb-4">{t('checkout.deliveryInformation')}</h2>
                  <div className="space-y-4">
                    <div>
                      <label className="block text-sm font-medium mb-1">{t('checkout.deliveryAddress')}</label>
                      <textarea
                        {...register('deliveryAddress')}
                        className={`w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 dark:bg-gray-700 dark:text-white ${isRTL ? 'text-right' : ''} ${errors.deliveryAddress ? 'border-red-500' : ''}`}
                        style={{}}
                        rows={3}
                        placeholder={t('checkout.deliveryAddress')}
                      />
                      {errors.deliveryAddress && (
                        <p className="text-red-500 text-sm mt-1">{errors.deliveryAddress.message}</p>
                      )}
                    </div>

                    <div>
                      <label className="block text-sm font-medium mb-1">{t('checkout.specialInstructions')}</label>
                      <textarea
                        {...register('specialInstructions')}
                        className={`w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 dark:bg-gray-700 dark:text-white ${isRTL ? 'text-right' : ''}`}
                        style={{}}
                        rows={2}
                        placeholder={t('checkout.specialInstructions')}
                      />
                    </div>
                  </div>
                </div>

                {/* Payment Method */}
                <div>
                  <h2 className="text-2xl font-bold mb-4">{t('checkout.paymentMethod')}</h2>
                  <div
                    className="p-4 border-2 rounded-lg bg-green-50 dark:bg-green-900/20"
                    style={{ borderColor: primaryColor }}
                  >
                    <input
                      {...register('paymentMethod')}
                      type="hidden"
                      value="cash"
                    />
                    <span className="font-semibold text-green-700 dark:text-green-400">
                      💵 {t('checkout.cash')} (Cash on Delivery)
                    </span>
                    <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                      Pay when your order arrives
                    </p>
                  </div>
                </div>

                {/* Submit Button */}
                <Button
                  type="submit"
                  disabled={isLoading}
                  className="w-full py-3 text-lg"
                  style={{ backgroundColor: primaryColor }}
                >
                  {isLoading ? (
                    <>
                      <Loader className="h-5 w-5 mr-2 animate-spin" />
                      {t('checkout.processing')}
                    </>
                  ) : (
                    t('checkout.placeOrder')
                  )}
                </Button>
              </form>
            </div>
          </div>

          {/* Order Summary (Sticky Sidebar) */}
          <div className={`lg:col-span-1 ${isRTL ? 'lg:order-1' : ''}`}>
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 lg:sticky lg:top-6">
              <h2 className="text-2xl font-bold mb-6">{t('checkout.orderSummary')}</h2>

              {/* Items List */}
              <div className="space-y-3 mb-6 max-h-96 overflow-y-auto">
                {cartItems.map((item) => (
                  <div key={item.id} className={`flex justify-between items-start pb-3 border-b border-gray-200 dark:border-gray-700 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
                    <div>
                      <p className="font-medium">{item.name}</p>
                      <p className="text-sm text-gray-600 dark:text-gray-400">x {item.quantity}</p>
                    </div>
                    <p className="font-semibold">${(item.total || item.price * item.quantity).toFixed(2)}</p>
                  </div>
                ))}
              </div>

              {/* Totals */}
              <div className="space-y-3 border-t border-gray-200 dark:border-gray-700 pt-4">
                <div className={`flex justify-between ${isRTL ? 'flex-row-reverse' : ''}`}>
                  <span className="text-gray-600 dark:text-gray-400">{t('common.subtotal')}</span>
                  <span className="font-semibold">${cartTotal.toFixed(2)}</span>
                </div>

                <div className={`flex justify-between ${isRTL ? 'flex-row-reverse' : ''}`}>
                  <span className="text-gray-600 dark:text-gray-400">{t('cart.deliveryFee')}</span>
                  <span className="font-semibold">
                    {deliveryFee === 0 ? (
                      <span className="text-green-600">{t('common.free')}</span>
                    ) : (
                      `$${deliveryFee.toFixed(2)}`
                    )}
                  </span>
                </div>

                <div className={`flex justify-between text-lg border-t border-gray-200 dark:border-gray-700 pt-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
                  <span className="font-bold">{t('common.grandTotal')}</span>
                  <span className="text-2xl font-bold" style={{ color: primaryColor }}>
                    ${grandTotal.toFixed(2)}
                  </span>
                </div>
              </div>

              {/* Back to Cart Link */}
              <Link href={`/${locale}/cart`} className="text-primary hover:underline text-sm mt-4 block" style={{ color: primaryColor }}>
                {t('cart.continueShopping')}
              </Link>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
