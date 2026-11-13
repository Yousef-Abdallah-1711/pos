"use client";

import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAdminLogout } from "@/hooks/useAdminLogout";
import {
  LayoutDashboard,
  Building2,
  Users,
  FileText,
  Webhook,
  Settings,
  LogOut,
} from "lucide-react";

const navigationItems = [
  {
    name: "Dashboard",
    href: "/dashboard",
    icon: LayoutDashboard,
  },
  {
    name: "Tenants",
    href: "/dashboard/tenants",
    icon: Building2,
  },
  {
    name: "Users",
    href: "/dashboard/users",
    icon: Users,
  },
  {
    name: "Audit Logs",
    href: "/dashboard/audit-logs",
    icon: FileText,
  },
  {
    name: "Payments",
    href: "/dashboard/payments",
    icon: Webhook,
  },
  {
    name: "Settings",
    href: "/dashboard/settings",
    icon: Settings,
  },
];

export function DashboardNav() {
  const pathname = usePathname();
  const { logout, pending, error } = useAdminLogout();

  return (
    <aside className="w-full bg-slate-900 text-white flex flex-col md:h-screen md:w-64 md:shrink-0">
      {/* Logo */}
      <div className="p-4 md:p-6 border-b border-slate-700">
        <h1 className="text-2xl font-bold">POS SaaS</h1>
        <p className="text-sm text-slate-400 mt-1">Super Admin</p>
      </div>

      {/* Navigation */}
      <nav className="flex gap-2 overflow-x-auto p-3 md:block md:flex-1 md:space-y-2 md:p-4">
        {navigationItems.map((item) => {
          const Icon = item.icon;
          const isActive = pathname === item.href || pathname.startsWith(item.href + "/");

          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex min-h-12 shrink-0 items-center gap-3 whitespace-nowrap px-4 py-2 rounded-lg transition-colors ${
                isActive
                  ? "bg-blue-600 text-white"
                  : "text-slate-300 hover:bg-slate-800"
              }`}
            >
              <Icon size={20} />
              <span>{item.name}</span>
            </Link>
          );
        })}
      </nav>

      {/* Logout */}
      <div className="p-4 border-t border-slate-700">
        <button
          type="button"
          onClick={logout}
          disabled={pending}
          className="flex min-h-12 w-full items-center gap-3 rounded-lg px-4 py-2 text-slate-300 transition-colors hover:bg-slate-800 disabled:opacity-50"
        >
          <LogOut size={20} />
          <span>{pending ? "Signing out..." : "Logout"}</span>
        </button>
        {error && <p role="alert" className="mt-2 text-sm text-red-300">{error}</p>}
      </div>
    </aside>
  );
}
