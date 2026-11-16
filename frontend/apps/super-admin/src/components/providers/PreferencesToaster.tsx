"use client";

import { Toaster } from "sonner";
import { useSuperAdminPreferences } from "@/components/providers/PreferencesProvider";

export function PreferencesToaster() {
  const { locale } = useSuperAdminPreferences();
  return <Toaster position={locale === "ar" ? "top-left" : "top-right"} richColors />;
}
