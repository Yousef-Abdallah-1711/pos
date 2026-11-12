import axios from 'axios'
import { getConfig } from './config'

const API_URL = getConfig().apiUrl

export const api = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add token to requests (but not for auth endpoints)
api.interceptors.request.use((config) => {
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
