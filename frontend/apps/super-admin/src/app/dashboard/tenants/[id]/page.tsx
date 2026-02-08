"use client";

import React, { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { superAdminAPI } from "@/lib/api";
import type { Tenant, User, PaginatedResponse, SubscriptionStatus } from "@/types";
import { ArrowLeft, Users as UsersIcon, AlertCircle, LogIn, Calendar, Clock, Lock, CheckCircle, Trash2, RefreshCw } from "lucide-react";
import { ImpersonationButton } from "@/components/impersonation/ImpersonationButton";
import { ExtendTrialModal } from "@/components/tenants/ExtendTrialModal";
import { DisableTenantModal } from "@/components/tenants/DisableTenantModal";
import { EnableTenantModal } from "@/components/tenants/EnableTenantModal";
import { DeleteTenantModal } from "@/components/tenants/DeleteTenantModal";
import { TenantStatusBadge } from "@/components/tenants/TenantStatusBadge";
import { TrialDaysCounter } from "@/components/tenants/TrialDaysCounter";
import { SubscriptionHistoryTable } from "@/components/tenants/SubscriptionHistoryTable";
import { toast } from "sonner";

interface SubscriptionInfo {
  tenant_id: number;
  subscription_status: SubscriptionStatus;
  trial_start_date: string | null;
  trial_end_date: string | null;
  trial_days: number;
  is_disabled: boolean;
  disabled_at: string | null;
  disabled_reason: string | null;
  enabled_at: string | null;
  days_remaining: number;
}

interface SubscriptionHistory {
  id: number;
  action: string;
  performed_by: number | null;
  old_status: string | null;
  new_status: string | null;
  old_trial_end_date: string | null;
  new_trial_end_date: string | null;
  days_added: number | null;
  reason: string | null;
  created_at: string;
}

export default function TenantDetailPage() {
  const params = useParams();
  const router = useRouter();
  const tenantId = parseInt(params.id as string);

  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [subscriptionInfo, setSubscriptionInfo] = useState<SubscriptionInfo | null>(null);
  const [subscriptionHistory, setSubscriptionHistory] = useState<SubscriptionHistory[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isDashboardLoading, setIsDashboardLoading] = useState(false);
  const [isRefreshingSubscription, setIsRefreshingSubscription] = useState(false);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showExtendModal, setShowExtendModal] = useState(false);
  const [showDisableModal, setShowDisableModal] = useState(false);
  const [showEnableModal, setShowEnableModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const limit = 10;

  // Fetch all data
  const fetchData = async () => {
    try {
      setIsLoading(true);
      console.log("[TenantDetailPage] Fetching tenant data for ID:", tenantId);

      const [tenantData, usersData, subData] = await Promise.all([
        superAdminAPI.getTenant(tenantId),
        superAdminAPI.getTenantUsers(tenantId, page, limit),
        superAdminAPI.getTenantSubscription(tenantId).catch(() => null),
      ]);

      setTenant(tenantData);
      setUsers(usersData.data);
      setTotal(usersData.total);

      if (subData) {
        setSubscriptionInfo(subData as SubscriptionInfo);
      }
    } catch (error) {
      console.error("[TenantDetailPage] Error fetching data:", error);
      toast.error("Failed to load tenant details");
    } finally {
      setIsLoading(false);
    }
  };

  // Fetch subscription history
  const fetchSubscriptionHistory = async () => {
    try {
      const response = await superAdminAPI.getSubscriptionHistory(tenantId, 5, 0);
      if (response?.items) {
        setSubscriptionHistory(response.items as SubscriptionHistory[]);
      }
    } catch (error) {
      console.error("[TenantDetailPage] Error fetching subscription history:", error);
    }
  };

  useEffect(() => {
    fetchData();
  }, [tenantId, page]);

  useEffect(() => {
    if (subscriptionInfo) {
      fetchSubscriptionHistory();
    }
  }, [subscriptionInfo]);

  const totalPages = Math.ceil(total / limit);

  const handleEnterDashboard = async () => {
    try {
      setIsDashboardLoading(true);
      console.log("[TenantDetailPage] Generating dashboard login token for tenant:", tenantId);
      const response = await superAdminAPI.generateDashboardLoginToken(tenantId);
      window.open(response.redirect_url, "_blank");
      toast.success("Opening tenant dashboard...");
    } catch (error) {
      console.error("[TenantDetailPage] Failed to generate dashboard login token:", error);
      toast.error("Failed to generate dashboard access token");
    } finally {
      setIsDashboardLoading(false);
    }
  };

  const handleRefreshSubscription = async () => {
    setIsRefreshingSubscription(true);
    try {
      const response = await superAdminAPI.getTenantSubscription(tenantId);
      if (response) {
        setSubscriptionInfo(response as SubscriptionInfo);
        toast.success("Subscription data refreshed");
        fetchSubscriptionHistory();
      }
    } catch (error) {
      console.error("Error refreshing subscription:", error);
      toast.error("Failed to refresh subscription data");
    } finally {
      setIsRefreshingSubscription(false);
    }
  };

  const handleAfterAction = () => {
    handleRefreshSubscription();
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading tenant details...</p>
        </div>
      </div>
    );
  }

  if (!tenant) {
    return (
      <div className="text-center py-12">
        <AlertCircle className="mx-auto text-gray-400" size={48} />
        <p className="text-gray-600 mt-4">Tenant not found</p>
        <button
          onClick={() => router.back()}
          className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          Go back
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <button
          onClick={() => router.back()}
          className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
        >
          <ArrowLeft size={24} className="text-gray-600" />
        </button>
        <div className="flex-1">
          <h1 className="text-3xl font-bold text-gray-900">{tenant.name}</h1>
          <p className="text-gray-600 mt-1">Tenant Management</p>
        </div>
      </div>

      {/* Tenant Info Card */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-900">Tenant Information</h2>
          <div className="flex items-center gap-2">
            {(!users || users.length === 0) && (
              <span className="text-sm text-amber-600 flex items-center gap-1">
                <AlertCircle size={14} />
                No admin users
              </span>
            )}
            <button
              onClick={handleEnterDashboard}
              disabled={isDashboardLoading || !users || users.length === 0}
              title={!users || users.length === 0 ? "Cannot enter dashboard: No admin users" : "Enter tenant dashboard"}
              className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <LogIn size={18} />
              {isDashboardLoading ? "Opening..." : "Enter Dashboard"}
            </button>
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="text-sm font-medium text-gray-600">Name</label>
            <p className="mt-1 text-lg text-gray-900">{tenant.name}</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-600">Slug</label>
            <p className="mt-1 text-lg text-gray-900">{tenant.slug}</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-600">Created At</label>
            <p className="mt-1 text-lg text-gray-900">
              {new Date(tenant.created_at).toLocaleDateString()}
            </p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-600">Updated At</label>
            <p className="mt-1 text-lg text-gray-900">
              {new Date(tenant.updated_at).toLocaleDateString()}
            </p>
          </div>
        </div>
      </div>

      {/* Subscription Management Card */}
      {subscriptionInfo && (
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-3">
              <Calendar size={20} className="text-blue-600" />
              <h2 className="text-lg font-semibold text-gray-900">Subscription & Trial Management</h2>
            </div>
            <button
              onClick={handleRefreshSubscription}
              disabled={isRefreshingSubscription}
              className="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors disabled:opacity-50"
              title="Refresh subscription data"
            >
              <RefreshCw size={18} className={isRefreshingSubscription ? "animate-spin" : ""} />
            </button>
          </div>

          {/* Subscription Status */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
            {/* Status Badge */}
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="text-sm font-medium text-gray-600 mb-2">Status</p>
              <TenantStatusBadge status={subscriptionInfo.subscription_status} />
            </div>

            {/* Trial Days Counter */}
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="text-sm font-medium text-gray-600 mb-2">Trial Days Remaining</p>
              <TrialDaysCounter trialEndDate={subscriptionInfo.trial_end_date} trialDays={subscriptionInfo.trial_days} showProgress />
            </div>

            {/* Trial End Date */}
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="text-sm font-medium text-gray-600 flex items-center gap-2 mb-2">
                <Clock size={14} />
                Trial End Date
              </p>
              <p className="text-lg font-semibold text-gray-900">
                {subscriptionInfo.trial_end_date
                  ? new Date(subscriptionInfo.trial_end_date).toLocaleDateString()
                  : "N/A"}
              </p>
            </div>

            {/* Disabled Status */}
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="text-sm font-medium text-gray-600 flex items-center gap-2 mb-2">
                <Lock size={14} />
                Account Status
              </p>
              {subscriptionInfo.is_disabled ? (
                <span className="inline-block px-3 py-1 bg-red-100 text-red-700 rounded-full text-sm font-medium">
                  Disabled
                </span>
              ) : (
                <span className="inline-block px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium flex items-center gap-1">
                  <CheckCircle size={14} />
                  Active
                </span>
              )}
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex flex-wrap gap-3 pt-4 border-t border-gray-200">
            <button
              onClick={() => setShowExtendModal(true)}
              disabled={subscriptionInfo.subscription_status === "expired"}
              className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              title={subscriptionInfo.subscription_status === "expired" ? "Cannot extend expired trial" : "Extend trial period"}
            >
              <RefreshCw size={16} />
              Extend Trial
            </button>

            {!subscriptionInfo.is_disabled ? (
              <button
                onClick={() => setShowDisableModal(true)}
                className="flex items-center gap-2 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
                title="Disable this tenant's access"
              >
                <Lock size={16} />
                Disable Tenant
              </button>
            ) : (
              <button
                onClick={() => setShowEnableModal(true)}
                className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
                title="Re-enable this tenant's access"
              >
                <CheckCircle size={16} />
                Enable Tenant
              </button>
            )}

            <button
              onClick={() => setShowDeleteModal(true)}
              className="flex items-center gap-2 px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
              title="Schedule this tenant for deletion"
            >
              <Trash2 size={16} />
              Delete Tenant
            </button>
          </div>

          {/* Disabled Reason */}
          {subscriptionInfo.is_disabled && subscriptionInfo.disabled_reason && (
            <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-sm font-medium text-red-900">Disable Reason:</p>
              <p className="text-sm text-red-700 mt-1">{subscriptionInfo.disabled_reason}</p>
            </div>
          )}
        </div>
      )}

      {/* Subscription History */}
      {subscriptionHistory.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Recent Subscription Changes</h2>
          <SubscriptionHistoryTable tenantId={tenantId} />
        </div>
      )}

      {/* Users Section */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
          <UsersIcon size={20} />
          Users in this Tenant
        </h2>

        {!users || users.length === 0 ? (
          <div className="text-center py-8">
            <UsersIcon className="mx-auto text-gray-400" size={48} />
            <p className="text-gray-600 mt-4">No users in this tenant</p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="bg-gray-50 border-b border-gray-200">
                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">Name</th>
                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">Email</th>
                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">Role</th>
                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">Created</th>
                    <th className="px-4 py-3 text-right text-sm font-semibold text-gray-900">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((user) => (
                    <tr key={user.id} className="border-b border-gray-200 hover:bg-gray-50 transition-colors">
                      <td className="px-4 py-3 text-sm text-gray-900 font-medium">
                        {user.first_name} {user.last_name}
                      </td>
                      <td className="px-4 py-3 text-sm text-gray-600">{user.email}</td>
                      <td className="px-4 py-3 text-sm">
                        <span className="px-2 py-1 bg-blue-100 text-blue-700 rounded text-xs font-medium">
                          {user.role}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-sm text-gray-600">
                        {new Date(user.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-3 text-right">
                        <ImpersonationButton user={user} tenantId={tenantId} size="sm" variant="secondary" />
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {totalPages > 1 && (
              <div className="mt-4 flex items-center justify-between">
                <p className="text-sm text-gray-600">
                  Showing {(page - 1) * limit + 1} to {Math.min(page * limit, total)} of {total} users
                </p>
                <div className="flex gap-2">
                  <button
                    onClick={() => setPage(Math.max(1, page - 1))}
                    disabled={page === 1}
                    className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50"
                  >
                    Previous
                  </button>
                  <span className="px-3 py-1 text-sm text-gray-600">
                    Page {page} of {totalPages}
                  </span>
                  <button
                    onClick={() => setPage(Math.min(totalPages, page + 1))}
                    disabled={page === totalPages}
                    className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50"
                  >
                    Next
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>

      {/* Modals */}
      <ExtendTrialModal
        tenantId={tenantId}
        tenantName={tenant.name}
        currentTrialEndDate={subscriptionInfo?.trial_end_date ?? null}
        isOpen={showExtendModal}
        onClose={() => setShowExtendModal(false)}
        onSuccess={handleAfterAction}
      />
      <DisableTenantModal
        tenantId={tenantId}
        tenantName={tenant.name}
        isOpen={showDisableModal}
        onClose={() => setShowDisableModal(false)}
        onSuccess={handleAfterAction}
      />
      <EnableTenantModal
        tenantId={tenantId}
        tenantName={tenant.name}
        disabledReason={subscriptionInfo?.disabled_reason ?? null}
        disabledAt={subscriptionInfo?.disabled_at ?? null}
        isOpen={showEnableModal}
        onClose={() => setShowEnableModal(false)}
        onSuccess={handleAfterAction}
      />
      <DeleteTenantModal
        tenantId={tenantId}
        tenantName={tenant.name}
        isOpen={showDeleteModal}
        onClose={() => setShowDeleteModal(false)}
        onSuccess={() => router.push("/dashboard/tenants")}
      />
    </div>
  );
}
