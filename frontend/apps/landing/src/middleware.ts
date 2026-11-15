import { NextRequest, NextResponse } from 'next/server'
import { defaultLocale, isLocale } from '@pos-saas/i18n'

const ROOT_METADATA_ROUTES = ['/opengraph-image', '/robots.txt', '/sitemap.xml']

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl

  // Root metadata routes must remain at root so Next can serve their generated files.
  if (ROOT_METADATA_ROUTES.includes(pathname)) return NextResponse.next()

  // Skip Next.js internals and static assets.
  if (
    pathname.startsWith('/_next') ||
    pathname.startsWith('/api') ||
    pathname.startsWith('/favicon') ||
    pathname.includes('.')
  ) {
    return NextResponse.next()
  }

  // Legacy shortcuts must resolve before locale-prefix handling.
  const shortcut = pathname.match(/^\/(?:(en|ar)\/)?(login|register)\/?$/)
  if (shortcut) {
    const url = request.nextUrl.clone()
    url.pathname = `/${shortcut[1] || defaultLocale}/auth/${shortcut[2]}`
    return NextResponse.redirect(url)
  }

  // Check if the path already has a valid locale prefix
  const pathnameSegments = pathname.split('/')
  const firstSegment = pathnameSegments[1]

  if (isLocale(firstSegment)) {
    // Already has a valid locale — proceed
    return NextResponse.next()
  }

  // Redirect to default locale
  const url = request.nextUrl.clone()
  url.pathname = `/${defaultLocale}${pathname === '/' ? '' : pathname}`
  return NextResponse.redirect(url)
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico).*)'],
}
