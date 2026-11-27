# Evidence index

- Landing baseline and browser evidence: `docs/landing-interactive-experience/artifacts/` and `operational-bento-showcase-report.md`.
- Existing product screenshot provenance: `frontend/apps/landing/public/images/product/servio/{screenshot-matrix.json,provenance.json}`.
- Authenticated dashboard baseline: [26 synthetic-user captures](baseline/README.md) from a separately created disposable test stack, including one synthetic POS cart and one staff user without role assignments. Earlier, an isolated browser session requesting the locally exposed dashboard at port 13112 redirected to the landing login; that pre-existing stack was not modified or seeded. Existing marketing images are not proof of the newly captured dashboard state. Additional restricted-role and populated synthetic-tenant scenarios remain outstanding.
- Source inventory method: `generate-inventories.py`; it reports source matches rather than validated accessibility, permissions, design roles, or runtime behavior.
