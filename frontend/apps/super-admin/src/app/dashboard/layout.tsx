"use client";

import React, { ReactNode } from "react";
import { useAuthStore } from "@/stores/authStore";
import { DashboardNav } from "@/components/layout/DashboardNav";
import { DashboardHeader } from "@/components/layout/DashboardHeader";

interface DashboardLayoutProps {
  children: ReactNode;
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  const superAdmin = useAuthStore((state) => state.superAdmin);

  if (!superAdmin) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-gray-500">Loading dashboard...</p>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col bg-gray-100 md:h-screen md:flex-row md:overflow-hidden">
      {/* Sidebar Navigation */}
      <DashboardNav />

      {/* Main Content */}
      <div className="min-w-0 flex-1 md:overflow-auto">
        {/* Header */}
        <DashboardHeader />

        {/* Content */}
        <main className="p-4 md:p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
