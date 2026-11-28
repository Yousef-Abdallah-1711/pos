import os, re, glob
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25\tests')
rep = [
    ('Proposed (root cause still PROBABLE, so step 0 is investigation): 0) trace the session lookup call chain from the auth middleware to backend/internal/service/cache.go and add a Redis-down fault-injection test that records per-call durations;',
     'Proposed (root cause CONFIRMED in RUN-005: middleware/auth.go:60 -> auth_session.go:166 -> cache.go:248 with no deadline, default go-redis options, server WriteTimeout 60 s): 0) add a Redis-down fault-injection test that records per-call durations (ph18_r5_redis_timing.mjs is the reference);'),
    ('Proposed (root cause PROBABLE, so step 0 is replay): 0) replay the supplier fixture transactions recording suppliers.current_balance, invoices_supplier.balance_due and account 2101 after each step to find the first divergence; 1) stop writing current_balance',
     'Proposed (root cause CONFIRMED in RUN-005: the payment path subtracts from current_balance at payment_acct_repo.go:302-306 and no invoice path ever adds to it): 0) keep the replay in ph06_r5_supplier_balance.mjs as the failing test; 1) stop writing current_balance'),
    ('Proposed (root cause PROBABLE; first run a human-paced repro with a network timeline): in the landing register page',
     'Proposed (root cause CONFIRMED in RUN-005 by source: nothing clears the manual slug error): in the landing register page'),
    ('Proposed (root cause PROBABLE: first trace the submit handlers of the four forms): add one shared hook',
     'Proposed (root cause CONFIRMED in RUN-005: parent rethrows, children and the contracts page do not catch): add one shared hook'),
]
n = 0
for f in glob.glob('finding_overrides_r4*.py'):
    s = open(f, encoding='utf-8').read(); o = s
    for a, b in rep:
        s = s.replace(a, b)
    if s != o:
        open(f, 'w', encoding='utf-8').write(s); n += 1
print('files changed', n)
