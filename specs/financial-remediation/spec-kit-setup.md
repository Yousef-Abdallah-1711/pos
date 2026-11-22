# Spec Kit setup and limitations

On 2026-11-22, Specify CLI 0.15.2 was initialized in a new, empty
temporary directory with `specify init <new-directory> --integration codex
--script ps --ignore-agent-tools`. Its generated `.specify` and `.agents`
directories were copied into this branch after checking that neither existed.
The generated constitution template was replaced with the financial recovery
principles. No existing repository files were overwritten.

`specify version` reported 0.15.2 and `specify check` reported ready. These
commands check installation and available tools, not specification quality.
The generated scripts use PowerShell; another device needs PowerShell to run
them. To resolve an existing issue without changing persistent feature state:

```powershell
$env:SPECIFY_FEATURE_DIRECTORY = (Resolve-Path 'specs/financial-remediation/issues/CFR-010').Path
& .\.specify\scripts\powershell\check-prerequisites.ps1 -PathsOnly -Json
Remove-Item Env:SPECIFY_FEATURE_DIRECTORY
```

This path-resolution check succeeded for CFR-010. It does not validate issue
content, acceptance, or cross-artifact consistency. The issue documents remain
recovery drafts; a full Spec Kit analysis, substantive issue-by-issue review,
and clean-clone verification remain outstanding.

The original tracker snapshot has Git blob ID
`a39e23ce5413c9fa597ad1329439243d14769f84`, identical to the tracker in
commit `fa2e301a53a34bb85daa3d81012510c56ae9a819`. Its 21 original issue
headings are present, but this comparison alone does not prove requirements
traceability or current engineering status.

## Fresh clone observation

At commit `183572bc40d15629f6109e1f49fa45f75d0319bd`, a normal clone
into a long Windows temporary path failed checkout on a long tracked filename.
A second fresh clone into `C:\w\fr117-clone` succeeded with
`git -c core.longpaths=true clone --filter=blob:none --single-branch --branch
codex/financial-recovery-specs-20261122 <repository-url> <short-path>`.
Its working tree was clean; all 21 issue directories contained each of
`spec.md`, `plan.md`, and `tasks.md`, and the path-resolution command above
worked when run from that clone. File presence is not completeness. Evidence
links, all remote branches, substantive acceptance criteria, handoff
sufficiency, and cross-artifact consistency were not verified by this check.
