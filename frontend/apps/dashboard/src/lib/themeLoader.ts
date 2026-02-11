/**
 * Theme Loader Utility - Loads and processes theme JSON files
 */

import { ThemeJson, BilingualText, ParsedThemeJson } from "@/types/themeJson"
import { ThemeData, ComponentConfig } from "@/types/theme"
import { debugError } from '@/lib/debug'

export function extractText(
  value: BilingualText | undefined, 
  locale: "en" | "ar" = "en"
): string {
  if (!value) return ""
  if (typeof value === "string") return value
  const obj = value as any
  return obj[locale] || obj.en || ""
}

export function extractComponents(theme: ThemeJson | any): ComponentConfig[] {
  if (!theme.components || !Array.isArray(theme.components)) {
    return []
  }

  return theme.components.map((comp: any, index: number) => ({
    id: comp.id || "component-" + index,
    type: comp.type || "custom",
    title: extractText(comp.title, "en"),
    displayOrder: comp.order || index,
    enabled: comp.enabled !== false,
    config: comp.config || {},
  }))
}

export function extractColors(theme: ThemeJson | any) {
  return theme.colors || {
    primary: "#3b82f6",
    secondary: "#10b981",
    accent: "#f59e0b",
    background: "#ffffff",
    text: "#000000",
    border: "#e5e7eb",
    shadow: "#000000",
  }
}

export function extractTypography(theme: ThemeJson | any) {
  return theme.typography || {
    fontFamily: "Inter, system-ui, sans-serif",
    baseFontSize: 16,
    lineHeight: 1.6,
    borderRadius: 8,
  }
}

export function loadThemeFromJson(jsonData: any): ParsedThemeJson | null {
  try {
    let themeJson = jsonData
    if (typeof jsonData === "string") {
      themeJson = JSON.parse(jsonData)
      return {
        ...themeJson,
        _rawConfig: jsonData,
      }
    }
    return themeJson as ParsedThemeJson
  } catch (error) {
    debugError("Failed to parse theme JSON:", error)
    return null
  }
}

export function getComponentCount(theme: ThemeJson | any): number {
  if (!theme.components || !Array.isArray(theme.components)) {
    return 0
  }
  return theme.components.filter((c: any) => c.enabled !== false).length
}

export function getThemeCategory(theme: ThemeJson | any): string {
  const meta = theme.meta || {}
  return meta.category || "professional"
}

export function getThemeTags(theme: ThemeJson | any): string[] {
  const meta = theme.meta || {}
  return meta.tags || []
}

export function isBilingual(theme: ThemeJson | any): boolean {
  const meta = theme.meta || {}
  return meta.bilingual === true
}

export default {
  extractText,
  extractComponents,
  extractColors,
  extractTypography,
  loadThemeFromJson,
  getComponentCount,
  getThemeCategory,
  getThemeTags,
  isBilingual,
}
