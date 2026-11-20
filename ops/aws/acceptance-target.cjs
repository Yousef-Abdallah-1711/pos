'use strict';

const childProcess = require('node:child_process');

const TEST_HOSTS = Object.freeze({
  landing: 'https://pos.servio.test',
  dashboard: 'https://pos.app.servio.test',
  admin: 'https://admin.pos.servio.test',
  storefront: 'https://restaurant-1.pos.servio.test',
  api: 'https://api.servio.test',
});
const STAGING_HOSTS = Object.freeze({
  landing: 'https://pos.rakmyat.online',
  dashboard: 'https://pos.app.rakmyat.online',
  admin: 'https://admin.pos.rakmyat.online',
  storefront: 'https://restaurant-1.pos.rakmyat.online',
  api: 'https://api.pos.rakmyat.online',
});
const STAGING_RELEASE_URL = `${STAGING_HOSTS.api}/_ops/release`;

function acceptanceMode() {
  const mode = process.env.SERVIO_ACCEPTANCE_ENV;
  if (mode !== 'ci' && mode !== 'staging') {
    throw new Error('Mutating acceptance requires disposable CI or verified staging');
  }
  return mode;
}

function verifyLiveStaging() {
  const result = childProcess.spawnSync('curl', [
    '--silent', '--show-error', '--max-time', '10', '--noproxy', '*',
    '--dump-header', '-', STAGING_RELEASE_URL,
  ], { encoding: 'utf8', timeout: 12000, maxBuffer: 65536 });
  if (result.error || result.status !== 0 || typeof result.stdout !== 'string') {
    throw new Error('Could not verify the live staging acceptance target');
  }
  const separator = result.stdout.lastIndexOf('\r\n\r\n');
  if (separator < 0) throw new Error('Staging release identity response is malformed');
  const headers = result.stdout.slice(0, separator);
  const body = result.stdout.slice(separator + 4);
  const statusLines = [...headers.matchAll(/^HTTP\/\S+\s+(\d{3})/gim)];
  const status = Number(statusLines.at(-1)?.[1]);
  if (status !== 200 || !/^x-servio-environment:\s*staging\s*$/im.test(headers)) {
    throw new Error('Acceptance API did not identify itself as staging');
  }
  let release;
  try {
    release = JSON.parse(body);
  } catch {
    throw new Error('Staging release identity response is malformed');
  }
  if (!release || release.environment !== 'staging' ||
      !['verified', 'candidate'].includes(release.status) ||
      !/^[0-9a-f]{40}$/.test(release.sha || '')) {
    throw new Error('Acceptance API release receipt is not for staging');
  }
}

function configuredHosts() {
  return acceptanceMode() === 'ci' ? TEST_HOSTS : STAGING_HOSTS;
}

function requireTarget(name, value, expectedPath = '/') {
  const mode = acceptanceMode();
  const expected = new URL((mode === 'ci' ? TEST_HOSTS : STAGING_HOSTS)[name]);
  if (!value) throw new Error(`Explicit ${name} acceptance URL is required`);
  const actual = new URL(value);
  if (actual.protocol !== 'https:' || actual.hostname !== expected.hostname ||
      actual.port && actual.port !== '443' || actual.username || actual.password ||
      actual.search || actual.hash ||
      actual.pathname.replace(/\/$/, '') !== expectedPath.replace(/\/$/, '')) {
    throw new Error(`Unsafe ${name} acceptance target`);
  }
  if (mode === 'staging') verifyLiveStaging();
  return value.replace(/\/$/, '');
}

function apiBase() {
  return requireTarget('api', process.env.SERVIO_ACCEPTANCE_API_BASE_URL, '/api/v1');
}

function uiOrigin(name) {
  const envName = `SERVIO_ACCEPTANCE_${name.toUpperCase().replaceAll('-', '_')}_URL`;
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
  const approvedHosts = new Set(Object.values(configuredHosts()).map(origin => new URL(origin).host));
  await context.route('**/*', route => {
    const url = new URL(route.request().url());
    const isDataResource = ['data:', 'blob:'].includes(url.protocol);
    const isApprovedService = approvedHosts.has(url.host);
    if ((url.protocol === 'https:' && isApprovedService) || isDataResource) return route.continue();
    return route.abort('blockedbyclient');
  });
}

module.exports = {
  STAGING_HOSTS,
  TEST_HOSTS,
  apiBase,
  apiRequestOptions,
  assertCompleteTargetConfig,
  guardContext,
  requireTarget,
  uiOrigin,
};
