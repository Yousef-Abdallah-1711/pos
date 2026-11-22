-- Keep global identity and membership data outside the ordinary tenant role.
-- This file is included after create_runtime_roles.sql's broad existing-table grants.
DO $$
DECLARE
    protected_table TEXT;
BEGIN
    FOREACH protected_table IN ARRAY ARRAY[
        'auth_identities',
        'identity_email_challenges',
        'tenant_memberships',
        'membership_restaurant_grants',
        'membership_role_grants'
    ] LOOP
        IF to_regclass(protected_table) IS NOT NULL THEN
            EXECUTE format('REVOKE ALL ON TABLE %I FROM PUBLIC, pos_app', protected_table);
        END IF;
    END LOOP;
END;
$$;
