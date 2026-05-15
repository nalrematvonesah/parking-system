// Package parking contains HTTP handlers for the Parking service.
//
// Owner: Tamerlan
package parking

import (
	"encoding/json"
	"net/http"

	parkingv1 "github.com/nalrematvonesah/parking.proto/gen/parking/v1"
)

type Handlers struct {
	parking parkingv1.ParkingServiceClient
}

func New(parking parkingv1.ParkingServiceClient) *Handlers {
	return &Handlers{parking: parking}
}

type availableResponse struct {
	Count int32 `json:"count"`
}

// GET /slots/available
func (h *Handlers) Available(w http.ResponseWriter, r *http.Request) {
	resp, err := h.parking.GetAvailableSlots(r.Context(), &parkingv1.Empty{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(availableResponse{Count: resp.GetCount()})
}
