SELECT p.proname, pg_get_functiondef(p.oid) FROM pg_proc p WHERE p.pronamespace='public'::regnamespace AND p.oid IN (SELECT tgfoid FROM pg_trigger WHERE NOT tgisinternal);
