'use client';

// Stub hooks for component-manager - these will be implemented later
// For now, returning default/empty values to allow build to complete

interface Theme {
  sections?: any[];
  [key: string]: any;
}

export function useTheme(restaurantId: number): { theme: Theme | null; isLoading: boolean; error: null } {
  return {
    theme: { sections: [] },
    isLoading: false,
    error: null,
  };
}

interface MutationHook {
  mutateAsync: (data: any) => Promise<any>;
  isPending: boolean;
  isError: boolean;
  error: null;
}

export function useSectionVisibility(restaurantId: number): MutationHook {
  return {
    mutateAsync: async ({ sectionId, isVisible }: { sectionId: number; isVisible: boolean }) => {
      // Stub implementation
      return { success: true };
    },
    isPending: false,
    isError: false,
    error: null,
  };
}

export function useSectionContent(restaurantId: number): MutationHook {
  return {
    mutateAsync: async (data: any) => {
      // Stub implementation
      return { success: true };
    },
    isPending: false,
    isError: false,
    error: null,
  };
}

export function useReorderSections(restaurantId: number): MutationHook {
  return {
    mutateAsync: async (data: any) => {
      // Stub implementation
      return { success: true };
    },
    isPending: false,
    isError: false,
    error: null,
  };
}
