'use client'

import { ReactNode, useEffect } from 'react';
import { usePreferencesStore } from '@/stores/preferencesStore';

// Auth provider with preferences initialization
const AuthProvider = ({ children }: { children: ReactNode }) => {
  useEffect(() => {
    // Initialize preferences from localStorage on app start
    const preferences = usePreferencesStore.getState().settings
    if (preferences) {
      const { setSettings } = usePreferencesStore.getState()
      // Re-apply persisted preference-driven UI state (colors and metadata)
      setSettings(preferences)
    }
  }, [])

  return <>{children}</>;
};

export default AuthProvider;
