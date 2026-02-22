import { NextRequest, NextResponse } from 'next/server'

const SUPPORTED_LOCALES = ['en', 'ar']
const DEFAULT_LOCALE = 'en'

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl

  // Skip Next.js internals and static files
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

  if (SUPPORTED_LOCALES.includes(firstSegment)) {
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
