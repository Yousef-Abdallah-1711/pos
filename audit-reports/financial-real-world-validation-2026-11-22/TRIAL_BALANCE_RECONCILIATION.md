# Trial Balance reconciliation

For the authenticated synthetic order, independently expected total debits and credits were each 5,500 cents. The API/database acceptance observed 5,500 cents for each; difference was zero. The source order total was also 5,500 cents. The accountant journal scenario independently balanced at 10,000 cents per side.

The 1,000-order performance seed reported 11,032,700 cents of journal debits and credits and the same totals in GL. Its benchmark JSON also records 5,229 journal lines/GL rows. The benchmark runner log reports a separate seed invariant at 11,000,000 cents for both journal sides. These logs are scenario-specific; the totals must not be conflated.

These PASS results cover the listed synthetic scenarios, not every fixture in `servio_test`, and do not establish original operating-data reconciliation. See `JOURNAL_GL_RECONCILIATION.md` and `financial-benchmark-results-1000.json`.
