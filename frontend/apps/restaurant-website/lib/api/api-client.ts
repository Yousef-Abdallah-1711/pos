import axios, { AxiosInstance, AxiosRequestConfig, AxiosError } from 'axios'
import { debugError } from '../debug'

/**
 * API Client for Restaurant Website
 * Used to fetch themes and other data from the backend
 * No authentication required for public endpoints
 */

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

/**
 * Create axios instance for API calls
 * This is used by the restaurant website to fetch public data
 */
export const apiClient = {
  /**
   * GET request
   * @param url - The endpoint URL
   * @param config - Optional axios config
   * @returns Promise with response data
   */
  async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await axios.get<T>(`${API_URL}${url}`, {
        ...config,
        timeout: config?.timeout || 10000,
        headers: {
          'Content-Type': 'application/json',
          ...config?.headers,
        },
      })
      return response.data
    } catch (error) {
      handleError(error, url, 'GET')
      throw error
    }
  },

  /**
   * POST request
   * @param url - The endpoint URL
   * @param data - Request body data
   * @param config - Optional axios config
   * @returns Promise with response data
   */
  async post<T = any>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    try {
      const response = await axios.post<T>(`${API_URL}${url}`, data, {
        ...config,
        timeout: config?.timeout || 10000,
        headers: {
          'Content-Type': 'application/json',
          ...config?.headers,
        },
      })
      return response.data
    } catch (error) {
      handleError(error, url, 'POST')
      throw error
    }
  },

  /**
   * PUT request
   * @param url - The endpoint URL
   * @param data - Request body data
   * @param config - Optional axios config
   * @returns Promise with response data
   */
  async put<T = any>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    try {
      const response = await axios.put<T>(`${API_URL}${url}`, data, {
        ...config,
        timeout: config?.timeout || 10000,
        headers: {
          'Content-Type': 'application/json',
          ...config?.headers,
        },
      })
      return response.data
    } catch (error) {
      handleError(error, url, 'PUT')
      throw error
    }
  },

  /**
   * DELETE request
   * @param url - The endpoint URL
   * @param config - Optional axios config
   * @returns Promise with response data
   */
  async delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await axios.delete<T>(`${API_URL}${url}`, {
        ...config,
        timeout: config?.timeout || 10000,
        headers: {
          'Content-Type': 'application/json',
          ...config?.headers,
        },
      })
      return response.data
    } catch (error) {
      handleError(error, url, 'DELETE')
      throw error
    }
  },
}

/**
 * Handle API errors with logging
 * @param error - The error that occurred
 * @param url - The endpoint URL
 * @param method - The HTTP method
 */
function handleError(error: unknown, url: string, method: string): void {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError
    debugError(`[API ${method}] ${url}`)
    debugError('Status:', axiosError.response?.status)
    debugError('Data:', axiosError.response?.data)
    debugError('Message:', axiosError.message)
  } else if (error instanceof Error) {
    debugError(`[API ${method}] ${url}`)
    debugError('Error:', error.message)
  } else {
    debugError(`[API ${method}] ${url}`)
    debugError('Unknown error:', error)
  }
}

/**
 * Utility type for API responses
 */
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
  message?: string
}
