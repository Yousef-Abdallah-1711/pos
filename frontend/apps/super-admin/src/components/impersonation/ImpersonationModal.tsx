"use client";

import React, { useState } from "react";
import { useAuthStore } from "@/stores/authStore";
import { superAdminAPI } from "@/lib/api";
import type { User } from "@/types";
import { X, AlertCircle, Copy, Check } from "lucide-react";
import { toast } from "sonner";

interface ImpersonationModalProps {
  isOpen: boolean;
  onClose: () => void;
  user: User;
  tenantId: number;
}

export function ImpersonationModal({
  isOpen,
  onClose,
  user,
  tenantId,
}: ImpersonationModalProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [impersonationToken, setImpersonationToken] = useState<string | null>(
    null
  );
  const [dashboardUrl, setDashboardUrl] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const { setImpersonation } = useAuthStore();

  const handleStartImpersonation = async () => {
    try {
      setIsLoading(true);
      const response = await superAdminAPI.startImpersonation(user.id);

      setImpersonationToken(response.token);

      // Construct dashboard URL
      const baseUrl = process.env.NEXT_PUBLIC_API_URL?.replace(
        "/api/v1",
        ""
      ) || "http://localhost:3002";
      const url = `${baseUrl}/dashboard?impersonation_token=${response.token}`;
      setDashboardUrl(url);

      // Update impersonation state
      setImpersonation({
        is_impersonating: true,
        target_user: user,
        target_tenant: { id: tenantId, name: "", slug: "", created_at: "", updated_at: "" },
        expires_at: response.expires_at,
      });

      toast.success("Impersonation session started");
    } catch (error) {
      toast.error("Failed to start impersonation");
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopyUrl = () => {
    if (dashboardUrl) {
      navigator.clipboard.writeText(dashboardUrl);
      setCopied(true);
      toast.success("URL copied to clipboard");
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleOpenDashboard = () => {
    if (dashboardUrl) {
      window.open(dashboardUrl, "_blank");
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <h2 className="text-xl font-semibold text-gray-900">
            Impersonate User
          </h2>
          <button
            onClick={onClose}
            className="p-1 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X size={24} className="text-gray-600" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* User Info */}
          <div className="bg-gray-50 rounded-lg p-4">
            <p className="text-sm text-gray-600 mb-1">User to impersonate</p>
            <p className="font-semibold text-gray-900">
              {user.first_name} {user.last_name}
            </p>
            <p className="text-sm text-gray-600">{user.email}</p>
            <p className="text-xs text-gray-500 mt-2">Role: {user.role}</p>
          </div>

          {/* Warning */}
          <div className="bg-amber-50 border border-amber-200 rounded-lg p-4 flex gap-3">
            <AlertCircle className="text-amber-600 flex-shrink-0 mt-0.5" size={20} />
            <div className="text-sm text-amber-900">
              <p className="font-medium mb-1">Important</p>
              <p>
                You will have full access to this user's dashboard. All actions
                will be logged for audit purposes.
              </p>
            </div>
          </div>

          {!impersonationToken ? (
            <button
              onClick={handleStartImpersonation}
              disabled={isLoading}
              className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white rounded-lg font-medium transition-colors"
            >
              {isLoading ? "Starting..." : "Start Impersonation"}
            </button>
          ) : (
            <>
              {/* Token Display */}
              <div className="bg-gray-50 rounded-lg p-4">
                <p className="text-sm text-gray-600 mb-2">
                  Impersonation Token
                </p>
                <div className="font-mono text-xs text-gray-900 break-all bg-white p-2 rounded border border-gray-200 mb-3">
                  {impersonationToken}
                </div>
                <button
                  onClick={handleCopyUrl}
                  className="w-full flex items-center justify-center gap-2 px-3 py-2 bg-gray-200 hover:bg-gray-300 text-gray-900 rounded-lg text-sm font-medium transition-colors"
                >
                  {copied ? (
                    <>
                      <Check size={16} />
                      Copied!
                    </>
                  ) : (
                    <>
                      <Copy size={16} />
                      Copy Dashboard URL
                    </>
                  )}
                </button>
              </div>

              {/* Action Buttons */}
              <div className="space-y-2">
                <button
                  onClick={handleOpenDashboard}
                  className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors"
                >
                  Open Dashboard in New Tab
                </button>
              </div>

              {/* Session Info */}
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 text-sm text-blue-900">
                <p className="font-medium mb-1">Session Information</p>
                <p>
                  This impersonation session will expire in 1 hour. You will need
                  to start a new session afterward.
                </p>
              </div>
            </>
          )}

          {/* Close Button */}
          <button
            onClick={onClose}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg font-medium text-gray-900 hover:bg-gray-50 transition-colors"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
}
