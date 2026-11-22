-- Keep global identity and membership data outside the ordinary tenant role.
-- This file is included after create_runtime_roles.sql's broad existing-table grants.
DO $$
DECLARE
    protected_table TEXT;
    target_schema TEXT := current_schema();
BEGIN
    FOREACH protected_table IN ARRAY ARRAY[
        'auth_identities',
        'identity_email_challenges',
        'tenant_memberships',
        'membership_restaurant_grants',
        'membership_role_grants'
    ] LOOP
        IF to_regclass(format('%I.%I', target_schema, protected_table)) IS NOT NULL THEN
            EXECUTE format('REVOKE ALL ON TABLE %I.%I FROM PUBLIC, pos_app', target_schema, protected_table);
        END IF;
    END LOOP;
END;
$$;
