import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('@/lib/config', () => ({ getConfig: () => ({ apiUrl: 'https://api.example.test/api/v1' }) }))

const successfulProfile = {
  id: 42,
  name: 'Sample Manager',
  email: 'manager@example.test',
  tenant_id: 8,
  restaurant_id: 3,
  restaurant_name: 'Sample Cafe',
  role: 'manager',
  is_super_admin: false,
}

function response(status: number, body?: unknown) {
  return { status, ok: status >= 200 && status < 300, json: async () => body }
}

async function loadSessionClient() {
  vi.resetModules()
  return import('./auth-session')
}

describe('landing session client', () => {
  beforeEach(() => vi.stubGlobal('fetch', vi.fn()))
  afterEach(() => {
    vi.useRealTimers()
    vi.unstubAllGlobals()
  })

  it('requests a minimal no-store session with API credentials and returns only display fields', async () => {
    const fetchMock = vi.mocked(fetch).mockResolvedValue(response(200, successfulProfile) as Response)
    const { discoverLandingSession } = await loadSessionClient()

    await expect(discoverLandingSession()).resolves.toEqual({
      status: 'authenticated',
      profile: { id: 42, name: 'Sample Manager', email: 'manager@example.test', restaurantName: 'Sample Cafe' },
    })
    expect(fetchMock).toHaveBeenCalledWith('https://api.example.test/api/v1/auth/session', expect.objectContaining({
      credentials: 'include',
      cache: 'no-store',
      headers: { Accept: 'application/json' },
      signal: expect.any(AbortSignal),
    }))
  })

  it('keeps expired sessions separate from temporary API failures', async () => {
    const fetchMock = vi.mocked(fetch)
      .mockResolvedValueOnce(response(401) as Response)
      .mockResolvedValueOnce(response(500) as Response)
    const { discoverLandingSession } = await loadSessionClient()

    await expect(discoverLandingSession(true)).resolves.toEqual({ status: 'guest' })
    await expect(discoverLandingSession(true)).resolves.toEqual({ status: 'temporarily-unavailable' })
    expect(fetchMock).toHaveBeenCalledTimes(2)
  })

  it('treats an aborted timeout as unavailable instead of an expired session', async () => {
    vi.useFakeTimers()
    const fetchMock = vi.mocked(fetch).mockImplementation((_url, init) => new Promise((_resolve, reject) => {
      init?.signal?.addEventListener('abort', () => reject(new DOMException('Aborted', 'AbortError')), { once: true })
    }))
    const { discoverLandingSession } = await loadSessionClient()
    const pending = discoverLandingSession(true)

    await vi.advanceTimersByTimeAsync(5_000)

    await expect(pending).resolves.toEqual({ status: 'temporarily-unavailable' })
    expect(fetchMock).toHaveBeenCalledOnce()
  })

  it('rejects privileged or malformed session data instead of showing it as a staff account', async () => {
    vi.mocked(fetch).mockResolvedValue(response(200, { ...successfulProfile, is_super_admin: true }) as Response)
    const { discoverLandingSession } = await loadSessionClient()

    await expect(discoverLandingSession(true)).resolves.toEqual({ status: 'temporarily-unavailable' })
  })

  it('rejects inconsistent user identifiers in the session DTO', async () => {
    vi.mocked(fetch).mockResolvedValue(response(200, { ...successfulProfile, user_id: 43 }) as Response)
    const { discoverLandingSession } = await loadSessionClient()

    await expect(discoverLandingSession(true)).resolves.toEqual({ status: 'temporarily-unavailable' })
  })

  it('coalesces concurrent session checks', async () => {
    let resolveResponse: (value: Response) => void = () => undefined
    const pending = new Promise<Response>(resolve => { resolveResponse = resolve })
    const fetchMock = vi.mocked(fetch).mockReturnValue(pending)
    const { discoverLandingSession } = await loadSessionClient()

    const firstCheck = discoverLandingSession(true)
    const secondCheck = discoverLandingSession(true)
    expect(fetchMock).toHaveBeenCalledOnce()

    resolveResponse(response(200, successfulProfile) as Response)
    await expect(Promise.all([firstCheck, secondCheck])).resolves.toEqual([
      { status: 'authenticated', profile: { id: 42, name: 'Sample Manager', email: 'manager@example.test', restaurantName: 'Sample Cafe' } },
      { status: 'authenticated', profile: { id: 42, name: 'Sample Manager', email: 'manager@example.test', restaurantName: 'Sample Cafe' } },
    ])
  })

  it('cancels an in-flight session request when its portal island unmounts', async () => {
    const fetchMock = vi.mocked(fetch).mockImplementation((_url, init) => new Promise((_resolve, reject) => {
      init?.signal?.addEventListener('abort', () => reject(new DOMException('Aborted', 'AbortError')), { once: true })
    }))
    const { cancelLandingSessionRequest, discoverLandingSession } = await loadSessionClient()
    const pending = discoverLandingSession(true)

    cancelLandingSessionRequest()
    await expect(pending).resolves.toEqual({ status: 'temporarily-unavailable' })
    expect(fetchMock.mock.calls[0][1]?.signal?.aborted).toBe(true)
  })

  it('only marks sign out complete after the API confirms revocation', async () => {
    const fetchMock = vi.mocked(fetch).mockResolvedValue(response(503) as Response)
    const { signOutLandingSession } = await loadSessionClient()

    await expect(signOutLandingSession()).rejects.toThrow('The server did not confirm sign out')
    expect(fetchMock).toHaveBeenCalledWith('https://api.example.test/api/v1/auth/logout', expect.objectContaining({
      method: 'POST',
      credentials: 'include',
      cache: 'no-store',
      signal: expect.any(AbortSignal),
    }))
  })
})
