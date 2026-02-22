import { NextRequest, NextResponse } from 'next/server';

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1';

/**
 * Middleware to extract restaurant context from request
 * - Extracts restaurant slug from subdomain
 * - Extracts locale from URL pathname
 * - For real subdomains: checks tenant status and redirects to maintenance if blocked
 * - For localhost (development): passes through without checks
 */
export async function middleware(request: NextRequest) {
  const host = request.headers.get('host') || '';
  const pathname = request.nextUrl.pathname;

  // Skip middleware for maintenance page (avoid redirect loop)
  if (pathname.startsWith('/maintenance')) {
    return NextResponse.next();
  }

  // Extract restaurant slug from subdomain
  const restaurantSlug = host.split('.')[0];
  const supportedLocales = ['en', 'ar'];

  // Redirect root path `/` to `/en` (default locale)
  const pathParts = pathname.split('/');
  const firstSegment = pathParts[1] || '';
  if (pathname === '/' || (!supportedLocales.includes(firstSegment.toLowerCase()) && !firstSegment.startsWith('_next') && firstSegment !== 'favicon.ico' && firstSegment !== 'public' && firstSegment !== 'api')) {
    const url = request.nextUrl.clone();
    url.pathname = `/en${pathname === '/' ? '' : pathname}`;
    return NextResponse.redirect(url);
  }

  // Extract locale from pathname
  const rawLocale = firstSegment || 'en';
  const validLocale = supportedLocales.includes(rawLocale.toLowerCase()) ? rawLocale.toLowerCase() : 'en';

  // Set headers
  const requestHeaders = new Headers(request.headers);
  requestHeaders.set('x-restaurant-slug', restaurantSlug);
  requestHeaders.set('x-restaurant-locale', validLocale);

  // For real subdomains (not localhost/www), check tenant status
  const isRealSubdomain = restaurantSlug !== 'localhost' && restaurantSlug !== 'www' && restaurantSlug !== '';
  if (isRealSubdomain) {
    try {
      const res = await fetch(`${API_BASE}/public/tenant-status/${restaurantSlug}`, {
        headers: { 'Content-Type': 'application/json' },
      });

      if (res.ok) {
        const json = await res.json();
        const data = json.data || json;

        if (!data.can_access) {
          // Tenant is blocked - redirect to maintenance page with contact info
          const blockReason = data.block_reason || 'maintenance';
          const maintenanceUrl = new URL('/maintenance', request.url);
          maintenanceUrl.searchParams.set('reason', blockReason);
          if (data.tenant_name) maintenanceUrl.searchParams.set('name', data.tenant_name);
          if (data.contact_email) maintenanceUrl.searchParams.set('email', data.contact_email);
          if (data.contact_phone) maintenanceUrl.searchParams.set('phone', data.contact_phone);
          return NextResponse.redirect(maintenanceUrl);
        }
      } else if (res.status === 404) {
        // Tenant not found - redirect to maintenance
        const maintenanceUrl = new URL('/maintenance', request.url);
        return NextResponse.redirect(maintenanceUrl);
      }
    } catch {
      // API unreachable - let the request through (graceful degradation)
    }
  }

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
