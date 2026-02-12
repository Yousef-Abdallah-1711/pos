import { redirect } from 'next/navigation';

/**
 * Root page (/) - redirects to /en (default locale)
 *
 * This page handles requests to the root path and ensures they're
 * directed to a valid [locale] segment route.
 *
 * Next.js redirect() is used instead of middleware redirect to ensure
 * proper routing through the [locale] dynamic segment.
 */
export default function RootPage() {
  // Use Next.js server-side redirect to /en
  // This happens during server-side rendering, not in middleware
  redirect('/en');
}
