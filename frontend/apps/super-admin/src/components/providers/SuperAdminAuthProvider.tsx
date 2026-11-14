"use client";

import React, { useEffect, ReactNode } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAuthStore } from "@/stores/authStore";
import { superAdminAPI } from "@/lib/api";

interface SuperAdminAuthProviderProps {
  children: ReactNode;
}

const PUBLIC_ROUTES = ["/login", "/"];
const PROTECTED_ROUTES = ["/dashboard"];

export function SuperAdminAuthProvider({ children }: SuperAdminAuthProviderProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { isLoading, setLoading } = useAuthStore();

  useEffect(() => {
    let active = true;
    setLoading(true);
    localStorage.removeItem("super-admin-token");

    superAdminAPI.getSession().then((session) => {
      if (active) {
        useAuthStore.setState({
          token: "cookie-session",
          superAdmin: { ...session.super_admin, email: session.email },
        });
      }
    }).catch(() => {
      if (active) useAuthStore.setState({ token: null });
    }).finally(() => {
      if (!active) return;
      const hasSession = useAuthStore.getState().isAuthenticated();
      const isPublicRoute = PUBLIC_ROUTES.includes(pathname);
      const isProtectedRoute = PROTECTED_ROUTES.some((route) => pathname.startsWith(route));

      if (isProtectedRoute && !hasSession) router.push("/login");
      else if (isPublicRoute && hasSession && pathname === "/login") router.push("/dashboard");
      setLoading(false);
    });

    return () => { active = false; };
  }, [pathname, router, setLoading]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return <>{children}</>;
}
