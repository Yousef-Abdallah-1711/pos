-- Privileged provisioning only; the API never receives schema-owner credentials.
SELECT 'CREATE ROLE pos_accounting NOLOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS'
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'pos_accounting') \gexec
\if :{?accounting_pw}
ALTER ROLE pos_accounting LOGIN PASSWORD :'accounting_pw';
\endif
ALTER ROLE pos_accounting NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS;
GRANT CONNECT ON DATABASE :"dbname" TO pos_accounting;
\connect :"dbname"
GRANT USAGE ON SCHEMA public TO pos_accounting;
