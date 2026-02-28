'use client';

import { useEffect, useState } from 'react';
import { usePathname } from 'next/navigation';
import { OrderStatusTracker, type OrderStatus } from '@/components/commerce/order-status-tracker';
import { Button } from '@pos-saas/ui';
import { Search, Package, MapPin, Clock, AlertCircle } from 'lucide-react';
import Link from 'next/link';
import { getLocaleFromPath, createTranslator } from '@/lib/translations';
import { apiClient } from '@/lib/api-client';

type StatusHistoryEntry = {
  new_status?: string;
  status?: string;
  created_at?: string;
  timestamp?: string;
};

type TrackedItem = {
  id?: number;
  name: string;
  quantity: number;
  price: number;
  notes?: string;
};

type TrackedOrder = {
  id?: number;
  orderNumber: string;
  status: OrderStatus['status'];
  totalAmount: number;
  createdAt: string;
  estimatedDelivery?: string;
  items: TrackedItem[];
  deliveryAddress?: string;
  customerName?: string;
  customerPhone?: string;
  paymentStatus?: string;
};

const TIMELINE: OrderStatus['status'][] = [
  'pending',
  'confirmed',
  'preparing',
  'ready',
  'out_for_delivery',
  'delivered',
];

const normalizeOrder = (raw: any): TrackedOrder => {
  const items = Array.isArray(raw?.items) ? raw.items : [];
  const normalizedItems: TrackedItem[] = items.map((item: any) => ({
    id: item?.id,
    name: item?.name || item?.product_name || item?.productName || '',
    quantity: Number(item?.quantity ?? 0),
    price: Number(item?.price ?? item?.unit_price ?? item?.unitPrice ?? 0),
    notes: item?.notes || item?.special_instructions,
  }));

  return {
    id: raw?.id,
    orderNumber: raw?.order_number || raw?.orderNumber || '',
    status: raw?.status,
    totalAmount: Number(raw?.total_amount ?? raw?.totalAmount ?? raw?.total ?? 0),
    createdAt: raw?.created_at || raw?.createdAt || '',
    estimatedDelivery:
      raw?.estimated_delivery ||
      raw?.estimated_delivery_time ||
      raw?.estimatedDelivery ||
      raw?.estimatedDeliveryTime,
    items: normalizedItems,
    deliveryAddress: raw?.delivery_address || raw?.deliveryAddress,
    customerName: raw?.customer_name || raw?.customerName,
    customerPhone: raw?.customer_phone || raw?.customerPhone,
    paymentStatus: raw?.payment_status || raw?.paymentStatus,
  };
};

const buildTimeline = (history: StatusHistoryEntry[]): OrderStatus[] => {
  const historyMap = new Map<string, string>();
  history.forEach((entry) => {
    const status = entry.new_status || entry.status;
    const timestamp = entry.created_at || entry.timestamp;
    if (status && timestamp && !historyMap.has(status)) {
      historyMap.set(status, timestamp);
    }
  });

  return TIMELINE.map((status) => ({
    status,
    timestamp: historyMap.get(status),
  }));
};

interface TrackPageContentProps {
  initialOrderNumber?: string;
}

export default function TrackPageContent({ initialOrderNumber }: TrackPageContentProps) {
  const pathname = usePathname();
  const locale = getLocaleFromPath(pathname);
  const isRTL = locale === 'ar';
  const t = createTranslator(locale);

  const [orderNumber, setOrderNumber] = useState('');
  const [searchSubmitted, setSearchSubmitted] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [trackedOrder, setTrackedOrder] = useState<TrackedOrder | null>(null);
  const [statusTimeline, setStatusTimeline] = useState<OrderStatus[]>([]);

  const fetchOrder = async (numberValue: string) => {
    const trimmed = numberValue.trim();
    if (!trimmed) return;

    setIsLoading(true);
    setError(null);

    try {
      const orderResponse = await apiClient.get<any>(
        `/public/orders/number/${encodeURIComponent(trimmed)}/lookup`,
      );
      const orderData = orderResponse?.data ?? orderResponse;
      const normalizedOrder = normalizeOrder(orderData);

      if (!normalizedOrder.orderNumber) {
        throw new Error('order_not_found');
      }

      setTrackedOrder(normalizedOrder);

      if (orderData?.id) {
        const trackResponse = await apiClient.get<any>(`/public/orders/${orderData.id}/track`);
        const trackData = trackResponse?.data ?? trackResponse;
        const history = Array.isArray(trackData?.status_history) ? trackData.status_history : [];
        setStatusTimeline(buildTimeline(history));
      } else {
        setStatusTimeline(buildTimeline([]));
      }
    } catch (err: any) {
      const message =
        err?.response?.data?.message ||
        err?.response?.data?.error ||
        t('errors.orderNotFound');
      setError(message);
      setTrackedOrder(null);
      setStatusTimeline([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orderNumber.trim()) return;

    setSearchSubmitted(true);
    await fetchOrder(orderNumber);
  };

  const handleReset = () => {
    setOrderNumber('');
    setSearchSubmitted(false);
    setTrackedOrder(null);
    setStatusTimeline([]);
    setError(null);
  };

  useEffect(() => {
    if (!initialOrderNumber) return;
    setOrderNumber(initialOrderNumber);
    setSearchSubmitted(true);
    fetchOrder(initialOrderNumber);
  }, [initialOrderNumber]);

  const getEstimatedTime = () => {
    if (!trackedOrder?.estimatedDelivery) {
      return t('orderTracking.estimatedDelivery');
    }

    const now = new Date();
    const estimated = new Date(trackedOrder.estimatedDelivery);
    const diffMs = estimated.getTime() - now.getTime();
    const diffMins = Math.round(diffMs / 60000);

    if (diffMins <= 0) return t('orderTracking.arrivingSoon');
    if (diffMins < 60) return `${diffMins} ${locale === 'ar' ? 'min' : 'minutes'}`;
    const hours = Math.round(diffMins / 60);
    return `${hours} ${locale === 'ar' ? 'hour' : 'hour'}${hours > 1 && locale !== 'ar' ? 's' : ''}`;
  };

  return (
    <>
      {/* Page Header */}
      <div className="bg-gradient-to-r from-primary to-secondary text-white py-12">
        <div className={`max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 ${isRTL ? 'text-right' : ''}`}>
          <h1 className="text-4xl font-bold mb-2">{t('orderTracking.title')}</h1>
          <p className="text-gray-100">{t('orderTracking.subtitle')}</p>
        </div>
      </div>

      {/* Main Content */}
      <div className={`max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 ${isRTL ? 'text-right' : ''}`}>
        {/* Search Form */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-8">
          <form onSubmit={handleSearch}>
            <div className={`flex flex-col sm:flex-row gap-3 ${isRTL ? 'sm:flex-row-reverse' : ''}`}>
              <input
                type="text"
                placeholder={t('orderTracking.searchPlaceholder')}
                value={orderNumber}
                onChange={(e) => setOrderNumber(e.target.value)}
                className={`flex-1 px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:text-white ${isRTL ? 'text-right' : ''}`}
              />
              <Button type="submit" disabled={isLoading} className={`px-6 py-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
                <Search className="h-5 w-5 mr-2" />
                {isLoading ? t('orderTracking.searching') : t('orderTracking.search')}
              </Button>
            </div>
          </form>
        </div>

        {/* Results */}
        {!searchSubmitted ? (
          <div className="text-center py-12">
            <Package className="h-16 w-16 text-gray-400 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-600 dark:text-gray-300 mb-2">
              {t('orderTracking.enterOrderNumber')}
            </h3>
            <p className="text-gray-500 dark:text-gray-400">{t('orderTracking.findInEmail')}</p>
          </div>
        ) : error ? (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6 flex items-start gap-3">
            <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 mt-0.5 flex-shrink-0" />
            <div className={isRTL ? 'text-right' : ''}>
              <p className="font-semibold text-red-900 dark:text-red-200">{t('errors.orderNotFound')}</p>
              <p className="text-sm text-red-700 dark:text-red-300">{error}</p>
            </div>
          </div>
        ) : (
          trackedOrder && (
            <div className="space-y-8">
            {/* Order Status Tracker */}
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-8">
              <h2 className="text-2xl font-bold mb-8">{t('orderTracking.status')}</h2>
              <OrderStatusTracker
                statuses={statusTimeline.length > 0 ? statusTimeline : TIMELINE.map((status) => ({ status }))}
                currentStatus={trackedOrder.status}
                locale={locale}
              />
            </div>

            {/* Delivery Information */}
            {trackedOrder.status === 'out_for_delivery' && (
              <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-blue-900 dark:text-blue-100 mb-4">
                  {t('orderTracking.outForDelivery')}
                </h3>
                <div className="space-y-4">
                  <div className={`flex items-start gap-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
                    <Clock className="h-5 w-5 text-blue-600 dark:text-blue-400 mt-0.5 flex-shrink-0" />
                    <div className={isRTL ? 'text-right' : ''}>
                      <p className="text-sm text-blue-700 dark:text-blue-300">
                        {t('orderTracking.estimatedDelivery')}
                      </p>
                      <p className="text-lg font-semibold text-blue-900 dark:text-blue-100">
                        {getEstimatedTime()}
                      </p>
                    </div>
                  </div>
                  <div className={`flex items-start gap-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
                    <MapPin className="h-5 w-5 text-blue-600 dark:text-blue-400 mt-0.5 flex-shrink-0" />
                    <div className={isRTL ? 'text-right' : ''}>
                      <p className="text-sm text-blue-700 dark:text-blue-300">
                        {t('orderTracking.deliveryAddress')}
                      </p>
                      <p className="text-blue-900 dark:text-blue-100">{trackedOrder.deliveryAddress || '-'}</p>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Driver Information */}
            {trackedOrder.status === 'out_for_delivery' && (
              <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
                <h3 className="text-lg font-semibold mb-4">{t('orderTracking.driver')}</h3>
                <div className={`grid grid-cols-1 md:grid-cols-2 gap-4 ${isRTL ? 'text-right' : ''}`}>
                  <div>
                    <p className="text-sm text-gray-600 dark:text-gray-400">{t('orderTracking.driverName')}</p>
                    <p className="text-lg font-semibold">{trackedOrder.customerName || '-'}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600 dark:text-gray-400">{t('orderTracking.contactDriver')}</p>
                    {trackedOrder.customerPhone ? (
                      <a href={`tel:${trackedOrder.customerPhone}`} className="text-primary hover:underline font-semibold">
                        {trackedOrder.customerPhone}
                      </a>
                    ) : (
                      <span className="text-gray-500 dark:text-gray-400">-</span>
                    )}
                  </div>
                </div>
              </div>
            )}

            {/* Order Details */}
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <h3 className="text-lg font-semibold mb-4">{t('orderTracking.orderDetails')}</h3>
              <div className="space-y-3 mb-4">
                {trackedOrder.items.map((item) => (
                  <div
                    key={item.id}
                    className={`flex justify-between items-center pb-3 border-b border-gray-200 dark:border-gray-700 last:border-0 last:pb-0 ${isRTL ? 'flex-row-reverse text-right' : ''}`}
                  >
                    <div>
                      <p className="font-medium">{item.name}</p>
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        {t('common.quantity')}: {item.quantity}
                      </p>
                    </div>
                    <p className="font-semibold">${(item.price * item.quantity).toFixed(2)}</p>
                  </div>
                ))}
              </div>
              <div className={`border-t border-gray-200 dark:border-gray-700 pt-4 flex justify-between items-center ${isRTL ? 'flex-row-reverse' : ''}`}>
                <p className="text-lg font-bold">{t('common.total')}:</p>
                <p className="text-2xl font-bold text-primary">${trackedOrder.totalAmount.toFixed(2)}</p>
              </div>
            </div>

            {/* Restaurant Info */}
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <h3 className="text-lg font-semibold mb-4">{t('orderTracking.restaurantInfo')}</h3>
              <p className="text-lg font-medium mb-4">{t('common.restaurant')}</p>
              <div className={`flex gap-3 ${isRTL ? 'flex-row-reverse' : ''}`}>
                <Button variant="outline" size="sm">
                  {t('orderTracking.viewMenu')}
                </Button>
                <Button variant="outline" size="sm">
                  {t('orderTracking.contactRestaurant')}
                </Button>
              </div>
            </div>

            {/* Actions */}
            <div className={`flex flex-col sm:flex-row gap-3 ${isRTL ? 'sm:flex-row-reverse' : ''}`}>
              <Button variant="outline" className="flex-1" onClick={handleReset}>
                {t('orderTracking.trackAnother')}
              </Button>
              <Link href={`/${locale}/orders`} className="flex-1">
                <Button className="w-full">{t('orderTracking.orderHistory')}</Button>
              </Link>
            </div>
          </div>
        )
        )}
      </div>
    </>
  );
}
