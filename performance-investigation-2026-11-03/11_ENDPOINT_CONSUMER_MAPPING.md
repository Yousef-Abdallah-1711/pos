# Endpoint Consumer Mapping

`08_COMPLETE_BACKEND_ENDPOINT_INVENTORY.csv` is a source registration inventory only. `05_PAGE_TO_API_DEPENDENCY_MATRIX_README.md` records that frontend consumers were not fully joined to routes. Existing backend route manifest has no verified consumer data for all entries. Therefore this report does not label endpoints unused based on absent evidence.

Useful static checks for the next audit pass:

1. Enumerate frontend Axios/fetch calls from each API client/hook with caller file and exported hook.
2. Normalize URL interpolation to method+path patterns, including aliases and URL prefix configuration.
3. Match patterns to route registrations, preserving conditional handlers and method aliases.
4. Separately inspect background workers, webhooks, tests and server components.
5. Mark each mapping as static or observed; store no auth header/cookie/body values.
