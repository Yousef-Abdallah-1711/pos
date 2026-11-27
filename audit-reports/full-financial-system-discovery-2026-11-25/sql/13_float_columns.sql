SELECT table_name,column_name,data_type FROM information_schema.columns WHERE table_schema='public' AND data_type IN ('real','double precision','money') ORDER BY 1,2;
