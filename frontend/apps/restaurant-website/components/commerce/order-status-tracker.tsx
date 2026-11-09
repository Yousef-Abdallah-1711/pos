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
      label: locale === 'en' ? 'Pending' : 'Ù‚ÙŠØ¯ Ø§Ù„Ø§Ù†ØªØ¸Ø§Ø±',
      icon: Clock,
    },
    confirmed: {
      label: locale === 'en' ? 'Confirmed' : 'ØªÙ… Ø§Ù„ØªØ£ÙƒÙŠØ¯',
      icon: CheckCircle2,
    },
    preparing: {
      label: locale === 'en' ? 'Preparing' : 'Ø¬Ø§Ø±ÙŠ Ø§Ù„ØªØ­Ø¶ÙŠØ±',
      icon: Clock,
    },
    ready: {
      label: locale === 'en' ? 'Ready' : 'Ø¬Ø§Ù‡Ø²',
      icon: CheckCircle2,
    },
    out_for_delivery: {
      label: locale === 'en' ? 'Out for Delivery' : 'ÙÙŠ Ø§Ù„Ø·Ø±ÙŠÙ‚',
      icon: Truck,
    },
    delivered: {
      label: locale === 'en' ? 'Delivered' : 'ØªÙ… Ø§Ù„ØªØ³Ù„ÙŠÙ…',
      icon: CheckCircle2,
    },
    cancelled: {
      label: locale === 'en' ? 'Cancelled' : 'Ù…Ù„ØºÙŠ',
      icon: AlertCircle,
    },
  }

  const timeline = ['pending', 'confirmed', 'preparing', 'ready', 'out_for_delivery', 'delivered'] as const

  const getStatusIndex = (status: string) => timeline.indexOf(status as any)
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
                {locale === 'en' ? 'This order has been cancelled' : 'ØªÙ… Ø¥Ù„ØºØ§Ø¡ Ù‡Ø°Ø§ Ø§Ù„Ø·Ù„Ø¨'}
              </p>
            </div>
          </div>
        </div>
      ) : (
        <div className={`${isRTL ? 'flex-row-reverse' : ''}`}>
          <div className={`flex items-start ${isRTL ? 'flex-row-reverse' : ''}`}>
            {timeline.map((status, index) => {
              const config = statusConfig[status as keyof typeof statusConfig]
              const Icon = config.icon
              const isCompleted = index < currentIndex
              const isActive = index === currentIndex

              return (
                <div key={status} className={`flex-1 relative ${isRTL ? 'text-right' : ''}`}>
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
                      className={`tracker-connector absolute top-6 h-1 w-full ${
                        isRTL ? 'right-full pr-3' : 'left-full pl-3'
                      } ${isCompleted ? `tracker-line-${status}` : 'tracker-connector-inactive'}`}
                    />
                  )}

                  <p
                    className={`tracker-label text-sm ${
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
                  : 'Ø·Ù„Ø¨Ùƒ Ù‚ÙŠØ¯ Ø§Ù„Ø§Ù†ØªØ¸Ø§Ø± Ù„Ù„ØªØ£ÙƒÙŠØ¯')}
              {currentStatus === 'confirmed' &&
                (locale === 'en'
                  ? 'Your order has been confirmed and will be prepared shortly'
                  : 'ØªÙ… ØªØ£ÙƒÙŠØ¯ Ø·Ù„Ø¨Ùƒ ÙˆØ³ÙŠØªÙ… ØªØ­Ø¶ÙŠØ±Ù‡ Ù‚Ø±ÙŠØ¨Ø§Ù‹')}
              {currentStatus === 'preparing' &&
                (locale === 'en'
                  ? 'Our team is preparing your delicious food'
                  : 'ÙØ±ÙŠÙ‚Ù†Ø§ ÙŠÙ‚ÙˆÙ… Ø¨ØªØ­Ø¶ÙŠØ± Ø·Ø¹Ø§Ù…Ùƒ Ø§Ù„Ù„Ø°ÙŠØ°')}
              {currentStatus === 'ready' &&
                (locale === 'en'
                  ? 'Your order is ready for pickup or out for delivery'
                  : 'Ø·Ù„Ø¨Ùƒ Ø¬Ø§Ù‡Ø² Ù„Ù„Ø§Ù„ØªÙ‚Ø§Ø· Ø£Ùˆ Ù„Ù„ØªØ³Ù„ÙŠÙ…')}
              {currentStatus === 'out_for_delivery' &&
                (locale === 'en' ? 'Your order is on its way to you' : 'Ø·Ù„Ø¨Ùƒ ÙÙŠ Ø§Ù„Ø·Ø±ÙŠÙ‚ Ø¥Ù„ÙŠÙƒ')}
              {currentStatus === 'delivered' &&
                (locale === 'en'
                  ? 'Your order has been delivered. Thank you for your order!'
                  : 'ØªÙ… ØªØ³Ù„ÙŠÙ… Ø·Ù„Ø¨Ùƒ. Ø´ÙƒØ±Ø§Ù‹ Ù„Ø·Ù„Ø¨Ùƒ!')}
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
