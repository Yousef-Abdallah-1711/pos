"""Phase 2A: small, behaviour-preserving dashboard type fixes (each verified by tsc)."""
import re

ROOT = "frontend/apps/dashboard/src/"


def rw(rel, pairs):
    p = ROOT + rel
    s = open(p, encoding="utf-8").read()
    crlf = "\r\n" in s
    s = s.replace("\r\n", "\n")
    for a, b in pairs:
        assert a in s, (rel, a[:70])
        s = s.replace(a, b, 1)
    if crlf:
        s = s.replace("\n", "\r\n")
    open(p, "w", encoding="utf-8", newline="").write(s)


# real bug: hexToHex was never defined (ReferenceError when generating a palette); rgbToHex is the helper
rw("stores/themeBuilderUtils.ts", [
    ("secondary: hexToHex(rgb.r * 0.8, rgb.g * 0.8, rgb.b * 0.8)",
     "secondary: rgbToHex(Math.round(rgb.r * 0.8), Math.round(rgb.g * 0.8), Math.round(rgb.b * 0.8))"),
    ("accent: hexToHex(rgb.r * 1.1, rgb.g * 0.6, rgb.b * 0.6)",
     "accent: rgbToHex(Math.min(255, Math.round(rgb.r * 1.1)), Math.round(rgb.g * 0.6), Math.round(rgb.b * 0.6))"),
])

# AuthState is re-exported from stores/index.ts but was not exported
rw("stores/authStore.ts", [("interface AuthState {", "export interface AuthState {")])

# validators returned `boolean | ""` (a string leaked out of the && chain)
s = open(ROOT + "stores/themeBuilderStore.ts", encoding="utf-8").read()
crlf = "\r\n" in s
s = s.replace("\r\n", "\n")
s = s.replace("""        return (
          typography.fontFamily &&
          typography.fontFamily.trim().length > 0 &&""", """        return (
          Boolean(typography.fontFamily) &&
          typography.fontFamily.trim().length > 0 &&""", 1)
open(ROOT + "stores/themeBuilderStore.ts", "w", encoding="utf-8", newline="").write(s.replace("\n", "\r\n") if crlf else s)

# MetricCard: an unset color looked up the key "gray", which does not exist (className became "undefined")
rw("components/monitoring/PerformanceMonitor.tsx", [
    ("""    red: 'text-red-600',
  }[color || 'gray']""", """    red: 'text-red-600',
  }[color ?? 'green' as const] ?? 'text-gray-600'"""),
])

# untyped call with a type argument -> cast the result instead
rw("hooks/useCache.ts", [
    ("const cached = (cache as any).get<T>(key)", "const cached = (cache as any).get(key) as T | undefined"),
    ("return (cache as any).get<T>(key)", "return (cache as any).get(key) as T | undefined"),
])

# Background Sync is not in lib.dom
BG = "(registration as ServiceWorkerRegistration & { sync?: { register(tag: string): Promise<void> } })"
rw("hooks/usePOSBackgroundSync.ts", [("await registration.sync.register(POS_ORDER_SYNC_TAG)",
    f"await {BG}.sync?.register(POS_ORDER_SYNC_TAG)")])
rw("hooks/useServiceWorker.ts", [("await registration.sync.register(tag)", f"await {BG}.sync?.register(tag)")])

# lib.dom lacks the non-standard LCP/navigation timing fields
rw("lib/web-vitals.ts", [
    ("value: lastEntry.renderTime || lastEntry.loadTime,", "value: lcpTime,"),
    ("rating: getRating('lcp', lastEntry.renderTime || lastEntry.loadTime),", "rating: getRating('lcp', lcpTime),"),
    ("domParsing: navigationEntry.domInteractive - navigationEntry.domLoading,",
     "domParsing: navigationEntry.domInteractive - (navigationEntry as PerformanceNavigationTiming & { domLoading?: number }).domLoading!,"),
])
s = open(ROOT + "lib/web-vitals.ts", encoding="utf-8").read()
crlf = "\r\n" in s
s = s.replace("\r\n", "\n")
s = s.replace("      const metric: WebVitalMetric = {\n        name: 'LCP',",
              "      const lcpEntry = lastEntry as PerformanceEntry & { renderTime?: number; loadTime?: number }\n      const lcpTime = lcpEntry.renderTime || lcpEntry.loadTime || 0\n      const metric: WebVitalMetric = {\n        name: 'LCP',", 1)
open(ROOT + "lib/web-vitals.ts", "w", encoding="utf-8", newline="").write(s.replace("\n", "\r\n") if crlf else s)

# translators take the Locale union, the components were typed `string`
for rel, iface in (("components/accounting/LedgerView.tsx", "LedgerViewProps"),
                   ("components/accounting/TrialBalanceTable.tsx", "TrialBalanceTableProps")):
    rw(rel, [(f"interface {iface} {{", f"interface {iface} {{")])
    s = open(ROOT + rel, encoding="utf-8").read()
    crlf = "\r\n" in s
    s = s.replace("\r\n", "\n")
    s = re.sub(rf"(interface {iface} \{{(?:.|\n)*?)locale: string", r"\1locale: Locale", s, count=1)
    if "type Locale" not in s and "Locale }" not in s and "Locale," not in s:
        lines = s.split("\n")
        idx = max(i for i, l in enumerate(lines) if l.startswith("import "))
        lines.insert(idx + 1, "import type { Locale } from '@/lib/translations'")
        s = "\n".join(lines)
    open(ROOT + rel, "w", encoding="utf-8", newline="").write(s.replace("\n", "\r\n") if crlf else s)

# optional callback
rw("components/hr/EmployeeList.tsx", [("await onDelete(deleteConfirm.id)", "await onDelete?.(deleteConfirm.id)")])
print("ok")
