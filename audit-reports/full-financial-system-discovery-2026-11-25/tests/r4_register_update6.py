import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
N8 = 'evidence/network/RUN-004/PH-08/'
N5 = 'evidence/network/RUN-004/PH-05/'
F = {
    'F-077': dict(
        title='Saving personal Appearance colours also rewrites the active public storefront theme, silently overwriting Theme Builder edits',
        category='Software defect (cross-feature coupling)', severity='Medium', evidence_status='REPRODUCED',
        affected_modules='Settings > Appearance, Theme Builder, public storefront',
        reproduction='Headed Chrome tenant 4: save #ff0000 on theme 13 in the editor and activate it; open Settings > Appearance and click Theme Colors; read themes_v2.config',
        expected='A user preference changes only that user settings',
        actual='The click sends PUT /users/4/settings/theme-colors and then PUT /admin/themes/{active id} with the merged colours (primary, secondary, accent from the user profile); the active theme colours and updated_at change (observed on themes 4 and 13), so edits made in the Theme Builder are overwritten whenever any user with settings access saves Appearance',
        monetary_impact_EGP='0',
        root_cause_first_break='CONFIRMED (source + runtime): syncBrandColorsToActiveTheme in frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx:38-95 copies the user colours into the active tenant theme via themeApi.updateTheme',
        source_refs='frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx:38; frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx:78',
        related_static_ids='F-059',
        evidence_paths=N5 + 'AUD-R4-APPEARANCE-THEME.json; ' + N8 + 'AUD-R4-THEME-SAVE-VS-ACTIVATE.json; tests/runtime/ph05_r4_appearance_theme.mjs',
        remediation_direction='Make brand colours a tenant-level setting edited in one place; stop writing the tenant theme from a per-user preference screen; require the theme-edit permission for any change to the public theme',
        regression_tests_needed='Saving Appearance leaves the active theme row unchanged; Theme Builder colours persist after an Appearance save'),
    'F-078': dict(
        title='Theme edits and activation do not refresh the public theme: the 5-minute server cache is never invalidated (and the storefront adds a 1-hour page cache)',
        category='Software defect (cache invalidation)', severity='Low', evidence_status='REPRODUCED',
        affected_modules='Theme Builder, public theme API, storefront',
        reproduction='Tenant 4: prime GET /public/restaurants/audit-c-diag/theme, save #00ff00 in the editor, activate the theme, poll the public API every 45 s',
        expected='The public theme reflects an activation or save immediately (or within seconds)',
        actual='Public API kept returning #3b82f6 for 182 s after activation although the database held #00ff00, and flipped only at 227 s when the 5-minute entry expired; the storefront page additionally revalidates every 3,600 s (app/[locale]/page.tsx:75,131) and the response carries Cache-Control max-age=300',
        monetary_impact_EGP='0',
        root_cause_first_break='CONFIRMED (source + runtime): invalidatePublicThemeCacheByRestaurant is defined (handler/http/theme_public_cache.go:45) but has no caller; GetThemeOnly caches the payload for publicThemeCacheTTL = 5 min (public_homepage_handler.go:117-144)',
        source_refs='backend/internal/handler/http/theme_public_cache.go:45; backend/internal/handler/http/public_homepage_handler.go:117; frontend/apps/restaurant-website/app/[locale]/page.tsx:75',
        related_static_ids='F-077',
        evidence_paths=N8 + 'AUD-R4-THEME-CACHE-TTL2.json; ' + N8 + 'AUD-R4-THEME-PAYLOAD.json; tests/runtime/ph08_r4_theme_cache_ttl2.mjs',
        remediation_direction='Call invalidatePublicThemeCacheByRestaurant from theme activate, update and delete; use on-demand revalidation for the storefront instead of a fixed 3,600 s',
        regression_tests_needed='After activate or save the public API returns the new colour within 2 s'),
}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
have = {r[0] for r in rows}
for fid, v in F.items():
    if fid not in have:
        rows.append([dict(v, finding_id=fid).get(k, '') for k in hd])
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
new = [('EX-RUN004-THEMEBUILDER-UI', 'UI-HEADED-CHROME', 'Theme Builder through the real UI: export dialog and download, duplicate, editor colour edit (invalid input blocked client-side, valid saved to config), activation, storefront and public API effect, Appearance coupling, restore', 'owner', 'PARTIAL',
        'export 1,612-byte JSON works; duplicate POST 201; invalid colour text sends no request (client validation); valid colour saved into config JSON not the primary_color column; activation works; public theme stale for about 4 min (F-078); Appearance save overwrites the active theme (F-077); history tab text only (API 501 known F-060); deletion of the active theme and import not exercised; evidence ' + N8)]
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    for r in new:
        if r[0] not in ex:
            csv.writer(h).writerow(r)
print(len(rows) - 1)
