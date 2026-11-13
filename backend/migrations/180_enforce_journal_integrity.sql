-- CFR-002: enforce posted journal and General Ledger invariants in PostgreSQL.
-- This migration deliberately rejects inconsistent historical rows. It never
-- repairs or rewrites financial history; run scripts/preflight_180_journal_integrity.sql first.

BEGIN;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM journal_entries je
        LEFT JOIN LATERAL (
            SELECT COUNT(*) AS line_count,
                   COALESCE(SUM(debit_amount), 0) AS debit_total,
                   COALESCE(SUM(credit_amount), 0) AS credit_total
            FROM journal_entry_lines jel
            WHERE jel.journal_entry_id = je.id
        ) totals ON TRUE
        WHERE (je.status = 'posted' OR (je.status = 'void' AND je.posted_at IS NOT NULL))
          AND (totals.line_count < 2
               OR totals.debit_total <> totals.credit_total
               OR totals.debit_total <> je.total_debit
               OR totals.credit_total <> je.total_credit)
    ) THEN
        RAISE EXCEPTION 'CFR-002 migration blocked: posted journals are missing lines or have inconsistent totals'
            USING ERRCODE = '23514',
                  HINT = 'Run scripts/preflight_180_journal_integrity.sql and remediate historical records under approved accounting procedure.';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM journal_entry_lines jel
        JOIN journal_entries je ON je.id = jel.journal_entry_id
        LEFT JOIN chart_of_accounts coa ON coa.id = jel.account_id
        JOIN restaurants r ON r.id = je.restaurant_id
        WHERE jel.tenant_id <> je.tenant_id
           OR r.tenant_id <> je.tenant_id
           OR coa.id IS NULL
           OR coa.tenant_id <> je.tenant_id
           OR coa.restaurant_id <> je.restaurant_id
    ) THEN
        RAISE EXCEPTION 'CFR-002 migration blocked: journal lines cross tenant or restaurant boundaries'
            USING ERRCODE = '23514',
                  HINT = 'Run scripts/preflight_180_journal_integrity.sql; do not reassign historical financial rows automatically.';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM general_ledger gl
        JOIN journal_entries je ON je.id = gl.journal_entry_id
        LEFT JOIN journal_entry_lines jel
          ON jel.id = gl.journal_line_id
         AND jel.journal_entry_id = gl.journal_entry_id
        LEFT JOIN chart_of_accounts coa ON coa.id = gl.account_id
        WHERE jel.id IS NULL
           OR gl.tenant_id <> je.tenant_id
           OR gl.restaurant_id <> je.restaurant_id
           OR gl.tenant_id <> jel.tenant_id
           OR gl.account_id <> jel.account_id
           OR coa.id IS NULL
           OR coa.tenant_id <> gl.tenant_id
           OR coa.restaurant_id <> gl.restaurant_id
           OR NOT (
                (gl.debit_amount = jel.debit_amount AND gl.credit_amount = jel.credit_amount)
                OR (je.status = 'void' AND je.posted_at IS NOT NULL
                    AND gl.debit_amount = jel.credit_amount
                    AND gl.credit_amount = jel.debit_amount)
           )
    ) THEN
        RAISE EXCEPTION 'CFR-002 migration blocked: General Ledger rows do not match their journal lines and scopes'
            USING ERRCODE = '23514',
                  HINT = 'Run scripts/preflight_180_journal_integrity.sql and reconcile historical ledger rows before retrying.';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM journal_entries je
        JOIN journal_entry_lines jel ON jel.journal_entry_id = je.id
        LEFT JOIN general_ledger gl
          ON gl.journal_entry_id = je.id
         AND gl.journal_line_id = jel.id
        GROUP BY je.id, je.status, je.posted_at, jel.id, jel.debit_amount, jel.credit_amount
        HAVING (je.status = 'posted' AND (
                    COUNT(gl.id) <> 1
                    OR COUNT(*) FILTER (WHERE gl.debit_amount = jel.debit_amount AND gl.credit_amount = jel.credit_amount) <> 1
                ))
            OR (je.status = 'void' AND je.posted_at IS NOT NULL AND (
                    COUNT(gl.id) <> 2
                    OR COUNT(*) FILTER (WHERE gl.debit_amount = jel.debit_amount AND gl.credit_amount = jel.credit_amount) <> 1
                    OR COUNT(*) FILTER (WHERE gl.debit_amount = jel.credit_amount AND gl.credit_amount = jel.debit_amount) <> 1
                ))
            OR ((je.status = 'draft' OR (je.status = 'void' AND je.posted_at IS NULL)) AND COUNT(gl.id) <> 0)
    ) THEN
        RAISE EXCEPTION 'CFR-002 migration blocked: historical journal lines have missing or duplicate ledger effects'
            USING ERRCODE = '23514',
                  HINT = 'Review posted and voided journals with the read-only preflight; no ledger effects were changed.';
    END IF;
END;
$$;

CREATE INDEX IF NOT EXISTS idx_general_ledger_journal_line
    ON general_ledger (journal_entry_id, journal_line_id);

CREATE OR REPLACE FUNCTION cfr002_guard_journal_header_write()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
DECLARE
    restaurant_tenant_id INT;
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.status <> 'draft' THEN
            RAISE EXCEPTION 'new journal entries must start as draft' USING ERRCODE = '23514';
        END IF;
        SELECT tenant_id INTO restaurant_tenant_id
        FROM public.restaurants WHERE id = NEW.restaurant_id FOR UPDATE;
        IF restaurant_tenant_id IS NULL OR restaurant_tenant_id <> NEW.tenant_id THEN
            RAISE EXCEPTION 'journal restaurant does not belong to its tenant' USING ERRCODE = '23514';
        END IF;
        RETURN NEW;
    END IF;

    IF TG_OP = 'DELETE' THEN
        IF OLD.status <> 'draft' THEN
            RAISE EXCEPTION 'posted or void journals cannot be deleted' USING ERRCODE = '23514';
        END IF;
        RETURN OLD;
    END IF;

    IF OLD.status = 'draft' THEN
        IF NEW.status NOT IN ('draft', 'posted', 'void') THEN
            RAISE EXCEPTION 'invalid journal status transition from draft to %', NEW.status USING ERRCODE = '23514';
        END IF;
        IF NEW.restaurant_id <> OLD.restaurant_id OR NEW.tenant_id <> OLD.tenant_id THEN
            SELECT tenant_id INTO restaurant_tenant_id
            FROM public.restaurants WHERE id = NEW.restaurant_id FOR UPDATE;
            IF restaurant_tenant_id IS NULL OR restaurant_tenant_id <> NEW.tenant_id THEN
                RAISE EXCEPTION 'journal restaurant does not belong to its tenant' USING ERRCODE = '23514';
            END IF;
        END IF;
        IF NEW.status = 'posted' AND (NEW.posted_at IS NULL OR NEW.posted_by IS NULL) THEN
            RAISE EXCEPTION 'posting requires posted_at and posted_by audit fields' USING ERRCODE = '23514';
        END IF;
        IF NEW.status = 'void' AND (NEW.voided_at IS NULL OR NEW.voided_by IS NULL OR btrim(COALESCE(NEW.void_reason, '')) = '') THEN
            RAISE EXCEPTION 'voiding a draft requires void audit fields and a reason' USING ERRCODE = '23514';
        END IF;
        RETURN NEW;
    END IF;

    IF OLD.status = 'posted' AND NEW.status = 'void' THEN
        IF ROW(NEW.id, NEW.tenant_id, NEW.restaurant_id, NEW.entry_number, NEW.entry_date,
               NEW.fiscal_period, NEW.description, NEW.description_ar, NEW.reference_type,
               NEW.reference_id, NEW.reference_number, NEW.entry_type, NEW.total_debit,
               NEW.total_credit, NEW.posted_at, NEW.posted_by, NEW.is_approved,
               NEW.approved_by, NEW.approved_at, NEW.notes, NEW.attachments,
               NEW.created_at, NEW.created_by)
           IS DISTINCT FROM
           ROW(OLD.id, OLD.tenant_id, OLD.restaurant_id, OLD.entry_number, OLD.entry_date,
               OLD.fiscal_period, OLD.description, OLD.description_ar, OLD.reference_type,
               OLD.reference_id, OLD.reference_number, OLD.entry_type, OLD.total_debit,
               OLD.total_credit, OLD.posted_at, OLD.posted_by, OLD.is_approved,
               OLD.approved_by, OLD.approved_at, OLD.notes, OLD.attachments,
               OLD.created_at, OLD.created_by) THEN
            RAISE EXCEPTION 'voiding cannot modify posted journal accounting history' USING ERRCODE = '23514';
        END IF;
        IF NEW.voided_at IS NULL OR NEW.voided_by IS NULL OR btrim(COALESCE(NEW.void_reason, '')) = '' THEN
            RAISE EXCEPTION 'voiding a posted journal requires void audit fields and a reason' USING ERRCODE = '23514';
        END IF;
        RETURN NEW;
    END IF;

    RAISE EXCEPTION 'journal status % is immutable; permitted transitions are draft to posted/void and posted to audited void', OLD.status
        USING ERRCODE = '23514';
END;
$$;

CREATE OR REPLACE FUNCTION cfr002_guard_journal_line_write()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
DECLARE
    target_entry_id INT;
    journal_status TEXT;
    journal_tenant_id INT;
    journal_restaurant_id INT;
    target_tenant_id INT;
    target_account_id INT;
BEGIN
    IF TG_OP = 'DELETE' THEN
        target_entry_id := OLD.journal_entry_id;
    ELSE
        target_entry_id := NEW.journal_entry_id;
    END IF;

    SELECT status, tenant_id, restaurant_id
    INTO journal_status, journal_tenant_id, journal_restaurant_id
    FROM public.journal_entries
    WHERE id = target_entry_id
    FOR UPDATE;

    IF NOT FOUND AND TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSIF NOT FOUND THEN
        RAISE EXCEPTION 'journal entry does not exist' USING ERRCODE = '23503';
    END IF;
    IF journal_status <> 'draft' THEN
        RAISE EXCEPTION 'journal lines are immutable after posting or voiding' USING ERRCODE = '23514';
    END IF;

    IF TG_OP <> 'DELETE' THEN
        IF NEW.tenant_id <> journal_tenant_id THEN
            RAISE EXCEPTION 'journal line tenant does not match its journal' USING ERRCODE = '23514';
        END IF;
        IF TG_OP = 'UPDATE' AND
           (NEW.journal_entry_id <> OLD.journal_entry_id OR NEW.tenant_id <> OLD.tenant_id) THEN
            RAISE EXCEPTION 'journal lines cannot be moved between journal entries or tenants' USING ERRCODE = '23514';
        END IF;
        SELECT tenant_id INTO target_tenant_id
        FROM public.chart_of_accounts
        WHERE id = NEW.account_id AND restaurant_id = journal_restaurant_id
        FOR UPDATE;
        IF target_tenant_id IS NULL OR target_tenant_id <> journal_tenant_id THEN
            RAISE EXCEPTION 'journal line account does not belong to its tenant and restaurant' USING ERRCODE = '23514';
        END IF;
        RETURN NEW;
    END IF;

    RETURN OLD;
END;
$$;

CREATE OR REPLACE FUNCTION cfr002_guard_account_scope_write()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
BEGIN
    IF NEW.tenant_id IS DISTINCT FROM OLD.tenant_id
       OR NEW.restaurant_id IS DISTINCT FROM OLD.restaurant_id THEN
        IF EXISTS (SELECT 1 FROM public.journal_entry_lines WHERE account_id = OLD.id)
           OR EXISTS (SELECT 1 FROM public.general_ledger WHERE account_id = OLD.id) THEN
            RAISE EXCEPTION 'an account referenced by journal history cannot move between tenant or restaurant scopes'
                USING ERRCODE = '23514';
        END IF;
    END IF;
    RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION cfr002_guard_restaurant_scope_write()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
BEGIN
    IF NEW.tenant_id IS DISTINCT FROM OLD.tenant_id
       AND EXISTS (SELECT 1 FROM public.journal_entries WHERE restaurant_id = OLD.id) THEN
        RAISE EXCEPTION 'a restaurant with journal history cannot move between tenants'
            USING ERRCODE = '23514';
    END IF;
    RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION cfr002_guard_general_ledger_insert()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
DECLARE
    journal_status TEXT;
    journal_tenant_id INT;
    journal_restaurant_id INT;
    journal_entry_date DATE;
    journal_fiscal_period VARCHAR(7);
    journal_voided_at TIMESTAMP;
    journal_void_reason TEXT;
    line_tenant_id INT;
    line_account_id INT;
    line_debit NUMERIC(15,2);
    line_credit NUMERIC(15,2);
    account_tenant_id INT;
    account_restaurant_id INT;
    existing_count INT;
    original_count INT;
    reverse_count INT;
BEGIN
    SELECT status, tenant_id, restaurant_id, entry_date, fiscal_period, voided_at, void_reason
    INTO journal_status, journal_tenant_id, journal_restaurant_id, journal_entry_date,
         journal_fiscal_period, journal_voided_at, journal_void_reason
    FROM public.journal_entries
    WHERE id = NEW.journal_entry_id
    FOR UPDATE;
    IF NOT FOUND OR journal_status NOT IN ('posted', 'void') THEN
        RAISE EXCEPTION 'General Ledger rows require a posted journal or its audited void transition' USING ERRCODE = '23514';
    END IF;
    IF NEW.tenant_id <> journal_tenant_id OR NEW.restaurant_id <> journal_restaurant_id THEN
        RAISE EXCEPTION 'General Ledger tenant or restaurant does not match its journal' USING ERRCODE = '23514';
    END IF;

    SELECT tenant_id, account_id, debit_amount, credit_amount
    INTO line_tenant_id, line_account_id, line_debit, line_credit
    FROM public.journal_entry_lines
    WHERE id = NEW.journal_line_id AND journal_entry_id = NEW.journal_entry_id;
    IF NOT FOUND OR NEW.tenant_id <> line_tenant_id OR NEW.account_id <> line_account_id THEN
        RAISE EXCEPTION 'General Ledger row does not match its journal line' USING ERRCODE = '23514';
    END IF;

    SELECT tenant_id, restaurant_id INTO account_tenant_id, account_restaurant_id
    FROM public.chart_of_accounts WHERE id = NEW.account_id;
    IF NOT FOUND OR account_tenant_id <> NEW.tenant_id OR account_restaurant_id <> NEW.restaurant_id THEN
        RAISE EXCEPTION 'General Ledger account is outside its tenant and restaurant' USING ERRCODE = '23514';
    END IF;

    SELECT COUNT(*),
           COUNT(*) FILTER (WHERE debit_amount = line_debit AND credit_amount = line_credit),
           COUNT(*) FILTER (WHERE debit_amount = line_credit AND credit_amount = line_debit)
    INTO existing_count, original_count, reverse_count
    FROM public.general_ledger
    WHERE journal_entry_id = NEW.journal_entry_id AND journal_line_id = NEW.journal_line_id;

    IF existing_count = 0 AND journal_status = 'posted' THEN
        IF NEW.debit_amount <> line_debit OR NEW.credit_amount <> line_credit THEN
            RAISE EXCEPTION 'initial General Ledger effect must match its journal line' USING ERRCODE = '23514';
        END IF;
        IF NEW.entry_date <> journal_entry_date OR
           NEW.fiscal_period <> COALESCE(journal_fiscal_period, to_char(journal_entry_date, 'YYYY-MM')) THEN
            RAISE EXCEPTION 'initial General Ledger period must match its journal' USING ERRCODE = '23514';
        END IF;
    ELSIF journal_status = 'void' AND journal_voided_at IS NOT NULL
          AND btrim(COALESCE(journal_void_reason, '')) <> ''
          AND existing_count = 1 AND original_count = 1 AND
          NEW.debit_amount = line_credit AND NEW.credit_amount = line_debit AND reverse_count = 0 THEN
        NULL; -- The repository transitions the journal to audited void before adding inverse effects.
    ELSE
        RAISE EXCEPTION 'duplicate or invalid General Ledger effect for journal line %', NEW.journal_line_id
            USING ERRCODE = '23514';
    END IF;

    RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION cfr002_reject_general_ledger_mutation()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
BEGIN
    RAISE EXCEPTION 'General Ledger rows are immutable; record corrections through an audited journal reversal'
        USING ERRCODE = '23514';
END;
$$;

CREATE OR REPLACE FUNCTION cfr002_assert_journal_integrity_deferred()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
DECLARE
    target_entry_id INT;
    journal_status TEXT;
    journal_posted_at TIMESTAMP;
    header_debit NUMERIC(15,2);
    header_credit NUMERIC(15,2);
    line_count BIGINT;
    line_debit NUMERIC;
    line_credit NUMERIC;
    invalid_line_id INT;
    ledger_effect_count BIGINT;
    matching_effect_count BIGINT;
BEGIN
    IF TG_TABLE_NAME = 'journal_entries' THEN
        IF TG_OP = 'DELETE' THEN target_entry_id := OLD.id; ELSE target_entry_id := NEW.id; END IF;
    ELSIF TG_TABLE_NAME = 'journal_entry_lines' THEN
        IF TG_OP = 'DELETE' THEN target_entry_id := OLD.journal_entry_id; ELSE target_entry_id := NEW.journal_entry_id; END IF;
    ELSE
        IF TG_OP = 'DELETE' THEN target_entry_id := OLD.journal_entry_id; ELSE target_entry_id := NEW.journal_entry_id; END IF;
    END IF;

    SELECT status, posted_at, total_debit, total_credit
    INTO journal_status, journal_posted_at, header_debit, header_credit
    FROM public.journal_entries WHERE id = target_entry_id;
    IF NOT FOUND THEN RETURN NULL; END IF;

    SELECT COUNT(*), COALESCE(SUM(debit_amount), 0), COALESCE(SUM(credit_amount), 0)
    INTO line_count, line_debit, line_credit
    FROM public.journal_entry_lines WHERE journal_entry_id = target_entry_id;

    IF EXISTS (
        SELECT 1
        FROM public.journal_entry_lines jel
        LEFT JOIN public.chart_of_accounts coa ON coa.id = jel.account_id
        WHERE jel.journal_entry_id = target_entry_id
          AND (jel.tenant_id <> (SELECT tenant_id FROM public.journal_entries WHERE id = target_entry_id)
               OR coa.id IS NULL
               OR coa.tenant_id <> (SELECT tenant_id FROM public.journal_entries WHERE id = target_entry_id)
               OR coa.restaurant_id <> (SELECT restaurant_id FROM public.journal_entries WHERE id = target_entry_id))
    ) THEN
        RAISE EXCEPTION 'journal % contains a line outside its tenant or restaurant scope', target_entry_id
            USING ERRCODE = '23514';
    END IF;

    IF journal_status = 'posted' THEN
        IF line_count < 2 OR line_debit <> line_credit OR
           line_debit <> header_debit OR line_credit <> header_credit THEN
            RAISE EXCEPTION 'posted journal % must contain at least two exactly balanced lines matching its header', target_entry_id
                USING ERRCODE = '23514';
        END IF;
        SELECT jel.id, COUNT(gl.id),
               COUNT(gl.id) FILTER (WHERE gl.debit_amount = jel.debit_amount AND gl.credit_amount = jel.credit_amount)
        INTO invalid_line_id, ledger_effect_count, matching_effect_count
            FROM public.journal_entry_lines jel
            LEFT JOIN public.general_ledger gl
              ON gl.journal_entry_id = target_entry_id AND gl.journal_line_id = jel.id
            WHERE jel.journal_entry_id = target_entry_id
            GROUP BY jel.id, jel.debit_amount, jel.credit_amount
            HAVING COUNT(gl.id) <> 1
                OR COUNT(gl.id) FILTER (WHERE gl.debit_amount = jel.debit_amount AND gl.credit_amount = jel.credit_amount) <> 1
            ORDER BY jel.id
            LIMIT 1;
        IF invalid_line_id IS NOT NULL THEN
            RAISE EXCEPTION 'posted journal % line % must have exactly one matching General Ledger row (found %, matching %)',
                target_entry_id, invalid_line_id, ledger_effect_count, matching_effect_count
                USING ERRCODE = '23514';
        END IF;
    ELSIF journal_status = 'void' AND journal_posted_at IS NOT NULL THEN
        IF line_count < 2 OR line_debit <> line_credit OR
           line_debit <> header_debit OR line_credit <> header_credit THEN
            RAISE EXCEPTION 'voided posted journal % must preserve its original balanced lines and header totals', target_entry_id
                USING ERRCODE = '23514';
        END IF;
        SELECT jel.id, COUNT(gl.id),
               COUNT(gl.id) FILTER (WHERE gl.debit_amount = jel.debit_amount AND gl.credit_amount = jel.credit_amount)
        INTO invalid_line_id, ledger_effect_count, matching_effect_count
            FROM public.journal_entry_lines jel
            LEFT JOIN public.general_ledger gl
              ON gl.journal_entry_id = target_entry_id AND gl.journal_line_id = jel.id
            WHERE jel.journal_entry_id = target_entry_id
            GROUP BY jel.id, jel.debit_amount, jel.credit_amount
            HAVING COUNT(gl.id) <> 2
                OR COUNT(gl.id) FILTER (WHERE gl.debit_amount = jel.debit_amount AND gl.credit_amount = jel.credit_amount) <> 1
                OR COUNT(gl.id) FILTER (WHERE gl.debit_amount = jel.credit_amount AND gl.credit_amount = jel.debit_amount) <> 1
            ORDER BY jel.id
            LIMIT 1;
        IF invalid_line_id IS NOT NULL THEN
            RAISE EXCEPTION 'voided journal % line % must retain one original and one inverse General Ledger row (found %, matching original %)',
                target_entry_id, invalid_line_id, ledger_effect_count, matching_effect_count
                USING ERRCODE = '23514';
        END IF;
    ELSE
        IF EXISTS (SELECT 1 FROM public.general_ledger WHERE journal_entry_id = target_entry_id) THEN
            RAISE EXCEPTION 'draft or unposted void journal % cannot have General Ledger effects', target_entry_id
                USING ERRCODE = '23514';
        END IF;
    END IF;
    RETURN NULL;
END;
$$;

DROP TRIGGER IF EXISTS cfr002_journal_header_write ON journal_entries;
CREATE TRIGGER cfr002_journal_header_write
BEFORE INSERT OR UPDATE OR DELETE ON journal_entries
FOR EACH ROW EXECUTE FUNCTION cfr002_guard_journal_header_write();

DROP TRIGGER IF EXISTS cfr002_journal_line_write ON journal_entry_lines;
CREATE TRIGGER cfr002_journal_line_write
BEFORE INSERT OR UPDATE OR DELETE ON journal_entry_lines
FOR EACH ROW EXECUTE FUNCTION cfr002_guard_journal_line_write();

DROP TRIGGER IF EXISTS cfr002_general_ledger_insert ON general_ledger;
CREATE TRIGGER cfr002_general_ledger_insert
BEFORE INSERT ON general_ledger
FOR EACH ROW EXECUTE FUNCTION cfr002_guard_general_ledger_insert();

DROP TRIGGER IF EXISTS cfr002_account_scope_write ON chart_of_accounts;
CREATE TRIGGER cfr002_account_scope_write
BEFORE UPDATE OF tenant_id, restaurant_id ON chart_of_accounts
FOR EACH ROW EXECUTE FUNCTION cfr002_guard_account_scope_write();

DROP TRIGGER IF EXISTS cfr002_restaurant_scope_write ON restaurants;
CREATE TRIGGER cfr002_restaurant_scope_write
BEFORE UPDATE OF tenant_id ON restaurants
FOR EACH ROW EXECUTE FUNCTION cfr002_guard_restaurant_scope_write();

DROP TRIGGER IF EXISTS cfr002_general_ledger_immutable ON general_ledger;
CREATE TRIGGER cfr002_general_ledger_immutable
BEFORE UPDATE OR DELETE ON general_ledger
FOR EACH ROW EXECUTE FUNCTION cfr002_reject_general_ledger_mutation();

DROP TRIGGER IF EXISTS cfr002_journal_integrity_header ON journal_entries;
CREATE CONSTRAINT TRIGGER cfr002_journal_integrity_header
AFTER INSERT OR UPDATE OR DELETE ON journal_entries
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW EXECUTE FUNCTION cfr002_assert_journal_integrity_deferred();

DROP TRIGGER IF EXISTS cfr002_journal_integrity_lines ON journal_entry_lines;
CREATE CONSTRAINT TRIGGER cfr002_journal_integrity_lines
AFTER INSERT OR UPDATE OR DELETE ON journal_entry_lines
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW EXECUTE FUNCTION cfr002_assert_journal_integrity_deferred();

DROP TRIGGER IF EXISTS cfr002_journal_integrity_ledger ON general_ledger;
CREATE CONSTRAINT TRIGGER cfr002_journal_integrity_ledger
AFTER INSERT OR UPDATE OR DELETE ON general_ledger
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW EXECUTE FUNCTION cfr002_assert_journal_integrity_deferred();

COMMIT;
