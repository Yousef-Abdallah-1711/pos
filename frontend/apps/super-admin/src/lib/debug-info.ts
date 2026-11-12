// Debug information and status checks for the Super Admin Dashboard
import { getSuperAdminApiBaseUrl } from "@/lib/config";

export function logDebugInfo() {
  console.log("=== POS SaaS Super Admin Dashboard - Debug Info ===");
  console.log("Environment Variables:");
  console.log("- API URL:", process.env.NEXT_PUBLIC_SUPER_ADMIN_API_URL);
  console.log("- App Name:", process.env.NEXT_PUBLIC_APP_NAME);

  console.log("\nBrowser authentication:");
  if (typeof window !== "undefined") {
    console.log("- Authentication uses an HttpOnly server session cookie");
  }

  console.log("\nAPI Endpoints that will be used:");
  let baseUrl = "";
  try {
    baseUrl = getSuperAdminApiBaseUrl();
  } catch (error) {
    console.warn("[SuperAdmin Debug] API base URL configuration error:", error);
  }
  console.log("- Dashboard stats:", `${baseUrl}/dashboard/stats`);
  console.log("- List tenants:", `${baseUrl}/tenants`);
  console.log("- List users:", `${baseUrl}/users`);
  console.log("- Audit logs:", `${baseUrl}/audit-logs`);

  console.log("\nOpen browser DevTools Console (F12) to see API calls");
  console.log("Look for [API] prefix in logs to track requests and responses");
}

// Call this on page load
if (typeof window !== "undefined") {
  logDebugInfo();
}
