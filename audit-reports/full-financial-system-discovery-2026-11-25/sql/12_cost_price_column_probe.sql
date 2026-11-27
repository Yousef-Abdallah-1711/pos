SELECT COALESCE(SUM((oi.quantity - oi.cancelled_quantity) * COALESCE(p.cost_price, 0)), 0) FROM order_items oi JOIN orders o ON o.id=oi.order_id LEFT JOIN products p ON p.id=oi.product_id;
