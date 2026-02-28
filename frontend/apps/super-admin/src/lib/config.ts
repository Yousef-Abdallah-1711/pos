const DEV_DEFAULTS = {
  superAdminApiUrl: "http://localhost:8080/api/v1/super-admin",
  dashboardUrl: "http://localhost:3002",
} as const;

function isProductionRuntime() {
  return process.env.NODE_ENV === "production";
}

function normalizeUrl(value: string) {
  return value.replace(/\/+$/, "");
}

function resolveUrlValue(
  label: string,
  configuredValue: string | undefined,
  devFallback: string
) {
  const trimmed = configuredValue?.trim();
  if (trimmed) {
    return normalizeUrl(trimmed);
  }

  if (!isProductionRuntime()) {
    return normalizeUrl(devFallback);
  }

  const key =
    label === "Super Admin API URL"
      ? "NEXT_PUBLIC_SUPER_ADMIN_API_URL"
      : "NEXT_PUBLIC_DASHBOARD_URL";

  throw new Error(`[super-admin/config] Missing ${label} in production. Set ${key}.`);
}

export function getSuperAdminApiBaseUrl() {
  return resolveUrlValue(
    "Super Admin API URL",
    process.env.NEXT_PUBLIC_SUPER_ADMIN_API_URL,
    DEV_DEFAULTS.superAdminApiUrl
  );
}

export function getDashboardBaseUrl() {
  return resolveUrlValue(
    "Dashboard URL",
    process.env.NEXT_PUBLIC_DASHBOARD_URL,
    DEV_DEFAULTS.dashboardUrl
  );
}
