-- Partial cancellation is permitted only before fulfillment/payment. Monetary refunds
-- remain separate and never restock. Physical returns use authorized stock adjustments.
ALTER TABLE order_items ADD COLUMN cancelled_quantity INT NOT NULL DEFAULT 0
    CHECK (cancelled_quantity >= 0 AND cancelled_quantity <= quantity);
CREATE TABLE order_item_cancellations (
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    restaurant_id BIGINT NOT NULL REFERENCES restaurants(id),
    order_id BIGINT NOT NULL REFERENCES orders(id),
    item_id BIGINT NOT NULL REFERENCES order_items(id),
    request_key UUID NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    reason TEXT NOT NULL,
    actor_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tenant_id, restaurant_id, order_id, request_key)
);
REVOKE ALL ON order_item_cancellations FROM PUBLIC, pos_app;
-- Runtime edits cannot bypass cancellation validation or its immutable receipt.
REVOKE UPDATE ON order_items FROM pos_app;

CREATE OR REPLACE FUNCTION guard_order_stock_items() RETURNS trigger
LANGUAGE plpgsql SECURITY DEFINER SET search_path = pg_catalog, public AS $$
DECLARE parent_status TEXT; source_row RECORD;
BEGIN
    IF TG_OP = 'DELETE' THEN source_row := OLD; ELSE source_row := NEW; END IF;
    IF TG_OP = 'INSERT' AND NEW.cancelled_quantity <> 0 THEN
        RAISE EXCEPTION 'New items cannot carry unaudited cancellations' USING ERRCODE = '23514';
    END IF;
    IF TG_OP = 'UPDATE' AND (NEW.order_id <> OLD.order_id OR NEW.tenant_id <> OLD.tenant_id OR NEW.restaurant_id <> OLD.restaurant_id) THEN
        RAISE EXCEPTION 'Order item scope is immutable' USING ERRCODE = '23514';
    END IF;
    SELECT status INTO parent_status FROM public.orders
    WHERE id = source_row.order_id AND tenant_id = source_row.tenant_id AND restaurant_id = source_row.restaurant_id FOR UPDATE;
    IF parent_status IS NOT NULL AND parent_status <> 'pending' THEN
        IF TG_OP <> 'UPDATE' OR current_user <> 'pos_migrator' OR parent_status <> 'confirmed'
           OR NEW.cancelled_quantity < OLD.cancelled_quantity
           OR (to_jsonb(NEW) - ARRAY['cancelled_quantity','total_price','updated_at']) <> (to_jsonb(OLD) - ARRAY['cancelled_quantity','total_price','updated_at']) THEN
            RAISE EXCEPTION 'Only audited pre-fulfillment item cancellations are permitted' USING ERRCODE = '23514';
        END IF;
    END IF;
    RETURN source_row;
END;
$$;
REVOKE ALL ON FUNCTION guard_order_stock_items() FROM PUBLIC, pos_app;

CREATE OR REPLACE FUNCTION maintain_order_stock() RETURNS trigger
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
    IF OLD.status = 'delivered' AND NEW.status = 'cancelled' THEN
        RAISE EXCEPTION 'Delivered stock requires an explicit physical return adjustment' USING ERRCODE = '23514';
    END IF;
    IF NEW.tenant_id <> OLD.tenant_id OR NEW.restaurant_id <> OLD.restaurant_id THEN
        RAISE EXCEPTION 'Order scope is immutable' USING ERRCODE = '23514';
    END IF;
    IF OLD.status = 'pending' AND NEW.status = 'confirmed' THEN
        -- Consistent lock order serializes competing orders without overselling.
        FOR item IN
            SELECT p.id, p.quantity_in_stock, q.quantity
            FROM public.products p
            JOIN (SELECT product_id, SUM(quantity - cancelled_quantity)::BIGINT AS quantity
                  FROM public.order_items
                  WHERE order_id = OLD.id AND tenant_id = OLD.tenant_id AND restaurant_id = OLD.restaurant_id
                  GROUP BY product_id HAVING SUM(quantity - cancelled_quantity) > 0) q ON q.product_id = p.id
            WHERE p.tenant_id = OLD.tenant_id AND p.restaurant_id = OLD.restaurant_id AND p.track_inventory
            ORDER BY p.id FOR UPDATE OF p
        LOOP
            IF item.quantity_in_stock IS NULL OR item.quantity <= 0 OR item.quantity_in_stock < item.quantity THEN
                RAISE EXCEPTION 'Insufficient inventory for product %', item.id USING ERRCODE = '23514';
            END IF;
            INSERT INTO public.order_stock_reservations(order_id, tenant_id, restaurant_id, product_id, quantity)
            VALUES (OLD.id, OLD.tenant_id, OLD.restaurant_id, item.id, item.quantity);
            UPDATE public.products SET quantity_in_stock = quantity_in_stock - item.quantity, updated_at = CURRENT_TIMESTAMP
            WHERE id = item.id AND tenant_id = OLD.tenant_id AND restaurant_id = OLD.restaurant_id;
            INSERT INTO public.inventory_movements(tenant_id,restaurant_id,product_id,order_id,reference_type,movement_type,quantity,notes)
            VALUES(OLD.tenant_id,OLD.restaurant_id,item.id,OLD.id,'order_stock','reservation',-item.quantity,'Stock reserved at order confirmation');
        END LOOP;
    ELSIF NEW.status = 'cancelled' AND OLD.status <> 'cancelled' THEN
        FOR item IN
            SELECT p.id, s.quantity - s.released_quantity AS quantity FROM public.order_stock_reservations s
            JOIN public.products p ON p.id = s.product_id AND p.tenant_id = s.tenant_id AND p.restaurant_id = s.restaurant_id
            WHERE s.order_id = OLD.id AND s.tenant_id = OLD.tenant_id AND s.restaurant_id = OLD.restaurant_id AND NOT s.released
            ORDER BY p.id FOR UPDATE OF p
        LOOP
            UPDATE public.products SET quantity_in_stock = quantity_in_stock + item.quantity, updated_at = CURRENT_TIMESTAMP
            WHERE id = item.id AND tenant_id = OLD.tenant_id AND restaurant_id = OLD.restaurant_id;
            INSERT INTO public.inventory_movements(tenant_id,restaurant_id,product_id,order_id,reference_type,movement_type,quantity,notes)
            VALUES(OLD.tenant_id,OLD.restaurant_id,item.id,OLD.id,'order_stock','reservation_release',item.quantity,'Unfulfilled order cancellation');
            UPDATE public.order_stock_reservations SET released = true, released_quantity = quantity, released_at = CURRENT_TIMESTAMP
            WHERE order_id = OLD.id AND product_id = item.id AND tenant_id = OLD.tenant_id AND restaurant_id = OLD.restaurant_id;
        END LOOP;
    END IF;
    RETURN NEW;
END;
$$;
REVOKE ALL ON FUNCTION maintain_order_stock() FROM PUBLIC, pos_app;

CREATE FUNCTION cancel_order_item_quantity(
    p_tenant BIGINT, p_restaurant BIGINT, p_order BIGINT, p_item BIGINT,
    p_quantity INT, p_request UUID, p_actor BIGINT, p_reason TEXT
) RETURNS JSONB LANGUAGE plpgsql SECURITY DEFINER SET search_path = pg_catalog, public AS $$
DECLARE parent RECORD; line RECORD;
        receipt RECORD; reservation RECORD;
        next_subtotal NUMERIC; next_tax NUMERIC; remaining INT; next_remaining BIGINT;
BEGIN
    IF p_quantity IS NULL OR p_quantity <= 0 OR p_reason IS NULL OR length(trim(p_reason)) NOT BETWEEN 1 AND 500 OR p_request IS NULL THEN
        RAISE EXCEPTION 'Invalid cancellation request' USING ERRCODE = '23514';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM public.users WHERE id = p_actor AND tenant_id = p_tenant
                   AND (restaurant_id = p_restaurant OR restaurant_id IS NULL) AND status = 'active') THEN
        RAISE EXCEPTION 'Cancellation actor is outside scope' USING ERRCODE = '42501';
    END IF;
    SELECT id,status,payment_status,discount_amount,subtotal,tax_amount INTO parent FROM public.orders
    WHERE id = p_order AND tenant_id = p_tenant AND restaurant_id = p_restaurant FOR UPDATE;
    IF NOT FOUND THEN RAISE EXCEPTION 'Order not found' USING ERRCODE = 'P0002'; END IF;
    SELECT item_id,quantity,actor_id,reason INTO receipt FROM public.order_item_cancellations
    WHERE tenant_id = p_tenant AND restaurant_id = p_restaurant AND order_id = p_order AND request_key = p_request;
    IF FOUND THEN
        IF receipt.item_id <> p_item OR receipt.quantity <> p_quantity OR receipt.actor_id <> p_actor OR receipt.reason <> p_reason THEN
            RAISE EXCEPTION 'Idempotency key payload conflict' USING ERRCODE = '23514';
        END IF;
        RETURN jsonb_build_object('order_id',p_order,'item_id',p_item,'replayed',true);
    END IF;
    IF parent.status NOT IN ('pending','confirmed') OR parent.payment_status NOT IN ('pending','failed') OR parent.discount_amount <> 0 THEN
        RAISE EXCEPTION 'Partial cancellation requires an unpaid, unfulfilled, undiscounted order' USING ERRCODE = '23514';
    END IF;
    SELECT quantity,cancelled_quantity,product_id INTO line FROM public.order_items
    WHERE id = p_item AND order_id = p_order AND tenant_id = p_tenant AND restaurant_id = p_restaurant FOR UPDATE;
    IF NOT FOUND THEN RAISE EXCEPTION 'Order item not found' USING ERRCODE = 'P0002'; END IF;
    remaining := line.quantity - line.cancelled_quantity;
    IF p_quantity > remaining THEN RAISE EXCEPTION 'Cancellation exceeds remaining quantity' USING ERRCODE = '23514'; END IF;
    SELECT quantity,released_quantity INTO reservation FROM public.order_stock_reservations
    WHERE order_id = p_order AND product_id = line.product_id AND tenant_id = p_tenant AND restaurant_id = p_restaurant FOR UPDATE;
    IF FOUND THEN
        IF reservation.quantity - reservation.released_quantity < p_quantity THEN
            RAISE EXCEPTION 'Cancellation exceeds reserved stock' USING ERRCODE = '23514';
        END IF;
        UPDATE public.products SET quantity_in_stock = quantity_in_stock + p_quantity, updated_at = CURRENT_TIMESTAMP
        WHERE id = line.product_id AND tenant_id = p_tenant AND restaurant_id = p_restaurant;
        INSERT INTO public.inventory_movements(tenant_id,restaurant_id,product_id,order_id,reference_type,movement_type,quantity,notes,created_by)
        VALUES(p_tenant,p_restaurant,line.product_id,p_order,'order_stock','reservation_release',p_quantity,p_reason,p_actor);
        UPDATE public.order_stock_reservations
        SET released_quantity = released_quantity + p_quantity,
            released = released_quantity + p_quantity = quantity,
            released_at = CURRENT_TIMESTAMP
        WHERE order_id = p_order AND product_id = line.product_id AND tenant_id = p_tenant AND restaurant_id = p_restaurant;
    END IF;
    UPDATE public.order_items SET cancelled_quantity = cancelled_quantity + p_quantity,
        total_price = round(total_price * (remaining - p_quantity) / remaining, 2), updated_at = CURRENT_TIMESTAMP
    WHERE id = p_item AND order_id = p_order AND tenant_id = p_tenant AND restaurant_id = p_restaurant;
    SELECT COALESCE(SUM(total_price),0), COALESCE(SUM(quantity-cancelled_quantity),0) INTO next_subtotal,next_remaining FROM public.order_items
    WHERE order_id = p_order AND tenant_id = p_tenant AND restaurant_id = p_restaurant;
    next_tax := CASE WHEN parent.subtotal > 0 THEN round(next_subtotal * parent.tax_amount / parent.subtotal, 2) ELSE 0 END;
    UPDATE public.orders SET subtotal = next_subtotal, tax_amount = next_tax,
        delivery_fee = CASE WHEN next_remaining = 0 THEN 0 ELSE delivery_fee END,
        total_amount = next_subtotal + next_tax + CASE WHEN next_remaining = 0 THEN 0 ELSE delivery_fee END,
        status = CASE WHEN next_remaining = 0 THEN 'cancelled' ELSE status END, updated_at = CURRENT_TIMESTAMP
    WHERE id = p_order AND tenant_id = p_tenant AND restaurant_id = p_restaurant;
    INSERT INTO public.order_item_cancellations(tenant_id,restaurant_id,order_id,item_id,request_key,quantity,reason,actor_id)
    VALUES(p_tenant,p_restaurant,p_order,p_item,p_request,p_quantity,p_reason,p_actor);
    INSERT INTO public.order_status_history(order_id,old_status,new_status,changed_by,change_reason)
    VALUES(p_order,parent.status,CASE WHEN next_remaining = 0 THEN 'cancelled' ELSE parent.status END,p_actor,'Partial item cancellation: ' || p_reason);
    RETURN jsonb_build_object('order_id',p_order,'item_id',p_item,'cancelled_quantity',p_quantity,'replayed',false);
END;
$$;
REVOKE ALL ON FUNCTION cancel_order_item_quantity(BIGINT,BIGINT,BIGINT,BIGINT,INT,UUID,BIGINT,TEXT) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION cancel_order_item_quantity(BIGINT,BIGINT,BIGINT,BIGINT,INT,UUID,BIGINT,TEXT) TO pos_app;

-- Rollback must retain the reservation/cancellation receipts for reconciliation.
-- Do not remove quantity columns or automatically release stock from paid/delivered orders.
-- A reviewed maintenance rollback can revoke EXECUTE on cancel_order_item_quantity;
-- the additive audit schema remains valid for the previous application version.

-- Preserve append-only stock history for the runtime identity.
REVOKE UPDATE,DELETE ON inventory_movements FROM pos_app;
