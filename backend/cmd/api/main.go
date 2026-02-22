package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"pos-saas/internal/config"
	handler "pos-saas/internal/handler/http"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/database"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
)

// In-memory storage for stub mode
type Role struct {
	ID               int       `json:"id"`
	TenantID         int       `json:"tenant_id"`
	RoleName         string    `json:"role_name"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	IsSystem         bool      `json:"is_system"`
	PermissionsCount int       `json:"permissions_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type InMemoryRBACStorage struct {
	mu        sync.RWMutex
	roles     map[int]*Role
	nextID    int
	userRoles map[int][]int // userId -> []roleId
}

var rbacStorage = &InMemoryRBACStorage{
	roles:     make(map[int]*Role),
	nextID:    2, // Start from 2 since Admin role is 1
	userRoles: make(map[int][]int),
}

func init() {
	// Initialize with Admin role
	rbacStorage.roles[1] = &Role{
		ID:               1,
		TenantID:         1,
		RoleName:         "Admin",
		Name:             "Admin",
		Description:      "System administrator with full access",
		IsSystem:         true,
		PermissionsCount: 6,
		CreatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	// Assign Admin role to user 7 (the mock user)
	rbacStorage.userRoles[7] = []int{1}
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	var db *sql.DB
	db, err = database.Connect(&cfg.Database)
	if err != nil {
		log.Printf("⚠️ Warning: Failed to connect to database: %v. Continuing without database (stub mode).", err)
		db = nil
	} else {
		defer db.Close()
	}

	// Initialize JWT service
	tokenService, err := jwt.NewTokenService(cfg.JWT.Secret, cfg.JWT.Expiry)
	if err != nil {
		log.Fatalf("Failed to create token service: %v", err)
	}

	// Initialize repositories
	authRepo := repository.NewAuthRepository(db)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	userRepo := repository.NewUserRepository(db)
	userSettingsRepo := repository.NewUserSettingsRepository(db)

	// Order Management repositories
	orderRepo := repository.NewOrderRepository(db)

	// Driver Management repository
	driverRepo := repository.NewDriverRepository(db)

	// HR Module repositories
	employeeRepo := repository.NewEmployeeRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	salaryRepo := repository.NewSalaryRepository(db)
	leaveRepo := repository.NewLeaveRepository(db)

	// Notification repository
	notificationRepo := repository.NewNotificationRepository(db)

	// RBAC repositories
	var rbacRoleRepo *repository.RBACRoleRepository
	var rolePermissionRepo *repository.RolePermissionRepository
	var userRoleRepo *repository.UserRoleRepository
	if db != nil {
		rbacRoleRepo = repository.NewRBACRoleRepository(db)
		rolePermissionRepo = repository.NewRolePermissionRepository(db)
		userRoleRepo = repository.NewUserRoleRepository(db)
	}

	// Theme repository (Phase 1)
	themeRepo := repository.NewThemeRepository(db)

	// Super Admin repositories (Phase 6 - Super Admin Dashboard)
	var superAdminRepo *repository.SuperAdminRepository
	var impersonationRepo *repository.ImpersonationRepository
	var superAdminAuditRepo *repository.SuperAdminAuditRepository
	if db != nil {
		superAdminRepo = repository.NewSuperAdminRepository(db)
		impersonationRepo = repository.NewImpersonationRepository(db)
		superAdminAuditRepo = repository.NewSuperAdminAuditRepository(db)
	}

	// Subscription repositories (Phase 1-8 - SaaS Trial/Subscription System)
	var subscriptionRepo repository.SubscriptionRepository
	if db != nil {
		subscriptionRepo = repository.NewSubscriptionRepository(db)
	}

	// Payment repository (Phase 1-3 - Paymob integration)
	var paymentRepo repository.PaymentRepository
	if db != nil {
		paymentRepo = repository.NewPaymentRepository(db)
	}

	// Subscription plans repository (Phase 12 - DB-driven plans)
	var planRepo repository.PlanRepository
	if db != nil {
		planRepo = repository.NewPlanRepository(db)
	}

	// Paymob client mode selection (mock by default).
	paymobMode := strings.ToLower(strings.TrimSpace(os.Getenv("PAYMOB_MODE")))
	var paymobClient service.PaymobClient
	if paymobMode == "live" {
		iframeID, _ := strconv.Atoi(strings.TrimSpace(os.Getenv("PAYMOB_IFRAME_ID")))
		paymobClient = service.NewRealPaymobClient(
			strings.TrimSpace(os.Getenv("PAYMOB_API_KEY")),
			iframeID,
			strings.TrimSpace(os.Getenv("PAYMOB_HMAC_SECRET")),
			strings.TrimSpace(os.Getenv("PAYMOB_BASE_URL")),
		)
		log.Println("[PaymentSystem] Running in LIVE mode")
	} else {
		paymobClient = service.NewMockPaymobClient()
		log.Println("[PaymentSystem] Running in MOCK mode")
	}

	// NOTE: V2 repositories will be implemented in Phase 2
	// themeRepo := repository.NewThemeRepositoryV2(db)
	// componentRepo := repository.NewComponentRepositoryV2(db)
	// libraryRepo := repository.NewComponentLibraryRepositoryV2(db)
	// themeHistoryRepo := repository.NewThemeHistoryRepositoryV2(db)

	// NOTE: Old repositories commented out - using v2 API instead
	// Theme repository (Restaurant Theme)
	// themeRepo := repository.NewRestaurantThemeRepository(db)

	// Homepage Builder repositories
	// homepageThemeRepo := repository.NewThemeRepository(db)
	// homepageSectionRepo := repository.NewHomepageSectionRepository(db)
	// websiteSettingsRepo := repository.NewWebsiteSettingsRepository(db)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(authRepo, tokenService)
	// Wire up super admin repo for token exchange (super admin dashboard access)
	if superAdminRepo != nil {
		authUseCase.SetSuperAdminRepo(superAdminRepo)
	}
	productUC := usecase.NewProductUseCase(productRepo, notificationRepo, "http://localhost:8080/uploads")
	// NOTE: User settings use case reserved for Phase 2
	notificationUC := usecase.NewNotificationUseCase(notificationRepo)
	orderUC := usecase.NewOrderUseCase(orderRepo, productRepo)

	// RBAC use cases (Phase 2 - commented out, not used in current phase)
	// 	// roleUC := usecase.NewRoleUseCase(roleRepo, auditLogRepo)
	// 	// permissionUC := usecase.NewPermissionUseCase(rolePermissionRepo, userRoleRepo, moduleDefRepo, auditLogRepo)
	// 	// userRoleUC := usecase.NewUserRoleUseCase(userRoleRepo, auditLogRepo)

	// Theme service (Phase 1)
	themeService := service.NewThemeService(themeRepo)

	// Super Admin token service (uses separate secret for extra security)
	superAdminTokenService, err := jwt.NewSuperAdminTokenService(cfg.JWT.Secret+"_super_admin", "15m")
	if err != nil {
		log.Printf("Warning: Failed to create super admin token service: %v", err)
	}

	// Super Admin use cases
	var superAdminAuthUC *usecase.SuperAdminAuthUseCase
	var impersonationUC *usecase.ImpersonationUseCase
	var superAdminTenantUC *usecase.SuperAdminTenantUseCase
	if db != nil && superAdminRepo != nil {
		superAdminAuthUC = usecase.NewSuperAdminAuthUseCase(superAdminRepo, userRepo, superAdminAuditRepo, superAdminTokenService)
		impersonationUC = usecase.NewImpersonationUseCase(impersonationRepo, superAdminRepo, superAdminAuditRepo)
		superAdminTenantUC = usecase.NewSuperAdminTenantUseCase(superAdminRepo, superAdminAuditRepo)
	}

	// Subscription use case (Phase 1-8 - SaaS Trial/Subscription System)
	var subscriptionUC *usecase.SubscriptionUseCase
	if db != nil && subscriptionRepo != nil && superAdminAuditRepo != nil {
		subscriptionUC = usecase.NewSubscriptionUseCase(subscriptionRepo, superAdminAuditRepo)
	}

	// Payment use case (Phase 1-3 - Paymob integration)
	var paymentUC *usecase.PaymentUseCase
	if paymobClient != nil {
		// Create in stub mode too so /subscription/plans remains available without DB.
		paymentUC = usecase.NewPaymentUseCase(paymentRepo, subscriptionRepo, planRepo, paymobClient, db)
		log.Printf("[PaymentSystem] PaymentUseCase initialized (db=%t)", db != nil)
	} else {
		log.Println("[PaymentSystem] PaymentUseCase skipped (paymobClient is nil)")
	}

	// NOTE: Old theme use cases commented out - replaced with v2 repository API
	// themeUC := usecase.NewRestaurantThemeUseCase(themeRepo)
	// homepageThemeUC := usecase.NewThemeUsecase(homepageThemeRepo)
	// homepageSectionUC := usecase.NewSectionUsecase(homepageSectionRepo)
	// websiteSettingsUC := usecase.NewSettingsUsecase(websiteSettingsRepo)

	// Initialize handlers
	restaurantRepo := repository.NewRestaurantRepository(db)

	authHandler := handler.NewAuthHandler(authUseCase)
	productHandler := handler.NewProductHandler(productUC)
	categoryHandler := handler.NewCategoryHandler(categoryRepo)
	publicMenuHandler := handler.NewPublicMenuHandler(productUC, restaurantRepo, categoryRepo)
	userSettingsHandler := handler.NewUserSettingsHandler(userSettingsRepo, userRepo)
	translationHandler := handler.NewTranslationHandler()

	// Order Management handlers
	publicOrderHandler := handler.NewPublicOrderHandler(orderUC, restaurantRepo)
	adminOrderHandler := handler.NewAdminOrderHandler(orderUC, driverRepo, userRoleRepo)

	// Driver Management handler
	adminDriverHandler := handler.NewAdminDriverHandler(driverRepo, orderUC)

	// HR Module handlers
	hrEmployeeHandler := handler.NewHREmployeeHandler(employeeRepo)
	hrAttendanceHandler := handler.NewHRAttendanceHandler(attendanceRepo, db)
	hrLeaveHandler := handler.NewHRLeaveHandler(leaveRepo, db)
	hrSalaryHandler := handler.NewHRSalaryHandler(salaryRepo, db)

	// Notification handler
	notificationHandler := handler.NewNotificationHandler(notificationUC)

	// RBAC handlers
	var rbacRoleHandler *handler.RBACRoleHandler
	var rbacPermissionHandler *handler.RBACPermissionHandler
	var rbacUserRoleHandler *handler.RBACUserRoleHandler
	if db != nil {
		rbacRoleHandler = handler.NewRBACRoleHandler(rbacRoleRepo, rolePermissionRepo)
		rbacPermissionHandler = handler.NewRBACPermissionHandler(rolePermissionRepo, userRoleRepo, rbacRoleRepo)
		rbacUserRoleHandler = handler.NewRBACUserRoleHandler(userRoleRepo, rbacRoleRepo, userRepo)
	}

	// Theme handler (Phase 1)
	themeHandler := handler.NewAdminThemeHandler(themeService, nil)
	// Theme handler V2 (Phase 5 - with presets support)
	themeHandlerV2 := handler.NewAdminThemeHandlerV2(themeService)
	// Public homepage handler (Phase 1 - now enabled)
	publicHomepageHandler := handler.NewPublicHomepageHandler(restaurantRepo, themeRepo)

	// Super Admin handlers (Phase 6 - Super Admin Dashboard)
	var superAdminHandler *handler.SuperAdminHandler
	var impersonationHandler *handler.ImpersonationHandler
	if superAdminAuthUC != nil && superAdminTenantUC != nil && impersonationUC != nil {
		superAdminHandler = handler.NewSuperAdminHandler(superAdminAuthUC, superAdminTenantUC)
		impersonationHandler = handler.NewImpersonationHandler(impersonationUC)
	}

	// Subscription handler (Phase 1-8 - SaaS Trial/Subscription System)
	var subscriptionHandler *handler.SubscriptionHandler
	if subscriptionUC != nil {
		subscriptionHandler = handler.NewSubscriptionHandler(subscriptionUC)
	}

	// Payment handler (Phase 1-3 - Paymob integration)
	var paymentHandler *handler.PaymentHandler
	var paymentMonitoringHandler *handler.PaymentMonitoringHandler
	var webhookRetryQueue *service.WebhookRetryQueue
	var webhookMonitor *service.PaymentWebhookMonitor
	var errorTracker service.ErrorTracker
	if paymentUC != nil {
		webhookMonitor = service.NewPaymentWebhookMonitor(200)
		errorTracker = service.NewErrorTrackerFromEnv()
		if errorTracker == nil {
			errorTracker = service.NewNoopTracker()
		}

		if db != nil {
			webhookRetryQueue = service.NewWebhookRetryQueue(db)

			// Worker retries transient webhook failures (1m -> 5m -> 30m -> 2h -> 12h).
			baseURL := "http://localhost:" + cfg.Server.Port
			retryWorker := service.NewWebhookRetryWorker(db, baseURL, webhookMonitor, errorTracker)
			retryWorker.Start(context.Background())

			// Worker marks active paid subscriptions as expired once period end is reached.
			expiryWorker := service.NewSubscriptionExpiryWorker(db)
			expiryWorker.Start(context.Background())
		}

		paymentHandler = handler.NewPaymentHandler(paymentUC).WithObservability(webhookMonitor, webhookRetryQueue, errorTracker)
		paymentMonitoringHandler = handler.NewPaymentMonitoringHandler(db, webhookMonitor, webhookRetryQueue)
	}

	// Mock Paymob handler (Phase 4 - enabled only in mock mode)
	var mockPaymobHandler *handler.MockPaymobHandler
	if paymobMode == "" || paymobMode == "mock" {
		mockPaymobHandler = handler.NewMockPaymobHandler()
	}

	// NOTE: Theme handlers commented out - old API being replaced with v2
	// Keeping minimal v2 API support
	// themeHandler := handler.NewAdminThemeHandler(themeUC)
	// adminHomepageHandler := handler.NewAdminHomepageHandler(...)

	// Setup routes
	mux := http.NewServeMux()

	// Public routes - Authentication
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/login/confirm", authHandler.LoginConfirm)
	mux.HandleFunc("POST /api/v1/auth/token-exchange", authHandler.TokenExchange)
	mux.HandleFunc("POST /api/v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("GET /api/v1/auth/check-subdomain", authHandler.CheckSubdomainAvailability)

	// Public routes - Menu API (no authentication required)
	mux.HandleFunc("GET /api/v1/public/restaurants/{slug}/menu", publicMenuHandler.GetRestaurantMenu)
	mux.HandleFunc("GET /api/v1/public/restaurants/{slug}/products", publicMenuHandler.GetRestaurantProducts)
	mux.HandleFunc("GET /api/v1/public/restaurants/{slug}/categories", publicMenuHandler.GetRestaurantCategories)
	mux.HandleFunc("GET /api/v1/public/restaurants/{slug}/info", publicMenuHandler.GetRestaurantInfo)
	mux.HandleFunc("GET /api/v1/public/restaurants/{slug}/categories/{categoryId}/products", publicMenuHandler.GetCategoryProducts)
	mux.HandleFunc("GET /api/v1/public/restaurants/{slug}/products/{productId}", publicMenuHandler.GetProductDetails)
	mux.HandleFunc("GET /api/v1/public/restaurants/{slug}/search", publicMenuHandler.SearchProducts)

	// Public routes - Homepage API (no authentication required)
	// NOTE: Phase 1 - now enabled for restaurant website preview
	mux.HandleFunc("GET /api/v1/public/restaurants/{slug}/homepage", publicHomepageHandler.GetHomepageData)
	mux.HandleFunc("GET /api/v1/public/restaurants/{slug}/theme", publicHomepageHandler.GetThemeOnly)
	mux.HandleFunc("GET /api/v1/public/restaurants/{slug}/sections", publicHomepageHandler.GetSectionsOnly)
	mux.HandleFunc("GET /api/v1/public/restaurants/{slug}/settings", publicHomepageHandler.GetSettingsOnly)

	// Public routes - Order API (minimal authentication required)
	// These routes require X-Tenant-ID and X-Restaurant-ID headers from middleware context
	// POST routes for order creation and validation
	mux.Handle("POST /api/v1/public/orders", middleware.TenantContextMiddleware(http.HandlerFunc(publicOrderHandler.CreateOrder)))
	mux.Handle("POST /api/v1/public/orders/validate", middleware.TenantContextMiddleware(http.HandlerFunc(publicOrderHandler.ValidateOrder)))

	mux.Handle("GET /api/v1/public/orders/number/{number}/lookup", middleware.TenantContextMiddleware(http.HandlerFunc(publicOrderHandler.GetOrderByNumber)))
	mux.Handle("GET /api/v1/public/orders/{id}/status", middleware.TenantContextMiddleware(http.HandlerFunc(publicOrderHandler.GetOrderStatus)))
	mux.Handle("GET /api/v1/public/orders/{id}/track", middleware.TenantContextMiddleware(http.HandlerFunc(publicOrderHandler.TrackOrder)))

	// GET route for order retrieval (supports both by ID and tracking)
	// Handlers can use request context or query params to determine behavior
	mux.Handle("GET /api/v1/public/orders/{id}", middleware.TenantContextMiddleware(http.HandlerFunc(publicOrderHandler.GetOrder)))

	// DELETE route for order cancellation
	mux.Handle("DELETE /api/v1/public/orders/{id}", middleware.TenantContextMiddleware(http.HandlerFunc(publicOrderHandler.CancelOrder)))

	// Translation routes (public - no authentication required)
	mux.HandleFunc("GET /api/v1/translations/health", translationHandler.HealthCheck)
	mux.HandleFunc("GET /api/v1/translations/languages", translationHandler.GetSupportedLanguages)
	mux.HandleFunc("GET /api/v1/translations/{language}", translationHandler.GetTranslations)
	mux.HandleFunc("GET /api/v1/translations/{language}/{namespace}/{key}", translationHandler.GetTranslationKey)
	mux.HandleFunc("POST /api/v1/translations/{language}/upload", translationHandler.UploadTranslations)

	// NOTE: Component library route commented out - using v2 library API
	// Component library route (public - no authentication required)
	// mux.HandleFunc("GET /api/v1/admin/components", themeHandler.GetComponentLibrary)

	// Health check
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Test route to verify routing
	mux.Handle("GET /api/v1/public/test-routing/{slug}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"message":"test route works","slug":"%s"}`, slug)))
	}))

	// Public Subscription Routes (Phase 1-8 - SaaS Trial/Subscription System)
	// Tenant status check by slug (no authentication)
	if subscriptionHandler != nil {
		log.Println("✅ Registering public subscription routes")
		mux.Handle("GET /api/v1/public/tenant-status/{slug}", http.HandlerFunc(subscriptionHandler.GetTenantStatusBySlug))
	} else {
		log.Println("⚠️ subscriptionHandler is nil, skipping public subscription route")
	}

	// Public payment routes (Phase 1-3 - Paymob integration)
	if paymentHandler != nil {
		log.Println("[PaymentSystem] Registering public payment routes")
		mux.HandleFunc("GET /api/v1/subscription/plans", paymentHandler.GetPlans)
		mux.Handle(
			"POST /api/v1/webhook/paymob",
			middleware.RateLimitByIP(60)(http.HandlerFunc(paymentHandler.HandlePaymobWebhook)),
		)
	} else {
		log.Println("[PaymentSystem] WARNING: paymentHandler is nil, /api/v1/subscription/plans will return 404")
	}

	// Mock Paymob routes (Phase 4 - local/dev testing only)
	if mockPaymobHandler != nil {
		mux.HandleFunc("GET /api/v1/mock/paymob/iframe", mockPaymobHandler.ServeIframe)
		mux.HandleFunc("POST /api/v1/mock/paymob/process", mockPaymobHandler.ProcessPayment)
		log.Println("[PaymentSystem] Running in MOCK mode - mock Paymob endpoints enabled")
	}

	// NOTE: Protected routes are now registered individually with middleware wrapping below.
	// This ensures proper execution order: Auth -> TenantContext -> Handler

	// Helper function to wrap handlers with auth and tenant middleware
	wrapProtected := func(next http.Handler) http.Handler {
		// Apply middleware in REVERSE order of execution
		// Auth must execute FIRST to set claims in context, then Tenant can read them
		// The last middleware applied is the outermost and executes first
		wrapped := middleware.TenantContextMiddleware(next)        // This will be innermost (execute 2nd)
		wrapped = middleware.AuthMiddleware(tokenService)(wrapped) // This will be outermost (execute 1st)
		return wrapped
	}

	// Helper function to wrap handlers with auth, tenant, and subscription status middleware
	wrapProtectedWithStatus := func(next http.Handler) http.Handler {
		// Apply middleware in REVERSE order of execution
		wrapped := middleware.TenantStatusMiddleware(subscriptionRepo)(next) // innermost (execute 3rd)
		wrapped = middleware.TenantContextMiddleware(wrapped)                // execute 2nd
		wrapped = middleware.AuthMiddleware(tokenService)(wrapped)           // outermost (execute 1st)
		return wrapped
	}

	// Helper function to wrap handlers with auth, tenant, and permission middleware
	// NOTE: Module IDs should match those created in database migrations (Phase 1)
	// Common module IDs: Products=1, HR=2, Notifications=3, Orders=4, etc.
	wrapWithPermission := func(next http.Handler, moduleID int64, permissionLevel string) http.Handler {
		// Apply middleware in REVERSE order of execution
		wrapped := middleware.WithRequiredPermission(moduleID, permissionLevel)(next) // innermost
		// 		wrapped = middleware.PermissionMiddleware(permissionUC)(wrapped)                       // permission check
		wrapped = middleware.TenantContextMiddleware(wrapped)      // tenant context
		wrapped = middleware.AuthMiddleware(tokenService)(wrapped) // outermost (auth first)
		return wrapped
	}

	// Payment routes: auth required, but no tenant status middleware (blocked tenants must still pay).
	if paymentHandler != nil {
		mux.Handle("POST /api/v1/subscription/initiate", wrapProtected(http.HandlerFunc(paymentHandler.InitiatePayment)))
		mux.Handle("GET /api/v1/subscription/payments", wrapProtected(http.HandlerFunc(paymentHandler.GetPaymentHistory)))
		mux.Handle("GET /api/v1/subscription/current", wrapProtected(http.HandlerFunc(paymentHandler.GetCurrentSubscription)))
		mux.Handle("POST /api/v1/subscription/cancel", wrapProtected(http.HandlerFunc(paymentHandler.CancelSubscription)))
	}

	// Register all protected routes with middleware wrapping
	// Each route handler is individually wrapped with auth + tenant middleware

	// User info endpoint
	mux.Handle("GET /api/v1/me", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetUserClaims(r)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"user_id":%d,"email":"%s"}`, claims.UserID, claims.Email)))
	})))

	// Product management endpoints (require authentication + RBAC permission)
	// Module ID 1 = Products (from migrations)
	mux.Handle("POST /api/v1/products", wrapWithPermission(http.HandlerFunc(productHandler.CreateProduct), 1, "WRITE"))
	mux.Handle("GET /api/v1/products", wrapWithPermission(http.HandlerFunc(productHandler.ListProducts), 1, "READ"))
	mux.Handle("GET /api/v1/products/alerts/low-stock", wrapWithPermission(http.HandlerFunc(productHandler.GetLowStockAlerts), 1, "READ"))
	mux.Handle("GET /api/v1/products/{id}", wrapWithPermission(http.HandlerFunc(productHandler.GetProduct), 1, "READ"))
	mux.Handle("PUT /api/v1/products/{id}", wrapWithPermission(http.HandlerFunc(productHandler.UpdateProduct), 1, "WRITE"))
	mux.Handle("DELETE /api/v1/products/{id}", wrapWithPermission(http.HandlerFunc(productHandler.DeleteProduct), 1, "DELETE"))

	// Category management endpoints (require authentication)
	mux.Handle("GET /api/v1/categories", wrapProtected(http.HandlerFunc(categoryHandler.ListCategories)))
	mux.Handle("GET /api/v1/categories/{id}", wrapProtected(http.HandlerFunc(categoryHandler.GetCategory)))
	mux.Handle("POST /api/v1/categories", wrapProtected(http.HandlerFunc(categoryHandler.CreateCategory)))
	mux.Handle("PUT /api/v1/categories/{id}", wrapProtected(http.HandlerFunc(categoryHandler.UpdateCategory)))
	mux.Handle("DELETE /api/v1/categories/{id}", wrapProtected(http.HandlerFunc(categoryHandler.DeleteCategory)))

	// Subscription endpoints (Phase 1-8 - SaaS Trial/Subscription System)
	// Tenant status check (requires authentication)
	// Note: TenantStatusMiddleware will skip this path internally (see shouldSkipStatusCheck)
	if subscriptionHandler != nil && subscriptionRepo != nil {
		mux.Handle("GET /api/v1/tenant/status", wrapProtectedWithStatus(http.HandlerFunc(subscriptionHandler.GetTenantStatus)))
	}

	// Tenant contact info endpoints (for maintenance page contact details)
	mux.Handle("GET /api/v1/tenant/contact-info", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := middleware.GetTenantID(r)
		if tenantID == 0 {
			http.Error(w, `{"error":"tenant context not found"}`, http.StatusBadRequest)
			return
		}
		var email, phone, name string
		err := db.QueryRowContext(r.Context(), `SELECT COALESCE(name,''), COALESCE(email,''), COALESCE(phone,'') FROM tenants WHERE id = $1`, tenantID).Scan(&name, &email, &phone)
		if err != nil {
			http.Error(w, `{"error":"tenant not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]string{
				"name":  name,
				"email": email,
				"phone": phone,
			},
		})
	})))

	mux.Handle("PUT /api/v1/tenant/contact-info", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := middleware.GetTenantID(r)
		if tenantID == 0 {
			http.Error(w, `{"error":"tenant context not found"}`, http.StatusBadRequest)
			return
		}
		var req struct {
			Email string `json:"email"`
			Phone string `json:"phone"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}
		_, err := db.ExecContext(r.Context(), `UPDATE tenants SET email = $1, phone = $2, updated_at = NOW() WHERE id = $3`, req.Email, req.Phone, tenantID)
		if err != nil {
			log.Printf("[TenantContactInfo] ERROR: Failed to update - tenant_id: %d, error: %v\n", tenantID, err)
			http.Error(w, `{"error":"failed to update contact info"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]string{
				"email": req.Email,
				"phone": req.Phone,
			},
		})
	})))

	// HR Module - Employee management endpoints (require authentication + RBAC permission)
	// Module ID 2 = HR (from migrations)
	mux.Handle("GET /api/v1/hr/employees", wrapWithPermission(http.HandlerFunc(hrEmployeeHandler.ListEmployees), 2, "READ"))
	mux.Handle("GET /api/v1/hr/employees/{id}", wrapWithPermission(http.HandlerFunc(hrEmployeeHandler.GetEmployee), 2, "READ"))
	mux.Handle("POST /api/v1/hr/employees", wrapWithPermission(http.HandlerFunc(hrEmployeeHandler.CreateEmployee), 2, "WRITE"))
	mux.Handle("PUT /api/v1/hr/employees/{id}", wrapWithPermission(http.HandlerFunc(hrEmployeeHandler.UpdateEmployee), 2, "WRITE"))
	mux.Handle("DELETE /api/v1/hr/employees/{id}", wrapWithPermission(http.HandlerFunc(hrEmployeeHandler.DeleteEmployee), 2, "DELETE"))

	// HR Module - Roles management endpoints (inline handlers using roles table)
	// Module ID 2 = HR (from migrations)
	mux.Handle("GET /api/v1/hr/roles", wrapWithPermission(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := middleware.GetTenantID(r)
		restaurantID := middleware.GetRestaurantID(r)
		log.Printf("[HR-Roles] ListRoles: tenantID=%d, restaurantID=%d\n", tenantID, restaurantID)
		rows, err := db.QueryContext(r.Context(),
			`SELECT id, tenant_id, restaurant_id, role_name, COALESCE(role_name_ar,''), COALESCE(description,''), COALESCE(description_ar,''), role_code, access_level, can_approve_leaves, can_approve_overtime, can_manage_payroll, can_view_reports, is_active, is_system_role, display_order, created_at, updated_at FROM roles WHERE tenant_id = $1 AND restaurant_id = $2 AND is_active = true ORDER BY display_order ASC`,
			tenantID, restaurantID)
		if err != nil {
			log.Printf("[HR-Roles] ListRoles ERROR: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Failed to fetch roles"})
			return
		}
		defer rows.Close()
		var roles []map[string]interface{}
		for rows.Next() {
			var id, tenantID2, restaurantID2, displayOrder int
			var roleName, roleNameAr, description, descriptionAr, roleCode, accessLevel string
			var canApproveLeaves, canApproveOvertime, canManagePayroll, canViewReports, isActive, isSystemRole bool
			var createdAt, updatedAt time.Time
			if err := rows.Scan(&id, &tenantID2, &restaurantID2, &roleName, &roleNameAr, &description, &descriptionAr, &roleCode, &accessLevel, &canApproveLeaves, &canApproveOvertime, &canManagePayroll, &canViewReports, &isActive, &isSystemRole, &displayOrder, &createdAt, &updatedAt); err != nil {
				log.Printf("[HR-Roles] ListRoles scan ERROR: %v\n", err)
				continue
			}
			roles = append(roles, map[string]interface{}{
				"id": id, "tenant_id": tenantID2, "restaurant_id": restaurantID2,
				"role_name": roleName, "role_name_ar": roleNameAr, "description": description,
				"description_ar": descriptionAr, "role_code": roleCode, "access_level": accessLevel,
				"can_approve_leaves": canApproveLeaves, "can_approve_overtime": canApproveOvertime,
				"can_manage_payroll": canManagePayroll, "can_view_reports": canViewReports,
				"is_active": isActive, "is_system_role": isSystemRole, "display_order": displayOrder,
				"created_at": createdAt, "updated_at": updatedAt,
			})
		}
		if roles == nil {
			roles = []map[string]interface{}{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": roles})
	}), 2, "READ"))

	mux.Handle("POST /api/v1/hr/roles", wrapWithPermission(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := middleware.GetTenantID(r)
		restaurantID := middleware.GetRestaurantID(r)
		var req struct {
			RoleName           string `json:"role_name"`
			RoleNameAr         string `json:"role_name_ar"`
			Description        string `json:"description"`
			DescriptionAr      string `json:"description_ar"`
			RoleCode           string `json:"role_code"`
			AccessLevel        string `json:"access_level"`
			CanApproveLeaves   bool   `json:"can_approve_leaves"`
			CanApproveOvertime bool   `json:"can_approve_overtime"`
			CanManagePayroll   bool   `json:"can_manage_payroll"`
			CanViewReports     bool   `json:"can_view_reports"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid request body"})
			return
		}
		if req.RoleName == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "role_name is required"})
			return
		}
		if req.RoleCode == "" {
			req.RoleCode = "ROLE_" + strconv.FormatInt(time.Now().UnixMilli(), 10)
		}
		if req.AccessLevel == "" {
			req.AccessLevel = "basic"
		}
		var id int
		err := db.QueryRowContext(r.Context(),
			`INSERT INTO roles (tenant_id, restaurant_id, role_name, role_name_ar, description, description_ar, role_code, access_level, can_approve_leaves, can_approve_overtime, can_manage_payroll, can_view_reports, permissions) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, '{}') RETURNING id`,
			tenantID, restaurantID, req.RoleName, req.RoleNameAr, req.Description, req.DescriptionAr, req.RoleCode, req.AccessLevel, req.CanApproveLeaves, req.CanApproveOvertime, req.CanManagePayroll, req.CanViewReports,
		).Scan(&id)
		if err != nil {
			log.Printf("[HR-Roles] CreateRole ERROR: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Failed to create role: " + err.Error()})
			return
		}
		log.Printf("[HR-Roles] CreateRole SUCCESS: id=%d\n", id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": map[string]interface{}{"id": id, "role_name": req.RoleName, "role_code": req.RoleCode}})
	}), 2, "WRITE"))

	mux.Handle("PUT /api/v1/hr/roles/{id}", wrapWithPermission(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := middleware.GetTenantID(r)
		restaurantID := middleware.GetRestaurantID(r)
		roleID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid role ID"})
			return
		}
		var req struct {
			RoleName           string `json:"role_name"`
			RoleNameAr         string `json:"role_name_ar"`
			Description        string `json:"description"`
			DescriptionAr      string `json:"description_ar"`
			AccessLevel        string `json:"access_level"`
			CanApproveLeaves   *bool  `json:"can_approve_leaves"`
			CanApproveOvertime *bool  `json:"can_approve_overtime"`
			CanManagePayroll   *bool  `json:"can_manage_payroll"`
			CanViewReports     *bool  `json:"can_view_reports"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid request body"})
			return
		}
		_, err = db.ExecContext(r.Context(),
			`UPDATE roles SET role_name = COALESCE(NULLIF($1,''), role_name), role_name_ar = $2, description = $3, description_ar = $4, access_level = COALESCE(NULLIF($5,''), access_level), updated_at = CURRENT_TIMESTAMP WHERE id = $6 AND tenant_id = $7 AND restaurant_id = $8`,
			req.RoleName, req.RoleNameAr, req.Description, req.DescriptionAr, req.AccessLevel, roleID, tenantID, restaurantID)
		if err != nil {
			log.Printf("[HR-Roles] UpdateRole ERROR: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Failed to update role"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	}), 2, "WRITE"))

	mux.Handle("DELETE /api/v1/hr/roles/{id}", wrapWithPermission(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := middleware.GetTenantID(r)
		restaurantID := middleware.GetRestaurantID(r)
		roleID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid role ID"})
			return
		}
		_, err = db.ExecContext(r.Context(),
			`UPDATE roles SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3 AND is_system_role = false`,
			roleID, tenantID, restaurantID)
		if err != nil {
			log.Printf("[HR-Roles] DeleteRole ERROR: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Failed to delete role"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	}), 2, "DELETE"))

	// RBAC endpoints
	if rbacRoleHandler != nil && rbacPermissionHandler != nil && rbacUserRoleHandler != nil {
		// Role management
		mux.Handle("GET /api/v1/rbac/roles", wrapProtected(http.HandlerFunc(rbacRoleHandler.ListRoles)))
		mux.Handle("GET /api/v1/rbac/roles/{id}", wrapProtected(http.HandlerFunc(rbacRoleHandler.GetRole)))
		mux.Handle("POST /api/v1/rbac/roles", wrapProtected(http.HandlerFunc(rbacRoleHandler.CreateRole)))
		mux.Handle("PUT /api/v1/rbac/roles/{id}", wrapProtected(http.HandlerFunc(rbacRoleHandler.UpdateRole)))
		mux.Handle("DELETE /api/v1/rbac/roles/{id}", wrapProtected(http.HandlerFunc(rbacRoleHandler.DeleteRole)))

		// Permission management
		mux.Handle("POST /api/v1/rbac/permissions/assign", wrapProtected(http.HandlerFunc(rbacPermissionHandler.AssignPermission)))
		mux.Handle("DELETE /api/v1/rbac/permissions/{roleId}/{moduleId}", wrapProtected(http.HandlerFunc(rbacPermissionHandler.RevokePermission)))
		mux.Handle("GET /api/v1/rbac/roles/{roleId}/permissions", wrapProtected(http.HandlerFunc(rbacPermissionHandler.GetRolePermissions)))
		mux.Handle("GET /api/v1/rbac/check-permission", wrapProtected(http.HandlerFunc(rbacPermissionHandler.CheckPermission)))
		mux.Handle("GET /api/v1/rbac/me/permissions", wrapProtected(http.HandlerFunc(rbacPermissionHandler.GetUserPermissions)))

		// User-role assignment
		mux.Handle("POST /api/v1/rbac/users/{userId}/roles", wrapProtected(http.HandlerFunc(rbacUserRoleHandler.AssignRoleToUser)))
		mux.Handle("DELETE /api/v1/rbac/users/{userId}/roles/{roleId}", wrapProtected(http.HandlerFunc(rbacUserRoleHandler.RemoveRoleFromUser)))
		mux.Handle("GET /api/v1/rbac/users/{userId}/roles", wrapProtected(http.HandlerFunc(rbacUserRoleHandler.GetUserRoles)))
		mux.Handle("GET /api/v1/rbac/me/roles", wrapProtected(http.HandlerFunc(rbacUserRoleHandler.GetMyRoles)))
		mux.Handle("POST /api/v1/rbac/users-by-email/{email}/roles", wrapProtected(http.HandlerFunc(rbacUserRoleHandler.AssignRoleToUserByEmail)))
		mux.Handle("DELETE /api/v1/rbac/users-by-email/{email}/roles/{roleId}", wrapProtected(http.HandlerFunc(rbacUserRoleHandler.RemoveRoleFromUserByEmail)))
		mux.Handle("GET /api/v1/rbac/users-by-email/{email}/roles", wrapProtected(http.HandlerFunc(rbacUserRoleHandler.GetUserRolesByEmail)))
	}

	// HR Module - Attendance management endpoints (require authentication + RBAC permission)
	// Module ID 2 = HR (from migrations)
	mux.Handle("GET /api/v1/hr/attendance", wrapWithPermission(http.HandlerFunc(hrAttendanceHandler.ListAttendance), 2, "READ"))
	mux.Handle("GET /api/v1/hr/attendance/{id}", wrapWithPermission(http.HandlerFunc(hrAttendanceHandler.GetAttendance), 2, "READ"))
	mux.Handle("POST /api/v1/hr/attendance", wrapWithPermission(http.HandlerFunc(hrAttendanceHandler.CreateAttendance), 2, "WRITE"))
	mux.Handle("POST /api/v1/hr/attendance/clock-in", wrapWithPermission(http.HandlerFunc(hrAttendanceHandler.ClockIn), 2, "WRITE"))
	mux.Handle("PUT /api/v1/hr/attendance/{id}", wrapWithPermission(http.HandlerFunc(hrAttendanceHandler.UpdateAttendance), 2, "WRITE"))
	mux.Handle("DELETE /api/v1/hr/attendance/{id}", wrapWithPermission(http.HandlerFunc(hrAttendanceHandler.DeleteAttendance), 2, "DELETE"))

	// HR Module - Leave management endpoints (require authentication + RBAC permission)
	// Module ID 2 = HR (from migrations)
	mux.Handle("GET /api/v1/hr/leaves", wrapWithPermission(http.HandlerFunc(hrLeaveHandler.ListLeaves), 2, "READ"))
	mux.Handle("GET /api/v1/hr/leaves/{id}", wrapWithPermission(http.HandlerFunc(hrLeaveHandler.GetLeave), 2, "READ"))
	mux.Handle("POST /api/v1/hr/leaves", wrapWithPermission(http.HandlerFunc(hrLeaveHandler.CreateLeave), 2, "WRITE"))
	mux.Handle("PUT /api/v1/hr/leaves/{id}", wrapWithPermission(http.HandlerFunc(hrLeaveHandler.UpdateLeave), 2, "WRITE"))
	mux.Handle("DELETE /api/v1/hr/leaves/{id}", wrapWithPermission(http.HandlerFunc(hrLeaveHandler.DeleteLeave), 2, "DELETE"))
	mux.Handle("POST /api/v1/hr/leaves/{id}/approve", wrapWithPermission(http.HandlerFunc(hrLeaveHandler.ApproveLeave), 2, "WRITE"))
	mux.Handle("POST /api/v1/hr/leaves/{id}/reject", wrapWithPermission(http.HandlerFunc(hrLeaveHandler.RejectLeave), 2, "WRITE"))

	// HR Module - Salary/Payroll management endpoints (require authentication + RBAC permission)
	// Module ID 2 = HR (from migrations)
	mux.Handle("GET /api/v1/hr/salaries", wrapWithPermission(http.HandlerFunc(hrSalaryHandler.ListSalaries), 2, "READ"))
	mux.Handle("GET /api/v1/hr/salaries/{id}", wrapWithPermission(http.HandlerFunc(hrSalaryHandler.GetSalary), 2, "READ"))
	mux.Handle("POST /api/v1/hr/salaries", wrapWithPermission(http.HandlerFunc(hrSalaryHandler.CreateSalary), 2, "WRITE"))
	mux.Handle("PUT /api/v1/hr/salaries/{id}", wrapWithPermission(http.HandlerFunc(hrSalaryHandler.UpdateSalary), 2, "WRITE"))
	mux.Handle("DELETE /api/v1/hr/salaries/{id}", wrapWithPermission(http.HandlerFunc(hrSalaryHandler.DeleteSalary), 2, "DELETE"))
	mux.Handle("POST /api/v1/hr/salaries/{id}/mark-paid", wrapWithPermission(http.HandlerFunc(hrSalaryHandler.MarkAsPaid), 2, "WRITE"))

	log.Println("✅ HR Module routes registered (employees, attendance, leaves, salaries, roles)")

	// User Settings Management endpoints (require authentication)
	// Use stub handlers when db is nil
	if db == nil {
		mux.Handle("GET /api/v1/users", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true,"data":[{"id":7,"email":"yoseabdallah866@gmail.com","name":"Admin User","role":"owner","created_at":"2024-01-01T00:00:00Z"}]}`))
		})))
	} else {
		mux.Handle("GET /api/v1/users", wrapProtected(http.HandlerFunc(userSettingsHandler.ListUsers)))
	}
	mux.Handle("GET /api/v1/users/{userId}/settings", wrapProtected(http.HandlerFunc(userSettingsHandler.GetUserSettings)))
	mux.Handle("PUT /api/v1/users/{userId}/settings/language", wrapProtected(http.HandlerFunc(userSettingsHandler.UpdateLanguage)))
	mux.Handle("PUT /api/v1/users/{userId}/settings/theme", wrapProtected(http.HandlerFunc(userSettingsHandler.UpdateTheme)))
	mux.Handle("PUT /api/v1/users/{userId}/settings/theme-colors", wrapProtected(http.HandlerFunc(userSettingsHandler.UpdateColors)))
	mux.Handle("POST /api/v1/users/{userId}/settings/change-password", wrapProtected(http.HandlerFunc(userSettingsHandler.ChangePassword)))
	mux.Handle("GET /api/v1/users/{userId}/profile", wrapProtected(http.HandlerFunc(userSettingsHandler.GetUserProfile)))
	mux.Handle("PUT /api/v1/users/{userId}/profile", wrapProtected(http.HandlerFunc(userSettingsHandler.UpdateProfile)))

	// User Management endpoints (create, update, delete users in tenant)
	if db == nil {
		// Stub handlers for user management when db is nil
		mux.Handle("POST /api/v1/users", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Phone    string `json:"phone"`
				Password string `json:"password"`
			}
			json.Unmarshal(body, &req)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id":         8,
					"email":      req.Email,
					"name":       req.Name,
					"phone":      req.Phone,
					"status":     "active",
					"created_at": time.Now(),
				},
			})
		})))

		mux.Handle("PUT /api/v1/users/{id}", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			body, _ := io.ReadAll(r.Body)
			var req struct {
				Name  string `json:"name"`
				Phone string `json:"phone"`
			}
			json.Unmarshal(body, &req)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id":         id,
					"name":       req.Name,
					"phone":      req.Phone,
					"updated_at": time.Now(),
				},
			})
		})))

		mux.Handle("DELETE /api/v1/users/{id}", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
			})
		})))
	} else {
		// Real database handlers for user management
		mux.Handle("POST /api/v1/users", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID := middleware.GetTenantID(r)
			restaurantID := middleware.GetRestaurantID(r)
			assignedBy := middleware.GetUserID(r)
			if tenantID == 0 {
				http.Error(w, "Missing tenant context", http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}

			var req struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Phone    string `json:"phone"`
				Password string `json:"password"`
			}
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Hash password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Failed to hash password", http.StatusInternalServerError)
				return
			}

			// Create user
			var restaurantIDPtr *int
			if restaurantID > 0 {
				rid := int(restaurantID)
				restaurantIDPtr = &rid
			}
			user, err := userRepo.CreateUserInTenant(r.Context(), tenantID, req.Name, req.Email, req.Phone, string(hashedPassword), restaurantIDPtr)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
				return
			}

			if assignedBy == 0 {
				assignedBy = int64(user.ID)
			}

			roleRestaurantID := restaurantID
			if roleRestaurantID == 0 && user.RestaurantID != nil {
				roleRestaurantID = int64(*user.RestaurantID)
			}
			if rbacRoleRepo != nil && userRoleRepo != nil && roleRestaurantID > 0 {
				role, err := rbacRoleRepo.GetByCode(r.Context(), tenantID, roleRestaurantID, "EMPLOYEE")
				if err == nil && role != nil {
					_ = userRoleRepo.Assign(r.Context(), tenantID, int64(user.ID), role.ID, assignedBy)
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id":         user.ID,
					"email":      user.Email,
					"name":       user.Name,
					"phone":      user.Phone,
					"status":     user.Status,
					"created_at": user.CreatedAt,
				},
			})
		})))

		mux.Handle("PUT /api/v1/users/{id}", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID := middleware.GetTenantID(r)
			if tenantID == 0 {
				http.Error(w, "Missing tenant context", http.StatusBadRequest)
				return
			}

			idStr := r.PathValue("id")
			userID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}

			var req struct {
				Name  string `json:"name"`
				Phone string `json:"phone"`
			}
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Update user
			var namePtr, phonePtr *string
			if req.Name != "" {
				namePtr = &req.Name
			}
			if req.Phone != "" {
				phonePtr = &req.Phone
			}

			user, err := userRepo.UpdateUser(r.Context(), userID, tenantID, namePtr, phonePtr)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to update user: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id":         user.ID,
					"email":      user.Email,
					"name":       user.Name,
					"phone":      user.Phone,
					"updated_at": user.UpdatedAt,
				},
			})
		})))

		mux.Handle("DELETE /api/v1/users/{id}", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID := middleware.GetTenantID(r)
			if tenantID == 0 {
				http.Error(w, "Missing tenant context", http.StatusBadRequest)
				return
			}

			idStr := r.PathValue("id")
			userID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			// Delete user
			err = userRepo.DeleteUser(r.Context(), userID, tenantID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to delete user: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
			})
		})))
	}

	// Admin Order Management endpoints (require authentication)
	// Module ID 10 = Orders (from migrations)
	mux.Handle("GET /api/v1/admin/orders/stats", wrapWithPermission(http.HandlerFunc(adminOrderHandler.GetOrderStats), 10, "READ"))
	mux.Handle("GET /api/v1/admin/orders/export", wrapWithPermission(http.HandlerFunc(adminOrderHandler.ExportOrders), 10, "READ"))
	mux.Handle("GET /api/v1/admin/orders/{id}/history", wrapWithPermission(http.HandlerFunc(adminOrderHandler.GetOrderHistory), 10, "READ"))
	mux.Handle("PUT /api/v1/admin/orders/{id}/status", wrapWithPermission(http.HandlerFunc(adminOrderHandler.UpdateOrderStatus), 10, "WRITE"))
	mux.Handle("PUT /api/v1/admin/orders/{id}/payment-status", wrapWithPermission(http.HandlerFunc(adminOrderHandler.UpdatePaymentStatus), 10, "WRITE"))
	mux.Handle("DELETE /api/v1/admin/orders/{id}", wrapWithPermission(http.HandlerFunc(adminOrderHandler.CancelOrder), 10, "DELETE"))
	mux.Handle("GET /api/v1/admin/orders/{id}", wrapWithPermission(http.HandlerFunc(adminOrderHandler.GetOrder), 10, "READ"))
	mux.Handle("GET /api/v1/admin/orders", wrapWithPermission(http.HandlerFunc(adminOrderHandler.ListOrders), 10, "READ"))

	// Admin Driver Management endpoints (require authentication)
	mux.Handle("POST /api/v1/admin/drivers", wrapProtected(http.HandlerFunc(adminDriverHandler.CreateDriver)))
	mux.Handle("GET /api/v1/admin/drivers", wrapProtected(http.HandlerFunc(adminDriverHandler.ListDrivers)))
	mux.Handle("GET /api/v1/admin/drivers/available", wrapProtected(http.HandlerFunc(adminDriverHandler.GetAvailableDrivers)))
	mux.Handle("GET /api/v1/admin/drivers/{id}", wrapProtected(http.HandlerFunc(adminDriverHandler.GetDriver)))
	mux.Handle("PUT /api/v1/admin/drivers/{id}", wrapProtected(http.HandlerFunc(adminDriverHandler.UpdateDriver)))
	mux.Handle("DELETE /api/v1/admin/drivers/{id}", wrapProtected(http.HandlerFunc(adminDriverHandler.DeleteDriver)))
	mux.Handle("PUT /api/v1/admin/drivers/{id}/status", wrapProtected(http.HandlerFunc(adminDriverHandler.UpdateDriverStatus)))
	mux.Handle("POST /api/v1/admin/orders/{orderId}/assign", wrapWithPermission(http.HandlerFunc(adminDriverHandler.AssignOrderToDriver), 10, "WRITE"))
	// 	mux.Handle("GET /api/v1/admin/drivers/{id}/stats", wrapProtected(http.HandlerFunc(adminDriverHandler.GetDriverStats)))
	// 	mux.Handle("POST /api/v1/admin/drivers/{id}/location", wrapProtected(http.HandlerFunc(adminDriverHandler.UpdateDriverLocation)))
	// 	mux.Handle("GET /api/v1/admin/drivers/{id}/location", wrapProtected(http.HandlerFunc(adminDriverHandler.GetDriverLocation)))

	// Notification management endpoints (require authentication)
	// Register specific routes BEFORE {id} pattern routes to avoid conflicts
	// Use stub handlers when db is nil
	if db == nil {
		mux.Handle("GET /api/v1/notifications", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true,"data":[],"total":0}`))
		})))
		mux.Handle("GET /api/v1/notifications/stats", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true,"data":{"total":0,"unread":0,"read":0}}`))
		})))
		mux.Handle("GET /api/v1/notifications/unread-count", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true,"count":0}`))
		})))
		mux.Handle("POST /api/v1/notifications", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true,"data":{"id":1}}`))
		})))
		mux.Handle("POST /api/v1/notifications/mark-all-read", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true}`))
		})))
		mux.Handle("POST /api/v1/notifications/{id}/read", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true}`))
		})))
		mux.Handle("POST /api/v1/notifications/{id}/unread", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true}`))
		})))
		mux.Handle("GET /api/v1/notifications/{id}", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true,"data":null}`))
		})))
		mux.Handle("DELETE /api/v1/notifications/{id}", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true}`))
		})))
	} else {
		mux.Handle("GET /api/v1/notifications", wrapProtected(http.HandlerFunc(notificationHandler.ListNotifications)))
		mux.Handle("GET /api/v1/notifications/stats", wrapProtected(http.HandlerFunc(notificationHandler.GetStats)))
		mux.Handle("GET /api/v1/notifications/unread-count", wrapProtected(http.HandlerFunc(notificationHandler.GetUnreadCount)))
		mux.Handle("POST /api/v1/notifications", wrapProtected(http.HandlerFunc(notificationHandler.CreateNotification)))
		mux.Handle("POST /api/v1/notifications/mark-all-read", wrapProtected(http.HandlerFunc(notificationHandler.MarkAllAsRead)))
		mux.Handle("POST /api/v1/notifications/{id}/read", wrapProtected(http.HandlerFunc(notificationHandler.MarkAsRead)))
		mux.Handle("POST /api/v1/notifications/{id}/unread", wrapProtected(http.HandlerFunc(notificationHandler.MarkAsUnread)))
		mux.Handle("GET /api/v1/notifications/{id}", wrapProtected(http.HandlerFunc(notificationHandler.GetNotification)))
		mux.Handle("DELETE /api/v1/notifications/{id}", wrapProtected(http.HandlerFunc(notificationHandler.DeleteNotification)))
	}

	// Theme management endpoints (require authentication)
	mux.Handle("GET /api/v1/admin/themes", wrapProtected(http.HandlerFunc(themeHandler.ListThemes)))
	mux.Handle("POST /api/v1/admin/themes", wrapProtected(http.HandlerFunc(themeHandler.CreateTheme)))
	mux.Handle("GET /api/v1/admin/themes/{id}", wrapProtected(http.HandlerFunc(themeHandler.GetTheme)))
	mux.Handle("PUT /api/v1/admin/themes/{id}", wrapProtected(http.HandlerFunc(themeHandler.UpdateTheme)))
	mux.Handle("DELETE /api/v1/admin/themes/{id}", wrapProtected(http.HandlerFunc(themeHandler.DeleteTheme)))
	mux.Handle("POST /api/v1/admin/themes/{id}/activate", wrapProtected(http.HandlerFunc(themeHandler.ActivateTheme)))
	mux.Handle("POST /api/v1/admin/themes/{id}/duplicate", wrapProtected(http.HandlerFunc(themeHandler.DuplicateTheme)))
	mux.Handle("GET /api/v1/admin/themes/{id}/history", wrapProtected(http.HandlerFunc(themeHandler.GetThemeHistory)))

	// Theme presets endpoints (V2 - with database support)
	mux.Handle("GET /api/v1/admin/theme-presets", wrapProtected(http.HandlerFunc(themeHandlerV2.ListThemePresets)))

	// RBAC stub endpoints (returns mock data for development)
	if db == nil {
		mux.Handle("GET /api/v1/rbac/roles", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rbacStorage.mu.RLock()
			defer rbacStorage.mu.RUnlock()

			roles := make([]*Role, 0, len(rbacStorage.roles))
			for _, role := range rbacStorage.roles {
				roles = append(roles, role)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    roles,
			})
		})))
		mux.Handle("POST /api/v1/rbac/roles", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}

			var req struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			rbacStorage.mu.Lock()
			defer rbacStorage.mu.Unlock()

			newRole := &Role{
				ID:               rbacStorage.nextID,
				TenantID:         1,
				RoleName:         req.Name,
				Name:             req.Name,
				Description:      req.Description,
				IsSystem:         false,
				PermissionsCount: 0,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}
			rbacStorage.roles[rbacStorage.nextID] = newRole
			rbacStorage.nextID++

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    newRole,
			})
		})))
		mux.Handle("GET /api/v1/rbac/roles/{id}", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idStr := r.PathValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid role ID", http.StatusBadRequest)
				return
			}

			rbacStorage.mu.RLock()
			defer rbacStorage.mu.RUnlock()

			role, exists := rbacStorage.roles[id]
			if !exists {
				http.Error(w, "Role not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    role,
			})
		})))
		mux.Handle("PUT /api/v1/rbac/roles/{id}", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idStr := r.PathValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid role ID", http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}

			var req struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			rbacStorage.mu.Lock()
			defer rbacStorage.mu.Unlock()

			role, exists := rbacStorage.roles[id]
			if !exists {
				http.Error(w, "Role not found", http.StatusNotFound)
				return
			}

			// Don't allow editing system roles
			if role.IsSystem {
				http.Error(w, "Cannot edit system roles", http.StatusForbidden)
				return
			}

			role.Name = req.Name
			role.RoleName = req.Name
			role.Description = req.Description
			role.UpdatedAt = time.Now()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    role,
			})
		})))
		mux.Handle("DELETE /api/v1/rbac/roles/{id}", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idStr := r.PathValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid role ID", http.StatusBadRequest)
				return
			}

			rbacStorage.mu.Lock()
			defer rbacStorage.mu.Unlock()

			role, exists := rbacStorage.roles[id]
			if !exists {
				http.Error(w, "Role not found", http.StatusNotFound)
				return
			}

			// Don't allow deleting system roles
			if role.IsSystem {
				http.Error(w, "Cannot delete system roles", http.StatusForbidden)
				return
			}

			delete(rbacStorage.roles, id)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
			})
		})))
		mux.Handle("GET /api/v1/rbac/me/roles", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			// Mock admin role for current user
			w.Write([]byte(`{"success":true,"data":[{"id":1,"tenant_id":1,"role_name":"Admin","description":"System administrator with full access","is_system":true,"created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}]}`))
		})))
		mux.Handle("GET /api/v1/rbac/users/{userId}/roles", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIdStr := r.PathValue("userId")
			userId, err := strconv.Atoi(userIdStr)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			rbacStorage.mu.RLock()
			defer rbacStorage.mu.RUnlock()

			roleIds := rbacStorage.userRoles[userId]
			roles := make([]*Role, 0, len(roleIds))
			for _, roleId := range roleIds {
				if role, exists := rbacStorage.roles[roleId]; exists {
					roles = append(roles, role)
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    roles,
			})
		})))
		mux.Handle("POST /api/v1/rbac/users/{userId}/roles", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIdStr := r.PathValue("userId")
			userId, err := strconv.Atoi(userIdStr)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}

			var req struct {
				RoleId int `json:"roleId"`
			}
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			rbacStorage.mu.Lock()
			defer rbacStorage.mu.Unlock()

			// Check if role exists
			if _, exists := rbacStorage.roles[req.RoleId]; !exists {
				http.Error(w, "Role not found", http.StatusNotFound)
				return
			}

			// Add role to user if not already assigned
			userRoles := rbacStorage.userRoles[userId]
			for _, roleId := range userRoles {
				if roleId == req.RoleId {
					// Already assigned
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]interface{}{
						"success": true,
					})
					return
				}
			}
			rbacStorage.userRoles[userId] = append(userRoles, req.RoleId)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
			})
		})))
		mux.Handle("DELETE /api/v1/rbac/users/{userId}/roles/{roleId}", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIdStr := r.PathValue("userId")
			userId, err := strconv.Atoi(userIdStr)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			roleIdStr := r.PathValue("roleId")
			roleId, err := strconv.Atoi(roleIdStr)
			if err != nil {
				http.Error(w, "Invalid role ID", http.StatusBadRequest)
				return
			}

			rbacStorage.mu.Lock()
			defer rbacStorage.mu.Unlock()

			userRoles := rbacStorage.userRoles[userId]
			newRoles := make([]int, 0, len(userRoles))
			for _, rId := range userRoles {
				if rId != roleId {
					newRoles = append(newRoles, rId)
				}
			}
			rbacStorage.userRoles[userId] = newRoles

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
			})
		})))
		mux.Handle("GET /api/v1/rbac/me/permissions", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			// Mock admin permissions for all modules
			w.Write([]byte(`{"success":true,"data":{"1":"ADMIN","2":"ADMIN","3":"ADMIN","4":"ADMIN","5":"ADMIN","6":"ADMIN"}}`))
		})))
		mux.Handle("GET /api/v1/rbac/check-permission", wrapProtected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			// Return proper permission check response with has_permission (snake_case)
			w.Write([]byte(`{"success":true,"has_permission":true,"module_id":6,"permission_level":"ADMIN"}`))
		})))
	}
	// ============================================================
	// SUPER ADMIN ROUTES (Phase 6 - Super Admin Dashboard)
	// These routes are completely isolated from tenant routes
	// ============================================================
	if superAdminHandler != nil && impersonationHandler != nil && superAdminTokenService != nil {
		// Helper function to wrap super admin routes with auth middleware
		wrapSuperAdmin := func(next http.Handler) http.Handler {
			wrapped := middleware.RequireSuperAdmin(superAdminRepo)(next)
			wrapped = middleware.ValidateIPWhitelist(superAdminRepo)(wrapped)
			wrapped = middleware.SuperAdminAuthMiddleware(superAdminTokenService)(wrapped)
			wrapped = middleware.RateLimitSuperAdmin(30)(wrapped) // 30 requests per minute
			return wrapped
		}

		// Super Admin Authentication (public - no auth required)
		mux.HandleFunc("POST /api/v1/super-admin/auth/login", superAdminHandler.Login)

		// Super Admin Protected Routes
		mux.Handle("POST /api/v1/super-admin/auth/logout", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.Logout)))
		mux.Handle("GET /api/v1/super-admin/dashboard/stats", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.GetDashboardStats)))
		if paymentMonitoringHandler != nil {
			mux.Handle("GET /api/v1/super-admin/payment-metrics", wrapSuperAdmin(http.HandlerFunc(paymentMonitoringHandler.GetPaymentMetrics)))
			mux.Handle("GET /api/v1/super-admin/webhook-health", wrapSuperAdmin(http.HandlerFunc(paymentMonitoringHandler.GetWebhookHealth)))
		} else {
			log.Println("[SuperAdmin] Payment monitoring routes skipped: paymentMonitoringHandler is nil")
		}
		if paymentHandler != nil {
			mux.Handle("POST /api/v1/super-admin/payments/{id}/refund", wrapSuperAdmin(http.HandlerFunc(paymentHandler.RefundPayment)))
		} else {
			log.Println("[SuperAdmin] Refund route skipped: paymentHandler is nil")
		}

		// Tenant Management
		mux.Handle("GET /api/v1/super-admin/tenants", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.ListTenants)))
		mux.Handle("GET /api/v1/super-admin/tenants/{id}", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.GetTenant)))
		mux.Handle("GET /api/v1/super-admin/tenants/{id}/users", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.GetTenantUsers)))
		mux.Handle("POST /api/v1/super-admin/tenants/{id}/dashboard-login", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.GenerateDashboardLoginToken)))
		mux.Handle("PUT /api/v1/super-admin/tenants/{id}/disable", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.DisableTenant)))
		mux.Handle("PUT /api/v1/super-admin/tenants/{id}/enable", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.EnableTenant)))

		// User Management (across all tenants)
		mux.Handle("GET /api/v1/super-admin/users", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.ListUsers)))
		mux.Handle("GET /api/v1/super-admin/users/{id}", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.GetUser)))

		// Impersonation
		mux.Handle("POST /api/v1/super-admin/impersonate/{userID}", wrapSuperAdmin(http.HandlerFunc(impersonationHandler.StartImpersonation)))
		mux.Handle("POST /api/v1/super-admin/impersonate/end", wrapSuperAdmin(http.HandlerFunc(impersonationHandler.EndImpersonation)))
		mux.Handle("GET /api/v1/super-admin/impersonate/status", wrapSuperAdmin(http.HandlerFunc(impersonationHandler.GetImpersonationStatus)))

		// Audit Logs
		mux.Handle("GET /api/v1/super-admin/audit-logs", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.GetAuditLogs)))

		// Settings
		mux.Handle("PUT /api/v1/super-admin/settings/ip-whitelist", wrapSuperAdmin(http.HandlerFunc(superAdminHandler.UpdateIPWhitelist)))

		// Public endpoint for validating impersonation token (called from tenant dashboard)
		mux.HandleFunc("GET /api/v1/impersonation/validate", impersonationHandler.ValidateImpersonationToken)
		mux.HandleFunc("POST /api/v1/impersonation/end", impersonationHandler.EndImpersonationByToken)
		mux.HandleFunc("GET /api/v1/impersonation/check", impersonationHandler.CheckImpersonation)

		log.Println("✅ Super Admin routes registered")
	} else {
		log.Printf("[SuperAdmin] Routes skipped: superAdminHandler=%t impersonationHandler=%t superAdminTokenService=%t", superAdminHandler != nil, impersonationHandler != nil, superAdminTokenService != nil)
	}

	// Super Admin Subscription Management Routes (Phase 1-8 - SaaS Trial/Subscription System)
	if superAdminHandler != nil && subscriptionHandler != nil && superAdminTokenService != nil {
		// Helper function to wrap super admin routes with auth middleware (reuse from above)
		wrapSuperAdmin := func(next http.Handler) http.Handler {
			// Apply middleware in REVERSE order of execution
			wrapped := middleware.RequireSuperAdmin(superAdminRepo)(next)
			wrapped = middleware.ValidateIPWhitelist(superAdminRepo)(wrapped)
			wrapped = middleware.SuperAdminAuthMiddleware(superAdminTokenService)(wrapped)
			wrapped = middleware.RateLimitSuperAdmin(30)(wrapped) // 30 requests per minute
			return wrapped
		}

		// Subscription management endpoints (super admin only)
		mux.Handle("GET /api/v1/super-admin/tenants/{id}/subscription", wrapSuperAdmin(http.HandlerFunc(subscriptionHandler.GetTenantSubscriptionByID)))
		mux.Handle("POST /api/v1/super-admin/tenants/{id}/extend-trial", wrapSuperAdmin(http.HandlerFunc(subscriptionHandler.ExtendTenantTrial)))
		mux.Handle("POST /api/v1/super-admin/tenants/{id}/disable", wrapSuperAdmin(http.HandlerFunc(subscriptionHandler.DisableTenant)))
		mux.Handle("POST /api/v1/super-admin/tenants/{id}/enable", wrapSuperAdmin(http.HandlerFunc(subscriptionHandler.EnableTenant)))
		mux.Handle("DELETE /api/v1/super-admin/tenants/{id}", wrapSuperAdmin(http.HandlerFunc(subscriptionHandler.DeleteTenant)))
		mux.Handle("GET /api/v1/super-admin/tenants/{id}/subscription-history", wrapSuperAdmin(http.HandlerFunc(subscriptionHandler.GetSubscriptionHistory)))

		log.Println("✅ Subscription management routes registered")
	}

	// Apply global middleware
	corsHandler := middleware.CORSMiddleware(mux)
	compressedHandler := middleware.GzipMiddleware(corsHandler)

	// Start server
	log.Printf("🚀 Server started on port %s", cfg.Server.Port)
	log.Printf("📍 API URL: http://localhost:%s/api/v1", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, compressedHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
