-- Reserve tracked stock at confirmation; release on cancellation, not monetary refund.
-- Existing orders are deliberately not backfilled: their historical stock cannot be inferred.
CREATE TABLE order_stock_reservations (
    order_id BIGINT NOT NULL REFERENCES orders(id),
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    restaurant_id BIGINT NOT NULL REFERENCES restaurants(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity BIGINT NOT NULL CHECK (quantity > 0),
    released BOOLEAN NOT NULL DEFAULT false,
    reserved_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    released_at TIMESTAMPTZ,
    PRIMARY KEY (order_id, product_id)
);
CREATE INDEX order_stock_reservations_scope ON order_stock_reservations(tenant_id, restaurant_id, order_id);
-- The runtime cannot forge or erase the idempotency ledger. Trigger-only writes use
-- the migration owner, a fixed search_path and scope from the actual locked order row.
REVOKE ALL ON order_stock_reservations FROM PUBLIC, pos_app;

CREATE FUNCTION maintain_order_stock() RETURNS trigger
LANGUAGE plpgsql SECURITY DEFINER SET search_path = pg_catalog, public AS $$
DECLARE item RECORD;
BEGIN
    IF TG_OP = 'DELETE' THEN
        IF EXISTS (SELECT 1 FROM public.order_stock_reservations
                   WHERE order_id = OLD.id AND tenant_id = OLD.tenant_id AND restaurant_id = OLD.restaurant_id) THEN
            RAISE EXCEPTION 'Orders with inventory history cannot be deleted' USING ERRCODE = '23514';
        END IF;
        RETURN OLD;
    END IF;
    IF NEW.tenant_id <> OLD.tenant_id OR NEW.restaurant_id <> OLD.restaurant_id THEN
        RAISE EXCEPTION 'Order scope is immutable' USING ERRCODE = '23514';
    END IF;
    IF OLD.status = 'pending' AND NEW.status = 'confirmed' THEN
        -- Consistent lock order serializes competing orders without overselling.
        FOR item IN
            SELECT p.id, p.quantity_in_stock, q.quantity
            FROM public.products p
            JOIN (SELECT product_id, SUM(quantity)::BIGINT AS quantity
                  FROM public.order_items
                  WHERE order_id = OLD.id AND tenant_id = OLD.tenant_id AND restaurant_id = OLD.restaurant_id
                  GROUP BY product_id) q ON q.product_id = p.id
            WHERE p.tenant_id = OLD.tenant_id AND p.restaurant_id = OLD.restaurant_id AND p.track_inventory
            ORDER BY p.id FOR UPDATE OF p
        LOOP
            IF item.quantity <= 0 OR item.quantity_in_stock < item.quantity THEN
                RAISE EXCEPTION 'Insufficient inventory for product %', item.id USING ERRCODE = '23514';
            END IF;
            INSERT INTO public.order_stock_reservations(order_id, tenant_id, restaurant_id, product_id, quantity)
            VALUES (OLD.id, OLD.tenant_id, OLD.restaurant_id, item.id, item.quantity);
            UPDATE public.products SET quantity_in_stock = quantity_in_stock - item.quantity, updated_at = CURRENT_TIMESTAMP
            WHERE id = item.id AND tenant_id = OLD.tenant_id AND restaurant_id = OLD.restaurant_id;
        END LOOP;
    ELSIF NEW.status = 'cancelled' AND OLD.status <> 'cancelled' THEN
        FOR item IN
            SELECT p.id, s.quantity FROM public.order_stock_reservations s
            JOIN public.products p ON p.id = s.product_id AND p.tenant_id = s.tenant_id AND p.restaurant_id = s.restaurant_id
            WHERE s.order_id = OLD.id AND s.tenant_id = OLD.tenant_id AND s.restaurant_id = OLD.restaurant_id AND NOT s.released
            ORDER BY p.id FOR UPDATE OF p
        LOOP
            UPDATE public.products SET quantity_in_stock = quantity_in_stock + item.quantity, updated_at = CURRENT_TIMESTAMP
            WHERE id = item.id AND tenant_id = OLD.tenant_id AND restaurant_id = OLD.restaurant_id;
            UPDATE public.order_stock_reservations SET released = true, released_at = CURRENT_TIMESTAMP
            WHERE order_id = OLD.id AND product_id = item.id AND tenant_id = OLD.tenant_id AND restaurant_id = OLD.restaurant_id;
        END LOOP;
    END IF;
    RETURN NEW;
END;
$$;
REVOKE ALL ON FUNCTION maintain_order_stock() FROM PUBLIC, pos_app;
CREATE TRIGGER order_stock_transition BEFORE UPDATE OR DELETE ON orders
FOR EACH ROW EXECUTE FUNCTION maintain_order_stock();

-- Prevent post-confirmation item edits from invalidating the reserved quantity.
-- Lock the parent before item writes, matching confirmation's order -> product lock order.
CREATE FUNCTION guard_order_stock_items() RETURNS trigger
LANGUAGE plpgsql SECURITY DEFINER SET search_path = pg_catalog, public AS $$
DECLARE parent_status TEXT; source_row RECORD;
BEGIN
    IF TG_OP = 'DELETE' THEN source_row := OLD; ELSE source_row := NEW; END IF;
    IF TG_OP = 'UPDATE' AND (NEW.order_id <> OLD.order_id OR NEW.tenant_id <> OLD.tenant_id OR NEW.restaurant_id <> OLD.restaurant_id) THEN
        RAISE EXCEPTION 'Order item scope is immutable' USING ERRCODE = '23514';
    END IF;
    SELECT status INTO parent_status FROM public.orders
    WHERE id = source_row.order_id AND tenant_id = source_row.tenant_id AND restaurant_id = source_row.restaurant_id FOR UPDATE;
    IF parent_status IS NOT NULL AND parent_status <> 'pending' THEN
        RAISE EXCEPTION 'Only pending order items may be changed' USING ERRCODE = '23514';
    END IF;
    RETURN source_row;
END;
$$;
REVOKE ALL ON FUNCTION guard_order_stock_items() FROM PUBLIC, pos_app;
CREATE TRIGGER order_stock_items BEFORE INSERT OR UPDATE OR DELETE ON order_items
FOR EACH ROW EXECUTE FUNCTION guard_order_stock_items();

-- Rollback (maintenance window only, after explicit reconciliation of reservations):
-- DROP TRIGGER order_stock_items ON order_items;
-- DROP TRIGGER order_stock_transition ON orders;
-- DROP FUNCTION guard_order_stock_items();
-- DROP FUNCTION maintain_order_stock();
-- DROP TABLE order_stock_reservations;
-- Do not automatically add reserved quantities back during rollback: delivered sales remain consumed.
