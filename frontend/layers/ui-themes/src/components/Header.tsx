/**
 * Data-Driven Header Component
 * Renders header based on theme configuration with no hardcoded content.
 * Header config includes logo, text, colors, navigation items, layout options.
 */

import React, { useState } from 'react'
import { Menu, X } from 'lucide-react'
import { HeaderConfig } from '../types'

export interface HeaderProps {
  config: HeaderConfig
  isArabic?: boolean
  onNavigationClick?: (href: string) => void
}

export const Header: React.FC<HeaderProps> = ({
  config,
  isArabic = false,
  onNavigationClick,
}) => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const textColor = config.textColor || '#111827'
  const backgroundColor = config.backgroundColor || '#ffffff'
  const headerHeight = config.height || 64
  const headerPadding = config.padding || 16
  const showLogo = config.showLogo !== false
  const showShadow = config.showShadow !== false
  const hasNavigation = Boolean(config.navigationItems && config.navigationItems.length > 0)
  const sortedNavigationItems = [...(config.navigationItems || [])].sort(
    (a, b) => (a.order || 0) - (b.order || 0)
  )
  const navPositionClass =
    config.navPosition === 'left'
      ? 'theme-header-nav-left'
      : config.navPosition === 'center'
        ? 'theme-header-nav-center'
        : 'theme-header-nav-right'
  const navAlignmentClass =
    config.navAlignment === 'vertical'
      ? 'theme-header-nav-vertical'
      : 'theme-header-nav-horizontal'
  const navLeftGap = config.navPosition === 'left' ? 32 : 0
  const navMarginLeft = isArabic ? 0 : navLeftGap
  const navMarginRight = isArabic ? navLeftGap : 0

  const logoContent = (
    <div className="theme-header-logo">
      {showLogo && (config.logoUrl || config.logoText) && (
        <>
          {config.logoUrl && (
            <img
              src={config.logoUrl}
              alt="Logo"
              className="theme-header-logo-image"
            />
          )}
          {config.logoText && <span className="theme-header-logo-text">{config.logoText}</span>}
        </>
      )}
    </div>
  )

  return (
    <header className="theme-header">
      <div className="theme-header-inner">
        {/* Logo Section */}
        {config.homeHref ? (
          <a href={config.homeHref} className="theme-header-logo-link">
            {logoContent}
          </a>
        ) : (
          logoContent
        )}

        {/* Desktop Navigation */}
        {hasNavigation && (
          <nav
            className={`theme-header-nav ${navPositionClass} ${navAlignmentClass} hidden md:flex`}
          >
            {sortedNavigationItems.map((item) => (
              <a
                key={item.id || item.label}
                href={item.href}
                onClick={(e) => {
                  if (onNavigationClick) {
                    e.preventDefault()
                    onNavigationClick(item.href)
                  }
                }}
                className="theme-header-link"
              >
                {item.label}
              </a>
            ))}
          </nav>
        )}

        {/* Mobile Menu Button */}
        {!config.hideNavOnMobile && hasNavigation && (
          <button
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="theme-header-mobile-button md:hidden flex"
            aria-label="Toggle menu"
          >
            {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
          </button>
        )}
      </div>

      {/* Mobile Navigation */}
      {mobileMenuOpen &&
        !config.hideNavOnMobile &&
        hasNavigation && (
          <nav className="theme-header-mobile-nav md:hidden">
            {sortedNavigationItems.map((item) => (
              <a
                key={item.id || item.label}
                href={item.href}
                onClick={(e) => {
                  if (onNavigationClick) {
                    e.preventDefault()
                    onNavigationClick(item.href)
                  }
                  setMobileMenuOpen(false)
                }}
                className="theme-header-mobile-link"
              >
                {item.label}
              </a>
            ))}
          </nav>
        )}

      <style>{`
        .theme-header {
          background-color: ${backgroundColor};
          color: ${textColor};
          min-height: ${headerHeight}px;
          padding: 0 ${headerPadding}px;
          border-bottom: 1px solid ${textColor}22;
          box-shadow: ${showShadow ? '0 12px 32px rgba(15, 23, 42, 0.12)' : 'none'};
          position: ${config.stickyHeader ? 'sticky' : 'relative'};
          top: 0;
          z-index: 60;
          direction: ${isArabic ? 'rtl' : 'ltr'};
          backdrop-filter: ${config.stickyHeader ? 'saturate(1.15) blur(10px)' : 'none'};
        }

        .theme-header-inner {
          height: ${headerHeight}px;
          display: flex;
          align-items: center;
          justify-content: space-between;
          gap: 12px;
          width: 100%;
          max-width: 1280px;
          margin: 0 auto;
        }

        .theme-header-logo-link {
          color: inherit;
          text-decoration: none;
          display: inline-flex;
          align-items: center;
          min-width: 0;
          transition: opacity 0.2s ease;
        }

        .theme-header-logo-link:hover {
          opacity: 0.92;
        }

        .theme-header-logo {
          display: flex;
          align-items: center;
          gap: 12px;
          min-width: 0;
        }

        .theme-header-logo-image {
          height: ${config.logoHeight || 40}px;
          object-fit: contain;
          max-width: 200px;
        }

        .theme-header-logo-text {
          font-weight: 700;
          font-size: 18px;
          letter-spacing: 0.01em;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .theme-header-nav {
          gap: 8px;
          font-size: 14px;
          flex: 1;
          align-items: center;
          margin-left: ${navMarginLeft}px;
          margin-right: ${navMarginRight}px;
        }

        .theme-header-nav-horizontal {
          flex-direction: row;
        }

        .theme-header-nav-vertical {
          flex-direction: column;
        }

        .theme-header-nav-left {
          justify-content: flex-start;
        }

        .theme-header-nav-center {
          justify-content: center;
        }

        .theme-header-nav-right {
          justify-content: flex-end;
        }

        .theme-header-link {
          color: inherit;
          text-decoration: none;
          transition: opacity 0.2s ease, background-color 0.2s ease, transform 0.2s ease;
          cursor: pointer;
          opacity: 0.92;
          padding: 9px 12px;
          border-radius: 999px;
          line-height: 1;
          font-weight: 500;
        }

        .theme-header-link:hover {
          opacity: 1;
          background-color: ${textColor}12;
          transform: translateY(-1px);
        }

        .theme-header-link:focus-visible {
          outline: 2px solid ${textColor}55;
          outline-offset: 2px;
          opacity: 1;
        }

        .theme-header-mobile-button {
          background: ${textColor}0f;
          border: 1px solid ${textColor}22;
          color: inherit;
          cursor: pointer;
          padding: 8px;
          width: 40px;
          height: 40px;
          border-radius: 12px;
          align-items: center;
          justify-content: center;
          transition: background-color 0.2s ease, border-color 0.2s ease, transform 0.2s ease;
        }

        .theme-header-mobile-button:hover {
          background: ${textColor}18;
          border-color: ${textColor}30;
          transform: translateY(-1px);
        }

        .theme-header-mobile-nav {
          margin: 8px auto 12px;
          width: 100%;
          max-width: 1280px;
          padding: 10px;
          border: 1px solid ${textColor}24;
          border-radius: 16px;
          background: ${backgroundColor};
          box-shadow: 0 14px 34px rgba(15, 23, 42, 0.14);
          display: flex;
          flex-direction: column;
          gap: 6px;
        }

        .theme-header-mobile-link {
          color: inherit;
          text-decoration: none;
          padding: 12px 14px;
          transition: opacity 0.2s ease, background-color 0.2s ease;
          cursor: pointer;
          border-radius: 12px;
          font-weight: 500;
          line-height: 1.15;
          opacity: 0.92;
        }

        .theme-header-mobile-link:hover {
          opacity: 1;
          background: ${textColor}10;
        }

        .theme-header-mobile-link:focus-visible {
          outline: 2px solid ${textColor}50;
          outline-offset: 2px;
        }

        @media (max-width: 767px) {
          .theme-header {
            padding: 0 max(12px, ${headerPadding - 4}px);
          }

          .theme-header-logo-text {
            font-size: 16px;
            max-width: 60vw;
          }
        }
      `}</style>
    </header>
  )
}

export default Header
