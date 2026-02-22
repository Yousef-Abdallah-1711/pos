/**
 * Data-Driven Footer Component
 * Renders footer based on theme configuration with no inline style attributes.
 * Footer config includes company info, sections, links, social, legal, colors.
 */

import React from 'react'
import { Mail, Phone, MapPin } from 'lucide-react'
import { FooterConfig } from '../types'

export interface FooterProps {
  config: FooterConfig
  isArabic?: boolean
  onLinkClick?: (href: string) => void
}

function handleBackToTop() {
  if (typeof window !== 'undefined') {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

export const Footer: React.FC<FooterProps> = ({ config, isArabic = false, onLinkClick }) => {
  const currentYear = new Date().getFullYear()
  const layout = config.layout || 'expanded'
  const textColor = config.textColor || '#f3f4f6'
  const linkColor = config.linkColor || textColor
  const showLegalLinks = config.showLegal !== false && Boolean(config.legalLinks && config.legalLinks.length > 0)
  const showSocialLinks = config.showLinks !== false && Boolean(config.socialLinks && config.socialLinks.length > 0)
  const showBackToTop = Boolean(config.showBackToTop)

  const footerColumns =
    config.columns === 2
      ? 'repeat(auto-fit, minmax(200px, 1fr))'
      : config.columns === 4
        ? 'repeat(auto-fit, minmax(150px, 1fr))'
        : 'repeat(auto-fit, minmax(200px, 1fr))'

  const handleLinkClick = (event: React.MouseEvent<HTMLAnchorElement>, href: string) => {
    if (onLinkClick) {
      event.preventDefault()
      onLinkClick(href)
    }
  }

  if (layout === 'minimal') {
    return (
      <footer className="theme-footer theme-footer-minimal">
        <div className="theme-footer-minimal-inner">
          <p className="theme-footer-minimal-copyright">
            {config.copyrightText || `Â© ${currentYear}. All rights reserved.`}
          </p>

          {showLegalLinks && (
            <div className="theme-footer-minimal-legal-links">
              {(config.legalLinks || []).map((link) => (
                <a
                  key={`${link.label}-${link.href}`}
                  href={link.href}
                  onClick={(event) => handleLinkClick(event, link.href)}
                  className="theme-footer-link"
                >
                  {link.label}
                </a>
              ))}
            </div>
          )}
        </div>

        <style>{`
          .theme-footer {
            background-color: ${config.backgroundColor || '#111827'};
            color: ${textColor};
            direction: ${isArabic ? 'rtl' : 'ltr'};
          }

          .theme-footer-minimal {
            padding: 24px;
          }

          .theme-footer-minimal-inner {
            max-width: 1280px;
            margin: 0 auto;
            display: flex;
            justify-content: space-between;
            align-items: center;
            gap: 16px;
            flex-wrap: wrap;
          }

          .theme-footer-minimal-copyright {
            margin: 0;
            font-size: 14px;
            opacity: 0.85;
          }

          .theme-footer-minimal-legal-links {
            display: flex;
            gap: 16px;
            flex-wrap: wrap;
          }

          .theme-footer-link {
            color: ${linkColor};
            text-decoration: none;
            opacity: 0.85;
            transition: opacity 0.2s;
            cursor: pointer;
          }

          .theme-footer-link:hover {
            opacity: 1;
          }
        `}</style>
      </footer>
    )
  }

  return (
    <footer className="theme-footer theme-footer-expanded">
      <div className="theme-footer-inner">
        {(config.companyName || config.companyDescription || config.address || config.phone || config.email) && (
          <div className="theme-footer-company">
            {config.companyName && <h3 className="theme-footer-company-name">{config.companyName}</h3>}
            {config.companyDescription && (
              <p className="theme-footer-company-description">{config.companyDescription}</p>
            )}
            <div className="theme-footer-contact-list">
              {config.address && (
                <div className="theme-footer-contact-item">
                  <MapPin size={16} className="theme-footer-icon theme-footer-icon-pin" />
                  <span>{config.address}</span>
                </div>
              )}
              {config.phone && (
                <div className="theme-footer-contact-item">
                  <Phone size={16} className="theme-footer-icon" />
                  <a
                    href={`tel:${config.phone}`}
                    onClick={(event) => handleLinkClick(event, `tel:${config.phone}`)}
                    className="theme-footer-contact-link"
                  >
                    {config.phone}
                  </a>
                </div>
              )}
              {config.email && (
                <div className="theme-footer-contact-item">
                  <Mail size={16} className="theme-footer-icon" />
                  <a
                    href={`mailto:${config.email}`}
                    onClick={(event) => handleLinkClick(event, `mailto:${config.email}`)}
                    className="theme-footer-contact-link"
                  >
                    {config.email}
                  </a>
                </div>
              )}
            </div>
          </div>
        )}

        {config.footerSections && config.footerSections.length > 0 && (
          <div className="theme-footer-sections">
            {config.footerSections.map((section) => (
              <div key={section.id || section.title}>
                <h4 className="theme-footer-section-title">{section.title}</h4>
                <ul className="theme-footer-section-links">
                  {(section.links || []).map((link) => (
                    <li key={`${section.id || section.title}-${link.label}-${link.href}`}>
                      <a
                        href={link.href}
                        onClick={(event) => handleLinkClick(event, link.href)}
                        className="theme-footer-link"
                      >
                        {link.label}
                      </a>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        )}

        {showSocialLinks && (
          <div className="theme-footer-social">
            <h4 className="theme-footer-social-title">{isArabic ? 'ØªØ§Ø¨Ø¹Ù†Ø§' : 'Follow Us'}</h4>
            <div className="theme-footer-social-links">
              {(config.socialLinks || []).map((link) => (
                <a
                  key={link.id || link.platform}
                  href={link.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  onClick={(event) => handleLinkClick(event, link.url)}
                  className="theme-footer-link theme-footer-social-link"
                >
                  {link.platform}
                </a>
              ))}
            </div>
          </div>
        )}

        {(showLegalLinks || showBackToTop) && (
          <div className="theme-footer-legal-row">
            {showLegalLinks ? (
              <div className="theme-footer-legal-links">
                {(config.legalLinks || []).map((link, idx) => (
                  <div key={`${link.label}-${link.href}`} className="theme-footer-legal-item">
                    {idx > 0 && <span className="theme-footer-legal-separator">|</span>}
                    <a
                      href={link.href}
                      onClick={(event) => handleLinkClick(event, link.href)}
                      className="theme-footer-link"
                    >
                      {link.label}
                    </a>
                  </div>
                ))}
              </div>
            ) : (
              <span />
            )}

            {showBackToTop && (
              <button type="button" onClick={handleBackToTop} className="theme-footer-top-button">
                {isArabic ? 'Ø§Ù„Ø¹ÙˆØ¯Ø© Ù„Ù„Ø£Ø¹Ù„Ù‰' : 'Back to top'}
              </button>
            )}
          </div>
        )}

        <div className="theme-footer-copyright">
          <p>{config.copyrightText || `Â© ${currentYear}. All rights reserved.`}</p>
        </div>
      </div>

      <style>{`
        .theme-footer {
          background-color: ${config.backgroundColor || '#111827'};
          color: ${textColor};
          direction: ${isArabic ? 'rtl' : 'ltr'};
        }

        .theme-footer-expanded {
          padding: ${config.padding || 48}px;
        }

        .theme-footer-inner {
          max-width: 1280px;
          margin: 0 auto;
        }

        .theme-footer-company {
          margin-bottom: 32px;
        }

        .theme-footer-company-name {
          font-size: 20px;
          font-weight: 700;
          margin-bottom: 16px;
          color: ${linkColor};
        }

        .theme-footer-company-description {
          opacity: 0.8;
          margin-bottom: 16px;
        }

        .theme-footer-contact-list {
          font-size: 14px;
          opacity: 0.8;
        }

        .theme-footer-contact-item {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 8px;
        }

        .theme-footer-icon {
          flex-shrink: 0;
        }

        .theme-footer-icon-pin {
          margin-top: 4px;
          align-self: flex-start;
        }

        .theme-footer-contact-link {
          color: inherit;
          text-decoration: none;
          cursor: pointer;
        }

        .theme-footer-sections {
          display: grid;
          grid-template-columns: ${footerColumns};
          gap: 32px;
          margin-bottom: 32px;
          border-top: 1px solid ${textColor}20;
          padding-top: 32px;
        }

        .theme-footer-section-title {
          font-weight: 700;
          margin-bottom: 16px;
          color: ${linkColor};
        }

        .theme-footer-section-links {
          list-style: none;
          padding: 0;
          margin: 0;
          display: flex;
          flex-direction: column;
          gap: 8px;
        }

        .theme-footer-link {
          color: ${linkColor};
          text-decoration: none;
          opacity: 0.8;
          transition: opacity 0.2s;
          cursor: pointer;
        }

        .theme-footer-link:hover {
          opacity: 1;
        }

        .theme-footer-social {
          border-top: 1px solid ${textColor}20;
          padding-top: 32px;
          margin-bottom: 32px;
        }

        .theme-footer-social-title {
          text-align: center;
          font-weight: 700;
          margin-bottom: 16px;
          color: ${linkColor};
        }

        .theme-footer-social-links {
          display: flex;
          justify-content: center;
          gap: 16px;
          flex-wrap: wrap;
        }

        .theme-footer-social-link {
          text-transform: capitalize;
        }

        .theme-footer-legal-row {
          border-top: 1px solid ${textColor}20;
          padding-top: 24px;
          margin-bottom: 24px;
          display: flex;
          justify-content: space-between;
          align-items: center;
          gap: 16px;
          flex-wrap: wrap;
          font-size: 12px;
          opacity: 0.8;
        }

        .theme-footer-legal-links {
          display: flex;
          align-items: center;
          gap: 8px;
          flex-wrap: wrap;
        }

        .theme-footer-legal-item {
          display: flex;
          align-items: center;
          gap: 8px;
        }

        .theme-footer-legal-separator {
          opacity: 0.7;
        }

        .theme-footer-top-button {
          border: none;
          background: transparent;
          color: ${linkColor};
          cursor: pointer;
          font-size: 12px;
          font-weight: 600;
          padding: 0;
        }

        .theme-footer-copyright {
          border-top: 1px solid ${textColor}20;
          padding-top: 24px;
          text-align: center;
          opacity: 0.8;
          font-size: 14px;
        }
      `}</style>
    </footer>
  )
}

export default Footer
