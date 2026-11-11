# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

Independent restaurant owners and operators, plus people managing multi-location restaurant groups. They use POOS to coordinate day-to-day restaurant work across front- and back-of-house operations.

## Product Purpose

POOS is a restaurant management platform. It brings restaurant operations into one product so owners and teams can manage connected work from a shared system.

## Positioning

All-in-one restaurant management system. Approved differentiator: one platform for front and back of house.

## Operating Context

Users access the web application from restaurant and office devices. Public onboarding supports owner name, email, password, restaurant name, restaurant subdomain, and optional phone. The landing application supports Arabic and English, with Arabic right-to-left and English left-to-right layouts.

## Capabilities and Constraints

- Public registration and restaurant subdomain availability checks are implemented through the existing authentication API.
- Login supports users associated with multiple tenants and preserves the existing tenant-selection flow.
- The codebase contains restaurant point-of-sale, order, kitchen, inventory, staff, and online storefront functionality; describe individual capabilities publicly only when their current end-to-end behavior is verified.
- Pricing, trial length, payment terms, and comparative performance claims are not approved for public display.
- Do not claim unverified customer counts, uptime, growth results, mobile apps, or other metrics and capabilities.
- The existing forgot-password entry points to a route that is not implemented in the landing app. The API handler exists, but public reset completion is not verified; do not present password recovery as a supported flow until it is confirmed end to end.

## Brand Commitments

- Name: POOS.
- Arabic-first public experience with English available.
- Primary purple `#6D28D9`, secondary purple `#8B5CF6`, soft lavender `#C4B5FD`, background `#F7F5FF`, dark ink `#171327`, and white `#FFFFFF`.
- Alexandria for Arabic headings and IBM Plex Sans Arabic for Arabic interface and body text; use suitable Latin fallbacks for English.
- Brand personality: premium, modern, trustworthy, clean, distinctive, technology-focused, and business-oriented.
- Use semantic design tokens and logical layout properties for RTL/LTR behavior.

## Evidence on Hand

- The repository contains application screenshots under `audit-reports/phase-1-system-discovery/evidence/screenshots/`, including current landing, storefront, and super-admin captures. These document the existing UI; they are not approved customer proof.
- No POOS customer logos, named testimonials, case studies, or independent performance evidence are approved for public use.
- Use real product UI captures only. Do not fabricate data or imply a synthetic dashboard is a live account.

## Product Principles

- Keep restaurant work connected across front- and back-of-house operations.
- Preserve the real authentication, tenant selection, and onboarding behavior.
- Make product and commercial claims only from verified functionality and approved evidence.
- Support Arabic and English without compromising either direction or language.
