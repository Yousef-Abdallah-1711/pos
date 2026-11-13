# 01 — Phase dependency map

~~~mermaid
flowchart TD
 P00[Phase 00 safety and integrity detection] --> P01[Phase 01 money and ledger contracts]
 P01 --> P02[Phase 02 posting correctness]
 P02 --> P03[Phase 03 transaction recovery]
 P03 --> P04[Phase 04 order and payment]
 P03 --> P05[Phase 05 inventory and AR/AP]
 P04 --> P06[Phase 06 reporting and close]
 P05 --> P06
 P06 --> P07[Phase 07 financial security]
 P06 --> P08[Phase 08 exports and UX]
 P07 --> P08
 P08 --> P09[Phase 09 measured scale]
 P09 --> P10[Phase 10 release verification]
~~~

Security and policy review are continuous gates. No phase auto-starts.