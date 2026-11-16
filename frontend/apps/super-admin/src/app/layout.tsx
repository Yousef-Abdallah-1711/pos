import type { Metadata } from "next";
import { SuperAdminAuthProvider } from "@/components/providers/SuperAdminAuthProvider";
import { PreferencesProvider } from "@/components/providers/PreferencesProvider";
import { PreferencesToaster } from "@/components/providers/PreferencesToaster";
import "./globals.css";

export const metadata: Metadata = {
  title: "POS SaaS - Super Admin Dashboard",
  description: "Super Admin Dashboard for POS SaaS Platform",
  icons: {
    icon: "/icon.svg",
    shortcut: "/icon.svg",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" dir="ltr" suppressHydrationWarning>
      <body>
        <PreferencesProvider>
          <SuperAdminAuthProvider>
            {children}
          </SuperAdminAuthProvider>
          <PreferencesToaster />
        </PreferencesProvider>
      </body>
    </html>
  );
}
