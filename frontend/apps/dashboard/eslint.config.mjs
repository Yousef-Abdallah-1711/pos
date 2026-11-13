import { defineConfig, globalIgnores } from 'eslint/config';
import nextVitals from 'eslint-config-next/core-web-vitals';
export default defineConfig([
  ...nextVitals,
  // React Compiler is not enabled. Preserve the enforced semantic Hooks rules;
  // report the newly introduced compiler-compatibility diagnostics as warnings.
  { rules: {
    'react-hooks/set-state-in-effect': 'warn',
    'react-hooks/preserve-manual-memoization': 'warn',
    'react-hooks/error-boundaries': 'warn',
  } },
  globalIgnores(['.next/**', 'node_modules/**', 'next-env.d.ts']),
]);
