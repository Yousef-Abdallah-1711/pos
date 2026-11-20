"use client";

import React, { useEffect, useState, ReactNode } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAuthStore } from "@/stores/authStore";
import { superAdminAPI } from "@/lib/api";
import { useSuperAdminText } from "@/components/providers/PreferencesProvider";

interface SuperAdminAuthProviderProps {
  children: ReactNode;
}

const PUBLIC_ROUTES = ["/login", "/"];
const PROTECTED_ROUTES = ["/dashboard"];

export function SuperAdminAuthProvider({ children }: SuperAdminAuthProviderProps) {
  const t = useSuperAdminText();
  const router = useRouter();
  const pathname = usePathname();
  const { isLoading, setLoading } = useAuthStore();
  const [sessionUnavailable, setSessionUnavailable] = useState(false);
  const [retryKey, setRetryKey] = useState(0);

  useEffect(() => {
    let active = true;
    let outcome: "authenticated" | "unauthorized" | "unavailable" = "unavailable";
    setLoading(true);
    setSessionUnavailable(false);
    localStorage.removeItem("super-admin-token");

    superAdminAPI.getSession().then((session) => {
      if (active) {
        outcome = "authenticated";
        useAuthStore.setState({
          token: "cookie-session",
          superAdmin: { ...session.super_admin, email: session.email },
        });
      }
    }).catch((error: unknown) => {
      if (!active) return;
      const status = typeof error === "object" && error !== null && "response" in error
        ? (error as { response?: { status?: unknown } }).response?.status
        : undefined;
      if (status === 401) {
        outcome = "unauthorized";
        useAuthStore.setState({ token: null, superAdmin: null });
      }
    }).finally(() => {
      if (!active) return;
      setSessionUnavailable(outcome === "unavailable");
      const hasSession = useAuthStore.getState().isAuthenticated();
      const isPublicRoute = PUBLIC_ROUTES.includes(pathname);
      const isProtectedRoute = PROTECTED_ROUTES.some((route) => pathname.startsWith(route));

      if (outcome === "unauthorized" && isProtectedRoute && !hasSession) router.push("/login");
      else if (outcome === "authenticated" && isPublicRoute && hasSession && pathname === "/login") router.push("/dashboard");
      setLoading(false);
    });

    return () => { active = false; };
  }, [pathname, retryKey, router, setLoading]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
          <p className="mt-4 text-gray-600">{t("Loading...")}</p>
        </div>
      </div>
    );
  }

  if (sessionUnavailable) {
    return (
      <div role="alert" className="flex min-h-screen flex-col items-center justify-center gap-4 text-center">
        <p>{t("Session verification is temporarily unavailable.")}</p>
        <button type="button" onClick={() => setRetryKey((value) => value + 1)}>
          {t("Retry")}
        </button>
      </div>
    );
  }

  return <>{children}</>;
}
