"use client";

import React, { useState } from "react";
import { useAuthStore } from "@/stores/authStore";
import { superAdminAPI } from "@/lib/api";
import { AlertTriangle, X, LogOut } from "lucide-react";
import { toast } from "sonner";

export function ImpersonationBanner() {
  const impersonation = useAuthStore((state) => state.impersonation);
  const { setImpersonation } = useAuthStore();
  const [isEnding, setIsEnding] = useState(false);

  if (!impersonation?.is_impersonating) {
    return null;
  }

  const handleEndImpersonation = async () => {
    try {
      setIsEnding(true);
      await superAdminAPI.endImpersonation();
      setImpersonation(null);
      toast.success("Impersonation session ended");
    } catch (error) {
      toast.error("Failed to end impersonation session");
      console.error(error);
    } finally {
      setIsEnding(false);
    }
  };

  return (
    <div className="bg-amber-50 border-b-2 border-amber-400 px-4 py-4">
      <div className="max-w-7xl mx-auto flex items-center justify-between gap-4">
        <div className="flex items-start gap-3">
          <AlertTriangle className="text-amber-600 flex-shrink-0 mt-0.5" size={20} />
          <div>
            <p className="font-semibold text-amber-900">
              Impersonating User
            </p>
            <p className="text-sm text-amber-800 mt-1">
              {impersonation?.target_user?.first_name}{" "}
              {impersonation?.target_user?.last_name} (
              {impersonation?.target_user?.email})
            </p>
            {impersonation?.expires_at && (
              <p className="text-xs text-amber-700 mt-1">
                Session expires:{" "}
                {new Date(impersonation.expires_at).toLocaleTimeString()}
              </p>
            )}
          </div>
        </div>

        <button
          onClick={handleEndImpersonation}
          disabled={isEnding}
          className="flex items-center gap-2 px-4 py-2 bg-amber-600 hover:bg-amber-700 disabled:bg-gray-400 text-white rounded-lg font-medium transition-colors flex-shrink-0"
        >
          <LogOut size={16} />
          {isEnding ? "Ending..." : "End Session"}
        </button>
      </div>
    </div>
  );
}
