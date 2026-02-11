"use client";

import { useState } from "react";
import { X, Eye } from "lucide-react";
import { debugError } from "@/lib/debug";
import { Button } from "@pos-saas/ui";
import { useThemeStore } from "@/lib/store/theme-store";
import { useThemeColors } from "@/lib/hooks/use-theme";
import { hexToHsl } from "@/lib/utils/theme-colors";

interface Theme {
  slug: string;
  name: string;
  description: string;
  colors: {
    primary: string;
    secondary: string;
    accent: string;
    background: string;
    text: string;
    border: string;
  };
}

interface ThemePreviewProps {
  theme: Theme;
  onApply?: (themeSlug: string) => void;
  onClose?: () => void;
  locale?: "en" | "ar";
  isOpen?: boolean;
}

/**
 * Theme preview modal showing theme colors and sample components
 */
export function ThemePreview({
  theme,
  onApply,
  onClose,
  locale = "en",
  isOpen = true,
}: ThemePreviewProps) {
  const { currentTheme, loadTheme } = useThemeStore();
  const [isApplying, setIsApplying] = useState(false);
  const isRTL = locale === "ar";

  const labels = {
    en: {
      preview: "Theme Preview",
      primaryColor: "Primary Color",
      secondaryColor: "Secondary Color",
      accentColor: "Accent Color",
      backgroundColor: "Background Color",
      textColor: "Text Color",
      borderColor: "Border Color",
      sampleComponents: "Sample Components",
      applyTheme: "Apply Theme",
      cancel: "Cancel",
      loading: "Applying...",
    },
    ar: {
      preview: "معاينة المظهر",
      primaryColor: "اللون الأساسي",
      secondaryColor: "اللون الثانوي",
      accentColor: "لون التمييز",
      backgroundColor: "لون الخلفية",
      textColor: "لون النص",
      borderColor: "لون الحدود",
      sampleComponents: "نماذج المكونات",
      applyTheme: "تطبيق المظهر",
      cancel: "إلغاء",
      loading: "جاري التطبيق...",
    },
  };

  const t = labels[locale as keyof typeof labels] || labels.en;

  const handleApply = async () => {
    setIsApplying(true);
    try {
      await loadTheme(theme.slug);
      onApply?.(theme.slug);
    } catch (error) {
      debugError("Failed to apply theme:", error);
    } finally {
      setIsApplying(false);
    }
  };

  if (!isOpen) return null;

  const bgColor = theme.colors.background;
  const textColor = theme.colors.text;
  const borderColor = theme.colors.border;
  const primaryColor = theme.colors.primary;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div
        className="bg-white dark:bg-gray-900 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto"
        style={{
          backgroundColor: bgColor,
          color: textColor,
        }}
      >
        {/* Header */}
        <div
          className="flex items-center justify-between p-6 border-b sticky top-0"
          style={{
            borderBottomColor: borderColor,
          }}
        >
          <h2 className="text-2xl font-bold">{t.preview}</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-black/10 rounded-lg transition-colors"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-8">
          {/* Theme Info */}
          <div>
            <h3
              className="text-xl font-semibold mb-2"
              style={{ color: primaryColor }}
            >
              {theme.name}
            </h3>
            <p
              className="text-sm"
              style={{ color: `${textColor}b3` }}
            >
              {theme.description}
            </p>
          </div>

          {/* Color Palette */}
          <div>
            <h4
              className="font-semibold mb-4"
              style={{ color: textColor }}
            >
              {t.primaryColor}
            </h4>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              <div>
                <div
                  className="h-16 rounded-lg mb-2 border"
                  style={{
                    backgroundColor: theme.colors.primary,
                    borderColor: borderColor,
                  }}
                />
                <p
                  className="text-xs font-mono"
                  style={{ color: `${textColor}99` }}
                >
                  {theme.colors.primary}
                </p>
              </div>

              <div>
                <div
                  className="h-16 rounded-lg mb-2 border"
                  style={{
                    backgroundColor: theme.colors.secondary,
                    borderColor: borderColor,
                  }}
                />
                <p
                  className="text-xs font-mono"
                  style={{ color: `${textColor}99` }}
                >
                  {theme.colors.secondary}
                </p>
              </div>

              <div>
                <div
                  className="h-16 rounded-lg mb-2 border"
                  style={{
                    backgroundColor: theme.colors.accent,
                    borderColor: borderColor,
                  }}
                />
                <p
                  className="text-xs font-mono"
                  style={{ color: `${textColor}99` }}
                >
                  {theme.colors.accent}
                </p>
              </div>

              <div>
                <div
                  className="h-16 rounded-lg mb-2 border"
                  style={{
                    backgroundColor: theme.colors.background,
                    borderColor: borderColor,
                    border: `2px solid ${borderColor}`,
                  }}
                />
                <p
                  className="text-xs font-mono"
                  style={{ color: `${textColor}99` }}
                >
                  {theme.colors.background}
                </p>
              </div>

              <div>
                <div
                  className="h-16 rounded-lg mb-2 flex items-center justify-center text-center"
                  style={{
                    backgroundColor: textColor,
                    color: bgColor,
                  }}
                >
                  Aa
                </div>
                <p
                  className="text-xs font-mono"
                  style={{ color: `${textColor}99` }}
                >
                  {theme.colors.text}
                </p>
              </div>

              <div>
                <div
                  className="h-16 rounded-lg mb-2 border-2"
                  style={{
                    borderColor: theme.colors.border,
                  }}
                />
                <p
                  className="text-xs font-mono"
                  style={{ color: `${textColor}99` }}
                >
                  {theme.colors.border}
                </p>
              </div>
            </div>
          </div>

          {/* Sample Components */}
          <div>
            <h4
              className="font-semibold mb-4"
              style={{ color: textColor }}
            >
              {t.sampleComponents}
            </h4>

            <div className="space-y-4">
              {/* Sample Button */}
              <div>
                <button
                  style={{
                    backgroundColor: primaryColor,
                    color: "#ffffff",
                    padding: "0.5rem 1rem",
                    borderRadius: "0.5rem",
                  }}
                  className="font-medium"
                >
                  Sample Button
                </button>
              </div>

              {/* Sample Card */}
              <div
                className="p-4 rounded-lg border"
                style={{
                  borderColor: borderColor,
                  backgroundColor: `${bgColor}cc`,
                }}
              >
                <h5
                  className="font-semibold mb-2"
                  style={{ color: primaryColor }}
                >
                  Sample Card
                </h5>
                <p
                  className="text-sm"
                  style={{ color: `${textColor}99` }}
                >
                  This is a sample card with the theme colors applied.
                </p>
              </div>

              {/* Sample Input */}
              <input
                type="text"
                placeholder="Sample input field"
                style={{
                  backgroundColor: `${textColor}08`,
                  color: textColor,
                  borderColor: borderColor,
                  borderWidth: "1px",
                  padding: "0.5rem",
                  borderRadius: "0.375rem",
                  width: "100%",
                }}
              />

              {/* Sample Badge */}
              <div className="flex gap-2 flex-wrap">
                {["New", "Featured", "Trending"].map((badge) => (
                  <span
                    key={badge}
                    style={{
                      backgroundColor: theme.colors.accent,
                      color: textColor,
                      padding: "0.25rem 0.75rem",
                      borderRadius: "9999px",
                      fontSize: "0.875rem",
                      fontWeight: "500",
                    }}
                  >
                    {badge}
                  </span>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div
          className="flex justify-end gap-3 p-6 border-t"
          style={{
            borderTopColor: borderColor,
          }}
        >
          <Button
            variant="outline"
            onClick={onClose}
            disabled={isApplying}
          >
            {t.cancel}
          </Button>
          <Button
            onClick={handleApply}
            disabled={isApplying || currentTheme?.slug === theme.slug}
            style={
              !isApplying && currentTheme?.slug !== theme.slug
                ? {
                    backgroundColor: primaryColor,
                    color: "#ffffff",
                    borderColor: primaryColor,
                  }
                : undefined
            }
          >
            {isApplying ? t.loading : t.applyTheme}
          </Button>
        </div>
      </div>
    </div>
  );
}
