import axios from 'axios'
import { getConfig } from './config'

const API_URL = getConfig().apiUrl

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add token to requests (but not for auth endpoints)
api.interceptors.request.use((config) => {
  // Don't add token to auth endpoints (like login, register)
  const isAuthEndpoint = config.url?.includes('/auth/')
  const token = localStorage.getItem('token')

  if (token && !isAuthEndpoint) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Handle 401 errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Don't redirect on auth endpoints to allow proper error handling
      const url = error.config?.url || '';
      if (!url.includes('/auth/')) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        // Redirect to the login page for non-auth requests
        import('./config').then(({ getConfig }) => {
          const { landingUrl } = getConfig()
          window.location.href = `${landingUrl}/auth/login`
        }).catch(() => {
          // Fallback to default URL if config import fails
          window.location.href = '/auth/login'
        })
      }
    }
    return Promise.reject(error)
  }
)

export default api
