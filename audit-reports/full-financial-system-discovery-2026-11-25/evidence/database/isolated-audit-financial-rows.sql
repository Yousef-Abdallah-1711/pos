--
-- PostgreSQL database dump
--

\restrict Rut16uWqioy1Nm2KLaD0DgQxc6JEWENpvwRyoNncgkwbeylzArcDtLVTrjM8th1

-- Dumped from database version 15.19
-- Dumped by pg_dump version 15.19

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.orders VALUES (1, 2, 2, 'ORD-2026-5f61f71a-a2e0-403a-9f73-4fce8a40ce0d', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'cash', 'paid', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:27:01.093643', '2026-11-25 22:27:01.181622', NULL, 'cc79f24394c79bd0f7a0396ce44ba31627be2a37ae3f579fb623b6d5c2081070');
INSERT INTO public.orders VALUES (2, 4, 4, 'ORD-2026-951d28ed-cf22-4ff5-bd2d-5d2091d0dae1', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'cash', 'paid', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:28:45.214093', '2026-11-25 22:28:45.391446', NULL, '9708604ca5c70be04080c89046ee269e4d8323094ddcf35fa019048bfd591e1f');
INSERT INTO public.orders VALUES (4, 4, 4, 'ORD-2026-8188c639-6477-4316-ab5c-69a029c3fc3a', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'wallet', 'pending', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:41:16.27453', '2026-11-25 22:41:16.311981', NULL, 'dac9d48d03f645b270982a61e262f32617413fe0053a99946b29dd9e17755d88');
INSERT INTO public.orders VALUES (12, 4, 4, 'ORD-2026-5bde08b4-2db6-453c-84d7-6425ff26a7cf', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'cash', 'pending', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:46:31.788601', '2026-11-25 22:46:31.800803', NULL, '8ef1886ff657bc12f37640323080660d1f7471067518e68c15cc242382211f27');
INSERT INTO public.orders VALUES (5, 4, 4, 'ORD-2026-ae041f8b-a711-4713-bbf5-9b69603824c6', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'cash', 'paid', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:46:25.455945', '2026-11-25 22:59:30.111041', NULL, 'cc4d133c7200e68e437efcd3555229533586d5ddfbe15563a2b5b2e00d2aa911');
INSERT INTO public.orders VALUES (3, 4, 4, 'ORD-2026-863e52a9-c776-494b-b453-6d34ee4e402a', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'card', 'paid', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:41:03.166192', '2026-11-25 22:45:43.86414', NULL, 'eec02ef275fb08222c8576c3ecdad6013761b7ca2775cbdfe4b59adb7d403d43');
INSERT INTO public.orders VALUES (6, 4, 4, 'ORD-2026-8a9b1ed5-6560-46a6-a966-8518e09df8e8', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'cash', 'pending', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:46:25.455973', '2026-11-25 22:46:25.478393', NULL, 'd12f2da1ffe41a7ede8f689855d1044602a44ef4435444c4b1330453678a64ee');
INSERT INTO public.orders VALUES (7, 4, 4, 'ORD-2026-c9bfcab3-87bf-42e7-be05-af732d30bcfd', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'cash', 'pending', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:46:25.463445', '2026-11-25 22:46:25.481796', NULL, '7da1c512b899d5d36a1713b3ce1ff63c6613aa9dba95d0bdc854bfb0a1133acb');
INSERT INTO public.orders VALUES (9, 4, 4, 'ORD-2026-b6d27ad4-4762-4a4c-9617-4e86b82f2cc9', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'cash', 'pending', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:46:25.540528', '2026-11-25 22:46:25.565557', NULL, 'a36facd636112854a9565b22ef216eb200388d2080f9772469017b94042b0395');
INSERT INTO public.orders VALUES (8, 4, 4, 'ORD-2026-23c1c796-f71b-4d76-a6c4-b487b5239996', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'cash', 'pending', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:46:25.531573', '2026-11-25 22:46:25.564803', NULL, '51d1a5ff2cc3b7a4cb7e009d562009efca8bb915815b8d2c475e5fef0820ad8a');
INSERT INTO public.orders VALUES (10, 4, 4, 'ORD-2026-5b4b441f-1e80-48e9-9ab9-c12ed61da71b', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'cash', 'pending', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:46:31.558551', '2026-11-25 22:46:31.574215', NULL, 'e23cfd049c491a0caeeb36e6749c389b522ac9a57413f79e9749e874044dbd4a');
INSERT INTO public.orders VALUES (11, 4, 4, 'ORD-2026-759cee42-7ae9-4bdc-9d3f-6c7600f196f1', 'Walk-in', NULL, '0000000000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 10000.00, 1000.00, 0.00, 0.00, 11000.00, 'cash', 'pending', 'confirmed', NULL, NULL, '', 'pos', '2026-11-25 22:46:31.667136', '2026-11-25 22:46:31.681106', NULL, '7a41744fa210b08e825414a019bb2a17d357513cda80285bebd20b8f7627d7ca');


--
-- Data for Name: accounting_event_outbox; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.accounting_event_outbox VALUES (2, 'paid-order:4:4:2', 'paid_order', 4, 4, 2, 4, 'completed', 1, 8, '2026-11-25 22:28:45.391446', NULL, NULL, '2026-11-25 22:28:45.391446', '2026-11-25 22:28:55.378759', '2026-11-25 22:28:55.378759', 1);
INSERT INTO public.accounting_event_outbox VALUES (3, 'paid-order:4:4:3', 'paid_order', 4, 4, 3, 4, 'completed', 1, 8, '2026-11-25 22:45:43.45249', NULL, NULL, '2026-11-25 22:45:43.45249', '2026-11-25 22:45:55.387383', '2026-11-25 22:45:55.387383', 1);
INSERT INTO public.accounting_event_outbox VALUES (1, 'paid-order:2:2:1', 'paid_order', 2, 2, 1, 2, 'completed', 6, 8, '2026-11-25 22:58:54.85484', NULL, NULL, '2026-11-25 22:27:01.181622', '2026-11-25 22:59:09.813216', '2026-11-25 22:59:09.813216', 6);
INSERT INTO public.accounting_event_outbox VALUES (23, 'paid-order:4:4:5', 'paid_order', 4, 4, 5, 4, 'completed', 1, 8, '2026-11-25 22:59:30.111041', NULL, NULL, '2026-11-25 22:59:30.111041', '2026-11-25 22:59:35.449617', '2026-11-25 22:59:35.449617', 1);


--
-- Data for Name: journal_entries; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.journal_entries VALUES (8, 4, 4, 'JE-2026-00008', '2026-11-25', '2026-11', 'Auto-posted sales for order ORD-2026-863e52a9-c776-494b-b453-6d34ee4e402a', '', 'order', 3, 'ORD-2026-863e52a9-c776-494b-b453-6d34ee4e402a', 'auto_sales', 11000.00, 11000.00, 'posted', '2026-11-25 22:45:54.933512', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:45:54.933512', '2026-11-25 22:45:54.933512', 4, 4);
INSERT INTO public.journal_entries VALUES (1, 4, 4, 'JE-2026-00001', '2026-11-25', '2026-11', 'Auto-posted sales for order ORD-2026-951d28ed-cf22-4ff5-bd2d-5d2091d0dae1', '', 'order', 2, 'ORD-2026-951d28ed-cf22-4ff5-bd2d-5d2091d0dae1', 'auto_sales', 11000.00, 11000.00, 'posted', '2026-11-25 22:28:54.977171', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:28:54.977171', '2026-11-25 22:28:54.977171', 4, 4);
INSERT INTO public.journal_entries VALUES (25, 4, 4, 'JE-2026-00017', '2026-11-25', '2026-11', 'Customer payment PMT-R-2026-00003', '', 'payment_received', 3, 'PMT-R-2026-00003', 'payment_received', 1000.00, 1000.00, 'posted', '2026-11-25 22:50:09.732584', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:50:09.732584', '2026-11-25 22:50:09.732584', 4, 4);
INSERT INTO public.journal_entries VALUES (9, 4, 4, 'JE-2026-00009', '2026-11-25', '2026-11', 'Auto-posted COGS for order ORD-2026-863e52a9-c776-494b-b453-6d34ee4e402a', '', 'order', 3, 'ORD-2026-863e52a9-c776-494b-b453-6d34ee4e402a', 'auto_cogs', 4000.00, 4000.00, 'posted', '2026-11-25 22:45:55.245533', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:45:55.245533', '2026-11-25 22:45:55.245533', 4, 4);
INSERT INTO public.journal_entries VALUES (2, 4, 4, 'JE-2026-00002', '2026-11-25', '2026-11', 'Auto-posted COGS for order ORD-2026-951d28ed-cf22-4ff5-bd2d-5d2091d0dae1', '', 'order', 2, 'ORD-2026-951d28ed-cf22-4ff5-bd2d-5d2091d0dae1', 'auto_cogs', 4000.00, 4000.00, 'posted', '2026-11-25 22:28:55.270987', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:28:55.270987', '2026-11-25 22:28:55.270987', 4, 4);
INSERT INTO public.journal_entries VALUES (10, 4, 4, 'JE-2026-00010', '2026-11-25', '2026-11', 'hdr', '', '', NULL, '', 'manual', 10.00, 10.00, 'draft', NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:47:01.692684', '2026-11-25 22:47:01.692684', 4, 4);
INSERT INTO public.journal_entries VALUES (4, 4, 4, 'JE-2026-00004', '2026-11-25', '2026-11', 'Rent paid in cash', '', '', NULL, '', 'manual', 3000.00, 3000.00, 'posted', '2026-11-25 22:31:28.729903', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:31:19.048836', '2026-11-25 22:31:28.729903', 4, 4);
INSERT INTO public.journal_entries VALUES (5, 4, 4, 'JE-2026-00005', '2026-11-25', '2026-11', 'Supplier invoice PINV-2026-00001', '', 'supplier_invoice', 1, 'PINV-2026-00001', 'supplier_invoice', 22800.00, 22800.00, 'posted', '2026-11-25 22:33:16.786267', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:33:16.786267', '2026-11-25 22:33:16.786267', 4, 4);
INSERT INTO public.journal_entries VALUES (6, 4, 4, 'JE-2026-00006', '2026-11-25', '2026-11', 'Supplier payment PMT-M-2026-00001', '', 'payment_made', 1, 'PMT-M-2026-00001', 'payment_made', 10000.00, 10000.00, 'posted', '2026-11-25 22:33:49.054066', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:33:49.054066', '2026-11-25 22:33:49.054066', 4, 4);
INSERT INTO public.journal_entries VALUES (16, 4, 4, 'JE-2026-00011', '2026-11-25', '2026-11', 'rounding', '', '', NULL, '', 'manual', 0.01, 0.01, 'draft', NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:47:14.83608', '2026-11-25 22:47:14.83608', 4, 4);
INSERT INTO public.journal_entries VALUES (7, 4, 4, 'JE-2026-00007', '2026-11-25', '2026-11', 'Auto-posted payroll for salary #1', '', 'salary', 1, 'SAL-1', 'auto_payroll', 5000.00, 5000.00, 'posted', '2026-11-25 22:37:56.655723', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:37:56.655723', '2026-11-25 22:37:56.655723', 4, 4);
INSERT INTO public.journal_entries VALUES (28, 4, 4, 'JE-2026-00020', '2026-11-25', '2026-11', 'Customer payment PMT-R-2026-00006', '', 'payment_received', 6, 'PMT-R-2026-00006', 'payment_received', 1000.00, 1000.00, 'posted', '2026-11-25 22:50:09.79576', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:50:09.79576', '2026-11-25 22:50:09.79576', 4, 4);
INSERT INTO public.journal_entries VALUES (26, 4, 4, 'JE-2026-00018', '2026-11-25', '2026-11', 'Customer payment PMT-R-2026-00004', '', 'payment_received', 4, 'PMT-R-2026-00004', 'payment_received', 1000.00, 1000.00, 'posted', '2026-11-25 22:50:09.763848', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:50:09.763848', '2026-11-25 22:50:09.763848', 4, 4);
INSERT INTO public.journal_entries VALUES (17, 4, 4, 'JE-2026-00012', '2026-11-25', '2026-11', 'rounding2', '', '', NULL, '', 'manual', 10.01, 10.01, 'draft', NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:47:14.981012', '2026-11-25 22:47:14.981012', 4, 4);
INSERT INTO public.journal_entries VALUES (22, 4, 4, 'JE-2026-00014', '2026-11-25', '2026-11', 'Customer invoice INV-2026-00001', '', 'customer_invoice', 2, 'INV-2026-00001', 'customer_invoice', 5200.00, 5200.00, 'posted', '2026-11-25 22:49:27.574748', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:49:27.574748', '2026-11-25 22:49:27.574748', 4, 4);
INSERT INTO public.journal_entries VALUES (18, 4, 4, 'JE-2027-00001', '2027-01-01', '2027-01', 'future', '', '', NULL, '', 'manual', 10.00, 10.00, 'draft', NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:47:15.144187', '2026-11-25 22:47:15.144187', 4, 4);
INSERT INTO public.journal_entries VALUES (19, 4, 4, 'JE-2020-00001', '2020-01-01', '2020-01', 'backdated', '', '', NULL, '', 'manual', 10.00, 10.00, 'draft', NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:47:15.298442', '2026-11-25 22:47:15.298442', 4, 4);
INSERT INTO public.journal_entries VALUES (20, 4, 4, 'JE-2026-00013', '2026-11-25', '2026-11', 'closed-period test', '', '', NULL, '', 'manual', 5.00, 5.00, 'draft', NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:47:38.548455', '2026-11-25 22:47:38.548455', 4, 4);
INSERT INTO public.journal_entries VALUES (3, 4, 4, 'JE-2026-00003', '2026-11-25', '2026-11', 'Capital contribution by owner', '', '', NULL, '', 'manual', 100000.00, 100000.00, 'void', '2026-11-25 22:31:03.921079', 4, '2026-11-25 22:47:38.696796', 4, 'audit closed period void test', false, NULL, NULL, '', NULL, '2026-11-25 22:29:54.515254', '2026-11-25 22:47:38.695827', 4, 4);
INSERT INTO public.journal_entries VALUES (31, 4, 4, 'JE-2026-00023', '2026-11-25', '2026-11', 'Supplier payment PMT-M-2026-00004', '', 'payment_made', 4, 'PMT-M-2026-00004', 'payment_made', 3000.00, 3000.00, 'posted', '2026-11-25 22:50:23.56892', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:50:23.56892', '2026-11-25 22:50:23.56892', 4, 4);
INSERT INTO public.journal_entries VALUES (23, 4, 4, 'JE-2026-00015', '2026-11-25', '2026-11', 'Customer payment PMT-R-2026-00001', '', 'payment_received', 1, 'PMT-R-2026-00001', 'payment_received', 2000.00, 2000.00, 'posted', '2026-11-25 22:49:50.010227', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:49:50.010227', '2026-11-25 22:49:50.010227', 4, 4);
INSERT INTO public.journal_entries VALUES (29, 4, 4, 'JE-2026-00021', '2026-11-25', '2026-11', 'Supplier payment PMT-M-2026-00002', '', 'payment_made', 2, 'PMT-M-2026-00002', 'payment_made', 3000.00, 3000.00, 'posted', '2026-11-25 22:50:23.524633', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:50:23.524633', '2026-11-25 22:50:23.524633', 4, 4);
INSERT INTO public.journal_entries VALUES (24, 4, 4, 'JE-2026-00016', '2026-11-25', '2026-11', 'Customer payment PMT-R-2026-00002', '', 'payment_received', 2, 'PMT-R-2026-00002', 'payment_received', 1000.00, 1000.00, 'posted', '2026-11-25 22:50:09.706039', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:50:09.706039', '2026-11-25 22:50:09.706039', 4, 4);
INSERT INTO public.journal_entries VALUES (27, 4, 4, 'JE-2026-00019', '2026-11-25', '2026-11', 'Customer payment PMT-R-2026-00005', '', 'payment_received', 5, 'PMT-R-2026-00005', 'payment_received', 1000.00, 1000.00, 'posted', '2026-11-25 22:50:09.790594', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:50:09.790594', '2026-11-25 22:50:09.790594', 4, 4);
INSERT INTO public.journal_entries VALUES (34, 2, 2, 'JE-2026-00002', '2026-11-25', '2026-11', 'Auto-posted COGS for order ORD-2026-5f61f71a-a2e0-403a-9f73-4fce8a40ce0d', '', 'order', 1, 'ORD-2026-5f61f71a-a2e0-403a-9f73-4fce8a40ce0d', 'auto_cogs', 4000.00, 4000.00, 'posted', '2026-11-25 22:59:09.782036', 2, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:59:09.782036', '2026-11-25 22:59:09.782036', 2, 2);
INSERT INTO public.journal_entries VALUES (30, 4, 4, 'JE-2026-00022', '2026-11-25', '2026-11', 'Supplier payment PMT-M-2026-00003', '', 'payment_made', 3, 'PMT-M-2026-00003', 'payment_made', 3000.00, 3000.00, 'posted', '2026-11-25 22:50:23.544114', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:50:23.544114', '2026-11-25 22:50:23.544114', 4, 4);
INSERT INTO public.journal_entries VALUES (33, 2, 2, 'JE-2026-00001', '2026-11-25', '2026-11', 'Auto-posted sales for order ORD-2026-5f61f71a-a2e0-403a-9f73-4fce8a40ce0d', '', 'order', 1, 'ORD-2026-5f61f71a-a2e0-403a-9f73-4fce8a40ce0d', 'auto_sales', 11000.00, 11000.00, 'posted', '2026-11-25 22:59:09.694097', 2, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:59:09.694097', '2026-11-25 22:59:09.694097', 2, 2);
INSERT INTO public.journal_entries VALUES (32, 4, 4, 'JE-2026-00024', '2026-11-25', '2026-11', 'Supplier payment PMT-M-2026-00005', '', 'payment_made', 5, 'PMT-M-2026-00005', 'payment_made', 3000.00, 3000.00, 'posted', '2026-11-25 22:50:23.595548', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:50:23.595548', '2026-11-25 22:50:23.595548', 4, 4);
INSERT INTO public.journal_entries VALUES (35, 4, 4, 'JE-2026-00025', '2026-11-25', '2026-11', 'Auto-posted sales for order ORD-2026-ae041f8b-a711-4713-bbf5-9b69603824c6', '', 'order', 5, 'ORD-2026-ae041f8b-a711-4713-bbf5-9b69603824c6', 'auto_sales', 11000.00, 11000.00, 'posted', '2026-11-25 22:59:35.321857', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:59:35.321857', '2026-11-25 22:59:35.321857', 4, 4);
INSERT INTO public.journal_entries VALUES (36, 4, 4, 'JE-2026-00026', '2026-11-25', '2026-11', 'Auto-posted COGS for order ORD-2026-ae041f8b-a711-4713-bbf5-9b69603824c6', '', 'order', 5, 'ORD-2026-ae041f8b-a711-4713-bbf5-9b69603824c6', 'auto_cogs', 4000.00, 4000.00, 'posted', '2026-11-25 22:59:35.408398', 4, NULL, NULL, NULL, false, NULL, NULL, '', NULL, '2026-11-25 22:59:35.408398', '2026-11-25 22:59:35.408398', 4, 4);


--
-- Data for Name: journal_entry_lines; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.journal_entry_lines VALUES (1, 1, 4, 70, '1110', 11000.00, 0.00, 'Order payment ORD-2026-951d28ed-cf22-4ff5-bd2d-5d2091d0dae1', '', 1, '', 0.00, '2026-11-25 22:28:54.977171');
INSERT INTO public.journal_entry_lines VALUES (2, 1, 4, 96, '4100', 0.00, 10000.00, 'Sales revenue ORD-2026-951d28ed-cf22-4ff5-bd2d-5d2091d0dae1', '', 2, '', 0.00, '2026-11-25 22:28:54.977171');
INSERT INTO public.journal_entry_lines VALUES (3, 1, 4, 87, '2120', 0.00, 1000.00, 'Sales tax ORD-2026-951d28ed-cf22-4ff5-bd2d-5d2091d0dae1', '', 3, '', 0.00, '2026-11-25 22:28:54.977171');
INSERT INTO public.journal_entry_lines VALUES (4, 2, 4, 103, '5110', 4000.00, 0.00, 'COGS for order ORD-2026-951d28ed-cf22-4ff5-bd2d-5d2091d0dae1', '', 1, '', 0.00, '2026-11-25 22:28:55.270987');
INSERT INTO public.journal_entry_lines VALUES (5, 2, 4, 76, '1130', 0.00, 4000.00, 'Inventory reduction for order ORD-2026-951d28ed-cf22-4ff5-bd2d-5d2091d0dae1', '', 2, '', 0.00, '2026-11-25 22:28:55.270987');
INSERT INTO public.journal_entry_lines VALUES (6, 3, 4, 70, '1110', 100000.00, 0.00, '', '', 1, '', 0.00, '2026-11-25 22:29:54.515254');
INSERT INTO public.journal_entry_lines VALUES (7, 3, 4, 93, '3100', 0.00, 100000.00, '', '', 2, '', 0.00, '2026-11-25 22:29:54.515254');
INSERT INTO public.journal_entry_lines VALUES (8, 4, 4, 111, '5310', 3000.00, 0.00, '', '', 1, '', 0.00, '2026-11-25 22:31:19.048836');
INSERT INTO public.journal_entry_lines VALUES (9, 4, 4, 70, '1110', 0.00, 3000.00, '', '', 2, '', 0.00, '2026-11-25 22:31:19.048836');
INSERT INTO public.journal_entry_lines VALUES (10, 5, 4, 76, '1130', 20000.00, 0.00, 'Supplier invoice expense PINV-2026-00001', '', 1, '', 0.00, '2026-11-25 22:33:16.786267');
INSERT INTO public.journal_entry_lines VALUES (11, 5, 4, 78, '1150', 2800.00, 0.00, 'Input VAT for supplier invoice PINV-2026-00001', '', 2, '', 0.00, '2026-11-25 22:33:16.786267');
INSERT INTO public.journal_entry_lines VALUES (12, 5, 4, 85, '2101', 0.00, 22800.00, 'Accounts payable for supplier invoice PINV-2026-00001', '', 3, '', 0.00, '2026-11-25 22:33:16.786267');
INSERT INTO public.journal_entry_lines VALUES (13, 6, 4, 85, '2101', 10000.00, 0.00, 'Accounts payable settlement PMT-M-2026-00001', '', 1, '', 0.00, '2026-11-25 22:33:49.054066');
INSERT INTO public.journal_entry_lines VALUES (14, 6, 4, 70, '1110', 0.00, 10000.00, 'Cash/bank payment PMT-M-2026-00001', '', 2, '', 0.00, '2026-11-25 22:33:49.054066');
INSERT INTO public.journal_entry_lines VALUES (15, 7, 4, 107, '5210', 5000.00, 0.00, 'Payroll expense salary #1', '', 1, '', 0.00, '2026-11-25 22:37:56.655723');
INSERT INTO public.journal_entry_lines VALUES (16, 7, 4, 71, '1111', 0.00, 5000.00, 'Salary payment #1', '', 2, '', 0.00, '2026-11-25 22:37:56.655723');
INSERT INTO public.journal_entry_lines VALUES (17, 8, 4, 73, '1113', 11000.00, 0.00, 'Order payment ORD-2026-863e52a9-c776-494b-b453-6d34ee4e402a', '', 1, '', 0.00, '2026-11-25 22:45:54.933512');
INSERT INTO public.journal_entry_lines VALUES (18, 8, 4, 96, '4100', 0.00, 10000.00, 'Sales revenue ORD-2026-863e52a9-c776-494b-b453-6d34ee4e402a', '', 2, '', 0.00, '2026-11-25 22:45:54.933512');
INSERT INTO public.journal_entry_lines VALUES (19, 8, 4, 87, '2120', 0.00, 1000.00, 'Sales tax ORD-2026-863e52a9-c776-494b-b453-6d34ee4e402a', '', 3, '', 0.00, '2026-11-25 22:45:54.933512');
INSERT INTO public.journal_entry_lines VALUES (20, 9, 4, 103, '5110', 4000.00, 0.00, 'COGS for order ORD-2026-863e52a9-c776-494b-b453-6d34ee4e402a', '', 1, '', 0.00, '2026-11-25 22:45:55.245533');
INSERT INTO public.journal_entry_lines VALUES (21, 9, 4, 76, '1130', 0.00, 4000.00, 'Inventory reduction for order ORD-2026-863e52a9-c776-494b-b453-6d34ee4e402a', '', 2, '', 0.00, '2026-11-25 22:45:55.245533');
INSERT INTO public.journal_entry_lines VALUES (22, 10, 4, 70, '1110', 10.00, 0.00, '', '', 1, '', 0.00, '2026-11-25 22:47:01.692684');
INSERT INTO public.journal_entry_lines VALUES (23, 10, 4, 93, '3100', 0.00, 10.00, '', '', 2, '', 0.00, '2026-11-25 22:47:01.692684');
INSERT INTO public.journal_entry_lines VALUES (24, 16, 4, 70, '1110', 0.01, 0.00, '', '', 1, '', 0.00, '2026-11-25 22:47:14.83608');
INSERT INTO public.journal_entry_lines VALUES (25, 16, 4, 93, '3100', 0.00, 0.01, '', '', 2, '', 0.00, '2026-11-25 22:47:14.83608');
INSERT INTO public.journal_entry_lines VALUES (26, 17, 4, 70, '1110', 10.01, 0.00, '', '', 1, '', 0.00, '2026-11-25 22:47:14.981012');
INSERT INTO public.journal_entry_lines VALUES (27, 17, 4, 93, '3100', 0.00, 10.01, '', '', 2, '', 0.00, '2026-11-25 22:47:14.981012');
INSERT INTO public.journal_entry_lines VALUES (28, 18, 4, 70, '1110', 10.00, 0.00, '', '', 1, '', 0.00, '2026-11-25 22:47:15.144187');
INSERT INTO public.journal_entry_lines VALUES (29, 18, 4, 93, '3100', 0.00, 10.00, '', '', 2, '', 0.00, '2026-11-25 22:47:15.144187');
INSERT INTO public.journal_entry_lines VALUES (30, 19, 4, 70, '1110', 10.00, 0.00, '', '', 1, '', 0.00, '2026-11-25 22:47:15.298442');
INSERT INTO public.journal_entry_lines VALUES (31, 19, 4, 93, '3100', 0.00, 10.00, '', '', 2, '', 0.00, '2026-11-25 22:47:15.298442');
INSERT INTO public.journal_entry_lines VALUES (32, 20, 4, 70, '1110', 5.00, 0.00, '', '', 1, '', 0.00, '2026-11-25 22:47:38.548455');
INSERT INTO public.journal_entry_lines VALUES (33, 20, 4, 93, '3100', 0.00, 5.00, '', '', 2, '', 0.00, '2026-11-25 22:47:38.548455');
INSERT INTO public.journal_entry_lines VALUES (37, 22, 4, 74, '1120', 5200.00, 0.00, 'Accounts receivable for invoice INV-2026-00001', '', 1, '', 0.00, '2026-11-25 22:49:27.574748');
INSERT INTO public.journal_entry_lines VALUES (38, 22, 4, 96, '4100', 0.00, 4500.00, 'Revenue recognition for invoice INV-2026-00001', '', 2, '', 0.00, '2026-11-25 22:49:27.574748');
INSERT INTO public.journal_entry_lines VALUES (39, 22, 4, 87, '2120', 0.00, 700.00, 'Output VAT for invoice INV-2026-00001', '', 3, '', 0.00, '2026-11-25 22:49:27.574748');
INSERT INTO public.journal_entry_lines VALUES (40, 23, 4, 70, '1110', 2000.00, 0.00, 'Collection for PMT-R-2026-00001', '', 1, '', 0.00, '2026-11-25 22:49:50.010227');
INSERT INTO public.journal_entry_lines VALUES (41, 23, 4, 74, '1120', 0.00, 2000.00, 'Receivable settlement for PMT-R-2026-00001', '', 2, '', 0.00, '2026-11-25 22:49:50.010227');
INSERT INTO public.journal_entry_lines VALUES (42, 24, 4, 70, '1110', 1000.00, 0.00, 'Collection for PMT-R-2026-00002', '', 1, '', 0.00, '2026-11-25 22:50:09.706039');
INSERT INTO public.journal_entry_lines VALUES (43, 24, 4, 74, '1120', 0.00, 1000.00, 'Receivable settlement for PMT-R-2026-00002', '', 2, '', 0.00, '2026-11-25 22:50:09.706039');
INSERT INTO public.journal_entry_lines VALUES (44, 25, 4, 70, '1110', 1000.00, 0.00, 'Collection for PMT-R-2026-00003', '', 1, '', 0.00, '2026-11-25 22:50:09.732584');
INSERT INTO public.journal_entry_lines VALUES (45, 25, 4, 74, '1120', 0.00, 1000.00, 'Receivable settlement for PMT-R-2026-00003', '', 2, '', 0.00, '2026-11-25 22:50:09.732584');
INSERT INTO public.journal_entry_lines VALUES (46, 26, 4, 70, '1110', 1000.00, 0.00, 'Collection for PMT-R-2026-00004', '', 1, '', 0.00, '2026-11-25 22:50:09.763848');
INSERT INTO public.journal_entry_lines VALUES (47, 26, 4, 74, '1120', 0.00, 1000.00, 'Receivable settlement for PMT-R-2026-00004', '', 2, '', 0.00, '2026-11-25 22:50:09.763848');
INSERT INTO public.journal_entry_lines VALUES (48, 27, 4, 70, '1110', 1000.00, 0.00, 'Collection for PMT-R-2026-00005', '', 1, '', 0.00, '2026-11-25 22:50:09.790594');
INSERT INTO public.journal_entry_lines VALUES (49, 27, 4, 74, '1120', 0.00, 1000.00, 'Receivable settlement for PMT-R-2026-00005', '', 2, '', 0.00, '2026-11-25 22:50:09.790594');
INSERT INTO public.journal_entry_lines VALUES (50, 28, 4, 70, '1110', 1000.00, 0.00, 'Collection for PMT-R-2026-00006', '', 1, '', 0.00, '2026-11-25 22:50:09.79576');
INSERT INTO public.journal_entry_lines VALUES (51, 28, 4, 74, '1120', 0.00, 1000.00, 'Receivable settlement for PMT-R-2026-00006', '', 2, '', 0.00, '2026-11-25 22:50:09.79576');
INSERT INTO public.journal_entry_lines VALUES (52, 29, 4, 85, '2101', 3000.00, 0.00, 'Accounts payable settlement PMT-M-2026-00002', '', 1, '', 0.00, '2026-11-25 22:50:23.524633');
INSERT INTO public.journal_entry_lines VALUES (53, 29, 4, 70, '1110', 0.00, 3000.00, 'Cash/bank payment PMT-M-2026-00002', '', 2, '', 0.00, '2026-11-25 22:50:23.524633');
INSERT INTO public.journal_entry_lines VALUES (54, 30, 4, 85, '2101', 3000.00, 0.00, 'Accounts payable settlement PMT-M-2026-00003', '', 1, '', 0.00, '2026-11-25 22:50:23.544114');
INSERT INTO public.journal_entry_lines VALUES (55, 30, 4, 70, '1110', 0.00, 3000.00, 'Cash/bank payment PMT-M-2026-00003', '', 2, '', 0.00, '2026-11-25 22:50:23.544114');
INSERT INTO public.journal_entry_lines VALUES (56, 31, 4, 85, '2101', 3000.00, 0.00, 'Accounts payable settlement PMT-M-2026-00004', '', 1, '', 0.00, '2026-11-25 22:50:23.56892');
INSERT INTO public.journal_entry_lines VALUES (57, 31, 4, 70, '1110', 0.00, 3000.00, 'Cash/bank payment PMT-M-2026-00004', '', 2, '', 0.00, '2026-11-25 22:50:23.56892');
INSERT INTO public.journal_entry_lines VALUES (58, 32, 4, 85, '2101', 3000.00, 0.00, 'Accounts payable settlement PMT-M-2026-00005', '', 1, '', 0.00, '2026-11-25 22:50:23.595548');
INSERT INTO public.journal_entry_lines VALUES (59, 32, 4, 70, '1110', 0.00, 3000.00, 'Cash/bank payment PMT-M-2026-00005', '', 2, '', 0.00, '2026-11-25 22:50:23.595548');
INSERT INTO public.journal_entry_lines VALUES (60, 33, 2, 132, '1110', 11000.00, 0.00, 'Order payment ORD-2026-5f61f71a-a2e0-403a-9f73-4fce8a40ce0d', '', 1, '', 0.00, '2026-11-25 22:59:09.694097');
INSERT INTO public.journal_entry_lines VALUES (61, 33, 2, 158, '4100', 0.00, 10000.00, 'Sales revenue ORD-2026-5f61f71a-a2e0-403a-9f73-4fce8a40ce0d', '', 2, '', 0.00, '2026-11-25 22:59:09.694097');
INSERT INTO public.journal_entry_lines VALUES (62, 33, 2, 149, '2120', 0.00, 1000.00, 'Sales tax ORD-2026-5f61f71a-a2e0-403a-9f73-4fce8a40ce0d', '', 3, '', 0.00, '2026-11-25 22:59:09.694097');
INSERT INTO public.journal_entry_lines VALUES (63, 34, 2, 165, '5110', 4000.00, 0.00, 'COGS for order ORD-2026-5f61f71a-a2e0-403a-9f73-4fce8a40ce0d', '', 1, '', 0.00, '2026-11-25 22:59:09.782036');
INSERT INTO public.journal_entry_lines VALUES (64, 34, 2, 138, '1130', 0.00, 4000.00, 'Inventory reduction for order ORD-2026-5f61f71a-a2e0-403a-9f73-4fce8a40ce0d', '', 2, '', 0.00, '2026-11-25 22:59:09.782036');
INSERT INTO public.journal_entry_lines VALUES (65, 35, 4, 70, '1110', 11000.00, 0.00, 'Order payment ORD-2026-ae041f8b-a711-4713-bbf5-9b69603824c6', '', 1, '', 0.00, '2026-11-25 22:59:35.321857');
INSERT INTO public.journal_entry_lines VALUES (66, 35, 4, 96, '4100', 0.00, 10000.00, 'Sales revenue ORD-2026-ae041f8b-a711-4713-bbf5-9b69603824c6', '', 2, '', 0.00, '2026-11-25 22:59:35.321857');
INSERT INTO public.journal_entry_lines VALUES (67, 35, 4, 87, '2120', 0.00, 1000.00, 'Sales tax ORD-2026-ae041f8b-a711-4713-bbf5-9b69603824c6', '', 3, '', 0.00, '2026-11-25 22:59:35.321857');
INSERT INTO public.journal_entry_lines VALUES (68, 36, 4, 103, '5110', 4000.00, 0.00, 'COGS for order ORD-2026-ae041f8b-a711-4713-bbf5-9b69603824c6', '', 1, '', 0.00, '2026-11-25 22:59:35.408398');
INSERT INTO public.journal_entry_lines VALUES (69, 36, 4, 76, '1130', 0.00, 4000.00, 'Inventory reduction for order ORD-2026-ae041f8b-a711-4713-bbf5-9b69603824c6', '', 2, '', 0.00, '2026-11-25 22:59:35.408398');


--
-- Data for Name: general_ledger; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.general_ledger VALUES (1, 4, 4, 70, 1, 1, '2026-11-25', '2026-11', 11000.00, 0.00, 11000.00, '2026-11-25 22:28:54.977171');
INSERT INTO public.general_ledger VALUES (2, 4, 4, 96, 1, 2, '2026-11-25', '2026-11', 0.00, 10000.00, 10000.00, '2026-11-25 22:28:54.977171');
INSERT INTO public.general_ledger VALUES (3, 4, 4, 87, 1, 3, '2026-11-25', '2026-11', 0.00, 1000.00, 1000.00, '2026-11-25 22:28:54.977171');
INSERT INTO public.general_ledger VALUES (4, 4, 4, 103, 2, 4, '2026-11-25', '2026-11', 4000.00, 0.00, 4000.00, '2026-11-25 22:28:55.270987');
INSERT INTO public.general_ledger VALUES (5, 4, 4, 76, 2, 5, '2026-11-25', '2026-11', 0.00, 4000.00, -4000.00, '2026-11-25 22:28:55.270987');
INSERT INTO public.general_ledger VALUES (6, 4, 4, 70, 3, 6, '2026-11-25', '2026-11', 100000.00, 0.00, 111000.00, '2026-11-25 22:31:03.921079');
INSERT INTO public.general_ledger VALUES (7, 4, 4, 93, 3, 7, '2026-11-25', '2026-11', 0.00, 100000.00, 100000.00, '2026-11-25 22:31:03.921079');
INSERT INTO public.general_ledger VALUES (8, 4, 4, 111, 4, 8, '2026-11-25', '2026-11', 3000.00, 0.00, 3000.00, '2026-11-25 22:31:28.729903');
INSERT INTO public.general_ledger VALUES (9, 4, 4, 70, 4, 9, '2026-11-25', '2026-11', 0.00, 3000.00, 108000.00, '2026-11-25 22:31:28.729903');
INSERT INTO public.general_ledger VALUES (10, 4, 4, 76, 5, 10, '2026-11-25', '2026-11', 20000.00, 0.00, 16000.00, '2026-11-25 22:33:16.786267');
INSERT INTO public.general_ledger VALUES (11, 4, 4, 78, 5, 11, '2026-11-25', '2026-11', 2800.00, 0.00, 2800.00, '2026-11-25 22:33:16.786267');
INSERT INTO public.general_ledger VALUES (12, 4, 4, 85, 5, 12, '2026-11-25', '2026-11', 0.00, 22800.00, 22800.00, '2026-11-25 22:33:16.786267');
INSERT INTO public.general_ledger VALUES (13, 4, 4, 85, 6, 13, '2026-11-25', '2026-11', 10000.00, 0.00, 12800.00, '2026-11-25 22:33:49.054066');
INSERT INTO public.general_ledger VALUES (14, 4, 4, 70, 6, 14, '2026-11-25', '2026-11', 0.00, 10000.00, 98000.00, '2026-11-25 22:33:49.054066');
INSERT INTO public.general_ledger VALUES (15, 4, 4, 107, 7, 15, '2026-11-25', '2026-11', 5000.00, 0.00, 5000.00, '2026-11-25 22:37:56.655723');
INSERT INTO public.general_ledger VALUES (16, 4, 4, 71, 7, 16, '2026-11-25', '2026-11', 0.00, 5000.00, -5000.00, '2026-11-25 22:37:56.655723');
INSERT INTO public.general_ledger VALUES (17, 4, 4, 73, 8, 17, '2026-11-25', '2026-11', 11000.00, 0.00, 11000.00, '2026-11-25 22:45:54.933512');
INSERT INTO public.general_ledger VALUES (18, 4, 4, 96, 8, 18, '2026-11-25', '2026-11', 0.00, 10000.00, 20000.00, '2026-11-25 22:45:54.933512');
INSERT INTO public.general_ledger VALUES (19, 4, 4, 87, 8, 19, '2026-11-25', '2026-11', 0.00, 1000.00, 2000.00, '2026-11-25 22:45:54.933512');
INSERT INTO public.general_ledger VALUES (20, 4, 4, 103, 9, 20, '2026-11-25', '2026-11', 4000.00, 0.00, 8000.00, '2026-11-25 22:45:55.245533');
INSERT INTO public.general_ledger VALUES (21, 4, 4, 76, 9, 21, '2026-11-25', '2026-11', 0.00, 4000.00, 12000.00, '2026-11-25 22:45:55.245533');
INSERT INTO public.general_ledger VALUES (22, 4, 4, 70, 3, 6, '2026-11-25', '2026-11', 0.00, 100000.00, -2000.00, '2026-11-25 22:47:38.695827');
INSERT INTO public.general_ledger VALUES (23, 4, 4, 93, 3, 7, '2026-11-25', '2026-11', 100000.00, 0.00, 0.00, '2026-11-25 22:47:38.695827');
INSERT INTO public.general_ledger VALUES (24, 4, 4, 74, 22, 37, '2026-11-25', '2026-11', 5200.00, 0.00, 5200.00, '2026-11-25 22:49:27.574748');
INSERT INTO public.general_ledger VALUES (25, 4, 4, 96, 22, 38, '2026-11-25', '2026-11', 0.00, 4500.00, 24500.00, '2026-11-25 22:49:27.574748');
INSERT INTO public.general_ledger VALUES (26, 4, 4, 87, 22, 39, '2026-11-25', '2026-11', 0.00, 700.00, 2700.00, '2026-11-25 22:49:27.574748');
INSERT INTO public.general_ledger VALUES (27, 4, 4, 70, 23, 40, '2026-11-25', '2026-11', 2000.00, 0.00, 0.00, '2026-11-25 22:49:50.010227');
INSERT INTO public.general_ledger VALUES (28, 4, 4, 74, 23, 41, '2026-11-25', '2026-11', 0.00, 2000.00, 3200.00, '2026-11-25 22:49:50.010227');
INSERT INTO public.general_ledger VALUES (29, 4, 4, 70, 24, 42, '2026-11-25', '2026-11', 1000.00, 0.00, 1000.00, '2026-11-25 22:50:09.706039');
INSERT INTO public.general_ledger VALUES (30, 4, 4, 74, 24, 43, '2026-11-25', '2026-11', 0.00, 1000.00, 2200.00, '2026-11-25 22:50:09.706039');
INSERT INTO public.general_ledger VALUES (31, 4, 4, 70, 25, 44, '2026-11-25', '2026-11', 1000.00, 0.00, 2000.00, '2026-11-25 22:50:09.732584');
INSERT INTO public.general_ledger VALUES (32, 4, 4, 74, 25, 45, '2026-11-25', '2026-11', 0.00, 1000.00, 1200.00, '2026-11-25 22:50:09.732584');
INSERT INTO public.general_ledger VALUES (33, 4, 4, 70, 26, 46, '2026-11-25', '2026-11', 1000.00, 0.00, 3000.00, '2026-11-25 22:50:09.763848');
INSERT INTO public.general_ledger VALUES (34, 4, 4, 74, 26, 47, '2026-11-25', '2026-11', 0.00, 1000.00, 200.00, '2026-11-25 22:50:09.763848');
INSERT INTO public.general_ledger VALUES (35, 4, 4, 70, 27, 48, '2026-11-25', '2026-11', 1000.00, 0.00, 4000.00, '2026-11-25 22:50:09.790594');
INSERT INTO public.general_ledger VALUES (36, 4, 4, 74, 27, 49, '2026-11-25', '2026-11', 0.00, 1000.00, -800.00, '2026-11-25 22:50:09.790594');
INSERT INTO public.general_ledger VALUES (37, 4, 4, 70, 28, 50, '2026-11-25', '2026-11', 1000.00, 0.00, 5000.00, '2026-11-25 22:50:09.79576');
INSERT INTO public.general_ledger VALUES (38, 4, 4, 74, 28, 51, '2026-11-25', '2026-11', 0.00, 1000.00, -1800.00, '2026-11-25 22:50:09.79576');
INSERT INTO public.general_ledger VALUES (39, 4, 4, 85, 29, 52, '2026-11-25', '2026-11', 3000.00, 0.00, 9800.00, '2026-11-25 22:50:23.524633');
INSERT INTO public.general_ledger VALUES (40, 4, 4, 70, 29, 53, '2026-11-25', '2026-11', 0.00, 3000.00, 2000.00, '2026-11-25 22:50:23.524633');
INSERT INTO public.general_ledger VALUES (41, 4, 4, 85, 30, 54, '2026-11-25', '2026-11', 3000.00, 0.00, 6800.00, '2026-11-25 22:50:23.544114');
INSERT INTO public.general_ledger VALUES (42, 4, 4, 70, 30, 55, '2026-11-25', '2026-11', 0.00, 3000.00, -1000.00, '2026-11-25 22:50:23.544114');
INSERT INTO public.general_ledger VALUES (43, 4, 4, 85, 31, 56, '2026-11-25', '2026-11', 3000.00, 0.00, 3800.00, '2026-11-25 22:50:23.56892');
INSERT INTO public.general_ledger VALUES (44, 4, 4, 70, 31, 57, '2026-11-25', '2026-11', 0.00, 3000.00, -4000.00, '2026-11-25 22:50:23.56892');
INSERT INTO public.general_ledger VALUES (45, 4, 4, 85, 32, 58, '2026-11-25', '2026-11', 3000.00, 0.00, 800.00, '2026-11-25 22:50:23.595548');
INSERT INTO public.general_ledger VALUES (46, 4, 4, 70, 32, 59, '2026-11-25', '2026-11', 0.00, 3000.00, -7000.00, '2026-11-25 22:50:23.595548');
INSERT INTO public.general_ledger VALUES (47, 2, 2, 132, 33, 60, '2026-11-25', '2026-11', 11000.00, 0.00, 11000.00, '2026-11-25 22:59:09.694097');
INSERT INTO public.general_ledger VALUES (48, 2, 2, 158, 33, 61, '2026-11-25', '2026-11', 0.00, 10000.00, 10000.00, '2026-11-25 22:59:09.694097');
INSERT INTO public.general_ledger VALUES (49, 2, 2, 149, 33, 62, '2026-11-25', '2026-11', 0.00, 1000.00, 1000.00, '2026-11-25 22:59:09.694097');
INSERT INTO public.general_ledger VALUES (50, 2, 2, 165, 34, 63, '2026-11-25', '2026-11', 4000.00, 0.00, 4000.00, '2026-11-25 22:59:09.782036');
INSERT INTO public.general_ledger VALUES (51, 2, 2, 138, 34, 64, '2026-11-25', '2026-11', 0.00, 4000.00, -4000.00, '2026-11-25 22:59:09.782036');
INSERT INTO public.general_ledger VALUES (52, 4, 4, 70, 35, 65, '2026-11-25', '2026-11', 11000.00, 0.00, 4000.00, '2026-11-25 22:59:35.321857');
INSERT INTO public.general_ledger VALUES (53, 4, 4, 96, 35, 66, '2026-11-25', '2026-11', 0.00, 10000.00, 34500.00, '2026-11-25 22:59:35.321857');
INSERT INTO public.general_ledger VALUES (54, 4, 4, 87, 35, 67, '2026-11-25', '2026-11', 0.00, 1000.00, 3700.00, '2026-11-25 22:59:35.321857');
INSERT INTO public.general_ledger VALUES (55, 4, 4, 103, 36, 68, '2026-11-25', '2026-11', 4000.00, 0.00, 12000.00, '2026-11-25 22:59:35.408398');
INSERT INTO public.general_ledger VALUES (56, 4, 4, 76, 36, 69, '2026-11-25', '2026-11', 0.00, 4000.00, 8000.00, '2026-11-25 22:59:35.408398');


--
-- Data for Name: invoices_customer; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.invoices_customer VALUES (2, 4, 4, 'INV-2026-00001', '2026-11-25', '2026-11-25', 1, NULL, 5000.00, 700.00, 500.00, 5200.00, 7000.00, 0.00, 'EGP', 'paid', NULL, NULL, NULL, 22, '2026-11-25 22:49:27.574748', '2026-11-25 22:50:09.79576', 4, 4);


--
-- Data for Name: suppliers; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.suppliers VALUES (1, 4, 4, 'SUP-0001', 'Audit Supplier Co', '', 'Audit Contact', 'supplier@example.test', '', '', '', '', '', '', '', '', 'Egypt', '', 'net_30', 0.00, -22000.00, 'EGP', 85, 'food', '', '', '', '', true, '2026-11-25 22:32:17.471356', '2026-11-25 22:50:23.595548', 4, 4);


--
-- Data for Name: invoices_supplier; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.invoices_supplier VALUES (1, 4, 4, 'PINV-2026-00001', NULL, '2026-11-25', '2026-11-25', 1, 20000.00, 2800.00, 0.00, 22800.00, 22000.00, 800.00, 'EGP', 'partially_paid', NULL, NULL, 5, '2026-11-25 22:32:47.808276', '2026-11-25 22:50:23.595548', 4, 4);


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.order_items VALUES (1, 1, 1, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:27:01.093643', '2026-11-25 22:27:01.093643', 2, 2, 0);
INSERT INTO public.order_items VALUES (2, 2, 2, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:28:45.214093', '2026-11-25 22:28:45.214093', 4, 4, 0);
INSERT INTO public.order_items VALUES (3, 3, 2, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:41:03.166192', '2026-11-25 22:41:03.166192', 4, 4, 0);
INSERT INTO public.order_items VALUES (4, 4, 2, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:41:16.27453', '2026-11-25 22:41:16.27453', 4, 4, 0);
INSERT INTO public.order_items VALUES (5, 5, 2, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:46:25.455945', '2026-11-25 22:46:25.455945', 4, 4, 0);
INSERT INTO public.order_items VALUES (6, 6, 2, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:46:25.455973', '2026-11-25 22:46:25.455973', 4, 4, 0);
INSERT INTO public.order_items VALUES (7, 7, 2, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:46:25.463445', '2026-11-25 22:46:25.463445', 4, 4, 0);
INSERT INTO public.order_items VALUES (8, 8, 2, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:46:25.531573', '2026-11-25 22:46:25.531573', 4, 4, 0);
INSERT INTO public.order_items VALUES (9, 9, 2, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:46:25.540528', '2026-11-25 22:46:25.540528', 4, 4, 0);
INSERT INTO public.order_items VALUES (10, 10, 2, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:46:31.558551', '2026-11-25 22:46:31.558551', 4, 4, 0);
INSERT INTO public.order_items VALUES (11, 11, 2, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:46:31.667136', '2026-11-25 22:46:31.667136', 4, 4, 0);
INSERT INTO public.order_items VALUES (12, 12, 2, NULL, 'Audit Burger', NULL, 1, 10000.00, 0.00, 10000.00, NULL, NULL, '2026-11-25 22:46:31.788601', '2026-11-25 22:46:31.788601', 4, 4, 0);


--
-- Data for Name: payments_made; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.payments_made VALUES (1, 4, 4, 'PMT-M-2026-00001', '2026-11-25', 1, 1, 10000.00, 'cash', NULL, NULL, NULL, 6, 'completed', '2026-11-25 22:33:49.054066', '2026-11-25 22:33:49.054066', 4);
INSERT INTO public.payments_made VALUES (2, 4, 4, 'PMT-M-2026-00002', '2026-11-25', 1, 1, 3000.00, 'cash', NULL, NULL, NULL, 29, 'completed', '2026-11-25 22:50:23.524633', '2026-11-25 22:50:23.524633', 4);
INSERT INTO public.payments_made VALUES (3, 4, 4, 'PMT-M-2026-00003', '2026-11-25', 1, 1, 3000.00, 'cash', NULL, NULL, NULL, 30, 'completed', '2026-11-25 22:50:23.544114', '2026-11-25 22:50:23.544114', 4);
INSERT INTO public.payments_made VALUES (4, 4, 4, 'PMT-M-2026-00004', '2026-11-25', 1, 1, 3000.00, 'cash', NULL, NULL, NULL, 31, 'completed', '2026-11-25 22:50:23.56892', '2026-11-25 22:50:23.56892', 4);
INSERT INTO public.payments_made VALUES (5, 4, 4, 'PMT-M-2026-00005', '2026-11-25', 1, 1, 3000.00, 'cash', NULL, NULL, NULL, 32, 'completed', '2026-11-25 22:50:23.595548', '2026-11-25 22:50:23.595548', 4);


--
-- Data for Name: payments_received; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.payments_received VALUES (1, 4, 4, 'PMT-R-2026-00001', '2026-11-25', 1, 2, 2000.00, 'cash', NULL, NULL, NULL, 23, 'completed', '2026-11-25 22:49:50.010227', '2026-11-25 22:49:50.010227', 4);
INSERT INTO public.payments_received VALUES (2, 4, 4, 'PMT-R-2026-00002', '2026-11-25', 1, 2, 1000.00, 'cash', NULL, NULL, NULL, 24, 'completed', '2026-11-25 22:50:09.706039', '2026-11-25 22:50:09.706039', 4);
INSERT INTO public.payments_received VALUES (3, 4, 4, 'PMT-R-2026-00003', '2026-11-25', 1, 2, 1000.00, 'cash', NULL, NULL, NULL, 25, 'completed', '2026-11-25 22:50:09.732584', '2026-11-25 22:50:09.732584', 4);
INSERT INTO public.payments_received VALUES (4, 4, 4, 'PMT-R-2026-00004', '2026-11-25', 1, 2, 1000.00, 'cash', NULL, NULL, NULL, 26, 'completed', '2026-11-25 22:50:09.763848', '2026-11-25 22:50:09.763848', 4);
INSERT INTO public.payments_received VALUES (5, 4, 4, 'PMT-R-2026-00005', '2026-11-25', 1, 2, 1000.00, 'cash', NULL, NULL, NULL, 27, 'completed', '2026-11-25 22:50:09.790594', '2026-11-25 22:50:09.790594', 4);
INSERT INTO public.payments_received VALUES (6, 4, 4, 'PMT-R-2026-00006', '2026-11-25', 1, 2, 1000.00, 'cash', NULL, NULL, NULL, 28, 'completed', '2026-11-25 22:50:09.79576', '2026-11-25 22:50:09.79576', 4);


--
-- Data for Name: salaries; Type: TABLE DATA; Schema: public; Owner: pos_migrator
--

INSERT INTO public.salaries VALUES (1, 4, 4, 4, '2026-11-01', '2026-11-30', NULL, 11, 2026, 5000.00, 'EGP', 0.00, 0.00, 0.00, 0.00, 0.00, 0.00, 0.00, 0.00, '{}', 0.00, 0.00, 0.00, 0.00, 0.00, 0.00, 0.00, '{}', 5000.00, 0.00, 5000.00, 20, 0, 160.00, 0.00, 'bank_transfer', 'PAY-1795646276249', NULL, NULL, 'paid', true, '2026-11-25 22:37:56.655723+00', 4, false, NULL, NULL, NULL, NULL, NULL, '2026-11-25 22:36:46.463516+00', '2026-11-25 22:37:56.655723+00', 4, NULL);


--
-- Name: accounting_event_outbox_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.accounting_event_outbox_id_seq', 23, true);


--
-- Name: general_ledger_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.general_ledger_id_seq', 56, true);


--
-- Name: invoices_customer_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.invoices_customer_id_seq', 2, true);


--
-- Name: invoices_supplier_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.invoices_supplier_id_seq', 1, true);


--
-- Name: journal_entries_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.journal_entries_id_seq', 36, true);


--
-- Name: journal_entry_lines_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.journal_entry_lines_id_seq', 69, true);


--
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.order_items_id_seq', 12, true);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.orders_id_seq', 12, true);


--
-- Name: payments_made_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.payments_made_id_seq', 5, true);


--
-- Name: payments_received_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.payments_received_id_seq', 6, true);


--
-- Name: salaries_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.salaries_id_seq', 1, true);


--
-- Name: suppliers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: pos_migrator
--

SELECT pg_catalog.setval('public.suppliers_id_seq', 1, true);


--
-- PostgreSQL database dump complete
--

\unrestrict Rut16uWqioy1Nm2KLaD0DgQxc6JEWENpvwRyoNncgkwbeylzArcDtLVTrjM8th1

