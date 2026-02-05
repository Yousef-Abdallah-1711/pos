'use client'
// src/lib/config.ts
export const getConfig = () => {
  // In production, we'll use relative URLs or same-domain approach
  // In development, we'll use the appropriate ports
  const isDevelopment = process.env.NODE_ENV === 'development';

  // Get origin safely (only in browser)
  const getOrigin = () => {
    if (typeof window === 'undefined') {
      return 'http://localhost:3002'
    }
    return window.location.origin
  }

  return {
    isDevelopment,
    landingUrl: isDevelopment
      ? 'http://localhost:3001'
      : process.env.NEXT_PUBLIC_LANDING_URL || getOrigin(),
    dashboardUrl: isDevelopment
      ? 'http://localhost:3002'
      : process.env.NEXT_PUBLIC_DASHBOARD_URL || getOrigin(),
    apiUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  };
};