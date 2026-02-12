import Link from 'next/link';

/**
 * Root page (/) - renders directly, no redirects
 */
export default function RootPage() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-white dark:bg-gray-900">
      <div className="text-center max-w-md p-8">
        <h1 className="text-3xl font-bold mb-4 text-gray-900 dark:text-white">
          Restaurant Website
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mb-6">
          The restaurant website is running. Access via your restaurant subdomain to see your restaurant.
        </p>
        <div className="space-y-2 text-sm text-gray-500">
          <p>Development mode: localhost</p>
          <p>
            <Link href="/en" className="text-blue-600 hover:underline">English</Link>
            {' | '}
            <Link href="/ar" className="text-blue-600 hover:underline">Arabic</Link>
          </p>
        </div>
      </div>
    </div>
  );
}
