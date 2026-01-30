"use client";

import React, { useState } from "react";
import type { User } from "@/types";
import { ImpersonationModal } from "./ImpersonationModal";
import { LogIn } from "lucide-react";

interface ImpersonationButtonProps {
  user: User;
  tenantId: number;
  variant?: "primary" | "secondary" | "ghost";
  size?: "sm" | "md" | "lg";
  className?: string;
}

export function ImpersonationButton({
  user,
  tenantId,
  variant = "primary",
  size = "md",
  className = "",
}: ImpersonationButtonProps) {
  const [isModalOpen, setIsModalOpen] = useState(false);

  const variantClasses = {
    primary:
      "bg-blue-600 hover:bg-blue-700 text-white",
    secondary:
      "bg-gray-200 hover:bg-gray-300 text-gray-900",
    ghost: "bg-transparent hover:bg-gray-100 text-gray-700 border border-gray-300",
  };

  const sizeClasses = {
    sm: "px-3 py-1 text-sm",
    md: "px-4 py-2 text-sm",
    lg: "px-6 py-3 text-base",
  };

  return (
    <>
      <button
        onClick={() => setIsModalOpen(true)}
        className={`flex items-center gap-2 rounded-lg font-medium transition-colors ${variantClasses[variant]} ${sizeClasses[size]} ${className}`}
      >
        <LogIn size={size === "sm" ? 14 : size === "md" ? 16 : 18} />
        Impersonate
      </button>

      <ImpersonationModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        user={user}
        tenantId={tenantId}
      />
    </>
  );
}
