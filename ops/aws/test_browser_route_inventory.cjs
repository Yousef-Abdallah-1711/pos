const { test } = require('node:test');
const assert = require('node:assert/strict');
const path = require('node:path');
const { discoverRoutes, urlPath } = require('./browser-route-inventory.cjs');
const root = path.resolve(__dirname, '../..');

test('route groups are invisible while locale and resource parameters remain', () => {
 assert.equal(urlPath(path.join('(entry)', 'page.tsx')), '/');
 assert.equal(urlPath(path.join('(entry)', 'auth', 'login', 'page.tsx')), '/auth/login');
 assert.equal(urlPath(path.join('[locale]', '(workspace)', 'orders', '[id]', 'page.tsx')), '/[locale]/orders/[id]');
});

test('landing traversal includes real grouped redirects and all seven feature URLs', () => {
 const routes = discoverRoutes(root, 'landing');
 for (const route of ['/', '/login', '/register', '/auth/login', '/auth/register']) assert.ok(routes.includes(route), route);
 const features = routes.filter(route => route.startsWith('/[locale]/features/'));
 assert.equal(features.length, 7);
 assert.ok(features.includes('/[locale]/features/pos'));
 assert.ok(features.includes('/[locale]/features/kitchen-display'));
 assert.ok(routes.every(route => !route.includes('(') && !route.includes('[slug]')));
 assert.deepEqual(routes, [...new Set(routes)].sort());
});

test('both languages expose the same feature routes', () => {
 const { getFeatures } = require(path.join(root, 'frontend/apps/landing/src/content/features.ts'));
 assert.deepEqual(getFeatures('ar').map(feature => feature.slug).sort(), getFeatures('en').map(feature => feature.slug).sort());
});

test('authenticated dynamic resources remain in the discovered inventory', () => {
 assert.ok(discoverRoutes(root, 'dashboard').includes('/[locale]/dashboard/orders/[id]'));
 assert.ok(discoverRoutes(root, 'super-admin').some(route => route.includes('[id]')));
 assert.ok(discoverRoutes(root, 'restaurant-website').some(route => route.includes('[orderNumber]')));
});
