#!/usr/bin/env python3
"""Staging API acceptance using pre-provisioned synthetic accounts and session-only writes."""
import http.cookiejar
import json
import os
from pathlib import Path
import re
import sys
import time
import urllib.error
import urllib.request

API_ORIGIN = 'https://api.pos.rakmyat.online'
API_BASE = API_ORIGIN + '/api/v1'
LANDING_ORIGIN = 'https://pos.rakmyat.online'
ADMIN_ORIGIN = 'https://admin.pos.rakmyat.online'
STORE_ORIGIN = 'https://restaurant-1.pos.rakmyat.online'
EXPECTED_SHA = os.environ.get('SERVIO_RELEASE_SHA', '')
RECEIPT_STATUS = os.environ.get('SERVIO_STAGING_RECEIPT_STATUS', 'candidate')
FIXTURE_DIR = Path(os.environ.get('SERVIO_STAGING_FIXTURE_DIR', '/opt/servio/ops/aws/evidence'))
EVIDENCE = Path(os.environ.get('SERVIO_STAGING_EVIDENCE', '/tmp/servio-staging-safe-acceptance.json'))
CHECKS = []
SESSIONS = []


class RejectRedirects(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, request, file_pointer, code, message, headers, new_url):
        raise RuntimeError('Acceptance requests cannot follow redirects')


def client():
    jar = http.cookiejar.CookieJar()
    opener = urllib.request.build_opener(RejectRedirects(), urllib.request.HTTPCookieProcessor(jar))
    return opener, jar


def request(opener, method, url, payload=None, origin=None, headers=None):
    values = {'Accept': 'application/json'}
    if payload is not None:
        values['Content-Type'] = 'application/json'
    if origin:
        values['Origin'] = origin
    values.update(headers or {})
    body = None if payload is None else json.dumps(payload).encode()
    req = urllib.request.Request(url, data=body, headers=values, method=method)
    try:
        response = opener.open(req, timeout=20)
    except urllib.error.HTTPError as error:
        response = error
    raw = response.read(1024 * 1024)
    try:
        data = json.loads(raw)
    except (ValueError, UnicodeDecodeError):
        data = {}
    return response.status, data, response.headers


def expect(name, passed, detail=None):
    CHECKS.append({'check': name, 'status': 'PASS' if passed else 'FAIL', **(detail or {})})
    print(('PASS ' if passed else 'FAIL ') + name, flush=True)
    if not passed:
        raise RuntimeError(name)


def read_fixture(name):
    path = FIXTURE_DIR / name
    if path.is_symlink() or not path.is_file() or path.stat().st_mode & 0o077:
        raise RuntimeError('Protected synthetic acceptance fixture is unavailable')
    if path.stat().st_size > 65536:
        raise RuntimeError('Synthetic acceptance fixture exceeds its size bound')
    return json.loads(path.read_text())


def verify_receipt(expected_status):
    opener, _ = client()
    status, receipt, _ = request(opener, 'GET', API_ORIGIN + '/_ops/release')
    expect('Release receipt matches expected SHA before acceptance',
           status == 200 and receipt.get('sha') == EXPECTED_SHA
           and receipt.get('status') == expected_status
           and receipt.get('environment') == 'staging')


def login(credentials, endpoint='/auth/login', origin=LANDING_ORIGIN):
    opener, jar = client()
    prefix = endpoint.rsplit('/auth/login', 1)[0]
    SESSIONS.append((opener, prefix, origin))
    status, _, _ = request(opener, 'POST', API_BASE + endpoint,
                           {'email': credentials['email'], 'password': credentials['password']}, origin)
    cookie_name = 'poos_super_admin_session' if endpoint.startswith('/super-admin/') else 'poos_session'
    cookies = [cookie for cookie in jar if cookie.name == cookie_name]
    expect('Approved synthetic account login succeeds', status == 200 and len(cookies) == 1)
    cookie = cookies[0]
    expect('Session cookie is HttpOnly, Secure, host-only, and session-only',
           cookie.secure and (cookie.has_nonstandard_attr('HttpOnly') or cookie.has_nonstandard_attr('httponly'))
           and not cookie.domain_specified and cookie.expires is None)
    return opener


def session(opener, prefix=''):
    return request(opener, 'GET', API_BASE + prefix + '/auth/session')


def check_staff_boundaries(staff):
    first = login(staff[0])
    second = login(staff[1])
    first_status, first_body, _ = session(first)
    second_status, second_body, _ = session(second)
    first_user = first_body.get('data', first_body)
    second_user = second_body.get('data', second_body)
    expect('Synthetic staff sessions authenticate independently', first_status == 200 and second_status == 200)
    expect('Synthetic staff identities belong to different tenants',
           first_user.get('tenant_id') != second_user.get('tenant_id')
           and first_user.get('user_id') != second_user.get('user_id'))
    tenant_headers = {'X-Tenant-ID': str(first_user.get('tenant_id'))}
    spoofed_status, spoofed_body, _ = request(second, 'GET', API_BASE + '/auth/session', headers=tenant_headers)
    spoofed_user = spoofed_body.get('data', spoofed_body)
    expect('Forged tenant header cannot switch the authenticated identity',
           spoofed_status == 200 and spoofed_user.get('tenant_id') == second_user.get('tenant_id'))
    profile_status, _, _ = request(second, 'GET', API_BASE + '/users/' + str(first_user.get('user_id')) + '/profile',
                                   headers=tenant_headers)
    expect('Cross-tenant profile access is denied', profile_status in (403, 404))
    for headers in ({'Origin': 'https://attacker.invalid'}, {}):
        csrf_status, _, _ = request(first, 'POST', API_BASE + '/auth/logout', {}, headers=headers)
        expect('Cookie-authenticated logout rejects untrusted or missing Origin', csrf_status == 403)
    still_status, _, _ = session(first)
    expect('Rejected CSRF requests preserve the synthetic session', still_status == 200)
    return first, second


def check_admin_isolation(admin):
    staff_opener, _, _ = SESSIONS[0]
    staff_admin_status, _, _ = session(staff_opener, '/super-admin')
    expect('Staff cookie cannot authenticate Super Admin APIs', staff_admin_status == 401)
    admin_opener = login(admin, '/super-admin/auth/login', ADMIN_ORIGIN)
    admin_status, _, _ = session(admin_opener, '/super-admin')
    staff_status, _, _ = session(admin_opener)
    expect('Super Admin session is isolated from staff session', admin_status == 200 and staff_status == 401)
    return admin_opener


def check_storefront(store_fixture):
    opener, _ = client()
    headers = {
        'X-Tenant-ID': str(store_fixture['tenant_id']),
        'X-Restaurant-ID': str(store_fixture['restaurant_id']),
    }
    status, body, response_headers = request(
        opener, 'GET', API_BASE + '/public/restaurants/restaurant-1/menu',
        origin=STORE_ORIGIN, headers=headers)
    cache = response_headers.get('Cache-Control', '')
    expect('Approved synthetic storefront catalog is readable anonymously',
           status == 200 and isinstance(body, (dict, list)))
    expect('Public storefront menu declares only a public cache policy',
           'public' in cache.lower() and 'private' not in cache.lower() and 'no-store' not in cache.lower())


def cleanup_sessions():
    failed = False
    for index, (opener, prefix, origin) in enumerate(SESSIONS, start=1):
        try:
            status, _, _ = request(opener, 'POST', API_BASE + prefix + '/auth/logout', {}, origin)
            current, _, _ = session(opener, prefix)
            revoked = status in (200, 204, 401) and current == 401
        except Exception:
            revoked = False
        name = 'Synthetic acceptance session ' + str(index) + ' is revoked and cannot be replayed'
        CHECKS.append({'check': name, 'status': 'PASS' if revoked else 'FAIL'})
        print(('PASS ' if revoked else 'FAIL ') + name, flush=True)
        failed = failed or not revoked
    return not failed


def run():
    if not re.fullmatch(r'[0-9a-f]{40}', EXPECTED_SHA):
        raise RuntimeError('A full expected release SHA is required')
    if RECEIPT_STATUS not in ('candidate', 'verified'):
        raise RuntimeError('Unsupported staging receipt state')
    verify_receipt(RECEIPT_STATUS)
    staff = read_fixture('browser-fixtures.json')
    admin = read_fixture('browser-admin.json')
    store = read_fixture('storefront-qa-fixture.json')
    if (not isinstance(staff, list) or len(staff) < 2
            or any(not item.get('email', '').endswith('@example.invalid') or not item.get('password') for item in staff[:2])
            or not admin.get('email', '').endswith('@example.invalid') or not admin.get('password')
            or not store.get('synthetic') or store.get('tenant_id') != 1 or store.get('restaurant_id') != 1):
        raise RuntimeError('Only provisioned synthetic staging fixtures are accepted')
    check_staff_boundaries(staff)
    check_admin_isolation(admin)
    check_storefront(store)
    verify_receipt(RECEIPT_STATUS)


def write_evidence(started, error=None):
    EVIDENCE.parent.mkdir(mode=0o700, parents=True, exist_ok=True)
    payload = {
        'started_utc': started,
        'completed_utc': time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime()),
        'expected_sha': EXPECTED_SHA,
        'receipt_status': RECEIPT_STATUS,
        'checks': CHECKS,
        'error_type': type(error).__name__ if error else None,
    }
    EVIDENCE.write_text(json.dumps(payload, indent=2) + '\n')
    EVIDENCE.chmod(0o600)


if __name__ == '__main__':
    started_utc = time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())
    failure = None
    try:
        run()
    except Exception as error:
        failure = error
        print('Staging-safe acceptance stopped: ' + type(error).__name__, file=sys.stderr)
        if not CHECKS or CHECKS[-1]['status'] == 'PASS':
            CHECKS.append({'check': 'Acceptance completed without an error', 'status': 'FAIL'})
    finally:
        if not cleanup_sessions():
            failure = failure or RuntimeError('Synthetic acceptance session cleanup failed')
        try:
            write_evidence(started_utc, failure)
        except OSError:
            print('Could not persist protected staging acceptance evidence', file=sys.stderr)
            failure = failure or RuntimeError('Evidence write failed')
    if failure or any(row['status'] != 'PASS' for row in CHECKS):
        raise SystemExit(1)
