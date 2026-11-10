package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

type JournalEntryHandler struct {
	journalRepo *repository.JournalEntryRepository
}

func NewJournalEntryHandler(journalRepo *repository.JournalEntryRepository) *JournalEntryHandler {
	return &JournalEntryHandler{journalRepo: journalRepo}
}

func (h *JournalEntryHandler) ListJournalEntries(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	filters := domain.JournalEntryFilters{
		Status:    strings.TrimSpace(r.URL.Query().Get("status")),
		EntryType: strings.TrimSpace(r.URL.Query().Get("entry_type")),
	}

	if startDate := strings.TrimSpace(r.URL.Query().Get("start_date")); startDate != "" {
		parsed, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   "Invalid start_date format",
			})
			return
		}
		filters.StartDate = &parsed
	}
	if endDate := strings.TrimSpace(r.URL.Query().Get("end_date")); endDate != "" {
		parsed, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   "Invalid end_date format",
			})
			return
		}
		filters.EndDate = &parsed
	}

	page := GetQueryParamInt(r, "page", 1)
	pageSize := GetQueryParamInt(r, "page_size", 20)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	entries, total, err := h.journalRepo.ListJournalEntriesPaginated(tenantID, restaurantID, filters, pageSize, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch journal entries",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    entries,
		"pagination": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *JournalEntryHandler) GetJournalEntry(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid journal entry ID",
		})
		return
	}

	entry, err := h.journalRepo.GetJournalEntryByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch journal entry",
		})
		return
	}
	if entry == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Journal entry not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    entry,
	})
}

func (h *JournalEntryHandler) CreateJournalEntry(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req domain.CreateJournalEntryRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	entryDate, err := parseOptionalDate(req.EntryDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid entry_date format",
		})
		return
	}

	entry := domain.JournalEntry{
		TenantID:        tenantID,
		RestaurantID:    restaurantID,
		EntryDate:       entryDate,
		Description:     strings.TrimSpace(req.Description),
		DescriptionAr:   strings.TrimSpace(req.DescriptionAr),
		ReferenceType:   strings.TrimSpace(req.ReferenceType),
		ReferenceID:     req.ReferenceID,
		ReferenceNumber: strings.TrimSpace(req.ReferenceNumber),
		EntryType:       strings.TrimSpace(req.EntryType),
		Notes:           strings.TrimSpace(req.Notes),
		Attachments:     req.Attachments,
		CreatedBy:       &userID,
		UpdatedBy:       &userID,
		Lines:           mapJournalLines(req.Lines),
	}

	id, err := h.journalRepo.CreateJournalEntry(&entry)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	created, err := h.journalRepo.GetJournalEntryByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"data":    map[string]int{"id": id},
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    created,
	})
}

func (h *JournalEntryHandler) UpdateJournalEntry(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid journal entry ID",
		})
		return
	}

	existing, err := h.journalRepo.GetJournalEntryByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch journal entry",
		})
		return
	}
	if existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Journal entry not found",
		})
		return
	}

	var req domain.UpdateJournalEntryRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	entryDate, err := parseOptionalDate(req.EntryDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid entry_date format",
		})
		return
	}

	entry := domain.JournalEntry{
		ID:              id,
		TenantID:        tenantID,
		RestaurantID:    restaurantID,
		EntryDate:       entryDate,
		Description:     strings.TrimSpace(req.Description),
		DescriptionAr:   strings.TrimSpace(req.DescriptionAr),
		ReferenceType:   strings.TrimSpace(req.ReferenceType),
		ReferenceID:     req.ReferenceID,
		ReferenceNumber: strings.TrimSpace(req.ReferenceNumber),
		EntryType:       strings.TrimSpace(req.EntryType),
		Notes:           strings.TrimSpace(req.Notes),
		Attachments:     req.Attachments,
		CreatedBy:       existing.CreatedBy,
		UpdatedBy:       &userID,
		Lines:           mapJournalLines(req.Lines),
	}

	if err := h.journalRepo.UpdateDraftEntry(&entry); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "journal entry not found" {
			statusCode = http.StatusNotFound
		}
		writeJSON(w, statusCode, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, err := h.journalRepo.GetJournalEntryByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    updated,
	})
}

func (h *JournalEntryHandler) PostJournalEntry(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid journal entry ID",
		})
		return
	}

	if err := h.journalRepo.PostEntry(tenantID, restaurantID, id, userID); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "journal entry not found" {
			statusCode = http.StatusNotFound
		}
		writeJSON(w, statusCode, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	entry, err := h.journalRepo.GetJournalEntryByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    entry,
	})
}

func (h *JournalEntryHandler) VoidJournalEntry(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid journal entry ID",
		})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if r.Body != nil {
		if err := decodeJSONStrictAllowEmpty(r, &req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   "Invalid request body",
			})
			return
		}
	}
	req.Reason = strings.TrimSpace(req.Reason)
	if req.Reason == "" {
		req.Reason = "Voided by user"
	}

	if err := h.journalRepo.VoidEntry(tenantID, restaurantID, id, userID, req.Reason); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "journal entry not found" {
			statusCode = http.StatusNotFound
		}
		writeJSON(w, statusCode, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	entry, err := h.journalRepo.GetJournalEntryByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    entry,
	})
}

func parseOptionalDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", value)
}

func mapJournalLines(inputs []domain.JournalEntryLineInput) []domain.JournalEntryLine {
	lines := make([]domain.JournalEntryLine, 0, len(inputs))
	for _, line := range inputs {
		lines = append(lines, domain.JournalEntryLine{
			AccountID:     line.AccountID,
			DebitAmount:   line.DebitAmount,
			CreditAmount:  line.CreditAmount,
			Description:   strings.TrimSpace(line.Description),
			DescriptionAr: strings.TrimSpace(line.DescriptionAr),
			LineOrder:     line.LineOrder,
			TaxCode:       strings.TrimSpace(line.TaxCode),
			TaxAmount:     line.TaxAmount,
		})
	}
	return lines
}
