// Package session contains HTTP handlers for the Session service.
//
// Owner: Aldiyar
package session

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gateway-service/internal/middleware"

	sessionv1 "github.com/nalrematvonesah/session.proto/gen/session/v1"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	sessions sessionv1.SessionServiceClient
}

func New(sessions sessionv1.SessionServiceClient) *Handlers {
	return &Handlers{sessions: sessions}
}

type startRequest struct {
	VehicleID int64  `json:"vehicle_id"`
	RequestID string `json:"request_id,omitempty"`
}

type sessionResponse struct {
	SessionID     int64   `json:"session_id"`
	SlotID        int64   `json:"slot_id"`
	VehicleID     int64   `json:"vehicle_id"`
	StartTimeUnix int64   `json:"start_time_unix"`
	EndTimeUnix   int64   `json:"end_time_unix,omitempty"`
	Amount        float64 `json:"amount,omitempty"`
}

type priceResponse struct {
	Amount         float64 `json:"amount"`
	ElapsedSeconds int64   `json:"elapsed_seconds"`
}

type activeSessionsResponse struct {
	Sessions []sessionResponse `json:"sessions"`
}

// POST /sessions/start
func (h *Handlers) Start(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserIDFrom(r.Context())
	var req startRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	resp, err := h.sessions.StartSession(r.Context(), &sessionv1.StartRequest{
		UserId: uid, VehicleId: req.VehicleID, RequestId: req.RequestID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, toSessionResponse(resp))
}

// POST /sessions/{id}/end
func (h *Handlers) End(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	resp, err := h.sessions.EndSession(r.Context(), &sessionv1.EndRequest{SessionId: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, toSessionResponse(resp))
}

// GET /sessions/{id}
func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	resp, err := h.sessions.GetSession(r.Context(), &sessionv1.SessionRequest{SessionId: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, toSessionResponse(resp))
}

// GET /sessions/{id}/price
func (h *Handlers) Price(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	resp, err := h.sessions.CalculatePrice(r.Context(), &sessionv1.SessionRequest{SessionId: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, priceResponse{
		Amount:         resp.GetAmount(),
		ElapsedSeconds: resp.GetElapsedSeconds(),
	})
}

// GET /sessions/active

func (h *Handlers) Active(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserIDFrom(r.Context())
	resp, err := h.sessions.GetActiveSessions(r.Context(), &sessionv1.UserRequest{UserId: uid})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessions := make([]sessionResponse, 0, len(resp.GetSessions()))
	for _, s := range resp.GetSessions() {
		sessions = append(sessions, toSessionResponse(s))
	}
	writeJSON(w, http.StatusOK, activeSessionsResponse{Sessions: sessions})
}

func (h *Handlers) History(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserIDFrom(r.Context())

	resp, err := h.sessions.GetUserSessions(
		r.Context(),
		&sessionv1.UserRequest{UserId: uid},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessions := make([]sessionResponse, 0, len(resp.GetSessions()))
	for _, s := range resp.GetSessions() {
		sessions = append(sessions, toSessionResponse(s))
	}

	writeJSON(w, http.StatusOK, activeSessionsResponse{
		Sessions: sessions,
	})
}

// ---------- helpers ----------

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	raw := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid session id", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func toSessionResponse(s *sessionv1.SessionResponse) sessionResponse {
	return sessionResponse{
		SessionID:     s.GetSessionId(),
		SlotID:        s.GetSlotId(),
		VehicleID:     s.GetVehicleId(),
		StartTimeUnix: s.GetStartTimeUnix(),
		EndTimeUnix:   s.GetEndTimeUnix(),
		Amount:        s.GetAmount(),
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
