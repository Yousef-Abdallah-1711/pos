import assert from 'node:assert/strict';
import { createRequire } from 'node:module';
import { readFileSync, realpathSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

const workspaceRoot = process.cwd();
const apps = ['apps/dashboard', 'apps/restaurant-website'];

function findNextIntlRoot(entryPath) {
  let directory = dirname(entryPath);
  while (directory !== dirname(directory)) {
    const packagePath = join(directory, 'package.json');
    try {
      const packageJson = JSON.parse(readFileSync(packagePath, 'utf8'));
      if (packageJson.name === 'next-intl') return directory;
    } catch {
      // Continue walking until the next-intl package root is found.
    }
    directory = dirname(directory);
  }
  throw new Error(`Unable to find next-intl package root from ${entryPath}`);
}

for (const app of apps) {
  const appRequire = createRequire(resolve(workspaceRoot, app, 'package.json'));
  const pluginPath = realpathSync(appRequire.resolve('next-intl/plugin'));
  assert.equal(typeof appRequire('next-intl/plugin'), 'function', `${app}: plugin loads`);

  const nextIntlRoot = findNextIntlRoot(pluginPath);
  const esmPluginPath = pathToFileURL(join(nextIntlRoot, 'dist/esm/development/plugin.js'));
  const { default: createNextIntlPlugin } = await import(esmPluginPath);
  assert.equal(typeof createNextIntlPlugin, 'function', `${app}: ESM plugin loads`);
  assert.throws(
    () => appRequire.resolve('@parcel/watcher', { paths: [nextIntlRoot] }),
    { code: 'MODULE_NOT_FOUND' },
    `${app}: patched next-intl no longer resolves @parcel/watcher`,
  );

  const watcherPath = join(nextIntlRoot, 'dist/esm/development/extractor/source/SourceFileWatcher.js');
  const { default: SourceFileWatcher } = await import(pathToFileURL(watcherPath));
  const watcher = new SourceFileWatcher([workspaceRoot], () => {});
  await assert.rejects(watcher.start(), /experimental\.extract file watching was removed/);
}

console.log('PASS: Dashboard and restaurant app package entries load without @parcel/watcher; extraction watcher fails explicitly.');
