/**
 * Error recovery suggestions
 * Provides context-aware suggestions for different error types
 */

import { ErrorCategory, ErrorSeverity } from './error-logger'

export interface RecoverySuggestion {
  title: string
  description: string
  action: {
    label: string
    onClick: () => void | Promise<void>
  }
  priority: number // Higher = more important
}

export interface ErrorRecoveryContext {
  errorMessage: string
  category: ErrorCategory
  severity: ErrorSeverity
  context?: Record<string, any>
}

/**
 * Get recovery suggestions for an error
 */
export function getRecoverySuggestions(
  context: ErrorRecoveryContext
): RecoverySuggestion[] {
  const suggestions: RecoverySuggestion[] = []

  // Add category-specific suggestions
  switch (context.category) {
    case 'validation':
      suggestions.push(...getValidationSuggestions(context))
      break
    case 'network':
      suggestions.push(...getNetworkSuggestions(context))
      break
    case 'authentication':
      suggestions.push(...getAuthenticationSuggestions(context))
      break
    case 'permission':
      suggestions.push(...getPermissionSuggestions(context))
      break
    case 'component':
      suggestions.push(...getComponentSuggestions(context))
      break
    case 'state':
      suggestions.push(...getStateSuggestions(context))
      break
    case 'storage':
      suggestions.push(...getStorageSuggestions(context))
      break
    default:
      suggestions.push(...getGeneralSuggestions(context))
  }

  // Add severity-based suggestions
  if (context.severity === 'critical') {
    suggestions.unshift({
      title: 'Critical Error',
      description: 'This is a critical error. Please refresh the page to continue.',
      action: {
        label: 'Refresh Page',
        onClick: () => window.location.reload(),
      },
      priority: 100,
    })
  }

  // Sort by priority (highest first)
  suggestions.sort((a, b) => b.priority - a.priority)

  return suggestions.slice(0, 5) // Return top 5 suggestions
}

function getValidationSuggestions(context: ErrorRecoveryContext): RecoverySuggestion[] {
  return [
    {
      title: 'Check Your Input',
      description: 'The error may be caused by invalid input. Please review the form fields and try again.',
      action: {
        label: 'Review Form',
        onClick: () => {
          // Scroll to form
          const form = document.querySelector('form')
          form?.scrollIntoView({ behavior: 'smooth' })
        },
      },
      priority: 80,
    },
    {
      title: 'Clear Form and Retry',
      description: 'Try clearing the form and filling it out again.',
      action: {
        label: 'Clear Form',
        onClick: () => {
          const form = document.querySelector('form') as HTMLFormElement
          form?.reset()
        },
      },
      priority: 60,
    },
  ]
}

function getNetworkSuggestions(context: ErrorRecoveryContext): RecoverySuggestion[] {
  return [
    {
      title: 'Check Internet Connection',
      description: 'Network errors are often caused by connectivity issues. Please check your internet connection.',
      action: {
        label: 'Retry Request',
        onClick: () => {
          // User can retry through error boundary
          location.reload()
        },
      },
      priority: 90,
    },
    {
      title: 'Retry Later',
      description: 'The server may be temporarily unavailable. Please try again in a few moments.',
      action: {
        label: 'Wait and Retry',
        onClick: async () => {
          await new Promise((resolve) => setTimeout(resolve, 3000))
          location.reload()
        },
      },
      priority: 70,
    },
    {
      title: 'Contact Support',
      description: 'If the problem persists, please contact our support team.',
      action: {
        label: 'Contact Support',
        onClick: () => {
          window.open('mailto:support@example.com', '_blank')
        },
      },
      priority: 50,
    },
  ]
}

function getAuthenticationSuggestions(context: ErrorRecoveryContext): RecoverySuggestion[] {
  return [
    {
      title: 'Login Again',
      description: 'Your session may have expired. Please log in again.',
      action: {
        label: 'Log In',
        onClick: () => {
          window.location.href = '/login'
        },
      },
      priority: 95,
    },
    {
      title: 'Clear Browser Cache',
      description: 'Cached authentication data may be invalid. Try clearing your browser cache.',
      action: {
        label: 'Instructions',
        onClick: () => {
          alert(
            'How to clear cache:\n\n' +
              'Chrome/Edge: Ctrl+Shift+Del (or Cmd+Shift+Del on Mac)\n' +
              'Firefox: Ctrl+Shift+Del (or Cmd+Shift+Del on Mac)\n' +
              'Safari: Develop > Empty Caches'
          )
        },
      },
      priority: 60,
    },
  ]
}

function getPermissionSuggestions(context: ErrorRecoveryContext): RecoverySuggestion[] {
  return [
    {
      title: 'Contact Administrator',
      description: 'You do not have permission to perform this action. Please contact your administrator.',
      action: {
        label: 'Contact Administrator',
        onClick: () => {
          window.open('mailto:admin@example.com', '_blank')
        },
      },
      priority: 85,
    },
    {
      title: 'Go Back',
      description: 'You can go back to the previous page.',
      action: {
        label: 'Go Back',
        onClick: () => {
          window.history.back()
        },
      },
      priority: 50,
    },
  ]
}

function getComponentSuggestions(context: ErrorRecoveryContext): RecoverySuggestion[] {
  return [
    {
      title: 'Refresh Component',
      description: 'Try refreshing the page to reload the component.',
      action: {
        label: 'Refresh',
        onClick: () => {
          window.location.reload()
        },
      },
      priority: 75,
    },
    {
      title: 'Clear Application State',
      description: 'The component state may be corrupted. Try clearing the application data.',
      action: {
        label: 'Clear Data',
        onClick: () => {
          if (confirm('This will clear all unsaved changes. Continue?')) {
            localStorage.clear()
            sessionStorage.clear()
            window.location.reload()
          }
        },
      },
      priority: 65,
    },
  ]
}

function getStateSuggestions(context: ErrorRecoveryContext): RecoverySuggestion[] {
  return [
    {
      title: 'Reset Application State',
      description: 'The application state may be inconsistent. Try resetting it.',
      action: {
        label: 'Reset',
        onClick: () => {
          if (confirm('This will reset the application. Continue?')) {
            // Clear relevant state
            sessionStorage.clear()
            window.location.reload()
          }
        },
      },
      priority: 80,
    },
    {
      title: 'Save Work and Refresh',
      description: 'Save any important work before refreshing the page.',
      action: {
        label: 'Refresh',
        onClick: () => {
          window.location.reload()
        },
      },
      priority: 60,
    },
  ]
}

function getStorageSuggestions(context: ErrorRecoveryContext): RecoverySuggestion[] {
  return [
    {
      title: 'Clear Browser Storage',
      description: 'Your browser storage may be full or corrupted. Clear it to free up space.',
      action: {
        label: 'Clear Storage',
        onClick: () => {
          try {
            localStorage.clear()
            sessionStorage.clear()
            alert('Browser storage cleared successfully')
            window.location.reload()
          } catch (e) {
            alert('Failed to clear storage')
          }
        },
      },
      priority: 85,
    },
    {
      title: 'Check Storage Quota',
      description: 'Your browser storage quota may be exceeded.',
      action: {
        label: 'Learn More',
        onClick: () => {
          window.open('https://developer.mozilla.org/en-US/docs/Web/API/Storage', '_blank')
        },
      },
      priority: 50,
    },
  ]
}

function getGeneralSuggestions(context: ErrorRecoveryContext): RecoverySuggestion[] {
  return [
    {
      title: 'Refresh Page',
      description: 'Try refreshing the page to see if the issue resolves.',
      action: {
        label: 'Refresh',
        onClick: () => {
          window.location.reload()
        },
      },
      priority: 80,
    },
    {
      title: 'Go to Dashboard',
      description: 'Return to the main dashboard to continue working.',
      action: {
        label: 'Dashboard',
        onClick: () => {
          window.location.href = '/en/dashboard'
        },
      },
      priority: 60,
    },
    {
      title: 'Contact Support',
      description: 'If the problem persists, please contact our support team.',
      action: {
        label: 'Contact',
        onClick: () => {
          window.open('mailto:support@example.com', '_blank')
        },
      },
      priority: 40,
    },
  ]
}

/**
 * Get error-specific advice text
 */
export function getErrorAdvice(category: ErrorCategory, severity: ErrorSeverity): string {
  const adviceMap: Record<string, Record<string, string>> = {
    validation: {
      low: 'Please check the form fields and correct any errors.',
      medium: 'Validation failed. Please review the required fields.',
      high: 'Please check the input fields and correct all validation errors before continuing.',
      critical: 'Critical validation error. Please refresh and try again.',
    },
    network: {
      low: 'Please check your internet connection and try again.',
      medium: 'Network connection failed. Please try again in a moment.',
      high: 'Server is not responding. Please check your connection and try again.',
      critical: 'Network connection lost. Please restore your connection and refresh.',
    },
    authentication: {
      low: 'Please ensure you are properly logged in.',
      medium: 'Authentication failed. Please log in again.',
      high: 'Your session has expired. Please log in again.',
      critical: 'Critical authentication error. Please log in again.',
    },
    permission: {
      low: 'You may not have permission for this action.',
      medium: 'You do not have permission to perform this action.',
      high: 'Permission denied. Contact your administrator.',
      critical: 'Access denied. Contact your administrator.',
    },
    component: {
      low: 'A component encountered an issue.',
      medium: 'Component error detected.',
      high: 'Critical component error. Refresh the page.',
      critical: 'Fatal component error. System requires restart.',
    },
    state: {
      low: 'Application state is inconsistent.',
      medium: 'State management error detected.',
      high: 'Critical state error. Please refresh the page.',
      critical: 'Fatal state error. System corruption detected.',
    },
    storage: {
      low: 'Storage access issue.',
      medium: 'Storage error detected.',
      high: 'Storage full or corrupted. Clear cache.',
      critical: 'Storage failure. Clear all data.',
    },
    unknown: {
      low: 'An unexpected issue occurred.',
      medium: 'An unexpected error occurred.',
      high: 'A critical error occurred. Please refresh.',
      critical: 'A fatal error occurred. System restart required.',
    },
  }

  return adviceMap[category]?.[severity] || 'An error occurred. Please try again.'
}

export default {
  getRecoverySuggestions,
  getErrorAdvice,
}
