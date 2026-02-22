import { AlertCircle, Home } from 'lucide-react';
import Link from 'next/link';

export default function NotFound() {
  // Use English as default for not-found page
  const locale = 'en';
  const isRTL = false;

  const homeUrl = `/${locale}`;
  const menuUrl = `/${locale}/menu`;

  return (
    <div
      className="flex items-center justify-center min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800"
      dir={isRTL ? 'rtl' : 'ltr'}
    >
      <div className="max-w-md w-full mx-auto px-4 text-center">
        {/* Icon */}
        <div className="mb-6">
          <AlertCircle className="h-20 w-20 text-red-500 mx-auto opacity-80" />
        </div>

        {/* Error Code */}
        <div className="mb-4">
          <h1 className="text-6xl font-bold text-gray-900 dark:text-white mb-2">
            404
          </h1>
          <p className="text-xl font-semibold text-gray-700 dark:text-gray-300">
            {locale === 'en' ? 'Page Not Found' : 'الصفحة غير موجودة'}
          </p>
        </div>

        {/* Description */}
        <p className="text-gray-600 dark:text-gray-400 mb-8 leading-relaxed">
          {locale === 'en'
            ? "The page you're looking for doesn't exist or has been moved. Let's get you back on track."
            : 'الصفحة التي تبحث عنها غير موجودة أو تم نقلها. دعنا نعيدك إلى المسار الصحيح.'}
        </p>

        {/* Action Buttons */}
        <div className="flex flex-col gap-3 sm:flex-row sm:justify-center">
          <Link
            href={homeUrl}
            className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition-colors"
          >
            <Home className="h-5 w-5" />
            {locale === 'en' ? 'Back to Home' : 'العودة للرئيسية'}
          </Link>

          <Link
            href={menuUrl}
            className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-gray-200 text-gray-900 font-medium rounded-lg hover:bg-gray-300 transition-colors dark:bg-gray-700 dark:text-white dark:hover:bg-gray-600"
          >
            {locale === 'en' ? 'View Menu' : 'اعرض القائمة'}
          </Link>
        </div>

        {/* Additional Help */}
        <div className="mt-12 pt-8 border-t border-gray-200 dark:border-gray-700">
          <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
            {locale === 'en' ? 'Need help?' : 'هل تحتاج إلى مساعدة؟'}
          </p>
          <Link
            href={`/${locale}#contact`}
            className="text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-medium transition-colors"
          >
            {locale === 'en' ? 'Contact our support' : 'تواصل مع الدعم الفني'}
          </Link>
        </div>
      </div>
    </div>
  );
}
