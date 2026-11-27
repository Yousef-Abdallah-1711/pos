R4 = {}
R4['F-076'] = dict(
    design='Proposed: either add frontend/apps/super-admin/src/app/dashboard/users/[id]/page.tsx (user profile, memberships, status, last login, using the existing super-admin users API) or remove the View link at frontend/apps/super-admin/src/app/dashboard/users/page.tsx:132 until the page exists; add a link-check test that visits every internal href rendered on the page.',
    alternatives='A: build the detail page (recommended if support staff need it). B: remove the link.',
    acceptance=['Every link on the Users page returns 200.', 'The detail page shows the same email, status and role data as the list row.'],
    proposed_files=['frontend/apps/super-admin/src/app/dashboard/users/[id]/page.tsx (new)'], api=['unchanged'], policy='None.', migration='No.', historical='None.')
