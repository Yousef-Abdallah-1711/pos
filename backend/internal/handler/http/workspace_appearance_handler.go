package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"pos-saas/internal/middleware"
	"pos-saas/internal/workspaceappearance"
)

type workspaceAppearanceStore interface {
	Get(context.Context, int64, int64) (workspaceappearance.Published, error)
	Save(context.Context, int64, int64, int64, json.RawMessage) (workspaceappearance.Published, error)
}

type WorkspaceAppearanceHandler struct{ store workspaceAppearanceStore }

func NewWorkspaceAppearanceHandler(store workspaceAppearanceStore) *WorkspaceAppearanceHandler {
	return &WorkspaceAppearanceHandler{store: store}
}

func workspaceIdentity(w http.ResponseWriter, r *http.Request) (int64, int64, bool) {
	claims, tenantID := middleware.GetUserClaims(r), middleware.GetTenantID(r)
	if claims == nil || claims.UserID <= 0 {
		respondError(w, http.StatusUnauthorized, "Authentication required")
		return 0, 0, false
	}
	if tenantID <= 0 || int64(claims.TenantID) != tenantID {
		respondError(w, http.StatusForbidden, "Tenant context required")
		return 0, 0, false
	}
	w.Header().Set("Cache-Control", "private, no-store")
	return tenantID, int64(claims.UserID), true
}

func (h *WorkspaceAppearanceHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenant, actor, ok := workspaceIdentity(w, r)
	if !ok {
		return
	}
	value, err := h.store.Get(r.Context(), tenant, actor)
	if err != nil {
		log.Printf("workspace appearance read: %v", err)
		respondError(w, http.StatusInternalServerError, "Unable to load workspace appearance")
		return
	}
	respondJSON(w, http.StatusOK, value)
}

func (h *WorkspaceAppearanceHandler) Put(w http.ResponseWriter, r *http.Request) {
	tenant, actor, ok := workspaceIdentity(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 32<<10)
	var request struct {
		ExpectedRevision *int64          `json:"expectedRevision"`
		Overrides        json.RawMessage `json:"overrides"`
	}
	if err := decodeJSONStrict(r, &request); err != nil || request.ExpectedRevision == nil || *request.ExpectedRevision < 0 {
		respondError(w, http.StatusBadRequest, "Invalid appearance request")
		return
	}
	if err := workspaceappearance.ValidateOverrides(request.Overrides); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid appearance overrides")
		return
	}
	value, err := h.store.Save(r.Context(), tenant, actor, *request.ExpectedRevision, request.Overrides)
	if errors.Is(err, workspaceappearance.ErrConflict) {
		respondError(w, http.StatusConflict, "Appearance changed; reload before saving")
		return
	}
	if err != nil {
		log.Printf("workspace appearance save: %v", err)
		respondError(w, http.StatusInternalServerError, "Unable to save workspace appearance")
		return
	}
	respondJSON(w, http.StatusOK, value)
}
