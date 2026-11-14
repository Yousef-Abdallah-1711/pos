// Filesystem routes must match Next.js URLs, including invisible route groups.
const fs = require('node:fs');
const path = require('node:path');

function discover(directory) {
 return fs.readdirSync(directory, { withFileTypes: true }).flatMap(entry =>
  entry.isDirectory() ? discover(path.join(directory, entry.name)) :
   entry.name === 'page.tsx' ? [path.join(directory, entry.name)] : []);
}

function urlPath(relativeFile) {
 return '/' + relativeFile.split(path.sep)
  .filter(segment => segment !== 'page.tsx' && !/^\([^)]*\)$/.test(segment)).join('/');
}

function discoverRoutes(root, app) {
 const base = path.join(root, 'frontend/apps', app, app === 'restaurant-website' ? 'app' : 'src/app');
 const raw = discover(base).map(file => urlPath(path.relative(base, file)));
 // Pinned Node 24 strips the content module's types; reuse the application's
 // canonical feature data rather than maintaining a second list of test slugs.
 const slugs = app === 'landing' ? require(path.join(root, 'frontend/apps/landing/src/content/features.ts'))
  .getFeatures('en').map(feature => feature.slug) : [];
 const expanded = raw.flatMap(route => app === 'landing' && route === '/[locale]/features/[slug]'
  ? slugs.map(slug => route.replace('[slug]', slug)) : [route]);
 return [...new Set(expanded)].sort();
}

module.exports = { discoverRoutes, urlPath };
