package user

import (
	"encoding/json"
	"net/http"

	"gateway-service/internal/auth"
	"gateway-service/internal/middleware"

	userv1 "github.com/nalrematvonesah/user.proto/gen/user/v1"
)

type Handlers struct {
	users userv1.UserServiceClient
	auth  *auth.Manager
}

func New(users userv1.UserServiceClient, auth *auth.Manager) *Handlers {
	return &Handlers{users: users, auth: auth}
}

type credentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type vehicleRequest struct {
	PlateNumber string `json:"plate_number"`
}

type vehiclesResponse struct {
	Vehicles []string `json:"vehicles"`
}

// POST /auth/register
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var req credentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	resp, err := h.users.Register(r.Context(), &userv1.RegisterRequest{
		Email: req.Email, Password: req.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.auth.Issue(resp.GetUserId())
	if err != nil {
		http.Error(w, "issue token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, tokenResponse{UserID: resp.GetUserId(), Token: token})
}

// POST /auth/login
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req credentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	resp, err := h.users.Login(r.Context(), &userv1.LoginRequest{
		Email: req.Email, Password: req.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	token, err := h.auth.Issue(resp.GetUserId())
	if err != nil {
		http.Error(w, "issue token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, tokenResponse{UserID: resp.GetUserId(), Token: token})
}

// POST /auth/logout
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserIDFrom(r.Context())
	if _, err := h.users.Logout(r.Context(), &userv1.UserRequest{UserId: uid}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// POST /vehicles
func (h *Handlers) AddVehicle(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserIDFrom(r.Context())
	var req vehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if _, err := h.users.AddVehicle(r.Context(), &userv1.AddVehicleRequest{
		UserId: uid, PlateNumber: req.PlateNumber,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// DELETE /vehicles
func (h *Handlers) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserIDFrom(r.Context())
	var req vehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if _, err := h.users.DeleteVehicle(r.Context(), &userv1.DeleteVehicleRequest{
		UserId: uid, PlateNumber: req.PlateNumber,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /vehicles
func (h *Handlers) ListVehicles(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserIDFrom(r.Context())
	resp, err := h.users.GetVehicles(r.Context(), &userv1.UserRequest{UserId: uid})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, vehiclesResponse{Vehicles: resp.GetVehicles()})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
