"use client";

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { translateSuperAdmin } from "@/lib/copy";

export type SuperAdminLocale = "en" | "ar";
export type SuperAdminTheme = "light" | "dark";

type Preferences = {
  locale: SuperAdminLocale;
  theme: SuperAdminTheme;
  setLocale: (locale: SuperAdminLocale) => void;
  toggleTheme: () => void;
  t: (text: string) => string;
};

const DEFAULT_PREFERENCES: Preferences = {
  locale: "en",
  theme: "light",
  setLocale: () => undefined,
  toggleTheme: () => undefined,
  t: (text) => text,
};
const PreferencesContext = createContext<Preferences>(DEFAULT_PREFERENCES);
const STORAGE_KEY = "servio-super-admin-preferences";

export function PreferencesProvider({ children }: { children: React.ReactNode }) {
  const [locale, setLocale] = useState<SuperAdminLocale>("en");
  const [theme, setTheme] = useState<SuperAdminTheme>("light");
  const [ready, setReady] = useState(false);

  useEffect(() => {
    try {
      const saved = JSON.parse(localStorage.getItem(STORAGE_KEY) || "null") as Partial<Pick<Preferences, "locale" | "theme">> | null;
      if (saved?.locale === "en" || saved?.locale === "ar") setLocale(saved.locale);
      if (saved?.theme === "light" || saved?.theme === "dark") setTheme(saved.theme);
    } catch {
      try {
        localStorage.removeItem(STORAGE_KEY);
      } catch {
        // Browser storage can be disabled; the in-memory preference still works.
      }
    } finally {
      setReady(true);
    }
  }, []);

  useEffect(() => {
    if (!ready) return;
    const root = document.documentElement;
    root.lang = locale;
    root.dir = locale === "ar" ? "rtl" : "ltr";
    root.classList.toggle("dark", theme === "dark");
    root.style.colorScheme = theme;
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({ locale, theme }));
    } catch {
      // Keep the current session usable when browser storage is unavailable.
    }
  }, [locale, theme, ready]);

  const toggleTheme = useCallback(() => {
    setTheme((current) => current === "dark" ? "light" : "dark");
  }, []);

  const t = useCallback((text: string) => translateSuperAdmin(locale, text), [locale]);
  const value = useMemo(() => ({ locale, theme, setLocale, toggleTheme, t }), [locale, theme, toggleTheme, t]);

  return <PreferencesContext.Provider value={value}>{children}</PreferencesContext.Provider>;
}

export function useSuperAdminText() {
  return useSuperAdminPreferences().t;
}

export function useSuperAdminPreferences() {
  return useContext(PreferencesContext);
}
