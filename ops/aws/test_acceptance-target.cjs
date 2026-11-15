'use strict';

const { test } = require('node:test');
const assert = require('node:assert/strict');
const {
  apiBase,
  apiRequestOptions,
  assertCompleteTargetConfig,
  guardContext,
  requireTarget,
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
