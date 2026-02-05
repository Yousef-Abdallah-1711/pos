'use client'
import axios from 'axios'
import { debugError, debugLog, debugWarn } from '@/lib/debug'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add token and tenant context to requests
api.interceptors.request.use((config) => {
  // Get token from Zustand store (which persists to localStorage under 'auth-storage')
  let token = null
  let user = null
  let tenantId = null
  let restaurantId = null

  try {
    const authStorage = localStorage.getItem('auth-storage')
    if (authStorage) {
      const parsed = JSON.parse(authStorage)
      token = parsed.state?.token || parsed.token
      user = parsed.state?.user || parsed.user

      // Extract tenant context with multiple fallback options
      tenantId = user?.tenant_id || user?.tenantId || null
      restaurantId = user?.restaurant_id || user?.restaurantId || null
    }
  } catch (e) {
    debugError('[API] Failed to parse auth-storage:', e)
  }

  // Fallback to direct localStorage keys
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

  // Don't add token to auth endpoints (like login, register)
  const isAuthEndpoint = config.url?.includes('/auth/')

  // Check for impersonation token
  let impersonationToken = null
  try {
    impersonationToken = sessionStorage.getItem('impersonation_token')
  } catch (e) {
    // sessionStorage might not be available in some contexts
  }

  if (impersonationToken && !isAuthEndpoint) {
    // If using impersonation token, add it as a custom header
    config.headers['X-Impersonation-Token'] = impersonationToken
  } else if (token && !isAuthEndpoint) {
    config.headers.Authorization = `Bearer ${token}`

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
    return response
  },
  (error) => {
    debugError('[API] Request failed:', error)

    // Check if it's a "Missing tenant information" error
    const errorMessage = error.response?.data?.error || ''
    const isTenantContextError = errorMessage.includes('Missing tenant') || errorMessage.includes('tenant')

    if (isTenantContextError) {
      debugWarn('[API] Tenant context error:', errorMessage)
    }

    // Only logout on actual authentication failures, not tenant context errors
    if (error.response?.status === 401 && !isTenantContextError) {
      debugWarn('[API] 401 auth error - logging out')
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      localStorage.removeItem('auth-storage')
      // Redirect to landing app login
      import('./config')
        .then(({ getConfig }) => {
          const { landingUrl } = getConfig()
          window.location.href = `${landingUrl}/auth/login`
        })
        .catch(() => {
          // Fallback to default URL if config import fails
          window.location.href = 'http://localhost:3001/auth/login'
        })
    }
    return Promise.reject(error)
  }
)

export default api
