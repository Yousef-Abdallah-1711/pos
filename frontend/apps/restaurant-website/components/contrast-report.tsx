"use client";

import { useMemo, useState } from "react";
import { CheckCircle2, AlertCircle, AlertTriangle } from "lucide-react";
import { getContrastRatio, meetsWCAGAA, meetsWCAGAAA } from "@/lib/utils/theme-colors";
import { useThemeColors } from "@/lib/hooks/use-theme";

interface ColorPair {
  foreground: string;
  background: string;
  name: string;
}

interface ContrastResult {
  name: string;
  ratio: number;
  level: "AAA" | "AA" | "Fail";
  percentage: number;
}

interface ContrastReportProps {
  themeName?: string;
  colors?: Record<string, string>;
  locale?: "en" | "ar";
  showDetails?: boolean;
}

function getLevelClass(level: ContrastResult["level"]): "level-aaa" | "level-aa" | "level-fail" {
  if (level === "AAA") return "level-aaa";
  if (level === "AA") return "level-aa";
  return "level-fail";
}

/**
 * Contrast Report Dashboard
 * Shows WCAG AA/AAA compliance for theme colors
 */
export function ContrastReport({
  themeName = "Current Theme",
  colors,
  locale = "en",
  showDetails = true,
}: ContrastReportProps) {
  const themeColors = useThemeColors();
  const [expandedPair, setExpandedPair] = useState<string | null>(null);

  const isRTL = locale === "ar";

  const finalColors = useMemo(
    () => (colors || themeColors || {}) as any,
    [colors, themeColors]
  );

  // Common color pairs to test
  const colorPairs: ColorPair[] = useMemo(
    () => [
      {
        foreground: finalColors.text || "#1f2937",
        background: finalColors.background || "#ffffff",
        name: locale === "en" ? "Text on Background" : "Ø§Ù„Ù†Øµ Ø¹Ù„Ù‰ Ø§Ù„Ø®Ù„ÙÙŠØ©",
      },
      {
        foreground: finalColors.primary || "#f97316",
        background: finalColors.background || "#ffffff",
        name: locale === "en" ? "Primary on Background" : "Ø§Ù„Ø£Ø³Ø§Ø³ÙŠ Ø¹Ù„Ù‰ Ø§Ù„Ø®Ù„ÙÙŠØ©",
      },
      {
        foreground: "#ffffff",
        background: finalColors.primary || "#f97316",
        name: locale === "en" ? "White on Primary" : "Ø£Ø¨ÙŠØ¶ Ø¹Ù„Ù‰ Ø§Ù„Ø£Ø³Ø§Ø³ÙŠ",
      },
      {
        foreground: finalColors.secondary || "#0ea5e9",
        background: finalColors.background || "#ffffff",
        name:
          locale === "en"
            ? "Secondary on Background"
            : "Ø§Ù„Ø«Ø§Ù†ÙˆÙŠ Ø¹Ù„Ù‰ Ø§Ù„Ø®Ù„ÙÙŠØ©",
      },
      {
        foreground: "#ffffff",
        background: finalColors.secondary || "#0ea5e9",
        name: locale === "en" ? "White on Secondary" : "Ø£Ø¨ÙŠØ¶ Ø¹Ù„Ù‰ Ø§Ù„Ø«Ø§Ù†ÙˆÙŠ",
      },
      {
        foreground: finalColors.accent || "#fbbf24",
        background: finalColors.background || "#ffffff",
        name: locale === "en" ? "Accent on Background" : "Ø§Ù„ØªÙ…ÙŠÙŠØ² Ø¹Ù„Ù‰ Ø§Ù„Ø®Ù„ÙÙŠØ©",
      },
      {
        foreground: finalColors.text || "#1f2937",
        background: finalColors.accent || "#fbbf24",
        name:
          locale === "en"
            ? "Text on Accent"
            : "Ø§Ù„Ù†Øµ Ø¹Ù„Ù‰ Ù„ÙˆÙ† Ø§Ù„ØªÙ…ÙŠÙŠØ²",
      },
      {
        foreground: finalColors.border || "#e5e7eb",
        background: finalColors.background || "#ffffff",
        name:
          locale === "en"
            ? "Border on Background"
            : "Ø§Ù„Ø­Ø¯ÙˆØ¯ Ø¹Ù„Ù‰ Ø§Ù„Ø®Ù„ÙÙŠØ©",
      },
    ],
    [finalColors, locale]
  );

  // Calculate contrast ratios and WCAG levels
  const results: ContrastResult[] = useMemo(
    () =>
      colorPairs.map((pair) => {
        const ratio = getContrastRatio(pair.foreground, pair.background);
        let level: "AAA" | "AA" | "Fail" = "Fail";

        if (meetsWCAGAAA(pair.foreground, pair.background)) {
          level = "AAA";
        } else if (meetsWCAGAA(pair.foreground, pair.background)) {
          level = "AA";
        }

        return {
          name: pair.name,
          ratio,
          level,
          percentage: 0,
        };
      }),
    [colorPairs]
  );

  // Calculate compliance summary
  const compliance = useMemo(
    () => ({
      AAA: results.filter((r) => r.level === "AAA").length,
      AA: results.filter((r) => r.level === "AA").length,
      fail: results.filter((r) => r.level === "Fail").length,
    }),
    [results]
  );

  const labels = {
    en: {
      contrastReport: "Contrast Report",
      wcagCompliance: "WCAG Compliance Summary",
      colorPairs: "Color Pairs",
      ratio: "Ratio",
      wcagLevel: "WCAG Level",
      passes: "Passes",
      passingPairs: "Passing Pairs",
      failingPairs: "Failing Pairs",
      meetsAAA: "Meets AAA (7:1)",
      meetsAA: "Meets AA (4.5:1)",
      doesNotMeet: "Does Not Meet",
      sampleText: "Sample Text",
      normal: "Normal Text",
      largeText: "Large Text (18pt+)",
    },
    ar: {
      contrastReport: "ØªÙ‚Ø±ÙŠØ± Ø§Ù„ØªØ¨Ø§ÙŠÙ†",
      wcagCompliance: "Ù…Ù„Ø®Øµ ØªÙˆØ§ÙÙ‚ WCAG",
      colorPairs: "Ø£Ø²ÙˆØ§Ø¬ Ø§Ù„Ø£Ù„ÙˆØ§Ù†",
      ratio: "Ø§Ù„Ù†Ø³Ø¨Ø©",
      wcagLevel: "Ù…Ø³ØªÙˆÙ‰ WCAG",
      passes: "ÙŠÙ…Ø±",
      passingPairs: "Ø§Ù„Ø£Ø²ÙˆØ§Ø¬ Ø§Ù„Ù†Ø§Ø¬Ø­Ø©",
      failingPairs: "Ø§Ù„Ø£Ø²ÙˆØ§Ø¬ Ø§Ù„ÙØ§Ø´Ù„Ø©",
      meetsAAA: "ÙŠØ³ØªÙˆÙÙŠ AAA (7:1)",
      meetsAA: "ÙŠØ³ØªÙˆÙÙŠ AA (4.5:1)",
      doesNotMeet: "Ù„Ø§ ÙŠØ³ØªÙˆÙÙŠ",
      sampleText: "Ù†Øµ Ø¹ÙŠÙ†Ø©",
      normal: "Ø§Ù„Ù†Øµ Ø§Ù„Ø¹Ø§Ø¯ÙŠ",
      largeText: "Ø§Ù„Ù†Øµ Ø§Ù„ÙƒØ¨ÙŠØ± (18pt+)",
    },
  };

  const t = labels[locale as keyof typeof labels] || labels.en;
  const allPassing = compliance.fail === 0;
  const baseTextColor = (themeColors as any)?.text || "#1f2937";
  const subtitleColor = "hsl(var(--theme-text) / 0.6)";
  const subtleColor = "hsl(var(--theme-text) / 0.5)";
  const mutedColor = "hsl(var(--theme-text) / 0.7)";

  const pairStyles = colorPairs
    .map(
      (pair, index) => `
      .result-${index} .sample-preview,
      .result-${index} .sample-preview-lg {
        color: ${pair.foreground};
        background-color: ${pair.background};
      }

      .result-${index} .foreground-swatch {
        background-color: ${pair.foreground};
      }

      .result-${index} .background-swatch {
        background-color: ${pair.background};
      }
    `
    )
    .join("\n");

  return (
    <div className={`contrast-report w-full space-y-6 ${isRTL ? "rtl" : "ltr"}`}>
      {/* Header */}
      <div>
        <h2 className="report-title text-2xl font-bold mb-2">{t.contrastReport}</h2>
        <p className="report-subtitle text-sm">{themeName}</p>
      </div>

      {/* Compliance Summary */}
      <div className="grid grid-cols-3 gap-4">
        <div className="summary-card summary-aaa p-4 rounded-lg border">
          <p className="summary-value text-2xl font-bold">{compliance.AAA}</p>
          <p className="summary-label text-xs">{t.meetsAAA}</p>
        </div>

        <div className="summary-card summary-aa p-4 rounded-lg border">
          <p className="summary-value text-2xl font-bold">{compliance.AA}</p>
          <p className="summary-label text-xs">{t.meetsAA}</p>
        </div>

        <div
          className={`summary-card summary-fail p-4 rounded-lg border ${
            compliance.fail === 0 ? "is-pass" : "is-fail"
          }`}
        >
          <p className="summary-value text-2xl font-bold">{compliance.fail}</p>
          <p className="summary-label text-xs">{t.failingPairs}</p>
        </div>
      </div>

      {/* Overall Status */}
      <div className={`overall-status p-4 rounded-lg border flex items-start gap-3 ${allPassing ? "pass" : "fail"}`}>
        {allPassing ? (
          <CheckCircle2 className="overall-icon pass h-5 w-5 flex-shrink-0 mt-0.5" />
        ) : (
          <AlertCircle className="overall-icon fail h-5 w-5 flex-shrink-0 mt-0.5" />
        )}
        <div>
          <h4 className={`overall-title font-semibold mb-1 ${allPassing ? "pass" : "fail"}`}>
            {allPassing
              ? locale === "en"
                ? "All Color Pairs Pass WCAG AA"
                : "Ø¬Ù…ÙŠØ¹ Ø£Ø²ÙˆØ§Ø¬ Ø§Ù„Ø£Ù„ÙˆØ§Ù† ØªØ³ØªÙˆÙÙŠ Ù…Ø¹Ø§ÙŠÙŠØ± WCAG AA"
              : locale === "en"
                ? "Some Color Pairs Do Not Meet WCAG AA"
                : "Ø¨Ø¹Ø¶ Ø£Ø²ÙˆØ§Ø¬ Ø§Ù„Ø£Ù„ÙˆØ§Ù† Ù„Ø§ ØªØ³ØªÙˆÙÙŠ Ù…Ø¹Ø§ÙŠÙŠØ± WCAG AA"}
          </h4>
          <p className={`overall-description text-sm ${allPassing ? "pass" : "fail"}`}>
            {allPassing
              ? locale === "en"
                ? "Your theme meets WCAG AA accessibility standards for all tested color combinations."
                : "ÙŠØ³ØªÙˆÙÙŠ Ù…ÙˆØ¶ÙˆØ¹Ùƒ Ù…Ø¹Ø§ÙŠÙŠØ± WCAG AA Ù„Ù„ÙˆØµÙˆÙ„ Ø§Ù„Ø´Ø§Ù…Ù„ Ù„Ø¬Ù…ÙŠØ¹ Ù…Ø¬Ù…ÙˆØ¹Ø§Øª Ø§Ù„Ø£Ù„ÙˆØ§Ù† Ø§Ù„Ù…Ø®ØªØ¨Ø±Ø©."
              : locale === "en"
                ? "Review the color pairs below and adjust colors to meet WCAG AA (4.5:1 contrast ratio)."
                : "Ø±Ø§Ø¬Ø¹ Ø£Ø²ÙˆØ§Ø¬ Ø§Ù„Ø£Ù„ÙˆØ§Ù† Ø£Ø¯Ù†Ø§Ù‡ ÙˆØ§Ø¶Ø¨Ø· Ø§Ù„Ø£Ù„ÙˆØ§Ù† Ù„Ø§Ø³ØªÙŠÙØ§Ø¡ Ù…Ø¹Ø§ÙŠÙŠØ± WCAG AA (Ù†Ø³Ø¨Ø© ØªØ¨Ø§ÙŠÙ† 4.5:1)."}
          </p>
        </div>
      </div>

      {/* Detailed Results */}
      {showDetails && (
        <div>
          <h3 className="details-title font-semibold mb-4">{t.colorPairs}</h3>

          <div className="space-y-2">
            {results.map((result, index) => {
              const pair = colorPairs[index];
              const isExpanded = expandedPair === result.name;
              const levelClass = getLevelClass(result.level);

              return (
                <div
                  key={result.name}
                  className={`result-item result-${index} ${levelClass} border rounded-lg overflow-hidden`}
                >
                  {/* Summary Row */}
                  <button
                    onClick={() => setExpandedPair(isExpanded ? null : result.name)}
                    className={`w-full p-4 flex items-center justify-between hover:bg-black/5 transition-colors ${
                      isRTL ? "text-right" : ""
                    }`}
                  >
                    <div className="flex items-center gap-3 flex-1">
                      {result.level === "Fail" ? (
                        <AlertTriangle className="status-icon fail h-5 w-5 flex-shrink-0" />
                      ) : (
                        <CheckCircle2
                          className={`status-icon ${
                            result.level === "AAA" ? "aaa" : "aa"
                          } h-5 w-5 flex-shrink-0`}
                        />
                      )}
                      <div className={isRTL ? "text-right" : "text-left"}>
                        <p className="result-name font-medium">{result.name}</p>
                        <p className="result-ratio text-xs">
                          {t.ratio}: {result.ratio.toFixed(2)}:1
                        </p>
                      </div>
                    </div>
                    <span className={`result-badge px-3 py-1 rounded-full text-xs font-semibold ${levelClass}`}>
                      {result.level}
                    </span>
                  </button>

                  {/* Details Row */}
                  {isExpanded && pair && (
                    <div className="details-row p-4 border-t space-y-4">
                      {/* Sample Text */}
                      <div>
                        <p className="details-label text-xs font-semibold mb-2">{t.sampleText}</p>
                        <p className="sample-preview text-sm mb-2">
                          {locale === "en"
                            ? "The quick brown fox jumps over the lazy dog"
                            : "Ø§Ù„Ø«Ø¹Ù„Ø¨ Ø§Ù„Ø¨Ù†ÙŠ Ø§Ù„Ø³Ø±ÙŠØ¹ ÙŠÙ‚ÙØ² ÙÙˆÙ‚ Ø§Ù„ÙƒÙ„Ø¨ Ø§Ù„ÙƒØ³ÙˆÙ„"}
                        </p>
                      </div>

                      {/* Large Text Sample */}
                      <div>
                        <p className="details-label text-xs font-semibold mb-2">{t.largeText}</p>
                        <p className="sample-preview-lg text-lg mb-2">Heading</p>
                      </div>

                      {/* Color Values */}
                      <div className="grid grid-cols-2 gap-4 text-xs">
                        <div>
                          <p className="details-label text-xs font-semibold mb-1">
                            {locale === "en" ? "Foreground" : "Ø§Ù„Ø£Ù…Ø§Ù…"}
                          </p>
                          <div className="color-swatch foreground-swatch h-8 rounded border" />
                          <p className="color-value mt-1 font-mono">{pair.foreground}</p>
                        </div>
                        <div>
                          <p className="details-label text-xs font-semibold mb-1">
                            {locale === "en" ? "Background" : "Ø§Ù„Ø®Ù„ÙÙŠØ©"}
                          </p>
                          <div className="color-swatch background-swatch h-8 rounded border" />
                          <p className="color-value mt-1 font-mono">{pair.background}</p>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>
      )}

      <style jsx>{`
        .contrast-report {
          --report-text: ${baseTextColor};
          --report-subtitle: ${subtitleColor};
          --report-subtle: ${subtleColor};
          --report-muted: ${mutedColor};
          --report-detail-border: hsl(var(--theme-text) / 0.12);
        }

        .contrast-report.rtl {
          direction: rtl;
        }

        .report-title,
        .details-title,
        .result-name {
          color: var(--report-text);
        }

        .report-subtitle {
          color: var(--report-subtitle);
        }

        .summary-aaa {
          background-color: #f0fdf4;
          border-color: #22c55e;
        }

        .summary-aaa .summary-value,
        .summary-aaa .summary-label {
          color: #16a34a;
        }

        .summary-aa {
          background-color: #fef3c7;
          border-color: #f59e0b;
        }

        .summary-aa .summary-value,
        .summary-aa .summary-label {
          color: #d97706;
        }

        .summary-fail.is-pass {
          background-color: #f0fdf4;
          border-color: #22c55e;
        }

        .summary-fail.is-pass .summary-value,
        .summary-fail.is-pass .summary-label {
          color: #16a34a;
        }

        .summary-fail.is-fail {
          background-color: #fee2e2;
          border-color: #ef4444;
        }

        .summary-fail.is-fail .summary-value,
        .summary-fail.is-fail .summary-label {
          color: #dc2626;
        }

        .overall-status.pass {
          background-color: #f0fdf4;
          border-color: #22c55e;
        }

        .overall-status.fail {
          background-color: #fee2e2;
          border-color: #ef4444;
        }

        .overall-icon.pass,
        .overall-title.pass,
        .overall-description.pass {
          color: #16a34a;
        }

        .overall-icon.fail,
        .overall-title.fail,
        .overall-description.fail {
          color: #dc2626;
        }

        .result-item.level-fail {
          border-color: #ef4444;
        }

        .result-item.level-aa {
          border-color: #f59e0b;
        }

        .result-item.level-aaa {
          border-color: #22c55e;
        }

        .status-icon.fail {
          color: #ef4444;
        }

        .status-icon.aa {
          color: #d97706;
        }

        .status-icon.aaa {
          color: #16a34a;
        }

        .result-ratio,
        .color-value {
          color: var(--report-subtle);
        }

        .result-badge.level-fail {
          background-color: #fee2e2;
          color: #dc2626;
        }

        .result-badge.level-aa {
          background-color: #fef3c7;
          color: #d97706;
        }

        .result-badge.level-aaa {
          background-color: #f0fdf4;
          color: #16a34a;
        }

        .details-row {
          background-color: hsl(var(--theme-background) / 0.3);
          border-top-color: var(--report-detail-border);
        }

        .details-label {
          color: var(--report-muted);
        }

        .sample-preview,
        .sample-preview-lg {
          padding: 12px;
          border-radius: 4px;
        }

        .color-swatch {
          border-color: var(--report-detail-border);
        }

        ${pairStyles}
      `}</style>
    </div>
  );
}
