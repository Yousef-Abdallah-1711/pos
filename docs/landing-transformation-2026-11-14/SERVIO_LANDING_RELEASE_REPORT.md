# Release report

## Delivery record

[PR #43](https://github.com/Yousef-Abdallah-1711/pos/pull/43) is the authoritative delivery record for merge state, exact approved main SHA, final CI run, signed release, deployed SHA, and live verification. These values are updated in the PR description after merge/deployment so the report does not pretend to know its own future merge hash. The code must not be called deployed until the runtime receipt agrees with that main release.

The built candidate at PR head `9e897f6ad2aa2bfbeaf0a539c37320d2be32ad4c` passed [all verified CI stages](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43575769431), including all five application builds and browser acceptance. See [test results](SERVIO_LANDING_PLAYWRIGHT_RESULTS.md) and [visual QA](SERVIO_LANDING_VISUAL_QA.md). The final gate additionally covers every public route at all requested widths and guards against replacing historical plan prices during a catalog outage.

## Delivered scope

Fourteen homepage sections; thirteen dedicated public routes per locale, including seven feature pages; EN/AR navigation and server-rendered document direction; authentic synthetic screenshots; contextual trial actions; product walkthrough; approved pricing; 14-day new-account entitlement; persisted billing preference; localized dashboard offers; unavailable paid checkout; metadata, sitemap, WebSite structured data, privacy-respecting local funnel events, and regression/performance evidence.

Canonical offer: 750 EGP/month, 5,000 EGP/year, 14-day trial. Annual saving is 4,000 EGP, approximately 44.4%. No existing trial timestamps, subscription references or historical payment amounts are rewritten. Callbacks validate the stored payment amount rather than today's catalog price.

## AWS prerequisite — not yet satisfied

The verified root device serial is `vol0ec6bfabc3638a90c`; EBS volume `vol-0ec6bfabc3638a90c` in eu-north-1 remains 30 GiB, with approximately 3.9 GiB free. The existing [runbook](../../ops/aws/RUNBOOK.md#prepared-disk-expansion--owner-approval-required) requires compressed release archives plus 4 GiB reserve. It prepares a 30→50 GiB expansion, preserving type/IOPS/throughput, explicitly requiring owner approval because storage cost increases. The host has no cloud-side AWS role. Approval/operator completion was requested; no paid AWS resource or disk expansion was performed by this task.

After the cloud device actually reports the approved new size, verify serial, partition and mounted filesystem before extending partition 1/ext4. Preserve backups, uploads, secrets, certificates and previous images. Do not prune working images/volumes or reduce the reserve to force deployment.

## Migration and deployment review

Migration `186_servio_trial_offer.sql` SHA-256:

`2368ef9595224b21d9a77febd6ea38e3427c8e03a9dfe672e05b6f4b47028fec`

It changes only future trial defaults, adds nullable monthly/annual preference, inserts two versioned offers, and retires old offers from new-sale discovery while retaining historical lookups. The registration transaction uses the database creation timestamp and approved duration. Forward migration and registration behavior passed disposable PostgreSQL integration. Pin this exact reviewed migration through the existing approved-migrations gate before the authorized main release deployment.

Release path: protected PR → all verified CI stages → main merge → signed immutable main release → existing outbound pull agent → health/receipt/SHA checks → live EN/AR and synthetic owner workflow verification. Current observed main protection requires zero approving reviews, strict ancestry, and backend/runtime/domain checks; this task additionally requires the entire verified workflow, including image/browser and aggregate checks. No protection is weakened.

Rollback uses the previous signed application images and existing recovery procedure. New nullable data/catalog entries can remain for application compatibility. Do not drop a data-bearing column after live writes or automatically reverse migration SQL. Preserve staging noindex and the paid-checkout gates.

## Scope limits

No real customer/provider charge or message was sent. Live demonstration booking is not advertised; the product walkthrough and trial are real destinations. No durable external analytics adapter, completed onboarding wizard, comprehensive finance certification, offline guarantee, or universal product-language completeness is claimed. Physical-device and NVDA/VoiceOver checks remain outside the recorded browser evidence.

AWS delivery remains **BLOCKED** until disk capacity is available and the exact signed main release passes live verification. A green PR build alone is not a COMPLETE verdict.
