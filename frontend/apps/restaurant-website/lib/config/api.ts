import { debugWarn } from '../debug';

/**
 * Centralized API Configuration
 * All API endpoints and configuration should be defined here
 */

/**
 * Get the base API URL from environment variables
 * Falls back to localhost development server if not configured
 */
export function getApiBaseUrl(): string {
  if (typeof window === 'undefined') {
    // Server-side
    return process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
  }

  // Client-side
  return (
    process.env.NEXT_PUBLIC_API_URL ||
    process.env.NEXT_PUBLIC_API_BASE_URL ||
    'http://localhost:8080/api/v1'
  );
}

/**
 * API timeout in milliseconds
 */
export const API_TIMEOUT = parseInt(process.env.NEXT_PUBLIC_API_TIMEOUT || '5000', 10);

/**
 * API Endpoints
 */
export const API_ENDPOINTS = {
  // Products
  PRODUCTS: (slug: string) => `/public/restaurants/${slug}/products`,
  CATEGORIES: (slug: string) => `/public/restaurants/${slug}/categories`,
  SEARCH_PRODUCTS: (slug: string) => `/public/restaurants/${slug}/products/search`,

  // Orders
  CREATE_ORDER: '/orders/create',
  GET_ORDER: (orderId: string) => `/orders/${orderId}`,
  GET_ORDERS: '/orders',

  // Themes
  THEME: (slug: string) => `/public/restaurants/${slug}/theme`,
  THEMES_LIST: '/themes',
} as const;

/**
 * Build full API URL with endpoint
 */
export function buildApiUrl(endpoint: string, baseUrl?: string): string {
  const base = baseUrl || getApiBaseUrl();
  return `${base}${endpoint}`;
}

/**
 * Validate that required environment variables are set
 */
export function validateApiConfig(): void {
  const apiUrl = getApiBaseUrl();

  if (!apiUrl) {
    debugWarn(
      'API_URL not configured. Set NEXT_PUBLIC_API_URL environment variable. ' +
      'Falling back to localhost:8080/api/v1'
    );
  }
}

// Validate on module load
if (typeof window !== 'undefined') {
  validateApiConfig();
}
