export type ApiFieldErrors = Record<string, string>

export interface ApiErrorFallbacks {
  validation: string
  invalidCredentials: string
  rateLimited: string
  unavailable: string
  network: string
  unexpected: string
  conflict: string
}

export interface ParsedApiError {
  code: string
  message: string
  fields: ApiFieldErrors
  requestId?: string
  status?: number
  retryable: boolean
}

const SAFE_MESSAGES: Record<string, keyof ApiErrorFallbacks> = {
  VALIDATION_ERROR: 'validation',
  INVALID_CREDENTIALS: 'invalidCredentials',
  RATE_LIMITED: 'rateLimited',
  SERVICE_UNAVAILABLE: 'unavailable',
  INTERNAL_ERROR: 'unexpected',
  INVALID_REQUEST: 'unexpected',
  REQUEST_TOO_LARGE: 'validation',
  EMAIL_ALREADY_EXISTS: 'conflict',
  SUBDOMAIN_TAKEN: 'conflict',
  RESTAURANT_NAME_CONFLICT: 'conflict',
}

function readResponse(error: unknown): { status?: number; data?: Record<string, unknown>; code?: string } {
  if (typeof error !== 'object' || error === null) return {}
  const candidate = error as { response?: { status?: unknown; data?: unknown }; code?: unknown }
  const status = typeof candidate.response?.status === 'number' ? candidate.response.status : undefined
  const data = typeof candidate.response?.data === 'object' && candidate.response.data !== null
    ? candidate.response.data as Record<string, unknown>
    : undefined
  const code = typeof candidate.code === 'string' ? candidate.code : undefined
  return { status, data, code }
}

export function parseApiError(error: unknown, fallback: ApiErrorFallbacks): ParsedApiError {
  const response = readResponse(error)
  const bodyCode = response.data?.code
  const code = typeof bodyCode === 'string' ? bodyCode : ''
  const messageKey = SAFE_MESSAGES[code]
  const fields: ApiFieldErrors = {}
  const rawFields = response.data?.fields
  if (typeof rawFields === 'object' && rawFields !== null) {
    for (const [name, value] of Object.entries(rawFields)) {
      if (typeof value === 'string' && value.length > 0) fields[name] = value
    }
  }

  const status = response.status
  const retryable = status === 429 || (status !== undefined && status >= 500) || !status
  const requestId = typeof response.data?.requestId === 'string' ? response.data.requestId : undefined
  const message = messageKey
    ? fallback[messageKey]
    : status === 429
      ? fallback.rateLimited
      : status === 401
        ? fallback.invalidCredentials
        : status && status >= 500
          ? fallback.unavailable
          : !status && (response.code === 'ERR_NETWORK' || response.code === 'ECONNABORTED' || response.code === 'ETIMEDOUT')
            ? fallback.network
            : fallback.unexpected

  return { code: code || 'UNKNOWN_ERROR', message, fields, requestId, status, retryable }
}

export function applyApiFieldErrors(
  fields: ApiFieldErrors,
  setError: (field: string, message: string) => void,
): void {
  const aliases: Record<string, string> = { restaurant_slug: 'restaurantSlug', restaurant_name: 'restaurantName' }
  for (const [field, message] of Object.entries(fields)) setError(aliases[field] ?? field, message)
}
