import { redirect } from 'next/navigation';

/**
 * Root page - redirects to /en (default locale)
 * This ensures accessing localhost:3003/ properly routes to the [locale] segment
 */
export default function RootPage() {
  redirect('/en');
}
