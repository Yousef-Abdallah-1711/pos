-- Preserve the interpretation of posted ledger history when an account's
-- natural balance direction is edited outside the application repositories.

CREATE OR REPLACE FUNCTION cfr003_guard_normal_balance_history()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
BEGIN
    IF NEW.normal_balance IS DISTINCT FROM OLD.normal_balance
       AND EXISTS (
           SELECT 1
           FROM public.general_ledger gl
           WHERE gl.tenant_id = OLD.tenant_id
             AND gl.restaurant_id = OLD.restaurant_id
             AND gl.account_id = OLD.id
       ) THEN
        RAISE EXCEPTION 'an account normal balance cannot change after ledger history exists'
            USING ERRCODE = '23514';
    END IF;

    RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION cfr003_lock_account_before_ledger_insert()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
DECLARE
    account_tenant_id INT;
    account_restaurant_id INT;
BEGIN
    SELECT tenant_id, restaurant_id
    INTO account_tenant_id, account_restaurant_id
    FROM public.chart_of_accounts
    WHERE id = NEW.account_id
    FOR UPDATE;

    IF NOT FOUND
       OR account_tenant_id <> NEW.tenant_id
       OR account_restaurant_id <> NEW.restaurant_id THEN
        RAISE EXCEPTION 'General Ledger account is outside its tenant and restaurant'
            USING ERRCODE = '23514';
    END IF;

    RETURN NEW;
END;
$$;

-- Keep the journal-then-account lock order used by repository posting while
-- serializing ledger inserts against normal_balance updates before row write.
DROP TRIGGER IF EXISTS cfr003_lock_account_after_cfr002_general_ledger_insert ON general_ledger;
CREATE TRIGGER cfr003_lock_account_after_cfr002_general_ledger_insert
BEFORE INSERT ON general_ledger
FOR EACH ROW EXECUTE FUNCTION cfr003_lock_account_before_ledger_insert();

DROP TRIGGER IF EXISTS cfr003_normal_balance_history_guard ON chart_of_accounts;
CREATE TRIGGER cfr003_normal_balance_history_guard
BEFORE UPDATE OF normal_balance ON chart_of_accounts
FOR EACH ROW EXECUTE FUNCTION cfr003_guard_normal_balance_history();
