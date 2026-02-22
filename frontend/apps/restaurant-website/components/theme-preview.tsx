"use client";

import { useState } from "react";
import { X } from "lucide-react";
import { debugError } from "@/lib/debug";
import { Button } from "@pos-saas/ui";
import { useThemeStore } from "@/lib/store/theme-store";

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
  const mutedTextColor = `${textColor}b3`;
  const subtleTextColor = `${textColor}99`;
  const softBgColor = `${bgColor}cc`;
  const inputBgColor = `${textColor}08`;
  const isApplied = currentTheme?.slug === theme.slug;

  return (
    <div className="theme-preview fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="preview-surface bg-white dark:bg-gray-900 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="preview-header flex items-center justify-between p-6 border-b sticky top-0">
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
            <h3 className="theme-name text-xl font-semibold mb-2">
              {theme.name}
            </h3>
            <p className="theme-description text-sm">
              {theme.description}
            </p>
          </div>

          {/* Color Palette */}
          <div>
            <h4 className="section-title font-semibold mb-4">
              {t.primaryColor}
            </h4>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              <div>
                <div className="palette-box palette-primary h-16 rounded-lg mb-2 border" />
                <p className="palette-value text-xs font-mono">
                  {theme.colors.primary}
                </p>
              </div>

              <div>
                <div className="palette-box palette-secondary h-16 rounded-lg mb-2 border" />
                <p className="palette-value text-xs font-mono">
                  {theme.colors.secondary}
                </p>
              </div>

              <div>
                <div className="palette-box palette-accent h-16 rounded-lg mb-2 border" />
                <p className="palette-value text-xs font-mono">
                  {theme.colors.accent}
                </p>
              </div>

              <div>
                <div className="palette-box palette-background h-16 rounded-lg mb-2 border" />
                <p className="palette-value text-xs font-mono">
                  {theme.colors.background}
                </p>
              </div>

              <div>
                <div className="palette-box palette-text h-16 rounded-lg mb-2 flex items-center justify-center text-center">
                  Aa
                </div>
                <p className="palette-value text-xs font-mono">
                  {theme.colors.text}
                </p>
              </div>

              <div>
                <div className="palette-box palette-border h-16 rounded-lg mb-2 border-2" />
                <p className="palette-value text-xs font-mono">
                  {theme.colors.border}
                </p>
              </div>
            </div>
          </div>

          {/* Sample Components */}
          <div>
            <h4 className="section-title font-semibold mb-4">
              {t.sampleComponents}
            </h4>

            <div className="space-y-4">
              {/* Sample Button */}
              <div>
                <button className="sample-button font-medium">
                  Sample Button
                </button>
              </div>

              {/* Sample Card */}
              <div className="sample-card p-4 rounded-lg border">
                <h5 className="sample-card-title font-semibold mb-2">
                  Sample Card
                </h5>
                <p className="sample-card-description text-sm">
                  This is a sample card with the theme colors applied.
                </p>
              </div>

              {/* Sample Input */}
              <input
                type="text"
                placeholder="Sample input field"
                className="sample-input"
              />

              {/* Sample Badge */}
              <div className="flex gap-2 flex-wrap">
                {["New", "Featured", "Trending"].map((badge) => (
                  <span key={badge} className="sample-badge">
                    {badge}
                  </span>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="preview-footer flex justify-end gap-3 p-6 border-t">
          <Button
            variant="outline"
            onClick={onClose}
            disabled={isApplying}
          >
            {t.cancel}
          </Button>
          <Button
            onClick={handleApply}
            disabled={isApplying || isApplied}
            className={!isApplying && !isApplied ? "theme-preview-apply" : undefined}
          >
            {isApplying ? t.loading : t.applyTheme}
          </Button>
        </div>
      </div>
      <style jsx>{`
        .theme-preview {
          --preview-bg: ${bgColor};
          --preview-text: ${textColor};
          --preview-border: ${borderColor};
          --preview-primary: ${primaryColor};
          --preview-accent: ${theme.colors.accent};
          --preview-secondary: ${theme.colors.secondary};
          --preview-muted: ${mutedTextColor};
          --preview-subtle: ${subtleTextColor};
          --preview-soft-bg: ${softBgColor};
          --preview-input-bg: ${inputBgColor};
        }

        .preview-surface {
          background-color: var(--preview-bg);
          color: var(--preview-text);
        }

        .preview-header {
          border-bottom-color: var(--preview-border);
          background-color: var(--preview-bg);
        }

        .preview-footer {
          border-top-color: var(--preview-border);
        }

        .theme-name,
        .sample-card-title {
          color: var(--preview-primary);
        }

        .theme-description,
        .palette-value,
        .sample-card-description {
          color: var(--preview-subtle);
        }

        .section-title {
          color: var(--preview-text);
        }

        .palette-box {
          border-color: var(--preview-border);
        }

        .palette-primary {
          background-color: var(--preview-primary);
        }

        .palette-secondary {
          background-color: var(--preview-secondary);
        }

        .palette-accent {
          background-color: var(--preview-accent);
        }

        .palette-background {
          background-color: var(--preview-bg);
          border-width: 2px;
        }

        .palette-text {
          background-color: var(--preview-text);
          color: var(--preview-bg);
        }

        .palette-border {
          border-color: var(--preview-border);
        }

        .sample-button {
          background-color: var(--preview-primary);
          color: #ffffff;
          padding: 0.5rem 1rem;
          border-radius: 0.5rem;
        }

        .sample-card {
          border-color: var(--preview-border);
          background-color: var(--preview-soft-bg);
        }

        .sample-input {
          width: 100%;
          border: 1px solid var(--preview-border);
          border-radius: 0.375rem;
          padding: 0.5rem;
          background-color: var(--preview-input-bg);
          color: var(--preview-text);
        }

        .sample-badge {
          background-color: var(--preview-accent);
          color: var(--preview-text);
          padding: 0.25rem 0.75rem;
          border-radius: 9999px;
          font-size: 0.875rem;
          font-weight: 500;
        }

        :global(.theme-preview-apply) {
          background-color: var(--preview-primary) !important;
          border-color: var(--preview-primary) !important;
          color: #ffffff !important;
        }
      `}</style>
    </div>
  );
}
