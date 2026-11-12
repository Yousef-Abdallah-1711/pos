export type ThemeLocale = "en" | "ar";

export type HeaderVariant = "glass" | "gradient" | "solid" | "floating" | "minimal";
export type FooterVariant = "expanded" | "compact" | "cards" | "social-heavy" | "minimal";

export type SectionType =
  | "hero"
  | "featured_products"
  | "categories"
  | "why_choose_us"
  | "gallery"
  | "testimonials"
  | "cta"
  | "contact"
  | "custom";

export interface ThemeNavItem {
  id?: string;
  label: string | Record<ThemeLocale, string>;
  href: string;
  order?: number;
}

export interface ThemeSocialLink {
  id?: string;
  platform: string;
  url: string;
}

export interface ThemeLegalLink {
  label: string | Record<ThemeLocale, string>;
  href: string;
}

export interface ThemeFooterSection {
  id: string;
  title: string | Record<ThemeLocale, string>;
  links: ThemeLegalLink[];
}

export interface ThemeHeaderConfig {
  headerVariant?: HeaderVariant;
  logoUrl?: string;
  logoText?: string | Record<ThemeLocale, string>;
  showLogo?: boolean;
  stickyHeader?: boolean;
  showShadow?: boolean;
  navigationItems?: ThemeNavItem[];
  backgroundColor?: string;
  textColor?: string;
  height?: number;
}

export interface ThemeFooterConfig {
  footerVariant?: FooterVariant;
  companyName?: string | Record<ThemeLocale, string>;
  companyDescription?: string | Record<ThemeLocale, string>;
  address?: string | Record<ThemeLocale, string>;
  phone?: string;
  email?: string;
  socialLinks?: ThemeSocialLink[];
  footerSections?: ThemeFooterSection[];
  legalLinks?: ThemeLegalLink[];
  backgroundColor?: string;
  textColor?: string;
}

export interface ThemeSectionConfig {
  id: string;
  type: SectionType;
  enabled: boolean;
  order: number;
  title?: string | Record<ThemeLocale, string>;
  settings?: Record<string, unknown>;
  componentOverrides?: Record<string, unknown>;
}

export interface ThemeConfig {
  id?: number | string;
  name: string;
  slug: string;
  templateId?: string;
  headerVariant?: HeaderVariant;
  footerVariant?: FooterVariant;
  colors: Record<string, string>;
  typography: Record<string, string | number>;
  identity?: Record<string, unknown>;
  header: ThemeHeaderConfig;
  footer: ThemeFooterConfig;
  components: ThemeSectionConfig[];
}
