'use client'

import type { ThemeLocale } from '@/types/theme'

export function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

export function toRecord(value: unknown): Record<string, unknown> {
  return isRecord(value) ? value : {}
}

export function toArray<T = Record<string, unknown>>(value: unknown): T[] {
  return Array.isArray(value) ? (value as T[]) : []
}

export function toStringValue(value: unknown, fallback: string): string {
  if (typeof value === 'string' && value.trim().length > 0) {
    return value.trim()
  }
  if (typeof value === 'number' && Number.isFinite(value)) {
    return String(value)
  }
  return fallback
}

export function toNumberValue(value: unknown, fallback: number): number {
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value
  }
  if (typeof value === 'string') {
    const parsed = Number(value)
    if (Number.isFinite(parsed)) {
      return parsed
    }
  }
  return fallback
}

export function toBooleanValue(value: unknown, fallback = false): boolean {
  if (typeof value === 'boolean') {
    return value
  }
  return fallback
}

export function readLocalized(
  value: unknown,
  locale: ThemeLocale,
  fallback = ''
): string {
  if (typeof value === 'string') {
    const trimmed = value.trim()
    return trimmed.length > 0 ? trimmed : fallback
  }

  if (!isRecord(value)) {
    return fallback
  }

  const localized = value[locale]
  if (typeof localized === 'string' && localized.trim().length > 0) {
    return localized.trim()
  }

  const english = value.en
  if (typeof english === 'string' && english.trim().length > 0) {
    return english.trim()
  }

  const arabic = value.ar
  if (typeof arabic === 'string' && arabic.trim().length > 0) {
    return arabic.trim()
  }

  return fallback
}

