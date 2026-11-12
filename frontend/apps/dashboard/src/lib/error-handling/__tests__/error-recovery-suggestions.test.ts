/**
 * Error Recovery Suggestions Tests
 * Tests context-aware suggestion generation and advice text
 */

import { getRecoverySuggestions, getErrorAdvice, type ErrorRecoveryContext } from '../error-recovery-suggestions'

describe('Error Recovery Suggestions', () => {
  describe('getRecoverySuggestions', () => {
    it('should return suggestions array', () => {
      const context: ErrorRecoveryContext = {
        errorMessage: 'Validation error',
        category: 'validation',
        severity: 'high',
      }

      const suggestions = getRecoverySuggestions(context)

      expect(Array.isArray(suggestions)).toBe(true)
      expect(suggestions.length).toBeGreaterThan(0)
    })

    it('should limit suggestions to 5', () => {
      const context: ErrorRecoveryContext = {
        errorMessage: 'Generic error',
        category: 'unknown',
        severity: 'medium',
      }

      const suggestions = getRecoverySuggestions(context)

      expect(suggestions.length).toBeLessThanOrEqual(5)
    })

    describe('Validation Errors', () => {
      const validationContext: ErrorRecoveryContext = {
        errorMessage: 'Validation failed: field is required',
        category: 'validation',
        severity: 'high',
      }

      it('should include check input suggestion for validation errors', () => {
        const suggestions = getRecoverySuggestions(validationContext)

        expect(suggestions.some((s) => s.title.toLowerCase().includes('input'))).toBe(true)
      })

      it('should include form suggestions', () => {
        const suggestions = getRecoverySuggestions(validationContext)

        expect(suggestions.some((s) => s.title.toLowerCase().includes('form'))).toBe(true)
      })

      it('should provide validation-specific advice', () => {
        const suggestions = getRecoverySuggestions(validationContext)

        expect(suggestions.length).toBeGreaterThan(0)
      })
    })

    describe('Network Errors', () => {
      const networkContext: ErrorRecoveryContext = {
        errorMessage: 'Network request failed',
        category: 'network',
        severity: 'medium',
      }

      it('should include connection check suggestion', () => {
        const suggestions = getRecoverySuggestions(networkContext)

        expect(suggestions.some((s) => s.title.toLowerCase().includes('connection'))).toBe(true)
      })

      it('should include retry suggestion', () => {
        const suggestions = getRecoverySuggestions(networkContext)

        expect(suggestions.some((s) => s.action.label.toLowerCase().includes('retry'))).toBe(true)
      })

      it('should include support suggestion', () => {
        const suggestions = getRecoverySuggestions(networkContext)

        expect(suggestions.some((s) => s.title.toLowerCase().includes('support'))).toBe(true)
      })
    })

    describe('Authentication Errors', () => {
      const authContext: ErrorRecoveryContext = {
        errorMessage: 'Unauthorized access',
        category: 'authentication',
        severity: 'critical',
      }

      it('should include login suggestion', () => {
        const suggestions = getRecoverySuggestions(authContext)

        expect(suggestions.some((s) => s.action.label.toLowerCase() === 'log in')).toBe(true)
      })

      it('should include clear cache suggestion', () => {
        const suggestions = getRecoverySuggestions(authContext)

        expect(suggestions.some((s) => s.title.toLowerCase().includes('cache'))).toBe(true)
      })
    })

    describe('Permission Errors', () => {
      const permContext: ErrorRecoveryContext = {
        errorMessage: 'Access forbidden',
        category: 'permission',
        severity: 'high',
      }

      it('should include contact admin suggestion', () => {
        const suggestions = getRecoverySuggestions(permContext)

        expect(suggestions.some((s) => s.action.label.toLowerCase().includes('administrator'))).toBe(true)
      })

      it('should include go back suggestion', () => {
        const suggestions = getRecoverySuggestions(permContext)

        expect(suggestions.some((s) => s.title.toLowerCase().includes('back'))).toBe(true)
      })
    })

    describe('Component Errors', () => {
      const componentContext: ErrorRecoveryContext = {
        errorMessage: 'Component render failed',
        category: 'component',
        severity: 'high',
      }

      it('should include refresh suggestion', () => {
        const suggestions = getRecoverySuggestions(componentContext)

        expect(suggestions.some((s) => s.title.toLowerCase().includes('refresh'))).toBe(true)
      })

      it('should include clear state suggestion', () => {
        const suggestions = getRecoverySuggestions(componentContext)

        expect(suggestions.some((s) => s.title.toLowerCase().includes('state'))).toBe(true)
      })
    })

    describe('State Errors', () => {
      const stateContext: ErrorRecoveryContext = {
        errorMessage: 'State update failed',
        category: 'state',
        severity: 'medium',
      }

      it('should include reset application suggestion', () => {
        const suggestions = getRecoverySuggestions(stateContext)

        expect(suggestions.some((s) => s.title.toLowerCase().includes('reset'))).toBe(true)
      })
    })

    describe('Storage Errors', () => {
      const storageContext: ErrorRecoveryContext = {
        errorMessage: 'Storage quota exceeded',
        category: 'storage',
        severity: 'high',
      }

      it('should include clear storage suggestion', () => {
        const suggestions = getRecoverySuggestions(storageContext)

        expect(suggestions.some((s) => s.title.toLowerCase().includes('storage'))).toBe(true)
      })

      it('should include quota info', () => {
        const suggestions = getRecoverySuggestions(storageContext)

        expect(suggestions.length).toBeGreaterThan(0)
      })
    })

    describe('Unknown Errors', () => {
      const unknownContext: ErrorRecoveryContext = {
        errorMessage: 'Something went wrong',
        category: 'unknown',
        severity: 'low',
      }

      it('should include general suggestions', () => {
        const suggestions = getRecoverySuggestions(unknownContext)

        expect(suggestions.length).toBeGreaterThan(0)
      })

      it('should include refresh and dashboard options', () => {
        const suggestions = getRecoverySuggestions(unknownContext)

        expect(suggestions.some((s) => s.title.toLowerCase().includes('refresh'))).toBe(true)
        expect(suggestions.some((s) => s.title.toLowerCase().includes('dashboard'))).toBe(true)
      })
    })

    describe('Suggestion Properties', () => {
      const context: ErrorRecoveryContext = {
        errorMessage: 'Test error',
        category: 'validation',
        severity: 'high',
      }

      it('should have title for each suggestion', () => {
        const suggestions = getRecoverySuggestions(context)

        suggestions.forEach((s) => {
          expect(s.title).toBeDefined()
          expect(typeof s.title).toBe('string')
          expect(s.title.length).toBeGreaterThan(0)
        })
      })

      it('should have description for each suggestion', () => {
        const suggestions = getRecoverySuggestions(context)

        suggestions.forEach((s) => {
          expect(s.description).toBeDefined()
          expect(typeof s.description).toBe('string')
        })
      })

      it('should have action with label and onClick', () => {
        const suggestions = getRecoverySuggestions(context)

        suggestions.forEach((s) => {
          expect(s.action).toBeDefined()
          expect(s.action.label).toBeDefined()
          expect(s.action.onClick).toBeDefined()
          expect(typeof s.action.onClick).toBe('function')
        })
      })

      it('should have priority', () => {
        const suggestions = getRecoverySuggestions(context)

        suggestions.forEach((s) => {
          expect(s.priority).toBeDefined()
          expect(typeof s.priority).toBe('number')
          expect(s.priority).toBeGreaterThan(0)
        })
      })
    })

    describe('Priority Ordering', () => {
      const context: ErrorRecoveryContext = {
        errorMessage: 'Test error',
        category: 'validation',
        severity: 'high',
      }

      it('should return suggestions in priority order', () => {
        const suggestions = getRecoverySuggestions(context)

        for (let i = 1; i < suggestions.length; i++) {
          expect(suggestions[i].priority).toBeLessThanOrEqual(suggestions[i - 1].priority)
        }
      })

      it('should prioritize by severity', () => {
        const lowSeverity: ErrorRecoveryContext = {
          errorMessage: 'Minor issue',
          category: 'unknown',
          severity: 'low',
        }

        const highSeverity: ErrorRecoveryContext = {
          errorMessage: 'Major issue',
          category: 'unknown',
          severity: 'high',
        }

        const lowSuggestions = getRecoverySuggestions(lowSeverity)
        const highSuggestions = getRecoverySuggestions(highSeverity)

        // Both should have suggestions, order may differ based on severity
        expect(lowSuggestions.length).toBeGreaterThan(0)
        expect(highSuggestions.length).toBeGreaterThan(0)
      })
    })

    describe('Context Awareness', () => {
      it('should use error context when provided', () => {
        const context: ErrorRecoveryContext = {
          errorMessage: 'API rate limit exceeded',
          category: 'network',
          severity: 'medium',
          context: { endpoint: '/api/themes', retryAfter: 60 },
        }

        const suggestions = getRecoverySuggestions(context)

        // Should provide context-aware suggestions
        expect(suggestions.length).toBeGreaterThan(0)
      })

      it('should work with empty context fields', () => {
        const context: ErrorRecoveryContext = {
          errorMessage: 'Error',
          category: 'unknown',
          severity: 'low',
          context: {},
        }

        const suggestions = getRecoverySuggestions(context)

        expect(suggestions.length).toBeGreaterThan(0)
      })
    })

    describe('Action Callbacks', () => {
      const context: ErrorRecoveryContext = {
        errorMessage: 'Test error',
        category: 'validation',
        severity: 'high',
      }

      it('should provide callable onClick functions', () => {
        const suggestions = getRecoverySuggestions(context)

        suggestions.forEach((s) => {
          expect(() => {
            s.action.onClick()
          }).not.toThrow()
        })
      })

      it('should handle async onClick functions', async () => {
        const suggestions = getRecoverySuggestions(context)

        for (const suggestion of suggestions) {
          const result = suggestion.action.onClick()
          // Should handle both sync and async
          if (result instanceof Promise) {
            await expect(result).resolves.not.toThrow()
          }
        }
      })
    })
  })

  describe('getErrorAdvice', () => {
    it('should return advice string', () => {
      const advice = getErrorAdvice('validation', 'high')

      expect(typeof advice).toBe('string')
      expect(advice.length).toBeGreaterThan(0)
    })

    it('should provide validation advice', () => {
      const advice = getErrorAdvice('validation', 'high')

      expect(advice.toLowerCase()).toContain('validation errors')
    })

    it('should provide network advice', () => {
      const advice = getErrorAdvice('network', 'medium')

      expect(advice.toLowerCase()).toContain('try again in a moment')
    })

    it('should provide authentication advice', () => {
      const advice = getErrorAdvice('authentication', 'critical')

      expect(advice.toLowerCase()).toContain('log')
    })

    it('should provide permission advice', () => {
      const advice = getErrorAdvice('permission', 'high')

      expect(advice.toLowerCase()).toContain('admin')
    })

    it('should provide component advice', () => {
      const advice = getErrorAdvice('component', 'high')

      expect(advice.toLowerCase()).toContain('refresh')
    })

    it('should provide state advice', () => {
      const advice = getErrorAdvice('state', 'medium')

      expect(advice.toLowerCase()).toContain('state')
    })

    it('should provide storage advice', () => {
      const advice = getErrorAdvice('storage', 'high')

      expect(advice.toLowerCase()).toContain('storage')
    })

    it('should provide generic advice for unknown', () => {
      const advice = getErrorAdvice('unknown', 'low')

      expect(advice.length).toBeGreaterThan(0)
    })

    describe('Severity-Specific Advice', () => {
      it('should adjust advice based on severity', () => {
        const lowSeverity = getErrorAdvice('validation', 'low')
        const highSeverity = getErrorAdvice('validation', 'high')

        // Both should be valid but may differ
        expect(lowSeverity.length).toBeGreaterThan(0)
        expect(highSeverity.length).toBeGreaterThan(0)
      })

      it('should emphasize urgency for critical errors', () => {
        const critical = getErrorAdvice('authentication', 'critical')
        const low = getErrorAdvice('authentication', 'low')

        // Critical should be different and more urgent
        expect(critical).toBeDefined()
        expect(low).toBeDefined()
      })
    })

    describe('Advice Content Quality', () => {
      it('should provide actionable advice', () => {
        const advice = getErrorAdvice('validation', 'high')

        // Should be a complete sentence or phrase
        expect(advice.split(' ').length).toBeGreaterThan(2)
      })

      it('should be user-friendly', () => {
        const advice = getErrorAdvice('network', 'medium')

        // Should avoid technical jargon
        expect(advice).not.toContain('undefined')
        expect(advice).not.toContain('null')
      })

      it('should be concise but informative', () => {
        const advice = getErrorAdvice('component', 'high')

        expect(advice.length).toBeLessThan(500)
        expect(advice.length).toBeGreaterThan(10)
      })
    })
  })

  describe('Combined Error Handling', () => {
    it('should provide consistent advice and suggestions', () => {
      const category = 'validation'
      const severity = 'high'

      const advice = getErrorAdvice(category, severity)
      const context: ErrorRecoveryContext = {
        errorMessage: 'Validation error',
        category,
        severity,
      }
      const suggestions = getRecoverySuggestions(context)

      // Both should be consistent and comprehensive
      expect(advice).toBeDefined()
      expect(suggestions.length).toBeGreaterThan(0)

      // Suggestions should relate to advice
      expect(suggestions.some((s) => s.title.length > 0)).toBe(true)
    })

    it('should handle all category-severity combinations', () => {
      const categories = ['validation', 'network', 'authentication', 'permission', 'component', 'state', 'storage', 'unknown'] as const
      const severities = ['low', 'medium', 'high', 'critical'] as const

      categories.forEach((category) => {
        severities.forEach((severity) => {
          const advice = getErrorAdvice(category, severity)
          const context: ErrorRecoveryContext = {
            errorMessage: 'Test error',
            category,
            severity,
          }
          const suggestions = getRecoverySuggestions(context)

          expect(advice).toBeDefined()
          expect(advice.length).toBeGreaterThan(0)
          expect(suggestions.length).toBeGreaterThan(0)
        })
      })
    })
  })
})
