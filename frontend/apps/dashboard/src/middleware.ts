import { NextRequest, NextResponse } from 'next/server'

// The root document does not receive nested locale params. Set trusted request
// metadata from the pathname, overwriting any incoming client-supplied header.
export function middleware(request: NextRequest) {
  const headers = new Headers(request.headers)
  headers.set('x-dashboard-locale', request.nextUrl.pathname.split('/')[1] === 'ar' ? 'ar' : 'en')
  return NextResponse.next({ request: { headers } })
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico).*)'],
}
