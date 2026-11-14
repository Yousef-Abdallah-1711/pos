-- The core order transaction enqueues work; its worker acknowledges fenced leases.
-- Financial posting still uses pos_accounting. Neither role receives new schema powers.
GRANT SELECT,INSERT,UPDATE ON accounting_event_outbox TO pos_app;
REVOKE DELETE,TRUNCATE,REFERENCES,TRIGGER ON accounting_event_outbox FROM pos_app;
GRANT USAGE,SELECT ON SEQUENCE accounting_event_outbox_id_seq TO pos_app;

-- Sale-history retries claim an immutable scoped receipt without touching stock again.
GRANT SELECT,INSERT ON accounting_inventory_movement_events TO pos_app;
REVOKE UPDATE,DELETE,TRUNCATE,REFERENCES,TRIGGER ON accounting_inventory_movement_events FROM pos_app;

-- Manual rollback with recovery workers stopped (never an automatic image rollback):
-- REVOKE INSERT,UPDATE ON accounting_event_outbox FROM pos_app;
-- REVOKE USAGE,SELECT ON SEQUENCE accounting_event_outbox_id_seq FROM pos_app;
-- REVOKE INSERT ON accounting_inventory_movement_events FROM pos_app;
-- Keep all event/history rows and the dedicated financial writer intact.
