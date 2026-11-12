'use client';

import { AlertCircle } from 'lucide-react';

interface SuperAdminErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function SuperAdminError({ reset }: SuperAdminErrorProps) {
  return (
    <main className="flex min-h-[60vh] items-center justify-center p-4 sm:p-6">
      <section role="alert" className="w-full max-w-xl space-y-4 rounded-xl border border-red-200 bg-white p-6 text-center shadow-sm dark:border-red-900 dark:bg-gray-900">
        <AlertCircle aria-hidden="true" className="mx-auto text-red-700 dark:text-red-300" size={28} />
        <h1 className="text-xl font-semibold text-gray-900 dark:text-gray-100">This page could not load.</h1>
        <p className="text-sm text-gray-700 dark:text-gray-300">Try again to continue. If the problem persists, contact your system administrator.</p>
        <button
          type="button"
          onClick={reset}
          className="min-h-12 min-w-28 rounded-md bg-blue-600 px-4 py-2 font-medium text-white hover:bg-blue-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
        >
          Try again
        </button>
      </section>
    </main>
  );
}
