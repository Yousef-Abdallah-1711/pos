'use client'

import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { RBACSettings } from '@/components/settings/RBACSettings'
import { Loader2 } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'

export default function AccessControlPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const user = useAuthStore((state) => state.user)

  if (!user) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
        <div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          {t('settings.sidebar.accessControlTitle', 'Access Control')}
        </h2>
        <p className="text-gray-600 dark:text-gray-400">
          {t('settings.sidebar.accessControlDesc', 'View your current roles and permissions')}
        </p>
      </div>

      <RBACSettings
        userId={user.id}
        tenantId={user.tenant_id || 1}
        userName={user.name}
      />
    </div>
  )
}
