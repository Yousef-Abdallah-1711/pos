# Defects and remediation

No source code was changed, committed, or pushed in this audit.

Confirmed/high-confidence findings:

1. All four clean-volume Next.js builds fail while prerendering `/404` because Next.js `Html` is imported outside `pages/_document`; dashboard additionally fails prerendering `/ar/dashboard/settings` with null `useContext`.
2. Browser client API URL targets `localhost:18080`, but the isolated test API was on host port `18082`; direct browser requests fail until an external temporary forward is supplied. This must be corrected in the supported runtime configuration and retested without a shim.
3. Landing emits a duplicate React key `#product` console error.
4. Authenticated dashboard at the requested mobile resize has horizontal overflow (tool-reported document 603px vs viewport 487px).
5. KDS reports `Disconnected` while API-backed data and a status update work; realtime transport acceptance is unresolved.
6. RLS is enabled on 15/92 public tables; core tenant-bearing tables listed in report 08 do not have DB RLS.

Remediation was paused because `C:\Users\OPT\Desktop\POS\CODING_STANDARDS.md`, required by AGENTS.md, is missing. The next coding pass must obtain the correct standards, reproduce/fix each build defect, and rerun the affected acceptance checks. Do not treat runtime pass evidence as a source-level security review.


7. The disposable restaurant storefront resolved to its maintenance page; storefront workflows could not be accepted. Super-admin page rendered but its API session check hit the unavailable localhost:18080 port after the temporary forward was stopped.

