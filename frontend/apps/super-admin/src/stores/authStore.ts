import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { SuperAdminUser, ImpersonationStatus } from "@/types";

interface AuthState {
  token: string | null;
  superAdmin: SuperAdminUser | null;
  impersonation: ImpersonationStatus | null;
  isLoading: boolean;
  error: string | null;

  // Actions
  setToken: (token: string | null) => void;
  setSuperAdmin: (superAdmin: SuperAdminUser | null) => void;
  setImpersonation: (impersonation: ImpersonationStatus | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  logout: () => void;
  reset: () => void;

  // Computed
  isAuthenticated: () => boolean;
  isImpersonating: () => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      superAdmin: null,
      impersonation: null,
      isLoading: false,
      error: null,

      setToken: (token) => set({ token }),
      setSuperAdmin: (superAdmin) => set({ superAdmin }),
      setImpersonation: (impersonation) => set({ impersonation }),
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),

      logout: () => {
        set({
          token: null,
          superAdmin: null,
          impersonation: null,
          error: null,
        });
      },

      reset: () => {
        set({
          token: null,
          superAdmin: null,
          impersonation: null,
          isLoading: false,
          error: null,
        });
      },

      isAuthenticated: () => {
        const { token } = get();
        return !!token;
      },

      isImpersonating: () => {
        const { impersonation } = get();
        return impersonation?.is_impersonating ?? false;
      },
    }),
    {
      name: "super-admin-auth-store",
      partialize: (state) => ({
        token: state.token,
        superAdmin: state.superAdmin,
        impersonation: state.impersonation,
      }),
    }
  )
);
