"use client";

import React, { useEffect, useState } from "react";
import { useAuthStore } from "@/stores/authStore";
import { superAdminAPI } from "@/lib/api";
import type { DashboardStats } from "@/types";
import { BarChart3, Users, Building2, Activity } from "lucide-react";
import { toast } from "sonner";

export default function DashboardPage() {
  const superAdmin = useAuthStore((state) => state.superAdmin);
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const data = await superAdminAPI.getDashboardStats();
        setStats(data);
      } catch (error) {
        toast.error("Failed to load dashboard statistics");
        console.error(error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchStats();
  }, []);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading statistics...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Welcome Section */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Welcome back!</h1>
        <p className="text-gray-600 mt-2">
          Here's an overview of your POS SaaS platform
        </p>
      </div>

      {/* Stats Grid */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {/* Total Tenants */}
          <StatCard
            title="Total Tenants"
            value={stats.total_tenants}
            icon={Building2}
            color="bg-blue-50"
            iconColor="text-blue-600"
          />

          {/* Total Users */}
          <StatCard
            title="Total Users"
            value={stats.total_users}
            icon={Users}
            color="bg-green-50"
            iconColor="text-green-600"
          />

          {/* Active Impersonations */}
          <StatCard
            title="Active Impersonations"
            value={stats.active_impersonations}
            icon={Activity}
            color="bg-orange-50"
            iconColor="text-orange-600"
          />

          {/* Total Audit Logs */}
          <StatCard
            title="Total Audit Logs"
            value={stats.total_audit_logs}
            icon={BarChart3}
            color="bg-purple-50"
            iconColor="text-purple-600"
          />
        </div>
      )}

      {/* Quick Actions */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">
          Quick Actions
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <QuickActionButton
            label="View All Tenants"
            href="/dashboard/tenants"
          />
          <QuickActionButton
            label="View All Users"
            href="/dashboard/users"
          />
          <QuickActionButton
            label="View Audit Logs"
            href="/dashboard/audit-logs"
          />
          <QuickActionButton
            label="Manage Settings"
            href="/dashboard/settings"
          />
        </div>
      </div>

      {/* Recent Activity */}
      {stats && stats.recent_activity.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            Recent Activity
          </h2>
          <div className="space-y-4">
            {stats.recent_activity.slice(0, 5).map((log) => (
              <div
                key={log.id}
                className="flex items-start justify-between p-4 bg-gray-50 rounded-lg"
              >
                <div>
                  <p className="font-medium text-gray-900">{log.action}</p>
                  <p className="text-sm text-gray-600">
                    from {log.ip_address}
                  </p>
                </div>
                <span className="text-xs text-gray-500">
                  {new Date(log.created_at).toLocaleDateString()}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Info Box */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-6">
        <h3 className="font-semibold text-blue-900 mb-2">
          📋 Platform Overview
        </h3>
        <p className="text-blue-800 text-sm">
          You have full control over the entire POS SaaS platform. Navigate using the sidebar to manage tenants, users, and monitor activities. Use the impersonation feature to access any tenant dashboard for support or troubleshooting.
        </p>
      </div>
    </div>
  );
}

function StatCard({
  title,
  value,
  icon: Icon,
  color,
  iconColor,
}: {
  title: string;
  value: number;
  icon: React.ReactNode;
  color: string;
  iconColor: string;
}) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className={`${color} w-12 h-12 rounded-lg flex items-center justify-center mb-4`}>
        <Icon className={`${iconColor}`} size={24} />
      </div>
      <p className="text-gray-600 text-sm font-medium">{title}</p>
      <p className="text-3xl font-bold text-gray-900 mt-2">{value}</p>
    </div>
  );
}

function QuickActionButton({
  label,
  href,
}: {
  label: string;
  href: string;
}) {
  return (
    <a
      href={href}
      className="block p-4 bg-gradient-to-br from-blue-50 to-blue-100 border border-blue-200 rounded-lg hover:shadow-lg transition-shadow text-center"
    >
      <span className="font-medium text-blue-900">{label}</span>
    </a>
  );
}
