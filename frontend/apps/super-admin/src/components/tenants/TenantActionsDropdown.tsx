'use client'

import { useState } from 'react'
import { ExtendTrialModal } from './ExtendTrialModal'
import { DisableTenantModal } from './DisableTenantModal'
import { EnableTenantModal } from './EnableTenantModal'
import { DeleteTenantModal } from './DeleteTenantModal'
import type { TenantSubscriptionInfo } from '@/types'

interface TenantActionsDropdownProps {
  tenant: TenantSubscriptionInfo & { id?: number }
  onImpersonate: () => void
  onActionsComplete: () => void
}

export function TenantActionsDropdown({
  tenant,
  onImpersonate,
  onActionsComplete,
}: TenantActionsDropdownProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [extendModalOpen, setExtendModalOpen] = useState(false)
  const [disableModalOpen, setDisableModalOpen] = useState(false)
  const [enableModalOpen, setEnableModalOpen] = useState(false)
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)

  const tenantId = tenant.id || tenant.tenant_id

  return (
    <>
      <div className="relative inline-block">
        <button
          onClick={() => setIsOpen(!isOpen)}
          className="px-3 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md transition-colors"
        >
          Actions ▼
        </button>

        {isOpen && (
          <div className="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded-md shadow-lg z-10">
            <button
              onClick={() => {
                setExtendModalOpen(true)
                setIsOpen(false)
              }}
              className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 border-b"
            >
              📅 Extend Trial
            </button>

            {!tenant.is_disabled ? (
              <button
                onClick={() => {
                  setDisableModalOpen(true)
                  setIsOpen(false)
                }}
                className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 border-b"
              >
                🔒 Disable Tenant
              </button>
            ) : (
              <button
                onClick={() => {
                  setEnableModalOpen(true)
                  setIsOpen(false)
                }}
                className="w-full text-left px-4 py-2 text-sm text-green-600 hover:bg-green-50 border-b"
              >
                ✓ Enable Tenant
              </button>
            )}

            <button
              onClick={() => {
                onImpersonate()
                setIsOpen(false)
              }}
              className="w-full text-left px-4 py-2 text-sm text-blue-600 hover:bg-blue-50 border-b"
            >
              👤 Impersonate
            </button>

            <button
              onClick={() => {
                setDeleteModalOpen(true)
                setIsOpen(false)
              }}
              className="w-full text-left px-4 py-2 text-sm text-red-700 hover:bg-red-50 rounded-b"
            >
              🗑️ Delete Tenant
            </button>
          </div>
        )}
      </div>

      <ExtendTrialModal
        tenantId={tenantId}
        tenantName={tenant.tenant_name}
        currentTrialEndDate={tenant.trial_end_date}
        isOpen={extendModalOpen}
        onClose={() => setExtendModalOpen(false)}
        onSuccess={() => {
          setExtendModalOpen(false)
          onActionsComplete()
        }}
      />

      <DisableTenantModal
        tenantId={tenantId}
        tenantName={tenant.tenant_name}
        isOpen={disableModalOpen}
        onClose={() => setDisableModalOpen(false)}
        onSuccess={() => {
          setDisableModalOpen(false)
          onActionsComplete()
        }}
      />

      <EnableTenantModal
        tenantId={tenantId}
        tenantName={tenant.tenant_name}
        disabledReason={tenant.disabled_reason}
        disabledAt={tenant.disabled_at}
        isOpen={enableModalOpen}
        onClose={() => setEnableModalOpen(false)}
        onSuccess={() => {
          setEnableModalOpen(false)
          onActionsComplete()
        }}
      />

      <DeleteTenantModal
        tenantId={tenantId}
        tenantName={tenant.tenant_name}
        isOpen={deleteModalOpen}
        onClose={() => setDeleteModalOpen(false)}
        onSuccess={() => {
          setDeleteModalOpen(false)
          onActionsComplete()
        }}
      />
    </>
  )
}
