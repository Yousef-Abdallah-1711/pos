"use client";

import Link from "next/link";
import { Button } from "@pos-saas/ui";
import { Eye, RefreshCw, AlertCircle } from "lucide-react";

export interface OrderItem {
  id: number;
  orderNumber: string;
  status: "pending" | "confirmed" | "preparing" | "ready" | "out_for_delivery" | "delivered" | "cancelled";
  totalAmount: number;
  createdAt: string;
  estimatedDelivery?: string;
  restaurantName?: string;
}

interface OrderListProps {
  orders: OrderItem[];
  locale?: "en" | "ar";
  onRefresh?: () => void;
  isLoading?: boolean;
}

/**
 * Order list component displaying all user orders
 */
export function OrderList({
  orders,
  locale = "en",
  onRefresh,
  isLoading = false,
}: OrderListProps) {
  const isRTL = locale === "ar";

  const statusColors: Record<string, { bg: string; text: string }> = {
    pending: { bg: "bg-yellow-100 dark:bg-yellow-900/30", text: "text-yellow-800 dark:text-yellow-200" },
    confirmed: { bg: "bg-blue-100 dark:bg-blue-900/30", text: "text-blue-800 dark:text-blue-200" },
    preparing: { bg: "bg-orange-100 dark:bg-orange-900/30", text: "text-orange-800 dark:text-orange-200" },
    ready: { bg: "bg-green-100 dark:bg-green-900/30", text: "text-green-800 dark:text-green-200" },
    out_for_delivery: { bg: "bg-purple-100 dark:bg-purple-900/30", text: "text-purple-800 dark:text-purple-200" },
    delivered: { bg: "bg-green-100 dark:bg-green-900/30", text: "text-green-800 dark:text-green-200" },
    cancelled: { bg: "bg-red-100 dark:bg-red-900/30", text: "text-red-800 dark:text-red-200" },
  };

  const statusLabels = {
    en: {
      pending: "Pending",
      confirmed: "Confirmed",
      preparing: "Preparing",
      ready: "Ready",
      out_for_delivery: "Out for Delivery",
      delivered: "Delivered",
      cancelled: "Cancelled",
      noOrders: "No orders yet",
      noOrdersDesc: "You haven't placed any orders. Start by browsing our menu!",
      viewDetails: "View Details",
      trackOrder: "Track Order",
      orderDate: "Order Date",
      estimatedDelivery: "Estimated Delivery",
      orderHistory: "Order History",
    },
    ar: {
      pending: "قيد الانتظار",
      confirmed: "تم التأكيد",
      preparing: "جاري التحضير",
      ready: "جاهز",
      out_for_delivery: "في الطريق",
      delivered: "تم التسليم",
      cancelled: "ملغي",
      noOrders: "لا توجد طلبات",
      noOrdersDesc: "لم تقم بتقديم أي طلبات بعد. ابدأ باستعراض قائمتنا!",
      viewDetails: "عرض التفاصيل",
      trackOrder: "تتبع الطلب",
      orderDate: "تاريخ الطلب",
      estimatedDelivery: "التسليم المتوقع",
      orderHistory: "سجل الطلبات",
    },
  };

  const t = statusLabels[locale as keyof typeof statusLabels] || statusLabels.en;

  if (orders.length === 0) {
    return (
      <div className="text-center py-12">
        <AlertCircle className="h-16 w-16 text-gray-400 mx-auto mb-4" />
        <h3 className="text-xl font-semibold text-gray-600 dark:text-gray-300 mb-2">
          {t.noOrders}
        </h3>
        <p className="text-gray-500 dark:text-gray-400">
          {t.noOrdersDesc}
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header with Refresh */}
      <div className={`flex justify-between items-center mb-6 ${isRTL ? "flex-row-reverse" : ""}`}>
        <h2 className="text-2xl font-bold">{t.orderHistory}</h2>
        <Button
          variant="outline"
          size="sm"
          onClick={onRefresh}
          disabled={isLoading}
          className={isRTL ? "flex-row-reverse" : ""}
        >
          <RefreshCw className={`h-4 w-4 ${isLoading ? "animate-spin" : ""}`} />
          {isLoading ? "..." : "Refresh"}
        </Button>
      </div>

      {/* Orders Grid */}
      <div className="space-y-4">
        {orders.map((order) => {
          const colors = statusColors[order.status];
          const statusLabel =
            t[order.status as keyof typeof t] || order.status;

          return (
            <div
              key={order.id}
              className={`bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow ${
                isRTL ? "text-right" : ""
              }`}
            >
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4 items-start mb-4">
                {/* Order Number */}
                <div>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {locale === "en" ? "Order #" : "الطلب #"}
                  </p>
                  <p className="text-lg font-bold text-primary">
                    {order.orderNumber}
                  </p>
                </div>

                {/* Amount */}
                <div>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {locale === "en" ? "Total Amount" : "المبلغ الإجمالي"}
                  </p>
                  <p className="text-lg font-bold">
                    ${order.totalAmount.toFixed(2)}
                  </p>
                </div>

                {/* Date */}
                <div>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {t.orderDate}
                  </p>
                  <p className="text-sm font-medium">
                    {new Date(order.createdAt).toLocaleDateString(
                      locale === "ar" ? "ar-EG" : "en-US"
                    )}
                  </p>
                </div>

                {/* Status */}
                <div className={isRTL ? "text-right" : ""}>
                  <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">
                    {locale === "en" ? "Status" : "الحالة"}
                  </p>
                  <span
                    className={`inline-block px-3 py-1 rounded-full text-sm font-semibold ${colors.bg} ${colors.text}`}
                  >
                    {statusLabel}
                  </span>
                </div>
              </div>

              {/* Estimated Delivery */}
              {order.estimatedDelivery && (
                <div className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                  <p>
                    {t.estimatedDelivery}:{" "}
                    {new Date(order.estimatedDelivery).toLocaleString(
                      locale === "ar" ? "ar-EG" : "en-US"
                    )}
                  </p>
                </div>
              )}

              {/* Restaurant Name */}
              {order.restaurantName && (
                <div className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                  <p>
                    {locale === "en" ? "Restaurant: " : "المطعم: "}
                    <span className="font-medium text-gray-900 dark:text-gray-100">
                      {order.restaurantName}
                    </span>
                  </p>
                </div>
              )}

              {/* Actions */}
              <div className={`flex gap-2 pt-4 border-t border-gray-200 dark:border-gray-700 ${isRTL ? "flex-row-reverse" : ""}`}>
                <Link href={`/${locale}/orders/${order.id}`}>
                  <Button variant="outline" size="sm" className={isRTL ? "flex-row-reverse" : ""}>
                    <Eye className="h-4 w-4" />
                    {t.viewDetails}
                  </Button>
                </Link>
                <Link href={`/${locale}/orders/track/${order.orderNumber}`}>
                  <Button variant="default" size="sm" className={isRTL ? "flex-row-reverse" : ""}>
                    <RefreshCw className="h-4 w-4" />
                    {t.trackOrder}
                  </Button>
                </Link>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
