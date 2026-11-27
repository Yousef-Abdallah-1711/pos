# Appearance API discovery

In `backend/cmd/api/main.go`, protected endpoints include `GET /api/v1/users/{userId}/settings`, `PUT .../settings/language`, `PUT .../settings/theme`, and `PUT .../settings/theme-colors`. `backend/internal/handler/http/user_settings_handler.go` and `backend/internal/repository/user_settings_repo.go` own personal settings. `GET/POST/PUT/DELETE /api/v1/admin/themes` and theme activation/history endpoints are wired through module 9 READ/WRITE/DELETE checks and concern the site theme builder. Public storefront theme retrieval includes `GET /api/v1/public/restaurants/{slug}/theme`.

Review `backend/internal/handler/http/admin_theme_handler_v2.go`, `backend/internal/service/theme_service.go`, the actual wired handler and repository, and corresponding isolation tests before extending them. Existing role/permission checks must be verified per action; this document does not establish that any new tenant-appearance endpoint is secure.
