package parking

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	parkingv1 "github.com/nalrematvonesah/parking.proto/gen/parking/v1"
)

type Handlers struct {
	parking parkingv1.ParkingServiceClient
}

type slotResponse struct {
	SlotID     int64 `json:"slot_id"`
	IsOccupied bool  `json:"is_occupied"`
}

type listSlotsResponse struct {
	Slots []slotResponse `json:"slots"`
}

type addSlotsRequest struct {
	Count int32 `json:"count"`
}

func New(parking parkingv1.ParkingServiceClient) *Handlers {
	return &Handlers{parking: parking}
}

type availableResponse struct {
	Count int32 `json:"count"`
}

func (h *Handlers) Available(w http.ResponseWriter, r *http.Request) {
	resp, err := h.parking.GetAvailableSlots(r.Context(), &parkingv1.Empty{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(availableResponse{Count: resp.GetCount()})
}

// GET /admin/slots
func (h *Handlers) ListAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.parking.ListAllSlots(r.Context(), &parkingv1.Empty{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slots := make([]slotResponse, 0, len(resp.GetSlots()))
	for _, s := range resp.GetSlots() {
		slots = append(slots, slotResponse{
			SlotID:     s.GetSlotId(),
			IsOccupied: s.GetIsOccupied(),
		})
	}

	writeJSON(w, http.StatusOK, listSlotsResponse{Slots: slots})
}

// GET /admin/slots/{id}
func (h *Handlers) GetSlot(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid slot id", http.StatusBadRequest)
		return
	}

	resp, err := h.parking.GetSlot(r.Context(), &parkingv1.GetSlotRequest{
		SlotId: id,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, slotResponse{
		SlotID:     resp.GetSlotId(),
		IsOccupied: resp.GetIsOccupied(),
	})
}

// POST /admin/slots
func (h *Handlers) AddSlots(w http.ResponseWriter, r *http.Request) {
	var req addSlotsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if req.Count <= 0 {
		http.Error(w, "count must be positive", http.StatusBadRequest)
		return
	}

	_, err := h.parking.AddSlots(r.Context(), &parkingv1.AddSlotsRequest{
		Count: req.Count,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
