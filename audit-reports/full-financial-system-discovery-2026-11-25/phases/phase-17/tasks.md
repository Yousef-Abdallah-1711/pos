# Phase 17 — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

- **T-P17-001** [EXECUTED-PARTIAL] Remaining security tests: SEC-013 (staging env copy), SEC-014 stub login, CSRF/CORS/cookies/rate limit — module security — expected: per 35/35a oracle or documented rule
- **T-P17-002** [EXECUTED-PARTIAL] RBAC matrix classified as observation unless a permission matrix is approved — module security — expected: per 35/35a oracle or documented rule

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.
