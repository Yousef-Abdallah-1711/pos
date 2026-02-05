"use client";

import Link from "next/link";
import { useState, memo } from "react";
import { Menu, X, ShoppingCart } from "lucide-react";
import { DarkModeToggle } from "./dark-mode-toggle";
import type { ServerThemeData } from "@/lib/api/get-theme-server";

interface HeaderSSRProps {
  themeData: ServerThemeData;
  cartItemsCount?: number;
  locale?: "en" | "ar";
}

/**
 * Header component for SSR pages
 * Receives theme data as props (no client-side fetching)
 * Works with home page and e-commerce pages consistently
 */
function HeaderSSRComponent({
  themeData,
  cartItemsCount = 0,
  locale = "en",
}: HeaderSSRProps) {
  const safeLocale = (locale || "en").toLowerCase() as "en" | "ar";
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const isRTL = safeLocale === "ar";

  // Use theme values directly from props with fallbacks
  const restaurantName = themeData?.identity?.site_title || "Restaurant";
  const logoUrl = themeData?.identity?.logo_url || "";
  const headerBgColor = themeData?.header?.background_color || "#ffffff";
  const headerTextColor = themeData?.header?.text_color || "#000000";
  const headerTagline =
    themeData?.header?.tagline ||
    (safeLocale === "en" ? "Order Food Online" : "اطلب الطعام أونلاين");

  // Use navigation items from theme, or default locale-aware links
  const themeLinkItems = themeData?.header?.navigation_items ?? [];
  const navLinks = themeLinkItems.length > 0
    ? themeLinkItems.map((item) => ({
        href: item.href.startsWith("/") ? `/${safeLocale}${item.href}` : item.href,
        label: item.label,
      }))
    : [
        { href: `/${safeLocale}/menu`, label: safeLocale === "en" ? "Menu" : "القائمة" },
        {
          href: `/${safeLocale}/orders`,
          label: safeLocale === "en" ? "Orders" : "الطلبات",
        },
        {
          href: `/${safeLocale}/settings`,
          label: safeLocale === "en" ? "Settings" : "الإعدادات",
        },
      ];

  return (
    <header
      className="sticky top-0 z-50 border-b"
      style={{
        backgroundColor: headerBgColor,
        color: headerTextColor,
        borderBottomColor: `${headerTextColor}20`,
      }}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo + Restaurant Name */}
          <Link
            href={`/${safeLocale}`}
            className="flex items-center gap-2 hover:opacity-80 transition-opacity"
          >
            {logoUrl ? (
              <img
                src={logoUrl}
                alt={restaurantName}
                className="h-10 w-10 rounded-lg object-contain"
              />
            ) : (
              <div
                className="w-10 h-10 rounded-lg bg-gradient-to-br flex items-center justify-center text-white font-bold text-lg"
                style={{
                  backgroundImage: `linear-gradient(135deg, ${headerBgColor} 0%, ${headerBgColor}dd 100%)`,
                }}
              >
                🍽️
              </div>
            )}
            <div className="hidden sm:block">
              <h1
                className="font-bold text-lg truncate"
                style={{ color: headerTextColor }}
              >
                {restaurantName}
              </h1>
              <p
                className="text-xs truncate"
                style={{ color: headerTextColor, opacity: 0.8 }}
              >
                {headerTagline}
              </p>
            </div>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center gap-6">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className="text-sm font-medium transition-opacity hover:opacity-80"
                style={{ color: headerTextColor }}
              >
                {link.label}
              </Link>
            ))}
          </nav>

          {/* Right Actions */}
          <div className="flex items-center gap-2">
            <DarkModeToggle />

            {/* Cart Icon */}
            <Link
              href={`/${safeLocale}/cart`}
              className="relative p-2 rounded-lg transition-all hover:opacity-80"
              style={{
                backgroundColor: `${headerTextColor}10`,
              }}
            >
              <ShoppingCart
                className="h-5 w-5"
                style={{ color: headerTextColor }}
              />
              {cartItemsCount > 0 && (
                <span
                  className="absolute top-0 right-0 bg-red-500 text-white text-xs font-bold rounded-full h-5 w-5 flex items-center justify-center"
                  style={{
                    transform: "translate(50%, -50%)",
                  }}
                >
                  {cartItemsCount > 99 ? "99+" : cartItemsCount}
                </span>
              )}
            </Link>

            {/* Mobile Menu Button */}
            <button
              className="md:hidden p-2 rounded-lg transition-all"
              style={{
                backgroundColor: `${headerTextColor}10`,
              }}
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              {mobileMenuOpen ? (
                <X
                  className="h-5 w-5"
                  style={{ color: headerTextColor }}
                />
              ) : (
                <Menu
                  className="h-5 w-5"
                  style={{ color: headerTextColor }}
                />
              )}
            </button>
          </div>
        </div>

        {/* Mobile Navigation */}
        {mobileMenuOpen && (
          <nav
            className="md:hidden border-t py-4 space-y-3"
            style={{
              borderTopColor: `${headerTextColor}20`,
            }}
          >
            {navLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className="block text-sm font-medium transition-opacity hover:opacity-80 py-2"
                style={{ color: headerTextColor }}
                onClick={() => setMobileMenuOpen(false)}
              >
                {link.label}
              </Link>
            ))}
          </nav>
        )}
      </div>
    </header>
  );
}

export const HeaderSSR = memo(HeaderSSRComponent);
