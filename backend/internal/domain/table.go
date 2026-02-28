package domain

import (
	"encoding/json"
	"time"
)

type TableStatus string

const (
	TableStatusAvailable   TableStatus = "available"
	TableStatusOccupied    TableStatus = "occupied"
	TableStatusReserved    TableStatus = "reserved"
	TableStatusMaintenance TableStatus = "maintenance"
)

type TableLocation string

const (
	TableLocationIndoor  TableLocation = "indoor"
	TableLocationOutdoor TableLocation = "outdoor"
	TableLocationVIP     TableLocation = "vip"
	TableLocationRooftop TableLocation = "rooftop"
	TableLocationPrivate TableLocation = "private"
)

type TableShape string

const (
	TableShapeSquare    TableShape = "square"
	TableShapeCircle    TableShape = "circle"
	TableShapeRectangle TableShape = "rectangle"
)

type RestaurantTable struct {
	ID           int64           `json:"id"`
	TenantID     int64           `json:"tenant_id"`
	RestaurantID int64           `json:"restaurant_id"`
	TableNumber  int             `json:"table_number"`
	Name         string          `json:"name,omitempty"`
	Capacity     int             `json:"capacity"`
	Location     TableLocation   `json:"location"`
	Status       TableStatus     `json:"status"`
	Description  string          `json:"description,omitempty"`
	QRCodeConfig json.RawMessage `json:"qr_code_config,omitempty"`
	PositionX    float64         `json:"position_x"`
	PositionY    float64         `json:"position_y"`
	Shape        TableShape      `json:"shape"`
	LayoutID     *int64          `json:"layout_id,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type TableLayout struct {
	ID              int64             `json:"id"`
	TenantID        int64             `json:"tenant_id"`
	RestaurantID    int64             `json:"restaurant_id"`
	Name            string            `json:"name"`
	FloorNumber     int               `json:"floor_number"`
	CanvasWidth     float64           `json:"canvas_width"`
	CanvasHeight    float64           `json:"canvas_height"`
	BackgroundImage string            `json:"background_image,omitempty"`
	IsDefault       bool              `json:"is_default"`
	Tables          []RestaurantTable `json:"tables,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type CreateTableRequest struct {
	TableNumber int    `json:"table_number"`
	Name        string `json:"name"`
	Capacity    int    `json:"capacity"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Shape       string `json:"shape"`
}

type UpdateTableRequest struct {
	TableNumber *int     `json:"table_number,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Capacity    *int     `json:"capacity,omitempty"`
	Location    *string  `json:"location,omitempty"`
	Status      *string  `json:"status,omitempty"`
	Description *string  `json:"description,omitempty"`
	PositionX   *float64 `json:"position_x,omitempty"`
	PositionY   *float64 `json:"position_y,omitempty"`
	Shape       *string  `json:"shape,omitempty"`
	LayoutID    *int64   `json:"layout_id,omitempty"`
}

type UpdateTableStatusRequest struct {
	Status string `json:"status"`
}

type UpdateTablePositionsRequest struct {
	Positions []TablePosition `json:"positions"`
}

type TablePosition struct {
	TableID   int64   `json:"table_id"`
	PositionX float64 `json:"position_x"`
	PositionY float64 `json:"position_y"`
}

type TableStats struct {
	TotalTables   int `json:"total_tables"`
	Available     int `json:"available"`
	Occupied      int `json:"occupied"`
	Reserved      int `json:"reserved"`
	Maintenance   int `json:"maintenance"`
	TotalCapacity int `json:"total_capacity"`
}

type CreateLayoutRequest struct {
	Name            string  `json:"name"`
	FloorNumber     int     `json:"floor_number"`
	CanvasWidth     float64 `json:"canvas_width"`
	CanvasHeight    float64 `json:"canvas_height"`
	BackgroundImage string  `json:"background_image"`
	IsDefault       bool    `json:"is_default"`
}

type TablePublicInfo struct {
	ID             int64  `json:"id"`
	TableNumber    int    `json:"table_number"`
	Name           string `json:"name,omitempty"`
	Capacity       int    `json:"capacity"`
	Location       string `json:"location"`
	Status         string `json:"status"`
	TenantID       int64  `json:"tenant_id"`
	RestaurantID   int64  `json:"restaurant_id"`
	RestaurantSlug string `json:"restaurant_slug"`
}

func ValidTableStatus(status string) bool {
	switch TableStatus(status) {
	case TableStatusAvailable, TableStatusOccupied, TableStatusReserved, TableStatusMaintenance:
		return true
	default:
		return false
	}
}

func ValidTableLocation(location string) bool {
	switch TableLocation(location) {
	case TableLocationIndoor, TableLocationOutdoor, TableLocationVIP, TableLocationRooftop, TableLocationPrivate:
		return true
	default:
		return false
	}
}

func ValidTableShape(shape string) bool {
	switch TableShape(shape) {
	case TableShapeSquare, TableShapeCircle, TableShapeRectangle:
		return true
	default:
		return false
	}
}
