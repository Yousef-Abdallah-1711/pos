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
	return u.RegisterContext(context.Background(), req)
}

// RegisterContext validates user input before persistence and carries request
// cancellation through the tenant onboarding transaction.
func (u *AuthUseCase) RegisterContext(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	if req == nil {
		return nil, &domain.ValidationErrors{Fields: map[string]string{"form": "Registration details are required."}}
	}
	if err := req.NormalizeAndValidate(); err != nil {
		return nil, err
	}

	validator := NewSubdomainValidator()
	slugProvided := strings.TrimSpace(req.RestaurantSlug) != ""
	if slugProvided {
		formattedSlug, err := validator.ValidateAndFormat(req.RestaurantSlug)
		if err != nil {
			return nil, &domain.ValidationErrors{Fields: map[string]string{"restaurant_slug": "Use a valid, unreserved website address."}}
		}
		req.RestaurantSlug = formattedSlug

		exists, err := u.repo.CheckTenantSlugExists(req.RestaurantSlug)
		if err != nil {
			return nil, fmt.Errorf("check tenant slug: %w", err)
		}
		if exists {
			return nil, domain.ErrSubdomainAlreadyExists
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
		return nil, domain.ErrEmailAlreadyExists
	}

	// Create tenant with user
	user, err := u.repo.CreateTenantWithUserContext(ctx, req)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return nil, domain.ErrEmailAlreadyExists
		}
		if errors.Is(err, repository.ErrTenantSlugExists) {
			return nil, domain.ErrSubdomainAlreadyExists
		}
		if errors.Is(err, repository.ErrTenantNameExists) {
			return nil, domain.ErrRestaurantNameConflict
		}
		if errors.Is(err, repository.ErrTenantDomainExists) {
			return nil, domain.ErrSubdomainAlreadyExists
		}
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	// Generate JWT token
	token, err := u.tokenService.GenerateToken(user.ID, user.TenantID, user.RestaurantID, user.Email, user.Role)
	if err != nil {
		return &domain.AuthResponse{User: user}, fmt.Errorf("registration committed but session token creation failed: %w", err)
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
	if req == nil {
		return nil, &domain.ValidationErrors{Fields: map[string]string{"email": "Email is required.", "password": "Password is required."}}
	}
	if err := req.NormalizeAndValidate(); err != nil {
		return nil, err
	}

	// Find ALL users with this email across all tenants
	users, err := u.repo.FindUserByEmailAllTenants(req.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrAuthDependencyUnavailable, err)
	}

	if len(users) == 0 {
		return nil, errors.New("invalid email or password")
	}

	// Verify password against the first user (all have same email, should have same password)
	firstUser := &users[0]

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(firstUser.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Case 1: Single tenant - auto-login
	if len(users) == 1 {
		user := &users[0]

		// Update last login
		_ = u.repo.UpdateLastLogin(user.ID)

		// Generate JWT token
		token, err := u.tokenService.GenerateToken(user.ID, user.TenantID, user.RestaurantID, user.Email, user.Role)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrTokenGeneration, err)
		}

		// Clear sensitive data
		user.PasswordHash = ""

		return &domain.AuthResponse{
			User:  user,
			Token: token,
		}, nil
	}

	// Case 2: Multiple tenants - return list for user to select
	tenantList := make([]domain.TenantInfo, 0, len(users))
	for _, user := range users {

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
		return nil, fmt.Errorf("%w: %v", domain.ErrTokenGeneration, err)
	}

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
		return nil, domain.ErrAuthDependencyUnavailable
	}

	users, err := u.confirmStore.FindUserByEmailAllTenants(proof.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrAuthDependencyUnavailable, err)
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
		return nil, fmt.Errorf("%w: %v", domain.ErrTokenGeneration, err)
	}

	selectedUser.PasswordHash = ""
	return &domain.AuthResponse{
		User:  selectedUser,
		Token: token,
	}, nil
}

// RefreshToken issues a replacement access token for an already authenticated
// browser session. The caller must validate the current session before calling
// this method; session ID rotation is handled atomically by the HTTP handler.
func (u *AuthUseCase) RefreshToken(claims *jwt.Claims) (string, error) {
	if u == nil || u.tokenService == nil || claims == nil || claims.UserID <= 0 || claims.TenantID <= 0 {
		return "", errors.New("authenticated session is required")
	}
	return u.tokenService.GenerateToken(claims.UserID, claims.TenantID, &claims.RestaurantID, claims.Email, claims.Role)
}

// ValidateCredential verifies a stored access credential before logout audit
// attribution. It does not authenticate requests by itself.
func (u *AuthUseCase) ValidateCredential(credential string) (*jwt.Claims, error) {
	if u == nil || u.tokenService == nil || credential == "" {
		return nil, errors.New("credential validation is unavailable")
	}
	return u.tokenService.ValidateToken(credential)
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
func (u *AuthUseCase) ExchangeSuperAdminToken(ctx context.Context, token *string) (*domain.AuthResponse, error) {
	if token == nil || *token == "" {
		return nil, errors.New("invalid token")
	}

	// Check if super admin repo is configured
	if u.superAdminRepo == nil {
		return nil, errors.New("super admin access not configured")
	}

	// 1. Validate the token and get the tenant ID
	tokenData, err := u.superAdminRepo.ValidateDashboardLoginToken(ctx, *token)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	// 2. Get the first admin user from that tenant
	user, err := u.superAdminRepo.GetFirstAdminUserByTenantID(ctx, tokenData.TenantID)
	if err != nil {
		return nil, fmt.Errorf("no admin user found in tenant %d - please create a user first", tokenData.TenantID)
	}

	// 2.5 CRITICAL: Verify the returned user belongs to the correct tenant
	// This is a safety check to prevent cross-tenant login issues
	if user.TenantID != int(tokenData.TenantID) {
		return nil, fmt.Errorf("security error: user tenant mismatch - expected %d, got %d", tokenData.TenantID, user.TenantID)
	}

	// 3. Generate JWT token for the admin user
	jwtToken, err := u.tokenService.GenerateToken(user.ID, user.TenantID, user.RestaurantID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.AuthResponse{
		User:  user,
		Token: jwtToken,
	}, nil
}
