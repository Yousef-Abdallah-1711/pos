"use client";

import React, { useCallback, useEffect, useMemo, useRef, useState } from "react";
import axios from "axios";
import { superAdminAPI } from "@/lib/api";
import type { PaymentMetrics, WebhookHealthResponse } from "@/types";
import { toast } from "sonner";
import { Activity, Clock3, RefreshCw, WalletCards, Webhook } from "lucide-react";

type Period = "today" | "week" | "month" | "all";

const PERIOD_OPTIONS: Array<{ value: Period; label: string }> = [
  { value: "today", label: "Today" },
  { value: "week", label: "Last 7 Days" },
  { value: "month", label: "Last 30 Days" },
  { value: "all", label: "All Time" },
];

function formatDateTime(value?: string) {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleString();
}

function statusBadgeClass(status: string) {
  const normalized = status.toLowerCase();
  if (normalized.includes("success")) return "bg-green-100 text-green-700";
  if (
    normalized.includes("failed") ||
    normalized.includes("invalid") ||
    normalized.includes("reject")
  ) {
    return "bg-red-100 text-red-700";
  }
  if (normalized.includes("processing") || normalized.includes("pending")) {
    return "bg-yellow-100 text-yellow-700";
  }
  return "bg-gray-100 text-gray-700";
}

function summaryCardClass(theme: "blue" | "green" | "orange" | "purple") {
  switch (theme) {
    case "green":
      return "bg-green-50 border-green-200";
    case "orange":
      return "bg-orange-50 border-orange-200";
    case "purple":
      return "bg-purple-50 border-purple-200";
    default:
      return "bg-blue-50 border-blue-200";
  }
}

function getErrorMessage(error: unknown) {
  if (axios.isAxiosError(error)) {
    if (!error.response) {
      return "Cannot connect to backend API. Start backend on http://localhost:8080 and try again.";
    }
    const apiMessage =
      (error.response.data as { message?: string; error?: string } | undefined)?.message ||
      (error.response.data as { message?: string; error?: string } | undefined)?.error;
    return apiMessage || `Request failed with status ${error.response.status}`;
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return "Failed to load payment monitoring data";
}

export default function PaymentsMonitoringPage() {
  const [period, setPeriod] = useState<Period>("month");
  const [metrics, setMetrics] = useState<PaymentMetrics | null>(null);
  const [health, setHealth] = useState<WebhookHealthResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [loadError, setLoadError] = useState<string | null>(null);
  const lastErrorRef = useRef<string | null>(null);

  const loadMonitoringData = useCallback(
    async (targetPeriod: Period, silent = false) => {
      if (!silent) {
        setIsLoading(true);
      } else {
        setIsRefreshing(true);
      }

      try {
        const [metricsData, healthData] = await Promise.all([
          superAdminAPI.getPaymentMetrics(targetPeriod),
          superAdminAPI.getWebhookHealth(targetPeriod, 50),
        ]);
        setMetrics(metricsData);
        setHealth(healthData);
        setLoadError(null);
        lastErrorRef.current = null;
      } catch (error) {
        const message = getErrorMessage(error);
        setLoadError(message);
        if (lastErrorRef.current !== message) {
          toast.error(message);
          lastErrorRef.current = message;
        }
        console.error("[PaymentsMonitoringPage] Failed to load data:", message);
      } finally {
        setIsLoading(false);
        setIsRefreshing(false);
      }
    },
    []
  );

  useEffect(() => {
    void loadMonitoringData(period);
  }, [period, loadMonitoringData]);

  const retryStats = useMemo(() => {
    return (
      health?.retry_queue || {
        pending: 0,
        processing: 0,
        success: 0,
        failed: 0,
      }
    );
  }, [health]);

  const totalRetryQueue =
    retryStats.pending + retryStats.processing + retryStats.success + retryStats.failed;

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[300px]">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-10 w-10 border-b-2 border-blue-600" />
          <p className="mt-3 text-gray-600">Loading payment monitoring...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Payments & Webhook Health</h1>
          <p className="text-gray-600 mt-1">
            Monitor payment outcomes, webhook processing, and retry queue reliability.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <select
            value={period}
            onChange={(e) => setPeriod(e.target.value as Period)}
            className="px-3 py-2 border border-gray-300 rounded-lg bg-white text-gray-800"
          >
            {PERIOD_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          <button
            onClick={() => loadMonitoringData(period, true)}
            disabled={isRefreshing}
            className="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 disabled:bg-gray-400 flex items-center gap-2"
          >
            <RefreshCw size={16} className={isRefreshing ? "animate-spin" : ""} />
            Refresh
          </button>
        </div>
      </div>

      {loadError && (
        <div className="rounded-lg border border-amber-300 bg-amber-50 p-4 text-amber-900">
          <p className="font-medium">Backend connection issue</p>
          <p className="text-sm mt-1">{loadError}</p>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
        <div className={`rounded-lg border p-4 ${summaryCardClass("blue")}`}>
          <div className="flex items-center gap-2 text-gray-700">
            <WalletCards size={18} />
            <span className="text-sm font-medium">Total Payments</span>
          </div>
          <p className="text-2xl font-bold text-gray-900 mt-2">{metrics?.total_payments ?? 0}</p>
        </div>

        <div className={`rounded-lg border p-4 ${summaryCardClass("green")}`}>
          <div className="flex items-center gap-2 text-gray-700">
            <Activity size={18} />
            <span className="text-sm font-medium">Success Rate</span>
          </div>
          <p className="text-2xl font-bold text-gray-900 mt-2">
            {(metrics?.success_rate ?? 0).toFixed(2)}%
          </p>
        </div>

        <div className={`rounded-lg border p-4 ${summaryCardClass("purple")}`}>
          <div className="flex items-center gap-2 text-gray-700">
            <WalletCards size={18} />
            <span className="text-sm font-medium">Revenue</span>
          </div>
          <p className="text-2xl font-bold text-gray-900 mt-2">
            {metrics?.total_revenue_formatted ?? "0 EGP"}
          </p>
        </div>

        <div className={`rounded-lg border p-4 ${summaryCardClass("orange")}`}>
          <div className="flex items-center gap-2 text-gray-700">
            <Clock3 size={18} />
            <span className="text-sm font-medium">Avg Processing</span>
          </div>
          <p className="text-2xl font-bold text-gray-900 mt-2">
            {(metrics?.avg_processing_time_ms ?? 0).toFixed(0)} ms
          </p>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-900">Webhook Retry Queue</h2>
          <span className="text-sm text-gray-500">Total records: {totalRetryQueue}</span>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
          <div className="rounded-lg border border-yellow-200 bg-yellow-50 p-3">
            <p className="text-xs text-yellow-700">Pending</p>
            <p className="text-xl font-semibold text-yellow-900">{retryStats.pending}</p>
          </div>
          <div className="rounded-lg border border-blue-200 bg-blue-50 p-3">
            <p className="text-xs text-blue-700">Processing</p>
            <p className="text-xl font-semibold text-blue-900">{retryStats.processing}</p>
          </div>
          <div className="rounded-lg border border-green-200 bg-green-50 p-3">
            <p className="text-xs text-green-700">Success</p>
            <p className="text-xl font-semibold text-green-900">{retryStats.success}</p>
          </div>
          <div className="rounded-lg border border-red-200 bg-red-50 p-3">
            <p className="text-xs text-red-700">Failed</p>
            <p className="text-xl font-semibold text-red-900">{retryStats.failed}</p>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center gap-2 mb-4">
          <Webhook size={18} className="text-gray-700" />
          <h2 className="text-lg font-semibold text-gray-900">Latest Webhook Events</h2>
        </div>

        {!health?.events?.length ? (
          <div className="rounded-lg border border-dashed border-gray-300 p-8 text-center text-gray-500">
            No webhook events recorded yet for this period.
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200 text-left text-gray-600">
                  <th className="py-2 pr-4">Time</th>
                  <th className="py-2 pr-4">Order ID</th>
                  <th className="py-2 pr-4">Status</th>
                  <th className="py-2 pr-4">HMAC</th>
                  <th className="py-2 pr-4">Processing (ms)</th>
                  <th className="py-2 pr-4">Source</th>
                  <th className="py-2">Error</th>
                </tr>
              </thead>
              <tbody>
                {health.events.map((event, idx) => (
                  <tr key={`${event.timestamp}-${event.paymob_order_id}-${idx}`} className="border-b border-gray-100">
                    <td className="py-2 pr-4 text-gray-700">{formatDateTime(event.timestamp)}</td>
                    <td className="py-2 pr-4 text-gray-900 font-medium">{event.paymob_order_id || "-"}</td>
                    <td className="py-2 pr-4">
                      <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${statusBadgeClass(event.status)}`}>
                        {event.status}
                      </span>
                    </td>
                    <td className="py-2 pr-4">
                      <span
                        className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${
                          event.hmac_valid ? "bg-green-100 text-green-700" : "bg-red-100 text-red-700"
                        }`}
                      >
                        {event.hmac_valid ? "valid" : "invalid"}
                      </span>
                    </td>
                    <td className="py-2 pr-4 text-gray-700">{event.processing_time_ms}</td>
                    <td className="py-2 pr-4 text-gray-700">{event.source}</td>
                    <td className="py-2 text-gray-600 max-w-[320px] truncate" title={event.error || ""}>
                      {event.error || "-"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
