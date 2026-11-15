# SERVIO full platform UI/UX audit baseline

This report starts the route by route audit from the repository snapshot at
`867a5721089d435e75248e33e89cb9b1dbcdbfb9` (15 November 2026). It records actual
Next.js page files as route patterns. Dynamic segments remain explicit so they
can be expanded into seeded-record cases during functional acceptance.

## Discovered scope

| Application | Route patterns | Locale baseline | Theme baseline |
| --- | ---: | --- | --- |
| Landing and authentication | 17 | English, Arabic | Light, dark |
| Restaurant dashboard | 79 | English, Arabic | Light, dark |
| Restaurant storefront | 12 | English, Arabic | Tenant configured |
| Super Admin | 9 | English only | OS preference, no explicit control |
| **Total** | **117** |  |  |

The repository also contains 201 component source files and 154 app-level
loading, error, and not-found boundary files. Those are component and state
surfaces, not 154 additional routes. The CSV records an authentication/role
baseline and marks every responsive, screenshot, accessibility, and functional
result `pending` until direct verification is attached.

## Verified baseline defect

Super Admin root CSS switches its background and foreground from the operating
system `prefers-color-scheme` media query. The app configures Tailwind `darkMode`
as `class`, but has no root theme provider or visible theme control. This makes
the document background follow OS dark mode while the application's `dark:`
component classes remain inactive. Super Admin pages also have no locale route,
locale provider, or language control; the root document is fixed to English.

## Audit limits at this checkpoint

The route list is exhaustive for discovered `page.tsx` files in the four frontend
apps. It is not a claim that every dynamic record, role, modal, theme, locale,
viewport, or interaction has been audited. The acceptance scripts currently
cover selected authenticated dashboard/POS/KDS and Super Admin routes at four
widths, and broad anonymous marketing/storefront samples. They do not establish
full route by locale by theme visual acceptance. The fields in the CSV remain
pending until evidence is captured and reviewed.

Source-derived baseline issue also observed in the dashboard overview: its
loading message is removed only after the request settles, but an API failure
can leave the screen without a visible error or retry state. Existing trial
registration screenshot capture checks subscription creation but can capture
the dashboard before overview data has settled. These require an authenticated
synthetic flow to verify and are listed as follow-up acceptance work.
