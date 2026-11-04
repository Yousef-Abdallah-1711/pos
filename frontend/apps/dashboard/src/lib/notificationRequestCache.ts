export interface NotificationRequestScope {
  tenantId: number | null | undefined
  restaurantId: number | null | undefined
  userId: number | null | undefined
}

const inFlightRequests = new Map<string, Promise<unknown>>()

export function createAuthenticatedScopeKey(scope: NotificationRequestScope): string {
  return `tenant:${scope.tenantId ?? 'none'}:restaurant:${scope.restaurantId ?? 'none'}:user:${scope.userId ?? 'none'}`
}

export function createNotificationCacheKey(
  scope: NotificationRequestScope,
  filters: unknown
): string {
  return `${createAuthenticatedScopeKey(scope)}:notifications:${JSON.stringify(filters)}`
}

export function dedupeInFlightRequest<T>(
  key: string,
  request: () => Promise<T>
): Promise<T> {
  const existing = inFlightRequests.get(key)
  if (existing) return existing as Promise<T>

  const pending = Promise.resolve().then(request)
  inFlightRequests.set(key, pending)
  void pending.then(
    () => clearCompletedRequest(key, pending),
    () => clearCompletedRequest(key, pending)
  )
  return pending
}

function clearCompletedRequest(key: string, pending: Promise<unknown>): void {
  if (inFlightRequests.get(key) === pending) inFlightRequests.delete(key)
}
