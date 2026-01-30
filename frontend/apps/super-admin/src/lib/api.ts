import axios, { AxiosInstance, AxiosError } from "axios";
import type {
  ApiResponse,
  ApiError,
  SuperAdminLoginRequest,
  SuperAdminLoginResponse,
  DashboardStats,
  PaginatedResponse,
  Tenant,
  User,
  AuditLog,
  ImpersonationToken,
} from "@/types";

class SuperAdminAPI {
  private client: AxiosInstance;

  constructor(baseURL: string = process.env.NEXT_PUBLIC_SUPER_ADMIN_API_URL) {
    this.client = axios.create({
      baseURL,
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
      return config;
    });

    // Handle errors
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError<ApiError>) => {
        if (error.response?.status === 401) {
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

  // Tenant endpoints
  async listTenants(
    page: number = 1,
    limit: number = 20,
    searchQuery?: string
  ): Promise<PaginatedResponse<Tenant>> {
    const response = await this.client.get<ApiResponse<PaginatedResponse<Tenant>>>(
      "/tenants",
      {
        params: {
          page,
          limit,
          search: searchQuery,
        },
      }
    );
    return response.data.data!;
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
    limit: number = 20
  ): Promise<PaginatedResponse<User>> {
    const response = await this.client.get<ApiResponse<PaginatedResponse<User>>>(
      `/tenants/${tenantId}/users`,
      {
        params: { page, limit },
      }
    );
    return response.data.data!;
  }

  async disableTenant(id: number): Promise<Tenant> {
    const response = await this.client.put<ApiResponse<Tenant>>(
      `/tenants/${id}/disable`
    );
    return response.data.data!;
  }

  async enableTenant(id: number): Promise<Tenant> {
    const response = await this.client.put<ApiResponse<Tenant>>(
      `/tenants/${id}/enable`
    );
    return response.data.data!;
  }

  // User endpoints
  async listUsers(
    page: number = 1,
    limit: number = 20,
    searchQuery?: string
  ): Promise<PaginatedResponse<User>> {
    const response = await this.client.get<ApiResponse<PaginatedResponse<User>>>(
      "/users",
      {
        params: {
          page,
          limit,
          search: searchQuery,
        },
      }
    );
    return response.data.data!;
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
    limit: number = 20,
    action?: string,
    startDate?: string,
    endDate?: string
  ): Promise<PaginatedResponse<AuditLog>> {
    const response = await this.client.get<ApiResponse<PaginatedResponse<AuditLog>>>(
      "/audit-logs",
      {
        params: {
          page,
          limit,
          action,
          start_date: startDate,
          end_date: endDate,
        },
      }
    );
    return response.data.data!;
  }

  // Settings endpoints
  async updateIPWhitelist(ipList: string[]): Promise<void> {
    await this.client.put("/settings/ip-whitelist", {
      ip_whitelist: ipList,
    });
  }
}

export const superAdminAPI = new SuperAdminAPI();
