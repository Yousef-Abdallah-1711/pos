# Dashboard color discovery

The [raw inventory](color-inventory.csv) contains 12,305 candidate matches across 424 dashboard source files: hex literals, CSS color functions, selected CSS variables and common Tailwind color utilities. It includes repeated and non-brand values. The generator is a source search, not a CSS evaluator; arbitrary utilities, computed colors, chart libraries, package consumers and runtime-injected values can be missed. There are 609 matched `blue-` utility occurrences, not 609 confirmed defects.

Most frequent matched utilities include `text-gray-900` (798), `bg-gray-700` (788), `text-gray-400` (730), `text-gray-300` (564), `text-gray-700` (506), `border-gray-200` (479), `bg-gray-200` (468), and `border-gray-700` (459). Counts include light and dark variants and duplicate usage. Red utilities such as `bg-red-50` and `text-red-700` also occur frequently and may indicate validation, danger or status; do not reinterpret them automatically as brand colors.

## Verified mechanisms and migration starting points

| Current mechanism | Evidence | Future review |
| --- | --- | --- |
| User-derived brand shades | `src/stores/preferencesStore.ts` writes `--primary-*`, `--secondary-*`, and `--accent-*`; `tailwind.config.js` maps these variables to color utilities. | Decide compatibility and precedence before applying tenant tokens. |
| Shell surfaces, dividers and text | `src/components/layout/DashboardLayout.tsx`, `Sidebar.tsx`, `src/app/globals.css` contain direct gray/white dark variants. | Map structural roles to light/dark workspace tokens; verify contrast and RTL. |
| Navigation active/hover | `Sidebar.tsx:433`, `:455`, `:475` combines gray hover and primary active utilities. | Preserve selected, hover and focus distinction in both modes and forced colors. |
| Semantic statuses and financial visuals | Candidate colored utilities occur throughout dashboard modules and charts. | Classify by use before migration; keep success/warning/danger/info and chart differentiation independently controlled. |
| Storefront branding | Storefront provider injects `--theme-*`; the appearance page currently synchronizes selected personal colors to an active site theme. | Preserve public-site scope and existing customer data during any future separation. |

`src/app/globals.css` defines initial `--primary-500: #3b82f6`, `--secondary-500: #6366f1`, and `--accent-500: #22c55e`; its dark class changes the base `--background`/`--foreground` HSL values, while `body` uses explicit gray Tailwind utilities. The original personal-settings migration and preferences-store fallback use `#10b981` for the accent instead of the initial CSS accent. The resulting first-paint and restored-state behavior needs runtime verification; do not assume the two defaults are interchangeable or that all structural colors inherit the CSS variables.

The preference store's shade generator retains the selected color's hue and saturation but sets the 500 shade's HSL lightness to 50%. Consequently, the generated `--primary-500`, `--secondary-500`, or `--accent-500` can differ from the selected input color. A future migration should explicitly decide whether to preserve this behavior or preserve the exact selected color.

The provisional [token registry](proposed-token-registry.json) is a design proposal, not a completed mapping or approved accessible palette. Complete source classification, computed-style checks, chart and third-party review, tenant isolation, and visual comparison before global rollout.
