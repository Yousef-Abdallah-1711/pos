import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export type DarkModePreference = 'light' | 'dark' | 'system';

export interface PreferencesSettings {
  language: 'en' | 'ar';
  theme: DarkModePreference;
  darkMode: DarkModePreference;
  primaryColor: string;
  secondaryColor: string;
  accentColor: string;
}

interface PreferencesState {
  settings: PreferencesSettings;
  setSettings: (settings: PreferencesSettings) => void;
  updateLanguage: (language: 'en' | 'ar') => void;
  updateTheme: (theme: DarkModePreference) => void;
  setDarkMode: (mode: DarkModePreference) => void;
  updateColors: (primaryColor: string, secondaryColor: string, accentColor: string) => void;
  applyTheme: () => void;
  applyColors: () => void;
}

const DEFAULT_SETTINGS: PreferencesSettings = {
  language: 'en',
  theme: 'system',
  darkMode: 'system',
  primaryColor: '#3b82f6',
  secondaryColor: '#10b981',
  accentColor: '#f59e0b',
};

/**
 * Apply theme to DOM
 */
function applyThemeToDom(theme: DarkModePreference) {
  if (typeof window === 'undefined') return;

  const htmlElement = document.documentElement;

  if (theme === 'system') {
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    htmlElement.classList.toggle('dark', isDark);
  } else {
    htmlElement.classList.toggle('dark', theme === 'dark');
  }
}

/**
 * Apply custom colors to DOM
 */
function applyColorsToDom(
  primaryColor: string,
  secondaryColor: string,
  accentColor: string
) {
  if (typeof window === 'undefined') return;

  const root = document.documentElement;
  root.style.setProperty('--color-primary', primaryColor);
  root.style.setProperty('--color-secondary', secondaryColor);
  root.style.setProperty('--color-accent', accentColor);
}

export const usePreferencesStore = create<PreferencesState>()(
  persist(
    (set, get) => ({
      settings: DEFAULT_SETTINGS,

      setSettings: (settings) => {
        set({ settings });
        applyThemeToDom(settings.darkMode ?? settings.theme);
        applyColorsToDom(settings.primaryColor, settings.secondaryColor, settings.accentColor);
      },

      updateLanguage: (language) => {
        const { settings } = get();
        const newSettings = { ...settings, language };
        get().setSettings(newSettings);
      },

      updateTheme: (theme) => {
        const { settings } = get();
        const newSettings = { ...settings, theme, darkMode: theme };
        get().setSettings(newSettings);
      },

      setDarkMode: (mode) => {
        const { settings } = get();
        const newSettings = { ...settings, darkMode: mode, theme: mode };
        get().setSettings(newSettings);
      },

      updateColors: (primaryColor, secondaryColor, accentColor) => {
        const { settings } = get();
        const newSettings = { ...settings, primaryColor, secondaryColor, accentColor };
        get().setSettings(newSettings);
      },

      applyTheme: () => {
        const { settings } = get();
        applyThemeToDom(settings.darkMode ?? settings.theme);
      },

      applyColors: () => {
        const { settings } = get();
        applyColorsToDom(settings.primaryColor, settings.secondaryColor, settings.accentColor);
      },
    }),
    {
      name: 'pos-restaurant-preferences',
    }
  )
);
