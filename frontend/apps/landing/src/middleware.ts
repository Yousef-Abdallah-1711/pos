import { NextRequest, NextResponse } from 'next/server'

const SUPPORTED_LOCALES = new Set(['en', 'ar'])
const DEFAULT_LOCALE = 'en'
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

  // Check if the path already has a valid locale prefix
  const pathnameSegments = pathname.split('/')
  const firstSegment = pathnameSegments[1]

  if (SUPPORTED_LOCALES.has(firstSegment)) {
    // Already has a valid locale — proceed
    return NextResponse.next()
  }

  // Redirect to default locale
  const url = request.nextUrl.clone()
  url.pathname = `/${DEFAULT_LOCALE}${pathname === '/' ? '' : pathname}`
  return NextResponse.redirect(url)
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico).*)'],
}
