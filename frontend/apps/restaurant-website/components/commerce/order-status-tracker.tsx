'use client'

import { CheckCircle2, Clock, Truck, AlertCircle } from 'lucide-react'

export interface OrderStatus {
  status:
    | 'pending'
    | 'confirmed'
    | 'preparing'
    | 'ready'
    | 'out_for_delivery'
    | 'delivered'
    | 'cancelled'
  timestamp?: string
}

interface OrderStatusTrackerProps {
  statuses: OrderStatus[]
  currentStatus: OrderStatus['status']
  locale?: 'en' | 'ar'
}

export function OrderStatusTracker({ statuses, currentStatus, locale = 'en' }: OrderStatusTrackerProps) {
  const isRTL = locale === 'ar'

  const statusConfig = {
    pending: {
      label: locale === 'en' ? 'Pending' : 'قيد الانتظار',
      icon: Clock,
    },
    confirmed: {
      label: locale === 'en' ? 'Confirmed' : 'تم التأكيد',
      icon: CheckCircle2,
    },
    preparing: {
      label: locale === 'en' ? 'Preparing' : 'جاري التحضير',
      icon: Clock,
    },
    ready: {
      label: locale === 'en' ? 'Ready' : 'جاهز',
      icon: CheckCircle2,
    },
    out_for_delivery: {
      label: locale === 'en' ? 'Out for Delivery' : 'في الطريق',
      icon: Truck,
    },
    delivered: {
      label: locale === 'en' ? 'Delivered' : 'تم التسليم',
      icon: CheckCircle2,
    },
    cancelled: {
      label: locale === 'en' ? 'Cancelled' : 'ملغي',
      icon: AlertCircle,
    },
  }

  const timeline = ['pending', 'confirmed', 'preparing', 'ready', 'out_for_delivery', 'delivered'] as const

  const getStatusIndex = (status: OrderStatus['status']) => timeline.findIndex(value => value === status)
  const currentIndex = getStatusIndex(currentStatus)
  const isCancelled = currentStatus === 'cancelled'

  return (
    <div className="w-full">
      {isCancelled ? (
        <div className="tracker-cancelled border rounded-lg p-6">
          <div className="flex items-center gap-3">
            <AlertCircle className="tracker-cancelled-icon h-6 w-6" />
            <div>
              <h3 className="tracker-cancelled-title font-semibold text-lg">{statusConfig.cancelled.label}</h3>
              <p className="tracker-cancelled-text text-sm">
                {locale === 'en' ? 'This order has been cancelled' : 'تم إلغاء هذا الطلب'}
              </p>
            </div>
          </div>
        </div>
      ) : (
        <div className={`${isRTL ? 'flex-row-reverse' : ''}`}>
          <div className="grid grid-cols-2 gap-x-4 gap-y-6 sm:grid-cols-3 lg:grid-cols-6">
            {timeline.map((status, index) => {
              const config = statusConfig[status as keyof typeof statusConfig]
              const Icon = config.icon
              const isCompleted = index < currentIndex
              const isActive = index === currentIndex

              return (
                <div key={status} className={`min-w-0 relative ${isRTL ? 'text-right' : ''}`}>
                  <div
                    className={`tracker-circle h-12 w-12 rounded-full flex items-center justify-center mb-3 transition-all border-2 ${
                      isCompleted || isActive ? `tracker-soft-${status}` : 'tracker-circle-inactive'
                    }`}
                  >
                    <Icon
                      className={`tracker-icon h-6 w-6 ${
                        isCompleted || isActive ? `tracker-text-${status}` : 'tracker-icon-inactive'
                      }`}
                    />
                  </div>

                  {index < timeline.length - 1 && (
                    <div
                      className={`tracker-connector hidden lg:block absolute top-6 h-1 w-full ${
                        isRTL ? 'right-full pr-3' : 'left-full pl-3'
                      } ${isCompleted ? `tracker-line-${status}` : 'tracker-connector-inactive'}`}
                    />
                  )}

                  <p
                    className={`tracker-label break-words text-sm ${
                      isActive
                        ? 'tracker-label-active'
                        : isCompleted
                          ? 'tracker-label-complete'
                          : 'tracker-label-pending'
                    }`}
                  >
                    {config.label}
                  </p>

                  {statuses[index]?.timestamp && (
                    <p className="tracker-time text-xs mt-1">
                      {new Date(statuses[index].timestamp!).toLocaleString(locale === 'ar' ? 'ar-EG' : 'en-US')}
                    </p>
                  )}
                </div>
              )
            })}
          </div>

          <div className={`tracker-current-box tracker-soft-${currentStatus} mt-8 p-4 border rounded-lg`}>
            <p className="tracker-current-title font-semibold mb-1">{statusConfig[currentStatus as keyof typeof statusConfig].label}</p>
            <p className="tracker-current-text text-sm">
              {currentStatus === 'pending' &&
                (locale === 'en'
                  ? 'Your order is waiting to be confirmed'
                  : 'طلبك قيد الانتظار للتأكيد')}
              {currentStatus === 'confirmed' &&
                (locale === 'en'
                  ? 'Your order has been confirmed and will be prepared shortly'
                  : 'تم تأكيد طلبك وسيتم تحضيره قريباً')}
              {currentStatus === 'preparing' &&
                (locale === 'en'
                  ? 'Our team is preparing your delicious food'
                  : 'فريقنا يقوم بتحضير طعامك اللذيذ')}
              {currentStatus === 'ready' &&
                (locale === 'en'
                  ? 'Your order is ready for pickup or out for delivery'
                  : 'طلبك جاهز للالتقاط أو للتسليم')}
              {currentStatus === 'out_for_delivery' &&
                (locale === 'en' ? 'Your order is on its way to you' : 'طلبك في الطريق إليك')}
              {currentStatus === 'delivered' &&
                (locale === 'en'
                  ? 'Your order has been delivered. Thank you for your order!'
                  : 'تم تسليم طلبك. شكراً لطلبك!')}
            </p>
          </div>
        </div>
      )}

      <style jsx>{`
        .tracker-cancelled {
          background-color: #fee2e2;
          border-color: #fecaca;
        }

        .tracker-cancelled-icon {
          color: #ef4444;
        }

        .tracker-cancelled-title {
          color: #991b1b;
        }

        .tracker-cancelled-text {
          color: hsl(var(--theme-text, 221 39% 11%) / 0.6);
        }

        .tracker-circle-inactive {
          background-color: hsl(var(--theme-background, 0 0% 100%) / 0.5);
          border-color: hsl(var(--theme-border, 0 0% 90%));
        }

        .tracker-icon-inactive {
          color: hsl(var(--theme-text, 221 39% 11%) / 0.3);
        }

        .tracker-connector-inactive {
          background-color: hsl(var(--theme-text, 221 39% 11%) / 0.1);
        }

        .tracker-label {
          font-weight: 500;
        }

        .tracker-label-active {
          color: hsl(var(--theme-primary, 24 95% 53%));
          font-weight: 700;
        }

        .tracker-label-complete {
          color: hsl(var(--theme-text, 221 39% 11%));
        }

        .tracker-label-pending,
        .tracker-time {
          color: hsl(var(--theme-text, 221 39% 11%) / 0.6);
        }

        .tracker-current-box {
          border-color: hsl(var(--theme-border, 0 0% 90%));
        }

        .tracker-current-title {
          color: hsl(var(--theme-text, 221 39% 11%));
        }

        .tracker-current-text {
          color: hsl(var(--theme-text, 221 39% 11%) / 0.7);
        }

        .tracker-soft-pending {
          background-color: #fef3c7;
          border-color: #fbbf24;
        }

        .tracker-soft-confirmed {
          background-color: #dbeafe;
          border-color: #3b82f6;
        }

        .tracker-soft-preparing {
          background-color: #fed7aa;
          border-color: #f97316;
        }

        .tracker-soft-ready {
          background-color: #dcfce7;
          border-color: #22c55e;
        }

        .tracker-soft-out_for_delivery {
          background-color: #e9d5ff;
          border-color: #a855f7;
        }

        .tracker-soft-delivered {
          background-color: #dcfce7;
          border-color: #16a34a;
        }

        .tracker-line-pending {
          background-color: #fbbf24;
        }

        .tracker-line-confirmed {
          background-color: #3b82f6;
        }

        .tracker-line-preparing {
          background-color: #f97316;
        }

        .tracker-line-ready {
          background-color: #22c55e;
        }

        .tracker-line-out_for_delivery {
          background-color: #a855f7;
        }

        .tracker-line-delivered {
          background-color: #16a34a;
        }

        .tracker-text-pending {
          color: #f59e0b;
        }

        .tracker-text-confirmed {
          color: #3b82f6;
        }

        .tracker-text-preparing {
          color: #f97316;
        }

        .tracker-text-ready {
          color: #22c55e;
        }

        .tracker-text-out_for_delivery {
          color: #a855f7;
        }

        .tracker-text-delivered {
          color: #16a34a;
        }
      `}</style>
    </div>
  )
}
