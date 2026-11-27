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
    routes.append({'route': route, 'source': page.relative_to(root).as_posix(),
                   'layouts': '; '.join(p.relative_to(root).as_posix() for p in layouts if p.is_file()),
                   'accessibility': 'Not verified', 'permissions': 'Inspect route and middleware; sidebar hint: ' + sidebar_routes.get(sidebar_path, 'not matched'),
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
            for match in re.finditer(r'#[0-9a-fA-F]{3,8}\b|\b(?:bg|text|border|ring|from|to|via)-blue-\d+\b|--(?:primary|secondary|accent|background|foreground)[\w-]*', line):
                value = match.group()
                category = ('Candidate brand token' if value.startswith(('--primary', '--secondary', '--accent')) else
                            'Candidate structural token' if value.startswith(('--background', '--foreground')) else
                            'Candidate hardcoded blue; may be semantic' if 'blue-' in value else 'Hex literal; role unknown')
                colors.append({'source': path, 'line': line_no, 'value': value, 'classification': category + ' — unreviewed',
                               'proposed_token': 'Needs review', 'risk': 'Needs review'})
            if re.search(r'font-family|fontFamily|next/font|\b(?:text-(?:xs|sm|base|lg|xl|\d+xl)|font-(?:normal|medium|semibold|bold))\b', line):
                typography.append({'source': path, 'line': line_no, 'excerpt': line.strip()[:240], 'classification': 'Needs review'})
write('color-inventory.csv', ['source', 'line', 'value', 'classification', 'proposed_token', 'risk'], colors)
write('typography-inventory.csv', ['source', 'line', 'excerpt', 'classification'], typography)

components = []
for component in sorted((dashboard / 'src/components').rglob('*.tsx')):
    if '__tests__' in component.parts or component.name.endswith('.test.tsx'): continue
    source = component.read_text(encoding='utf-8', errors='replace')
    components.append({'source': component.relative_to(root).as_posix(),
                       'primary_token_mentions': len(re.findall(r'(?:primary-\d+|--primary)', source)),
                       'blue_utility_mentions': len(re.findall(r'\b(?:bg|text|border|ring)-blue-\d+', source)),
                       'dark_utility_mentions': len(re.findall(r'\bdark:', source)),
                       'font_utility_mentions': len(re.findall(r'\bfont-(?:normal|medium|semibold|bold)', source)),
                       'review_status': 'Unreviewed; inspect actual role and accessibility'})
write('component-appearance-inventory.csv', ['source', 'primary_token_mentions', 'blue_utility_mentions', 'dark_utility_mentions', 'font_utility_mentions', 'review_status'], components)

impact = [{'route': row['route'], 'source': row['source'], 'current_behavior': 'Not individually verified',
           'proposed_behavior': 'Inherit validated tenant workspace tokens where applicable',
           'dependencies': row['layouts'], 'risk': 'Needs review', 'tests': 'Locale, theme, keyboard, permissions and workflow',
           'capture_impact': 'Determine active marketing usage'} for row in routes]
write('migration-impact-map.csv', ['route', 'source', 'current_behavior', 'proposed_behavior', 'dependencies', 'risk', 'tests', 'capture_impact'], impact)

matrix = root / 'frontend/apps/landing/public/images/product/servio'
metadata = json.loads((matrix / 'screenshot-matrix.json').read_text(encoding='utf-8'))
captures = {item['file']: item for item in metadata['captures']}
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
    assets.append({'path': image.relative_to(root).as_posix(), 'absolute_path': str(image), 'bytes': image.stat().st_size,
                   'feature': record.get('name', 'Unspecified'), 'locale': record.get('locale', 'Unknown'),
                   'theme': record.get('theme', 'Unknown'), 'width': width,
                   'height': height, 'aspect_ratio': round(width / height, 4) if isinstance(width, int) and isinstance(height, int) and height else 'Unknown',
                   'source_app': ('Storefront checkout' if record.get('name') == 'checkout' else 'Dashboard' if record else 'Unknown'),
                   'source_route': 'Not recorded in matrix; verify', 'capture_method': 'Playwright, synthetic/redacted per provenance' if record else 'Unknown',
                   'source_release': record.get('sourceRelease', 'Unknown'), 'legacy_branding': 'Review visual',
                   'usage': 'See CaptureMedia.tsx and captures.json; per-route usage unverified', 'regeneration': 'Review after rebrand'})
write('marketing-screenshot-inventory.csv', ['path', 'absolute_path', 'bytes', 'feature', 'locale', 'theme', 'width', 'height', 'aspect_ratio', 'source_app', 'source_route', 'capture_method', 'source_release', 'legacy_branding', 'usage', 'regeneration'], assets)
print(f'routes={len(routes)} components={len(components)} color_occurrences={len(colors)} typography_lines={len(typography)} product_image_assets={len(assets)}')
