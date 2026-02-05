// Debug information and status checks for the Super Admin Dashboard

export function logDebugInfo() {
  console.log("=== POS SaaS Super Admin Dashboard - Debug Info ===");
  console.log("Environment Variables:");
  console.log("- API URL:", process.env.NEXT_PUBLIC_SUPER_ADMIN_API_URL);
  console.log("- App Name:", process.env.NEXT_PUBLIC_APP_NAME);

  console.log("\nLocalStorage:");
  if (typeof window !== "undefined") {
    const token = localStorage.getItem("super-admin-token");
    const authStore = localStorage.getItem("super-admin-auth-store");
    console.log("- Token exists:", !!token);
    console.log("- Token length:", token?.length);
    console.log("- Auth store exists:", !!authStore);
  }

  console.log("\nAPI Endpoints that will be used:");
  const baseUrl = process.env.NEXT_PUBLIC_SUPER_ADMIN_API_URL || "http://localhost:8080/api/v1/super-admin";
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
