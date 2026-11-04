import { apiClient } from './api-client'
import { getThemeBySlug } from './theme-api'
import type { ThemeData } from '@/lib/store/theme-store'

jest.mock('./api-client', () => ({
  apiClient: {
    get: jest.fn(),
  },
}))

describe('getThemeBySlug', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('shares one in-flight read for concurrent consumers of the same theme', async () => {
    let resolveRequest!: (theme: ThemeData) => void
    const theme = {
      name: 'Warm Comfort',
      slug: 'warm-comfort',
      colors: {},
    } as ThemeData
    const pendingRequest = new Promise<ThemeData>((resolve) => {
      resolveRequest = resolve
    })
    const getMock = apiClient.get as jest.Mock
    getMock.mockReturnValue(pendingRequest)

    const first = getThemeBySlug('warm-comfort')
    const second = getThemeBySlug('warm-comfort')

    expect(first).toBe(second)
    expect(apiClient.get).toHaveBeenCalledTimes(1)

    resolveRequest(theme)
    await expect(Promise.all([first, second])).resolves.toEqual([theme, theme])
  })
})
