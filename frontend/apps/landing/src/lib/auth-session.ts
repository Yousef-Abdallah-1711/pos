import { getConfig } from '@/lib/config'

export interface LandingSessionProfile {
  id: number
  name: string
  email: string
  restaurantName?: string
}

export type LandingSessionState =
  | { status: 'checking' }
  | { status: 'guest' }
  | { status: 'authenticated'; profile: LandingSessionProfile }
  | { status: 'temporarily-unavailable'; reason?: 'session-check' | 'sign-out' }

type ResolvedSessionState = Exclude<LandingSessionState, { status: 'checking' }>

const SESSION_FRESHNESS_MS = 30_000
const SESSION_REQUEST_TIMEOUT_MS = 5_000
const AUTH_SESSION_CHANNEL = 'servio-auth-session'
let cachedState: ResolvedSessionState | null = null
let cachedAt = 0
let inFlightRequest: { promise: Promise<ResolvedSessionState>; controller: AbortController } | null = null

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null
}

function parseSessionProfile(value: unknown): LandingSessionProfile | null {
  if (!isRecord(value) || value.is_super_admin !== false) return null
  const userID = typeof value.id === 'number' ? value.id : value.user_id
  if (typeof userID !== 'number' || !Number.isSafeInteger(userID) || userID <= 0) return null
  if (typeof value.id === 'number' && typeof value.user_id === 'number' && value.id !== value.user_id) return null
  if (typeof value.email !== 'string' || value.email.length === 0) return null
  const name = typeof value.name === 'string' && value.name.trim() ? value.name.trim() : value.email
  const restaurantName = typeof value.restaurant_name === 'string' ? value.restaurant_name.trim() : ''

  return {
    id: userID,
    name,
    email: value.email,
    ...(restaurantName ? { restaurantName } : {}),
  }
}

async function requestSession(signal: AbortSignal): Promise<ResolvedSessionState> {
  try {
    const response = await fetchWithTimeout(`${getConfig().apiUrl}/auth/session`, {
      credentials: 'include',
      cache: 'no-store',
      headers: { Accept: 'application/json' },
    }, signal)
    if (response.status === 401) return { status: 'guest' }
    if (!response.ok) return { status: 'temporarily-unavailable' }

    const profile = parseSessionProfile(await response.json())
    return profile ? { status: 'authenticated', profile } : { status: 'temporarily-unavailable' }
  } catch {
    return { status: 'temporarily-unavailable' }
  }
}

async function fetchWithTimeout(url: string, init: RequestInit, externalSignal?: AbortSignal): Promise<Response> {
  const controller = new AbortController()
  const abortFromCaller = () => controller.abort()
  if (externalSignal?.aborted) controller.abort()
  else externalSignal?.addEventListener('abort', abortFromCaller, { once: true })
  const timeout = window.setTimeout(() => controller.abort(), SESSION_REQUEST_TIMEOUT_MS)
  try {
    return await fetch(url, { ...init, signal: controller.signal })
  } finally {
    window.clearTimeout(timeout)
    externalSignal?.removeEventListener('abort', abortFromCaller)
  }
}

export async function discoverLandingSession(force = false): Promise<ResolvedSessionState> {
  if (!force && cachedState && Date.now() - cachedAt < SESSION_FRESHNESS_MS) return cachedState
  if (inFlightRequest) return inFlightRequest.promise

  const controller = new AbortController()
  const pendingRequest = { promise: requestSession(controller.signal), controller }
  inFlightRequest = pendingRequest
  try {
    const state = await pendingRequest.promise
    if (controller.signal.aborted) return state
    if (state.status !== 'temporarily-unavailable') {
      cachedState = state
      cachedAt = Date.now()
    } else {
      cachedState = null
    }
    return state
  } finally {
    if (inFlightRequest === pendingRequest) inFlightRequest = null
  }
}

export function cancelLandingSessionRequest() {
  const activeRequest = inFlightRequest
  inFlightRequest = null
  activeRequest?.controller.abort()
}

export async function signOutLandingSession(): Promise<void> {
  cancelLandingSessionRequest()
  cachedState = null
  cachedAt = 0
  const response = await fetchWithTimeout(`${getConfig().apiUrl}/auth/logout`, {
    method: 'POST',
    credentials: 'include',
    cache: 'no-store',
  })
  if (!response.ok) throw new Error('The server did not confirm sign out')

  cachedState = { status: 'guest' }
  cachedAt = Date.now()
  publishSessionChange()
}

export function subscribeToSessionChanges(onChange: () => void): () => void {
  if (typeof BroadcastChannel === 'undefined') return () => undefined
  const channel = new BroadcastChannel(AUTH_SESSION_CHANNEL)
  channel.onmessage = event => {
    if (isRecord(event.data) && event.data.type === 'revalidate') onChange()
  }
  return () => channel.close()
}

function publishSessionChange() {
  if (typeof BroadcastChannel === 'undefined') return
  const channel = new BroadcastChannel(AUTH_SESSION_CHANNEL)
  channel.postMessage({ type: 'revalidate' })
  channel.close()
}
