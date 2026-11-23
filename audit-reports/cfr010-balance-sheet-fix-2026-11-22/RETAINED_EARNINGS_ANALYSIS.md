# Retained Earnings Analysis

The disposable PostgreSQL integration test posts an explicit synthetic close: debit revenue 500.00, credit expense 300.00, and credit retained earnings 200.00. After that explicit close, assets and equity are both 1,150.00; retained earnings are 200.00 and the current-year earnings test account is 0.00. The report reconciles because a journal exists in the fixture.

This demonstrates report arithmetic after an explicit close only. It is not permission to create closing entries or to infer retained earnings from unclosed activity. No production close or retained-earnings mutation was performed.

**Decision:** CLOSE-01 remains PENDING; retained-earnings presentation is not approved for the open source case.
