'use client'
import axios from 'axios'
import { getConfig } from './config'
import { debugError, debugLog, debugWarn } from '@/lib/debug'
import { recordPerformanceMetric } from '@/lib/performance'

const API_URL = getConfig().apiUrl
const DEFAULT_API_TIMEOUT_MS = 6000
const rawApiTimeout = Number.parseInt(process.env.NEXT_PUBLIC_API_TIMEOUT_MS || '', 10)
const API_TIMEOUT_MS = Number.isFinite(rawApiTimeout) && rawApiTimeout > 0
  ? rawApiTimeout
  : DEFAULT_API_TIMEOUT_MS
const AUTH_CACHE_TTL_MS = 1500
const REQUEST_DEDUP_WINDOW_MS = 100

type AuthContextSnapshot = {
  token: string | null
  user: any
  tenantId: number | null
  restaurantId: number | null
  cachedAt: number
}

let authContextCache: AuthContextSnapshot | null = null
const pendingRequests = new Map<string, { promise: Promise<any>; startedAt: number }>()

const buildApiMetricKey = (url?: string) => {
  if (!url) {
    return 'api.unknown'
  }

  const normalized = url
    .replace(/\d+/g, ':id')
    .replace(/[?].*$/, '')
    .replace(/^\/+/, '')

  return `api.${normalized || 'root'}`
}

export const api = axios.create({
  baseURL: API_URL,
  timeout: API_TIMEOUT_MS,
  headers: {
    'Content-Type': 'application/json',
  },
})

const buildDedupKey = (config: any) => {
  const method = String(config?.method || 'get').toLowerCase()
  if (method !== 'get') {
    return null
  }
  const url = String(config?.url || '')
  const params = config?.params ? JSON.stringify(config.params) : ''
  const data = config?.data ? (typeof config.data === 'string' ? config.data : JSON.stringify(config.data)) : ''
  return `${method}:${url}:${params}:${data}`
}

const originalRequest = api.request.bind(api)
api.request = ((config: any) => {
  const key = buildDedupKey(config)
  if (!key) {
    return originalRequest(config)
  }
  const existing = pendingRequests.get(key)
  const now = Date.now()
  if (existing && now - existing.startedAt <= REQUEST_DEDUP_WINDOW_MS) {
    return existing.promise
  }

  const promise = originalRequest(config)
    .finally(() => {
      pendingRequests.delete(key)
    })
  pendingRequests.set(key, { promise, startedAt: now })
  return promise
}) as typeof api.request

function readAuthContext(): AuthContextSnapshot {
  const now = Date.now()
  if (authContextCache && now-authContextCache.cachedAt < AUTH_CACHE_TTL_MS) {
    return authContextCache
  }

  let token: string | null = null
  let user: any = null
  let tenantId: number | null = null
  let restaurantId: number | null = null

  try {
    const authStorage = localStorage.getItem('auth-storage')
    if (authStorage) {
      const parsed = JSON.parse(authStorage)
      token = parsed.state?.token || parsed.token
      user = parsed.state?.user || parsed.user
      tenantId = user?.tenant_id || user?.tenantId || null
      restaurantId = user?.restaurant_id || user?.restaurantId || null
    }
  } catch (e) {
    debugError('[API] Failed to parse auth-storage:', e)
  }

  if (!token) {
    token = localStorage.getItem('token')
  }

  if (!user) {
    try {
      const userStr = localStorage.getItem('user')
      if (userStr) {
        user = JSON.parse(userStr)
        tenantId = user?.tenant_id || user?.tenantId
        restaurantId = user?.restaurant_id || user?.restaurantId
      }
    } catch (e) {
      debugError('[API] Failed to parse user from localStorage:', e)
    }
  }

  authContextCache = {
    token,
    user,
    tenantId: tenantId ?? null,
    restaurantId: restaurantId ?? null,
    cachedAt: now,
  }
  return authContextCache
}

// Add token and tenant context to requests
api.interceptors.request.use((config) => {
  const timedConfig = config as typeof config & {
    __perf?: {
      startedAt: number
      metricKey: string
    }
  }

  if (typeof config.timeout !== 'number' || config.timeout <= 0) {
    config.timeout = API_TIMEOUT_MS
  }
  if (typeof window !== 'undefined' && window.performance) {
    timedConfig.__perf = {
      startedAt: window.performance.now(),
      metricKey: buildApiMetricKey(config.url),
    }
  }

  const { token, tenantId, restaurantId } = readAuthContext()

  // Don't add token to auth endpoints (like login, register)
  const isAuthEndpoint = config.url?.includes('/auth/')

  // Check for impersonation token
  let impersonationToken = null
  try {
    impersonationToken = sessionStorage.getItem('impersonation_token')
  } catch (e) {
    // sessionStorage might not be available in some contexts
  }

  debugLog('[API] Request interceptor - url:', config.url, 'impersonationToken:', !!impersonationToken, 'token:', !!token)

  if (impersonationToken && !isAuthEndpoint) {
    // If using impersonation token, add it as a custom header
    config.headers['X-Impersonation-Token'] = impersonationToken
    debugLog('[API] Added X-Impersonation-Token header')

    // Also add tenant context from user if available
    if (tenantId) {
      config.headers['X-Tenant-ID'] = tenantId.toString()
    }
    if (restaurantId) {
      config.headers['X-Restaurant-ID'] = restaurantId.toString()
    }
  } else if (token && !isAuthEndpoint) {
    config.headers.Authorization = `Bearer ${token}`
    debugLog('[API] Added Authorization header')

    // Add tenant context headers from user object
    if (tenantId) {
      config.headers['X-Tenant-ID'] = tenantId.toString()
    } else {
      debugWarn('[API] tenant_id is missing or zero - backend may reject this request')
    }

    if (restaurantId) {
      config.headers['X-Restaurant-ID'] = restaurantId.toString()
    } else {
      debugWarn('[API] restaurant_id is missing or zero - backend may reject this request')
    }
  }

  // If sending FormData, let the browser set the Content-Type with boundary
  if (config.data instanceof FormData) {
    delete config.headers['Content-Type']
  }

  return config
})

// Handle 401 errors
api.interceptors.response.use(
  (response) => {
    const timedConfig = response.config as typeof response.config & {
      __perf?: {
        startedAt: number
        metricKey: string
      }
    }

    if (typeof window !== 'undefined' && timedConfig.__perf?.startedAt) {
      recordPerformanceMetric(
        timedConfig.__perf.metricKey,
        window.performance.now() - timedConfig.__perf.startedAt
      )
    }

    return response
  },
  (error) => {
    const timedConfig = error?.config as
      | (typeof error.config & {
          __perf?: {
            startedAt: number
            metricKey: string
          }
        })
      | undefined

    if (typeof window !== 'undefined' && timedConfig?.__perf?.startedAt) {
      recordPerformanceMetric(
        timedConfig.__perf.metricKey,
        window.performance.now() - timedConfig.__perf.startedAt
      )
    }

    debugError('[API] Request failed:', error)

    // Check if it's a "Missing tenant information" error
    const errorMessage = error.response?.data?.error || ''
    const isTenantContextError = errorMessage.includes('Missing tenant') || errorMessage.includes('tenant')

    if (isTenantContextError) {
      debugWarn('[API] Tenant context error:', errorMessage)
    }

    // Only logout on actual authentication failures, not tenant context errors
    // Also don't logout if using impersonation - let the provider handle it
    let isUsingImpersonation = false
    try {
      isUsingImpersonation = !!sessionStorage.getItem('impersonation_token')
    } catch (e) {
      // sessionStorage might not be available
    }

    debugLog('[API] Error handler - status:', error.response?.status, 'isTenantContextError:', isTenantContextError, 'isUsingImpersonation:', isUsingImpersonation)

    if (error.response?.status === 401 && !isTenantContextError && !isUsingImpersonation) {
      debugWarn('[API] 401 auth error - logging out')
      debugLog('[API] Logging out due to 401')
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      localStorage.removeItem('auth-storage')
      authContextCache = null
      // Redirect to landing app login
      import('./config')
        .then(({ getConfig }) => {
          const { landingUrl } = getConfig()
          window.location.href = `${landingUrl}/auth/login`
        })
        .catch(() => {
          window.location.href = '/auth/login'
        })
    } else if (error.response?.status === 401 && isUsingImpersonation) {
      debugLog('[API] Got 401 while impersonating - NOT logging out, backend should handle impersonation token')
    }
    return Promise.reject(error)
  }
)

// API function to check tenant status
export async function getTenantStatus() {
  try {
    const response = await api.get('/tenant/status')
    return response.data
  } catch (error) {
    debugError('[API] Failed to get tenant status:', error)
    throw error
  }
}

export interface SubscriptionPlan {
  id: string
  name: string
  name_ar?: string
  price: number
  price_cents: number
  currency: string
  period: string
  features?: string[]
  features_ar?: string[]
}

export interface CurrentSubscriptionResponse {
  subscription_status: string
  trial_end_date?: string | null
  subscription_expires_at?: string | null
  subscription_cancelled_at?: string | null
  cancel_at_period_end?: boolean
  last_payment_at?: string | null
  days_remaining: number
  can_access: boolean
  block_reason?: string | null
  plan?: {
    id: string
    name: string
    price: number
    price_cents?: number
    currency: string
    period?: string
  }
  current_plan?: {
    id: string
    name: string
    price: number
    price_cents?: number
    currency: string
    period?: string
  }
  available_plans?: Array<{
    id: string
    name: string
    price: number
    price_cents?: number
    currency: string
    period?: string
  }>
}

export interface PaymentHistoryItem {
  id: number
  amount_cents: number
  currency: string
  status: 'pending' | 'success' | 'failed' | 'refunded'
  payment_method?: string | null
  card_last_four?: string | null
  card_brand?: string | null
  merchant_order_id: string
  created_at: string
  paid_at?: string | null
}

export interface InitiatePaymentResponse {
  payment_id: number
  iframe_url: string
  paymob_order_id: number
  merchant_order_id: string
}

export interface PaymentHistoryResponse {
  payments: PaymentHistoryItem[]
  total: number
  page: number
  page_size: number
}

export interface CancelSubscriptionResponse {
  subscription_expires_at?: string | null
  cancelled_at: string
  message: string
}

export async function getSubscriptionPlans(): Promise<SubscriptionPlan[]> {
  try {
    const response = await api.get('/subscription/plans')
    const payload = response.data?.data ?? response.data
    return Array.isArray(payload?.plans) ? payload.plans : []
  } catch (error) {
    debugError('[API] Failed to get subscription plans:', error)
    throw error
  }
}

export async function getCurrentSubscription(): Promise<CurrentSubscriptionResponse> {
  try {
    const response = await api.get('/subscription/current')
    const payload = response.data?.data ?? response.data
    return payload as CurrentSubscriptionResponse
  } catch (error) {
    debugError('[API] Failed to get current subscription:', error)
    throw error
  }
}

export async function getPaymentHistory(page = 1, pageSize = 20): Promise<PaymentHistoryResponse> {
  try {
    const response = await api.get('/subscription/payments', {
      params: {
        page,
        page_size: pageSize,
      },
    })
    const payload = response.data?.data ?? response.data
    return {
      payments: Array.isArray(payload?.payments) ? payload.payments : [],
      total: typeof payload?.total === 'number' ? payload.total : 0,
      page: typeof payload?.page === 'number' ? payload.page : page,
      page_size: typeof payload?.page_size === 'number' ? payload.page_size : pageSize,
    }
  } catch (error) {
    debugError('[API] Failed to get payment history:', error)
    throw error
  }
}

export async function initiateSubscriptionPayment(planID: string): Promise<InitiatePaymentResponse> {
  try {
    const response = await api.post('/subscription/initiate', { plan_id: planID })
    const payload = response.data?.data ?? response.data
    if (!payload?.iframe_url) {
      debugWarn('[API] Initiate payment response missing iframe_url', payload)
    }
    return payload as InitiatePaymentResponse
  } catch (error) {
    debugError('[API] Failed to initiate subscription payment:', error)
    throw error
  }
}

export async function cancelSubscription(reason?: string): Promise<CancelSubscriptionResponse> {
  try {
    const payload = reason ? { reason } : {}
    const response = await api.post('/subscription/cancel', payload)
    const data = response.data?.data ?? response.data
    return data as CancelSubscriptionResponse
  } catch (error) {
    debugError('[API] Failed to cancel subscription:', error)
    throw error
  }
}

export default api
