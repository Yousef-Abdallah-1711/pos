'use client';

import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { AlertCircle, AlertTriangle, Loader, RefreshCw } from 'lucide-react';
import { ProductCard } from '@/components/commerce/product-card';
import { useCartStore } from '@/lib/store/cart-store';
import { useProducts } from '@/lib/hooks/use-api-queries';
import { SectionRenderer, ThemeProvider as UiThemeProvider } from '@pos-saas/ui-themes';
import { debugError } from '@/lib/debug';

interface ThemeComponent {
  id: string;
  type: string;
  title: string;
  enabled: boolean;
  displayOrder: number;
  config?: Record<string, any>;
}

interface HomepageData {
  theme: any;
  components: ThemeComponent[];
  restaurant?: { slug: string };
}

interface HomePageContentProps {
  restaurantSlug: string;
  locale: 'en' | 'ar';
  initialHomepageData?: any;
}

function extractBilingualText(value: any, locale: 'en' | 'ar' = 'en'): string {
  if (!value) return '';
  if (typeof value === 'string') return value;
  if (typeof value === 'object' && value[locale]) return value[locale];
  if (typeof value === 'object' && value.en) return value.en;
  if (typeof value === 'object' && value.ar) return value.ar;
  return String(value);
}

function mapComponentType(type: string): string {
  const typeMapping: Record<string, string> = {
    'featured-items': 'products',
    featured_items: 'products',
    'why-choose-us': 'why_us',
    why_choose_us: 'why_us',
    'why-us': 'why_us',
    'call-to-action': 'cta',
  };
  return typeMapping[type.toLowerCase()] || type;
}

class SectionErrorBoundary extends React.Component<
  { children: React.ReactNode; locale: 'en' | 'ar'; section?: string },
  { hasError: boolean; error: Error | null }
> {
  constructor(props: { children: React.ReactNode; locale: 'en' | 'ar'; section?: string }) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error) {
    debugError(`[section:${this.props.section || 'unknown'}]`, error);
  }

  reset = () => this.setState({ hasError: false, error: null });

  render() {
    if (!this.state.hasError) {
      return this.props.children;
    }

    const isArabic = this.props.locale === 'ar';
    return (
      <div className={`rounded-lg border border-red-200 bg-red-50 p-4 ${isArabic ? 'text-right' : 'text-left'}`}>
        <div className="flex items-start gap-2">
          <AlertTriangle className="mt-0.5 h-5 w-5 text-red-600" />
          <div>
            <p className="font-semibold text-red-900">
              {isArabic ? '??? ??? ????? ????? ?????' : 'Section failed to load'}
            </p>
            <p className="mt-1 text-sm text-red-800">{this.state.error?.message}</p>
            <button
              className="mt-3 inline-flex items-center gap-1 rounded bg-red-600 px-3 py-1 text-xs font-medium text-white"
              onClick={this.reset}
            >
              <RefreshCw className="h-3 w-3" />
              {isArabic ? '????? ????????' : 'Retry'}
            </button>
          </div>
        </div>
      </div>
    );
  }
}

function FeaturedProducts({ restaurantSlug, locale }: { restaurantSlug: string; locale: 'en' | 'ar' }) {
  const addToCart = useCartStore((state) => state.addItem);
  const { data: allProducts = [], isLoading, error } = useProducts(restaurantSlug);
  const products = useMemo(() => (Array.isArray(allProducts) ? allProducts.slice(0, 6) : []), [allProducts]);

  const handleAddToCart = useCallback(
    (product: any, quantity: number) => {
      for (let i = 0; i < quantity; i += 1) {
        addToCart({
          id: product.id,
          name: product.name || product.name_en || 'Product',
          price: product.price || 0,
          image: product.image || product.main_image_url || '',
          track_inventory: product.track_inventory,
          quantity_in_stock: product.quantity_in_stock,
        });
      }
    },
    [addToCart]
  );

  if (isLoading) {
    return (
      <section className="w-full px-4 py-16">
        <div className="mx-auto flex min-h-48 max-w-6xl items-center justify-center">
          <Loader className="h-8 w-8 animate-spin text-gray-400" />
        </div>
      </section>
    );
  }

  if (error && products.length === 0) {
    return (
      <section className="w-full px-4 py-16">
        <div className="mx-auto max-w-6xl rounded-lg border border-yellow-200 bg-yellow-50 p-4 text-sm text-yellow-900">
          {error instanceof Error ? error.message : 'Failed to load products'}
        </div>
      </section>
    );
  }

  if (products.length === 0) {
    return null;
  }

  return (
    <section className="w-full bg-gray-50 px-4 py-16 dark:bg-gray-900/30">
      <div className="mx-auto max-w-6xl">
        <h2 className="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
          {locale === 'ar' ? '?????? ?????' : 'Featured Products'}
        </h2>
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {products.map((product) => (
            <ProductCard
              key={product.id}
              product={product}
              locale={locale}
              onAddToCart={handleAddToCart}
            />
          ))}
        </div>
      </div>
    </section>
  );
}

export function HomePageContent({
  restaurantSlug,
  locale,
  initialHomepageData,
}: HomePageContentProps) {
  const [homepageData] = useState<HomepageData | null>(() => {
    if (!initialHomepageData) return null;

    const innerData = initialHomepageData.data?.data || initialHomepageData.data || initialHomepageData;
    if (!innerData?.theme) return null;

    const rawTheme = innerData.theme;
    const configData = rawTheme.config || {};
    const componentsFromConfig = (configData.components || []).map((comp: any) => ({
      ...comp,
      type: mapComponentType(comp.type || 'hero'),
      title: extractBilingualText(comp.title, locale),
    }));

    const theme = {
      colors: {
        primary: configData.colors?.primary || rawTheme.primary_color || '#3b82f6',
        secondary: configData.colors?.secondary || rawTheme.secondary_color || '#1e40af',
        accent: configData.colors?.accent || rawTheme.accent_color || '#0ea5e9',
        background: configData.colors?.background || rawTheme.background_color || '#ffffff',
        text: configData.colors?.text || rawTheme.text_color || '#1f2937',
        border: configData.colors?.border || rawTheme.border_color || '#e5e7eb',
        shadow: configData.colors?.shadow || rawTheme.shadow_color || '#1f2937',
      },
      typography: {
        fontFamily: configData.typography?.fontSans || rawTheme.font_family || 'Inter',
        baseFontSize: configData.typography?.baseSize || rawTheme.base_font_size || 14,
        borderRadius: configData.typography?.borderRadius || rawTheme.border_radius || 8,
        lineHeight: configData.typography?.lineHeight || rawTheme.line_height || 1.5,
      },
      identity: {
        siteTitle: configData.identity?.siteTitle || rawTheme.site_title || 'Restaurant',
      },
      customCss: configData.customCss || configData.custom_css || '',
    };

    return {
      theme,
      components: componentsFromConfig,
      restaurant: innerData.restaurant ? { slug: innerData.restaurant.slug } : undefined,
    };
  });

  useEffect(() => {
    if (!homepageData?.theme?.customCss || typeof document === 'undefined') {
      return undefined;
    }

    const styleId = 'restaurant-theme-custom-css';
    let styleTag = document.getElementById(styleId) as HTMLStyleElement | null;
    if (!styleTag) {
      styleTag = document.createElement('style');
      styleTag.id = styleId;
      document.head.appendChild(styleTag);
    }

    styleTag.textContent = homepageData.theme.customCss;
    return () => {
      const current = document.getElementById(styleId);
      current?.remove();
    };
  }, [homepageData?.theme?.customCss]);

  if (!homepageData?.theme) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="rounded-lg border border-red-200 bg-red-50 p-6">
          <div className="flex items-center gap-2 text-red-800">
            <AlertCircle className="h-5 w-5" />
            <span>{locale === 'ar' ? '???? ????? ?????? ??????' : 'Failed to load restaurant data'}</span>
          </div>
        </div>
      </div>
    );
  }

  const visibleComponents = homepageData.components
    .filter((component) => component.enabled)
    .sort((a, b) => {
      const orderA = typeof a.displayOrder === 'number' ? a.displayOrder : 0;
      const orderB = typeof b.displayOrder === 'number' ? b.displayOrder : 0;
      return orderA - orderB;
    });

  return (
    <UiThemeProvider theme={{ colors: homepageData.theme.colors, typography: homepageData.theme.typography }}>
      <div dir={locale === 'ar' ? 'rtl' : 'ltr'}>
        <main>
          {visibleComponents.map((component) => (
            <SectionErrorBoundary key={component.id} locale={locale} section={component.type}>
              <SectionRenderer component={component as any} isArabic={locale === 'ar'} />
            </SectionErrorBoundary>
          ))}
          <SectionErrorBoundary locale={locale} section="featured_products">
            <FeaturedProducts restaurantSlug={restaurantSlug} locale={locale} />
          </SectionErrorBoundary>
        </main>
      </div>
    </UiThemeProvider>
  );
}
