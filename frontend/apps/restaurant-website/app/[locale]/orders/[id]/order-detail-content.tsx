'use client';

import { usePathname } from 'next/navigation';
import { OrderStatusTracker } from '@/components/order-status-tracker';
import { Button } from '@pos-saas/ui';
import { Package, AlertCircle, Loader } from 'lucide-react';
import Link from 'next/link';
import { getLocaleFromPath, createTranslator } from '@/lib/translations';
import { useOrder, useCancelOrder } from '@/lib/hooks/use-api-queries';
import { useState } from 'react';

export default function OrderDetailContent() {
  const pathname = usePathname();
  const locale = getLocaleFromPath(pathname);
  const isRTL = locale === 'ar';
  const t = createTranslator(locale);

  // Extract id from pathname: /[locale]/orders/[id]
  const pathParts = pathname.split('/').filter(Boolean);
  const idFromPath = pathParts[pathParts.length - 1] || '0';

  const orderId = parseInt(idFromPath);
  const { data: order, isLoading, error } = useOrder(orderId);
  const { mutate: cancelOrder, isPending: isCancelling } = useCancelOrder();
  const [cancelConfirm, setCancelConfirm] = useState(false);
  const [cancelError, setCancelError] = useState<string | null>(null);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <Loader className="h-12 w-12 animate-spin text-primary mx-auto mb-4" />
          <p className="text-gray-600 dark:text-gray-400">{t('common.loading')}</p>
        </div>
      </div>
    );
  }

  if (error || !order) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className={`text-center max-w-md ${isRTL ? 'text-right' : ''}`}>
          <AlertCircle className="h-16 w-16 text-red-400 mx-auto mb-4" />
          <h1 className="text-2xl font-bold mb-2">{t('errors.orderNotFound')}</h1>
          <Link href={`/${locale}/orders`}>
            <Button className="mt-4">{t('orders.title')}</Button>
          </Link>
        </div>
      </div>
    );
  }

  const orderData = order as any;

  const handleCancelOrder = () => {
    cancelOrder(orderId, {
      onSuccess: () => {
        setCancelConfirm(false);
      },
      onError: (err: any) => {
        setCancelError(err.response?.data?.message || t('errors.somethingWentWrong'));
      },
    });
  };

  const canCancelOrder = ['pending', 'confirmed', 'preparing'].includes(orderData.status);

  // Create status array for OrderStatusTracker
  const allStatuses = [
    { status: 'pending' as const },
    { status: 'confirmed' as const },
    { status: 'preparing' as const },
    { status: 'ready' as const },
    { status: 'out_for_delivery' as const },
    { status: 'delivered' as const },
  ];

  return (
    <>
      {/* Page Header */}
      <div className="bg-gradient-to-r from-primary to-secondary text-white py-12">
        <div className={`max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 ${isRTL ? 'text-right' : ''}`}>
          <h1 className="text-4xl font-bold mb-2">{t('orders.title')}</h1>
          <p className="text-gray-100">
            {t('common.orderNumber')}: {orderData.orderNumber}
          </p>
        </div>
      </div>

      {/* Main Content */}
      <div className={`max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 ${isRTL ? 'text-right' : ''}`}>
        <div className="space-y-8">
          {/* Error Message */}
          {cancelError && (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex items-start gap-3">
              <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
              <p className="text-red-800 dark:text-red-200">{cancelError}</p>
            </div>
          )}

          {/* Order Status */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-8">
            <h2 className="text-2xl font-bold mb-8">{t('orderTracking.status')}</h2>
            <OrderStatusTracker statuses={allStatuses} currentStatus={orderData.status} locale={locale} />
          </div>

          {/* Order Information */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Basic Info */}
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <h3 className="text-lg font-semibold mb-4">{t('checkout.orderSummary')}</h3>
              <div className={`space-y-3 ${isRTL ? 'text-right' : ''}`}>
                <div className={`flex justify-between ${isRTL ? 'flex-row-reverse' : ''}`}>
                  <span className="text-gray-600 dark:text-gray-400">{t('common.orderNumber')}</span>
                  <span className="font-semibold">{orderData.orderNumber}</span>
                </div>
                <div className={`flex justify-between ${isRTL ? 'flex-row-reverse' : ''}`}>
                  <span className="text-gray-600 dark:text-gray-400">{t('common.total')}</span>
                  <span className="font-semibold text-lg">${orderData.totalAmount.toFixed(2)}</span>
                </div>
                <div className={`flex justify-between ${isRTL ? 'flex-row-reverse' : ''}`}>
                  <span className="text-gray-600 dark:text-gray-400">{t('orders.orderDate')}</span>
                  <span className="font-semibold">
                    {new Date(orderData.createdAt).toLocaleDateString(locale === 'ar' ? 'ar-EG' : 'en-US')}
                  </span>
                </div>
                {orderData.estimatedDelivery && (
                  <div className={`flex justify-between ${isRTL ? 'flex-row-reverse' : ''}`}>
                    <span className="text-gray-600 dark:text-gray-400">{t('orders.estimatedDelivery')}</span>
                    <span className="font-semibold">
                      {new Date(orderData.estimatedDelivery).toLocaleDateString(
                        locale === 'ar' ? 'ar-EG' : 'en-US'
                      )}
                    </span>
                  </div>
                )}
              </div>
            </div>

            {/* Customer Info */}
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <h3 className="text-lg font-semibold mb-4">{t('checkout.customerInformation')}</h3>
              <div className={`space-y-3 ${isRTL ? 'text-right' : ''}`}>
                {orderData.customerName && (
                  <div>
                    <p className="text-sm text-gray-600 dark:text-gray-400">{t('checkout.fullName')}</p>
                    <p className="font-semibold">{orderData.customerName}</p>
                  </div>
                )}
                {orderData.customerEmail && (
                  <div>
                    <p className="text-sm text-gray-600 dark:text-gray-400">{t('checkout.email')}</p>
                    <p className="font-semibold">{orderData.customerEmail}</p>
                  </div>
                )}
                {orderData.customerPhone && (
                  <div>
                    <p className="text-sm text-gray-600 dark:text-gray-400">{t('checkout.phone')}</p>
                    <p className="font-semibold">{orderData.customerPhone}</p>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Delivery Address */}
          {orderData.deliveryAddress && (
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <h3 className="text-lg font-semibold mb-4">{t('checkout.deliveryAddress')}</h3>
              <p className={isRTL ? 'text-right' : ''}>{orderData.deliveryAddress}</p>
            </div>
          )}

          {/* Order Items */}
          {orderData.items && orderData.items.length > 0 && (
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <h3 className="text-lg font-semibold mb-4">{t('orderTracking.orderDetails')}</h3>
              <div className="space-y-3">
                {orderData.items.map((item, idx) => (
                  <div
                    key={idx}
                    className={`flex justify-between items-center pb-3 border-b border-gray-200 dark:border-gray-700 last:border-0 last:pb-0 ${isRTL ? 'flex-row-reverse text-right' : ''}`}
                  >
                    <div>
                      <p className="font-medium">{item.name}</p>
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        {t('common.quantity')}: {item.quantity}
                      </p>
                      {item.notes && (
                        <p className="text-sm text-gray-500 italic">Note: {item.notes}</p>
                      )}
                    </div>
                    <p className="font-semibold">${(item.price * item.quantity).toFixed(2)}</p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Actions */}
          <div className={`flex flex-col sm:flex-row gap-4 ${isRTL ? 'sm:flex-row-reverse' : ''}`}>
            <Link href={`/${locale}/orders`} className="flex-1 sm:flex-none">
              <Button variant="outline" className="w-full sm:w-auto">
                {t('orders.title')}
              </Button>
            </Link>

            <Link href={`/${locale}/orders/track/${orderData.orderNumber}`} className="flex-1 sm:flex-none">
              <Button variant="outline" className="w-full sm:w-auto">
                {t('orders.trackOrder')}
              </Button>
            </Link>

            {canCancelOrder && (
              <div className="flex-1 sm:flex-none">
                {cancelConfirm ? (
                  <div className={`flex gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
                    <Button
                      variant="destructive"
                      onClick={handleCancelOrder}
                      disabled={isCancelling}
                      className="flex-1 sm:flex-none"
                    >
                      {isCancelling ? t('common.loading') : t('common.delete')}
                    </Button>
                    <Button
                      variant="outline"
                      onClick={() => setCancelConfirm(false)}
                      className="flex-1 sm:flex-none"
                    >
                      {t('common.cancel')}
                    </Button>
                  </div>
                ) : (
                  <Button
                    variant="destructive"
                    onClick={() => setCancelConfirm(true)}
                    className="w-full sm:w-auto"
                  >
                    {t('common.delete')}
                  </Button>
                )}
              </div>
            )}

            <Link href={`/${locale}/menu`} className="flex-1 sm:flex-none">
              <Button className="w-full sm:w-auto">{t('navigation.menu')}</Button>
            </Link>
          </div>
        </div>
      </div>
    </>
  );
}
