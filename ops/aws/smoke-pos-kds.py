#!/usr/bin/env python3
import base64,http.cookiejar,json,secrets,socket,ssl,struct,time,urllib.request,urllib.error
from pathlib import Path
BASE='https://api.pos.rakmyat.online/api/v1'; ORIGIN='https://pos.app.rakmyat.online'
rows=[]
def check(name,actual,expected):
    expected=expected if isinstance(expected,list) else [expected]
    rows.append(dict(check=name,actual=actual,expected=expected,status='PASS' if actual in expected else 'FAIL'))
def req(op,method,path,body=None,multipart=False):
    hdr={'Origin':ORIGIN,'Content-Type':'application/json'}
    data=None if body is None else json.dumps(body).encode()
    if multipart:
        boundary=secrets.token_hex(16)
        data=(f'--{boundary}\r\nContent-Disposition: form-data; name="data"\r\n\r\n'+json.dumps(body)+f'\r\n--{boundary}--\r\n').encode()
        hdr['Content-Type']='multipart/form-data; boundary='+boundary
    try:r=op.open(urllib.request.Request(BASE+path,data=data,headers=hdr,method=method),timeout=20)
    except urllib.error.HTTPError as e:r=e
    raw=r.read()
    try:parsed=json.loads(raw)
    except ValueError:parsed={}
    if isinstance(parsed,dict):parsed=parsed.get('data',parsed)
    return r.status,parsed

def ws_connect(jar):
    ws=ssl.create_default_context().wrap_socket(socket.create_connection(('api.pos.rakmyat.online',443),timeout=15),server_hostname='api.pos.rakmyat.online')
    key=base64.b64encode(secrets.token_bytes(16)).decode()
    cookies='; '.join(c.name+'='+c.value for c in jar)
    request=f'GET /api/v1/kds/ws HTTP/1.1\r\nHost: api.pos.rakmyat.online\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: {key}\r\nSec-WebSocket-Version: 13\r\nOrigin: {ORIGIN}\r\nCookie: {cookies}\r\n\r\n'
    ws.sendall(request.encode());response=b''
    while b'\r\n\r\n' not in response:response+=ws.recv(1)
    code=int(response.split(b' ')[1]);check('Authenticated WSS through Caddy',code,101)
    if code!=101:ws.close();return None
    ack=ws_event(ws);check('WSS connected acknowledgement',ack.get('type') if ack else None,'connected')
    return ws

def ws_event(ws):
    ws.settimeout(15)
    head=ws.recv(2)
    if len(head)!=2:return None
    length=head[1]&127
    if length==126:length=struct.unpack('!H',ws.recv(2))[0]
    elif length==127:length=struct.unpack('!Q',ws.recv(8))[0]
    data=b''
    while len(data)<length:data+=ws.recv(length-len(data))
    try:return json.loads(data)
    except ValueError:return None

def run():
    fixture=json.loads((Path(__file__).resolve().parents[2]/'ops/aws/evidence/browser-fixtures.json').read_text())[0]
    jar=http.cookiejar.CookieJar();op=urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar))
    code,_=req(op,'POST','/auth/login',{'email':fixture['email'],'password':fixture['password']});check('Flow login',code,200)
    stamp=str(int(time.time()))
    code,category=req(op,'POST','/categories',{'name':'AWS POS flow '+stamp,'is_active':True});check('Flow category',code,201)
    if code!=201:return
    code,product=req(op,'POST','/products',{'category_id':category['id'],'name_en':'AWS disposable POS item','name_ar':'اختبار','sku':'aws-'+stamp,'price':25,'cost':5,'is_available':True,'track_inventory':True,'quantity_in_stock':20,'allergens':[],'available_days':[]},True);check('Create synthetic product',code,201)
    if code!=201:return
    ws=ws_connect(jar)
    code,order=req(op,'POST','/pos/orders',{'customer_name':'AWS Disposable Test','customer_phone':'+201000000001','delivery_address':'Staging only','payment_method':'cash','order_source':'pos','notes':'Disposable AWS acceptance fixture','items':[{'product_id':product['id'],'quantity':2}]});check('Create cash POS order',code,201)
    if code!=201:
        if ws:ws.close()
        return
    check('Server-calculated total with repository 10 percent tax',order.get('total_amount'),55)
    code,current=req(op,'GET','/products/'+str(product['id']));check('Inventory after order',current.get('quantity_in_stock'),18)
    code,queue=req(op,'GET','/kds/queue');check('KDS queue HTTP',code,200)
    ticket=next((q for q in (queue or []) if q.get('order_id')==order['id']),None)
    check('POS ticket appears in KDS',ticket is not None,True)
    if ws:
        event=ws_event(ws);check('WSS order_new event received',event.get('type') if event else None,'order_new')
        # Close only the synthetic client connection.
        ws.close()
    if ticket:
        for step in range(2):
            code,_=req(op,'PUT','/kds/orders/'+str(ticket['id'])+'/bump',{});check('KDS bump '+str(step+1),code,200)
    code,detail=req(op,'GET','/admin/orders/'+str(order['id']));check('Read own POS order',code,200)
    code,_=req(op,'POST','/auth/logout',{});check('Flow logout',code,[200,204])
try:run()
except Exception as e:rows.append(dict(check='POS/KDS flow execution',status='FAIL',error=str(e)))
(Path(__file__).resolve().parents[2]/'ops/aws/evidence/pos-kds-smoke.json').write_text(json.dumps(rows,indent=2)+'\n')
for row in rows:print(row['status']+' '+row['check'])
raise SystemExit(any(r['status']=='FAIL' for r in rows))
