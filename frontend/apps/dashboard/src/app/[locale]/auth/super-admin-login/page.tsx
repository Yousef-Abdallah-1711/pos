"use client";

import React, { useEffect, useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/stores/authStore";
import { api } from "@/lib/api";
import { AlertCircle, CheckCircle } from "lucide-react";
import { debugLog, debugError } from '@/lib/debug';

export default function SuperAdminLoginPage() {
  const router = useRouter();
  const setAuth = useAuthStore((state) => state.setAuth);

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const hasExecuted = useRef(false);

  useEffect(() => {
    const exchangeToken = async () => {
      // Prevent double execution
      if (hasExecuted.current) return;
      hasExecuted.current = true;

      try {
        // Get token from URL query parameters using window.location
        const params = new URLSearchParams(window.location.search);
        const token = params.get("token");
        debugLog("[SuperAdminLogin] Component mounted, token from URL:", token ? "present" : "missing");

        if (!token) {
          debugError("[SuperAdminLogin] No token provided");
          setError("No token provided. Invalid access link.");
          setIsLoading(false);
          return;
        }

        // Exchange the one-time token for a regular JWT
        debugLog("[SuperAdminLogin] Exchanging token for JWT...");
        const response = await api.post("/auth/token-exchange", { token });

        debugLog("[SuperAdminLogin] Token exchange response:", response.status);
        const data = response.data;
        debugLog("[SuperAdminLogin] Token exchange successful", data);

        // Store the JWT token and user
        const authData = data.data || data;
        if (authData?.token && authData?.user) {
          debugLog("[SuperAdminLogin] Setting auth with user:", authData.user.email);
          setAuth(authData.user, authData.token);

          // Store in localStorage for persistence
          localStorage.setItem('token', authData.token)
          localStorage.setItem('user', JSON.stringify(authData.user))
          debugLog("[SuperAdminLogin] ✅ Auth stored in localStorage");

          setSuccess(true);

          // Redirect to dashboard after 2 seconds
          setTimeout(() => {
            debugLog("[SuperAdminLogin] Redirecting to dashboard...");
            router.replace("/en/dashboard");
          }, 2000);
        } else {
          debugError("[SuperAdminLogin] Missing token or user in response:", authData);
          throw new Error("No JWT token or user received");
        }
      } catch (err) {
        debugError("[SuperAdminLogin] Error:", err);
        setError(err instanceof Error ? err.message : "An error occurred");
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
              onClick={() => router.replace("/en/login")}
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
              Welcome!
            </h1>
            <p className="text-gray-600 mb-2">
              Your access has been granted.
            </p>
            <p className="text-gray-600 text-sm">
              Redirecting to dashboard...
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
