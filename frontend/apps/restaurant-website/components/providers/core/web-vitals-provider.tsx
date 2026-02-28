'use client';

import { useWebVitals } from '@/lib/hooks/use-web-vitals';

/**
 * Web Vitals Provider Component
 * Initializes Core Web Vitals tracking
 * Must be used in a client component
 */
export function WebVitalsProvider() {
  // Initialize Web Vitals tracking
  useWebVitals();

  // This component doesn't render anything
  return null;
}
