'use client';

import { useEffect } from 'react';
import { AlertCircle, RotateCcw, Home } from 'lucide-react';
import { debugError } from '@/lib/debug';
import Link from 'next/link';

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
  params?: {
    locale: 'en' | 'ar';
  };
}

export default function Error({ error, reset, params }: ErrorProps) {
  const locale = params?.locale || 'en';
  const isRTL = locale === 'ar';

  // Log error to console for debugging
  useEffect(() => {
    debugError('Page Error:', error);
  }, [error]);

  const homeUrl = `/${locale}`;
  const isDevelopment = process.env.NODE_ENV === 'development';

  return (
    <div
      className="flex items-center justify-center min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800"
      dir={isRTL ? 'rtl' : 'ltr'}
    >
      <div className="max-w-md w-full mx-auto px-4 text-center">
        {/* Icon */}
        <div className="mb-6">
          <AlertCircle className="h-20 w-20 text-amber-500 mx-auto opacity-80" />
        </div>

        {/* Error Code */}
        <div className="mb-4">
          <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-2">
            {locale === 'en' ? 'Oops! Something went wrong' : 'عذراً! حدث خطأ ما'}
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-400">
            {locale === 'en'
              ? 'We encountered an unexpected error.'
              : 'واجهنا خطأ غير متوقع.'}
          </p>
        </div>

        {/* Error Details (Dev Only) */}
        {isDevelopment && (
          <div className="mt-6 mb-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-start">
            <p className="text-xs font-mono text-red-700 dark:text-red-400 break-all">
              {error.message || 'Unknown error'}
            </p>
            {error.digest && (
              <p className="text-xs text-red-600 dark:text-red-500 mt-2">
                Digest: {error.digest}
              </p>
            )}
          </div>
        )}

        {/* Description */}
        <p className="text-gray-600 dark:text-gray-400 mb-8 leading-relaxed">
          {locale === 'en'
            ? 'Please try again or contact our support team if the problem persists.'
            : 'يرجى المحاولة مرة أخرى أو التواصل مع فريق الدعم الفني إذا استمرت المشكلة.'}
        </p>

        {/* Action Buttons */}
        <div className="flex flex-col gap-3 sm:flex-row sm:justify-center">
          <button
            onClick={() => reset()}
            className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition-colors"
          >
            <RotateCcw className="h-5 w-5" />
            {locale === 'en' ? 'Try Again' : 'حاول مرة أخرى'}
          </button>

          <Link
            href={homeUrl}
            className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-gray-200 text-gray-900 font-medium rounded-lg hover:bg-gray-300 transition-colors dark:bg-gray-700 dark:text-white dark:hover:bg-gray-600"
          >
            <Home className="h-5 w-5" />
            {locale === 'en' ? 'Go Home' : 'اذهب للرئيسية'}
          </Link>
        </div>

        {/* Additional Help */}
        <div className="mt-12 pt-8 border-t border-gray-200 dark:border-gray-700">
          <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
            {locale === 'en' ? 'Need assistance?' : 'هل تحتاج مساعدة؟'}
          </p>
          <Link
            href={`/${locale}#contact`}
            className="text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-medium transition-colors"
          >
            {locale === 'en' ? 'Contact Support' : 'تواصل مع الدعم'}
          </Link>
        </div>
      </div>
    </div>
  );
}
