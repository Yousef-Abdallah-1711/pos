# Evidence index

- Landing baseline and browser evidence: `docs/landing-interactive-experience/artifacts/` and `operational-bento-showcase-report.md`.
- Existing product screenshot provenance: `frontend/apps/landing/public/images/product/servio/{screenshot-matrix.json,provenance.json}`.
- Authenticated dashboard baseline screenshots for this discovery: **not captured**. A fresh isolated browser session requesting the locally exposed dashboard at port 13112 was redirected to the landing login at port 13111; no authenticated sidebar or product view was available in that session. This existing disposable stack may be in use by other work, so it was not modified or seeded. Do not use existing marketing images as proof of current dashboard state. Capture synthetic tenant scenarios in a dedicated environment after verifying the environment and redactions.
- Source inventory method: `generate-inventories.py`; it reports source matches rather than validated accessibility, permissions, design roles, or runtime behavior.
