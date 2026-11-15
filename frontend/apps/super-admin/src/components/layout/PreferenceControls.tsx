"use client";

import { Languages, Moon, Sun } from "lucide-react";
import { useSuperAdminPreferences } from "@/components/providers/PreferencesProvider";

export function PreferenceControls({ compact = false }: { compact?: boolean }) {
  const { locale, theme, setLocale, toggleTheme } = useSuperAdminPreferences();
  const nextLocale = locale === "en" ? "ar" : "en";
  const nextTheme = theme === "light" ? "dark" : "light";

  return (
    <div className={`flex items-center gap-2 ${compact ? "" : "shrink-0"}`}>
      <button
        type="button"
        onClick={toggleTheme}
        aria-label={locale === "ar" ? (nextTheme === "dark" ? "تفعيل الوضع الداكن" : "تفعيل الوضع الفاتح") : `Switch to ${nextTheme} mode`}
        title={locale === "ar" ? (nextTheme === "dark" ? "الوضع الداكن" : "الوضع الفاتح") : `${nextTheme} mode`}
        className="inline-flex min-h-12 min-w-12 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-700 transition-colors hover:bg-gray-50 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100 dark:hover:bg-gray-700"
      >
        {theme === "light" ? <Moon aria-hidden="true" size={18} /> : <Sun aria-hidden="true" size={18} />}
      </button>
      <button
        type="button"
        onClick={() => setLocale(nextLocale)}
        aria-label={locale === "ar" ? "التغيير إلى الإنجليزية" : "Switch to Arabic"}
        className="inline-flex min-h-12 items-center justify-center gap-2 rounded-lg border border-gray-200 bg-white px-3 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100 dark:hover:bg-gray-700"
      >
        <Languages aria-hidden="true" size={17} />
        <span>{nextLocale === "ar" ? "العربية" : "English"}</span>
      </button>
    </div>
  );
}
