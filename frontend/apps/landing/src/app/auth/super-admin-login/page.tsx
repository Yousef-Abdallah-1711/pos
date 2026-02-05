"use client";

import React, { useEffect, useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/stores/authStore";
import { api } from "@/lib/api";
import { AlertCircle, CheckCircle, Loader2 } from "lucide-react";

export default function SuperAdminLoginPage() {
  const router = useRouter();
  const setAuth = useAuthStore((state) => state.setAuth);

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [redirecting, setRedirecting] = useState(false);
  const hasExecuted = useRef(false);

  useEffect(() => {
    const exchangeToken = async () => {
      // Prevent double execution
      if (hasExecuted.current) return;
      hasExecuted.current = true;

      try {
        // Get token from URL query parameters
        const params = new URLSearchParams(window.location.search);
        const token = params.get("token");

        if (!token) {
          setError("No token provided. Invalid access link.");
          setIsLoading(false);
          return;
        }

        // Exchange the one-time token for a regular JWT
        const response = await api.post("/auth/token-exchange", { token });
        const data = response.data;

        // Store the JWT token and user
        const authData = data.data || data;
        if (authData?.token && authData?.user) {
          setAuth(authData.user, authData.token);

          // Build dashboard URL with auth params
          const userParam = encodeURIComponent(JSON.stringify(authData.user));
          const tokenParam = encodeURIComponent(authData.token);
          const dashboardUrl = `http://localhost:3002/en/dashboard?auth_token=${tokenParam}&auth_user=${userParam}`;

          setSuccess(true);
          setIsLoading(false);

          // Auto-redirect to dashboard after a short delay
          setRedirecting(true);
          setTimeout(() => {
            window.location.href = dashboardUrl;
          }, 1500);
        } else {
          throw new Error("No JWT token or user received");
        }
      } catch (err: unknown) {
        // Extract error message from API response or error object
        let errorMessage = "An error occurred during authentication";

        if (err && typeof err === 'object') {
          const axiosError = err as { response?: { data?: { error?: string } }; message?: string };
          if (axiosError.response?.data?.error) {
            errorMessage = axiosError.response.data.error;
          } else if (axiosError.message) {
            errorMessage = axiosError.message;
          }
        }

        setError(errorMessage);
        setIsLoading(false);
      }
    };

    exchangeToken();
  }, [setAuth, router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="bg-white rounded-lg shadow-xl p-8 max-w-md w-full">
        {isLoading && !error && !success && (
          <div className="text-center">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">
              Accessing Dashboard
            </h1>
            <p className="text-gray-600">
              Please wait while we set up your access...
            </p>
          </div>
        )}

        {error && (
          <div className="text-center">
            <div className="mb-4 inline-block">
              <AlertCircle className="text-red-500" size={48} />
            </div>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">
              Access Failed
            </h1>
            <p className="text-red-600 mb-6">{error}</p>
            <button
              onClick={() => router.replace("/auth/login")}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              Go to Login
            </button>
          </div>
        )}

        {success && (
          <div className="text-center">
            <div className="mb-4 inline-block">
              <CheckCircle className="text-green-500" size={48} />
            </div>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">
              Access Granted!
            </h1>
            <p className="text-gray-600 mb-4">
              Your super admin access has been verified.
            </p>

            {redirecting && (
              <div className="flex items-center justify-center gap-2 text-blue-600">
                <Loader2 className="animate-spin" size={20} />
                <span>Redirecting to dashboard...</span>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
