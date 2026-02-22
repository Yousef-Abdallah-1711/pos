import type { Metadata } from "next";
import { Toaster } from "sonner";
import { SuperAdminAuthProvider } from "@/components/providers/SuperAdminAuthProvider";
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
    <html lang="en">
      <body>
        <SuperAdminAuthProvider>
          {children}
        </SuperAdminAuthProvider>
        <Toaster />
      </body>
    </html>
  );
}
