"use client";

import React, { useState } from "react";
import { useAuthStore } from "@/stores/authStore";
import { superAdminAPI } from "@/lib/api";
import { toast } from "sonner";
import { Save, Plus, Trash2, AlertCircle } from "lucide-react";

export default function SettingsPage() {
  const superAdmin = useAuthStore((state) => state.superAdmin);
  const [ipWhitelist, setIpWhitelist] = useState<string[]>(
    superAdmin?.ip_whitelist || []
  );
  const [newIp, setNewIp] = useState("");
  const [isSaving, setIsSaving] = useState(false);

  const handleAddIP = () => {
    if (!newIp.trim()) {
      toast.error("Please enter an IP address");
      return;
    }

    if (ipWhitelist.includes(newIp.trim())) {
      toast.error("This IP is already in the whitelist");
      return;
    }

    setIpWhitelist([...ipWhitelist, newIp.trim()]);
    setNewIp("");
    toast.success("IP added to whitelist");
  };

  const handleRemoveIP = (ip: string) => {
    setIpWhitelist(ipWhitelist.filter((item) => item !== ip));
    toast.success("IP removed from whitelist");
  };

  const handleSave = async () => {
    try {
      setIsSaving(true);
      await superAdminAPI.updateIPWhitelist(ipWhitelist);
      toast.success("IP whitelist updated successfully");
    } catch (error) {
      toast.error("Failed to update IP whitelist");
      console.error(error);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Settings</h1>
        <p className="text-gray-600 mt-2">Manage your super admin account settings</p>
      </div>

      {/* Account Information */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">
          Account Information
        </h2>
        <div className="space-y-4">
          <div>
            <label className="text-sm font-medium text-gray-600">Email</label>
            <p className="mt-1 px-4 py-2 bg-gray-50 rounded-lg text-gray-900">
              {superAdmin?.email}
            </p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-600">Account Status</label>
            <div className="mt-1 flex items-center gap-2">
              <span
                className={`px-3 py-1 rounded-full text-sm font-medium ${
                  superAdmin?.is_active
                    ? "bg-green-100 text-green-700"
                    : "bg-red-100 text-red-700"
                }`}
              >
                {superAdmin?.is_active ? "Active" : "Inactive"}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* IP Whitelist */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">
          IP Whitelist
        </h2>
        <p className="text-sm text-gray-600 mb-4">
          Add IP addresses that are allowed to access this super admin account. Leave empty to allow all IPs.
        </p>

        {/* Info Box */}
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6 flex gap-3">
          <AlertCircle className="text-blue-600 flex-shrink-0 mt-0.5" size={20} />
          <div className="text-sm text-blue-900">
            <p className="font-medium mb-1">IP Whitelist Information</p>
            <p>Only requests from whitelisted IP addresses will be allowed. If the list is empty, all IP addresses are permitted.</p>
          </div>
        </div>

        {/* Add IP Form */}
        <div className="mb-6 p-4 bg-gray-50 rounded-lg">
          <div className="flex gap-2">
            <input
              type="text"
              value={newIp}
              onChange={(e) => setNewIp(e.target.value)}
              onKeyPress={(e) => {
                if (e.key === "Enter") {
                  handleAddIP();
                }
              }}
              placeholder="e.g., 192.168.1.1 or 10.0.0.0/8"
              className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
            />
            <button
              onClick={handleAddIP}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2"
            >
              <Plus size={18} />
              Add IP
            </button>
          </div>
        </div>

        {/* IP List */}
        {ipWhitelist.length === 0 ? (
          <div className="text-center py-8 bg-gray-50 rounded-lg">
            <p className="text-gray-600">No IP addresses in whitelist</p>
            <p className="text-sm text-gray-500 mt-1">
              All IP addresses are currently allowed
            </p>
          </div>
        ) : (
          <div className="space-y-2 mb-6">
            {ipWhitelist.map((ip) => (
              <div
                key={ip}
                className="flex items-center justify-between p-3 bg-gray-50 border border-gray-200 rounded-lg"
              >
                <span className="font-mono text-gray-900">{ip}</span>
                <button
                  onClick={() => handleRemoveIP(ip)}
                  className="p-2 hover:bg-red-50 text-red-600 rounded-lg transition-colors"
                  title="Remove IP"
                >
                  <Trash2 size={18} />
                </button>
              </div>
            ))}
          </div>
        )}

        {/* Save Button */}
        <div className="flex justify-end">
          <button
            onClick={handleSave}
            disabled={isSaving}
            className="px-6 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white rounded-lg transition-colors flex items-center gap-2"
          >
            <Save size={18} />
            {isSaving ? "Saving..." : "Save Changes"}
          </button>
        </div>
      </div>

      {/* Security Information */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-6">
        <h2 className="text-lg font-semibold text-blue-900 mb-3">
          🔒 Security Best Practices
        </h2>
        <ul className="space-y-2 text-sm text-blue-800">
          <li>✓ Always enable IP whitelist for production environments</li>
          <li>✓ Use static IP addresses for maximum security</li>
          <li>✓ Monitor audit logs regularly for suspicious activity</li>
          <li>✓ Change your password frequently</li>
          <li>✓ Never share your login credentials</li>
          <li>✓ Log out when finished using the super admin dashboard</li>
        </ul>
      </div>
    </div>
  );
}
