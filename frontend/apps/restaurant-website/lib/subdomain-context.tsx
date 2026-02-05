"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { parseSubdomain } from "./subdomain";
import type { SubdomainContext } from "./subdomain";

interface SubdomainContextType extends SubdomainContext {
  loading: boolean;
  error: string | null;
}

const SubdomainContext = createContext<SubdomainContextType | undefined>(undefined);

/**
 * Provider component to make subdomain context available to child components
 * Reads subdomain info from headers (server-side) and cookies (client-side)
 */
export function SubdomainProvider({ children }: { children: React.ReactNode }) {
  const [context, setContext] = useState<SubdomainContextType>({
    subdomain: null,
    isRestaurantSite: false,
    isMainSite: true,
    slug: null,
    host: "",
    loading: true,
    error: null,
  });

  useEffect(() => {
    async function initContext() {
      try {
        // Get host from window.location
        const host = typeof window !== "undefined" ? window.location.host : "";
        const parsed = parseSubdomain(host);

        // Get subdomain info from cookies set by middleware (fallback to host parsing)
        const cookieSlug = getCookie("restaurant-slug");
        const cookieIsRestaurantSite = getCookie("is-restaurant-site") === "true";
        const restaurantSlug = cookieSlug || parsed.slug;
        const isRestaurantSite = cookieSlug ? cookieIsRestaurantSite : parsed.isRestaurantSite;

        // If we have a restaurant slug, fetch the restaurant info to get IDs
        if (restaurantSlug && isRestaurantSite && typeof window !== "undefined") {
          try {
            const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";
            const response = await fetch(`${apiBase}/public/restaurants/${restaurantSlug}/info`);
            if (response.ok) {
              const data = await response.json();
              // Store tenant and restaurant IDs for API client to use
              // The /info endpoint returns fields directly: { tenant_id, id, name, ... }
              if (data?.tenant_id) {
                localStorage.setItem("tenant-id", String(data.tenant_id));
              }
              if (data?.id) {
                localStorage.setItem("restaurant-id", String(data.id));
              }
            }
          } catch (fetchError) {
            console.warn("[SubdomainProvider] Failed to fetch restaurant info:", fetchError);
          }
        }

        setContext({
          subdomain: restaurantSlug || null,
          isRestaurantSite,
          isMainSite: !isRestaurantSite,
          slug: restaurantSlug || null,
          host,
          loading: false,
          error: null,
        });
      } catch (error) {
        setContext((prev) => ({
          ...prev,
          loading: false,
          error: error instanceof Error ? error.message : "Failed to parse subdomain",
        }));
      }
    }

    initContext();
  }, []);

  return (
    <SubdomainContext.Provider value={context}>
      {children}
    </SubdomainContext.Provider>
  );
}

/**
 * Hook to access subdomain context in components
 */
export function useSubdomain(): SubdomainContextType {
  const context = useContext(SubdomainContext);
  if (context === undefined) {
    throw new Error("useSubdomain must be used within SubdomainProvider");
  }
  return context;
}

/**
 * Helper to get cookie value
 */
function getCookie(name: string): string | null {
  if (typeof document === "undefined") return null;

  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) {
    return parts.pop()?.split(";").shift() || null;
  }
  return null;
}
