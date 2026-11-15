"use client";

import React, { useState } from "react";
import { isAxiosError } from "axios";
import { useAuthStore } from "@/stores/authStore";
import { superAdminAPI } from "@/lib/api";
import { toast } from "sonner";
import { PreferenceControls } from "@/components/layout/PreferenceControls";
import { useSuperAdminText } from "@/components/providers/PreferencesProvider";

export default function LoginPage() {
  const t = useSuperAdminText();
  const { setToken, setSuperAdmin, setError: setAuthError } =
    useAuthStore();
  const [email, setEmail] = useState("superadmin@pos-saas.com");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      if (!email || !password) {
        setError(t("Please enter both email and password"));
        return;
      }

      const response = await superAdminAPI.login({
        email,
        password,
      });

      setToken("cookie-session");
      setSuperAdmin(response.super_admin);

      toast.success(t("Login successful!"));

      // Start a new app document after establishing the cookie-backed session.
      window.location.assign("/dashboard");
    } catch (err: unknown) {
      const serverError = isAxiosError<{ error?: unknown }>(err)
        ? err.response?.data?.error
        : undefined;
      const errorMessage = typeof serverError === "string"
        ? serverError
        : t("Failed to login. Please check your credentials.");
      setError(String(errorMessage));
      setAuthError(String(errorMessage));
      toast.error(String(errorMessage));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="relative flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-950 via-slate-900 to-violet-950 px-4 py-20 sm:px-6 lg:px-8">
      <div className="absolute right-4 top-4 sm:right-6 sm:top-6">
        <PreferenceControls compact />
      </div>
      <div className="w-full max-w-md space-y-8">
        {/* Header */}
        <div className="text-center">
          <h1 className="text-4xl font-bold text-white mb-2">
            SERVIO
          </h1>
          <p className="text-lg text-slate-300">
            {t("Super Admin Dashboard")}
          </p>
          <p className="text-sm text-slate-300 mt-2">
            {t("Secure access for platform administrators")}
          </p>
        </div>

        {/* Login Card */}
        <div className="bg-white rounded-lg shadow-xl p-8 space-y-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Email Field */}
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                {t("Email Address")}
              </label>
              <input
                type="email"
                id="email"
                dir="ltr"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="superadmin@pos-saas.com"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
                disabled={isLoading}
              />
            </div>

            {/* Password Field */}
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
                {t("Password")}
              </label>
              <input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder={t("Enter your password")}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
                disabled={isLoading}
              />
            </div>

            {/* Error Message */}
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
                {error}
              </div>
            )}

            {/* Login Button */}
            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg transition duration-200"
            >
              {isLoading ? t("Logging in...") : t("Login")}
            </button>
          </form>

        </div>

        {/* Footer */}
        <div className="text-center text-sm text-slate-400">
          <p>© 2024 POS SaaS. All rights reserved.</p>
          <p className="mt-1">
            This is a secure area. Unauthorized access is prohibited.
          </p>
        </div>
      </div>
    </div>
  );
}
