import { debugWarn } from '../debug';

/**
 * Centralized API Configuration
 * All API endpoints and configuration should be defined here
 */

const DEV_API_BASE_URL = 'http://localhost:8080/api/v1';
let hasWarnedAboutDevFallback = false;

function normalizeBaseUrl(url: string): string {
  return url.replace(/\/+$/, '');
}

function getConfiguredApiBaseUrl(): string | null {
  const raw =
    process.env.NEXT_PUBLIC_API_BASE_URL ||
    process.env.NEXT_PUBLIC_API_URL ||
    '';

  const trimmed = raw.trim();
  if (!trimmed) {
    return null;
  }

  return normalizeBaseUrl(trimmed);
}

/** Resolve an optional server-only URL for Server Components and route handlers. */
export function getServerApiBaseUrl(): string {
  const internalUrl = process.env.API_INTERNAL_URL?.trim();
  return internalUrl ? normalizeBaseUrl(internalUrl) : getApiBaseUrl();
}

export function isProductionRuntime(): boolean {
  return process.env.NODE_ENV === 'production';
}

/**
 * Resolve the base API URL.
 * - Prefers env configuration
 * - Allows localhost fallback in development only
 * - Throws in production when missing to avoid silently calling localhost
 */
export function getApiBaseUrl(): string {
  const configured = getConfiguredApiBaseUrl();
  if (configured) {
    return configured;
  }

  if (!isProductionRuntime()) {
    if (!hasWarnedAboutDevFallback) {
      hasWarnedAboutDevFallback = true;
      debugWarn(
        '[config/api] NEXT_PUBLIC_API_BASE_URL (or NEXT_PUBLIC_API_URL) is not set. ' +
        `Using development fallback: ${DEV_API_BASE_URL}`
      );
    }
    return DEV_API_BASE_URL;
  }

  throw new Error(
    '[config/api] Missing required API base URL in production. ' +
    'Set NEXT_PUBLIC_API_BASE_URL (preferred) or NEXT_PUBLIC_API_URL.'
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
  getApiBaseUrl();
}
