#!/usr/bin/env python3
"""Add idempotent synthetic menu fixtures only to the repository's guarded demo tenant."""
import json
from pathlib import Path
import subprocess

ROOT=Path(__file__).resolve().parents[2]
SQL="""
BEGIN;
DO $$ BEGIN
 IF NOT EXISTS (SELECT 1 FROM tenants t JOIN restaurants r ON r.tenant_id=t.id
   WHERE t.id=1 AND t.name='Test Tenant' AND t.slug='restaurant-1' AND r.id=1 AND r.slug='restaurant-1') THEN
   RAISE EXCEPTION 'Refusing QA fixtures outside the exact seeded demo restaurant';
 END IF;
END $$;
INSERT INTO categories(tenant_id,restaurant_id,name,description,is_active)
SELECT 1,1,'SERVIO QA / اختبار','Synthetic staging browser acceptance only',true
WHERE NOT EXISTS(SELECT 1 FROM categories WHERE tenant_id=1 AND restaurant_id=1 AND name='SERVIO QA / اختبار');
WITH created AS (
 INSERT INTO products(tenant_id,restaurant_id,category_id,name,name_en,name_ar,description_en,description_ar,
   sku,price,cost,available,is_available,status,track_inventory,quantity_in_stock,allergens,available_days)
 SELECT 1,1,c.id,'SERVIO staging QA item','SERVIO staging QA item','منتج اختبار سيرفيو',
   'Synthetic staging acceptance. No real sale.','بيانات اختبار فقط',
   'servio-qa-staging-20261113',25,5,true,true,'active',true,20,'[]','[]'
 FROM categories c WHERE c.tenant_id=1 AND c.restaurant_id=1 AND c.name='SERVIO QA / اختبار'
 AND NOT EXISTS(SELECT 1 FROM products p WHERE p.tenant_id=1 AND p.restaurant_id=1 AND p.sku='servio-qa-staging-20261113')
 RETURNING id,tenant_id,restaurant_id
)
INSERT INTO inventory_movements(tenant_id,restaurant_id,product_id,reference_type,movement_type,quantity,notes)
SELECT tenant_id,restaurant_id,id,'qa_fixture','adjustment',20,'Synthetic staging catalog bootstrap' FROM created;
INSERT INTO restaurant_tables(tenant_id,restaurant_id,table_number,name,capacity,location,status,description,shape)
SELECT 1,1,901,'SERVIO staging QA',4,'indoor','available','Synthetic browser acceptance only','square'
WHERE NOT EXISTS(SELECT 1 FROM restaurant_tables WHERE tenant_id=1 AND restaurant_id=1 AND table_number=901);
SELECT json_build_object('tenant_id',1,'restaurant_id',1,'category_id',p.category_id,'product_id',p.id,
 'price',p.price,'quantity_in_stock',p.quantity_in_stock,'table_id',t.id,'synthetic',true)
FROM products p JOIN restaurant_tables t ON t.tenant_id=p.tenant_id AND t.restaurant_id=p.restaurant_id
WHERE p.tenant_id=1 AND p.restaurant_id=1 AND p.sku='servio-qa-staging-20261113' AND t.table_number=901;
COMMIT;
"""

if __name__=='__main__':
    result=subprocess.run([str(ROOT/'ops/aws/servio'),'exec','-T','postgres','psql','-U','postgres','-d','servio_staging','-qAt','-v','ON_ERROR_STOP=1'],input=SQL,text=True,stdout=subprocess.PIPE,check=True,timeout=30)
    fixture=json.loads(result.stdout.strip())
    assert fixture['synthetic'] and fixture['tenant_id']==1 and fixture['restaurant_id']==1
    target=ROOT/'ops/aws/evidence/storefront-qa-fixture.json'
    target.write_text(json.dumps(fixture,indent=2)+'\n');target.chmod(0o600)
    print('PASS guarded synthetic storefront catalog and table fixture; existing stock was not reset')
