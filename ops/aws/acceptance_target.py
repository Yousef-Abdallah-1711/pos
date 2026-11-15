"""Fail-closed target validation for mutating acceptance helpers."""
import os
import urllib.request
from urllib.parse import urlsplit

TEST_HOSTS = frozenset({
    'api.servio.test',
    'pos.servio.test',
    'pos.app.servio.test',
    'admin.pos.servio.test',
    'restaurant-1.pos.servio.test',
})


def require_acceptance_url(value, expected_host=None, expected_path=None):
    if os.environ.get('SERVIO_ACCEPTANCE_ENV') != 'ci':
        raise RuntimeError('Mutating acceptance is allowed only in disposable CI')
    if not value:
        raise RuntimeError('Explicit disposable CI target is required')
    parsed = urlsplit(value)
    if (parsed.scheme != 'https' or parsed.hostname not in TEST_HOSTS
            or parsed.username or parsed.password or parsed.port not in (None, 443)
            or parsed.query or parsed.fragment):
        raise RuntimeError('Acceptance target must use an approved servio.test host')
    if expected_host and parsed.hostname != expected_host:
        raise RuntimeError('Acceptance target host does not match the expected service')
    if expected_path and parsed.path.rstrip('/') != expected_path.rstrip('/'):
        raise RuntimeError('Acceptance target path does not match the expected service')
    return value.rstrip('/')


def api_base():
    return require_acceptance_url(
        os.environ.get('SERVIO_ACCEPTANCE_API_BASE_URL'),
        expected_host='api.servio.test',
        expected_path='/api/v1',
    )


def ui_origin(name):
    env_name = 'SERVIO_ACCEPTANCE_' + name.upper().replace('-', '_') + '_URL'
    expected_host = {
        'landing': 'pos.servio.test',
        'dashboard': 'pos.app.servio.test',
        'admin': 'admin.pos.servio.test',
        'storefront': 'restaurant-1.pos.servio.test',
    }[name]
    return require_acceptance_url(os.environ.get(env_name), expected_host, '/')


class RejectRedirects(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, request, file_pointer, code, message, headers, new_url):
        raise RuntimeError('Acceptance requests cannot follow redirects')


def safe_opener(*handlers):
    return urllib.request.build_opener(RejectRedirects(), *handlers)
