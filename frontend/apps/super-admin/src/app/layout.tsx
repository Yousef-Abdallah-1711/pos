import type { Metadata } from "next";
import { Toaster } from "sonner";
import { SuperAdminAuthProvider } from "@/components/providers/SuperAdminAuthProvider";
import "./globals.css";
import "@/lib/debug-info";

export const metadata: Metadata = {
  title: "POS SaaS - Super Admin Dashboard",
  description: "Super Admin Dashboard for POS SaaS Platform",
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
