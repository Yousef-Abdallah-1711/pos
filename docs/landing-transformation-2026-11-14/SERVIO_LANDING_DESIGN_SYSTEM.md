# SERVIO landing design system

The site leads with the cashier-to-kitchen handoff, then explains the supported tasks with actual synthetic application captures. SERVIO is the public brand; violet is the action color.

## Visual language

Manrope supports English; IBM Plex Sans Arabic supports Arabic. Headings use fluid 32–78px sizing and restrained weight. Logical CSS properties preserve RTL composition. The page alternates paper, pale violet, dark kitchen, and saturated final-action sections. Maximum content width is 1280px. Mobile has 20px side gutters; tablet and desktop have 40px gutters. Interactive controls target 48px minimum height. Focus uses a visible violet outline. Motion is limited to action feedback and respects reduced motion.

The hero combines authentic cashier and kitchen screens. Editorial rows, an expandable journey, a kitchen state ledger, paired guest screens, a stock quantity ledger, employee records, and a tabbed walkthrough provide different reading rhythms. Screens remain original captures; CSS framing does not fabricate product functionality.

## Research

Primary references inspected: [Linear](https://linear.app), [Stripe](https://stripe.com), [Vercel](https://vercel.com), [Toast](https://pos.toasttab.com), [Square Restaurants](https://squareup.com/us/en/point-of-sale/restaurants). Adopted principles: clear category positioning, contextual product evidence, deliberate light/dark rhythm, a meaningful mobile composition, and one subscription with two frequencies. No reference assets or branding were copied. Toast's browser request was blocked; only its public text extraction was available. Reference screenshots are research material, not shipped assets.

## Implementation

`frontend/apps/landing/src/app/marketing.css` defines the scoped visual system. Marketing content is server-rendered. Client components implement navigation, billing selection, walkthrough tabs, and local funnel events. Next Image receives actual capture dimensions. Below-fold captures lazy load. Public capture provenance is recorded beside the assets.
