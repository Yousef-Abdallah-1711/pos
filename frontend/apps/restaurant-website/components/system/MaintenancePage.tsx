'use client'

import { AlertCircle, Clock, Lock, Phone, Mail } from 'lucide-react'
import { createTranslator } from '@/lib/translations'
import type { Locale } from '@/lib/translations'

interface MaintenancePageProps {
  restaurantName?: string
  reason?: 'trial_expired' | 'tenant_disabled' | 'maintenance' | 'error'
  contactEmail?: string
  contactPhone?: string
  estimatedReturnTime?: string
  locale?: Locale
}

export function MaintenancePage({
  restaurantName = 'Our Restaurant',
  reason = 'maintenance',
  contactEmail,
  contactPhone,
  estimatedReturnTime,
  locale = 'en',
}: MaintenancePageProps) {
  const t = createTranslator(locale)

  const getContent = () => {
    switch (reason) {
      case 'trial_expired':
        return {
          icon: <Clock className="h-16 w-16 text-orange-500" />,
          title: t('maintenance.trialExpired.title'),
          subtitle: t('maintenance.trialExpired.subtitle'),
          description: t('maintenance.trialExpired.description'),
          color: 'orange',
        }
      case 'tenant_disabled':
        return {
          icon: <Lock className="h-16 w-16 text-red-500" />,
          title: t('maintenance.tenantDisabled.title'),
          subtitle: t('maintenance.tenantDisabled.subtitle'),
          description: t('maintenance.tenantDisabled.description'),
          color: 'red',
        }
      case 'error':
        return {
          icon: <AlertCircle className="h-16 w-16 text-red-500" />,
          title: t('maintenance.error.title'),
          subtitle: t('maintenance.error.subtitle'),
          description: t('maintenance.error.description'),
          color: 'red',
        }
      default:
        return {
          icon: <Clock className="h-16 w-16 text-blue-500" />,
          title: t('maintenance.default.title'),
          subtitle: t('maintenance.default.subtitle'),
          description: t('maintenance.default.description'),
          color: 'blue',
        }
    }
  }

  const content = getContent()
  const bgGradient =
    content.color === 'orange'
      ? 'from-orange-50 to-white dark:from-orange-950/20 dark:to-gray-900'
      : content.color === 'red'
        ? 'from-red-50 to-white dark:from-red-950/20 dark:to-gray-900'
        : 'from-blue-50 to-white dark:from-blue-950/20 dark:to-gray-900'

  return (
    <div className={`min-h-screen bg-gradient-to-b ${bgGradient} flex items-center justify-center px-4 py-8`}>
      <div className="max-w-md w-full">
        {/* Icon */}
        <div className="flex justify-center mb-6">{content.icon}</div>

        {/* Title */}
        <h1 className="text-3xl md:text-4xl font-bold text-center text-gray-900 dark:text-white mb-2">
          {content.title}
        </h1>

        {/* Subtitle */}
        <p className="text-center text-sm text-gray-500 dark:text-gray-400 mb-6 font-medium">{content.subtitle}</p>

        {/* Description */}
        <p className="text-center text-gray-600 dark:text-gray-300 mb-8 leading-relaxed">{content.description}</p>

        {/* Estimated Return Time */}
        {estimatedReturnTime && (
          <div className="bg-white dark:bg-gray-800 rounded-lg p-4 mb-8 text-center border border-gray-200 dark:border-gray-700">
            <p className="text-sm text-gray-600 dark:text-gray-300 mb-1">{t('maintenance.estimatedReturnTime')}</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">{estimatedReturnTime}</p>
          </div>
        )}

        {/* Contact Information */}
        {(contactPhone || contactEmail) && (
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 mb-8">
            <h3 className="font-semibold text-gray-900 dark:text-white mb-4 text-center">{t('maintenance.contactUs')}</h3>

            <div className="space-y-3">
              {/* Phone */}
              {contactPhone && (
                <a
                  href={`tel:${contactPhone.replace(/\s+/g, '')}`}
                  className="flex items-center justify-center gap-3 text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition"
                >
                  <Phone className="h-5 w-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                  <span className="text-sm font-medium">{contactPhone}</span>
                </a>
              )}

              {/* Email */}
              {contactEmail && (
                <a
                  href={`mailto:${contactEmail}`}
                  className="flex items-center justify-center gap-3 text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition"
                >
                  <Mail className="h-5 w-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                  <span className="text-sm font-medium">{contactEmail}</span>
                </a>
              )}
            </div>
          </div>
        )}

        {/* Restaurant Name */}
        {restaurantName && (
          <p className="text-center text-gray-500 dark:text-gray-400 text-sm">
            {t('maintenance.thankYou')} <span className="font-semibold text-gray-900 dark:text-white">{restaurantName}</span>
          </p>
        )}
      </div>
    </div>
  )
}
