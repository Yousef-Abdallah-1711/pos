"use client";

import { useState } from "react";
import { Palette, Check } from "lucide-react";
import { debugError } from "@/lib/debug";
import { useThemeStore } from "@/lib/store/theme-store";
import { useThemeColors } from "@/lib/hooks/use-theme";

interface Theme {
  slug: string;
  name: string;
  description: string;
  colors: {
    primary: string;
    secondary: string;
    accent: string;
  };
}

interface ThemeSelectorProps {
  themes?: Theme[];
  locale?: "en" | "ar";
  onThemeChange?: (themeSlug: string) => void;
  compact?: boolean;
}

/**
 * Theme selector component for choosing between different themes
 * Displays theme previews and handles theme switching
 */
export function ThemeSelector({
  themes = [],
  locale = "en",
  onThemeChange,
  compact = false,
}: ThemeSelectorProps) {
  const { currentTheme, loadTheme } = useThemeStore();
  const themeColors = useThemeColors() as any;
  const isRTL = locale === "ar";

  const [loading, setLoading] = useState(false);

  const handleThemeSelect = async (themeSlug: string) => {
    setLoading(true);
    try {
      await loadTheme(themeSlug);
      onThemeChange?.(themeSlug);
    } catch (error) {
      debugError("Failed to load theme:", error);
    } finally {
      setLoading(false);
    }
  };

  // Fallback colors
  const bgColor = themeColors?.background || "#ffffff";
  const textColor = themeColors?.text || "#1f2937";
  const borderColor = themeColors?.border || "#e5e7eb";
  const mutedTextColor = "hsl(var(--theme-text) / 0.6)";

  const labels = {
    en: {
      selectTheme: "Select Theme",
      themes: "Themes",
      description: "Choose a theme to personalize your experience",
    },
    ar: {
      selectTheme: "اختر المظهر",
      themes: "المظاهر",
      description: "اختر مظهراً لتخصيص تجربتك",
    },
  };

  const t = labels[locale as keyof typeof labels] || labels.en;
  const themeCardStyles = themes
    .map(
      (theme, index) => `
      .theme-card-${index}.is-selected {
        border-color: ${theme.colors.primary};
      }
      .theme-card-${index} .selected-badge {
        background-color: ${theme.colors.primary};
      }
      .theme-card-${index} .swatch-primary {
        background-color: ${theme.colors.primary};
      }
      .theme-card-${index} .swatch-secondary {
        background-color: ${theme.colors.secondary};
      }
      .theme-card-${index} .swatch-accent {
        background-color: ${theme.colors.accent};
      }
      .theme-card-${index} .loading-spinner {
        border-top-color: ${theme.colors.primary};
        border-right-color: ${theme.colors.primary};
      }
    `
    )
    .join("\n");

  if (compact) {
    // Compact dropdown view
    return (
      <div
        className={`theme-selector compact flex items-center gap-2 ${
          isRTL ? "rtl" : "ltr"
        }`}
      >
        <Palette className="selector-icon h-5 w-5" />
        <select
          value={currentTheme?.slug || ""}
          onChange={(e) => handleThemeSelect(e.target.value)}
          disabled={loading}
          className="theme-select px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 disabled:opacity-50"
        >
          <option value="">{t.selectTheme}</option>
          {themes.map((theme) => (
            <option key={theme.slug} value={theme.slug}>
              {theme.name}
            </option>
          ))}
        </select>
        <style jsx>{`
          .theme-selector {
            --selector-bg: ${bgColor};
            --selector-text: ${textColor};
            --selector-border: ${borderColor};
          }

          .theme-selector.rtl {
            direction: rtl;
          }

          .selector-icon {
            color: var(--selector-text);
          }

          .theme-select {
            background-color: var(--selector-bg);
            color: var(--selector-text);
            border-color: var(--selector-border);
          }
        `}</style>
      </div>
    );
  }

  // Full grid view
  return (
    <div className={`theme-selector full w-full ${isRTL ? "rtl" : "ltr"}`}>
      <div className="mb-6">
        <h2 className="selector-title text-2xl font-bold mb-2">
          {t.selectTheme}
        </h2>
        <p className="selector-description text-sm">
          {t.description}
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {themes.map((theme, index) => {
          const isSelected = currentTheme?.slug === theme.slug;

          return (
            <button
              key={theme.slug}
              onClick={() => handleThemeSelect(theme.slug)}
              disabled={loading || isSelected}
              className={`theme-card theme-card-${index} relative p-4 rounded-lg border-2 transition-all text-left hover:scale-105 ${
                isSelected ? "is-selected" : ""
              }`}
            >
              {/* Selected Badge */}
              {isSelected && (
                <div className="selected-badge absolute top-2 right-2 rounded-full p-1">
                  <Check className="selected-icon h-4 w-4" />
                </div>
              )}

              {/* Color Preview */}
              <div className="mb-3 flex gap-2">
                <div className="swatch swatch-primary h-8 w-8 rounded" title="Primary" />
                <div className="swatch swatch-secondary h-8 w-8 rounded" title="Secondary" />
                <div className="swatch swatch-accent h-8 w-8 rounded" title="Accent" />
              </div>

              {/* Theme Name */}
              <h3 className="theme-name font-semibold mb-1">
                {theme.name}
              </h3>

              {/* Theme Description */}
              <p className="theme-description text-sm">
                {theme.description}
              </p>

              {/* Loading state */}
              {loading && isSelected && (
                <div className="absolute inset-0 bg-black/10 rounded-lg flex items-center justify-center">
                  <div className="loading-spinner w-5 h-5 border-2 border-transparent rounded-full animate-spin" />
                </div>
              )}
            </button>
          );
        })}
      </div>
      <style jsx>{`
        .theme-selector {
          --selector-bg: ${bgColor};
          --selector-text: ${textColor};
          --selector-muted: ${mutedTextColor};
          --selector-border: ${borderColor};
        }

        .theme-selector.rtl {
          direction: rtl;
        }

        .selector-title,
        .theme-name {
          color: var(--selector-text);
        }

        .selector-description,
        .theme-description {
          color: var(--selector-muted);
        }

        .theme-card {
          border-color: var(--selector-border);
          background-color: var(--selector-bg);
          box-shadow: none;
        }

        .theme-card:not(.is-selected):not(:disabled):hover {
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }

        .theme-card:disabled {
          cursor: not-allowed;
        }

        .selected-icon {
          color: #ffffff;
        }

        ${themeCardStyles}
      `}</style>
    </div>
  );
}
