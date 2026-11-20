'use strict';

const { test } = require('node:test');
const assert = require('node:assert/strict');
const childProcess = require('node:child_process');
const {
  STAGING_HOSTS,
  apiBase,
  apiRequestOptions,
  assertCompleteTargetConfig,
  guardContext,
  requireTarget,
  uiOrigin,
} = require('./acceptance-target.cjs');

const targetEnvironment = {
  SERVIO_ACCEPTANCE_ENV: 'ci',
  SERVIO_ACCEPTANCE_API_BASE_URL: 'https://api.servio.test/api/v1',
  SERVIO_ACCEPTANCE_LANDING_URL: 'https://pos.servio.test',
  SERVIO_ACCEPTANCE_DASHBOARD_URL: 'https://pos.app.servio.test',
  SERVIO_ACCEPTANCE_ADMIN_URL: 'https://admin.pos.servio.test',
  SERVIO_ACCEPTANCE_STOREFRONT_URL: 'https://restaurant-1.pos.servio.test',
};

test('acceptance targets require explicit disposable CI configuration', () => {
  const original = { ...process.env };
  try {
    Object.assign(process.env, targetEnvironment);
    assertCompleteTargetConfig();
    assert.equal(apiBase(), targetEnvironment.SERVIO_ACCEPTANCE_API_BASE_URL);
  } finally {
    for (const key of Object.keys(process.env)) {
      if (!(key in original)) delete process.env[key];
    }
    Object.assign(process.env, original);
  }
});

test('production, staging, unknown hosts, and unsafe URL forms are rejected', () => {
  const original = process.env.SERVIO_ACCEPTANCE_ENV;
  process.env.SERVIO_ACCEPTANCE_ENV = 'ci';
  const unsafe = [
    'https://api.pos.rakmyat.online/api/v1',
    'https://api.example.invalid/api/v1',
    'http://api.servio.test/api/v1',
    'https://user:pass@api.servio.test/api/v1',
    'https://api.servio.test:444/api/v1',
    'https://api.servio.test/api/v1?redirect=https://example.invalid',
    'https://api.servio.test/api/v1#fragment',
  ];
  try {
    for (const target of unsafe) assert.throws(() => requireTarget('api', target, '/api/v1'));
  } finally {
    if (original === undefined) delete process.env.SERVIO_ACCEPTANCE_ENV;
    else process.env.SERVIO_ACCEPTANCE_ENV = original;
  }
});

test('staging targets require an exact host and a live staging marker', () => {
  const original = { ...process.env };
  const originalSpawnSync = childProcess.spawnSync;
  const stagingEnvironment = {
    SERVIO_ACCEPTANCE_ENV: 'staging',
    SERVIO_ACCEPTANCE_API_BASE_URL: `${STAGING_HOSTS.api}/api/v1`,
    SERVIO_ACCEPTANCE_LANDING_URL: STAGING_HOSTS.landing,
    SERVIO_ACCEPTANCE_DASHBOARD_URL: STAGING_HOSTS.dashboard,
    SERVIO_ACCEPTANCE_ADMIN_URL: STAGING_HOSTS.admin,
    SERVIO_ACCEPTANCE_STOREFRONT_URL: STAGING_HOSTS.storefront,
  };
  try {
    Object.assign(process.env, stagingEnvironment);
    childProcess.spawnSync = () => ({
      status: 0,
      stdout: 'HTTP/2 200\r\nx-servio-environment: staging\r\n\r\n{"environment":"staging","status":"verified","sha":"0000000000000000000000000000000000000000"}',
    });
    assertCompleteTargetConfig();
    assert.equal(apiBase(), stagingEnvironment.SERVIO_ACCEPTANCE_API_BASE_URL);
    assert.equal(uiOrigin('dashboard'), STAGING_HOSTS.dashboard);

    childProcess.spawnSync = () => ({
      status: 0,
      stdout: 'HTTP/2 200\r\nx-servio-environment: production\r\n\r\n{"environment":"staging","status":"verified","sha":"0000000000000000000000000000000000000000"}',
    });
    assert.throws(() => apiBase(), /did not identify itself as staging/);
  } finally {
    childProcess.spawnSync = originalSpawnSync;
    for (const key of Object.keys(process.env)) {
      if (!(key in original)) delete process.env[key];
    }
    Object.assign(process.env, original);
  }
});

test('API request helpers refuse redirect following', () => {
  assert.equal(apiRequestOptions().maxRedirects, 0);
});

test('browser request routing rejects public targets before navigation', async () => {
  const original = { ...process.env };
  const routed = [];
  try {
    Object.assign(process.env, targetEnvironment);
    await guardContext({ route: async (pattern, handler) => routed.push({ pattern, handler }) });
    const route = routed[0];
    assert.equal(route.pattern, '**/*');
    const result = await route.handler({
      request: () => ({ url: () => 'https://api.pos.rakmyat.online/api/v1/orders' }),
      abort: () => 'aborted',
      continue: () => 'continued',
    });
    assert.equal(result, 'aborted');
  } finally {
    for (const key of Object.keys(process.env)) {
      if (!(key in original)) delete process.env[key];
    }
    Object.assign(process.env, original);
  }
});
