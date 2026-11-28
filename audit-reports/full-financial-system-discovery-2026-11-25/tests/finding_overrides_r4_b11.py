R4 = {}
R4['F-082'] = dict(
    design='Proposed (root cause PROBABLE: first isolate the 803 px element with the browser inspector at 768 px): wrap the report table or summary grid in frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx with an overflow-x-auto container, or let the grid collapse to one column at the md breakpoint; add a viewport test at 768, 1024 and 1440 px asserting scrollWidth <= innerWidth.',
    alternatives='A: overflow wrapper (recommended, smallest). B: responsive grid collapse.',
    acceptance=['At 768 px the page has no horizontal scroll (document scrollWidth <= 768).', 'The table remains fully readable by scrolling inside its container.', 'No regression at 360, 390, 1024 and 1440 px, English and Arabic.'],
    proposed_files=[], api=['unchanged'], policy='None.', migration='No.', historical='None.')
