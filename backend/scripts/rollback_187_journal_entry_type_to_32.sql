-- Downgrade for migration 187. Refuse to truncate or rewrite journal history.
BEGIN;
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM journal_entries
        WHERE length(entry_type) > 20
    ) THEN
        RAISE EXCEPTION 'cannot downgrade journal_entries.entry_type: values longer than 20 characters exist';
    END IF;
END;
$$;
ALTER TABLE journal_entries
    ALTER COLUMN entry_type TYPE VARCHAR(20);
COMMIT;

