"""One-off codemod (Phase 2A C005): public order object-level authorization."""


def edit(path, pairs):
    s = open(path, encoding="utf-8").read()
    for old, new in pairs:
        assert s.count(old) >= 1, (path, old[:60])
        s = s.replace(old, new, 1)
    open(path, "w", encoding="utf-8", newline="").write(s)


# 1. domain: non-persisted token carried on the create response
edit("backend/internal/domain/order.go", [(
    '	CustomerPhone string `json:"customer_phone"`',
    '	CustomerPhone string `json:"customer_phone"`\n'
    '	// TrackingToken is the raw public access token. It is only populated on the response of\n'
    '	// order creation (the database stores its hash) and is never loaded from the database.\n'
    '	TrackingToken string `json:"tracking_token,omitempty"`',
)])

# 2. repository
repo = open("backend/internal/repository/order_repo.go", encoding="utf-8").read()
repo += '''

// SetPublicTokenHash stores the hash of the public tracking token for an order.
func (r *OrderRepository) SetPublicTokenHash(tenantID, restaurantID, orderID int64, hash string) error {
	res, err := r.db.Exec(
		`UPDATE orders SET public_token_hash = $1 WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4`,
		hash, orderID, tenantID, restaurantID,
	)
	if err != nil {
		return fmt.Errorf("failed to store public token: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("order not found")
	}
	return nil
}

// GetPublicAccessInfo returns the stored public token hash (empty for legacy orders) and the
// customer phone of an order, scoped to tenant and restaurant.
func (r *OrderRepository) GetPublicAccessInfo(tenantID, restaurantID, orderID int64) (hash string, phone string, err error) {
	var h sql.NullString
	err = r.db.QueryRow(
		`SELECT public_token_hash, customer_phone FROM orders WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`,
		orderID, tenantID, restaurantID,
	).Scan(&h, &phone)
	if err == sql.ErrNoRows {
		return "", "", errors.New("order not found")
	}
	if err != nil {
		return "", "", fmt.Errorf("failed to read order access info: %w", err)
	}
	return h.String, phone, nil
}
'''
open("backend/internal/repository/order_repo.go", "w", encoding="utf-8", newline="").write(repo)

# 3. usecase: issue token on create, authorize public access
edit("backend/internal/usecase/order_usecase.go", [
    ('''	// Table integration: mark table as occupied for dine-in orders.
	if finalOrder.TableID != nil && uc.tableRepo != nil {''',
     '''	// Public access credential: the raw token is returned once; only its hash is stored.
	trackingToken, tokenHash, tokenErr := domain.GeneratePublicTrackingToken()
	if tokenErr == nil {
		tokenErr = uc.orderRepo.SetPublicTokenHash(tenantID, restaurantID, finalOrder.ID, tokenHash)
	}
	if tokenErr != nil {
		// The order stays valid; it is simply not publicly trackable.
		log.Printf("[ORDERS] failed to issue public tracking token: tenant=%d restaurant=%d order=%d err=%v", tenantID, restaurantID, finalOrder.ID, tokenErr)
		trackingToken = ""
	}

	// Table integration: mark table as occupied for dine-in orders.
	if finalOrder.TableID != nil && uc.tableRepo != nil {'''),
    ('''	return finalOrder, nil
}

// GetOrder retrieves a single order by ID''',
     '''	finalOrder.TrackingToken = trackingToken
	return finalOrder, nil
}

// AuthorizePublicOrderAccess decides whether an unauthenticated caller may read or cancel an order.
// Access requires the order's tracking token or, as a fallback for the manual order-number lookup,
// the customer phone used at checkout. Legacy orders (no stored token) are only reachable by phone.
// Any failure, including a missing order, returns false so callers can answer an identical 404.
func (uc *OrderUseCase) AuthorizePublicOrderAccess(tenantID, restaurantID, orderID int64, token, phone string) bool {
	storedHash, orderPhone, err := uc.orderRepo.GetPublicAccessInfo(tenantID, restaurantID, orderID)
	if err != nil {
		return false
	}
	if domain.VerifyPublicTrackingToken(token, storedHash) {
		return true
	}
	return domain.PhonesMatch(phone, orderPhone)
}

// GetOrder retrieves a single order by ID'''),
])

# 4. handlers
h = "backend/internal/handler/http/public_order_handler.go"
s = open(h, encoding="utf-8").read()

helper = '''
// publicOrderProof extracts the caller's proof of ownership. The tracking token travels in a
// header (never the URL) so it does not end up in access logs or browser history.
func publicOrderProof(r *http.Request) (token, phone string) {
	return strings.TrimSpace(r.Header.Get("X-Order-Token")), strings.TrimSpace(r.Header.Get("X-Customer-Phone"))
}

// requirePublicOrderAccess answers an identical 404 for missing orders and for orders the caller
// cannot prove ownership of, so order ids and numbers cannot be enumerated.
func (h *PublicOrderHandler) requirePublicOrderAccess(w http.ResponseWriter, r *http.Request, tenantID, restaurantID, orderID int64, allowPhone bool) bool {
	token, phone := publicOrderProof(r)
	if !allowPhone {
		phone = ""
	}
	if !h.orderUC.AuthorizePublicOrderAccess(tenantID, restaurantID, orderID, token, phone) {
		respondError(w, http.StatusNotFound, "Order not found")
		return false
	}
	return true
}
'''
anchor = "// GetOrderByNumber retrieves order details by order number"
assert anchor in s
s = s.replace(anchor, helper.lstrip("\n") + "\n" + anchor, 1)

# by id handlers: guard right before the usecase fetch
def guard_before(s, marker, count_expected=1):
    return s

# GetOrder
old = '''	// Get order from usecase
	order, err := h.orderUC.GetOrder(tenantID, restaurantID, orderID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}

	// Return order with items'''
new = '''	if !h.requirePublicOrderAccess(w, r, tenantID, restaurantID, orderID, false) {
		return
	}

	// Get order from usecase
	order, err := h.orderUC.GetOrder(tenantID, restaurantID, orderID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}

	// Return order with items'''
assert old in s
s = s.replace(old, new, 1)

# GetOrderByNumber: authorize after lookup by number (token or phone)
old = '''	// Get order from usecase
	order, err := h.orderUC.GetOrderByNumber(tenantID, restaurantID, orderNumber)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}
'''
new = '''	// Get order from usecase
	order, err := h.orderUC.GetOrderByNumber(tenantID, restaurantID, orderNumber)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}
	if !h.requirePublicOrderAccess(w, r, tenantID, restaurantID, order.ID, true) {
		return
	}
'''
assert old in s
s = s.replace(old, new, 1)

# GetOrderStatus
old = '''	// Get order from usecase
	order, err := h.orderUC.GetOrder(tenantID, restaurantID, orderID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order status")
		return
	}
'''
new = '''	if !h.requirePublicOrderAccess(w, r, tenantID, restaurantID, orderID, false) {
		return
	}

	// Get order from usecase
	order, err := h.orderUC.GetOrder(tenantID, restaurantID, orderID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order status")
		return
	}
'''
assert old in s
s = s.replace(old, new, 1)

# TrackOrder
old = '''	// Get order from usecase
	order, err := h.orderUC.GetOrder(tenantID, restaurantID, orderID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}

	// Get status history'''
new = '''	if !h.requirePublicOrderAccess(w, r, tenantID, restaurantID, orderID, false) {
		return
	}

	// Get order from usecase
	order, err := h.orderUC.GetOrder(tenantID, restaurantID, orderID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "Order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}

	// Get status history'''
assert old in s
s = s.replace(old, new, 1)

# CancelOrder
old = '''	// Parse optional cancellation reason'''
new = '''	if !h.requirePublicOrderAccess(w, r, tenantID, restaurantID, orderID, false) {
		return
	}

	// Parse optional cancellation reason'''
assert old in s
s = s.replace(old, new, 1)
open(h, "w", encoding="utf-8", newline="").write(s)
print("ok")
