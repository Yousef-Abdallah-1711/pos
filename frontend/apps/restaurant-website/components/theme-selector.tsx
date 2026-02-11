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

  if (compact) {
    // Compact dropdown view
    return (
      <div className="flex items-center gap-2">
        <Palette className="h-5 w-5" style={{ color: textColor }} />
        <select
          value={currentTheme?.slug || ""}
          onChange={(e) => handleThemeSelect(e.target.value)}
          disabled={loading}
          style={{
            backgroundColor: bgColor,
            color: textColor,
            borderColor: borderColor,
          }}
          className="px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 disabled:opacity-50"
        >
          <option value="">{t.selectTheme}</option>
          {themes.map((theme) => (
            <option key={theme.slug} value={theme.slug}>
              {theme.name}
            </option>
          ))}
        </select>
      </div>
    );
  }

  // Full grid view
  return (
    <div className="w-full">
      <div className="mb-6">
        <h2
          className="text-2xl font-bold mb-2"
          style={{ color: textColor }}
        >
          {t.selectTheme}
        </h2>
        <p
          className="text-sm"
          style={{ color: `hsl(var(--theme-text) / 0.6)` }}
        >
          {t.description}
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {themes.map((theme) => {
          const isSelected = currentTheme?.slug === theme.slug;

          return (
            <button
              key={theme.slug}
              onClick={() => handleThemeSelect(theme.slug)}
              disabled={loading || isSelected}
              className="relative p-4 rounded-lg border-2 transition-all text-left group hover:scale-105"
              style={{
                borderColor: isSelected ? theme.colors.primary : borderColor,
                backgroundColor: bgColor,
              }}
              onMouseEnter={(e) => {
                if (!isSelected) {
                  e.currentTarget.style.boxShadow =
                    "0 4px 12px rgba(0, 0, 0, 0.1)";
                }
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.boxShadow = "none";
              }}
            >
              {/* Selected Badge */}
              {isSelected && (
                <div
                  className="absolute top-2 right-2 rounded-full p-1"
                  style={{ backgroundColor: theme.colors.primary }}
                >
                  <Check className="h-4 w-4" style={{ color: "#ffffff" }} />
                </div>
              )}

              {/* Color Preview */}
              <div className="mb-3 flex gap-2">
                <div
                  className="h-8 w-8 rounded"
                  style={{ backgroundColor: theme.colors.primary }}
                  title="Primary"
                />
                <div
                  className="h-8 w-8 rounded"
                  style={{ backgroundColor: theme.colors.secondary }}
                  title="Secondary"
                />
                <div
                  className="h-8 w-8 rounded"
                  style={{ backgroundColor: theme.colors.accent }}
                  title="Accent"
                />
              </div>

              {/* Theme Name */}
              <h3
                className="font-semibold mb-1"
                style={{ color: textColor }}
              >
                {theme.name}
              </h3>

              {/* Theme Description */}
              <p
                className="text-sm"
                style={{ color: `hsl(var(--theme-text) / 0.6)` }}
              >
                {theme.description}
              </p>

              {/* Loading state */}
              {loading && isSelected && (
                <div className="absolute inset-0 bg-black/10 rounded-lg flex items-center justify-center">
                  <div
                    className="w-5 h-5 border-2 border-transparent rounded-full animate-spin"
                    style={{
                      borderTopColor: theme.colors.primary,
                      borderRightColor: theme.colors.primary,
                    }}
                  />
                </div>
              )}
            </button>
          );
        })}
      </div>
    </div>
  );
}
