const SENSITIVE_DETAIL_KEY =
  /(password|secret|token|authorization|cookie|credential|api[_-]?key)/i;

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function sanitizeAuditValue(value: unknown): unknown {
  if (Array.isArray(value)) {
    return value.map(sanitizeAuditValue);
  }

  if (!isRecord(value)) {
    return value;
  }

  return Object.fromEntries(
    Object.entries(value).map(([key, nestedValue]) => [
      key,
      SENSITIVE_DETAIL_KEY.test(key)
        ? "[REDACTED]"
        : sanitizeAuditValue(nestedValue),
    ])
  );
}

export function getSafeAuditDetails(
  details: unknown
): Record<string, unknown> | null {
  if (!isRecord(details)) {
    return null;
  }

  return sanitizeAuditValue(details) as Record<string, unknown>;
}
