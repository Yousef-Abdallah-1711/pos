'use strict';

const TEST_HOSTS = Object.freeze({
  landing: 'https://pos.servio.test',
  dashboard: 'https://pos.app.servio.test',
  admin: 'https://admin.pos.servio.test',
  storefront: 'https://restaurant-1.pos.servio.test',
  api: 'https://api.servio.test',
});

function requireCI() {
  if (process.env.SERVIO_ACCEPTANCE_ENV !== 'ci') {
    throw new Error('Mutating acceptance is allowed only in disposable CI');
  }
}

function requireTarget(name, value, expectedPath = '/') {
  requireCI();
  const expected = new URL(TEST_HOSTS[name]);
  if (!value) throw new Error(`Explicit disposable ${name} URL is required`);
  const actual = new URL(value);
  if (actual.protocol !== 'https:' || actual.host !== expected.host ||
      actual.username || actual.password || actual.search || actual.hash ||
      actual.pathname.replace(/\/$/, '') !== expectedPath.replace(/\/$/, '')) {
    throw new Error(`Unsafe ${name} acceptance target`);
  }
  return value.replace(/\/$/, '');
}

function apiBase() {
  return requireTarget('api', process.env.SERVIO_ACCEPTANCE_API_BASE_URL, '/api/v1');
}

function uiOrigin(name) {
  const envName = `SERVIO_ACCEPTANCE_${name.toUpperCase()}_URL`;
  return requireTarget(name, process.env[envName]);
}

function assertCompleteTargetConfig() {
  apiBase();
  for (const name of ['landing', 'dashboard', 'admin', 'storefront']) uiOrigin(name);
}

function apiRequestOptions(options = {}) {
  return { ...options, maxRedirects: 0 };
}

async function guardContext(context) {
  assertCompleteTargetConfig();
  await context.route('**/*', route => {
    const url = new URL(route.request().url());
    const isDataResource = ['data:', 'blob:'].includes(url.protocol);
    const isApprovedService = Object.values(TEST_HOSTS).some(origin => new URL(origin).host === url.host);
    if ((url.protocol === 'https:' && isApprovedService) || isDataResource) return route.continue();
    return route.abort('blockedbyclient');
  });
}

module.exports = { TEST_HOSTS, apiBase, apiRequestOptions, assertCompleteTargetConfig, guardContext, requireTarget, uiOrigin };
