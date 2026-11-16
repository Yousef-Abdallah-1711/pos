# Landing image inventory

The detailed inventory is in [`assets.csv`](assets.csv). It records the file
format, dimensions, byte size, SHA-256, classification, and available
provenance for all 39 raster assets under `frontend/apps/landing/public/images`.

## Findings

- 35 files are product or authenticated application screenshots. The 28-image
  `servio/matrix` covers seven dashboard/storefront surfaces in English/Arabic
  and light/dark themes. Its provenance identifies source release
  `e192f30df18a85d2aeb65d554199f27e1ff7fc13`, synthetic data, and locator-based
  redactions. These screenshots are genuine captures; the redaction masks are
  visibly black and can obscure content. Re-capture and review them after the
  final approved release is deployed.
- The roles/permissions pair is documented as authenticated synthetic capture
  evidence from 11 November 2026.
- The access-control image is a crop of a supplied real UI capture. Its
  provenance says the original full-size screenshot is not retained in this
  worktree, so it cannot be reconstructed from this repository alone.
- One café-interior photo is now used in a separate, responsive home-page
  storytelling section. Its source page, photographer, license reference,
  downloaded dimensions, hash, crop, and intended use are recorded in
  [`cafe-interior.provenance.md`](../../frontend/apps/landing/public/images/marketing/cafe-interior.provenance.md).

The image files are not regenerated or altered by this inventory.
