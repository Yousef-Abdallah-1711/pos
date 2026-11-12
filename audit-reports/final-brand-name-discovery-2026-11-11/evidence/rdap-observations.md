# Direct RDAP response log

Date: 2026-11-11. Direct HTTP GET requests from the local research environment.

| URL | HTTP result | Classification |
|---|---:|---|
| https://rdap.verisign.com/com/v1/domain/ularo.com | 200 | Domain registration object returned; registered |
| https://rdap.verisign.com/com/v1/domain/cendro.com | 200 | Domain registration object returned; registered |
| https://rdap.verisign.com/com/v1/domain/flarivo.com | 200 | Domain registration object returned; registered |
| https://rdap.nic.io/domain/ularo.io and corresponding Cendro/Flarivo URLs | No conclusive status | Unknown |
| https://rdap.nic.google/domain/ularo.app and corresponding Cendro/Flarivo URLs | No conclusive status | Unknown |
| https://rdap.nic.co/domain/ularo.co and corresponding Cendro/Flarivo URLs | No conclusive status | Unknown |

No WHOIS owner inference, registrar checkout, purchase, or reservation was performed.

## Follow-up checks — 2026-11-11

The three `.com` records were rechecked against Verisign RDAP and returned HTTP 200. The earlier probes for `.io`, `.app` and `.co` used invalid/unresolvable hostnames (`rdap.nic.io`, `rdap.nic.google`, `rdap.nic.co`) and should be treated as endpoint failures, not RDAP evidence. Corrected follow-up queries were made through the `rdap.org` bootstrap endpoint, which returns the registry response:

| Domain | HTTP | Observation |
|---|---:|---|
| ularo.io | 404 | No RDAP object returned at query time |
| ularo.app | 200 | Registration object returned |
| ularo.co | 404 | No RDAP object returned at query time |
| cendro.io | 404 | No RDAP object returned at query time |
| cendro.app | 200 | Registration object returned |
| cendro.co | 404 | No RDAP object returned at query time |
| flarivo.io | 404 | No RDAP object returned at query time |
| flarivo.app | 404 | No RDAP object returned at query time |
| flarivo.co | 404 | No RDAP object returned at query time |

Queries used `https://rdap.org/domain/{name}`. A 404 does not establish that a domain is purchasable, registrable by this owner, or free from conflicting rights. `.com.eg` was not queried. No purchase or reservation was made.
