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
  const backgroundColor = config.backgroundColor || '#0f172a'
  const textColor = config.textColor || '#f3f4f6'
  const linkColor = config.linkColor || textColor
  const basePadding = Math.max(16, config.padding || 48)
  const sectionGap = layout === 'compact' ? 20 : 28
  const shellRadius = layout === 'compact' ? 20 : 24
  const innerPadding = layout === 'compact' ? 20 : 28
  const minimalPadding = layout === 'minimal' ? Math.max(16, Math.round(basePadding * 0.5)) : basePadding

  const showLegalLinks =
    config.showLegal !== false && Boolean(config.legalLinks && config.legalLinks.length > 0)
  const showSocialLinks =
    config.showLinks !== false && Boolean(config.socialLinks && config.socialLinks.length > 0)
  const showBackToTop = Boolean(config.showBackToTop)
  const hasCompanyBlock = Boolean(
    config.companyName || config.companyDescription || config.address || config.phone || config.email
  )
  const hasFooterSections = Boolean(config.footerSections && config.footerSections.length > 0)

  const footerColumns =
    config.columns === 2
      ? 'repeat(auto-fit, minmax(220px, 1fr))'
      : config.columns === 4
        ? 'repeat(auto-fit, minmax(160px, 1fr))'
        : 'repeat(auto-fit, minmax(190px, 1fr))'

  const followUsLabel = isArabic ? '\u062a\u0627\u0628\u0639\u0646\u0627' : 'Follow Us'
  const backToTopLabel = isArabic ? '\u0627\u0644\u0639\u0648\u062f\u0629 \u0644\u0644\u0623\u0639\u0644\u0649' : 'Back to top'
  const defaultCopyright = `\u00A9 ${currentYear}. All rights reserved.`

  const handleLinkClick = (event: React.MouseEvent<HTMLAnchorElement>, href: string) => {
    if (onLinkClick) {
      event.preventDefault()
      onLinkClick(href)
    }
  }

  if (layout === 'minimal') {
    return (
      <footer className="theme-footer theme-footer-minimal">
        <div className="theme-footer-minimal-shell">
          <div className="theme-footer-minimal-inner">
            <p className="theme-footer-minimal-copyright">
              {config.copyrightText || defaultCopyright}
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
        </div>

        <style>{`
          .theme-footer {
            background: linear-gradient(180deg, ${backgroundColor}, ${backgroundColor});
            color: ${textColor};
            direction: ${isArabic ? 'rtl' : 'ltr'};
            border-top: 1px solid ${textColor}18;
          }

          .theme-footer-minimal {
            padding: ${minimalPadding}px;
          }

          .theme-footer-minimal-shell {
            max-width: 1280px;
            margin: 0 auto;
            border: 1px solid ${textColor}16;
            background: ${textColor}08;
            border-radius: 18px;
            box-shadow: 0 16px 40px rgba(2, 6, 23, 0.18);
          }

          .theme-footer-minimal-inner {
            padding: 16px 18px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 12px 18px;
            flex-wrap: wrap;
          }

          .theme-footer-minimal-copyright {
            margin: 0;
            font-size: 13px;
            opacity: 0.88;
          }

          .theme-footer-minimal-legal-links {
            display: flex;
            gap: 8px;
            flex-wrap: wrap;
          }

          .theme-footer-link {
            color: ${linkColor};
            text-decoration: none;
            opacity: 0.9;
            transition: opacity 0.2s ease, background-color 0.2s ease, transform 0.2s ease;
            cursor: pointer;
            padding: 8px 10px;
            border-radius: 999px;
            line-height: 1;
          }

          .theme-footer-link:hover {
            opacity: 1;
            background: ${textColor}12;
            transform: translateY(-1px);
          }

          .theme-footer-link:focus-visible {
            outline: 2px solid ${linkColor}66;
            outline-offset: 2px;
            opacity: 1;
          }

          @media (max-width: 640px) {
            .theme-footer-minimal-inner {
              align-items: flex-start;
              justify-content: flex-start;
            }
          }
        `}</style>
      </footer>
    )
  }

  return (
    <footer
      className={`theme-footer ${
        layout === 'compact' ? 'theme-footer-compact' : 'theme-footer-expanded'
      }`}
    >
      <div className="theme-footer-shell">
        <div className="theme-footer-inner">
          {hasCompanyBlock && (
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

          {hasFooterSections && (
            <div className="theme-footer-sections">
              {(config.footerSections || []).map((section) => (
                <div key={section.id || section.title} className="theme-footer-section-card">
                  <h4 className="theme-footer-section-title">{section.title}</h4>
                  <ul className="theme-footer-section-links">
                    {(section.links || []).map((link) => (
                      <li key={`${section.id || section.title}-${link.label}-${link.href}`}>
                        <a
                          href={link.href}
                          onClick={(event) => handleLinkClick(event, link.href)}
                          className="theme-footer-link theme-footer-link-inline"
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
              <h4 className="theme-footer-social-title">{followUsLabel}</h4>
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
                        className="theme-footer-link theme-footer-link-inline"
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
                  {backToTopLabel}
                </button>
              )}
            </div>
          )}

          <div className="theme-footer-copyright">
            <p>{config.copyrightText || defaultCopyright}</p>
          </div>
        </div>
      </div>

      <style>{`
        .theme-footer {
          background: linear-gradient(180deg, ${backgroundColor}, ${backgroundColor});
          color: ${textColor};
          direction: ${isArabic ? 'rtl' : 'ltr'};
          padding: ${basePadding}px;
          border-top: 1px solid ${textColor}16;
        }

        .theme-footer-compact {
          padding: ${Math.max(16, Math.round(basePadding * 0.75))}px;
        }

        .theme-footer-shell {
          max-width: 1280px;
          margin: 0 auto;
          border: 1px solid ${textColor}14;
          background:
            radial-gradient(circle at 12% 8%, ${textColor}08 0%, transparent 44%),
            radial-gradient(circle at 90% 12%, ${linkColor}10 0%, transparent 42%),
            ${textColor}06;
          border-radius: ${shellRadius}px;
          box-shadow: 0 24px 60px rgba(2, 6, 23, 0.22);
          overflow: hidden;
        }

        .theme-footer-inner {
          padding: ${innerPadding}px;
        }

        .theme-footer-company {
          padding: ${layout === 'compact' ? 16 : 20}px;
          border-radius: ${layout === 'compact' ? 16 : 18}px;
          border: 1px solid ${textColor}14;
          background: ${textColor}06;
          margin-bottom: ${sectionGap}px;
        }

        .theme-footer-company-name {
          font-size: ${layout === 'compact' ? 18 : 20}px;
          font-weight: 700;
          margin: 0 0 10px;
          color: ${linkColor};
          letter-spacing: 0.01em;
        }

        .theme-footer-company-description {
          margin: 0 0 14px;
          opacity: 0.86;
          line-height: 1.65;
        }

        .theme-footer-contact-list {
          display: grid;
          gap: 8px;
          font-size: 14px;
          opacity: 0.92;
        }

        .theme-footer-contact-item {
          display: flex;
          align-items: center;
          gap: 8px;
          min-width: 0;
        }

        .theme-footer-icon {
          flex-shrink: 0;
          opacity: 0.88;
        }

        .theme-footer-icon-pin {
          margin-top: 2px;
          align-self: flex-start;
        }

        .theme-footer-contact-link {
          color: inherit;
          text-decoration: none;
          cursor: pointer;
          opacity: 0.95;
          transition: opacity 0.2s ease;
        }

        .theme-footer-contact-link:hover {
          opacity: 1;
        }

        .theme-footer-sections {
          display: grid;
          grid-template-columns: ${footerColumns};
          gap: ${sectionGap}px;
          margin-bottom: ${sectionGap}px;
        }

        .theme-footer-section-card {
          border: 1px solid ${textColor}12;
          background: ${textColor}05;
          border-radius: 16px;
          padding: 16px;
          min-height: 100%;
        }

        .theme-footer-section-title {
          margin: 0 0 12px;
          font-weight: 700;
          color: ${linkColor};
          font-size: 14px;
          letter-spacing: 0.02em;
          text-transform: uppercase;
        }

        .theme-footer-section-links {
          list-style: none;
          padding: 0;
          margin: 0;
          display: flex;
          flex-direction: column;
          gap: 6px;
        }

        .theme-footer-link {
          color: ${linkColor};
          text-decoration: none;
          opacity: 0.9;
          transition: opacity 0.2s ease, background-color 0.2s ease, transform 0.2s ease;
          cursor: pointer;
          border-radius: 999px;
        }

        .theme-footer-link:hover {
          opacity: 1;
        }

        .theme-footer-link:focus-visible {
          outline: 2px solid ${linkColor}66;
          outline-offset: 2px;
          opacity: 1;
        }

        .theme-footer-link-inline {
          display: inline-flex;
          align-items: center;
          padding: 6px 8px;
          margin-left: -8px;
          margin-right: -8px;
          max-width: calc(100% + 16px);
        }

        .theme-footer-link-inline:hover {
          background: ${textColor}10;
          transform: translateX(${isArabic ? '-1px' : '1px'});
        }

        .theme-footer-social {
          border-top: 1px solid ${textColor}12;
          padding-top: ${sectionGap}px;
          margin-bottom: ${sectionGap}px;
          display: flex;
          align-items: center;
          justify-content: space-between;
          gap: 12px 20px;
          flex-wrap: wrap;
        }

        .theme-footer-social-title {
          margin: 0;
          font-weight: 700;
          color: ${linkColor};
          font-size: 14px;
          letter-spacing: 0.02em;
          text-transform: uppercase;
        }

        .theme-footer-social-links {
          display: flex;
          gap: 8px;
          flex-wrap: wrap;
        }

        .theme-footer-social-link {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          padding: 8px 12px;
          border: 1px solid ${textColor}14;
          background: ${textColor}06;
          text-transform: capitalize;
          font-size: 13px;
          line-height: 1;
        }

        .theme-footer-social-link:hover {
          background: ${textColor}10;
          transform: translateY(-1px);
        }

        .theme-footer-legal-row {
          border-top: 1px solid ${textColor}12;
          padding-top: ${layout === 'compact' ? 16 : 20}px;
          margin-bottom: ${layout === 'compact' ? 16 : 20}px;
          display: flex;
          justify-content: space-between;
          align-items: center;
          gap: 12px 16px;
          flex-wrap: wrap;
          font-size: 12px;
          opacity: 0.92;
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
          opacity: 0.55;
        }

        .theme-footer-top-button {
          border: 1px solid ${textColor}18;
          background: ${textColor}07;
          color: ${linkColor};
          cursor: pointer;
          font-size: 12px;
          font-weight: 600;
          line-height: 1;
          padding: 10px 12px;
          border-radius: 999px;
          transition: background-color 0.2s ease, border-color 0.2s ease, transform 0.2s ease;
        }

        .theme-footer-top-button:hover {
          background: ${textColor}12;
          border-color: ${textColor}24;
          transform: translateY(-1px);
        }

        .theme-footer-top-button:focus-visible {
          outline: 2px solid ${linkColor}66;
          outline-offset: 2px;
        }

        .theme-footer-copyright {
          border-top: 1px solid ${textColor}12;
          padding-top: ${layout === 'compact' ? 16 : 20}px;
          text-align: center;
          opacity: 0.86;
          font-size: 13px;
        }

        .theme-footer-copyright p {
          margin: 0;
        }

        @media (max-width: 768px) {
          .theme-footer {
            padding: ${Math.max(14, Math.round(basePadding * 0.6))}px;
          }

          .theme-footer-compact {
            padding: ${Math.max(12, Math.round(basePadding * 0.5))}px;
          }

          .theme-footer-inner {
            padding: 16px;
          }

          .theme-footer-sections {
            grid-template-columns: 1fr;
          }

          .theme-footer-social {
            align-items: flex-start;
            justify-content: flex-start;
          }

          .theme-footer-legal-row {
            flex-direction: column;
            align-items: flex-start;
          }
        }
      `}</style>
    </footer>
  )
}

export default Footer
