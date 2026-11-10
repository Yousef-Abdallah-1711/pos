'use client'

import { useEffect, useRef, useState } from 'react'
import { Loader2, X } from 'lucide-react'
import { useAuthStore } from '@/stores/authStore'
import { getConfig } from '@/lib/config'
import { toast } from 'sonner'
import { useLanguage } from '@/contexts/LanguageContext'

interface Tenant {
  tenant_id: number
  tenant_name: string
  user_id: number
  roles: string[]
  default_restaurant_id?: number
}

interface TenantSelectorModalProps {
  isOpen: boolean
  email: string
  selectionToken: string
  tenants: Tenant[]
  onTenantSelected?: (tenantId: number) => void
  onCancel?: () => void
}

export function TenantSelectorModal({
  isOpen,
  email,
  selectionToken,
  tenants,
  onTenantSelected,
  onCancel,
}: TenantSelectorModalProps) {
  const setAuth = useAuthStore((state) => state.setAuth)
  const { t, dir, language } = useLanguage()
  const [selectedTenantId, setSelectedTenantId] = useState<number | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const dialogRef = useRef<HTMLDialogElement>(null)

  useEffect(() => {
    const dialog = dialogRef.current
    if (!dialog) return
    if (isOpen && !dialog.open) {
      if (typeof dialog.showModal === 'function') dialog.showModal()
      else dialog.setAttribute('open', '')
    }
    if (!isOpen && dialog.open) {
      if (typeof dialog.close === 'function') dialog.close()
      else dialog.removeAttribute('open')
    }
  }, [isOpen])

  if (!isOpen || !tenants.length) return null

  const handleSelectTenant = async (tenantId: number) => {
    setSelectedTenantId(tenantId)
    setIsLoading(true)
    setError(null)

    try {
      const response = await fetch(`${getConfig().apiUrl}/auth/login/confirm`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, tenant_id: tenantId, selection_token: selectionToken }),
      })

      if (!response.ok) {
        const errorData: unknown = await response.json()
        const apiMessage = typeof errorData === 'object' && errorData !== null && 'error' in errorData && typeof errorData.error === 'string'
          ? errorData.error
          : t.auth.tenantSelector.failedToSelect
        throw new Error(apiMessage)
      }

      const data: { user?: Parameters<typeof setAuth>[0] } = await response.json()
      if (!data.user) throw new Error(t.auth.tenantSelector.failedToSelect)

      setAuth(data.user, '')
      toast.success(t.auth.login.loginSuccess)
      onTenantSelected?.(tenantId)
      const { dashboardUrl } = getConfig()
      window.location.href = `${dashboardUrl}/${language}/auth/callback`
    } catch (requestError: unknown) {
      const message = requestError instanceof Error ? requestError.message : t.auth.tenantSelector.anError
      setError(message)
      toast.error(message)
      setIsLoading(false)
      setSelectedTenantId(null)
    }
  }

  return (
    <dialog ref={dialogRef} className="tenant-dialog" dir={dir} aria-labelledby="tenant-dialog-title" onCancel={(event) => { event.preventDefault(); onCancel?.() }}>
      <div className="tenant-dialog__content">
        <header className="tenant-dialog__header">
          <div>
            <h2 id="tenant-dialog-title">{t.auth.tenantSelector.title}</h2>
            <p>{t.auth.tenantSelector.description}</p>
          </div>
          <button type="button" className="tenant-dialog__close" onClick={onCancel} aria-label={language === 'ar' ? 'العودة إلى تسجيل الدخول' : 'Return to sign in'} disabled={isLoading}>
            <X aria-hidden="true" size={20} />
          </button>
        </header>

        {error && <p className="tenant-dialog__error" role="alert">{error}</p>}

        <div className="tenant-dialog__list" aria-label={t.auth.tenantSelector.title}>
          {tenants.map((tenant) => (
            <button
              key={tenant.tenant_id}
              type="button"
              onClick={() => handleSelectTenant(tenant.tenant_id)}
              disabled={isLoading}
              aria-pressed={selectedTenantId === tenant.tenant_id}
              className={`tenant-dialog__option ${selectedTenantId === tenant.tenant_id ? 'is-selected' : ''}`}
            >
              <span>
                <strong>{tenant.tenant_name}</strong>
                <span className="tenant-dialog__roles">
                  {tenant.roles?.length
                    ? tenant.roles.map((role) => <span className="tenant-dialog__role" key={role}>{role}</span>)
                    : (language === 'ar' ? 'لا توجد أدوار معروضة' : 'No roles listed')}
                </span>
              </span>
              {selectedTenantId === tenant.tenant_id && isLoading && <Loader2 aria-hidden="true" className="animate-spin" size={20} />}
            </button>
          ))}
        </div>

        <p className="tenant-dialog__account">{t.auth.tenantSelector.loggedInAs} <strong dir="ltr">{email}</strong></p>
      </div>
    </dialog>
  )
}
