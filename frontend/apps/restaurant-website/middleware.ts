import { NextRequest, NextResponse } from 'next/server';

/**
 * Middleware to extract restaurant context from request
 * - Extracts restaurant slug from subdomain (e.g., demo.localhost:3000 → demo)
 * - Extracts locale from URL pathname (e.g., /en/menu → en)
 * - Stores both in request headers for use in server components
 */
export function middleware(request: NextRequest) {
  const host = request.headers.get('host') || '';
  const pathname = request.nextUrl.pathname;

  // Extract restaurant slug from subdomain
  // Examples:
  // - demo.localhost:3000 → demo
  // - myrestaurant.com → myrestaurant
  // - restaurant.example.com → restaurant
  const restaurantSlug = host.split('.')[0];

  // Extract locale from pathname
  // Examples:
  // - /en/menu → en
  // - /ar/menu → ar
  // - /en → en
  const pathParts = pathname.split('/');
  const rawLocale = pathParts[1] || 'en'; // Default to 'en' if not specified
  const normalizedLocale = rawLocale.toLowerCase();

  // Validate locale is one of the supported ones
  const supportedLocales = ['en', 'ar'];
  const validLocale = supportedLocales.includes(normalizedLocale) ? normalizedLocale : 'en';

  // Redirect uppercase or unsupported locales to a canonical, lowercase path
  if (rawLocale !== validLocale) {
    const normalizedParts = [...pathParts];
    normalizedParts[1] = validLocale;
    const normalizedPath = normalizedParts.join('/') || '/';
    const redirectURL = request.nextUrl.clone();
    redirectURL.pathname = normalizedPath.startsWith('/') ? normalizedPath : `/${normalizedPath}`;
    return NextResponse.redirect(redirectURL, 308);
  }

  // Store in request headers (server-side, cannot be tampered by client)
  const requestHeaders = new Headers(request.headers);
  requestHeaders.set('x-restaurant-slug', restaurantSlug);
  requestHeaders.set('x-restaurant-locale', validLocale);

  // Continue to the next middleware/route with updated headers
  return NextResponse.next({
    request: {
      headers: requestHeaders,
    },
  });
}

// Configure which routes the middleware should run on
export const config = {
  matcher: [
    /*
     * Match all request paths except:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder files
     */
    '/((?!_next/static|_next/image|favicon.ico|public).*)',
  ],
};
