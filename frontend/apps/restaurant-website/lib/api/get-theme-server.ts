import { headers } from 'next/headers';
import { debugError, debugWarn } from '../debug';

/**
 * Server-side theme data structure
 * Matches the API response format from theme endpoint
 */
export interface ServerThemeData {
  id?: string | number;
  name?: string;
  slug?: string;
  colors?: Record<string, any>;
  typography?: Record<string, any>;
  identity?: Record<string, any>;
  header?: Record<string, any>;
  footer?: Record<string, any>;
  components?: Array<any>;
  customCss?: string;
  [key: string]: any;
}

function parseJsonIfString(value: any): any {
  if (typeof value !== 'string') {
    return value;
  }

  try {
    return JSON.parse(value);
  } catch {
    return null;
  }
}

function isThemeConfigCandidate(value: any): value is Record<string, any> {
  if (!value || typeof value !== 'object' || Array.isArray(value)) {
    return false;
  }

  const themeKeys = [
    'colors',
    'typography',
    'identity',
    'header',
    'footer',
    'components',
    'customCss',
    'custom_css',
  ];

  return themeKeys.some((key) => key in value);
}

/**
 * Extracts the actual theme config object from inconsistent API payload shapes.
 * Handles nested wrappers like:
 * - { success, data: { data: {...themeConfig} } }
 * - { success, data: {...themeConfig} }
 * - { success, data: { theme: { config: {...themeConfig} } } }
 */
export function extractThemeConfig(payload: any): ServerThemeData | null {
  const candidates = [
    payload?.data?.data?.config,
    payload?.data?.config,
    payload?.data?.data,
    payload?.data?.theme?.config,
    payload?.data?.theme,
    payload?.theme?.config,
    payload?.theme,
    payload?.config,
    payload?.data,
    payload,
  ];

  for (const candidate of candidates) {
    const parsed = parseJsonIfString(candidate);
    if (isThemeConfigCandidate(parsed)) {
      return parsed as ServerThemeData;
    }
  }

  return null;
}

/**
 * Fetch theme data server-side for a restaurant
 * Uses server headers to get restaurant slug
 * Caches theme for 1 hour (ISR friendly)
 *
 * @param restaurantSlug - Optional restaurant slug. If not provided, reads from x-restaurant-slug header
 * @returns Theme data or null if not found
 */
export async function getThemeDataServer(
  restaurantSlug?: string
): Promise<ServerThemeData | null> {
  try {
    // Get restaurant slug from parameter or headers
    let slug = restaurantSlug;
    if (!slug) {
      const headersList = await headers();
      slug = headersList.get('x-restaurant-slug');
    }

    if (!slug) {
      debugError('Restaurant slug not found in parameter or headers');
      return null;
    }

    // Fetch theme data from API
    const apiUrl =
      process.env.NEXT_PUBLIC_API_BASE_URL ||
      process.env.NEXT_PUBLIC_API_URL ||
      'http://localhost:8080/api/v1';
    const response = await fetch(
      `${apiUrl}/public/restaurants/${slug}/theme`,
      {
        next: { revalidate: 3600 }, // Cache for 1 hour (rarely changes)
      }
    );

    if (!response.ok) {
      if (response.status === 404) {
        debugWarn(`Theme not found for restaurant: ${slug}`);
      } else {
        debugError(`Failed to fetch theme: ${response.status}`);
      }
      return null;
    }

    const data = await response.json();
    const theme = extractThemeConfig(data);

    // Ensure theme is an object before returning
    if (typeof theme !== 'object' || theme === null) {
      debugWarn('Theme data is not a valid object');
      return null;
    }

    // Return theme even if some fields are missing - components have defensive fallbacks
    return theme as ServerThemeData;
  } catch (error) {
    debugError('Error fetching theme:', error);
    return null;
  }
}

/**
 * Get default fallback theme data when API fails
 * Used for graceful degradation
 */
export function getDefaultThemeData(): ServerThemeData {
  return {
    id: 'default',
    name: 'Default Theme',
    slug: 'default',
    colors: {
      primary: '#f97316',
      secondary: '#0ea5e9',
      accent: '#fbbf24',
      background: '#ffffff',
      text: '#1f2937',
      border: '#e5e7eb',
      shadow: '#1f2937',
    },
    typography: {
      font_family: 'Inter, sans-serif',
      base_font_size: 14,
      border_radius: 8,
      line_height: 1.5,
    },
    identity: {
      site_title: 'Restaurant',
      logo_url: '',
      favicon_url: '',
    },
    header: {
      background_color: '#ffffff',
      text_color: '#1f2937',
      tagline: 'Order Food Online',
      height: 64,
      navigation_items: [
        { label: 'Menu', href: '/menu', order: 1 },
        { label: 'Orders', href: '/orders', order: 2 },
        { label: 'Settings', href: '/settings', order: 3 },
      ],
    },
    footer: {
      background_color: '#111827',
      text_color: '#f3f4f6',
      company_description: 'Delivering delicious food to your doorstep',
      address: '',
      phone: '',
      email: '',
      copyright_text: `© ${new Date().getFullYear()} Restaurant. All rights reserved.`,
      social_links: [],
      footer_sections: [],
      legal_links: [],
    },
  };
}
