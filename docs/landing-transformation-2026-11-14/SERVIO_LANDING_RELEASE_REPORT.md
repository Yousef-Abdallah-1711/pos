# Release report

Status: implementation and verification in progress. Branch: `feat/servio-product-led-website`. Baseline main: `0d8fea4`; audit preservation commit: `ecfa8e9`.

The authorized release path is a pull request, all required verified CI jobs, protected main merge, signed immutable main image release, existing outbound AWS pull agent, then live receipt/SHA and workflow checks. No merge, successful CI, deployment, or final release verdict is claimed yet.

## Migration review

Migration 186 changes only the future trial default, adds a nullable constrained billing preference, inserts two versioned commercial offers, and removes legacy offers from new-sale discovery. Existing tenant dates, plan references, payment amounts and snapshots remain intact. Review the exact migration hash before adding it to the existing approved-migrations gate. Preserve database volumes, uploads, secrets, TLS state, backups and previous release images. Paid checkout and staging noindex remain disabled/protected.

Rollback uses the prior signed image release. The new nullable column and versioned catalog entries can remain during an application rollback; do not drop data-bearing columns after live writes. Catalog retirement reversal requires an explicit reviewed action. Never bypass required CI, force-push main, or substitute an unsigned feature image for the requested main release.

Pending: full integration/build CI, final visual/accessibility/performance evidence, PR/merge identifiers, signed release receipt and live acceptance. The final verdict must remain IN PROGRESS until those gates have actual evidence.
