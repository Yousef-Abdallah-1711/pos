import axios, { AxiosInstance } from "axios";
import { getSuperAdminApiBaseUrl } from "@/lib/config";
import type {
  ApiResponse,
  ApiError,
  SuperAdminLoginRequest,
  SuperAdminLoginResponse,
  DashboardStats,
  PaymentMetrics,
  WebhookHealthResponse,
  RefundPaymentResponse,
  PaginatedResponse,
  Tenant,
  User,
  AuditLog,
  ImpersonationToken,
} from "@/types";

const API_TIMEOUT_MS = 10000;

type NormalizedClientError = {
  url: string;
  baseURL: string;
  method: string;
  status: number | null;
  code: string | null;
  message: string;
  data: unknown;
};

function normalizeClientError(
  error: unknown,
  fallbackBaseURL: string
): NormalizedClientError {
  const axiosError = axios.isAxiosError<ApiError>(error) ? error : null;
  const isNetworkError = Boolean(axiosError && !axiosError.response);

  let message = "Unknown client error";
  if (error instanceof SyntaxError && error.message.includes("Unexpected end of JSON input")) {
    message = "Backend returned an empty or malformed JSON response.";
  } else if (axiosError) {
    message = isNetworkError
      ? `Network error: backend may be unavailable at ${axiosError.config?.baseURL ?? fallbackBaseURL}`
      : axiosError.message;
  } else if (error instanceof Error) {
    message = error.message;
  }

  return {
    url: axiosError?.config?.url ?? "unknown",
    baseURL: axiosError?.config?.baseURL ?? fallbackBaseURL,
    method: axiosError?.config?.method?.toUpperCase() ?? "GET",
    status: axiosError?.response?.status ?? null,
    code: axiosError?.code ?? (error instanceof Error ? error.name : null),
    message,
    data: axiosError?.response?.data ?? null,
  };
}

class SuperAdminAPI {
  private client: AxiosInstance;

  constructor(baseURL: string = "") {
    const resolvedBaseURL = baseURL || getSuperAdminApiBaseUrl();
    this.client = axios.create({
      baseURL: resolvedBaseURL,
      timeout: API_TIMEOUT_MS,
      headers: {
        "Content-Type": "application/json",
      },
    });

    // Add token to requests
    this.client.interceptors.request.use((config) => {
      if (typeof window !== "undefined") {
        const token = localStorage.getItem("super-admin-token");
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
      }
      if (process.env.NODE_ENV === "development") {
        console.log(`[API] ${config.method?.toUpperCase()} ${config.url}`, {
          baseURL: config.baseURL,
          params: config.params,
        });
      }
      return config;
    });

    // Handle errors
    this.client.interceptors.response.use(
      (response) => {
        if (process.env.NODE_ENV === "development") {
          console.log(`[API] Response:`, {
            url: response.config.url,
            status: response.status,
          });
        }
        return response;
      },
      (error: unknown) => {
        const normalizedError = normalizeClientError(error, resolvedBaseURL);

        if (process.env.NODE_ENV === "development") {
          // Use warn to avoid Next.js dev "Console Error" overlay for handled API failures.
          console.warn("[API] Request failed:", normalizedError);
        }

        const isUnauthorized =
          axios.isAxiosError<ApiError>(error) && error.response?.status === 401;

        if (isUnauthorized) {
          // Clear auth data on 401
          if (typeof window !== "undefined") {
            localStorage.removeItem("super-admin-token");
            localStorage.removeItem("super-admin-auth-store");
            window.location.href = "/login";
          }
        }
        return Promise.reject(error);
      }
    );
  }

  // Auth endpoints
  async login(
    credentials: SuperAdminLoginRequest
  ): Promise<SuperAdminLoginResponse> {
    const response = await this.client.post<ApiResponse<SuperAdminLoginResponse>>(
      "/auth/login",
      credentials
    );
    return response.data.data!;
  }

  async logout(): Promise<void> {
    await this.client.post("/auth/logout");
  }

  // Dashboard endpoints
  async getDashboardStats(): Promise<DashboardStats> {
    const response = await this.client.get<ApiResponse<DashboardStats>>(
      "/dashboard/stats"
    );
    return response.data.data!;
  }

  async getPaymentMetrics(
    period: "today" | "week" | "month" | "all" = "month"
  ): Promise<PaymentMetrics> {
    const response = await this.client.get<ApiResponse<PaymentMetrics>>(
      "/payment-metrics",
      { params: { period } }
    );
    return response.data.data!;
  }

  async getWebhookHealth(
    period: "today" | "week" | "month" | "all" = "month",
    limit: number = 50
  ): Promise<WebhookHealthResponse> {
    const response = await this.client.get<ApiResponse<WebhookHealthResponse>>(
      "/webhook-health",
      { params: { period, limit } }
    );
    return response.data.data!;
  }

  // Tenant endpoints
  async listTenants(
    page: number = 1,
    pageSize: number = 20,
    searchQuery?: string
  ): Promise<PaginatedResponse<Tenant>> {
    const response = await this.client.get(
      "/tenants",
      {
        params: {
          page,
          page_size: pageSize,
          search: searchQuery,
        },
      }
    );
    // Backend returns { success, data: [...tenants], meta: {...} }
    const meta = response.data.meta || {};
    return {
      data: response.data.data || [],
      total: meta.total || 0,
      page: meta.page || page,
      limit: meta.page_size || pageSize,
      pages: meta.total_pages || 1,
    };
  }

  async getTenant(id: number): Promise<Tenant> {
    const response = await this.client.get<ApiResponse<Tenant>>(
      `/tenants/${id}`
    );
    return response.data.data!;
  }

  async getTenantUsers(
    tenantId: number,
    page: number = 1,
    pageSize: number = 20
  ): Promise<PaginatedResponse<User>> {
    const response = await this.client.get(
      `/tenants/${tenantId}/users`,
      {
        params: { page, page_size: pageSize },
      }
    );
    // Backend returns { success, data: [...users], meta: {...} }
    // Transform to match PaginatedResponse<User> structure
    const meta = response.data.meta || {};
    return {
      data: response.data.data || [],
      total: meta.total || 0,
      page: meta.page || page,
      limit: meta.page_size || pageSize,
      pages: meta.total_pages || 1,
    };
  }

  // User endpoints
  async listUsers(
    page: number = 1,
    pageSize: number = 20,
    searchQuery?: string
  ): Promise<PaginatedResponse<User>> {
    const response = await this.client.get(
      "/users",
      {
        params: {
          page,
          page_size: pageSize,
          search: searchQuery,
        },
      }
    );
    // Backend returns { success, data: [...users], meta: {...} }
    const meta = response.data.meta || {};
    return {
      data: response.data.data || [],
      total: meta.total || 0,
      page: meta.page || page,
      limit: meta.page_size || pageSize,
      pages: meta.total_pages || 1,
    };
  }

  async getUser(id: number): Promise<User> {
    const response = await this.client.get<ApiResponse<User>>(
      `/users/${id}`
    );
    return response.data.data!;
  }

  // Impersonation endpoints
  async startImpersonation(userId: number): Promise<ImpersonationToken> {
    const response = await this.client.post<ApiResponse<ImpersonationToken>>(
      `/impersonate/${userId}`
    );
    return response.data.data!;
  }

  async endImpersonation(): Promise<void> {
    await this.client.post("/impersonate/end");
  }

  async validateImpersonationToken(
    token: string
  ): Promise<{ valid: boolean }> {
    const response = await this.client.post<ApiResponse<{ valid: boolean }>>(
      "/impersonate/validate",
      { token }
    );
    return response.data.data!;
  }

  async getImpersonationStatus() {
    const response = await this.client.get(
      "/impersonate/status"
    );
    return response.data.data;
  }

  // Audit log endpoints
  async getAuditLogs(
    page: number = 1,
    pageSize: number = 20,
    action?: string,
    startDate?: string,
    endDate?: string
  ): Promise<PaginatedResponse<AuditLog>> {
    const response = await this.client.get(
      "/audit-logs",
      {
        params: {
          page,
          page_size: pageSize,
          action,
          start_date: startDate,
          end_date: endDate,
        },
      }
    );
    // Backend returns { success, data: [...logs], meta: {...} }
    const meta = response.data.meta || {};
    return {
      data: response.data.data || [],
      total: meta.total || 0,
      page: meta.page || page,
      limit: meta.page_size || pageSize,
      pages: meta.total_pages || 1,
    };
  }

  // Dashboard Login endpoint
  async generateDashboardLoginToken(tenantId: number): Promise<{ token: string; redirect_url: string; expires_in_seconds: number }> {
    const response = await this.client.post<ApiResponse<{ token: string; dashboard_url: string; redirect_url: string; expires_in_seconds: number }>>(`/tenants/${tenantId}/dashboard-login`);
    return response.data.data!;
  }

  // Settings endpoints
  async updateIPWhitelist(ipList: string[]): Promise<void> {
    await this.client.put("/settings/ip-whitelist", {
      ip_whitelist: ipList,
    });
  }

  // Subscription endpoints
  async getTenantSubscription(tenantId: number) {
    const response = await this.client.get<ApiResponse>(`/tenants/${tenantId}/subscription`);
    return response.data.data;
  }

  async extendTenantTrial(
    tenantId: number,
    days: number,
    reason?: string
  ): Promise<{ new_trial_end_date: string; days_added: number; total_trial_days: number }> {
    const response = await this.client.post<ApiResponse>(`/tenants/${tenantId}/extend-trial`, {
      days,
      reason,
    });
    return response.data.data as any;
  }

  async disableTenant(tenantId: number, reason: string): Promise<{ disabled_at: string; previous_status: string }> {
    const response = await this.client.post<ApiResponse>(`/tenants/${tenantId}/disable`, {
      reason,
    });
    return response.data.data as any;
  }

  async enableTenant(tenantId: number): Promise<{ enabled_at: string; restored_status: string }> {
    const response = await this.client.post<ApiResponse>(`/tenants/${tenantId}/enable`);
    return response.data.data as any;
  }

  async deleteTenant(tenantId: number, reason?: string): Promise<{ message: string; deleted_at: string }> {
    const response = await this.client.delete<ApiResponse>(`/tenants/${tenantId}?confirm=true`, {
      data: { reason },
    });
    return response.data.data as any;
  }

  async getSubscriptionHistory(
    tenantId: number,
    page: number = 1,
    pageSize: number = 20
  ): Promise<{ items: any[]; total: number; page: number; page_size: number }> {
    const response = await this.client.get<ApiResponse>(`/tenants/${tenantId}/subscription-history`, {
      params: { page, page_size: pageSize },
    });
    return response.data.data as any;
  }

  async refundPayment(paymentId: number, reason?: string): Promise<RefundPaymentResponse> {
    const response = await this.client.post<ApiResponse<RefundPaymentResponse>>(
      `/payments/${paymentId}/refund`,
      reason ? { reason } : {}
    );
    return response.data.data!;
  }
}

export const superAdminAPI = new SuperAdminAPI();
