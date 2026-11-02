package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
)

// loginConfirmStore is the persistence needed by LoginConfirm.
type loginConfirmStore interface {
	FindUserByEmailAllTenants(email string) ([]domain.User, error)
	UpdateLastLogin(userID int) error
}

// selectionTokenTTL bounds how long a multi-tenant login proof stays valid.
const selectionTokenTTL = 5 * time.Minute

type AuthUseCase struct {
	repo           *repository.AuthRepository
	superAdminRepo *repository.SuperAdminRepository
	tokenService   *jwt.TokenService
	confirmStore   loginConfirmStore

	// usedSelections remembers consumed selection-token IDs until they expire.
	// It is per-process; running several API replicas needs a shared store.
	selectionMu    sync.Mutex
	usedSelections map[string]time.Time
}

func NewAuthUseCase(repo *repository.AuthRepository, tokenService *jwt.TokenService) *AuthUseCase {
	uc := &AuthUseCase{
		repo:         repo,
		tokenService: tokenService,
	}
	if repo != nil {
		uc.confirmStore = repo
	}
	return uc
}

// consumeSelection marks a selection token ID as used; false means it was already used.
func (u *AuthUseCase) consumeSelection(id string, expires time.Time) bool {
	u.selectionMu.Lock()
	defer u.selectionMu.Unlock()
	if u.usedSelections == nil {
		u.usedSelections = make(map[string]time.Time)
	}
	now := time.Now()
	for k, exp := range u.usedSelections {
		if now.After(exp) {
			delete(u.usedSelections, k)
		}
	}
	if _, used := u.usedSelections[id]; used {
		return false
	}
	u.usedSelections[id] = expires
	return true
}

// SetSuperAdminRepo sets the super admin repository for token exchange
func (u *AuthUseCase) SetSuperAdminRepo(repo *repository.SuperAdminRepository) {
	u.superAdminRepo = repo
}

// Register creates a new tenant, restaurant, and admin user
func (u *AuthUseCase) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Validate input
	if req.Name == "" || req.Email == "" || req.Password == "" || req.RestaurantName == "" {
		return nil, errors.New("all fields are required")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	validator := NewSubdomainValidator()
	slugProvided := strings.TrimSpace(req.RestaurantSlug) != ""
	if slugProvided {
		formattedSlug, err := validator.ValidateAndFormat(req.RestaurantSlug)
		if err != nil {
			return nil, err
		}
		req.RestaurantSlug = formattedSlug

		exists, err := u.repo.CheckTenantSlugExists(req.RestaurantSlug)
		if err != nil {
			return nil, fmt.Errorf("failed to validate subdomain: %w", err)
		}
		if exists {
			return nil, errors.New("subdomain already exists")
		}
	} else {
		req.RestaurantSlug = validator.ToSlug(req.RestaurantName)
	}

	// Check if email already exists
	exists, err := u.repo.CheckEmailExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Create tenant with user
	user, err := u.repo.CreateTenantWithUser(req)
	if err != nil {
		if errors.Is(err, repository.ErrTenantSlugExists) {
			return nil, errors.New("subdomain already exists")
		}
		if errors.Is(err, repository.ErrTenantNameExists) {
			return nil, errors.New("restaurant name already exists")
		}
		if errors.Is(err, repository.ErrTenantDomainExists) {
			return nil, errors.New("restaurant domain already exists")
		}
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	// Generate JWT token
	token, err := u.tokenService.GenerateToken(user.ID, user.TenantID, user.RestaurantID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Clear sensitive data
	user.PasswordHash = ""

	return &domain.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// Login authenticates a user and detects if they have multiple tenant accounts
func (u *AuthUseCase) Login(req *domain.LoginRequest) (interface{}, error) {
	fmt.Println("=== LOGIN USE CASE DEBUG START ===")
	fmt.Printf("Login attempt for email: %s\n", req.Email)

	// Validate input
	if req.Email == "" || req.Password == "" {
		fmt.Println("ERROR: Email or password is empty")
		return nil, errors.New("email and password are required")
	}

	// Find ALL users with this email across all tenants
	fmt.Printf("Looking up all users with email: %s\n", req.Email)
	users, err := u.repo.FindUserByEmailAllTenants(req.Email)
	if err != nil {
		fmt.Printf("ERROR: Database query failed: %v\n", err)
		return nil, errors.New("invalid email or password")
	}

	if len(users) == 0 {
		fmt.Printf("ERROR: No users found with email: %s\n", req.Email)
		return nil, errors.New("invalid email or password")
	}

	fmt.Printf("Found %d users with email %s\n", len(users), req.Email)

	// Verify password against the first user (all have same email, should have same password)
	firstUser := &users[0]
	fmt.Printf("Verifying password for first user - ID: %d, Email: %s\n", firstUser.ID, firstUser.Email)

	// Verify password
	fmt.Println("Comparing password with bcrypt...")
	err = bcrypt.CompareHashAndPassword([]byte(firstUser.PasswordHash), []byte(req.Password))
	if err != nil {
		fmt.Printf("ERROR: Password comparison failed: %v\n", err)
		return nil, errors.New("invalid email or password")
	}
	fmt.Println("SUCCESS: Password verified successfully!")

	// Case 1: Single tenant - auto-login
	if len(users) == 1 {
		fmt.Println("Single tenant found - auto-logging in")
		user := &users[0]

		// Update last login
		fmt.Printf("Updating last login for user ID: %d\n", user.ID)
		_ = u.repo.UpdateLastLogin(user.ID)

		// Generate JWT token
		fmt.Println("Generating JWT token...")
		token, err := u.tokenService.GenerateToken(user.ID, user.TenantID, user.RestaurantID, user.Email, user.Role)
		if err != nil {
			fmt.Printf("ERROR: Failed to generate token: %v\n", err)
			return nil, fmt.Errorf("failed to generate token: %w", err)
		}
		fmt.Printf("Token generated successfully (first 20 chars): %s...\n", token[:min(20, len(token))])

		// Clear sensitive data
		user.PasswordHash = ""

		fmt.Println("=== LOGIN USE CASE DEBUG END (SUCCESS - SINGLE TENANT) ===")
		return &domain.AuthResponse{
			User:  user,
			Token: token,
		}, nil
	}

	// Case 2: Multiple tenants - return list for user to select
	fmt.Printf("Multiple tenants found (%d) - returning list for selection\n", len(users))

	tenantList := make([]domain.TenantInfo, 0, len(users))
	for _, user := range users {
		fmt.Printf("Processing user in tenant %d\n", user.TenantID)

		// For now, use the role from the user record
		// In production, you might want to fetch all assigned roles from user_roles table
		tenantInfo := domain.TenantInfo{
			TenantID:            user.TenantID,
			TenantName:          fmt.Sprintf("Tenant %d", user.TenantID), // TODO: Fetch actual tenant name
			UserID:              user.ID,
			Roles:               []string{user.Role},
			DefaultRestaurantID: user.RestaurantID,
		}
		tenantList = append(tenantList, tenantInfo)
	}

	tenantIDs := make([]int, 0, len(users))
	for _, user := range users {
		tenantIDs = append(tenantIDs, user.TenantID)
	}
	selectionToken, err := u.tokenService.GenerateSelectionToken(firstUser.Email, tenantIDs, selectionTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to create selection token: %w", err)
	}

	fmt.Println("=== LOGIN USE CASE DEBUG END (SUCCESS - MULTIPLE TENANTS) ===")
	return &domain.MultiTenantLoginResponse{
		Success:         true,
		MultipleTenants: true,
		Tenants:         tenantList,
		Message:         "Multiple organizations found. Please select one.",
		SelectionToken:  selectionToken,
	}, nil
}

// LoginConfirm exchanges a multi-tenant login proof (issued by Login after the
// password was verified) plus the chosen tenant for a session token. Without a
// valid, unexpired, unused proof covering the requested tenant and email, no
// session is issued.
func (u *AuthUseCase) LoginConfirm(req *domain.TenantSelectionRequest) (*domain.AuthResponse, error) {
	if req.Email == "" || req.TenantID == 0 {
		return nil, errors.New("email and tenant_id are required")
	}
	if req.SelectionToken == "" {
		return nil, errors.New("login verification required: please sign in again")
	}

	proof, err := u.tokenService.ValidateSelectionToken(req.SelectionToken)
	if err != nil {
		return nil, errors.New("login verification expired or invalid: please sign in again")
	}
	if !strings.EqualFold(strings.TrimSpace(req.Email), proof.Email) || !proof.AllowsTenant(req.TenantID) {
		return nil, errors.New("invalid email or tenant selection")
	}
	if u.confirmStore == nil {
		return nil, errors.New("login is unavailable")
	}

	users, err := u.confirmStore.FindUserByEmailAllTenants(proof.Email)
	if err != nil {
		return nil, errors.New("invalid email or tenant selection")
	}

	var selectedUser *domain.User
	for i := range users {
		if users[i].TenantID == req.TenantID {
			selectedUser = &users[i]
			break
		}
	}
	if selectedUser == nil {
		return nil, errors.New("user not found in selected organization")
	}

	// Consume the proof only once the selection is known to be valid.
	if !u.consumeSelection(proof.ID, proof.ExpiresAt.Time) {
		return nil, errors.New("login verification already used: please sign in again")
	}

	_ = u.confirmStore.UpdateLastLogin(selectedUser.ID)

	token, err := u.tokenService.GenerateToken(selectedUser.ID, selectedUser.TenantID, selectedUser.RestaurantID, selectedUser.Email, selectedUser.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	selectedUser.PasswordHash = ""
	return &domain.AuthResponse{
		User:  selectedUser,
		Token: token,
	}, nil
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CheckSubdomainAvailability checks if a subdomain slug is available and valid
// Returns (available bool, message string)
func (u *AuthUseCase) CheckSubdomainAvailability(slug string) (bool, string) {
	validator := NewSubdomainValidator()

	formattedSlug, err := validator.ValidateAndFormat(slug)
	if err != nil {
		return false, err.Error()
	}

	exists, err := u.repo.CheckTenantSlugExists(formattedSlug)
	if err != nil {
		return false, "unable to check subdomain availability"
	}
	if exists {
		return false, "subdomain is already taken"
	}

	return true, "Subdomain is available"
}

// ForgotPassword initiates password reset process
func (u *AuthUseCase) ForgotPassword(req *domain.ForgotPasswordRequest) error {
	// Validate input
	if req.Email == "" {
		return errors.New("email is required")
	}

	// Check if user exists
	_, err := u.repo.FindUserByEmail(req.Email)
	if err != nil {
		// Don't reveal if user exists or not
		return nil
	}

	// TODO: Generate reset token and send email
	// For now, just return success
	return nil
}

// ExchangeSuperAdminToken exchanges a super admin one-time token for a regular JWT
// This validates the token, retrieves the tenant info, finds an admin user, and generates a JWT
func (u *AuthUseCase) ExchangeSuperAdminToken(token *string) (*domain.AuthResponse, error) {
	if token == nil || *token == "" {
		return nil, errors.New("invalid token")
	}

	fmt.Printf("\n========================================\n")
	fmt.Printf("🔄 TOKEN EXCHANGE DEBUG START\n")
	fmt.Printf("========================================\n")
	fmt.Printf("Token received (first 32 chars): %s\n", (*token)[:min(32, len(*token))])

	// Check if super admin repo is configured
	if u.superAdminRepo == nil {
		fmt.Println("❌ ERROR: Super admin repository not configured")
		return nil, errors.New("super admin access not configured")
	}
	fmt.Println("✅ Super admin repository configured")

	ctx := context.Background()

	// 1. Validate the token and get the tenant ID
	fmt.Println("\n→ Step 1: Validating dashboard login token...")
	tokenData, err := u.superAdminRepo.ValidateDashboardLoginToken(ctx, *token)
	if err != nil {
		fmt.Printf("❌ ERROR: Token validation failed: %v\n", err)
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	fmt.Printf("✅ Token validated successfully\n")
	fmt.Printf("   TenantID: %d\n", tokenData.TenantID)
	fmt.Printf("   SuperAdminID: %d\n", tokenData.SuperAdminID)
	fmt.Printf("   ExpiresAt: %s\n", tokenData.ExpiresAt)

	// 2. Get the first admin user from that tenant
	fmt.Printf("\n→ Step 2: Looking up admin/owner user for tenant %d...\n", tokenData.TenantID)
	user, err := u.superAdminRepo.GetFirstAdminUserByTenantID(ctx, tokenData.TenantID)
	if err != nil {
		fmt.Printf("❌ ERROR: Failed to get admin user: %v\n", err)
		return nil, fmt.Errorf("no admin user found in tenant %d - please create a user first", tokenData.TenantID)
	}

	fmt.Printf("✅ Found admin user:\n")
	fmt.Printf("   ID: %d\n", user.ID)
	fmt.Printf("   Email: %s\n", user.Email)
	fmt.Printf("   Name: %s\n", user.Name)
	fmt.Printf("   Role: %s\n", user.Role)
	fmt.Printf("   TenantID: %d\n", user.TenantID)
	if user.RestaurantID != nil {
		fmt.Printf("   RestaurantID: %d\n", *user.RestaurantID)
	}

	// 2.5 CRITICAL: Verify the returned user belongs to the correct tenant
	// This is a safety check to prevent cross-tenant login issues
	if user.TenantID != int(tokenData.TenantID) {
		fmt.Printf("❌ CRITICAL ERROR: Tenant ID mismatch!\n")
		fmt.Printf("   Expected TenantID: %d\n", tokenData.TenantID)
		fmt.Printf("   Actual TenantID: %d\n", user.TenantID)
		return nil, fmt.Errorf("security error: user tenant mismatch - expected %d, got %d", tokenData.TenantID, user.TenantID)
	}
	fmt.Printf("✅ Tenant ID verification passed\n")

	// 3. Generate JWT token for the admin user
	fmt.Println("\n→ Step 3: Generating JWT token...")
	jwtToken, err := u.tokenService.GenerateToken(user.ID, user.TenantID, user.RestaurantID, user.Email, user.Role)
	if err != nil {
		fmt.Printf("❌ ERROR: Failed to generate JWT: %v\n", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	fmt.Printf("✅ JWT generated successfully\n")
	fmt.Printf("   Token length: %d\n", len(jwtToken))
	fmt.Printf("   User: %s\n", user.Email)

	fmt.Printf("\n========================================\n")
	fmt.Printf("✅ TOKEN EXCHANGE COMPLETED SUCCESSFULLY\n")
	fmt.Printf("========================================\n\n")

	return &domain.AuthResponse{
		User:  user,
		Token: jwtToken,
	}, nil
}
