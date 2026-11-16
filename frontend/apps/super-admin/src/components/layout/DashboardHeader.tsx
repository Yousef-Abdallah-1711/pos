"use client";

import React from "react";
import { useAuthStore } from "@/stores/authStore";
import { AlertCircle, Clock } from "lucide-react";
import { PreferenceControls } from "@/components/layout/PreferenceControls";
import { useSuperAdminPreferences, useSuperAdminText } from "@/components/providers/PreferencesProvider";

export function DashboardHeader() {
  const t = useSuperAdminText();
  const { locale } = useSuperAdminPreferences();
  const superAdmin = useAuthStore((state) => state.superAdmin);
  const impersonation = useAuthStore((state) => state.impersonation);

  return (
    <header className="border-b border-gray-200 bg-white px-4 py-4 shadow-sm dark:border-gray-700 dark:bg-gray-900 sm:px-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-semibold text-gray-900">
            {t("Super Admin Dashboard")}
          </h2>
          <p className="text-sm text-gray-600 mt-1">
            {t("Logged in as:")} <span dir="ltr" className="font-medium">{superAdmin?.email}</span>
          </p>
        </div>

        {/* Impersonation Alert */}
        {impersonation?.is_impersonating && (
          <div className="bg-amber-50 border border-amber-200 rounded-lg px-4 py-3 flex items-start gap-3">
            <AlertCircle className="text-amber-600 mt-0.5 flex-shrink-0" size={20} />
            <div>
              <p className="text-sm font-medium text-amber-900">
                {t("Impersonating User")}
              </p>
              <p className="text-xs text-amber-800 mt-1">
                {impersonation?.target_user?.first_name} {impersonation?.target_user?.last_name}
              </p>
              {impersonation?.expires_at && (
                <p className="text-xs text-amber-700 mt-1 flex items-center gap-1">
                  <Clock size={12} />
                  {t("Session expires:")} {new Date(impersonation.expires_at).toLocaleTimeString(locale)}
                </p>
              )}
            </div>
          </div>
        )}
        <PreferenceControls />
      </div>
    </header>
  );
}
