"""Create source inventories for the dashboard appearance handoff; no product files are changed."""
import csv
import json
import re
import struct
from pathlib import Path

root = Path(__file__).resolve().parents[2]
out = Path(__file__).resolve().parent
dashboard = root / 'frontend/apps/dashboard'

def write(name, fields, rows):
    with (out / name).open('w', encoding='utf-8', newline='') as handle:
        writer = csv.DictWriter(handle, fieldnames=fields)
        writer.writeheader()
        writer.writerows(rows)

sidebar = (dashboard / 'src/components/layout/Sidebar.tsx').read_text(encoding='utf-8')
observation_file = out / 'runtime-route-observations.json'
observations = {(item.get('route_template') or item['path']): item for item in json.loads(observation_file.read_text(encoding='utf-8'))['observations']} if observation_file.exists() else {}
sidebar_routes = {}
for match in re.finditer(r'href:\s*`([^`]+)`([^\n]*)', sidebar):
    module = re.search(r'moduleId:\s*MODULE_IDS\.([A-Z_]+)', match.group(2))
    sidebar_routes[match.group(1).split('?')[0]] = module.group(1) if module else 'See parent group'
routes = []
for page in sorted((dashboard / 'src/app').rglob('page.tsx')):
    source_text = page.read_text(encoding='utf-8', errors='replace')
    imports = re.findall(r"import\s+(?:[^;\n]*?\s+from\s+)?['\"]([^'\"]+)['\"]", source_text)
    components_used = [value for value in imports if '/components/' in value or value.startswith('./components/')]
    styles_used = [value for value in imports if value.endswith(('.css', '.scss'))]
    rel = page.relative_to(dashboard / 'src/app')
    segments = [part for part in rel.parts[:-1] if not (part.startswith('(') and part.endswith(')'))]
    layouts = [parent / 'layout.tsx' for parent in [page.parent, *page.parent.parents] if parent == dashboard / 'src/app' or dashboard / 'src/app' in parent.parents]
    route = '/' + '/'.join(segments)
    sidebar_path = route.replace('/[locale]', '')
    observed = observations.get(route) or observations.get(route.replace('/[locale]', '/en'))
    access = ('Synthetic owner EN: HTTP ' + str(observed.get('status', 'unknown')) + '; heading: ' + str(observed.get('heading') or 'none') + '; see runtime observations' if observed and not observed.get('error') else 'Not verified')
    routes.append({'route': route, 'source': page.relative_to(root).as_posix(),
                   'layouts': '; '.join(p.relative_to(root).as_posix() for p in layouts if p.is_file()),
                   'accessibility': access, 'permissions': 'Inspect route and middleware; sidebar hint: ' + sidebar_routes.get(sidebar_path, 'not matched'),
                   'component_imports': '; '.join(components_used) or 'No directly imported component path detected',
                   'style_imports': '; '.join(styles_used) or 'No direct stylesheet import detected; inspect shared layouts and utilities',
                   'locale_theme_responsive': 'Not individually verified', 'migration_impact': 'Review shared shell and hardcoded styles'})
write('dashboard-routes.csv', ['route', 'source', 'layouts', 'component_imports', 'style_imports', 'accessibility', 'permissions', 'locale_theme_responsive', 'migration_impact'], routes)

colors, typography = [], []
for directory in [dashboard / 'src', dashboard / 'components', dashboard / 'lib']:
    if not directory.exists(): continue
    for file in directory.rglob('*'):
        if file.suffix not in {'.ts', '.tsx', '.css', '.js'} or not file.is_file(): continue
        for line_no, line in enumerate(file.read_text(encoding='utf-8', errors='replace').splitlines(), 1):
            path = file.relative_to(root).as_posix()
            for match in re.finditer(r'#[0-9a-fA-F]{3,8}\b|\b(?:rgba?|hsla?)\([^)]*\)|\b(?:bg|text|border|ring|from|to|via)-(?:slate|gray|zinc|neutral|stone|red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose)-\d+\b|--(?:primary|secondary|accent|background|foreground)[\w-]*', line):
                value = match.group()
                category = ('Candidate brand token' if value.startswith(('--primary', '--secondary', '--accent')) else
                            'Candidate structural token' if value.startswith(('--background', '--foreground')) else
                            'Candidate hardcoded blue; may be semantic' if 'blue-' in value else
                            'Color function; role unknown' if re.match(r'^(?:rgba?|hsla?)\(', value) else
                            'Color literal or utility; role unknown')
                colors.append({'source': path, 'line': line_no, 'value': value, 'classification': category + ' — unreviewed',
                               'proposed_token': 'Needs review', 'risk': 'Needs review'})
            if re.search(r'font-family|font-size|font-weight|line-height|letter-spacing|fontFamily|fontSize|fontWeight|lineHeight|letterSpacing|next/font|\b(?:text-(?:xs|sm|base|lg|xl|\d+xl)|font-(?:normal|medium|semibold|bold)|leading-[\w\[\]./-]+|tracking-[\w\[\]./-]+)\b|\btext-\[(?:\d|clamp|var)', line):
                typography.append({'source': path, 'line': line_no, 'excerpt': line.strip()[:240], 'classification': 'Needs review'})
write('color-inventory.csv', ['source', 'line', 'value', 'classification', 'proposed_token', 'risk'], colors)
write('typography-inventory.csv', ['source', 'line', 'excerpt', 'classification'], typography)

components = []
component_dirs = [dashboard / 'src/components', dashboard / 'src/app', dashboard / 'components']
component_files = sorted(component for directory in component_dirs if directory.exists() for component in directory.rglob('*.tsx'))
for component in component_files:
    if '__tests__' in component.parts or component.name.endswith(('.test.tsx', '.spec.tsx')) or component.name in {'page.tsx', 'layout.tsx', 'loading.tsx', 'error.tsx', 'not-found.tsx'}: continue
    source = component.read_text(encoding='utf-8', errors='replace')
    components.append({'source': component.relative_to(root).as_posix(),
                       'scope': ('Shared src/components' if dashboard / 'src/components' in component.parents else 'Route-local src/app' if dashboard / 'src/app' in component.parents else 'Legacy/root components'),
                       'primary_token_mentions': len(re.findall(r'(?:primary-\d+|--primary)', source)),
                       'blue_utility_mentions': len(re.findall(r'\b(?:bg|text|border|ring)-blue-\d+', source)),
                       'dark_utility_mentions': len(re.findall(r'\bdark:', source)),
                       'font_utility_mentions': len(re.findall(r'\bfont-(?:normal|medium|semibold|bold)', source)),
                       'review_status': 'Unreviewed; inspect actual role and accessibility'})
write('component-appearance-inventory.csv', ['source', 'scope', 'primary_token_mentions', 'blue_utility_mentions', 'dark_utility_mentions', 'font_utility_mentions', 'review_status'], components)

impact = [{'route': row['route'], 'source': row['source'], 'current_behavior': 'Not individually verified',
           'proposed_behavior': 'Inherit validated tenant workspace tokens where applicable',
           'dependencies': row['layouts'], 'risk': 'Needs review', 'tests': 'Locale, theme, keyboard, permissions and workflow',
           'capture_impact': 'Determine active marketing usage'} for row in routes]
write('migration-impact-map.csv', ['route', 'source', 'current_behavior', 'proposed_behavior', 'dependencies', 'risk', 'tests', 'capture_impact'], impact)

migration_rows = []
for migration in sorted((root / 'backend/migrations').glob('*.sql')):
    sql = migration.read_text(encoding='utf-8', errors='replace')
    terms = sorted({match.group().lower() for match in re.finditer(r'\b(?:themes_v2|themes|user_settings|tenant_settings|restaurant_settings|branding|theme_history)\b', sql, re.I)})
    if terms:
        migration_rows.append({'source': migration.relative_to(root).as_posix(), 'matched_terms': '; '.join(terms),
                               'review_status': 'Candidate only; inspect statements, order, policies and active usage'})
write('appearance-migration-inventory.csv', ['source', 'matched_terms', 'review_status'], migration_rows)

api_source = root / 'backend/cmd/api/main.go'
api_rows = []
for line_no, line in enumerate(api_source.read_text(encoding='utf-8').splitlines(), 1):
    match = re.search(r'mux\.Handle(?:Func)?\("(GET|POST|PUT|PATCH|DELETE)\s+([^\"]+)"', line)
    if match and re.search(r'theme|setting|appearance|brand', match.group(2), re.I):
        api_rows.append({'method': match.group(1), 'path': match.group(2), 'source': api_source.relative_to(root).as_posix(),
                         'line': line_no, 'registration': line.strip(), 'review_status': 'Candidate; trace active handler, authorization and data scope'})
write('appearance-api-inventory.csv', ['method', 'path', 'source', 'line', 'registration', 'review_status'], api_rows)

matrix = root / 'frontend/apps/landing/public/images/product/servio'
metadata = json.loads((matrix / 'screenshot-matrix.json').read_text(encoding='utf-8'))
captures = {item['file']: item for item in metadata['captures']}
feature_items = json.loads((root / 'frontend/apps/landing/src/i18n/messages/en/features.json').read_text(encoding='utf-8'))
feature_routes = {item['image']: '/{locale}/features/' + item['slug'] for item in feature_items}
homepage_usage = {
    'dashboard': 'hero; operational bento; order story; management',
    'pos': 'operational bento; layered typography',
    'kitchen': 'operational bento; order story; kitchen workflow',
    'checkout': 'order story; online workflow',
    'inventory': 'management',
    'staff': 'management',
    'theme': 'feature showcase',
}
assets = []
product_assets = root / 'frontend/apps/landing/public/images/product'
external_capture = root / 'frontend/apps/landing/public/images/poos-access-control-capture.png'
for image in sorted([*product_assets.rglob('*'), external_capture]):
    if not image.is_file() or image.suffix.lower() not in {'.png', '.jpg', '.jpeg', '.webp', '.avif'}: continue
    rel = image.relative_to(matrix).as_posix() if matrix in image.parents else ''
    record = captures.get(rel, {})
    width, height = record.get('width', 'Unknown'), record.get('height', 'Unknown')
    if image.suffix.lower() == '.png':
        with image.open('rb') as source:
            header = source.read(24)
        if header[:8] == b'\x89PNG\r\n\x1a\n':
            width, height = struct.unpack('>II', header[16:24])
    name = record.get('name', '')
    usage = ('Homepage: feature showcase; ' + homepage_usage.get(name, 'other usage unverified') + '; feature index and route: ' + feature_routes.get(name, 'unverified') if record else 'Historical or fallback; active route usage not verified')
    assets.append({'path': image.relative_to(root).as_posix(), 'absolute_path': str(image), 'bytes': image.stat().st_size,
                   'feature': record.get('name', 'Unspecified'), 'locale': record.get('locale', 'Unknown'),
                   'theme': record.get('theme', 'Unknown'), 'width': width,
                   'height': height, 'aspect_ratio': round(width / height, 4) if isinstance(width, int) and isinstance(height, int) and height else 'Unknown',
                   'source_app': ('Storefront checkout' if record.get('name') == 'checkout' else 'Dashboard' if record else 'Unknown'),
                   'source_route': 'Not recorded in matrix; verify', 'capture_method': 'Playwright, synthetic/redacted per provenance' if record else 'Unknown',
                   'source_release': record.get('sourceRelease', 'Unknown'), 'legacy_branding': 'Review visual',
                   'usage': usage, 'rendering_component': 'ProductCapture.tsx → CaptureMedia.tsx' if record else 'Needs verification',
                   'animation_geometry': 'Review hero, bento, story and showcase frame and scroll geometry' if record else 'Needs verification',
                   'mobile_crop': 'Review responsive frame; dedicated crop not yet established', 'regeneration': 'Review after rebrand'})
write('marketing-screenshot-inventory.csv', ['path', 'absolute_path', 'bytes', 'feature', 'locale', 'theme', 'width', 'height', 'aspect_ratio', 'source_app', 'source_route', 'capture_method', 'source_release', 'legacy_branding', 'usage', 'rendering_component', 'animation_geometry', 'mobile_crop', 'regeneration'], assets)
print(f'routes={len(routes)} components={len(components)} color_occurrences={len(colors)} typography_lines={len(typography)} product_image_assets={len(assets)} candidate_migrations={len(migration_rows)} candidate_api_routes={len(api_rows)}')
