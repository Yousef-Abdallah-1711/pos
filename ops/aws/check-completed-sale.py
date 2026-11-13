#!/usr/bin/env python3
import http.cookiejar,json,urllib.request,urllib.error
from pathlib import Path
jar=http.cookiejar.CookieJar();op=urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar));rows=[]
def call(method,path,data=None):
    req=urllib.request.Request('https://api.pos.rakmyat.online/api/v1'+path,data=None if data is None else json.dumps(data).encode(),headers={'Content-Type':'application/json','Origin':'https://pos.app.rakmyat.online'},method=method)
    try:r=op.open(req,timeout=20)
    except urllib.error.HTTPError as e:r=e
    raw=r.read();body=json.loads(raw) if raw else {};return r.status,body.get('data',body)
def record(name,actual,expected):rows.append(dict(check=name,actual=actual,expected=expected,status='PASS' if actual==expected else 'FAIL'))
fixture=json.loads(Path('/opt/servio/ops/aws/evidence/browser-fixtures.json').read_text())[0]
code,_=call('POST','/auth/login',{'email':fixture['email'],'password':fixture['password']});assert code==200
code,body=call('GET','/admin/orders?limit=20');assert code==200
order=next(o for o in body['orders'] if o['customer_name']=='AWS Disposable Test')
order_id=order['id'];path='/admin/orders/'+str(order_id)
code,detail=call('GET',path);assert code==200
product_id=detail['items'][0]['product_id']
if detail['payment_status']!='paid':
    code,_=call('PUT',path+'/payment-status',{'payment_status':'paid'});record('Mark synthetic cash sale paid',code,200)
code,detail=call('GET',path);record('Stored payment status',detail.get('payment_status'),'paid')
statuses=['confirmed','preparing','ready','out_for_delivery','delivered']
current=detail['status'];assert current in statuses
for status in statuses[statuses.index(current)+1:]:
    code,_=call('PUT',path+'/status',{'status':status,'reason':'Disposable AWS lifecycle verification'});record('Order transition '+status,code,200)
code,detail=call('GET',path);record('Completed order status',detail.get('status'),'delivered')
code,product=call('GET','/products/'+str(product_id));record('Tracked inventory after paid completed sale',product.get('quantity_in_stock'),18)
call('POST','/auth/logout',{})
Path('/opt/servio/ops/aws/evidence/completed-sale-inventory.json').write_text(json.dumps({'order_id':order_id,'product_id':product_id,'checks':rows},indent=2)+'\n')
for row in rows:print(row['status']+' '+row['check'])
raise SystemExit(any(row['status']=='FAIL' for row in rows))
