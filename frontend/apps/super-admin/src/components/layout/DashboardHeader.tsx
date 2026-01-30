"use client";

import React from "react";
import { useAuthStore } from "@/stores/authStore";
import { AlertCircle, Clock } from "lucide-react";

export function DashboardHeader() {
  const superAdmin = useAuthStore((state) => state.superAdmin);
  const impersonation = useAuthStore((state) => state.impersonation);

  return (
    <header className="bg-white shadow border-b border-gray-200 px-6 py-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold text-gray-900">
            Super Admin Dashboard
          </h2>
          <p className="text-sm text-gray-600 mt-1">
            Logged in as: <span className="font-medium">{superAdmin?.email}</span>
          </p>
        </div>

        {/* Impersonation Alert */}
        {impersonation?.is_impersonating && (
          <div className="bg-amber-50 border border-amber-200 rounded-lg px-4 py-3 flex items-start gap-3">
            <AlertCircle className="text-amber-600 mt-0.5 flex-shrink-0" size={20} />
            <div>
              <p className="text-sm font-medium text-amber-900">
                Impersonating User
              </p>
              <p className="text-xs text-amber-800 mt-1">
                {impersonation?.target_user?.first_name} {impersonation?.target_user?.last_name}
              </p>
              {impersonation?.expires_at && (
                <p className="text-xs text-amber-700 mt-1 flex items-center gap-1">
                  <Clock size={12} />
                  Session expires: {new Date(impersonation.expires_at).toLocaleTimeString()}
                </p>
              )}
            </div>
          </div>
        )}
      </div>
    </header>
  );
}
