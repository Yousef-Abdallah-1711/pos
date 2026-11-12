import { getDashboardLoginUrl, hasUsableBearerToken } from '@/lib/authSession'
import { useAuthStore } from '@/stores/authStore'

const SAFE_METHODS = new Set(['GET', 'HEAD', 'OPTIONS'])

const readCookie = (name: string): string | null => {
  if (typeof document === 'undefined') {
    return null
  }

  const cookies = document.cookie.split(';').map((part) => part.trim())
  for (const cookie of cookies) {
    if (cookie.startsWith(`${name}=`)) {
      return decodeURIComponent(cookie.slice(name.length + 1))
    }
  }

  return null
}

export const buildAuthHeaders = (headers?: HeadersInit, method?: string) => {
  const mergedHeaders = new Headers(headers)
  const { token, user } = useAuthStore.getState()
  const impersonationToken =
    typeof window !== 'undefined' ? window.sessionStorage.getItem('impersonation_token') : null

  if (impersonationToken && !mergedHeaders.has('X-Impersonation-Token')) {
    mergedHeaders.set('X-Impersonation-Token', impersonationToken)
  }

  if (hasUsableBearerToken(token) && !mergedHeaders.has('Authorization')) {
    mergedHeaders.set('Authorization', `Bearer ${token}`)
  }

  const tenantId = user?.tenant_id || user?.tenantId
  const restaurantId = user?.restaurant_id || user?.restaurantId

  if (tenantId && !mergedHeaders.has('X-Tenant-ID')) {
    mergedHeaders.set('X-Tenant-ID', String(tenantId))
  }

  if (restaurantId && !mergedHeaders.has('X-Restaurant-ID')) {
    mergedHeaders.set('X-Restaurant-ID', String(restaurantId))
  }

  const normalizedMethod = String(method || 'GET').toUpperCase()
  if (!SAFE_METHODS.has(normalizedMethod)) {
    const csrfToken = readCookie('csrf_token')
    if (csrfToken && !mergedHeaders.has('X-CSRF-Token')) {
      mergedHeaders.set('X-CSRF-Token', csrfToken)
    }
  }

  return mergedHeaders
}

export const authFetch = async (input: RequestInfo | URL, init: RequestInit = {}) => {
  const headers = buildAuthHeaders(init.headers, init.method)

  if (init.body instanceof FormData) {
    headers.delete('Content-Type')
  }

  const response = await fetch(input, {
    ...init,
    credentials: init.credentials ?? 'include',
    headers,
  })

  const usingImpersonation =
    typeof window !== 'undefined' && !!window.sessionStorage.getItem('impersonation_token')

  if (response.status === 401 && !usingImpersonation) {
    useAuthStore.getState().logout()

    if (typeof window !== 'undefined') {
      window.location.href = getDashboardLoginUrl()
    }
  }

  return response
}
