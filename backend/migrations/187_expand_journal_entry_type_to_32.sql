-- Supplier auto-posting entry types use 21 characters and must fit journal_entries.entry_type.
-- This widens the catalog value only; existing journal rows are not changed.
ALTER TABLE journal_entries
    ALTER COLUMN entry_type TYPE VARCHAR(32);
