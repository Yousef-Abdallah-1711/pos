SELECT rolname,rolsuper,rolcanlogin,rolbypassrls FROM pg_roles WHERE rolname !~ '^pg_' ORDER BY 1;
SELECT grantee,table_name,string_agg(privilege_type,',' ORDER BY privilege_type) FROM information_schema.role_table_grants WHERE table_schema='public' AND grantee NOT IN ('postgres') GROUP BY 1,2 ORDER BY 2,1;
