# Current Earnings Analysis

The disposable journal/report case shows revenue 5,000¢, expenses 0¢, and net income 5,000¢. The Balance Sheet query classifies only posted asset, liability, and equity account types; it does not include that nominal-period result in equity. That explains the observed arithmetic but does not establish whether SERVIO should present current-period earnings in equity.

The integration fixture also tests a separate current period with revenue 100.00, expense 120.00, net loss 20.00, and an open-period Balance Sheet residual of -20.00. This is evidence that the existing report does not roll net loss into equity either.

**Decision:** EARN-01 remains PENDING. No diagnostic is activated as production equity, no balancing entry is posted, and no report status is forced to pass.
