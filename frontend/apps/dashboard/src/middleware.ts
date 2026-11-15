import { NextRequest, NextResponse } from 'next/server'
import { defaultLocale } from '@pos-saas/i18n'

// The root document does not receive nested locale params. Set trusted request
// metadata from the pathname, overwriting any incoming client-supplied header.
export function middleware(request: NextRequest) {
  const headers = new Headers(request.headers)
  const requestedLocale = request.nextUrl.pathname.split('/')[1]
  headers.set('x-dashboard-locale', requestedLocale === 'en' || requestedLocale === 'ar' ? requestedLocale : defaultLocale)
  return NextResponse.next({ request: { headers } })
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico).*)'],
}
