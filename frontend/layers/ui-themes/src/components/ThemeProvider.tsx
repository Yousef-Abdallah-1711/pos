'use client'

import React from 'react'
import type { ThemeProviderProps } from '../types'

/**
 * ThemeProvider: Injects CSS variables from theme configuration
 *
 * This component creates CSS custom properties that can be used throughout
 * the application for dynamic theming. It handles both the root level and
 * inline styles for compatibility.
 *
 * CSS Variables available:
 * - --color-primary
 * - --color-secondary
 * - --color-accent
 * - --color-background
 * - --color-text
 * - --color-border
 * - --color-shadow
 * - --font-family
 * - --font-size-base
 * - --border-radius
 * - --line-height
 */
export const ThemeProvider: React.FC<ThemeProviderProps> = ({ theme, children }) => {
  return (
    <div>
      <style>{`
        :root {
          --color-primary: ${theme.colors.primary};
          --color-secondary: ${theme.colors.secondary};
          --color-accent: ${theme.colors.accent};
          --color-background: ${theme.colors.background};
          --color-text: ${theme.colors.text};
          --color-border: ${theme.colors.border};
          --color-shadow: ${theme.colors.shadow};
          --font-family: ${theme.typography.fontFamily};
          --font-size-base: ${theme.typography.baseFontSize}px;
          --border-radius: ${theme.typography.borderRadius}px;
          --line-height: ${theme.typography.lineHeight};
        }

        body {
          background-color: var(--color-background);
          color: var(--color-text);
          font-family: var(--font-family), system-ui, sans-serif;
          font-size: var(--font-size-base);
          line-height: var(--line-height);
        }

        a {
          color: var(--color-primary);
          transition: color 0.2s;
        }

        a:hover {
          color: var(--color-secondary);
        }

        button {
          font-family: var(--font-family), system-ui, sans-serif;
        }

        /* Button styles */
        .btn-primary {
          background-color: var(--color-primary);
          color: white;
          border: none;
          border-radius: var(--border-radius);
          padding: 0.75rem 1.5rem;
          font-size: var(--font-size-base);
          cursor: pointer;
          transition: background-color 0.3s, opacity 0.3s;
        }

        .btn-primary:hover {
          background-color: var(--color-secondary);
          opacity: 0.9;
        }

        .btn-secondary {
          background-color: var(--color-secondary);
          color: white;
          border: none;
          border-radius: var(--border-radius);
          padding: 0.75rem 1.5rem;
          font-size: var(--font-size-base);
          cursor: pointer;
          transition: background-color 0.3s;
        }

        .btn-secondary:hover {
          opacity: 0.9;
        }

        .btn-accent {
          background-color: var(--color-accent);
          color: white;
          border: none;
          border-radius: var(--border-radius);
          padding: 0.75rem 1.5rem;
          font-size: var(--font-size-base);
          cursor: pointer;
          transition: background-color 0.3s;
        }

        .btn-accent:hover {
          opacity: 0.9;
        }

        /* Text styles */
        .text-primary {
          color: var(--color-primary);
        }

        .text-secondary {
          color: var(--color-secondary);
        }

        .text-accent {
          color: var(--color-accent);
        }

        /* Background styles */
        .bg-primary {
          background-color: var(--color-primary);
        }

        .bg-secondary {
          background-color: var(--color-secondary);
        }

        .bg-accent {
          background-color: var(--color-accent);
        }

        /* Border styles */
        .border-primary {
          border-color: var(--color-primary);
        }

        .border-secondary {
          border-color: var(--color-secondary);
        }

        .rounded {
          border-radius: var(--border-radius);
        }
      `}</style>
      {children}
    </div>
  )
}

export default ThemeProvider
