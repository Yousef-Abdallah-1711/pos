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
          background-color: ${config.backgroundColor || '#ffffff'};
          color: ${textColor};
          height: ${config.height || 64}px;
          padding: 0 ${config.padding || 16}px;
          box-shadow: ${showShadow ? '0 1px 3px rgba(0,0,0,0.1)' : 'none'};
          position: ${config.stickyHeader ? 'sticky' : 'relative'};
          top: 0;
          z-index: 50;
          direction: ${isArabic ? 'rtl' : 'ltr'};
        }

        .theme-header-inner {
          height: 100%;
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .theme-header-logo-link {
          color: inherit;
          text-decoration: none;
        }

        .theme-header-logo {
          display: flex;
          align-items: center;
          gap: 16px;
        }

        .theme-header-logo-image {
          height: ${config.logoHeight || 40}px;
          object-fit: contain;
        }

        .theme-header-logo-text {
          font-weight: 700;
          font-size: 18px;
        }

        .theme-header-nav {
          gap: 24px;
          font-size: 14px;
          flex: 1;
          margin-left: ${navLeftGap}px;
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
          transition: opacity 0.2s;
          cursor: pointer;
        }

        .theme-header-link:hover {
          opacity: 0.8;
        }

        .theme-header-mobile-button {
          background: none;
          border: none;
          color: inherit;
          cursor: pointer;
          padding: 8px;
        }

        .theme-header-mobile-nav {
          margin-top: 16px;
          padding-top: 16px;
          border-top: 1px solid ${textColor}40;
          display: flex;
          flex-direction: column;
          gap: 8px;
        }

        .theme-header-mobile-link {
          color: inherit;
          text-decoration: none;
          padding: 8px 16px;
          transition: opacity 0.2s;
          cursor: pointer;
        }

        .theme-header-mobile-link:hover {
          opacity: 0.8;
        }
      `}</style>
    </header>
  )
}

export default Header
