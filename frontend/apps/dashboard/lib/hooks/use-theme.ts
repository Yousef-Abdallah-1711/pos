import { useState, useCallback } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient from '@/lib/api/themeApi';

export interface Theme {
  id: number;
  restaurant_id: number;
  tenant_id: number;
  primary_color: string;
  secondary_color: string;
  accent_color: string;
  logo_url: string;
  font_family: string;
  sections: ThemeSection[];
  created_at: string;
  updated_at: string;
}

export interface ThemeSection {
  id: number;
  theme_id: number;
  section_type: string;
  order: number;
  is_visible: boolean;
  title: string;
  subtitle: string;
  description: string;
  background_image?: string;
  button_text: string;
  button_link: string;
  content: Record<string, string | number | string[] | undefined>;
  created_at: string;
  updated_at: string;
}

export interface UpdateThemeRequest {
  primary_color?: string;
  secondary_color?: string;
  accent_color?: string;
  logo_url?: string;
  font_family?: string;
  sections?: ThemeSection[];
}

interface ThemeResponse {
  success: boolean;
  data: Theme;
  message?: string;
}

const queryKeys = {
  theme: (restaurantId: number) => ['theme', restaurantId],
  all: () => ['themes'],
};

/**
 * Hook to fetch and manage restaurant theme
 */
export function useTheme(restaurantId: number) {
  const queryClient = useQueryClient();
  const [error, setError] = useState<string | null>(null);

  // Fetch theme
  const {
    data: themeData,
    isLoading,
    error: fetchError,
  } = useQuery({
    queryKey: queryKeys.theme(restaurantId),
    queryFn: async () => {
      try {
        const response = await apiClient.get<ThemeResponse>(
          `/admin/restaurants/${restaurantId}/theme`
        );
        return response.data.data;
      } catch (err) {
        throw new Error(err instanceof Error ? err.message : 'Failed to fetch theme');
      }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 30 * 60 * 1000, // 30 minutes
  });

  // Update theme mutation
  const updateMutation = useMutation({
    mutationFn: async (updateData: UpdateThemeRequest) => {
      const response = await apiClient.put<ThemeResponse>(
        `/admin/restaurants/${restaurantId}/theme`,
        updateData
      );
      return response.data.data;
    },
    onSuccess: (updatedTheme) => {
      // Update cache
      queryClient.setQueryData(queryKeys.theme(restaurantId), updatedTheme);
    },
    onError: (err) => {
      setError(err instanceof Error ? err.message : 'Failed to update theme');
    },
  });

  const updateTheme = useCallback(
    async (data: UpdateThemeRequest) => {
      setError(null);
      try {
        await updateMutation.mutateAsync(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to update theme');
        throw err;
      }
    },
    [updateMutation]
  );

  return {
    theme: themeData,
    isLoading,
    error: error || (fetchError ? (fetchError as Error).message : null),
    updateTheme,
    isSaving: updateMutation.isPending,
  };
}

/**
 * Hook to update section visibility
 */
export function useSectionVisibility(restaurantId: number) {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: async ({ sectionId, isVisible }: { sectionId: number; isVisible: boolean }) => {
      const response = await apiClient.put(
        `/admin/restaurants/${restaurantId}/theme/sections/${sectionId}/visibility`,
        { is_visible: isVisible }
      );
      return response.data;
    },
    onSuccess: () => {
      // Invalidate theme query to refetch
      queryClient.invalidateQueries({ queryKey: queryKeys.theme(restaurantId) });
    },
  });

  return mutation;
}

/**
 * Hook to update section content
 */
export function useSectionContent(restaurantId: number) {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: async ({
      sectionId,
      title,
      subtitle,
      description,
      content,
    }: {
      sectionId: number;
      title: string;
      subtitle: string;
      description: string;
      content: Record<string, any>;
    }) => {
      const response = await apiClient.put(
        `/admin/restaurants/${restaurantId}/theme/sections/${sectionId}/content`,
        {
          title,
          subtitle,
          description,
          content,
        }
      );
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.theme(restaurantId) });
    },
  });

  return mutation;
}

/**
 * Hook to reorder sections
 */
export function useReorderSections(restaurantId: number) {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: async (orders: Record<number, number>) => {
      const response = await apiClient.put(
        `/admin/restaurants/${restaurantId}/theme/sections/reorder`,
        { orders }
      );
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.theme(restaurantId) });
    },
  });

  return mutation;
}

/**
 * Hook to fetch component library
 */
export function useComponentLibrary() {
  const queryKeys = {
    all: () => ['components'],
  };

  return useQuery({
    queryKey: queryKeys.all(),
    queryFn: async () => {
      const response = await apiClient.get('/admin/components');
      return response.data.data;
    },
    staleTime: 60 * 60 * 1000, // 1 hour (components rarely change)
    gcTime: 24 * 60 * 60 * 1000, // 24 hours
  });
}
