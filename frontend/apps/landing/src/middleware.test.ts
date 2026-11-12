import { describe, expect, it } from 'vitest'
import { NextRequest } from 'next/server'
import { middleware } from './middleware'

describe('locale middleware', () => {
  it.each(['/opengraph-image', '/robots.txt', '/sitemap.xml'])('%s remains at the root', (pathname) => {
    const response = middleware(new NextRequest(`https://www.poos.example${pathname}`))

    expect(response.headers.get('location')).toBeNull()
  })

  it('redirects the root to the default locale', () => {
    const response = middleware(new NextRequest('https://www.poos.example/'))

    expect(response.headers.get('location')).toBe('https://www.poos.example/en')
  })

  it('does not rewrite locale-prefixed pages', () => {
    const response = middleware(new NextRequest('https://www.poos.example/ar'))

    expect(response.headers.get('location')).toBeNull()
  })
})
