type RequestConfigForDedup = {
  method?: string
  url?: string
  params?: unknown
  data?: unknown
  headers?: Record<string, unknown> & { get?: (name: string) => unknown }
}

function identityValue(value: unknown): string {
  return typeof value === 'string' || typeof value === 'number' ? String(value) : ''
}

export function buildDedupIdentityScope(user: unknown): string | null {
  if (!user || typeof user !== 'object' || Array.isArray(user)) return null

  const candidate = user as Record<string, unknown>
  const userId = identityValue(candidate.id ?? candidate.user_id)
  if (!userId) return null

  return JSON.stringify([
    userId,
    identityValue(candidate.tenant_id ?? candidate.tenantId),
    identityValue(candidate.restaurant_id ?? candidate.restaurantId),
    identityValue(candidate.role),
    candidate.is_super_admin === true,
  ])
}

function readHeader(headers: RequestConfigForDedup['headers'], name: string): string {
  if (!headers) return ''
  const fromGetter = headers.get?.(name)
  if (fromGetter !== undefined && fromGetter !== null) return String(fromGetter)
  const key = Object.keys(headers).find((candidate) => candidate.toLowerCase() === name.toLowerCase())
  return key ? String(headers[key] ?? '') : ''
}

export function buildScopedGetRequestDedupKey(
  config: RequestConfigForDedup,
  identityScope: string | null,
): string | null {
  const method = String(config.method || 'get').toLowerCase()
  if (method !== 'get' || identityScope === null) return null

  // A caller-supplied bearer credential may represent an identity that is not
  // reflected in localStorage, so avoid sharing its response across requests.
  if (readHeader(config.headers, 'Authorization')) return null

  const tenantId = readHeader(config.headers, 'X-Tenant-ID')
  const restaurantId = readHeader(config.headers, 'X-Restaurant-ID')
  const params = config.params ? JSON.stringify(config.params) : ''
  const data = config.data
    ? typeof config.data === 'string' ? config.data : JSON.stringify(config.data)
    : ''
  return JSON.stringify([method, config.url ?? '', params, data, identityScope, tenantId, restaurantId])
}
