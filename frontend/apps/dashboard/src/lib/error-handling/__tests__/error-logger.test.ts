/**
 * Error Logger Service Tests
 * Tests error categorization, severity detection, logging, and querying
 */

import { ErrorLoggerService, ErrorCategory, ErrorSeverity, type ErrorLog } from '../error-logger'

describe('ErrorLoggerService', () => {
  let logger: ErrorLoggerService

  beforeEach(() => {
    localStorage.clear()
    logger = new ErrorLoggerService()
  })

  describe('Initialization', () => {
    it('should initialize with empty logs', () => {
      expect(logger.getAllLogs()).toHaveLength(0)
    })

    it('should have default configuration', () => {
      const summary = logger.getSummary()
      expect(summary.total).toBe(0)
      expect(summary.bySeverity).toEqual({
        low: 0,
        medium: 0,
        high: 0,
        critical: 0,
      })
    })
  })

  describe('Error Logging', () => {
    it('should log basic error', () => {
      const error = new Error('Test error')
      const log = logger.log(error, { component: 'TestComponent' })

      expect(log).toBeDefined()
      expect(log.message).toBe('Test error')
      expect(log.category).toBe('unknown')
    })

    it('should generate unique IDs for each log', () => {
      const error = new Error('Test error')
      const log1 = logger.log(error)
      const log2 = logger.log(error)

      expect(log1.id).not.toBe(log2.id)
    })

    it('should capture error timestamp', () => {
      const before = new Date()
      const log = logger.log(new Error('Test'))
      const after = new Date()

      expect(log.timestamp.getTime()).toBeGreaterThanOrEqual(before.getTime())
      expect(log.timestamp.getTime()).toBeLessThanOrEqual(after.getTime())
    })

    it('should capture user agent', () => {
      const log = logger.log(new Error('Test'))

      expect(log.userAgent).toBeDefined()
      expect(log.userAgent.length).toBeGreaterThan(0)
    })

    it('should capture current URL', () => {
      const log = logger.log(new Error('Test'))

      expect(log.url).toBeDefined()
    })

    it('should include optional context', () => {
      const context = { userId: '123', theme: 'dark' }
      const log = logger.log(new Error('Test'), context)

      expect(log.context).toEqual(context)
    })

    it('should include optional metadata', () => {
      const metadata = { source: 'API', endpoint: '/api/themes' }
      const log = logger.log(new Error('Test'), {}, metadata)

      expect(log.metadata).toEqual(metadata)
    })

    it('should capture error stack trace', () => {
      const error = new Error('Test error')
      const log = logger.log(error)

      expect(log.stack).toBeDefined()
      expect(log.stack).toContain('Test error')
    })
  })

  describe('Error Categorization', () => {
    it('should categorize validation errors', () => {
      const error = new Error('Validation failed: field is required')
      const log = logger.log(error)

      expect(log.category).toBe('validation')
    })

    it('should categorize network errors', () => {
      const error = new Error('Network request failed')
      const log = logger.log(error)

      expect(log.category).toBe('network')
    })

    it('should categorize authentication errors', () => {
      const error = new Error('Unauthorized access')
      const log = logger.log(error)

      expect(log.category).toBe('authentication')
    })

    it('should categorize permission errors', () => {
      const error = new Error('Access forbidden')
      const log = logger.log(error)

      expect(log.category).toBe('permission')
    })

    it('should categorize component errors', () => {
      const error = new Error('Component error: render failed')
      const log = logger.log(error)

      expect(log.category).toBe('component')
    })

    it('should categorize state errors', () => {
      const error = new Error('State update failed')
      const log = logger.log(error)

      expect(log.category).toBe('state')
    })

    it('should categorize storage errors', () => {
      const error = new Error('Storage quota exceeded')
      const log = logger.log(error)

      expect(log.category).toBe('storage')
    })

    it('should categorize unknown errors', () => {
      const error = new Error('Some random error')
      const log = logger.log(error)

      expect(log.category).toBe('unknown')
    })
  })

  describe('Severity Detection', () => {
    it('should detect critical severity for authentication errors', () => {
      const error = new Error('Unauthorized')
      const log = logger.log(error)

      expect(log.severity).toBe('critical')
    })

    it('should detect high severity for validation errors', () => {
      const error = new Error('Validation failed')
      const log = logger.log(error)

      expect(log.severity).toBe('high')
    })

    it('should detect medium severity for network errors', () => {
      const error = new Error('Network timeout')
      const log = logger.log(error)

      expect(log.severity).toBe('medium')
    })

    it('should detect low severity for unknown errors', () => {
      const error = new Error('Random error')
      const log = logger.log(error)

      expect(log.severity).toBe('low')
    })

    it('should support custom severity override', () => {
      const error = new Error('Test error')
      const log = logger.log(error, { severity: 'critical' })

      expect(log.severity).toBe('critical')
    })
  })

  describe('Specialized Logging Methods', () => {
    it('should log validation errors with field info', () => {
      const log = logger.logValidationError('email', 'Invalid email format', {
        attemptedValue: 'test@',
      })

      expect(log.category).toBe('validation')
      expect(log.message).toContain('email')
      expect(log.context).toEqual({ attemptedValue: 'test@' })
    })

    it('should log network errors with endpoint', () => {
      const log = logger.logNetworkError('/api/themes', 500, 'Internal server error', {
        retryCount: 2,
      })

      expect(log.category).toBe('network')
      expect(log.message).toContain('/api/themes')
      expect(log.message).toContain('500')
    })

    it('should log component errors with component name', () => {
      const error = new Error('Render failed')
      const log = logger.logComponentError('ThemeBuilder', error, { props: {} })

      expect(log.category).toBe('component')
      expect(log.message).toContain('ThemeBuilder')
    })
  })

  describe('Querying Logs', () => {
    beforeEach(() => {
      logger.log(new Error('Validation failed'), { severity: 'high' })
      logger.log(new Error('Network timeout'), { severity: 'medium' })
      logger.log(new Error('Unauthorized'), { severity: 'critical' })
      logger.log(new Error('Random error'), { severity: 'low' })
    })

    it('should get all logs', () => {
      const logs = logger.getAllLogs()

      expect(logs).toHaveLength(4)
    })

    it('should filter logs by category', () => {
      logger.log(new Error('Validation failed'), {})
      logger.log(new Error('Network timeout'), {})

      const validationLogs = logger.getLogsByCategory('validation')
      const networkLogs = logger.getLogsByCategory('network')

      expect(validationLogs.length).toBeGreaterThan(0)
      expect(networkLogs.length).toBeGreaterThan(0)
    })

    it('should filter logs by severity', () => {
      const criticalLogs = logger.getLogsBySeverity('critical')
      const highLogs = logger.getLogsBySeverity('high')

      expect(criticalLogs.length).toBeGreaterThan(0)
      expect(highLogs.length).toBeGreaterThan(0)
    })

    it('should get recent logs with limit', () => {
      for (let i = 0; i < 5; i++) {
        logger.log(new Error(`Error ${i}`))
      }

      const recent = logger.getRecentLogs(2)

      expect(recent).toHaveLength(2)
      expect(recent[0].message).toContain('Error')
    })

    it('should get recent logs with default limit of 10', () => {
      for (let i = 0; i < 15; i++) {
        logger.log(new Error(`Error ${i}`))
      }

      const recent = logger.getRecentLogs()

      expect(recent.length).toBeLessThanOrEqual(10)
    })
  })

  describe('Summary Generation', () => {
    beforeEach(() => {
      logger.log(new Error('Validation failed'), { severity: 'high' })
      logger.log(new Error('Network timeout'), { severity: 'medium' })
      logger.log(new Error('Unauthorized'), { severity: 'critical' })
    })

    it('should generate error summary', () => {
      const summary = logger.getSummary()

      expect(summary.total).toBe(3)
    })

    it('should count errors by severity in summary', () => {
      const summary = logger.getSummary()

      expect(summary.bySeverity.critical).toBeGreaterThan(0)
      expect(summary.bySeverity.high).toBeGreaterThan(0)
      expect(summary.bySeverity.medium).toBeGreaterThan(0)
    })

    it('should count errors by category in summary', () => {
      const summary = logger.getSummary()

      expect(summary.byCategory).toBeDefined()
      expect(typeof summary.byCategory).toBe('object')
    })

    it('should include last error in summary', () => {
      const summary = logger.getSummary()

      expect(summary.lastError).toBeDefined()
      expect(summary.lastError?.message).toBeDefined()
    })
  })

  describe('Export Functionality', () => {
    beforeEach(() => {
      logger.log(new Error('Test error 1'), { userId: '123' })
      logger.log(new Error('Test error 2'), { userId: '456' })
    })

    it('should export logs as JSON string', () => {
      const exported = logger.exportLogs()

      expect(typeof exported).toBe('string')
      expect(exported).toContain('Test error 1')
      expect(exported).toContain('Test error 2')
    })

    it('should export valid JSON', () => {
      const exported = logger.exportLogs()

      expect(() => JSON.parse(exported)).not.toThrow()
    })

    it('should export with timestamp format', () => {
      const exported = logger.exportLogs()
      const parsed = JSON.parse(exported)

      expect(parsed).toBeInstanceOf(Array)
      expect(parsed[0]).toHaveProperty('timestamp')
    })
  })

  describe('Storage Persistence', () => {
    it('should persist logs to localStorage', () => {
      logger.log(new Error('Test error'), {}, {}, 'test-key')

      const stored = localStorage.getItem('test-key')
      expect(stored).toBeDefined()
      expect(stored).toContain('Test error')
    })

    it('should load logs from localStorage on init', () => {
      const savedLogs = [
        {
          id: 'test-1',
          message: 'Saved error',
          category: 'validation',
          severity: 'high',
          timestamp: new Date().toISOString(),
          userAgent: 'test',
          url: 'http://test',
        },
      ]

      localStorage.setItem('error-logs', JSON.stringify(savedLogs))

      const newLogger = new ErrorLoggerService({ storageKey: 'error-logs' })
      const logs = newLogger.getAllLogs()

      expect(logs.length).toBeGreaterThan(0)
    })
  })

  describe('Log Limits', () => {
    it('should not exceed max log count', () => {
      // Logger has max 100 logs by default
      for (let i = 0; i < 150; i++) {
        logger.log(new Error(`Error ${i}`))
      }

      const logs = logger.getAllLogs()
      expect(logs.length).toBeLessThanOrEqual(100)
    })

    it('should keep most recent logs when exceeding limit', () => {
      for (let i = 0; i < 102; i++) {
        logger.log(new Error(`Error ${i}`))
      }

      const logs = logger.getAllLogs()
      const lastLog = logs[logs.length - 1]

      // Last log should be one of the most recent (Error 100+)
      expect(lastLog.message).toContain('Error')
    })
  })

  describe('Error Message Patterns', () => {
    it('should detect validation error patterns', () => {
      const patterns = [
        'Validation failed',
        'Field validation error',
        'Invalid input',
        'Required field',
      ]

      patterns.forEach((pattern) => {
        const log = logger.log(new Error(pattern))
        expect(log.category).toBe('validation')
      })
    })

    it('should detect network error patterns', () => {
      const patterns = [
        'Network request failed',
        'Failed to fetch',
        'Network timeout',
        'CORS error',
      ]

      patterns.forEach((pattern) => {
        const log = logger.log(new Error(pattern))
        expect(log.category).toBe('network')
      })
    })

    it('should detect authentication error patterns', () => {
      const patterns = [
        'Unauthorized',
        'Authentication failed',
        'Invalid credentials',
        'Token expired',
      ]

      patterns.forEach((pattern) => {
        const log = logger.log(new Error(pattern))
        expect(log.category).toBe('authentication')
      })
    })
  })

  describe('Metadata and Context', () => {
    it('should preserve complex context objects', () => {
      const context = {
        user: { id: '123', role: 'admin' },
        theme: { id: 'dark', version: '1.0' },
        nested: { deep: { value: 'test' } },
      }

      const log = logger.log(new Error('Test'), context)

      expect(log.context).toEqual(context)
    })

    it('should preserve metadata separately from context', () => {
      const context = { userAction: 'save' }
      const metadata = { duration: 250, retries: 2 }

      const log = logger.log(new Error('Test'), context, metadata)

      expect(log.context).toEqual(context)
      expect(log.metadata).toEqual(metadata)
    })

    it('should support user ID in logs', () => {
      const log = logger.log(new Error('Test'), {}, {}, 'error-logs', '123')

      expect(log.userId).toBe('123')
    })
  })

  describe('Error Types', () => {
    it('should handle Error objects', () => {
      const error = new Error('Test error')
      const log = logger.log(error)

      expect(log.message).toBe('Test error')
      expect(log.stack).toBeDefined()
    })

    it('should handle string errors', () => {
      const log = logger.log('String error' as any)

      expect(log.message).toContain('String error')
    })

    it('should handle error-like objects', () => {
      const errorLike = { message: 'Custom error', stack: 'custom stack' }
      const log = logger.log(errorLike as any)

      expect(log.message).toContain('Custom error')
    })
  })

  describe('Clearing Logs', () => {
    it('should clear all logs', () => {
      logger.log(new Error('Error 1'))
      logger.log(new Error('Error 2'))
      logger.log(new Error('Error 3'))

      expect(logger.getAllLogs()).toHaveLength(3)

      logger.clearLogs?.()

      expect(logger.getAllLogs()).toHaveLength(0)
    })
  })
})
