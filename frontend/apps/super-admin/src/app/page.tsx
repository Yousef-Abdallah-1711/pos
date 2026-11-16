"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/stores/authStore";
import { useSuperAdminText } from "@/components/providers/PreferencesProvider";

export default function Home() {
  const t = useSuperAdminText();
  const router = useRouter();
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated());

  useEffect(() => {
    if (isAuthenticated) {
      router.push("/dashboard");
    } else {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <h1 className="text-3xl font-bold mb-4">SERVIO · {t("Super Admin")}</h1>
        <p className="text-gray-500">{t("Redirecting...")}</p>
      </div>
    </div>
  );
}
