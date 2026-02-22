// src/lib/config.ts
//
// Production values should be set via environment variables:
//   NEXT_PUBLIC_LANDING_URL       - e.g. https://landing.example.com
//   NEXT_PUBLIC_DASHBOARD_URL     - e.g. https://dashboard.example.com
//   NEXT_PUBLIC_API_URL           - e.g. https://api.example.com/api/v1
//   NEXT_PUBLIC_RESTAURANT_DOMAIN - e.g. example.com (used for tenant subdomains)
//
// Localhost fallbacks are provided for local development only.

export const getConfig = () => {
  const isDevelopment = process.env.NODE_ENV === 'development';

  // Get origin safely (only in browser)
  const getOrigin = () => {
    if (typeof window === 'undefined') {
      return 'http://localhost:3001'
    }
    return window.location.origin
  }

  return {
    isDevelopment,
    landingUrl: process.env.NEXT_PUBLIC_LANDING_URL || (isDevelopment ? 'http://localhost:3001' : getOrigin()),
    dashboardUrl: process.env.NEXT_PUBLIC_DASHBOARD_URL || (isDevelopment ? 'http://localhost:3002' : getOrigin()),
    apiUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
    restaurantDomain: process.env.NEXT_PUBLIC_RESTAURANT_DOMAIN || 'localhost:3003',
  };
};