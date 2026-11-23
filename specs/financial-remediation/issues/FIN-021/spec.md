# FIN-021: Root HTML language/direction does not reflect the active locale

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P3. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-021 — Root HTML language/direction does not reflect the active locale

- **Severity/status:** P3 — REPRODUCED in a current-run browser check; localized report content itself did render RTL.
- **Problem / impact:** The browser rendered an Arabic financial report under a nested locale wrapper with `dir="rtl"`, but the document root retained `lang="en"` and no `dir`. Assistive technology, browser text handling, and document-level direction-sensitive behavior can therefore use English/LTR semantics around Arabic content. This is an accessibility/localization defect, not evidence of incorrect financial values.
- **Current / expected:** `frontend/apps/dashboard/src/app/layout.tsx:20-21` hardcodes `<html lang="en">` and leaves `<body>` without a direction. The locale subtree controls its own direction. Expected: root document language and direction correspond to the selected locale, or the application has a documented locale-root strategy that ensures document semantics are correct.
- **Root cause:** Root layout is not locale-aware while localized routes nest under `[locale]`; a nested element cannot change the document `<html>` attributes.
- **Reproduction/evidence:** In a local dashboard browser session on 2026-11-12, Arabic Balance Sheet content appeared, the localized wrapper computed RTL, but `document.documentElement.lang` was `en`, `document.documentElement.dir` was empty, and body direction was LTR. No page/console errors; mobile viewport had no horizontal overflow; PDF download completed. API/auth responses were mocked, so this did not exercise a real authenticated backend workflow.
- **Baseline scope:** Source at fetched `origin/main` and current worktree contains the hardcoded root layout; no independent clean-main browser build was run. Treat as present in both by source comparison, with runtime behavior reproduced on the current worktree.
- **Recommended solution / components:** Make the root layout locale-aware using the framework-supported locale routing pattern and set document-level `lang`/`dir`; verify English and Arabic routes, hydration, nested direction, assistive semantics, and report export.
- **Migration:** None.
- **Dependencies/policy:** Supported locale list and language/direction mapping; frontend routing architecture.
- **Acceptance:** Browser assertions verify root lang/dir for each supported locale; Arabic report and export remain usable; no hydration warnings or horizontal overflow at agreed viewports.
- **Checklist:** OPEN.


## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
