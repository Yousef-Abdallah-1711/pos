package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

// AdminDriverHandler handles admin driver management endpoints.
type AdminDriverHandler struct {
	driverRepo repository.DriverRepository
	orderUC    *usecase.OrderUseCase
}

// NewAdminDriverHandler creates a new admin driver handler.
func NewAdminDriverHandler(driverRepo repository.DriverRepository, orderUC *usecase.OrderUseCase) *AdminDriverHandler {
	return &AdminDriverHandler{
		driverRepo: driverRepo,
		orderUC:    orderUC,
	}
}

type driverPayload struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	LicenseNumber string `json:"license_number"`
	VehicleType   string `json:"vehicle_type"`
	VehicleNumber string `json:"vehicle_number"`
	DateOfBirth   string `json:"date_of_birth,omitempty"`
	Address       string `json:"address,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	ZipCode       string `json:"zip_code,omitempty"`
	Notes         string `json:"notes,omitempty"`
	Status        string `json:"status,omitempty"`
	Availability  string `json:"availability_status,omitempty"`
}

type driverResponse struct {
	ID                  int       `json:"id"`
	TenantID            int       `json:"tenant_id"`
	RestaurantID        int       `json:"restaurant_id"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Email               string    `json:"email"`
	PhoneNumber         string    `json:"phone_number"`
	LicenseNumber       string    `json:"license_number"`
	VehicleType         string    `json:"vehicle_type,omitempty"`
	VehicleNumber       string    `json:"vehicle_number,omitempty"`
	Status              string    `json:"status"`
	AvailabilityStatus  string    `json:"availability_status"`
	TotalDeliveries     int       `json:"total_deliveries"`
	CompletedDeliveries int       `json:"completed_deliveries"`
	CancelledDeliveries int       `json:"cancelled_deliveries"`
	AverageRating       float64   `json:"average_rating"`
	ActiveOrdersCount   int       `json:"active_orders_count"`
	IsVerified          bool      `json:"is_verified"`
	DateOfBirth         string    `json:"date_of_birth,omitempty"`
	Address             string    `json:"address,omitempty"`
	City                string    `json:"city,omitempty"`
	State               string    `json:"state,omitempty"`
	ZipCode             string    `json:"zip_code,omitempty"`
	Notes               string    `json:"notes,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func formatDateOnly(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("2006-01-02")
}

func parseDateOnly(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func toDriverResponse(driver domain.Driver) driverResponse {
	return driverResponse{
		ID:                  driver.ID,
		TenantID:            driver.TenantID,
		RestaurantID:        driver.RestaurantID,
		FirstName:           driver.FirstName,
		LastName:            driver.LastName,
		Email:               driver.Email,
		PhoneNumber:         driver.PhoneNumber,
		LicenseNumber:       driver.LicenseNumber,
		VehicleType:         driver.VehicleType,
		VehicleNumber:       driver.VehicleNumber,
		Status:              driver.Status,
		AvailabilityStatus:  driver.AvailabilityStatus,
		TotalDeliveries:     driver.TotalDeliveries,
		CompletedDeliveries: driver.CompletedDeliveries,
		CancelledDeliveries: driver.CancelledDeliveries,
		AverageRating:       driver.AverageRating,
		ActiveOrdersCount:   driver.ActiveOrdersCount,
		IsVerified:          driver.IsVerified,
		DateOfBirth:         formatDateOnly(driver.DateOfBirth),
		Address:             driver.Address,
		City:                driver.City,
		State:               driver.State,
		ZipCode:             driver.ZipCode,
		Notes:               driver.Notes,
		CreatedAt:           driver.CreatedAt,
		UpdatedAt:           driver.UpdatedAt,
	}
}

// CreateDriver creates a new driver.
// POST /api/v1/admin/drivers
func (h *AdminDriverHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	var req driverPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(req.Email)
	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
	req.LicenseNumber = strings.TrimSpace(req.LicenseNumber)

	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.PhoneNumber == "" || req.LicenseNumber == "" {
		respondError(w, http.StatusBadRequest, "Missing required driver fields")
		return
	}

	dateOfBirth, err := parseDateOnly(req.DateOfBirth)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid date_of_birth format")
		return
	}

	driver := domain.Driver{
		TenantID:           int(tenantID),
		RestaurantID:       int(restaurantID),
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		Email:              req.Email,
		PhoneNumber:        req.PhoneNumber,
		LicenseNumber:      req.LicenseNumber,
		VehicleType:        strings.TrimSpace(req.VehicleType),
		VehicleNumber:      strings.TrimSpace(req.VehicleNumber),
		Status:             "active",
		AvailabilityStatus: "available",
		DateOfBirth:        dateOfBirth,
		Address:            strings.TrimSpace(req.Address),
		City:               strings.TrimSpace(req.City),
		State:              strings.TrimSpace(req.State),
		ZipCode:            strings.TrimSpace(req.ZipCode),
		Notes:              strings.TrimSpace(req.Notes),
	}

	if err := h.driverRepo.Create(r.Context(), &driver); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create driver")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"driver": toDriverResponse(driver),
	})
}

// ListDrivers lists all drivers for the current tenant/restaurant.
// GET /api/v1/admin/drivers
func (h *AdminDriverHandler) ListDrivers(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	drivers, err := h.driverRepo.GetByTenantAndRestaurant(r.Context(), int(tenantID), int(restaurantID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list drivers")
		return
	}

	response := make([]driverResponse, 0, len(drivers))
	for _, driver := range drivers {
		response = append(response, toDriverResponse(driver))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"drivers": response,
	})
}

// GetAvailableDrivers lists available drivers.
// GET /api/v1/admin/drivers/available
func (h *AdminDriverHandler) GetAvailableDrivers(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	drivers, err := h.driverRepo.GetAvailableDrivers(r.Context(), int(tenantID), int(restaurantID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load drivers")
		return
	}

	response := make([]driverResponse, 0, len(drivers))
	for _, driver := range drivers {
		response = append(response, toDriverResponse(driver))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"drivers": response,
	})
}

// GetDriver retrieves a single driver.
// GET /api/v1/admin/drivers/{id}
func (h *AdminDriverHandler) GetDriver(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	driverID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid driver ID")
		return
	}

	driver, err := h.driverRepo.GetByID(r.Context(), driverID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Driver not found")
		return
	}

	if driver.TenantID != int(tenantID) || driver.RestaurantID != int(restaurantID) {
		respondError(w, http.StatusForbidden, "Driver not found")
		return
	}

	respondJSON(w, http.StatusOK, toDriverResponse(*driver))
}

// UpdateDriver updates an existing driver.
// PUT /api/v1/admin/drivers/{id}
func (h *AdminDriverHandler) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	driverID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid driver ID")
		return
	}

	existing, err := h.driverRepo.GetByID(r.Context(), driverID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Driver not found")
		return
	}
	if existing.TenantID != int(tenantID) || existing.RestaurantID != int(restaurantID) {
		respondError(w, http.StatusForbidden, "Driver not found")
		return
	}

	var req driverPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	dateOfBirth, err := parseDateOnly(req.DateOfBirth)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid date_of_birth format")
		return
	}

	existing.FirstName = strings.TrimSpace(req.FirstName)
	existing.LastName = strings.TrimSpace(req.LastName)
	existing.Email = strings.TrimSpace(req.Email)
	existing.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
	existing.LicenseNumber = strings.TrimSpace(req.LicenseNumber)
	existing.VehicleType = strings.TrimSpace(req.VehicleType)
	existing.VehicleNumber = strings.TrimSpace(req.VehicleNumber)
	existing.DateOfBirth = dateOfBirth
	existing.Address = strings.TrimSpace(req.Address)
	existing.City = strings.TrimSpace(req.City)
	existing.State = strings.TrimSpace(req.State)
	existing.ZipCode = strings.TrimSpace(req.ZipCode)
	existing.Notes = strings.TrimSpace(req.Notes)

	if err := h.driverRepo.Update(r.Context(), existing); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update driver")
		return
	}

	updated, _ := h.driverRepo.GetByID(r.Context(), driverID)
	if updated == nil {
		updated = existing
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"driver": toDriverResponse(*updated),
	})
}

// DeleteDriver deletes a driver.
// DELETE /api/v1/admin/drivers/{id}
func (h *AdminDriverHandler) DeleteDriver(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	driverID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid driver ID")
		return
	}

	driver, err := h.driverRepo.GetByID(r.Context(), driverID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Driver not found")
		return
	}
	if driver.TenantID != int(tenantID) || driver.RestaurantID != int(restaurantID) {
		respondError(w, http.StatusForbidden, "Driver not found")
		return
	}

	if err := h.driverRepo.Delete(r.Context(), driverID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete driver")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// UpdateDriverStatus updates a driver's status.
// PUT /api/v1/admin/drivers/{id}/status
func (h *AdminDriverHandler) UpdateDriverStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	driverID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid driver ID")
		return
	}

	driver, err := h.driverRepo.GetByID(r.Context(), driverID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Driver not found")
		return
	}
	if driver.TenantID != int(tenantID) || driver.RestaurantID != int(restaurantID) {
		respondError(w, http.StatusForbidden, "Driver not found")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		respondError(w, http.StatusBadRequest, "Status is required")
		return
	}

	if err := h.driverRepo.UpdateStatus(r.Context(), driverID, status); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update driver status")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// AssignOrderToDriver assigns an order to a driver.
// POST /api/v1/admin/orders/{orderId}/assign
func (h *AdminDriverHandler) AssignOrderToDriver(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	orderID, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req domain.AssignOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.DriverID == 0 {
		respondError(w, http.StatusBadRequest, "driver_id is required")
		return
	}

	if _, err := h.orderUC.GetOrder(tenantID, restaurantID, int64(orderID)); err != nil {
		respondError(w, http.StatusNotFound, "Order not found")
		return
	}

	driver, err := h.driverRepo.GetByID(r.Context(), req.DriverID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Driver not found")
		return
	}
	if driver.TenantID != int(tenantID) || driver.RestaurantID != int(restaurantID) {
		respondError(w, http.StatusForbidden, "Driver not found")
		return
	}

	existingAssignment, err := h.driverRepo.GetAssignmentByOrderID(r.Context(), orderID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to check existing assignments")
		return
	}
	if existingAssignment != nil {
		respondError(w, http.StatusConflict, "Order already assigned to a driver")
		return
	}

	assignment := domain.DriverAssignment{
		OrderID:         orderID,
		DriverID:        req.DriverID,
		Status:          "pending",
		AssignmentNotes: strings.TrimSpace(req.AssignmentNotes),
	}

	if err := h.driverRepo.AssignOrder(r.Context(), &assignment); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to assign order")
		return
	}

	_ = h.driverRepo.UpdateAvailability(r.Context(), req.DriverID, "busy")
	_ = h.driverRepo.UpdateActiveOrders(r.Context(), req.DriverID, driver.ActiveOrdersCount+1)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"assignment": assignment,
	})
}
