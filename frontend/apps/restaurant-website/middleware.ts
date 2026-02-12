import { NextRequest, NextResponse } from 'next/server';

/**
 * Middleware to extract restaurant context from request
 * - Extracts restaurant slug from subdomain
 * - Extracts locale from URL pathname
 * - Stores both in request headers for use in server components
 * - NO REDIRECTS - just passes through with headers set
 */
export async function middleware(request: NextRequest) {
  const host = request.headers.get('host') || '';
  const pathname = request.nextUrl.pathname;

  // Extract restaurant slug from subdomain
  const restaurantSlug = host.split('.')[0];

  // Extract locale from pathname
  const pathParts = pathname.split('/');
  const rawLocale = pathParts[1] || 'en';
  const supportedLocales = ['en', 'ar'];
  const validLocale = supportedLocales.includes(rawLocale.toLowerCase()) ? rawLocale.toLowerCase() : 'en';

  // Set headers and pass through - no redirects
  const requestHeaders = new Headers(request.headers);
  requestHeaders.set('x-restaurant-slug', restaurantSlug);
  requestHeaders.set('x-restaurant-locale', validLocale);

  return NextResponse.next({
    request: {
      headers: requestHeaders,
    },
  });
}

// Configure which routes the middleware should run on
export const config = {
  matcher: [
    '/((?!_next/static|_next/image|favicon.ico|public).*)',
  ],
};
