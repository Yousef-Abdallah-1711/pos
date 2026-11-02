/**
 * Error logging service for tracking and reporting application errors
 * Provides centralized error logging, categorization, and reporting
 */

import { debugError, debugWarn } from '@/lib/debug'

export type ErrorSeverity = 'low' | 'medium' | 'high' | 'critical'

export type ErrorCategory =
  | 'validation'
  | 'network'
  | 'authentication'
  | 'permission'
  | 'component'
  | 'state'
  | 'storage'
  | 'unknown'

export interface ErrorLog {
  id: string
  message: string
  category: ErrorCategory
  severity: ErrorSeverity
  timestamp: Date
  userAgent: string
  url: string
  stack?: string
  context?: Record<string, any>
  userId?: string
  metadata?: Record<string, any>
}

interface ErrorLoggerConfig {
  maxLogs?: number
  enableLocalStorage?: boolean
  enableRemoteLogging?: boolean
  remoteEndpoint?: string
  environment?: 'development' | 'production' | 'staging'
}

export class ErrorLoggerService {
  private logs: ErrorLog[] = []
  private config: Required<ErrorLoggerConfig>
  private maxLogs: number = 100

  constructor(config: ErrorLoggerConfig = {}) {
    this.config = {
      maxLogs: config.maxLogs || 100,
      enableLocalStorage: config.enableLocalStorage !== false,
      enableRemoteLogging: config.enableRemoteLogging === true,
      remoteEndpoint: config.remoteEndpoint || '/api/v1/logs',
      environment: config.environment || 'development',
    }
    this.maxLogs = this.config.maxLogs

    this.loadFromStorage()
  }

  /**
   * Categorize error based on message and type
   */
  private categorizeError(error: Error | string): ErrorCategory {
    const message = typeof error === 'string' ? error : error.message
    const lowerMessage = message.toLowerCase()

    if (
      lowerMessage.includes('validation') ||
      lowerMessage.includes('invalid') ||
      lowerMessage.includes('required')
    ) {
      return 'validation'
    }

    if (
      lowerMessage.includes('network') ||
      lowerMessage.includes('fetch') ||
      lowerMessage.includes('timeout') ||
      lowerMessage.includes('connection')
    ) {
      return 'network'
    }

    if (
      lowerMessage.includes('auth') ||
      lowerMessage.includes('unauthorized') ||
      lowerMessage.includes('token')
    ) {
      return 'authentication'
    }

    if (
      lowerMessage.includes('permission') ||
      lowerMessage.includes('forbidden') ||
      lowerMessage.includes('access denied')
    ) {
      return 'permission'
    }

    if (lowerMessage.includes('component')) {
      return 'component'
    }

    if (
      lowerMessage.includes('state') ||
      lowerMessage.includes('cannot') ||
      lowerMessage.includes('undefined')
    ) {
      return 'state'
    }

    if (
      lowerMessage.includes('storage') ||
      lowerMessage.includes('localstorage') ||
      lowerMessage.includes('quota')
    ) {
      return 'storage'
    }

    return 'unknown'
  }

  /**
   * Determine error severity
   */
  private getSeverity(category: ErrorCategory, message: string): ErrorSeverity {
    const lowerMessage = message.toLowerCase()

    // Critical errors
    if (
      category === 'authentication' ||
      category === 'permission' ||
      lowerMessage.includes('critical') ||
      lowerMessage.includes('fatal')
    ) {
      return 'critical'
    }

    // High severity
    if (
      category === 'network' ||
      category === 'storage' ||
      lowerMessage.includes('error')
    ) {
      return 'high'
    }

    // Medium severity
    if (category === 'component' || category === 'state') {
      return 'medium'
    }

    // Low severity
    if (category === 'validation') {
      return 'low'
    }

    return 'medium'
  }

  /**
   * Log an error
   */
  log(
    error: Error | string,
    context?: Record<string, any>,
    metadata?: Record<string, any>
  ): ErrorLog {
    const message = typeof error === 'string' ? error : error.message
    const stack = typeof error === 'string' ? undefined : error.stack
    const category = this.categorizeError(error)
    const severity = this.getSeverity(category, message)

    const errorLog: ErrorLog = {
      id: this.generateId(),
      message,
      category,
      severity,
      timestamp: new Date(),
      userAgent: navigator.userAgent,
      url: window.location.href,
      stack,
      context,
      metadata,
    }

    this.logs.push(errorLog)

    // Keep only recent logs
    if (this.logs.length > this.maxLogs) {
      this.logs = this.logs.slice(-this.maxLogs)
    }

    // Save to local storage
    if (this.config.enableLocalStorage) {
      this.saveToStorage()
    }

    // Send to remote server if configured
    if (this.config.enableRemoteLogging && severity === 'high') {
      this.sendToRemote(errorLog)
    }

    // Log to console in development
    if (this.config.environment === 'development') {
      debugError(`[${category.toUpperCase()}] ${message}`, {
        severity,
        context,
        stack,
      })
    }

    return errorLog
  }

  /**
   * Log validation error
   */
  logValidationError(
    field: string,
    message: string,
    context?: Record<string, any>
  ): ErrorLog {
    return this.log(`Validation error: ${field} - ${message}`, {
      field,
      ...context,
    })
  }

  /**
   * Log network error
   */
  logNetworkError(
    endpoint: string,
    status: number,
    message: string,
    context?: Record<string, any>
  ): ErrorLog {
    return this.log(`Network error (${status}): ${endpoint} - ${message}`, {
      endpoint,
      status,
      ...context,
    })
  }

  /**
   * Log component error
   */
  logComponentError(
    componentName: string,
    error: Error,
    context?: Record<string, any>
  ): ErrorLog {
    return this.log(error, {
      component: componentName,
      ...context,
    })
  }

  /**
   * Get all logs
   */
  getAllLogs(): ErrorLog[] {
    return [...this.logs]
  }

  /**
   * Get logs by category
   */
  getLogsByCategory(category: ErrorCategory): ErrorLog[] {
    return this.logs.filter((log) => log.category === category)
  }

  /**
   * Get logs by severity
   */
  getLogsBySeverity(severity: ErrorSeverity): ErrorLog[] {
    return this.logs.filter((log) => log.severity === severity)
  }

  /**
   * Get recent logs
   */
  getRecentLogs(count: number = 10): ErrorLog[] {
    return this.logs.slice(-count)
  }

  /**
   * Clear all logs
   */
  clearLogs(): void {
    this.logs = []
    this.removeFromStorage()
  }

  /**
   * Get error summary
   */
  getSummary(): {
    total: number
    bySeverity: Record<ErrorSeverity, number>
    byCategory: Record<ErrorCategory, number>
    lastError?: ErrorLog
  } {
    const bySeverity: Record<ErrorSeverity, number> = {
      low: 0,
      medium: 0,
      high: 0,
      critical: 0,
    }

    const byCategory: Record<ErrorCategory, number> = {
      validation: 0,
      network: 0,
      authentication: 0,
      permission: 0,
      component: 0,
      state: 0,
      storage: 0,
      unknown: 0,
    }

    this.logs.forEach((log) => {
      bySeverity[log.severity]++
      byCategory[log.category]++
    })

    return {
      total: this.logs.length,
      bySeverity,
      byCategory,
      lastError: this.logs[this.logs.length - 1],
    }
  }

  /**
   * Export logs as JSON
   */
  exportLogs(): string {
    return JSON.stringify(this.logs, null, 2)
  }

  /**
   * Private methods
   */

  private generateId(): string {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
  }

  private saveToStorage(): void {
    try {
      localStorage.setItem('errorLogs', JSON.stringify(this.logs))
    } catch (e) {
      debugWarn('Failed to save error logs to localStorage', e)
    }
  }

  private loadFromStorage(): void {
    try {
      const stored = localStorage.getItem('errorLogs')
      if (stored) {
        this.logs = JSON.parse(stored)
      }
    } catch (e) {
      debugWarn('Failed to load error logs from localStorage', e)
    }
  }

  private removeFromStorage(): void {
    try {
      localStorage.removeItem('errorLogs')
    } catch (e) {
      debugWarn('Failed to remove error logs from localStorage', e)
    }
  }

  private sendToRemote(errorLog: ErrorLog): void {
    // Only send critical errors to remote
    if (errorLog.severity === 'critical' || errorLog.severity === 'high') {
      fetch(this.config.remoteEndpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(errorLog),
      }).catch((e) => {
        debugWarn('Failed to send error log to remote', e)
      })
    }
  }
}

// Create singleton instance
let loggerInstance: ErrorLoggerService | null = null

/**
 * Get or create error logger instance
 */
export function getErrorLogger(config?: ErrorLoggerConfig): ErrorLoggerService {
  if (!loggerInstance) {
    loggerInstance = new ErrorLoggerService(config)
  }
  return loggerInstance
}

/**
 * Initialize error logger with config
 */
export function initializeErrorLogger(config: ErrorLoggerConfig): ErrorLoggerService {
  loggerInstance = new ErrorLoggerService(config)
  return loggerInstance
}

export default ErrorLoggerService
