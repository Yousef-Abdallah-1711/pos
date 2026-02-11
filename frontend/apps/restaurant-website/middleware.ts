import { NextRequest, NextResponse } from 'next/server';
import { checkTenantStatus, shouldSkipStatusCheck } from '@/lib/checkTenantStatus';

const DEBUG = process.env.DEBUG_LOGS === 'true';
const debugLog = (...args: unknown[]) => { if (DEBUG) console.log(...args); };
const debugError = (...args: unknown[]) => { if (DEBUG) console.error(...args); };

/**
 * Middleware to extract restaurant context from request
 * - Extracts restaurant slug from subdomain (e.g., demo.localhost:3000 → demo)
 * - Extracts locale from URL pathname (e.g., /en/menu → en)
 * - Checks tenant status (trial expired, disabled, etc.)
 * - Redirects to maintenance page if tenant is blocked
 * - Stores both in request headers for use in server components
 */
export async function middleware(request: NextRequest) {
  const host = request.headers.get('host') || '';
  const pathname = request.nextUrl.pathname;

  debugLog(`[Middleware] Processing request - host: ${host}, path: ${pathname}`);

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
    debugLog(`[Middleware] Redirecting to normalized locale: ${validLocale}`);
    const normalizedParts = [...pathParts];
    normalizedParts[1] = validLocale;
    const normalizedPath = normalizedParts.join('/') || '/';
    const redirectURL = request.nextUrl.clone();
    redirectURL.pathname = normalizedPath.startsWith('/') ? normalizedPath : `/${normalizedPath}`;
    return NextResponse.redirect(redirectURL, 308);
  }

  // Check if this path should skip status check
  if (shouldSkipStatusCheck(pathname)) {
    debugLog(`[Middleware] Skipping status check for path: ${pathname}`);
    const requestHeaders = new Headers(request.headers);
    requestHeaders.set('x-restaurant-slug', restaurantSlug);
    requestHeaders.set('x-restaurant-locale', validLocale);

    return NextResponse.next({
      request: {
        headers: requestHeaders,
      },
    });
  }

  // Check tenant status only for restaurant sites with valid slugs
  let blockReason = null;
  if (restaurantSlug && restaurantSlug !== 'localhost' && restaurantSlug !== 'www') {
    debugLog(`[Middleware] Checking tenant status for slug: ${restaurantSlug}`);

    try {
      const status = await checkTenantStatus(restaurantSlug);

      debugLog(
        `[Middleware] Tenant status for ${restaurantSlug}: canAccess=${status.canAccess}, reason=${status.blockReason}`
      );

      if (!status.canAccess) {
        debugLog(`[Middleware] Tenant blocked with reason: ${status.blockReason} - redirecting to maintenance`);
        blockReason = status.blockReason;

        // Redirect to maintenance page
        const maintenanceUrl = request.nextUrl.clone();
        maintenanceUrl.pathname = '/maintenance';

        // Pass block reason as header for maintenance page to use
        const response = NextResponse.redirect(maintenanceUrl);
        response.headers.set('x-block-reason', blockReason || 'maintenance');

        return response;
      }
    } catch (error) {
      debugError(`[Middleware] Error checking tenant status for ${restaurantSlug}:`, error);
      // On error, allow access (fail open to avoid blocking legitimate requests)
    }
  }

  // Store in request headers (server-side, cannot be tampered by client)
  const requestHeaders = new Headers(request.headers);
  requestHeaders.set('x-restaurant-slug', restaurantSlug);
  requestHeaders.set('x-restaurant-locale', validLocale);
  if (blockReason) {
    requestHeaders.set('x-block-reason', blockReason);
  }

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
