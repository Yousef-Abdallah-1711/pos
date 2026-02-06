package domain

import "time"

// Driver represents a delivery driver
type Driver struct {
	ID                  int        `json:"id"`
	TenantID            int        `json:"tenant_id"`
	RestaurantID        int        `json:"restaurant_id"`
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	Email               string     `json:"email"`
	PhoneNumber         string     `json:"phone_number"`
	LicenseNumber       string     `json:"license_number"`
	VehicleType         string     `json:"vehicle_type,omitempty"`
	VehicleNumber       string     `json:"vehicle_number,omitempty"`
	Status              string     `json:"status"`              // active, inactive, on_duty, off_duty
	AvailabilityStatus  string     `json:"availability_status"` // available, busy, offline
	CurrentLatitude     *float64   `json:"current_latitude,omitempty"`
	CurrentLongitude    *float64   `json:"current_longitude,omitempty"`
	TotalDeliveries     int        `json:"total_deliveries"`
	CompletedDeliveries int        `json:"completed_deliveries"`
	CancelledDeliveries int        `json:"cancelled_deliveries"`
	AverageRating       float64    `json:"average_rating"`
	ActiveOrdersCount   int        `json:"active_orders_count"`
	DateOfBirth         *time.Time `json:"date_of_birth,omitempty"`
	Address             string     `json:"address,omitempty"`
	City                string     `json:"city,omitempty"`
	State               string     `json:"state,omitempty"`
	ZipCode             string     `json:"zip_code,omitempty"`
	IsVerified          bool       `json:"is_verified"`
	VerificationDate    *time.Time `json:"verification_date,omitempty"`
	JoinedDate          time.Time  `json:"joined_date"`
	LastActive          *time.Time `json:"last_active,omitempty"`
	Notes               string     `json:"notes,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// DriverAssignment represents assignment of an order to a driver
type DriverAssignment struct {
	ID              int        `json:"id"`
	OrderID         int        `json:"order_id"`
	DriverID        int        `json:"driver_id"`
	Driver          *Driver    `json:"driver,omitempty"`
	AssignedAt      time.Time  `json:"assigned_at"`
	AcceptedAt      *time.Time `json:"accepted_at,omitempty"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	Status          string     `json:"status"` // pending, accepted, in_progress, completed, failed
	AssignmentNotes string     `json:"assignment_notes,omitempty"`
	Rating          *int       `json:"rating,omitempty"`
	RatingComment   string     `json:"rating_comment,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// DriverLocation represents driver location history
type DriverLocation struct {
	ID         int       `json:"id"`
	DriverID   int       `json:"driver_id"`
	OrderID    *int      `json:"order_id,omitempty"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Accuracy   *int      `json:"accuracy,omitempty"`
	RecordedAt time.Time `json:"recorded_at"`
}

// DriverStats represents driver performance statistics
type DriverStats struct {
	DriverID            int     `json:"driver_id"`
	FirstName           string  `json:"first_name"`
	LastName            string  `json:"last_name"`
	TotalDeliveries     int     `json:"total_deliveries"`
	CompletedDeliveries int     `json:"completed_deliveries"`
	CancelledDeliveries int     `json:"cancelled_deliveries"`
	AverageRating       float64 `json:"average_rating"`
	CompletionRate      float64 `json:"completion_rate"`
	ActiveOrders        int     `json:"active_orders"`
}

// CreateDriverRequest represents the request to create a new driver
type CreateDriverRequest struct {
	FirstName     string     `json:"first_name" validate:"required"`
	LastName      string     `json:"last_name" validate:"required"`
	Email         string     `json:"email" validate:"required,email"`
	PhoneNumber   string     `json:"phone_number" validate:"required"`
	LicenseNumber string     `json:"license_number" validate:"required"`
	VehicleType   string     `json:"vehicle_type"`
	VehicleNumber string     `json:"vehicle_number"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`
	Address       string     `json:"address,omitempty"`
	City          string     `json:"city,omitempty"`
	State         string     `json:"state,omitempty"`
	ZipCode       string     `json:"zip_code,omitempty"`
}

// UpdateDriverRequest represents the request to update a driver
type UpdateDriverRequest struct {
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Email         string `json:"email,omitempty"`
	PhoneNumber   string `json:"phone_number,omitempty"`
	VehicleType   string `json:"vehicle_type,omitempty"`
	VehicleNumber string `json:"vehicle_number,omitempty"`
	Status        string `json:"status,omitempty"`
	Address       string `json:"address,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	ZipCode       string `json:"zip_code,omitempty"`
	Notes         string `json:"notes,omitempty"`
}

// AssignOrderRequest represents the request to assign an order to a driver
type AssignOrderRequest struct {
	DriverID        int    `json:"driver_id" validate:"required"`
	AssignmentNotes string `json:"assignment_notes,omitempty"`
}

// UpdateAssignmentStatusRequest represents the request to update assignment status
type UpdateAssignmentStatusRequest struct {
	Status string `json:"status" validate:"required"`
	Notes  string `json:"notes,omitempty"`
}

// RateDeliveryRequest represents the request to rate a delivery
type RateDeliveryRequest struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment,omitempty"`
}

// UpdateLocationRequest represents the request to update driver location
type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Accuracy  *int    `json:"accuracy,omitempty"`
}
