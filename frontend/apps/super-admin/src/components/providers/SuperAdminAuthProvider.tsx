"use client";

import React, { useEffect, ReactNode } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAuthStore } from "@/stores/authStore";

interface SuperAdminAuthProviderProps {
  children: ReactNode;
}

const PUBLIC_ROUTES = ["/login", "/"];
const PROTECTED_ROUTES = ["/dashboard"];

export function SuperAdminAuthProvider({ children }: SuperAdminAuthProviderProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { token, isLoading, setLoading } = useAuthStore();
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated());

  useEffect(() => {
    setLoading(true);

    // Check if we have a stored token
    const storedToken = localStorage.getItem("super-admin-token");

    if (storedToken && !token) {
      // Token in storage but not in state - restore it
      useAuthStore.setState({ token: storedToken });
    }

    // Determine if route protection is needed
    const isPublicRoute = PUBLIC_ROUTES.includes(pathname);
    const isProtectedRoute = PROTECTED_ROUTES.some((route) =>
      pathname.startsWith(route)
    );

    if (isProtectedRoute && !isAuthenticated && !storedToken) {
      // Trying to access protected route without auth
      router.push("/login");
    } else if (isPublicRoute && isAuthenticated) {
      // Already authenticated, redirect from login to dashboard
      if (pathname === "/login") {
        router.push("/dashboard");
      }
    }

    setLoading(false);
  }, [pathname, token, isAuthenticated, router, setLoading]);

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
