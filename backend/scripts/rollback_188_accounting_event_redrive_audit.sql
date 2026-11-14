-- Manual pre-enable rollback only. Stop the redrive endpoint first.
-- Preserve migration receipts and never remove persisted operator audit history.
BEGIN;
LOCK TABLE accounting_event_redrive_audit IN ACCESS EXCLUSIVE MODE;
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM accounting_event_redrive_audit) THEN
        RAISE EXCEPTION 'cannot downgrade accounting redrive: audit history exists';
    END IF;
END;
$$;
DROP TABLE accounting_event_redrive_audit;
DROP INDEX idx_accounting_event_outbox_failed_scope;
COMMIT;
