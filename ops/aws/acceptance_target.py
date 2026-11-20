"""Fail-closed target validation for mutating acceptance helpers."""
import json
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
TEST_UI_HOSTS = {
    'landing': 'pos.servio.test',
    'dashboard': 'pos.app.servio.test',
    'admin': 'admin.pos.servio.test',
    'storefront': 'restaurant-1.pos.servio.test',
}
STAGING_API_RELEASE_URL = 'https://api.pos.rakmyat.online/_ops/release'
STAGING_HOSTS = {
    'api': 'api.pos.rakmyat.online',
    'landing': 'pos.rakmyat.online',
    'dashboard': 'pos.app.rakmyat.online',
    'admin': 'admin.pos.rakmyat.online',
    'storefront': 'restaurant-1.pos.rakmyat.online',
}


def _verify_live_staging():
    request = urllib.request.Request(
        STAGING_API_RELEASE_URL,
        headers={'Accept': 'application/json', 'User-Agent': 'servio-staging-acceptance'},
    )
    try:
        with safe_opener().open(request, timeout=10) as response:
            if response.status != 200:
                raise RuntimeError('Staging release identity endpoint did not return HTTP 200')
            if response.headers.get('X-Servio-Environment', '').lower() != 'staging':
                raise RuntimeError('Acceptance API did not identify itself as staging')
            release = json.load(response)
    except RuntimeError:
        raise
    except Exception as error:
        raise RuntimeError('Could not verify the live staging acceptance target') from error
    if (not isinstance(release, dict) or release.get('environment') != 'staging'
            or release.get('status') not in {'verified', 'candidate'}
            or not isinstance(release.get('sha'), str)
            or len(release['sha']) != 40
            or any(character not in '0123456789abcdef' for character in release['sha'])):
        raise RuntimeError('Acceptance API release receipt is not for staging')


def require_acceptance_url(value, expected_host=None, expected_path=None):
    mode = os.environ.get('SERVIO_ACCEPTANCE_ENV')
    if mode not in {'ci', 'staging'}:
        raise RuntimeError('Mutating acceptance requires disposable CI or verified staging')
    if not value:
        raise RuntimeError('Explicit acceptance target is required')
    try:
        parsed = urlsplit(value)
        port = parsed.port
    except ValueError as error:
        raise RuntimeError('Acceptance target URL is invalid') from error
    allowed_hosts = TEST_HOSTS if mode == 'ci' else frozenset(STAGING_HOSTS.values())
    if (parsed.scheme != 'https' or parsed.hostname not in allowed_hosts
            or parsed.username or parsed.password or port not in (None, 443)
            or parsed.query or parsed.fragment):
        raise RuntimeError('Acceptance target must use an approved HTTPS host')
    if expected_host and parsed.hostname != expected_host:
        raise RuntimeError('Acceptance target host does not match the expected service')
    if expected_path and parsed.path.rstrip('/') != expected_path.rstrip('/'):
        raise RuntimeError('Acceptance target path does not match the expected service')
    if mode == 'staging':
        _verify_live_staging()
    return value.rstrip('/')


def api_base():
    mode = os.environ.get('SERVIO_ACCEPTANCE_ENV')
    expected_host = 'api.servio.test' if mode == 'ci' else STAGING_HOSTS['api']
    return require_acceptance_url(
        os.environ.get('SERVIO_ACCEPTANCE_API_BASE_URL'),
        expected_host=expected_host,
        expected_path='/api/v1',
    )


def ui_origin(name):
    env_name = 'SERVIO_ACCEPTANCE_' + name.upper().replace('-', '_') + '_URL'
    hosts = TEST_UI_HOSTS if os.environ.get('SERVIO_ACCEPTANCE_ENV') == 'ci' else STAGING_HOSTS
    return require_acceptance_url(os.environ.get(env_name), hosts[name], '/')


class RejectRedirects(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, request, file_pointer, code, message, headers, new_url):
        raise RuntimeError('Acceptance requests cannot follow redirects')


def safe_opener(*handlers):
    return urllib.request.build_opener(RejectRedirects(), *handlers)
