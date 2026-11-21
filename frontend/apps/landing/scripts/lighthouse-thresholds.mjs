export function evaluateLighthouseThresholds(results) {
 return results.map(({ locale, performance, accessibility, lcpMs, cls }) => ({
  locale,
  performance,
  performancePass: Number.isFinite(performance) && performance >= .90,
  accessibility,
  accessibilityPass: Number.isFinite(accessibility) && accessibility >= .95,
  lcpMs,
  cls,
  layoutStabilityPass: Number.isFinite(cls) && cls < .1,
 }))
}

export function failedLighthouseLocales(results) {
 return results.filter(result => !result.performancePass || !result.accessibilityPass || !result.layoutStabilityPass)
  .map(({ locale, performance, accessibility, cls }) => ({ locale, performance, accessibility, cls }))
}
