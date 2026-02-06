// Super Admin Authentication
export interface SuperAdminClaims {
  user_id: number;
  email: string;
  super_admin_id: number;
  is_super_admin: boolean;
  exp: number;
}

export interface SuperAdminLoginRequest {
  email: string;
  password: string;
}

export interface SuperAdminLoginResponse {
  token: string;
  expires_at: string;
  super_admin: SuperAdminUser;
}

export interface SuperAdminUser {
  id: number;
  user_id: number;
  email: string;
  secret_key: string;
  ip_whitelist: string[];
  is_active: boolean;
  last_login_at: string | null;
  failed_login_attempts: number;
  locked_until: string | null;
  created_at: string;
  updated_at: string;
}

// Tenant Management
export interface Tenant {
  id: number;
  name: string;
  slug: string;
  created_at: string;
  updated_at: string;
}

export interface TenantWithStats extends Tenant {
  user_count: number;
  order_count: number;
  revenue: number;
  is_active: boolean;
}

// User Management
export interface User {
  id: number;
  tenant_id: number;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  is_super_admin: boolean;
  created_at: string;
  updated_at: string;
}

export interface UserWithTenant extends User {
  tenant_name: string;
}

// Impersonation
export interface ImpersonationSession {
  id: number;
  super_admin_id: number;
  target_user_id: number;
  target_tenant_id: number;
  token: string;
  expires_at: string;
  created_at: string;
}

export interface ImpersonationToken {
  token: string;
  expires_at: string;
  target_user: User;
  target_tenant: Tenant;
}

export interface ImpersonationStatus {
  is_impersonating: boolean;
  target_user: User | null;
  target_tenant: Tenant | null;
  expires_at: string | null;
}

// Audit Logs
export interface AuditLog {
  id: number;
  super_admin_id: number;
  action: string;
  target_user_id: number | null;
  target_tenant_id: number | null;
  ip_address: string;
  user_agent: string;
  details: Record<string, unknown>;
  created_at: string;
}

// Dashboard Stats
export interface DashboardStats {
  total_tenants: number;
  total_users: number;
  active_impersonations: number;
  total_audit_logs: number;
  recent_activity: AuditLog[];
}

// Pagination
export interface PaginationParams {
  page: number;
  limit: number;
  sort_by?: string;
  sort_order?: "asc" | "desc";
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  pages: number;
}

// API Responses
export interface ApiResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface ApiError {
  success: false;
  error: string;
  message?: string;
}

// Subscription Management
export type SubscriptionStatus = 'trial' | 'active' | 'expired' | 'suspended' | 'disabled';

export interface TenantSubscriptionInfo {
  tenant_id: number;
  tenant_name: string;
  slug: string;
  subscription_status: SubscriptionStatus;
  trial_start_date: string | null;
  trial_end_date: string | null;
  trial_days: number;
  days_remaining: number;
  is_disabled: boolean;
  disabled_at: string | null;
  disabled_reason: string | null;
  can_access: boolean;
  block_reason: string | null;
}

export interface SubscriptionHistoryItem {
  id: number;
  tenant_id: number;
  action: string;
  performed_by: number | null;
  old_status: string | null;
  new_status: string | null;
  old_trial_end_date: string | null;
  new_trial_end_date: string | null;
  days_added: number | null;
  reason: string | null;
  ip_address: string | null;
  user_agent: string | null;
  metadata: Record<string, unknown> | null;
  created_at: string;
}

export interface SubscriptionHistoryResponse {
  items: SubscriptionHistoryItem[];
  total: number;
  page: number;
  page_size: number;
}
