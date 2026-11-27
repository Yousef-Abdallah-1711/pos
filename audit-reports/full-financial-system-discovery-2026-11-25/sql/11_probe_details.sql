SELECT 'p31', tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance FROM chart_of_accounts WHERE (account_type IN ('asset','expense') AND normal_balance<>'debit') OR (account_type IN ('liability','equity','revenue') AND normal_balance<>'credit');
SELECT 'p33', r.id, r.tenant_id, r.name FROM restaurants r WHERE NOT EXISTS (SELECT 1 FROM chart_of_accounts c WHERE c.restaurant_id=r.id);
SELECT 'coa_per_restaurant', tenant_id, restaurant_id, count(*) FROM chart_of_accounts GROUP BY 2,3;
SELECT 'tenants', count(*) FROM tenants; SELECT 'restaurants', count(*) FROM restaurants;
