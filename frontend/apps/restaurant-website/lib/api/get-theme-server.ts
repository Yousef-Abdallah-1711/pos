import { headers } from 'next/headers';
import { debugError, debugWarn } from '../debug';

/**
 * Server-side theme data structure
 * Matches the API response format from theme endpoint
 */
export interface ServerThemeData {
  id: string;
  name: string;
  slug: string;
  colors: {
    primary: string;
    secondary: string;
    accent: string;
    background: string;
    text: string;
    border: string;
    shadow: string;
  };
  typography: {
    font_family: string;
    base_font_size: number;
    border_radius: number;
    line_height: number;
  };
  identity: {
    site_title: string;
    logo_url: string;
    favicon_url: string;
  };
  header: {
    background_color: string;
    text_color: string;
    tagline: string;
    height: number;
    navigation_items: Array<{
      id?: string;
      label: string;
      href: string;
      order?: number;
    }>;
  };
  footer: {
    background_color: string;
    text_color: string;
    company_description: string;
    address: string;
    phone: string;
    email: string;
    copyright_text: string;
    social_links: Array<{
      platform: string;
      url: string;
    }>;
    footer_sections: Array<{
      id: string;
      title: string;
      links?: Array<{ label: string; href: string }>;
    }>;
    legal_links: Array<{ label: string; href: string }>;
  };
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
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
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

    // Extract theme from nested response structure
    const theme = data.data?.data?.config || data.data?.config || data;

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
