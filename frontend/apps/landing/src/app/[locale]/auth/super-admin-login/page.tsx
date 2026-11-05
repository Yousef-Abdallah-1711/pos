"use client";

import React, { useEffect, useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/stores/authStore";
import { api } from "@/lib/api";
import { getConfig } from "@/lib/config";
import { AlertCircle, CheckCircle, Loader2 } from "lucide-react";
import { useLanguage } from "@/contexts/LanguageContext";

export default function SuperAdminLoginPage() {
  const router = useRouter();
  const setAuth = useAuthStore((state) => state.setAuth);
  const { t, dir, language } = useLanguage();
  const locale = language;

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [redirecting, setRedirecting] = useState(false);
  const hasExecuted = useRef(false);

  useEffect(() => {
    const exchangeToken = async () => {
      if (hasExecuted.current) return;
      hasExecuted.current = true;

      try {
        const response = await fetch(`${getConfig().apiUrl}/auth/session`, { credentials: 'include', cache: 'no-store' });
        if (!response.ok) throw new Error(t.auth.superAdmin.noToken);
        const sessionUser: unknown = await response.json();
        setAuth(sessionUser as Parameters<typeof setAuth>[0], '');
        window.history.replaceState({}, document.title, window.location.pathname);
        setSuccess(true);
        setIsLoading(false);
        setRedirecting(true);
        setTimeout(() => {
          window.location.href = `${getConfig().dashboardUrl}/${locale}/dashboard`;
        }, 500);
      } catch (err: unknown) {
        let errorMessage = t.auth.superAdmin.authError;

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
  }, [locale, setAuth, t]);

  return (
    <div dir={dir} className="min-h-screen flex items-center justify-center bg-gradient-to-br from-orange-50 to-red-50 dark:from-gray-900 dark:to-gray-800">
      <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-700 p-8 max-w-md w-full mx-4">
        {isLoading && !error && !success && (
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-br from-orange-500 to-red-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
              <Loader2 className="h-8 w-8 text-white animate-spin" />
            </div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
              {t.auth.superAdmin.accessing}
            </h1>
            <p className="text-gray-600 dark:text-gray-400">
              {t.auth.superAdmin.pleaseWait}
            </p>
          </div>
        )}

        {error && (
          <div className="text-center">
            <div className="mb-4 inline-block">
              <AlertCircle className="text-red-500" size={48} />
            </div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
              {t.auth.superAdmin.accessFailed}
            </h1>
            <p className="text-red-600 dark:text-red-400 mb-6">{error}</p>
            <button
              onClick={() => router.replace(`/${locale}/auth/login`)}
              className="px-6 py-2.5 bg-gradient-to-r from-orange-500 to-red-600 text-white rounded-lg hover:from-orange-600 hover:to-red-700 font-medium"
            >
              {t.auth.superAdmin.goToLogin}
            </button>
          </div>
        )}

        {success && (
          <div className="text-center">
            <div className="mb-4 inline-block">
              <CheckCircle className="text-green-500" size={48} />
            </div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
              {t.auth.superAdmin.accessGranted}
            </h1>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              {t.auth.superAdmin.verified}
            </p>

            {redirecting && (
              <div className="flex items-center justify-center gap-2 text-orange-600 dark:text-orange-400">
                <Loader2 className="animate-spin" size={20} />
                <span>{t.auth.superAdmin.redirecting}</span>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
