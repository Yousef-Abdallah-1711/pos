"use client";

import Link from "next/link";
import { memo } from "react";
import { Facebook, Instagram, Twitter, Phone, Mail, MapPin } from "lucide-react";
import type { ServerThemeData } from "@/lib/api/get-theme-server";

interface FooterSSRProps {
  themeData: ServerThemeData;
  locale?: "en" | "ar";
}

/**
 * Footer component for SSR pages
 * Receives theme data as props (no client-side fetching)
 * Works with home page and e-commerce pages consistently
 */
function FooterSSRComponent({
  themeData,
  locale = "en",
}: FooterSSRProps) {
  const safeLocale = (locale || "en").toLowerCase() as "en" | "ar";
  const isRTL = safeLocale === "ar";

  // Use theme values directly from props with fallbacks
  const footerBgColor = themeData?.footer?.background_color || "#000000";
  const footerTextColor = themeData?.footer?.text_color || "#ffffff";
  const restaurantName = themeData?.identity?.site_title || "Restaurant";
  const companyDescription =
    themeData?.footer?.company_description ||
    (safeLocale === "en"
      ? "Delivering delicious food right to your doorstep with quality and care."
      : "توصيل الطعام اللذيذ إلى باب منزلك بجودة واهتمام.");
  const address = themeData?.footer?.address;
  const phone = themeData?.footer?.phone;
  const email = themeData?.footer?.email;
  const copyrightText =
    themeData?.footer?.copyright_text ||
    (safeLocale === "en"
      ? `© ${new Date().getFullYear()} ${restaurantName}. All rights reserved.`
      : `© ${new Date().getFullYear()} ${restaurantName}. جميع الحقوق محفوظة.`);

  const sections = {
    en: {
      about: "About Us",
      menu: "Menu",
      orders: "Orders",
      contact: "Contact",
      followUs: "Follow Us",
      company: "Company",
      legal: "Legal",
      hours: "Opening Hours",
      address: "Address",
      phone: "Phone",
      email: "Email",
      copyright: "© 2024 {name}. All rights reserved.",
      links: {
        about: "About",
        faq: "FAQ",
        blog: "Blog",
        careers: "Careers",
      },
      legal_links: {
        terms: "Terms of Service",
        privacy: "Privacy Policy",
        cookies: "Cookie Policy",
      },
    },
    ar: {
      about: "معلومات عنا",
      menu: "القائمة",
      orders: "الطلبات",
      contact: "الاتصال",
      followUs: "تابعنا",
      company: "الشركة",
      legal: "القانوني",
      hours: "ساعات العمل",
      address: "العنوان",
      phone: "الهاتف",
      email: "البريد الإلكتروني",
      copyright: "© 2024 {name}. جميع الحقوق محفوظة.",
      links: {
        about: "عن الشركة",
        faq: "الأسئلة الشائعة",
        blog: "المدونة",
        careers: "الوظائف",
      },
      legal_links: {
        terms: "شروط الخدمة",
        privacy: "سياسة الخصوصية",
        cookies: "سياسة ملفات تعريف الارتباط",
      },
    },
  };

  const t = sections[safeLocale as keyof typeof sections] || sections.en;

  return (
    <footer
      className="pt-16 pb-8"
      style={{
        backgroundColor: footerBgColor,
        color: footerTextColor,
      }}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Footer Grid */}
        <div
          className={`grid grid-cols-1 md:grid-cols-4 gap-8 mb-12 ${
            isRTL ? "text-right" : "text-left"
          }`}
        >
          {/* About Section */}
          <div>
            <h3
              className="text-lg font-semibold mb-4"
              style={{ color: footerTextColor }}
            >
              {restaurantName}
            </h3>
            <p
              className="text-sm mb-4"
              style={{ color: footerTextColor, opacity: 0.8 }}
            >
              {companyDescription}
            </p>
            <div className="flex gap-4 justify-start">
              {themeData?.footer?.social_links &&
                themeData?.footer?.social_links?.map((link) => {
                  const IconComponent =
                    link.platform === "facebook"
                      ? Facebook
                      : link.platform === "instagram"
                        ? Instagram
                        : link.platform === "twitter"
                          ? Twitter
                          : null;

                  if (!IconComponent) return null;

                  return (
                    <Link
                      key={link.platform}
                      href={link.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="transition-opacity hover:opacity-80"
                      style={{ color: footerTextColor, opacity: 0.8 }}
                    >
                      <IconComponent className="h-5 w-5" />
                    </Link>
                  );
                })}
            </div>
          </div>

          {/* Quick Links */}
          <div>
            <h4
              className="text-lg font-semibold mb-4"
              style={{ color: footerTextColor }}
            >
              {t.company}
            </h4>
            <ul className="space-y-2">
              {Object.entries(t.links).map(([key, label]) => (
                <li key={key}>
                  <Link
                    href="#"
                    className="text-sm transition-opacity hover:opacity-100"
                    style={{ color: footerTextColor, opacity: 0.8 }}
                  >
                    {label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Legal Links */}
          <div>
            <h4
              className="text-lg font-semibold mb-4"
              style={{ color: footerTextColor }}
            >
              {t.legal}
            </h4>
            <ul className="space-y-2">
              {Object.entries(t.legal_links).map(([key, label]) => (
                <li key={key}>
                  <Link
                    href="#"
                    className="text-sm transition-opacity hover:opacity-100"
                    style={{ color: footerTextColor, opacity: 0.8 }}
                  >
                    {label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Contact Info */}
          <div>
            <h4
              className="text-lg font-semibold mb-4"
              style={{ color: footerTextColor }}
            >
              {t.contact}
            </h4>
            <ul className="space-y-3">
              {address && (
                <li className="flex items-start gap-2">
                  <MapPin className="h-5 w-5 mt-0.5 flex-shrink-0" />
                  <span
                    className="text-sm"
                    style={{ color: footerTextColor, opacity: 0.8 }}
                  >
                    {address}
                  </span>
                </li>
              )}
              {phone && (
                <li className="flex items-start gap-2">
                  <Phone className="h-5 w-5 mt-0.5 flex-shrink-0" />
                  <a
                    href={`tel:${phone.replace(/\D/g, "")}`}
                    className="text-sm transition-opacity hover:opacity-80"
                    style={{ color: footerTextColor, opacity: 0.8 }}
                  >
                    {phone}
                  </a>
                </li>
              )}
              {email && (
                <li className="flex items-start gap-2">
                  <Mail className="h-5 w-5 mt-0.5 flex-shrink-0" />
                  <a
                    href={`mailto:${email}`}
                    className="text-sm transition-opacity hover:opacity-80"
                    style={{ color: footerTextColor, opacity: 0.8 }}
                  >
                    {email}
                  </a>
                </li>
              )}
            </ul>
          </div>
        </div>

        {/* Divider */}
        <div
          className="mb-8"
          style={{
            borderTop: `1px solid ${footerTextColor}`,
            opacity: 0.3,
          }}
        ></div>

        {/* Copyright */}
        <div
          className={`text-center text-sm ${isRTL ? "text-right" : "text-left"}`}
          style={{ color: footerTextColor, opacity: 0.7 }}
        >
          <p>{copyrightText}</p>
        </div>
      </div>
    </footer>
  );
}

export const FooterSSR = memo(FooterSSRComponent);
