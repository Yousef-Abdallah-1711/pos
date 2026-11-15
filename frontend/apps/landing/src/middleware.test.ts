import { describe, expect, it } from 'vitest'
import { NextRequest } from 'next/server'
import { middleware } from './middleware'

describe('locale middleware', () => {
  it.each(['/opengraph-image', '/robots.txt', '/sitemap.xml'])('%s remains at the root', (pathname) => {
    const response = middleware(new NextRequest(`https://www.poos.example${pathname}`))

    expect(response.headers.get('location')).toBeNull()
  })

  it('redirects new public visitors to Arabic by default', () => {
    const response = middleware(new NextRequest('https://www.poos.example/'))

    expect(response.headers.get('location')).toBe('https://www.poos.example/ar')
  })

  it.each(['/ar', '/en'])('does not rewrite an explicit locale route %s', (pathname) => {
    const response = middleware(new NextRequest(`https://www.poos.example${pathname}`))
    expect(response.headers.get('location')).toBeNull()
  })
})

describe('authentication shortcuts', () => {
  it.each([
    ['/login', '/ar/auth/login'], ['/register', '/ar/auth/register'],
    ['/en/login', '/en/auth/login'], ['/ar/login', '/ar/auth/login'],
    ['/en/register', '/en/auth/register'], ['/ar/register', '/ar/auth/register'],
  ])('routes %s to the actual localized authentication page', (pathname, target) => {
    const response = middleware(new NextRequest(`https://pos.rakmyat.online${pathname}?next=dashboard`))
    expect(response.headers.get('location')).toBe(`https://pos.rakmyat.online${target}?next=dashboard`)
  })
})
