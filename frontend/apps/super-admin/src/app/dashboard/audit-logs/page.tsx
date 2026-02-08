"use client";

import React, { useEffect, useState } from "react";
import { superAdminAPI } from "@/lib/api";
import type { AuditLog, PaginatedResponse } from "@/types";
import { AlertCircle, Shield } from "lucide-react";
import { toast } from "sonner";

const ACTION_COLORS: Record<string, string> = {
  LOGIN_SUCCESS: "bg-green-100 text-green-700",
  LOGIN_FAILED: "bg-red-100 text-red-700",
  LOGOUT: "bg-blue-100 text-blue-700",
  IMPERSONATE_START: "bg-orange-100 text-orange-700",
  IMPERSONATE_END: "bg-purple-100 text-purple-700",
  TENANT_DISABLED: "bg-red-100 text-red-700",
  TENANT_ENABLED: "bg-green-100 text-green-700",
  IP_WHITELIST_UPDATED: "bg-yellow-100 text-yellow-700",
};

export default function AuditLogsPage() {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [filterAction, setFilterAction] = useState("");
  const limit = 50;

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        setIsLoading(true);
        const data = await superAdminAPI.getAuditLogs(
          page,
          limit,
          filterAction || undefined
        );
        const logsArray = Array.isArray(data) ? data : (data?.data || []);
        const totalCount = typeof data?.total === 'number' ? data.total : 0;
        setLogs(logsArray);
        setTotal(totalCount);
      } catch (error) {
        console.error("[AuditLogsPage] Error fetching logs:", error);
        toast.error("Failed to load audit logs");
        setLogs([]);
      } finally {
        setIsLoading(false);
      }
    };

    fetchLogs();
  }, [page, filterAction]);

  const totalPages = Math.ceil(total / limit);
  const uniqueActions = logs && logs.length > 0 ? Array.from(new Set(logs.map((log) => log.action))) : [];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Audit Logs</h1>
        <p className="text-gray-600 mt-2">
          Complete activity log of all super admin actions
        </p>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-4">
        <div className="flex items-center gap-4">
          <label className="text-sm font-medium text-gray-700">
            Filter by Action:
          </label>
          <select
            value={filterAction}
            onChange={(e) => {
              setFilterAction(e.target.value);
              setPage(1);
            }}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
          >
            <option value="">All Actions</option>
            {uniqueActions.map((action) => (
              <option key={action} value={action}>
                {action}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Audit Logs List */}
      {isLoading ? (
        <div className="flex items-center justify-center p-12">
          <div className="text-center">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">Loading audit logs...</p>
          </div>
        </div>
      ) : !logs || logs.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-12 text-center">
          <AlertCircle className="mx-auto text-gray-400" size={48} />
          <p className="text-gray-600 mt-4">No audit logs found</p>
        </div>
      ) : (
        <>
          <div className="space-y-4">
            {logs.map((log) => (
              <div
                key={log.id}
                className="bg-white rounded-lg shadow p-6 border-l-4 border-blue-500"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <Shield className="text-blue-600" size={24} />
                    <div>
                      <h3 className="font-semibold text-gray-900">
                        {log.action}
                      </h3>
                      <span
                        className={`inline-block mt-1 px-2 py-1 rounded text-xs font-medium ${
                          ACTION_COLORS[log.action] ||
                          "bg-gray-100 text-gray-700"
                        }`}
                      >
                        {log.action}
                      </span>
                    </div>
                  </div>
                  <span className="text-sm text-gray-500">
                    {new Date(log.created_at).toLocaleString()}
                  </span>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 bg-gray-50 rounded p-4">
                  <div>
                    <p className="text-xs text-gray-600 font-medium uppercase">
                      IP Address
                    </p>
                    <p className="text-sm text-gray-900 mt-1 font-mono">
                      {log.ip_address}
                    </p>
                  </div>
                  {log.target_user_id && (
                    <div>
                      <p className="text-xs text-gray-600 font-medium uppercase">
                        Target User ID
                      </p>
                      <p className="text-sm text-gray-900 mt-1">
                        {log.target_user_id}
                      </p>
                    </div>
                  )}
                  {log.target_tenant_id && (
                    <div>
                      <p className="text-xs text-gray-600 font-medium uppercase">
                        Target Tenant ID
                      </p>
                      <p className="text-sm text-gray-900 mt-1">
                        {log.target_tenant_id}
                      </p>
                    </div>
                  )}
                  {Object.keys(log.details).length > 0 && (
                    <div>
                      <p className="text-xs text-gray-600 font-medium uppercase">
                        Details
                      </p>
                      <p className="text-xs text-gray-900 mt-1 font-mono">
                        {JSON.stringify(log.details)}
                      </p>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between bg-white rounded-lg shadow p-4">
              <p className="text-sm text-gray-600">
                Showing {(page - 1) * limit + 1} to{" "}
                {Math.min(page * limit, total)} of {total} logs
              </p>
              <div className="flex gap-2">
                <button
                  onClick={() => setPage(Math.max(1, page - 1))}
                  disabled={page === 1}
                  className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Previous
                </button>
                <span className="px-4 py-2 text-sm font-medium text-gray-900">
                  Page {page} of {totalPages}
                </span>
                <button
                  onClick={() => setPage(Math.min(totalPages, page + 1))}
                  disabled={page === totalPages}
                  className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
