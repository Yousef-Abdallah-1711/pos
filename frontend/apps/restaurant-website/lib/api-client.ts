import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from "axios";
import { getApiBaseUrl } from './config/api';
import { debugLog, debugError } from './debug';
import { isCreateOrderRequest, rememberOrderToken, tokenForOrderUrl } from './order-access';

/**
 * API Client with Axios
 * Includes interceptors for public order capabilities and request/response logging
 */
class ApiClient {
  private instance: AxiosInstance;
  private baseURL: string;

  constructor(baseURL?: string) {
    this.baseURL = baseURL || getApiBaseUrl();

    this.instance = axios.create({
      baseURL: this.baseURL,
      timeout: 30000,
    });

    this.setupInterceptors();
  }

  /**
   * Setup request/response interceptors
   */
  private setupInterceptors() {
    // Request interceptor
    this.instance.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        // Remove credentials left by older storefront builds. Public storefront
        // requests use order-specific capabilities, not staff bearer tokens.
        if (typeof window !== "undefined") {
          try {
            window.localStorage.removeItem("auth-token");
          } catch {
            // Storage can be disabled; public requests should still proceed.
          }
        }
        config.headers.delete("Authorization");

        if (config.method && !['get', 'head', 'options'].includes(config.method.toLowerCase())) {
          config.headers['Content-Type'] = 'application/json'
        }

        // Prove ownership of public orders placed from this device (see order-access.ts)
        const orderToken = tokenForOrderUrl(config.url);
        if (orderToken && !config.headers["X-Order-Token"]) {
          config.headers["X-Order-Token"] = orderToken;
        }

        // Add tenant and restaurant headers if available
        const tenantId = typeof window !== "undefined" ? localStorage.getItem("tenant-id") : null;
        const restaurantId = typeof window !== "undefined" ? localStorage.getItem("restaurant-id") : null;

        if (tenantId) {
          config.headers["X-Tenant-ID"] = tenantId;
        }
        if (restaurantId) {
          config.headers["X-Restaurant-ID"] = restaurantId;
        }

        // Log request in development
        if (process.env.NODE_ENV === "development") {
          debugLog(`[API] ${config.method?.toUpperCase()} ${config.url}`, {
            tenantId,
            restaurantId,
          });
        }

        return config;
      },
      (error: AxiosError) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor
    this.instance.interceptors.response.use(
      (response) => {
        // Order creation returns the public tracking token exactly once: keep it for later lookups
        if (isCreateOrderRequest(response.config?.method, response.config?.url)) {
          rememberOrderToken(response.data);
        }

        // Log successful responses in development
        if (process.env.NODE_ENV === "development") {
          debugLog(`[API] Response from ${response.config.url}`, { status: response.status });
        }
        return response;
      },
      (error: AxiosError) => {
        // Public 401 responses pass through without redirecting visitors to a staff login page.
        if (error.response?.status === 403) {
          // Forbidden
          debugError("[API] Forbidden:", { url: error.config?.url });
        }

        if (error.response?.status === 404) {
          // Not found
          debugError("[API] Not found:", error.config?.url);
        }

        if (error.response?.status === 500) {
          // Server error
          debugError("[API] Server error:", { url: error.config?.url });
        }

        // Log error in development
        if (process.env.NODE_ENV === "development") {
          debugError(`[API] Error from ${error.config?.url}`, {
            status: error.response?.status,
            message: error.message,
          });
        }

        return Promise.reject(error);
      }
    );
  }

  /**
   * GET request with AbortController support
   */
  async get<T>(url: string, config?: any): Promise<T> {
    const response = await this.instance.get<T>(url, config);
    return response.data;
  }

  /**
   * GET request with abort signal for cancellation support
   */
  async getCancellable<T>(url: string, signal?: AbortSignal, config?: any): Promise<T> {
    const response = await this.instance.get<T>(url, {
      ...config,
      signal,
    });
    return response.data;
  }

  /**
   * POST request
   */
  async post<T>(url: string, data?: any, config?: any): Promise<T> {
    const response = await this.instance.post<T>(url, data, config);
    return response.data;
  }

  /**
   * PUT request
   */
  async put<T>(url: string, data?: any, config?: any): Promise<T> {
    const response = await this.instance.put<T>(url, data, config);
    return response.data;
  }

  /**
   * PATCH request
   */
  async patch<T>(url: string, data?: any, config?: any): Promise<T> {
    const response = await this.instance.patch<T>(url, data, config);
    return response.data;
  }

  /**
   * DELETE request
   */
  async delete<T>(url: string, config?: any): Promise<T> {
    const response = await this.instance.delete<T>(url, config);
    return response.data;
  }

  /**
   * Get axios instance for advanced usage
   */
  getAxiosInstance(): AxiosInstance {
    return this.instance;
  }
}

// Export singleton instance
export const apiClient = new ApiClient();

export default apiClient;
