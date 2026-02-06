'use client'

import { AlertCircle, Clock, Lock, Phone, Mail } from 'lucide-react'

interface MaintenancePageProps {
  restaurantName?: string
  reason?: 'trial_expired' | 'tenant_disabled' | 'maintenance' | 'error'
  contactEmail?: string
  contactPhone?: string
  estimatedReturnTime?: string
}

export function MaintenancePage({
  restaurantName = 'Our Restaurant',
  reason = 'maintenance',
  contactEmail,
  contactPhone,
  estimatedReturnTime,
}: MaintenancePageProps) {
  const getContent = () => {
    switch (reason) {
      case 'trial_expired':
        return {
          icon: <Clock className="h-16 w-16 text-orange-500" />,
          title: "We'll Be Back Soon!",
          subtitle: 'Trial Period Ended',
          description:
            'Our online ordering system is temporarily unavailable as we update our service. The restaurant is currently working to restore access. Please check back soon or contact us directly for orders.',
          color: 'orange',
        }
      case 'tenant_disabled':
        return {
          icon: <Lock className="h-16 w-16 text-red-500" />,
          title: 'Service Temporarily Unavailable',
          subtitle: 'Account Disabled',
          description:
            'We apologize for the inconvenience. Our online ordering system is temporarily unavailable. Please contact the restaurant directly to place your order.',
          color: 'red',
        }
      case 'error':
        return {
          icon: <AlertCircle className="h-16 w-16 text-red-500" />,
          title: 'Something Went Wrong',
          subtitle: 'Technical Issue',
          description:
            'Our online ordering system is experiencing technical difficulties. Please try again later or contact us directly for assistance.',
          color: 'red',
        }
      default:
        return {
          icon: <Clock className="h-16 w-16 text-blue-500" />,
          title: 'We Are Currently Offline',
          subtitle: 'Scheduled Maintenance',
          description:
            'We are performing scheduled maintenance to improve your experience. We will be back online shortly. Thank you for your patience!',
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
            <p className="text-sm text-gray-600 dark:text-gray-300 mb-1">Estimated return time:</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">{estimatedReturnTime}</p>
          </div>
        )}

        {/* Contact Information */}
        {(contactPhone || contactEmail) && (
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 mb-8">
            <h3 className="font-semibold text-gray-900 dark:text-white mb-4 text-center">Contact Us</h3>

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
            Thank you for choosing <span className="font-semibold text-gray-900 dark:text-white">{restaurantName}</span>
          </p>
        )}
      </div>
    </div>
  )
}
