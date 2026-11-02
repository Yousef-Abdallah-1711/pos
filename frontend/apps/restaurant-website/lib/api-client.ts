import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from "axios";
import { getApiBaseUrl } from './config/api';
import { debugLog, debugError } from './debug';
import { isCreateOrderRequest, rememberOrderToken, tokenForOrderUrl } from './order-access';

/**
 * API Client with Axios
 * Includes interceptors for error handling, auth, and request/response logging
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
        // Add auth token if available
        const token = typeof window !== "undefined" ? localStorage.getItem("auth-token") : null;
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }

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
            headers: config.headers
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
          debugLog(`[API] Response from ${response.config.url}`, response.data);
        }
        return response;
      },
      (error: AxiosError) => {
        // Handle specific error cases
        if (error.response?.status === 401) {
          // Unauthorized - clear auth and redirect to login
          if (typeof window !== "undefined") {
            localStorage.removeItem("auth-token");
            window.location.href = "/login";
          }
        }

        if (error.response?.status === 403) {
          // Forbidden
          debugError("[API] Forbidden:", error.response.data);
        }

        if (error.response?.status === 404) {
          // Not found
          debugError("[API] Not found:", error.config?.url);
        }

        if (error.response?.status === 500) {
          // Server error
          debugError("[API] Server error:", error.response.data);
        }

        // Log error in development
        if (process.env.NODE_ENV === "development") {
          debugError(`[API] Error from ${error.config?.url}`, {
            status: error.response?.status,
            data: error.response?.data,
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
   * Set authorization token
   */
  setAuthToken(token: string) {
    if (typeof window !== "undefined") {
      localStorage.setItem("auth-token", token);
    }
    this.instance.defaults.headers.common["Authorization"] = `Bearer ${token}`;
  }

  /**
   * Clear authorization token
   */
  clearAuthToken() {
    if (typeof window !== "undefined") {
      localStorage.removeItem("auth-token");
    }
    delete this.instance.defaults.headers.common["Authorization"];
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
